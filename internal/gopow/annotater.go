package gopow

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/freetype"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font"

	"github.com/dhogborg/rtl-gopow/internal/resources"
)

// font configuration
const (
	dpi      float64 = 72
	fontfile string  = "resources/fonts/luxisr.ttf"
	hinting  string  = "none"
	size     float64 = 15
	spacing  float64 = 1.1
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

	luxisr, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// Initialize the context.
	fg := image.White

	a.context = freetype.NewContext()
	a.context.SetDPI(dpi)
	a.context.SetFont(luxisr)
	a.context.SetFontSize(size)

	a.context.SetClip(a.image.Bounds())
	a.context.SetDst(a.image)
	a.context.SetSrc(fg)

	switch hinting {
	default:
		a.context.SetHinting(font.HintingNone)
	case "full":
		a.context.SetHinting(font.HintingFull)
	}

	return nil
}

func (a *Annotator) DrawXScale() error {

	log.WithFields(log.Fields{
		"hzHigh": humanize.SI(a.table.HzHigh, "Hz"),
		"hzLow":  humanize.SI(a.table.HzLow, "Hz"),
	}).Debug("annotate X scale")

	// how many samples?
	count := int(math.Floor(float64(a.table.Bins) / float64(350)))

	hzPerLabel := float64(a.table.HzHigh-a.table.HzLow) / float64(count)
	pxPerLabel := int(math.Floor(float64(a.table.Bins) / float64(count)))

	log.WithFields(log.Fields{
		"labels":     count,
		"hzPerLabel": humanize.SI(hzPerLabel, "Hz"),
		"pxPerLabel": pxPerLabel,
	}).Debug("annotate X scale")

	for si := 0; si < count; si++ {

		hz := a.table.HzLow + (float64(si) * hzPerLabel)
		px := si * pxPerLabel

		fract, suffix := humanize.ComputeSI(hz)
		str := fmt.Sprintf("%0.2f %sHz", fract, suffix)

		// draw a guideline on the exact frequency
		for i := 0; i < 30; i++ {
			a.image.Set(px, i, image.White)
		}

		// draw the text
		pt := freetype.Pt(px+5, 17)
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

func (a *Annotator) DrawInfoBox() error {

	tStart, tEnd := a.table.TimeStart, a.table.TimeEnd
	// tDuration := humanize.RelTime(*tStart, *tEnd, "", "")
	tPixel := (tEnd.Unix() - tStart.Unix()) / int64(a.table.Integrations)

	fStart, fEnd := a.table.HzLow, a.table.HzHigh
	fBandwidth := fEnd - fStart
	fPixel := fBandwidth / float64(a.table.Bins)

	perPixel := fmt.Sprintf("%s x %d seconds", a.humanHz(fPixel), tPixel)

	// positioning
	imgSize := a.table.Image().Bounds().Size()
	top, left := imgSize.Y-75, 3

	strings := []string{
		"Scan start: " + tStart.String(),
		"Scan end: " + tEnd.String(),
		// "Scan duration: " + tDuration,
		fmt.Sprintf("Band: %s to %s", a.humanHz(fStart), a.humanHz(fEnd)),
		fmt.Sprintf("Bandwidth: %s", a.humanHz(fBandwidth)),
		"1 pixel = " + perPixel,
	}

	// drawing
	pt := freetype.Pt(left, top)
	for _, s := range strings {
		_, _ = a.context.DrawString(s, pt)
		pt.Y += a.context.PointToFixed(size * spacing)
	}

	return nil
}

func (a *Annotator) humanHz(hz float64) string {
	fpxSI, fpxSuffix := humanize.ComputeSI(hz)
	return fmt.Sprintf("%0.2f %sHz", fpxSI, fpxSuffix)
}
