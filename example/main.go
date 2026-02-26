package main

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/xraph/forgeui"
	"github.com/xraph/forgeui/bridge"
	"github.com/xraph/forgeui/router"
)

func main() {
	// Initialize ForgeUI App
	app := forgeui.New(
		forgeui.WithDebug(true),
		forgeui.WithAssetPublicDir("example/static"),
	)

	// Initialize Bridge System
	b := bridge.New(
		bridge.WithTimeout(30),
		bridge.WithCSRF(false),
	)

	// Register bridge functions
	_ = b.Register("Greet", Greet)

	// Start dev server with hot reload in development mode
	if app.IsDev() {
		go func() {
			log.Println("Starting dev server with hot reload...")
			if err := app.Assets.StartDevServer(context.Background()); err != nil {
				log.Printf("Dev server error: %v\n", err)
			}
		}()
	}

	// Serve static files through asset pipeline
	http.Handle("/static/", app.Assets.Handler())

	// Bridge HTTP endpoints
	http.Handle("/api/bridge/call", b.Handler())
	http.Handle("/api/bridge/stream/", b.StreamHandler())

	// SSE endpoint for hot reload
	if app.IsDev() {
		if handler := app.Assets.SSEHandler(); handler != nil {
			http.Handle("/_forgeui/reload", handler.(http.Handler))
		}
	}

	// Setup routes using ForgeUI router
	app.Get("/", homePage)
	app.Get("/theme", themePage)
	app.Get("/icons", iconsPage)
	app.Get("/alpine", alpinePage)

	log.Println("ForgeUI Example running at http://localhost:8080")
	log.Println("  Home:     http://localhost:8080")
	log.Println("  Theme:    http://localhost:8080/theme")
	log.Println("  Icons:    http://localhost:8080/icons")
	log.Println("  Alpine:   http://localhost:8080/alpine")

	if err := http.ListenAndServe(":8080", app.Handler()); err != nil {
		log.Fatal(err)
	}
}

// Bridge function
func Greet(name string) string {
	return "Hello, " + name + "! Welcome to ForgeUI."
}

// Page handlers
func homePage(ctx *router.PageContext) (templ.Component, error) {
	return HomePageView(), nil
}

func themePage(ctx *router.PageContext) (templ.Component, error) {
	return ThemePageView(), nil
}

func iconsPage(ctx *router.PageContext) (templ.Component, error) {
	return IconsPageView(), nil
}

func alpinePage(ctx *router.PageContext) (templ.Component, error) {
	return AlpinePageView(), nil
}
