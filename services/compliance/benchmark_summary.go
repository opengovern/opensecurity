package compliance

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (h *HttpHandler) GetBenchmarkTreeIDs(ctx context.Context, rootID string) ([]string, error) {
	ids := []string{rootID}
	// tracer :
	_, span2 := tracer.Start(ctx, "new_GetBenchmark", trace.WithSpanKind(trace.SpanKindServer))
	span2.SetName("new_GetBenchmark")
	h.logger.Info(fmt.Sprintf("RootId: %s", rootID))
	root, err := h.db.GetFramework(ctx, rootID)
	if err != nil {
		return nil, err
	}
	h.logger.Info(fmt.Sprintf("Root: %v", root))
	span2.AddEvent("information", trace.WithAttributes(
		attribute.String("benchmark ID", root.ID),
	))
	span2.End()

	for _, child := range root.Children {
		h.logger.Info(fmt.Sprintf("Child: %v", child))
		cids, err := h.GetBenchmarkTreeIDs(ctx, child.ID)
		if err != nil {
			return nil, err
		}

		ids = append(ids, cids...)
	}
	return ids, nil
}
