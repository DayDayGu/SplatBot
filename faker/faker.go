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

func init() {

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path += "/tmp/"
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func TempPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return path + "/tmp/"
}

func resourcePath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return path + "/resource/"
}

// Download download file
func Download(url string, name string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	file, err := os.Create(TempPath() + name)
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

	file := TempPath() + name
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Get(name string) image.Image {
	file := TempPath() + name
	buf, err := os.Open(file)
	defer buf.Close()
	if err != nil {
		return nil
	}
	img, _ := png.Decode(buf)
	return img

}
func GetResrc(name string) image.Image {
	file := resourcePath() + name
	buf, err := os.Open(file)
	defer buf.Close()
	if err != nil {
		return nil
	}
	img, _ := png.Decode(buf)
	return img

}
