package verifycode

import (
	"12306/opencv"
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

	//_, currentfile, _, _ := runtime.Caller(0)
	//filename := path.Join(path.Dir(currentfile), imagepath)

	image := opencv.LoadImage(imagepath)
	if image == nil {
		panic("LoadImage test.jpg fail  ")
	}
	defer image.Release()

	mat := image.GetMat()
	for i := 0; i < mat.Rows(); i++ {
		for j := 0; j < mat.Cols(); j++ {
			s := mat.Get2D(i, j)
			fmt.Sprintf("%d ", s.Val()[0])
		}
	}

	win := opencv.NewWindow("Go-OpenCV")
	defer win.Destroy()

	win.SetMouseCallback(func(event, x, y, flags int) {
		//fmt.Printf("event = %d, x = %d, y = %d, flags = %d\n",
		//	event, x, y, flags,
		//)
	})
	win.CreateTrackbar("Thresh", 1, 100, func(pos int) {
		fmt.Printf("pos = %d\n", pos)
	})

	win.ShowImage(image)

	opencv.WaitKey(0)

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
