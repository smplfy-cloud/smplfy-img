package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	http.HandleFunc("/", handleServer)
	err := http.ListenAndServe(":3232", nil)
	if err != nil {
		return
	}
}

type imageParse struct {
	protocol string
	host     string
	fileName string
	fileType string
}

func getImageUrlParse(imageUrl string) imageParse {
	imageUrlParse, err := url.Parse(imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	imageUrlRegexPattern := regexp.MustCompile(`(?m)([^\/]+)(\.\w+$)`)
	parsedUrl := imageUrlRegexPattern.FindString(imageUrlParse.String())
	splitUrl := strings.Split(parsedUrl, ".")

	return imageParse{
		protocol: imageUrlParse.Scheme,
		host:     imageUrlParse.Hostname(),
		fileName: splitUrl[0],
		fileType: splitUrl[1],
	}
}

func getImageUrlSyntaxValidation(imageUrl string) bool {
	golangciYamlExample, _ := os.ReadFile("config.yml")

	fmt.Println("protocol : ", golangciYamlExample)

	//getImageUrl := getImageUrlParse(imageUrl)

	return true
}

func downloadImage(imageUrl string) bool {
	response, e := http.Get(imageUrl)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	_url := getImageUrlParse(imageUrl)

	fmt.Println("protocol : ", _url.protocol)
	fmt.Println("host : ", _url.host)
	fmt.Println("fileName : ", _url.fileName)
	fmt.Println("fileType : ", _url.fileType)

	//open a file for writing
	file, err := os.Create("warehouse/asdf.jpg")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func handleServer(w http.ResponseWriter, r *http.Request) {
	getImageUrl := r.URL.Query().Get("url")
	if getImageUrl != "" {
		isOk := downloadImage(getImageUrl)
		if isOk {
			fmt.Fprintf(w, "Resim indirildi.")
		} else {
			fmt.Fprintf(w, "Resim indirilemedi.")
		}
	} else {
		fmt.Fprintf(w, "Url mevcut değil.")
	}
}
