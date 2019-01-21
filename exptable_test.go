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

func TestTableComputeRandom(t *testing.T) {
	var b, m, max big.Int
	for i := 2; i < 10; i++ {
		max.SetUint64(1)
		max.Lsh(&max, 1<<uint(i))
		m.Rand(rnd, &max)
		b.Rand(rnd, &m)
		testTableCompute(t, &b, &m, 100)
	}
}

func TestTableCompute1307_4425(t *testing.T) {
	var l, m, b big.Int
	l.SetUint64(4425)
	m.SetUint64(1)
	m.Lsh(&m, 1307)
	m.Sub(&m, &l)
	b.SetUint64(2)
	testTableCompute(t, &b, &m, 100)
}

func TestTableCompute4682_190265(t *testing.T) {
	var l, m, b big.Int
	l.SetUint64(190265)
	m.SetUint64(1)
	m.Lsh(&m, 4682)
	m.Sub(&m, &l)
	b.SetUint64(2)
	testTableCompute(t, &b, &m, 10)
}

func testTableCompute(t *testing.T, b, m *big.Int, nIters int) {
	var table Table
	table.Compute(b, m, 4)

	for i := 0; i < nIters; i++ {
		var s, r1, r2 big.Int
		s.Rand(rnd, m)
		table.Exp(&r1, &s)
		r2.Exp(b, &s, m)
		if r1.Cmp(&r2) != 0 {
			t.Fatalf("%v^%v mod  %v = %v != %v", b, &s, m, &r2, &r1)
		}
	}
}

func BenchmarkTableExp2BmC_w8(b *testing.B) { benchmarkTableExp2BmC(b, 8) }
func BenchmarkTableExp2BmC_w7(b *testing.B) { benchmarkTableExp2BmC(b, 7) }
func BenchmarkTableExp2BmC_w6(b *testing.B) { benchmarkTableExp2BmC(b, 6) }
func BenchmarkTableExp2BmC_w5(b *testing.B) { benchmarkTableExp2BmC(b, 5) }
func BenchmarkTableExp2BmC_w4(b *testing.B) { benchmarkTableExp2BmC(b, 4) }
func BenchmarkTableExp2BmC_w3(b *testing.B) { benchmarkTableExp2BmC(b, 3) }
func BenchmarkTableExp2BmC_w2(b *testing.B) { benchmarkTableExp2BmC(b, 2) }
func BenchmarkTableExp2BmC_w1(b *testing.B) { benchmarkTableExp2BmC(b, 1) }

func BenchmarkTableExp_w8(b *testing.B) { benchmarkTableExp(b, 8) }
func BenchmarkTableExp_w7(b *testing.B) { benchmarkTableExp(b, 7) }
func BenchmarkTableExp_w6(b *testing.B) { benchmarkTableExp(b, 6) }
func BenchmarkTableExp_w5(b *testing.B) { benchmarkTableExp(b, 5) }
func BenchmarkTableExp_w4(b *testing.B) { benchmarkTableExp(b, 4) }
func BenchmarkTableExp_w3(b *testing.B) { benchmarkTableExp(b, 3) }
func BenchmarkTableExp_w2(b *testing.B) { benchmarkTableExp(b, 2) }
func BenchmarkTableExp_w1(b *testing.B) { benchmarkTableExp(b, 1) }

func benchmarkTableExp2BmC(b *testing.B, w uint) {
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
	table.mIsTwoBMinusC = false
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		table.Exp(&r, &s)
	}
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
