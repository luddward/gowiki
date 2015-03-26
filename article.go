package gowiki

import (
	"io/ioutil"
)

//Data struct for article
type Page struct {
	Title string
	Body  []byte
}

// --------------- SAVE AND LOAD ------------- //
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("data/"+filename, p.Body, 0600)
}

//load function
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("data/" + filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

// --------------- END SAVE AND LOAD ------------ //
