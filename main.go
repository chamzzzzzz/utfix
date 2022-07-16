package main

import (
	"bytes"
	"fmt"
	chardetect "github.com/djykissyou/chardetect"
	"io"
	"os"
)

var (
	destEnc = "utf-8"
)

func DetectEncoding(data []byte) string {
	enc := chardetect.Mostlike(data)
	if enc == "gbk" || enc == "utf-16be" || enc == "utf-16le" {
		enc = "gbk"
	}
	return enc
}

func UTFix(filePath string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	srcBuf := bytes.NewBuffer(src)
	src, err = io.ReadAll(srcBuf)
	if err != nil {
		return err
	}

	srcEnc := DetectEncoding(src)
	if srcEnc != destEnc {
		r, err := chardetect.NewReader(srcBuf, srcEnc, src)
		if err != nil {
			return err
		}

		var destBuf bytes.Buffer
		w, err := chardetect.NeWriter(&destBuf, destEnc, false)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}

		if err = os.WriteFile(filePath, destBuf.Bytes(), 0666); err != nil {
			return err
		}
		fmt.Printf("%s %s -> %s\n", filePath, srcEnc, destEnc)
	}
	return nil
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("usage: utfix file")
		os.Exit(1)
	}

	for _, filePath := range os.Args[1:] {
		if err := UTFix(filePath); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}
