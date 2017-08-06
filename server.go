package main

import (

	"fmt"
	"net/http"
	"encoding/json"

	core "github.com/made2591/go-battleship/core"
	util "github.com/made2591/go-battleship/util"

)

const (

	HOST_NAME = "localhost"
	HOST_PORT = "8080"

	START_REQUEST = "/start"
	SHOT_REQUEST  = "/shot"
	EXIT_REQUEST  = "/exit"

)

func startRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> new game generation...")
	g := core.PrepareGame(10, 0, "Matteo", 5, 9999, "HAL", 5, 9999)
	core.PrettyPrintGame(&g, 1)
	json.NewEncoder(w).Encode(g)
}

// ServerGunShot from p Player to t Player in p Coordinates
func serverShotRequest(w http.ResponseWriter, r *http.Request) {

	// decode game
	d := json.NewDecoder(r.Body)
	g := core.Game{}
	err := d.Decode(&g)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// debug print of json
	fmt.Println(core.DebugGame(&g))


	util.ConsolePause(util.PAUSE_MEX)

	s := util.Random(0, len(g.SecondPlayer.Sea.Ships)-1)
	p := util.Random(0, len(g.SecondPlayer.Sea.Ships[s].Positions)-1)
	GunShot(&g.SecondPlayer, &g.FirstPlayer, &g.SecondPlayer.Sea.Ships[s].Positions[p])

	fmt.Printf(">>> gun shot coordinates [%d, %d]\n",
		g.SecondPlayer.Sea.Ships[s].Positions[p].Abscissa,
		g.SecondPlayer.Sea.Ships[s].Positions[p].Ordinate)
	fmt.Printf(">>> press ENTER to go on...\n")
	reader.ReadString('\n')

	PrettyPrintGame(&g, 1)

	fmt.Printf(">>> press ENTER to go on...\n")
	reader.ReadString('\n')

	PrettyPrintGame(&g, 1)

	json.NewEncoder(w).Encode(g)
}

func main() {
	util.CleanScreen()
	http.HandleFunc("/start", startRequest)
	http.HandleFunc("/shot", serverShotRequest)
	http.HandleFunc("/exit", util.Exit)
	http.ListenAndServe(":8080", nil)
}