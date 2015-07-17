package gopow

import (
	"fmt"
	"image"
	"math"
	"time"

	"code.google.com/p/freetype-go/freetype"
	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"

	"github.com/dhogborg/rtl-gopow/internal/resources"
)

// font configuration
const (
	dpi      float64 = 72
	fontfile string  = "resources/fonts/luxisr.ttf"
	hinting  string  = "none"
	size     float64 = 18
)

type Annotator struct {
	image *image.RGBA
	table *TableComplex

	context *freetype.Context
}

func NewAnnotator(img *image.RGBA, table *TableComplex) (*Annotator, error) {

	a := &Annotator{
		image: img,
		table: table,
	}

	err := a.init()
	if err != nil {
		return nil, err
	}

	return a, nil

}

func (a *Annotator) init() error {

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

	a.context = freetype.NewContext()
	a.context.SetDPI(dpi)
	a.context.SetFont(font)
	a.context.SetFontSize(size)

	a.context.SetClip(a.image.Bounds())
	a.context.SetDst(a.image)
	a.context.SetSrc(fg)

	switch hinting {
	default:
		a.context.SetHinting(freetype.NoHinting)
	case "full":
		a.context.SetHinting(freetype.FullHinting)
	}

	return nil
}

func (a *Annotator) DrawXScale() error {

	log.WithFields(log.Fields{
		"hzLow":  a.table.HzLow,
		"hzHigh": a.table.HzHigh,
	}).Debug("annotate X scale")

	// how many samples?
	count := int(math.Floor(float64(a.table.Bins) / float64(350)))

	log.WithFields(log.Fields{
		"labels": count,
	}).Debug("annotate X scale")

	hzPerLabel := float64(a.table.HzHigh-a.table.HzLow) / float64(count)
	pxPerLabel := int(math.Floor(float64(a.table.Bins) / float64(count)))

	for si := 0; si < count; si++ {

		hz := a.table.HzLow + (float64(si) * hzPerLabel)
		px := si * pxPerLabel

		fract, suffix := humanize.ComputeSI(hz)
		str := fmt.Sprintf("%0.2f %sHz", fract, suffix)

		// draw a guideline on the exact frequency
		for i := 0; i < 50; i++ {
			a.image.Set(px, i, image.White)
		}

		// draw the text
		pt := freetype.Pt(px+10, 30)
		_, _ = a.context.DrawString(str, pt)

	}

	return nil
}

func (a *Annotator) DrawYScale() error {

	log.WithFields(log.Fields{
		"timestart": a.table.TimeStart.String(),
		"timeend":   a.table.TimeEnd.String(),
	}).Debug("annotate Y scale")

	start, end := a.table.TimeStart, a.table.TimeEnd

	// how many samples?
	count := int(math.Floor(float64(a.table.Integrations) / float64(100)))

	uStart := start.Unix()
	uEnd := end.Unix()

	secsPerLabel := int(math.Floor(float64(uEnd-uStart) / float64(count)))
	pxPerLabel := int(math.Floor(float64(a.table.Integrations) / float64(count)))

	log.WithFields(log.Fields{
		"labels":       count,
		"secsPerLabel": secsPerLabel,
		"pxPerLabel":   pxPerLabel,
	}).Debug("annotate Y scale")

	for si := 0; si < count; si++ {

		secs := time.Duration(secsPerLabel * si * int(time.Second))
		px := si * pxPerLabel

		var str string = ""

		if si == 0 {
			str = start.String()
		} else {
			point := start.Add(secs)
			str = point.Format("15:04:05")
		}

		// draw a guideline on the exact time
		for i := 0; i < 75; i++ {
			a.image.Set(i, px, image.White)
		}

		// draw the text, 3 px margin to the line
		pt := freetype.Pt(3, px-3)
		_, _ = a.context.DrawString(str, pt)

	}

	return nil

}
