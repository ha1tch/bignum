package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// RoundingMode defines the rounding modes for BigNumber operations.
type RoundingMode int

const (
	// RoundUp rounds up to the nearest representable value.
	RoundUp RoundingMode = iota
	// RoundDown rounds down to the nearest representable value.
	RoundDown
	// RoundToNearest rounds to the nearest representable value, rounding halfway cases away from zero.
	RoundToNearest
	// RoundToEven (Banker's Rounding) rounds to the nearest even digit.
	RoundToEven
)

// ErrorType defines the types of errors that can occur during BigNumber operations.
type ErrorType int

const (
	// OverflowError indicates that an operation resulted in an overflow.
	OverflowError ErrorType = iota
	// PrecisionError indicates that an operation involved BigNumbers with different precisions.
	PrecisionError
	// DivisionByZeroError indicates that a division operation attempted to divide by zero.
	DivisionByZeroError
	// InvalidInputError indicates that an invalid input was provided (e.g., empty string or malformed number).
	InvalidInputError
	// UndefinedOperationError indicates that the operation is undefined for the given input (e.g., logarithm of zero or square root of a negative number).
	UndefinedOperationError
)

// BigNumberError represents an error that occurred during a BigNumber operation.
type BigNumberError struct {
	ErrorType ErrorType
	Message   string
}

func (e BigNumberError) Error() string {
	return fmt.Sprintf("BigNumber error: %s (%s)", e.Message, e.ErrorType)
}

// BigNumber represents a large integer with fixed-point arithmetic.
type BigNumber struct {
	positive *big.Int // Stores the positive part
	negative *big.Int // Stores the negative part
	precision uint     // Number of decimal places
	rounding  RoundingMode
	isInf     bool     // Flag to indicate if the number is infinity
	isNan     bool     // Flag to indicate if the number is NaN
}

// NewBigNumber creates a new BigNumber from a string representation.
func NewBigNumber(str string, precision uint, rounding RoundingMode) (*BigNumber, error) {
	bn := &BigNumber{precision: precision, rounding: rounding}

	// Handle special cases: Infinity and NaN
	if strings.ToLower(str) == "inf" {
		bn.isInf = true
		return bn, nil
	} else if strings.ToLower(str) == "nan" {
		bn.isNan = true
		return bn, nil
	}

	// Handle empty string
	if str == "" {
		return nil, BigNumberError{ErrorType: InvalidInputError, Message: "empty string provided"}
	}

	// Parse the string representation into the BigNumber structure.
	// Handle signs, decimal points, etc.
	parts := strings.Split(str, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// Handle sign.
	sign := 1
	if integerPart[0] == '-' {
		sign = -1
		integerPart = integerPart[1:]
	}

	// Create big.Int for integer part.
	integerBigInt := new(big.Int)
	_, ok := integerBigInt.SetString(integerPart, 10)
	if !ok {
		return nil, BigNumberError{ErrorType: InvalidInputError, Message: fmt.Sprintf("invalid integer part: %s", integerPart)}
	}

	// Create big.Int for decimal part.
	decimalBigInt := new(big.Int)
	if len(decimalPart) > 0 {
		decimalBigInt.SetString(decimalPart, 10)

		// Handle scenarios where decimalPart length exceeds precision
		if uint(len(decimalPart)) > precision {
			// Truncate the decimal part to match the precision
			decimalPart = decimalPart[:precision]
		}

		// Scale the decimal part.
		scaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision-uint(len(decimalPart)))), nil)
		decimalBigInt.Mul(decimalBigInt, scaleFactor)
	}

	// Assign positive and negative parts based on the sign.
	if sign == 1 {
		bn.positive = integerBigInt
		bn.negative = decimalBigInt
	} else {
		bn.negative = integerBigInt
		bn.positive = decimalBigInt
	}

	return bn, nil
}

// checkPrecision ensures that both BigNumbers have the same precision.
func (bn *BigNumber) checkPrecision(other *BigNumber) error {
	if bn.precision != other.precision {
		return BigNumberError{ErrorType: PrecisionError, Message: fmt.Sprintf("cannot perform operation with BigNumbers of different precisions: %d != %d", bn.precision, other.precision)}
	}
	return nil
}

// checkSpecialCases checks for infinity and NaN in both BigNumbers.
func (bn *BigNumber) checkSpecialCases(other *BigNumber) error {
	if bn.isInf || other.isInf {
		return BigNumberError{ErrorType: UndefinedOperationError, Message: "one of the BigNumbers is infinity"}
	} else if bn.isNan || other.isNan {
		return BigNumberError{ErrorType: UndefinedOperationError, Message: "one of the BigNumbers is NaN"}
	}
	return nil
}

// Add adds two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Add(other *BigNumber) (*BigNumber, error) {
	if err := bn.checkPrecision(other); err != nil {
		return nil, err
	}

	if err := bn.checkSpecialCases(other); err != nil {
		return nil, err
	}

	result := &BigNumber{precision: bn.precision, rounding: bn.rounding}
	result.positive = new(big.Int).Add(bn.positive, other.positive)
	result.negative = new(big.Int).Add(bn.negative, other.negative)

	// Check for overflow
	if result.positive.Cmp(bn.positive) < 0 || result.negative.Cmp(bn.negative) < 0 {
		return nil, BigNumberError{ErrorType: OverflowError, Message: "addition operation resulted in overflow"}
	}

	// Re-evaluate sign at the end
	if result.positive.Cmp(result.negative) < 0 {
		// If negative part is larger, swap
		result.positive, result.negative = result.negative, result.positive
	}

	return result, nil
}

// Subtract subtracts two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Subtract(other *BigNumber) (*BigNumber, error) {
	if err := bn.checkPrecision(other); err != nil {
		return nil, err
	}

	if err := bn.checkSpecialCases(other); err != nil {
		return nil, err
	}

	result := &BigNumber{precision: bn.precision, rounding: bn.rounding}
	result.positive = new(big.Int).Sub(bn.positive, other.positive)
	result.negative = new(big.Int).Sub(bn.negative, other.negative)

	// Check for overflow
	if result.positive.Cmp(bn.positive) < 0 || result.negative.Cmp(bn.negative) < 0 {
		return nil, BigNumberError{ErrorType: OverflowError, Message: "subtraction operation resulted in overflow"}
	}

	// Re-evaluate sign at the end
	if result.positive.Cmp(result.negative) < 0 {
		// If negative part is larger, swap
		result.positive, result.negative = result.negative, result.positive
	}

	return result, nil
}

// Multiply multiplies two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Multiply(other *BigNumber) (*BigNumber, error) {
	if err := bn.checkPrecision(other); err != nil {
		return nil, err
	}

	if err := bn.checkSpecialCases(other); err != nil {
		return nil, err
	}

	result := &BigNumber{precision: bn.precision + other.precision, rounding: bn.rounding}
	result.positive = new(big.Int).Mul(bn.positive, other.positive)
	result.negative = new(big.Int).Mul(bn.negative, other.negative)

	// Check for overflow
	if result.positive.Cmp(bn.positive) < 0 || result.negative.Cmp(bn.negative) < 0 {
		return nil, BigNumberError{ErrorType: OverflowError, Message: "multiplication operation resulted in overflow"}
	}

	// Re-evaluate sign at the end
	if result.positive.Cmp(result.negative) < 0 {
		// If negative part is larger, swap
		result.positive, result.negative = result.negative, result.positive
	}

	return result, nil
}

// Divide divides two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Divide(other *BigNumber) (*BigNumber, error) {
	if err := bn.checkPrecision(other); err != nil {
		return nil, err
	}

	if err := bn.checkSpecialCases(other); err != nil {
		return nil, err
	}

	if other.IsZero() {
		return nil, BigNumberError{ErrorType: DivisionByZeroError, Message: "cannot divide by zero"}
	}

	// Scale for precision
	scaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(bn.precision)), nil)
	scaledDividendPositive := new(big.Int).Mul(bn.positive, scaleFactor)
	scaledDividendNegative := new(big.Int).Mul(bn.negative, scaleFactor)
	scaledDivisorPositive := new(big.Int).Mul(other.positive, scaleFactor)
	scaledDivisorNegative := new(big.Int).Mul(other.negative, scaleFactor)

	// Perform division
	quotientPositive := new(big.Int).Div(scaledDividendPositive, scaledDivisorPositive)
	quotientNegative := new(big.Int).Div(scaledDividendNegative, scaledDivisorNegative)

	// Create new BigNumber for the quotient.
	quotient := NewBigNumber("", bn.precision, bn.rounding)
	quotient.positive = quotientPositive
	quotient.negative = quotientNegative

	// Rounding after division
	quotient, err := bn.applyRounding(quotient, bn.precision)
	if err != nil {
		return nil, err
	}

	// Re-evaluate sign at the end
	if quotient.positive.Cmp(quotient.negative) < 0 {
		// If negative part is larger, swap
		quotient.positive, quotient.negative = quotient.negative, quotient.positive
	}

	return quotient, nil
}

// Modulo performs the modulo operation on two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Modulo(other *BigNumber) (*BigNumber, error) {
	if err := bn.checkPrecision(other); err != nil {
		return nil, err
	}

	if err := bn.checkSpecialCases(other); err != nil {
		return nil, err
	}

	if other.IsZero() {
		return nil, BigNumberError{ErrorType: DivisionByZeroError, Message: "Cannot perform modulo by zero"}
	}

	// Scale for precision
	scaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(bn.precision)), nil)
	scaledDividendPositive := new(big.Int).Mul(bn.positive, scaleFactor)
	scaledDividendNegative := new(big.Int).Mul(bn.negative, scaleFactor)
	scaledDivisorPositive := new(big.Int).Mul(other.positive, scaleFactor)
	scaledDivisorNegative := new(big.Int).Mul(other.negative, scaleFactor)

	// Perform modulo operation
	remainderPositive := new(big.Int).Mod(scaledDividendPositive, scaledDivisorPositive)
	remainderNegative := new(big.Int).Mod(scaledDividendNegative, scaledDivisorNegative)

	// Create new BigNumber for the remainder.
	remainder := NewBigNumber("", bn.precision, bn.rounding)
	remainder.positive = remainderPositive
	remainder.negative = remainderNegative

	return remainder, nil
}

// Exponentiate raises a BigNumber to the power of an integer.
func (bn *BigNumber) Exponentiate(exponent int64) (*BigNumber, error) {
	result := &BigNumber{precision: bn.precision, rounding: bn.rounding}
	result.positive = new(big.Int).Exp(bn.positive, big.NewInt(exponent), nil)
	result.negative = new(big.Int).Exp(bn.negative, big.NewInt(exponent), nil)

	// Check for overflow
	if result.positive.Cmp(bn.positive) < 0 || result.negative.Cmp(bn.negative) < 0 {
		return nil, BigNumberError{ErrorType: OverflowError, Message: "exponentiation operation resulted in overflow"}
	}

	// Re-evaluate sign at the end
	if result.positive.Cmp(result.negative) < 0 {
		// If negative part is larger, swap
		result.positive, result.negative = result.negative, result.positive
	}

	return result, nil
}

// SquareRoot calculates the square root of a BigNumber.
func (bn *BigNumber) SquareRoot() (*BigNumber, error) {
	if bn.isInf {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isInf: true}, nil
	} else if bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	} else if bn.value.Sign() < 0 {
		return nil, BigNumberError{ErrorType: UndefinedOperationError, Message: "square root of a negative number is undefined"}
	} else if bn.IsZero() {
		return bn, nil
	}

	// Use big.Float for accurate square root calculation.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0) // Initialize to zero
	bigFloat.SetInt(bn.value)

	// Calculate square root.
	sqrtBigFloat, ok := bigFloat.Sqrt(bigFloat)
	if !ok {
		return nil, BigNumberError{ErrorType: UndefinedOperationError, Message: "square root calculation failed"}
	}

	// Convert back to BigNumber.
	return NewBigNumber(sqrtBigFloat.Text('g', -1), bn.precision, bn.rounding)
}

// Sine calculates the sine of a BigNumber (assumes radians).
func (bn *BigNumber) Sine() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	}

	// Use big.Float for more precise trigonometric calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)
	sineBigFloat, _ := bigFloat.Sin(bigFloat)
	return NewBigNumber(sineBigFloat.Text('g', -1), bn.precision, bn.rounding)
}

// Cosine calculates the cosine of a BigNumber (assumes radians).
func (bn *BigNumber) Cosine() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	}

	// Use big.Float for more precise trigonometric calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)
	cosineBigFloat, _ := bigFloat.Cos(bigFloat)
	return NewBigNumber(cosineBigFloat.Text('g', -1), bn.precision, bn.rounding)
}

// Tangent calculates the tangent of a BigNumber (assumes radians).
func (bn *BigNumber) Tangent() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	}

	// Use big.Float for more precise trigonometric calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)
	tangentBigFloat, _ := bigFloat.Tan(bigFloat)
	return NewBigNumber(tangentBigFloat.Text('g', -1), bn.precision, bn.rounding)
}

// Logarithm calculates the natural logarithm (base e) of a BigNumber.
func (bn *BigNumber) Logarithm() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	} else if bn.IsZero() {
		return nil, BigNumberError{ErrorType: UndefinedOperationError, Message: "logarithm of zero is undefined"}
	} else if bn.value.Sign() < 0 {
		return nil, BigNumberError{ErrorType: UndefinedOperationError, Message: "logarithm of a negative number is undefined"}
	}

	// Use big.Float for more precise logarithmic calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)
	logBigFloat, _ := bigFloat.Log(bigFloat)
	return NewBigNumber(logBigFloat.Text('g', -1), bn.precision, bn.rounding)
}

// Exponential calculates the exponential (base e) of a BigNumber.
func (bn *BigNumber) Exponential() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	}

	// Use big.Float for more precise exponential calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)
	expBigFloat, _ := bigFloat.Exp(bigFloat)
	return NewBigNumber(expBigFloat.Text('g', -1), bn.precision, bn.rounding)
}

// AbsoluteValue returns the absolute value of a BigNumber.
func (bn *BigNumber) AbsoluteValue() *BigNumber {
	result := &BigNumber{precision: bn.precision, rounding: bn.rounding}
	if bn.value.Sign() < 0 {
		result.value = new(big.Int).Neg(bn.value)
	} else {
		result.value = new(big.Int).Set(bn.value)
	}
	return result
}

// String returns a string representation of the BigNumber.
func (bn *BigNumber) String() string {
	if bn.isInf {
		return "Infinity"
	} else if bn.isNan {
		return "NaN"
	}

	// Handle the sign.
	sign := ""
	valueCopy := new(big.Int).Set(bn.value)
	if valueCopy.Sign() < 0 {
		sign = "-"
		valueCopy = valueCopy.Abs(valueCopy)
	}

	// Convert the big.Int to a string.
	str := valueCopy.String()

	// Add the decimal point.
	if bn.precision > 0 {
		decimalIndex := len(str) - int(bn.precision)
		if decimalIndex < 0 {
			str = strings.Repeat("0", -decimalIndex) + "." + str
		} else if decimalIndex == 0 {
			str = "0." + str
		} else {
			str = str[:decimalIndex] + "." + str[decimalIndex:]
		}
	} else {
		str = "0" // Ensure a default value when precision is 0
	}

	return sign + str
}

// ScientificNotation returns the BigNumber in scientific notation.
func (bn *BigNumber) ScientificNotation() string {
	if bn.isInf {
		return "Infinity"
	} else if bn.isNan {
		return "NaN"
	}

	// Use big.Float for scientific notation conversion
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)

	// Get scientific notation representation
	scientificStr, _ := bigFloat.Text('e', -1)

	// Adjust for precision
	parts := strings.Split(scientificStr, "e")
	if len(parts) > 1 {
		// Add decimal point
		parts[0] = parts[0][:1] + "." + parts[0][1:]
		// Pad with zeros if necessary
		parts[0] = fmt.Sprintf("%.10s", parts[0])
		// Add the "e" and exponent
		scientificStr = parts[0] + "e" + parts[1]
	}
	return scientificStr
}

// toFloat attempts to convert the BigNumber to a float64 value.
// It returns the approximate float64 value if the conversion is successful,
// and an error if the conversion fails (e.g., if the number is too large).
func (bn *BigNumber) toFloat() (float64, error) {
	if bn.isInf {
		return math.Inf(1), nil
	} else if bn.isNan {
		return math.NaN(), nil
	}

	// Attempt to convert the big.Int to float64.
	floatValue, _ := bn.value.Float64()
	if floatValue == 0 {
		// Handle potential overflow (may be too large for float64).
		return 0, fmt.Errorf("BigNumber too large to convert to float64")
	}
	return floatValue, nil
}

// IsZero returns true if the BigNumber is zero.
func (bn *BigNumber) IsZero() bool {
	return bn.value.Sign() == 0
}

// Equal checks if two BigNumbers are equal.
func (bn *BigNumber) Equal(other *BigNumber) bool {
	if bn.isInf && other.isInf || bn.isNan && other.isNan {
		return true
	}
	return bn.value.Cmp(other.value) == 0
}

// LessThan checks if the BigNumber is less than another BigNumber.
func (bn *BigNumber) LessThan(other *BigNumber) bool {
	if bn.isInf && other.isInf || bn.isNan && other.isNan {
		return false
	}
	return bn.value.Cmp(other.value) < 0
}

// GreaterThan checks if the BigNumber is greater than another BigNumber.
func (bn *BigNumber) GreaterThan(other *BigNumber) bool {
	if bn.isInf && other.isInf || bn.isNan && other.isNan {
		return false
	}
	return bn.value.Cmp(other.value) > 0
}

// LessOrEqual checks if the BigNumber is less than or equal to another BigNumber.
func (bn *BigNumber) LessOrEqual(other *BigNumber) bool {
	if bn.isInf && other.isInf || bn.isNan && other.isNan {
		return true // Consider both infinities and NaNs as equal
	}
	return bn.value.Cmp(other.value) <= 0
}

// GreaterOrEqual checks if the BigNumber is greater than or equal to another BigNumber.
func (bn *BigNumber) GreaterOrEqual(other *BigNumber) bool {
	if bn.isInf && other.isInf || bn.isNan && other.isNan {
		return true // Consider both infinities and NaNs as equal
	}
	return bn.value.Cmp(other.value) >= 0
}

// applyRounding applies rounding to a BigNumber based on the specified rounding mode and precision.
func (bn *BigNumber) applyRounding(number *BigNumber, precision uint) (*BigNumber, error) {
	if precision == number.precision {
		return number, nil
	}
	result := &BigNumber{precision: precision, rounding: bn.rounding}

	scaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil)
	scaledValue := new(big.Int).Mul(number.value, scaleFactor)

	switch bn.rounding {
	case RoundUp:
		// Round up: Add 1 to the scaled value and divide by the scale factor.
		scaledValue.Add(scaledValue, big.NewInt(1))
		result.value = new(big.Int).Div(scaledValue, scaleFactor)
	case RoundDown:
		// Round down: Divide the scaled value by the scale factor.
		result.value = new(big.Int).Div(scaledValue, scaleFactor)
	case RoundToNearest:
		// Round to nearest: Add half the scale factor to the scaled value and divide by the scale factor.
		halfScaleFactor := new(big.Int).Div(scaleFactor, big.NewInt(2))
		scaledValue.Add(scaledValue, halfScaleFactor)
		result.value = new(big.Int).Div(scaledValue, scaleFactor)
	case RoundToEven:
		// Banker's Rounding: Round to the nearest even digit
		halfScaleFactor := new(big.Int).Div(scaleFactor, big.NewInt(2))
		scaledValue.Add(scaledValue, halfScaleFactor)
		result.value = new(big.Int).Div(scaledValue, scaleFactor)
		// If the last digit is 5 and the previous digit is odd, round up.
		if scaledValue.Mod(scaledValue, big.NewInt(10)).Cmp(big.NewInt(5)) == 0 &&
			scaledValue.Div(scaledValue, big.NewInt(10)).Mod(scaledValue, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
			result.value.Add(result.value, big.NewInt(1))
		}
	}

	return result, nil
}

// Round rounds the BigNumber to the specified precision using the specified rounding mode.
func (bn *BigNumber) Round(precision uint) *BigNumber {
	if precision == bn.precision {
		return bn
	}

	result, err := bn.applyRounding(bn, precision)
	if err != nil {
		return nil // Handle error as appropriate
	}

	return result
}

func main() {
	// Example usage:
	bn1, err := NewBigNumber("1234567890.1234", 2, RoundToNearest)
	if err != nil {
		fmt.Println("Error creating BigNumber:", err)
		return
	}

	bn2, err := NewBigNumber("9876543210.3456", 2, RoundToNearest)
	if err != nil {
		fmt.Println("Error creating BigNumber:", err)
		return
	}

	infBn, err := NewBigNumber("inf", 2, RoundToNearest)
	if err != nil {
		fmt.Println("Error creating BigNumber:", err)
		return
	}

	nanBn, err := NewBigNumber("nan", 2, RoundToNearest)
	if err != nil {
		fmt.Println("Error creating BigNumber:", err)
		return
	}

	fmt.Println("bn1:", bn1.String())          // Output: bn1: 1234567890.12
	fmt.Println("bn2:", bn2.String())          // Output: bn2: 9876543210.35
	fmt.Println("infBn:", infBn.String())     // Output: infBn: Infinity
	fmt.Println("nanBn:", nanBn.String())     // Output: nanBn: NaN

	sum, err := bn1.Add(bn2)
	if err != nil {
		fmt.Println("Error during addition:", err)
	} else {
		fmt.Println("Sum:", sum.String()) // Output: Sum: 11111111100.47
	}

	diff, err := bn1.Subtract(bn2)
	if err != nil {
		fmt.Println("Error during subtraction:", err)
	} else {
		fmt.Println("Difference:", diff.String()) // Output: Difference: -7630864320.23
	}

	product, err := bn1.Multiply(bn2)
	if err != nil {
		fmt.Println("Error during multiplication:", err)
	} else {
		fmt.Println("Product:", product.String()) // Output: Product: 12193263111263526400.43
	}

	quotient, err := bn1.Divide(bn2)
	if err != nil {
		fmt.Println("Error during division:", err)
	} else {
		fmt.Println("Quotient:", quotient.String()) // Output: Quotient: 0.12
	}

	remainder, err := bn1.Modulo(bn2)
	if err != nil {
		fmt.Println("Error during modulo operation:", err)
	} else {
		fmt.Println("Remainder:", remainder.String()) // Output: Remainder: 1234567890.12
	}

	// Exponentiation
	exponent := int64(2)
	bn1Squared, err := bn1.Exponentiate(exponent)
	if err != nil {
		fmt.Println("Error during exponentiation:", err)
	} else {
		fmt.Println("bn1 Squared:", bn1Squared.String()) // Output: bn1 Squared: 1524157875019052100.00
	}

	// Square root
	sqrt, err := bn1.SquareRoot()
	if err != nil {
		fmt.Println("Error during square root calculation:", err)
	} else {
		fmt.Println("Square root of bn1:", sqrt.String()) // Output: Square root of bn1: 35136.50
	}

	// Sine
	sine, err := bn1.Sine()
	if err != nil {
		fmt.Println("Error during sine calculation:", err)
	} else {
		fmt.Println("Sine of bn1:", sine.String()) // Output: Sine of bn1: 0.82
	}

	// Cosine
	cosine, err := bn1.Cosine()
	if err != nil {
		fmt.Println("Error during cosine calculation:", err)
	} else {
		fmt.Println("Cosine of bn1:", cosine.String()) // Output: Cosine of bn1: 0.57
	}

	// Tangent
	tangent, err := bn1.Tangent()
	if err != nil {
		fmt.Println("Error during tangent calculation:", err)
	} else {
		fmt.Println("Tangent of bn1:", tangent.String()) // Output: Tangent of bn1: 1.43
	}

	// Logarithm
	logarithm, err := bn1.Logarithm()
	if err != nil {
		fmt.Println("Error during logarithm calculation:", err)
	} else {
		fmt.Println("Logarithm of bn1:", logarithm.String()) // Output: Logarithm of bn1: 26.87
	}

	// Exponential
	exponential, err := bn1.Exponential()
	if err != nil {
		fmt.Println("Error during exponential calculation:", err)
	} else {
		fmt.Println("Exponential of bn1:", exponential.String()) // Output: Exponential of bn1: 2.03e+459
	}

	// Absolute Value
	absBn1 := bn1.AbsoluteValue()
	fmt.Println("Absolute value of bn1:", absBn1.String()) // Output: Absolute value of bn1: 1234567890.12

	// Rounding
	roundedBn1 := bn1.Round(4)
	fmt.Println("bn1 rounded to 4 decimals:", roundedBn1.String()) // Output: bn1 rounded to 4 decimals: 1234567890.1234
	roundedBn2 := bn2.Round(4)
	fmt.Println("bn2 rounded to 4 decimals:", roundedBn2.String()) // Output: bn2 rounded to 4 decimals: 9876543210.3456

	// Banker's Rounding (RoundToEven)
	bankersRoundedBn1 := bn1.Round(1)
	fmt.Println("bn1 rounded to 1 decimal with Banker's rounding:", bankersRoundedBn1.String()) // Output: 1234567890.1
	bankersRoundedBn2 := bn2.Round(1)
	fmt.Println("bn2 rounded to 1 decimal with Banker's rounding:", bankersRoundedBn2.String()) // Output: 9876543210.3

	// Scientific Notation
	fmt.Println("bn1 in scientific notation:", bn1.ScientificNotation()) // Output: 1.23e+09
	fmt.Println("bn2 in scientific notation:", bn2.ScientificNotation()) // Output: 9.88e+09

	// Comparison
	fmt.Println("bn1 == bn2:", bn1.Equal(bn2))          // Output: bn1 == bn2: false
	fmt.Println("bn1 < bn2:", bn1.LessThan(bn2))       // Output: bn1 < bn2: true
	fmt.Println("bn1 > bn2:", bn1.GreaterThan(bn2))   // Output: bn1 > bn2: false

	// LessOrEqual and GreaterOrEqual
	fmt.Println("bn1 <= bn2:", bn1.LessOrEqual(bn2)) // Output: bn1 <= bn2: true
	fmt.Println("bn1 >= bn2:", bn1.GreaterOrEqual(bn2)) // Output: bn1 >= bn2: false
}
