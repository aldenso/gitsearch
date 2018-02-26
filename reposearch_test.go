package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"net/http/httptest"
)

var requestsOK = []byte(`{"total_count": 1,
	"incomplete_results": false,
	"items": [
	  {	"name": "pattern",
		"owner": {"login": "myself"},
		"html_url": "https://github.com/klepek/zfssa-zabbix",
		"description": "Zabbix template + scripts for ZFSSA monitoring",
		"language": "Go",
		"stargazers_count": 3 }]}`)

func Test_searchRepo(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(requestsOK)
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
		if r.URL.RequestURI() != "/repositories?q=pattern+fork:true&per_page=100" {
			t.Errorf("Expected request to '/repositories?q=pattern+fork:true&per_page=100', got '%s'", r.URL.RequestURI())
		}
	}))
	defer ts1.Close()
	perpage := "100"
	pattern := "pattern"
	resp1 := searchRepo(ts1.URL+"/", pattern, perpage)
	if resp1.Count != 1 {
		fmt.Printf("Using pattern '%s', language '%s', and login '%s'", pattern, language, login)
		t.Errorf("Count mismatch, expected '1', got '%d'", resp1.Count)
	}
	// test
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(requestsOK)
		if r.URL.RequestURI() != "/repositories?q=pattern+user:aldenso+fork:true&per_page=100" {
			t.Errorf("Expected request to '/repositories?q=pattern+user:aldenso+fork:true&per_page=100', got '%s'", r.URL.RequestURI())
		}
	}))
	defer ts2.Close()
	login = "aldenso"
	searchRepo(ts2.URL+"/", pattern, perpage)
	// test
	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(requestsOK)
		if r.URL.RequestURI() != "/repositories?q=pattern+user:aldenso+fork:true+language:Go&per_page=100" {
			t.Errorf("Expected request to '/repositories?q=pattern+user:aldenso+fork:true+language:Go&per_page=100', got '%s'", r.URL.RequestURI())
		}
	}))
	defer ts3.Close()
	language = "Go"
	searchRepo(ts3.URL+"/", pattern, perpage)
	// test
	ts4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(requestsOK)
		if r.URL.RequestURI() != "/repositories?q=pattern+fork:true+language:Go&per_page=100" {
			t.Errorf("Expected request to '/repositories?q=pattern+fork:true+language:Go&per_page=100', got '%s'", r.URL.RequestURI())
		}
	}))
	defer ts4.Close()
	login = ""
	searchRepo(ts4.URL+"/", pattern, perpage)
}

func Test_showRepoResult(t *testing.T) {
	items := ItemRepo{
		Name:            "MyName",
		Owner:           ItemOwner{Login: "Owner"},
		Description:     "MyDescription",
		HTMLURL:         "http://notfound.com",
		Language:        "Go",
		StargazersCount: 30,
	}
	requests := RespRepo{NextURL: "http://notfound.com/next",
		Count:      1,
		Incomplete: false,
		Items:      []ItemRepo{items},
	}
	result1 := requests.showRepoResult()
	for _, v := range result1.Items {
		if v.ID != 0 || v.HTMLURL != "http://notfound.com" {
			t.Errorf("showRepoResult wrong, expected ID == '0' HTMLURL == 'http://notfound.com', but got '%d' and '%s'", v.ID, v.HTMLURL)
		}
	}
	ospager = "more"
	result2 := requests.showRepoResult()
	for _, v := range result2.Items {
		if v.ID != 0 || v.HTMLURL != "http://notfound.com" {
			t.Errorf("showRepoResult wrong, expected ID == '0' HTMLURL == 'http://notfound.com', but got '%d' and '%s'", v.ID, v.HTMLURL)
		}
	}
}

func Test_continueSearchRepo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(requestsOK)
	}))
	defer ts.Close()
	resp := continueSearchRepo(ts.URL + "/")
	if resp.Count != 1 {
		t.Errorf("Count mismatch, expected '1', got '%d'", resp.Count)
	}
}

func Test_getRepoConfirmNextNotEmpty(t *testing.T) {
	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts0.Close()
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts1.Close()
	items := ItemRepo{
		Name:            "MyName",
		Owner:           ItemOwner{Login: "Owner"},
		Description:     "MyDescription",
		HTMLURL:         "http://notfound.com",
		Language:        "Go",
		StargazersCount: 30,
	}
	requests := RespRepo{NextURL: ts0.URL + "/",
		Count:      1,
		Incomplete: false,
		Items:      []ItemRepo{items},
	}
	choices := ItemChoices{[]ItemChoice{ItemChoice{ID: 0,
		HTMLURL: "https://fakeurl.com"}}}

	test := [][]byte{[]byte("s"), []byte("r"), []byte("n"), []byte("p")}
	for _, x := range test {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(x); err != nil {
			log.Fatal(err)
		}

		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }() // Restore original Stdin

		os.Stdin = tmpfile
		if err := getRepoConfirmNextNotEmpty(&requests, &choices); err != nil {
			t.Errorf("getUserConfirm failed: %v", err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	}
	// test previous url not empty
	requests = RespRepo{PreviousURL: ts0.URL + "/",
		Count:      1,
		Incomplete: false,
		Items:      []ItemRepo{items},
	}
	test = [][]byte{[]byte("p")}
	for _, x := range test {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(x); err != nil {
			log.Fatal(err)
		}

		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }() // Restore original Stdin

		os.Stdin = tmpfile
		if err := getRepoConfirmNextNotEmpty(&requests, &choices); err != nil {
			t.Errorf("getUserConfirm failed: %v", err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func Test_getRepoConfirmNextEmpty(t *testing.T) {
	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts0.Close()
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts1.Close()
	items := ItemRepo{
		Name:            "MyName",
		Owner:           ItemOwner{Login: "Owner"},
		Description:     "MyDescription",
		HTMLURL:         "http://notfound.com",
		Language:        "Go",
		StargazersCount: 30,
	}
	requests := RespRepo{NextURL: ts0.URL + "/",
		Count:      1,
		Incomplete: false,
		Items:      []ItemRepo{items},
	}
	choices := ItemChoices{[]ItemChoice{ItemChoice{ID: 0,
		HTMLURL: "https://fakeurl.com"}}}

	test := [][]byte{[]byte("s"), []byte("r"), []byte("p")}
	for _, x := range test {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(x); err != nil {
			log.Fatal(err)
		}

		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }() // Restore original Stdin

		os.Stdin = tmpfile
		if err := getRepoConfirmNextEmpty(&requests, &choices); err != nil {
			t.Errorf("getUserConfirm failed: %v", err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	}
	// test previous url not empty
	requests = RespRepo{PreviousURL: ts0.URL + "/",
		Count:      1,
		Incomplete: false,
		Items:      []ItemRepo{items},
	}
	test = [][]byte{[]byte("p")}
	for _, x := range test {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(x); err != nil {
			log.Fatal(err)
		}

		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }() // Restore original Stdin

		os.Stdin = tmpfile
		if err := getRepoConfirmNextEmpty(&requests, &choices); err != nil {
			t.Errorf("getUserConfirm failed: %v", err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
