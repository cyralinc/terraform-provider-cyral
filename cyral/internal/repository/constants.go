package repository

const (
	resourceName   = "cyral_repository"
	dataSourceName = "cyral_repository"
)

const (
	// Schema keys.
	RepoIDKey     = "id"
	RepoTypeKey   = "type"
	RepoNameKey   = "name"
	RepoLabelsKey = "labels"
	// Connection draining keys.
	RepoConnDrainingKey         = "connection_draining"
	RepoConnDrainingAutoKey     = "auto"
	RepoConnDrainingWaitTimeKey = "wait_time"
	// Repo node keys.
	RepoNodesKey       = "repo_node"
	RepoHostKey        = "host"
	RepoPortKey        = "port"
	RepoNodeDynamicKey = "dynamic"
	// MongoDB settings keys.
	RepoMongoDBSettingsKey       = "mongodb_settings"
	RepoMongoDBReplicaSetNameKey = "replica_set_name"
	RepoMongoDBServerTypeKey     = "server_type"
	RepoMongoDBSRVRecordName     = "srv_record_name"
	RepoMongoDBFlavorKey         = "flavor"
)

const (
	Denodo          = "denodo"
	Dremio          = "dremio"
	DynamoDB        = "dynamodb"
	DynamoDBStreams = "dynamodbstreams"
	Galera          = "galera"
	MariaDB         = "mariadb"
	MongoDB         = "mongodb"
	MySQL           = "mysql"
	Oracle          = "oracle"
	PostgreSQL      = "postgresql"
	Redshift        = "redshift"
	S3              = "s3"
	Snowflake       = "snowflake"
	SQLServer       = "sqlserver"
)

func RepositoryTypes() []string {
	return []string{
		Denodo,
		Dremio,
		DynamoDB,
		DynamoDBStreams,
		Galera,
		MariaDB,
		MongoDB,
		MySQL,
		Oracle,
		PostgreSQL,
		Redshift,
		S3,
		Snowflake,
		SQLServer,
	}
}

const (
	ReplicaSet = "replicaset"
	Standalone = "standalone"
	Sharded    = "sharded"
)

const (
	MongoDBFlavorMongoDB    = "mongodb"
	MongoDBFlavorDocumentDB = "documentdb"
)

func mongoServerTypes() []string {
	return []string{
		ReplicaSet,
		Standalone,
		Sharded,
	}
}

func mongoFlavors() []string {
	return []string{
		MongoDBFlavorMongoDB,
		MongoDBFlavorDocumentDB,
	}
}
