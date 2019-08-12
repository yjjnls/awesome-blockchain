package main

import (
	"fmt"
	"math/big"
	"testing"

	"./zkproofs"
	"github.com/stretchr/testify/assert"
)

func TestCTx(t *testing.T) {
	// value
	x := new(big.Int).SetInt64(30)
	y := new(big.Int).SetInt64(20)
	z := new(big.Int).SetInt64(10)

	blindFactorX, tX, hX, pX, proofX, _ := zkproofs.GetZkrp().GenerateProof(x)
	var ok bool
	verifier := zkproofs.GetVerifier(tX, hX, pX)
	ok, _ = verifier.Verify(proofX)
	assert.True(t, ok)

	blindFactorY, tY, hY, pY, proofY, _ := zkproofs.GetZkrp().GenerateProof(y)
	verifier = zkproofs.GetVerifier(tY, hY, pY)
	ok, _ = verifier.Verify(proofY)
	assert.True(t, ok)

	blindFactorZ, tZ, hZ, pZ, proofZ, _ := zkproofs.GetZkrp().GenerateProof(z)
	verifier = zkproofs.GetVerifier(tZ, hZ, pZ)
	ok, _ = verifier.Verify(proofZ)
	assert.True(t, ok)

	// 佩德森承诺检查输入之和与输出之和是否相等
	blindOut := new(big.Int).Add(blindFactorY, blindFactorZ)
	blindDiff := new(big.Int).Sub(blindFactorX, blindOut)

	check := zkproofs.VerifyPedersenCommitment([]*zkproofs.PedersenCommitment{proofX.V}, []*zkproofs.PedersenCommitment{proofY.V, proofZ.V}, blindDiff)
	fmt.Println("pedersen verify result:", check)
	assert.True(t, check)
}

func TestCTx2(t *testing.T) {
	// value
	x := new(big.Int).SetInt64(30)
	y := new(big.Int).SetInt64(40)
	z := new(big.Int).SetInt64(10)

	blindFactorX, tX, hX, pX, proofX, _ := zkproofs.GetZkrp().GenerateProof(x)
	var ok bool
	verifier := zkproofs.GetVerifier(tX, hX, pX)
	ok, _ = verifier.Verify(proofX)
	assert.True(t, ok)

	blindFactorY, tY, hY, pY, proofY, _ := zkproofs.GetZkrp().GenerateProof(y)
	verifier = zkproofs.GetVerifier(tY, hY, pY)
	ok, _ = verifier.Verify(proofY)
	assert.True(t, ok)

	blindFactorZ, tZ, hZ, pZ, proofZ, _ := zkproofs.GetZkrp().GenerateProof(z)
	verifier = zkproofs.GetVerifier(tZ, hZ, pZ)
	ok, _ = verifier.Verify(proofZ)
	assert.True(t, ok)

	// 佩德森承诺检查输入之和与输出之和是否相等
	blindOut := new(big.Int).Add(blindFactorY, blindFactorZ)
	blindDiff := new(big.Int).Sub(blindFactorX, blindOut)

	check := zkproofs.VerifyPedersenCommitment([]*zkproofs.PedersenCommitment{proofX.V}, []*zkproofs.PedersenCommitment{proofY.V, proofZ.V}, blindDiff)
	fmt.Println("pedersen verify result:", check)
	assert.False(t, check)
}

func TestCTx3(t *testing.T) {
	// value
	x := new(big.Int).SetInt64(30)
	y := new(big.Int).SetInt64(40)
	z := new(big.Int).SetInt64(-10)

	blindFactorX, tX, hX, pX, proofX, _ := zkproofs.GetZkrp().GenerateProof(x)
	var ok bool
	verifier := zkproofs.GetVerifier(tX, hX, pX)
	ok, _ = verifier.Verify(proofX)
	assert.True(t, ok)

	blindFactorY, tY, hY, pY, proofY, _ := zkproofs.GetZkrp().GenerateProof(y)
	verifier = zkproofs.GetVerifier(tY, hY, pY)
	ok, _ = verifier.Verify(proofY)
	assert.True(t, ok)

	blindFactorZ, tZ, hZ, pZ, proofZ, _ := zkproofs.GetZkrp().GenerateProof(z)
	verifier = zkproofs.GetVerifier(tZ, hZ, pZ)
	ok, _ = verifier.Verify(proofZ)
	assert.False(t, ok)

	// 佩德森承诺检查输入之和与输出之和是否相等
	blindOut := new(big.Int).Add(blindFactorY, blindFactorZ)
	blindDiff := new(big.Int).Sub(blindFactorX, blindOut)

	check := zkproofs.VerifyPedersenCommitment([]*zkproofs.PedersenCommitment{proofX.V}, []*zkproofs.PedersenCommitment{proofY.V, proofZ.V}, blindDiff)
	fmt.Println("pedersen verify result:", check)
	assert.True(t, check)
}
