# io-rate

## limit io write speed

```golang
limit := NewLimiter(512)
// w write speed to 512B/s
w := NewWriter(ioutil.Discard, limit)
```

## limit http response write speed

```golang
func download(w http.ResponseWriter, r *http.Request)
    // rw write speed to 1MB/s
    rw = iorate.NewResponseWrite(rw, 1024*1024)
    ...
}
```

## limit http server global speed

```golang
// 1M/s speed shared by all clients
limit := NewLimiter(1024*1024)
func download(w http.ResponseWriter, r *http.Request)
    rw = iorate.NewResponseWriteByLimit(rw, limit)
    ...
}
```
