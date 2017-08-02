package main

import (
	"encoding/json"
	"fmt"
	core "github.com/made2591/go-battleship/core"
	//util "github.com/made2591/go-battleship/util"
	"bufio"
	"bytes"
	"net/http"
	"os"
	"strconv"
	//"time"
)

const (
	HOST_NAME = "localhost"
	HOST_PORT = "8080"
)

func start(g *core.Game) {

	fmt.Println(">>> request new game...")
	res, _ := http.Get("http://" + HOST_NAME + ":" + HOST_PORT + "/start")
	json.NewDecoder(res.Body).Decode(g)
	core.NetPrintGame(g, 0)

}

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

func menu() int {

	menuOption := `Enter next move:
	1) Start
	2) GunShot
	3) Exit`

	choiceError := `Choice not available: press enter to continue...`
	i := 0

	for {
		//util.CleanScreen()
		fmt.Println(menuOption)
		_, err := fmt.Scanf("%d", &i)
		fmt.Println(i)
		if err != nil {
			fmt.Println(choiceError)
		} else {
			if checkChoice(i) {
				break
			} else {
				fmt.Println(choiceError)
			}
		}
	}
	return i

}

func play(a int, g *core.Game) {

	reader := bufio.NewReader(os.Stdin)

	switch a {
	case 1:
		start(g)
	case 2:
		fmt.Printf("x: ")
		x, _ := reader.ReadString('\n')
		ix, _ := strconv.Atoi(x)
		fmt.Printf("y: ")
		y, _ := reader.ReadString('\n')
		iy, _ := strconv.Atoi(y)
		g.GunShot(&g.FirstPlayer, &g.SecondPlayer, core.Coordinates{Abscissa: int(ix), Ordinate: int(iy)})
		bbb, _ := json.Marshal(g)
		fmt.Println(string(bbb))
		fmt.Printf(">>> press ENTER to go on...\n")
		reader.ReadString('\n')
		core.NetPrintGame(g, 0)
		fmt.Printf(">>> press ENTER to go on...\n")
		reader.ReadString('\n')
		js, _ := json.Marshal(g)
		res, _ := http.Post("http://"+HOST_NAME+":"+HOST_PORT+"/gunshot", "application/json", bytes.NewBuffer(js))
		json.NewDecoder(res.Body).Decode(g)
		core.NetPrintGame(g, 0)
	case 3:
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
