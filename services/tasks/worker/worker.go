package worker

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgtype"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"github.com/opengovern/opensecurity/services/tasks/config"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"github.com/opengovern/opensecurity/services/tasks/worker/consts"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

var (
	ESAddress  = os.Getenv("ELASTICSEARCH_ADDRESS")
	ESUsername = os.Getenv("ELASTICSEARCH_USERNAME")
	ESPassword = os.Getenv("ELASTICSEARCH_PASSWORD")
	ESIsOnAks  = os.Getenv("ELASTICSEARCH_ISONAKS")

	InventoryBaseURL = os.Getenv("CORE_BASEURL")
)

func CreateWorker(ctx context.Context, cfg config.Config, kubeClient client.Client, taskConfig *models.Task, namespace string) error {
	soNatsUrl, _ := os.LookupEnv("SCALED_OBJECT_NATS_URL")

	env, err := defaultEnvs(cfg, taskConfig)
	if err != nil {
		return err
	}

	var deployment appsv1.Deployment
	err = kubeClient.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      taskConfig.ID,
	}, &deployment)
	if err != nil {
		deployment = appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      taskConfig.ID,
				Namespace: namespace,
				Labels: map[string]string{
					"app": taskConfig.ID,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: aws.Int32(0),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": taskConfig.ID,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": taskConfig.ID,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  taskConfig.ID,
								Image: taskConfig.ImageUrl,
								Command: []string{
									taskConfig.Command,
								},
								ImagePullPolicy: corev1.PullAlways,
								Env:             env,
							},
						},
					},
				},
			},
		}
		err := kubeClient.Create(ctx, &deployment)
		if err != nil {
			return err
		}
	}

	var scaleConfig ScaleConfig
	if taskConfig.ScaleConfig.Status == pgtype.Present {
		if err := json.Unmarshal(taskConfig.ScaleConfig.Bytes, &scaleConfig); err != nil {
			return err
		}
	}

	// scaled-object
	var scaledObject kedav1alpha1.ScaledObject
	err = kubeClient.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      taskConfig.ID + "-scaled-object",
	}, &scaledObject)
	if err != nil {
		scaledObject = kedav1alpha1.ScaledObject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      taskConfig.ID + "-scaled-object",
				Namespace: namespace,
			},
			Spec: kedav1alpha1.ScaledObjectSpec{
				ScaleTargetRef: &kedav1alpha1.ScaleTarget{
					Name:       taskConfig.ID,
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				PollingInterval: aws.Int32(scaleConfig.PollingInterval),
				CooldownPeriod:  aws.Int32(scaleConfig.CooldownPeriod),
				MinReplicaCount: aws.Int32(scaleConfig.MinReplica),
				MaxReplicaCount: aws.Int32(scaleConfig.MaxReplica),
				Fallback: &kedav1alpha1.Fallback{
					FailureThreshold: 1,
					Replicas:         1,
				},
				Triggers: []kedav1alpha1.ScaleTriggers{
					{
						Type: "nats-jetstream",
						Metadata: map[string]string{
							"account":                      "$G",
							"natsServerMonitoringEndpoint": soNatsUrl,
							"stream":                       scaleConfig.Stream,
							"consumer":                     scaleConfig.Consumer + "-service",
							"lagThreshold":                 scaleConfig.LagThreshold,
							"useHttps":                     "false",
						},
					},
				},
			},
		}
		err = kubeClient.Create(ctx, &scaledObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func defaultEnvs(cfg config.Config, taskConfig *models.Task) ([]corev1.EnvVar, error) {
	var natsConfig NatsConfig
	if taskConfig.NatsConfig.Status == pgtype.Present {
		if err := json.Unmarshal(taskConfig.ScaleConfig.Bytes, &natsConfig); err != nil {
			return nil, err
		}
	}

	return []corev1.EnvVar{
		{
			Name:  consts.NatsURLEnv,
			Value: cfg.NATS.URL,
		},
		{
			Name:  consts.NatsConsumerEnv,
			Value: natsConfig.Consumer,
		},
		{
			Name:  consts.NatsStreamNameEnv,
			Value: natsConfig.Stream,
		},
		{
			Name:  consts.NatsTopicNameEnv,
			Value: natsConfig.Topic,
		},
		{
			Name:  consts.NatsResultTopicNameEnv,
			Value: natsConfig.ResultTopic,
		},
		{
			Name:  consts.ElasticSearchAddressEnv,
			Value: ESAddress,
		},
		{
			Name:  consts.ElasticSearchUsernameEnv,
			Value: ESUsername,
		},
		{
			Name:  consts.ElasticSearchPasswordEnv,
			Value: ESPassword,
		},
		{
			Name:  consts.ElasticSearchIsOnAksNameEnv,
			Value: ESIsOnAks,
		},
		{
			Name:  consts.ElasticSearchIsOpenSearch,
			Value: strconv.FormatBool(cfg.ElasticSearch.IsOpenSearch),
		},
		{
			Name:  consts.ElasticSearchAwsRegionEnv,
			Value: "",
		},
		{
			Name:  consts.ElasticSearchAssumeRoleArnEnv,
			Value: "",
		},
		{
			Name:  consts.InventoryBaseURL,
			Value: InventoryBaseURL,
		},
	}, nil
}
