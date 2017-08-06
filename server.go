package main

import (

	"fmt"
	"time"
	"net/http"
	"encoding/json"

	core "github.com/made2591/go-battleship/core"
	util "github.com/made2591/go-battleship/util"

)

const (

	PROTOCOLL = "http"
	HOST_NAME = "localhost"
	HOST_PORT = "8080"

	START_REQUEST = "/start"
	SHOT_REQUEST  = "/shot"
	EXIT_REQUEST  = "/exit"

	SIMULATION_THINKING_TIME = 2000 //milliseconds

)

// ###########################################################################################################
// ############################################# SERVER LOGIC ################################################
// ###########################################################################################################

func SleepRequest() {
	time.Sleep(SIMULATION_THINKING_TIME * time.Millisecond)
}

func StartRequest(w http.ResponseWriter, r *http.Request) {

	// create request
	fmt.Println(">>> new game generation...")

	// init game
	g := core.PrepareGame(core.PC_GRID, core.PC_MODE, "Matteo", core.PC_SHIPS, core.PC_SHOTS, "", -1, -1)

	// pretty print game
	core.PrettyPrintGame(&g)

	// IMPORTANT STEP
	core.SwitchPointOfView(&g)

	// send back to client
	json.NewEncoder(w).Encode(g)

}

// ServerGunShot from p Player to t Player in p Coordinates
func ServerShotRequest(w http.ResponseWriter, r *http.Request) {

	// decode received game AND SWITCH IT
	g := core.Decode(r)

	// IMPORTANT STEP
	core.SwitchPointOfView(g)

	// debug print of json
	// fmt.Println(core.DebugGame(g))

	// debug pause
	util.ConsolePause(util.PAUSE_MEX)

	// create Random SHOT
	// TODO: IMPLEMENT STRATEGY
	s := util.Random(0, len(g.SecondPlayer.Sea.Ships)-1)
	p := util.Random(0, len(g.SecondPlayer.Sea.Ships[s].Positions)-1)
	core.GunShot(&g.SecondPlayer, &g.FirstPlayer, &g.SecondPlayer.Sea.Ships[s].Positions[p])

	// debug pause
	util.ConsolePause(util.PAUSE_MEX)

	// print game
	fmt.Printf(core.PrettyPrintGame(g))

	// IMPORTANT STEP
	core.SwitchPointOfView(g)

	// send back to client
	json.NewEncoder(w).Encode(g)

}

func main() {

	// clean screen
	util.CleanScreen()

	// handle routes
	http.HandleFunc(START_REQUEST, StartRequest)
	http.HandleFunc(SHOT_REQUEST, ServerShotRequest)
	http.HandleFunc(EXIT_REQUEST, util.Exit)

	// start serving
	http.ListenAndServe(PROTOCOLL+"://"+HOST_NAME+":"+HOST_PORT, nil)

}