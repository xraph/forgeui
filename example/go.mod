module example

go 1.23

replace github.com/xraph/forgeui => ../

require (
	github.com/xraph/forgeui v0.0.0-00010101000000-000000000000
	maragu.dev/gomponents v1.2.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	nhooyr.io/websocket v1.8.17 // indirect
)
