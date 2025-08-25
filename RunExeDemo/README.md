# RunExeDemo (Android)

Android app that runs a bundled native executable (.so) and streams its output in real time. Includes an in-app JSON config editor and a bottom navigation to switch between Run and Config screens.

## Features
- Run native binary from app `nativeLibraryDir` (packaged under `app/src/main/jniLibs/...`).
- Real-time stdout/stderr streaming with auto-scroll to latest line.
- Output text is selectable (long-press to copy).
- Pick which `.so` to run from a dropdown populated from `nativeLibraryDir` (filters out `libandroidx*`).
- JSON config stored in app private storage (`files/config/config.json`), auto-saves while editing.
- Binary invoked with `--config <absolute_path>` automatically.

## Project layout (key parts)
- `app/src/main/java/com/example/runexedemo/MainActivity.kt` – UI (Run/Config tabs, streaming, dropdown, auto-scroll).
- `app/src/main/java/com/example/runexedemo/BinaryRunner.kt` – Process launcher (runs selected `.so`, injects `--config`).
- `app/src/main/java/com/example/runexedemo/ConfigManager.kt` – JSON config read/write in app private dir.
- `app/src/main/jniLibs/arm64-v8a/libexample_arm64.so` – Example native binary (packaged as a .so so it resides in `nativeLibraryDir`).
- `example_go/main.go` – Sample Go program demonstrating `--config`, progressive output.

## Requirements
- Android Studio Giraffe+ (or recent), Android Gradle Plugin per `gradle/libs.versions.toml`.
- Device with arm64-v8a (recommended). Min SDK 24.

## Building & Running
1. Place your native binary as a PIE ET_DYN executable named like `lib<name>.so` under:
   - `app/src/main/jniLibs/arm64-v8a/`
2. Build and run from Android Studio.
3. In the app:
   - Go to Config tab to edit `config.json` (auto-saved).
   - Go to Run tab, pick the `.so`, press Run. Output streams live; press Stop to terminate.

## Go cross-compile quick guide (arm64)
```bash
# Windows PowerShell example
$env:CGO_ENABLED="0"
$env:GOOS="android"
$env:GOARCH="arm64"
# Build a PIE ET_DYN executable (.so name for packaging under jniLibs)
go build -buildmode=pie -ldflags "-s -w" -o libexample_arm64.so ./example_go
```
Notes:
- Android 5.0+ requires PIE.
- Provide arm64-v8a at minimum for modern devices.

## Network
- Manifest includes INTERNET permission.
- If your binary uses Go `net/http` and DNS fails on some devices (e.g., attempts `[::1]:53`), implement a custom `net.Resolver` with a public DNS (1.1.1.1 / 8.8.8.8) or allow the app side to resolve and pass an IP.

## Permissions & Execution constraints
- App cannot execute from `/data/local/tmp` or external storage. Binaries run from the app’s `nativeLibraryDir`.
- Executing from `filesDir` may be blocked on some OEMs; this project avoids that by packaging under `jniLibs`.

## Troubleshooting
- Permission denied (EACCES): ensure you are executing from `nativeLibraryDir` and the binary is a valid PIE executable for the device ABI.
- Exec format error: architecture mismatch; rebuild for arm64-v8a.
- No such file: ensure your `.so` exists in `app/src/main/jniLibs/arm64-v8a/` and is included by Gradle sync.
- DNS errors in Go: use a custom `net.Resolver` or configurable `--dns` option.

## Customization tips
- Rename the example `.so` and adjust default in `BinaryRunner.kt` if desired.
- The dropdown filters out `libandroidx*`; extend the filter to hide other system libs.
- UI built with Compose Material3; tweak colors in `NavigationBarItemDefaults.colors`.

## License
This project is for demonstration purposes. Add your preferred license here.