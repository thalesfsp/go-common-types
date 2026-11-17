# SafeSet

## Overview

SafeSet is a thread-safe, generic set implementation for Go, which provides safe concurrent access to its elements. It is a generic implementation that accepts values of any `comparable` type. ***It's a SafeSlice but with unique elements.***

Use `New` before concurrent access. Equality uses Go `==` for comparable values
(including pointer identity and dynamic interface types), and `reflect.DeepEqual`
for noncomparable values. Formatting methods do not define identity. Nil and empty
slices are distinct. Membership and insertion scan the ordered values in O(n) time.
The container protects its own storage; callers must synchronize mutations of
referenced elements. Callbacks run on a snapshot without holding collection locks.

## Installation

Use `go get` to add the `safeset` package to your project:

```sh
go get github.com/thalesfsp/go-common-types/safeset
```

## Usage

Example:

```go
package main

import (
	"fmt"
	"github.com/thalesfsp/go-common-types/safeset"
)

func main() {
	ss := safeset.New[int]()
	ss.Add(1)
	ss.Add(1)
	ss.Add(1)
	ss.Add(2)
	ss.Add(3)

	fmt.Println(ss) // [1 2 3]
}
```

## License

See [`LICENSE`](LICENSE) file for more details.

## Contributing

Feel free to open issues or submit pull requests with improvements or bug fixes. Please ensure that your code follows the coding standards.
