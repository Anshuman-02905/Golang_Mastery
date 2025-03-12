package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	path, _ := os.Getwd()
	fmt.Println(path)
	ConnectDatabase("../database/Library.db")
	assert.NotNil(t, DB, "Database connection should not be nil")
	sqlDB, err := DB.DB()
	assert.NoError(t, err, "Should be able to access the underlying database")
	assert.NoError(t, sqlDB.Ping(), "Database should be reachable")

}
