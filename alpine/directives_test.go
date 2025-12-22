package alpine

import (
	"bytes"
	"strings"
	"testing"
)

func TestXData(t *testing.T) {
	tests := []struct {
		name  string
		state map[string]any
		want  string
	}{
		{
			name:  "nil state",
			state: nil,
			want:  `x-data=""`,
		},
		{
			name:  "empty state",
			state: map[string]any{},
			want:  `x-data=""`,
		},
		{
			name:  "simple state",
			state: map[string]any{"open": false},
			want:  `x-data=`,
		},
		{
			name: "complex state",
			state: map[string]any{
				"count": 0,
				"items": []any{},
			},
			want: `x-data=`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			XData(tt.state).Render(&buf)
			got := buf.String()

			if !strings.Contains(got, tt.want) {
				t.Errorf("XData() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestXShow(t *testing.T) {
	var buf bytes.Buffer
	XShow("isVisible").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-show="isVisible"`) {
		t.Errorf("XShow() = %v, want x-show attribute", got)
	}
}

func TestXIf(t *testing.T) {
	var buf bytes.Buffer
	XIf("count > 5").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-if="count &gt; 5"`) {
		t.Errorf("XIf() = %v, want x-if attribute", got)
	}
}

func TestXFor(t *testing.T) {
	var buf bytes.Buffer
	XFor("item in items").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-for="item in items"`) {
		t.Errorf("XFor() = %v, want x-for attribute", got)
	}
}

func TestXForKeyed(t *testing.T) {
	nodes := XForKeyed("item in items", "item.id")

	if len(nodes) != 2 {
		t.Errorf("XForKeyed() returned %d nodes, want 2", len(nodes))
	}

	var buf bytes.Buffer
	for _, node := range nodes {
		node.Render(&buf)
	}
	got := buf.String()

	if !strings.Contains(got, `x-for="item in items"`) {
		t.Errorf("XForKeyed() missing x-for attribute")
	}
	if !strings.Contains(got, `:key="item.id"`) {
		t.Errorf("XForKeyed() missing :key attribute")
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
			want: `:disabled="loading"`,
		},
		{
			name: "href binding",
			attr: "href",
			expr: "'/users/' + userId",
			want: `:href="&#39;/users/&#39; + userId"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			XBind(tt.attr, tt.expr).Render(&buf)
			got := buf.String()

			if !strings.Contains(got, `:`+tt.attr+`=`) {
				t.Errorf("XBind() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestXModel(t *testing.T) {
	var buf bytes.Buffer
	XModel("name").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-model="name"`) {
		t.Errorf("XModel() = %v, want x-model attribute", got)
	}
}

func TestXModelNumber(t *testing.T) {
	var buf bytes.Buffer
	XModelNumber("age").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-model.number="age"`) {
		t.Errorf("XModelNumber() = %v, want x-model.number attribute", got)
	}
}

func TestXModelLazy(t *testing.T) {
	var buf bytes.Buffer
	XModelLazy("description").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-model.lazy="description"`) {
		t.Errorf("XModelLazy() = %v, want x-model.lazy attribute", got)
	}
}

func TestXModelDebounce(t *testing.T) {
	var buf bytes.Buffer
	XModelDebounce("searchQuery", 300).Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-model.debounce.300ms="searchQuery"`) {
		t.Errorf("XModelDebounce() = %v, want debounced model", got)
	}
}

func TestXOn(t *testing.T) {
	tests := []struct {
		name    string
		event   string
		handler string
		want    string
	}{
		{
			name:    "click event",
			event:   "click",
			handler: "count++",
			want:    `@click="count&#43;&#43;"`,
		},
		{
			name:    "submit prevent",
			event:   "submit.prevent",
			handler: "handleSubmit()",
			want:    `@submit.prevent="handleSubmit()"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			XOn(tt.event, tt.handler).Render(&buf)
			got := buf.String()

			if !strings.Contains(got, `@`+tt.event) {
				t.Errorf("XOn() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestXClick(t *testing.T) {
	var buf bytes.Buffer
	XClick("doSomething()").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `@click=`) {
		t.Errorf("XClick() = %v, want @click attribute", got)
	}
}

func TestXSubmit(t *testing.T) {
	var buf bytes.Buffer
	XSubmit("submit()").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `@submit.prevent=`) {
		t.Errorf("XSubmit() = %v, want @submit.prevent attribute", got)
	}
}

func TestXInput(t *testing.T) {
	var buf bytes.Buffer
	XInput("update()").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `@input=`) {
		t.Errorf("XInput() = %v, want @input attribute", got)
	}
}

func TestXChange(t *testing.T) {
	var buf bytes.Buffer
	XChange("onChange()").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `@change=`) {
		t.Errorf("XChange() = %v, want @change attribute", got)
	}
}

func TestXText(t *testing.T) {
	var buf bytes.Buffer
	XText("user.name").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-text="user.name"`) {
		t.Errorf("XText() = %v, want x-text attribute", got)
	}
}

func TestXHtml(t *testing.T) {
	var buf bytes.Buffer
	XHtml("content").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-html="content"`) {
		t.Errorf("XHtml() = %v, want x-html attribute", got)
	}
}

func TestXRef(t *testing.T) {
	var buf bytes.Buffer
	XRef("emailInput").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-ref="emailInput"`) {
		t.Errorf("XRef() = %v, want x-ref attribute", got)
	}
}

func TestXInit(t *testing.T) {
	var buf bytes.Buffer
	XInit("loadData()").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-init="loadData()"`) {
		t.Errorf("XInit() = %v, want x-init attribute", got)
	}
}

func TestXInitFetch(t *testing.T) {
	var buf bytes.Buffer
	XInitFetch("/api/users", "users").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-init=`) {
		t.Errorf("XInitFetch() = %v, want x-init attribute", got)
	}
	if !strings.Contains(got, `/api/users`) {
		t.Errorf("XInitFetch() missing URL")
	}
	if !strings.Contains(got, `users`) {
		t.Errorf("XInitFetch() missing target variable")
	}
}

func TestXCloak(t *testing.T) {
	var buf bytes.Buffer
	XCloak().Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-cloak=""`) {
		t.Errorf("XCloak() = %v, want x-cloak attribute", got)
	}
}

func TestXIgnore(t *testing.T) {
	var buf bytes.Buffer
	XIgnore().Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-ignore=""`) {
		t.Errorf("XIgnore() = %v, want x-ignore attribute", got)
	}
}
