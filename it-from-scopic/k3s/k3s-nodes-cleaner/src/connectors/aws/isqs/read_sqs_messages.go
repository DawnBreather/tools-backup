package isqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

func ReadSqsMessages() []*sqs.Message{
	receive_resp, err := _sqsSession.ReceiveMessage(receive_params)
	if err != nil {
		logrus.Errorf("ERROR reading message from {%s}: %v", receive_params.QueueUrl, err)
	}

	return receive_resp.Messages
}