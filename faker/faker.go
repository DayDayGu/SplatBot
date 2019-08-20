// Package faker provides image download
package faker

import (
	"image"
	"image/png"
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
}

func Exist(name string) bool {
	path, err := os.Getwd()
	if err != nil {
		return false
	}
	file := path + "/tmp/" + name
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Get(name string) image.Image {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file := path + "/tmp/" + name
	buf, err := os.Open(file)
	defer buf.Close()
	if err != nil {
		return nil
	}
	img, _ := png.Decode(buf)
	return img

}
func GetResrc(name string) image.Image {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file := path + "/resource/" + name
	buf, err := os.Open(file)
	defer buf.Close()
	if err != nil {
		return nil
	}
	img, _ := png.Decode(buf)
	return img

}
