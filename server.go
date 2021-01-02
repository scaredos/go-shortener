package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var urls []string   // URLs to forward to
var shorts []string // Shortened URLs
var port string     // Port to run service on
var shorten string

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func index(slice []string, item string) int {
	for i, _ := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func banner() {
	fmt.Printf("go-shortener - lightweight url shortener\n")
	fmt.Printf("usage: ./%s <options>\n", os.Args[0])
	fmt.Println("---\toptions\t---")
	fmt.Println("-p, --port\t Port for service to run on")
}

func main() {
	if len(os.Args) <= 1 {
		banner()
		os.Exit(0)
	}
	for i, ch := range os.Args {
		if i == 0 {
			continue
		}
		if strings.Contains(ch, "-p") || strings.Contains(ch, "--port") {
			port = os.Args[i+1]
		}
	}
	fmt.Printf("Running on :%s\n", port)
	http.HandleFunc("/", urlShortener)
	http.HandleFunc("/api/v1/addUrl", addUrl)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

// This checks the slice shorts for the shortened URL
// If it exists, it then redirects user to the URL from slice urls
func urlShortener(w http.ResponseWriter, r *http.Request) {
	if index(shorts, r.URL.Path[1:]) == -1 {
		fmt.Fprintf(w, "404 not found")
	} else {
		ii := index(shorts, r.URL.Path[1:])
		http.Redirect(w, r, urls[ii], http.StatusPermanentRedirect)
	}
}

// API for adding URLs to the shortening service
// host:port/api/v1/addUrl?url=
// You can also set custom shortened URLs by specifying parameter short
func addUrl(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"]
	short := r.URL.Query()["short"]
	if len(url) < 1 {
		fmt.Fprintf(w, "you must supply the parameter 'url'\n")
		return
	}
	if len(short) < 1 {
		shorten = RandStringBytes(7)
	} else {
		shorten = short[0]
	}
	shorts = append(shorts, shorten)
	urls = append(urls, url[0])
	fmt.Fprintf(w, fmt.Sprintf("now forwarding %s/%s to %s\n", r.Host, shorten, url[0]))
}
