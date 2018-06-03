package main

import "os"
import "image"
import "image/png"
import "log"
import "encoding/json"
import "histogram"
import "image/color"

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

func sampleCentre(picture image.Image) *image.Rectangle {
	bounds := picture.Bounds()

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	sampleArea := image.Rectangle{
		image.Point{bounds.Min.X + int(float64(width) * 0.1),
			bounds.Min.Y + int(float64(height) * 0.1)},
		image.Point{bounds.Max.X - int(float64(width) * 0.1),
			bounds.Max.Y - int(float64(height) * 0.1)}}

	return &sampleArea
}

func histograms(picture image.Image, sampleArea *image.Rectangle) (*histogram.Channels) {

	channels := new(histogram.Channels)

	for x := sampleArea.Min.X; x < sampleArea.Max.X; x++ {
		for y := sampleArea.Min.Y; y < sampleArea.Max.Y; y++ {
			r, g, b, _ := color.RGBAModel.Convert(picture.At(x, y)).RGBA()
			channels.Red.Values[r] += 1
			channels.Green.Values[g] += 1
			channels.Blue.Values[b] += 1
			channels.Red.Cumulative += 1
			channels.Green.Cumulative += 1
			channels.Blue.Cumulative += 1
		}
	}
	return channels
}

func main() {
	var summary histogram.Channels
	for i := 1; i < len(os.Args); i++ {
		picture := load(os.Args[i])
		sampleArea := sampleCentre(picture)
		sampled := histograms(picture, sampleArea)
		summary.Red.Merge(&sampled.Red)
		summary.Green.Merge(&sampled.Green)
		summary.Blue.Merge(&sampled.Blue)
	}

	b, merr := json.Marshal(summary)
	if merr == nil {
		os.Stdout.Write(b)
	} else {
		log.Fatal(merr)
	}
}
