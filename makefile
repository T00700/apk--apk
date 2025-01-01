build:
	rm -f release/notsigned.apk && rm -f release/not4.apk && rm -f release/signed.apk
	go run main.go -o release/not4.apk release/app-release.apk
	/Users/parapeng/Library/Android/sdk/build-tools/30.0.3/zipalign  -v 4 release/not4.apk release/notsigned.apk
	/Users/parapeng/Library/Android/sdk/build-tools/30.0.3/apksigner sign --ks test.keystore --ks-key-alias "key0" --ks-pass pass:123456 --key-pass pass:123456 --out release/signed.apk  release/notsigned.apk
	adb install ./release/signed.apk

