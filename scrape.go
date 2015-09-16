package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bentranter/chalk"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// CourseMatcher looks for any `<tr>` element with the
// class `timetable-course-one` or `timetable-course-two`
var CourseMatcher = func(n *html.Node) bool {
	if n.DataAtom == atom.Tr {
		return scrape.Attr(n, "class") == "timetable-course-two" || scrape.Attr(n, "class") == "timetable-course-one"
	}
	return false
}

// URLs are all the absolute URLs that need to be scraped
var URLs = []string{
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

// Course holds all the data for a single course at
// Lakehead. Some of these will need to change into maps
// in order to have more usable data (ie, `Dates` should
// become `StartDate` and `EndDate`).
type Course struct {
	Code         string  // PHYS-1030-FA
	Title        string  // Intro Appl Phys I (Mechanics)
	Type         string  // LEC
	RoomNo       string  // RB 1046
	Weekdays     string  // MWF
	Times        string  // 1:00PM - 2:00PM
	Dates        string  // 09/14/15 - 03/12/15
	Weight       float64 // 0.5
	Synonym      int     // 643369
	Instructor   string  // Dr. Mark C. Gallagher
	TextBookLink string  // https://textbooks.lakeheadu.ca...
	TextBook     string  // some struct with textbook info
}

// HTTPResponse is the struct that holds our response
// data
type HTTPResponse struct {
	url      string
	response *http.Response
	err      error
}

// Scrape finds and serializes the data from Lakehead's
// site. Eventually, it should return an array of
// `Match` reponses.
func Scrape(resp *HTTPResponse) {

	root, err := html.Parse(resp.response.Body)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	// Yhat's package expects atomic values for tags, see
	// https://godoc.org/golang.org/x/net/html/atom if you
	// need a different tag.
	data := scrape.FindAll(root, CourseMatcher)

	// Current course flag for looping, and counter for
	// total courses.
	currentCourse := ""
	counter := 0

	// currentCourseCounter counts which row of the `<tr>`
	// element we're in, so that we don't scrape out of
	// bounds.
	currentCourseCounter := 0

	// This actually works
	for _, match := range data {

		thisCourse := scrape.Attr(match, "class")

		// Set the courrentCourse flag
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
			currentCourseCounter = 0
		}

		// This actually works. Each case is a number that
		// matches a row (a `<tr>` element) in the current
		// course. Since attempting to scrape data that
		// doesn't exist in a row results in an out of
		// bounds error, this makes sure that we're never
		// scraping where there is no data. Keep in mind
		// scraping text from an element that has no value
		// but does exist won't throw an error, but trying
		// to scrape a non-existent attribute will throw an
		// error.
		switch currentCourseCounter {
		case 1:
			// Get the course code (string, since there's
			// no benefit from grabbing the int from it)
			courseCode := scrape.Text(match.FirstChild.NextSibling)
			fmt.Println(chalk.White(courseCode))

			// Get the course title
			courseTitle := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling)
			fmt.Println(courseTitle)

			// Get the course type. It _should_ be only one
			// of either `LEC`, `LAB`, or `WEB`... is a
			// boolean a better idea?
			courseType := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling)
			fmt.Println(courseType)

			// Get the course room number. It might be
			// worth it to parse the first two numbers to
			// get the building code, and then grab the
			// rest of the entry to get the room number.
			// Not sure how consistent all that is though.
			courseRoomNo := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling)
			fmt.Println(courseRoomNo)

			// Get the course days. This is the hardest one
			// to parse (probably) because it needs to be
			// seperated into an array of week days.
			//
			// @TODO: Parse into array of week days.
			courseWeekdays := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling)
			fmt.Println(courseWeekdays)

			// Get the course time. Needs to be parsed into
			// a start time and end time - shouldn't be
			// hard at all using `trim`.
			//
			// @TODO: Parse into start and end time.
			courseTimes := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling)
			fmt.Println(courseTimes)

			// Get the course dates. Needs to be parsed
			// into a start and end date.
			//
			// @TODO: Parse into start and end date.
			courseDates := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling)
			fmt.Println(courseDates)

			// Get the credit weight. Gets parsed into a
			// a float.
			courseWeightRaw := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling)
			courseWeight, err := strconv.ParseFloat(courseWeightRaw, 64)
			if err != nil {
				fmt.Println("Couldn't parse float: ", err)
			}
			fmt.Println(courseWeight)

		case 2:
			// Get the synonym, and turn it into an int
			courseSynonymRaw := scrape.Text(match.FirstChild.NextSibling)
			synonymNumStr := strings.TrimLeft(courseSynonymRaw, "Synonym: ")
			courseSynonym, err := strconv.ParseInt(synonymNumStr, 0, 64)
			if err != nil {
				fmt.Println("Couldn't convert string")
			}
			fmt.Println(courseSynonym)

			// Get the instructor name
			courseInstructorRaw := scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling)
			courseInstructor := strings.TrimLeft(courseInstructorRaw, "Instructor: ")
			fmt.Println(courseInstructor)

		case 3:
			// Get the course textbook link. This one is
			// annoying because you have to check to see
			// if there is actually anything there for you
			// to scrape.
			if scrape.Text(match.FirstChild.NextSibling.NextSibling.NextSibling) == "BOOKS" {
				courseTextbookLink := scrape.Attr(match.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild, "href")
				fmt.Println(courseTextbookLink)

				// Fun fact: Right here, you'll need to
				// make a `GET` request to that URL and
				// scrape the textbook page to get the
				// textbook info.
			} else {
				fmt.Println("No textbook exists for this course yet.")
			}
		case 4:
			fmt.Println("\n\n")
			// There's a fourth row, but never anything in
			// it
		default:
			// Do nothing
		}
	}

	return
}
