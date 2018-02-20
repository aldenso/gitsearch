package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

//showRepoResult function to print results values
func (results *RespRepo) showRepoResult() ItemChoices {
	var output []string
	var items ItemChoices
	var item ItemChoice
	var match string
	count := fmt.Sprintf("%s\nResults Count: %v\n%s", line, results.Count, line)
	output = append(output, count)
	counter := 0
	for _, i := range results.Items {
		item.ID, item.HTMLURL = counter, i.HTMLURL
		items.Items = append(items.Items, item)
		if ospager != "more" {
			match = fmt.Sprintf("Repo: %v\t\tOwner: %v\nDescription: %v\nURL: %v\nLanguage: %v\t\tStars: %v\n%s\nto Git Clone Choose: %s\n%s",
				color.YellowString(i.Name), i.Owner.Login, i.Description, i.HTMLURL, i.Language, i.StargazersCount, linesmall, color.YellowString(strconv.Itoa(counter)), line)
		} else {
			match = fmt.Sprintf("Repo: %v\t\tOwner: %v\nDescription: %v\nURL: %v\nLanguage: %v\t\tStars: %v\n%s\nto Git Clone Choose: %d\n%s",
				i.Name, i.Owner.Login, i.Description, i.HTMLURL, i.Language, i.StargazersCount, linesmall, counter, line)

		}
		output = append(output, match)
		counter++
	}
	pager(strings.Join(output, "\n"))
	return items
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
	if items.Count > 0 {
		choices := items.showRepoResult()
		// loop over next page url
		for items.NextURL != "" {
			fmt.Println(line)
			fmt.Println("Next Page ==>", items.NextURL)
			var answer string
			fmt.Println("Git Clone repo # or Go to next page? (Y/N):")
			n, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				fmt.Println(n, err)
			}
			switch {
			case answer == "Y" || answer == "y":
				items = continueSearchRepo(items.NextURL)
				items.showRepoResult()
			case answer == "N" || answer == "n":
				fmt.Println("Stopping")
				os.Exit(0)
			case answer == "S" || answer == "s":
				items.showRepoResult()
			default:
				var clone string
				for _, v := range choices.Items {
					if answer == strconv.Itoa(v.ID) {
						clone = v.HTMLURL
					}
				}
				if clone == "" {
					fmt.Printf("%s\n--- Option '%s' no available ---\n%s\n", linebig, answer, linebig)
				} else {
					fmt.Printf("%s\n+++ Cloning repo %s +++\n%s\n", linebig, clone, linebig)
					cmd := exec.Command("git", "clone", clone)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Run()
					fmt.Println(line)
				}
			}
		}
		if items.NextURL == "" {

			fmt.Println(line)
			var answer string
			fmt.Println("Git Clone repo #, show or quit? (#/S/Q):")
			n, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				fmt.Println(n, err)
			}
			switch {
			case answer == "Q" || answer == "q":
				fmt.Println("Stopping")
				os.Exit(0)
			case answer == "S" || answer == "s":
				items.showRepoResult()
			default:
				var clone string
				for _, v := range choices.Items {
					if answer == strconv.Itoa(v.ID) {
						clone = v.HTMLURL
					}
				}
				if clone == "" {
					fmt.Printf("%s\n--- Option '%s' no available ---\n%s\n", linebig, answer, linebig)
				} else {
					fmt.Printf("%s\n+++ Cloning repo %s +++\n%s\n", linebig, clone, linebig)
					cmd := exec.Command("git", "clone", clone)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Run()
					fmt.Println(line)
				}
			}
		}
		fmt.Println(line)
		os.Exit(0)
	}
}
