package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"

	"github.com/TangliziGit/dolphindb-jaeger/pkg/jaeger"
	"github.com/TangliziGit/dolphindb-jaeger/pkg/spans"
)

type Options struct {
	Host       string `value-name:"host" short:"h" long:"host" default:"localhost" description:"The jaeger host"`
	ThriftPort string `value-name:"thrift-port" long:"thrift-port" default:"14268" description:"The jaeger thrift port"`
	HttpPort   string `value-name:"http-port" long:"http-port" default:"16686" description:"The jaeger http port"`
	Args       struct {
		LogPath string `description:"The merged trace log file"`
	} `positional-args:"true" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		if flags.WroteHelp(err) {
			return
		} else {
			panic(err)
		}
	}

	if err := sendTrace(opts); err != nil {
		fmt.Println(err.Error())
	}
}

func sendTrace(opts Options) error {
	spanMap, tid, err := spans.BuildSpanMap(opts.Args.LogPath)
	if err != nil {
		return err
	}

	client := jaeger.NewClient(opts.Host, opts.ThriftPort, opts.HttpPort)

	exists, err := client.IsTraceExists(tid)
	if err != nil {
		return err
	} else if exists {
		return fmt.Errorf("this trace has already reported: traceID=%s", tid.HexString())
	}

	batches := spans.BuildSpanBatch(spanMap)
	for _, batches := range batches {
		err := client.Upload(context.Background(), batches)
		if err != nil {
			return err
		}
	}

	return err
}
