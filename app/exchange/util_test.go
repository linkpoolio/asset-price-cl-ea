package exchange

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSupportedExchanges(t *testing.T) {
	exchanges := GetSupportedExchanges()
	assert.True(t, len(exchanges) > 0, "no supported exchanges were returned")
}
