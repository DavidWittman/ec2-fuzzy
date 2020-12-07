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
	"github.com/spf13/viper"
)

func getNameTag(i *ec2.Instance) string {
	for _, keys := range i.Tags {
		if *keys.Key == "Name" {
			return *keys.Value
		}
	}
	return ""
}

func findInstance(reservations []*ec2.Reservation) (*ec2.Instance, error) {
	displayFn := func(i int) string {
		name := getNameTag(reservations[i].Instances[0])
		if name == "" {
			return *reservations[i].Instances[0].InstanceId
		}
		return fmt.Sprintf("%s (%s)",
			name,
			*reservations[i].Instances[0].InstanceId,
		)
	}

	previewWindow := func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		instance := reservations[i].Instances[0]
		name := getNameTag(instance)

		publicIp := "None"
		if instance.PublicIpAddress != nil {
			publicIp = *instance.PublicIpAddress
		}

		values := []string{
			name,
			*instance.InstanceId,
			publicIp,
			*instance.PrivateIpAddress,
			*instance.InstanceType,
			*instance.ImageId,
			instance.LaunchTime.String(),
		}

		// force preview window to bottom
		margin := 3
		newlines := strings.Repeat("\n", h-margin-len(values))

		// Convert []string to []interface so we can pass it to Sprintf
		// https://stackoverflow.com/a/12334902
		formattedValues := make([]interface{}, len(values)+1)
		formattedValues[0] = newlines
		for i, v := range values {
			formattedValues[i+1] = v
		}

		return fmt.Sprintf("%sName: %s\nID: %s\nPublic IP: %s\nPrivate IP: %s\nType: %s\nAMI: %s\nLaunch Time: %s",
			formattedValues...,
		)

	}

	i, err := fzf.Find(reservations, displayFn, fzf.WithPreviewWindow(previewWindow))
	if err != nil {
		return nil, err
	}

	return reservations[i].Instances[0], nil
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

	instance, err := findInstance(resp.Reservations)
	if err != nil {
		fmt.Println("No selection made. Exiting.")
		return
	}

	var ip string
	if viper.GetBool("private") {
		ip = *instance.PrivateIpAddress
	} else {
		if instance.PublicIpAddress == nil {
			log.Fatal("No public IP address found for instance. Set --private flag to connect to the Private IP.")
		}
		ip = *instance.PublicIpAddress
	}

	ssh, err := exec.LookPath("ssh")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: use a builder
	var cmd []string
	username := viper.GetString("user")
	if username != "" {
		cmd = []string{ssh, "-l", username, ip}
	} else {
		cmd = []string{ssh, ip}
	}

	fmt.Println(strings.Join(cmd, " "))
	err = syscall.Exec(ssh, cmd, os.Environ())
	fmt.Println(err)
	return
}
