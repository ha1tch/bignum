package bignum

import (
	"fmt"
	"math"
	"math/big"
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
	positive  *big.Int // Stores the positive part
	negative  *big.Int // Stores the negative part
	precision uint     // Number of decimal places
	rounding  RoundingMode
	isInf     bool     // Flag to indicate if the number is infinity
	isNan     bool     // Flag to indicate if the number is NaN
	value     *big.Int // Stores the actual big integer value
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

	// Combine positive and negative parts with the sign.
	bn.value = new(big.Int).Sub(bn.positive, bn.negative)

	return bn, nil
}

// checkPrecision ensures that both BigNumbers have the same precision.
func (bn *BigNumber) checkPrecision(other *BigNumber) error {
	if bn.precision != other.precision {
		return BigNumberError{ErrorType: PrecisionError, Message: fmt.Sprintf("cannot perform operation with BigNumbers of different precisions: %d != %d", bn.precision, other.precision)}
	}
	return nil
}

// checkSpecialCases checks for infinity and NaN in both BigNumbers and returns an error if found.
func checkSpecialCases(bn, other *BigNumber) error {
	if bn.isInf || other.isInf {
		return BigNumberError{ErrorType: UndefinedOperationError, Message: "one of the BigNumbers is infinity"}
	} else if bn.isNan || other.isNan {
		return BigNumberError{ErrorType: UndefinedOperationError, Message: "one of the BigNumbers is NaN"}
	} else if bn.IsZero() && other.IsZero() {
		// Both BigNumbers are zero, return an error as appropriate (e.g., for division)
		return BigNumberError{ErrorType: UndefinedOperationError, Message: "cannot perform operation with both BigNumbers being zero"}
	}
	return nil
}

// checkOperands checks if two BigNumbers have the same precision and handles special cases (Infinity and NaN).
func checkOperands(bn, other *BigNumber) error {
	if err := bn.checkPrecision(other); err != nil {
		return err
	}

	if err := checkSpecialCases(bn, other); err != nil {
		return err
	}

	return nil
}

// Add adds two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Add(other *BigNumber) (*BigNumber, error) {
	if err := checkOperands(bn, other); err != nil {
		return nil, err
	}

	result := &BigNumber{precision: bn.precision, rounding: bn.rounding}
	result.positive = new(big.Int).Set(bn.positive) // Copying positive part
	result.negative = new(big.Int).Set(bn.negative) // Copying negative part

	result.positive.Add(result.positive, other.positive)
	result.negative.Add(result.negative, other.negative)

	// Update the value
	result.value = new(big.Int).Sub(result.positive, result.negative)

	// Apply rounding
	result.value = bn.applyRounding(result.value)

	return result, nil
}

// Subtract subtracts two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Subtract(other *BigNumber) (*BigNumber, error) {
	if err := checkOperands(bn, other); err != nil {
		return nil, err
	}

	result := &BigNumber{precision: bn.precision, rounding: bn.rounding}
	result.positive = new(big.Int).Set(bn.positive) // Copying positive part
	result.negative = new(big.Int).Set(bn.negative) // Copying negative part

	result.positive.Sub(result.positive, other.positive)
	result.negative.Sub(result.negative, other.negative)

	// Update the value
	result.value = new(big.Int).Sub(result.positive, result.negative)

	// Apply rounding
	result.value = bn.applyRounding(result.value)

	return result, nil
}

// Multiply multiplies two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Multiply(other *BigNumber) (*BigNumber, error) {
	if err := checkOperands(bn, other); err != nil {
		return nil, err
	}

	result := &BigNumber{precision: bn.precision + other.precision, rounding: bn.rounding}
	result.positive = new(big.Int).Set(bn.positive) // Copying positive part
	result.negative = new(big.Int).Set(bn.negative) // Copying negative part

	result.positive.Mul(result.positive, other.positive)
	result.negative.Mul(result.negative, other.negative)

	// Update the value
	result.value = new(big.Int).Sub(result.positive, result.negative)

	// Apply rounding
	result.value = bn.applyRounding(result.value)

	return result, nil
}

// Divide divides two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Divide(other *BigNumber) (*BigNumber, error) {
	if err := bn.checkPrecision(other); err != nil {
		return nil, err
	}

	if err := checkSpecialCases(bn, other); err != nil {
		return nil, err
	}

	if other.IsZero() {
		return nil, BigNumberError{ErrorType: DivisionByZeroError, Message: "cannot divide by zero"}
	}

	// Scale for precision
	scaleFactor := bn.scaleForPrecision()
	scaledDividendPositive := new(big.Int).Mul(bn.positive, scaleFactor)
	scaledDividendNegative := new(big.Int).Mul(bn.negative, scaleFactor)
	scaledDivisorPositive := new(big.Int).Mul(other.positive, scaleFactor)
	scaledDivisorNegative := new(big.Int).Mul(other.negative, scaleFactor)

	// Perform division
	quotientPositive := new(big.Int).Div(scaledDividendPositive, scaledDivisorPositive)
	quotientNegative := new(big.Int).Div(scaledDividendNegative, scaledDivisorNegative)

	// Create new BigNumber for the quotient.
	quotient, _ := NewBigNumber("", bn.precision, bn.rounding)
	quotient.positive = quotientPositive
	quotient.negative = quotientNegative

	// Rounding after division
	quotient.value = new(big.Int).Sub(quotient.positive, quotient.negative) // Calculate the value
	quotient.value = bn.applyRounding(quotient.value)                       // Apply rounding

	// Re-evaluate sign at the end
	if quotient.positive.Cmp(quotient.negative) < 0 {
		// If negative part is larger, swap
		quotient.positive, quotient.negative = quotient.negative, quotient.positive
	}

	// Update the 'value' field based on the sign
	quotient.value = new(big.Int).Sub(quotient.positive, quotient.negative)

	return quotient, nil
}

// Modulo performs the modulo operation on two BigNumbers and returns a new BigNumber.
func (bn *BigNumber) Modulo(other *BigNumber) (*BigNumber, error) {
	if err := bn.checkPrecision(other); err != nil {
		return nil, err
	}

	if err := checkSpecialCases(bn, other); err != nil {
		return nil, err
	}

	if other.IsZero() {
		return nil, BigNumberError{ErrorType: DivisionByZeroError, Message: "Cannot perform modulo by zero"}
	}

	// Scale for precision
	scaleFactor := bn.scaleForPrecision()
	scaledDividendPositive := new(big.Int).Mul(bn.positive, scaleFactor)
	scaledDividendNegative := new(big.Int).Mul(bn.negative, scaleFactor)
	scaledDivisorPositive := new(big.Int).Mul(other.positive, scaleFactor)
	scaledDivisorNegative := new(big.Int).Mul(other.negative, scaleFactor)

	// Perform modulo operation
	remainderPositive := new(big.Int).Mod(scaledDividendPositive, scaledDivisorPositive)
	remainderNegative := new(big.Int).Mod(scaledDividendNegative, scaledDivisorNegative)

	// Create new BigNumber for the remainder.
	remainder, _ := NewBigNumber("", bn.precision, bn.rounding)
	remainder.positive = remainderPositive
	remainder.negative = remainderNegative

	// Update the 'value' field based on the sign
	remainder.value = new(big.Int).Sub(remainder.positive, remainder.negative)

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

	// Update the 'value' field based on the sign
	result.value = new(big.Int).Sub(result.positive, result.negative)

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
	sqrtBigFloat := bigFloat.Sqrt(bigFloat) // sqrtBigFloat is of type *big.Float

	// Convert back to BigNumber
	sqrtBn, err := NewBigNumber(sqrtBigFloat.Text('g', -1), bn.precision, bn.rounding)
	if err != nil {
		return nil, err
	}
	return sqrtBn, nil // Return the new BigNumber
}

// Sine calculates the sine of a BigNumber (assumes radians).
func (bn *BigNumber) Sine() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		//return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
		return nil, fmt.Errorf("cannot perform Sine operation: value is Infinity or NaN")
	}

	// Use big.Float for more precise trigonometric calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)

	// Convert to float64 for math.Sin, but check for errors
	floatVal, accuracy := bigFloat.Float64()
	if accuracy != big.Exact { // **Check for accuracy:**
		return nil, fmt.Errorf("cannot perform Sine operation: loss of precision during conversion to float64")
	}
	sine := math.Sin(floatVal) // Calculate sine using math.Sin

	// Convert back to *big.Float
	bigFloat.SetFloat64(sine)

	// Convert back to BigNumber
	sineBn, err := NewBigNumber(bigFloat.Text('g', -1), bn.precision, bn.rounding)
	if err != nil {
		return nil, err
	}
	return sineBn, nil // Return the new BigNumber
}

// Cosine calculates the cosine of a BigNumber (assumes radians).
func (bn *BigNumber) Cosine() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		// return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
		return nil, fmt.Errorf("cannot perform Cosine operation: value is Infinity or NaN")
	}

	// Use big.Float for more precise trigonometric calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)

	// Convert to float64 for math.Cos, but check for errors
	floatVal, accuracy := bigFloat.Float64()
	if accuracy != big.Exact { // **Check for accuracy:**
		return nil, fmt.Errorf("cannot perform Cosine operation: loss of precision during conversion to float64")
	}
	cosine := math.Cos(floatVal) // Calculate cosine using math.Cos

	// Convert back to *big.Float
	bigFloat.SetFloat64(cosine)

	// Convert back to BigNumber
	cosineBn, err := NewBigNumber(bigFloat.Text('g', -1), bn.precision, bn.rounding)
	if err != nil {
		return nil, err
	}
	return cosineBn, nil // Return the new BigNumber
}

// Tangent calculates the tangent of a BigNumber (assumes radians).
func (bn *BigNumber) Tangent() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		// return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
		return nil, fmt.Errorf("cannot perform Tangent operation: value is Infinity or NaN")
	}

	// Use big.Float for more precise trigonometric calculations.
	bigFloat := new(big.Float)
	bigFloat.SetFloat64(0)
	bigFloat.SetInt(bn.value)

	// Convert to float64 for math.Tan, but check for errors
	floatVal, accuracy := bigFloat.Float64()
	if accuracy != big.Exact { // **Check for accuracy:**
		return nil, fmt.Errorf("cannot perform Tangent operation: loss of precision during conversion to float64")
	}
	tangent := math.Tan(floatVal) // Calculate tangent using math.Tan

	// Convert back to *big.Float
	bigFloat.SetFloat64(tangent)

	// Convert back to BigNumber
	tangentBn, err := NewBigNumber(bigFloat.Text('g', -1), bn.precision, bn.rounding)
	if err != nil {
		return nil, err
	}
	return tangentBn, nil // Return the new BigNumber
}

// Log approximates the natural logarithm (base e) of a BigNumber using Newton's method.
func (bn *BigNumber) Log() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	} else if bn.IsZero() {
		return nil, BigNumberError{ErrorType: UndefinedOperationError, Message: "logarithm of zero is undefined"}
	} else if bn.value.Sign() < 0 {
		return nil, BigNumberError{ErrorType: UndefinedOperationError, Message: "logarithm of a negative number is undefined"}
	}

	// Convert BigNumber to big.Int for calculations
	xInt := new(big.Int).Set(bn.value)

	// Start with an initial guess, e.g., y = 1.0 (for big.Int, we use 1)
	yInt := big.NewInt(1)
	deltaInt := big.NewInt(0)
	thresholdInt := big.NewInt(1) // We'll use 1 as a simple threshold (can be adjusted)

	// Calculate e using Taylor series (from Exp function)
	expYInt, _ := bn.Exp()
	expYInt.value = bn.applyRounding(expYInt.value)

	for {
		// deltaInt = (expYInt.value - xInt) / expYInt.value
		deltaInt.Sub(expYInt.value, xInt)
		deltaInt.Div(deltaInt, expYInt.value)

		// yInt = yInt - deltaInt
		yInt.Sub(yInt, deltaInt)

		// Stop if deltaInt is smaller than thresholdInt (can be adjusted)
		if deltaInt.Cmp(thresholdInt) < 0 {
			break
		}

		// Recalculate expYInt for next iteration
		expYInt, _ = bn.Exp()
		expYInt.value = bn.applyRounding(expYInt.value)
	}

	// Create new BigNumber with the result and apply rounding
	result, _ := NewBigNumber(yInt.String(), bn.precision, bn.rounding)
	result.value = bn.applyRounding(result.value)

	return result, nil
}

// Exp approximates the exponential function (base e) of a BigNumber using Taylor series.
func (bn *BigNumber) Exp() (*BigNumber, error) {
	if bn.isInf || bn.isNan {
		return &BigNumber{precision: bn.precision, rounding: bn.rounding, isNan: true}, nil
	}

	// Convert BigNumber to big.Int for calculations
	xInt := new(big.Int).Set(bn.value)

	// Calculate Taylor series approximation
	resultInt := new(big.Int).SetInt64(1) // e^0 = 1
	termInt := new(big.Int).SetInt64(1)   // Current term in series (starts at 1)
	factorialInt := big.NewInt(1)         // Current factorial value

	i := 1
	for {
		// termInt *= xInt / i
		termInt.Mul(termInt, xInt)
		termInt.Div(termInt, big.NewInt(int64(i)))

		// resultInt += termInt
		resultInt.Add(resultInt, termInt)

		// Break if term is small enough to stop (close enough to precision)
		if termInt.Cmp(big.NewInt(1)) < 0 {
			break
		}

		// Update factorial for next iteration
		factorialInt.Mul(factorialInt, big.NewInt(int64(i)))
		i++
	}

	// Create new BigNumber with the result and apply rounding
	result, _ := NewBigNumber(resultInt.String(), bn.precision, bn.rounding)
	result.value = bn.applyRounding(result.value)

	return result, nil
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
	scientificStr := bigFloat.Text('e', -1) // scientificStr is of type string

	// Convert back to BigNumber (not necessary, but following the pattern)
	sciBn, err := NewBigNumber(scientificStr, bn.precision, bn.rounding)
	if err != nil {
		return "" // Handle error as appropriate
	}
	return sciBn.String() // Return the new BigNumber
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
func (bn *BigNumber) applyRounding(value *big.Int) *big.Int {
	// Rounding logic based on rounding mode
	switch bn.rounding {
	case RoundToNearest:
		// Round to nearest: Add half the scale factor to the scaled value and divide by the scale factor.
		halfScaleFactor := new(big.Int).Div(bn.scaleForPrecision(), big.NewInt(2))
		value.Add(value, halfScaleFactor)
		value.Div(value, bn.scaleForPrecision())
	case RoundToEven:
		// Banker's Rounding: Round to the nearest even digit
		halfScaleFactor := new(big.Int).Div(bn.scaleForPrecision(), big.NewInt(2))
		value.Add(value, halfScaleFactor)
		value.Div(value, bn.scaleForPrecision())
		// If the last digit is 5 and the previous digit is odd, round up.
		if value.Mod(value, big.NewInt(10)).Cmp(big.NewInt(5)) == 0 &&
			value.Div(value, big.NewInt(10)).Mod(value, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
			value.Add(value, big.NewInt(1))
		}
	}

	return value
}

// scaleForPrecision returns a big.Int representing the scale factor for the specified precision.
func (bn *BigNumber) scaleForPrecision() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(bn.precision)), nil)
}

// Round rounds the BigNumber to the specified precision using the specified rounding mode.
func (bn *BigNumber) Round(precision uint) *BigNumber {
	if precision == bn.precision {
		return bn
	}

	result := &BigNumber{precision: precision, rounding: bn.rounding}
	result.value = new(big.Int).Set(bn.value) // Copy the value

	// Apply rounding to the copied value
	result.value = bn.applyRounding(result.value)

	return result
}

