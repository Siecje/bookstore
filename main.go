package main

import (
	"fmt"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

var testURL = []string{"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/anth.html"}

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
	// Current course flag for looping, and counter for
	// total courses
	currentCourse := ""
	counter := 0

	// currentCourseCounter counts which row of the `<tr>`
	// element we're in, so that we don't scrape out of
	// bounds
	currentCourseCounter := 0

	// This actually works
	for _, match := range data {

		thisCourse := scrape.Attr(match, "class")

		switch scrape.Attr(match, "class") {
		case "timetable-course-two":
			currentCourse = "timetable-course-two"
		case "timetable-course-one":
			currentCourse = "timetable-course-two"
		default:
			// Do nothing
		}

		currentCourseCounter++

		if currentCourse != thisCourse {
			counter++
			fmt.Println("Counter:", counter)
			currentCourseCounter = 0
		}

		// This actually works
		switch currentCourseCounter {
		case 1:
			courseCode := scrape.Text(match.FirstChild.NextSibling)
			fmt.Println(courseCode)
		case 2:
			synonym := scrape.Text(match)
			fmt.Println(synonym)
		case 3:
			fmt.Println(3)
		case 4:
			fmt.Println(4)
		default:
			// Do nothing
		}

		// fmt.Println(
		// 	i,
		// 	scrape.Text(match),
		// 	scrape.Text(match.FirstChild.NextSibling),             // ANTH-1032-FA
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling), // Synonym: 64549
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling),
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling),
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling),
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling),
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling),
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling),
		// 	scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling),
		// )
	}

	return
}

func main() {
	Fetch(testURL)
}
