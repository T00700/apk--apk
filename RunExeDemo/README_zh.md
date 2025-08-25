# RunExeDemo（Android）

一个用于在 Android 上运行随应用打包的原生可执行文件（.so）的示例应用，并实时流式展示其输出。应用内提供 JSON 配置编辑器，通过底部导航在 Run 与 Config 两个页面间切换。

## 功能特性
- 从应用的 `nativeLibraryDir` 运行原生二进制（通过 `app/src/main/jniLibs/...` 打包进入 APK）。
- 实时读取 stdout/stderr 并显示，支持自动滚动至最新行。
- 输出区域支持长按选择复制。
- 从 `nativeLibraryDir` 自动枚举可运行文件，以下拉框选择（过滤 `libandroidx*` 等系统库）。
- 在应用私有目录保存 JSON 配置（`files/config/config.json`），编辑时自动保存。
- 运行时自动注入 `--config <绝对路径>` 参数传给二进制。

## 目录与关键文件
- `app/src/main/java/com/example/runexedemo/MainActivity.kt`：界面（Run/Config 页）、流式输出、下拉选择、自动滚动。
- `app/src/main/java/com/example/runexedemo/BinaryRunner.kt`：进程启动器（按所选 `.so` 启动，自动注入 `--config`）。
- `app/src/main/java/com/example/runexedemo/ConfigManager.kt`：应用私有目录内 JSON 配置的读写。
- `app/src/main/jniLibs/arm64-v8a/libexample_arm64.so`：示例原生“可执行”（以 .so 命名，便于放入 `nativeLibraryDir`）。
- `example_go/main.go`：示例 Go 程序，演示 `--config` 与逐秒输出。

## 环境要求
- Android Studio（较新版本），与本仓库 Gradle/AGP 版本匹配。
- 设备 ABI：建议 arm64-v8a；minSdk 24。

## 构建与运行
1. 将你的原生二进制（PIE，ET_DYN/可执行）以 `lib<name>.so` 命名，放到：
   - `app/src/main/jniLibs/arm64-v8a/`
2. 在 Android Studio 中 Build & Run。
3. 在应用中：
   - Config 页编辑 `config.json`（自动保存）。
   - Run 页选择要运行的 `.so`，点击 Run；输出会实时滚动显示，可随时 Stop。

## Go 交叉编译快速指引（arm64）
```bash
# Windows PowerShell 示例
$env:CGO_ENABLED="0"
$env:GOOS="android"
$env:GOARCH="arm64"
# 生成 PIE 可执行（以 .so 命名，便于打包到 jniLibs）
go build -buildmode=pie -ldflags "-s -w" -o libexample_arm64.so ./example_go
```
注意：
- Android 5.0+ 要求 PIE。
- 现代设备优先提供 arm64-v8a。

## 网络
- Manifest 已包含 `INTERNET` 权限。
- 若 Go `net/http` 在部分设备上 DNS 解析异常（如退回 `[::1]:53`），建议自定义 `net.Resolver`（如 1.1.1.1/8.8.8.8），或在 App 侧解析后传入 IP。

## 权限与执行限制
- 不能从 `/data/local/tmp` 或外部存储执行；本项目从应用 `nativeLibraryDir` 运行。
- 某些厂商策略可能限制 `filesDir` 的可执行权限，因此采用 `jniLibs` 打包方案更稳妥。

## 常见问题
- EACCES（权限被拒）：确认在 `nativeLibraryDir` 执行，且为设备 ABI 的有效 PIE 可执行。
- Exec format error：架构不匹配；请针对 arm64-v8a 重编译。
- No such file：确认 `.so` 已放在 `app/src/main/jniLibs/arm64-v8a/` 并同步到构建。
- DNS 问题：使用自定义 `net.Resolver` 或在程序中支持 `--dns` 之类参数。

## 自定义建议
- 可重命名示例 `.so` 并在 `BinaryRunner.kt` 中调整默认值。
- 下拉过滤已忽略 `libandroidx*`；可按需扩展过滤规则。
- 界面使用 Compose Material3，可在 `NavigationBarItemDefaults.colors` 调整主题色与样式。

## 许可
演示示例项目。可按你的需要添加许可证文件。