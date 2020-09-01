// Test function for learning how to upload/use Google cloud funcitons
//

// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

type Result struct {
	Message string `json:"Message"`
}

// Will I win the OUBL this season?
func Summarize(w http.ResponseWriter, r *http.Request) {

	var d struct {
		Url string `json:"url"`
	}

	//reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error!")
		return
	}
	if d.Url == "" {
		fmt.Fprint(w, "No logs sent. Please send Showdown logs in 'body' ")
		return
	}

	fmt.Fprintf(w, "The Url sent is %s", html.EscapeString(d.Url))
}
