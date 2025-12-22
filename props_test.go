package forgeui

import "testing"

type testProps struct {
	Value string
	Count int
}

func TestApplyOptions(t *testing.T) {
	props := &testProps{
		Value: "initial",
		Count: 0,
	}

	opts := []Option[testProps]{
		func(p *testProps) { p.Value = "updated" },
		func(p *testProps) { p.Count = 5 },
	}

	ApplyOptions(props, opts)

	if props.Value != "updated" {
		t.Errorf("Value = %v, want 'updated'", props.Value)
	}

	if props.Count != 5 {
		t.Errorf("Count = %v, want 5", props.Count)
	}
}

func TestApplyOptions_Empty(t *testing.T) {
	props := &testProps{
		Value: "initial",
		Count: 0,
	}

	ApplyOptions(props, []Option[testProps]{})

	if props.Value != "initial" {
		t.Error("Props should not change when no options applied")
	}
}

func TestBaseProps(t *testing.T) {
	props := BaseProps{
		Class:    "custom-class",
		Disabled: true,
		ID:       "test-id",
	}

	if props.Class != "custom-class" {
		t.Error("Class not set correctly")
	}

	if !props.Disabled {
		t.Error("Disabled not set correctly")
	}

	if props.ID != "test-id" {
		t.Error("ID not set correctly")
	}
}
