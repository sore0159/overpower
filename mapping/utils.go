package mapping

import (
	"mule/hexagon"
	"mule/overpower"
)

func GetVP(mv overpower.MapView) *hexagon.Viewport {
	vp := hexagon.MakeViewport(float64(mv.Zoom()), false, true)
	center := mv.Center()
	vp.SetAnchor(center[0], center[1], float64(MAPW)/2, float64(MAPH)/2)
	vp.SetFrame(0, 0, MAPW, MAPH)
	return vp
}
