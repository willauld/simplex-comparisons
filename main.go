package main

import (
	"fmt"
	"log"
	"os"
)

const (
	PlusInfinity  = 999
	NinusInfinity = -999
)

// This program is to compare and contrast different simplex methods.
// Using the same utility code accross the methods allows us to have
// similar environments so we can see the real differences with less
// distraction.
// The current input assumption is:
//
//	Maximize z = cx
//		s.t.
//		Ax <= b
//		x>=0
//
// Simplex method here reliese on b>=0 (does not use artificial variable
// to overcome negitive b[i]). The daul simplex seems to need negitive b[i]
func main() {
	var T modernTableau
	var f *os.File
	var err error
	var fname string
	//fname = "expl_1.txt" // dual example
	//fname = "expl_2.txt" // dual example
	//fname = "expl_3.txt" // simplex example
	fname = "expl_4.txt" // dual example
	//fname = "expl_5.txt" // simplex example
	//f = os.Stdin
	f, err = os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	T = loaddata(f)
	fmt.Printf("\nSolve with Dual Simplex:")
	dodual(T, f)
	dX := T.getAnswerVector()
	f.Close()

	f, err = os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	T = loaddata(f)
	fmt.Printf("\nSolve with Simplex:")
	dosimplex(T, f)
	pX := T.getAnswerVector()
	f.Close()

	for i, _ := range dX {
		fmt.Printf("dx%d: %10.2f  :: px%d: %10.2f\n", i, dX[i], i, pX[i])
	}

	/*
		//f = os.Stdin
		f, err = os.Open("expl_3.txt")
		//f, err = os.Open("expl_5.txt")
		if err != nil {
			log.Fatal(err)
		}
		T = loaddata(f)
		dosimplex(T, f)
		// */

}
