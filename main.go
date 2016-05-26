package main

import (
	"flag"
	"fmt"
	"regexp"
)

var (
	apiUrl       = "https://api.github.com/search/"
	user         = flag.Bool("user", false, "search for a user, (do not use with -repo flag)")
	repo         = flag.Bool("repo", false, "search for a repo, (do not use with -user flag)")
	searchString = flag.String("pattern", "", "indicate the pattern you are looking for")
	language     = flag.String("lang", "", "indicate a language for search")
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
	if *user && *searchString != "" {
		RunSearchUser(*searchString)
	} else if *repo && *searchString != "" {
		RunSearchRepo(*searchString)
	} else {
		fmt.Println("You must use an option:")
		fmt.Println("./gitsearch -user -pattern pattern")
		fmt.Println("./gitsearch -repo -pattern pattern")
	}
}
