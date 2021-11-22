package cwemf

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/pkg/errors"
	"os"
)

func LogsApi(accountId, region string) (cloudwatchlogsiface.CloudWatchLogsAPI, error) {
	baseSess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, os.Getenv("ROLE_NAME"))
	creds := stscreds.NewCredentials(baseSess, roleArn)
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigStateFromEnv,
		Config: *aws.NewConfig().
			WithLogLevel(aws.LogDebugWithHTTPBody).
			WithCredentials(creds).
			WithRegion(region),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cloudwatchlogs.New(sess), nil
}

func Subscribe(ctx context.Context, logs cloudwatchlogsiface.CloudWatchLogsAPI, region, logGroupName string) error {
	destinationArn := fmt.Sprintf("arn:aws:logs:%s:%s:destination:%s", region, os.Getenv("DESTINATION_ACCOUNT_ID"), os.Getenv("DESTINATION_NAME"))

	logs.DeleteSubscriptionFilter(&cloudwatchlogs.DeleteSubscriptionFilterInput{
		FilterName:   aws.String("cwemf-to-honeycomb"),
		LogGroupName: &logGroupName,
	})

	_, err := logs.PutSubscriptionFilterWithContext(ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		DestinationArn: aws.String(destinationArn),
		FilterName:     aws.String(FilterNamePrefix + region),
		FilterPattern:  aws.String(`{ $._aws.Timestamp > 0 }`),
		LogGroupName:   &logGroupName,
	})
	return errors.WithStack(err)
}

var FilterNamePrefix = "cwemf-to-honeycomb-"
