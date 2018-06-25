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
	basic []int // m basic variables
	n     int   // number of variable given
	m     int   // number of constraints given
	nm    int   // number of variables including slack / surplus variables
}

func (m modernTableau) c(index int) float64 {
	// TODO convert the current tableau to single matrix
	// with access methods for convienience and readability
	return 0.0
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

	fmt.Printf("\nEnter number of terms in objective function\n")
	scanner.Scan()
	text := scanner.Text()
	n, _ := strconv.ParseInt(text, 10, 32)
	T.n = int(n)
	fmt.Printf("\nEnter number of constraints\n")
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

	fmt.Printf("\nEnter the object function co-efficients\n")
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
	fmt.Printf("\nEnter the co-efficient of constraints\n")
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
	fmt.Printf("\nEnter values of bi's\n")
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
	PrintTableau(Tab, 0)
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
	return Tab
}

func (T modernTableau) ObjMostNegativeCoefficient() int {
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
func (T modernTableau) MinRatioTest(column int) int {
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
func (T modernTableau) pivot(row, col int) {
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
	fmt.Printf("%10s", "obj (z)") // no basic in object row
	for j := 0; j < T.nm; j++ {
		fmt.Printf("%10.2f", T.Mat[T.m][j]) // Cj
	}
	fmt.Printf("%10.2f", T.Mat[T.m][T.nm]) // z
	fmt.Printf(" <= Max(z)\n")
}

func calctemp(T tableau) {
	// calc zi, needed for _Zj-Cj.go codes
	// see page 80 of Introduction to Management Science by Thomas Cook 2Ed
	var i, j int
	for i = 0; i < T.nm; i++ {
		T.temp[i] = 0
		for j = 0; j < T.m; j++ {
			T.temp[i] = T.temp[i] + T.c[T.basic[j]]*T.A[j][i]
		}
		T.temp[i] = T.temp[i] - T.c[i]
	}
}

func maximum(arr []float64, arrmaxpos *int, n int) {
	var i int
	var arrmax float64
	arrmax = arr[0]
	*arrmaxpos = 0
	for i = 0; i < n; i++ {
		if arr[i] > arrmax {
			arrmax = arr[i]
			*arrmaxpos = i
		}
	}
}

func minimum(arr []float64, arrminpos *int, n int) {
	var i int
	var arrmin float64
	arrmin = arr[0]
	*arrminpos = 0
	for i = 0; i < n; i++ {
		if arr[i] < arrmin {
			arrmin = arr[i]
			*arrminpos = i
		}
	}
}

func pivot(T tableau, row, col int) {

	key := T.A[row][col]
	if key == 1.0 {
		fmt.Printf("**** Pivit value is 1.0! No need to multiply PivRow\n")
	}
	T.b[row] = T.b[row] / key
	for i := 0; i < T.nm; i++ {
		T.A[row][i] = T.A[row][i] / key
	}
	for i := 0; i < T.m; i++ {
		if row == i {
			continue
		}
		key = T.A[i][col]
		if key == 0.0 {
			fmt.Printf("**** Non-Pivit value is 0.0! No need to multiply Non-PivRow\n")
		}
		for j := 0; j < T.nm; j++ {
			T.A[i][j] = T.A[i][j] - T.A[row][j]*key
		}
		T.b[i] = T.b[i] - T.b[row]*key
	}
}

func display(T tableau) {
	var i, j int
	displayframe(T.c, T.nm)
	for i = 0; i < T.m; i++ {
		fmt.Printf("\n%0.3g\tX%d\t%0.3g\t", T.c[T.basic[i]], T.basic[i]+1, T.b[i])
		for j = 0; j < T.nm; j++ {
			fmt.Printf("%0.3g\t", T.A[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\tZj-Cj\t\t")
	for i = 0; i < T.nm; i++ {
		fmt.Printf("%0.3g\t", T.temp[i])
	}
	fmt.Printf("\n\n")
}

func displayframe(c []float64, nm int) {
	var i int
	fmt.Printf("\n\t\tc[j]\t")
	for i = 0; i < nm; i++ {
		fmt.Printf("%0.2g\t", c[i])
	}
	fmt.Printf("\n")
	fmt.Printf("\nc[B]\tB\tb\t")
	for i = 0; i < nm; i++ {
		fmt.Printf("a[%d]\t", i+1)
	}
	fmt.Printf("\n")
}
