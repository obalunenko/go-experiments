package masksensitive

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayer_Mask(t *testing.T) {
	p := &player{
		Name:     "John",
		LastName: "Doe",
		Phone:    "1234567890",
		Card:     "1234567890123456",
	}

	b, err := json.Marshal(p)
	if err != nil {
		require.NoError(t, err)
	}

	expected := `{"Name":"John","LastName":"Doe","Phone":"1234567***","Card":"********"}`

	assert.Equal(t, expected, string(b))
}
