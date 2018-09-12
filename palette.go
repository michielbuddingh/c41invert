package main

import (
	"math"
	"image"
	"image/color"
)

// Channel stores a histogram of all values present in
// and image for a specific colour channel
type Channel struct {
	Shades   [65536]uint32
	Min, Max, Total uint32
}

// Add adds a colour value to the channel
func (c *Channel) Add(val uint32) {
	c.Shades[val] += 1
	c.Total += 1
	if val < c.Min {
		c.Min = val
	}
	if val > c.Max {
		c.Max = val
	}
}

// Percentile finds the value at percentile pct in
// all the values stored for the channel.
func (c Channel) Percentile(pct float64) uint16 {
	total := float64(c.Total)

	var sum uint32
	for i := c.Min; i < c.Max; i++ {
		sum += c.Shades[i]
		if float64(sum)/total > pct {
			return uint16(i)
		}
	}

	return 65535
}

// Merge adds up the values in two channels.
func (c *Channel) Merge(other *Channel) {
	if other.Min < c.Min {
		c.Min = other.Min
	}

	if other.Max > c.Max {
		c.Max = other.Max
	}

	for i := c.Min; i < c.Max; i++ {
		c.Shades[i] += other.Shades[i]
	}

	c.Total += other.Total
}

// Palette groups a Red, a Green and a Blue channel
type Palette struct {
	Red, Green, Blue Channel
	Total            int
}

// Merge merges another Palette with this one.
func (p *Palette) Merge(other Palette) {
	p.Red.Merge(&other.Red)
	p.Green.Merge(&other.Green)
	p.Blue.Merge(&other.Blue)
	p.Total += other.Total
}

// Add adds a new colour value to a Palette
func (p *Palette) Add(c color.RGBA64) {
	r, g, b, _ := c.RGBA()
	p.Red.Add(r)
	p.Green.Add(g)
	p.Blue.Add(b)
	p.Total += 1
}

type Range struct {
	Low, High uint16
}

type Transformation struct {
	Red, Green, Blue Range
	Contrast float64
}

type Mapping func (color.RGBA64) (out color.RGBA64)

func (m Mapping) Apply(input image.Image) image.Image {
	bounds := input.Bounds()
	copy := image.NewRGBA64(bounds)
	
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			in := color.RGBA64Model.Convert(input.At(x, y))
			out := m(in.(color.RGBA64))
			copy.Set(x, y, out)
		}
	}

	return copy
}

func (t Transformation) Sigmoid() Mapping {
	rmin, rmax := float64(t.Red.Low), float64(t.Red.High)
	gmin, gmax := float64(t.Green.Low), float64(t.Green.High)
	bmin, bmax := float64(t.Blue.Low), float64(t.Blue.High)

	rdiff := float64(rmax - rmin)
	gdiff := float64(gmax - gmin)
	bdiff := float64(bmax - bmin)

	return func (in color.RGBA64) (out color.RGBA64) {
		valr := (float64(in.R) - rmin) - (rdiff / 2)
		valg := (float64(in.G) - gmin) - (gdiff / 2)
		valb := (float64(in.B) - bmin) - (bdiff / 2)

		valr *= (math.Pi / (rdiff / t.Contrast))
		valg *= (math.Pi / (gdiff / t.Contrast))
		valb *= (math.Pi / (bdiff / t.Contrast))
		
		out.R = uint16((math.Erf(valr) + 1) * 32767)
		out.G = uint16((math.Erf(valg) + 1) * 32767)
		out.B = uint16((math.Erf(valb) + 1) * 32767)

		out.A = in.A
		return out
	}
}

func (t Transformation) Linear() Mapping {
	rmin, rmax := float64(t.Red.Low), float64(t.Red.High)
	gmin, gmax := float64(t.Green.Low), float64(t.Green.High)
	bmin, bmax := float64(t.Blue.Low), float64(t.Blue.High)

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
