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
	"log"

	"github.com/spf13/cobra"

	"github.com/digitaljanitors/policia/aws"
	"github.com/digitaljanitors/policia/output"
)

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Police EC2 instances",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var ec2ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List EC2 Instances",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := aws.GetEC2Instances(cmd, args)
		if err != nil {
			log.Println(err)
		}

		table := output.NewInstancesTable()
		err = table.Render(resp)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(ec2Cmd)
	ec2Cmd.AddCommand(ec2ListCmd)

	// Flags for all EC2 subcommands
	ec2Cmd.PersistentFlags().BoolP("dry-run", "", false, "Do not make any changes, just show what would happen")

	// Flags for ec2ListCmd
	ec2ListCmd.Flags().BoolP("tagged", "t", false, "Show only tagged instances")
	ec2ListCmd.Flags().BoolP("show-stopped", "", false, "Show stopped instances also")
}
