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

	// true if m is of the form 2^b-c for small c.  This allows for faster
	// modular reduction in exponentation.
	mIsTwoBMinusC bool
	mC            big.Int // value of c in that case
	mB            uint    // value of b in that case
	mBMask        big.Int // (1 << b) - 1
}

// Set r to b^s mod m, where b and m are the arguments given to t.Compute().
func (t *Table) Exp(r, s *big.Int) {
	if t.mIsTwoBMinusC {
		t.expTwoBMinusC(r, s)
		return
	}
	t.expDefault(r, s)
}

func (t *Table) expTwoBMinusC(r, s *big.Int) {
	var s2, tmp, carry big.Int
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

		// Internally an a.Mul(b, c) does a big allocation.  It stores this
		// big allocation as capacity in a.  It can reuse the capacity of a
		// if b and c don't overlap with a.  So we use a tmp big.Int to prevent
		// an allocation on every iteration of the loop.
		tmp.Mul(r, &t.v[(i*t.wm)+(ws-1)])
		carry.Rsh(&tmp, t.mB)
		r.And(&tmp, &t.mBMask)

		tmp.Mul(&carry, &t.mC)
		r.Add(r, &tmp)

		carry.Rsh(r, t.mB)
		r.And(r, &t.mBMask)
		tmp.Mul(&carry, &t.mC)
		r.Add(r, &tmp)
	}
	r.Mod(r, &t.m)
}

func (t *Table) expDefault(r, s *big.Int) {
	var s2, tmp big.Int
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
		tmp.Mul(r, &t.v[(i*t.wm)+(ws-1)])
		r.Mod(&tmp, &t.m)
	}
}

// Fills the table for base b and modulus m using window width w.
//
// Memory usage is exponential in w and modular exponentiation speed
// is proportional to 1/w.  w=4 is probably fine.
func (t *Table) Compute(b, m *big.Int, w uint) {
	// Check whether m = 2^b-c for small c.
	bl := m.BitLen()
	t.m.Set(m)
	t.mC.SetUint64(1)
	t.mC.Lsh(&t.mC, uint(bl))
	t.mC.Sub(&t.mC, &t.m)
	if t.mC.Sign() == 1 && t.mC.BitLen() < 64 { // TODO figure out cutoff
		var one big.Int
		one.SetUint64(1)
		t.mIsTwoBMinusC = true
		t.mB = uint(bl)
		t.mBMask.SetUint64(1)
		t.mBMask.Lsh(&t.mBMask, t.mB)
		t.mBMask.Sub(&t.mBMask, &one)
	} else {
		t.mC.SetUint64(0) // free memory
	}

	// Compute limb size, etc.
	t.w = w
	t.n = uint(bl-1)/t.w + 1
	t.wm = (uint(1) << t.w) - 1
	lenV := t.n * t.wm
	t.v = make([]big.Int, lenV)

	// Compute table
	var x, rb, tmp big.Int
	x.Set(b)
	rb.Set(b)
	for i := uint(0); i < lenV; i += t.wm {
		for j := uint(0); j < t.wm; j++ {
			t.v[i+j].Set(&x)
			tmp.Mul(&x, &rb)
			x.Mod(&tmp, &t.m)
		}
		rb.Set(&x)
	}
}
