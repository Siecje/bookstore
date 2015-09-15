package main

import (
	"fmt"
	"net/http"
)

var testURL = []string{"http://timetable.lakeheadu.ca/2015FW_UG_TBAY/anth.html"}

func main() {

	// Create a new `app` instance
	app := &App{}

	// Route handlers
	http.Handle("/", app.Handle(HomeHandler))

	// Test the fetching script
	Fetch(testURL)

	// Listen and serve
	fmt.Println("Serving HTTP on port 3000")
	http.ListenAndServe(":3000", nil)
}
