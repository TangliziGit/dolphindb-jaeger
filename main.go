package main

import (
	"context"
	"fmt"

	"github.com/jessevdk/go-flags"

	"github.com/TangliziGit/dolphindb-jaeger/pkg/jaeger"
	"github.com/TangliziGit/dolphindb-jaeger/pkg/spans"
)

var opts struct {
	Host string `value-name:"host" short:"h" long:"host" default:"localhost" description:"The jaeger host"`
	Port string `value-name:"port" short:"p" long:"port" default:"14268" description:"The jaeger thrift port"`
	Args struct {
		LogPath string `description:"The merged trace log file"`
	} `positional-args:"true" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if flags.WroteHelp(err) {
			return
		} else {
			panic(err)
		}
	}

	if err := sendTrace(opts.Args.LogPath, opts.Host, opts.Port); err != nil {
		fmt.Println(err.Error())
	}
}

func sendTrace(path string, host string, port string) error {
	spanMap, err := spans.BuildSpanMap(path)
	if err != nil {
		return err
	}

	batches := spans.BuildSpanBatch(spanMap)
	uploader := jaeger.NewUploader(fmt.Sprintf("http://%s:%s/api/traces", host, port))
	for _, batches := range batches {
		err := uploader.Upload(context.Background(), batches)
		if err != nil {
			return err
		}
	}

	return err
}
