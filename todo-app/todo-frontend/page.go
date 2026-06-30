package main

import (
	"fmt"
	"html/template"
)

type PageData struct {
	Title      string
	Message    string
	Image      string
	Tasks      []Task
	BackendURL string
}

func Print(t *Task) string {
	var state string
	switch t.State {
	case StateTodo:
		state = "[ ]"
	case StateDone:
		state = "[X]"
	}
	return fmt.Sprintf("%s %s\n", state, t.Title)
}

var tmplFuncs = template.FuncMap{
	"Print": Print,
}

func newTemplateCache() (*template.Template, error) {
	// The name must match the base name of the file you are parsing
	ts := template.New("index.tmpl").Funcs(tmplFuncs)

	// Because of fs.Sub, "ui/index.tmpl" is now just "index.tmpl"
	ts, err := ts.ParseFS(templateFS, "index.tmpl")
	if err != nil {
		return nil, err
	}

	return ts, nil
}
