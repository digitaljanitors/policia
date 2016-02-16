// Copyright © 2016 Chris McNabb <raizyr@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tagged  bool
	stopped bool
)

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		var filters []*ec2.Filter

		if !stopped {
			filters = append(filters, &ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			})
		}

		if tagged {
			filters = append(filters, &ec2.Filter{
				Name: aws.String("tag-key"),
				Values: []*string{
					aws.String("UnaMordida"),
				},
			})
		}

		svc := ec2.New(session.New(), &aws.Config{Region: aws.String(viper.GetString("aws.default_region"))})

		params := &ec2.DescribeInstancesInput{Filters: filters}
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			log.Println(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Instance ID", "Instance Type", "Availability Zone", "Instance State", "Tagged?"})

		// resp has all of the response data, pull out instance IDs:
		for idx, _ := range resp.Reservations {
			for _, inst := range resp.Reservations[idx].Instances {
				table.Append([]string{
					getInstanceLabel(inst.Tags),
					*inst.InstanceId,
					*inst.InstanceType,
					*inst.Placement.AvailabilityZone,
					*inst.State.Name,
					isTagged(inst)})
			}
		}

		table.Render()
	},
}

func init() {
	listCmd.AddCommand(ec2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ec2Cmd.Flags().BoolVarP(&tagged, "tagged", "t", false, "Show only tagged instances")
	ec2Cmd.Flags().BoolVarP(&stopped, "show-stopped", "", false, "Show stopped instances also")

}

func isTagged(inst *ec2.Instance) string {
	_, err := getTag(inst.Tags, "UnaMordida")
	if err == nil {
		return "\u2713"
	}
	return ""
}

func getInstanceLabel(tags []*ec2.Tag) (label string) {
	tag, err := getTag(tags, "Name")
	if err == nil {
		label = *tag.Value
	}
	return
}

func getTag(tags []*ec2.Tag, name string) (tag *ec2.Tag, err error) {
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
