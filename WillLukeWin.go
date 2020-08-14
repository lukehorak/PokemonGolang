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

// Will I win the OUBL this season?
func WillLukeWin(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Message string `json:"message"`
		Number  int    `json:"number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Of course he wont, he's terrible at Pokemon!")
		return
	}
	if d.Message == "" {
		fmt.Fprint(w, "Of course he wont, he's terrible at Pokemon!")
		return
	}
	fmt.Fprint(w, html.EscapeString(d.Message))
}
