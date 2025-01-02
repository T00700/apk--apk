package app

import (
	"fmt"
	webview "github.com/webview/webview_go"
	"log"
	"runtime"
)

func main() {
	w := webview.New(true)
	w.SetTitle("Bind Example")
	w.SetSize(480*2, 320*2, webview.HintNone)
	defer w.Destroy()
	run(w)
}
func run(w webview.WebView) {
	tls := runtime.GOOS != "windows"
	server, addr, err := runHttp(tls)
	if err != nil {
		return
	}
	bind(w)
	scheme := "http"
	if tls {
		scheme += "s"
	}
	w.Navigate(fmt.Sprintf("%s://localhost:%d/", scheme, addr.Port))
	w.Run()
	server.Close()
}

func bind(w webview.WebView) {
	// A binding that increments a value and immediately returns the new value.
	w.Bind("log", func(logText string) {
		log.Println(logText)
	})
}
