package verifycode

import (
	"strconv"
	"strings"
)

const (
	Image_Size   = 67
	Image_SpaceW = 5
	Image_SpaceH = 5
	Pic_HeadHigh = 36
)

func GetImageRowNum(w, h int) int {
	h = h - Pic_HeadHigh - Image_SpaceH
	if h > 0 {
		ret := h / (Image_Size + Image_SpaceH)
		if h%(Image_Size+Image_SpaceH) == 0 {
			return ret
		}
		return ret + 1
	}
	return 0
}

func GetImageColNum(w, h int) int {
	w = w - Image_SpaceW
	if w > 0 {
		ret := w / (Image_Size + Image_SpaceW)
		if w%(Image_Size+Image_SpaceW) == 0 {
			return ret
		}
		return ret + 1
	}
	return 0
}

func GetImageRange(picW, picH, indexR, indexC int) (startW, startH, endW, endH int) {
	rows := GetImageRowNum(picW, picH)
	cols := GetImageColNum(picW, picH)
	if indexR <= 0 || indexC <= 0 || rows < indexR || cols < indexC {
		return -1, -1, -1, -1
	}
	if indexC <= 1 {
		startW = Image_SpaceW
	} else {
		startW = (indexC-1)*Image_Size + indexC*Image_SpaceW
	}
	if indexR <= 1 {
		startH = Pic_HeadHigh + Image_SpaceH
	} else {
		startH = Pic_HeadHigh + (indexR-1)*Image_Size + indexR*Image_SpaceH
	}
	endW += startW + Image_Size
	endH += startH + Image_Size
	return
}

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
