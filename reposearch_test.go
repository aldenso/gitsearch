package main

import (
	"fmt"
	"net/http"
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
			fmt.Println()
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
