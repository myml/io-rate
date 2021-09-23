package iorate

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWriter(t *testing.T) {
	assert := require.New(t)
	limit := NewLimiter(1024 * 1024)
	w := NewWriter(ioutil.Discard, limit)
	now := time.Now()
	n, err := io.Copy(w, io.LimitReader(rand.Reader, 1024*1024*5))
	assert.NoError(err)
	d := int(time.Since(now) / time.Second)
	// 实际限速会有稍微的差距，但应该相差不能超过一秒
	assert.Contains([]int{d - 1, d, d + 1}, int(n)/1024/1024)
}

// 测试多个流共用一个limiter
func TestNewWriterMulti(t *testing.T) {
	assert := require.New(t)
	data := make([]byte, 1024*2)
	_, err := io.ReadFull(rand.Reader, data)
	assert.NoError(err)
	limit := NewLimiter(512)
	now := time.Now()
	for i := 0; i < 3; i++ {
		t.Run(fmt.Sprintf("client_%d", i+1), func(t *testing.T) {
			t.Parallel()
			var buff bytes.Buffer
			w := NewWriter(&buff, limit)
			_, err := io.Copy(w, bytes.NewReader(data))
			assert.NoError(err)
			assert.EqualValues(buff.Bytes(), data)
		})
	}
	t.Cleanup(func() {
		d := int(time.Since(now) / time.Second)
		// 实际限速会有稍微的差距，但应该相差不能超过一秒
		assert.Contains([]int{d - 1, d, d + 1}, len(data)/512*3)
	})
}
