package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	fzf "github.com/ktr0731/go-fuzzyfinder"
)

func getNameTag(i *ec2.Instance) string {
	for _, keys := range i.Tags {
		if *keys.Key == "Name" {
			return *keys.Value
		}
	}
	return ""
}

func main() {
	svc := ec2.New(session.Must(session.NewSession()))

	// Only return running instances
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		log.Fatal(err)
	}

	displayFn := func(i int) string {
		name := getNameTag(resp.Reservations[i].Instances[0])
		if name == "" {
			return *resp.Reservations[i].Instances[0].InstanceId
		}
		return fmt.Sprintf("%s (%s)",
			name,
			*resp.Reservations[i].Instances[0].InstanceId,
		)
	}

	previewWindow := func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		instance := resp.Reservations[i].Instances[0]
		name := getNameTag(instance)

		// force preview window to bottom. TODO: don't hard-code value length :-\
		values := 7
		margin := 3
		newlines := strings.Repeat("\n", h-margin-values)

		publicIp := "None"
		if instance.PublicIpAddress != nil {
			publicIp = *instance.PublicIpAddress
		}

		return fmt.Sprintf("%sName: %s\nID: %s\nPublic IP: %s\nPrivate IP: %s\nType: %s\nAMI: %s\nLaunch Time: %s",
			newlines,
			name,
			*instance.InstanceId,
			publicIp,
			*instance.PrivateIpAddress,
			*instance.InstanceType,
			*instance.ImageId,
			*instance.LaunchTime,
		)

	}

	// TODO: not sure if a reservation can have multiple instances
	i, err := fzf.Find(resp.Reservations, displayFn, fzf.WithPreviewWindow(previewWindow))
	if err != nil {
		fmt.Println("No selection made. Exiting.")
		return
	}

	ssh, err := exec.LookPath("ssh")
	if err != nil {
		log.Fatal(err)
	}
	err = syscall.Exec(ssh, []string{ssh, *resp.Reservations[i].Instances[0].PrivateIpAddress}, os.Environ())
	fmt.Println(err)
	return
}
