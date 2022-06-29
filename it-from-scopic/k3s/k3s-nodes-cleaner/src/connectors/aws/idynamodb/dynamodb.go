package idynamodb

import (
	"connectors/aws/iauth"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Client clientWrapper

func Initilize() {

	Client = clientWrapper{
		dynamodbClient: dynamodb.New(iauth.GetCurrentAuthSession()),
	}
}

type clientWrapper struct{
	dynamodbClient *dynamodb.DynamoDB
}