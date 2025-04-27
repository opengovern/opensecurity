package worker

type WorkloadType string

type NatsConfig struct {
	Stream         string `json:"stream" yaml:"stream"`
	Topic          string `json:"topic" yaml:"topic"`
	Consumer       string `json:"consumer" yaml:"consumer"`
	ResultTopic    string `json:"result_topic" yaml:"result_topic"`
	ResultConsumer string `json:"result_consumer" yaml:"result_consumer"`
}

type ScaleConfig struct {
	Stream       string `json:"stream" yaml:"stream"`
	Consumer     string `json:"consumer" yaml:"consumer"`
	LagThreshold string `json:"lag_threshold" yaml:"lag_threshold"`
	MinReplica   int32  `json:"min_replica" yaml:"min_replica"`
	MaxReplica   int32  `json:"max_replica" yaml:"max_replica"`

	PollingInterval int32 `json:"polling_interval" yaml:"polling_interval"`
	CooldownPeriod  int32 `json:"cooldown_period" yaml:"cooldown_period"`
}

type Interval struct {
	Months  int32 `yaml:"months,omitempty"`
	Days    int32 `yaml:"days,omitempty"`
	Hours   int32 `yaml:"hours,omitempty"`
	Minutes int32 `yaml:"minutes,omitempty"`
}

type TaskRunSchedule struct {
	ID        string         `yaml:"id"`
	Params    map[string]any `yaml:"params"`
	Frequency string         `yaml:"frequency"`
}

type Task struct {
	Type                string            `json:"type" yaml:"type"`
	ID                  string            `json:"id" yaml:"id"`
	Name                string            `json:"name" yaml:"name"`
	Description         string            `json:"description" yaml:"description"`
	IsEnabled           bool              `json:"is_enabled" yaml:"is_enabled"`
	ImageURL            string            `json:"image_url" yaml:"image_url"`
	ArtifactsURL        string            `json:"artifacts_url" yaml:"artifacts_url"`
	SteampipePluginName string            `json:"steampipe_plugin_name" yaml:"steampipe_plugin_name"`
	Command             string            `json:"command" yaml:"command"`
	Timeout             string            `json:"timeout" yaml:"timeout"`
	NatsConfig          NatsConfig        `json:"nats_config" yaml:"nats_config"`
	ScaleConfig         ScaleConfig       `json:"scale_config" yaml:"scale_config"`
	RunSchedule         []TaskRunSchedule `json:"run_schedule" yaml:"run_schedule"`
	Params              []string          `json:"params" yaml:"params"`
	Configs             []string          `json:"configs" yaml:"configs"`
}
