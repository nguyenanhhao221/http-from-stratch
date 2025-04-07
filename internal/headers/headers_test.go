package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra white space
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069                           \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 57, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n a bunch of other stuff")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Empty(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Uppercase field-name should output lowercase
	headers = NewHeaders()
	data = []byte("HoSt: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	_, ok := headers["host"]
	assert.True(t, ok)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: field-name with number still work
	headers = NewHeaders()
	data = []byte("H1st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid field name length
	headers = NewHeaders()
	data = []byte(": localhost:42069\r\n\r\n")
	_, _, err = headers.Parse(data)
	require.Error(t, err)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: valid headers that match exist header field name
	// Set-Person: lane-loves-go;
	// Set-Person: prime-loves-zig;
	// Set-Person: tj-loves-ocaml;
	// This is valid and should become set-person: lane-loves-go, prime-loves-zig, tj-loves-ocaml
	headers = Headers{"set-person": "lane-loves-go"}
	data = []byte("Set-Person: prime-loves-zig\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	// assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Equal(t, "lane-loves-go, prime-loves-zig", headers["set-person"])
}

func TestHeaders_Get(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		key   string
		want  string
		want2 bool
	}{
		{name: "valid get", key: "hello", want: "world", want2: true},
		{name: "valid get with uppercase key", key: "HELLO", want: "world", want2: true},
		{name: "get non exist key", key: "non-exist", want: "", want2: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHeaders()
			h.Set("hello", "world")
			got, got2 := h.Get(tt.key)
			if tt.want != got {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
			if tt.want2 != got2 {
				t.Errorf("Get() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
