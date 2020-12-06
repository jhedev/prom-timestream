package adapter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/prometheus/prometheus/prompb"
)

func seriesToRecords(series *prompb.TimeSeries) (*timestreamwrite.Record, []*timestreamwrite.Record, error) {
	dims := make([]*timestreamwrite.Dimension, len(series.Labels))

	for i, label := range series.Labels {
		dims[i] = &timestreamwrite.Dimension{
			DimensionValueType: aws.String("VARCHAR"),
			Name:               aws.String(label.Name),
			Value:              aws.String(label.Value),
		}
	}

	common := &timestreamwrite.Record{
		Dimensions:       dims,
		MeasureName:      aws.String("metric"),
		MeasureValueType: aws.String("DOUBLE"),
	}

	samples := make([]*timestreamwrite.Record, len(series.Samples))

	for i, sample := range series.Samples {
		samples[i] = &timestreamwrite.Record{
			MeasureValue: aws.String(fmt.Sprintf("%f", sample.Value)),
		}
	}

	return common, samples, nil
}
