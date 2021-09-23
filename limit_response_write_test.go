package iorate

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResponseWrite(t *testing.T) {
	assert := require.New(t)
	recorder := httptest.NewRecorder()
	now := time.Now()
	_, err := io.Copy(recorder, bytes.NewReader(make([]byte, 5)))
	assert.NoError(err)
	assert.True(time.Since(now) < time.Second)

	w := NewResponseWrite(recorder, 2)
	now = time.Now()
	_, err = io.Copy(w, bytes.NewReader(make([]byte, 6)))
	assert.NoError(err)
	assert.True(time.Since(now) > time.Second)
}
