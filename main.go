package main

import (
	"embed"
	"github.com/pzx521521/apk-editor/editor"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed release/*
var embedFiles embed.FS

func main() {
	checkErr := func(err error) {
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	}
	if len(os.Args) != 2 {
		app := filepath.Base(os.Args[0])
		log.Printf("Usage: %s https://www.example.com\n", app)
		log.Printf("or:    %s <yourpath>/index.html\n", app)
		log.Printf("or:    %s <your-dir>\n", app)
		log.Printf("or:    %s <your-dir>/demo.zip\n", app)
		return
	}
	crt, err := embedFiles.ReadFile("release/signing.crt")
	checkErr(err)
	apk, err := embedFiles.ReadFile("release/app-release.apk")
	checkErr(err)
	key, err := embedFiles.ReadFile("release/signing.key")
	checkErr(err)
	apkEditor := editor.NewApkEditor(apk, key, crt)
	stat, err := os.Stat(os.Args[1])
	if os.IsNotExist(err) {
		if strings.HasPrefix(os.Args[1], "http") {
			apkEditor.Url = os.Args[1]
		} else {
			log.Println("file '" + os.Args[1] + "' does not exist")
			return
		}
	} else {
		if stat.IsDir() {
			apkEditor.Url = os.Args[1]
		} else {
			file, err := os.ReadFile(os.Args[1])
			checkErr(err)
			if strings.HasSuffix(os.Args[1], ".zip") {
				apkEditor.HtmlZip = file
			} else {
				apkEditor.IndexHtml = file
			}
		}
	}

	edit, err := apkEditor.Edit()
	checkErr(err)
	abs, err := filepath.Abs("/Users/parapeng/Downloads/webview.apk")
	checkErr(err)
	err = os.WriteFile(abs, edit, 0644)
	checkErr(err)
	log.Printf("success save at:%s\n", abs)
}
