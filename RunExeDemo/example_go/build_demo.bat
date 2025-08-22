SET CGO_ENABLED=0
SET GOOS=android
SET GOARCH=arm64
go build -buildmode=pie -ldflags="-s -w" -o ..\app\src\main\jniLibs\arm64-v8a\app_arm64.so
pause