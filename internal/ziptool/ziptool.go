package ziptool

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"log"
)

func Unbase64(str string) []byte {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func DecompressData(data []byte) []byte {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
    if err != nil {
        log.Fatal(err)
    }
	filesCount := len(zipReader.File)
	if filesCount > 1 {
		log.Fatalln("More then 1 file in zip archive!")
	}
	zipFile := zipReader.File[0]
	log.Println("Decompress file: ", zipFile.Name)
	f, err := zipFile.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	unzippedDataBytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	return unzippedDataBytes
}