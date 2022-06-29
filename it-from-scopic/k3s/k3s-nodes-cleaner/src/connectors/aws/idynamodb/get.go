package idynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"
)

func (sw clientWrapper) GetEntryByKey(tableName, keyName, keyValue string) (item interface{}){
	result, err := sw.dynamodbClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S:    aws.String(keyValue),
			},
		},
	})

	if err != nil {
		logrus.Errorf(`ERROR getting item by key {%s}=={%s} from DynamoDB -> table {%s}: %v`, keyName, keyValue, tableName, err)
		return nil
	}

	logrus.Infof("INFO got item found in DynamoDB -> table {%s} by key {%s}=={%s}: %v", tableName, keyName, keyValue, result.Item)
	return convertMapToStructInterface(result, tableName)
}

func (sw clientWrapper) GetAllEntries(tableName string) (items interface{}){

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := sw.dynamodbClient.Scan(params)

	if err != nil {
		logrus.Errorf(`ERROR getting all entries from DynamoDB -> table {%s}: %v`, tableName, err)
		return nil
	}

	return convertMapsToStructsInterface(result, tableName)

}