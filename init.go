/* Main file for NanoWiki,
 * this is like "The Powerhouse"
 * of the NanoWiki (or alternatively
 * the bane of your existence).
 */

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
    "strings"
    "runtime"
    "path/filepath"
)

func PackagePath() string {
    _, b, _, _ := runtime.Caller(0)
    basepath := filepath.Dir(b)
    splitpath := strings.Split(basepath, "/")
    path := strings.Join(splitpath[:len(splitpath)], "/")

    return path
}

/* Page struct, we need 
 * the body to be a []byte
 * because that's what the
 * modules expect of us.
 */
type Page struct {
	Title string
	Body  []byte
}

/* Page attributes. */
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

/* Page/Web Handlers */
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/* Template handler for all of the different modes. */
var templates = template.Must(template.ParseFiles(
    PackagePath() + "/tmpl/edit.html",
    PackagePath() + "/tmpl/view.html")) /* External files. */

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

/* This function is basically glue for all of the other functions. */
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

/* Main function; this straps all of is 
 * trash together, to create something functional.
 */
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	fmt.Println("Server has been initiated.")
	http.ListenAndServe(":8080", nil)
}
