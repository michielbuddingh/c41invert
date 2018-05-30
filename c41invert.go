package main

import "os"
import "image"
import "image/png"
import "image/jpeg"
import "image/color"
import "log"
import "math"

func load(filename string) (image.Image) {
	f, oerr := os.Open(filename)
	if oerr != nil {
		log.Fatal(oerr)
	}
	defer f.Close()

	p, derr := png.Decode(f)
	if derr != nil {
		log.Fatal(derr)
	}

	return p
}

type histogram [65536]int32

func (h histogram) min() uint16 {
	for i := 0; i < len(h); i++ {
		if h[i] > 0 {
			log.Printf("min %d", i)
			return uint16(i)
		}
	}
	return 0
}

func (h histogram) max() uint16 {
	for i := len(h) - 1; i >= 0; i-- {
		if h[i] > 0 {
			log.Printf("max %d", i)
			return uint16(i)
		}
	}
	return 65535
}

func mapping(red, green, blue *histogram) (func (color.RGBA64) (out color.RGBA64))  {
	rmin, rmax := float64(red.min()), float64(red.max())
	gmin, gmax := float64(green.min()), float64(green.max())
	bmin, bmax := float64(blue.min()), float64(blue.max())

	rdiff := float64(rmax - rmin)
	gdiff := float64(gmax - gmin)
	bdiff := float64(bmax - bmin)

	return func (in color.RGBA64) (out color.RGBA64) {
		valr := (float64(in.R) - rmin) - (rdiff / 2)
		valg := (float64(in.G) - gmin) - (gdiff / 2)
		valb := (float64(in.B) - bmin) - (bdiff / 2)

		valr *= (math.Pi / (-rdiff / 1))
		valg *= (math.Pi / (-gdiff / 1))
		valb *= (math.Pi / (-bdiff / 1))
		

		out.R = uint16((math.Erf(valr) + 1) * 32768)
		out.G = uint16((math.Erf(valg) + 1) * 32768)
		out.B = uint16((math.Erf(valb) + 1) * 32768)

//		log.Printf("%d - %f, %d, %f, %f\n", in.R, rmax, out.R, valr, rdiff)
		out.A = in.A
		return out
	}
}


func histograms(picture image.Image, sampleArea image.Rectangle) (*histogram, *histogram, *histogram) {

	red := new(histogram)
	green := new(histogram)
	blue := new(histogram)

	for x := sampleArea.Min.X; x < sampleArea.Max.X; x++ {
		for y := sampleArea.Min.Y; y < sampleArea.Max.Y; y++ {
			r, g, b, _ := picture.At(x, y).RGBA()
			red[r] += 1
			green[g] += 1
			blue[b] += 1
		}
	}
	return red, green, blue
}


func main() {
	if len(os.Args) == 1 {
		log.Fatal("Need a path to an image file")
	}

	picture := load(os.Args[1])
	bounds := picture.Bounds()

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	log.Printf("W x H : %d x %d",
		width,
		height)

	sampleArea := image.Rectangle{
		image.Point{bounds.Min.X + int(float64(width) * 0.1),
			bounds.Min.Y + int(float64(height) * 0.1)},
		image.Point{bounds.Max.X - int(float64(width) * 0.1),
			bounds.Max.Y - int(float64(height) * 0.1)}}

	colormap := mapping(histograms(picture, sampleArea))

	copy := image.NewRGBA64(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			in := picture.At(x, y)
			out := colormap(in.(color.RGBA64))
			copy.Set(x, y, out)
		}
	}

	jpeg.Encode(os.Stdout, copy, &jpeg.Options{Quality:95})

}
