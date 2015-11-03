package metric

import (
	"fmt"
)

type MeasurementType string

const (
	GAUGE      = MeasurementType("gauge")
	CUMULATIVE = MeasurementType("cumulative")
)


type MetricPuller interface {
    Pull() <-[]Measurement  // Returns a channel of measurements
}
