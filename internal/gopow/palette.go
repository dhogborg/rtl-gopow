package gopow

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

type Palette interface {
	ColorAt(table *TableComplex, x, y int) color.Color
}

type YellowPalette struct {
}

func (p *YellowPalette) ColorAt(table *TableComplex, x, y int) color.Color {
	cell := table.Rows[y].Sample(x)

	hueStart := 0.0
	hueEnd := 1.0

	span := (*table.Config.MinPower - *table.Config.MaxPower) * -1
	hPerDeg := (hueStart - hueEnd) / span
	powNormalized := cell - *table.Config.MinPower
	powDegrees := powNormalized * hPerDeg
	hue := hueStart - powDegrees

	if hue < hueStart {
		hue = hueStart
	}
	if hue > hueEnd {
		hue = hueEnd
	}

	return colorful.Color{hue, hue, 0}
}

type SpectrumPalette struct {
}

func (p *SpectrumPalette) ColorAt(table *TableComplex, x, y int) color.Color {
	cell := table.Rows[y].Sample(x)

	hueStart := 236.0
	hueEnd := 0.0

	span := (*table.Config.MinPower - *table.Config.MaxPower) * -1
	hPerDeg := (hueStart - hueEnd) / span
	powNormalized := cell - *table.Config.MinPower
	powDegrees := powNormalized * hPerDeg
	hue := hueStart - powDegrees

	if hue > hueStart {
		hue = hueStart
	}
	if hue < hueEnd {
		hue = hueEnd
	}

	return colorful.Hsv(hue, 1, 0.90)
}
