package smocker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	URL string
}

func NewAPI(url string) *API {
	return &API{URL: url}
}

// Reset Clear the mocks and the history of calls.
func (a *API) Reset(ctx context.Context) error {
	request := Request{
		Method:  Post,
		BaseURL: fmt.Sprintf("%s/reset?force=true", a.URL),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	res, err := SendWithContext(ctx, request)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("reset failed: %s", res.Body)
	}
	return nil
}

// AddMocks Add a mock to the mocks list.
func (a *API) AddMocks(ctx context.Context, mock []*Mock) error {
	body, err := json.Marshal(mock)
	if err != nil {
		return fmt.Errorf("failed to marshal mocks request: %w", err)
	}
	request := Request{
		Method:  Post,
		BaseURL: fmt.Sprintf("%s/mocks", a.URL),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}
	res, err := SendWithContext(ctx, request)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add mocks: %s", res.Body)
	}
	return nil
}
