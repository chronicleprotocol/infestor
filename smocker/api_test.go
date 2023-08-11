package smocker

import (
	"context"
	"os"
	"testing"
)

const defaultURI = "http://localhost:8081"

func TestResetAPICall(t *testing.T) {
	uri := defaultURI
	if os.Getenv("SMOCKER_HOST") != "" {
		uri = os.Getenv("SMOCKER_HOST")
	}
	api := NewAPI(uri)
	err := api.Reset(context.Background())
	if err != nil {
		t.Errorf("Error resetting API call: %s", err)
	}
}
