// Copyright 2018 Brian Fitzgerald. All rights reserved.

package core

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/neurocline/gscat/pkg/console"
	"github.com/neurocline/gscat/pkg/test"

	"github.com/spf13/cobra"
)

var (
	// CommitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string
)

func Run() {
	cmd := BuildCommand()
	cmd.cmd.ExecuteC()
}

type gscatCmd struct {
	cmd *cobra.Command

	stats bool
	read bool
}

func BuildCommand() *gscatCmd {
	c := &gscatCmd{}

	c.cmd = &cobra.Command{
		Use: "gscat",
		Short: "catalog, archive, and search across multiple systems",
		Long: `distributed catalog and search.

Scans, analyzes and archives data across multiple systems, including search
by multiple criteria.`,
		RunE: c.gscat,
	}

	c.cmd.Flags().BoolVarP(&c.stats, "stats", "", false, "scan for stats")
	c.cmd.Flags().BoolVarP(&c.read, "read", "", false, "read files")

	return c
}

func (c *gscatCmd) gscat(cmd *cobra.Command, args []string) error {
	fmt.Printf("BuildDate=%s CommitHash=%s\n", BuildDate, CommitHash)
	fmt.Printf("runtime.GOMAXPROCS(0)=%d\n", runtime.GOMAXPROCS(0))
	info, err := console.GetConsoleScreenBufferInfo(0)
	if err != nil {
		panic("omg")
	}
	fmt.Printf("console info = %s\n", info)

	basepath := "./"
	if len(args) > 0 {
		basepath = filepath.Clean(args[0])
	}

	if (c.stats) {
		test.StatsTest(basepath)
	}
	if (c.read) {
		test.ReadTest(basepath)
	}

	return nil
}
