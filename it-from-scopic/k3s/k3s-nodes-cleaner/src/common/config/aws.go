package config

type aws struct{
	Auth awsAuth
	DynamoDb awsDynamoDb
	//Ec2 awsEc2
}

type awsAuth struct {
	AwsAccessKeyId     string
	AwsSecretAccessKey string
	AwsDefaultRegion   string
}

type awsDynamoDb struct {
	Table string
}

/*type awsEc2 struct {

}*/