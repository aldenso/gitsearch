package main

import (
	"testing"
)

func Test_checkUsage(t *testing.T) {
	checkUsage()
	// Output: You must use an option like:
	// gitsearch -help
	// gitsearch -h
	// gitsearch -user -pattern pattern
	// gitsearch -repo -pattern pattern
	// gitsearch -u -p pattern
	// gitsearch -r -login username
	// gitsearch -r -p pattern
	// gitsearch -r -p pattern -l language -login username
	// gitsearch -r -p pattern -paging=10
	// gitsearch -r -p docker -fork only
}
