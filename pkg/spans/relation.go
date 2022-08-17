package spans

import (
	"dolphindb-jaeger/pkg/uuid"
	"fmt"

	gen "dolphindb-jaeger/pkg/jaeger/gen-go/jaeger"
)

func NewRelation(tid *uuid.UUID, relType string, parentID string) (ref *gen.SpanRef, err error) {
	parentUUID, err := uuid.NewUUID(parentID)
	if err != nil {
		return nil, err
	}

	ref = &gen.SpanRef{
		RefType:     0,
		TraceIdLow:  tid.Low,
		TraceIdHigh: tid.High,
		SpanId:      parentUUID.Squash(),
	}

	switch relType {
	case "Root":
		ref.RefType = gen.SpanRefType_CHILD_OF
		ref.SpanId = 0
	case "ChildOf":
		ref.RefType = gen.SpanRefType_CHILD_OF
	case "FollowsFrom":
		ref.RefType = gen.SpanRefType_FOLLOWS_FROM
	default:
		err = fmt.Errorf("no such relation type: %s", relType)
	}

	return
}
