package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
)

// run this fucntion will the server is running to test the uncompression on the server
func test_compression() {
	file := "test.json"
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Falied to read: %s", err)
	}

	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	if _, err := gzipWriter.Write(data); err != nil {
		fmt.Println("Error writiing to gzip", err)
	}

	if err := gzipWriter.Close(); err != nil {
		fmt.Println("Error closing gzip writer:", err)
		return
	}

	request, err := http.NewRequest("POST", "http://localhost:4000/compressed/drawing", &buffer)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer response.Body.Close()
	fmt.Println("Statuse Code", response.StatusCode)
}
