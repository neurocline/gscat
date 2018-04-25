// Copyright 2018 Brian Fitzgerald. All rights reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/neurocline/gscat/pkg/test"
)

var (
	// CommitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string
)

func main() {
	fmt.Printf("BuildDate=%s CommitHash=%s\n", BuildDate, CommitHash)
	fmt.Printf("runtime.GOMAXPROCS(0)=%d\n", runtime.GOMAXPROCS(0))

	basepath := "./"
	if len(os.Args) > 1 {
		basepath = filepath.Clean(os.Args[1])
	}

	// test.StatsTest(basepath)
	test.ReadTest(basepath)
}
