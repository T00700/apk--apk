rm:
	rm -f release/signed.apk.idsig && rm -f release/notsigned.apk&& rm -f release/notsigned_my.apk && rm -f release/not4.apk && rm -f release/signed.apk


build:
	rm -f release/signed4.apk && rm -f release/notsigned.apk&& rm -f release/notsigned_my.apk && rm -f release/not4.apk && rm -f release/signed.apk
	go run main.go -o release/not4.apk release/app-release.apk

	#go run zipalign.go  release/not4.apk  release/notsigned.apk
	/Users/parapeng/Library/Android/sdk/build-tools/30.0.3/zipalign  -v 4 release/not4.apk release/notsigned.apk
	#/Users/parapeng/Library/Android/sdk/build-tools/30.0.3/zipalign -c -v 4 release/notsigned.apk
	#/Users/parapeng/Library/Android/sdk/build-tools/30.0.3/apksigner sign --ks test.keystore --ks-key-alias "key0" --ks-pass pass:123456 --key-pass pass:123456 --out release/signed.apk  release/notsigned.apk
	#adb install ./release/signed.apk

install:
	adb install release/signed.apk