package mapping

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"mule/overpower"
	"net/http"
)

func ServeMap(w http.ResponseWriter, mv overpower.MapView, fid int, facList []overpower.Faction, pvList []overpower.PlanetView, svList []overpower.ShipView, orders []overpower.Order) {
	final := image.NewRGBA(image.Rect(0, 0, 600, 600))
	black := color.RGBA{0x66, 0x00, 0x00, 255}
	//black := color.RGBA{0x00, 0x00, 0x00, 255}
	draw.Draw(final, final.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)
	png.Encode(w, final)
}
