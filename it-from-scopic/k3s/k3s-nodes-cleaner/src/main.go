package main

import (
	"common/config"
	"connectors/aws/iauth"
	"connectors/aws/iec2"
	"connectors/ik8s"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"time"
)


// TODO Implement Sentry.IO for exceptions tracing
//
func main() {

	// Initialize application's configuration parameters
	config.Initialize()

	// Authenticate in AWS
	iauth.Authenticate()
	// Initialize EC2 manager
	iec2.InitializeEc2()
	// Initialize k8s manager
	ik8s.Initialize()


	for {

		// This slice will be containing the names of all the nodes we are to remove from k8s
		var nodesToDelete []string

		// Fetch all nodes from k8s with NotReady status
		nodes := ik8s.GetAllNodes()

		if len(nodes) > 0 {

			// Fetch all EC2 instances in the region
			instances, ok := iec2.GetInstances()

			if ok {

				// For each of NotReady node
				for _, n := range nodes {

					// Let's use flag defining whether the nodes physically exists in EC2 or not
					var nodePhysicallyExists = false

					// Important note: the names of our nodes in k8s correspond to the PrivateDnsNames of the EC2 instances
					// So let's filter EC2 instances and identify whether there is an EC2 instance with the same PrivateDnsName as our unready k8s node's name
					foundInstances := funk.Filter(instances, func(instance *ec2.Instance) bool {
						if *instance.PrivateDnsName == n.Name {
							if *instance.State.Name == "running" {
								return true
							}
						}
						return false
					}).([]*ec2.Instance)

					// if positive - that means NotReady k8s node physically exists and we shouldn't remove it
					if len(foundInstances) > 0 {
						nodePhysicallyExists = true
					}

					// but otherwise - we just should remove it
					if !nodePhysicallyExists {
						nodesToDelete = append(nodesToDelete, n.Name)
						logrus.Infof("INFO k8s node {%s} is marked for removal", n.Name)
					} else {
						logrus.Warnf("WARNING k8s node {%s} is NotReady, but physically existing. We suggest checking the EC2 instance ID: %s", n.Name, *foundInstances[0].InstanceId)
					}

					if len(nodesToDelete) == 0 {
						logrus.Infof("INFO no k8s nodes marked for removal")
					}
				}
			}
		}

		// ok, finally removing the inexistent nodes
		ik8s.DeleteNodesByNames(nodesToDelete)

		// repeat the whole workflow after specified interval
		time.Sleep(config.WORKFLOW.CheckInterval)
	}



	//
	// Fancy implementation with SQS and DynamoDB (yep, more granular and proper, but more complicated either)
	//

	//isqs.InitializeSqs()
	//idynamodb.Initilize()

	/*const (
		QueueUrl    = "https://sqs.us-east-1.amazonaws.com/073283464250/k3s-develop-lifecycle-hook"
	)*/

	/*for {

		var event = processors.GetNextEventFromSqs()

		if event != nil {

			// if the event is: launching ec2 instance
			if event.State == variables.EC2_LIFECYCLE_LAUNCHING {

				// then let's put internal dns of that ec2 instance into dynamodb
				if idynamodb.Client.PutEntry(models.InstanceIdToInternalDnsMapping{
					InstanceId:  event.InstanceId,
					InternalDns: event.InternalDns,
				},

					database.INSTANCE_ID_TO_INTERNAL_DNS_TABLE) {

					// and if we succeed - just remove the event from SQS
					event.DeleteEventFromSqs()

				}
			}

			// otherwise if the event is about terminating of some ec2 instance
			if event.State == variables.EC2_LIFECYCLE_TERMINATING {

				// let's pull it's internal dns from the dynamodb
				mappingInterface := idynamodb.Client.GetEntryByKey(database.INSTANCE_ID_TO_INTERNAL_DNS_TABLE, "instance-id", event.InstanceId)
				if mappingInterface != nil {
					var mapping = mappingInterface.(models.InstanceIdToInternalDnsMapping)

					// delete that k8s node from the cluster
					// and remove the corresponding entry from DynamoDB
					if ik8s.DeleteNodeByName(mapping.InternalDns) && idynamodb.Client.DeleteEntryByKey(database.INSTANCE_ID_TO_INTERNAL_DNS_TABLE, "instance-id", event.InstanceId) {
						// if we succeeded - we are fine to remove that event from SQS
						event.DeleteEventFromSqs()
					}
				}

			}

		}

		// check SQS for new events every 5 seconds
		time.Sleep(5 * time.Second)

	}*/

}




