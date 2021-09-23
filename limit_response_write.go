package iorate

import (
	"io"
	"net/http"

	"golang.org/x/time/rate"
)

// NewResponseWrite 封装ResponseWriter为下载限速模式，rate为限速值，单位byte
//
// 例子 w = Newiorate(w, 100*1024) 限制客户端下载速度为100kb
func NewResponseWrite(w http.ResponseWriter, r int) http.ResponseWriter {
	limit := NewLimiter(r)
	return &iorate{w, NewWriter(w, limit)}
}

// NewResponseWriteByLimit 类似 NewResponseWrite
// 可传入limiter，便于多个ResponseWrite共用一个限速器，达到全局限速
func NewResponseWriteByLimit(w http.ResponseWriter, limit *rate.Limiter) http.ResponseWriter {
	return &iorate{w, NewWriter(w, limit)}
}

type iorate struct {
	http.ResponseWriter
	io.Writer
}

func (w *iorate) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}
