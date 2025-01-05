package main

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
	server, _, err := runHttp(tls)
	if err != nil {
		return
	}
	bind(w)
	index := fmt.Sprintf("%s/html-to-apk", vwPort)
	fmt.Printf("%v\n", index)
	w.Navigate(index)
	w.Run()
	server.Close()
}

func bind(w webview.WebView) {
	// A binding that increments a value and immediately returns the new value.
	w.Bind("log", func(logText string) {
		log.Println(logText)
	})

	w.Bind("wvPort", func() string {
		return vwPort
	})
}
