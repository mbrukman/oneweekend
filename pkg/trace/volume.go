package trace

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

type Volume struct {
	boundary HitBounder
	density  float64
	phase    Material
}

func NewVolume(boundary HitBounder, density float64, phase Material) *Volume {
	return &Volume{
		boundary: boundary,
		density:  density,
		phase:    phase,
	}
}

func (v *Volume) Hit(r Ray, dMin, dMax float64) *Hit {
	hit1 := v.boundary.Hit(r, -math.MaxFloat64, math.MaxFloat64)
	if hit1 == nil {
		return nil
	}
	hit2 := v.boundary.Hit(r, hit1.Dist+bias, math.MaxFloat64)
	if hit2 == nil {
		return nil
	}
	if hit1.Dist < dMin {
		hit1.Dist = dMin
	}
	if hit2.Dist > dMax {
		hit2.Dist = dMax
	}
	if hit1.Dist > hit2.Dist {
		return nil
	}
	if hit1.Dist < 0 {
		hit1.Dist = 0
	}
	dInside := hit2.Dist - hit1.Dist
	dHit := -(1 / v.density) * math.Log(rand.Float64())
	if dHit >= dInside {
		return nil
	}
	d := hit1.Dist + dHit
	return &Hit{
		Dist: d,
		Norm: geom.NewUnit(1, 0, 0),
		UV:   geom.NewVec(0, 0, 0),
		Pt:   r.At(d),
		Mat:  v.phase,
	}
}

func (v *Volume) Bounds(t0, t1 float64) *AABB {
	return v.boundary.Bounds(t0, t1)
}
