package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//ShowUserResult function to print results values
func (results *RespUser) ShowUserResult() {
	var output []string
	count := fmt.Sprintf("%s\nResults Count: %v\n%s", line, results.Count, line)
	output = append(output, count)
	for _, i := range results.Items {
		match := fmt.Sprintf("User: %v\nURL: %v\nScore: %v\n%s", i.Login, i.HTMLURL, i.Score, line)
		output = append(output, match)
	}
	pager(strings.Join(output, "\n"))
}

//searchUser function to make search for a particular user pattern
func searchUser(pattern, paging string) RespUser {
	var data RespUser
	var linkHeader string
	var url string
	switch {
	case len(language) > 0:
		url = (apiURL + "users?q=" + pattern + "+language:" + language + "&per_page=" + paging)
	case len(language) == 0:
		url = (apiURL + "users?q=" + pattern + "&per_page=" + paging)
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

// function to continue with next url
func continueSearchUser(url string) RespUser {
	var data RespUser
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

//RunSearchUser function to run the main process for user search
func RunSearchUser(user, paging string) {
	items := searchUser(user, paging)
	if items.Count > 0 {
		items.ShowUserResult()
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
				items = continueSearchUser(items.NextURL)
				items.ShowUserResult()
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
