// Copyright 2013 Sonia Keys
// License MIT: http://www.opensource.org/licenses/MIT

package perm

import (
	"fmt"
	"math/big"
	"sort"
)

// LexNextInt takes an int slice p and reorders it in place to generate the
// next permutation in lexicographic order.
//
// The item at index 0 is considered most significant.
//
// The slice can be a ZPerm but note the slice does not have to be a ZPerm
// and items do not even have to be distinct.  For a slice with duplicate
// values, distinct multiset permutations are produced.
//
// The function returns true when it produces a new permutation.
// If p represents the last permutation in lexicographic order, it is
// left unmodified and the function returns false.
func LexNextInt(p []int) bool {
	if len(p) <= 1 {
		return false
	}
	last := len(p) - 1
	k := last - 1
	for ; p[k] >= p[k+1]; k-- {
		if k == 0 {
			return false
		}
	}
	l := last
	for ; p[k] >= p[l]; l-- {
	}
	p[k], p[l] = p[l], p[k]
	for l, r := k+1, last; l < r; l, r = l+1, r-1 {
		p[l], p[r] = p[r], p[l]
	}
	return true
}

// LexNextSort takes a value satisfying sort.Interface and reorders it in place
// to generate the next permutation in lexicographic order.
//
// The item at index 0 is considered most significant.
//
// Argument elements do not have to be distinct.  For an argument with duplicate
// element values, distinct multiset permutations are produced.
//
// The function returns true when it produces a new permutation.
// If p represents the last permutation in lexicographic order, it is
// left unmodified and the function returns false.
func LexNextSort(p sort.Interface) bool {
	if p.Len() <= 1 {
		return false
	}
	last := p.Len() - 1
	k := last - 1
	for ; !p.Less(k, k+1); k-- {
		if k == 0 {
			return false
		}
	}
	l := last
	for ; !p.Less(k, l); l-- {
	}
	p.Swap(k, l)
	for l, r := k+1, last; l < r; l, r = l+1, r-1 {
		p.Swap(l, r)
	}
	return true
}

// log2 returns ⌈log₂ x⌉ for x > 0.
func log2(x int) (n uint) {
	x--
	for ; x > 0xff; x >>= 8 {
		n += 8
	}
	for ; x > 0; x >>= 1 {
		n++
	}
	return n
}

// LexRank returns rank of p relative to lexicographic order.
//
// Time complexity claimed by algorithm authors is O(n log n) in len(p).
func (p ZPerm) LexRank() *big.Int {
	// Ref. Blai Bonet. "Efficient Algorithms to Rank and Unrank Permutations
	// in Lexicographic Order", Blai Bonet.
	k := log2(len(p))
	t := make([]int, 1<<(k+1))
	var r, b big.Int
	for i, c := range p {
		nd := 1<<k + c
		for j := uint(0); j < k; j++ {
			if nd&1 == 1 {
				c -= t[nd&^1]
			}
			t[nd]++
			nd >>= 1
		}
		t[nd]++
		r.Mul(&r, b.SetInt64(int64(len(p)-i)))
		r.Add(&r, b.SetInt64(int64(c)))
	}
	return &r
}

// LexPerm creates the permutation of n integers with rank i relative to
// lexicographic order.
//
// Time complexity claimed by algorithm authors is O(n log n) in len(p).
func LexPerm(i *big.Int, n int) (ZPerm, bool) {
	// Ref. Blai Bonet. "Efficient Algorithms to Rank and Unrank Permutations
	// in Lexicographic Order", Blai Bonet.
	fmt.Println("LexPerm i, n:", i, n)
	f, ok := NewFact(i, n)
	if !ok {
		return nil, false
	}
	fmt.Println("LexPerm f:", f)
	p := make(ZPerm, n)
	k := log2(n)
	k2 := 1 << k
	t := make([]int, 1<<(k+1))
	for i := uint(0); i <= k; i++ {
		for j, end := uint(1), uint(1)<<i; j <= end; j++ {
			t[end+j-1] = 1 << (k - i)
		}
	}
	for i := len(f) - 1; i >= 0; i-- {
		d := f[i]
		fmt.Println("d:", d)
		nd := 1
		for j := uint(0); j < k; j++ {
			t[nd]--
			nd <<= 1
			if d >= t[nd] {
				d -= t[nd]
				nd++
			}
		}
		t[nd] = 0
		p[len(f)-i-1] = nd - k2
	}
	nd := 1
	for j := uint(0); j < k; j++ {
		nd <<= 1
		nd += 1 - t[nd]
	}
	p[len(f)] = nd - k2
	return p, true
}
