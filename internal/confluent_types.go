package internal

type Environment struct {
	Id   string `json:"id"`   // id of environment
	Name string `json:"name"` // name of environment
}

type SchemaRegistry struct {
	Id         string `json:"id"`
	EnpointURL string `json:"endpoint_url"`
}

type FullSchemaRegistry struct {
	Name                string `json:"name"`
	ClusterId           string `json:"cluster_id"`
	EndpointURL         string `json:"endpoint_url"`
	UsedSchemas         string `json:"used_schemas"`
	AvailableSchemas    string `json:"available_schemas"`
	GlobalCompatibility string `json:"global_compatibility"`
	Mode                string `json:"mode"`
	ServiceProvider     string `json:"service_provider"`
}

type ClusterKafka struct {
	Id           string `json:"id"`            //"id": "lkc-57qx6n",
	Name         string `json:"name"`          //"name": "demo-kafka",
	Type         string `json:"type"`          //"type": "BASIC",
	Ingress      int16  `json:"ingress"`       //"ingress": 100,
	Egress       int16  `json:"egress"`        //"egress": 100,
	Storage      string `json:"storage"`       //"storage": "5 TB",
	Provider     string `json:"provider"`      //"provider": "aws",
	Region       string `json:"region"`        //"region": "us-east-2",
	Availability string `json:"availability"`  //"availability": "single-zone",
	Status       string `json:"status"`        //"status": "UP",
	Endpoint     string `json:"endpoint"`      //"endpoint": "SASL_SSL://pkc-ymrq7.us-east-2.aws.confluent.cloud:9092",
	RestEndpoint string `json:"rest_endpoint"` //"rest_endpoint": "https://pkc-ymrq7.us-east-2.aws.confluent.cloud:443"
}

type ApiKey struct {
	Api    string `json:"key"`
	Secret string `json:"secret"`
}

/*
 *{
 *  "id": "sa-xmvjm1",
 *  "name": "CadastralServiceAccount",
 *  "description": "This is a text"
 *}
 */
type ServiceAccount struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
