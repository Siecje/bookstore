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
	http.HandleFunc("/scrape", HomeHandler)
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Hello, Benny world!")
	Scrape("http://timetable.lakeheadu.ca/2015FW_UG_TBAY/phys.html", r)
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
func Scrape(url string, r *http.Request) (Matches []*Match, err error) {

	// Required for logging
	c := appengine.NewContext(r)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error :", err)
		return nil, err
	}

	root, err := html.Parse(res.Body)
	if err != nil {
		fmt.Println("Error :", err)
		return nil, err
	}

	// Yhat's package expects atomic values for tags, see
	// https://godoc.org/golang.org/x/net/html/atom if you
	// need a different tag.
	data := scrape.FindAll(root, scrape.ByTag(0x10502))
	// matches := make([]*Match, len(data))
	for _, match := range data {
		c.Infof("Match: %v", match)
	}

	return nil, nil
}
