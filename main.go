package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	apiURL                        = "https://api.github.com/search/"
	user, repo                    bool
	paging                        int
	searchString, language, login string
	line                          = "==============================================================================="
	linesmall                     = "-------------------------------------------------------------------------------"
	linebig                       = "###############################################################################"
	ospager                       string
	less, lessEXE                 = "/usr/bin/less", "less.exe"
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
	fmt.Println("You must use an option like:")
	fmt.Println("./gitsearch -help")
	fmt.Println("./gitsearch -h")
	fmt.Println("./gitsearch -user -pattern pattern")
	fmt.Println("./gitsearch -repo -pattern pattern")
	fmt.Println("./gitsearch -u -p pattern")
	fmt.Println("./gitsearch -r -login username")
	fmt.Println("./gitsearch -r -p pattern")
	fmt.Println("./gitsearch -r -p pattern -l language -login username")
	fmt.Println("./gitsearch -r -p pattern -paging=10")
}

//Regexp function go get the url from Link in header
func Regexp(input string) string {
	var url string
	re0 := regexp.MustCompile("next")
	next := re0.FindString(input)
	if next != "" {
		re1 := regexp.MustCompile("<https://[^/]*/[^>]*>; rel=\"next\"")
		url = re1.FindString(input)
		url = strings.Replace(strings.Replace(strings.Split(url, ";")[0], "<", "", -1), ">", "", -1)
		return url
	}
	return url
}

func pager(input string) {
	if ospager == less || ospager == lessEXE || ospager == "less" {
		cmd := exec.Command(ospager, "-X", "-F", "-R")
		cmd.Stdin = strings.NewReader(input)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cmd := exec.Command(ospager)
		cmd.Stdin = strings.NewReader(input)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	ospager = os.Getenv("PAGER")
	if ospager == "" {
		if _, err := os.Stat(less); os.IsNotExist(err) {
			if os.Getenv("github_shell") == "true" {
				ospager = lessEXE
			} else {
				ospager = "more"
			}

		} else {
			ospager = less
		}
	}
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
