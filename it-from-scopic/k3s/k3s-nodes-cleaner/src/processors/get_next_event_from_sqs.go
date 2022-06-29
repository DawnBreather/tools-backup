package processors

import (
	"common/models"
	"connectors/aws/iec2"
	"connectors/aws/isqs"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func GetNextEventFromSqs() (e *models.Event) {
	sqsMessages := isqs.ReadSqsMessages()

	for _, msg := range sqsMessages {

		var content map[string]interface{}

		err := json.Unmarshal([]byte(*msg.Body), &content)

		if err != nil {
			logrus.Errorf(`ERROR parsing SQS message body: %v`, err)
			continue
		} else {

			var instance_id string
			var state string

			if val, ok := content["LifecycleTransition"]; ok {
				//fmt.Printf(`Life cycle: %s`+"\n", val)

				state = val.(string)
			}

			if val, ok := content["EC2InstanceId"]; ok {
				//fmt.Printf(`Instace ID: %s`+"\n", val)

				if val.(string) != "" {
					instance_id = val.(string)
				} else {
					logrus.Errorf("ERROR detected: {instance_id} is empty")
					continue
				}
			}

			e = &models.Event{
				InstanceId:  instance_id,
				InternalDns: iec2.GetInstancePrivateDnsByInstanceId(instance_id),
				State:       state,
				SqsMessage:  *msg,
			}

			logrus.Infof("INFO ( %s ) ( INSTANCE_ID: %s ) ( INSTANCE_PRIVATE_DNS: %s )", e.GetPrintableState(), instance_id, e.InternalDns)

			return e
		}
	}

	return nil
}