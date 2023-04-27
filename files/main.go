package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const fileName = "changing_file.txt"

var fileMutex sync.Mutex

type Job struct {
	id       int
	filename string
	m        *sync.Mutex
	wg       *sync.WaitGroup
}

func NewJob(id int, filename string, m *sync.Mutex, wg *sync.WaitGroup) *Job {
	return &Job{
		id:       id,
		filename: fileName,
		m:        m,
		wg:       wg,
	}
}

func (j *Job) run() {
	log.Printf("\t%d strated\n", j.id)

	defer log.Printf("\t%d finished\n", j.id)

	// Ensure the WaitGroup counter is decremented when the function returns
	defer j.wg.Done()

	content := fmt.Sprintf("Entry %d", j.id)

	for i := 0; i < 10; i++ {
		c := fmt.Sprintf("%s_%s\n", content, time.Now().Format(time.RFC3339))

		if err := writeFile(c, j.filename, j.id); err != nil {
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func writeFile(content string, filename string, jobid int) error {
	log.Printf("\t\t%d Acquired lock\n", jobid)
	// Acquire the mutex before opening and writing to the file
	fileMutex.Lock()
	defer func() {
		fileMutex.Unlock()

		log.Printf("\t\t%d Returned lock\n", jobid)
	}()

	// Open the file in append mode or create it if it does not exist
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	defer f.Close()

	// Write the content to the file
	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("write string: %w", err)
	}

	return nil
}

func main() {
	// Create a WaitGroup to track concurrent writes
	var wg sync.WaitGroup

	fpath := filepath.Join("testdata", fileName)

	// Simulate constant changes by writing to the file every second
	for i := 0; i < 10; i++ {
		wg.Add(1)

		j := NewJob(i+1, fpath, &fileMutex, &wg)

		go j.run()
	}

	// Wait for all writes to complete
	wg.Wait()

	// Read and print the file contents
	fileData, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("File contents:")
	fmt.Println(string(fileData))
}
