package main

import (
	"fmt"
	"log"

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

	kong.Parse(
		&CLI,
		kong.Vars{"version": version},
	)

	return &CLI.Options
}

func main() {
	options := parseArgs()
	task := qube.NewTask(options)
	report, err := task.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(report.JSON())
}
