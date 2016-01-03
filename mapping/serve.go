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
	starC := color.RGBA{25, 25, 25, 255}
	for i := 0; i < 12000; i++ {
		if i == 6000 {
			starC = color.RGBA{50, 50, 50, 255}
		} else if i == 9000 {
			starC = color.RGBA{100, 100, 100, 255}
		} else if i == 11000 {
			starC = color.RGBA{200, 200, 200, 255}
		}
		x, y := rand.Intn(width), rand.Intn(height)
		final.Set(x, y, starC)
	}
	// -------- //
	gridC := color.RGBA{0x3F, 0x3F, 0x9F, 0xFF}
	target1C := color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
	target2C := color.RGBA{0x99, 0x99, 0x00, 0xFF}
	orderC := color.RGBA{0x0F, 0xAF, 0xAF, 0xFF}
	ownTrailDotC := color.RGBA{0x0F, 0x4F, 0x0F, 0x3F}
	enTrailDotC := color.RGBA{0x4F, 0x0F, 0x0F, 0x3F}
	destLineC := orderC
	//trailDotC := color.RGBA{0x39, 0x39, 0x39, 0x3F}
	//destLineC := color.RGBA{0x0F, 0xFF, 0x0F, 0xFF}
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
	target1, target2 := mv.Target1(), mv.Target2()
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
		if target2.Valid && !target2.Eq(target1) && visMap[target2.Coord] {
			gc.SetStrokeColor(target2C)
			DrawHex(gc, vp, target2.Coord)
		}
		if target1.Valid && visMap[target1.Coord] {
			gc.SetStrokeColor(target1C)
			DrawHex(gc, vp, target1.Coord)
		}
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
	trailDots := map[hexagon.Coord]bool{}
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
					if trailDots[test] {
						continue
					}
					trailDots[test] = sv.Controller() != fid
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
			for c, enemy := range trailDots {
				if enemy {
					gc.SetFillColor(enTrailDotC)
				} else {
					gc.SetFillColor(ownTrailDotC)
				}
				DrawTrailDot(gc, vp, c)
			}
		}
	} else {
		gc.SetLineWidth(1)
		if zoom > 5 {
			gc.SetLineDash([]float64{8, 8}, 0)
		} else if zoom > 1 {
			gc.SetLineDash([]float64{3, 3}, 0)
		}
		for i, pts := range trailToDraw {
			cont := trailFids[i]
			if pts[0] != pts[1] {
				if cont == fid {
					gc.SetStrokeColor(color.RGBA{0x0F, 0xFF, 0x0F, 0xFF})
				} else {
					gc.SetStrokeColor(color.RGBA{0xFF, 0x0F, 0x0F, 0xFF})
				}
				DrawLine(gc, vp, pts[0], pts[1])
			} else {
				gc.SetLineWidth(0)
				if cont == fid {
					gc.SetFillColor(ownTrailDotC)
				} else {
					gc.SetFillColor(enTrailDotC)
				}
				DrawTrailDot(gc, vp, pts[0])
				gc.SetLineWidth(1)
			}
		}
		gc.SetLineDash([]float64{}, 0)
	}
	gc.SetLineWidth(1)
	gc.SetStrokeColor(destLineC)
	for _, pts := range destToDraw {
		DrawDestLine(gc, vp, pts[0], pts[1], showGrid)
	}
	gc.SetStrokeColor(color.Black)
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
	if zoom > 14 {
		gc.SetLineWidth(2)
	} else if zoom > 5 {
		gc.SetLineWidth(1)
	} else {
		gc.SetLineWidth(.25)
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
		availMap[src.Pid()] -= o.Size()
		DrawOrderLine(gc, vp, src.Loc(), tar.Loc(), showGrid)
	}
	// -------------------- PLANETS --------------------- //
	gc.SetLineWidth(.25)
	gc.SetStrokeColor(color.Black)
	for _, pv := range plToDraw {
		loc := pv.Loc()
		isTar1 := target1.IsCoord(loc)
		isTar2 := target2.IsCoord(loc)
		var avail int
		if pv.Controller() == fid {
			avail = availMap[pv.Pid()]
		}
		DrawPlanet(gc, vp, fid, avail, showGrid, isTar1, isTar2, pv)
	}

	// -------- //
	png.Encode(w, final)
}
