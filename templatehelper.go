package templatehelper

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	DefTmplExt = ".tmpl"
)

// ParseDirWithDelims parses all template files in the given dir and subdirs recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the relative path of the template file.
// The path contains the "dir" argument as a prefix.
// dir: the dir contains template files.
// ext: extend name of template file(e.g. ".tmpl"). It'll use ".tmpl" as default extend name if ext is empty.
// leftDelim / rightDelim: left / right delimiter of the template.
// It'll use default delimiters("{{" and "}}") of Golang if any of them is empty.
func ParseDirWithDelims(dir, ext, leftDelim, rightDelim string) ([]*template.Template, error) {
	var tmpls []*template.Template

	// Set ext to default template ext if it's empty.
	if ext == "" {
		ext = DefTmplExt
	}
	ext = strings.ToLower(ext)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Parse templates in sub dir recursively.
			subDir := filepath.Join(dir, entry.Name())

			tmplsInSubDir, err := ParseDirWithDelims(subDir, ext, leftDelim, rightDelim)
			if err != nil {
				return nil, err
			}

			tmpls = append(tmpls, tmplsInSubDir...)
		} else {
			filename := entry.Name()

			if strings.ToLower(filepath.Ext(filename)) != ext {
				continue
			}

			path := filepath.Join(dir, filename)
			absPath, err := filepath.Abs(path)
			if err != nil {
				return nil, err
			}

			// Read the content from the template file.
			data, err := os.ReadFile(absPath)
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
			if leftDelim != "" && rightDelim != "" {
				t = t.Delims(leftDelim, rightDelim)
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

// ParseDir parses all template files in the given dir and subdirs recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the relative path of the template file.
// The path contains the "dir" argument as a prefix.
// dir: the dir contains template files.
// ext: extend name of template file(e.g. ".tmpl"). It'll use ".tmpl" as default extend name if ext is empty.
func ParseDir(dir, ext string) ([]*template.Template, error) {
	return ParseDirWithDelims(dir, ext, "", "")
}

// ParseFSDirWithDelims parses all template files in the FS and named dir recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the relative path of the template file.
// The path contains the "dir" argument as a prefix.
// fsys: file system.
// dir: the dir contains template files.
// ext: extend name of template file(e.g. ".tmpl"). It'll use ".tmpl" as default extend name if ext is empty.
// leftDelim / rightDelim: left / right delimiter of the template.
// It'll use default delimiters("{{" and "}}") of Golang if any of them is empty.
func ParseFSDirWithDelims(fsys fs.FS, dir, ext, leftDelim, rightDelim string) ([]*template.Template, error) {
	var tmpls []*template.Template

	// Set ext to default template ext if it's empty.
	if ext == "" {
		ext = DefTmplExt
	}
	ext = strings.ToLower(ext)

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Parse templates in sub dir recursively.
			subDir := filepath.Join(dir, entry.Name())

			tmplsInSubDir, err := ParseFSDirWithDelims(fsys, subDir, ext, leftDelim, rightDelim)
			if err != nil {
				return nil, err
			}

			tmpls = append(tmpls, tmplsInSubDir...)
		} else {
			filename := entry.Name()

			if strings.ToLower(filepath.Ext(filename)) != ext {
				continue
			}

			path := filepath.Join(dir, filename)
			if err != nil {
				return nil, err
			}

			// Read the content from the template file.
			data, err := fs.ReadFile(fsys, path)
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
			if leftDelim != "" && rightDelim != "" {
				t = t.Delims(leftDelim, rightDelim)
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

// ParseFSDir parses all template files in the file system and named dir recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the relative path of the template file.
// The path contains the "dir" argument as a prefix.
// fsys: file system.
// dir: the dir contains template files.
// ext: extend name of template file(e.g. ".tmpl"). It'll use ".tmpl" as default extend name if ext is empty.
func ParseFSDir(fsys fs.FS, dir, ext string) ([]*template.Template, error) {
	return ParseFSDirWithDelims(fsys, dir, ext, "", "")
}
