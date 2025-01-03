package main

import (
	"embed"
	"flag"
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
	versionCode := flag.Int("versionCode", 111, "应用的版本代码 (111)")
	versionName := flag.String("versionName", "111.111.111", "应用的版本名称 (111.111.111)")
	label := flag.String("label", "WebViewDemo", "应用的标签 (WebViewDemo)")
	packageName := flag.String("package", "com.parap.webview", "应用的包名 (com.parap.webview)")
	output := flag.String("o", "webview.apk", "输出文件路径")
	// 解析命令行参数
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		app := filepath.Base(os.Args[0])
		log.Printf("Usage: %s https://www.example.com\n", app)
		log.Printf("or:    %s <yourpath>/index.html\n", app)
		log.Printf("or:    %s <your-dir>\n", app)
		log.Printf("or:    %s <your-dir>/demo.zip\n", app)
		log.Printf("or:    %s <your-dir>/demo.apk\n", app)
		return
	}
	inputPath := args[0]
	abs, err := filepath.Abs(*output)
	checkErr(err)
	crt, err := embedFiles.ReadFile("release/signing.crt")
	checkErr(err)
	var apk []byte
	if filepath.Ext(inputPath) == ".apk" {
		apk, err = os.ReadFile(inputPath)
		checkErr(err)
	} else {
		apk, err = embedFiles.ReadFile("release/app-release.apk")
		checkErr(err)
	}
	key, err := embedFiles.ReadFile("release/signing.key")
	checkErr(err)
	apkEditor := editor.NewApkEditor(apk, key, crt)
	stat, err := os.Stat(inputPath)
	if os.IsNotExist(err) {
		if strings.HasPrefix(inputPath, "http") {
			apkEditor.Url = inputPath
		} else {
			log.Println("file '" + inputPath + "' does not exist")
			return
		}
	} else {
		if stat.IsDir() {
			apkEditor.Url = inputPath
		} else {
			file, err := os.ReadFile(inputPath)
			checkErr(err)
			if strings.HasSuffix(inputPath, ".zip") {
				apkEditor.HtmlZip = file
			} else {
				apkEditor.IndexHtml = file
			}
		}
	}
	if len(*versionName+*label+*packageName) > 0 || *versionCode > 0 {
		apkEditor.Manifest = &editor.Manifest{
			VersionCode: uint32(*versionCode),
			VersionName: *versionName,
			Label:       *label,
			Package:     *packageName,
		}
	}
	edit, err := apkEditor.Edit()
	checkErr(err)
	err = os.WriteFile(abs, edit, 0644)
	checkErr(err)
	log.Printf("success save at:%s\n", abs)
}
