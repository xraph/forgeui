package alpine

import (
	"encoding/json"
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Store represents an Alpine.js global store definition.
//
// Example:
//
//	Store{
//	    Name: "notifications",
//	    State: map[string]any{
//	        "items": []any{},
//	        "count": 0,
//	    },
//	    Methods: `
//	        add(item) {
//	            this.items.push(item);
//	            this.count++;
//	        },
//	        clear() {
//	            this.items = [];
//	            this.count = 0;
//	        }
//	    `,
//	}
type Store struct {
	// Name is the store identifier (accessed via $store.name)
	Name string

	// State is the initial state as a map
	State map[string]any

	// Methods is raw JavaScript code defining store methods
	// Will be merged with the state object
	Methods string
}

// RegisterStores creates a script tag that registers Alpine stores.
// This should be placed before the Alpine.js script tag.
//
// Example:
//
//	alpine.RegisterStores(
//	    Store{
//	        Name: "cart",
//	        State: map[string]any{"items": []any{}, "total": 0},
//	        Methods: `
//	            addItem(item) { this.items.push(item); },
//	            clear() { this.items = []; this.total = 0; }
//	        `,
//	    },
//	)
func RegisterStores(stores ...Store) g.Node {
	if len(stores) == 0 {
		return g.Group(nil)
	}

	var js strings.Builder
	js.WriteString("document.addEventListener('alpine:init', () => {\n")

	for _, store := range stores {
		stateJSON, _ := json.Marshal(store.State)

		// If methods are provided, merge them with state
		if store.Methods != "" {
			js.WriteString(fmt.Sprintf("  Alpine.store('%s', { ...%s, %s });\n",
				store.Name, stateJSON, strings.TrimSpace(store.Methods)))
		} else {
			js.WriteString(fmt.Sprintf("  Alpine.store('%s', %s);\n",
				store.Name, stateJSON))
		}
	}

	js.WriteString("});")

	return html.Script(g.Raw(js.String()))
}

// StoreAccess returns the JavaScript expression to access a store property.
//
// Example:
//
//	alpine.StoreAccess("cart", "items") // Returns: "$store.cart.items"
func StoreAccess(storeName, key string) string {
	return fmt.Sprintf("$store.%s.%s", storeName, key)
}

// StoreMethod returns the JavaScript expression to call a store method.
//
// Example:
//
//	alpine.StoreMethod("cart", "addItem", "product") // Returns: "$store.cart.addItem(product)"
func StoreMethod(storeName, method, args string) string {
	return fmt.Sprintf("$store.%s.%s(%s)", storeName, method, args)
}

// XStore creates an x-data directive that initializes with a reference to a store.
// This allows components to reactively use store data.
//
// Example:
//
//	alpine.XStore("cart") // x-data="$store.cart"
func XStore(storeName string) g.Node {
	return g.Attr("x-data", fmt.Sprintf("$store.%s", storeName))
}

