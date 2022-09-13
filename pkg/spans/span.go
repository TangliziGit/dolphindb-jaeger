package spans

import (
	"bufio"
	"fmt"
	gen "github.com/TangliziGit/dolphindb-jaeger/pkg/jaeger/report/gen-go/jaeger"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/TangliziGit/dolphindb-jaeger/pkg/uuid"
)

var filenameRegex = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}-merged.log$")

func NewSpan(tid *uuid.UUID, logType string, tokens []string) (span *gen.Span, err error) {
	spanId, err := uuid.NewUUID(tokens[0])
	if err != nil {
		return nil, err
	}

	timestamp, err := strconv.ParseInt(tokens[1], 10, 64)
	if err != nil {
		return nil, err
	}

	switch logType {
	case "Start":
		ref, err := NewRelation(tid, tokens[3], tokens[4])
		if err != nil {
			return nil, err
		}

		span = &gen.Span{
			TraceIdLow:    tid.Low,
			TraceIdHigh:   tid.High,
			SpanId:        spanId.Squash(),
			ParentSpanId:  ref.SpanId,
			OperationName: tokens[5],
			References:    []*gen.SpanRef{ref},
			StartTime:     timestamp / 1000,
			Duration:      0,
			Tags:          NewTags(tokens),
			Logs:          []*gen.Log{},
		}
	case "End":
		span = &gen.Span{
			TraceIdLow:  tid.Low,
			TraceIdHigh: tid.High,
			SpanId:      spanId.Squash(),
			StartTime:   timestamp / 1000,
		}
	case "Done":
		span = nil
	default:
		err = fmt.Errorf("no such log type: %s", logType)
	}

	return span, err
}

func BuildSpanMap(path string) (spanMap map[int64]*gen.Span, tid *uuid.UUID, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	filename := filepath.Base(path)
	if filenameRegex.MatchString(filename) == false {
		return nil, nil, fmt.Errorf("invalid filename: %s", filename)
	}

	tid, err = uuid.NewUUID(filename[:36])
	if err != nil {
		return nil, nil, err
	}

	spanMap = map[int64]*gen.Span{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), ",")
		logType := tokens[2]

		span, err := NewSpan(tid, logType, tokens)
		if err != nil {
			return nil, nil, err
		}

		switch logType {
		case "Start":
			spanMap[span.SpanId] = span
		case "End":
			originSpan := spanMap[span.SpanId]
			originSpan.Duration = span.StartTime - originSpan.StartTime
		case "Done":
			// do nothing
		default:
			err = fmt.Errorf("no such log type: %s", logType)
		}
	}

	return
}

func BuildSpanBatch(spanMap map[int64]*gen.Span) map[string]*gen.Batch {
	batches := map[string]*gen.Batch{}

	getNode := func(span *gen.Span) string {
		for _, tag := range span.GetTags() {
			if tag.Key == "node" {
				return *tag.VStr
			}
		}
		return ""
	}

	for _, span := range spanMap {
		node := getNode(span)
		if _, ok := batches[node]; !ok {
			batches[node] = &gen.Batch{
				Spans: []*gen.Span{},
				Process: &gen.Process{
					ServiceName: node,
				},
			}
		}

		batches[node].Spans = append(batches[node].Spans, span)
	}

	return batches
}
