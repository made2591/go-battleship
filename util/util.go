package util

import (
	"os"
	"fmt"
	"time"
	"bufio"
	"net/http"
	"math/rand"
)

const (

	ROW_NUMBER = 1024
	COL_NUMBER = 1024

	BLANK_SPAC = " "
	PAUSE_MEX  = ">>> press ENTER to go on..."

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

func ConsolePause(m string) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(m+"\n")
	reader.ReadString('\n')

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