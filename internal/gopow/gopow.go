package gopow

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
)

type RunConfig struct {
	InputFile   string
	OutputFile  string
	Format      string
	Annotations bool
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

	log.Debug("staring render")
	g.timestamp = time.Now()

	table, err := NewTable(g.config.InputFile)
	if err != nil {
		return err
	}

	g.image = table.Image()

	for y, row := range table.Rows {
		for x, _ := range row.Samples {
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

	err = png.Encode(out, g.image)
	if err != nil {
		return err
	}

	duration := humanize.RelTime(g.timestamp, time.Now(), "", "")
	log.Info("GoPow finished in " + duration)

	return nil
}
