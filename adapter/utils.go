package adapter

import (
	"fmt"
	"math"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/prometheus/prometheus/prompb"
)

func seriesToRecords(prefix, suffix string, series *prompb.TimeSeries) (*timestreamwrite.Record, []*timestreamwrite.Record, error) {
	dims := make([]*timestreamwrite.Dimension, 0)
	measureName := "metric"
	for _, label := range series.Labels {
		// the actual name of the series will be used as measure_name
		// and as such not added to the dimensions
		if label.Name == "__name__" {
			measureName = label.Value
			continue
		}
		name := prefix + "_" + label.Name + "_" + suffix
		dims = append(dims, &timestreamwrite.Dimension{
			DimensionValueType: aws.String("VARCHAR"),
			Name:               aws.String(name),
			Value:              aws.String(label.Value),
		})
	}

	common := &timestreamwrite.Record{
		Dimensions:       dims,
		MeasureName:      aws.String(measureName),
		MeasureValueType: aws.String("DOUBLE"),
		TimeUnit:         aws.String("MILLISECONDS"),
	}

	samples := make([]*timestreamwrite.Record, 0)
	for _, sample := range series.Samples {
		if math.IsNaN(sample.Value) {
			continue
		}
		value := fmt.Sprintf("%f", sample.Value)
		ts := fmt.Sprintf("%d", sample.Timestamp)
		samples = append(samples, &timestreamwrite.Record{
			MeasureValue: aws.String(value),
			Time:         aws.String(ts),
		})
	}

	return common, samples, nil
}
