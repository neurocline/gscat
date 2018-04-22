// Copyright 2018 Brian Fitzgerald. All rights reserved.

package main

import (
	"fmt"
)

var (
	// CommitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string
)

func main() {
	fmt.Println("Nothing to see here yet")
	fmt.Printf("BuildDate=%s CommitHash=%s\n", BuildDate, CommitHash)
}
