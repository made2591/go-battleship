package main

import (

	"os"
	"fmt"
	"bufio"
	"bytes"
	"net/http"
	"encoding/json"

	core "github.com/made2591/go-battleship/core"
	util "github.com/made2591/go-battleship/util"

)

func start(g *core.Game) {

	fmt.Println(">>> request new game...")
	res, _ := http.Get("http://" + HOST_NAME + ":" + HOST_PORT + START_REQUEST)
	json.NewDecoder(res.Body).Decode(g)
	core.PrettyPrintGame(g, 0)

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
		var x, y int

		for {
			fmt.Printf("x: ")
			_, e1 := fmt.Scanf("%d", &x)
			fmt.Printf("y: ")
			_, e2 := fmt.Scanf("%d", &y)
			if e1 == nil || e2 == nil {
				break
			}
		}

		core.GunShot(&g.FirstPlayer, &g.SecondPlayer, &core.Coordinates{Abscissa: x, Ordinate: y})

//		bbb, _ := json.Marshal(g)
//		fmt.Println(string(bbb))

//		fmt.Printf(">>> press ENTER to go on...\n")
//		reader.ReadString('\n')

		core.PrettyPrintGame(g, 0)

		fmt.Printf(">>> press ENTER to shot...\n")
		reader.ReadString('\n')

		js, _ := json.Marshal(g)
		res, _ := http.Post("http://"+HOST_NAME+":"+HOST_PORT+SHOT_REQUEST, "application/json", bytes.NewBuffer(js))
		json.NewDecoder(res.Body).Decode(g)

//		fmt.Printf(">>> shot received in coordinates [%d, %d]\n",
//			g.FirstPlayer.Suffered[len(g.FirstPlayer.Suffered)-1].Abscissa,
//			g.FirstPlayer.Suffered[len(g.FirstPlayer.Suffered)-1].Ordinate)
//		fmt.Printf(">>> press ENTER to go on...\n")
		reader.ReadString('\n')

		core.PrettyPrintGame(g, 0)
	case 3:
		fmt.Println(">>> exit game...")
		http.Get("http://" + HOST_NAME + ":" + HOST_PORT + EXIT_REQUEST)
		util.CleanScreen()
		os.Exit(1)
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
