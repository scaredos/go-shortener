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
var passes []string // Passwords to edit URLs
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
	for i := range slice {
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
	http.HandleFunc("/api/v1/editUrl", editUrl)
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
	passw := r.URL.Query()["password"]
	if len(url) < 1 {
		fmt.Fprintf(w, "{'error': 'url parameter required', 'status': 'fail'}")
		return
	}
	if len(passw) < 1 {
		fmt.Fprintf(w, "{'error': ' password parameter required', 'status': 'fail'}")
		return
	}
	if len(short) < 1 {
		shorten = RandStringBytes(7)
	} else {
		shorten = short[0]
	}
	if index(shorts, shorten) != -1 {
		fmt.Fprintf(w, fmt.Sprintf("{'error': '%s already exists', 'status': 'fail'}", shorten))
		return
	}
	shorts = append(shorts, shorten)
	urls = append(urls, url[0])
	passes = append(passes, passw[0])
	fmt.Fprintf(w, fmt.Sprintf("{'url': '%s', 'short_link': 'http://%s/%s', 'password': '%s', 'status': 'good'}", url[0], r.Host, shorten, passw[0]))
}

// API for editng destinations of the shortened URL
// host:port/api/v1/editUrl?url=&short=
// You must specify the new url as "url" and the shortened url as "short"
// Example: http://domain.com:1337/api/v1/editUrl?url=https://github.com/&short=jWhhgeM
func editUrl(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"]
	short := r.URL.Query()["short"]
	passw := r.URL.Query()["password"]
	if len(url) < 1 {
		fmt.Fprintf(w, "{'error': 'url parameter required', 'status': 'fail'}\n")
		return
	}
	if len(short) < 1 {
		fmt.Fprintf(w, "{'error': 'short parameter required', 'status': 'fail'}\n")
		return
	}
	if len(passw) < 1 {
		fmt.Fprintf(w, "{'error': 'password parameter required', 'status': 'fail'}")
		return
	}
	if index(shorts, short[0]) == -1 {
		fmt.Fprintf(w, "{'error': 'short does not exist', 'status': 'fail'}\n")
		return
	}
	shortPlace := index(shorts, short[0])
	if index(passes, passw[0]) != shortPlace {
		fmt.Fprintf(w, "{'error': 'password incorrect', 'status': 'fail'}")
		return
	}
	urls[shortPlace] = url[0]
	fmt.Fprintf(w, fmt.Sprintf("{'url': '%s', 'short_link': 'http://%s/%s', 'status': 'good'}", url[0], r.Host, shorten))
}
