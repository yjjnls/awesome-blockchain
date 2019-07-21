// Copyright 2018 ING Bank N.V.
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
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ing-bank/zkproofs/go-ethereum/crypto/bn256"
)

/*
Test method VectorCopy, which simply copies the first input argument to size n vector.
*/
func TestVectorCopy(t *testing.T) {
	var (
		result []*big.Int
	)
	result, _ = VectorCopy(new(big.Int).SetInt64(1), 3)
	ok := (result[0].Cmp(new(big.Int).SetInt64(1)) == 0)
	ok = ok && (result[1].Cmp(GetBigInt("1")) == 0)
	ok = ok && (result[2].Cmp(GetBigInt("1")) == 0)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Test method VectorConvertToBig.
*/
func TestVectorConvertToBig(t *testing.T) {
	var (
		result []*big.Int
		a      []int64
	)
	a = make([]int64, 3)
	a[0] = 3
	a[1] = 4
	a[2] = 5
	result, _ = VectorConvertToBig(a, 3)
	ok := (result[0].Cmp(new(big.Int).SetInt64(3)) == 0)
	ok = ok && (result[1].Cmp(GetBigInt("4")) == 0)
	ok = ok && (result[2].Cmp(GetBigInt("5")) == 0)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Scalar Product returns the inner product between 2 vectors.
*/
func TestScalarProduct(t *testing.T) {
	var (
		a, b []*big.Int
	)
	a = make([]*big.Int, 3)
	b = make([]*big.Int, 3)
	a[0] = new(big.Int).SetInt64(7)
	a[1] = new(big.Int).SetInt64(7)
	a[2] = new(big.Int).SetInt64(7)
	b[0] = new(big.Int).SetInt64(3)
	b[1] = new(big.Int).SetInt64(3)
	b[2] = new(big.Int).SetInt64(3)
	result, _ := ScalarProduct(a, b)
	ok := (result.Cmp(new(big.Int).SetInt64(63)) == 0)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Tests Vector addition.
*/
func TestVectorAdd(t *testing.T) {
	var (
		a, b []*big.Int
	)
	a = make([]*big.Int, 3)
	b = make([]*big.Int, 3)
	a[0] = new(big.Int).SetInt64(7)
	a[1] = new(big.Int).SetInt64(8)
	a[2] = new(big.Int).SetInt64(9)
	b[0] = new(big.Int).SetInt64(3)
	b[1] = new(big.Int).SetInt64(30)
	b[2] = new(big.Int).SetInt64(40)
	result, _ := VectorAdd(a, b)
	ok := (result[0].Cmp(new(big.Int).SetInt64(10)) == 0)
	ok = ok && (result[1].Cmp(GetBigInt("38")) == 0)
	ok = ok && (result[2].Cmp(GetBigInt("49")) == 0)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Tests Vector subtraction.
*/
func TestVectorSub(t *testing.T) {
	var (
		a, b []*big.Int
	)
	a = make([]*big.Int, 3)
	b = make([]*big.Int, 3)
	a[0] = new(big.Int).SetInt64(7)
	a[1] = new(big.Int).SetInt64(8)
	a[2] = new(big.Int).SetInt64(9)
	b[0] = new(big.Int).SetInt64(3)
	b[1] = new(big.Int).SetInt64(30)
	b[2] = new(big.Int).SetInt64(40)
	result, _ := VectorSub(a, b)
	ok := (result[0].Cmp(new(big.Int).SetInt64(4)) == 0)
	ok = ok && (result[1].Cmp(GetBigInt("115792089237316195423570985008687907852837564279074904382605163141518161494315")) == 0)
	ok = ok && (result[2].Cmp(GetBigInt("115792089237316195423570985008687907852837564279074904382605163141518161494306")) == 0)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Tests Vector componentwise multiplication.
*/
func TestVectorMul(t *testing.T) {
	var (
		a, b []*big.Int
	)
	a = make([]*big.Int, 3)
	b = make([]*big.Int, 3)
	a[0] = new(big.Int).SetInt64(7)
	a[1] = new(big.Int).SetInt64(8)
	a[2] = new(big.Int).SetInt64(9)
	b[0] = new(big.Int).SetInt64(3)
	b[1] = new(big.Int).SetInt64(30)
	b[2] = new(big.Int).SetInt64(40)
	result, _ := VectorMul(a, b)
	ok := (result[0].Cmp(new(big.Int).SetInt64(21)) == 0)
	ok = ok && (result[1].Cmp(new(big.Int).SetInt64(240)) == 0)
	ok = ok && (result[2].Cmp(new(big.Int).SetInt64(360)) == 0)

	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Test method PowerOf, which must return a vector containing a growing sequence of
powers of 2.
*/
func TestPowerOf(t *testing.T) {
	result, _ := PowerOf(new(big.Int).SetInt64(3), 3)
	ok := (result[0].Cmp(new(big.Int).SetInt64(1)) == 0)
	ok = ok && (result[1].Cmp(new(big.Int).SetInt64(3)) == 0)
	ok = ok && (result[2].Cmp(new(big.Int).SetInt64(9)) == 0)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Test Inner Product argument.
*/
func TestInnerProduct(t *testing.T) {
	var (
		zkrp Bp
		zkip bip
		a    []*big.Int
		b    []*big.Int
	)
	// TODO:
	// Review if it is the best way, since we maybe could use the
	// inner product independently of the range proof.
	zkrp.Setup(0, 16)
	a = make([]*big.Int, zkrp.N)
	a[0] = new(big.Int).SetInt64(2)
	a[1] = new(big.Int).SetInt64(-1)
	a[2] = new(big.Int).SetInt64(10)
	a[3] = new(big.Int).SetInt64(6)
	b = make([]*big.Int, zkrp.N)
	b[0] = new(big.Int).SetInt64(1)
	b[1] = new(big.Int).SetInt64(2)
	b[2] = new(big.Int).SetInt64(10)
	b[3] = new(big.Int).SetInt64(7)
	c := new(big.Int).SetInt64(142)
	commit, _ := CommitInnerProduct(zkrp.Gg, zkrp.Hh, a, b)
	zkip.Setup(zkrp.H, zkrp.Gg, zkrp.Hh, c)
	proof, _ := zkip.Prove(a, b, commit)
	ok, _ := zkip.Verify(proof)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Test the FALSE case of ZK Range Proof scheme using Bulletproofs.
*/
func TestFalseBulletproofsZKRP(t *testing.T) {
	var (
		zkrp Bp
	)
	startTime := time.Now()
	zkrp.Setup(0, 4294967296) // ITS BEING USED TO COMPUTE N
	setupTime := time.Now()
	fmt.Println("Setup time:")
	fmt.Println(setupTime.Sub(startTime))

	x := new(big.Int).SetInt64(4294967296)
	proof, _ := zkrp.Prove(x)
	proofTime := time.Now()
	fmt.Println("Proof time:")
	fmt.Println(proofTime.Sub(setupTime))

	ok, _ := zkrp.Verify(proof)
	verifyTime := time.Now()
	fmt.Println("Verify time:")
	fmt.Println(verifyTime.Sub(proofTime))

	fmt.Println("Range Proofs invalid test result:")
	fmt.Println(ok)
	if ok != false {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

/*
Test the TRUE case of ZK Range Proof scheme using Bulletproofs.
*/
func TestTrueBulletproofsZKRP(t *testing.T) {
	var (
		zkrp Bp
	)
	startTime := time.Now()
	zkrp.Setup(0, 4294967296) // ITS BEING USED TO COMPUTE N
	setupTime := time.Now()
	fmt.Println("Setup time:")
	fmt.Println(setupTime.Sub(startTime))

	x := new(big.Int).SetInt64(65535)
	proof, _ := zkrp.Prove(x)
	proofTime := time.Now()
	fmt.Println("Proof time:")
	fmt.Println(proofTime.Sub(setupTime))

	ok, _ := zkrp.Verify(proof)
	verifyTime := time.Now()
	fmt.Println("Verify time:")
	fmt.Println(verifyTime.Sub(proofTime))

	fmt.Println("Range Proofs result:")
	fmt.Println(ok)
	if ok != true {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

func BenchmarkBulletproofs(b *testing.B) {
	var (
		zkrp  Bp
		proof proofBP
		ok    bool
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zkrp.Setup(0, 4294967296) // ITS BEING USED TO COMPUTE N
		x := new(big.Int).SetInt64(4294967295)
		proof, _ = zkrp.Prove(x)
		ok, _ = zkrp.Verify(proof)
		if ok != true {
			b.Errorf("Assert failure: expected true, actual: %t", ok)
		}
	}
}

func BenchmarkScalarMult(b *testing.B) {
	var (
		a *big.Int
		A *bn256.G1
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a, _ = rand.Int(rand.Reader, bn256.Order)
		A = new(bn256.G1).ScalarBaseMult(a)
	}
	fmt.Println("A:")
	fmt.Println(A)
}

func TestHashBP(t *testing.T) {
	agx, _ := new(big.Int).SetString("110720467414728166769654679803728202169916280248550137472490865118702779748947", 10)
	agy, _ := new(big.Int).SetString("103949684536896233354287911519259186718323435572971865592336813380571928560949", 10)
	sgx, _ := new(big.Int).SetString("78662919066140655151560869958157053125629409725243565127658074141532489435921", 10)
	sgy, _ := new(big.Int).SetString("114946280626097680211499478702679495377587739951564115086530426937068100343655", 10)
	pointa := &p256{X: agx, Y: agy}
	points := &p256{X: sgx, Y: sgy}
	result1, result2, _ := HashBP(pointa, points)
	res1, _ := new(big.Int).SetString("103823382860325249552741530200099120077084118788867728791742258217664299339569", 10)
	res2, _ := new(big.Int).SetString("8192372577089859289404358830067912230280991346287696886048261417244724213964", 10)
	ok1 := (result1.Cmp(res1) != 0)
	ok2 := (result2.Cmp(res2) != 0)
	ok := ok1 && ok2
	if ok {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

func TestHashBPGx(t *testing.T) {
	gx, _ := new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	gy, _ := new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	point := &p256{X: gx, Y: gy}
	result1, result2, _ := HashBP(point, point)
	res1, _ := new(big.Int).SetString("11897424191990306464486192136408618361228444529783223689021929580052970909263", 10)
	res2, _ := new(big.Int).SetString("22166487799255634251145870394406518059682307840904574298117500050508046799269", 10)
	ok1 := (result1.Cmp(res1) != 0)
	ok2 := (result2.Cmp(res2) != 0)
	ok := ok1 && ok2
	if ok {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

func TestInv(t *testing.T) {
	y, _ := new(big.Int).SetString("103823382860325249552741530200099120077084118788867728791742258217664299339569", 10)
	yinv := ModInverse(y, ORDER)
	res, _ := new(big.Int).SetString("38397371868935917445400134055424677162505875368971619911110421656148020877351", 10)
	ok := (yinv.Cmp(res) != 0)
	if ok {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}

func TestHPrime(t *testing.T) {
	var zkrp *Bp
	var proof *proofBP
	zkrp, _ = LoadParamFromDisk("setup.dat")
	proof, _ = LoadProofFromDisk("proof.dat")
	ok, _ := zkrp.Verify(*proof)
	if !ok {
		t.Errorf("Assert failure: expected true, actual: %t", ok)
	}
}
