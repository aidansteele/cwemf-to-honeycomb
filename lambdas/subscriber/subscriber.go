package main

import (
	"context"
	"encoding/json"
	"github.com/aidansteele/cwemf-to-honeycomb/cwemf"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/glassechidna/go-emf/emf"
	"github.com/pkg/errors"
)

func main() {
	s := &subscriber{}
	lambda.Start(s.handle)
}

type subscriber struct {
}

func (s *subscriber) handle(ctx context.Context, input *events.CloudWatchEvent) error {
	detail := cloudTrailDetail{}
	err := json.Unmarshal(input.Detail, &detail)
	if err != nil {
		return errors.WithStack(err)
	}

	region := detail.AwsRegion
	accountId := detail.RecipientAccountId
	logGroupName := detail.RequestParameters.LogGroupName
	emf.Emit(emf.MSI{
		"AccountId":    region,
		"Region":       region,
		"LogGroupName": logGroupName,
	})

	logs, err := cwemf.LogsApi(accountId, region)
	if err != nil {
		return err
	}

	err = cwemf.Subscribe(ctx, logs, region, logGroupName)
	return err
}

type cloudTrailDetail struct {
	RecipientAccountId string `json:"recipientAccountId"`
	AwsRegion          string `json:"awsRegion"`
	RequestParameters  struct {
		LogGroupName string `json:"logGroupName"`
	} `json:"requestParameters"`
}
