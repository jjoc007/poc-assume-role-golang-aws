package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jjoc007/poc-assume-role-golang-aws/assume/aws/iam/assumerole"
)

func main()  {
	c := &assumerole.AssumeRole{
		AccessKeyID:     "XXXXX",
		SecretAccessKey: "XXXXXXX",
		Region:          "us-east-2",
		RoleARN:         "arn:aws:iam::167186109795:role/adl-rockstars-sts-assume-role-test",
		ExternalID:      "123456789",
	}
	assumeRoleService := assumerole.New(c)
	config, err := assumeRoleService.NewCredentialsV2()
	if err != nil{
		panic(err)
	}

	s3ClientwithAssume:= s3.New(*config)
	req:= s3ClientwithAssume.ListBucketsRequest(&s3.ListBucketsInput{})
	response, err := req.Send(context.Background())
	if err != nil{
		panic(err)
	}

	fmt.Println("Buckets con assume role")
	for _, b:= range response.Buckets{
		fmt.Println(*b.Name)
	}

	config2, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err)
	}

	s3ClientwithoutAssume:= s3.New(config2)
	req2:= s3ClientwithoutAssume.ListBucketsRequest(&s3.ListBucketsInput{})
	response2, err := req2.Send(context.Background())
	if err != nil{
		panic(err)
	}

	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("Buckets sin assume role")
	for _, b:= range response2.Buckets{
		fmt.Println(*b.Name)
	}
}
