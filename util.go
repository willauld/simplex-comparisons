package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type tableau struct {
	c        []float64
	A        [][]float64
	b        []float64
	x        []float64
	temp     []float64
	origC    []float64
	basic    []int
	n, m, nm int
}
type modernTableau struct {
	Mat [][]float64 // (m+1) x (nm +1) matrix
	// columns + 1 for RHS / bi values
	// row + 1 for object function
	basic []int  // m basic variables
	n     int    // number of variable given
	m     int    // number of constraints given
	nm    int    // number of variables including slack / surplus variables
	zstr  string // the z or -z to be printed with the tableau
}

func loaddata(f *os.File) (Tab modernTableau) {
	var i, j int
	var T tableau

	waitforuser := false
	if f == os.Stdin {
		waitforuser = true
	}
	scanner := bufio.NewScanner(f)
	/*** Input requisite amount of data ***/

	if waitforuser {
		fmt.Printf("\nEnter number of terms in objective function\n")
	}
	scanner.Scan()
	text := scanner.Text()
	n, _ := strconv.ParseInt(text, 10, 32)
	T.n = int(n)
	if waitforuser {
		fmt.Printf("\nEnter number of constraints\n")
	}
	scanner.Scan()
	text = scanner.Text()
	m, _ := strconv.ParseInt(text, 10, 32)
	T.m = int(m)
	T.nm = T.n + T.m
	T.c = make([]float64, T.nm)
	T.x = make([]float64, T.nm)
	T.temp = make([]float64, T.nm)
	T.basic = make([]int, T.nm)
	T.b = make([]float64, T.m)
	T.A = make([][]float64, T.m)
	for i = 0; i < T.m; i++ {
		T.A[i] = make([]float64, T.nm)
	}

	if waitforuser {
		fmt.Printf("\nEnter the object function co-efficients\n")
	}
	for i = 0; i < T.n; i++ {
		scanner.Scan()
		text = scanner.Text()
		T.c[i], _ = strconv.ParseFloat(text, 64)
	}
	fmt.Printf("\nYou have entered the function as follows:-\n")
	fmt.Printf("\nMax z = ")
	for i = 0; i < T.n; i++ {
		if i == 0 {
			fmt.Printf("%g*x%d", T.c[i], i+1)
		} else {
			fmt.Printf(" + %g*x%d", T.c[i], i+1)
		}
	}
	if waitforuser {
		fmt.Printf("\nEnter the co-efficient of constraints\n")
	}
	for i = 0; i < T.m; i++ {
		for j = 0; j < T.n; j++ {
			scanner.Scan()
			text = scanner.Text()
			T.A[i][j], _ = strconv.ParseFloat(text, 64)
		}
	}
	for i = 0; i < T.m; i++ {
		T.A[i][j] = 1
		j++
	}
	if waitforuser {
		fmt.Printf("\nEnter values of bi's\n")
	}
	for i = 0; i < T.m; i++ {
		scanner.Scan()
		text = scanner.Text()
		T.b[i], _ = strconv.ParseFloat(text, 64)
	}
	for i = 0; i < T.nm; i++ {
		T.basic[i] = i + T.n
	}
	fmt.Printf("\nYou have entered the function as follows:-\n")
	for i = 0; i < T.m; i++ {
		for j = 0; j < T.n; j++ {
			if j == 0 {
				fmt.Printf(" %g*x%d ", T.A[i][j], j+1)
			} else {
				fmt.Printf(" + %g*x%d ", T.A[i][j], j+1)
			}
		}
		fmt.Printf(" <= %g\n", T.b[i])
	}
	if waitforuser {
		scanner.Scan()
		text = scanner.Text()
	}
	Tab = buildTableauFromLoadedData(T)
	//PrintTableau(Tab, 0)
	return Tab
}

func buildTableauFromLoadedData(T tableau) (Tab modernTableau) {
	Tab.n = T.n
	Tab.m = T.m
	Tab.nm = T.nm
	matCols := T.nm + 1
	matRows := T.m + 1
	Tab.Mat = make([][]float64, matRows)
	for i := 0; i < matRows; i++ {
		Tab.Mat[i] = make([]float64, matCols)
	}
	Tab.basic = make([]int, T.m)
	// Fill in Mat with Aij
	for i := 0; i < Tab.m; i++ {
		for j := 0; j < Tab.nm; j++ {
			Tab.Mat[i][j] = T.A[i][j]
		}
		Tab.Mat[i][Tab.nm] = T.b[i]
	}
	for j := 0; j < Tab.nm; j++ {
		Tab.Mat[Tab.m][j] = -T.c[j]
	}
	for i := 0; i < Tab.m; i++ {
		Tab.basic[i] = Tab.n + i
	}
	Tab.zstr = "z"
	return Tab
}

func (T *modernTableau) setUseNegitiveZ() {
	// force z to -z
	for j := 0; j <= T.nm; j++ {
		T.Mat[T.m][j] = -T.Mat[T.m][j]
	}
	T.zstr = "-z"
}

func (T *modernTableau) getAnswerVector() []float64 {
	X := make([]float64, T.nm+1)
	for i := 0; i < T.m; i++ {
		basici := T.basic[i]
		X[basici+1] = T.Mat[i][T.nm]
	}
	X[0] = T.Mat[T.m][T.nm]
	return X
}

func (T *modernTableau) ObjMostNegativeCoefficient() int {
	minj := -1
	min := 10.0
	for j := 0; j < T.nm; j++ {
		val := T.Mat[T.m][j]
		if val < 0 && val < min {
			min = val
			minj = j
		}
	}
	return minj
}
func (T *modernTableau) MinRatioTest(column int) int {
	var minratio float64
	minrow := -1
	for i := 0; i < T.m; i++ {
		aval := T.Mat[i][column]
		if aval > 0 {
			ratio := T.Mat[i][T.nm] / aval
			if ratio < minratio || minrow < 0 {
				minratio = ratio
				minrow = i
			}
		}
	}
	return minrow
}

func (T *modernTableau) RHSMinBLessThanZero() int {
	minrow := -1
	min := 0.0
	for i := 0; i < T.m; i++ {
		bi := T.Mat[i][T.nm]
		if bi < 0 {
			if bi < min {
				min = bi
				minrow = i
			}

		}
	}
	return minrow
}

func (T *modernTableau) ObjFuncMinRatio(row int) int {
	mincol := -1
	min := 0.0
	for j := 0; j < T.nm; j++ {
		// check that Xj is not already basic
		jAlreadyBasic := false
		for i := 0; i < T.m; i++ {
			if T.basic[i] == j {
				jAlreadyBasic = true
				break
			}
		}
		if !jAlreadyBasic {
			aij := T.Mat[row][j]
			if aij < 0 {
				ratio := T.Mat[T.m][j] / aij
				if ratio < min || mincol < 0 {
					min = ratio
					mincol = j
				}
			}
		}
	}
	return mincol
}

func (T *modernTableau) pivot(row, col int) {
	T.basic[row] = col // FIXME is this true for dual as well as primal?
	key := T.Mat[row][col]
	if key == 1.0 {
		fmt.Printf("**** Pivit value is 1.0! No need to multiply PivRow %d\n", row+1)
	}
	for j := 0; j <= T.nm; j++ {
		T.Mat[row][j] = T.Mat[row][j] / key
	}
	for i := 0; i <= T.m; i++ {
		if row == i {
			continue
		}
		key = T.Mat[i][col]
		if key == 0.0 {
			fmt.Printf("**** Non-Pivit value is 0.0! No need to multiply Non-PivRow %d\n", i+1)
		}
		for j := 0; j <= T.nm; j++ {
			T.Mat[i][j] = T.Mat[i][j] - T.Mat[row][j]*key
		}
	}
}

func PrintTableau(T modernTableau, itr int) {
	fmt.Printf("\n============= Iteration %d ================\n", itr)
	fmt.Printf("%10s", "Basic")
	for j := 0; j < T.nm; j++ {
		s := fmt.Sprintf("x%d", j+1)
		fmt.Printf("%10s", s) // Xi start at 1 not zero
	}
	fmt.Printf("%10s\n", "b") // bi start at 1 not zero
	for j := 0; j < T.nm+2; j++ {
		fmt.Printf("%10s", "----------")
	}
	fmt.Printf("\n")
	for i := 0; i < T.m; i++ {
		s := fmt.Sprintf("x%d", T.basic[i]+1)
		fmt.Printf("%10s", s) // var start at 1 not zero
		for j := 0; j < T.nm; j++ {
			fmt.Printf("%10.2f", T.Mat[i][j]) // Aij
		}
		fmt.Printf("%10.2f\n", T.Mat[i][T.nm]) // bi
	}
	for j := 0; j < T.nm+2; j++ {
		fmt.Printf("%10s", "----------")
	}
	fmt.Printf("\n")
	s := fmt.Sprintf("obj (%s)", T.zstr)
	fmt.Printf("%10s", s) // no basic in object row
	for j := 0; j < T.nm; j++ {
		fmt.Printf("%10.2f", T.Mat[T.m][j]) // Cj
	}
	fmt.Printf("%10.2f", T.Mat[T.m][T.nm]) // z
	fmt.Printf(" <= Max(%s)\n", T.zstr)
}
