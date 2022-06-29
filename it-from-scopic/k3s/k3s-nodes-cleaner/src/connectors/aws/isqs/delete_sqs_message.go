package isqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

func DeleteSqsMessage(sqsMessage sqs.Message) (ok bool){

	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(*receive_params.QueueUrl),  // Required
		ReceiptHandle: sqsMessage.ReceiptHandle, // Required
	}

	_, err := _sqsSession.DeleteMessage(deleteParams) // No response returned when successed.
	if err != nil {
		logrus.Errorf("ERROR removing message {%s} from SQS URL {%s}: %v", sqsMessage.MessageId, *receive_params.QueueUrl, err)
		return false
	}

	return true
}