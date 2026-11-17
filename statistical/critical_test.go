package statistical_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/thalesfsp/go-common-types/statistical"
)

func TestPercentileRejectsInvalidFraction(t *testing.T) {
	for _, p := range []float64{-0.1, 1.01, 90, math.NaN(), math.Inf(1), math.Inf(-1)} {
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Errorf("invalid percentile %v panicked: %v", p, recovered)
				}
			}()

			input := []int{5, 1, 3}

			_, err := statistical.Percentile(input, p)
			if err == nil || !reflect.DeepEqual(input, []int{5, 1, 3}) {
				t.Errorf("invalid percentile %v accepted or mutated input: err=%v input=%v", p, err, input)
			}
		}()
	}

	for _, p := range []float64{0, 0.5, 1} {
		got, err := statistical.Percentile([]int{5, 1, 3}, p)

		want := int(1 + 4*p)
		if err != nil || got != want {
			t.Errorf("supported fraction %v: got=%d want=%d err=%v", p, got, want, err)
		}
	}
}
