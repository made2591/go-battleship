package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Game struct {
	Sea_a Sea
	Sea_b Sea
	Moves_vs_a []Coordinates
	Moves_vs_b []Coordinates
}

type Sea struct {
	Dimension int
	Grid      [][]int
	Ships     []Ship
	Moves     []Coordinates
}

type Ship struct {
	Dimension int
	Positions []Coordinates
	Available []Coordinates
}

type Coordinates struct {
	Abscissa int
	Ordinate int
	Status   int
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

// PrepareSea return a Sea struct given Dimension: n
func PrepareSea(n int, s int) (sea Sea) {

	ss := make([]Ship, s)

	for i := 0; i < s; i++ {

		st := PrepareShip(i+1, n)
		if !CheckCollisions(st, ss) {
			ss = append(ss, PrepareShip(i+1, n))
		} else {
			i--
		}

	}

	sea = Sea{Dimension: n, Grid: make([][]int, n), Ships: ss, Moves: []Coordinates{}}

	return

}

func PrepareShip(n int, m int) (s Ship) {

	h := random(0, 1) == 1

	p := make([]Coordinates, n)

	a := make([]Coordinates, n)

	if n == 1 {

		x := random(1, m)
		y := random(1, m)
		p[0] = Coordinates{Abscissa: x, Ordinate: y}
		a[0] = Coordinates{Abscissa: x, Ordinate: y}

	} else {

		x := random(1, m-n)
		y := random(1, m)

		for t := 0; t < n; t++ {

			if h {

				p[t] = Coordinates{Abscissa: x + t, Ordinate: y}
				a[t] = Coordinates{Abscissa: x + t, Ordinate: y}

			} else {

				p[t] = Coordinates{Abscissa: y, Ordinate: x + t}
				a[t] = Coordinates{Abscissa: y, Ordinate: x + t}

			}
		}

	}

	s = Ship{Dimension: n, Positions: p, Available: a}
	return

}

func CheckCollisions(a Ship, b []Ship) bool {

	for _, sb := range b {
		if CheckCollision(a, sb) {
			return true
		}
	}
	return false

}

func CheckCollision(a Ship, b Ship) bool {

	for _, av := range a.Positions {
		for _, bv := range b.Positions {
			if av.Abscissa == bv.Abscissa && av.Ordinate == bv.Ordinate {
				return true
			}
		}
	}
	return false

}

func CheckPosition(x int, y int, s Sea) bool {

	for _, sv := range s.Ships {
		for _, cv := range sv.Positions {
			if x == cv.Abscissa && y == cv.Ordinate {
				return true
			}
		}
	}
	return false

}

func CoordinatesPrettyInfo(c Coordinates) (cs string) {

	cs = "(" + strconv.Itoa(c.Abscissa) + "; " + strconv.Itoa(c.Ordinate) + ")"
	return

}

func ShipPrettyInfo(s Ship) (ss string) {

	ss = "\tShip dimensions: " + strconv.Itoa(s.Dimension) + "\n\t\t["
	for _, pv := range s.Positions {
		ss += CoordinatesPrettyInfo(pv) + " "
	}
	ss += "]"
	return

}

func SeaPrettyInfo(s Sea) (ss string) {

	ss = "Sea dimensions: " + strconv.Itoa(s.Dimension) + "\n"
	for _, sv := range s.Ships {
		if sv.Dimension != 0 {
			ss += ShipPrettyInfo(sv) + "\n"
		}
	}
	return

}

func StringfySea(s Sea) (ss string) {

	ss = "|"
	for r := 0; r < s.Dimension-1; r++ {
		ss += "-----"
	}
	ss += "----|\n"

	for r := 0; r < s.Dimension; r++ {

		ss += "|"
		for c := 0; c < s.Dimension; c++ {
			if CheckPosition(r+1, c+1, s) {
				ss += " ** |"
			} else {
				ss += "    |"
			}
		}
		ss += "\n"
		ss += "|"
		for c := 0; c < s.Dimension-1; c++ {
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
