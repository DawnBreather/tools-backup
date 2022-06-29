package idynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"
)

func (sw clientWrapper) DeleteEntryByKey(tableName, keyName, keyValue string) (ok bool){
	_, err := sw.dynamodbClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S:    aws.String(keyValue),
			},
		},
	})

	if err != nil {
		logrus.Errorf(`ERROR deleting item by key {%s}=={%s} from DynamoDB -> table {%s}: %v`, keyName, keyValue, tableName, err)
		return false
	}

	logrus.Infof("INFO delete item from DynamoDB -> table {%s} by key {%s}=={%s}: succeeded", tableName, keyName, keyValue)
	return true
}