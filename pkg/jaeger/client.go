package jaeger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gen "github.com/TangliziGit/dolphindb-jaeger/pkg/jaeger/report/gen-go/jaeger"
	"github.com/TangliziGit/dolphindb-jaeger/pkg/uuid"
	"github.com/apache/thrift/lib/go/thrift"
	ui "github.com/jaegertracing/jaeger/model/json"
	"io"
	"net/http"
)

// Why not use gRPC Trace retrieval API v2?
// Here is the problem occurred: https://github.com/jaegertracing/jaeger/issues/3662

type Client struct {
	thriftEndpoint string
	httpEndpoint   string
	client         *http.Client
}

func NewClient(host string, thriftPort string, httpPort string) *Client {
	return &Client{
		thriftEndpoint: fmt.Sprintf("http://%s:%s/api/traces", host, thriftPort),
		httpEndpoint:   fmt.Sprintf("http://%s:%s/api/traces", host, httpPort),
		client:         http.DefaultClient,
	}
}

func (c *Client) IsTraceExists(tid *uuid.UUID) (bool, error) {
	url := fmt.Sprintf("%s/%s", c.httpEndpoint, tid.HexString())
	r, err := c.client.Get(url)
	if err != nil {
		return false, err
	}

	var resp TraceResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return len(resp.Data) > 0, nil
}

func (c *Client) Upload(ctx context.Context, batch *gen.Batch) error {
	body, err := serialize(batch)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.thriftEndpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-thrift")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	if err = resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to report traces; HTTP status code: %d", resp.StatusCode)
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

type TraceResponse struct {
	Data   []*ui.Trace  `json:"data"`
	Total  int          `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
	Errors []TraceError `json:"errors"`
}

type TraceError struct {
	Code    int        `json:"code,omitempty"`
	Msg     string     `json:"msg"`
	TraceID ui.TraceID `json:"traceID,omitempty"`
}
