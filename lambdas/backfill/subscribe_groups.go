package main

import (
	"context"
	"github.com/aidansteele/cwemf-to-honeycomb/cwemf"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pkg/errors"
)

type subscribeGroupsInput struct {
	environment
	NextToken *string `json:",omitempty"`
}

type subscribeGroupsOutput struct {
	subscribeGroupsInput
	LastLogGroupName string
}

func (b *backfiller) subscribeGroups(ctx context.Context, input *subscribeGroupsInput) (*subscribeGroupsOutput, error) {
	logs, err := cwemf.LogsApi(input.AccountId, input.Region)
	if err != nil {
		return nil, err
	}

	describe, err := logs.DescribeLogGroupsWithContext(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		NextToken: input.NextToken,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	lastLogGroupName := ""
	for _, group := range describe.LogGroups {
		lastLogGroupName = *group.LogGroupName
		err = cwemf.Subscribe(ctx, logs, input.Region, lastLogGroupName)
		if err != nil {
			return nil, err
		}
	}

	return &subscribeGroupsOutput{
		subscribeGroupsInput: subscribeGroupsInput{
			environment: input.environment,
			NextToken:   describe.NextToken,
		},
		LastLogGroupName: lastLogGroupName,
	}, nil
}
