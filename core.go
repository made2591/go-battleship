package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Game struct {
	sea_a Sea
	sea_b Sea
	moves []Coordinates
}

type Sea struct {
	dimension int
	grid      [][]int
	ships     []Ship
	moves     []Coordinates
}

type Ship struct {
	dimension int
	positions []Coordinates
	available []Coordinates
}

type Coordinates struct {
	abscissa int
	ordinate int
}

func random(min, max int) int {
	max = max + 1
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func search(a int, b []int) bool {
	for _, v := range b {
		if v == a {
			return true
		}
	}
	return false
}

// PrepareSea return a Sea struct given dimension n
func PrepareSea(n int, s int) (sea Sea) {

	ss := make([]Ship, s)

	for i := 0; i < s; i++ {

		st := PrepareShip(i+1, n)
		if !checkCollisions(st, ss) {
			ss = append(ss, PrepareShip(i+1, n))
		} else {
			i--
		}

	}

	sea = Sea{dimension: n, grid: make([][]int, n), ships: ss, moves: []Coordinates{}}

	return

}

func PrepareShip(n int, m int) (s Ship) {

	h := random(0, 1) == 1

	p := make([]Coordinates, n)

	a := make([]Coordinates, n)

	if n == 1 {

		x := random(1, m)
		y := random(1, m)
		p[0] = Coordinates{abscissa: x, ordinate: y}
		a[0] = Coordinates{abscissa: x, ordinate: y}

	} else {

		x := random(1, m-n)
		y := random(1, m)

		for t := 0; t < n; t++ {

			if h {

				p[t] = Coordinates{abscissa: x + t, ordinate: y}
				a[t] = Coordinates{abscissa: x + t, ordinate: y}

			} else {

				p[t] = Coordinates{abscissa: y, ordinate: x + t}
				a[t] = Coordinates{abscissa: y, ordinate: x + t}

			}
		}

	}

	s = Ship{dimension: n, positions: p, available: a}
	return

}

func checkCollisions(a Ship, b []Ship) bool {

	for _, sb := range b {
		if checkCollision(a, sb) {
			return true
		}
	}
	return false

}

func checkCollision(a Ship, b Ship) bool {

	for _, av := range a.positions {
		for _, bv := range b.positions {
			if av.abscissa == bv.abscissa && av.ordinate == bv.ordinate {
				return true
			}
		}
	}
	return false

}

func checkPosition(x int, y int, s Sea) bool {

	for _, sv := range s.ships {
		for _, cv := range sv.positions {
			if x == cv.abscissa && y == cv.ordinate {
				return true
			}
		}
	}
	return false

}

func CoordinatesPrettyInfo(c Coordinates) (cs string) {

	cs = "(" + strconv.Itoa(c.abscissa) + "; " + strconv.Itoa(c.ordinate) + ")"
	return

}

func ShipPrettyInfo(s Ship) (ss string) {

	ss = "\tShip dimensions: " + strconv.Itoa(s.dimension) + "\n\t\t["
	for _, pv := range s.positions {
		ss += CoordinatesPrettyInfo(pv) + " "
	}
	ss += "]"
	return

}

func SeaPrettyInfo(s Sea) (ss string) {

	ss = "Sea dimensions: " + strconv.Itoa(s.dimension) + "\n"
	for _, sv := range s.ships {
		if sv.dimension != 0 {
			ss += ShipPrettyInfo(sv) + "\n"
		}
	}
	return

}

func StringfySea(s Sea) (ss string) {

	ss = "|"
	for r := 0; r < s.dimension-1; r++ {
		ss += "-----"
	}
	ss += "----|\n"

	for r := 0; r < s.dimension; r++ {

		ss += "|"
		for c := 0; c < s.dimension; c++ {
			if checkPosition(r+1, c+1, s) {
				ss += " ** |"
			} else {
				ss += "    |"
			}
		}
		ss += "\n"
		ss += "|"
		for c := 0; c < s.dimension-1; c++ {
			ss += "-----"
		}
		ss += "----|\n"

	}

	return ss

}

func main() {

	s := PrepareSea(10, 5)
	fmt.Println(SeaPrettyInfo(s))
	fmt.Println(StringfySea(s))

}
