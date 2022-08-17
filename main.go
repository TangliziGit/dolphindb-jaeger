package main

import (
	"context"
	"fmt"
	"os"

	"github.com/TangliziGit/dolphindb-jaeger/pkg/jaeger"
	"github.com/TangliziGit/dolphindb-jaeger/pkg/spans"
)

func main() {
	if err := sendTrace(os.Args[1], os.Args[2]); err != nil {
		fmt.Println(err.Error())
	}
}

func sendTrace(path string, host string) error {
	spanMap, err := spans.BuildSpanMap(path)
	if err != nil {
		return err
	}

	batches := spans.BuildSpanBatch(spanMap)
	uploader := jaeger.NewUploader("http://" + host + ":14268/api/traces")
	for _, batches := range batches {
		err := uploader.Upload(context.Background(), batches)
		if err != nil {
			return err
		}
	}

	return err
}
