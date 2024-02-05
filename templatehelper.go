package templatehelper

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

func ParseDir(dir, ext, leftDelim, rightDelim string) (map[string]*template.Template, error) {
	m := make(map[string]*template.Template)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) != ext {
			return nil
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		t := template.New("").Delims(leftDelim, rightDelim)

		if t, err = t.ParseFiles(absPath); err != nil {
			return err
		}

		if tmpls := t.Templates(); err != nil {
			if len(tmpls) != 1 {
				return fmt.Errorf("associated templates num is not 1")
			}
			m[absPath] = tmpls[0]
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}
