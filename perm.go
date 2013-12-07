// Copyright 2013 Sonia Keys
// License MIT: http://www.opensource.org/licenses/MIT

// Package permute has functions to generate permutations.  And other related
// functions
package perm

//import "fmt"

// Ints returns a slice of ints containing sequential integers 0..n-1
//
// Simply a little utility function.
func Ints(n int) []int {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return p
}

// SJTRecursive uses a recursive method to generate permutations in the order
// of the Steinhaus-Johnson-Trotter algorithm, or "plain changes".
//
// The algorithm is recursive but "loopless" in the sense that there are no
// iterative loops per permuation.
//
// It takes a slice p and returns an iterator function.  The iterator
// permutes p in place and returns true for each permutation.  After all
// permutations have been generated, the iterator returns false and p is left
// in its initial order.  The values in the slice p are not considered.  The
// generator permutes the slice regardless of its contents.
func SJTRecursive(p []int) func() bool {
	f := sjtr(len(p))
	return func() bool {
		return f(p)
	}
}

// Recursive function used by perm, returns a chain of closures that
// implement a loopless recursive generator.  Successive permutations are
// generated by a single in-place swap of two adjacent items.  Information
// directing the swap is stored in the variables of the closures so that
// each permution is generated in a single descent into the chain and a
// single swap.
func sjtr(n int) func([]int) bool {
	perm := true
	switch n {
	case 0, 1:
		return func([]int) (r bool) {
			r = perm
			perm = false
			return
		}
	default:
		p0 := sjtr(n - 1)
		i := n
		var d int
		return func(p []int) bool {
			switch {
			case !perm:
			case i == n:
				i--
				perm = p0(p[:i])
				d = -1
			case i == 0:
				i++
				perm = p0(p[1:])
				d = 1
				if !perm {
					p[0], p[1] = p[1], p[0]
				}
			default:
				p[i], p[i-1] = p[i-1], p[i]
				i += d
			}
			return perm
		}
	}
}

// LexNext takes a slice p and reorders it in place to generate the next
// permutation in lexicographic order.  For a slice with duplicate values,
// distinct multiset permutations are produced.  The function returns true
// when it produces a new permutation.  If p represents the last permutation
// in lexicographic order, it is left unmodified and the function returns false.
func LexNext(p []int) bool {
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

// SJTE implements the iterative Steinhaus-Johnson-Trotter algorithm with
// Even's speedup.
//
// Given n, a number of items to permute, it returns a slice of n integers
// in order from 0 to n-1 and a function that will permute the slice in-place.
// The function permutes the slice and returns true until the permutation
// rolls over to the original order, then it returns false.  You can continue
// to call the function at this point, and the cycle of permutations repeats.
func SJTE(n int) ([]int, func() bool) {
	p := make([]int, n+2)
	d := make([]int, n+2)
	p[0] = n
	for i := range p[1:] {
		p[i+1] = i
		d[i] = -1
	}
	return p[1 : n+1 : n+1], func() bool {
		var k int
		max := -1
		for i := 1; i <= n; i++ {
			if v := p[i]; v > max && v > p[i+d[i]] {
				max = v
				k = i
			}
		}
		if k == 0 {
			p[1], p[2] = 0, 1
			for i := 3; i <= n; i++ {
				d[i] = -1
			}
			return false
		}
		nx := k + d[k]
		p[k], p[nx] = p[nx], max
		d[k], d[nx] = d[nx], d[k]
		for i := 1; i <= n; i++ {
			if p[i] > max {
				d[i] = -d[i]
			}
		}
		return true
	}
}