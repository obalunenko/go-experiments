package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	testDataDir = "testdata"
	ext         = ".zip"
	baseName    = "ambank_audit_files"
)

// compressFiles Compress csv file into zip archive
func compressFiles(files []file, datePostfix string) (file, error) {

	if len(files) == 0 {
		return file{}, errors.New("no files to compress")
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zw := zip.NewWriter(buf)

	for _, f := range files {
		fw, err := zw.Create(f.name)
		if err != nil {
			return file{}, errors.New("failed to create file in archive")
		}

		_, err = fw.Write(f.body)
		if err != nil {
			return file{}, errors.New("failed to write content of file to archive")
		}
	}

	err := zw.Close()
	if err != nil {
		return file{}, errors.New("failed to close archive writer")
	}

	return file{
		name: baseName + datePostfix + ext,
		body: buf.Bytes(),
	}, nil
}

type file struct {
	name string
	body []byte
}

// openTestFile reads testdata file and return a file object
func openTestFile(path string) (file, error) {
	path = filepath.Join(testDataDir, path)
	f, err := os.Open(path)
	if err != nil {
		return file{}, err
	}

	info, err := f.Stat()
	if err != nil {
		return file{}, err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return file{}, err
	}

	return file{
		name: info.Name(),
		body: content,
	}, nil

}

func main() {
	testFiles := []string{"file1.txt", "file2.txt"}
	var input []file
	for _, fPath := range testFiles {
		f, err := openTestFile(fPath)
		if err != nil {
			log.Fatal(err)
		}
		input = append(input, f)

	}
	arch, err := compressFiles(input, "_golden")
	if err != nil {
		log.Fatal(err)
	}
	osF, err := os.Create(filepath.Join(testDataDir, arch.name))
	if err != nil {
		log.Fatal(err)
	}
	_, err = osF.Write(arch.body)
	if err != nil {
		log.Fatal(err)
	}

	if err = osF.Close(); err != nil {
		log.Fatal(err)
	}

}
