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
	root := BuildCommand()
	root.Execute()
}

func BuildCommand() *gscatCmd {

	root := buildRootCommand()

	root.cmd.AddCommand(buildInfoCommand(root).cmd)
	root.cmd.AddCommand(buildReadCommand(root).cmd)

	return root
}

func (c *gscatCmd) Execute() {
	c.cmd.ExecuteC()
}

type gscatCmd struct {
	cmd *cobra.Command
	rootFlags
}

type rootFlags struct {
	verbose bool
	debug bool
}

func buildRootCommand() *gscatCmd {
	c := &gscatCmd{}

	c.cmd = &cobra.Command{
		Use: "gscat",
		Short: "catalog, archive, and search across multiple systems",
		Long: `distributed catalog and search.

Scans, analyzes and archives data across multiple systems, including search
by multiple criteria.`,
		RunE: nil,
	}

	c.cmd.PersistentFlags().BoolVarP(&c.rootFlags.verbose, "verbose", "v", false, "verbose output")
	c.cmd.PersistentFlags().BoolVarP(&c.rootFlags.debug, "debug", "", false, "debug output")

	return c
}

type gscatInfoCmd struct {
	cmd *cobra.Command
	*rootFlags
}

func buildInfoCommand(parent *gscatCmd) *gscatInfoCmd {
	c := &gscatInfoCmd{}
	c.rootFlags = &parent.rootFlags

	c.cmd = &cobra.Command{
		Use: "info",
		Short: "show information about local catalogs",
		RunE: c.info,
	}

	return c
}

func (c *gscatInfoCmd) info(cmd *cobra.Command, args []string) error {
	fmt.Printf("BuildDate=%s CommitHash=%s\n", BuildDate, CommitHash)
	fmt.Printf("runtime.GOMAXPROCS(0)=%d\n", runtime.GOMAXPROCS(0))
	info, err := console.GetConsoleScreenBufferInfo(0)
	if err != nil {
		panic("omg")
	}
	fmt.Printf("console info = %s\n", info)

	return nil
}

type gscatReadCmd struct {
	cmd *cobra.Command
	*rootFlags

	stats bool
}

func buildReadCommand(parent *gscatCmd) *gscatReadCmd {
	c := &gscatReadCmd{}
	c.rootFlags = &parent.rootFlags

	c.cmd = &cobra.Command{
		Use: "read",
		Short: "scan disk and calculate content-ids",
		RunE: c.read,
	}

	c.cmd.Flags().BoolVarP(&c.stats, "stats", "", false, "just get cheap metadata")

	return c
}

func (c *gscatReadCmd) read(cmd *cobra.Command, args []string) error {
	basepath := "./"
	if len(args) > 0 {
		basepath = filepath.Clean(args[0])
	}

	if (c.stats) {
		test.StatsTest(basepath)
	} else {
		test.ReadTest(basepath)
	}

	return nil
}
