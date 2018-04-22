// +build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// mage install - stamp version number and build
// mage vendor - sync vendored dependencies

// --------------------------------------------------------------------------------------

// Install gscat binary
func Install() error {
	s := flagEnv()
	fmt.Println(s)
	return sh.RunWith(flagEnv(), "go", "install", "-ldflags", ldflags, packageName)
}

// Remember that -X takes the import path. For variables in the main package,
// this is just "main", but for variables in a package, this would be the full
// import path to the package

const packageName  = "github.com/neurocline/gscat"
const noGitLdflags = "-X main.BuildDate=$BUILD_DATE"

var ldflags = "-X main.CommitHash=$COMMIT_HASH -X main.BuildDate=$BUILD_DATE"

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     packageName,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

// --------------------------------------------------------------------------------------

// Install Go Dep and sync vendored dependencies
func Vendor() error {
	mg.Deps(getDep)
	return sh.Run("dep", "ensure")
}

func getDep() error {
	return sh.Run("go", "get", "-u", "github.com/golang/dep/cmd/dep")
}
