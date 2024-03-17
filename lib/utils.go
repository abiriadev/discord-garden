package lib

import (
	"path/filepath"
	"strings"
	"text/template"
)

func useTemplate(path string, data any) (string, error) {
	var buf strings.Builder

	if tmpl, err := template.New(filepath.Base(path)).ParseFiles(
		path,
	); err != nil {
		return buf.String(), err
	} else if err := tmpl.Execute(&buf, data); err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}
