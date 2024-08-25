# templatefsdir
Package templatefsdir provides functions to parse and render the templates in a file system dir.

## Documentation
* <https://pkg.go.dev/github.com/northbright/templatehelper/templatefsdir>

## Security
templatefsdir uses [text/template](https://pkg.go.dev/text/template) but not [html/template](https://pkg.go.dev/html/template) to make it possible to output raw HTML / JS / CSS code.

To secure HTML output, you may need to sanitize the input before execute the templates(e.g. using [bluemonday](https://github.com/microcosm-cc/bluemonday)).
