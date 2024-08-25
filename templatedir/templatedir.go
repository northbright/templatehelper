package templatedir

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/northbright/copy/copyfile"
	"github.com/northbright/pathelper"
	"github.com/northbright/templatehelper"
)

// Dir represents the template dir which contains templates.
type Dir struct {
	path       string
	ext        string
	leftDelim  string
	rightDelim string
}

// Option is optional parameters to create a [Dir].
type Option func(d *Dir)

// Ext returns the option to set the file name extension of the template files.
func Ext(ext string) Option {
	return func(d *Dir) {
		d.ext = strings.ToLower(ext)
	}
}

// Delims returns the option to set the left and right delimiters of the template files.
// Default Golang template delimiters: {{ }}.
func Delims(left, right string) Option {
	return func(d *Dir) {
		d.leftDelim = left
		d.rightDelim = right
	}
}

// New creates a [Dir].
// path: the dir contains template files.
func New(path string, options ...Option) *Dir {
	d := &Dir{
		path: path,
		ext:  templatehelper.DefaultTmplExt,
	}

	for _, option := range options {
		option(d)
	}

	return d
}

// Parse parses all template files in the dir and subdirs recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the full path of the template file.
// Unlike [text/template.ParseFiles] and [text/template.ParseGlob],
// it will keep all parsed templates
// when parsing multiple files with the same name in different directories.
// e.g. The template dir contains "DIR/a/foo.tmpl" and "DIR/b/foo.tmpl".
// Parse() stores "DIR/a/foo.tmpl" as the template named "DIR/a/foo.tmpl",
// while stores "DIR/b/foo.tmpl" as the template named "DIR/b/foo.tmpl".
func (d *Dir) Parse() ([]*template.Template, error) {
	var tmpls []*template.Template

	err := filepath.WalkDir(d.path, func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		// Check file extension name.
		filename := entry.Name()
		if strings.ToLower(filepath.Ext(filename)) != d.ext {
			return nil
		}

		// Read the content from the template file.
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Convert content from []byte to string via strings.Builder.
		var b strings.Builder
		if _, err = b.Write(data); err != nil {
			return err
		}

		// Create a new empty template which name is path.
		t := template.New(path)

		// Set delimiters if need.
		if d.leftDelim != "" && d.rightDelim != "" {
			t = t.Delims(d.leftDelim, d.rightDelim)
		}

		// Parse the template.
		if t, err = t.Parse(b.String()); err != nil {
			return err
		}

		tmpls = append(tmpls, t)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tmpls, nil
}

// Render parses all template files in the dir and subdirs recursively then applies all templates to the specific data and write to the files in the output dir.
// It uses template file name with extension name cut as the name of output file(e.g. src/xx.md.tmpl -> dst/xx.md).
// For the files in the template dir whose extension is not ".tmpl" or specified extension(e.g. ".jpg" or other assets), it copies the files to the output dir.
func (d *Dir) Render(outputDir string, data any) error {
	err := filepath.WalkDir(d.path, func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		// Check file extension name.
		filename := entry.Name()
		if strings.ToLower(filepath.Ext(filename)) != d.ext {
			// Not a template file, just copy it to output dir.

			// Make output file name.
			dst, _ := strings.CutPrefix(path, d.path)
			dst = filepath.Join(outputDir, dst)

			// Copy file.
			ctx := context.Background()
			return copyfile.Do(ctx, path, dst)
		}

		// Read the content from the template file.
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Convert content from []byte to string via strings.Builder.
		var b strings.Builder
		if _, err = b.Write(content); err != nil {
			return err
		}

		// Create a new empty template which name is path.
		t := template.New(path)

		// Set delimiters if need.
		if d.leftDelim != "" && d.rightDelim != "" {
			t = t.Delims(d.leftDelim, d.rightDelim)
		}

		// Parse the template.
		if t, err = t.Parse(b.String()); err != nil {
			return err
		}

		// Make output file name.
		// 1. Cut prefix: template dir root.
		dst, _ := strings.CutPrefix(path, d.path)
		// 2. Cut suffix: template files extension name(e.g. xx.md.tmpl -> xx.md).
		dst, _ = strings.CutSuffix(dst, d.ext)
		// 3. Add prefix: output dir.
		dst = filepath.Join(outputDir, dst)

		// Create output file's parent dir if need.
		dir := filepath.Dir(dst)
		if err = pathelper.CreateDirIfNotExists(dir, 0755); err != nil {
			return err
		}

		// Open output file.
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer f.Close()

		return t.Execute(f, data)
	})

	return err
}
