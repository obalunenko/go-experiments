package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonEqual(t *testing.T) {
	t.Run("test respAPI", func(t *testing.T) {
		b1, err := testRespAPI()
		require.NoError(t, err)

		b2, err := testObj()
		require.NoError(t, err)

		assert.JSONEq(t, string(b1), string(b2))
	})
}
