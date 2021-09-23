package iorate

import (
	"bytes"
	"context"
	"io"

	"golang.org/x/time/rate"
)

// Write 限速的写流
type Write struct {
	w io.Writer
	b *rate.Limiter
}

// Write 实现writer接口
func (lw *Write) Write(data []byte) (int, error) {
	b := lw.b.Burst()
	// 如果data大于桶的大小，则分成多次写入，每次写入长度和桶大小相等
	if len(data) > b {
		buff := make([]byte, b)
		// MultiReader用于隐藏bytes.Reader的WriteTo接口
		n, err := io.CopyBuffer(lw, io.MultiReader(bytes.NewReader(data)), buff)
		return int(n), err
	}
	n, err := lw.w.Write(data)
	if err != nil {
		// 如果发生错误则不消耗令牌数量，可能会导致轻微的限流不准确
		return n, err
	}
	return n, lw.b.WaitN(context.Background(), n)
}

// NewWriter 初始化一个限速的写流
func NewWriter(w io.Writer, b *rate.Limiter) *Write {
	return &Write{w: w, b: b}
}

// NewLimiter 初始化一个限速令牌桶
func NewLimiter(r int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(r), r)
}
