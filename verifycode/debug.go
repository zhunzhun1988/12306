package verifycode

import (
	"fmt"
)

type DebugVerify struct {
}

func NewDebugVerify() Interface {
	return &DebugVerify{}
}
func checkIndexToPos(image string, i, j int) (x, y int) {
	switch j {
	case 1:
		x = 31
	case 2:
		x = 107
	case 3:
		x = 173
	case 4:
		x = 252
	}
	switch i {
	case 1:
		y = 49
	case 2:
		y = 109
	}
	return
}

func (dv *DebugVerify) GetAnswer(imagepath string) VerifyPosList {
	n := 0
	var poss VerifyPosList
	fmt.Printf("num:")
	fmt.Scanf("%d", &n)
	poss = make(VerifyPosList, n)
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Printf("[%d pos]:", i)
		fmt.Scanf("%d %d", &x, &y)
		poss[i].X, poss[i].Y = checkIndexToPos(imagepath, x, y)
	}
	return poss
}
