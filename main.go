package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	charset   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	urlLength = 6
)

type PageData struct {
	Result        string
	Error         string
	ShortenedUrls map[string]string
	page          string
}

var pageData PageData

func main() {
	pageData.ShortenedUrls = make(map[string]string)

	page, _ := os.ReadFile("index.html")
	pageData.page = string(page)

	http.HandleFunc("/", handleIndexAndRedirects)

	http.HandleFunc("/submit", handleUrlShorteningRequest)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func handleIndexAndRedirects(writer http.ResponseWriter, request *http.Request) {
	shortUrl := strings.TrimPrefix(request.URL.Path, "/")

	if shortUrl != "" {
		handleWebPageRedirect(writer, request, shortUrl)
	} else {
		displayMainPage(writer)
	}
}

func handleWebPageRedirect(writer http.ResponseWriter, request *http.Request, shortUrl string) {
	value, ok := pageData.ShortenedUrls[shortUrl]
	if !ok {
		pageData.Error = "Shortened URL (" + shortUrl + ") not found"
		defer func() { pageData.Error = "" }()
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	} else {
		http.Redirect(writer, request, value, http.StatusMovedPermanently)
	}
}

func displayMainPage(writer http.ResponseWriter) {
	templ, _ := template.New("index").Parse(pageData.page)
	err := templ.Execute(writer, pageData)
	if err != nil {
		return
	}
}

func handleUrlShorteningRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		url := extendUrlIfNeeded(request.FormValue("url"))
		shortUrl := getUniqueShortUrl()
		pageData.Result = "localhost:8080/" + shortUrl
		pageData.ShortenedUrls[shortUrl] = url
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	} else {
		http.Error(writer, "Unsupported method", http.StatusBadRequest)
	}
}

func getUniqueShortUrl() string {
	for {
		s := randomString()
		_, found := pageData.ShortenedUrls[s]
		if !found {
			return s
		}
	}
}

func extendUrlIfNeeded(value string) string {
	var prefix string
	if !strings.HasPrefix(value, "http") {
		prefix = "https://"
	}
	return prefix + value
}

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, urlLength)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}
