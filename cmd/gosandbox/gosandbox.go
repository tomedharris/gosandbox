/*
Gosandbox manages sandbox instances in AWS.
It allows you to bring a sandbox instance up and down and connect to it.

Usage:

	gosandbox [flags] [action ...]

The flags are:

	-c
		Connect to the sandbox when bringing the sandbox up.

The actions are:

	up
		Create the sandbox if one does not already exist.

	down
		Terminate the sandbox if it exists.

	refresh
		Terminate the sandbox if it exists and bring up a new one.

	connect
		Connect to an already existing sandbox.
*/
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const tagName = "SandboxName"
const tagValue = "Sandbox"

var region = "eu-west-1"
var newInstanceAmi = "ami-0ebfed9ccce07b642"
var newInstanceType = types.InstanceTypeT2Micro
var newInstanceKeyName = "mac-eu-west-1"

func main() {
	// TODO: Implement up (with -c)
	// TODO: Implement down
	// TODO: Implement refresh (with -c)
	// TODO: Implement connect

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := ec2.NewFromConfig(cfg)

	// check for existing sandbox...
	log.Printf("Checking for existing sandbox in %s region...\n", region)
	err = clearAnyExistingSandboxes(client)
	if err != nil {
		log.Fatalf("error clearing existing sandbox: %s", err)
	}

	// launch new sandbox...
	log.Println("Launching new sandbox in eu-west-1 region...")
	err = runNewSandbox(client)
	if err != nil {
		log.Fatalf("error launching new sandbox: %s", err)
	}

	// start new session...

	// on exit, terminate sandbox...

}

// Run a new Sandbox Instance with the correct tags for later identification.
// This will return any errors encountered during the operation.
// If a new instance is created, the Instance ID will be logged.
func runNewSandbox(client *ec2.Client) error {
	response, err := client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		ImageId:      aws.String(newInstanceAmi),
		InstanceType: newInstanceType,
		KeyName:      aws.String(newInstanceKeyName),
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeInstance,
			Tags: []types.Tag{{
				Key:   aws.String(tagName),
				Value: aws.String(tagValue),
			}},
		}},
	})
	if err != nil {
		return fmt.Errorf("failed to run instances, %v", err)
	}

	for _, instance := range response.Instances {
		log.Printf("created new sandbox instance: %s\n", *instance.InstanceId)
	}

	return nil
}

// Clear existing Sandbox Instances according to the assigned tags.
// This will return any errors encountered during the operation.
// If any instances are deleted, the Instance ID will be logged.
func clearAnyExistingSandboxes(client *ec2.Client) error {
	response, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", tagName)),
				Values: []string{tagValue},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to describe instances, %v", err)
	}

	// if existing, delete sandbox...
	if len(response.Reservations) > 0 {
		for _, reservation := range response.Reservations {
			for _, instance := range reservation.Instances {
				log.Printf("found existing sandbox, %s clearing up...\n", *instance.InstanceId)

				_, err := client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
					InstanceIds: []string{*instance.InstanceId},
				})
				if err != nil {
					return fmt.Errorf("failed to terminate instance %s: %w", *instance.InstanceId, err)
				}
			}
		}
	}

	return nil
}

// TODO: Instance Info
// func printInstanceInfo(instance types.Instance) {
// 	fmt.Printf("Instance ID: %s\n", *instance.InstanceId)
// 	// fmt.Printf("Instance Type: %s\n", *instance.InstanceType)
// 	// fmt.Printf("State: %s\n", *instance.State.Name)
// 	// fmt.Printf("Public IP: %s\n", aws.StringValue(instance.PublicIpAddress))
// 	// fmt.Printf("Private IP: %s\n", aws.StringValue(instance.PrivateIpAddress))
// 	// fmt.Println("Tags:")
// 	// for _, tag := range instance.Tags {
// 	// fmt.Printf("  %s: %s\n", *tag.Key, *tag.Value)
// 	// }
// 	fmt.Println()
// }
