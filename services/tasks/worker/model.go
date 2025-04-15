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
	Params    map[string]any `yaml:"params"`
	Frequency string         `yaml:"frequency"`
}

type Task struct {
	ID                  string            `yaml:"id"`
	Name                string            `yaml:"name"`
	Description         string            `yaml:"description"`
	IsEnabled           bool              `yaml:"is_enabled"`
	ImageURL            string            `yaml:"image_url"`
	ArtifactsURL        string            `yaml:"artifacts_url"`
	SteampipePluginName string            `yaml:"steampipe_plugin_name"`
	Command             string            `yaml:"command"`
	Timeout             string            `yaml:"timeout"`
	NatsConfig          NatsConfig        `yaml:"nats_config"`
	ScaleConfig         ScaleConfig       `yaml:"scale_config"`
	RunSchedule         []TaskRunSchedule `yaml:"run_schedule"`
}
