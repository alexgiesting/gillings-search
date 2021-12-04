package paths

const (
	ENV_SERVER_PORT            = "GS_SERVER_PORT"         //
	ENV_QUERY_PORT             = "GS_QUERY_PORT"          // local query
	ENV_UPDATE_PORT            = "GS_UPDATE_PORT"         // local       update
	ENV_MONGODB_HOST           = "MONGODB_SERVICE_HOST"   // local              TODO
	ENV_MONGODB_PORT           = "MONGODB_SERVICE_PORT"   // local              TODO
	ENV_MONGODB_USERNAME       = "MONGODB_USER"           //
	ENV_MONGODB_PASSWORD       = "MONGODB_PASSWORD"       //                    TODO
	ENV_MONGODB_ADMIN_PASSWORD = "MONGODB_ADMIN_PASSWORD" // local query update TODO
	ENV_MONGODB_NAME           = "MONGODB_DATABASE"       // local query update
	ENV_SOLR_HOST              = "SOLR_SERVICE_HOST"      // local       update TODO
	ENV_SOLR_PORT              = "SOLR_SERVICE_PORT"      // local       update TODO

	SECRET_DIR                   = "/run/gs-secrets/"
	SECRET_UPDATE_KEY            = SECRET_DIR + "update.key"     // local       update
	SECRET_SCOPUS_API_KEY        = SECRET_DIR + "scopus.key"     // local       update
	SECRET_SCOPUS_CLIENT_ADDRESS = SECRET_DIR + "subscriber.key" // local

	PATH_QUERY  = "/query/"
	PATH_UPDATE = "/update/"
)
