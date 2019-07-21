// Copyright 2017 ING Bank N.V.
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package byteconversion

import (
  "errors"
  "math/big"
)

/**
 * Returns the big.Int based on the passed byte-array assuming the byte-array contains 
 * the two's-complement representation of this big.Int. The byte array will be in big-endian byte-order: 
 * the most significant byte is in the zeroth element. The array will contain the minimum number of bytes 
 * required to represent this BigInteger, including at least one sign bit, which is (ceil((this.bitLength() + 1)/8)). 
 */
func FromByteArray(bytesIn [] byte) (*big.Int, error) {

	const MINUS_ONE = -1

	if len(bytesIn) == 0 {
		err := errors.New("Cannot convert empty array to big.Int.")
		return nil, err
	}

	highestByte := bytesIn[0]
	isNegative := (highestByte & 128) !=0
	var convertedBytes []byte

    if isNegative {

		tmpInt := new(big.Int).SetBytes(bytesIn)
		tmpInt = tmpInt.Sub(tmpInt, big.NewInt(1))
		tmpBytes := tmpInt.Bytes()

		if tmpBytes[0] == 255 {
			convertedBytes = FlipBytes(tmpBytes)[1:]
		} else {
			convertedBytes = tmpBytes
			copy(convertedBytes, FlipBytes(tmpBytes))
		}
		tmp := new(big.Int).SetBytes(convertedBytes)
		return tmp.Mul(tmp, big.NewInt(MINUS_ONE)), nil
    } else {
    	// if positive leave unchanged (additional 0-bytes will be ignored)
    	return new(big.Int).SetBytes(bytesIn), nil
    }
}

/**
 * Returns a byte array containing the two's-complement representation of this big.Int.
 * The byte array will be in big-endian byte-order: the most significant byte is in the 
 * zeroth element. The array will contain the minimum number of bytes required to represent 
 * this BigInteger, including at least one sign bit, which is (ceil((this.bitLength() + 1)/8)). 
 */
func ToByteArray(in *big.Int) []byte {

	isNegative := in.Cmp(new(big.Int).SetInt64(0)) < 0

	bytes := in.Bytes()
	length := len(bytes)

	if length == 0 {
		return []byte { 0 };
	}

	highestByte := bytes[0]
	var convertedBytes []byte

	if !isNegative {
		if (highestByte & 128) !=0 {

			convertedBytes = make([]byte, length + 1)
			convertedBytes[0] = 0
			copy(convertedBytes[1:], bytes)
			return convertedBytes
		} else {
			return bytes
		}
	} else {
		if (highestByte & 128) !=0 {

			convertedBytes = make([]byte, length + 1)
			convertedBytes[0] = 255
			copy(convertedBytes[1:], FlipBytes(bytes))
		} else {
			convertedBytes = FlipBytes(bytes)
		}

		convertedInt := new(big.Int).SetBytes(convertedBytes)
		convertedInt.Add(convertedInt, big.NewInt(1))
		return convertedInt.Bytes()
	}
}

/**
 * Flips all bytes in each of the array's elements.
 * Returns the flipped elements.
 */
func FlipBytes(bytesIn [] byte) []byte {

	length := len(bytesIn)
	flippedBytes := make([]byte, length)

	for i := 0; i < length; i++ {
		flippedBytes[i] = bytesIn[i] ^ 255;
	}
	return flippedBytes
}