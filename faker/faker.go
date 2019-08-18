// Package faker provides image download
package faker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Download download file
func Download(url string, name string) {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	file, err := os.Create(path + "/tmp/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
