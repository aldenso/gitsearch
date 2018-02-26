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

var userOK = []byte(`{"total_count": 1,
	"items": [
	  {"login": "user1",
	   "url": "https://api.github.com/users/user1",
		"score": 30}]}`)

func Test_searchUser(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
			fmt.Println()
		}
		if r.URL.RequestURI() != "/users?q=user1&per_page=100" {
			t.Errorf("Expected request to '/users?q=user1&per_page=100', got '%s'", r.URL.RequestURI())
		}
	}))
	defer ts1.Close()
	perpage := "100"
	pattern := "user1"
	language = ""
	resp1 := searchUser(ts1.URL+"/", pattern, perpage)
	if resp1.Count != 1 {
		fmt.Printf("Using pattern '%s' and login '%s'", pattern, login)
		t.Errorf("Count mismatch, expected '1', got '%d'", resp1.Count)
	}
	// test
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
			fmt.Println()
		}
		if r.URL.RequestURI() != "/users?q=user1+language:Go&per_page=100" {
			t.Errorf("Expected request to '/users?q=user1+language:Go&per_page=100', got '%s'", r.URL.RequestURI())
		}
	}))
	defer ts2.Close()
	language = "Go"
	searchUser(ts2.URL+"/", pattern, perpage)
}

func Test_showUserResult(t *testing.T) {
	items := ItemUser{
		Login:   "user1",
		HTMLURL: "http://notfound.com",
		Score:   30,
	}
	requests := RespUser{NextURL: "http://notfound.com/next",
		Count:             1,
		IncompleteResults: false,
		Items:             []ItemUser{items},
	}
	ospager = "less"
	requests.ShowUserResult()
	ospager = "more"
	requests.ShowUserResult()
}

func Test_continueSearchUser(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", "<https://www2.com/next>; rel=\"next\", <https://www2.com/prev>; rel=\"prev\"")
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts1.Close()
	resp1 := continueSearchUser(ts1.URL + "/")
	if resp1.Count != 1 {
		t.Errorf("Count mismatch, expected '1', got '%d'", resp1.Count)
	}
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts2.Close()
	resp2 := continueSearchUser(ts2.URL + "/")
	if resp2.Count != 1 {
		t.Errorf("Count mismatch, expected '1', got '%d'", resp2.Count)
	}
}

func Test_getUserConfirm(t *testing.T) {
	ts0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts0.Close()
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//urls := fmt.Sprintf("<%s>; rel=\"next\", <%s>; rel=\"prev\"", ts0.URL+"/", ts0.URL+"/")
		//w.Header().Set("Link", urls)
		w.WriteHeader(http.StatusOK)
		w.Write(userOK)
	}))
	defer ts1.Close()
	items := ItemUser{
		Login:   "user1",
		HTMLURL: "http://notfound.com",
		Score:   30,
	}
	requests := RespUser{NextURL: ts1.URL + "/",
		PreviousURL:       "",
		Count:             1,
		IncompleteResults: false,
		Items:             []ItemUser{items},
	}

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
		if err := getUserConfirm(&requests); err != nil {
			t.Errorf("getUserConfirm failed: %v", err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	}
	// test previous url not empty
	requests = RespUser{PreviousURL: ts1.URL + "/",
		Count:             1,
		IncompleteResults: false,
		Items:             []ItemUser{items},
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
		if err := getUserConfirm(&requests); err != nil {
			t.Errorf("getUserConfirm failed: %v", err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
