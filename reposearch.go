package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// function to print results values
func (i *ItemRepo) ShowRepoResult() {
	fmt.Println("Repo:", i.Name, "\t\tOwner:", i.Owner.Login)
	fmt.Println("Description:", i.Description)
	fmt.Println("URL:", i.Html_url)
	fmt.Println("Language:", i.Language, "\t\tStars:", i.Stargazers_count)
}

// function to make search for a particular user pattern
func searchRepo(pattern string) (RespRepo, string) {
	var data RespRepo
	var link string
	var url string
	switch {
	case len(language) > 0 && len(login) > 0:
		url = (apiUrl + "repositories?q=" + pattern + "+user:" + login + "+language:" + language)
	case len(language) > 0 && len(login) == 0:
		url = (apiUrl + "repositories?q=" + pattern + "+language:" + language)
	case len(language) == 0 && len(login) > 0:
		url = (apiUrl + "repositories?q=" + pattern + "+user:" + login)
	default:
		url = (apiUrl + "repositories?q=" + pattern)
	}
	lines()
	fmt.Println("using url:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if link = response.Header.Get("Link"); len(link) == 0 {
		fmt.Println("Showing results in one page")
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data, link
}

// function to get next item and maybe next page url
func NextUrlRepo(url string) (RespRepo, string) {
	var data RespRepo
	var link string
	lines()
	fmt.Println("using url:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if link = response.Header.Get("Link"); len(link) == 0 {
		lines()
		fmt.Println("No more results")
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data, link
}

// function to run the main process for user search
func RunSearchRepo(repo string) {
	items, Link := searchRepo(repo)
	lines()
	fmt.Println("Results Count:", items.Count)
	if items.Count > 0 {
		for _, item := range items.Items {
			lines()
			item.ShowRepoResult()
		}
		// loop over next page url
		for len(Link) > 0 {
			nexturl := Regexp(Link)
			lines()
			fmt.Println("Next Page ==>", nexturl)
			var answer string
			fmt.Println("Go to next page? (Y/N):")
			n, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				fmt.Println(n, err)
			}
			switch {
			case answer == "Y" || answer == "y":
				items, Link = NextUrlRepo(nexturl)
				for _, item := range items.Items {
					lines()
					item.ShowRepoResult()
				}
			case answer == "N" || answer == "n":
				fmt.Println("Stopping")
				os.Exit(0)
			default:
				fmt.Println("*** You must indicate \"Y\" or \"N\" ***")
			}
		}
		lines()
		fmt.Println("No more results")
		os.Exit(0)
	}
}
