// Cloud function for summarizing Pokemon Showdown matches for draft-league.com/oubl
// Package p contains an HTTP Cloud Function.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
// Structs & Methods
///////////////////////////////////////////////////////////////////////////////////////////////////
type MatchLog struct {
	Id  string
	P1  string
	P2  string
	Log string
}

/////////////////////////////////////////////
// Player type
/////////////////////////////////////////////
type Player struct {
	Username string
	Team     map[string]Pokemon
	Id       string
}

func newPlayer(id string, username string) Player {
	p := Player{Username: username, Id: id, Team: make(map[string]Pokemon)}
	return p
}
func (p Player) addMon(species string) {
	poke := newPokemon(species)
	p.Team[species] = poke
	return
}

/////////////////////////////////////////////
// Status type
/////////////////////////////////////////////
type Status struct {
	Condition string
	SetBy     string
}

// set status on Pokemon and record who set it
func (s Status) setStatus(status string, setBy string) {
	s.Condition = status
	s.SetBy = setBy
}

/////////////////////////////////////////////
// Pokemon type
/////////////////////////////////////////////
type Pokemon struct {
	Name     string
	Species  string
	Alive    bool
	Direct   bool
	KilledBy string
	Status   Status
}

func newPokemon(species string) Pokemon {
	p := Pokemon{Species: species, Alive: true, Status: Status{Condition: "none", SetBy: "none"}}
	return p
}

func (p *Pokemon) kill(killer string) {
	p.Alive = false
	p.KilledBy = killer
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Functions
///////////////////////////////////////////////////////////////////////////////////////////////////

func parseLine(line string, teams map[string]Player) {
	pipes := strings.Split(line, "|")
	keyword := pipes[1]

	switch keyword {
	case "poke":
		species := strings.Split(pipes[3], ",")[0]
		player := pipes[2]

		fmt.Printf("Player --> %s \n\n", player)
		teams[player].addMon(species)
		break
	}
	return

}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Main Function that's not actually main()
///////////////////////////////////////////////////////////////////////////////////////////////////

func summ(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Do the JSON
	var matchlog MatchLog
	errr := json.Unmarshal(body, &matchlog)
	if errr != nil {
		log.Println(errr)
	}

	teams := make(map[string]Player)
	teams["p1"] = newPlayer("p1", matchlog.P1)
	teams["p2"] = newPlayer("p2", matchlog.P2)

	// // Separate for parsing
	lines := strings.Split(matchlog.Log, "\n")

	for i := 0; i < 30; i++ {
		parseLine(lines[i], teams)
	}

	//Print all
	//pretty, _ := json.MarshalIndent(matchlog, "", "	")
	fmt.Printf("%v", teams["p1"].Team)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Literally main()
///////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	// Get Log (testing)
	u := "https://replay.pokemonshowdown.com/gen8nationaldex-1149941851.json"

	summ(u)

}
