#!/bin/sh
set -eux

export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=0

for func in lambdas/*; do
  cd "$func"
  go build -ldflags="-s -w -buildid=" -trimpath -o bootstrap
  cd -
done

sam deploy --template example-deployment.yml

for region in ap-southeast-2 us-east-1 us-west-2
do
  stackit up \
    --stack-name cwemf-to-honeycomb-destination \
    --template destinations.yml \
    --region $region \
    FirehoseArn=arn:aws:firehose:ap-southeast-2:000000000:deliverystream/cwemf-to-honeycomb-Central-ISbULrZYW1Ur
done

