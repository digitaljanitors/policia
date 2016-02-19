package output

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/digitaljanitors/policia/aws"
	"github.com/olekukonko/tablewriter"
)

type AsciiTable interface {
	Render(interface{})
}

type InstancesTable struct {
	table *tablewriter.Table
}

func NewInstancesTable() (table *InstancesTable) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Instance ID", "Instance Type", "Availability Zone", "Instance State", "Tagged?"})
	table = &InstancesTable{t}
	return
}

func (i *InstancesTable) Render(any interface{}) (err error) {
	if data, ok := any.(*ec2.DescribeInstancesOutput); ok {
		for idx, _ := range data.Reservations {
			for _, inst := range data.Reservations[idx].Instances {
				i.table.Append([]string{
					getInstanceLabel(inst.Tags),
					*inst.InstanceId,
					*inst.InstanceType,
					*inst.Placement.AvailabilityZone,
					*inst.State.Name,
					taggedCheckmark(inst)})
			}
		}
		i.table.Render()
		return
	}
	return fmt.Errorf("InstancesTable.Render requires *ec2.DescribeInstancesOutput")
}

func taggedCheckmark(inst *ec2.Instance) string {
	if aws.IsTagged(inst) {
		return "\u2713"
	}
	return ""
}

func getInstanceLabel(tags []*ec2.Tag) (label string) {
	tag, err := aws.GetTag(tags, "Name")
	if err == nil {
		label = *tag.Value
	}
	return
}
