package mapping

import (
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
	//"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"mule/mylog"
	"mule/planetattack/attack"
	"os"
	"strconv"
)

var (
	Log      = mylog.Err
	FONTFILE = "DroidSansMono.ttf"
	MAPDIR   = ""
)

func SetFont(fileName string) {
	FONTFILE = fileName
}
func SetMapDir(dirName string) {
	MAPDIR = dirName
}

func MakeMap(f *attack.Faction) (string, error) {
	sv := f.View
	//numF := len(f.OtherNames)
	galaxyRad := sv.Size
	hexR := 10
	size := 2 * (galaxyRad + 5) * hexR * 2
	Plane := func(in [2]int) (out [2]int) {
		out = attack.Hex2Plane(hexR, in)
		out = [2]int{out[0] + size/2, size/2 - out[1]}
		return
	}
	starMap := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{size, size}})
	// Void //
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(starMap, starMap.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)
	// Stars //
	numStars := galaxyRad * galaxyRad / 3
	for i := 0; i < numStars; i++ {
		var starColor color.RGBA
		switch rand.Intn(3) {
		case 1:
			starColor = color.RGBA{uint8(rand.Intn(55)), uint8(0), 50 + uint8(rand.Intn(155)), 255}
		case 2:
			bright := 150 + uint8(rand.Intn(105))
			starColor = color.RGBA{bright, bright, bright, 255}
		default:
			starColor = color.RGBA{50 + uint8(rand.Intn(155)), uint8(0), uint8(rand.Intn(55)), 255}
		}
		starX := rand.Intn(size)
		starY := rand.Intn(size)
		starRadius := 1 + rand.Intn(1)
		DrawBlip(starMap, starColor, [2]int{starX, starY}, starRadius, size)
	}
	// ShipTrails //
	trailR := 7
	for loc, trList := range sv.TrailGrid {
		_ = trList
		tColor := color.RGBA{55, 55, 55, 255}
		DrawBlip(starMap, tColor, Plane(loc), trailR, size)
	}
	// Planets //
	plR := 5
	var plColor color.RGBA
	for loc, plv := range sv.PlanetGrid {
		if plv.Yours {
			plColor = Color(true, false)
		} else if plv.Inhabitants[0] != 0 {
			plColor = Color(false, true)
		} else {
			plColor = Color(false, false)
		}
		DrawBlip(starMap, plColor, Plane(loc), plR, size)
	}
	// Planet Names //
	typer := freetype.NewContext()
	typer.SetClip(starMap.Bounds())
	typer.SetDst(starMap)
	fontBytes, err := ioutil.ReadFile(FONTFILE)
	if err != nil {
		return "", Log(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return "", Log(err)
	}
	typer.SetFont(font)
	typer.SetFontSize(10)
	for loc, plv := range sv.PlanetGrid {
		if plv.Yours {
			plColor = Color(true, false)
		} else if plv.Inhabitants[0] != 0 {
			plColor = Color(false, true)
		} else {
			plColor = Color(false, false)
		}
		typer.SetSrc(&image.Uniform{plColor})
		pixels := Plane(loc)
		pixels[1] -= (hexR + 2)
		pt := freetype.Pt(pixels[0], pixels[1])
		_, err := typer.DrawString(fmt.Sprintf("%s (%d)", plv.Name, plv.ID), pt)
		if err != nil {
			return "", Log(err)
		}
	}
	// Ships //
	shipR := 5
	for loc, sList := range sv.ShipGrid {
		var useR int
		if len(sList) > 1 {
			useR = shipR + 2
		} else {
			useR = shipR
		}
		sColor := Color(true, false)
		shStr := ""
		for _, sh := range sList {
			if !sh.Yours {
				shStr += f.OtherNames[sh.FactionID] + ":" + strconv.Itoa(sh.Size) + " "
				sColor = Color(false, true)
			} else {
				shStr += "Yours:" + strconv.Itoa(sh.Size) + " "
			}
		}
		DrawBlip(starMap, sColor, Plane(loc), useR, size)
		typer.SetSrc(&image.Uniform{sColor})
		pixels := Plane(loc)
		pixels[1] += (2 * hexR)
		pt := freetype.Pt(pixels[0], pixels[1])
		_, err := typer.DrawString(shStr, pt)
		if err != nil {
			return "", Log(err)
		}
	}
	fileName := Filename(f.Name, f.View.Turn)
	starfile, err := os.Create(fileName)
	if err != nil {
		Log(err)
		return "", err
	}
	png.Encode(starfile, starMap)
	//jpeg.Encode(starfile, starMap, nil)
	return fileName, nil
}

func DrawBlip(img *image.RGBA, color color.RGBA, center [2]int, rad, imgSize int) {
	for x := center[0] - rad; x < center[0]+rad; x++ {
		if x < 0 || x > imgSize {
			continue
		}
		for y := center[1] - rad; y < center[1]+rad; y++ {
			if y < 0 || y > imgSize {
				continue
			}
			/*if rad > 3 && abs(y-center[1])+abs(x-center[0]) >= rad {
				continue
			}*/
			img.Set(x, y, color)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	} else {
		return x
	}
}

func Color(yours bool, enemy bool) color.RGBA {
	if yours {
		return color.RGBA{105, 255, 105, 255}
	}
	if enemy {
		return color.RGBA{255, 105, 105, 255}
	}
	return color.RGBA{255, 255, 255, 255}
}

func Filename(name string, turn int) string {
	return fmt.Sprintf("%s%03d_%s.png", MAPDIR, turn, name)
	//return fmt.Sprintf("%s%03d_%s.jpeg", MAPDIR, turn, name)
}
