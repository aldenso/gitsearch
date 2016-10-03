package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

//showRepoResult function to print results values
func (i *ItemRepo) showRepoResult() {
	fmt.Println("Repo:", i.Name, "\t\tOwner:", i.Owner.Login)
	fmt.Println("Description:", i.Description)
	fmt.Println("URL:", i.HTMLURL)
	fmt.Println("Language:", i.Language, "\t\tStars:", i.StargazersCount)
}

//searchRepo function to make search for a particular user pattern
func searchRepo(pattern, paging string) RespRepo {
	var data RespRepo
	var linkHeader string
	var url string
	switch {
	case len(language) > 0 && len(login) > 0:
		url = (apiURL + "repositories?q=" + pattern + "+user:" + login + "+language:" + language + "&per_page=" + paging)
	case len(language) > 0 && len(login) == 0:
		url = (apiURL + "repositories?q=" + pattern + "+language:" + language + "&per_page=" + paging)
	case len(language) == 0 && len(login) > 0:
		url = (apiURL + "repositories?q=" + pattern + "+user:" + login + "&per_page=" + paging)
	default:
		url = (apiURL + "repositories?q=" + pattern + "&per_page=" + paging)
	}
	lines()
	fmt.Println("using url:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if linkHeader = response.Header.Get("Link"); len(linkHeader) == 0 {
		fmt.Println("Showing results in one page")
	}
	data.NextURL = Regexp(linkHeader)
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

//continueSearchRepo function to make search for a particular user pattern
func continueSearchRepo(url string) RespRepo {
	var data RespRepo
	var linkHeader string
	lines()
	fmt.Println("using url:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if linkHeader = response.Header.Get("Link"); len(linkHeader) == 0 {
		fmt.Println("No more results")
	}
	data.NextURL = Regexp(linkHeader)
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

//RunSearchRepo function to run the main process for user search
func RunSearchRepo(repo, paging string) {
	items := searchRepo(repo, paging)
	lines()
	fmt.Println("Results Count:", items.Count)
	if items.Count > 0 {
		for _, item := range items.Items {
			lines()
			item.showRepoResult()
		}
		// loop over next page url
		for items.NextURL != "" {
			lines()
			fmt.Println("Next Page ==>", items.NextURL)
			var answer string
			fmt.Println("Go to next page? (Y/N):")
			n, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				fmt.Println(n, err)
			}
			switch {
			case answer == "Y" || answer == "y":
				items = continueSearchRepo(items.NextURL)
				for _, item := range items.Items {
					lines()
					item.showRepoResult()
				}
			case answer == "N" || answer == "n":
				fmt.Println("Stopping")
				os.Exit(0)
			default:
				fmt.Println("*** You must indicate \"Y\" or \"N\" ***")
			}
		}
		lines()
		//fmt.Println("No more results")
		os.Exit(0)
	}
}
