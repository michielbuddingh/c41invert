package histogram

import "image/color"
import "math"

type Histogram [65536]int32

func (h Histogram) Min() uint16 {
	for i := 0; i < len(h); i++ {
		if h[i] > 0 {
			return uint16(i)
		}
	}
	return 0
}

func (h Histogram) Max() uint16 {
	for i := len(h) - 1; i >= 0; i-- {
		if h[i] > 0 {
			return uint16(i)
		}
	}
	return 65535
}


func (h *Histogram) Merge(other *Histogram) {
	for i := 0; i < len(other); i++ {
		h[i] += other[i]
	}
}

type Channels struct {
	Red, Green, Blue Histogram
}

type Mapping func (color.RGBA64) (out color.RGBA64)

func (c Channels) Sigmoid(contrast float64) Mapping  {
	rmin, rmax := float64(c.Red.Min()), float64(c.Red.Max())
	gmin, gmax := float64(c.Green.Min()), float64(c.Green.Max())
	bmin, bmax := float64(c.Blue.Min()), float64(c.Blue.Max())

	rdiff := float64(rmax - rmin)
	gdiff := float64(gmax - gmin)
	bdiff := float64(bmax - bmin)

	return func (in color.RGBA64) (out color.RGBA64) {
		valr := (float64(in.R) - rmin) - (rdiff / 2)
		valg := (float64(in.G) - gmin) - (gdiff / 2)
		valb := (float64(in.B) - bmin) - (bdiff / 2)

		valr *= (math.Pi / (-rdiff / contrast))
		valg *= (math.Pi / (-gdiff / contrast))
		valb *= (math.Pi / (-bdiff / contrast))
		

		out.R = uint16((math.Erf(valr) + 1) * 32768)
		out.G = uint16((math.Erf(valg) + 1) * 32768)
		out.B = uint16((math.Erf(valb) + 1) * 32768)

		out.A = in.A
		return out
	}
}







