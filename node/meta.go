package node

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/euforia/metermaid/types"
)

const (
	// SpotTag is the name of the tag on an ec2 instance to determine
	// if the instance is a spot instance
	SpotTag = "aws:ec2spot:fleet-request-id"
)

// MetaProvider implements a node metadata provider
type MetaProvider interface {
	Meta() (types.Meta, error)
}

// NewMetaProvider returns a new node metadata provider based on the given
// platform
func NewMetaProvider(platform string) MetaProvider {
	switch platform {
	case PlatformAmazon:
		return &AWSNodeMeta{}
	}
	return nil
}

// AWSNodeMeta ...
type AWSNodeMeta struct{}

// NewAWSNodeMeta ...
func NewAWSNodeMeta() *AWSNodeMeta {
	return &AWSNodeMeta{}
}

// Meta returns metadata for the node
func (nodemeta *AWSNodeMeta) Meta() (types.Meta, error) {
	meta, err := getInstanceMeta()
	if err == nil {
		var reserve *ec2.Reservation
		reserve, err = describeInstance(meta["Region"], meta["InstanceID"])
		if err == nil {
			instance := reserve.Instances[0]
			//"InstanceLifecycle": "spot",
			for _, tag := range instance.Tags {
				meta[*tag.Key] = *tag.Value
			}
		}
	}

	return meta, err
}

func getInstanceMeta() (map[string]string, error) {
	sess, err := session.NewSession(&aws.Config{})
	if err == nil {
		svc := ec2metadata.New(sess)
		var ident ec2metadata.EC2InstanceIdentityDocument
		if ident, err = svc.GetInstanceIdentityDocument(); err == nil {
			meta := make(map[string]string)
			meta["InstanceType"] = ident.InstanceType
			meta["InstanceID"] = ident.InstanceID
			meta["AvailabilityZone"] = ident.AvailabilityZone
			meta["Region"] = ident.Region
			return meta, nil
		}
	}
	return nil, err
}

func describeInstance(region, instanceID string) (*ec2.Reservation, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err == nil {
		svc := ec2.New(sess)
		var resp *ec2.DescribeInstancesOutput
		resp, err = svc.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{aws.String(instanceID)}})
		if err == nil {
			return resp.Reservations[0], nil
		}
	}
	return nil, err
}
