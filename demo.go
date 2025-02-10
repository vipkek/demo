package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
	"time"
)

const RTP = 0.97
const MinCrashValue = 1.01
const houseEdge = 0.03

var globalCasinoBank = 0.0

type PlayerBet struct {
	PlayerID    string
	BetAmount   float64
	AutoCashOut float64 // value 0 for manual cashout
	WonAmount   float64
	CashedOut   bool
}

type GameRound struct {
	ID              int
	Players         []*PlayerBet
	CrashMultiplier float64
	CashoutChannel  chan string // for handling manual cashout
}

func generateCrashMultiplier() float64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}
	randomNumber := float64(n.Int64()) / 1000000.0
	fmt.Println("randomNumber", randomNumber)

	const adjustedRTP = (1.0 - RTP) / RTP //1.0 - (1.0 - RTP - houseEdge)
	fmt.Println("adjustedRTP", adjustedRTP)
	crashMultiplier := 1.0 / (1.0 - adjustedRTP*randomNumber) // balance RTP

	fmt.Println("crashMultiplier", crashMultiplier)
	return math.Max(crashMultiplier, MinCrashValue) // crash rate more than 1.01x
}
func totalBet(players []*PlayerBet) float64 {
	totalBet := 0.0
	for _, bet := range players {
		totalBet += bet.BetAmount
	}
	return totalBet
}

func generateCrashMultiplierBasedOnBets(players []*PlayerBet) float64 {
	fmt.Println("totalBet: ", totalBet(players))

	riskFactor := 1.0 - math.Min(totalBet(players)/1000.0, 0.1) // Caps at 10% adjustment

	fmt.Println("risk factor: ", riskFactor)
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}
	randomNumber := float64(n.Int64()) / 1000000.0
	adjustedFactor := 1.0 - (1.0 - RTP - houseEdge)

	crashMultiplier := (1.0 / (1.0 - adjustedFactor*randomNumber)) * riskFactor
	return math.Max(crashMultiplier, MinCrashValue)
}

func newGameRound(id int, players []*PlayerBet) *GameRound {
	return &GameRound{
		ID:             id,
		Players:        players,
		CashoutChannel: make(chan string),
	}
}

func (g *GameRound) StartRound() {
	fmt.Printf("\nRound %d started! Waiting for bets...\n", g.ID)

	time.Sleep(2 * time.Second)

	g.CrashMultiplier = generateCrashMultiplierBasedOnBets(g.Players)
	fmt.Printf("Round %d started! Crash at %.2fx\n", g.ID, g.CrashMultiplier)

	gameValidator(g.ID, g.CrashMultiplier, g.Players)

}

func (g *GameRound) listenForCashouts() {
	for {
		var playerID string
		fmt.Print("Enter Player ID or skip (press Enter)")
		fmt.Scanln(&playerID)

		if playerID != "" {
			g.CashoutChannel <- playerID
		}
	}
}

func main() {
	// testPlayers := []*PlayerBet{
	// 	{"Bogdan", 10, 2.5, 0, false},
	// 	{"Slavik", 50, 3.15, 0, false},
	// 	{"Olya", 20, 5.5, 0, false},
	// 	{"Katya", 12, 1.4, 0, false}, // manual cashout
	// 	{"Lida", 25, 6.4, 0, false},
	// 	{"Max", 13, 1.25, 0, false},
	// 	{"Alex", 44, 2.5, 0, false},
	// 	{"Grisha", 78, 1.2, 0, false}, // manual cashout
	// }

	generatedPlayers := testPlayersGenerator(15, 0)

	for roundID := 1; roundID <= 5; roundID++ {
		newGame := newGameRound(roundID, generatedPlayers)
		go newGame.StartRound()
		// fmt.Println("Waiting 10 sec for new game...")
		time.Sleep(10 * time.Second)
	}

	fmt.Println("TOTAL CASINO BANK: ", globalCasinoBank)

}

func gameValidator(ID int, crashMultiplier float64, players []*PlayerBet) {
	totalBet := totalBet(players)
	totalWin := 0.0
	fmt.Printf("GAME VALIDATION %d -- CRASH AT %.2fx -- TOTAL BETS: %.2f \n", ID, crashMultiplier, totalBet)
	for _, player := range players {
		if player.AutoCashOut <= crashMultiplier {
			totalWin += player.BetAmount * player.AutoCashOut
			fmt.Printf("WINNER: %s -- bet: %.3f, win: %.3f \n", player.PlayerID, player.BetAmount, player.AutoCashOut*player.BetAmount)
		}
	}

	fmt.Printf("TOTAL PLAYERS WIN: %f, TOTAL PLAYERS LOST: %.3f \n", totalWin, totalBet-totalWin)
	globalCasinoBank += totalBet - totalWin
}

func testPlayersGenerator(amount int, numOfAutoCashout int) []*PlayerBet {

	players := []*PlayerBet{
		{"BogdanFreeman777", 10, 2.5, 0, false},
	}

	for i := 1; i <= amount; i++ {
		player := &PlayerBet{
			generateUsername(),
			randomFloatInRange(10.00, 1000.00),
			randomFloatInRange(1.01, 10.00),
			0,
			false,
		}
		players = append(players, player)
		fmt.Printf("Created player: %s -- amount: %.2f -- autocashout: %.2f \n", player.PlayerID, player.BetAmount, player.AutoCashOut)
	}
	return players
}

func generateUsername() string {
	names := []string{"Stive", "Brad", "Cloue", "Harry", "Max", "Clay", "Andrey", "Loly", "Faust", "George"}
	surnames := []string{"Trueman", "Musk", "Jackson", "Sinatra", "Perry", "Baseman", "Gates", "Lenon", "Joe", "Clarson"}

	name := names[mrand.Intn(len(names))]
	surname := surnames[mrand.Intn(len(surnames))]
	num := mrand.Intn(1000)
	return fmt.Sprintf("%s%s%d", name, surname, num)
}

func randomFloatInRange(min, max float64) float64 {
	return min + mrand.Float64()*(max-min)
}
