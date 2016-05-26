package main

import (
	"flag"
	"fmt"
	"regexp"
)

var (
	apiUrl       = "https://api.github.com/search/"
	user         = flag.Bool("user", false, "indicate is you are looking for a user")
	repo         = flag.Bool("repo", false, "indicate is you are looking for a repo")
	searchString = flag.String("pattern", "github", "indicate the pattern you are looking for")
)

func lines() {
	fmt.Println("===============================================================")
}

// function go get the url from Link in header
func Regexp(input string) string {
	var url string
	re1 := regexp.MustCompile("next")
	next := re1.FindString(input)
	if next != "" {
		re2 := regexp.MustCompile("[[:alnum:]]+[[:graph:]]+[[:alnum:]]")
		url = re2.FindString(input)
		return url
	}
	return url
}

func main() {
	flag.Parse()
	if *user {
		RunSearchUser(*searchString)
	}
}
