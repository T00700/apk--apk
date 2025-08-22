# Android 运行原生二进制文件示例
这个项目演示了如何在Android应用中运行交叉编译的原生二进制文件。碰到的坑记录一下

## 编写简单程序 交叉编译并且adb 测试
编译:
```shell
SET CGO_ENABLED=0
SET GOOS=android
SET GOARCH=arm64
go build -buildmode=pie -ldflags="-s -w" -o ..\app\src\main\jniLibs\arm64-v8a\app_arm64.so
```

adb push 到 `/data/local/tmp`(只有此目录能运行二进制文件) 注意加权限
```shell
adb push app\src\main\jniLibs\arm64-v8a\app_arm64.so /data/local/tmp
adb shell chmod 777 /data/local/tmp/app_arm64.so 
adb shell /data/local/tmp/app_arm64.so
```
## android app 调用二进制文件
+ 普通应用（app uid，SELinux domain=untrusted_app）默认无法访问/data/local/tmp（该目录属主shell:shell，一般771），没有读写权限，只有 adb shell（uid=shell）或 root 才能用
+ Android 5.0+ 必须 PIE：-buildmode=pie
+ 从 Android 10 (API level 29) 开始，Google 加强了安全策略，明确禁止从应用数据目录 (/data/data/... 或 /data/user/0/...) 直接执行二进制文件。[execute-permission](https://developer.android.com/about/versions/10/behavior-changes-10?utm_source=chatgpt.com&hl=zh-cn#execute-permission)[stackoverflow](https://stackoverflow.com/questions/62391811/android-10-alternative-to-launching-executable-as-subproccess-stored-in-app-ho?utm_source=chatgpt.com)
+ 执行报错:`error=13, Permission denied`
+ 需要从 nativeLibraryDir 运行 libexample_arm64.so(为了打包方便加了后缀.so)

## 其他问题
### golang编译到android DNS 解析出问题
Post "xxx": dial tcp: lookup h5.if.qidian.com on [::1]:53: dial udp [::1]:53: socket: operation not permitted  
原因是在 Android 上交叉编译的 Go 程序默认用 netdns=go，它会尝试直接用 127.0.0.1 或 [::1] 上的 DNS，结果因为 Android 系统上根本没监听这个端口，所以解析失败  
最简单的方法是 自定义dns解析  
