package gopow

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"code.google.com/p/draw2d/draw2d"
	"code.google.com/p/freetype-go/freetype"
	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/dhogborg/rtl-gopow/internal/resources"
)

// font configuration
var (
	dpi      float64 = 72
	fontfile string  = "font/luxisr.ttf"
	hinting  string  = "none"
	size     float64 = 24
	spacing  float64 = 1.5
	wonb     bool    = true
)

type TableComplex struct {
	File string // our input file

	Rows []*LineComplex

	Min float64 // minimum power value, used for color rendering
	Max float64 // maximum dito

	Bins         int // horizontal slots, columns, bandwidth
	Integrations int // vertical slots, rows

	HzLow  float64 // X Scale start
	HzHigh float64 // X Scale end

	TimeStart *time.Time // real time, Y Scale
	TimeEnd   *time.Time
}

func NewTable(file string) (*TableComplex, error) {

	log.Debug("creating table")

	t := &TableComplex{}

	err := t.Load(file)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TableComplex) Load(file string) error {

	log.Debug("loading table")

	t.File = file

	buff, err := ioutil.ReadFile(t.File)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"bytes": len(buff),
		"size":  humanize.Bytes(uint64(len(buff))),
	}).Debug("file loaded")

	t.Rows = t.parseBuffer(buff)

	return nil
}

func (t *TableComplex) parseBuffer(filebuffer []byte) []*LineComplex {

	t.Max = float64(-99999999)
	t.Min = float64(99999999)

	block := string(filebuffer)
	lines := strings.Split(block, "\n")

	table := map[string][]*LineComplex{}

	for _, l := range lines {

		cells := strings.Split(l, ",")
		line := NewLineComplex(cells)

		if table[line.Hash] == nil {
			table[line.Hash] = []*LineComplex{}
		}

		table[line.Hash] = append(table[line.Hash], line)
	}

	rows := []*LineComplex{}

	// loop over hash keys with lines
	for _, lines := range table {

		row := t.IntegrateLines(lines)

		if row != nil {

			rows = append(rows, row)

			if t.Min > row.LowSample() {
				t.Min = row.LowSample()
			}
			if t.Max < row.HighSample() {
				t.Max = row.HighSample()
			}

			if t.HzLow > row.HzLow {
				t.HzLow = row.HzLow
			}
			if t.HzHigh < row.HzHigh {
				t.HzHigh = row.HzHigh
			}

			if row.Time != nil {

				if t.TimeStart == nil {
					t.TimeStart = row.Time
				}

				if t.TimeEnd == nil {
					t.TimeEnd = row.Time
				}

				if t.TimeStart.Unix() > row.Time.Unix() {
					t.TimeStart = row.Time
				}

				if t.TimeEnd.Unix() < row.Time.Unix() {
					t.TimeEnd = row.Time
				}
			}

		}
	}

	sort.Sort(LineSort(rows))

	log.WithFields(log.Fields{
		"pMax": t.Max,
		"pMin": t.Min,
	}).Debug("integrated lines")

	t.Integrations = len(rows)

	if t.Integrations > 0 {
		t.Bins = len(rows[0].Samples)
	} else {
		log.Fatal("no samples found")
	}

	log.WithFields(log.Fields{
		"bins":         t.Bins,
		"integrations": t.Integrations,
	}).Debug("parsed table")

	return rows
}

func (t *TableComplex) Image() *image.RGBA {

	log.WithFields(log.Fields{
		"width":  t.Bins,
		"height": t.Integrations,
	}).Debug("create image")

	return image.NewRGBA(image.Rect(0, 0, int(t.Bins), int(t.Integrations)))
}

func (t *TableComplex) IntegrateLines(lines []*LineComplex) *LineComplex {

	if len(lines) == 0 {
		return nil
	}

	masterline := lines[0]
	for i, l := range lines {
		if i > 0 {
			masterline.AddSamples(l)
		}

	}

	return masterline
}

func (t *TableComplex) ColorAt(x, y int) color.Color {

	cell := t.Rows[y].Sample(x)

	hue_start := 236.0
	hue_end := 0.0

	span := (t.Min - t.Max) * -1
	h_per_deg := (hue_start - hue_end) / span
	pow_normalized := cell - t.Min
	pow_degrees := pow_normalized * h_per_deg
	hue := hue_start - pow_degrees

	return colorful.Hsv(hue, 1, 0.8)

}

func (t *TableComplex) AnnotateXScale(img *image.RGBA) error {

	// how many samples?
	// const samples = 10

	fontBytes, err := resources.Asset("resources/fonts/luxisr.ttf")
	if err != nil {
		return err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	gc := draw2d.NewGraphicContext(img)

	draw2d.RoundRect(gc, 5, 5, 95, 95, 10, 10)
	gc.FillStroke()
	gc.SetFontSize(18)
	gc.MoveTo(10, 52)

	gc.SetFont(font)
	// gc.SetFontData(draw2d.FontData{"luxi", draw2d.FontFamilyMono, draw2d.FontStyleBold | draw2d.FontStyleItalic})
	gc.SetFont(font)

	width := gc.FillString("cou")
	fmt.Printf("width: %f\n", width)
	gc.RMoveTo(width+1, 0)
	gc.FillString("cou")

	return nil
}

func (t *TableComplex) AnnotateYScale(img *image.RGBA) error {

	return nil
}
