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
			match = fmt.Sprintf("Repo: %v\t\tOwner: %v\nDescription: %v\nURL: %v\nLanguage: %v\t\tStars: %v\tfork: %v\n%s\nto Git Clone Choose: %s\n%s",
				color.YellowString(i.Name), i.Owner.Login, i.Description, i.HTMLURL, i.Language, i.StargazersCount, i.Fork, linesmall, color.YellowString(strconv.Itoa(counter)), line)
		} else {
			match = fmt.Sprintf("Repo: %v\t\tOwner: %v\nDescription: %v\nURL: %v\nLanguage: %v\t\tStars: %v\tfork: %v\n%s\nto Git Clone Choose: %d\n%s",
				i.Name, i.Owner.Login, i.Description, i.HTMLURL, i.Language, i.StargazersCount, i.Fork, linesmall, counter, line)

		}
		output = append(output, match)
		counter++
	}
	pager(strings.Join(output, "\n"))
	return items
}

//searchRepo function to make search for a particular user pattern
func searchRepo(apiurl, pattern, paging string) RespRepo {
	var data RespRepo
	var linkHeader string
	var url string
	switch {
	case len(language) > 0 && len(login) > 0:
		url = (apiurl + "repositories?q=" + pattern + "+user:" + login + "+fork:" + fork + "+language:" + language + "&per_page=" + paging)
	case len(language) > 0 && len(login) == 0:
		url = (apiurl + "repositories?q=" + pattern + "+fork:" + fork + "+language:" + language + "&per_page=" + paging)
	case len(language) == 0 && len(login) > 0:
		url = (apiurl + "repositories?q=" + pattern + "+user:" + login + "+fork:" + fork + "&per_page=" + paging)
	default:
		url = (apiurl + "repositories?q=" + pattern + "+fork:" + fork + "&per_page=" + paging)
	}
	fmt.Println(line)
	fmt.Printf("using url:\n%s\n", url)
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

//continueSearchRepo function to make search for a particular user pattern
func continueSearchRepo(url string) RespRepo {
	var data RespRepo
	var linkHeader string
	fmt.Println(line)
	fmt.Printf("using url:\n%s\n", url)
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

func getRepoConfirmNextNotEmpty(items *RespRepo, choices *ItemChoices) error {
	fmt.Println(line)
	fmt.Println("Next Page ==>", items.NextURL)
	fmt.Println("Previous Page ==>", items.PreviousURL)
	var answer string
	fmt.Println("Git Clone repo #, Go to next/previous page, Show again or Quit? (#/N/P/S/Q):")
	n, err := fmt.Scanf("%s\n", &answer)
	if err != nil {
		fmt.Println(n, err)
	}
	switch {
	case answer == "N" || answer == "n":
		*items = continueSearchRepo(items.NextURL)
		*choices = items.showRepoResult()
		return nil
	case answer == "P" || answer == "p":
		if items.PreviousURL != "" {
			*items = continueSearchRepo(items.PreviousURL)
			*choices = items.showRepoResult()
		} else {
			fmt.Println("No previous page")
		}
		return nil
	case answer == "Q" || answer == "q":
		fmt.Println("Stopping")
		os.Exit(0)
	case answer == "S" || answer == "s":
		*choices = items.showRepoResult()
		return nil
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
			if ospager != "more" {
				colormsg := fmt.Sprintf("%s\n%s %s\n%s\n",
					color.CyanString(linebig), color.CyanString("Cloning repo:"), color.MagentaString(clone), color.CyanString(linebig))
				pager(colormsg)
			} else {
				fmt.Printf("%s\nCloning repo %s\n%s\n", linebig, clone, linebig)
			}
			cmd := exec.Command("git", "clone", clone)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			fmt.Println(line)
		}
	}
	return nil
}

func getRepoConfirmNextEmpty(items *RespRepo, choices *ItemChoices) error {
	fmt.Println(line)
	fmt.Println("Next Page ==>", items.NextURL)
	fmt.Println("Previous Page ==>", items.PreviousURL)
	var answer string
	fmt.Println("Git Clone repo #, Go to previous page, Show again or Quit? (#/P/S/Q):")
	n, err := fmt.Scanf("%s\n", &answer)
	if err != nil {
		fmt.Println(n, err)
	}
	switch {
	case answer == "P" || answer == "p":
		if items.PreviousURL != "" {
			*items = continueSearchRepo(items.PreviousURL)
			*choices = items.showRepoResult()
		} else {
			fmt.Println("No previous page")
		}
		return nil
	case answer == "Q" || answer == "q":
		fmt.Println("Stopping")
		os.Exit(0)
	case answer == "S" || answer == "s":
		*choices = items.showRepoResult()
		return nil
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
			if ospager != "more" {
				colormsg := fmt.Sprintf("%s\n%s %s\n%s\n",
					color.CyanString(linebig), color.CyanString("Cloning repo:"), color.MagentaString(clone), color.CyanString(linebig))
				pager(colormsg)
			} else {
				fmt.Printf("%s\nCloning repo %s\n%s\n", linebig, clone, linebig)
			}
			cmd := exec.Command("git", "clone", clone)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			fmt.Println(line)
		}
	}
	return nil
}

//RunSearchRepo function to run the main process for user search
func RunSearchRepo(apiurl, repo, paging string) {
	items := searchRepo(apiurl, repo, paging)
	if items.Count > 0 {
		choices := items.showRepoResult()
		// loop over next page url
		for items.NextURL != "" {
			getRepoConfirmNextNotEmpty(&items, &choices)
		}
		for items.NextURL == "" {
			fmt.Println(line)
			getRepoConfirmNextEmpty(&items, &choices)
		}
		fmt.Println(line)
		os.Exit(0)
	}
}
