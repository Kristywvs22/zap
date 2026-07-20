package zap

import (
    "bytes"
    "io"
    "net/http"
)

// NewRequest creates an http.Request with proper handling for retries.
func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
    // Create the initial request
    req, err := http.NewRequest(method, url, body)
    if err!= nil {
        return nil, err
    }

    // Check if the body is already seekable/rewindable
    if req.GetBody!= nil {
        return req, nil
    }

    // Buffer the body if it's not seekable
    var buf bytes.Buffer
    if _, err := io.Copy(&buf, body); err!= nil {
        return nil, err
    }

    // Replace the body with a bytes.Reader
    req.Body = io.NopCloser(&buf)
    req.ContentLength = int64(buf.Len())
    req.GetBody = func() (io.ReadCloser, error) {
        return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
    }

    return req, nil
}