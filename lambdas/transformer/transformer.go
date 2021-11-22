package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aidansteele/cwemf-to-honeycomb/cwemf"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"regexp"
	"strings"
)

func main() {
	var p = &transformer{}
	lambda.Start(p.handle)
}

type transformer struct {
}

func (p *transformer) handle(ctx context.Context, input *events.KinesisFirehoseEvent) (*events.KinesisFirehoseResponse, error) {
	j, _ := json.Marshal(input)
	fmt.Println(string(j))

	records := []events.KinesisFirehoseResponseRecord{}
	for _, record := range input.Records {
		records = append(records, p.transformRecord(record))
	}

	return &events.KinesisFirehoseResponse{Records: records}, nil
}

var lambdaStreamNameRegexp = regexp.MustCompile(`\d+/\d+/\d+/\[([^]]+)].+`)

func (p *transformer) transformRecord(record events.KinesisFirehoseEventRecord) (resp events.KinesisFirehoseResponseRecord) {
	resp = events.KinesisFirehoseResponseRecord{
		RecordID: record.RecordID,
		Result:   events.KinesisFirehoseTransformedStateOk,
		Data:     record.Data,
	}

	defer func() {
		if rerr := recover(); rerr != nil {
			resp.Result = events.KinesisFirehoseTransformedStateProcessingFailed
			fmt.Printf("error: %+v\n", rerr)
		}
	}()

	gzr, err := gzip.NewReader(bytes.NewReader(record.Data))
	if err != nil {
		panic(err)
	}

	d := events.CloudwatchLogsData{}
	err = json.NewDecoder(gzr).Decode(&d)
	if err != nil {
		panic(err)
	}

	if d.MessageType != "DATA_MESSAGE" {
		resp.Result = events.KinesisFirehoseTransformedStateDropped
		return
	}

	buf := &bytes.Buffer{}
	for _, event := range d.LogEvents {
		payload := map[string]interface{}{}
		err = json.Unmarshal([]byte(event.Message), &payload)
		if err != nil {
			panic(err)
		}

		emf := map[string]interface{}{
			"group":   d.LogGroup,
			"stream":  d.LogStream,
			"account": d.Owner,
		}

		// i find this helpful quite often
		match := lambdaStreamNameRegexp.FindStringSubmatch(d.LogStream)
		if len(match) > 0 {
			emf["functionversion"] = match[1]
		}

		for _, filter := range d.SubscriptionFilters {
			if strings.HasPrefix(filter, cwemf.FilterNamePrefix) {
				region := strings.TrimPrefix(filter, cwemf.FilterNamePrefix)
				emf["region"] = region
			}
		}

		delete(payload, "_aws")
		payload["emf"] = emf

		jsonline, err := json.Marshal(map[string]interface{}{
			"data": payload,
			"time": fmt.Sprintf("%d", event.Timestamp),
		})
		if err != nil {
			panic(err)
		}

		buf.Write(jsonline)
		buf.Write([]byte("\n"))
	}

	resp.Data = buf.Bytes()
	return
}
