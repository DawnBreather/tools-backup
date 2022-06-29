package idynamodb

import (
	"common/database"
	"common/models"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

func convertMapToStructInterface(item *dynamodb.GetItemOutput, tableName string) interface{}{

	var m map[string]interface{}

	err := dynamodbattribute.UnmarshalMap(item.Item, &m)
	if err != nil {
		logrus.Errorf("ERROR unmarshalling item {%v} from table {%s}: %v", item.Item, tableName, err)
		return nil
	}

	jsonByte, _ := json.Marshal(m)

	switch tableName {
	case database.INSTANCE_ID_TO_INTERNAL_DNS_TABLE:
		var instanceIdToIntenralDnsMapping models.InstanceIdToInternalDnsMapping
		json.Unmarshal(jsonByte, &instanceIdToIntenralDnsMapping)

		return instanceIdToIntenralDnsMapping

	default:
		logrus.Errorf("ERROR matching structure in DynamoDb to table {%s}", tableName)
		return nil
	}
}

func convertMapsToStructsInterface(item *dynamodb.ScanOutput, tableName string) interface{}{

	var m []map[string]interface{}

		err := dynamodbattribute.UnmarshalListOfMaps(item.Items, &m)
		if err != nil {
			logrus.Errorf("ERROR unmarshalling items from table {%s}: %v", tableName, err)
			return nil
		}

		jsonByte, _ := json.Marshal(m)

		switch tableName {
		case database.INSTANCE_ID_TO_INTERNAL_DNS_TABLE:
			var instanceIdToIntenralDnsMappings models.InstanceIdToIntenralDnsMappings
			json.Unmarshal(jsonByte, &instanceIdToIntenralDnsMappings)

			return instanceIdToIntenralDnsMappings

		default:
			logrus.Errorf("ERROR matching structure in DynamoDb to table {%s}", tableName)
			return nil
		}
}