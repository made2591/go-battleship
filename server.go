package main

import (
	"encoding/json"
	"fmt"
	core "github.com/made2591/go-battleship/core"
	util "github.com/made2591/go-battleship/util"
	"net/http"
)

const (
	HOST_NAME = "localhost"
	HOST_PORT = "8080"
)

func startRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> new game generation...")
	g := core.PrepareGame(10, 0, "Matteo", 5, 9999, "HAL", 5, 9999)
	core.PrettyPrintGame(&g, 1)
	json.NewEncoder(w).Encode(g)
}

func main() {
	util.CleanScreen()
	http.HandleFunc("/start", startRequest)
	http.HandleFunc("/gunshot", core.ServerGunShot)
	http.HandleFunc("/exit", util.Exit)
	http.ListenAndServe(":8080", nil)
}