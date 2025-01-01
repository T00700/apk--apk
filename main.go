// Copyright 2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Zipmerge merges the content of many zip files,
// without decompressing and recompressing the data.
//
// Usage:
//
//	zipmerge [-o out.zip] a.zip b.zip ...
//
// By default, zipmerge appends the content of the second and subsequent zip files
// to the first, rewriting the first in place.
// If the -o option is given, zipmerge creates a new output file containing
// the content of all the zip files, without modifying any of the source zip files.
package main

import (
	"apkEditor/zip"
	"flag"
	"log"
	"os"
)

var outputFile = flag.String("o", "", "write to `file`")

func main() {
	flag.Parse()
	if *outputFile == "" {
		return
	}
	var err error
	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	w := zip.NewWriter(f)

	var files = []struct {
		Name, Body string
	}{
		//{"assets/url.txt", "www.baidu.com"},
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
	for _, name := range flag.Args() {
		rc, err := zip.OpenReader(name)
		if err != nil {
			log.Print(err)
			continue
		}
		for _, file := range rc.File {
			if err := w.Copy(file); err != nil {
				log.Printf("copying from %s (%s): %v", name, file.Name, err)
			}
		}
	}
	if err := w.Close(); err != nil {
		log.Fatal("finishing zip file: %v", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal("finishing zip file: %v", err)
	}
}
