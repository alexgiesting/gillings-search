package paths

const (
	ENV_SERVER_PORT            = "GS_SERVER_PORT"           //
	ENV_QUERY_PORT             = "GS_QUERY_PORT"            // local query
	ENV_UPDATE_PORT            = "GS_UPDATE_PORT"           // local       update
	ENV_UPDATE_KEY             = "GS_UPDATE_KEY"            // local       update
	ENV_MONGODB_HOST           = "MONGODB_SERVICE_HOST"     // local              TODO
	ENV_MONGODB_PORT           = "MONGODB_SERVICE_PORT"     // local              TODO
	ENV_MONGODB_USERNAME       = "MONGODB_USER"             //
	ENV_MONGODB_PASSWORD       = "MONGODB_PASSWORD"         //
	ENV_MONGODB_ADMIN_PASSWORD = "MONGODB_ADMIN_PASSWORD"   // local query update
	ENV_MONGODB_NAME           = "MONGODB_DATABASE"         // local query update
	ENV_SCOPUS_API_KEY         = "GS_SCOPUS_API_KEY"        // local       update
	ENV_SCOPUS_CLIENT_ADDRESS  = "GS_SCOPUS_CLIENT_ADDRESS" // local

	PATH_QUERY  = "/query/"
	PATH_UPDATE = "/update/"
)
