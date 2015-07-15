package gopow

import (
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type LineComplex struct {
	Time *time.Time
	Hash string // a unique hash for this line in time

	HzLow       float64
	HzHigh      float64
	HzStep      float64
	SampleCount int

	Samples []float64
}

type LineSort []*LineComplex

func (a LineSort) Len() int {
	return len(a)
}

func (a LineSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a LineSort) Less(i, j int) bool {
	if a[i].Time == nil {
		return false
	}

	if a[j].Time == nil {
		return true
	}

	return a[i].Time.Unix() < a[j].Time.Unix()
}

func NewLineComplex(cells []string) *LineComplex {

	// bail early if there is something wrong with the line
	if len(cells) < 7 {
		return &LineComplex{}
	}

	date := cells[0]
	clock := cells[1]

	const format = "2006-01-02 15:04:05"
	datetime, err := time.Parse(format, date+clock)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"string": date + clock,
		}).Fatal("date parsing failure")
	}

	hzLow, _ := strconv.ParseFloat(cells[2], 64)
	hzHigh, _ := strconv.ParseFloat(cells[3], 64)
	hzStep, _ := strconv.ParseFloat(cells[3], 64)
	sc, _ := strconv.ParseInt(cells[4], 10, 64)

	samples := []float64{}
	for _, s := range cells[6:] {
		sf64, err := strconv.ParseFloat(strings.Trim(s, " "), 64)
		if err != nil {
			samples = append(samples, 0)
		} else {
			samples = append(samples, sf64)
		}
	}

	return &LineComplex{
		Time:        &datetime,
		Hash:        cells[0] + cells[1],
		HzLow:       hzLow,
		HzHigh:      hzHigh,
		HzStep:      hzStep,
		SampleCount: int(sc),

		Samples: samples, // the rest of the cells end up as samples
	}
}

func (l *LineComplex) AddSamples(line *LineComplex) {

	if line.HzHigh > l.HzHigh {
		l.HzHigh = line.HzHigh
	}

	if line.HzLow < l.HzLow {
		l.HzLow = line.HzLow
	}

	if l.Samples == nil {
		l.Samples = []float64{}
	}

	l.Samples = append(l.Samples, line.Samples...)

}

func (l *LineComplex) HighSample() float64 {

	high := float64(-99999)
	for _, sample := range l.Samples {
		if sample > high {
			high = sample
		}
	}

	return high
}

func (l *LineComplex) LowSample() float64 {
	low := float64(99999)
	for _, sample := range l.Samples {
		if sample < low {
			low = sample
		}
	}

	return low
}

func (l *LineComplex) Sample(x int) float64 {
	return l.Samples[x]
}
