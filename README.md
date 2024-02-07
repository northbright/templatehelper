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
It returns a template slice and each template name is the relative path(with "dir" argument as prefix") of the template file. So the slice may contains multiple templates with same base file names(for [ParseFiles](https://pkg.go.dev/text/template#ParseFiles) of [text/template](https://pkg.go.dev/text/template), the last one will be the one that results).

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
                fmt.Printf("ParseDir() error: %%v\n", err)
        }

        // List the parsed temlates.
        fmt.Printf("Parsed templates:\n")
        for _, tmpl := range tmpls {
                fmt.Printf("%%v\n", tmpl.Name())
        }

        // Output:
        //Parsed templates:
        //templates/markdown/chapters/00-about.md
        //templates/markdown/chapters/01-installation.md
        //templates/markdown/chapters/02-usage.md
        //templates/markdown/title.md
}
```
