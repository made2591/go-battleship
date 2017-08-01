package main

import (
	"encoding/json"
	"fmt"
	core "github.com/made2591/go-battleship/core"
	"net/http"
)

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> new game generation...")
	g := core.PrepareGame(10, 0, "Matteo", 5, 9999, "HAL", 5, 9999)
	fmt.Println(">>> Matteo")
	fmt.Println(core.StringfySea(g.FirstPlayer.Sea))
	fmt.Println(">>> HAL")
	fmt.Println(core.StringfySea(g.SecondPlayer.Sea))
	json.NewEncoder(w).Encode(g)
}

func main() {
	http.HandleFunc("/start", start)
	//http.HandleFunc("/gunshot", gunshot)
	http.ListenAndServe(":8080", nil)
}
