package sim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMatrix(t *testing.T) {
	mat := NewMatrix(10, 20)

	assert.Equal(t, mat.cols, 20, "The cols should be 20")
	assert.Equal(t, mat.rows, 10, "The rows should be 10")
	assert.Equal(t, len(mat.slice), 20*10, "The size should be 200")
	assert.Equal(t, mat.slice[0], float32(0.0), "first element should be zero after init")

	mat.slice[9*19] = 100.0
	assert.Equal(t, mat.slice[9*19], float32(100.0), "we should get our vaue back")
}

func TestNewRandomMatrix(t *testing.T) {
	mat := NewRandomMatrix(1, 2, 0)

	assert.Equal(t, 2, mat.cols, "The cols should be 2")
	assert.Equal(t, 1, mat.rows, "The rows should be 1")
	assert.Equal(t, 2, len(mat.slice), "The size should be 2")
	assert.Equal(t, float32(945.19617), mat.slice[0], "1st element should be 0.0")
	assert.Equal(t, float32(244.96509), mat.slice[1], "1st element should be 0.0")
}

func TestAdd(t *testing.T) {
	mat1 := &Matrix{1, 2, make([]float32, 1*2)}
	mat1.slice[0] = 10.0
	mat1.slice[1] = 20.0

	assert.Equal(t, 2, mat1.cols, "The cols should be 2")
	assert.Equal(t, 1, mat1.rows, "The rows should be 1")

	mat2 := Matrix{1, 2, make([]float32, 1*2)}
	mat2.slice[0] = 11.0
	mat2.slice[1] = 22.0

	mat1.Add(mat2)

	assert.Equal(t, float32(21.0), mat1.slice[0], "we should have 10 + 11")
	assert.Equal(t, float32(42.0), mat1.slice[1], "we should have 20 + 22")

	mat3 := Matrix{1, 3, make([]float32, 1*3)}
	mat3.slice[0] = 12.0
	mat3.slice[1] = 24.0
	mat3.slice[2] = 36.0

	mat1.Add(mat3)

	assert.Equal(t, 1, mat1.rows, "The rows should be 1")
	assert.Equal(t, 3, mat1.cols, "The cols should be 3")
	assert.Equal(t, float32(33.0), mat1.slice[0], "we should have 10 + 11 + 12")
	assert.Equal(t, float32(66.0), mat1.slice[1], "we should have 20 + 22 + 24")
	assert.Equal(t, float32(36.0), mat1.slice[2], "we should have 36")

	mat4 := Matrix{0, 0, make([]float32, 0*0)}
	mat1.Add(mat4)

	assert.Equal(t, 1, mat1.rows, "The rows should be 1")
	assert.Equal(t, 3, mat1.cols, "The cols should be 3")
	assert.Equal(t, float32(33.0), mat1.slice[0], "we should have 10 + 11 + 12")
	assert.Equal(t, float32(66.0), mat1.slice[1], "we should have 20 + 22 + 24")
	assert.Equal(t, float32(36.0), mat1.slice[2], "we should have 36")
}

func testMult(t *testing.T) {
	mat := NewMatrix(1, 2)
	mat.slice[0] = 10.0
	mat.slice[1] = 20.0
	mat.Mult(3)
	assert.Equal(t, float32(30.0), mat.slice[0], "we should have 10*3=30")
	assert.Equal(t, float32(60.0), mat.slice[1], "we should have 20*3=60")
}

func testGet(t *testing.T) {
	mat := NewMatrix(1, 2)
	mat.slice[0] = 10.0
	mat.slice[1] = 20.0
	val1, _ := mat.Get(0, 0)
	val2, _ := mat.Get(0, 1)
	assert.Equal(t, float32(10.0), val1, "we should have 10")
	assert.Equal(t, float32(20.0), val2, "we should have 20")
}

func benchmarkNewMatrix(rows int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMatrix(rows, 100)
	}
}
func BenchmarkNewMatrix10(b *testing.B)   { benchmarkNewMatrix(10, b) }
func BenchmarkNewMatrix100(b *testing.B)  { benchmarkNewMatrix(100, b) }
func BenchmarkNewMatrix1000(b *testing.B) { benchmarkNewMatrix(1000, b) }
