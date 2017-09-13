package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//showRepoResult function to print results values
func (results *RespRepo) showRepoResult() {
	var output []string
	for _, i := range results.Items {
		match := fmt.Sprintf("Repo: %v\t\tOwner: %v\nDescription: %v\nURL: %v\nLanguage: %v\t\tStars: %v\n%s",
			i.Name, i.Owner.Login, i.Description, i.HTMLURL, i.Language, i.StargazersCount, line)
		output = append(output, match)
	}
	pager(strings.Join(output, "\n"))
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
	fmt.Println(line)
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
	fmt.Println(line)
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
	fmt.Println(line)
	fmt.Println("Results Count:", items.Count)
	if items.Count > 0 {
		items.showRepoResult()
		// loop over next page url
		for items.NextURL != "" {
			fmt.Println(line)
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
				items.showRepoResult()
				/*for _, item := range items.Items {
					lines()
					item.showRepoResult()
				}*/
			case answer == "N" || answer == "n":
				fmt.Println("Stopping")
				os.Exit(0)
			default:
				fmt.Println("*** You must indicate \"Y\" or \"N\" ***")
			}
		}
		fmt.Println(line)
		//fmt.Println("No more results")
		os.Exit(0)
	}
}
