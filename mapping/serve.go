package mapping

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"mule/hexagon"
	"mule/overpower"
	"net/http"
	"strings"
)

func ServeMap(w http.ResponseWriter, mv overpower.MapView, fid int, facList []overpower.Faction, pvList []overpower.PlanetView, svList []overpower.ShipView, orders []overpower.Order) {
	planetGrid := make(map[hexagon.Coord]overpower.PlanetView, len(pvList))
	plidGrid := make(map[int]overpower.PlanetView, len(pvList))
	shipGrid := make(map[hexagon.Coord][]overpower.ShipView, len(svList))
	trailGrid := make(map[hexagon.Coord][]overpower.ShipView, len(svList))
	names := make(map[int]string, len(facList))
	for _, fac := range facList {
		names[fac.Fid()] = fac.Name()
	}
	for _, pv := range pvList {
		planetGrid[pv.Loc()] = pv
		plidGrid[pv.Pid()] = pv
	}
	for _, sv := range svList {
		if l, ok := sv.Loc(); ok {
			if list, ok := shipGrid[l]; ok {
				shipGrid[l] = append(list, sv)
			} else {
				shipGrid[l] = []overpower.ShipView{sv}
			}
		}
		for _, t := range sv.Trail() {
			if list, ok := trailGrid[t]; ok {
				trailGrid[t] = append(list, sv)
			} else {
				trailGrid[t] = []overpower.ShipView{sv}
			}
		}
	}
	// -------- //
	width, height := 800, 600
	frame := image.Rect(0, 0, width, height)
	final := image.NewRGBA(frame)
	draw.Draw(final, final.Bounds(), image.Black, image.ZP, draw.Src)
	starC := color.RGBA{50, 50, 50, 255}
	for i := 0; i < 12000; i++ {
		if i == 6000 {
			starC = color.RGBA{100, 100, 100, 255}
		} else if i == 11000 {
			starC = color.RGBA{200, 200, 200, 255}
		}
		x, y := rand.Intn(width), rand.Intn(height)
		final.Set(x, y, starC)
	}
	// -------- //
	gc := draw2dimg.NewGraphicContext(final)
	zoom := mv.Zoom()
	if zoom > 100 {
		zoom = 100
	} else if zoom < 1 {
		zoom = 1
	}
	if zoom > 40 {
		gc.SetLineWidth(.5)
	} else {
		gc.SetLineWidth(.25)
	}
	if zoom < 10 {
		gc.SetFontSize(8)
	} else {
		gc.SetFontSize(10)
	}
	gc.SetStrokeColor(color.RGBA{0x3F, 0x3F, 0x9F, 0xFF})
	draw2d.SetFontFolder("DATA")
	gc.SetFontData(draw2d.FontData{Name: "DroidSansMono", Family: draw2d.FontFamilyMono})
	//
	rad := float64(zoom)
	vp := GetVP(mv)
	plToDraw := []overpower.PlanetView{}
	shToDraw := map[hexagon.Coord][]overpower.ShipView{}
	trToDraw := map[hexagon.Coord][]overpower.ShipView{}
	focus, focValid := mv.Focus()
	var focusVis bool
	if zoom > 14 {
		for _, h := range vp.VisList() {
			if focValid && !focusVis && h == focus {
				focusVis = true
			}
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
			if pv, ok := planetGrid[h]; ok {
				plToDraw = append(plToDraw, pv)
			}
			if list, ok := shipGrid[h]; ok {
				shToDraw[h] = list
			}
			if list, ok := trailGrid[h]; ok {
				trToDraw[h] = list
			}
		}
		if focusVis {
			gc.SetLineWidth(1)
			gc.SetStrokeColor(color.RGBA{0xFF, 0xFF, 0x00, 0xFF})
			corners := vp.CornersOf(focus)
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
			if zoom > 40 {
				gc.SetLineWidth(.5)
			} else {
				gc.SetLineWidth(.25)
			}
		}
	} else {
		for _, h := range vp.VisList() {
			if pv, ok := planetGrid[h]; ok {
				plToDraw = append(plToDraw, pv)
			}
			if list, ok := shipGrid[h]; ok {
				shToDraw[h] = list
			}
			if list, ok := trailGrid[h]; ok {
				trToDraw[h] = list
			}
		}
	}
	gc.SetStrokeColor(color.Black)
	for c, list := range trToDraw {
		DrawTrails(gc, vp, fid, rad, names, c, list)
	}
	for c, list := range shToDraw {
		DrawShips(gc, vp, fid, rad, names, c, list)
	}
	//gc.SetStrokeColor(color.RGBA{0x9F, 0x9F, 0x3F, 0xFF})
	gc.SetStrokeColor(color.RGBA{0x0F, 0xAF, 0xAF, 0xFF})
	if zoom > 40 {
		gc.SetLineWidth(2)
	} else if zoom > 10 {
		gc.SetLineWidth(1)
	} else {
		gc.SetLineWidth(.5)
	}
	for _, o := range orders {
		src, okS := plidGrid[o.Source()]
		tar, okT := plidGrid[o.Target()]
		if !okS || !okT {
			// Log(BAD ORDER: PLANETS NOT FOUND", o)
			continue
		}
		if src.Controller() != fid {
			continue
		}
		srcP := vp.CenterOf(src.Loc())
		tarP := vp.CenterOf(tar.Loc())
		var yadj float64
		if zoom > 10 {
			yadj = rad * .25
		}
		gc.MoveTo(srcP[0], srcP[1]-yadj)
		gc.LineTo(tarP[0], tarP[1]-yadj)
		gc.Stroke()
	}
	if zoom > 40 {
		gc.SetLineWidth(.5)
	} else {
		gc.SetLineWidth(.25)
	}
	gc.SetStrokeColor(color.Black)
	for _, pv := range plToDraw {
		if zoom < 15 && focValid && pv.Loc() == focus {
			gc.SetFillColor(color.RGBA{0xFF, 0xFF, 0x00, 0xFF})
		} else if pv.Turn() > 0 && pv.Controller() != 0 {
			if pv.Controller() == fid {
				gc.SetFillColor(color.RGBA{0x0F, 0xFF, 0x0F, 0xFF})
			} else {
				gc.SetFillColor(color.RGBA{0xFF, 0x0F, 0x0F, 0xFF})
			}

		} else {
			gc.SetFillColor(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
		}
		h := pv.Loc()
		c := vp.CenterOf(h)
		if zoom > 10 {
			gc.FillStringAt(fmt.Sprintf("%s (%d,%d)", pv.Name(), h[0], h[1]), c[0]+(rad*.25), c[1]-(.75*rad))
			gc.ArcTo(c[0], c[1]-rad*.25, rad*.45, rad*.45, 0, -math.Pi*2)
		} else if zoom >= 5 {
			gc.FillStringAt(fmt.Sprintf("%s", pv.Name()), c[0]+(rad*.25), c[1]-(.75*rad))
			gc.ArcTo(c[0], c[1], 5, 5, 0, -math.Pi*2)
		} else {
			gc.ArcTo(c[0], c[1], 3, 3, 0, -math.Pi*2)
		}
		gc.Close()
		gc.FillStroke()
	}

	// -------- //
	png.Encode(w, final)
}

func DrawTrails(gc draw2d.GraphicContext, vp *hexagon.Viewport, fid int, rad float64, names map[int]string, c hexagon.Coord, list []overpower.ShipView) {
	trailC := color.RGBA{0x39, 0x39, 0x39, 0x3F}
	gc.SetFillColor(trailC)
	p := vp.CenterOf(c)
	if rad > 10 {
		gc.ArcTo(p[0], p[1], rad*.45, rad*.25, 0, -math.Pi*2)
	} else if rad >= 5 {
		gc.ArcTo(p[0], p[1], 5, 3, 0, -math.Pi*2)
	} else {
		gc.ArcTo(p[0], p[1], 3, 1, 0, -math.Pi*2)
	}
	gc.Close()
	gc.FillStroke()
}

func DrawShips(gc draw2d.GraphicContext, vp *hexagon.Viewport, fid int, rad float64, names map[int]string, c hexagon.Coord, list []overpower.ShipView) {
	var anyYours bool

	parts := make([]string, len(list))
	for i, sv := range list {
		if sv.Controller() == fid {
			anyYours = true
			parts[i] = fmt.Sprintf("YOURS(%d)", sv.Size())
		} else {
			parts[i] = fmt.Sprintf("%d(%d)", names[sv.Controller()], sv.Size())
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
	if rad > 10 {
		gc.ArcTo(p[0], p[1], rad*.45, rad*.25, 0, -math.Pi*2)
	} else if rad >= 5 {
		gc.ArcTo(p[0], p[1], 5, 3, 0, -math.Pi*2)
	} else {
		gc.ArcTo(p[0], p[1], 3, 1, 0, -math.Pi*2)
	}
	gc.Close()
	gc.FillStroke()
	if rad >= 5 {
		gc.FillStringAt(name, p[0]+(rad*.55), p[1]+3)
	}
}

func GetVP(mv overpower.MapView) *hexagon.Viewport {
	vp := hexagon.MakeViewport(float64(mv.Zoom()), false, true)
	center := mv.Center()
	vp.SetAnchor(center[0], center[1], 400.0, 300.0)
	vp.SetFrame(0, 0, 800, 600)
	return vp
}
