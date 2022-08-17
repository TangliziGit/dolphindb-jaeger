package jaeger

import (
	"bytes"
	"context"
	gen "dolphindb-jaeger/pkg/jaeger/gen-go/jaeger"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"io"
	"net/http"
)

// code copied from
//`https://github.com/open-telemetry/opentelemetry-go/blob/main/exporters/jaeger/uploader.go`
type Uploader struct {
	endpoint   string
	httpClient *http.Client
}

func NewUploader(endpoint string) *Uploader {
	return &Uploader{
		endpoint:   endpoint,
		httpClient: http.DefaultClient,
	}
}

func (c *Uploader) Upload(ctx context.Context, batch *gen.Batch) error {
	body, err := serialize(batch)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-thrift")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	if err = resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to upload traces; HTTP status code: %d", resp.StatusCode)
	}
	return nil
}

func serialize(obj thrift.TStruct) (*bytes.Buffer, error) {
	buf := thrift.NewTMemoryBuffer()
	if err := obj.Write(context.Background(), thrift.NewTBinaryProtocolConf(buf, &thrift.TConfiguration{})); err != nil {
		return nil, err
	}
	return buf.Buffer, nil
}
