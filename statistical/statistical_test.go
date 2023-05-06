package statistical

import (
	"fmt"
	"math"
	"testing"
)

func approxEqual(t *testing.T, a, b, tolerance float64) bool {
	t.Helper()

	return math.Abs(a-b) < tolerance
}

func TestFrequency(t *testing.T) {
	items := []string{"a", "b", "b", "c", "c", "c", "d"}

	freq := Frequency(items)

	if freq["a"] != 1 {
		t.Errorf("Expected frequency of 'a' to be 1, got %d", freq["a"])
	}

	if freq["b"] != 2 {
		t.Errorf("Expected frequency of 'b' to be 2, got %d", freq["b"])
	}

	if freq["c"] != 3 {
		t.Errorf("Expected frequency of 'c' to be 3, got %d", freq["c"])
	}

	if freq["d"] != 1 {
		t.Errorf("Expected frequency of 'd' to be 1, got %d", freq["d"])
	}
}

func TestMedian(t *testing.T) {
	s := []float64{1, 2, 3, 4}

	median, err := Median(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if median != 2.5 {
		t.Errorf("Expected median to be 2.5, got %v", median)
	}

	s = []float64{1, 2, 3}

	median, err = Median(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if median != 2 {
		t.Errorf("Expected median to be 2, got %v", median)
	}

	s = []float64{1, 2}

	median, err = Median(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if median != 1.5 {
		t.Errorf("Expected median to be 1.5, got %v", median)
	}

	s = []float64{}

	if _, err = Median(s); err == nil {
		t.Errorf("Expected error calculating median of empty slice")
	}
}

func TestRange(t *testing.T) {
	s := []float64{1, 2, 3, 4}

	min, max, err := Range(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if min != 1 || max != 4 {
		t.Errorf("Expected range to be (1, 4), got (%v, %v)", min, max)
	}

	s = []float64{1}

	if _, _, err = Range(s); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	s = []float64{}

	if _, _, err = Range(s); err == nil {
		t.Errorf("Expected error calculating range of empty slice")
	}
}

func TestMean(t *testing.T) {
	s := []float64{1, 2, 3, 4}

	mean := Mean(s)

	if mean != 2.5 {
		t.Errorf("Expected mean to be 2.5, got %v", mean)
	}

	s = []float64{1, 2}

	mean = Mean(s)

	if mean != 1.5 {
		t.Errorf("Expected mean to be 1.5, got %v", mean)
	}

	s = []float64{}

	mean = Mean(s)

	if !math.IsNaN(mean) {
		t.Errorf("Expected mean of empty slice to be NaN, got %v", mean)
	}
}

func TestVariance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []float64
		want    float64
		wantErr bool
	}{
		{
			name:  "empty slice",
			input: []float64{},
			want:  0,
			// Variance requires at least two elements
			wantErr: true,
		},
		{
			name:  "single element slice",
			input: []float64{1},
			want:  0,
			// Variance requires at least two elements
			wantErr: true,
		},
		{
			name:    "two element slice",
			input:   []float64{1, 2},
			want:    0.5,
			wantErr: false,
		},
		{
			name:    "three element slice",
			input:   []float64{1, 2, 3},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add the t.Parallel() call here
			t.Parallel()

			got, err := Variance(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Variance() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !approxEqual(t, got, tt.want, 1e-6) {
				t.Errorf("Variance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStandardDeviation(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5}

	expected := 1.5811388300841898

	result, err := StandardDeviation(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !approxEqual(t, result, expected, 0.0001) {
		t.Errorf("Result (%v) does not match expected (%v)", result, expected)
	}
}

func TestPercentile(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	expected := 3

	result, err := Percentile(input, 0.5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("Result (%v) does not match expected (%v)", result, expected)
	}
}

func ExampleFrequency() {
	items := []string{"a", "b", "b", "c", "c", "c", "d"}
	freq := Frequency(items)
	fmt.Printf("Frequency: %v\n", freq)
	// Output: Frequency: map[a:1 b:2 c:3 d:1]
}

func ExampleMedian() {
	s := []float64{1, 2, 3, 4, 5}
	median, _ := Median(s)
	fmt.Printf("Median: %v\n", median)
	// Output: Median: 3
}

func ExampleRange() {
	s := []float64{1, 2, 3, 4, 5}
	min, max, _ := Range(s)
	fmt.Printf("Range: (%v, %v)\n", min, max)
	// Output: Range: (1, 5)
}

func ExampleMean() {
	s := []float64{1, 2, 3, 4, 5}
	mean := Mean(s)
	fmt.Printf("Mean: %v\n", mean)
	// Output: Mean: 3
}

func ExampleVariance() {
	s := []float64{1, 2, 3, 4, 5}
	variance, _ := Variance(s)
	fmt.Printf("Variance: %v\n", variance)
	// Output: Variance: 2.5
}

func ExampleStandardDeviation() {
	s := []float64{1, 2, 3, 4, 5}
	stdDev, _ := StandardDeviation(s)
	fmt.Printf("Standard Deviation: %v\n", stdDev)
	// Output: Standard Deviation: 1.5811388300841898
}

func ExamplePercentile() {
	s := []int{1, 2, 3, 4, 5}
	percentile, _ := Percentile(s, 0.5)
	fmt.Printf("Percentile: %v\n", percentile)
	// Output: Percentile: 3
}
