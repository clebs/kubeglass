package fetch

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseList(t *testing.T) {
	data, err := os.ReadFile("../../test/testdata/discovery-dir-ls.json")
	require.NoError(t, err)

	out := files{}
	require.NoError(t, json.Unmarshal(data, &out))
	require.Len(t, out, 54)
	require.Equal(t, "aggregated_v2beta1.json", out[0].Name)
}
