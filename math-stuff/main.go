package main

import (
	"flag"
	"fmt"
	"math/big"
)

// Computes square roots at high precision fixed-point using
// Newton's method.  -dga 8 Nov 2011
// From http://www.angio.net/pi/pi-programs.html

func sqrt(n, unity *big.Int) *big.Int {
	// Initial guess = 2^(log_2(n)/2)
	guess := big.NewInt(2)
	guess.Exp(guess, big.NewInt(int64(n.BitLen()/2)), nil)

	n_big := big.NewInt(0).Mul(n, unity)
	// Now refine using Newton's method
	prev_guess := big.NewInt(0)
	for {
		prev_guess.Set(guess)
		guess.Add(guess, big.NewInt(0).Div(n_big, guess))
		guess.Div(guess, big.NewInt(2))
		if guess.Cmp(prev_guess) == 0 {
			break
		}
	}
	return guess
}

type testval struct {
	n      string
	digits int64
}

func init() {

	flag.StringVar(option, "o", "", "option")
	flag.Parse()
}

func main() {

	switch option {
	case "newt-sqrt":
		tests := []testval{{"4", 0}, {"2", 10}, {"28352352393", 30}}
		for _, t := range tests {
			unity := big.NewInt(10)
			unity.Exp(unity, big.NewInt(t.digits), nil)
			n, _ := big.NewInt(0).SetString(t.n, 10)
			n.Mul(n, unity)
			sqrt_n := sqrt(n, unity)
			fmt.Println("Sqrt ", n, " = ", sqrt_n)
		}
	default:
		fmt.Println("that's not a valid option")
	}

}
