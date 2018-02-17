# gitsearch

Small golang program to help you search for users or repos in github, useful when you are using only a shell and don't want to deal with long inputs and output from curl.

Works with GitHub API v3.

```sh
go get github.com/aldenso/gitsearch
```

To show you the usage just type the command alone (shows examples) or use the option help.

```sh
gitsearch
gitsearch -h
```

```txt
You must use an option like:
./gitsearch -help
./gitsearch -h
./gitsearch -user -pattern pattern
./gitsearch -repo -pattern pattern
./gitsearch -u -p pattern
./gitsearch -r -login username
./gitsearch -r -p pattern
./gitsearch -r -p pattern -l language -login username
./gitsearch -r -p pattern -paging=10


```

```txt
Usage of gitsearch:
  -l string
        shorthand for -lang.
  -lang string
        indicate a language for search.
  -login string
        indicate username for a repo search.
  -p string
        shorthand for -pattern.
  -paging int
        set per page limit. (default 100)
  -pattern string
        indicate the pattern you are looking for.
  -r    shorthand for -repo.
  -repo
        indicate the pattern you are looking for (don't combine with -user|-u).
  -u    shorthand for -user.
  -user
        search for a user.
```

Search repos for username "aldenso".

```sh
gitsearch -r -login aldenso
```

```txt
===============================================================================
Results Count: 23
===============================================================================
Repo: tools             Owner: aldenso
Description: This Repo is to backup and share some tools (scripts in bash or python)
URL: https://github.com/aldenso/tools
Language: Python                Stars: 3
-------------------------------------------------------------------------------
to Git Clone Choose: 0
===============================================================================
Repo: gitsearch         Owner: aldenso
Description: Small golang program to help you search for users or repos in github
URL: https://github.com/aldenso/gitsearch
Language: Go            Stars: 3
-------------------------------------------------------------------------------
to Git Clone Choose: 1
===============================================================================
.
.
.
```

Search repos with pattern "gitsearch".

```sh
gitsearch -r -p gitsearch
```

```txt
===============================================================================
Results Count: 151
===============================================================================
Repo: gitsearch         Owner: Tom-Alexander
Description: git repository indexer and client
URL: https://github.com/Tom-Alexander/gitsearch
Language: JavaScript            Stars: 4
-------------------------------------------------------------------------------
to Git Clone Choose: 0
===============================================================================
Repo: gitsearch         Owner: aldenso
Description: Small golang program to help you search for users or repos in github
URL: https://github.com/aldenso/gitsearch
Language: Go            Stars: 3
-------------------------------------------------------------------------------
to Git Clone Choose: 1
===============================================================================
Repo: gitSearch         Owner: Zivi
Description:
URL: https://github.com/Zivi/gitSearch
Language: JavaScript            Stars: 0
-------------------------------------------------------------------------------
.
.
.
```

Search repos with pattern "gitsearch" and paging equals 4 (paging is convenient for smaller show groups and while paging you'll have the possibility to clone a repo indicating the given number in the show).

```sh
gitsearch -r -p gitsearch -paging=4
```

```txt
===============================================================================
Results Count: 151
===============================================================================
Repo: gitsearch         Owner: Tom-Alexander
Description: git repository indexer and client
URL: https://github.com/Tom-Alexander/gitsearch
Language: JavaScript            Stars: 4
-------------------------------------------------------------------------------
to Git Clone Choose: 0
===============================================================================
Repo: gitsearch         Owner: aldenso
Description: Small golang program to help you search for users or repos in github
URL: https://github.com/aldenso/gitsearch
Language: Go            Stars: 3
-------------------------------------------------------------------------------
to Git Clone Choose: 1
===============================================================================
Repo: gitSearch         Owner: Zivi
Description:
URL: https://github.com/Zivi/gitSearch
Language: JavaScript            Stars: 0
-------------------------------------------------------------------------------
.
.
.
```

Search repos with pattern "go", associated with a username "aldenso".

```sh
gitsearch -r -p go -login aldenso
```

```txt
===============================================================================
Results Count: 3
===============================================================================
Repo: gotoolbackup              Owner: aldenso
Description: Program to create backups using toml file, where you indicate origin/destiny directories and retention period in days
URL: https://github.com/aldenso/gotoolbackup
Language: Go            Stars: 0
-------------------------------------------------------------------------------
to Git Clone Choose: 0
===============================================================================
Repo: golang-examples           Owner: aldenso
Description: My Golang examples to practice go programming
URL: https://github.com/aldenso/golang-examples
Language: Go            Stars: 0
-------------------------------------------------------------------------------
to Git Clone Choose: 1
===============================================================================
Repo: golang-mongodb-microservice-example               Owner: aldenso
Description: Small example for a microservice using golang, mongodb and gorilla
URL: https://github.com/aldenso/golang-mongodb-microservice-example
Language: Go            Stars: 0
-------------------------------------------------------------------------------
to Git Clone Choose: 2
===============================================================================
 Showing results in one page
===============================================================================
```

Search repos with pattern "zfs", associated with a username "aldenso" and language python.

```sh
gitsearch -r -p zfs -l python -login aldenso
```

```txt
===============================================================================
Results Count: 4
===============================================================================
Repo: prtgZFSSAhealth           Owner: aldenso
Description: PRTG Python Advanced script to get health status from ZFSSA using Rest service.
URL: https://github.com/aldenso/prtgZFSSAhealth
Language: Python                Stars: 0
-------------------------------------------------------------------------------
to Git Clone Choose: 0
===============================================================================
Repo: prtgZFSSAmetrics          Owner: aldenso
Description: PRTG Python Advanced script to get metrics from analytics datasets using ZFSSA Rest api
URL: https://github.com/aldenso/prtgZFSSAmetrics
Language: Python                Stars: 0
-------------------------------------------------------------------------------
to Git Clone Choose: 1
===============================================================================
Repo: zfssa-scripts             Owner: aldenso
Description: My scripts for ZFS Storage Appliance
URL: https://github.com/aldenso/zfssa-scripts
Language: Python                Stars: 0
-------------------------------------------------------------------------------
to Git Clone Choose: 2
===============================================================================
Repo: zfssa_utils               Owner: aldenso
Description: Command Line utility to handle most common tasks with ZFS Storage Appliance **Work In progress**
URL: https://github.com/aldenso/zfssa_utils
Language: Python                Stars: 1
-------------------------------------------------------------------------------
to Git Clone Choose: 3
===============================================================================
```

Search users with pattern "alde".

```sh
gitsearch -u -p alde
```

```txt
===============================================================================
Results Count: 1290
===============================================================================
User: alde
URL: https://github.com/alde
Score: 74.62327
===============================================================================
User: iObsidian
URL: https://github.com/iObsidian
Score: 46.332863
===============================================================================
User: aldeed
URL: https://github.com/aldeed
Score: 33.927265
===============================================================================
User: aldeka
URL: https://github.com/aldeka
Score: 25.134708
===============================================================================
User: aldesantis
URL: https://github.com/aldesantis
Score: 23.03635
===============================================================================
.
.
.
```

Search users with pattern "alden" associated with language "go" ("The API doesn't return an accurate result!").

```sh
gitsearch -u -p alden -l go
```

```txt
===============================================================================
Results Count: 3
===============================================================================
User: aotimme
URL: https://github.com/aotimme
Score: 26.87695
===============================================================================
User: aldenso
URL: https://github.com/aldenso
Score: 17.838873
===============================================================================
User: CampinCarl
URL: https://github.com/CampinCarl
Score: 5.132159
===============================================================================
```
