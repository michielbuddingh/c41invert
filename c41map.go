package main

import "os"
import "io"
import "flag"
import "encoding/json"
import "bytes"
import "histogram"
import "log"
import "image/jpeg"
import "image/png"
import "image"
import "image/color"

var options struct {
	Histogram string
	Input string
	Output string
	Contrast float64
}

func readHistogram() histogram.Channels {
	var buf bytes.Buffer
	f, err := os.Open(options.Histogram)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(&buf, f)

	var channels histogram.Channels

	jerr := json.Unmarshal(buf.Bytes(), &channels)

	if jerr != nil {
		log.Fatal(jerr)
	}

	return channels
}

func loadImage(filename string) (image.Image) {
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


func main() {
	flag.StringVar(&options.Histogram, "histogram", "", "")
	flag.StringVar(&options.Input, "input", "", "")
	flag.StringVar(&options.Output, "output", "", "")
	flag.Float64Var(&options.Contrast, "contrast", 1.0, "tweak contrast, 1.0 = normal")

	flag.Parse()

	channels := readHistogram()

	picture := loadImage(options.Input)

	colormap := channels.Sigmoid(options.Contrast)
	//colormap := channels.Linear()

	bounds := picture.Bounds()
	copy := image.NewRGBA64(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			in := color.RGBA64Model.Convert(picture.At(x, y))
			out := colormap(in.(color.RGBA64))
			copy.Set(x, y, out)
		}
	}

	f, ferr := os.Create(options.Output)
	if ferr != nil {
		log.Fatal(ferr)
	}
	defer f.Close()
	jpeg.Encode(f, copy, &jpeg.Options{Quality:95})
}
