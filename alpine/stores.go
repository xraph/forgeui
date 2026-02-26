package alpine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
)

// Store represents an Alpine.js global store definition.
type Store struct {
	// Name is the store identifier (accessed via $store.name)
	Name string

	// State is the initial state as a map
	State map[string]any

	// Methods is raw JavaScript code defining store methods
	// Will be merged with the state object
	Methods string
}

// RegisterStores creates a templ.Component that renders a script tag registering Alpine stores.
// This should be placed before the Alpine.js script tag.
func RegisterStores(stores ...Store) templ.Component {
	if len(stores) == 0 {
		return templ.ComponentFunc(func(_ context.Context, _ io.Writer) error {
			return nil
		})
	}

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		var js strings.Builder
		js.WriteString("<script>document.addEventListener('alpine:init', () => {\n")

		for _, store := range stores {
			stateJSON, err := json.Marshal(store.State)
			if err != nil {
				stateJSON = []byte("{}")
			}

			if store.Methods != "" {
				js.WriteString(fmt.Sprintf("  Alpine.store('%s', { ...%s, %s });\n",
					store.Name, stateJSON, strings.TrimSpace(store.Methods)))
			} else {
				js.WriteString(fmt.Sprintf("  Alpine.store('%s', %s);\n",
					store.Name, stateJSON))
			}
		}

		js.WriteString("});</script>")

		_, err := io.WriteString(w, js.String())
		return err
	})
}

// StoreAccess returns the JavaScript expression to access a store property.
func StoreAccess(storeName, key string) string {
	return fmt.Sprintf("$store.%s.%s", storeName, key)
}

// StoreMethod returns the JavaScript expression to call a store method.
func StoreMethod(storeName, method, args string) string {
	return fmt.Sprintf("$store.%s.%s(%s)", storeName, method, args)
}

// XStore creates an x-data directive that initializes with a reference to a store.
func XStore(storeName string) templ.Attributes {
	return templ.Attributes{"x-data": "$store." + storeName}
}
