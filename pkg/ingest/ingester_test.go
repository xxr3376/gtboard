package ingest

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIngest(t *testing.T) {
	path := os.Getenv("INGEST_PATH")
	if path == "" {
		t.Skip("INGEST_PATH not set")
	}
	ingester, err := NewIngester("test", path)
	require.NoError(t, err)
	defer ingester.Close()

	in, err := ingester.FetchUpdates(context.Background())
	require.NoError(t, err)
	require.NotZero(t, in)

	run := ingester.GetRun()
	require.NotNil(t, run)
	t.Logf("run contains %d scalars", len(run.Scalars))
	t.Logf("run contains %d texts", len(run.Texts))
}
