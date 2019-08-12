package main

import (
	"fmt"
	"math/big"

	"./zkproofs"
)

func main() {

	// value
	x := new(big.Int).SetInt64(30)
	y := new(big.Int).SetInt64(20)
	z := new(big.Int).SetInt64(10)
	// get blind factor, pedersen commit and zkproof
	blindFactorX, tX, hX, pX, proofX, _ := zkproofs.GetZkrp().GenerateProof(x)

	proofDataX, _ := zkproofs.DumpProof(tX, hX, pX, &proofX)

	verifier, proof, _ := zkproofs.LoadProof(proofDataX)

	var ok bool
	ok, err := verifier.Verify(*proof)
	if !ok {
		fmt.Println("proofX failed!!!")
		fmt.Println(ok)
		fmt.Println(err)
	} else {
		fmt.Println("proofX verified >0.")
	}

	blindFactorY, tY, hY, pY, proofY, _ := zkproofs.GetZkrp().GenerateProof(y)

	proofDataY, _ := zkproofs.DumpProof(tY, hY, pY, &proofY)

	verifier, proof, _ = zkproofs.LoadProof(proofDataY)
	ok, _ = verifier.Verify(proofY)
	if !ok {
		fmt.Println("proofY failed!!!")
	} else {
		fmt.Println("proofY verified >0.")
	}

	blindFactorZ, tZ, hZ, pZ, proofZ, _ := zkproofs.GetZkrp().GenerateProof(z)

	proofDataZ, _ := zkproofs.DumpProof(tZ, hZ, pZ, &proofZ)

	verifier, proof, _ = zkproofs.LoadProof(proofDataZ)
	ok, _ = verifier.Verify(proofZ)
	if !ok {
		fmt.Println("proofZ failed!!!")
	} else {
		fmt.Println("proofZ verified >0.")
	}

	// 佩德森承诺检查输入之和与输出之和是否相等
	blindOut := new(big.Int).Add(blindFactorY, blindFactorZ)
	blindDiff := new(big.Int).Sub(blindFactorX, blindOut)

	check := zkproofs.VerifyPedersenCommitment([]*zkproofs.PedersenCommitment{proofX.V}, []*zkproofs.PedersenCommitment{proofY.V, proofZ.V}, blindDiff)
	fmt.Println("pedersen verify result:", check)

}

// 最好能抽象出一个独立的 zkrp.Bp 出来，然后
// 主要有两个 challenge值 y 和 z 应该是由 verifier 那边根据A 和 S 生成的，这边把这个交互过程省略了，所以最后应该把
// zkrp.Zkip.Hh = hprime
// zkrp.Zkip.Cc = tprime
// 这两个值也放在 proof 中发过去，有个问题：G 和 H 用户是否应该知道？
