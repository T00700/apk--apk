package editor

import (
	"bytes"
	"errors"
	"github.com/pzx521521/apk-editor/editor/signv2"
	"github.com/pzx521521/apk-editor/editor/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const ASSETS_DIR = "assets/"

type Manifest struct {
	VersionCode uint32
	VersionName string
	Label       string
	Package     string
}

var DefaultManifest = &Manifest{
	VersionCode: 111,
	VersionName: "111.111.111",
	Label:       "WebViewDemo",
	Package:     "com.parap.webview",
}

func (m *Manifest) Modify(manifest []byte) ([]byte, error) {
	var ms []ModifyInfo[string]
	if m.Label != "" && m.Label != DefaultManifest.Label {
		ms = append(ms, ModifyInfo[string]{DefaultManifest.Label, m.Label})
	}
	if m.Package != "" && m.Package != DefaultManifest.Package {
		ms = append(ms, ModifyInfo[string]{DefaultManifest.Package, m.Package})
	}
	if m.VersionName != "" && m.VersionName != DefaultManifest.VersionName {
		ms = append(ms, ModifyInfo[string]{DefaultManifest.VersionName, m.VersionName})
	}
	if m.VersionCode != 0 && m.VersionCode != DefaultManifest.VersionCode {
		return Modify(manifest, ModifyInfo[uint32]{DefaultManifest.VersionCode, m.VersionCode}, ms)
	} else {
		return Modify(manifest, ms)
	}
}

type MergeEntry struct {
	Name string
	Data []byte
}

type ApkEditor struct {
	Url       string    `json:"url,omitempty"`
	IndexHtml []byte    `json:"index_html,omitempty"`
	HtmlZip   []byte    `json:"html_zip,omitempty"`
	Manifest  *Manifest `json:"manifest,omitempty"`
	apkRaw    []byte
	keyBytes  []byte
	certBytes []byte
	apkSign   *signv2.ApkSign
}

func NewApkEditor(apk, keyBytes, certBytes []byte) *ApkEditor {
	return &ApkEditor{apkRaw: apk, keyBytes: keyBytes, certBytes: certBytes}
}

func (a *ApkEditor) Edit() ([]byte, error) {
	modifyContent, err := a.modifyContent()
	if err != nil {
		return nil, err
	}
	if len(modifyContent) == 0 {
		return nil, errors.New("no content to modify")
	}
	r, err := zip.NewReader(bytes.NewReader(a.apkRaw), int64(len(a.apkRaw)))
	if err != nil {
		return nil, err
	}
	aBuf := new(bytes.Buffer)
	aBuf.Write(a.apkRaw[:r.AppendOffset()])
	w := r.Append(aBuf, a.Manifest != nil)
	err = merge(w, modifyContent...)
	if err != nil {
		return nil, err
	}
	err = a.manifest(r, w)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return sign(aBuf.Bytes(), a.keyBytes, a.certBytes)
}
func (a *ApkEditor) modifyContent() ([]*MergeEntry, error) {
	var mergeEntries []*MergeEntry
	if a.Url != "" {
		if strings.HasPrefix(a.Url, "http") {
			mergeEntries = append(mergeEntries, &MergeEntry{ASSETS_DIR + "url.txt", []byte(a.Url)})
		} else {
			c, err := dirContent(filepath.Clean(a.Url))
			if err != nil {
				return nil, err
			}
			mergeEntries = c
		}
	} else if a.IndexHtml != nil && len(a.IndexHtml) > 0 {
		mergeEntries = append(mergeEntries, &MergeEntry{ASSETS_DIR + "index.html", []byte(a.IndexHtml)})
	} else if a.HtmlZip != nil && len(a.HtmlZip) > 0 {
		content, err := zipContent(a.HtmlZip)
		if err != nil {
			return nil, err
		}
		mergeEntries = append(mergeEntries, content...)
	}
	return mergeEntries, nil
}

func (a *ApkEditor) manifest(r *zip.Reader, w *zip.Writer) error {
	manifest, err := readManifest(r)
	if err != nil {
		return err
	}

	err = merge(w, &MergeEntry{zip.ANDROIDMANIFEST, manifest})
	if err != nil {
		return err
	}
	return nil
}
func zipContent(zipData []byte) ([]*MergeEntry, error) {
	var mergeEntries []*MergeEntry
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}
	for _, f := range reader.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		err = rc.Close()
		if err != nil {
			return nil, err
		}
		mergeEntries = append(mergeEntries, &MergeEntry{ASSETS_DIR + f.Name, b})
	}
	return mergeEntries, nil
}
func dirContent(dir string) ([]*MergeEntry, error) {
	mergeEntrys := []*MergeEntry{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			path, _ = filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			mergeEntrys = append(mergeEntrys, &MergeEntry{ASSETS_DIR + path, file})
		}
		return nil
	})
	return mergeEntrys, nil
}
func sign(apk, keyBytes, certBytes []byte) ([]byte, error) {
	var keys = []*signv2.SigningCert{
		{SigningKey: signv2.SigningKey{
			KeyBytes: keyBytes,
			Type:     signv2.RSA,
			Hash:     signv2.SHA256,
		},
			CertBytes: certBytes,
		},
	}
	z, err := signv2.NewApkSign(apk)
	if err != nil {
		return nil, err
	}
	return z.SignV2(keys)
}
func merge(w *zip.Writer, mf ...*MergeEntry) error {
	for _, file := range mf {
		header := &zip.FileHeader{
			Name:   file.Name,
			Method: zip.Deflate,
		}
		header.SetMode(0o666)
		f, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = f.Write(file.Data)
		if err != nil {
			return err
		}
	}
	return nil
}
func readManifest(r *zip.Reader) ([]byte, error) {
	//读取源数据
	for _, f := range r.File {
		if f.Name == zip.ANDROIDMANIFEST {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			b, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}
			err = rc.Close()
			if err != nil {
				return nil, err
			}
			return b, nil
		}
	}
	return nil, errors.New("no manifest found")
}
