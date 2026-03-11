package bridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"time"
)

// ExecuteResult holds the result of a function execution
type ExecuteResult struct {
	Result any
	Error  *Error
}

// execute runs a registered function with the given parameters
func (b *Bridge) execute(ctx Context, funcName string, params json.RawMessage) ExecuteResult {
	startTime := time.Now()

	// Get the function
	fn, err := b.GetFunction(funcName)
	if err != nil {
		return ExecuteResult{
			Error: NewError(ErrCodeMethodNotFound, fmt.Sprintf("Function '%s' not found", funcName)),
		}
	}

	// Trigger before hook
	b.hooks.Trigger(BeforeCall, ctx, HookData{
		FunctionName: funcName,
		Params:       params,
	})

	var paramValue reflect.Value

	if fn.HasInput {
		// Parse parameters
		var parseErr error
		paramValue, parseErr = parseParams(params, fn.InputType)
		if parseErr != nil {
			var bridgeErr *Error
			if errors.As(parseErr, &bridgeErr) {
				return ExecuteResult{Error: bridgeErr}
			}

			return ExecuteResult{
				Error: NewError(ErrCodeInvalidParams, "Invalid parameters", parseErr.Error()),
			}
		}

		// Validate parameters
		if fn.LaxValidation {
			if validateErr := validateParamsLax(paramValue, fn.InputType); validateErr != nil {
				var bridgeErr *Error
				if errors.As(validateErr, &bridgeErr) {
					return ExecuteResult{Error: bridgeErr}
				}

				return ExecuteResult{
					Error: NewError(ErrCodeInvalidParams, "Parameter validation failed", validateErr.Error()),
				}
			}
		} else {
			if validateErr := validateParams(paramValue, fn.InputType); validateErr != nil {
				var bridgeErr *Error
				if errors.As(validateErr, &bridgeErr) {
					return ExecuteResult{Error: bridgeErr}
				}

				return ExecuteResult{
					Error: NewError(ErrCodeInvalidParams, "Parameter validation failed", validateErr.Error()),
				}
			}
		}
	}

	// Execute with timeout
	result := b.executeWithTimeout(ctx, fn, paramValue)

	// Calculate duration
	duration := time.Since(startTime).Microseconds()

	// Trigger after hook
	hookData := HookData{
		FunctionName: funcName,
		Duration:     duration,
		Result:       result.Result,
	}

	if fn.HasInput {
		hookData.Params = paramValue.Interface()
	}

	if result.Error != nil {
		hookData.Error = result.Error
		b.hooks.Trigger(OnError, ctx, hookData)
	} else {
		b.hooks.Trigger(OnSuccess, ctx, hookData)
	}

	b.hooks.Trigger(AfterCall, ctx, hookData)

	return result
}

// executeDirect runs a registered function with a pre-parsed parameter value.
// This avoids the JSON round-trip used by execute() and is used by the HTMX handler.
func (b *Bridge) executeDirect(ctx Context, fn *Function, paramValue reflect.Value) ExecuteResult {
	startTime := time.Now()

	// Trigger before hook
	b.hooks.Trigger(BeforeCall, ctx, HookData{
		FunctionName: fn.Name,
	})

	// Validate parameters if needed
	if fn.HasInput && !fn.LaxValidation {
		if validateErr := validateParams(paramValue, fn.InputType); validateErr != nil {
			var bridgeErr *Error
			if errors.As(validateErr, &bridgeErr) {
				return ExecuteResult{Error: bridgeErr}
			}

			return ExecuteResult{
				Error: NewError(ErrCodeInvalidParams, "Parameter validation failed", validateErr.Error()),
			}
		}
	} else if fn.HasInput && fn.LaxValidation {
		if validateErr := validateParamsLax(paramValue, fn.InputType); validateErr != nil {
			var bridgeErr *Error
			if errors.As(validateErr, &bridgeErr) {
				return ExecuteResult{Error: bridgeErr}
			}

			return ExecuteResult{
				Error: NewError(ErrCodeInvalidParams, "Parameter validation failed", validateErr.Error()),
			}
		}
	}

	// Execute with timeout
	result := b.executeWithTimeout(ctx, fn, paramValue)

	// Calculate duration
	duration := time.Since(startTime).Microseconds()

	// Trigger after hook
	hookData := HookData{
		FunctionName: fn.Name,
		Duration:     duration,
		Result:       result.Result,
	}

	if fn.HasInput {
		hookData.Params = paramValue.Interface()
	}

	if result.Error != nil {
		hookData.Error = result.Error
		b.hooks.Trigger(OnError, ctx, hookData)
	} else {
		b.hooks.Trigger(OnSuccess, ctx, hookData)
	}

	b.hooks.Trigger(AfterCall, ctx, hookData)

	return result
}

// executeWithTimeout executes a function with timeout and panic recovery
func (b *Bridge) executeWithTimeout(ctx Context, fn *Function, paramValue reflect.Value) ExecuteResult {
	// Create context with timeout
	// Use function-specific timeout if set, otherwise use bridge config timeout
	timeout := fn.Timeout
	if timeout == 0 {
		timeout = b.config.Timeout
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
	defer cancel()

	// Create a channel to receive the result
	resultChan := make(chan ExecuteResult, 1)

	// Execute in a goroutine
	go func() {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				resultChan <- ExecuteResult{
					Error: NewError(ErrCodeInternal, "Function panicked", map[string]any{
						"panic": fmt.Sprint(r),
						"stack": string(stack),
					}),
				}
			}
		}()

		// Build args based on signature type
		var args []reflect.Value
		switch fn.SignatureType {
		case SigInputOutput, SigInputOnly:
			args = []reflect.Value{reflect.ValueOf(ctx), paramValue}
		case SigOutput, SigVoid:
			args = []reflect.Value{reflect.ValueOf(ctx)}
		}

		// Call the function
		results := fn.Handler.Call(args)

		// Extract result and error based on signature type
		switch fn.SignatureType {
		case SigInputOutput, SigOutput:
			// 2 returns: (result, error)
			var (
				result any
				fnErr  error
			)

			if canBeNil(results[0]) {
				if !results[0].IsNil() {
					result = results[0].Interface()
				}
			} else {
				result = results[0].Interface()
			}

			if !results[1].IsNil() {
				fnErr = results[1].Interface().(error)
			}

			if fnErr != nil {
				var bridgeErr *Error
				if errors.As(fnErr, &bridgeErr) {
					resultChan <- ExecuteResult{Error: bridgeErr}
				} else {
					resultChan <- ExecuteResult{
						Error: NewError(ErrCodeInternal, fnErr.Error()),
					}
				}
			} else {
				resultChan <- ExecuteResult{Result: result}
			}

		case SigInputOnly, SigVoid:
			// 1 return: (error)
			if !results[0].IsNil() {
				fnErr := results[0].Interface().(error)
				var bridgeErr *Error
				if errors.As(fnErr, &bridgeErr) {
					resultChan <- ExecuteResult{Error: bridgeErr}
				} else {
					resultChan <- ExecuteResult{
						Error: NewError(ErrCodeInternal, fnErr.Error()),
					}
				}
			} else {
				resultChan <- ExecuteResult{Result: nil}
			}
		}
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result
	case <-timeoutCtx.Done():
		return ExecuteResult{
			Error: NewError(ErrCodeTimeout, fmt.Sprintf("Function execution timed out after %v", timeout)),
		}
	}
}

// canBeNil checks if a reflect.Value's kind supports IsNil()
func canBeNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return true
	}
	return false
}

// Call executes a function and returns the result
// This is a high-level method that handles all execution logic
func (b *Bridge) Call(ctx Context, funcName string, params json.RawMessage) (any, error) {
	result := b.execute(ctx, funcName, params)
	if result.Error != nil {
		return nil, result.Error
	}

	return result.Result, nil
}

// CallBatch executes multiple functions in parallel
func (b *Bridge) CallBatch(ctx Context, requests []Request) []Response {
	// Limit batch size
	if len(requests) > b.config.MaxBatchSize {
		return []Response{{
			JSONRPC: "2.0",
			Error:   NewError(ErrCodeBadRequest, fmt.Sprintf("Batch size exceeds maximum of %d", b.config.MaxBatchSize)),
		}}
	}

	responses := make([]Response, len(requests))

	// Execute all requests in parallel
	done := make(chan bool)

	for i, req := range requests {
		go func(index int, request Request) {
			result := b.execute(ctx, request.Method, request.Params)

			responses[index] = Response{
				JSONRPC: "2.0",
				ID:      request.ID,
			}

			if result.Error != nil {
				responses[index].Error = result.Error
			} else {
				responses[index].Result = result.Result
			}

			done <- true
		}(i, req)
	}

	// Wait for all executions to complete
	for range requests {
		<-done
	}

	return responses
}
