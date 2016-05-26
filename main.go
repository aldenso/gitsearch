package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

// function to make search for a particular user pattern
func searchUser(pattern string) (Resp, string) {
	var data Resp
	var link string
	url := (apiUrl + "users?q=" + pattern)
	lines()
	fmt.Println("using url:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if link = response.Header.Get("Link"); len(link) > 0 {
		fmt.Println("No more results")
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data, link
}

// function to get next item and maybe next page url
func NextUrl(url string) (Resp, string) {
	var data Resp
	var link string
	lines()
	fmt.Println("using url:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if link = response.Header.Get("Link"); len(link) > 0 {
		lines()
		fmt.Println("No more results")
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data, link
}

// function go get the url from Link in header
func Regexp(input string) string {
	re := regexp.MustCompile("rel=next")
	next := re.FindString(input)
	if len(next) > 0 {
		re = regexp.MustCompile("[[:alnum:]]+[[:graph:]]+[[:alnum:]]")
		url := re.FindString(input)
		return url
	} else {
		return ""
	}
}

func RunSearchUser(user string) {
	items, Link := searchUser(user)
	lines()
	fmt.Println("Results Count:", items.Count)
	if items.Count > 0 {
		for _, item := range items.Items {
			lines()
			fmt.Println("User: ", item.Login)
			fmt.Println("URL", item.Html_url)
		}
		// loop over next page url
		for len(Link) > 0 {
			nexturl := Regexp(Link)
			fmt.Println(nexturl)
			var answer string
			fmt.Println("Go to next page? (y/N):")
			n, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				fmt.Println(n, err)
			}
			switch answer {
			case "Y":
				items, Link = NextUrl(nexturl)
				for _, item := range items.Items {
					lines()
					fmt.Println("User: ", item.Login)
					fmt.Println("URL", item.Html_url)
				}
			case "N":
				fmt.Println("Stopping")
				os.Exit(0)
			}
		}
		lines()
		fmt.Println("No more results")
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	if *user {
		RunSearchUser(*searchString)
	}
}
