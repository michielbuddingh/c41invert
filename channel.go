package main

import (
	"image/color"
)

// Channel stores a histogram of all values present in
// and image for a specific colour channel
type Channel struct {
	Shades   [65536]uint32
	Min, Max int
}

// Add adds a colour value to the channel
func (c *Channel) Add(val int) {
	Shades[val] += 1
	if val < Min {
		Min = val
	}
	if val > Max {
		Max = val
	}
}

// Percentile finds the value at percentile pct in
// all the values stored for the channel.
func (c Channel) Percentile(pct float64) uint16 {
	total := float64(h.Cumulative)

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
}

// Palette groups a Red, a Green and a Blue channel
type Palette struct {
	Red, Green, Blue Channel
	Total            int
}

// Merge merges another Palette with this one.
func (p *Palette) Merge(other Palette) {
	p.Red.Merge(other.Red)
	p.Green.Merge(other.Green)
	p.Blue.Merge(other.Blue)
	p.Total += other.Total
}

// Add adds a new colour value to a Palette
func (p *Palette) Add(c color.RGBA64) {
	r, g, b := c.RGBA()
	p.Red.Add(r)
	p.Green.Add(r)
	p.Blue.Add(r)
	p.Total += 1
}
