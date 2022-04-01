package compliance_report

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api/auth/kubernetes"
	"gitlab.com/keibiengine/keibi-engine/pkg/internal/queue"
	"gitlab.com/keibiengine/keibi-engine/pkg/internal/vault"
	"gopkg.in/Shopify/sarama.v1"
	"strings"
)

type Worker struct {
	id             string
	jobQueue       queue.Interface
	jobResultQueue queue.Interface
	config         Config
	vault          vault.SourceConfig
	kfkProducer    sarama.SyncProducer
	kfkTopic       string
}

func InitializeWorker(
	id string,
	config Config,
	complianceReportJobQueue, complianceReportJobResultQueue string,
) (w *Worker, err error) {
	if id == "" {
		return nil, fmt.Errorf("'id' must be set to a non empty string")
	}

	w = &Worker{id: id, kfkTopic: config.Kafka.Topic}
	defer func() {
		if err != nil && w != nil {
			w.Stop()
		}
	}()

	qCfg := queue.Config{}
	qCfg.Server.Username = config.RabbitMQ.Username
	qCfg.Server.Password = config.RabbitMQ.Password
	qCfg.Server.Host = config.RabbitMQ.Host
	qCfg.Server.Port = config.RabbitMQ.Port
	qCfg.Queue.Name = complianceReportJobQueue
	qCfg.Queue.Durable = true
	qCfg.Consumer.ID = w.id
	reportJobQueue, err := queue.New(qCfg)
	if err != nil {
		return nil, err
	}

	w.jobQueue = reportJobQueue

	qCfg = queue.Config{}
	qCfg.Server.Username = config.RabbitMQ.Username
	qCfg.Server.Password = config.RabbitMQ.Password
	qCfg.Server.Host = config.RabbitMQ.Host
	qCfg.Server.Port = config.RabbitMQ.Port
	qCfg.Queue.Name = complianceReportJobResultQueue
	qCfg.Queue.Durable = true
	qCfg.Producer.ID = w.id
	reportResultQueue, err := queue.New(qCfg)
	if err != nil {
		return nil, err
	}

	w.jobResultQueue = reportResultQueue

	producer, err := newKafkaProducer(strings.Split(config.Kafka.Addresses, ","))
	if err != nil {
		return nil, err
	}
	w.kfkProducer = producer
	w.config = config

	k8sAuth, err := kubernetes.NewKubernetesAuth(
		config.Vault.RoleName,
		kubernetes.WithServiceAccountToken(config.Vault.Token),
	)
	if err != nil {
		return nil, err
	}

	// setup vault
	v, err := vault.NewSourceConfig(config.Vault.Address, k8sAuth)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to vault:", config.Vault.Address)
	w.vault = v

	return w, nil
}

func (w *Worker) Run() error {
	msgs, err := w.jobQueue.Consume()
	if err != nil {
		return err
	}

	fmt.Printf("Waiting indefinitly for messages. To exit press CTRL+C\n")
	for msg := range msgs {
		var job Job
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			fmt.Printf("Failed to unmarshal task: %s", err.Error())
			msg.Nack(false, false)
			continue
		}

		result := job.Do(w.vault, w.kfkProducer, w.kfkTopic, w.config)

		err := w.jobResultQueue.Publish(result)
		if err != nil {
			fmt.Printf("Failed to send results to queue: %s", err.Error())
		}

		fmt.Printf("A job is done and result is published into the result queue, result: %v\n", result)
		msg.Ack(false)
	}

	return fmt.Errorf("report jobs channel is closed")
}

func (w *Worker) Stop() {
	if w.jobQueue != nil {
		w.jobQueue.Close()
		w.jobQueue = nil
	}

	if w.jobResultQueue != nil {
		w.jobResultQueue.Close()
		w.jobResultQueue = nil
	}
}
func newKafkaProducer(kafkaServers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Retry.Max = 3
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Version = sarama.V2_1_0_0

	producer, err := sarama.NewSyncProducer(kafkaServers, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
