package main

import (
	"context"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/pkg/errors"
)

type listEnvironmentsInput struct {
	Regions []string
}

type listEnvironmentsOutput struct {
	Environments []environment
}

type environment struct {
	AccountId string
	Region    string
}

func (b *backfiller) listEnvironments(ctx context.Context, input *listEnvironmentsInput) (*listEnvironmentsOutput, error) {
	accountIds := []string{}

	err := b.orgs.ListAccountsPagesWithContext(ctx, &organizations.ListAccountsInput{}, func(page *organizations.ListAccountsOutput, lastPage bool) bool {
		for _, account := range page.Accounts {
			accountIds = append(accountIds, *account.Id)
		}
		return !lastPage
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	regions := input.Regions
	if len(regions) == 0 {
		// TODO
	}

	environments := []environment{}
	for _, accountId := range accountIds {
		for _, region := range regions {
			environments = append(environments, environment{AccountId: accountId, Region: region})
		}
	}

	return &listEnvironmentsOutput{Environments: environments}, nil
}

