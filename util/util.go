package util

import (
	"math/rand"
	"time"
	"fmt"
	"net/http"
	"os"	
)

const (
	ROW_NUMBER = 1024
	COL_NUMBER = 1024
	BLANK_SPAC = " "
)

func Random(min, max int) int {
	max = max + 1
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func Search(a int, b []int) bool {
	for _, v := range b {
		if v == a {
			return true
		}
	}
	return false
}

func CleanScreen() {
	r, c := 0, 0
	for r < ROW_NUMBER {
		for c < COL_NUMBER {
			fmt.Printf(BLANK_SPAC)
			c++
		}
		fmt.Printf(BLANK_SPAC+"\n")
		r++
	}
	fmt.Printf("\033[0;0H")

}

func Exit(w http.ResponseWriter, r *http.Request) {
	CleanScreen()
	os.Exit(1)
}