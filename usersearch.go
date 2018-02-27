package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

//ShowUserResult function to print results values
func (results *RespUser) ShowUserResult() {
	var output []string
	count := fmt.Sprintf("%s\nResults Count: %v\n%s", line, results.Count, line)
	output = append(output, count)
	for _, i := range results.Items {
		if ospager != "more" {
			match := fmt.Sprintf("User: %v\nURL: %v\nScore: %v\n%s", color.YellowString(i.Login), i.HTMLURL, i.Score, line)
			output = append(output, match)
		} else {
			match := fmt.Sprintf("User: %v\nURL: %v\nScore: %v\n%s", i.Login, i.HTMLURL, i.Score, line)
			output = append(output, match)
		}
	}
	pager(strings.Join(output, "\n"))
}

//searchUser function to make search for a particular user pattern
func searchUser(apiurl, pattern, paging string) RespUser {
	var data RespUser
	var linkHeader string
	var url string
	switch {
	case len(language) > 0:
		url = (apiurl + "users?q=" + pattern + "+language:" + language + "&per_page=" + paging)
	case len(language) == 0:
		url = (apiurl + "users?q=" + pattern + "&per_page=" + paging)
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
	data.NextURL, data.PreviousURL = Regexp(linkHeader)
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
	data.NextURL, data.PreviousURL = Regexp(linkHeader)
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func getUserConfirm(items *RespUser) error {
	var answer string
	fmt.Println("Go to next/previous page, Show again or Quit? (N/P/S/Q):")
	n, err := fmt.Scanf("%s\n", &answer)
	if err != nil {
		fmt.Println(n, err)
		return err
	}
	switch {
	case answer == "N" || answer == "n":
		if items.NextURL != "" {
			*items = continueSearchUser(items.NextURL)
			items.ShowUserResult()
		} else {
			fmt.Println("No next page found")
		}
		return nil
	case answer == "P" || answer == "p":
		if items.PreviousURL != "" {
			*items = continueSearchUser(items.PreviousURL)
			items.ShowUserResult()
		} else {
			fmt.Println("No previous page found")
		}
		return nil
	case answer == "Q" || answer == "q":
		fmt.Println("Stopping")
		os.Exit(0)
	case answer == "S" || answer == "s":
		items.ShowUserResult()
		return nil
	default:
		fmt.Printf("%s\n--- Option '%s' no available ---\n%s\n", linebig, answer, linebig)
	}
	return nil
}

//RunSearchUser function to run the main process for user search
func RunSearchUser(apiurl, user, paging string) {
	items := searchUser(apiurl, user, paging)
	if items.Count > 0 {
		items.ShowUserResult()
		// loop over next page url
		for items.NextURL != "" || items.PreviousURL != "" {
			fmt.Println(line)
			fmt.Println("Next Page ==>", items.NextURL)
			fmt.Println("Previous Page ==>", items.PreviousURL)
			getUserConfirm(&items)
		}
		fmt.Println(line)
		//os.Exit(0)
	}
}
