package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/jq"
	essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/utils"
	"go.uber.org/zap"
	"time"
)

const (
	StreamName     = "es-sink"
	SinkQueueTopic = "es-sink-queue"
	ConsumerGroup  = "es-sink-consumer"
)

type EsSinkService struct {
	logger        *zap.Logger
	elasticSearch essdk.Client
	nats          *jq.JobQueue
	esSinkModule  *EsSinkModule
}

func NewEsSinkService(ctx context.Context, logger *zap.Logger, elasticSearch essdk.Client, nats *jq.JobQueue) (*EsSinkService, error) {
	service := EsSinkService{
		logger:        logger,
		elasticSearch: elasticSearch,
		nats:          nats,
	}

	esSinkModule, err := NewEsSinkModule(ctx, logger, elasticSearch)
	if err != nil {
		logger.Error("failed to create es sink module", zap.Error(err))
		return nil, err
	}
	service.esSinkModule = esSinkModule

	err = service.nats.StreamWithConfig(ctx, StreamName, "es sink stream", []string{SinkQueueTopic}, jetstream.StreamConfig{
		//Name:                 "",
		//Description:          "",
		//Subjects:             nil,
		Retention:    jetstream.WorkQueuePolicy,
		MaxConsumers: -1,
		MaxMsgs:      15000000,
		MaxBytes:     1024 * 15000000,
		Discard:      jetstream.DiscardNew,
		MaxAge:       time.Hour * 48,
		MaxMsgSize:   50 * 1024 * 1024,
		Storage:      jetstream.FileStorage,
		Replicas:     1,
		Duplicates:   100 * time.Millisecond,
		Compression:  jetstream.S2Compression,
		ConsumerLimits: jetstream.StreamConsumerLimits{
			MaxAckPending: 1000,
		},
	})
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *EsSinkService) Start(ctx context.Context) {
	s.logger.Info("starting es sink service")
	defer s.logger.Info("es sink service stopped")

	utils.EnsureRunGoroutine(func() {
		s.ConsumeCycle(ctx)
	})

	s.esSinkModule.Start(ctx)
}

func (s *EsSinkService) ConsumeCycle(ctx context.Context) {
	consumeCtx, err := s.nats.Consume(ctx, "es-sink", StreamName, []string{SinkQueueTopic}, ConsumerGroup, func(msg jetstream.Msg) {
		var doc es.DocBase
		err := json.Unmarshal(msg.Data(), &doc)
		if err != nil {
			s.logger.Error("failed to unmarshal doc", zap.Error(err), zap.Any("msg", msg))
			return
		}

		s.esSinkModule.QueueDoc(doc)

		err = msg.Ack()
		if err != nil {
			s.logger.Error("failed to ack message", zap.Error(err), zap.Any("msg", msg))
		}
	})
	if err != nil {
		s.logger.Fatal("failed to consume", zap.Error(err))
	}

	s.logger.Info("consuming", zap.String("stream", StreamName), zap.String("topic", SinkQueueTopic))

	<-ctx.Done()
	consumeCtx.Drain()
	consumeCtx.Stop()
}

type FailedDoc struct {
	Doc es.DocBase `json:"doc"`
	Err string     `json:"err"`
}

func (s *EsSinkService) Ingest(ctx context.Context, docs []es.DocBase) ([]FailedDoc, error) {
	failedDocs := make([]FailedDoc, 0)
	for _, doc := range docs {
		id, idx := doc.GetIdAndIndex()
		docJson, err := json.Marshal(doc)
		if err != nil {
			s.logger.Error("failed to marshal doc", zap.Error(err))
			failedDocs = append(failedDocs, FailedDoc{
				Doc: doc,
				Err: err.Error(),
			})
		}
		fullId := fmt.Sprintf(fmt.Sprintf("%s:::%s:::%s", idx, id, uuid.New().String()))
		h := sha256.New()
		h.Write([]byte(fullId))
		hashedId := fmt.Sprintf("%x", h.Sum(nil))
		_, err = s.nats.Produce(ctx, SinkQueueTopic, docJson, fmt.Sprintf("%s:::%s", uuid.New().String(), hashedId))
		if err != nil {
			s.logger.Error("failed to produce message", zap.Error(err), zap.Any("doc", doc))
			failedDocs = append(failedDocs, FailedDoc{
				Doc: doc,
				Err: err.Error(),
			})
		}
	}

	return failedDocs, nil
}
