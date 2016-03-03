package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetEC2Instances(cmd *cobra.Command, args []string) (resp map[string]*ec2.DescribeInstancesOutput, err error) {
	var (
		filters []*ec2.Filter
		stopped bool
		tagged  bool
		r       *ec2.DescribeInstancesOutput
	)

	resp = make(map[string]*ec2.DescribeInstancesOutput)

	stopped, err = cmd.Flags().GetBool("show-stopped")
	if !stopped {
		filters = append(filters, &ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
			},
		})
	}

	tagged, err = cmd.Flags().GetBool("tagged")
	if tagged {
		filters = append(filters, &ec2.Filter{
			Name: aws.String("tag-key"),
			Values: []*string{
				aws.String(viper.GetString("TagName")),
			},
		})
	}

	regions, err := GetRegions(cmd, args)
	if err != nil {
		return
	}

	for _, region := range regions.Regions {
		svc := ec2.New(session.New(), aws.NewConfig().WithRegion(*region.RegionName))

		params := &ec2.DescribeInstancesInput{Filters: filters}
		r, err = svc.DescribeInstances(params)
		if err != nil {
			return
		}

		resp[*region.RegionName] = r
	}

	return
}

func StopInstances(cmd *cobra.Command, args []string) (resp *ec2.StopInstancesOutput, err error) {
	var (
		instanceIds []*string
		dryrun      bool
		data        map[string]*ec2.DescribeInstancesOutput
	)

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(viper.GetString("DefaultEC2Region"))})

	dryrun, err = cmd.Flags().GetBool("dryrun")

	data, err = GetEC2Instances(cmd, args)
	if err != nil {
		for region, _ := range data {
			for idx, _ := range data[region].Reservations {
				for _, inst := range data[region].Reservations[idx].Instances {
					if IsTagged(inst) && IsRunning(inst) {
						instanceIds = append(instanceIds, inst.InstanceId)
					}
				}
			}
		}

		params := &ec2.StopInstancesInput{
			InstanceIds: instanceIds,
			DryRun:      aws.Bool(dryrun),
		}

		resp, err = svc.StopInstances(params)
	}

	return
}

func IsTagged(inst *ec2.Instance) (is bool) {
	if _, ok := GetTag(inst.Tags, viper.GetString("TagName")); ok == nil {
		is = true
	}
	return
}

func IsRunning(inst *ec2.Instance) (is bool) {
	is = (*inst.State.Code == int64(16))
	return
}

func GetTag(tags []*ec2.Tag, name string) (tag *ec2.Tag, err error) {
	for _, t := range tags {
		if *t.Key == name {
			tag = t
		}
	}
	if tag == nil {
		err = fmt.Errorf("Tag not found")
	}
	return
}

func GetRegions(cmd *cobra.Command, args []string) (resp *ec2.DescribeRegionsOutput, err error) {
	var (
		svc         *ec2.EC2
		input       *ec2.DescribeRegionsInput
		regions     []string
		regionNames []*string
	)

	regions, err = cmd.Flags().GetStringSlice("region")
	if err != nil {
		return
	}

	for _, v := range regions {
		regionNames = append(regionNames, aws.String(v))
	}

	if len(regionNames) > 0 {
		input = &ec2.DescribeRegionsInput{RegionNames: regionNames}
	} else {
		input = &ec2.DescribeRegionsInput{}
	}

	svc = ec2.New(session.New(), aws.NewConfig().WithRegion(viper.GetString("DefaultEC2Region")))

	resp, err = svc.DescribeRegions(input)
	return
}
