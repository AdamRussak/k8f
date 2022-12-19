package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSAssumeRoleAPI interface {
	AssumeRole(ctx context.Context,
		params *sts.AssumeRoleInput,
		optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}
