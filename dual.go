package main

import (
	"bufio"
	"fmt"
	"os"
)

/*************************************************/
/***** Solves the LPP by Dual Simplex method *****/
/*************************************************/

func dodual(T modernTableau, f *os.File) {
	flag := 0 /* Terminating variable */

	scanner := bufio.NewScanner(f)
	waitforuser := false
	if f == os.Stdin {
		waitforuser = true
	}

	/*** Do the dual algorithm ***/

	T.setUseNegitiveZ()

	PrintTableau(T, 0)
	itr := 1
	for flag == 0 {
		/*** Determining the outgoing variable ***/
		RowExitingVar := T.RHSMinBLessThanZero()

		/*** Terminating condition ***/
		if RowExitingVar < 0 {
			// No bi < 0
			flag = 0
			break
		}

		/*** Determing the entering variable / column ***/
		ColEnteringVar := T.ObjFuncMinRatio(RowExitingVar)

		T.pivot(RowExitingVar, ColEnteringVar)

		PrintTableau(T, itr)
		itr++
		if waitforuser {
			scanner.Scan()
		}
	}

	if waitforuser {
		fmt.Printf("\nPress any key to exit...\n")
		scanner.Scan()
	}
}
