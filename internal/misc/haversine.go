package misc

import (
	"math"
)

type Havesine struct {
}

func (h *Havesine) CalcDistanceFromDegree(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	var EarthRadius = 6371.00
	var bx, lx float64
	var by, ly float64
	bx = latitude1
	lx = longitude1
	by = latitude2
	ly = longitude2
	//Bogenmaß
	bx = bx * math.Pi / 180
	by = by * math.Pi / 180
	lx = lx * math.Pi / 180
	ly = ly * math.Pi / 180
	a := math.Pow(math.Sin(((by-bx)/2)), 2) + math.Cos(bx)*math.Cos(by)*
		math.Pow(math.Sin(((ly-lx)/2)), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := EarthRadius * c
	return d
}

func NewHaversine() *Havesine {
	res := Havesine{}
	return &res
}
