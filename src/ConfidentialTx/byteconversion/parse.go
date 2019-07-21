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
	"math/big"
	"errors"
)

var (
	errInvalidBigInteger = errors.New("invalid ASCII for big integer")
)

// Decodes a byte array (ASCII encoding of comma-separated integers) into an array of big integers
func ParseInput(in []byte) ([]*big.Int, error) {

	prevIndex := 0
	var output[] *big.Int

	for index, element := range in {

		if element == 44 {
		    newInt, err := ConvertToBigInt(in[prevIndex:index])
		    if err != nil {
		        return nil, err
		    }

			output = append(output, newInt)
			prevIndex = index + 1
		}
	}

    newInt, err := ConvertToBigInt(in[prevIndex:])
    if err != nil {
        return nil, err
    }
	return append(output, newInt), nil
}

// Decodes a byte array (ASCII encoding of a signed integer) into a big integer
func ConvertToBigInt(in []byte) (*big.Int, error) {

    // Validate
    for index, element := range in {
        if !((element >= 48 && element <= 57) || (index == 0 && element == 45)) {
            return nil, errInvalidBigInteger
        }
    }

	s := string(in)
	i := new(big.Int)
	i.SetString(s, 10)
	return i, nil
}
