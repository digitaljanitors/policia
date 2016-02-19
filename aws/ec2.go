package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetEC2Instances(cmd *cobra.Command, args []string) (resp *ec2.DescribeInstancesOutput, err error) {
	var (
		filters []*ec2.Filter
		stopped bool
		tagged  bool
	)

	stopped, err = cmd.Flags().GetBool("stopped")
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

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(viper.GetString("DefaultEC2Region"))})

	params := &ec2.DescribeInstancesInput{Filters: filters}
	resp, err = svc.DescribeInstances(params)

	return
}

//func StopInstances(DescribeInstancesOutput) {

//}

func IsTagged(inst *ec2.Instance) (is bool) {
	if _, ok := GetTag(inst.Tags, viper.GetString("TagName")); ok == nil {
		is = true
	}
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
