package bridge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// parseHTTPParams extracts parameters from an HTTP request into the target type.
// It detects the content type and parses accordingly:
//   - application/json → JSON body
//   - application/x-www-form-urlencoded or multipart/form-data → form data
//   - default (GET requests) → query parameters
func parseHTTPParams(r *http.Request, targetType reflect.Type) (reflect.Value, error) {
	if targetType == nil {
		return reflect.Value{}, nil
	}

	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		return parseJSONBody(r, targetType)
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"),
		strings.HasPrefix(contentType, "multipart/form-data"):
		return parseFormData(r, targetType)
	default:
		return parseQueryParams(r, targetType)
	}
}

// parseQueryParams parses URL query parameters into the target struct type
func parseQueryParams(r *http.Request, targetType reflect.Type) (reflect.Value, error) {
	return mapToStruct(r.URL.Query(), targetType)
}

// parseFormData parses form data (URL-encoded or multipart) into the target struct type
func parseFormData(r *http.Request, targetType reflect.Type) (reflect.Value, error) {
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// ParseMultipartForm reads multipart body; 32 MB max memory
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return reflect.Value{}, &Error{
				Code:    ErrCodeInvalidParams,
				Message: "Failed to parse multipart form data",
				Data:    err.Error(),
			}
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return reflect.Value{}, &Error{
				Code:    ErrCodeInvalidParams,
				Message: "Failed to parse form data",
				Data:    err.Error(),
			}
		}
	}

	return mapToStruct(r.Form, targetType)
}

// parseJSONBody parses a JSON request body into the target type
func parseJSONBody(r *http.Request, targetType reflect.Type) (reflect.Value, error) {
	valuePtr := reflect.New(targetType)
	if err := json.NewDecoder(r.Body).Decode(valuePtr.Interface()); err != nil {
		return reflect.Value{}, &Error{
			Code:    ErrCodeInvalidParams,
			Message: "Failed to parse JSON body",
			Data:    err.Error(),
		}
	}

	return valuePtr.Elem(), nil
}

// mapToStruct converts url.Values to a struct via reflection, using JSON tag names as keys
func mapToStruct(values url.Values, targetType reflect.Type) (reflect.Value, error) {
	if targetType.Kind() != reflect.Struct {
		return reflect.New(targetType).Elem(), nil
	}

	valuePtr := reflect.New(targetType)
	elem := valuePtr.Elem()

	for i := range targetType.NumField() {
		field := targetType.Field(i)
		if !field.IsExported() {
			continue
		}

		// Determine the key: use json tag name if present, otherwise field name
		key := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		if jsonTag != "" {
			parts := strings.SplitN(jsonTag, ",", 2)
			if parts[0] != "" {
				key = parts[0]
			}
		}

		vals, ok := values[key]
		if !ok || len(vals) == 0 {
			continue
		}

		fieldValue := elem.Field(i)
		if err := setFieldFromString(fieldValue, vals[0]); err != nil {
			return reflect.Value{}, &Error{
				Code:    ErrCodeInvalidParams,
				Message: fmt.Sprintf("Invalid value for field '%s': %s", key, err.Error()),
			}
		}
	}

	return elem, nil
}

// setFieldFromString sets a reflect.Value from a string, handling type conversions
func setFieldFromString(field reflect.Value, s string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as integer", s)
		}
		field.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as unsigned integer", s)
		}
		field.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as float", s)
		}
		field.SetFloat(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("cannot parse %q as boolean", s)
		}
		field.SetBool(v)
	default:
		// For complex types, try JSON unmarshal
		ptr := reflect.New(field.Type())
		if err := json.Unmarshal([]byte(s), ptr.Interface()); err != nil {
			return fmt.Errorf("cannot parse %q into %s", s, field.Type())
		}
		field.Set(ptr.Elem())
	}

	return nil
}
