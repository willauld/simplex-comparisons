package main

import (
	"bufio"
	"fmt"
	"os"
)

/****************************************************/
/***** Solves the LPP by Primal Simplex method *****/
/***************************************************/

func dosimplex(T modernTableau, f *os.File) {
	var colOfEnteringVar int
	var rowOfExitingVar int
	flag := 0 // Terminating variable
	itr := 1

	scanner := bufio.NewScanner(f)
	waitforuser := false
	if f == os.Stdin {
		waitforuser = true
	}

	/*** Calculation for actual table ***/
	for flag == 0 {
		/*** Determining the incoming variable / column ***/
		colOfEnteringVar = T.ObjMostNegativeCoefficient()
		if colOfEnteringVar < 0 {
			// we are done
			flag = 0
			break
		}

		/*** Determining the outgoing column ***/
		rowOfExitingVar = T.MinRatioTest(colOfEnteringVar)

		/*** Performing the operations to bring similar expressions in
		  in-coming variable as out-going variable by row operations ***/

		T.pivot(rowOfExitingVar, colOfEnteringVar)

		PrintTableau(T, itr)
		itr++
		if itr > 10 {
			os.Exit(1)
		}
		if waitforuser {
			scanner.Scan()
		}
	}
	fmt.Printf("\nPress any key to exit...\n")
	if waitforuser {
		scanner.Scan()
	}
}
