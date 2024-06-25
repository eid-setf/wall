package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

type Data struct {
	Id        string
	Url       string
	Short_url string
	Path      string
}

type ImageJson struct {
	Data []Data
}

func makeQueryUrl(searchUrl, query string) string {
	return searchUrl + "?q=" + query
}

func getImagePaths(body io.ReadCloser, num int) []string {
	var images ImageJson
	if err := 	json.NewDecoder(body).Decode(&images); err != nil {
		panic(err)
	}

	var result = make([]string, 0, num)
	for i := 0; i < num; i++ {
		result = append(result, images.Data[i].Path)
	}

	return result
}

func MakeFilename(url string) string {
	re := regexp.MustCompile("/wallhaven-.*\\.(png|jpg|jpeg)")
	result := re.Find([]byte(url))
	result = result[1:]
	return string(result)
}

func DownloadImage(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	filename := MakeFilename(url)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	io.Copy(file, resp.Body)
}

func main() {
	var searchUrl = "https://wallhaven.cc/api/v1/search"
	var numberOfImages = 10
	var query = "luffy"

	queryUrl := makeQueryUrl(searchUrl, query)

	resp, err := http.Get(queryUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	imagePaths := getImagePaths(resp.Body, numberOfImages)

	for _, v := range imagePaths {
		fmt.Printf("Downloading %v\n", v)
		DownloadImage(v)
	}
}
