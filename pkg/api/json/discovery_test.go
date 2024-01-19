package json

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToAPIMap(t *testing.T) {
	data, err := os.ReadFile("../../../test/testdata/k8s-api-versions-1.28.json")
	require.NoError(t, err)

	m, err := ToAPIMap(data)
	require.NoError(t, err)
	require.Len(t, m.Groups, 21)

	data, err = os.ReadFile("../../../test/testdata/k8s-api-versions-1.29.json")
	require.NoError(t, err)

	m, err = ToAPIMap(data)
	require.NoError(t, err)
	require.Len(t, m.Groups, 21)

}
