# Safe Ordered Map

## Overview

Safe Ordered Map is a thread-safe, ordered map implementation for Go, which maintains the insertion order while providing safe concurrent access. It is a generic implementation that accepts keys as strings and values of any type.

## Features

- **Thread-safe**: Safe concurrent access with read-write mutexes for synchronization.
- **Ordered**: Maintains the insertion order, allowing ordered iterations.
- **Generics**: Supports any value type, thanks to Go generics.
- **JSON Serialization**: Implements `MarshalJSON` and `UnmarshalJSON` for easy JSON serialization and deserialization.
- **Rich Functional API**: Provides a collection of functional methods for easy manipulation of map elements.

## Table for the CRUD Operations

| Method | Description                                     | Input                     | Output               |
|--------|-------------------------------------------------|---------------------------|----------------------|
| Set    | Sets a value in the map.                        | Key (string), Value (T)    | None                 |
| Get    | Gets a value from the map.                       | Key (string)              | Value (T) or error   |
| Delete | Deletes a value from the map.                    | Key (string)              | None                 |

## Table for the Operations on Keys and Values

| Method | Description                                  | Input | Output              |
|--------|----------------------------------------------|-------|---------------------|
| Keys   | Returns a list of all keys.                   | None  | List of keys (strings) |
| Values | Returns a list of all values in the same order as the keys. | None  | List of values (T) |

## Table for the Meta Operations

| Method | Description                                           | Input | Output                                |
|--------|-------------------------------------------------------|-------|---------------------------------------|
| Size   | Returns the number of elements in the map.             | None  | Number of elements (int)              |
| Empty  | Checks if the map is empty and returns a boolean value. | None  | Boolean (true if map is empty)        |
| Clone  | Creates a deep copy of the map and returns it.          | None  | New SafeOrderedMap with same elements |
| Index  | Returns the index and value of the given key.           | Key   | Index (int), Value (T), bool (true if key exists) |

## Table Regarding Collection Operations (Higher-Order Functions)

| Method    | Description                                                                                               | Input                          | Output                                            | Modifies Original Map |
|-----------|-----------------------------------------------------------------------------------------------------------|--------------------------------|---------------------------------------------------|-----------------------|
| All       | Checks if all elements in the map satisfy a given condition (predicate) and returns a boolean value.      | Predicate (key, value)         | Boolean (true if all meet condition)              | No                    |
| Map       | Applies a given function to all elements in the map and creates a new map containing the results.        | Function (key, value)          | New map with transformed elements                 | No                    |
| Filter    | Creates a new map containing only the elements that satisfy a given condition (predicate).                | Predicate (key, value)         | New map with filtered elements                    | No                    |
| Each      | Iterates over all elements and applies a given function to each element without returning any result.    | Function (key, value)          | None                                              | No                    |
| Reduce    | Accumulates the elements in the map using a given binary function.                                        | Binary function, initial value | Accumulated single value                          | No                    |
| Find      | Returns the first element that satisfies a given predicate.                                              | Predicate (key, value)         | Key, value, boolean (true if found)               | No                    |
| Any       | Checks if any element in the map satisfies a given predicate.                                            | Predicate (key, value)         | Boolean (true if any element meets condition)     | No                    |
| TakeWhile | Returns a new ordered map containing the longest prefix of elements that satisfy a given predicate.     | Predicate (key, value)         | New map with elements that meet condition         | No                    |
| DropWhile | Returns a new ordered map with all elements after (and not including) the first element that does not satisfy a given predicate. | Predicate (key, value)         | New map with elements after not meeting condition | No                    |

## Table Regarding Set Operations

| Method    | Description                                                                                               | Input                          | Output                                            |
|-----------|-----------------------------------------------------------------------------------------------------------|--------------------------------|---------------------------------------------------|
| Union     | Returns a new map containing all elements present in the original map and the other map.                 | Another ordered map            | New map with all elements from both maps           |
| Difference| Returns a new map containing elements present in the original map but not in the other map.              | Another ordered map            | New map with elements present in original but not other|
| Subset    | Checks if all elements in the map are present in the other map.                                           | Another ordered map            | Boolean (true if all elements are present in other) |
| Superset  | Checks if all elements in the other map are present in the map.                                           | Another ordered map            | Boolean (true if all elements are present in map)   |

These set operations provide a powerful way to compare and combine the elements of two ordered maps, allowing developers to easily manipulate and transform data.

## Installation

Use `go get` to add the `safeorderedmap` package to your project:

```sh
go get github.com/thalesfsp/go-common-types/safeorderedmap
```

## Usage

Example:

```go
package main

import (
	"fmt"
	"github.com/thalesfsp/go-common-types/safeorderedmap"
)

func main() {
	som := safeorderedmap.New[int]()
	som.Set("one", 1)
	som.Set("two", 2)
	som.Set("three", 3)

	somProcessed := som.Each(func(k string, e int) {e * 2})

	fmt.Println(somProcessed) // {one: 2, two: 4, three: 6}
}
```

## License

See `LICENSE` file for more details.

## Contributing

Feel free to open issues or submit pull requests with improvements or bug fixes. Please ensure that your code follows the coding standards