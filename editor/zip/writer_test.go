// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
)

// TODO(adg): a more sophisticated test suite

type WriteTest struct {
	Name   string
	Data   []byte
	Method uint16
	Mode   os.FileMode
}

var writeTests = []WriteTest{
	{
		Name:   "foo.txt",
		Data:   []byte("Rabbits, guinea pigs, gophers, marsupial rats, and quolls."),
		Method: Deflate,
		Mode:   0666,
	},
	{
		Name:   "setuid",
		Data:   []byte("setuid file"),
		Method: Deflate,
		Mode:   0755 | os.ModeSetuid,
	},
	{
		Name:   "setgid",
		Data:   []byte("setgid file"),
		Method: Deflate,
		Mode:   0755 | os.ModeSetgid,
	},
}

func TestWriter(t *testing.T) {
	// write a zip file
	buf := new(bytes.Buffer)
	w := NewWriter(buf)

	for _, wt := range writeTests {
		testCreate(t, w, &wt)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.WriteFile("../release/test.zip", buf.Bytes(), 0666)
	// read it back
	r, err := NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	for i, wt := range writeTests {
		testReadFile(t, r.File[i], &wt)
	}
}

func TestWriterCopy(t *testing.T) {

	f, err := os.Create("../release/test4copy.zip")
	if err != nil {
		log.Fatal(err)
	}
	w := zip.NewWriter(f)

	var files = []struct {
		Name, Body string
	}{
		{"assets/url.txt", "www.baidu.com"},
		{"assets/index.html", "<h1>Hello World Local Html</h1>"},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}

	rc, err := zip.OpenReader("../release/test4.zip")
	if err != nil {
		log.Print(err)
	}
	for _, file := range rc.File {
		if err := w.Copy(file); err != nil {
			log.Printf("copying (%s): %v", file.Name, err)
		}
	}

	if err := w.Close(); err != nil {
		log.Printf("%v\n", err)
	}
}

func TestAppend(t *testing.T) {
	// write a zip file
	buf := new(bytes.Buffer)
	w := NewWriter(buf)

	for _, wt := range writeTests {
		testCreate(t, w, &wt)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	// read it back
	r, err := NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	// append a file to it.
	abuf := new(bytes.Buffer)
	fmt.Printf("%v:%v\n", buf.Len(), r.AppendOffset())
	abuf.Write(buf.Bytes()[:r.AppendOffset()])
	w = r.Append(abuf)

	wt := WriteTest{
		Name:   "append",
		Data:   []byte("Badgers, canines, weasels, owls, and snakes"),
		Method: Deflate,
		Mode:   0755,
	}
	testCreate(t, w, &wt)

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	// read the whole thing back.
	allBytes := abuf.Bytes()
	os.WriteFile("../release/test4.zip", allBytes, 0666)
	r, err = NewReader(bytes.NewReader(allBytes), int64(len(allBytes)))
	if err != nil {
		t.Fatal(err)
	}

	writeTests := append(writeTests, wt)
	for i, wt := range writeTests {
		testReadFile(t, r.File[i], &wt)
	}
}

func TestAppendFile(t *testing.T) {
	checkError := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	inputData, err := os.ReadFile("../release/app-release.apk")
	checkError(err)
	r, err := NewReader(bytes.NewReader(inputData), int64(len(inputData)))
	checkError(err)

	aBuf := new(bytes.Buffer)
	aBuf.Write(inputData[:r.AppendOffset()])
	fmt.Printf("%v:%v\n", len(inputData), r.AppendOffset())
	w := r.Append(aBuf)
	wt := WriteTest{
		Name:   "append.txt",
		Data:   []byte("append, append, append, append, and append"),
		Method: Deflate,
		Mode:   0666,
	}
	testCreate(t, w, &wt)
	w.Close()
	os.WriteFile("../release/test4.zip", aBuf.Bytes(), 0666)
}
func TestWriterOffset(t *testing.T) {
	largeData := make([]byte, 1<<17)
	for i := range largeData {
		largeData[i] = byte(rand.Int())
	}
	writeTests[1].Data = largeData
	defer func() {
		writeTests[1].Data = nil
	}()

	// write a zip file
	buf := new(bytes.Buffer)
	existingData := []byte{1, 2, 3, 1, 2, 3, 1, 2, 3}
	n, _ := buf.Write(existingData)
	w := NewWriter(buf)
	w.SetOffset(int64(n))

	for _, wt := range writeTests {
		testCreate(t, w, &wt)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	// read it back
	r, err := NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	for i, wt := range writeTests {
		testReadFile(t, r.File[i], &wt)
	}
}

func TestWriterFlush(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(struct{ io.Writer }{&buf})
	_, err := w.Create("foo")
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() > 0 {
		t.Fatalf("Unexpected %d bytes already in buffer", buf.Len())
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("No bytes written after Flush")
	}
}

func testCreate(t *testing.T, w *Writer, wt *WriteTest) {
	header := &FileHeader{
		Name:   wt.Name,
		Method: wt.Method,
	}
	if wt.Mode != 0 {
		header.SetMode(wt.Mode)
	}
	f, err := w.CreateHeader(header)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write(wt.Data)
	if err != nil {
		t.Fatal(err)
	}
}

func testReadFile(t *testing.T, f *File, wt *WriteTest) {
	if f.Name != wt.Name {
		t.Fatalf("File name: got %q, want %q", f.Name, wt.Name)
	}
	testFileMode(t, wt.Name, f, wt.Mode)
	rc, err := f.Open()
	if err != nil {
		t.Fatal("opening:", err)
	}
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal("reading:", err)
	}
	err = rc.Close()
	if err != nil {
		t.Fatal("closing:", err)
	}
	if !bytes.Equal(b, wt.Data) {
		t.Errorf("File contents %q, want %q", b, wt.Data)
	}
	iLen := len(wt.Data)
	if iLen > 100 {
		iLen = 100
	}
	fmt.Printf("%s:%s\n", wt.Name, wt.Data[:iLen])
}

func BenchmarkCompressedZipGarbage(b *testing.B) {
	b.ReportAllocs()
	var buf bytes.Buffer
	bigBuf := bytes.Repeat([]byte("a"), 1<<20)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		zw := NewWriter(&buf)
		for j := 0; j < 3; j++ {
			w, _ := zw.CreateHeader(&FileHeader{
				Name:   "foo",
				Method: Deflate,
			})
			w.Write(bigBuf)
		}
		zw.Close()
	}
}
