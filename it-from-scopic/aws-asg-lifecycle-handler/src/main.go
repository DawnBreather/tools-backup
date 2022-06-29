package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

const (
	QueueUrl    = "https://sqs.us-east-1.amazonaws.com/073283464250/k3s-develop-lifecycle-hook"
	Region      = "ap-northeast-1"
)

var receive_params = &sqs.ReceiveMessageInput{
	QueueUrl:            aws.String(QueueUrl),
	MaxNumberOfMessages: aws.Int64(1),  // 一次最多取幾個 message
	VisibilityTimeout:   aws.Int64(30), // 如果這個 message 沒刪除，下次再被取出來的時間
	WaitTimeSeconds:     aws.Int64(20), // long polling 方式取，會建立一條長連線並且等在那邊，直到 SQS 收到新 message 回傳給這條連線才中斷
}

type queue sqs.SQS

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("AKIARCEAI5Q5FERNRYBX", "FQ8zbX6p+giEE2+tREuOJTVaZepITBh8LcIgkqSj", ""),
		MaxRetries:  aws.Int(5),
	})

	if err != nil {
		log.Panicf("FATAL authentication in AWS failed: %v", err)
	}

	var svc queue
	svc = queue(*sqs.New(sess))

	// Send message
	/*send_params := &sqs.SendMessageInput{
		MessageBody:  aws.String("message body"), // Required
		QueueUrl:     aws.String(QueueUrl),       // Required
		DelaySeconds: aws.Int64(3),               // (optional) 傳進去的 message 延遲 n 秒才會被取出, 0 ~ 900s (15 minutes)
	}
	send_resp, err := svc.SendMessage(send_params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[Send message] \n%v \n\n", send_resp)*/

	var counter = 0

	for {
		counter ++
		// Receive message
		receive_resp, err := svc.ReceiveMessage(receive_params)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("[Receive message] \n%v \n\n", receive_resp)

		if counter == 20 {
			break
		}

	}

	// Delete message
	/*for _, message := range receive_resp.Messages {
		delete_params := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(QueueUrl),  // Required
			ReceiptHandle: message.ReceiptHandle, // Required

		}
		_, err := svc.DeleteMessage(delete_params) // No response returned when successed.
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
	}*/
}

func (q queue) receiveMessage() *sqs.ReceiveMessageOutput{
	receive_resp, err := sqs.SQS(q).ReceiveMessage(receive_params)
	if err != nil {
		logrus.Errorf("ERROR reading message: %v", err)
		return nil
	}

	return receive_resp
}