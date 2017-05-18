package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
)

var (
	apiURL                        = "https://api.github.com/search/"
	user, repo                    bool
	paging                        int
	searchString, language, login string
)

func init() {
	flag.BoolVar(&user, "user", false, "search for a user.")
	flag.BoolVar(&user, "u", false, "shorthand for -user.")
	flag.BoolVar(&repo, "repo", false, "indicate the pattern you are looking for (don't combine with -user|-u).")
	flag.BoolVar(&repo, "r", false, "shorthand for -repo.")
	flag.StringVar(&searchString, "pattern", "", "indicate the pattern you are looking for.")
	flag.StringVar(&searchString, "p", "", "shorthand for -pattern.")
	flag.StringVar(&language, "lang", "", "indicate a language for search.")
	flag.StringVar(&language, "l", "", "shorthand for -lang.")
	flag.StringVar(&login, "login", "", "indicate username for a repo search.")
	flag.IntVar(&paging, "paging", 100, "set per page limit.")
}

func checkUsage() {
	fmt.Println("You must use an option:")
	fmt.Println("./gitsearch -help")
	fmt.Println("./gitsearch -h")
	fmt.Println("./gitsearch -user -pattern pattern")
	fmt.Println("./gitsearch -repo -pattern pattern")
	fmt.Println("./gitsearch -u -p pattern")
	fmt.Println("./gitsearch -r -p pattern")
	fmt.Println("./gitsearch -r -p pattern -l language -login username")
	fmt.Println("./gitsearch -r -p pattern -paging=10")
}

func lines() {
	fmt.Println("===============================================================================")
}

//Regexp function go get the url from Link in header
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
	perpage := strconv.Itoa(paging)
	if user && searchString != "" {
		RunSearchUser(searchString, perpage)
	} else if repo {
		switch {
		case searchString == "" && login == "":
			checkUsage()
		case searchString == "" && login != "":
			RunSearchRepo(searchString, perpage)
		case searchString != "":
			RunSearchRepo(searchString, perpage)
		}
	} else {
		checkUsage()
	}
}
