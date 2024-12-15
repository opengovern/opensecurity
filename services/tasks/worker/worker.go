package worker

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	"github.com/opengovern/opencomply/services/tasks/config"
	"github.com/opengovern/opencomply/services/tasks/worker/consts"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateWorker(ctx context.Context, cfg config.Config, kubeClient client.Client, taskConfig *Task, namespace string) error {
	soNatsUrl, _ := os.LookupEnv("SCALED_OBJECT_NATS_URL")

	var env []corev1.EnvVar
	for k, v := range taskConfig.EnvVars {
		env = append(env, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	env = append(env, []corev1.EnvVar{
		{
			Name:  consts.NatsURLEnv,
			Value: cfg.NATS.URL,
		},
		{
			Name:  consts.NatsConsumerEnv,
			Value: taskConfig.NatsConfig.Consumer,
		},
		{
			Name:  consts.NatsStreamNameEnv,
			Value: taskConfig.NatsConfig.Stream,
		},
		{
			Name:  consts.NatsTopicNameEnv,
			Value: taskConfig.NatsConfig.Topic,
		},
		{
			Name:  consts.NatsResultTopicNameEnv,
			Value: taskConfig.NatsConfig.ResultTopic,
		},
	}...)
	switch taskConfig.WorkloadType {
	case WorkloadTypeJob:
		// job
		var job v1.Job
		err := kubeClient.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      taskConfig.ID,
		}, &job)
		if err != nil {
			job = v1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      taskConfig.ID,
					Namespace: namespace,
					Labels: map[string]string{
						"app": taskConfig.ID,
					},
				},
				Spec: v1.JobSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": taskConfig.ID,
							},
						},
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyNever,
							Containers: []corev1.Container{
								{
									Name:  taskConfig.ID,
									Image: taskConfig.ImageURL,
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
			err := kubeClient.Create(ctx, &job)
			if err != nil {
				return err
			}
		}

		// scaled-job
		var scaledObject kedav1alpha1.ScaledJob
		err = kubeClient.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      taskConfig.ID + "-scaled-job",
		}, &scaledObject)
		if err != nil {
			scaledObject = kedav1alpha1.ScaledJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      taskConfig.ID + "-scaled-job",
					Namespace: namespace,
				},
				Spec: kedav1alpha1.ScaledJobSpec{
					JobTargetRef: &v1.JobSpec{
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
								RestartPolicy: corev1.RestartPolicyNever,
								Containers: []corev1.Container{
									{
										Name:  taskConfig.ID,
										Image: taskConfig.ImageURL,
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
					PollingInterval: aws.Int32(30),
					MinReplicaCount: aws.Int32(taskConfig.ScaleConfig.MinReplica),
					MaxReplicaCount: aws.Int32(taskConfig.ScaleConfig.MaxReplica),
					Triggers: []kedav1alpha1.ScaleTriggers{
						{
							Type: "nats-jetstream",
							Metadata: map[string]string{
								"account":                      "$G",
								"natsServerMonitoringEndpoint": soNatsUrl,
								"stream":                       taskConfig.ScaleConfig.Stream,
								"consumer":                     taskConfig.ScaleConfig.Consumer,
								"lagThreshold":                 taskConfig.ScaleConfig.LagThreshold,
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
	case WorkloadTypeDeployment:
		// deployment
		var deployment appsv1.Deployment
		err := kubeClient.Get(ctx, client.ObjectKey{
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
									Image: taskConfig.ImageURL,
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
						APIVersion: appsv1.SchemeGroupVersion.Version,
					},
					PollingInterval: aws.Int32(30),
					CooldownPeriod:  aws.Int32(300),
					MinReplicaCount: aws.Int32(taskConfig.ScaleConfig.MinReplica),
					MaxReplicaCount: aws.Int32(taskConfig.ScaleConfig.MaxReplica),
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
								"stream":                       taskConfig.ScaleConfig.Stream,
								"consumer":                     taskConfig.ScaleConfig.Consumer,
								"lagThreshold":                 taskConfig.ScaleConfig.LagThreshold,
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
	default:
		return fmt.Errorf("invalid workload type: %s", taskConfig.WorkloadType)
	}

	return nil
}
