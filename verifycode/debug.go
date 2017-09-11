package verifycode

import (
	"12306/opencv"
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
	image := opencv.LoadImage(imagepath)
	if image == nil {
		panic("LoadImage test.jpg fail  ")
	}
	imageSave := opencv.LoadImage(imagepath)
	defer image.Release()
	defer imageSave.Release()

	mat := image.GetMat()
	savemat := imageSave.GetMat()

	selected := make(map[int]bool)
	selectOk := false
	win := opencv.NewWindow("12306 Verify Code")
	defer win.Destroy()

	okStartW, okStartH, okEndW, okEndH := GetOKRange(mat.Cols(), mat.Rows())
	win.SetMouseCallback(func(event, x, y, flags int) {
		if event == 1 {
			if x >= okStartW && x <= okEndW && y >= okStartH && y <= okEndH {
				selectOk = true
			}
			r, c := GetIndexByPos(mat.Cols(), mat.Rows(), x, y)
			if r <= 0 || c <= 0 {
				return
			}
			sw, sh, ew, eh := GetImageRange(mat.Cols(), mat.Rows(), r, c)
			chosed := selected[r*10+c]
			for i := sw + (ew-sw)/2 - 10; i < sw+(ew-sw)/2+10; i++ {
				for j := sh + (eh-sh)/2 - 10; j < sh+(eh-sh)/2+10; j++ {
					if chosed == false {
						mat.Set2D(j, i, opencv.NewScalar(255, 0, 0, 0))
					} else {
						s := savemat.Get2D(j, i)
						mat.Set2D(j, i, opencv.NewScalar(s.Val()[0], s.Val()[1], s.Val()[2], 0))
					}
				}

			}
			selected[r*10+c] = !chosed
			win.ShowImage(image)
		}
	})
	win.AddText(image, "OK", mat.Cols()-30, 20)
	win.ShowImage(image)
	for {
		opencv.WaitKey(100)
		if selectOk == true {
			break
		}
	}
	var poss VerifyPosList
	poss = make(VerifyPosList, len(selected))
	i := 0
	for s := range selected {
		x, y := s/10, s%10
		poss[i].X, poss[i].Y = checkIndexToPos(imagepath, x, y)
		i++
	}
	return poss
}
