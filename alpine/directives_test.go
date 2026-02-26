package alpine

import (
	"testing"
)

func TestXData(t *testing.T) {
	tests := []struct {
		name  string
		state map[string]any
		key   string
	}{
		{
			name:  "nil state",
			state: nil,
			key:   "x-data",
		},
		{
			name:  "empty state",
			state: map[string]any{},
			key:   "x-data",
		},
		{
			name:  "simple state",
			state: map[string]any{"open": false},
			key:   "x-data",
		},
		{
			name: "complex state",
			state: map[string]any{
				"count": 0,
				"items": []any{},
			},
			key: "x-data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := XData(tt.state)

			if _, ok := attrs[tt.key]; !ok {
				t.Errorf("XData() missing key %q", tt.key)
			}
		})
	}
}

func TestXShow(t *testing.T) {
	attrs := XShow("isVisible")

	if v, ok := attrs["x-show"]; !ok || v != "isVisible" {
		t.Errorf("XShow() = %v, want x-show=isVisible", attrs)
	}
}

func TestXIf(t *testing.T) {
	attrs := XIf("count > 5")

	if v, ok := attrs["x-if"]; !ok || v != "count > 5" {
		t.Errorf("XIf() = %v, want x-if=\"count > 5\"", attrs)
	}
}

func TestXFor(t *testing.T) {
	attrs := XFor("item in items")

	if v, ok := attrs["x-for"]; !ok || v != "item in items" {
		t.Errorf("XFor() = %v, want x-for=\"item in items\"", attrs)
	}
}

func TestXForKeyed(t *testing.T) {
	attrs := XForKeyed("item in items", "item.id")

	if v, ok := attrs["x-for"]; !ok || v != "item in items" {
		t.Errorf("XForKeyed() missing x-for attribute, got %v", attrs)
	}

	if v, ok := attrs[":key"]; !ok || v != "item.id" {
		t.Errorf("XForKeyed() missing :key attribute, got %v", attrs)
	}
}

func TestXBind(t *testing.T) {
	tests := []struct {
		name string
		attr string
		expr string
		want string
	}{
		{
			name: "disabled binding",
			attr: "disabled",
			expr: "loading",
			want: "loading",
		},
		{
			name: "href binding",
			attr: "href",
			expr: "'/users/' + userId",
			want: "'/users/' + userId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := XBind(tt.attr, tt.expr)
			key := ":" + tt.attr

			if v, ok := attrs[key]; !ok || v != tt.want {
				t.Errorf("XBind() = %v, want %s=%s", attrs, key, tt.want)
			}
		})
	}
}

func TestXModel(t *testing.T) {
	attrs := XModel("name")

	if v, ok := attrs["x-model"]; !ok || v != "name" {
		t.Errorf("XModel() = %v, want x-model=name", attrs)
	}
}

func TestXModelNumber(t *testing.T) {
	attrs := XModelNumber("age")

	if v, ok := attrs["x-model.number"]; !ok || v != "age" {
		t.Errorf("XModelNumber() = %v, want x-model.number=age", attrs)
	}
}

func TestXModelLazy(t *testing.T) {
	attrs := XModelLazy("description")

	if v, ok := attrs["x-model.lazy"]; !ok || v != "description" {
		t.Errorf("XModelLazy() = %v, want x-model.lazy=description", attrs)
	}
}

func TestXModelDebounce(t *testing.T) {
	attrs := XModelDebounce("searchQuery", 300)

	if v, ok := attrs["x-model.debounce.300ms"]; !ok || v != "searchQuery" {
		t.Errorf("XModelDebounce() = %v, want x-model.debounce.300ms=searchQuery", attrs)
	}
}

func TestXOn(t *testing.T) {
	tests := []struct {
		name    string
		event   string
		handler string
	}{
		{
			name:    "click event",
			event:   "click",
			handler: "count++",
		},
		{
			name:    "submit prevent",
			event:   "submit.prevent",
			handler: "handleSubmit()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := XOn(tt.event, tt.handler)
			key := "@" + tt.event

			if v, ok := attrs[key]; !ok || v != tt.handler {
				t.Errorf("XOn() = %v, want %s=%s", attrs, key, tt.handler)
			}
		})
	}
}

func TestXClick(t *testing.T) {
	attrs := XClick("doSomething()")

	if v, ok := attrs["@click"]; !ok || v != "doSomething()" {
		t.Errorf("XClick() = %v, want @click=doSomething()", attrs)
	}
}

func TestXSubmit(t *testing.T) {
	attrs := XSubmit("submit()")

	if v, ok := attrs["@submit.prevent"]; !ok || v != "submit()" {
		t.Errorf("XSubmit() = %v, want @submit.prevent=submit()", attrs)
	}
}

func TestXInput(t *testing.T) {
	attrs := XInput("update()")

	if v, ok := attrs["@input"]; !ok || v != "update()" {
		t.Errorf("XInput() = %v, want @input=update()", attrs)
	}
}

func TestXChange(t *testing.T) {
	attrs := XChange("onChange()")

	if v, ok := attrs["@change"]; !ok || v != "onChange()" {
		t.Errorf("XChange() = %v, want @change=onChange()", attrs)
	}
}

func TestXText(t *testing.T) {
	attrs := XText("user.name")

	if v, ok := attrs["x-text"]; !ok || v != "user.name" {
		t.Errorf("XText() = %v, want x-text=user.name", attrs)
	}
}

func TestXHtml(t *testing.T) {
	attrs := XHtml("content")

	if v, ok := attrs["x-html"]; !ok || v != "content" {
		t.Errorf("XHtml() = %v, want x-html=content", attrs)
	}
}

func TestXRef(t *testing.T) {
	attrs := XRef("emailInput")

	if v, ok := attrs["x-ref"]; !ok || v != "emailInput" {
		t.Errorf("XRef() = %v, want x-ref=emailInput", attrs)
	}
}

func TestXInit(t *testing.T) {
	attrs := XInit("loadData()")

	if v, ok := attrs["x-init"]; !ok || v != "loadData()" {
		t.Errorf("XInit() = %v, want x-init=loadData()", attrs)
	}
}

func TestXInitFetch(t *testing.T) {
	attrs := XInitFetch("/api/users", "users")

	v, ok := attrs["x-init"]
	if !ok {
		t.Fatal("XInitFetch() missing x-init attribute")
	}

	val, isStr := v.(string)
	if !isStr {
		t.Fatal("XInitFetch() x-init value is not a string")
	}

	if len(val) == 0 {
		t.Error("XInitFetch() x-init value is empty")
	}
}

func TestXCloak(t *testing.T) {
	attrs := XCloak()

	if _, ok := attrs["x-cloak"]; !ok {
		t.Errorf("XCloak() = %v, want x-cloak attribute", attrs)
	}
}

func TestXIgnore(t *testing.T) {
	attrs := XIgnore()

	if _, ok := attrs["x-ignore"]; !ok {
		t.Errorf("XIgnore() = %v, want x-ignore attribute", attrs)
	}
}
