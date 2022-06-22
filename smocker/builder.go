package smocker

import (
	"net/http"
	"time"
)

type MockBuilder struct {
	mock *Mock
}

func NewMockBuilder() *MockBuilder {
	return &MockBuilder{mock: &Mock{
		Request: MockRequest{
			Method: ShouldMatch(".*"),
			Path:   ShouldMatch(".*"),
		},
		Response: &MockResponse{
			Status: http.StatusOK,
		},
	}}
}

func (mb *MockBuilder) SetRequestPath(path StringMatcher) *MockBuilder {
	mb.mock.Request.Path = path
	return mb
}

func (mb *MockBuilder) SetRequestMethod(method StringMatcher) *MockBuilder {
	mb.mock.Request.Method = method
	return mb
}

func (mb *MockBuilder) SetRequestBodyString(body StringMatcher) *MockBuilder {
	if mb.mock.Request.Body == nil {
		mb.mock.Request.Body = &BodyMatcher{}
	}
	mb.mock.Request.Body.BodyString = &body
	return mb
}

func (mb *MockBuilder) AddRequestBodyJSON(key string, body StringMatcher) *MockBuilder {
	if mb.mock.Request.Body == nil {
		mb.mock.Request.Body = &BodyMatcher{}
	}
	if mb.mock.Request.Body.BodyJSON == nil {
		mb.mock.Request.Body.BodyJSON = make(map[string]StringMatcher)
	}
	mb.mock.Request.Body.BodyJSON[key] = body
	return mb
}

func (mb *MockBuilder) SetRequestQueryParams(queryParams MultiMapMatcher) *MockBuilder {
	if mb.mock.Request.QueryParams == nil {
		mb.mock.Request.QueryParams = make(MultiMapMatcher)
	}
	mb.mock.Request.QueryParams = queryParams
	return mb
}

func (mb *MockBuilder) AddRequestQueryParam(name string, value StringMatcher) *MockBuilder {
	if mb.mock.Request.QueryParams == nil {
		mb.mock.Request.QueryParams = make(MultiMapMatcher)
	}
	mb.mock.Request.QueryParams[name] = append(mb.mock.Request.QueryParams[name], value)
	return mb
}

func (mb *MockBuilder) SetRequestHeaders(headers MultiMapMatcher) *MockBuilder {
	if mb.mock.Request.Headers == nil {
		mb.mock.Request.Headers = make(MultiMapMatcher)
	}
	mb.mock.Request.Headers = headers
	return mb
}

func (mb *MockBuilder) AddRequestHeader(name string, value StringMatcher) *MockBuilder {
	if mb.mock.Request.Headers == nil {
		mb.mock.Request.Headers = make(MultiMapMatcher)
	}
	mb.mock.Request.Headers[name] = append(mb.mock.Request.Headers[name], value)
	return mb
}

func (mb *MockBuilder) SetResponseStatus(status int) *MockBuilder {
	mb.mock.Response.Status = status
	return mb
}

func (mb *MockBuilder) SetResponseHeaders(headers MapStringSlice) *MockBuilder {
	if mb.mock.Response.Headers == nil {
		mb.mock.Response.Headers = make(MapStringSlice)
	}
	mb.mock.Response.Headers = headers
	return mb
}

func (mb *MockBuilder) AddResponseHeader(name string, value string) *MockBuilder {
	if mb.mock.Response.Headers == nil {
		mb.mock.Response.Headers = make(MapStringSlice)
	}
	mb.mock.Response.Headers[name] = append(mb.mock.Response.Headers[name], value)
	return mb
}

func (mb *MockBuilder) SetResponseDelay(min, max time.Duration) *MockBuilder {
	mb.mock.Response.Delay = Delay{Min: min, Max: max}
	return mb
}

func (mb *MockBuilder) SetResponseBody(body string) *MockBuilder {
	mb.mock.Response.Body = body
	return mb
}

func (mb *MockBuilder) SetContextTimes(times uint) *MockBuilder {
	if mb.mock.Context == nil {
		mb.mock.Context = &MockContext{}
	}
	mb.mock.Context.Times = times
	return mb
}

func (mb *MockBuilder) Mock() *Mock {
	return mb.mock
}
