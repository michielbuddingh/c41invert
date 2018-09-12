package main

import "os"
import "log"
import "context"
import "flag"
import "image"
import _ "image/png"
import "image/jpeg"
import _ "golang.org/x/image/tiff"
import _ "golang.org/x/image/webp"
import "image/color"
import "github.com/google/subcommands"

func load(filename string) (image.Image, error) {
	f, oerr := os.Open(filename)
	if oerr != nil {
		return nil, oerr
	}
	defer f.Close()

	p, _, derr := image.Decode(f)
	if derr != nil {
		return nil, derr
	}

	return p, nil
}

func samplePalette(picture image.Image, sampleArea image.Rectangle) *Palette {
	var palette Palette
	for x := sampleArea.Min.X; x < sampleArea.Max.X; x++ {
		for y := sampleArea.Min.Y; y < sampleArea.Max.Y; y++ {
			palette.Add(color.RGBA64Model.Convert(picture.At(x, y)).(color.RGBA64))
		}
	}
	return &palette
}

type convertCmd struct {
	infile                string
	outfile               string
	sampleFraction        float64
	lowlights, highlights float64
	linear                bool
}

func (*convertCmd) Name() string {
	return "convert"
}

func (*convertCmd) Synopsis() string {
	return "Invert input image, normalize colors and output a file"
}

func (*convertCmd) Usage() string {
	return ""
}

func (c *convertCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.infile, "i", "", "Input file")
	f.StringVar(&c.outfile, "o", "", "Output file")
	f.Float64Var(&c.sampleFraction,
		"sample_fraction", 0.8,
		"Sample palette from a fraction crop of the center, 0 < fraction < 1 (default 0.8)")
	f.Float64Var(&c.lowlights,
		"lowlights", 0.01,
		"Shadows start here, lower values save more shadows")
	f.Float64Var(&c.highlights,
		"highlights", 0.99,
		"Highlights start here, lower values saves more highlights")
	f.BoolVar(&c.linear,
		"linear", false,
		"Use linear mapping instead of sigmoid function")
}

func (c *convertCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.infile == "" {
		log.Fatal("Must specify input file")
	}
	if c.outfile == "" {
		log.Fatal("Must specify output file")
	}

	picture, load_err := load(c.infile)

	if load_err != nil {
		log.Fatalf("Could not load input file `%s`: %v",
			c.infile,
			load_err)
	}

	sampleArea := sampleBounds(c.sampleFraction, picture)
	palette := samplePalette(picture, sampleArea)

	t := Transformation{
		Range{Low: palette.Red.Percentile(c.lowlights), High: palette.Red.Percentile(c.highlights)},
		Range{Low: palette.Green.Percentile(c.lowlights), High: palette.Green.Percentile(c.highlights)},
		Range{Low: palette.Blue.Percentile(c.lowlights), High: palette.Blue.Percentile(c.highlights)},
		c.lowlights - c.highlights,
	}

	mapping := t.Sigmoid()
	if c.linear {
		mapping = t.Linear()
	}

	copy := mapping.Apply(picture)

	of, ferr := os.Create(c.outfile)
	if ferr != nil {
		log.Fatal(ferr)
	}
	defer of.Close()
	jpeg.Encode(of, copy, &jpeg.Options{Quality: 95})

	return subcommands.ExitSuccess
}

// sampleBounds gets a bounding box for a center fraction of the image, based
// on parameter fraction
func sampleBounds(fraction float64, picture image.Image) image.Rectangle {
	bounds := picture.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	border := (1 - fraction) / 2

	sampleArea := image.Rectangle{
		image.Point{bounds.Min.X + int(float64(width)*border),
			bounds.Min.Y + int(float64(height)*border)},
		image.Point{bounds.Max.X - int(float64(width)*border),
			bounds.Max.Y - int(float64(height)*border)}}

	return sampleArea
}

func main() {
	subcommands.Register(&convertCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
