package bookstore

import (
	"fmt"
	"net/http"

	"github.com/mattbaird/elastigo"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

// urls are all the absolute URLs that need to be scraped
var urls = []string{
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/anth.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/apbi.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/biol.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/busi.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/chem.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/clas.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/comp.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/crim.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/econ.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/educ.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/engi.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/finn.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/fren.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/geoa.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/geog.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/geol.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/gero.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/gsci.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/hist.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/indi.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/intd.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/ital.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/kine.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/lang.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/laws.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/ling.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/math.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/mdst.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/meds.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/musi.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/nacc.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/nort.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/nrmt.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/nurs.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/ojib.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/outd.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/phil.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/phys.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/poli.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/psyc.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/reli.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/soci.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/sowk.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/span.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/visu.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/wate.html",
	"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/wome.html",
}

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
	http.HandleFunc("/", AuthHandler)
	http.HandleFunc("/scrape", HomeHandler)
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Home")
}

// AuthHandler allows a user to login and greets them
func AuthHandler(w http.ResponseWriter, r *http.Request) {
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

// HTTPResponse is the struct that holds our response
// data
type HTTPResponse struct {
	url      string
	response *http.Response
	err      error
}

// Scrape finds and serializes the data from Lakehead's
// site.
func Scrape(resp *HTTPResponse) {

	root, err := html.Parse(resp.response.Body)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	// Looks for any `<tr>` element with the class
	// `timetable-course-one` or `timetable-course-two`
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Tr {
			return scrape.Attr(n, "class") == "timetable-course-two" || scrape.Attr(n, "class") == "timetable-course-one"
		}
		return false
	}

	// Yhat's package expects atomic values for tags, see
	// https://godoc.org/golang.org/x/net/html/atom if you
	// need a different tag.
	data := scrape.FindAll(root, matcher)
	for _, match := range data {
		fmt.Println(scrape.Text(match))
		// scrape.Text(match.FirstChild.NextSibling),             // ANTH-1032-FA
		// scrape.Text(match.FirstChild.NextSibling.NextSibling), // Synonym: 64549

		//)
	}

	return
}

func asyncHTTPGet(urls []string) []*HTTPResponse {

	ch := make(chan *HTTPResponse)
	responses := []*HTTPResponse{}

	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s \n", url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error: ", err)
				ch <- &HTTPResponse{url, nil, err}
			} else {
				ch <- &HTTPResponse{url, resp, err}
			}
		}(url)
	}

	for {
		select {
		case r := <-ch:

			// Scrape
			Scrape(r)

			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}

			defer r.response.Body.Close()
		}
	}

	return responses
}
