package yotemplate

import (
	"html/template"
	"io"
	"log"
	"os"
	"time"
)

// YoTemplate is an usual html template but file modification aware
type YoTemplate struct {
	TemplatePath string
	Template     *template.Template
	ModTime      time.Time
}

// InitYoTemplate initializes YoTemplate structure
func InitYoTemplate(path string) (*YoTemplate, error) {
	stat, err := os.Stat(path)
	if err != nil {
		log.Printf("Error. Cannot open template file '%v': %v\n", path, err)
		return nil, err
	}

	templ, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Error. Cannot parse template file '%v': %v\n", path, err)
		return nil, err
	}

	t := &YoTemplate{
		TemplatePath: path,
		Template:     templ,
		ModTime:      stat.ModTime()}
	return t, nil
}

// Execute checks if template is up to date and executes it
func (yt *YoTemplate) Execute(wr io.Writer, data interface{}) error {
	stat, err := os.Stat(yt.TemplatePath)
	if err != nil {
		log.Printf("Error. Cannot open template file '%v': %v.\n", yt.TemplatePath, err)
		return err
	}

	if yt.ModTime.Before(stat.ModTime()) {
		var templ *template.Template
		templ, err = template.ParseFiles(yt.TemplatePath)
		if err != nil {
			log.Printf("Error. Cannot parse template file '%v': %v.\n", yt.TemplatePath, err)
			return err
		}
		log.Printf("Info. Template is updated from file '%v'.\n", yt.TemplatePath)
		yt.Template = templ
	}

	err = yt.Template.Execute(wr, data)

	return err
}
