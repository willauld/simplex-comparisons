package main

import (
	"bufio"
	"fmt"
	"os"
)

/***************************************************/
/***** Solves the LPP by "DUAL SIMPLEX" method *****/
/***************************************************/

func dodual(T tableau, f *os.File) {
	var bminpos int /* Stores the minimum valued position
	   of {Zj-Cj} i.e. coming in variable */
	var maxratio []float64 /* Stores the value of the ratio Zj-Cj/a[i][j] */
	var maxratiomaxpos int /* Stores the minimum valued position of
	   b[i]/a[i][j] i.e. going out variable */
	//var key float64  /* Stores the key element */
	var gooutcol int /* Stores the column number which goes out */
	var incomingcol int
	var z float64 /* Stores the value of the objective function */
	var i int     /* Loop variables */
	flag := 0     /* Terminating variable */

	maxratio = make([]float64, T.nm)

	scanner := bufio.NewScanner(f)
	waitforuser := false
	if f == os.Stdin {
		waitforuser = true
	}

	/*** Do the dual algorithm ***/
	for flag == 0 {
		/*** Terminating condition ***/
		for i = 0; i < T.m; i++ {
			flag = 1
			//if(b[i]<=0) // original WGA
			if T.b[i] < 0 {
				flag = 0
				break
			}
		}

		z = 0
		calctemp(T.temp, T.A, T.c, T.basic, T.n, T.m)
		fmt.Printf("\n")

		display(T)

		/*** Determining the outgoing column ***/

		minimum(T.b, &bminpos, T.m) // TODO FIXME bminpos should only return a value if it is strictly negative; because we didn't terminate above this will do so
		gooutcol = T.basic[bminpos]

		/*** Determining the incoming column ***/

		for i = 0; i < T.nm; i++ {
			if T.A[bminpos][i] >= 0 {
				maxratio[i] = NinusInfinity
				continue
			}
			maxratio[i] = T.temp[i] / T.A[bminpos][i]
		}

		maximum(maxratio, &maxratiomaxpos, 2*T.m)
		incomingcol = maxratiomaxpos
		for i = 0; i < T.nm; i++ {
			T.x[i] = 0
		}
		for i = 0; i < T.m; i++ {
			T.x[T.basic[i]] = T.b[i]
			fmt.Printf("x[%d]=%0.3g\n", T.basic[i]+1, T.b[i])
		}
		for i = 0; i < T.m; i++ { // seems like this should be updated by the pivot
			z = z + T.c[i]*T.x[i]
		}
		fmt.Printf("Max(z) = %g\n", z)
		fmt.Printf("Outgoing variable = X%d\n", gooutcol+1)
		fmt.Printf("Incoming in variable = X%d\n", incomingcol+1)

		/*** Changing the basic and non-basic variable ***/

		T.basic[bminpos] = incomingcol

		pivot(T, bminpos, incomingcol)

		if waitforuser {
			scanner.Scan()
		}
	}

	fmt.Printf("\nPress any key to exit...\n")
	if waitforuser {
		scanner.Scan()
	}
}
