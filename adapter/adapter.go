package adapter

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-sdk-go/service/timestreamquery/timestreamqueryiface"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/aws/aws-sdk-go/service/timestreamwrite/timestreamwriteiface"
	"github.com/prometheus/prometheus/prompb"
	log "github.com/sirupsen/logrus"
)

const (
	defaultTimeout = 30 * time.Second
)

type Adapter struct {
	logger  log.FieldLogger
	timeout time.Duration

	databaseName string
	tableName    string
	writeClient  timestreamwriteiface.TimestreamWriteAPI
	queryClient  timestreamqueryiface.TimestreamQueryAPI
}

type Option func(*Adapter) error

func New(databaseName, tableName string, sess *session.Session, opts ...Option) (*Adapter, error) {
	a := &Adapter{
		logger:  log.New().WithField("component", "adapter"),
		timeout: defaultTimeout,

		databaseName: databaseName,
		tableName:    tableName,
		writeClient:  timestreamwrite.New(sess),
		queryClient:  timestreamquery.New(sess),
	}

	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	_, err := a.writeClient.DescribeTableWithContext(ctx, &timestreamwrite.DescribeTableInput{
		DatabaseName: aws.String(a.databaseName),
		TableName:    aws.String(a.tableName),
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Adapter) writeSeries(ctx context.Context, series *prompb.TimeSeries) error {
	common, records, err := seriesToRecords(series)
	if err != nil {
		return err
	}

	ctxx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	_, err = a.writeClient.WriteRecordsWithContext(ctxx, &timestreamwrite.WriteRecordsInput{
		DatabaseName:     aws.String(a.databaseName),
		TableName:        aws.String(a.tableName),
		CommonAttributes: common,
		Records:          records,
	})
	return err
}

func (a *Adapter) Write(ctx context.Context, req *prompb.WriteRequest) error {
	for _, series := range req.Timeseries {
		if err := a.writeSeries(ctx, series); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) Read(ctx context.Context, req *prompb.ReadRequest) (*prompb.ReadResponse, error) {
	return nil, nil
}
