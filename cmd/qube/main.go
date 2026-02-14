package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/qube"
)

var version string

func init() {
	log.SetFlags(0)
}

func parseArgs() *qube.Options {
	var CLI struct {
		qube.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&CLI, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &CLI.Options
}

func main() {
	options := parseArgs()
	task := qube.NewTask(options)
	report, err := task.Run()

	if err != nil {
		log.Fatal(err)
	}

	report.Version = version
	fmt.Println(report.JSON())
}
