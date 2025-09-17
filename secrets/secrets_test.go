package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	// Valid names
	assert.NoError(t, ValidateName("valid_name"))
	assert.NoError(t, ValidateName("valid-name"))
	assert.NoError(t, ValidateName("validname123"))
	assert.NoError(t, ValidateName("valid/name"))
	assert.NoError(t, ValidateName("valid\\name"))

	// Invalid names
	assert.Error(t, ValidateName(""))
	assert.Error(t, ValidateName("invalid name"))
	assert.Error(t, ValidateName("invalid.name"))
}
