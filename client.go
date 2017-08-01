package main

import (
	"encoding/json"
	"fmt"
	core "github.com/made2591/go-battleship/core"
	"net/http"
	//"time"
)

func main() {

	g := core.Game{}
	fmt.Println(">>> request new game...")
	res, _ := http.Get("http://localhost:8080/start")
	json.NewDecoder(res.Body).Decode(&g)
	fmt.Println(">>> Matteo")
	fmt.Println(core.StringfySea(g.FirstPlayer.Sea))
	fmt.Println(">>> HAL")
	fmt.Println(core.StringfySea(g.SecondPlayer.Sea))

}
