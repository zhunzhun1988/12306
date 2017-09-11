package verifycode

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Image_Size   = 67
	Image_SpaceW = 5
	Image_SpaceH = 5
	Pic_HeadHigh = 36
)

type ImageRange struct {
	sw, sh int
	ew, eh int
}

func GetImageRowNum(w, h int) int {
	h = h - Pic_HeadHigh - Image_SpaceH*2
	if h > 0 {
		ret := h / (Image_Size + Image_SpaceH)
		if h%(Image_Size+Image_SpaceH) == 0 {
			return ret
		}
		fmt.Printf("%d, %d\n", h, Image_Size+Image_SpaceH)
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
func GetOKRange(picW, picH int) (startW, startH, endW, endH int) {
	startW = picW - 30
	endW = picW
	startH = 0
	endH = Pic_HeadHigh - Image_SpaceH
	return
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

func GetImageRanges(picW, picH int) [][]ImageRange {
	rows := GetImageRowNum(picW, picH)
	cols := GetImageColNum(picW, picH)
	ranges := make([][]ImageRange, rows+1)
	for i := 1; i <= rows; i++ {
		ranges[i] = make([]ImageRange, cols+1)
		for j := 1; j <= cols; j++ {
			ranges[i][j].sw, ranges[i][j].sh, ranges[i][j].ew, ranges[i][j].eh =
				GetImageRange(picW, picH, i, j)
		}
	}
	return ranges
}

func GetIndexByPos(picW, picH, x, y int) (r, c int) {
	ranges := GetImageRanges(picW, picH)
	for i, row := range ranges {
		for j, cel := range row {
			if cel.sw < x && x < cel.ew && cel.sh < y && y < cel.eh {
				return i, j
			}
		}
	}
	return -1, -1
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
