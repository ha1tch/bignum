package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"bignum"

	"github.com/shopspring/decimal"
)

const (
	// Set these based on your system's capabilities and desired test range
	maxIntegerDigits = 9 // Maximum number of integer digits to test
	maxDecimalDigits = 9 // Maximum number of decimal digits to test
)

// Generate a random number with the given integer and decimal digits.
func generateRandomNumber(integerDigits int, decimalDigits int) (decimal.Decimal, *bignum.BigNumber) {
	rand.Seed(time.Now().UnixNano())

	// Generate integer part (10^integerDigits - 1) to ensure maximum length
	integer := rand.Intn(int(math.Pow10(integerDigits))) + 1

	// Generate decimal part (10^decimalDigits - 1) to ensure maximum length
	decimalPart := rand.Intn(int(math.Pow10(decimalDigits))) + 1

	// Build string representations
	integerStr := fmt.Sprintf("%d", integer)
	decimalStr := fmt.Sprintf(".%0*d", decimalDigits, decimalPart)

	// Create decimal.Decimal and bignum.BigNumber
	d, _ := decimal.NewFromString(integerStr + decimalStr)
	bn, _ := bignum.NewBigNumber(integerStr+decimalStr, uint(decimalDigits), bignum.RoundToNearest)

	return d, bn
}

// Benchmark for addition
func BenchmarkAddition(b *testing.B) {
	// Define these variables outside the loop to avoid repeated initialization
	integerDigits1 := 3
	decimalDigits1 := 3
	integerDigits2 := 3
	decimalDigits2 := 3

	d1, bn1 := generateRandomNumber(integerDigits1, decimalDigits1)
	d2, bn2 := generateRandomNumber(integerDigits2, decimalDigits2)

	b.ResetTimer()

	// Decimal addition
	for i := 0; i < b.N; i++ {
		d1.Add(d2)
	}

	// Bignum addition
	for i := 0; i < b.N; i++ {
		bn1.Add(bn2)
	}
}

// Benchmark for multiplication
func BenchmarkMultiplication(b *testing.B) {
	// Define these variables outside the loop to avoid repeated initialization
	integerDigits1 := 3
	decimalDigits1 := 3
	integerDigits2 := 3
	decimalDigits2 := 3

	d1, bn1 := generateRandomNumber(integerDigits1, decimalDigits1)
	d2, bn2 := generateRandomNumber(integerDigits2, decimalDigits2)

	b.ResetTimer()

	// Decimal multiplication
	for i := 0; i < b.N; i++ {
		d1.Mul(d2)
	}

	// Bignum multiplication
	for i := 0; i < b.N; i++ {
		bn1.Multiply(bn2)
	}
}

// Benchmark for division
func BenchmarkDivision(b *testing.B) {
	// Define these variables outside the loop to avoid repeated initialization
	integerDigits1 := 3
	decimalDigits1 := 3
	integerDigits2 := 3
	decimalDigits2 := 3

	d1, bn1 := generateRandomNumber(integerDigits1, decimalDigits1)
	d2, bn2 := generateRandomNumber(integerDigits2, decimalDigits2)

	b.ResetTimer()

	// Decimal division
	for i := 0; i < b.N; i++ {
		d1.Div(d2)
	}

	// Bignum division
	for i := 0; i < b.N; i++ {
		bn1.Divide(bn2)
	}
}

// Run benchmarks for all combinations of integer and decimal digits
func runBenchmarks(b *testing.B, benchmarkFunc func(*testing.B)) {
	for integerDigits := 3; integerDigits <= maxIntegerDigits; integerDigits += 3 {
		for decimalDigits := 3; decimalDigits <= maxDecimalDigits; decimalDigits += 3 {
			b.Run(fmt.Sprintf("Integer%d_Decimal%d", integerDigits, decimalDigits), func(b *testing.B) {
				benchmarkFunc(b)
			})
		}
	}
}

// Test cases for positive numbers
func BenchmarkPositive(b *testing.B) {
	b.Run("Addition", func(b *testing.B) {
		runBenchmarks(b, BenchmarkAddition)
	})

	b.Run("Multiplication", func(b *testing.B) {
		runBenchmarks(b, BenchmarkMultiplication)
	})

	b.Run("Division", func(b *testing.B) {
		runBenchmarks(b, BenchmarkDivision)
	})
}

// Test cases for negative numbers
func BenchmarkNegative(b *testing.B) {
	b.Run("Addition", func(b *testing.B) {
		for integerDigits := 3; integerDigits <= maxIntegerDigits; integerDigits += 3 {
			for decimalDigits := 3; decimalDigits <= maxDecimalDigits; decimalDigits += 3 {
				b.Run(fmt.Sprintf("Integer%d_Decimal%d", integerDigits, decimalDigits), func(b *testing.B) {
					d1, bn1 := generateRandomNumber(integerDigits, decimalDigits)
					d2, bn2 := generateRandomNumber(integerDigits, decimalDigits)

					d1 = d1.Neg()             // Make d1 negative
					bn1 = bn1.AbsoluteValue() // Make bn1 negative (bignum doesn't have a Neg function, so we use Abs)

					// Decimal addition
					for i := 0; i < b.N; i++ {
						d1.Add(d2)
					}

					// Bignum addition
					for i := 0; i < b.N; i++ {
						bn1.Add(bn2)
					}
				})
			}
		}
	})

	b.Run("Multiplication", func(b *testing.B) {
		for integerDigits := 3; integerDigits <= maxIntegerDigits; integerDigits += 3 {
			for decimalDigits := 3; decimalDigits <= maxDecimalDigits; decimalDigits += 3 {
				b.Run(fmt.Sprintf("Integer%d_Decimal%d", integerDigits, decimalDigits), func(b *testing.B) {
					d1, bn1 := generateRandomNumber(integerDigits, decimalDigits)
					d2, bn2 := generateRandomNumber(integerDigits, decimalDigits)

					d1 = d1.Neg()             // Make d1 negative
					bn1 = bn1.AbsoluteValue() // Make bn1 negative (bignum doesn't have a Neg function, so we use Abs)

					// Decimal multiplication
					for i := 0; i < b.N; i++ {
						d1.Mul(d2)
					}

					// Bignum multiplication
					for i := 0; i < b.N; i++ {
						bn1.Multiply(bn2)
					}
				})
			}
		}
	})

	b.Run("Division", func(b *testing.B) {
		for integerDigits := 3; integerDigits <= maxIntegerDigits; integerDigits += 3 {
			for decimalDigits := 3; decimalDigits <= maxDecimalDigits; decimalDigits += 3 {
				b.Run(fmt.Sprintf("Integer%d_Decimal%d", integerDigits, decimalDigits), func(b *testing.B) {
					d1, bn1 := generateRandomNumber(integerDigits, decimalDigits)
					d2, bn2 := generateRandomNumber(integerDigits, decimalDigits)

					d1 = d1.Neg()             // Make d1 negative
					bn1 = bn1.AbsoluteValue() // Make bn1 negative (bignum doesn't have a Neg function, so we use Abs)

					// Decimal division
					for i := 0; i < b.N; i++ {
						d1.Div(d2)
					}

					// Bignum division
					for i := 0; i < b.N; i++ {
						bn1.Divide(bn2)
					}
				})
			}
		}
	})
}

// Test cases for mixed positive and negative numbers
func BenchmarkMixed(b *testing.B) {
	b.Run("Addition", func(b *testing.B) {
		for integerDigits1 := 3; integerDigits1 <= maxIntegerDigits; integerDigits1 += 3 {
			for decimalDigits1 := 3; decimalDigits1 <= maxDecimalDigits; decimalDigits1 += 3 {
				for integerDigits2 := 3; integerDigits2 <= maxIntegerDigits; integerDigits2 += 3 {
					for decimalDigits2 := 3; decimalDigits2 <= maxDecimalDigits; decimalDigits2 += 3 {
						b.Run(fmt.Sprintf("Integer%d_Decimal%d_vs_Integer%d_Decimal%d", integerDigits1, decimalDigits1, integerDigits2, decimalDigits2), func(b *testing.B) {
							d1, bn1 := generateRandomNumber(integerDigits1, decimalDigits1)
							d2, bn2 := generateRandomNumber(integerDigits2, decimalDigits2)

							// Randomly choose to negate one of the numbers
							if rand.Intn(2) == 0 {
								d1 = d1.Neg()
								bn1 = bn1.AbsoluteValue() // Make bn1 negative (bignum doesn't have a Neg function, so we use Abs)
							} else {
								d2 = d2.Neg()
								bn2 = bn2.AbsoluteValue() // Make bn2 negative (bignum doesn't have a Neg function, so we use Abs)
							}

							// Decimal addition
							for i := 0; i < b.N; i++ {
								d1.Add(d2)
							}

							// Bignum addition
							for i := 0; i < b.N; i++ {
								bn1.Add(bn2)
							}
						})
					}
				}
			}
		}
	})

	b.Run("Multiplication", func(b *testing.B) {
		for integerDigits1 := 3; integerDigits1 <= maxIntegerDigits; integerDigits1 += 3 {
			for decimalDigits1 := 3; decimalDigits1 <= maxDecimalDigits; decimalDigits1 += 3 {
				for integerDigits2 := 3; integerDigits2 <= maxIntegerDigits; integerDigits2 += 3 {
					for decimalDigits2 := 3; decimalDigits2 <= maxDecimalDigits; decimalDigits2 += 3 {
						b.Run(fmt.Sprintf("Integer%d_Decimal%d_vs_Integer%d_Decimal%d", integerDigits1, decimalDigits1, integerDigits2, decimalDigits2), func(b *testing.B) {
							d1, bn1 := generateRandomNumber(integerDigits1, decimalDigits1)
							d2, bn2 := generateRandomNumber(integerDigits2, decimalDigits2)

							// Randomly choose to negate one of the numbers
							if rand.Intn(2) == 0 {
								d1 = d1.Neg()
								bn1 = bn1.AbsoluteValue() // Make bn1 negative (bignum doesn't have a Neg function, so we use Abs)
							} else {
								d2 = d2.Neg()
								bn2 = bn2.AbsoluteValue() // Make bn2 negative (bignum doesn't have a Neg function, so we use Abs)
							}

							// Decimal multiplication
							for i := 0; i < b.N; i++ {
								d1.Mul(d2)
							}

							// Bignum multiplication
							for i := 0; i < b.N; i++ {
								bn1.Multiply(bn2)
							}
						})
					}
				}
			}
		}
	})

	b.Run("Division", func(b *testing.B) {
		for integerDigits1 := 3; integerDigits1 <= maxIntegerDigits; integerDigits1 += 3 {
			for decimalDigits1 := 3; decimalDigits1 <= maxDecimalDigits; decimalDigits1 += 3 {
				for integerDigits2 := 3; integerDigits2 <= maxIntegerDigits; integerDigits2 += 3 {
					for decimalDigits2 := 3; decimalDigits2 <= maxDecimalDigits; decimalDigits2 += 3 {
						b.Run(fmt.Sprintf("Integer%d_Decimal%d_vs_Integer%d_Decimal%d", integerDigits1, decimalDigits1, integerDigits2, decimalDigits2), func(b *testing.B) {
							d1, bn1 := generateRandomNumber(integerDigits1, decimalDigits1)
							d2, bn2 := generateRandomNumber(integerDigits2, decimalDigits2)

							// Randomly choose to negate one of the numbers
							if rand.Intn(2) == 0 {
								d1 = d1.Neg()
								bn1 = bn1.AbsoluteValue() // Make bn1 negative (bignum doesn't have a Neg function, so we use Abs)
							} else {
								d2 = d2.Neg()
								bn2 = bn2.AbsoluteValue() // Make bn2 negative (bignum doesn't have a Neg function, so we use Abs)
							}

							// Decimal division
							for i := 0; i < b.N; i++ {
								d1.Div(d2)
							}

							// Bignum division
							for i := 0; i < b.N; i++ {
								bn1.Divide(bn2)
							}
						})
					}
				}
			}
		}
	})
}
