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
	*tablewriter.Table
}

func NewInstancesTable() (table *InstancesTable) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Instance ID", "Instance Type", "Availability Zone", "Instance State", "Tagged?"})
	table = &InstancesTable{t}
	return
}

func (i *InstancesTable) Append(any interface{}) (err error) {
	if data, ok := any.(*ec2.DescribeInstancesOutput); ok {
		for idx, _ := range data.Reservations {
			for _, inst := range data.Reservations[idx].Instances {
				i.Table.Append([]string{
					getInstanceLabel(inst.Tags),
					*inst.InstanceId,
					*inst.InstanceType,
					*inst.Placement.AvailabilityZone,
					*inst.State.Name,
					taggedCheckmark(inst)})
			}
		}
		return
	}
	return fmt.Errorf("InstancesTable.Render requires *ec2.DescribeInstancesOutput")
}

type StateChangeTable struct {
	*tablewriter.Table
}

func NewStateChangeTable() (table *StateChangeTable) {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Region", "InstanceId", "Previous State", "Current State"})
	table = &StateChangeTable{t}
	return
}

func (i *StateChangeTable) Append(region string, changes []*ec2.InstanceStateChange) (err error) {
	for _, c := range changes {
		i.Table.Append([]string{region, *c.InstanceId, *c.PreviousState.Name, *c.CurrentState.Name})
	}
	return
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
