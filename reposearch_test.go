package main

import (
	"fmt"
	"testing"
)

func Test_searchRepo(t *testing.T) {
	pattern := "time"
	language, login = "Go", "golang"
	resp := searchRepo(pattern)
	if resp.Count != 1 {
		fmt.Printf("Using pattern '%s', language '%s', and login '%s'", pattern, language, login)
		t.Errorf("Count mismatch, expected '1', got '%d'", resp.Count)
	}
}

func Test_continueSearchRepo(t *testing.T) {
	url := "https://api.github.com/search/repositories?q=go"
	resp := continueSearchRepo(url)
	if resp.NextURL == "" {
		t.Errorf("continueSearchRepo failed, expected a NextURL but got '%s'", resp.NextURL)
	}
}
