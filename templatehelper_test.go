package templatehelper_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/northbright/templatehelper"
)

// Manual represents the manual of templatehelper.
type Manual struct {
	Title        string
	Author       string
	About        string
	Installation string
	ExampleCode  string
}

var (
	// Create a manual for templatehelper.
	manual = Manual{
		Title:        "templatehelper Manual",
		Author:       "Frank Xu",
		About:        "A Go Library provides helper functions for template package.",
		Installation: `go get -u github.com/northbright/templatehelper`,
		ExampleCode:  exampleCode,
	}
)

func ExampleParseDirWithDelims() {
	// Create and parse templates from .tex file.
	// LaTex uses "{}" for command parameter syntax,
	// which conflicts with default golang template delimiters("{{" and "}}")
	// To use Golang template package to parse .tex file,
	// we use new delimiters: "\{\{" and "\}\}" in the .tex template file.

	dir := "templates/latex"
	tmpls, err := templatehelper.ParseDirWithDelims(dir, ".tex", "\\{\\{", "\\}\\}")
	if err != nil {
		log.Printf("ParseDirWithDelims() error: %v", err)
		return
	}

	if len(tmpls) == 0 {
		log.Printf("No template parsed")
		return
	}

	// List the parsed temlates.
	fmt.Printf("Parsed templates:\n")
	for _, tmpl := range tmpls {
		fmt.Printf("%v\n", strings.ReplaceAll(tmpl.Name(), string(os.PathSeparator), ">"))
	}

	// Execute the templates.
	for _, tmpl := range tmpls {
		log.Printf("execute template: %v\n", tmpl.Name())
		tmpl.Execute(os.Stderr, manual)
	}

	// Output:
	//Parsed templates:
	//templates>latex>chapters>00-about.tex
	//templates>latex>chapters>01-installation.tex
	//templates>latex>chapters>02-usage.tex
	//templates>latex>manual.tex
	//templates>latex>title.tex
}

func ExampleParseDir() {
	dir := "templates/markdown"
	tmpls, err := templatehelper.ParseDir(dir, ".md")
	if err != nil {
		log.Printf("ParseDir() error: %v", err)
	}

	if len(tmpls) == 0 {
		log.Printf("No template parsed")
		return
	}

	// List the parsed temlates.
	fmt.Printf("Parsed templates:\n")
	for _, tmpl := range tmpls {
		fmt.Printf("%v\n", strings.ReplaceAll(tmpl.Name(), string(os.PathSeparator), ">"))
	}

	// Execute the templates.
	for _, tmpl := range tmpls {
		log.Printf("execute template: %v\n", tmpl.Name())
		tmpl.Execute(os.Stderr, manual)
	}

	// Output:
	//Parsed templates:
	//templates>markdown>chapters>00-about.md
	//templates>markdown>chapters>01-installation.md
	//templates>markdown>chapters>02-usage.md
	//templates>markdown>title.md
}

var (
	exampleCode = `
package main

import (
        "fmt"

        "github.com/northbright/templatehelper"
)

func main() {
        dir := "templates/markdown"
        tmpls, err := templatehelper.ParseDir(dir, ".md")
        if err != nil {
                fmt.Printf("ParseDir() error: %%v\n", err)
        }

        // List the parsed temlates.
        fmt.Printf("Parsed templates:\n")
        for _, tmpl := range tmpls {
                fmt.Printf("%%v\n", strings.ReplaceAll(tmpl.Name(), string(os.PathSeparator), ">"))
        }

        // Output:
        //Parsed templates:
        //templates>markdown>chapters>00-about.md
        //templates>markdown>chapters>01-installation.md
        //templates>markdown>chapters>02-usage.md
        //templates>markdown>title.md
}
`
)
