package alpine

import (
	"bytes"
	"strings"
	"testing"
)

func TestRegisterStores(t *testing.T) {
	tests := []struct {
		name   string
		stores []Store
		want   []string
	}{
		{
			name:   "no stores",
			stores: nil,
			want:   make([]string, 0),
		},
		{
			name: "single store with state only",
			stores: []Store{
				{
					Name:  "cart",
					State: map[string]any{"items": make([]any, 0), "total": 0},
				},
			},
			want: []string{"alpine:init", "Alpine.store", "cart"},
		},
		{
			name: "store with methods",
			stores: []Store{
				{
					Name:  "notifications",
					State: map[string]any{"items": make([]any, 0)},
					Methods: `
						add(msg) { this.items.push(msg); },
						clear() { this.items = []; }
					`,
				},
			},
			want: []string{"alpine:init", "Alpine.store", "notifications", "add", "clear"},
		},
		{
			name: "multiple stores",
			stores: []Store{
				{Name: "store1", State: map[string]any{"value": 1}},
				{Name: "store2", State: map[string]any{"value": 2}},
			},
			want: []string{"alpine:init", "store1", "store2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			node := RegisterStores(tt.stores...)

			if len(tt.stores) == 0 {
				// With no stores, function returns g.Group(nil) which can't be rendered directly
				// Just verify it doesn't panic
				return
			}

		if err := node.Render(&buf); err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		got := buf.String()

			// Check all expected strings are present
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("RegisterStores() = %v, want to contain %v", got, w)
				}
			}

			// Check it's wrapped in a script tag
			if !strings.Contains(got, "<script>") {
				t.Errorf("RegisterStores() missing script tag")
			}
		})
	}
}

func TestStoreAccess(t *testing.T) {
	tests := []struct {
		name      string
		storeName string
		key       string
		want      string
	}{
		{
			name:      "simple access",
			storeName: "cart",
			key:       "items",
			want:      "$store.cart.items",
		},
		{
			name:      "nested access",
			storeName: "user",
			key:       "profile.name",
			want:      "$store.user.profile.name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StoreAccess(tt.storeName, tt.key)
			if got != tt.want {
				t.Errorf("StoreAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreMethod(t *testing.T) {
	tests := []struct {
		name      string
		storeName string
		method    string
		args      string
		want      string
	}{
		{
			name:      "method with single arg",
			storeName: "cart",
			method:    "addItem",
			args:      "product",
			want:      "$store.cart.addItem(product)",
		},
		{
			name:      "method with multiple args",
			storeName: "notifications",
			method:    "add",
			args:      "'Error', 'danger'",
			want:      "$store.notifications.add('Error', 'danger')",
		},
		{
			name:      "method with no args",
			storeName: "cart",
			method:    "clear",
			args:      "",
			want:      "$store.cart.clear()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StoreMethod(tt.storeName, tt.method, tt.args)
			if got != tt.want {
				t.Errorf("StoreMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXStore(t *testing.T) {
	var buf bytes.Buffer
	XStore("cart").Render(&buf)
	got := buf.String()

	if !strings.Contains(got, `x-data="$store.cart"`) {
		t.Errorf("XStore() = %v, want x-data with store reference", got)
	}
}
