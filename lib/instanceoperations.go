package lib

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const sandboxNameTagName = "SandboxName"
const isSandboxTagName = "IsSandbox"
const isSandboxTagValue = "True"

type Instance struct {
	InstanceId string
	Name       string
}

func ListInstances(client *ec2.Client) ([]Instance, error) {
	response, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", isSandboxTagName)),
				Values: []string{isSandboxTagValue},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to describe instances, %v", err)
	}

	var result []Instance
	if len(response.Reservations) > 0 {
		for _, reservation := range response.Reservations {
			for _, instance := range reservation.Instances {
				result = append(result, Instance{
					*instance.InstanceId,
					"TODO",
				})
			}
		}
	}

	return result, nil
}
