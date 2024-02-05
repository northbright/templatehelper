package templatehelper_test

import (
	"log"
	"os"

	"github.com/northbright/templatehelper"
)

func ExampleParseFilesWithDelims() {
	// Manual represents the manual of tex-go.
	type Manual struct {
		Title        string
		Author       string
		About        string
		Installation string
		ExampleCode  string
	}

	// Create a manual for templatehelper.
	manual := Manual{
		Title:        "templatehelper Manual",
		Author:       "Frank Xu",
		About:        "A Go Library provides helper functions for template package.",
		Installation: `go get -u github.com/northbright/templatehelper`,
		ExampleCode: `
package main

import (
    "github.com/northbright/templatehelper"
)
`,
	}

	// Create and parse templates from .tex file.
	// LaTex uses "{}" for command parameter syntax,
	// which conflicts with default golang template delimiters("{{" and "}}")
	// To use Golang template package to parse .tex file,
	// we use new delimiters: "\{\{" and "\}\}" in the .tex template file.
	filenames := []string{
		"templates/manual.tex",
		"templates/title.tex",
		"templates/chapters/00-about.tex",
		"templates/chapters/01-installation.tex",
		"templates/chapters/02-usage.tex",
	}

	t, err := templatehelper.ParseFilesWithDelims("\\{\\{", "\\}\\}", filenames...)
	if err != nil {
		log.Printf("ParseFilesWithDelims() error: %v", err)
		return
	}

	tmpls := t.Templates()
	for _, tmpl := range tmpls {
		tmpl.Execute(os.Stdout, manual)
	}

	if err := t.Execute(os.Stdout, manual); err != nil {
		log.Printf("t.Execute() error: %v", err)
		return
	}

	// Output:
}
