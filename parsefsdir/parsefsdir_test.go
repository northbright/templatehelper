package parsefsdir_test

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/northbright/templatehelper/parsefsdir"
)

var (
	//go:embed parsefsdir_test.go
	exampleCode string

	//go:embed templates
	embededFiles embed.FS
)

// Manual represents the manual of parsefsdir.
type Manual struct {
	Title        string
	Author       string
	About        string
	Installation string
	ExampleCode  string
}

var (
	// Create a manual for parsefsdir.
	manual = Manual{
		Title:        "parsefsdir package Manual",
		Author:       "Frank Xu",
		About:        "parsefsdir is a golang package which parses templates in dir of fs.FS.",
		Installation: `go get -u github.com/northbright/templatehelper/parsefsdir`,
		// Embed "dir_test.go" into a string as the example code.
		ExampleCode: exampleCode,
	}
)

func ExampleNew() {
	// Parse markdown template files in a dir of embed FS.
	dir := "templates/markdown"
	ext := ".md"

	// Create a Parser.
	p := parsefsdir.New(embededFiles, dir, parsefsdir.Ext(ext))

	ctx := context.Background()

	log.Printf("p.Parse() starts...\ndir: %v\next: %v", dir, ext)
	// Parse the templates in the dir.
	tmpls, err := p.Parse(ctx)
	if err != nil {
		log.Printf("p.Parse() error: %v", err)
		return
	}

	log.Printf("p.Parse() successfully")
	if len(tmpls) == 0 {
		log.Printf("no template parsed")
		return
	}

	// List the parsed templates.
	fmt.Printf("parsed templates:\n")
	for _, tmpl := range tmpls {
		fmt.Printf("%v\n", tmpl.Name())
	}

	// Execute the templates.
	for _, tmpl := range tmpls {
		log.Printf("execute template: %v\n", tmpl.Name())
		tmpl.Execute(os.Stderr, manual)
	}

	// Parse LaTex template files in a dir.
	dir = "templates/latex"
	ext = ".tex"

	// Create a Parser with specified delimiters.
	// For templates(.tex) in templates/latex, use \{\{ and \}\} as delimiters
	// because '{' and '}' already used by LaTex.
	p = parsefsdir.New(embededFiles, dir, parsefsdir.Ext(ext), parsefsdir.Delims("\\{\\{", "\\}\\}"))

	ctx = context.Background()

	log.Printf("p.Parse() starts...\ndir: %v\next: %v", dir, ext)
	// Parse the templates in the dir.
	tmpls, err = p.Parse(ctx)
	if err != nil {
		log.Printf("p.Parse() error: %v", err)
		return
	}

	log.Printf("p.Parse() successfully")
	if len(tmpls) == 0 {
		log.Printf("no template parsed")
		return
	}

	// List the parsed templates.
	fmt.Printf("parsed templates:\n")
	for _, tmpl := range tmpls {
		fmt.Printf("%v\n", tmpl.Name())
	}

	// Execute the templates.
	for _, tmpl := range tmpls {
		log.Printf("execute template: %v\n", tmpl.Name())
		tmpl.Execute(os.Stderr, manual)
	}

	// Output:
	//parsed templates:
	//templates/markdown/chapters/00-about.md
	//templates/markdown/chapters/01-installation.md
	//templates/markdown/chapters/02-usage.md
	//templates/markdown/title.md
	//parsed templates:
	//templates/latex/chapters/00-about.tex
	//templates/latex/chapters/01-installation.tex
	//templates/latex/chapters/02-usage.tex
	//templates/latex/manual.tex
	//templates/latex/title.tex
}
