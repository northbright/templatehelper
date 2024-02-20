package templatehelper

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	DefTmplExt = ".tmpl"
)

// parseDirWithDelims parses all template files in the named dir and subdirs recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the relative path of the template file.
// The path contains the "dir" argument as a prefix.
// fsys: The file system(fs.FS) contains the named dir.
// Set it to nil if dir is not in a fs.FS.
// dir: the dir contains template files.
// ext: extend name of template file(e.g. ".tmpl"). It'll use ".tmpl" as default extend name if ext is empty.
// leftDelim / rightDelim: left / right delimiter of the template.
// It'll use default delimiters("{{" and "}}") of Golang if any of them is empty.
func parseDirWithDelims(fsys fs.FS, dir, ext, leftDelim, rightDelim string) ([]*template.Template, error) {
	var (
		err     error
		tmpls   []*template.Template
		entries []fs.DirEntry
		data    []byte
	)

	// Set ext to default template ext if it's empty.
	if ext == "" {
		ext = DefTmplExt
	}
	ext = strings.ToLower(ext)

	if fsys != nil {
		// dir is in a fs.FS.
		entries, err = fs.ReadDir(fsys, dir)
	} else {
		// dir is not in a fs.FS.
		entries, err = os.ReadDir(dir)
	}

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subDir := ""
			if fsys != nil {
				subDir = path.Join(dir, entry.Name())
			} else {
				subDir = filepath.Join(dir, entry.Name())
			}

			// Parse templates in sub dir recursively.
			tmplsInSubDir, err := parseDirWithDelims(fsys, subDir, ext, leftDelim, rightDelim)
			if err != nil {
				return nil, err
			}

			tmpls = append(tmpls, tmplsInSubDir...)
		} else {
			filename := entry.Name()

			if strings.ToLower(filepath.Ext(filename)) != ext {
				continue
			}

			p := ""

			if fsys != nil {
				// dir is in a fs.FS:
				// always use forward slash('/') as path separtor.
				p = path.Join(dir, filename)

				// Read the content from the template file.
				data, err = fs.ReadFile(fsys, p)
			} else {
				// dir is not in a fs.FS:
				// use OS specified path separator.
				p = filepath.Join(dir, filename)

				// Read the content from the template file.
				data, err = os.ReadFile(p)
			}

			if err != nil {
				return nil, err
			}

			// Convert content from []byte to string via strings.Builder.
			var b strings.Builder
			if _, err = b.Write(data); err != nil {
				return nil, err
			}

			// Create a new empty template which name is path.
			t := template.New(p)

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

// ParseDirWithDelims parses all template files in the given dir and subdirs recursively.
// It returns a slice contains parsed templates.
// The name of each parsed template is set to the relative path of the template file.
// The path contains the "dir" argument as a prefix.
// dir: the dir contains template files.
// ext: extend name of template file(e.g. ".tmpl"). It'll use ".tmpl" as default extend name if ext is empty.
// leftDelim / rightDelim: left / right delimiter of the template.
// It'll use default delimiters("{{" and "}}") of Golang if any of them is empty.
func ParseDirWithDelims(dir, ext, leftDelim, rightDelim string) ([]*template.Template, error) {
	return parseDirWithDelims(nil, dir, ext, leftDelim, rightDelim)
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
	return parseDirWithDelims(fsys, dir, ext, leftDelim, rightDelim)
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
