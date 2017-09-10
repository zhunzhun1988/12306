package verifycode

import (
	"strconv"
	"strings"
)

type VerifyPos struct {
	X int
	Y int
}

type VerifyPosList []VerifyPos

func (vpl VerifyPosList) Len() int {
	return len(vpl)
}
func (vpl VerifyPosList) Swap(i, j int) {
	vpl[i], vpl[j] = vpl[j], vpl[i]
}
func (vpl VerifyPosList) Less(i, j int) bool {
	if vpl[i].X != vpl[j].X {
		return vpl[i].X < vpl[j].X
	}
	return vpl[i].Y < vpl[j].Y
}

func (vpl VerifyPosList) ToString() string {
	strs := make([]string, 0, len(vpl)*2)
	for _, pv := range vpl {
		strs = append(strs, strconv.Itoa(pv.X))
		strs = append(strs, strconv.Itoa(pv.Y))
	}
	return strings.Join(strs, ",")
}

type Interface interface {
	GetAnswer(imagepath string) VerifyPosList
}

type VerifierList []Interface

func (vl VerifierList) GetAnswer(imagepath string) VerifyPosList {
	for _, verifier := range vl {
		ret := verifier.GetAnswer(imagepath)
		if len(ret) > 0 {
			return ret
		}
	}
	return VerifyPosList{}
}
