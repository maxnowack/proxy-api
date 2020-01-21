package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"h12.io/socks"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", proxyServer)
	http.ListenAndServe(":"+port, nil)
}

func getSocksClient(url string) *http.Client {
	dialSocksProxy := socks.Dial(url)
	tr := &http.Transport{Dial: dialSocksProxy}
	return &http.Client{Transport: tr}
}

func proxyServer(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	proxyURL := r.URL.Query().Get("proxyUrl")

	if url == "" || proxyURL == "" {
		fmt.Fprintln(os.Stderr, "empty options")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	httpClient := getSocksClient(proxyURL)
	// create a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't GET page:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading body:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(b))
}
