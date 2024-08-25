package templatedir_test

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/northbright/templatehelper/templatedir"
)

//go:embed templatedir_test.go
var exampleCode string

// Manual represents the manual of templatedir.
type Manual struct {
	Title        string
	Author       string
	About        string
	Installation string
	ExampleCode  string
}

var (
	// Create a manual for templatedir.
	manual = Manual{
		Title:        "templatedir package Manual",
		Author:       "Frank Xu",
		About:        "templatedir is a golang package which parses templates in a dir.",
		Installation: `go get -u github.com/northbright/templatehelper/templatedir`,
		// Embed "templatedir_test.go" into a string as the example code.
		ExampleCode: exampleCode,
	}
)

func ExampleDir_Parse() {
	// Part I. Parse markdown templates.
	dir := "templates/markdown"

	// Create a template dir.
	d := templatedir.New(dir)

	log.Printf("Parse template files in a dir starts...\ndir: %v\n", dir)
	// Parse the templates in the dir.
	tmpls, err := d.Parse()
	if err != nil {
		log.Printf("d.Parse() error: %v", err)
		return
	}

	log.Printf("Parse template files in a dir successfully")
	if len(tmpls) == 0 {
		log.Printf("no template parsed")
		return
	}

	// List the parsed templates.
	fmt.Printf("parsed templates:\n")
	for _, tmpl := range tmpls {
		fmt.Printf("%v\n", strings.ReplaceAll(tmpl.Name(), string(os.PathSeparator), ">"))
	}

	// Execute the templates.
	for _, tmpl := range tmpls {
		log.Printf("execute template: %v\n", tmpl.Name())
		tmpl.Execute(os.Stderr, manual)
	}

	// Part II. Parse LaTex templates.
	dir = "templates/latex"

	// Create a template dir with specified delimiters.
	// For template files in templates/latex,
	// use \{\{ and \}\} as delimiters because '{' and '}' already used by LaTex.
	// It makes the template files can be opened / edited by LaTex editors successfully.
	d = templatedir.New(dir, templatedir.Delims("\\{\\{", "\\}\\}"))

	log.Printf("Parse template files in a dir starts...\ndir: %v\n", dir)
	// Parse the templates in the dir.
	tmpls, err = d.Parse()
	if err != nil {
		log.Printf("d.Parse() error: %v", err)
		return
	}

	log.Printf("Parse template files in a dir successfully")
	if len(tmpls) == 0 {
		log.Printf("no template parsed")
		return
	}

	// List the parsed templates.
	fmt.Printf("parsed templates:\n")
	for _, tmpl := range tmpls {
		fmt.Printf("%v\n", strings.ReplaceAll(tmpl.Name(), string(os.PathSeparator), ">"))
	}

	// Execute the templates.
	for _, tmpl := range tmpls {
		log.Printf("execute template: %v\n", tmpl.Name())
		tmpl.Execute(os.Stderr, manual)
	}

	// Output:
	//parsed templates:
	//templates>markdown>chapters>00-about.md.tmpl
	//templates>markdown>chapters>01-installation.md.tmpl
	//templates>markdown>chapters>02-usage.md.tmpl
	//templates>markdown>title.md.tmpl
	//parsed templates:
	//templates>latex>chapters>00-about.tex.tmpl
	//templates>latex>chapters>01-installation.tex.tmpl
	//templates>latex>chapters>02-usage.tex.tmpl
	//templates>latex>manual.tex.tmpl
	//templates>latex>title.tex.tmpl
}

func ExampleDir_Render() {
	// Part I. Render markdown templates.
	dir := "templates/markdown"
	outputDir := filepath.Join(os.TempDir(), "templates/markdown")

	// Create a template dir.
	d := templatedir.New(dir)

	log.Printf("Render template files in a dir starts...\ndir: %v\noutput dir: %v\n", dir, outputDir)
	// Render the templates in the dir.
	if err := d.Render(outputDir, manual); err != nil {
		log.Printf("d.Render() error: %v", err)
		return
	}

	log.Printf("Render template files in a dir successfully")

	// Part II. Render LaTex templates.
	dir = "templates/latex"
	outputDir = filepath.Join(os.TempDir(), "templates/latex")

	// Create a template dir with specified delimiters.
	// For template files in templates/latex,
	// use \{\{ and \}\} as delimiters because '{' and '}' already used by LaTex.
	// It makes the template files can be opened / edited by LaTex editors successfully.
	d = templatedir.New(dir, templatedir.Delims("\\{\\{", "\\}\\}"))

	log.Printf("Render template files in a dir starts...\ndir: %v\noutput dir: %v\n", dir, outputDir)
	// Render the templates in the dir.
	if err := d.Render(outputDir, manual); err != nil {
		log.Printf("d.Render() error: %v", err)
		return
	}

	log.Printf("Render template files in a dir successfully")

	// Output:
}
