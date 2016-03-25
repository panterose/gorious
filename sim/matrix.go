package sim

import (
	"fmt"
	"math/rand"
)

//Matrix simple matrix object
type Matrix struct {
	rows, cols int
	slice      []float32
}

//NewMatrix constructor
func NewMatrix(rows int, cols int) Matrix {
	return Matrix{rows, cols, make([]float32, rows*cols)}
}

func NewRandomMatrix(rows int, cols int, seed int64) Matrix {
	sd := rand.NewSource(seed)
	rng := rand.New(sd)
	array := make([]float32, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			array[r*cols+c] = rng.Float32() * 1000
		}
	}
	return Matrix{rows, cols, array}
}

//Dims  returns the dimension of the matrix
func (m *Matrix) Dims() (int, int) {
	return m.rows, m.cols
}

//Copy the matrix in a newly assign matrix
func (m *Matrix) Clone() Matrix {
	ar, ac := m.Dims()
	copy := NewMatrix(ar, ac)
	for r := 0; r < ar; r++ {
		for c := 0; c < ac; c++ {
			copy.slice[r*ac+c] = m.slice[r*ac+c]
		}
	}
	return copy
}

//Redims change the dimension and copy previous state on new array double
func (m *Matrix) Redims(nr, nc int) error {
	ar, ac := m.Dims()
	mr, mc := min(ar, nr), min(ac, nc)
	oldSlice := m.slice
	//fmt.Printf("Will copy only: %d / %d \n", mr, mc)

	m.slice = make([]float32, nr*nc)
	m.rows = nr
	m.cols = nc
	//copy the old values
	for r := 0; r < mr; r++ {
		for c := 0; c < mc; c++ {
			//fmt.Printf("Assigning %v to %v / %v \n", oldSlice[r*ar+c], r, c)
			m.slice[r*m.cols+c] = oldSlice[r*ac+c]
		}
	}
	return nil
}

//GET change the dimension and copy previous state on new array double
func (m *Matrix) Get(r, c int) (float32, error) {
	if r > m.rows {
		return 0, fmt.Errorf("matrix.get: rows %d > %d", r, m.rows)
	}

	if c > m.cols {
		return 0, fmt.Errorf("matrix.get: cols %d > %d", c, m.cols)
	}

	return m.slice[r*m.cols+c], nil
}

//Mult multiply the all element of the matrix by a scalar
func (a *Matrix) Mult(s float32) (Matrix, error) {
	m := NewMatrix(a.rows, a.cols)
	for r := 0; r < a.rows; r++ {
		for c := 0; c < a.cols; c++ {
			m.slice[r*a.cols+c] = a.slice[r*a.cols+c] * s
		}
	}
	return m, nil
}

//Add a matrix to current one, while growing if necessary
func (a *Matrix) Add(b Matrix) error {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if br == 0 && bc == 0 {
		//nothing to add
		//fmt.Println("Nothing to add")
		return nil
	}

	// check if you need to grow
	mr, mc := max(ar, br), max(ac, bc)
	if ar < br || ac < bc {
		//fmt.Printf("Need to grow: %d / %d \n", mr, mc)
		a.Redims(mr, mc)
	}

	for r := 0; r < br; r++ {
		for c := 0; c < bc; c++ {
			//fmt.Printf("Adding %v to %v / %v\n", b.array[r*br+c], r, c)
			a.slice[r*a.cols+c] = a.slice[r*a.cols+c] + b.slice[r*bc+c]
		}
	}

	return nil
}

// simple max function on uint16
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// simple max function on uint16
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
