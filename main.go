package main

import (
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
	var T tableau
	var f *os.File
	var err error
	//*
	//f = os.Stdin
	//f, err = os.Open("expl_1.txt")
	//f, err = os.Open("expl_2.txt")
	f, err = os.Open("expl_4.txt")
	if err != nil {
		log.Fatal(err)
	}
	T = loaddata(f)
	dodual(T, f)
	/* /
	//f = os.Stdin
	f, err = os.Open("expl_3.txt")
	if err != nil {
		log.Fatal(err)
	}
	T = loaddata(f)
	dosimplex(T, f)
	*/

}
