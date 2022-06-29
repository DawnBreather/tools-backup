package models

type InstanceIdToInternalDnsMapping struct{
	InstanceId string `json:"instance-id"`
	InternalDns string `json:"internal-dns"`
}

type InstanceIdToIntenralDnsMappings []InstanceIdToInternalDnsMapping