package smocker

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type API struct {
	Host string
	Port int
}

func NewAPI(host string, port int) *API {
	return &API{
		Host: host,
		Port: port,
	}
}

func (a *API) basePath() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

// Reset Clear the mocks and the history of calls.
func (a *API) Reset(ctx context.Context) error {
	request := Request{
		Method:  Post,
		BaseURL: fmt.Sprintf("%s/reset", a.basePath()),
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

// AddMock Add a mock to the mocks list.
func (a *API) AddMock(ctx context.Context, mock Mock) error {
	request := Request{
		Method:  Post,
		BaseURL: fmt.Sprintf("%s/mocks", a.basePath()),
		Headers: map[string]string{
			"Content-Type": "application/x-yaml",
		},
		Body: mock.Body,
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
