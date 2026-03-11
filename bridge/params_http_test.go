package bridge

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type httpTestParams struct {
	Name   string  `json:"name"`
	Count  int     `json:"count,omitempty"`
	Active bool    `json:"active,omitempty"`
	Score  float64 `json:"score,omitempty"`
}

func TestParseHTTPParams_QueryParams(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test?name=hello&count=5&active=true&score=3.14", nil)
	targetType := reflect.TypeFor[httpTestParams]()

	val, err := parseHTTPParams(r, targetType)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "hello" {
		t.Errorf("Name = %q, want %q", result.Name, "hello")
	}
	if result.Count != 5 {
		t.Errorf("Count = %d, want 5", result.Count)
	}
	if !result.Active {
		t.Error("Active = false, want true")
	}
	if result.Score != 3.14 {
		t.Errorf("Score = %f, want 3.14", result.Score)
	}
}

func TestParseHTTPParams_FormData(t *testing.T) {
	body := strings.NewReader("name=world&count=10&active=false")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	targetType := reflect.TypeFor[httpTestParams]()

	val, err := parseHTTPParams(r, targetType)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "world" {
		t.Errorf("Name = %q, want %q", result.Name, "world")
	}
	if result.Count != 10 {
		t.Errorf("Count = %d, want 10", result.Count)
	}
	if result.Active {
		t.Error("Active = true, want false")
	}
}

func TestParseHTTPParams_JSONBody(t *testing.T) {
	body := strings.NewReader(`{"name":"json-test","count":42,"active":true,"score":2.5}`)
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	r.Header.Set("Content-Type", "application/json")
	targetType := reflect.TypeFor[httpTestParams]()

	val, err := parseHTTPParams(r, targetType)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "json-test" {
		t.Errorf("Name = %q, want %q", result.Name, "json-test")
	}
	if result.Count != 42 {
		t.Errorf("Count = %d, want 42", result.Count)
	}
	if !result.Active {
		t.Error("Active = false, want true")
	}
	if result.Score != 2.5 {
		t.Errorf("Score = %f, want 2.5", result.Score)
	}
}

func TestParseHTTPParams_MissingFields(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test?name=partial", nil)
	targetType := reflect.TypeFor[httpTestParams]()

	val, err := parseHTTPParams(r, targetType)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "partial" {
		t.Errorf("Name = %q, want %q", result.Name, "partial")
	}
	if result.Count != 0 {
		t.Errorf("Count = %d, want 0 (zero value)", result.Count)
	}
}

func TestParseHTTPParams_InvalidInt(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test?count=not-a-number", nil)
	targetType := reflect.TypeFor[httpTestParams]()

	_, err := parseHTTPParams(r, targetType)
	if err == nil {
		t.Fatal("parseHTTPParams() expected error for invalid int, got nil")
	}
}

func TestParseHTTPParams_NilTargetType(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	val, err := parseHTTPParams(r, nil)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	if val.IsValid() {
		t.Error("expected invalid reflect.Value for nil targetType")
	}
}

func TestMapToStruct_JSONTagNames(t *testing.T) {
	type tagTest struct {
		FieldOne string `json:"field_one"`
		FieldTwo int    `json:"field_two,omitempty"`
	}

	values := map[string][]string{
		"field_one": {"value1"},
		"field_two": {"99"},
	}

	val, err := mapToStruct(values, reflect.TypeFor[tagTest]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	result := val.Interface().(tagTest)
	if result.FieldOne != "value1" {
		t.Errorf("FieldOne = %q, want %q", result.FieldOne, "value1")
	}
	if result.FieldTwo != 99 {
		t.Errorf("FieldTwo = %d, want 99", result.FieldTwo)
	}
}

// --- setFieldFromString tests ---

func TestSetFieldFromString_AllTypes(t *testing.T) {
	type allTypes struct {
		S   string  `json:"s"`
		I   int     `json:"i"`
		I8  int8    `json:"i8"`
		I16 int16   `json:"i16"`
		I32 int32   `json:"i32"`
		I64 int64   `json:"i64"`
		U   uint    `json:"u"`
		U8  uint8   `json:"u8"`
		U16 uint16  `json:"u16"`
		U32 uint32  `json:"u32"`
		U64 uint64  `json:"u64"`
		F32 float32 `json:"f32"`
		F64 float64 `json:"f64"`
		B   bool    `json:"b"`
	}

	values := map[string][]string{
		"s":   {"hello"},
		"i":   {"-42"},
		"i8":  {"127"},
		"i16": {"32000"},
		"i32": {"2000000"},
		"i64": {"9000000000"},
		"u":   {"100"},
		"u8":  {"255"},
		"u16": {"65000"},
		"u32": {"4000000000"},
		"u64": {"18000000000000000000"},
		"f32": {"3.14"},
		"f64": {"2.718281828"},
		"b":   {"true"},
	}

	val, err := mapToStruct(values, reflect.TypeFor[allTypes]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	result := val.Interface().(allTypes)
	if result.S != "hello" {
		t.Errorf("S = %q, want %q", result.S, "hello")
	}
	if result.I != -42 {
		t.Errorf("I = %d, want -42", result.I)
	}
	if result.I8 != 127 {
		t.Errorf("I8 = %d, want 127", result.I8)
	}
	if result.U8 != 255 {
		t.Errorf("U8 = %d, want 255", result.U8)
	}
	if result.F64 != 2.718281828 {
		t.Errorf("F64 = %f, want 2.718281828", result.F64)
	}
	if !result.B {
		t.Error("B = false, want true")
	}
}

func TestSetFieldFromString_InvalidFloat(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test?score=not-float", nil)
	targetType := reflect.TypeFor[httpTestParams]()

	_, err := parseHTTPParams(r, targetType)
	if err == nil {
		t.Fatal("expected error for invalid float, got nil")
	}
}

func TestSetFieldFromString_InvalidBool(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test?active=maybe", nil)
	targetType := reflect.TypeFor[httpTestParams]()

	_, err := parseHTTPParams(r, targetType)
	if err == nil {
		t.Fatal("expected error for invalid bool, got nil")
	}
}

func TestSetFieldFromString_InvalidUint(t *testing.T) {
	type uintStruct struct {
		Val uint `json:"val"`
	}

	r := httptest.NewRequest(http.MethodGet, "/test?val=-5", nil)
	targetType := reflect.TypeFor[uintStruct]()

	_, err := parseHTTPParams(r, targetType)
	if err == nil {
		t.Fatal("expected error for negative uint, got nil")
	}
}

func TestSetFieldFromString_ComplexTypeJSON(t *testing.T) {
	type nested struct {
		Inner struct {
			X int `json:"x"`
		} `json:"inner"`
	}

	values := map[string][]string{
		"inner": {`{"x":42}`},
	}

	val, err := mapToStruct(values, reflect.TypeFor[nested]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	result := val.Interface().(nested)
	if result.Inner.X != 42 {
		t.Errorf("Inner.X = %d, want 42", result.Inner.X)
	}
}

func TestSetFieldFromString_InvalidComplexType(t *testing.T) {
	type nested struct {
		Inner struct {
			X int `json:"x"`
		} `json:"inner"`
	}

	values := map[string][]string{
		"inner": {`not-json`},
	}

	_, err := mapToStruct(values, reflect.TypeFor[nested]())
	if err == nil {
		t.Fatal("expected error for invalid JSON in complex type, got nil")
	}
}

// --- mapToStruct edge cases ---

func TestMapToStruct_UnexportedFieldsSkipped(t *testing.T) {
	type withPrivate struct {
		Public  string `json:"public"`
		private string //nolint:unused
	}

	values := map[string][]string{
		"public":  {"value"},
		"private": {"should-not-set"},
	}

	val, err := mapToStruct(values, reflect.TypeFor[withPrivate]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	result := val.Interface().(withPrivate)
	if result.Public != "value" {
		t.Errorf("Public = %q, want %q", result.Public, "value")
	}
}

func TestMapToStruct_JSONDashIgnored(t *testing.T) {
	type withIgnored struct {
		Name    string `json:"name"`
		Ignored string `json:"-"`
	}

	values := map[string][]string{
		"name":    {"hello"},
		"Ignored": {"should-not-set"},
	}

	val, err := mapToStruct(values, reflect.TypeFor[withIgnored]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	result := val.Interface().(withIgnored)
	if result.Name != "hello" {
		t.Errorf("Name = %q, want %q", result.Name, "hello")
	}
	if result.Ignored != "" {
		t.Errorf("Ignored = %q, want empty (json:\"-\" field)", result.Ignored)
	}
}

func TestMapToStruct_NoJSONTag_UsesFieldName(t *testing.T) {
	type noTag struct {
		MyField string
	}

	values := map[string][]string{
		"MyField": {"value"},
	}

	val, err := mapToStruct(values, reflect.TypeFor[noTag]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	result := val.Interface().(noTag)
	if result.MyField != "value" {
		t.Errorf("MyField = %q, want %q", result.MyField, "value")
	}
}

func TestMapToStruct_EmptyValues(t *testing.T) {
	values := map[string][]string{}

	val, err := mapToStruct(values, reflect.TypeFor[httpTestParams]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "" {
		t.Errorf("Name = %q, want empty", result.Name)
	}
	if result.Count != 0 {
		t.Errorf("Count = %d, want 0", result.Count)
	}
}

func TestMapToStruct_NonStructType(t *testing.T) {
	values := map[string][]string{"key": {"val"}}

	val, err := mapToStruct(values, reflect.TypeFor[string]())
	if err != nil {
		t.Fatalf("mapToStruct() error = %v", err)
	}

	// Should return zero value for non-struct types
	result := val.Interface().(string)
	if result != "" {
		t.Errorf("result = %q, want empty string", result)
	}
}

// --- parseHTTPParams content type routing ---

func TestParseHTTPParams_MultipartFormData(t *testing.T) {
	body := strings.NewReader("name=multipart-test&count=7")
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	targetType := reflect.TypeFor[httpTestParams]()

	val, err := parseHTTPParams(r, targetType)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "multipart-test" {
		t.Errorf("Name = %q, want %q", result.Name, "multipart-test")
	}
}

func TestParseHTTPParams_InvalidJSON(t *testing.T) {
	body := strings.NewReader(`{invalid json}`)
	r := httptest.NewRequest(http.MethodPost, "/test", body)
	r.Header.Set("Content-Type", "application/json")
	targetType := reflect.TypeFor[httpTestParams]()

	_, err := parseHTTPParams(r, targetType)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParseHTTPParams_NoContentType_FallsToQuery(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/test?name=query-fallback", nil)
	// No Content-Type header
	targetType := reflect.TypeFor[httpTestParams]()

	val, err := parseHTTPParams(r, targetType)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "query-fallback" {
		t.Errorf("Name = %q, want %q", result.Name, "query-fallback")
	}
}

// --- multipart form tests ---

func TestParseHTTPParams_MultipartForm(t *testing.T) {
	body := &strings.Builder{}
	boundary := "testboundary"
	body.WriteString("--" + boundary + "\r\n")
	body.WriteString("Content-Disposition: form-data; name=\"name\"\r\n\r\n")
	body.WriteString("multipart-test\r\n")
	body.WriteString("--" + boundary + "\r\n")
	body.WriteString("Content-Disposition: form-data; name=\"count\"\r\n\r\n")
	body.WriteString("42\r\n")
	body.WriteString("--" + boundary + "--\r\n")

	r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body.String()))
	r.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

	targetType := reflect.TypeFor[httpTestParams]()
	val, err := parseHTTPParams(r, targetType)
	if err != nil {
		t.Fatalf("parseHTTPParams() error = %v", err)
	}

	result := val.Interface().(httpTestParams)
	if result.Name != "multipart-test" {
		t.Errorf("Name = %q, want %q", result.Name, "multipart-test")
	}
	if result.Count != 42 {
		t.Errorf("Count = %d, want 42", result.Count)
	}
}

func TestParseFormData_InvalidMultipart(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("not valid multipart"))
	r.Header.Set("Content-Type", "multipart/form-data; boundary=missing")

	targetType := reflect.TypeFor[httpTestParams]()
	_, err := parseFormData(r, targetType)
	if err == nil {
		t.Error("expected error for invalid multipart data")
	}
}
