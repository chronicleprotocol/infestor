package smocker

import (
	"encoding/json"
	"net/http"
	"time"
)

type Mock struct {
	Request  MockRequest   `json:"request,omitempty" yaml:"request"`
	Response *MockResponse `json:"response,omitempty" yaml:"response,omitempty"`
	Context  *MockContext  `json:"context,omitempty" yaml:"context,omitempty"`
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

type MockContext struct {
	Times uint `json:"times" yaml:"times"`
}

func ShouldEqualJSON(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldEqualJSON", Value: value}
}

func ShouldEqual(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldEqual", Value: value}
}

func ShouldNotEqual(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldNotEqual", Value: value}
}

func ShouldResemble(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldResemble", Value: value}
}

func ShouldNotResemble(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldNotResemble", Value: value}
}

func ShouldContainSubstring(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldContainSubstring", Value: value}
}

func ShouldNotContainSubstring(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldNotContainSubstring", Value: value}
}

func ShouldStartWith(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldStartWith", Value: value}
}

func ShouldNotStartWith(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldNotStartWith", Value: value}
}

func ShouldEndWith(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldEndWith", Value: value}
}

func ShouldNotEndWith(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldNotEndWith", Value: value}
}

func ShouldMatch(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldMatch", Value: value}
}

func ShouldNotMatch(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldNotMatch", Value: value}
}

func ShouldBeEmpty(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldBeEmpty", Value: value}
}

func ShouldNotBeEmpty(value string) StringMatcher {
	return StringMatcher{Matcher: "ShouldNotBeEmpty", Value: value}
}
