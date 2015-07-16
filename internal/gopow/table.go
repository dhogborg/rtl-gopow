package gopow

import (
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"sort"
	"strings"
	"time"

	"code.google.com/p/freetype-go/freetype"
	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/dhogborg/rtl-gopow/internal/resources"
)

// font configuration
var (
	dpi      float64 = 72
	fontfile string  = "resources/fonts/luxisr.ttf"
	hinting  string  = "none"
	size     float64 = 34
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

	t.Max = float64(math.MaxFloat64 * -1)
	t.Min = float64(math.MaxFloat64)

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

			t.HzLow = row.HzLow
			t.HzHigh = row.HzHigh

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

	return colorful.Hsv(hue, 1, 0.90)

}

func (t *TableComplex) AnnotateXScale(img *image.RGBA) error {

	log.WithFields(log.Fields{
		"hzLow":  t.HzLow,
		"hzHigh": t.HzHigh,
	}).Debug("annotate X scale")

	// load the font
	fontBytes, err := resources.Asset(fontfile)
	if err != nil {
		return err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// Initialize the context.
	fg := image.White
	ruler := image.White

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(size)

	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)

	switch hinting {
	default:
		c.SetHinting(freetype.NoHinting)
	case "full":
		c.SetHinting(freetype.FullHinting)
	}

	// how many samples?
	count := int(math.Floor(float64(t.Bins) / float64(500)))

	log.WithFields(log.Fields{
		"labels": count,
	}).Debug("annotate X scale")

	hzPerLable := float64(t.HzHigh-t.HzLow) / float64(count)
	pxPerLable := int(math.Floor(float64(t.Bins) / float64(count)))

	for si := 0; si < count; si++ {

		hz := t.HzLow + (float64(si) * hzPerLable)
		px := si * pxPerLable

		str := humanize.SI(hz, "Hz")

		// draw a guideline on the exact frequency
		for i := 0; i < 50; i++ {
			img.Set(px, i, ruler)
		}

		// draw the text
		pt := freetype.Pt(px+10, 30)
		_, _ = c.DrawString(str, pt)

	}

	return nil
}

func (t *TableComplex) AnnotateYScale(img *image.RGBA) error {

	log.WithFields(log.Fields{
		"timestart": t.TimeStart.String(),
		"timeend":   t.TimeEnd.String(),
	}).Debug("annotate Y scale")

	// how many samples?
	const count = 10

	log.WithFields(log.Fields{
		"labels": count,
	}).Debug("annotate Y scale")

	return nil
}
