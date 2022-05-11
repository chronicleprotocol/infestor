package smocker

import (
	"encoding/json"
	"net/http"
	"time"
)

// OriginMock struct that represents a mock request to smocker API.
// - Reset: Optional (defaults to false), used to reset on Smocker before adding mocks.
// - Session: Optional, the name of the new session to start.
// - Body: Required, the yaml content of the mock.
type OriginMock struct {
	Reset   bool
	Session string
	Mocks   []*Mock
}

func (om *OriginMock) Body() ([]byte, error) {
	return json.Marshal(om.Mocks)
}

type Mock struct {
	Request  MockRequest   `json:"request,omitempty" yaml:"request"`
	Response *MockResponse `json:"response,omitempty" yaml:"response,omitempty"`
}

func (bm *Mock) Validate() error {
	if bm.Response == nil {
		return nil
	}
	// Clean up body in case of non Success Stats code
	if bm.Response.Status >= http.StatusBadRequest {
		bm.Response.Body = ""
	}
	return nil
}

type StringMatcher struct {
	Matcher string `json:"matcher" yaml:"matcher,flow"`
	Value   string `json:"value" yaml:"value,flow"`
}

type Delay struct {
	Min time.Duration `json:"min,omitempty" yaml:"min,omitempty"`
	Max time.Duration `json:"max,omitempty" yaml:"max,omitempty"`
}

type BodyMatcher struct {
	BodyString *StringMatcher
	BodyJSON   map[string]StringMatcher
}

func (bm BodyMatcher) MarshalJSON() ([]byte, error) {
	if bm.BodyString != nil {
		return json.Marshal(bm.BodyString)
	}
	return json.Marshal(bm.BodyJSON)
}

type MultiMapMatcher map[string]StringMatcherSlice

type StringMatcherSlice []StringMatcher

type StringSlice []string

type MapStringSlice map[string]StringSlice

type MockRequest struct {
	Path        StringMatcher   `json:"path" yaml:"path"`
	Method      StringMatcher   `json:"method" yaml:"method"`
	Body        *BodyMatcher    `json:"body,omitempty" yaml:"body,omitempty"`
	QueryParams MultiMapMatcher `json:"query_params,omitempty" yaml:"query_params,omitempty"`
	Headers     MultiMapMatcher `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type MockResponse struct {
	Body    string         `json:"body,omitempty" yaml:"body,omitempty"`
	Status  int            `json:"status" yaml:"status"`
	Delay   Delay          `json:"delay,omitempty" yaml:"delay,omitempty"`
	Headers MapStringSlice `json:"headers,omitempty" yaml:"headers,omitempty"`
}

// NewStringMatcher helps to create matchers easy
func NewStringMatcher(value string) StringMatcher {
	return StringMatcher{
		Matcher: "ShouldEqual",
		Value:   value,
	}
}

func NewSubstringMatcher(value string) StringMatcher {
	return StringMatcher{
		Matcher: "ShouldContainSubstring",
		Value:   value,
	}
}
