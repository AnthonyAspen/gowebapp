package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// the struct is representing page
type Page struct {
	ImagePath string
	Title     string
	Body      []byte
}

func (p *Page) save() error {
	path := "./data/"
	filename := path + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	path := "./data/"
	filename := path + title + ".txt"
	imagePath := "/assets/" + title + ".jpg"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body, ImagePath: imagePath}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploadFile(w, r, title)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// a function to upload an image from the form with name="image"
func uploadFile(w http.ResponseWriter, r *http.Request, pageName string) {

	// UPLOADING AN IMAGE PROCESS IS HERE
	// 1. parse input, type multipart/form-data.
	r.ParseMultipartForm(10000000)
	// 2. retrieve file from posted form-data

	file, _, err := r.FormFile("image")
	if file == nil {
		fmt.Println("empty file")
		return
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("./assets/"+pageName+".jpg", fileBytes, 0644)
}
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {

	path := r.URL.Path
	if strings.HasPrefix(path, "/view/") || strings.HasPrefix(path, "/save/") || strings.HasPrefix(path, "/edit/") {
		validPath := path[6:]
		return validPath, nil
	}
	return "", errors.New("Invalid Page Request")
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this is a normal way to find the valid path, but I changed it to
// comprehend it better
//var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
