package functions

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudsql "github.com/GoogleCloudPlatform/golang-samples/cloudsql/postgres/database-sql"
)

// TABS vs. SPACES App for Cloud Functions
func init() {
	functions.HTTP("Votes", cloudsql.Votes)
}
