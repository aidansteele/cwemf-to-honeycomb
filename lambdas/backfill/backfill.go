package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
	"os"
)

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigStateFromEnv})
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}

	b := &backfiller{
		orgs:     organizations.New(sess),
		roleName: os.Getenv("ROLE_NAME"),
	}

	switch os.Getenv("_HANDLER") {
	case "ListEnvironments":
		lambda.Start(b.listEnvironments)
	case "SubscribeGroups":
		lambda.Start(b.subscribeGroups)
	default:
		panic("unexpected _HANDLER")
	}
}

type backfiller struct {
	orgs     organizationsiface.OrganizationsAPI
	roleName string
}
