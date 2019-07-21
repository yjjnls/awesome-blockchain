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


func TestCalculateHash(t * testing.T) {

	a := GetBigInt("20905485153255974750600830283139712767405035066172127447413526262122898097752829902691919420016794244099612526431387099905077116995490485444167190551980224865082320241670546533063409921082864323224863076823319894932240914571396941354556281385023649535909639921239646795929610627460276589386330363348840105387073757406261480377763345436612442076323102518362946991582624513737241437269968051355243751819094759669539075841991633425362795570590507959822047022497500292880734028347273355847985904992235033659931679254742902977502890883426551960403450937665750386501228142099266824028488862959626463948822181376617128628357")
	b := GetBigInt("5711912074763938920844020768820827016918638588776093786691324830937965710562669998102969607754216881533101753509522661181935679768137553251696427895001308210043958162362474454915118307661021406997989560047755201343617470288619030784987198511772840498354380632474664457429003510207310347179884080000294301502325103527312780599913053243627156705417875172756769585807691558680079741149166677442267851492473670184071199725213912264373214980177804010561543807969309223405291240876888702197126709861726023144260487044339708816278182396486957437256069194438047922679665536060592545457448379589893428429445378466414731324407")

	expectedResult := GetBigInt("-19913561841364303941087968013056854925409568225408501509608065500928998362191")
	actualResult, _ := CalculateHash(a, b)
	actualResult2, _ := CalculateHash(a, b)

	if expectedResult.Cmp(actualResult) != 0 {
		t.Errorf("Assert failure: hashed is: %s", actualResult)
	}
	if expectedResult.Cmp(actualResult2) != 0 {
		t.Errorf("Assert failure: hashed 2 is: %s", actualResult2)
	}
}

func TestModPow1(t *testing.T) {

	base := big.NewInt(10)
	exponent := big.NewInt(3)
	modulo := big.NewInt(7)

	result := ModPow(base, exponent, modulo)

	if result.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("Assert failure: expected 6, actual: %s", result)
	}
}

func TestModPow2(t *testing.T) {

	base := big.NewInt(30)
	exponent := big.NewInt(2)
	modulo := big.NewInt(7)

	var result = ModPow(base, exponent, modulo)

	if result.Cmp(big.NewInt(4)) != 0 {
		t.Errorf("Assert failure: expected 4, actual: %s", result)
	}
}

func TestModPowNegativeExp1(t *testing.T) {

	result := ModPow(big.NewInt(16), big.NewInt(-1), big.NewInt(7))

	if result.Cmp(big.NewInt(4)) != 0 {
		t.Errorf("Assert failure: expected 4, actual: %s", result)
	}
}

func TestModPowNegativeExp2(t *testing.T) {

	result := ModPow(big.NewInt(34), big.NewInt(-2), big.NewInt(9))

	if result.Cmp(big.NewInt(7)) != 0 {
		t.Errorf("Assert failure: expected 7, actual: %s", result)
	}
}

func TestModInverse1(t *testing.T) {

	base := big.NewInt(5)
	modulo := big.NewInt(1)

	var result = ModInverse(base, modulo)

	if result.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Assert failure: expected 0, actual: %s", result)
	}
}

func TestModInverse2(t *testing.T) {

	base := big.NewInt(3)
	modulo := big.NewInt(7)

	var result = ModInverse(base, modulo)

	if result.Cmp(big.NewInt(5)) != 0 {
		t.Errorf("Assert failure: expected 5, actual: %s", result)
	}
}

func TestMultiply(t *testing.T) {

	factor1 := big.NewInt(3)
	factor2 := big.NewInt(7)

	var result = Multiply(factor1, factor2)
	if result.Cmp(big.NewInt(21)) != 0 {
		t.Errorf("Assert failure: expected 21, actual: %s", result)
	}
}

func TestMod(t *testing.T) {

	result := Mod(big.NewInt(16), big.NewInt(7))

	if result.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("Assert failure: expected 2, actual: %s", result)
	}
}

