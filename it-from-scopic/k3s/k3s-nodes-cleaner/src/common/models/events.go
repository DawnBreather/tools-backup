package models

import (
	"connectors/aws/isqs"
	"github.com/aws/aws-sdk-go/service/sqs"
	"variables"
)

type Event struct{
	InstanceId  string
	InternalDns string
	State       string

	SqsMessage sqs.Message
}

func(e Event) GetPrintableState() string{
	if e.State == variables.EC2_LIFECYCLE_LAUNCHING {
		return "LAUNCHING"
	}

	if e.State == variables.EC2_LIFECYCLE_TERMINATING {
		return "TERMINATING"
	}

	return ""
}

func (e Event) DeleteEventFromSqs() (ok bool){
	return isqs.DeleteSqsMessage(e.SqsMessage)
}