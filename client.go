package main

import (
	"encoding/json"
	"fmt"
	core "github.com/made2591/go-battleship/core"
	//util "github.com/made2591/go-battleship/util"
	"net/http"
	"os"
	"bytes"
	"bufio"
	"strconv"
	//"time"
)

const (
	HOST_NAME = "localhost"
	HOST_PORT = "8080"
)

func start(g *core.Game) {

	fmt.Println(">>> request new game...")
	res, _ := http.Get("http://"+HOST_NAME+":"+HOST_PORT+"/start")
	json.NewDecoder(res.Body).Decode(g)
	core.NetPrintGame(g)

}

func checkChoice(c int) bool {
	switch c {
		case 1:
			return true
		case 2:
			return true
		case 3:
			return true
		default:
			return false
	}
}

func menu() string {

	menuOption := `Enter next move:
	1) Start
	2) GunShot
	3) Exit`

	choiceError := `Choice not available: press enter to continue...`
	reader := bufio.NewReader(os.Stdin)
	text := ""
	
	for {
		//util.CleanScreen()
		fmt.Println(menuOption)
		text, _ = reader.ReadString('\n')
		fmt.Println(text)
		itext, err := strconv.Atoi(text)
		if err != nil {
			if checkChoice(itext) {
				break
			} else {
				fmt.Println(choiceError)
			}
		} else {
			fmt.Println(choiceError)
		}
	}
	
	return text

}

func play(a string, g *core.Game) {

	reader := bufio.NewReader(os.Stdin)

	switch a {
		case "1":
			start(g)
		case "2":
			fmt.Println("x: ")
			x, _ := reader.ReadString('\n')
			ix, _ := strconv.Atoi(x)
			fmt.Println("y: ")
			y, _ := reader.ReadString('\n')
			iy, _ := strconv.Atoi(y)
			g.GunShot(&g.FirstPlayer, &g.SecondPlayer, core.Coordinates{Abscissa: int(ix), Ordinate: int(iy)})
			core.NetPrintGame(g)
			jsonValue, _ := json.Marshal(g)
			res, _ := http.Post("http://"+HOST_NAME+":"+HOST_PORT+"/gunshot", "application/json", bytes.NewBuffer(jsonValue))
			json.NewDecoder(res.Body).Decode(g)
			core.NetPrintGame(g)
		case "3":
			break
	}

}

func main() {

	g := core.Game{}
	for {
		a := menu()
		play(a, &g)
	}

}