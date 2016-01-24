package gopow

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
)

const (
	PowerConfigAuto = -9813
)

type RunConfig struct {
	InputFile   string
	OutputFile  string
	Format      string
	Annotations bool
	MaxPower    float64
	MinPower    float64
}

type GoPow struct {
	config    *RunConfig
	image     *image.RGBA
	timestamp time.Time
}

func NewGoPow(c *cli.Context) (*GoPow, error) {

	config := &RunConfig{
		InputFile:   c.String("input"),
		OutputFile:  c.String("output"),
		Format:      c.String("format"),
		Annotations: !c.Bool("no-annotations"),
		MaxPower:    c.Float64("max-power"),
		MinPower:    c.Float64("min-power"),
	}

	if !c.IsSet("max-power") {
		config.MaxPower = PowerConfigAuto
	}
	if !c.IsSet("min-power") {
		config.MinPower = PowerConfigAuto
	}

	if config.InputFile == "" {
		return nil, fmt.Errorf("missing input file")
	}

	if config.Format == "" {
		config.Format = "png"
	}

	if config.OutputFile == "" {
		config.OutputFile = config.InputFile + "." + config.Format
	}

	log.WithFields(log.Fields{
		"input": config.InputFile,
	}).Info("GoPow init")
	log.WithFields(log.Fields{
		"output": config.OutputFile,
	}).Info("GoPow init")
	log.WithFields(log.Fields{
		"format": config.Format,
	}).Info("GoPow init")

	g := &GoPow{
		config: config,
	}

	return g, nil
}

func (g *GoPow) Render() error {

	conf := &RenderConfig{}

	if g.config.MaxPower != PowerConfigAuto {
		conf.MaxPower = &g.config.MaxPower
	}

	if g.config.MinPower != PowerConfigAuto {
		conf.MinPower = &g.config.MinPower
	}

	log.Debug("staring render")
	g.timestamp = time.Now()

	table, err := NewTable(g.config.InputFile, conf)
	if err != nil {
		return err
	}

	g.image = table.Image()

	for y, row := range table.Rows {
		for x := range row.Samples {
			g.image.Set(x, y, table.ColorAt(x, y))
		}
	}

	if g.config.Annotations {

		annotator, err := NewAnnotator(g.image, table)
		if err != nil {
			return err
		}

		// add some frequency and time annotation
		annotator.DrawXScale()
		annotator.DrawYScale()
		annotator.DrawInfoBox()
	}

	return nil
}

func (g *GoPow) Write() error {

	log.WithFields(log.Fields{
		"file": g.config.OutputFile,
	}).Debug("staring output write")

	out, err := os.Create(g.config.OutputFile)
	if err != nil {
		return err
	}

	switch g.config.Format {
	case "png":
		err = png.Encode(out, g.image)
		break

	case "jpeg", "jpg":
		opt := &jpeg.Options{
			Quality: 98,
		}
		err = jpeg.Encode(out, g.image, opt)
		break

	default:
		return fmt.Errorf("unsupported format: %s", g.config.Format)
	}

	if err != nil {
		return err
	}

	duration := humanize.RelTime(g.timestamp, time.Now(), "", "")
	log.Info("GoPow finished in " + duration)

	return nil
}
