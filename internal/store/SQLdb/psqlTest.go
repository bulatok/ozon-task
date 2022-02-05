package SQLdb

import (
	"database/sql"
	"github.com/bulatok/ozon-task/configs"
)

// CreateTEST is needed to test PostgreSQL database only
func CreateTEST(config *configs.Config) *PostgreDB {
	return &PostgreDB{
		DB : &sql.DB{},
		DbURL: config.DatabaseURL,
	}
}