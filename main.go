package main

import (
	"context"
	"dolphindb-jaeger/pkg/jaeger"
	"dolphindb-jaeger/pkg/spans"
	"fmt"
	"os"
)

func main() {
	if err := sendTrace(os.Args[1]); err != nil {
		fmt.Println(err.Error())
	}
}

func sendTrace(path string) error {
	spanMap, err := spans.BuildSpanMap(path)
	if err != nil {
		return err
	}

	batches := spans.BuildSpanBatch(spanMap)
	uploader := jaeger.NewUploader("http://localhost:14268/api/traces")
	for _, batches := range batches {
		err := uploader.Upload(context.Background(), batches)
		if err != nil {
			return err
		}
	}

	return err
}
