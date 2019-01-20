package exptable

import (
	"math/big"
	"math/rand"
	"os"
	"testing"
)

var (
	rnd *rand.Rand
)

func TestTableCompute(t *testing.T) {
	var table Table
	var b, m, l big.Int
	l.SetUint64(190265)
	m.SetUint64(1)
	m.Lsh(&m, 468)
	m.Sub(&m, &l)
	b.SetUint64(2)
	table.Compute(&b, &m, 4)

	for i := 0; i < 100; i++ {
		var s, r1, r2 big.Int
		s.Rand(rnd, &m)
		table.Exp(&r1, &s)
		r2.Exp(&b, &s, &m)
		if r1.Cmp(&r2) != 0 {
			t.Fatalf("%v^%v mod  %v = %v != %v", &b, &s, &m, &r2, &r1)
		}
	}
}

func benchmarkTableExp(b *testing.B, w uint) {
	var table Table
	var r, s, g, m, l big.Int
	l.SetUint64(190265)
	m.SetUint64(1)
	m.Lsh(&m, 4682)
	m.Sub(&m, &l)
	g.SetUint64(2)
	table.Compute(&g, &m, w)
	s.Rand(rnd, &m)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.Exp(&r, &s)
	}
}

func BenchmarkTableExp_w8(b *testing.B) {
	benchmarkTableExp(b, 8)
}
func BenchmarkTableExp_w7(b *testing.B) {
	benchmarkTableExp(b, 7)
}
func BenchmarkTableExp_w6(b *testing.B) {
	benchmarkTableExp(b, 6)
}
func BenchmarkTableExp_w5(b *testing.B) {
	benchmarkTableExp(b, 5)
}
func BenchmarkTableExp_w4(b *testing.B) {
	benchmarkTableExp(b, 4)
}
func BenchmarkTableExp_w3(b *testing.B) {
	benchmarkTableExp(b, 3)
}
func BenchmarkTableExp_w2(b *testing.B) {
	benchmarkTableExp(b, 2)
}
func BenchmarkTableExp_w1(b *testing.B) {
	benchmarkTableExp(b, 1)
}

func BenchmarkBigIntExp(b *testing.B) {
	var r, s, g, m, l big.Int
	l.SetUint64(190265)
	m.SetUint64(1)
	m.Lsh(&m, 4682)
	m.Sub(&m, &l)
	g.SetUint64(2)
	s.Rand(rnd, &m)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Exp(&g, &s, &m)
	}
}

func TestMain(m *testing.M) {
	rnd = rand.New(rand.NewSource(37))
	os.Exit(m.Run())
}
