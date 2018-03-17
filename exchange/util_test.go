package exchange

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetSupportedExchanges(t *testing.T) {
	exchanges := GetSupportedExchanges()
	assert.True(t, len(exchanges) > 0, "no supported exchanges were returned")
}