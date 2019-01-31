// +build integration

package db

import (
	"testing"

	"github.com/thomaspoignant/user-microservice/testingUtils"

	"github.com/stretchr/testify/assert"
)

// TestGettingDynamoDBTable_emptyName test if we are throwing
//  an error when we ask for an emtpy table
func TestGettingDynamoDBTable_emptyName(t *testing.T) {
	testingUtils.PrepareTest()
	table, err := GetDynamodbTable("")
	assert.NotNil(t, err)
	assert.Nil(t, table)
}

func TestGettingDynamoDBTable_notEmptyName(t *testing.T) {
	testingUtils.PrepareTest()
	table, err := GetDynamodbTable("test")
	assert.Nil(t, err)
	assert.NotNil(t, table)
}
