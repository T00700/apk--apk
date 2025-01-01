package signv2

import (
	"os"
	"testing"
)

const (
	keyPath  = "signing.key"
	certPath = "signing.crt"
)

var keys = []*SigningCert{
	{SigningKey: SigningKey{
		KeyPath: keyPath,
		Type:    RSA,
		Hash:    SHA256,
	},
		CertPath: certPath,
	},
}

func loadFile(name string) ([]byte, error) {
	var err error
	var b []byte
	if b, err = os.ReadFile(name); err != nil {
		return nil, err
	}
	return b, err
}

func saveFile(name string, b []byte) error {
	var f *os.File
	var err error
	if f, err = os.Create(name); err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	if _, err = f.Write(b); err != nil {
		return err
	}
	return nil
}

func TestGenerateSignedFile(t *testing.T) {
	var z *ApkSign
	var err error
	var b []byte

	if b, err = loadFile("/Users/parapeng/Desktop/apkEditor/release/notsigned.apk"); err != nil {
		t.Log("error loading file", err)
		t.FailNow()
	}
	if z, err = NewZip(b); err != nil {
		t.Log("error parsing zip", err)
		t.FailNow()
	}
	if z, err = z.SignV2(keys); err != nil {
		t.Log("error signing zip", err)
		t.FailNow()
	}
	if err = saveFile("/Users/parapeng/Desktop/apkEditor/release/signed.apk", z.Bytes()); err != nil {
		t.Error("error signing zip", err)
	}
}
