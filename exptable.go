// Speed up fixed-base modular exponentation by precomputing tables.
package exptable

import (
	"math/big"
)

// Table to hold precomputations for a certain base.
//
// Fill a table using Table.Compute().  Then use Table.Exp() to compute
// modular exponents.
//
//     var table exptable.Table
//     table.Compute(&base, &modulus, 4)
//     table.Exp(&result, &exponent)
type Table struct {
	m  big.Int   // modulus
	w  uint      // window width
	wm uint      // window mask
	n  uint      // number of limbs
	v  []big.Int // the table itself
}

// Set r to b^s mod m, where b and m are the arguments given to t.Compute().
func (t *Table) Exp(r, s *big.Int) {
	var s2 big.Int
	r.SetUint64(1)
	s2.Set(s)
	for i := uint(0); i < t.n; i++ {
		if len(s2.Bits()) == 0 {
			break
		}
		ws := uint(s2.Bits()[0]) & t.wm
		s2.Rsh(&s2, t.w)
		if ws == 0 {
			continue
		}
		r.Mul(r, &t.v[(i*t.wm)+(ws-1)])
		r.Mod(r, &t.m)
	}
}

// Fills the table for base b and modulus m using window width w.
//
// Memory usage is exponential in w and modular exponentiation speed
// is proportional to 1/w.  w=4 is probably fine.
func (t *Table) Compute(b, m *big.Int, w uint) {
	t.w = w
	t.n = uint(m.BitLen()-1)/t.w + 1
	t.m.Set(m)
	t.wm = (uint(1) << t.w) - 1
	lenV := t.n * t.wm
	t.v = make([]big.Int, lenV)

	var x big.Int
	var rb big.Int
	x.Set(b)
	rb.Set(b)
	for i := uint(0); i < lenV; i += t.wm {
		for j := uint(0); j < t.wm; j++ {
			t.v[i+j].Set(&x)
			x.Mul(&x, &rb)
			x.Mod(&x, &t.m)
		}
		rb.Set(&x)
	}
}
