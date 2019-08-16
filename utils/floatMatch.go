package utils

import (
	"math"
	"math/big"
)

//是否相等
func Float64IsEqual(f1, f2 float64, precise float64) bool {
	return math.Dim(f1, f2) > precise
}

//获取指定精度浮点 最大精度15位
func FloatGetPrec(precNum int) float64 {
	switch precNum {
	case 1:
		return 0.1
	case 2:
		return 0.01
	case 3:
		return 0.001
	case 4:
		return 0.0001
	case 5:
		return 0.00001
	case 6:
		return 0.000001
	case 7:
		return 0.0000001
	case 8:
		return 0.00000001
	case 9:
		return 0.000000001
	case 10:
		return 0.0000000001
	case 11:
		return 0.00000000001
	case 12:
		return 0.000000000001
	case 13:
		return 0.0000000000001
	case 14:
		return 0.00000000000001
	case 15:
		return 0.000000000000001
	}
	return 0
}

//加法
func Float64Add(f1, f2 float64) float64 {
	fb1 := big.NewFloat(f1)
	fb2 := big.NewFloat(f2)
	fb3 := big.NewFloat(0)
	fb3.Add(fb1, fb2)
	cont, _ := fb3.Float64()

	return cont
}

//减法
func Float64Sub(f1, f2 float64) float64 {
	fb1 := big.NewFloat(f1)
	fb2 := big.NewFloat(f2)
	fb3 := big.NewFloat(0)
	fb3.Sub(fb1, fb2)
	cont, _ := fb3.Float64()
	return cont
}

//乘法
func Float64Mul(f1, f2 float64) float64 {
	fb1 := big.NewFloat(f1)
	fb2 := big.NewFloat(f2)
	fb3 := big.NewFloat(0)
	fb3.Mul(fb1, fb2)
	cont, _ := fb3.Float64()
	return cont
}

//除法
func Float64Quo(f1, f2 float64) float64 {
	if f2 == 0 {
		return 0
	}
	fb1 := big.NewFloat(f1)
	fb2 := big.NewFloat(f2)
	fb3 := big.NewFloat(0)
	fb3.Quo(fb1, fb2)
	cont, _ := fb3.Float64()
	return cont
}
