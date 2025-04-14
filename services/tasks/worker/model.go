package worker

type WorkloadType string

type NatsConfig struct {
	Stream         string `yaml:"stream"`
	Topic          string `yaml:"topic"`
	Consumer       string `yaml:"consumer"`
	ResultTopic    string `yaml:"result_topic"`
	ResultConsumer string `yaml:"result_consumer"`
}

type ScaleConfig struct {
	Stream       string `yaml:"stream"`
	Consumer     string `yaml:"consumer"`
	LagThreshold string `yaml:"lag_threshold"`
	MinReplica   int32  `yaml:"min_replica"`
	MaxReplica   int32  `yaml:"max_replica"`

	PollingInterval int32 `yaml:"polling_interval"`
	CooldownPeriod  int32 `yaml:"cooldown_period"`
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
	ID          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	IsEnabled   bool              `yaml:"is_enabled"`
	ImageURL    string            `yaml:"image_url"`
	Command     string            `yaml:"command"`
	Timeout     string            `yaml:"timeout"`
	NatsConfig  NatsConfig        `yaml:"nats_config"`
	ScaleConfig ScaleConfig       `yaml:"scale_config"`
	RunSchedule []TaskRunSchedule `yaml:"run_schedule"`
}
