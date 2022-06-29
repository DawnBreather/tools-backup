package iec2

import (
	"connectors/aws/iauth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

var _ec2Session *ec2.EC2

func InitializeEc2() *ec2.EC2{
	_ec2Session = ec2.New(iauth.GetCurrentAuthSession())

	return _ec2Session
}

func GetInstances() (instances []*ec2.Instance, ok bool){
	resp, err := _ec2Session.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		logrus.Errorf("ERROR fetching description for EC2 instances: %v", err)
		return nil, false
	}

	for _, reservation := range resp.Reservations{
		instances = append(instances, reservation.Instances...)
	}

	logrus.Infof("INFO fetched {%d} EC2 instances", len(instances))


	return instances, true
}

func GetInstancePrivateDnsByInstanceId(instanceId string) string{
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(`instance-id`),
				Values: []*string{aws.String(instanceId)},
			},
		},
	}

	resp, err := _ec2Session.DescribeInstances(params)
	if err != nil {
		logrus.Errorf("ERROR fetching description for EC2 instance {%s}: %v", instanceId, err)
		return ""
	}

	if len (resp.Reservations) > 0 {
		if len(resp.Reservations[0].Instances) > 0 {
			return *resp.Reservations[0].Instances[0].PrivateDnsName
		}
	}

	return ""
}