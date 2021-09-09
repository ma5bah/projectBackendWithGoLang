package main

import (
	"io/ioutil"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (i Page) Save() {
	fileName := i.Title + ".txt"
	ioutil.WriteFile(os.Getenv("templateDirectory")+fileName, i.Body, 0600)
	//body,_:=ioutil.ReadFile(fileName)
}
func LoadPage(Title string) (*Page, error) {
	fileName := Title + ".txt"
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Page{Title: Title, Body: body}, nil
}
