package infestor

import (
	"bytes"
	"context"

	"github.com/chronicleprotocol/infestor/origin"

	"github.com/chronicleprotocol/infestor/smocker"
)

type MocksBuilder struct {
	reset     bool
	exchanges []origin.Exchange
}

func NewMocksBuilder() *MocksBuilder {
	return &MocksBuilder{}
}

// Reset clears all previously created in smocker
func (mb *MocksBuilder) Reset() *MocksBuilder {
	mb.reset = true
	return mb
}

// Add adds new exchange/pair mock to smoker
func (mb *MocksBuilder) Add(e *origin.Exchange) *MocksBuilder {
	mb.exchanges = append(mb.exchanges, *e)
	return mb
}

func (mb *MocksBuilder) Deploy(api smocker.API) error {
	var yaml bytes.Buffer
	for _, e := range mb.exchanges {
		part, err := origin.BuildMock(e)
		if err != nil {
			return err
		}
		yaml.Write(part)
		yaml.WriteString("\n")
	}
	mock := smocker.Mock{
		Body: yaml.Bytes(),
	}
	return api.AddMock(context.Background(), mock)
}
