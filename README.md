# templatehelper
A Go Library provides helper functions for template package.

## Installation
```shell
go get -u github.com/northbright/templatehelper
```

## Documentation
* <https://pkg.go.dev/github.com/northbright/templatehelper>

## Usage
#### Parse All Template Files in a Directory Recursively
Use [ParseDir](https://pkg.go.dev/github.com/northbright/templatehelper#ParseDir) to parse all files in a directory.

It returns a [template.Template](https://pkg.go.dev/text/template#Template) slice and each template name is the relative path(with "dir" argument as prefix") of the template file. So the slice may contains multiple templates with same base file names. e.g. "dir/foo.tmpl", "dir/a/foo.tmpl".

```golang
package main

import (
        "fmt"

        "github.com/northbright/templatehelper"
)

func main() {
        dir := "templates/markdown"
        tmpls, err := templatehelper.ParseDir(dir, ".md")
        if err != nil {
                fmt.Printf("ParseDir() error: %v\n", err)
        }

        // List the parsed temlates.
        fmt.Printf("Parsed templates:\n")
        for _, tmpl := range tmpls {
                fmt.Printf("%v\n", strings.ReplaceAll(tmpl.Name(), string(os.PathSeparator), ">"))
        }

        // Output:
        //Parsed templates:
        //templates>markdown>chapters>00-about.md
        //templates>markdown>chapters>01-installation.md
        //templates>markdown>chapters>02-usage.md
        //templates>markdown>title.md
}
```

## Security
templatehelper uses [text/template](https://pkg.go.dev/text/template) but not [html/template](https://pkg.go.dev/html/template) to make it possible to output raw HTML / JS / CSS code.

To secure HTML output, you may need to sanitize the input before execute the templates(e.g. using [bluemonday](https://github.com/microcosm-cc/bluemonday)).
