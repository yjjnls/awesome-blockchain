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
	"crypto/sha256"
	"encoding/json"
	"math/big"

	"../byteconversion"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/google"
	// "../crypto/bn256"
)

//Constants that are going to be used frequently, then we just need to compute them once.
var (
	G1 = new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(1))
	G2 = new(bn256.G2).ScalarBaseMult(new(big.Int).SetInt64(1))
	E  = bn256.Pair(G1, G2)
)

/*
Decompose receives as input a bigint x and outputs an array of integers such that
x = sum(xi.u^i), i.e. it returns the decomposition of x into base u.
*/
func Decompose(x *big.Int, u int64, l int64) ([]int64, error) {
	var (
		result []int64
		i      int64
	)
	result = make([]int64, l, l)
	i = 0
	for i < l {
		result[i] = Mod(x, new(big.Int).SetInt64(u)).Int64()
		x = new(big.Int).Div(x, new(big.Int).SetInt64(u))
		i = i + 1
	}
	return result, nil
}

/*
Commit method corresponds to the Pedersen commitment scheme. Namely, given input
message x, and randomness r, it outputs g^x.h^r.
*/
func Commit(x, r *big.Int, h *bn256.G2) (*bn256.G2, error) {
	var (
		C *bn256.G2
	)
	C = new(bn256.G2).ScalarBaseMult(x)
	C.Add(C, new(bn256.G2).ScalarMult(h, r))
	return C, nil
}

/*
CommitG1 method corresponds to the Pedersen commitment scheme. Namely, given input
message x, and randomness r, it outputs g^x.h^r.
*/
func CommitG1(x, r *big.Int, h *p256) (*p256, error) {
	var (
		C *p256
	)
	C = new(p256).ScalarBaseMult(x)
	Hr := new(p256).ScalarMult(h, r)
	C.Add(C, Hr)
	return C, nil
}

func Mult(a *p256, n *big.Int) *p256 {
	return new(p256).ScalarMult(a, n)
}

/*
HashSet is responsible for the computing a Zp element given elements from GT and G2.
*/
func HashSet(a *bn256.GT, D *bn256.G2) (*big.Int, error) {
	digest := sha256.New()
	digest.Write([]byte(a.String()))
	digest.Write([]byte(D.String()))
	output := digest.Sum(nil)
	tmp := output[0:len(output)]
	return byteconversion.FromByteArray(tmp)
}

/*
Hash is responsible for the computing a Zp element given elements from GT and G2.
*/
func Hash(a []*bn256.GT, D *bn256.G2) (*big.Int, error) {
	digest := sha256.New()
	for i := range a {
		digest.Write([]byte(a[i].String()))
	}
	digest.Write([]byte(D.String()))
	output := digest.Sum(nil)
	tmp := output[0:len(output)]
	return byteconversion.FromByteArray(tmp)
}

/*
Read big integer in base 10 from string.
*/
func GetBigInt(value string) *big.Int {
	i := new(big.Int)
	i.SetString(value, 10)
	return i
}

/*
Get common base
*/
func GetZkrp() *Bp {
	var zkrp Bp
	zkrp.Setup(0, 4294967296)
	return &zkrp
}

/*
Get zkrp verifier
*/
func GetVerifier(t *big.Int, h []*p256, p *p256) *Bp {
	zkrp := GetZkrp()
	zkrp.Zkip.Cc = t
	zkrp.Zkip.Hh = h
	zkrp.Zkip.P = p
	return zkrp
}

/*
Pedersen Commitment verification, check input_sum == output_sum ?
*/
func VerifyPedersenCommitment(input, output []*PedersenCommitment, blindDiff *big.Int) bool {
	inputCommitment := new(PedersenCommitment)
	outputCommitment := new(PedersenCommitment)
	for _, p := range input {
		inputCommitment = inputCommitment.Add(inputCommitment, p)
	}
	for _, p := range output {
		outputCommitment = outputCommitment.Add(outputCommitment, p)
	}
	// 计算佩德森承诺输入输出之差
	diffCommitment := inputCommitment.Add(inputCommitment, outputCommitment.Neg(outputCommitment))
	// 根据盲因子计算理论值
	H, _ := MapToGroup(SEEDH)
	checkCommitment := Mult(H, blindDiff)
	// 比较是否相等
	return checkCommitment.X.Cmp(diffCommitment.X) == 0 && checkCommitment.Y.Cmp(diffCommitment.Y) == 0

}

type ProofData struct {
	Proof *proofBP
	T     *big.Int
	Hh    []*p256
	P     *p256
}

func DumpProof(t *big.Int, h []*p256, p *p256, proof *proofBP) ([]byte, error) {
	zkrproof := &ProofData{proof, t, h, p}
	data, err := json.Marshal(zkrproof)
	return data, err
}

func LoadProof(data []byte) (*Bp, *proofBP, error) {
	var p ProofData
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, nil, err
	}
	return GetVerifier(p.T, p.Hh, p.P), p.Proof, nil
}
