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

/*
This file contains the implementation of the Bulletproofs scheme proposed in the paper:
Bulletproofs: Short Proofs for Confidential Transactions and More
Benedikt Bunz, Jonathan Bootle, Dan Boneh, Andrew Poelstra, Pieter Wuille and Greg Maxwell
Asiacrypt 2008
*/

package zkproofs

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/big"

	"../byteconversion"
)

var (
	ORDER = CURVE.N
	SEEDH = "BulletproofsDoesNotNeedTrustedSetupH"
	SEEDU = "BulletproofsDoesNotNeedTrustedSetupU"
	SAVE  = true
)

/*
Bulletproofs parameters.
*/
type Bp struct {
	N    int64 // n 位
	G    *p256 // 曲线上的点 G 和 H
	H    *p256
	Gg   []*p256
	Hh   []*p256
	Zkip bip
}

/*
Bulletproofs proof.
*/
type proofBP struct {
	V       *p256
	A       *p256
	S       *p256
	T1      *p256
	T2      *p256
	Taux    *big.Int
	Mu      *big.Int
	Tprime  *big.Int
	Proofip proofBip
	Commit  *p256
}

type (
	pstring struct {
		X string
		Y string
	}
)

type (
	ipstring struct {
		N  int64
		A  string
		B  string
		U  pstring
		P  pstring
		Gg pstring
		Hh pstring
		Ls []pstring
		Rs []pstring
	}
)

func (p *proofBP) MarshalJSON() ([]byte, error) {
	type Alias proofBP
	var iLs []pstring
	var iRs []pstring
	var i int
	logn := len(p.Proofip.Ls)
	iLs = make([]pstring, logn)
	iRs = make([]pstring, logn)
	i = 0
	for i < logn {
		iLs[i] = pstring{X: p.Proofip.Ls[i].X.String(), Y: p.Proofip.Ls[i].Y.String()}
		iRs[i] = pstring{X: p.Proofip.Rs[i].X.String(), Y: p.Proofip.Rs[i].Y.String()}
		i = i + 1
	}
	return json.Marshal(&struct {
		V       pstring  `json:"V"`
		A       pstring  `json:"A"`
		S       pstring  `json:"S"`
		T1      pstring  `json:"T1"`
		T2      pstring  `json:"T2"`
		Taux    string   `json:"Taux"`
		Mu      string   `json:"Mu"`
		Tprime  string   `json:"Tprime"`
		Commit  pstring  `json:"Commit"`
		Proofip ipstring `json:"Proofip"`
		*Alias
	}{
		V:      pstring{X: p.V.X.String(), Y: p.V.Y.String()},
		A:      pstring{X: p.A.X.String(), Y: p.A.Y.String()},
		S:      pstring{X: p.S.X.String(), Y: p.S.Y.String()},
		T1:     pstring{X: p.T1.X.String(), Y: p.T1.Y.String()},
		T2:     pstring{X: p.T2.X.String(), Y: p.T2.Y.String()},
		Mu:     p.Mu.String(),
		Taux:   p.Taux.String(),
		Tprime: p.Tprime.String(),
		Commit: pstring{X: p.Commit.X.String(), Y: p.Commit.Y.String()},
		Proofip: ipstring{
			N:  p.Proofip.N,
			A:  p.Proofip.A.String(),
			B:  p.Proofip.B.String(),
			U:  pstring{X: p.Proofip.U.X.String(), Y: p.Proofip.U.Y.String()},
			P:  pstring{X: p.Proofip.P.X.String(), Y: p.Proofip.P.Y.String()},
			Gg: pstring{X: p.Proofip.Gg.X.String(), Y: p.Proofip.Gg.Y.String()},
			Hh: pstring{X: p.Proofip.Hh.X.String(), Y: p.Proofip.Hh.Y.String()},
			Ls: iLs,
			Rs: iRs,
		},
		Alias: (*Alias)(p),
	})
}

func (p *proofBP) UnmarshalJSON(data []byte) error {
	type Alias proofBP
	aux := &struct {
		V       pstring  `json:"V"`
		A       pstring  `json:"A"`
		S       pstring  `json:"S"`
		T1      pstring  `json:"T1"`
		T2      pstring  `json:"T2"`
		Taux    string   `json:"Taux"`
		Mu      string   `json:"Mu"`
		Tprime  string   `json:"Tprime"`
		Commit  pstring  `json:"Commit"`
		Proofip ipstring `json:"Proofip"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	valVX, _ := new(big.Int).SetString(aux.V.X, 10)
	valVY, _ := new(big.Int).SetString(aux.V.Y, 10)
	valAX, _ := new(big.Int).SetString(aux.A.X, 10)
	valAY, _ := new(big.Int).SetString(aux.A.Y, 10)
	valSX, _ := new(big.Int).SetString(aux.S.X, 10)
	valSY, _ := new(big.Int).SetString(aux.S.Y, 10)
	valT1X, _ := new(big.Int).SetString(aux.T1.X, 10)
	valT1Y, _ := new(big.Int).SetString(aux.T1.Y, 10)
	valT2X, _ := new(big.Int).SetString(aux.T2.X, 10)
	valT2Y, _ := new(big.Int).SetString(aux.T2.Y, 10)
	valCommitX, _ := new(big.Int).SetString(aux.Commit.X, 10)
	valCommitY, _ := new(big.Int).SetString(aux.Commit.Y, 10)
	valN := aux.Proofip.N
	valA, _ := new(big.Int).SetString(aux.Proofip.A, 10)
	valB, _ := new(big.Int).SetString(aux.Proofip.B, 10)
	valUx, _ := new(big.Int).SetString(aux.Proofip.U.X, 10)
	valUy, _ := new(big.Int).SetString(aux.Proofip.U.Y, 10)
	valPx, _ := new(big.Int).SetString(aux.Proofip.P.X, 10)
	valPy, _ := new(big.Int).SetString(aux.Proofip.P.Y, 10)
	valGgx, _ := new(big.Int).SetString(aux.Proofip.Gg.X, 10)
	valGgy, _ := new(big.Int).SetString(aux.Proofip.Gg.Y, 10)
	valHhx, _ := new(big.Int).SetString(aux.Proofip.Hh.X, 10)
	valHhy, _ := new(big.Int).SetString(aux.Proofip.Hh.Y, 10)
	p.V = &p256{
		X: valVX,
		Y: valVY,
	}
	p.A = &p256{
		X: valAX,
		Y: valAY,
	}
	p.S = &p256{
		X: valSX,
		Y: valSY,
	}
	p.T1 = &p256{
		X: valT1X,
		Y: valT1Y,
	}
	p.T2 = &p256{
		X: valT2X,
		Y: valT2Y,
	}
	p.Commit = &p256{
		X: valCommitX,
		Y: valCommitY,
	}
	valU := &p256{
		X: valUx,
		Y: valUy,
	}
	valP := &p256{
		X: valPx,
		Y: valPy,
	}
	valGg := &p256{
		X: valGgx,
		Y: valGgy,
	}
	valHh := &p256{
		X: valHhx,
		Y: valHhy,
	}
	p.Taux, _ = new(big.Int).SetString(aux.Taux, 10)
	p.Mu, _ = new(big.Int).SetString(aux.Mu, 10)
	p.Tprime, _ = new(big.Int).SetString(aux.Tprime, 10)
	logn := len(aux.Proofip.Ls)
	valLs := make([]*p256, logn)
	valRs := make([]*p256, logn)
	var (
		i      int
		valLsx *big.Int
		valLsy *big.Int
		valRsx *big.Int
		valRsy *big.Int
	)
	i = 0
	for i < logn {
		valLsx, _ = new(big.Int).SetString(aux.Proofip.Ls[i].X, 10)
		valLsy, _ = new(big.Int).SetString(aux.Proofip.Ls[i].Y, 10)
		valLs[i] = &p256{X: valLsx, Y: valLsy}
		valRsx, _ = new(big.Int).SetString(aux.Proofip.Rs[i].X, 10)
		valRsy, _ = new(big.Int).SetString(aux.Proofip.Rs[i].Y, 10)
		valRs[i] = &p256{X: valRsx, Y: valRsy}
		i = i + 1
	}
	p.Proofip = proofBip{
		N:  valN,
		A:  valA,
		B:  valB,
		U:  valU,
		P:  valP,
		Gg: valGg,
		Hh: valHh,
		Ls: valLs,
		Rs: valRs,
	}
	return nil
}

type (
	ipgenstring struct {
		N  int64
		Cc string
		Uu pstring
		H  pstring
		Gg []pstring
		Hh []pstring
		P  pstring
	}
)

func (s *Bp) MarshalJSON() ([]byte, error) {
	type Alias Bp
	var iHh []pstring
	var iGg []pstring

	var i int
	n := len(s.Gg)
	iGg = make([]pstring, n)
	iHh = make([]pstring, n)
	i = 0
	for i < n {
		iGg[i] = pstring{X: s.Zkip.Gg[i].X.String(), Y: s.Zkip.Gg[i].Y.String()}
		iHh[i] = pstring{X: s.Zkip.Hh[i].X.String(), Y: s.Zkip.Hh[i].Y.String()}
		i = i + 1
	}
	return json.Marshal(&struct {
		Zkip ipgenstring `json:"Zkip"`
		*Alias
	}{
		Zkip: ipgenstring{
			N:  s.N,
			Cc: s.Zkip.Cc.String(),
			Uu: pstring{X: s.Zkip.Uu.X.String(), Y: s.Zkip.Uu.Y.String()},
			H:  pstring{X: s.Zkip.H.X.String(), Y: s.Zkip.H.Y.String()},
			Gg: iGg,
			Hh: iHh,
			P:  pstring{X: s.Zkip.P.X.String(), Y: s.Zkip.P.Y.String()},
		},
		Alias: (*Alias)(s),
	})
}

func (s *Bp) UnmarshalJSON(data []byte) error {
	type Alias Bp
	aux := &struct {
		Zkip ipgenstring `json:"Zkip"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	n := aux.N
	valGg := make([]*p256, n)
	valHh := make([]*p256, n)
	var (
		i      int64
		valGgx *big.Int
		valGgy *big.Int
		valHhx *big.Int
		valHhy *big.Int
	)
	i = 0
	for i < n {
		valGgx, _ = new(big.Int).SetString(aux.Zkip.Gg[i].X, 10)
		valGgy, _ = new(big.Int).SetString(aux.Zkip.Gg[i].Y, 10)
		valGg[i] = &p256{X: valGgx, Y: valGgy}
		valHhx, _ = new(big.Int).SetString(aux.Zkip.Hh[i].X, 10)
		valHhy, _ = new(big.Int).SetString(aux.Zkip.Hh[i].Y, 10)
		valHh[i] = &p256{X: valHhx, Y: valHhy}
		i = i + 1
	}
	valN := aux.N
	valCc, _ := new(big.Int).SetString(aux.Zkip.Cc, 10)
	valUux, _ := new(big.Int).SetString(aux.Zkip.Uu.X, 10)
	valUuy, _ := new(big.Int).SetString(aux.Zkip.Uu.Y, 10)
	valHx, _ := new(big.Int).SetString(aux.Zkip.H.X, 10)
	valHy, _ := new(big.Int).SetString(aux.Zkip.H.Y, 10)
	valPx, _ := new(big.Int).SetString(aux.Zkip.P.X, 10)
	valPy, _ := new(big.Int).SetString(aux.Zkip.P.Y, 10)
	valUu := &p256{
		X: valUux,
		Y: valUuy,
	}
	valH := &p256{
		X: valHx,
		Y: valHy,
	}
	valP := &p256{
		X: valPx,
		Y: valPy,
	}
	s.Zkip = bip{
		N:  valN,
		Cc: valCc,
		Uu: valUu,
		H:  valH,
		Gg: valGg,
		Hh: valHh,
		P:  valP,
	}
	return nil
}

/*
VectorCopy returns a vector composed by copies of a.
*/
func VectorCopy(a *big.Int, n int64) ([]*big.Int, error) {
	var (
		i      int64
		result []*big.Int
	)
	result = make([]*big.Int, n)
	i = 0
	for i < n {
		result[i] = a
		i = i + 1
	}
	return result, nil
}

/*
VectorCopy returns a vector composed by copies of a.
*/
func VectorG1Copy(a *p256, n int64) ([]*p256, error) {
	var (
		i      int64
		result []*p256
	)
	result = make([]*p256, n)
	i = 0
	for i < n {
		result[i] = a
		i = i + 1
	}
	return result, nil
}

/*
VectorConvertToBig converts an array of int64 to an array of big.Int.
*/
func VectorConvertToBig(a []int64, n int64) ([]*big.Int, error) {
	var (
		i      int64
		result []*big.Int
	)
	result = make([]*big.Int, n)
	i = 0
	for i < n {
		result[i] = new(big.Int).SetInt64(a[i])
		i = i + 1
	}
	return result, nil
}

/*
VectorAdd computes vector addition componentwisely.
*/
func VectorAdd(a, b []*big.Int) ([]*big.Int, error) {
	var (
		result  []*big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("Size of first argument is different from size of second argument.")
	}
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = Add(a[i], b[i])
		result[i] = Mod(result[i], ORDER)
		i = i + 1
	}
	return result, nil
}

/*
VectorSub computes vector addition componentwisely.
*/
func VectorSub(a, b []*big.Int) ([]*big.Int, error) {
	var (
		result  []*big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("Size of first argument is different from size of second argument.")
	}
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = Sub(a[i], b[i])
		result[i] = Mod(result[i], ORDER)
		i = i + 1
	}
	return result, nil
}

/*
VectorScalarMul computes vector scalar multiplication componentwisely.
*/
func VectorScalarMul(a []*big.Int, b *big.Int) ([]*big.Int, error) {
	var (
		result []*big.Int
		i, n   int64
	)
	n = int64(len(a))
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = Multiply(a[i], b)
		result[i] = Mod(result[i], ORDER)
		i = i + 1
	}
	return result, nil
}

/*
VectorMul computes vector multiplication componentwisely.
*/
func VectorMul(a, b []*big.Int) ([]*big.Int, error) {
	var (
		result  []*big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("Size of first argument is different from size of second argument.")
	}
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = Multiply(a[i], b[i])
		result[i] = Mod(result[i], ORDER)
		i = i + 1
	}
	return result, nil
}

/*
VectorECMul computes vector EC addition componentwisely.
*/
func VectorECAdd(a, b []*p256) ([]*p256, error) {
	var (
		result  []*p256
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("Size of first argument is different from size of second argument.")
	}
	result = make([]*p256, n)
	i = 0
	for i < n {
		result[i] = new(p256).Multiply(a[i], b[i])
		i = i + 1
	}
	return result, nil
}

/*
ScalarProduct return the inner product between a and b.
*/
func ScalarProduct(a, b []*big.Int) (*big.Int, error) {
	var (
		result  *big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("Size of first argument is different from size of second argument.")
	}
	i = 0
	result = GetBigInt("0")
	for i < n {
		ab := Multiply(a[i], b[i])
		result.Add(result, ab)
		result = Mod(result, ORDER)
		i = i + 1
	}
	return result, nil
}

/*
VectorExp computes Prod_i^n{a[i]^b[i]}.
*/
func VectorExp(a []*p256, b []*big.Int) (*p256, error) {
	var (
		result  *p256
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("Size of first argument is different from size of second argument.")
	}
	i = 0
	result = new(p256).SetInfinity()
	for i < n {
		result.Multiply(result, new(p256).ScalarMult(a[i], b[i]))
		i = i + 1
	}
	return result, nil
}

/*
VectorScalarExp computes a[i]^b for each i.
*/
func VectorScalarExp(a []*p256, b *big.Int) ([]*p256, error) {
	var (
		result []*p256
		i, n   int64
	)
	n = int64(len(a))
	result = make([]*p256, n)
	i = 0
	for i < n {
		result[i] = new(p256).ScalarMult(a[i], b)
		i = i + 1
	}
	return result, nil
}

/*
PowerOf returns a vector composed by powers of x.
*/
func PowerOf(x *big.Int, n int64) ([]*big.Int, error) {
	var (
		i      int64
		result []*big.Int
	)
	result = make([]*big.Int, n)
	current := GetBigInt("1")
	i = 0
	for i < n {
		result[i] = current
		current = Multiply(current, x)
		current = Mod(current, ORDER)
		i = i + 1
	}
	return result, nil
}

/*
aR = aL - 1^n
*/
func ComputeAR(x []int64) ([]int64, error) {
	var (
		i      int64
		result []int64
	)
	result = make([]int64, len(x))
	i = 0
	for i < int64(len(x)) {
		if x[i] == 0 {
			result[i] = -1
		} else if x[i] == 1 {
			result[i] = 0
		} else {
			return nil, errors.New("input contains non-binary element")
		}
		i = i + 1
	}
	return result, nil
}

/*
Hash is responsible for the computing a Zp element given elements from GT and G1.
*/
func HashBP(A, S *p256) (*big.Int, *big.Int, error) {

	digest1 := sha256.New()
	var buffer bytes.Buffer
	buffer.WriteString(A.X.String())
	buffer.WriteString(A.Y.String())
	buffer.WriteString(S.X.String())
	buffer.WriteString(S.Y.String())
	digest1.Write([]byte(buffer.String()))
	output1 := digest1.Sum(nil)
	tmp1 := output1[0:len(output1)]
	result1 := new(big.Int).SetBytes(tmp1)

	digest2 := sha256.New()
	var buffer2 bytes.Buffer
	buffer2.WriteString(A.X.String())
	buffer2.WriteString(A.Y.String())
	buffer2.WriteString(S.X.String())
	buffer2.WriteString(S.Y.String())
	buffer2.WriteString(result1.String())
	digest2.Write([]byte(buffer2.String()))
	output2 := digest2.Sum(nil)
	tmp2 := output2[0:len(output2)]
	result2 := new(big.Int).SetBytes(tmp2)

	return result1, result2, nil
}

/*
Commitvector computes a commitment to the bit of the secret.
*/
func CommitVector(aL, aR []int64, alpha *big.Int, G, H *p256, g, h []*p256, n int64) (*p256, error) {
	var (
		i int64
		R *p256
	)
	// Compute h^alpha.vg^aL.vh^aR
	R = new(p256).ScalarMult(H, alpha)
	i = 0
	for i < n {
		gaL := new(p256).ScalarMult(g[i], new(big.Int).SetInt64(aL[i]))
		haR := new(p256).ScalarMult(h[i], new(big.Int).SetInt64(aR[i]))
		R.Multiply(R, gaL)
		R.Multiply(R, haR)
		i = i + 1
	}
	return R, nil
}

/*

 */
func CommitVectorBig(aL, aR []*big.Int, alpha *big.Int, G, H *p256, g, h []*p256, n int64) (*p256, error) {
	var (
		i int64
		R *p256
	)
	// Compute h^alpha.vg^aL.vh^aR
	R = new(p256).ScalarMult(H, alpha)
	i = 0
	for i < n {
		R.Multiply(R, new(p256).ScalarMult(g[i], aL[i]))
		R.Multiply(R, new(p256).ScalarMult(h[i], aR[i]))
		i = i + 1
	}
	return R, nil
}

/*
SaveToDisk is responsible for saving the generator to disk, such it is possible
to then later.
*/
func (zkrp *Bp) SaveToDisk(s string, p *proofBP) error {
	data, err := json.Marshal(zkrp)
	errw := ioutil.WriteFile(s, data, 0644)
	if p != nil {
		datap, errp := json.Marshal(p)
		errpw := ioutil.WriteFile("proof.dat", datap, 0644)
		if errp != nil || errpw != nil {
			return errors.New("proof not saved to disk.")
		}
	}
	if err != nil || errw != nil {
		return errors.New("parameters not saved to disk.")
	}
	return nil
}

/*
LoadGenFromDisk reads the generator from a file.
*/
func LoadParamFromDisk(s string) (*Bp, error) {
	var result Bp
	c, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}
	if len(c) > 0 {
		json.Unmarshal(c, &result)
		return &result, nil
	}
	return nil, errors.New("Could not load generators.")
}

/*
LoadProofFromDisk reads the generator from a file.
*/
func LoadProofFromDisk(s string) (*proofBP, error) {
	var result proofBP
	c, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}
	if len(c) > 0 {
		json.Unmarshal(c, &result)
		return &result, nil
	}
	return nil, errors.New("Could not load proof.")
}

/*
delta(y,z) = (z-z^2) . < 1^n, y^n > - z^3 . < 1^n, 2^n >
*/
func (zkrp *Bp) Delta(y, z *big.Int) (*big.Int, error) {
	var (
		result *big.Int
	)
	// delta(y,z) = (z-z^2) . < 1^n, y^n > - z^3 . < 1^n, 2^n >
	z2 := Multiply(z, z)
	z2 = Mod(z2, ORDER)
	z3 := Multiply(z2, z)
	z3 = Mod(z3, ORDER)

	// < 1^n, y^n >
	v1, _ := VectorCopy(new(big.Int).SetInt64(1), zkrp.N)
	vy, _ := PowerOf(y, zkrp.N)
	sp1y, _ := ScalarProduct(v1, vy)

	// < 1^n, 2^n >
	p2n, _ := PowerOf(new(big.Int).SetInt64(2), zkrp.N)
	sp12, _ := ScalarProduct(v1, p2n)

	result = Sub(z, z2)
	result = Mod(result, ORDER)
	result = Multiply(result, sp1y)
	result = Mod(result, ORDER)
	result = Sub(result, Multiply(z3, sp12))
	result = Mod(result, ORDER)

	return result, nil
}

/*
SetupPre is responsible for computing the common parameters.
*/
func (zkrp *Bp) SetupPre(a, b int64) {
	res, _ := LoadParamFromDisk("setup.json")
	zkrp = res
	// Setup Inner Product
	zkrp.Zkip.Setup(zkrp.H, zkrp.Gg, zkrp.Hh, new(big.Int).SetInt64(0))
}

/*
Setup is responsible for computing the common parameters.
*/
func (zkrp *Bp) Setup(a, b int64) {
	var (
		i int64
	)
	// 计算 G 和 H
	zkrp.G = new(p256).ScalarBaseMult(new(big.Int).SetInt64(1))
	zkrp.H, _ = MapToGroup(SEEDH)
	// 有 n 位
	zkrp.N = int64(math.Log2(float64(b)))
	zkrp.Gg = make([]*p256, zkrp.N)
	zkrp.Hh = make([]*p256, zkrp.N)
	i = 0
	for i < zkrp.N {
		zkrp.Gg[i], _ = MapToGroup(SEEDH + "g" + string(i))
		zkrp.Hh[i], _ = MapToGroup(SEEDH + "h" + string(i))
		i = i + 1
	}

	// Setup Inner Product
	zkrp.Zkip.Setup(zkrp.H, zkrp.Gg, zkrp.Hh, new(big.Int).SetInt64(0))
	// zkrp.SaveToDisk("setup.json", nil)
}

/*
Prove computes the ZK proof.
*/
func (zkrp *Bp) GenerateProof(secret *big.Int) (*big.Int, *big.Int, []*p256, *p256, proofBP, error) {
	var (
		i     int64
		sL    []*big.Int
		sR    []*big.Int
		proof proofBP
	)
	//////////////////////////////////////////////////////////////////////////////
	// First phase
	//////////////////////////////////////////////////////////////////////////////

	// commitment to v and gamma
	gamma, _ := rand.Int(rand.Reader, ORDER)
	V, _ := CommitG1(secret, gamma, zkrp.H)

	// aL, aR and commitment: (A, alpha)
	// 因式分解得到 aL
	aL, _ := Decompose(secret, 2, zkrp.N)
	// aR = aL - 1
	aR, _ := ComputeAR(aL)
	// 盲因子 alpha
	alpha, _ := rand.Int(rand.Reader, ORDER)
	// A 为 aL 的佩德森承诺
	A, _ := CommitVector(aL, aR, alpha, zkrp.G, zkrp.H, zkrp.Gg, zkrp.Hh, zkrp.N)

	// sL, sR and commitment: (S, rho)
	rho, _ := rand.Int(rand.Reader, ORDER)
	sL = make([]*big.Int, zkrp.N)
	sR = make([]*big.Int, zkrp.N)
	i = 0
	for i < zkrp.N {
		sL[i], _ = rand.Int(rand.Reader, ORDER)
		sR[i], _ = rand.Int(rand.Reader, ORDER)
		i = i + 1
	}
	// S 为 aR 的佩德森承诺
	S, _ := CommitVectorBig(sL, sR, rho, zkrp.G, zkrp.H, zkrp.Gg, zkrp.Hh, zkrp.N)

	// Fiat-Shamir heuristic to compute challenges y, z
	y, z, _ := HashBP(A, S)

	//////////////////////////////////////////////////////////////////////////////
	// Second phase
	//////////////////////////////////////////////////////////////////////////////
	tau1, _ := rand.Int(rand.Reader, ORDER) // page 20 from eprint version
	tau2, _ := rand.Int(rand.Reader, ORDER)

	// compute t1: < aL - z.1^n, y^n . sR > + < sL, y^n . (aR + z . 1^n) >
	vz, _ := VectorCopy(z, zkrp.N)
	vy, _ := PowerOf(y, zkrp.N)

	// aL - z.1^n
	naL, _ := VectorConvertToBig(aL, zkrp.N)
	aLmvz, _ := VectorSub(naL, vz)

	// y^n .sR
	ynsR, _ := VectorMul(vy, sR)

	// scalar prod: < aL - z.1^n, y^n . sR >
	sp1, _ := ScalarProduct(aLmvz, ynsR)

	// scalar prod: < sL, y^n . (aR + z . 1^n) >
	naR, _ := VectorConvertToBig(aR, zkrp.N)
	aRzn, _ := VectorAdd(naR, vz)
	ynaRzn, _ := VectorMul(vy, aRzn)

	// Add z^2.2^n to the result
	// z^2 . 2^n
	p2n, _ := PowerOf(new(big.Int).SetInt64(2), zkrp.N)
	zsquared := Multiply(z, z)
	z22n, _ := VectorScalarMul(p2n, zsquared)
	ynaRzn, _ = VectorAdd(ynaRzn, z22n)
	sp2, _ := ScalarProduct(sL, ynaRzn)

	// sp1 + sp2
	t1 := Add(sp1, sp2)
	t1 = Mod(t1, ORDER)

	// compute t2: < sL, y^n . sR >
	t2, _ := ScalarProduct(sL, ynsR)
	t2 = Mod(t2, ORDER)

	// compute T1
	T1, _ := CommitG1(t1, tau1, zkrp.H)

	// compute T2
	T2, _ := CommitG1(t2, tau2, zkrp.H)

	// Fiat-Shamir heuristic to compute 'random' challenge x
	x, _, _ := HashBP(T1, T2)

	//////////////////////////////////////////////////////////////////////////////
	// Third phase                                                              //
	//////////////////////////////////////////////////////////////////////////////

	// compute bl
	sLx, _ := VectorScalarMul(sL, x)
	bl, _ := VectorAdd(aLmvz, sLx)

	// compute br
	// y^n . ( aR + z.1^n + sR.x )
	sRx, _ := VectorScalarMul(sR, x)
	aRzn, _ = VectorAdd(aRzn, sRx)
	ynaRzn, _ = VectorMul(vy, aRzn)
	// y^n . ( aR + z.1^n sR.x ) + z^2 . 2^n
	br, _ := VectorAdd(ynaRzn, z22n)

	// Compute t` = < bl, br >
	tprime, _ := ScalarProduct(bl, br)

	// Compute taux = tau2 . x^2 + tau1 . x + z^2 . gamma
	taux := Multiply(tau2, Multiply(x, x))
	taux = Add(taux, Multiply(tau1, x))
	taux = Add(taux, Multiply(Multiply(z, z), gamma))
	taux = Mod(taux, ORDER)

	// Compute mu = alpha + rho.x
	mu := Multiply(rho, x)
	mu = Add(mu, alpha)
	mu = Mod(mu, ORDER)

	// Inner Product over (g, h', P.h^-mu, tprime)
	// Compute h'
	hprime := make([]*p256, zkrp.N)
	// Switch generators
	yinv := ModInverse(y, ORDER)
	expy := yinv
	hprime[0] = zkrp.Hh[0]
	i = 1
	for i < zkrp.N {
		hprime[i] = new(p256).ScalarMult(zkrp.Hh[i], expy)
		expy = Multiply(expy, yinv)
		i = i + 1
	}

	// Update Inner Product Proof Setup
	zkrp.Zkip.Hh = hprime
	zkrp.Zkip.Cc = tprime

	commit, _ := CommitInnerProduct(zkrp.Gg, hprime, bl, br)
	proofip, _ := zkrp.Zkip.GenerateProof(bl, br, commit)

	proof.V = V
	proof.A = A
	proof.S = S
	proof.T1 = T1
	proof.T2 = T2
	proof.Taux = taux
	proof.Mu = mu
	proof.Tprime = tprime
	proof.Proofip = proofip
	proof.Commit = commit

	// zkrp.SaveToDisk("setup.json", &proof)
	return gamma, tprime, hprime, zkrp.Zkip.P, proof, nil
}

/*
Verify returns true if and only if the proof is valid.
*/
func (zkrp *Bp) Verify(proof proofBP) (bool, error) {
	var (
		i      int64
		hprime []*p256
	)
	hprime = make([]*p256, zkrp.N)
	y, z, _ := HashBP(proof.A, proof.S)
	x, _, _ := HashBP(proof.T1, proof.T2)

	// Switch generators
	yinv := ModInverse(y, ORDER)
	expy := yinv
	hprime[0] = zkrp.Hh[0]
	i = 1
	for i < zkrp.N {
		hprime[i] = new(p256).ScalarMult(zkrp.Hh[i], expy)
		expy = Multiply(expy, yinv)
		i = i + 1
	}

	//////////////////////////////////////////////////////////////////////////////
	// Check that tprime  = t(x) = t0 + t1x + t2x^2  ----------  Condition (65) //
	//////////////////////////////////////////////////////////////////////////////

	// Compute left hand side
	lhs, _ := CommitG1(proof.Tprime, proof.Taux, zkrp.H)

	// Compute right hand side
	z2 := Multiply(z, z)
	z2 = Mod(z2, ORDER)
	x2 := Multiply(x, x)
	x2 = Mod(x2, ORDER)

	rhs := new(p256).ScalarMult(proof.V, z2)

	delta, _ := zkrp.Delta(y, z)

	gdelta := new(p256).ScalarBaseMult(delta)

	rhs.Multiply(rhs, gdelta)

	T1x := new(p256).ScalarMult(proof.T1, x)
	T2x2 := new(p256).ScalarMult(proof.T2, x2)

	rhs.Multiply(rhs, T1x)
	rhs.Multiply(rhs, T2x2)

	// Subtract lhs and rhs and compare with poitn at infinity
	lhs.Neg(lhs)
	rhs.Multiply(rhs, lhs)
	c65 := rhs.IsZero() // Condition (65), page 20, from eprint version

	// Compute P - lhs  #################### Condition (66) ######################

	// S^x
	Sx := new(p256).ScalarMult(proof.S, x)
	// A.S^x
	ASx := new(p256).Add(proof.A, Sx)

	// g^-z
	mz := Sub(ORDER, z)
	vmz, _ := VectorCopy(mz, zkrp.N)
	gpmz, _ := VectorExp(zkrp.Gg, vmz)
	//fmt.Println("############## gpmz ###############")
	//fmt.Println(gpmz);

	// z.y^n
	vz, _ := VectorCopy(z, zkrp.N)
	vy, _ := PowerOf(y, zkrp.N)
	zyn, _ := VectorMul(vy, vz)

	p2n, _ := PowerOf(new(big.Int).SetInt64(2), zkrp.N)
	zsquared := Multiply(z, z)
	z22n, _ := VectorScalarMul(p2n, zsquared)

	// z.y^n + z^2.2^n
	zynz22n, _ := VectorAdd(zyn, z22n)

	lP := new(p256)
	lP.Add(ASx, gpmz)

	// h'^(z.y^n + z^2.2^n)
	hprimeexp, _ := VectorExp(hprime, zynz22n)

	lP.Add(lP, hprimeexp)

	// Compute P - rhs  #################### Condition (67) ######################

	// h^mu
	rP := new(p256).ScalarMult(zkrp.H, proof.Mu)
	rP.Multiply(rP, proof.Commit)

	// Subtract lhs and rhs and compare with poitn at infinity
	lP = lP.Neg(lP)
	rP.Add(rP, lP)
	c67 := rP.IsZero()

	// Verify Inner Product Proof ################################################
	ok, _ := zkrp.Zkip.Verify(proof.Proofip)

	result := c65 && c67 && ok

	return result, nil
}

//////////////////////////////////// Inner Product ////////////////////////////////////

/*
Base struct for the Inner Product Argument.
*/
type bip struct {
	N  int64
	Cc *big.Int
	Uu *p256
	H  *p256
	Gg []*p256
	Hh []*p256
	P  *p256
}

/*
Struct that contains the Inner Product Proof.
*/
type proofBip struct {
	Ls []*p256
	Rs []*p256
	U  *p256
	P  *p256
	Gg *p256
	Hh *p256
	A  *big.Int
	B  *big.Int
	N  int64
}

/*
HashIP is responsible for the computing a Zp element given elements from GT and G1.
*/
func HashIP(g, h []*p256, P *p256, c *big.Int, n int64) (*big.Int, error) {
	var (
		i int64
	)

	digest := sha256.New()
	digest.Write([]byte(P.String()))

	i = 0
	for i < n {
		digest.Write([]byte(g[i].String()))
		digest.Write([]byte(h[i].String()))
		i = i + 1
	}

	digest.Write([]byte(c.String()))
	output := digest.Sum(nil)
	tmp := output[0:len(output)]
	result, err := byteconversion.FromByteArray(tmp)

	return result, err
}

/*
CommitInnerProduct is responsible for calculating g^a.h^b.
*/
func CommitInnerProduct(g, h []*p256, a, b []*big.Int) (*p256, error) {
	var (
		result *p256
	)

	ga, _ := VectorExp(g, a)
	hb, _ := VectorExp(h, b)
	result = new(p256).Multiply(ga, hb)
	return result, nil
}

/*
Setup is responsible for computing the inner product basic parameters that are common to both
Prove and Verify algorithms.
*/
func (zkip *bip) Setup(H *p256, g, h []*p256, c *big.Int) (bip, error) {
	var (
		params bip
	)

	zkip.Gg = make([]*p256, zkip.N)
	zkip.Hh = make([]*p256, zkip.N)
	zkip.Uu, _ = MapToGroup(SEEDU)
	zkip.H = H
	zkip.Gg = g
	zkip.Hh = h
	zkip.Cc = c
	zkip.P = new(p256).SetInfinity()

	return params, nil
}

/*
Prove is responsible for the generation of the Inner Product Proof.
*/
func (zkip *bip) GenerateProof(a, b []*big.Int, P *p256) (proofBip, error) {
	var (
		proof proofBip
		n, m  int64
		Ls    []*p256
		Rs    []*p256
	)

	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return proof, errors.New("Size of first array argument must be equal to the second")
	} else {
		// Fiat-Shamir:
		// x = Hash(g,h,P,c)
		x, _ := HashIP(zkip.Gg, zkip.Hh, P, zkip.Cc, zkip.N)
		// Pprime = P.u^(x.c)
		ux := new(p256).ScalarMult(zkip.Uu, x)
		uxc := new(p256).ScalarMult(ux, zkip.Cc)
		PP := new(p256).Multiply(P, uxc)
		// Execute Protocol 2 recursively
		zkip.P = PP
		proof, err := BIP(a, b, zkip.Gg, zkip.Hh, ux, zkip.P, n, Ls, Rs)
		proof.P = PP
		return proof, err
	}

	return proof, nil
}

/*
BIP is the main recursive function that will be used to compute the inner product argument.
*/
func BIP(a, b []*big.Int, g, h []*p256, u, P *p256, n int64, Ls, Rs []*p256) (proofBip, error) {
	var (
		proof                            proofBip
		cL, cR, x, xinv, x2, x2inv       *big.Int
		L, R, Lh, Rh, Pprime             *p256
		gprime, hprime, gprime2, hprime2 []*p256
		aprime, bprime, aprime2, bprime2 []*big.Int
	)

	if n == 1 {
		// recursion end
		proof.A = a[0]
		proof.B = b[0]
		proof.Gg = g[0]
		proof.Hh = h[0]
		proof.P = P
		proof.U = u
		proof.Ls = Ls
		proof.Rs = Rs

	} else {
		// recursion

		// nprime := n / 2
		nprime := n / 2

		// Compute cL = < a[:n'], b[n':] >
		cL, _ = ScalarProduct(a[:nprime], b[nprime:])
		// Compute cR = < a[n':], b[:n'] >
		cR, _ = ScalarProduct(a[nprime:], b[:nprime])
		// Compute L = g[n':]^(a[:n']).h[:n']^(b[n':]).u^cL
		L, _ = VectorExp(g[nprime:], a[:nprime])
		Lh, _ = VectorExp(h[:nprime], b[nprime:])
		L.Multiply(L, Lh)
		L.Multiply(L, new(p256).ScalarMult(u, cL))

		// Compute R = g[:n']^(a[n':]).h[n':]^(b[:n']).u^cR
		R, _ = VectorExp(g[:nprime], a[nprime:])
		Rh, _ = VectorExp(h[nprime:], b[:nprime])
		R.Multiply(R, Rh)
		R.Multiply(R, new(p256).ScalarMult(u, cR))

		// Fiat-Shamir:
		x, _, _ = HashBP(L, R)
		xinv = ModInverse(x, ORDER)

		// Compute g' = g[:n']^(x^-1) * g[n':]^(x)
		gprime, _ = VectorScalarExp(g[:nprime], xinv)
		gprime2, _ = VectorScalarExp(g[nprime:], x)
		gprime, _ = VectorECAdd(gprime, gprime2)
		// Compute h' = h[:n']^(x)    * h[n':]^(x^-1)
		hprime, _ = VectorScalarExp(h[:nprime], x)
		hprime2, _ = VectorScalarExp(h[nprime:], xinv)
		hprime, _ = VectorECAdd(hprime, hprime2)

		// Compute P' = L^(x^2).P.R^(x^-2)
		x2 = Mod(Multiply(x, x), ORDER)
		x2inv = ModInverse(x2, ORDER)
		Pprime = new(p256).ScalarMult(L, x2)
		Pprime.Multiply(Pprime, P)
		Pprime.Multiply(Pprime, new(p256).ScalarMult(R, x2inv))

		// Compute a' = a[:n'].x      + a[n':].x^(-1)
		aprime, _ = VectorScalarMul(a[:nprime], x)
		aprime2, _ = VectorScalarMul(a[nprime:], xinv)
		aprime, _ = VectorAdd(aprime, aprime2)
		// Compute b' = b[:n'].x^(-1) + b[n':].x
		bprime, _ = VectorScalarMul(b[:nprime], xinv)
		bprime2, _ = VectorScalarMul(b[nprime:], x)
		bprime, _ = VectorAdd(bprime, bprime2)

		Ls = append(Ls, L)
		Rs = append(Rs, R)
		// recursion BIP(g',h',u,P'; a', b')
		proof, _ = BIP(aprime, bprime, gprime, hprime, u, Pprime, nprime, Ls, Rs)
	}
	proof.N = n
	return proof, nil
}

/*
Verify is responsible for the verification of the Inner Product Proof.
*/
func (zkip *bip) Verify(proof proofBip) (bool, error) {

	logn := len(proof.Ls)
	var (
		i                                    int64
		x, xinv, x2, x2inv                   *big.Int
		ngprime, nhprime, ngprime2, nhprime2 []*p256
	)

	i = 0
	gprime := zkip.Gg
	hprime := zkip.Hh
	Pprime := zkip.P
	nprime := proof.N
	for i < int64(logn) {
		nprime = nprime / 2
		x, _, _ = HashBP(proof.Ls[i], proof.Rs[i])
		xinv = ModInverse(x, ORDER)
		// Compute g' = g[:n']^(x^-1) * g[n':]^(x)
		ngprime, _ = VectorScalarExp(gprime[:nprime], xinv)
		ngprime2, _ = VectorScalarExp(gprime[nprime:], x)
		gprime, _ = VectorECAdd(ngprime, ngprime2)
		// Compute h' = h[:n']^(x)    * h[n':]^(x^-1)
		nhprime, _ = VectorScalarExp(hprime[:nprime], x)
		nhprime2, _ = VectorScalarExp(hprime[nprime:], xinv)
		hprime, _ = VectorECAdd(nhprime, nhprime2)
		// Compute P' = L^(x^2).P.R^(x^-2)
		x2 = Mod(Multiply(x, x), ORDER)
		x2inv = ModInverse(x2, ORDER)
		Pprime.Multiply(Pprime, new(p256).ScalarMult(proof.Ls[i], x2))
		Pprime.Multiply(Pprime, new(p256).ScalarMult(proof.Rs[i], x2inv))
		i = i + 1
	}

	// c == a*b
	ab := Multiply(proof.A, proof.B)
	ab = Mod(ab, ORDER)

	rhs := new(p256).ScalarMult(gprime[0], proof.A)
	hb := new(p256).ScalarMult(hprime[0], proof.B)
	rhs.Multiply(rhs, hb)
	rhs.Multiply(rhs, new(p256).ScalarMult(proof.U, ab))

	nP := Pprime.Neg(Pprime)
	nP.Multiply(nP, rhs)
	c := nP.IsZero()

	return c, nil
}
