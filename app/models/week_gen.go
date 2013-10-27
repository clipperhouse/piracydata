// gen *models.Week
// this file was auto-generated using github.com/clipperhouse/gen
// Sun, 27 Oct 2013 15:35:30 UTC

package models

type Weeks []*Week

func (rcv Weeks) AggregateInt(fn func(*Week, int) int) (result int) {
	for _, _item := range rcv {
		result = fn(_item, result)
	}
	return result
}

func (rcv Weeks) AggregateString(fn func(*Week, string) string) (result string) {
	for _, _item := range rcv {
		result = fn(_item, result)
	}
	return result
}

func (rcv Weeks) All(fn func(*Week) bool) bool {
	for _, _item := range rcv {
		if !fn(_item) {
			return false
		}
	}
	return true
}

func (rcv Weeks) Any(fn func(*Week) bool) bool {
	for _, _item := range rcv {
		if fn(_item) {
			return true
		}
	}
	return false
}

func (rcv Weeks) Count(fn func(*Week) bool) int {
	var count = func(_item *Week, acc int) int {
		if fn(_item) {
			acc++
		}
		return acc
	}
	return rcv.AggregateInt(count)
}

func (rcv Weeks) Each(fn func(*Week)) {
	for _, _item := range rcv {
		fn(_item)
	}
}

func (rcv Weeks) JoinString(fn func(*Week) string, delimiter string) string {
	var join = func(_item *Week, acc string) string {
		if _item != rcv[0] {
			acc += delimiter
		}
		return acc + fn(_item)
	}
	return rcv.AggregateString(join)
}

func (rcv Weeks) SumInt(fn func(*Week) int) int {
	var sum = func(_item *Week, acc int) int {
		return acc + fn(_item)
	}
	return rcv.AggregateInt(sum)
}

func (rcv Weeks) Where(fn func(*Week) bool) (result Weeks) {
	for _, _item := range rcv {
		if fn(_item) {
			result = append(result, _item)
		}
	}
	return result
}

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func (rcv Weeks) Sort(less func(Weeks, int, int) bool) Weeks {
	result := make(Weeks, len(rcv))
	copy(result, rcv)

	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(result)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortWeeks(result, less, 0, n, maxDepth)
	return result
}

// IsSorted reports whether rcv is sorted.
func (rcv Weeks) IsSorted(less func(Weeks, int, int) bool) bool {
	n := len(rcv)
	for i := n - 1; i > 0; i-- {
		if less(rcv, i, i-1) {
			return false
		}
	}
	return true
}

func swapWeeks(rcv Weeks, a, b int) {
	rcv[a], rcv[b] = rcv[b], rcv[a]
}

// Insertion sort
func insertionSortWeeks(rcv Weeks, less func(Weeks, int, int) bool, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(rcv, j, j-1); j-- {
			swapWeeks(rcv, j, j-1)
		}
	}
}

// siftDown implements the heap property on rcv[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDownWeeks(rcv Weeks, less func(Weeks, int, int) bool, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(rcv, first+child, first+child+1) {
			child++
		}
		if !less(rcv, first+root, first+child) {
			return
		}
		swapWeeks(rcv, first+root, first+child)
		root = child
	}
}

func heapSortWeeks(rcv Weeks, less func(Weeks, int, int) bool, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownWeeks(rcv, less, i, hi, first)
	}

	// Pop elements, largest first, into end of rcv.
	for i := hi - 1; i >= 0; i-- {
		swapWeeks(rcv, first, first+i)
		siftDownWeeks(rcv, less, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// Engineering a Sort Function, SP&E November 1993.

// medianOfThree moves the median of the three values rcv[a], rcv[b], rcv[c] into rcv[a].
func medianOfThreeWeeks(rcv Weeks, less func(Weeks, int, int) bool, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if less(rcv, m1, m0) {
		swapWeeks(rcv, m1, m0)
	}
	if less(rcv, m2, m1) {
		swapWeeks(rcv, m2, m1)
	}
	if less(rcv, m1, m0) {
		swapWeeks(rcv, m1, m0)
	}
	// now rcv[m0] <= rcv[m1] <= rcv[m2]
}

func swapRangeWeeks(rcv Weeks, a, b, n int) {
	for i := 0; i < n; i++ {
		swapWeeks(rcv, a+i, b+i)
	}
}

func doPivotWeeks(rcv Weeks, less func(Weeks, int, int) bool, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's Ninther, median of three medians of three.
		s := (hi - lo) / 8
		medianOfThreeWeeks(rcv, less, lo, lo+s, lo+2*s)
		medianOfThreeWeeks(rcv, less, m, m-s, m+s)
		medianOfThreeWeeks(rcv, less, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeWeeks(rcv, less, lo, m, hi-1)

	// Invariants are:
	//	rcv[lo] = pivot (set up by ChoosePivot)
	//	rcv[lo <= i < a] = pivot
	//	rcv[a <= i < b] < pivot
	//	rcv[b <= i < c] is unexamined
	//	rcv[c <= i < d] > pivot
	//	rcv[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if less(rcv, b, pivot) { // rcv[b] < pivot
				b++
			} else if !less(rcv, pivot, b) { // rcv[b] = pivot
				swapWeeks(rcv, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if less(rcv, pivot, c-1) { // rcv[c-1] > pivot
				c--
			} else if !less(rcv, c-1, pivot) { // rcv[c-1] = pivot
				swapWeeks(rcv, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// rcv[b] > pivot; rcv[c-1] < pivot
		swapWeeks(rcv, b, c-1)
		b++
		c--
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	n := min(b-a, a-lo)
	swapRangeWeeks(rcv, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRangeWeeks(rcv, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortWeeks(rcv Weeks, less func(Weeks, int, int) bool, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortWeeks(rcv, less, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotWeeks(rcv, less, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSortWeeks(rcv, less, a, mlo, maxDepth)
			a = mhi // i.e., quickSortWeeks(rcv, mhi, b)
		} else {
			quickSortWeeks(rcv, less, mhi, b, maxDepth)
			b = mlo // i.e., quickSortWeeks(rcv, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSortWeeks(rcv, less, a, b)
	}
}
