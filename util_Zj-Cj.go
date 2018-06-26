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

func loaddata(f *os.File) (Tab tableau) {
	var i, j int

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
	return T
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
