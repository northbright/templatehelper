package parsedir

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/northbright/templatehelper"
)

// Parser represents the parser to parse templates in a dir.
type Parser struct {
	dir        string
	ext        string
	leftDelim  string
	rightDelim string
}

// Option is optional parameters to create a Parser.
type Option func(p *Parser)

// Ext returns the option to set the file name extension for the parser.
func Ext(ext string) Option {
	return func(p *Parser) {
		p.ext = strings.ToLower(ext)
	}
}

// Delims returns the option to set the left and right delimiters for the parser.
// Default Golang template delimiters: {{ }}.
func Delims(left, right string) Option {
	return func(p *Parser) {
		p.leftDelim = left
		p.rightDelim = right
	}
}

// New creates a parser.
// dir: the dir contains template files.
func New(dir string, options ...Option) *Parser {
	p := &Parser{
		dir: dir,
		ext: templatehelper.DefaultTmplExt,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// Parse parses all template files in the named dir and subdirs recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the path of the template file.
// The path contains the "dir" argument as a prefix.
func (p *Parser) Parse(ctx context.Context) ([]*template.Template, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
		entries, err := os.ReadDir(p.dir)
		if err != nil {
			return nil, err
		}

		var tmpls []*template.Template
		for _, entry := range entries {
			if entry.IsDir() {
				subDir := filepath.Join(p.dir, entry.Name())

				// Parse templates in sub dir recursively.
				parser := New(
					subDir,
					Ext(p.ext),
					Delims(p.leftDelim, p.rightDelim),
				)

				tmplsInSubDir, err := parser.Parse(ctx)
				if err != nil {
					return nil, err
				}

				tmpls = append(tmpls, tmplsInSubDir...)
			} else {
				filename := entry.Name()

				if strings.ToLower(filepath.Ext(filename)) != p.ext {
					continue
				}

				// use OS specified path separator.
				path := filepath.Join(p.dir, filename)

				// Read the content from the template file.
				data, err := os.ReadFile(path)

				if err != nil {
					return nil, err
				}

				// Convert content from []byte to string via strings.Builder.
				var b strings.Builder
				if _, err = b.Write(data); err != nil {
					return nil, err
				}

				// Create a new empty template which name is path.
				t := template.New(path)

				// Set delimiters if need.
				if p.leftDelim != "" && p.rightDelim != "" {
					t = t.Delims(p.leftDelim, p.rightDelim)
				}

				// Parse the template.
				if t, err = t.Parse(b.String()); err != nil {
					return nil, err
				}

				tmpls = append(tmpls, t)
			}
		}

		return tmpls, nil
	}
}