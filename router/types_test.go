package router

import "testing"

func TestParams_Get(t *testing.T) {
	params := Params{
		"id":   "123",
		"name": "test",
	}

	if params.Get("id") != "123" {
		t.Errorf("Expected id=123, got %s", params.Get("id"))
	}

	if params.Get("missing") != "" {
		t.Error("Expected empty string for missing key")
	}
}

func TestParams_Has(t *testing.T) {
	params := Params{
		"id": "123",
	}

	if !params.Has("id") {
		t.Error("Expected Has(id) to be true")
	}

	if params.Has("missing") {
		t.Error("Expected Has(missing) to be false")
	}
}
