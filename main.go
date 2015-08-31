package bookstore

import (
	"fmt"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

// Match is a result returned from scraping
type Match struct {
	CourseCode string // PHYS-1030-FA
	Synonym    string // 643369
	Title      string // Intro Appl Phys I (Mechanics)
	Instructor string // Dr. Mark C. Gallagher
	Books      string // Link?
	Term       string // Fall
	Department string // Physics
	YearLevel  string // 1
}

func init() {
	http.HandleFunc("/", auth)
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func auth(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Hello, %v!", u)
}

// Scrape finds and serializes the data from Lakehead's
// site.
func Scrape(url string) []*Match {
	return nil
}
