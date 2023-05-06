package statistical

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

//////
// Const, vars, and types.
//////

// Numbers is a constraint that permits any numeric type.
type Numbers interface {
	constraints.Signed | constraints.Float
}

//////
// Exported functionalities.
//////

// Frequency calculates the frequency of elements in a slice and returns a map
// with the counts of each element.
func Frequency[T comparable](items []T) map[T]int {
	freq := make(map[T]int)

	for _, item := range items {
		freq[item]++
	}

	return freq
}

// Median calculates the median value of a slice of numbers.
func Median[T Numbers](s []T) (T, error) {
	n := len(s)

	if n == 0 {
		return *new(T), fmt.Errorf("cannot calculate median of empty slice")
	}

	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })

	if n%2 == 0 {
		i := n / 2

		return (s[i-1] + s[i]) / 2, nil
	}

	i := (n - 1) / 2

	return s[i], nil
}

// Range calculates the range of a slice of numbers (maximum value minus the
// minimum value).
func Range[T Numbers](s []T) (T, T, error) {
	n := len(s)

	if n == 0 {
		return *new(T), *new(T), fmt.Errorf("cannot calculate range of empty slice")
	}

	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })

	return s[0], s[n-1], nil
}

// Mean calculates the mean (average) value of a slice of numbers.
func Mean(s []float64) float64 {
	sum := 0.0

	for _, x := range s {
		sum += x
	}

	return sum / float64(len(s))
}

// Variance calculates the variance of a slice of numbers.
func Variance(s []float64) (float64, error) {
	n := len(s)

	if n < 2 {
		return 0, fmt.Errorf("variance requires at least two elements")
	}

	mean := Mean(s)

	variance := 0.0

	for _, x := range s {
		variance += (x - mean) * (x - mean)
	}

	variance /= float64(n - 1)

	return variance, nil
}

// StandardDeviation calculates the standard deviation of a slice of numbers.
func StandardDeviation(s []float64) (float64, error) {
	variance, err := Variance(s)
	if err != nil {
		return 0, err
	}

	return math.Sqrt(variance), nil
}

// Percentile calculates the percentile value of a slice of numbers for a given
// percentage (between 0 and 100).
func Percentile[T Numbers](s []T, p float64) (T, error) {
	n := len(s)

	if n == 0 {
		return *new(T), fmt.Errorf("cannot calculate percentile of empty slice")
	}

	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })

	rank := float64(n-1) * p

	idx := int(rank)

	if idx == n-1 {
		return s[n-1], nil
	}

	frac := rank - float64(idx)

	result := float64(s[idx])*(1-frac) + float64(s[idx+1])*frac

	return T(result), nil
}
