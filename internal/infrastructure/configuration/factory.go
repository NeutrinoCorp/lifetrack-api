package configuration

// NewConfiguration reads and returns a kernel configuration
func NewConfiguration() Configuration {
	stage := "dev"
	return Configuration{
		Version: "0.1.0-alpha",
		Stage:   stage,
		DynamoTable: dynamoTable{
			Name:   "lifetrack-" + stage,
			Region: "us-east-1",
		},
		Cassandra: cassandra{
			Keyspace: "lifetrack_" + stage,
			Cluster:  []string{"127.0.0.1"},
		},
	}
}
