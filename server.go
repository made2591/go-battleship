package main

import (

	"fmt"
	"time"
	"net/http"
	"encoding/json"

	core "github.com/made2591/go-battleship/core"
	util "github.com/made2591/go-battleship/util"

)


// ###########################################################################################################
// ############################################# SERVER LOGIC ################################################
// ###########################################################################################################

func SleepRequest() {
	time.Sleep(core.SIMULATION_THINKING_TIME * time.Millisecond)
}

func StartRequest(w http.ResponseWriter, r *http.Request) {

	// create request
	fmt.Println(">>> new game generation...")

	// init game
	g := core.PrepareGame(core.PC_GRID, core.PC_MODE, "Matteo", core.PC_SHIPS, core.PC_SHOTS, "", -1, -1)

	// ROTATE POINT OF VIEW FOR PRINTING: SERVER BECOME FIRST PLAYER
	core.SwitchPointOfView(&g)

	// pretty print game
	fmt.Println(core.PrettyPrintGame(&g))

	// RELOAD ORIGINAL PERSPECTIVE: CLIENT BECOME FIRST PLAYER
	core.SwitchPointOfView(&g)

	// send back to client
	json.NewEncoder(w).Encode(g)

}

// ServerGunShot from p Player to t Player in p Coordinates
func ServerShotRequest(w http.ResponseWriter, r *http.Request) {

	// decode received game AND SWITCH IT
	g := core.GameDecoder(r)

	// ROTATE POINT OF VIEW FOR PRINTING: SERVER BECOME FIRST PLAYER
	core.SwitchPointOfView(&g)

	// debug print of json
	//fmt.Println(core.DebugGame(&g))

	// debug pause
	util.ConsolePause(util.PAUSE_MEX)

	// create Random SHOT on second PLAYER
	// TODO: IMPLEMENT STRATEGY
	s := util.Random(0, len(g.SecondPlayer.Sea.Ships)-1)
	p := util.Random(0, len(g.SecondPlayer.Sea.Ships[s].Positions)-1)
	core.GunShot(&g.FirstPlayer, &g.SecondPlayer, &g.SecondPlayer.Sea.Ships[s].Positions[p])

	// debug pause
	util.ConsolePause(util.PAUSE_MEX)

	// print game
	fmt.Println(core.PrettyPrintGame(&g))

	// RELOAD ORIGINAL PERSPECTIVE: CLIENT BECOME FIRST PLAYER
	core.SwitchPointOfView(&g)

	// send back to client
	json.NewEncoder(w).Encode(g)

}

func main() {

	// clean screen
	util.CleanScreen()

	// handle routes
	http.HandleFunc(core.START_REQUEST, StartRequest)
	http.HandleFunc(core.SHOT_REQUEST, ServerShotRequest)
	http.HandleFunc(core.EXIT_REQUEST, util.Exit)

	// start serving
	http.ListenAndServe(":"+core.HOST_PORT, nil)

}