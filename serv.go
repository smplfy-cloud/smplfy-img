package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/chai2010/webp"
)

func main() {
	http.HandleFunc("/", handleServer)
	err := http.ListenAndServe(":3232", nil)
	if err != nil {
		return
	}
}

type readImageUrlType struct {
	protocol string
	host     string
	fileName string
	fileType string
}

func readImageUrl(imageUrl string) readImageUrlType {
	imageUrlParsed, err := url.Parse(imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	imageUrlRegexPattern := regexp.MustCompile(`(?m)([^\/]+)(\.\w+$)`)
	parsedUrl := imageUrlRegexPattern.FindString(imageUrlParsed.String())
	splitUrl := strings.Split(parsedUrl, ".")

	return readImageUrlType{
		protocol: imageUrlParsed.Scheme,
		host:     imageUrlParsed.Hostname(),
		fileName: splitUrl[0],
		fileType: splitUrl[1],
	}
}

func downloadImage(imageUrl string) []byte {
	var buf bytes.Buffer
	var data []byte
	var err error

	response, e := http.Get(imageUrl)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	// Load file data
	if data, err = ioutil.ReadAll(response.Body); err != nil {
		log.Println(err)
	}

	// Decode webp
	m, text, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Got string format %s", text)

	// Encode lossless webp
	if err = webp.Encode(&buf, m, &webp.Options{Lossless: false}); err != nil {
		log.Println(err)
	}

	if err = ioutil.WriteFile(getImageFilePath(imageUrl), buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	return buf.Bytes()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func returnImage(w http.ResponseWriter, fileBytes []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}

func getImageFilePath(getImageUrl string) string {
	imageUrl := readImageUrl(getImageUrl)
	return "warehouse/" + imageUrl.fileName + ".webp"
}

func handleServer(w http.ResponseWriter, r *http.Request) {
	getImageUrl := r.URL.Query().Get("url")
	if getImageUrl != "" {
		var fileBytes []byte
		imageFilePath := getImageFilePath(getImageUrl)
		isExtis := fileExists(imageFilePath)
		if isExtis {
			fileBytes, _ = getFile(imageFilePath)
		} else {
			fileBytes = downloadImage(getImageUrl)
		}
		returnImage(w, fileBytes)
	} else {
		fmt.Fprintf(w, "Url mevcut değil.")
	}
}
