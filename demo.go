package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

const RTP = 0.97
const MinCrashValue = 1.01
const houseEdge = 0.003

func generateCrashMultiplier() float64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		panic(err)
	}
	randomNumber := float64(n.Int64()) / 1000000.0
	fmt.Println("randomNumber", randomNumber)

	const adjustedRTP = RTP - houseEdge
	crashMultiplier := 1.0 / (1.0 - (1.0-adjustedRTP)*randomNumber) // balance RTP

	fmt.Println("crashMultiplier", crashMultiplier)
	return math.Max(crashMultiplier, MinCrashValue) // crash rate more than 1.01x
}

func main() {
	fmt.Print("final: ", generateCrashMultiplier())
}
