package smocker

import (
	"context"
	"testing"
)

func TestResetAPICall(t *testing.T) {
	api := NewAPI("http://localhost:8081")
	err := api.Reset(context.Background())
	if err != nil {
		t.Errorf("Error resetting API call: %s", err)
	}
}
