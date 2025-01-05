uiBuild:
	go build  -ldflags "-s -w" -o bin/webview.app/Contents/MacOS/webview
uiBuild-windows:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H windowsgui" -o bin/webview.exe
