package bignum

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestAbsoluteValue(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("123.45", 2, RoundToNearest)
		result := bn.AbsoluteValue()
		expected, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("-123.45", 2, RoundToNearest)
		result := bn.AbsoluteValue()
		expected, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Zero", func(t *testing.T) {
		bn, _ := NewBigNumber("0", 2, RoundToNearest)
		result := bn.AbsoluteValue()
		expected, _ := NewBigNumber("0", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		result := bn.AbsoluteValue()
		expected, _ := NewBigNumber("inf", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		result := bn.AbsoluteValue()
		expected, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})
}

func TestString(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if bn.String() != "123.45" {
			t.Errorf("Expected 123.45, got %s", bn.String())
		}
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("-123.45", 2, RoundToNearest)
		if bn.String() != "-123.45" {
			t.Errorf("Expected -123.45, got %s", bn.String())
		}
	})

	t.Run("Zero", func(t *testing.T) {
		bn, _ := NewBigNumber("0", 2, RoundToNearest)
		if bn.String() != "0" {
			t.Errorf("Expected 0, got %s", bn.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		if bn.String() != "Infinity" {
			t.Errorf("Expected Infinity, got %s", bn.String())
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if bn.String() != "NaN" {
			t.Errorf("Expected NaN, got %s", bn.String())
		}
	})
}

func TestScientificNotation(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("1234567890.1234567890", 10, RoundToNearest)
		if bn.ScientificNotation() != "1.2345678901234568e+09" {
			t.Errorf("Expected 1.2345678901234568e+09, got %s", bn.ScientificNotation())
		}
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("-1234567890.1234567890", 10, RoundToNearest)
		if bn.ScientificNotation() != "-1.2345678901234568e+09" {
			t.Errorf("Expected -1.2345678901234568e+09, got %s", bn.ScientificNotation())
		}
	})

	t.Run("Zero", func(t *testing.T) {
		bn, _ := NewBigNumber("0", 2, RoundToNearest)
		if bn.ScientificNotation() != "0" {
			t.Errorf("Expected 0, got %s", bn.ScientificNotation())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		if bn.ScientificNotation() != "Infinity" {
			t.Errorf("Expected Infinity, got %s", bn.ScientificNotation())
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if bn.ScientificNotation() != "NaN" {
			t.Errorf("Expected NaN, got %s", bn.ScientificNotation())
		}
	})
}

func TestIsZero(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		bn, _ := NewBigNumber("0", 2, RoundToNearest)
		if !bn.IsZero() {
			t.Errorf("Expected true for IsZero, got false")
		}
	})

	t.Run("NonZero", func(t *testing.T) {
		bn, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if bn.IsZero() {
			t.Errorf("Expected false for IsZero, got true")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		if bn.IsZero() {
			t.Errorf("Expected false for IsZero, got true")
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if bn.IsZero() {
			t.Errorf("Expected false for IsZero, got true")
		}
	})
}

func TestEqual(t *testing.T) {
	t.Run("EqualNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if !bn1.Equal(bn2) {
			t.Errorf("Expected true for Equal, got false")
		}
	})

	t.Run("DifferentNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if bn1.Equal(bn2) {
			t.Errorf("Expected false for Equal, got true")
		}
	})

	t.Run("DifferentPrecisions", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.450", 3, RoundToNearest)
		if bn1.Equal(bn2) {
			t.Errorf("Expected false for Equal, got true")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("inf", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		if !bn1.Equal(bn2) {
			t.Errorf("Expected true for Equal, got false")
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("NaN", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !bn1.Equal(bn2) {
			t.Errorf("Expected true for Equal, got false")
		}
	})
}

func TestLessThan(t *testing.T) {
	t.Run("SmallerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if !bn1.LessThan(bn2) {
			t.Errorf("Expected true for LessThan, got false")
		}
	})

	t.Run("LargerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if bn1.LessThan(bn2) {
			t.Errorf("Expected false for LessThan, got true")
		}
	})

	t.Run("EqualNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if bn1.LessThan(bn2) {
			t.Errorf("Expected false for LessThan, got true")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		if !bn1.LessThan(bn2) {
			t.Errorf("Expected true for LessThan, got false")
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !bn1.LessThan(bn2) {
			t.Errorf("Expected true for LessThan, got false")
		}
	})
}

func TestGreaterThan(t *testing.T) {
	t.Run("LargerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if !bn1.GreaterThan(bn2) {
			t.Errorf("Expected true for GreaterThan, got false")
		}
	})

	t.Run("SmallerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if bn1.GreaterThan(bn2) {
			t.Errorf("Expected false for GreaterThan, got true")
		}
	})

	t.Run("EqualNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if bn1.GreaterThan(bn2) {
			t.Errorf("Expected false for GreaterThan, got true")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("inf", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if !bn1.GreaterThan(bn2) {
			t.Errorf("Expected true for GreaterThan, got false")
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !bn1.GreaterThan(bn2) {
			t.Errorf("Expected true for GreaterThan, got false")
		}
	})
}

func TestLessOrEqual(t *testing.T) {
	t.Run("SmallerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if !bn1.LessOrEqual(bn2) {
			t.Errorf("Expected true for LessOrEqual, got false")
		}
	})

	t.Run("LargerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if bn1.LessOrEqual(bn2) {
			t.Errorf("Expected false for LessOrEqual, got true")
		}
	})

	t.Run("EqualNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if !bn1.LessOrEqual(bn2) {
			t.Errorf("Expected true for LessOrEqual, got false")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		if !bn1.LessOrEqual(bn2) {
			t.Errorf("Expected true for LessOrEqual, got false")
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !bn1.LessOrEqual(bn2) {
			t.Errorf("Expected true for LessOrEqual, got false")
		}
	})
}

func TestGreaterOrEqual(t *testing.T) {
	t.Run("LargerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if !bn1.GreaterOrEqual(bn2) {
			t.Errorf("Expected true for GreaterOrEqual, got false")
		}
	})

	t.Run("SmallerNumber", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		if bn1.GreaterOrEqual(bn2) {
			t.Errorf("Expected false for GreaterOrEqual, got true")
		}
	})

	t.Run("EqualNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if !bn1.GreaterOrEqual(bn2) {
			t.Errorf("Expected true for GreaterOrEqual, got false")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("inf", 2, RoundToNearest)
		bn2, _ := NewBigNumber("123.45", 2, RoundToNearest)
		if !bn1.GreaterOrEqual(bn2) {
			t.Errorf("Expected true for GreaterOrEqual, got false")
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !bn1.GreaterOrEqual(bn2) {
			t.Errorf("Expected true for GreaterOrEqual, got false")
		}
	})
}

func TestRound(t *testing.T) {
	t.Run("RoundToNearest", func(t *testing.T) {
		bn, _ := NewBigNumber("123.456789", 5, RoundToNearest)
		rounded := bn.Round(2)
		expected, _ := NewBigNumber("123.46", 2, RoundToNearest)
		if !rounded.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), rounded.String())
		}
	})

	t.Run("RoundToEven", func(t *testing.T) {
		bn, _ := NewBigNumber("123.455", 3, RoundToEven)
		rounded := bn.Round(2)
		expected, _ := NewBigNumber("123.46", 2, RoundToEven)
		if !rounded.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), rounded.String())
		}
	})

	t.Run("RoundUp", func(t *testing.T) {
		bn, _ := NewBigNumber("123.456789", 5, RoundUp)
		rounded := bn.Round(2)
		expected, _ := NewBigNumber("123.46", 2, RoundUp)
		if !rounded.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), rounded.String())
		}
	})

	t.Run("RoundDown", func(t *testing.T) {
		bn, _ := NewBigNumber("123.456789", 5, RoundDown)
		rounded := bn.Round(2)
		expected, _ := NewBigNumber("123.45", 2, RoundDown)
		if !rounded.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), rounded.String())
		}
	})

	t.Run("SamePrecision", func(t *testing.T) {
		bn, _ := NewBigNumber("123.456789", 5, RoundToNearest)
		rounded := bn.Round(5)
		if !rounded.Equal(bn) {
			t.Errorf("Expected %s, got %s", bn.String(), rounded.String())
		}
	})
}

func TestToFloat(t *testing.T) {
	t.Run("ValidNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("123.45", 2, RoundToNearest)
		floatVal, err := bn.toFloat()
		if err != nil {
			t.Errorf("Error converting to float: %v", err)
		}
		if floatVal != 123.45 {
			t.Errorf("Expected 123.45, got %f", floatVal)
		}
	})

	t.Run("LargeNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("1e+308", 2, RoundToNearest)
		_, err := bn.toFloat()
		if err == nil {
			t.Error("Expected error for large number, got nil")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		floatVal, err := bn.toFloat()
		if err != nil {
			t.Errorf("Error converting to float: %v", err)
		}
		if math.IsInf(floatVal, 1) {
			t.Errorf("Expected positive infinity, got %f", floatVal)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		floatVal, err := bn.toFloat()
		if err != nil {
			t.Errorf("Error converting to float: %v", err)
		}
		if !math.IsNaN(floatVal) {
			t.Errorf("Expected NaN, got %f", floatVal)
		}
	})
}

func TestApplyRounding(t *testing.T) {
	t.Run("RoundToNearest", func(t *testing.T) {
		bn := &BigNumber{precision: 2, rounding: RoundToNearest}
		value := new(big.Int).Set(big.NewInt(12345))
		rounded := bn.applyRounding(value)
		expected := new(big.Int).Set(big.NewInt(12345))
		expected.Div(expected, big.NewInt(100))
		if rounded.Cmp(expected) != 0 {
			t.Errorf("Expected %s, got %s", expected.String(), rounded.String())
		}
	})

	t.Run("RoundToEven", func(t *testing.T) {
		bn := &BigNumber{precision: 2, rounding: RoundToEven}
		value := new(big.Int).Set(big.NewInt(12345))
		rounded := bn.applyRounding(value)
		expected := new(big.Int).Set(big.NewInt(12346))
		expected.Div(expected, big.NewInt(100))
		if rounded.Cmp(expected) != 0 {
			t.Errorf("Expected %s, got %s", expected.String(), rounded.String())
		}
	})
}

func TestScaleForPrecision(t *testing.T) {
	bn := &BigNumber{precision: 2}
	scaleFactor := bn.scaleForPrecision()
	if scaleFactor.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected scale factor 100, got %s", scaleFactor.String())
	}
}

func TestNewBigNumber(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		bn, err := NewBigNumber("123.45", 2, RoundToNearest)
		if err != nil {
			t.Errorf("Error creating BigNumber: %v", err)
		}
		abn, err := NewBigNumber("123.45", 2, RoundToNearest)
		if !bn.Equal(abn) {
			t.Errorf("Expected %s, got %s", "123.45", bn.String())
		}
	})

	t.Run("EmptyInput", func(t *testing.T) {
		_, err := NewBigNumber("", 2, RoundToNearest)
		if err == nil {
			t.Error("Expected error for empty string, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("InvalidIntegerPart", func(t *testing.T) {
		_, err := NewBigNumber("abc", 2, RoundToNearest)
		if err == nil {
			t.Error("Expected error for invalid integer part, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, err := NewBigNumber("inf", 2, RoundToNearest)
		if err != nil {
			t.Errorf("Error creating BigNumber: %v", err)
		}
		if !bn.isInf {
			t.Error("Expected BigNumber to be infinity")
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, err := NewBigNumber("NaN", 2, RoundToNearest)
		if err != nil {
			t.Errorf("Error creating BigNumber: %v", err)
		}
		if !bn.isNan {
			t.Error("Expected BigNumber to be NaN")
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("PositiveNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		result, _ := bn1.Add(bn2)
		expected, _ := NewBigNumber("191.34", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("-123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Add(bn2)
		expected, _ := NewBigNumber("-191.34", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("MixedSigns", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Add(bn2)
		expected, _ := NewBigNumber("55.56", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("DifferentPrecisions", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.890", 3, RoundToNearest)
		_, err := bn1.Add(bn2)
		if err == nil {
			t.Error("Expected error for different precisions, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		_, err := bn1.Add(bn2)
		if err == nil {
			t.Error("Expected error for adding with infinity, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn1.Add(bn2)
		if err == nil {
			t.Error("Expected error for adding with NaN, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})
}

func TestSubtract(t *testing.T) {
	t.Run("PositiveNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		result, _ := bn1.Subtract(bn2)
		expected, _ := NewBigNumber("55.56", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("-123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Subtract(bn2)
		expected, _ := NewBigNumber("-55.56", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("MixedSigns", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Subtract(bn2)
		expected, _ := NewBigNumber("191.34", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("DifferentPrecisions", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.890", 3, RoundToNearest)
		_, err := bn1.Subtract(bn2)
		if err == nil {
			t.Error("Expected error for different precisions, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		_, err := bn1.Subtract(bn2)
		if err == nil {
			t.Error("Expected error for subtracting with infinity, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn1.Subtract(bn2)
		if err == nil {
			t.Error("Expected error for subtracting with NaN, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})
}

func TestMultiply(t *testing.T) {
	t.Run("PositiveNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		result, _ := bn1.Multiply(bn2)
		expected, _ := NewBigNumber("8388.60", 4, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("-123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Multiply(bn2)
		expected, _ := NewBigNumber("8388.60", 4, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("MixedSigns", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Multiply(bn2)
		expected, _ := NewBigNumber("-8388.60", 4, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("DifferentPrecisions", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.890", 3, RoundToNearest)
		result, _ := bn1.Multiply(bn2)
		expected, _ := NewBigNumber("8388.6065", 5, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		_, err := bn1.Multiply(bn2)
		if err == nil {
			t.Error("Expected error for multiplying with infinity, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn1.Multiply(bn2)
		if err == nil {
			t.Error("Expected error for multiplying with NaN, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})
}

func TestDivide(t *testing.T) {
	t.Run("PositiveNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		result, _ := bn1.Divide(bn2)
		expected, _ := NewBigNumber("1.81", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("-123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Divide(bn2)
		expected, _ := NewBigNumber("1.81", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("MixedSigns", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Divide(bn2)
		expected, _ := NewBigNumber("-1.81", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("DifferentPrecisions", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.890", 3, RoundToNearest)
		result, _ := bn1.Divide(bn2)
		expected, _ := NewBigNumber("1.81", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("DivideByZero", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("0", 2, RoundToNearest)
		_, err := bn1.Divide(bn2)
		if err == nil {
			t.Error("Expected error for division by zero, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		result, _ := bn1.Divide(bn2)
		expected, _ := NewBigNumber("0", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn1.Divide(bn2)
		if err == nil {
			t.Error("Expected error for dividing with NaN, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})
}

func TestModulo(t *testing.T) {
	t.Run("PositiveNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.89", 2, RoundToNearest)
		result, _ := bn1.Modulo(bn2)
		expected, _ := NewBigNumber("55.56", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeNumbers", func(t *testing.T) {
		bn1, _ := NewBigNumber("-123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Modulo(bn2)
		expected, _ := NewBigNumber("-55.56", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("MixedSigns", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("-67.89", 2, RoundToNearest)
		result, _ := bn1.Modulo(bn2)
		expected, _ := NewBigNumber("55.56", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("DifferentPrecisions", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("67.890", 3, RoundToNearest)
		result, _ := bn1.Modulo(bn2)
		expected, _ := NewBigNumber("55.56", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("ModuloByZero", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("0", 2, RoundToNearest)
		_, err := bn1.Modulo(bn2)
		if err == nil {
			t.Error("Expected error for modulo by zero, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("inf", 2, RoundToNearest)
		_, err := bn1.Modulo(bn2)
		if err == nil {
			t.Error("Expected error for modulo with infinity, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn1, _ := NewBigNumber("123.45", 2, RoundToNearest)
		bn2, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn1.Modulo(bn2)
		if err == nil {
			t.Error("Expected error for modulo with NaN, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})
}

func TestExponentiate(t *testing.T) {
	t.Run("PositiveExponent", func(t *testing.T) {
		bn, _ := NewBigNumber("2.5", 2, RoundToNearest)
		result, _ := bn.Exponentiate(3)
		expected, _ := NewBigNumber("15.63", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeExponent", func(t *testing.T) {
		bn, _ := NewBigNumber("2.5", 2, RoundToNearest)
		result, _ := bn.Exponentiate(-2)
		expected, _ := NewBigNumber("0.16", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("ZeroExponent", func(t *testing.T) {
		bn, _ := NewBigNumber("2.5", 2, RoundToNearest)
		result, _ := bn.Exponentiate(0)
		expected, _ := NewBigNumber("1", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Overflow", func(t *testing.T) {
		bn, _ := NewBigNumber("2.5", 2, RoundToNearest)
		_, err := bn.Exponentiate(1000)
		if err == nil {
			t.Error("Expected error for overflow, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})
}

func TestSquareRoot(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("9", 2, RoundToNearest)
		result, _ := bn.SquareRoot()
		expected, _ := NewBigNumber("3", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("-9", 2, RoundToNearest)
		_, err := bn.SquareRoot()
		if err == nil {
			t.Error("Expected error for square root of a negative number, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("Zero", func(t *testing.T) {
		bn, _ := NewBigNumber("0", 2, RoundToNearest)
		result, _ := bn.SquareRoot()
		expected, _ := NewBigNumber("0", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		result, _ := bn.SquareRoot()
		expected, _ := NewBigNumber("inf", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		result, _ := bn.SquareRoot()
		expected, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})
}

func TestSine(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		bn, _ := NewBigNumber("0.5", 10, RoundToNearest)
		result, _ := bn.Sine()
		expected, _ := NewBigNumber(fmt.Sprintf("%f", math.Sin(0.5)), 10, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		_, err := bn.Sine()
		if err == nil {
			t.Error("Expected error for sine of infinity, got nil")
		}
		if _, ok := err.(error); !ok {
			t.Errorf("Expected error, got %T", err)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn.Sine()
		if err == nil {
			t.Error("Expected error for sine of NaN, got nil")
		}
		if _, ok := err.(error); !ok {
			t.Errorf("Expected error, got %T", err)
		}
	})
}

func TestCosine(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		bn, _ := NewBigNumber("0.5", 10, RoundToNearest)
		result, _ := bn.Cosine()
		expected, _ := NewBigNumber(fmt.Sprintf("%f", math.Cos(0.5)), 10, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		_, err := bn.Cosine()
		if err == nil {
			t.Error("Expected error for cosine of infinity, got nil")
		}
		if _, ok := err.(error); !ok {
			t.Errorf("Expected error, got %T", err)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn.Cosine()
		if err == nil {
			t.Error("Expected error for cosine of NaN, got nil")
		}
		if _, ok := err.(error); !ok {
			t.Errorf("Expected error, got %T", err)
		}
	})
}

func TestTangent(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		bn, _ := NewBigNumber("0.5", 10, RoundToNearest)
		result, _ := bn.Tangent()
		expected, _ := NewBigNumber(fmt.Sprintf("%f", math.Tan(0.5)), 10, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		_, err := bn.Tangent()
		if err == nil {
			t.Error("Expected error for tangent of infinity, got nil")
		}
		if _, ok := err.(error); !ok {
			t.Errorf("Expected error, got %T", err)
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		_, err := bn.Tangent()
		if err == nil {
			t.Error("Expected error for tangent of NaN, got nil")
		}
		if _, ok := err.(error); !ok {
			t.Errorf("Expected error, got %T", err)
		}
	})
}

func TestLog(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("2.71828", 5, RoundToNearest)
		result, _ := bn.Log()
		expected, _ := NewBigNumber("1", 5, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Zero", func(t *testing.T) {
		bn, _ := NewBigNumber("0", 2, RoundToNearest)
		_, err := bn.Log()
		if err == nil {
			t.Error("Expected error for logarithm of zero, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("-2.71828", 5, RoundToNearest)
		_, err := bn.Log()
		if err == nil {
			t.Error("Expected error for logarithm of a negative number, got nil")
		}
		if _, ok := err.(BigNumberError); !ok {
			t.Errorf("Expected BigNumberError, got %T", err)
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		result, _ := bn.Log()
		expected, _ := NewBigNumber("inf", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		result, _ := bn.Log()
		expected, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})
}

func TestExp(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		bn, _ := NewBigNumber("1", 5, RoundToNearest)
		result, _ := bn.Exp()
		expected, _ := NewBigNumber("2.71828", 5, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Zero", func(t *testing.T) {
		bn, _ := NewBigNumber("0", 2, RoundToNearest)
		result, _ := bn.Exp()
		expected, _ := NewBigNumber("1", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		bn, _ := NewBigNumber("inf", 2, RoundToNearest)
		result, _ := bn.Exp()
		expected, _ := NewBigNumber("inf", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})

	t.Run("NaN", func(t *testing.T) {
		bn, _ := NewBigNumber("NaN", 2, RoundToNearest)
		result, _ := bn.Exp()
		expected, _ := NewBigNumber("NaN", 2, RoundToNearest)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, got %s", expected.String(), result.String())
		}
	})
}
