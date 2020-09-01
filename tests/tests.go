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

func getName(pipe string) (string, string) {
	name := strings.Split(pipe, "a: ")
	return name[0], name[1]
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
func newStatus(status string, setBy string) Status {
	s := Status{Condition: status, SetBy: setBy}
	return s
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
	p := Pokemon{Species: species, Alive: true, Status: Status{Condition: "", SetBy: ""}}
	return p
}

func (p *Pokemon) setStatus(status string, setBy string) {
	p.Status = newStatus(status, setBy)
}

func (p *Pokemon) kill(killer string) {
	p.Alive = false
	p.KilledBy = killer
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Functions
///////////////////////////////////////////////////////////////////////////////////////////////////

func parseLine(lines *[]string, i int, teams map[string]Player, ref map[string]string) {
	pipes := strings.Split((*lines)[i], "|")
	keyword := pipes[1]

	switch keyword {
	case "poke":
		species := strings.Split(pipes[3], ",")[0]
		player := pipes[2]
		teams[player].addMon(species)
		break
	case "switch":
		nickname := strings.Split(pipes[2], "a: ")[1]
		species := strings.Split(pipes[3], ",")[0]
		// check if nickname was taken, take it if not
		if _, ok := ref[nickname]; ok {
			fmt.Printf("switched to %s\n", nickname)
		} else {
			ref[nickname] = species
		}
		break
	case "-status":
		// TODO - Handle self-inflicting status (e.g. Toxic Orb) (currently throws an error on line 46 of logs)

		previous := strings.Split((*lines)[i-1], "|")
		_, setByName := getName(previous[1])
		setBy := ref[setByName]
		status := previous[2]
		defTrainer, targetName := getName(pipes[1])
		target := ref[targetName]
		t := teams[defTrainer].Team[target]
		t.setStatus(status, setBy)
		break
	case "-damage":
		//next := strings.Split((*lines)[i-1], "|")
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

	reference := make(map[string]string)

	teams := make(map[string]Player)
	teams["p1"] = newPlayer("p1", matchlog.P1)
	teams["p2"] = newPlayer("p2", matchlog.P2)

	// // Separate for parsing
	lines := strings.Split(matchlog.Log, "\n")

	// Parse
	for i := 0; i < 522; i++ {
		parseLine(&lines, i, teams, reference)
	}

	//Print all
	pretty, _ := json.MarshalIndent(reference, "", "	")
	fmt.Printf("%s", pretty)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Literally main()
///////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	// Get Log (testing)
	u := "https://replay.pokemonshowdown.com/gen8nationaldex-1149941851.json"

	summ(u)

}
