package isqs

import (
	"connectors/aws/iauth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var _sqsSession *sqs.SQS

// Receive message
var receive_params = &sqs.ReceiveMessageInput{
	QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/073283464250/k3s-develop-lifecycle-hook"),
	MaxNumberOfMessages: aws.Int64(1),  // 一次最多取幾個 message
	VisibilityTimeout:   aws.Int64(5), // 如果這個 message 沒刪除，下次再被取出來的時間
	WaitTimeSeconds:     aws.Int64(5), // long polling 方式取，會建立一條長連線並且等在那邊，直到 SQS 收到新 message 回傳給這條連線才中斷
}

func InitializeSqs() *sqs.SQS{
	_sqsSession = sqs.New(iauth.GetCurrentAuthSession())

	return _sqsSession
}

