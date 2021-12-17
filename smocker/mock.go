package smocker

// Mock struct that represents a mock request to smocker API.
// - Reset: Optional (defaults to false), used to reset on Smocker before adding mocks.
// - Session: Optional, the name of the new session to start.
// - Body: Required, the yaml content of the mock.
type Mock struct {
	Reset   bool
	Session string
	Body    []byte
}
