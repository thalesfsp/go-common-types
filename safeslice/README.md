# SafeSlice

## Overview

SafeSlice is a thread-safe, generic slice implementation for Go, which provides safe concurrent access to its elements. It is a generic implementation that accepts values of any `comparable` type.

## Features

- **Thread-safe**: Safe concurrent access with read-write mutexes for synchronization.
- **Generics**: Supports any value type, thanks to Go generics.
- **JSON Serialization**: Implements `MarshalJSON` and `UnmarshalJSON` for easy JSON serialization and deserialization.
- **Rich Functional API**: Provides a collection of functional methods for easy manipulation of slice elements.

## Table for the CRUD Operations

| Method | Description                                      | Input   | Output |
|--------|--------------------------------------------------|---------|--------|
| Add    | Appends a new element to the end of the slice.    | Element | None   |
| Get    | Retrieves an element from the slice at the index. | Index   | Element|
| Delete | Removes an element from the slice at the index.   | Index   | None   |
**| First | First return the first element.   | None   | Element   |
| Last | Last return the last element.   | None   | Element   |**

## Table for the Meta Operations

| Method  | Description                                                                                        | Input   | Output                                     |
|---------|----------------------------------------------------------------------------------------------------|---------|--------------------------------------------|
| Contains| Checks if the given element is present in the slice.                                              | Element | Boolean                                    |
| Size    | Returns the number of elements in the slice.                                                       | None    | Number of elements                          |
| Empty   | Checks if the slice is empty.                                                                       | None    | Boolean                                    |
| Clone   | Returns a new copy of the slice.                                                                   | None    | New SafeSlice with same elements as original|
| Index   | Returns the index of the first occurrence of the given element in the slice. If not found, returns -1 and false.| Element | Index and Boolean                          |
| Unique  | Returns a new SafeSlice with all duplicates removed.                                              | None    | New SafeSlice with unique elements         |

Note: This is not a complete list of methods. Please refer to the documentation for a full list of methods and their descriptions.

## Table for the Collection Operations (Higher-Order Functions)

| Method     | Description                                                                                              | Input                           | Output                          |
|------------|----------------------------------------------------------------------------------------------------------|---------------------------------|---------------------------------|
| All        | Checks if all elements in the slice satisfy a given condition (predicate) and returns a boolean value. | Predicate (element)             | Boolean                         |
| Map        | Applies a given function to all elements in the slice and returns a new slice containing the results.  | Function (element)              | New slice                       |
| Filter     | Creates a new slice containing only the elements that satisfy a given condition (predicate).           | Predicate (element)             | New slice                       |
| Each       | Iterates over all elements in the slice and applies a given function to each element.                   | Function (element)              | None                            |
| Reduce     | Applies a given function to the elements and returns a single result.                                   | Function (accumulator, element) | Accumulator or error            |
| Find       | Finds the first element that satisfies a given condition (predicate).                                  | Predicate (element)             | Element or error                |
| Any        | Checks if any element in the slice satisfies a given condition (predicate) and returns a boolean value. | Predicate (element)             | Boolean                         |
| TakeWhile  | Takes elements from the slice until a given condition (predicate) returns false.                        | Predicate (element)             | New slice containing elements   |
| DropWhile  | Drops elements from the slice until a given condition (predicate) returns false.                        | Predicate (element)             | New slice containing elements   |

## Table for the Set Operations

| Method     | Description                                                                                      | Input         | Output                                          |
|------------|--------------------------------------------------------------------------------------------------|---------------|-------------------------------------------------|
| Union      | Creates a new slice containing all elements from both input slices.                             | SafeSlice (T) | New slice containing all unique elements.       |
| Difference | Creates a new slice containing only the elements that exist in the first slice but not in the second. | SafeSlice (T) | New slice containing elements unique to first. |
| Subset     | Checks if all elements in the first slice exist in the second slice.                            | SafeSlice (T) | Boolean                                         |
| Superset   | Checks if all elements in the second slice exist in the first slice.                            | SafeSlice (T) | Boolean    
                                     |

## Installation

Use `go get` to add the `safeslice` package to your project:

```sh
go get github.com/thalesfsp/go-common-types/safeslice
```

## Usage

Example:

```go
package main

import (
	"fmt"
	"github.com/thalesfsp/go-common-types/safeslice"
)

func main() {
	ss := safeslice.New[int]()
	ss.Add(1)
	ss.Add(2)
	ss.Add(3)

	ssProcessed := ss.Each(func(e int) { e * 2 })
	
	fmt.Println(ssProcessed) // [2, 4, 6]
}
```

## License

See [`LICENSE`](LICENSE) file for more details.

## Contributing

Feel free to open issues or submit pull requests with improvements or bug fixes. Please ensure that your code follows the coding standards.