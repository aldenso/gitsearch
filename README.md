gitsearch
=========

Small golang program to help you search for users or repos in github, useful when you are using only a shell and don't want to deal with long inputs and output from curl.

Works with GitHub API v3.

Usage:

Show Program .

    $ ./gitsearch -h
    Usage of ./gitsearch:
    -l string
      shorthand for -lang
    -lang string
      indicate a language for search
    -login string
      indicate username for a repo search
    -p string
      shorthand for -pattern
    -pattern string
      indicate the pattern you are looking for
    -r	shorthand for -repo
    -repo
      indicate the pattern you are looking for (don't combine with -user|-u)
    -u	shorthand for -user
    -user
      search for a user.


Search for repos for username "aldenso".

    $ ./gitsearch -r -login aldenso
    ===============================================================================
    using url: https://api.github.com/search/repositories?q=+user:aldenso
    Showing results in one page
    ===============================================================================
    Results Count: 12
    ===============================================================================
    .
    .
    .
    ===============================================================================
    Repo: pyconverter 		Owner: aldenso
    Description: Utility to convert units to different units and measurement systems.
    URL: https://github.com/aldenso/pyconverter
    Language: Python 		Stars: 0
    ===============================================================================


Search for repos with pattern "go" associated with a username "aldenso".

    $ ./gitsearch -r -p go -login aldenso
    ===============================================================================
    using url: https://api.github.com/search/repositories?q=go+user:aldenso
    Showing results in one page
    ===============================================================================
    Results Count: 3
    ===============================================================================
    Repo: gotoolbackup 		Owner: aldenso
    Description: Program to create backups using toml file, where you indicate origi
    n/destiny directories and retention period in days
    URL: https://github.com/aldenso/gotoolbackup
    Language: Go 		Stars: 0
    ===============================================================================
    Repo: golang-examples 		Owner: aldenso
    Description: My Golang examples to practice go programming
    URL: https://github.com/aldenso/golang-examples
    Language: Go 		Stars: 0
    ===============================================================================
    Repo: golang-mongodb-microservice-example 		Owner: aldenso
    Description: Small example for a microservice using golang, mongodb and gorilla
    URL: https://github.com/aldenso/golang-mongodb-microservice-example
    Language: Go 		Stars: 0
    ===============================================================================
    No more results


Search for users with pattern "alde".

    $ ./gitsearch -u -p alde
    ===============================================================================
    using url: https://api.github.com/search/users?q=alde
    ===============================================================================
    Results Count: 724
    ===============================================================================
    User:  alde
    URL https://github.com/alde
    Score: 63.394615
    ===============================================================================
    .
    .
    .
    ===============================================================================
    User:  Alderian
    URL https://github.com/Alderian
    Score: 9.959794
    ===============================================================================
    Next Page ==> https://api.github.com/search/users?q=alde&page=2
    Go to next page? (Y/N):


Search for users with pattern "aldens" associated with language "shell" ("The API doesn't return an accurate result!").

    $ ./gitsearch -u -p aldens -l shell
    ===============================================================================
    using url: https://api.github.com/search/users?q=aldens+language:shell
    Showing results in one page
    ===============================================================================
    Results Count: 1
    ===============================================================================
    User:  aldenso
    URL https://github.com/aldenso
    Score: 13.826229
    ===============================================================================
    No more results
