package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
)

func main() {
	var site Site
	site.TemplateRegexpPattern = `\.template$`
	site.PageRegexpPattern = `\.(?:md|markdown|html)$`

	flag.BoolVar(&site.Verbose, "verbose", true, "Verbose output")
	flag.StringVar(&site.SourceDirPath, "source", ".", "Source dir")
	flag.StringVar(&site.TargetDirPath, "target", "rigid-result", "Target dir")
	flag.Parse()

	switch command := flag.Arg(0); command {
	case "init":
		// @todo Create default config & blueprint
		fmt.Println("Not implemented yet")
	case "build":
		if err := build(&site); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("Usage: rigid build")
		os.Exit(0)
	}
}
