package bookstore

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", HomeHandler)
}

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
