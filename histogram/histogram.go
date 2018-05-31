package histogram

import "image/color"
import "math"

type Histogram struct {
	Values [65536]uint32
	Cumulative int32
}

func (h Histogram) Min() uint16 {
	for i := 0; i < len(h.Values); i++ {
		if h.Values[i] > 0 {
			return uint16(i)
		}
	}
	return 0
}

func (h Histogram) Pctile(pct float64) uint16 {
	total := float64(h.Cumulative)

	var sum uint32
	for i := 0; i < len(h.Values); i++ {
		sum += h.Values[i]
		if float64(sum) / total > pct {
			return uint16(i)
		}
	}

	return 65535
}

func (h Histogram) Max() uint16 {
	for i := len(h.Values) - 1; i >= 0; i-- {
		if h.Values[i] > 0 {
			return uint16(i)
		}
	}
	return 65535
}

func (h *Histogram) Merge(other *Histogram) {
	for i := 0; i < len(other.Values); i++ {
		h.Values[i] += other.Values[i]
	}
	h.Cumulative += other.Cumulative
}

type Channels struct {
	Red, Green, Blue Histogram
}

type Mapping func (color.RGBA64) (out color.RGBA64)

func (c Channels) Linear() Mapping {
	rmin, rmax := float64(c.Red.Pctile(0.1)), float64(c.Red.Pctile(0.9))
	gmin, gmax := float64(c.Green.Pctile(0.1)), float64(c.Green.Pctile(0.9))
	bmin, bmax := float64(c.Blue.Pctile(0.1)), float64(c.Blue.Pctile(0.9))

	rdiff := float64(rmax - rmin)
	gdiff := float64(gmax - gmin)
	bdiff := float64(bmax - bmin)

	unit := func(v float64) float64 {
		if v < 0 {
			return 1
		}
		if v > 1 {
			return 0
		}
		return 1 - v
	}
	
	return func (in color.RGBA64) (out color.RGBA64) {
		// normalize to 0..1
		valr := unit((float64(in.R) - rmin) / rdiff)
		valg := unit((float64(in.G) - gmin) / gdiff)
		valb := unit((float64(in.B) - bmin) / bdiff)

		out.R = uint16(valr * 65535.0)
		out.G = uint16(valg * 65535.0)
		out.B = uint16(valb * 65535.0)
		out.A = in.A

		return out
	}
}

func (c Channels) Sigmoid(contrast float64) Mapping  {
	rmin, rmax := float64(c.Red.Pctile(0.1)), float64(c.Red.Pctile(0.9))
	gmin, gmax := float64(c.Green.Pctile(0.1)), float64(c.Green.Pctile(0.9))
	bmin, bmax := float64(c.Blue.Pctile(0.1)), float64(c.Blue.Pctile(0.9))

	
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







