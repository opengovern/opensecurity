package steampipe

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/keibiengine/keibi-engine/pkg/internal/queue"
)

type Worker struct {
	id             string
	jobQueue       queue.Interface
	jobResultQueue queue.Interface
	s3Client       *s3.S3
	config         Config
}

func InitializeWorker(
	id string,
	config Config,
	steampipeJobQueue, steampipeJobResultQueue string,
) (w *Worker, err error) {
	if id == "" {
		return nil, fmt.Errorf("'id' must be set to a non empty string")
	}

	w = &Worker{id: id}
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
	qCfg.Queue.Name = steampipeJobQueue
	qCfg.Queue.Durable = true
	qCfg.Consumer.ID = w.id
	describeQueue, err := queue.New(qCfg)
	if err != nil {
		return nil, err
	}

	w.jobQueue = describeQueue

	qCfg = queue.Config{}
	qCfg.Server.Username = config.RabbitMQ.Username
	qCfg.Server.Password = config.RabbitMQ.Password
	qCfg.Server.Host = config.RabbitMQ.Host
	qCfg.Server.Port = config.RabbitMQ.Port
	qCfg.Queue.Name = steampipeJobResultQueue
	qCfg.Queue.Durable = true
	qCfg.Producer.ID = w.id
	describeResultsQueue, err := queue.New(qCfg)
	if err != nil {
		return nil, err
	}

	w.jobResultQueue = describeResultsQueue

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(config.S3Client.Key, config.S3Client.Secret, ""),
		Endpoint:    aws.String(config.S3Client.Endpoint),
		Region:      aws.String(config.S3Client.Region),
	}

	newSession := session.New(s3Config)
	w.s3Client = s3.New(newSession)
	w.config = config

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

		result := job.Do(w.s3Client, w.config)

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
}
