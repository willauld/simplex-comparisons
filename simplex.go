package main

import (
	"bufio"
	"fmt"
	"os"
)

/************************************************************/
/***** Solves the LPP by "SIMPLEX" method i.e. by table *****/
/************************************************************/
/*--------------------------------------------------------------------------*\
         Cj        5        4        3        0        0        0        miniRatio
cB       xB        b        a1       a2       a3       a4       a5       a6    bi/aij
0        x4        5        2        3        1        1        0        0        2.5
0        x5        11       4        1        2        0        1        0        2.75
0        x6        8        3        4        2        0        0        1        2.66
----------------------------------------------------------------------------
Zj-Cj                       -5       -4       -3       0        0        0
----------------------------------------------------------------------------
5        x1        2.5      1        1.5      0.5      0.5      0        0        5
0        x5        1        0        -5       0        -2       1        0        infinity
0        x6        105      0        -0.5     0.5      -1.5     0        1        1
----------------------------------------------------------------------------
Zj-Cj                       0        3.5      -0.5     2.5      0        0
----------------------------------------------------------------------------
5        x1        2        1        2        0        2        0        -1
0        x5        1        0        -5       0        -2       1        0
3        x3        1        0        -1       1        -3       0        2
----------------------------------------------------------------------------
Zj-Cj                       0        3        0        1        0        1
----------------------------------------------------------------------------
So the solution is :-
x1=2        x2=0        x3=1        x4=0        x5=1        x6=0
max(z) = 5*2 + 4*0 + 3*1 = 13.
\*--------------------------------------------------------------------------*/

func dosimplex(T tableau, f *os.File) {
	/*
	  float c[M]={5,4,3,0,0,0};
	  float a[N][M]={
	    {2,3,1,1,0,0},
	    {4,1,2,0,1,0},
	    {3,4,2,0,0,1}
	  };
	  float b[N]={5,11,8};
	  float temp[M]={0,0,0,0,0,0};
	*/
	var tempminpos int /* Stores the minimum valued position
	   of {Zj-Cj} i.e. coming in variable */
	var miniratio []float64 /* N Stores the value of the ratio b[i]/a[i][j] */
	var miniratiominpos int /* Stores the minimum valued position of
	   b[i]/a[i][j] i.e. going out variable */
	//var key float64  /* Stores the key element */
	var gooutcol int /* Stores the column number which goes out */
	var z float64    /* Stores the value of the objective function */
	//var x []float64  /* M Stores the value of the variables */
	var i int /* Loop variables */
	//int basic[N];        /* Stores the basic variable */
	var nonbasic []int /* N Stores the non-basic variable */
	flag := 0          /* Terminating variable */

	/*** Initializing basic variables to 3,4,5 i.e. x4,x5,x6 ***/

	scanner := bufio.NewScanner(f)
	waitforuser := false
	if f == os.Stdin {
		waitforuser = true
	}

	miniratio = make([]float64, T.nm)
	nonbasic = make([]int, T.nm)
	for i = 0; i < T.nm; i++ {
		T.basic[i] = (i + T.n)
		nonbasic[i] = i
	}

	/*** Calculation for actual table ***/
	for flag == 0 {
		z = 0
		calctemp(T.temp, T.A, T.c, T.basic, T.n, T.m)
		fmt.Printf("\n")

		/*** Determining the incoming column ***/

		minimum(T.temp, &tempminpos, T.nm)
		display(T)

		for i = 0; i < T.m; i++ {
			basici := T.basic[i]
			bi := T.b[i]
			T.x[basici] = bi
			nonbasici := nonbasic[i]
			T.x[nonbasici] = 0
			fmt.Printf("x[%d]=%g\n", T.basic[i]+1, T.b[i])
		}
		for i = 0; i < T.nm; i++ {
			z = z + T.c[i]*T.x[i]
		}
		fmt.Printf("Max(z) = %g", z)

		/*** Determining the outgoing column ***/

		for i = 0; i < T.m; i++ {
			if T.A[i][tempminpos] <= 0 {
				miniratio[i] = PlusInfinity
				continue
			}
			bi := T.b[i]
			aTempminpos := T.A[i][tempminpos]
			miniratio[i] = bi / aTempminpos
		}
		minimum(miniratio, &miniratiominpos, T.m)
		for i = 0; i < T.nm; i++ {
			if miniratiominpos == i { // why does not not just set the go without for loop? wga
				gooutcol = T.basic[i]
			}
		}
		fmt.Printf("\nComing in variable = X%d\t", tempminpos+1)
		fmt.Printf("Going out variable = X%d\n", gooutcol+1)

		/*** Changing the basic and non-basic variable ***/

		T.basic[miniratiominpos] = tempminpos
		nonbasic[tempminpos] = gooutcol

		/*** Performing the operations to bring similar expressions in
		  in-coming variable as out-going variable by row operations ***/

		pivot(T, miniratiominpos, tempminpos)

		if waitforuser {
			scanner.Scan()
		}

		/*** Terminating condition ***/

		for i = 0; i < T.nm; i++ {
			flag = 1
			if T.temp[i] < 0 {
				flag = 0
				break
			}
		}
	}
	fmt.Printf("\nPress any key to exit...\n")
	if waitforuser {
		scanner.Scan()
	}
}
