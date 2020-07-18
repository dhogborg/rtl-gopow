package gopow

import (
	"bytes"
	"image"
	"io/ioutil"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

type TableComplex struct {
	File   string // our input file
	Config *RenderConfig

	Rows []*LineComplex

	Bins         int // horizontal slots, columns, bandwidth
	Integrations int // vertical slots, rows

	HzLow  float64 // X Scale start
	HzHigh float64 // X Scale end

	TimeStart *time.Time // real time, Y Scale
	TimeEnd   *time.Time
}

// RenderConfig overrides automaticly calculated defaults
type RenderConfig struct {
	MinPower *float64 // minimum power value, used for color rendering
	MaxPower *float64 // maximum dito
}

func NewTable(file string, conf *RenderConfig) (*TableComplex, error) {
	log.Debug("creating table")

	t := &TableComplex{
		Config: conf,
	}

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
	var max = float64(math.MaxFloat64 * -1)
	var min = float64(math.MaxFloat64)

	lines := bytes.Split(filebuffer, []byte("\n"))

	table := map[string][]*LineComplex{}

	for _, l := range lines {
		cells := strings.Split(string(l), ",")
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

			if min > row.LowSample() {
				min = row.LowSample()
			}
			if max < row.HighSample() {
				max = row.HighSample()
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

	if t.Config.MaxPower == nil {
		t.Config.MaxPower = &max
	}

	if t.Config.MinPower == nil {
		t.Config.MinPower = &min
	}

	log.WithFields(log.Fields{
		"pMax": *t.Config.MaxPower,
		"pMin": *t.Config.MinPower,
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
