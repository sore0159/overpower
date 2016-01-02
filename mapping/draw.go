package mapping

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"image/color"
	"math"
	"mule/hexagon"
	"mule/overpower"
	"strings"
)

func DrawPlanet(gc draw2d.GraphicContext, vp *hexagon.Viewport, fid int, avail int, showGrid, isFocus, isCenter bool, pv overpower.PlanetView) {
	if !showGrid && isCenter {
		gc.SetFillColor(color.RGBA{0xFF, 0xFF, 0x00, 0xFF})
	} else if !showGrid && isFocus {
		gc.SetFillColor(color.RGBA{0x99, 0x99, 0x00, 0xFF})
	} else if cont := pv.Controller(); pv.Turn() > 0 && cont != 0 {
		if cont == fid {
			gc.SetFillColor(color.RGBA{0x0F, 0xFF, 0x0F, 0xFF})
		} else {
			gc.SetFillColor(color.RGBA{0xFF, 0x0F, 0x0F, 0xFF})
		}
	} else {
		gc.SetFillColor(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	}
	c := vp.CenterOf(pv.Loc())
	rad := vp.HexR
	if showGrid {
		c[1] -= rad * .25
	}
	var size float64
	if rad > 10 {
		size = rad * .45
	} else if rad > 4 {
		size = 5
	} else {
		size = 3
	}
	var plStr string
	if avail > 0 {
		plStr = fmt.Sprintf("(%d)%s", avail, pv.Name())
	} else {
		plStr = fmt.Sprintf("%s", pv.Name())
	}
	//
	if rad > 4 {
		gc.FillStringAt(plStr, c[0]+(rad*.25), c[1]-(.75*rad))
	} else if rad > 2 && avail > 0 {
		plStr := fmt.Sprintf("(%d)", avail)
		gc.FillStringAt(plStr, c[0]+(rad*.25), c[1]-(.75*rad))
	}
	gc.ArcTo(c[0], c[1], size, size, 0, -math.Pi*2)
	gc.Close()
	gc.FillStroke()
}

func DrawLine(gc draw2d.GraphicContext, vp *hexagon.Viewport, h1, h2 hexagon.Coord) {
	p1 := vp.CenterOf(h1)
	p2 := vp.CenterOf(h2)
	gc.MoveTo(p1[0], p1[1])
	gc.LineTo(p2[0], p2[1])
	gc.Stroke()

}

func DrawDestLine(gc draw2d.GraphicContext, vp *hexagon.Viewport, h1, h2 hexagon.Coord, showGrid bool) {
	p1 := vp.CenterOf(h1)
	p2 := vp.CenterOf(h2)
	var yadj float64
	if showGrid {
		yadj = vp.HexR * .25
	}
	gc.MoveTo(p1[0], p1[1])
	gc.LineTo(p2[0], p2[1]-yadj)
	gc.Stroke()

}

func DrawOrderLine(gc draw2d.GraphicContext, vp *hexagon.Viewport, h1, h2 hexagon.Coord, showGrid bool) {
	p1 := vp.CenterOf(h1)
	p2 := vp.CenterOf(h2)
	var yadj float64
	if showGrid {
		yadj = vp.HexR * .25
	}
	gc.MoveTo(p1[0], p1[1]-yadj)
	gc.LineTo(p2[0], p2[1]-yadj)
	gc.Stroke()

}

func DrawHex(gc draw2d.GraphicContext, vp *hexagon.Viewport, h hexagon.Coord) {
	corners := vp.CornersOf(h)
	gc.MoveTo(corners[0][0], corners[0][1])
	for i, _ := range corners {
		var px, py float64
		if i == 5 {
			px, py = corners[0][0], corners[0][1]
		} else {
			px, py = corners[i+1][0], corners[i+1][1]
		}
		gc.LineTo(px, py)
	}
	gc.Close()
	gc.Stroke()
}

func DrawTrailDot(gc draw2d.GraphicContext, vp *hexagon.Viewport, c hexagon.Coord) {
	p := vp.CenterOf(c)
	rad := vp.HexR
	gc.ArcTo(p[0], p[1], rad*.45, rad*.25, 0, -math.Pi*2)
	gc.Close()
	gc.FillStroke()
}

func DrawShips(gc draw2d.GraphicContext, vp *hexagon.Viewport, fid int, names map[int]string, c hexagon.Coord, list []overpower.ShipView) {
	rad := vp.HexR
	var anyYours bool

	parts := make([]string, len(list))
	for i, sv := range list {
		if sv.Controller() == fid {
			anyYours = true
			parts[i] = fmt.Sprintf("YOURS(%d)", sv.Size())
		} else {
			parts[i] = fmt.Sprintf("%s(%d)", names[sv.Controller()], sv.Size())
		}
	}
	var name string
	if len(list) > 2 {
		name = fmt.Sprintf("%d SHIPS", len(list))
	} else {
		name = strings.Join(parts, "||")
	}
	var shipC color.RGBA
	if anyYours {
		shipC = color.RGBA{0x09, 0xF9, 0x09, 0xFF}
	} else {
		shipC = color.RGBA{0xF9, 0x09, 0x09, 0xFF}
	}
	gc.SetFillColor(shipC)
	p := vp.CenterOf(c)
	var sizeX, sizeY float64
	if rad > 10 {
		sizeX, sizeY = rad*.45, rad*.25
	} else if rad > 4 {
		sizeX, sizeY = 5, 3
	} else {
		sizeX, sizeY = 3, 1
	}
	gc.ArcTo(p[0], p[1], sizeX, sizeY, 0, -math.Pi*2)
	gc.Close()
	gc.FillStroke()
	if rad >= 5 {
		gc.FillStringAt(name, p[0]+(rad*.55), p[1]+3)
	}
}
