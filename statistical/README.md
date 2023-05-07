# Statistical

## Overview

Statistical is a Go package that provides a collection of statistical operations for slices of `Numbers`.

## Features

- **Generics**: Supports any value type, thanks to Go generics.
- **Rich Functional API**: Provides a collection of functional methods for easy manipulation of map elements.

## Table for the Stastitical Operations

| Method | Description                                     | Input                     | Output               |
|--------|-------------------------------------------------|---------------------------|----------------------|
| `Frequency` | Calculates the frequency of elements in a slice and returns a map with the counts of each element. | `items []T` | `map[T]int` |
| `Median` | Calculates the median value of a slice of numbers. | `s []T` | `T, error` |
| `Range` | Calculates the range of a slice of numbers (maximum value minus the minimum value). | `s []T` | `T, T, error` |
| `Mean` | Calculates the mean (average) value of a slice of numbers. | `s []float64` | `float64` |
| `Variance` | Calculates the variance of a slice of numbers. | `s []float64` | `float64, error` |
| `StandardDeviation` | Calculates the standard deviation of a slice of numbers. | `s []float64` | `float64, error` |
| `Percentile` | Calculates the percentile value of a slice of numbers for a given percentage (between 0 and 100). | `s []T, p float64` | `T, error` |


## Installation

Use `go get` to add the `statistical` package to your project:

```sh
go get github.com/thalesfsp/go-common-types/statistical
```

## Usage

Example:

```go
package main

import (
	"fmt"
	"github.com/thalesfsp/go-common-types/statistical"
)

func main() {
	items := []string{"a", "b", "b", "c", "c", "c", "d"}
	
	freq := statistical.Frequency(items)

	fmt.Println(freq) // map[a:1 b:2 c:3 d:1]
}
```

## License

See [`LICENSE`](LICENSE) file for more details.

## Contributing

Feel free to open issues or submit pull requests with improvements or bug fixes. Please ensure that your code follows the coding standards
