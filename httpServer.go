package gowiki

//Imports
import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
)

//Struct for the filesystem
type justFilesFilesystem struct {
	fs http.FileSystem
}

//
type neuteredReaddirFile struct {
	http.File
}

// --------------- GLOBALS --------------- //
var templates = template.Must(template.ParseFiles("templ/edit.html", "templ/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

// --------------- START HANDLERS --------------- //
//view function
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

//Edit function.
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

//Save function
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}

	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/*
	Validates and creates a handler if the request is valid!

	If the request is valid then continue to the handler that
	is passed in the parameter.
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

//Handles root requests.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

// --------------- END HANDLERS --------------- //

//Renders a given HTML template
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Added for security reasons
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)

	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

//Main function
func main() {
	fmt.Println("Webserver started at: 127.0.0.1:8080")
	flag.Parse()

	//Setup the Handlers.
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	fs := justFilesFilesystem{http.Dir("templ/")}
	http.Handle("/templ/", http.StripPrefix("/templ/", http.FileServer(fs)))

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")

		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)

		if err != nil {
			log.Fatal(err)
		}

		s := &http.Server{}
		s.Serve(l)
		return
	}
	//Listen to port 8080
	http.ListenAndServe(":8080", nil)
}
