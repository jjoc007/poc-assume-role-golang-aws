package assumerole

import (
	aws_v2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	stscreds_v2 "github.com/aws/aws-sdk-go-v2/aws/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog/log"
)

type Operations interface {
	NewCredentialsV2() (*aws_v2.Config, error)
}

// New creates and returns a new aws assume role service instance
func New(configAssumeRole *AssumeRole) Operations {
	return &assumeRoleService{
		assumeRoleConfig: configAssumeRole,
	}
}

type assumeRoleService struct {
	assumeRoleConfig *AssumeRole
}

func (a *assumeRoleService) NewCredentialsV2() (config *aws_v2.Config, err error) {
	config, err = a.validateKeysV2()
	if err != nil {
		log.Err(err)
		return
	}

	stsSvc := sts.New(*config)
	stsCredProvider := stscreds_v2.NewAssumeRoleProvider(stsSvc,
		a.assumeRoleConfig.RoleARN,
		func(options *stscreds_v2.AssumeRoleProviderOptions) {
			options.ExternalID = aws_v2.String(a.assumeRoleConfig.ExternalID)
		})
	config.Credentials = aws_v2.CredentialsProvider(stsCredProvider)
	return
}

func (a *assumeRoleService) validateKeysV2() (*aws_v2.Config, error) {
	if a.assumeRoleConfig.AccessKeyID != "" && a.assumeRoleConfig.SecretAccessKey != "" {
		config, err := external.LoadDefaultAWSConfig(
			external.WithCredentialsProvider{
				CredentialsProvider: aws_v2.StaticCredentialsProvider{
					Value: aws_v2.Credentials{
						AccessKeyID:     a.assumeRoleConfig.AccessKeyID,
						SecretAccessKey: a.assumeRoleConfig.SecretAccessKey,
					},
				},
			},
		)
		if err != nil {
			return nil, err
		}
		return &config, err
	}
	config, err := external.LoadDefaultAWSConfig()
	config.Region = a.assumeRoleConfig.Region
	return &config, err
}