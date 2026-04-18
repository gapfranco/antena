//go:build headless

package main

func runDesktop(app *application, addr string) {
	runServer(app, addr)
}