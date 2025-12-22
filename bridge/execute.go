package bridge

import (
	"context"
	"encoding/json"
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

	// Parse parameters
	paramValue, err := parseParams(params, fn.InputType)
	if err != nil {
		if bridgeErr, ok := err.(*Error); ok {
			return ExecuteResult{Error: bridgeErr}
		}
		return ExecuteResult{
			Error: NewError(ErrCodeInvalidParams, "Invalid parameters", err.Error()),
		}
	}

	// Validate parameters
	if err := validateParams(paramValue, fn.InputType); err != nil {
		if bridgeErr, ok := err.(*Error); ok {
			return ExecuteResult{Error: bridgeErr}
		}
		return ExecuteResult{
			Error: NewError(ErrCodeInvalidParams, "Parameter validation failed", err.Error()),
		}
	}

	// Execute with timeout
	result := b.executeWithTimeout(ctx, fn, paramValue)

	// Calculate duration
	duration := time.Since(startTime).Microseconds()

	// Trigger after hook
	hookData := HookData{
		FunctionName: funcName,
		Params:       paramValue.Interface(),
		Result:       result.Result,
		Duration:     duration,
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

		// Call the function
		results := fn.Handler.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			paramValue,
		})

		// Extract result and error
		var result any
		var err error

		// Check if result can be nil before calling IsNil()
		// Only pointer, interface, slice, map, chan, and func types can be nil
		if results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || 
		   results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Map || 
		   results[0].Kind() == reflect.Chan || results[0].Kind() == reflect.Func {
			if !results[0].IsNil() {
				result = results[0].Interface()
			}
		} else {
			// For non-nillable types (int, string, bool, struct, etc.), just get the value
			result = results[0].Interface()
		}

		// Error is always an interface, so safe to check IsNil
		if !results[1].IsNil() {
			err = results[1].Interface().(error)
		}

		if err != nil {
			// Check if it's already a bridge error
			if bridgeErr, ok := err.(*Error); ok {
				resultChan <- ExecuteResult{Error: bridgeErr}
			} else {
				resultChan <- ExecuteResult{
					Error: NewError(ErrCodeInternal, err.Error()),
				}
			}
		} else {
			resultChan <- ExecuteResult{Result: result}
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

