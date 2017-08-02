package core

import (
	"fmt"
	//	"math/rand"
	"bufio"
	"encoding/json"
	util "github.com/made2591/go-battleship/util"
	"net/http"
	"os"
	"strconv"
	//	"time"
)

const (
	PC_NAME  = "HAL"
	PC_SHOTS = 9999

	STATUS_OK_BASE      = "00"
	STATUS_DESTROY_BASE = "XX"
	STATUS_FIRE_BASE    = "/\\"
	SEA_BASE            = "  "
	SEA_SHOT_BASE       = "~~"

	STATUS_OK      = iota // 0
	STATUS_DESTROY        // 1
	STATUS_FIRE           // 2
	SEA                   // 3
	SEA_SHOT              // 4

)

type Game struct {
	FirstPlayer  Player `json:"FirstPlayer"`
	SecondPlayer Player `json:"SecondPlayer"`
}

type Player struct {
	Name     string        `json:"Name"`
	Sea      Sea           `json:"Sea"`
	GunShot  int           `json:"GunShot"`
	Moves    []Coordinates `json:"Moves"`
	Suffered []Coordinates `json:"Suffered"`
}

type Sea struct {
	Dimension int    `json:"Dimension"`
	Ships     []Ship `json:"Ships"`
}

type Ship struct {
	Dimension int           `json:"Dimension"`
	Positions []Coordinates `json:"Positions"`
}

type Coordinates struct {
	Abscissa int `json:"Abscissa"`
	Ordinate int `json:"Ordinate"`
	// 0 ok 1 hit
	Status int `json:"Status"`
}

func PrepareGame(d int, m int, na string, sa int, ga int, nb string, sb int, gb int) (g Game) {

	pf := Player{}
	ps := Player{}

	if m == 0 {
		pf = Player{Name: na, GunShot: ga, Sea: PrepareSea(d, sa)}
		ps = Player{Name: PC_NAME, GunShot: PC_SHOTS, Sea: PrepareSea(d, sb)}
	}
	g = Game{FirstPlayer: pf, SecondPlayer: ps}
	return

}

// PrepareSea return a Sea struct given Dimension: n
func PreparePlayer(n string, g int, s Sea) (p Player) {

	p = Player{Name: n, GunShot: g, Sea: s, Moves: []Coordinates{}, Suffered: []Coordinates{}}

	return

}

// PrepareSea return a Sea struct given Dimension: n
func PrepareSea(n int, s int) (sea Sea) {

	ss := make([]Ship, s)

	for i := 0; i < s; i++ {

		st := PrepareShip(i+1, n)
		if !CheckCollisions(st, ss) {
			ss[i] = PrepareShip(i+1, n)
		} else {
			i--
		}

	}

	sea = Sea{Dimension: n, Ships: ss}

	return

}

func PrepareShip(n int, m int) (s Ship) {

	h := util.Random(0, 1) == 1

	p := make([]Coordinates, n)

	if n == 1 {

		x := util.Random(1, m)
		y := util.Random(1, m)
		p[0] = Coordinates{Abscissa: x, Ordinate: y}

	} else {

		x := util.Random(1, m-n)
		y := util.Random(1, m)

		for t := 0; t < n; t++ {

			if h {

				p[t] = Coordinates{Abscissa: x + t, Ordinate: y}

			} else {

				p[t] = Coordinates{Abscissa: y, Ordinate: x + t}

			}
		}

	}

	s = Ship{Dimension: n, Positions: p}
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

func CheckShot(p Coordinates, s Sea) (bool, int, int) {

	for si, sv := range s.Ships {
		for ci, cv := range sv.Positions {
			if p.Abscissa == cv.Abscissa && p.Ordinate == cv.Ordinate {
				return true, si, ci
			}
		}
	}
	return false, -1, -1

}

func CheckShipPosition(x int, y int, s Sea) (bool, int, int) {

	for si, sv := range s.Ships {
		for ci, cv := range sv.Positions {
			if x == cv.Abscissa && y == cv.Ordinate {
				return true, si, ci
			}
		}
	}
	return false, -1, -1

}

func CheckSeaPosition(x int, y int, p Player) (bool, int) {

	for pi, pv := range p.Suffered {
		if x == pv.Abscissa && y == pv.Ordinate {
			return true, pi
		}
	}
	return false, -1

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

func StringfySea(p Player) (ss string) {

	ss = "|"
	for r := 0; r < p.Sea.Dimension-1; r++ {
		ss += "-----"
	}
	ss += "----|\n"

	s := p.Sea

	for r := 0; r < s.Dimension; r++ {

		ss += "|"
		for c := 0; c < s.Dimension; c++ {
			rp, si, ci := CheckShipPosition(r+1, c+1, s)
			if rp {
				switch s.Ships[si].Positions[ci].Status {
				case STATUS_DESTROY:
					ss += " " + STATUS_DESTROY_BASE + " |"
				case STATUS_FIRE:
					ss += " " + STATUS_FIRE_BASE + " |"
				default:
					ss += " " + STATUS_OK_BASE + " |"
				}
			} else {
				pp, pi := CheckSeaPosition(r+1, c+1, p)
				if pp {
					switch p.Suffered[pi].Status {
					case SEA_SHOT:
						ss += " " + SEA_SHOT_BASE + " |"
					default:
						ss += " " + SEA_BASE + " |"
					}
				} else {
					ss += " " + SEA_BASE + " |"
				}
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

func (g Game) GunShot(f *Player, t *Player, p Coordinates) {

	if f.GunShot > 0 {
		f.GunShot--
	}
	rs, si, ci := CheckShot(p, t.Sea)
	if rs {
		f.Moves = append(f.Moves, p)
		f.Sea.Ships[si].Positions[ci].Status = STATUS_DESTROY
		p.Status = SEA_SHOT
	} else {
		p.Status = SEA
	}
	t.Suffered = append(t.Suffered, p)

}

func NetPrintGame(g *Game, m int) {

	util.CleanScreen()
	if m == 0 {
		fmt.Printf(">>> %s's sea\n", g.FirstPlayer.Name)
		fmt.Println(StringfySea(g.FirstPlayer))
		//fmt.Printf(">>> %s's sea\n", g.SecondPlayer.Name)
		//fmt.Println(StringfySea(g.SecondPlayer))
	}
	if m == 1 {
		fmt.Printf(">>> %s's sea\n", g.SecondPlayer.Name)
		fmt.Println(StringfySea(g.SecondPlayer))
		//fmt.Printf(">>> %s's sea\n", g.FirstPlayer.Name)
		//fmt.Println(StringfySea(g.FirstPlayer))
	}

}

func ServerGunShot(w http.ResponseWriter, r *http.Request) {
	reader := bufio.NewReader(os.Stdin)

	d := json.NewDecoder(r.Body)
	g := Game{}

	err := d.Decode(&g)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	bbb, _ := json.Marshal(g)
	fmt.Println(string(bbb))

	fmt.Printf(">>> press ENTER to go on...\n")
	reader.ReadString('\n')

	fmt.Printf(">>> shot received in coordinates [%d, %d]\n",
		g.SecondPlayer.Suffered[len(g.SecondPlayer.Suffered)-1].Abscissa,
		g.SecondPlayer.Suffered[len(g.SecondPlayer.Suffered)-1].Abscissa)
	fmt.Printf(">>> press ENTER to go on...\n")
	reader.ReadString('\n')

	s := util.Random(0, len(g.SecondPlayer.Sea.Ships)-1)
	p := util.Random(0, len(g.SecondPlayer.Sea.Ships[s].Positions)-1)
	g.GunShot(&g.SecondPlayer, &g.FirstPlayer, g.SecondPlayer.Sea.Ships[s].Positions[p])

	fmt.Printf(">>> gun shot coordinates [%d, %d]\n",
		g.SecondPlayer.Sea.Ships[s].Positions[p].Abscissa,
		g.SecondPlayer.Sea.Ships[s].Positions[p].Ordinate)
	fmt.Printf(">>> press ENTER to go on...\n")
	reader.ReadString('\n')

	NetPrintGame(&g, 1)

	fmt.Printf(">>> press ENTER to go on...\n")
	reader.ReadString('\n')

	NetPrintGame(&g, 1)

	json.NewEncoder(w).Encode(g)
}

func main() {

	//s := PrepareSea(10, 5)
	//fmt.Println(SeaPrettyInfo(s))
	//fmt.Println(StringfySea(s))
	g := PrepareGame(10, 0, "Matteo", 5, 9999, "HAL", 5, 9999)
	fmt.Println(StringfySea(g.FirstPlayer))
	fmt.Println(StringfySea(g.SecondPlayer))
	//fmt.Println(g.SecondPlayer.Sea.Ships)
	g.GunShot(&g.FirstPlayer, &g.SecondPlayer, g.SecondPlayer.Sea.Ships[0].Positions[0])
	fmt.Println(StringfySea(g.SecondPlayer))

}
