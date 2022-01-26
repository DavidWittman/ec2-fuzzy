package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	fzf "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/viper"
)

func getNameTag(i *types.Instance) string {
	for _, keys := range i.Tags {
		if *keys.Key == "Name" {
			return *keys.Value
		}
	}
	return ""
}

func findInstance(instances []*types.Instance, id string) (*types.Instance, error) {
	var findById bool
	// Assume that we're looking for an instance ID if the value starts with `i-`
	if strings.HasPrefix(id, "i-") {
		findById = true
	}

	for _, i := range instances {
		if findById {
			if *i.InstanceId == id {
				return i, nil
			}
		} else if getNameTag(i) == id {
			return i, nil
		}
	}

	return nil, fmt.Errorf("Could not find instance matching %s", id)
}

func fuzzyFindInstance(instances []*types.Instance) (*types.Instance, error) {
	displayFn := func(i int) string {
		name := getNameTag(instances[i])
		if name == "" {
			return *instances[i].InstanceId
		}
		return fmt.Sprintf("%s (%s)",
			name,
			*instances[i].InstanceId,
		)
	}

	previewWindow := func(i, w, h int) string {
		if i == -1 {
			return ""
		}
		instance := instances[i]
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
			string(instance.InstanceType),
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

	i, err := fzf.Find(instances, displayFn, fzf.WithPreviewWindow(previewWindow))
	if err != nil {
		return nil, err
	}

	return instances[i], nil
}

func flattenReservations(reservations []types.Reservation) []*types.Instance {
	var instances []*types.Instance
	for _, r := range reservations {
		for _, instance := range r.Instances {
			instances = append(instances, &instance)
		}
	}
	return instances
}

func main() {
	var instance *types.Instance

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	svc := ec2.NewFromConfig(cfg)

	// Only return running instances
	params := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running", "pending"},
			},
		},
	}

	resp, err := svc.DescribeInstances(context.TODO(), params)
	if err != nil {
		log.Fatal(err)
	}
	reservations := flattenReservations(resp.Reservations)

	if id := viper.GetString("instance"); id != "" {
		instance, err = findInstance(reservations, id)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println(viper.GetString("instance"))
		instance, err = fuzzyFindInstance(reservations)
		if err != nil {
			fmt.Println("No selection made. Exiting.")
			return
		}
	}

	var ip string
	if viper.GetBool("private") {
		ip = *instance.PrivateIpAddress
	} else if instance.PublicIpAddress == nil {
		log.Info("No public IP address found for instance. Trying to connect to the Private IP.")
		ip = *instance.PrivateIpAddress
	} else {
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
