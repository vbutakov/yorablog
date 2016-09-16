package yotemplate

import (
	"html/template"
	"io"
	"log"
	"os"
	"time"
)

// Template is an usual html template but file modification aware
type Template struct {
	template *template.Template
	pathes   []string
	modTimes []time.Time
}

// InitYoTemplate initializes YoTemplate structure
func InitTemplate(pathes ...string) (*Template, error) {

	modTimes := make([]time.Time, len(pathes))
	p := make([]string, len(pathes))

	for i, path := range pathes {
		stat, err := os.Stat(path)
		if err != nil {
			log.Printf("Error. Cannot open template file '%v': %v\n", path, err)
			return nil, err
		}

		modTimes[i] = stat.ModTime()
		p[i] = path
	}

	templ, err := template.ParseFiles(pathes...)
	if err != nil {
		log.Printf("Error. Cannot parse template files: %v\n", err)
		return nil, err
	}

	t := &Template{
		template: templ,
		pathes:   p,
		modTimes: modTimes}
	return t, nil
}

// Execute checks if template is up to date and executes it
func (yt *Template) Execute(wr io.Writer, data interface{}) error {

	needReRead := false
	var err error

	for i, modTime := range yt.modTimes {
		path := yt.pathes[i]
		stat, err := os.Stat(path)
		if err != nil {
			log.Printf("Error. Cannot open template file '%v': %v.\n", path, err)
			return err
		}

		if modTime.Before(stat.ModTime()) {
			needReRead = true
		}
	}

	if needReRead {
		var templ *template.Template
		templ, err = template.ParseFiles(yt.pathes...)
		if err != nil {
			log.Printf("Error. Cannot parse template files: %v.\n", err)
			return err
		}
		log.Printf("Info. Templates are updated.\n")
		yt.template = templ
	}

	err = yt.template.Execute(wr, data)

	return err
}
