package infestor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chronicleprotocol/infestor/origin"
	"github.com/chronicleprotocol/infestor/smocker"
)

type MocksBuilder struct {
	debug bool
	reset bool
	mocks map[string][]origin.ExchangeMock
}

func NewMocksBuilder() *MocksBuilder {
	return &MocksBuilder{
		debug: false,
		reset: false,
		mocks: make(map[string][]origin.ExchangeMock),
	}
}

// Debug sets debug flag and mocks yaml file will be created
func (mb *MocksBuilder) Debug() *MocksBuilder {
	mb.debug = true
	return mb
}

// Reset clears all previously created in smocker
func (mb *MocksBuilder) Reset() *MocksBuilder {
	mb.reset = true
	return mb
}

// Add adds new exchange/pair mock to smoker
func (mb *MocksBuilder) Add(e *origin.ExchangeMock) *MocksBuilder {
	mocks, ok := mb.mocks[e.Name]
	if !ok {
		mb.mocks[e.Name] = []origin.ExchangeMock{*e}
		return mb
	}

	mb.mocks[e.Name] = append(mocks, *e)
	return mb
}

func (mb *MocksBuilder) Deploy(api smocker.API) error {
	ctx := context.Background()
	var result []*smocker.Mock
	for name, mocks := range mb.mocks {
		part, err := origin.BuildMocksForExchanges(name, mocks)
		if err != nil {
			return err
		}
		result = append(result, part...)
	}

	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to build mocks body: %w", err)
	}
	if mb.debug {
		err := os.WriteFile("./mocks.json", body, 0644) //nolint:gosec,gomnd
		if err != nil {
			return fmt.Errorf("failed to write debug mocks.yaml: %w", err)
		}
	}

	if mb.reset {
		err := api.Reset(ctx)
		if err != nil {
			return fmt.Errorf("failed to reset mocks before pushing new: %w", err)
		}
	}

	// If we don't have mocks available - do nothing...
	if len(mb.mocks) == 0 {
		return nil
	}

	return api.AddMocks(ctx, result)
}
