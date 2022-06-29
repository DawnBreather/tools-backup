package idynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
	"reflect"
)

func (sw clientWrapper) PutEntry(o interface{}, tableName string) (ok bool){

	attributeValue, err := dynamodbattribute.MarshalMap(o)
	if err != nil {
		logrus.Errorf(`ERROR marshalling object into DynamoDB attributeValue: %v`, err)
		return false
	} else {

		inputEntry := &dynamodb.PutItemInput{
			Item:      attributeValue,
			TableName: aws.String(tableName),
		}

		_, err := sw.dynamodbClient.PutItem(inputEntry)
		if err != nil {
			logrus.Errorf(`ERROR putting entry into DynamoDB table {%s}: %v`, tableName, err)
			return false
		} else {
			logrus.Infof(`INFO entry {%s} {%v} added to table {%s}`, reflect.TypeOf(o).Name(), o, tableName)
			return true
		}
	}
}