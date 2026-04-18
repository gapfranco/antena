//go:build !headless

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	webview "github.com/webview/webview_go"
)

func runDesktop(app *application, addr string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Printf("failed to bind desktop port: %v — falling back to server mode", err)
		runServer(app, addr)
		return
	}
	port := ln.Addr().(*net.TCPAddr).Port

	srv := &http.Server{
		Handler:  app.routes(),
		ErrorLog: log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("desktop server error: %v", err)
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	log.Printf("desktop server listening on %s", url)

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Antena — Monitoramento de Alarmes")
	w.SetSize(1280, 800, webview.HintNone)
	w.Navigate(url)
	w.Run()
}