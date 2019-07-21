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

package zkproofs

import (
	"testing"
	"math/big"
)

func TestKeyGen(t *testing.T) {
	kp, _ := keygen()
	signature, _ := sign(big.NewInt(42), kp.privk)	
	res, _ := verify(signature, big.NewInt(42), kp.pubk)
	if res != true {
		t.Errorf("Assert failure: expected true, actual: %t", res)
		t.Fail()
	}
}

