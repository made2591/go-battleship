package main

import (

	"os"
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"

	core "github.com/made2591/go-battleship/core"
	util "github.com/made2591/go-battleship/util"

)

const (

	ERROR_CHOICE = `Choice not available: press enter to continue...`

)

// ###########################################################################################################
// ############################################# CLIENT LOGIC ################################################
// ###########################################################################################################

// CreateNewGame make a request to start route for game creation
func CreateNewGame(g *core.Game) {

	// create request
	res, _ := http.Get(PROTOCOLL+"://" + HOST_NAME + ":" + HOST_PORT + START_REQUEST)

	// game decoding and SWITCH
	g = core.Decode(res.Request)

	return

}

// DoGunShot make a request to gun shot and receive response
func DoGunShot(g *core.Game) {

	// prepare JSON
	js, _ := json.Marshal(g)

	// make gun shot request
	res, _ := http.Post("http://"+HOST_NAME+":"+HOST_PORT+SHOT_REQUEST, "application/json", bytes.NewBuffer(js))

	// decode game response
	json.NewDecoder(res.Body).Decode(g)

}

// DoGunShot make an exit request
func ExitGame() {

	http.Get(PROTOCOLL+"://" + HOST_NAME + ":" + HOST_PORT + EXIT_REQUEST)
	util.CleanScreen()
	os.Exit(1)

}

// ###########################################################################################################
// ############################################## CLIENT CLI #################################################
// ###########################################################################################################

// startNewGame make a request to start route for game creation
func startNewGame(g *core.Game) {

	// cli output
	fmt.Println(">>> request new game...")
	// create new game
	CreateNewGame(g)

}

// checkChoice available in CLI
func checkChoice(c int) (b bool) {
	switch c {
		case 1:
			b = true
		case 2:
			b = true
		case 3:
			b = true
		default:
			b = false
	}
	return
}

// menu
func menu() int {

	menuOption := `
	1) Start
	2) GunShot
	3) Exit
Enter next move`

	i := 0

	for {

		fmt.Println(menuOption)
		_, err := fmt.Scanf("%d", &i)

		if err != nil {
			fmt.Println(ERROR_CHOICE)

		} else {

			if checkChoice(i) {
				break
			} else {
				fmt.Println(ERROR_CHOICE)
			}

		}
	}
	return i

}

func getCoordinates() (c core.Coordinates) {

	var x, y int

	for {
		fmt.Printf("x: ")
		_, e1 := fmt.Scanf("%d", &x)
		fmt.Printf("y: ")
		_, e2 := fmt.Scanf("%d", &y)
		if e1 == nil && e2 == nil {
			c = core.Coordinates{}
			break
		} else {
			fmt.Println(ERROR_CHOICE)
		}
	}

	return

}

func play(a int, g *core.Game) {

	// switch choices
	switch a {

		// start new game
		case 1:
			startNewGame(g)

		// do a gunshot
		case 2:

			// get coordinates
			c := getCoordinates()

			// gun shot
			core.GunShot(&g.FirstPlayer, &g.SecondPlayer, &c)

			// pretty print game
			core.PrettyPrintGame(g)

			// gun shot
			DoGunShot(g)

			// pretty print game
			core.PrettyPrintGame(g)

		// exit game
		case 3:

			fmt.Println(">>> exit game...")
			ExitGame()

		}

}

func main() {
	util.CleanScreen()
	g := core.Game{}
	for {
		a := menu()
		play(a, &g)
	}
}
