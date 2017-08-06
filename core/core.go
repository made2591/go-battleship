package core

import (

	"fmt"
	"math"
	"strconv"
	"encoding/json"

	util "github.com/made2591/go-battleship/util"

	"net/http"
)

// Constants for default and game config
const (

	PROTOCOL = "http"
	HOST_NAME = "localhost"
	HOST_PORT = "8080"

	START_REQUEST = "/start"
	SHOT_REQUEST  = "/shot"
	EXIT_REQUEST  = "/exit"

	SIMULATION_THINKING_TIME = 2000 //milliseconds
	PRINT_CALLERSEA_MODE	 = 2
	PRINT_OPPONENTSEA_MODE   = 1

	GUN_SHOT_COST 		= 10
	FIRE_SHOT_COST		= 100
	NUCLEAR_SHOT_COST	= 1000

	PC_NAME  = "HAL"
	PC_SHOTS = 9999
	PC_SHIPS = 5
	PC_GRID  = 10
	PC_MODE  = 0

	GAME_GRID_BORDER   = "|"
	STR_SHIP_OK        = "00"
	STR_SHIP_STRICKEN  = "++"
	STR_SHIP_BURNING   = "§§"
	STR_SHIP_DESTROYED = "XX"
	STR_SEA_OK         = "  "
	STR_SEA_STRICKEN   = "~~"
	STR_STATUS_ERROR   = "??"

	STATUS_SHIP_OK 			= iota
	STATUS_SHIP_STRICKEN
	STATUS_SHIP_BURNING
	STATUS_SHIP_DESTROYED
	STATUS_SEA_OK
	STATUS_SEA_STRICKEN

)

// ###########################################################################################################
// ################################################### GAME ##################################################
// ###########################################################################################################

type Game struct {
	FirstPlayer  Player `json:"FirstPlayer"`
	SecondPlayer Player `json:"SecondPlayer"`
}

type Player struct {
	Name     string        `json:"Name"`
	Sea      Sea           `json:"Sea"`
	GunShot  int           `json:"GunShot"`
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
	Status int `json:"Status"`
}

// ###########################################################################################################
// ############################################### CONSTRUCTORS ##############################################
// ###########################################################################################################

// Game init
//	[d:int]		grid dimension		[m:int]	game mode 0 1:PC 1:1
//	[nf:int]	first player name	[sf:int]	first player number of ship	[gf:int]	first player gun shots
//	[ns:string]	first player name	[ss:int]	first player number of ship	[gs:int]	first player gun shots
func PrepareGame(d int, m int, nf string, sf int, gf int, ns string, ss int, gs int) (g Game) {

	// create First Player
	fp := Player{Name: nf, GunShot: gf, Sea: PrepareSea(d, sf)}

	// create First Player
	sp := Player{}
	if m == PC_MODE {
		sp = Player{Name: PC_NAME, GunShot: PC_SHOTS, Sea: PrepareSea(PC_GRID, PC_SHIPS)}
	} else {
		sp = Player{Name: ns, GunShot: gs, Sea: PrepareSea(d, ss)}
	}

	g = Game{FirstPlayer: fp, SecondPlayer: sp}

	return

}

// Sea init
//	[d:int]		grid dimension		[s:int]	number of ship
// TODO CREATE MORE EFFICIENT ALGORITHM FOR RANDOM GEN OF SHIPS
func PrepareSea(d int, n int) (s Sea) {

	// prepare array of Ship
	ss := make([]Ship, n)

	// create n Ship with incremental dimension
	for i := 0; i < n; i++ {

		// create Ship
		st := PrepareShip(i+1, d)

		// if it doesn't collide with other ships
		if !CheckCollisions(&st, ss) {
			// add to Sea
			ss[i] = PrepareShip(i+1, d)
		} else {
			// retry
			i--
		}

	}

	// create Sea
	s = Sea{Dimension: d, Ships: ss}
	return

}

// Ship init
//	[sd:int]	ship dimension		[gd:int]		grid dimension
func PrepareShip(sd int, gd int) (s Ship) {

	// choose if horizontal
	h := util.Random(0, 1) == 1

	// create Ship coordinates
	p := make([]Coordinates, sd)

	// if Ship dimension is 1
	if sd == 1 {

		// create Random coordinate
		x := util.Random(1, gd)
		y := util.Random(1, gd)
		// add unique Coordinate
		p[0] = Coordinates{Abscissa: x, Ordinate: y, Status: STATUS_SHIP_OK}

	} else {

		// create x coordinate no more than grid dimension
		x := util.Random(1, gd-sd)
		// create y coordinate no more than grid dimension
		y := util.Random(1, gd)

		// create Coordinates
		for t := 0; t < sd; t++ {

			// offset on x
			if h {
				p[t] = Coordinates{Abscissa: x + t, Ordinate: y, Status: STATUS_SHIP_OK}
			// offset on y
			} else {
				p[t] = Coordinates{Abscissa: y, Ordinate: x + t, Status: STATUS_SHIP_OK}
			}
		}

	}

	// create Ship
	s = Ship{Dimension: sd, Positions: p}
	return

}

// ###########################################################################################################
// ########################################### GAME LOGIC METHODS ############################################
// ###########################################################################################################

// GameDecoder decode game in HTTP request
func GameDecoder(r *http.Request) (g Game) {

	// create Game
	g = Game{}

	// decode game
	d := json.NewDecoder(r.Body)
	err := d.Decode(&g)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return

}

// SwitchPointOfView switch point of view
//	[g:*Game]	Game pointer
func SwitchPointOfView(g *Game) {

	t := g.FirstPlayer
	g.FirstPlayer = g.SecondPlayer
	g.SecondPlayer = t

}

// CheckCollisions check if a collides with at least one of b ships
//	[a:*Ship]	ship pointer		[b:array of Ships]		array of Ships
func CheckCollisions(a *Ship, b []Ship) bool {

	for _, sb := range b {
		if CheckCollision(a, &sb) {
			return true
		}
	}
	return false

}

// CheckCollisions check if a collides with b
//	[a:*Ship]	a ship pointer		[b:*Ship]		b Ship pointer
func CheckCollision(a *Ship, b *Ship) bool {

	for _, av := range a.Positions {
		for _, bv := range b.Positions {
			if av.Abscissa == bv.Abscissa && av.Ordinate == bv.Ordinate {
				return true
			}
		}
	}
	return false

}

// CheckShipPosition check if in p coordinates in given Sea there's a Ship
//	[p:*Coordinates]	Coordinate point pointer		[s:*Sea]		b Sea pointer
//	[return]	bool (collision), ship index, coordinate index
func CheckShipPosition(p *Coordinates, s *Sea) (bool, int, int) {

	// for each Ships
	for si, sv := range s.Ships {
		// for each Positions occuped by
		for ci, cv := range sv.Positions {
			// if coordinates == positions
			if p.Abscissa == cv.Abscissa && p.Ordinate == cv.Ordinate {
				// return true, ship index, positions index in ship struct
				return true, si, ci
			}
		}
	}
	return false, -1, -1

}

// CheckSufferedMoves check p coordinates in given Sea's Player
//	[p:*Coordinates]	Coordinate point pointer		[pp:*Player]		b Player pointer
//	[return]	bool (collision), ship index, coordinate index
func CheckSufferedMoves(p *Coordinates, pp *Player) (bool, int) {

	// for each suffered
	for pi, pv := range pp.Suffered {
		// if coordinates == positions
		if p.Abscissa == pv.Abscissa && p.Ordinate == pv.Ordinate {
			// return true, positions index in Player Suffered Moves
			return true, pi
		}
	}
	return false, -1

}

// GunShot from p Player to t Player in p Coordinates
//	[f:*Player]			from Player	pointer		//	[t:*Player]			to Player pointer
//	[p:*Coordinates]	Coordinate point pointer
func GunShot(f *Player, t *Player, p *Coordinates) {

	// check if f player has sufficient coins to shot
	if f.GunShot > GUN_SHOT_COST {
		f.GunShot = f.GunShot - GUN_SHOT_COST
	}

	// check if f player hit t player in position p
	rs, si, ci := CheckShipPosition(p, &t.Sea)

	np := Coordinates{}

	// if Ship hit
	if rs {
		// TODO CHECK IF SHIP IS DESTROYED
		// t player ship stricken
		t.Sea.Ships[si].Positions[ci].Status = STATUS_SHIP_STRICKEN
		np = Coordinates{Abscissa: p.Abscissa, Ordinate: p.Ordinate, Status: STATUS_SHIP_STRICKEN}
	} else {
		p.Status = STATUS_SEA_STRICKEN
		np = Coordinates{Abscissa: p.Abscissa, Ordinate: p.Ordinate, Status: STATUS_SEA_STRICKEN}
	}

	// add to suffered
	t.Suffered = append(t.Suffered, np)

}

// ###########################################################################################################
// ########################################## STRUCTS STRINGIFIER ############################################
// ###########################################################################################################

// StatusToString print Status of Coordinates
func StatusToString(s int) (string) {

	// check Ship status in specific position
	switch s {

		// ship && sea status
		case STATUS_SHIP_STRICKEN:
			return STR_SHIP_STRICKEN
		case STATUS_SHIP_BURNING:
			return STR_SHIP_BURNING
		case STATUS_SHIP_DESTROYED:
			return STR_SHIP_DESTROYED
		case STATUS_SHIP_OK:
			return STR_SHIP_OK
		case STATUS_SEA_STRICKEN:
			return STR_SEA_STRICKEN
		case STATUS_SEA_OK:
			return STR_SEA_OK

	}
	return STR_STATUS_ERROR

}

// SeaToString print Player Sea
//	[p:*Player]	Player
//	[return]	string
func SeaToString(p *Player, h int) (ss string) {

	// create first separation line
	ss = GAME_GRID_BORDER
	for r := 0; r < p.Sea.Dimension-1; r++ {
		ss += "-----"
	}
	ss += "----"+GAME_GRID_BORDER+"\n"

	// for each row
	for r := 0; r < p.Sea.Dimension; r++ {

		// start with grid border
		ss += GAME_GRID_BORDER

		// for each column
		for c := 0; c < p.Sea.Dimension; c++ {

			// if we are drawing caller's Sea
			if h == PRINT_CALLERSEA_MODE {

				// check ShipPosition in Sea
				rp, si, ci := CheckShipPosition(&Coordinates{Abscissa: r + 1, Ordinate: c + 1}, &p.Sea)

				// if there's a Sea in position
				if rp {

					// add correct status representation
					ss += " " + StatusToString(p.Sea.Ships[si].Positions[ci].Status) + " " + GAME_GRID_BORDER

				} else {

					// check SufferedMoves in Sea
					pp, pi := CheckSufferedMoves(&Coordinates{Abscissa: r + 1, Ordinate: c + 1}, p)

					// if opponent shot in the cell
					if pp {

						// add correct status representation
						ss += " " + StatusToString(p.Suffered[pi].Status) + " " + GAME_GRID_BORDER

					} else {

						ss += " " + STR_SEA_OK + " " + GAME_GRID_BORDER

					}

				}

			}

			// if we are drawing caller's Sea
			if h == PRINT_OPPONENTSEA_MODE {

				// check SufferedMoves in Sea
				pp, pi := CheckSufferedMoves(&Coordinates{Abscissa: r + 1, Ordinate: c + 1}, p)

				// if opponent shot in the cell
				if pp {

					// add correct status representation
					ss += " " + StatusToString(p.Suffered[pi].Status) + " " + GAME_GRID_BORDER

				} else {

					ss += " " + STR_SEA_OK + " " + GAME_GRID_BORDER

				}

			}

		}

		// create separation line
		ss += "\n"+GAME_GRID_BORDER
		for c := 0; c < p.Sea.Dimension-1; c++ {
			ss += "-----"
		}
		ss += "----"+GAME_GRID_BORDER+"\n"

	}

	return ss

}

// PrettyPrintCoordinatesInfo return String rappresentation of Coordinates
//	[p:*Coordinates]	Coordinates point pointer
//	[return]	string
func PrettyPrintCoordinatesInfo(p *Coordinates) (ps string) {

	ps = "[" + strconv.Itoa(p.Abscissa) + "; " + strconv.Itoa(p.Ordinate) + "]"
	return

}

func LastShotInToString(g *Game) (s string) {
	s = ">>> shot out in coordinates "+PrettyPrintCoordinatesInfo(&g.FirstPlayer.Suffered[len(g.FirstPlayer.Suffered)-1])+"\n"
	return
}

func LastShotOutToString(g *Game) (s string) {
	s = ">>> shot in  in coordinates "+PrettyPrintCoordinatesInfo(&g.SecondPlayer.Suffered[len(g.SecondPlayer.Suffered)-1])+"\n"
	return
}

func ShotHistory(g *Game) (s string) {

	sm := int(math.Max(float64(len(g.FirstPlayer.Suffered)), float64(len(g.SecondPlayer.Suffered))))
	sp := len(g.FirstPlayer.Suffered) > len(g.SecondPlayer.Suffered)

	for i := 0; i < sm; i++ {
		if sp && g.SecondPlayer.Suffered != nil && i < len(g.SecondPlayer.Suffered) {
			s += ">>> shot out coordinates "+PrettyPrintCoordinatesInfo(&g.SecondPlayer.Suffered[i])+"\n"
		}
		if g.FirstPlayer.Suffered != nil && i < len(g.FirstPlayer.Suffered) {
			s += ">>> shot in  coordinates " + PrettyPrintCoordinatesInfo(&g.FirstPlayer.Suffered[i]) + "\n"
		}
		if !sp && g.SecondPlayer.Suffered != nil && i < len(g.SecondPlayer.Suffered) {
			s += ">>> shot out coordinates " + PrettyPrintCoordinatesInfo(&g.SecondPlayer.Suffered[i]) + "\n"
		}
	}
	return

}

// PrettyPrintShipInfo return Ship string info
//	[s:*Ship]	Ship point pointer
//	[return]	string
func PrettyPrintShipInfo(s *Ship) (ss string) {

	ss = "\tShip dimensions: " + strconv.Itoa(s.Dimension) + "\n\t\t["
	for _, pv := range s.Positions {
		ss += PrettyPrintCoordinatesInfo(&pv) + " "
	}
	ss += "]"
	return

}

// PrettyPrintGame from p Player to t Player in p Coordinates
//	[g:*Game]			Game pointer		[m:int]	game mode 0 1:PC 1:1
func PrettyPrintGame(g *Game) (gs string) {

	// clean tty screen
	util.CleanScreen()

	gs += ">>> "+g.FirstPlayer.Name+"'s sea\n"
	gs += SeaToString(&g.FirstPlayer, PRINT_CALLERSEA_MODE)
	gs += ">>> "+g.SecondPlayer.Name+"'s sea\n"
	gs += SeaToString(&g.SecondPlayer, PRINT_OPPONENTSEA_MODE)
	gs += "\n"+ShotHistory(g)

	return

}

// PrettyPrintSeaInfo return Sea string info
//	[s:*Ship]	Ship point pointer
//	[return]	string
func PrettyPrintSeaInfo(s *Sea) (ss string) {

	ss = "Sea dimensions: " + strconv.Itoa(s.Dimension) + "\n"
	for _, sv := range s.Ships {
		if sv.Dimension != 0 {
			ss += PrettyPrintShipInfo(&sv) + "\n"
		}
	}
	return

}

func DebugGame(g *Game) string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// ###########################################################################################################
// ######################################### TEST METHODS STRINGIFIER ########################################
// ###########################################################################################################

func main() {

	s := PrepareSea(10, 5)
	fmt.Println(PrettyPrintSeaInfo(&s))
	g := PrepareGame(10, 0, "Matteo", 5, 9999, "HAL", 5, 9999)
	fmt.Println(SeaToString(&g.FirstPlayer, -1))
	fmt.Println(SeaToString(&g.SecondPlayer, -1))
	GunShot(&g.FirstPlayer, &g.SecondPlayer, &g.SecondPlayer.Sea.Ships[0].Positions[0])
	fmt.Println(SeaToString(&g.SecondPlayer, -1))

}
