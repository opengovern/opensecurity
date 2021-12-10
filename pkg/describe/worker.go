package describe

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api/auth/kubernetes"
	"gitlab.com/keibiengine/keibi-engine/pkg/internal/queue"
	"gitlab.com/keibiengine/keibi-engine/pkg/internal/vault"
	"gopkg.in/Shopify/sarama.v1"
)

type Worker struct {
	id             string
	jobQueue       queue.Interface
	jobResultQueue queue.Interface
	kfkProducer    sarama.SyncProducer
	kfkTopic       string
	vault          vault.Keibi
}

func InitializeWorker(
	id string,
	rabbitMQUsername string,
	rabbitMQPassword string,
	rabbitMQHost string,
	rabbitMQPort int,
	describeJobQueue string,
	describeJobResultQueue string,
	kafkaBrokers []string,
	kafkaTopic string,
	vaultAddress string,
	vaultRoleName string,
	vaultToken string,
) (w *Worker, err error) {
	if id == "" {
		return nil, fmt.Errorf("'id' must be set to a non empty string")
	} else if kafkaTopic == "" {
		return nil, fmt.Errorf("'kfkTopic' must be set to a non empty string")
	}

	w = &Worker{id: id, kfkTopic: kafkaTopic}
	defer func() {
		if err != nil && w != nil {
			w.Stop()
		}
	}()

	qCfg := queue.Config{}
	qCfg.Server.Username = rabbitMQUsername
	qCfg.Server.Password = rabbitMQPassword
	qCfg.Server.Host = rabbitMQHost
	qCfg.Server.Port = rabbitMQPort
	qCfg.Queue.Name = describeJobQueue
	qCfg.Queue.Durable = true
	qCfg.Consumer.ID = w.id
	describeQueue, err := queue.New(qCfg)
	if err != nil {
		return nil, err
	}

	w.jobQueue = describeQueue

	qCfg = queue.Config{}
	qCfg.Server.Username = rabbitMQUsername
	qCfg.Server.Password = rabbitMQPassword
	qCfg.Server.Host = rabbitMQHost
	qCfg.Server.Port = rabbitMQPort
	qCfg.Queue.Name = describeJobResultQueue
	qCfg.Queue.Durable = true
	qCfg.Producer.ID = w.id
	describeResultsQueue, err := queue.New(qCfg)
	if err != nil {
		return nil, err
	}

	w.jobResultQueue = describeResultsQueue

	producer, err := newKafkaProducer(strings.Split(KafkaService, ","))
	if err != nil {
		return nil, err
	}

	w.kfkProducer = producer

	k8sAuth, err := kubernetes.NewKubernetesAuth(
		vaultRoleName,
		kubernetes.WithServiceAccountToken(vaultToken),
	)
	if err != nil {
		return nil, err
	}

	// setup vault
	v, err := vault.NewVault(vaultAddress, k8sAuth)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to vault:", vaultAddress)
	w.vault = v

	return w, nil
}

func (w *Worker) Run() error {
	msgs, err := w.jobQueue.Consume()
	if err != nil {
		return err
	}

	fmt.Printf("Waiting indefinitly for messages. To exit press CTRL+C")
	for msg := range msgs {
		var job Job
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			fmt.Printf("Failed to unmarshal task: %s", err.Error())
			msg.Nack(false, false)
			continue
		}

		result := job.Do(w.vault, w.kfkProducer, w.kfkTopic)

		err := w.jobResultQueue.Publish(result)
		if err != nil {
			fmt.Printf("Failed to send results to queue: %s", err.Error())
		}

		msg.Ack(false)
	}

	return fmt.Errorf("descibe jobs channel is closed")
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

	if w.kfkProducer != nil {
		w.kfkProducer.Close()
		w.kfkProducer = nil
	}
}

func newKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Retry.Max = 3
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Version = sarama.V2_1_0_0

	producer, err := sarama.NewSyncProducer(strings.Split(KafkaService, ","), cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
