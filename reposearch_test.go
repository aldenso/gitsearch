package main

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_searchRepo(t *testing.T) {
	perpage := strconv.Itoa(paging)
	pattern := "time"
	language, login = "Go", "golang"
	resp1 := searchRepo(pattern, perpage)
	if resp1.Count != 1 {
		fmt.Printf("Using pattern '%s', language '%s', and login '%s'", pattern, language, login)
		t.Errorf("Count mismatch, expected '1', got '%d'", resp1.Count)
	}
	language, login = "Go", ""
	resp2 := searchRepo(pattern, perpage)
	if resp2.Count == 0 {
		fmt.Printf("Using pattern '%s', language '%s'", pattern, language)
		t.Errorf("Count mismatch, expected '!=0', got '%d'", resp2.Count)
	}
	pattern = "zfssa"
	language, login = "", "aldenso"
	resp3 := searchRepo(pattern, perpage)
	if resp3.Count == 0 {
		fmt.Printf("Using pattern '%s', login '%s'", pattern, login)
		t.Errorf("Count mismatch, expected '!=0', got '%d'", resp3.Count)
	}
	pattern = "zfssa"
	language, login = "", ""
	resp4 := searchRepo(pattern, perpage)
	if resp4.Count == 0 {
		fmt.Printf("Using pattern '%s'", pattern)
		t.Errorf("Count mismatch, expected '!=0', got '%d'", resp4.Count)
	}
}

func Test_continueSearchRepo(t *testing.T) {
	url := "https://api.github.com/search/repositories?q=go"
	resp := continueSearchRepo(url)
	if resp.NextURL == "" {
		t.Errorf("continueSearchRepo failed, expected a NextURL but got '%s'", resp.NextURL)
	}
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
	result := requests.showRepoResult()
	for _, v := range result.Items {
		if v.ID != 0 || v.HTMLURL != "http://notfound.com" {
			t.Errorf("showRepoResult wrong, expected ID == '0' HTMLURL == 'http://notfound.com', but got '%d' and '%s'", v.ID, v.HTMLURL)
		}
	}
}
