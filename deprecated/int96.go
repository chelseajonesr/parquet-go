package deprecated

import (
	"math/big"
	"math/bits"
	"time"
	"unsafe"
)

// Int96 is an implementation of the deprecated INT96 parquet type.
type Int96 [3]uint32

// Int32ToInt96 converts a int32 value to a Int96.
func Int32ToInt96(value int32) (i96 Int96) {
	if value < 0 {
		i96[2] = 0xFFFFFFFF
		i96[1] = 0xFFFFFFFF
	}
	i96[0] = uint32(value)
	return
}

// Int64ToInt96 converts a int64 value to Int96.
func Int64ToInt96(value int64) (i96 Int96) {
	if value < 0 {
		i96[2] = 0xFFFFFFFF
	}
	i96[1] = uint32(value >> 32)
	i96[0] = uint32(value)
	return
}

// IsZero returns true if i is the zero-value.
func (i Int96) IsZero() bool { return i == Int96{} }

// Negative returns true if i is a negative value.
func (i Int96) Negative() bool {
	return (i[2] >> 31) != 0
}

// Less returns true if i < j.
//
// The method implements a signed comparison between the two operands.
func (i Int96) Less(j Int96) bool {
	if i.Negative() {
		if !j.Negative() {
			return true
		}
	} else {
		if j.Negative() {
			return false
		}
	}
	for k := 2; k >= 0; k-- {
		a, b := i[k], j[k]
		switch {
		case a < b:
			return true
		case a > b:
			return false
		}
	}
	return false
}

// Int converts i to a big.Int representation.
func (i Int96) Int() *big.Int {
	z := new(big.Int)
	z.Or(z, big.NewInt(int64(i[2])<<32|int64(i[1])))
	z.Lsh(z, 32)
	z.Or(z, big.NewInt(int64(i[0])))
	return z
}

// Int32 converts i to a int32, potentially truncating the value.
func (i Int96) Int32() int32 {
	return int32(i[0])
}

// Int64 converts i to a int64, potentially truncating the value.
func (i Int96) Int64() int64 {
	return int64(i[1])<<32 | int64(i[0])
}

const (
	EPOCH_AS_JULIAN_DAY        = 2440588
	NANOS_PER_DAY       uint64 = 86400 * 1000 * 1000 * 1000
)

// Time converts i to a Time
// More information on the use of Int96 to store timestamps can be found here:
// https://github.com/apache/parquet-format/pull/49
func (i Int96) Time() time.Time {
	nanos := uint64(i[1])<<32 | uint64(i[0])
	days := i[2]

	totalNanos := uint64(days-EPOCH_AS_JULIAN_DAY)*NANOS_PER_DAY + nanos
	return time.Unix(0, int64(totalNanos)).UTC()
}

// String returns a string representation of i.
func (i Int96) String() string {
	return i.Int().String()
}

// Len returns the minimum length in bits required to store the value of i.
func (i Int96) Len() int {
	switch {
	case i[2] != 0:
		return 64 + bits.Len32(i[2])
	case i[1] != 0:
		return 32 + bits.Len32(i[1])
	default:
		return bits.Len32(i[0])
	}
}

// Int96ToBytes converts the slice of Int96 values to a slice of bytes sharing
// the same backing array.
func Int96ToBytes(data []Int96) []byte {
	return unsafe.Slice(*(**byte)(unsafe.Pointer(&data)), 12*len(data))
}

// BytesToInt96 converts the byte slice passed as argument to a slice of Int96
// sharing the same backing array.
//
// When the number of bytes in the input is not a multiple of 12, the function
// truncates it in the returned slice.
func BytesToInt96(data []byte) []Int96 {
	return unsafe.Slice(*(**Int96)(unsafe.Pointer(&data)), len(data)/12)
}

func MaxLenInt96(data []Int96) int {
	max := 0
	for i := range data {
		n := data[i].Len()
		if n > max {
			max = n
		}
	}
	return max
}

func MinInt96(data []Int96) (min Int96) {
	if len(data) > 0 {
		min = data[0]
		for _, v := range data[1:] {
			if v.Less(min) {
				min = v
			}
		}
	}
	return min
}

func MaxInt96(data []Int96) (max Int96) {
	if len(data) > 0 {
		max = data[0]
		for _, v := range data[1:] {
			if max.Less(v) {
				max = v
			}
		}
	}
	return max
}

func MinMaxInt96(data []Int96) (min, max Int96) {
	if len(data) > 0 {
		min = data[0]
		max = data[0]
		for _, v := range data[1:] {
			if v.Less(min) {
				min = v
			}
			if max.Less(v) {
				max = v
			}
		}
	}
	return min, max
}

func OrderOfInt96(data []Int96) int {
	if len(data) > 1 {
		if int96AreInAscendingOrder(data) {
			return +1
		}
		if int96AreInDescendingOrder(data) {
			return -1
		}
	}
	return 0
}

func int96AreInAscendingOrder(data []Int96) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i].Less(data[i-1]) {
			return false
		}
	}
	return true
}

func int96AreInDescendingOrder(data []Int96) bool {
	for i := len(data) - 1; i > 0; i-- {
		if data[i-1].Less(data[i]) {
			return false
		}
	}
	return true
}
