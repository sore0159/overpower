package mapping

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"mule/hexagon"
	"mule/overpower"
	"net/http"
)

const (
	MAPW = 800
	MAPH = 600
)

func ServeMap(w http.ResponseWriter, mv overpower.MapView, fid int, facList []overpower.Faction, pvList []overpower.PlanetView, svList []overpower.ShipView, orders []overpower.Order) {
	width, height := MAPW, MAPH
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
	gridC := color.RGBA{0x3F, 0x3F, 0x9F, 0xFF}
	focusC := color.RGBA{0x99, 0x99, 0x00, 0xFF}
	selectC := color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
	orderC := color.RGBA{0x0F, 0xAF, 0xAF, 0xFF}
	trailDotC := color.RGBA{0x39, 0x39, 0x39, 0x3F}
	trailLineC := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	//destLineC := color.RGBA{0x0F, 0xFF, 0x0F, 0xFF}
	destLineC := orderC
	_ = trailDotC
	gc := draw2dimg.NewGraphicContext(final)
	draw2d.SetFontFolder("DATA")
	gc.SetFontData(draw2d.FontData{Name: "DroidSansMono", Family: draw2d.FontFamilyMono})
	zoom := mv.Zoom()
	if zoom > 100 {
		zoom = 100
	} else if zoom < 1 {
		zoom = 1
	}
	if zoom < 10 {
		gc.SetFontSize(8)
	} else {
		gc.SetFontSize(10)
	}
	if zoom > 40 {
		gc.SetLineWidth(.5)
	} else {
		gc.SetLineWidth(.25)
	}
	//
	showGrid := zoom > 14
	vp := GetVP(mv)
	center := mv.Center()
	focus, focValid := mv.Focus()
	visList := vp.VisList()
	visMap := make(map[hexagon.Coord]bool, len(visList))
	gc.SetStrokeColor(gridC)
	// ------------ GRID + SELECT + FOCUS ------------- //
	for _, h := range vp.VisList() {
		visMap[h] = true
		if showGrid {
			DrawHex(gc, vp, h)
		}
	}
	if showGrid {
		gc.SetLineWidth(1)
		if focValid && visMap[focus] {
			gc.SetStrokeColor(focusC)
			DrawHex(gc, vp, focus)
		}
		gc.SetStrokeColor(selectC)
		DrawHex(gc, vp, center)
	}
	// -------- PLANETS PREP ----------- //
	plidGrid := make(map[int]overpower.PlanetView, len(pvList))
	plToDraw := make([]overpower.PlanetView, 0, len(pvList))
	availMap := map[int]int{}
	for _, pv := range pvList {
		pid := pv.Pid()
		plidGrid[pid] = pv
		if visMap[pv.Loc()] {
			plToDraw = append(plToDraw, pv)
		}
		if pv.Controller() == fid {
			availMap[pid] = pv.Parts()
		}
	}
	// ------------------- SHIPS ---------------------- //
	names := make(map[int]string, len(facList))
	for _, fac := range facList {
		names[fac.Fid()] = fac.Name()
	}
	shToDraw := make(map[hexagon.Coord][]overpower.ShipView, 0)
	destToDraw := make([][2]hexagon.Coord, 0)
	trailToDraw := make([][2]hexagon.Coord, 0)
	trailFids := make([]int, 0)
	trailDots := make([]hexagon.Coord, 0)
	for _, sv := range svList {
		var locVis, locOk bool
		loc, locOk := sv.Loc()
		if locOk {
			locVis = visMap[loc]
			if locVis {
				if list, ok := shToDraw[loc]; ok {
					shToDraw[loc] = append(list, sv)
				} else {
					shToDraw[loc] = []overpower.ShipView{sv}
				}
			}
			if dest, ok := sv.Dest(); ok && (locVis || visMap[dest]) {
				destToDraw = append(destToDraw, [2]hexagon.Coord{loc, dest})
			}
		}
		trail := sv.Trail()
		for _, test := range trail {
			if visMap[test] {
				if showGrid {
					trailDots = append(trailDots, test)
				} else {
					var end hexagon.Coord
					if locOk {
						end = loc
					} else {
						end = trail[len(trail)-1]
					}
					trailToDraw = append(trailToDraw, [2]hexagon.Coord{trail[0], end})
					trailFids = append(trailFids, sv.Controller())
					break
				}
			}
		}
	}
	if showGrid {
		if len(trailDots) > 0 {
			gc.SetLineWidth(.25)
			gc.SetStrokeColor(color.Black)
			gc.SetFillColor(trailDotC)
			for _, c := range trailDots {
				DrawTrailDot(gc, vp, c)
			}
		}
	} else {
		gc.SetLineWidth(1)
		gc.SetFillColor(trailLineC)
		for i, pts := range trailToDraw {
			if pts[0] != pts[1] {
				cont := trailFids[i]
				if cont == fid {
					gc.SetStrokeColor(color.RGBA{0x0F, 0xFF, 0x0F, 0xFF})
				} else {
					gc.SetStrokeColor(color.RGBA{0xFF, 0x0F, 0x0F, 0xFF})
				}
				DrawLine(gc, vp, pts[0], pts[1])
			} else {
				gc.SetLineWidth(0)
				DrawTrailDot(gc, vp, pts[0])
				gc.SetLineWidth(1)
			}
		}
	}
	gc.SetLineWidth(1)
	gc.SetStrokeColor(destLineC)
	for _, pts := range destToDraw {
		DrawDestLine(gc, vp, pts[0], pts[1], showGrid)
	}
	for c, list := range shToDraw {
		DrawShips(gc, vp, fid, names, c, list)
	}
	// ------------ ORDERS ------------- //
	gc.SetStrokeColor(orderC)
	/*
		if zoom > 40 {
			gc.SetLineWidth(2)
		} else if zoom > 10 {
			gc.SetLineWidth(1)
		} else {
			gc.SetLineWidth(.5)
		}
	*/
	gc.SetLineWidth(2)
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
		availMap[src.Pid()] -= o.Size()
		DrawOrderLine(gc, vp, src.Loc(), tar.Loc(), showGrid)
	}
	// -------------------- PLANETS --------------------- //
	gc.SetLineWidth(.25)
	gc.SetStrokeColor(color.Black)
	for _, pv := range plToDraw {
		loc := pv.Loc()
		isFocus := focValid && (focus == loc)
		isCenter := loc == center
		var avail int
		if pv.Controller() == fid {
			avail = availMap[pv.Pid()]
		} else {
			avail = pv.Inhabitants()
		}
		DrawPlanet(gc, vp, fid, avail, showGrid, isFocus, isCenter, pv)
	}

	// -------- //
	png.Encode(w, final)
}
