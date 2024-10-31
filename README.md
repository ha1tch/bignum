# bignum: Arbitrary Precision Fixed-Point Arithmetic in Go

This Go package provides a `BigNumber` type that allows you to work with arbitrary precision numbers with fixed-point arithmetic. It supports a variety of operations, including:

* **Basic Arithmetic:** Addition, subtraction, multiplication, division, modulo
* **Exponentiation:** Raising a BigNumber to an integer power
* **Square Root:** Calculating the square root of a BigNumber
* **Trigonometric Functions:** Sine, cosine, tangent
* **Logarithm:** Approximating the natural logarithm (base e)
* **Exponential Function:** Approximating the exponential function (base e)
* **Rounding Modes:**  Round to Nearest, Round to Even (Banker's Rounding), Round Up, Round Down
* **Error Handling:** Handles overflow, division by zero, and invalid input

## Usage

```go
package main

import (
	"fmt"
	"github.com/ha1tch/bignum"
)

func main() {
	// Create two BigNumbers with a precision of 4 decimal places.
	bn1, _ := bignum.NewBigNumber("123.4567", 4, bignum.RoundToNearest)
	bn2, _ := bignum.NewBigNumber("8.9012", 4, bignum.RoundToNearest)

	// Perform addition.
	result, _ := bn1.Add(bn2)
	fmt.Println(result.String()) // Output: 132.3579

	// Calculate the square root of bn1.
	sqrt, _ := bn1.SquareRoot()
	fmt.Println(sqrt.String()) // Output: 11.1111
}
```

## Installation

```bash
go get github.com/yourusername/bignum
```

## Example

```go
package main

import (
	"fmt"
	"github.com/ha1tch/bignum"
)

func main() {
	// Create a BigNumber with a precision of 2 decimal places.
	bn, _ := bignum.NewBigNumber("123.4567", 2, bignum.RoundToNearest)

	// Perform basic operations.
	fmt.Println("Original BigNumber:", bn.String()) 
	fmt.Println("Rounded to 2 decimal places:", bn.Round(2).String())
	fmt.Println("Absolute value:", bn.AbsoluteValue().String()) 

	// Work with special cases.
	inf, _ := bignum.NewBigNumber("inf", 2, bignum.RoundToNearest)
	nan, _ := bignum.NewBigNumber("nan", 2, bignum.RoundToNearest)
	fmt.Println("Infinity:", inf.String())
	fmt.Println("NaN:", nan.String())
}
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.
