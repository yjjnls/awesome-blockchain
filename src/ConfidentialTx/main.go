package main

import (
	"fmt"
	"math/big"

	"./zkproofs"
)

func main() {
	var (
		zkrpX zkproofs.Bp
		zkrpY zkproofs.Bp
		zkrpZ zkproofs.Bp
	)
	// common setup
	zkrpX.Setup(0, 4294967296)
	zkrpY.Setup(0, 4294967296)
	zkrpZ.Setup(0, 4294967296)

	// value
	x := new(big.Int).SetInt64(30)
	y := new(big.Int).SetInt64(20)
	z := new(big.Int).SetInt64(10)
	// get blind factor, pedersen commit and zkproof
	// 这里会更新 zkrpX，所以 zkrpX 和 proofX 都需要保存下来
	blindFactorX, commitmentX, proofX, _ := zkrpX.Prove(x)

	// zkrp, _ := zkproofs.LoadParamFromDisk("setup.dat")
	var ok bool
	// proofX中包含commitmentX，如果修改了proofX.V 验证也不会通过
	ok, err := zkrpX.Verify(proofX)
	if !ok {
		fmt.Println("proofX failed!!!")
		fmt.Println(ok)
		fmt.Println(err)
	}

	blindFactorY, commitmentY, proofY, _ := zkrpY.Prove(y)

	// zkrp, _ = zkproofs.LoadParamFromDisk("setup.dat")
	ok, _ = zkrpY.Verify(proofY)
	if !ok {
		fmt.Println("proofY failed!!!")
	}

	blindFactorZ, commitmentZ, proofZ, _ := zkrpZ.Prove(z)
	// zkrp, _ = zkproofs.LoadParamFromDisk("setup.dat")
	ok, err = zkrpZ.Verify(proofZ)
	if !ok {
		fmt.Println("proofZ failed!!!")
	}
	blindOut := new(big.Int).Add(blindFactorY, blindFactorZ)
	blindDiff := new(big.Int).Sub(blindFactorX, blindOut)

	commitmentOut := commitmentY.Add(commitmentY, commitmentZ)
	commitmentDiff := commitmentX.Add(commitmentX, commitmentOut.Neg(commitmentOut))

	fmt.Printf("blind diff: %d\n", blindDiff)
	fmt.Printf("pedersen commitment diff: ( %d , %d )\n", commitmentDiff.X, commitmentDiff.Y)

	check := zkproofs.Mult(zkrpX.H, blindDiff)
	fmt.Printf("check               diff: ( %d , %d )\n", check.X, check.Y)
	fmt.Println(check.X.Cmp(commitmentDiff.X) == 0)
	fmt.Println(check.Y.Cmp(commitmentDiff.Y) == 0)

}
