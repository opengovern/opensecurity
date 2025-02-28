package worker

type WorkloadType string

const (
	WorkloadTypeDeployment WorkloadType = "deployment"
)

type NatsConfig struct {
	Stream         string `yaml:"Stream"`
	Topic          string `yaml:"Topic"`
	Consumer       string `yaml:"Consumer"`
	ResultTopic    string `yaml:"ResultTopic"`
	ResultConsumer string `yaml:"ResultConsumer"`
}

type ScaleConfig struct {
	Stream       string `yaml:"Stream"`
	Consumer     string `yaml:"Consumer"`
	LagThreshold string `yaml:"LagThreshold"`
	MinReplica   int32  `yaml:"MinReplica"`
	MaxReplica   int32  `yaml:"MaxReplica"`

	PollingInterval int32 `yaml:"PollingInterval"`
	CooldownPeriod  int32 `yaml:"CooldownPeriod"`
}

type Task struct {
	ID           string            `yaml:"ID"`
	Name         string            `yaml:"Name"`
	Description  string            `yaml:"Description"`
	ImageURL     string            `yaml:"ImageURL"`
	Command      string            `yaml:"Command"`
	ResultType   string            `yaml:"ResultType"`
	WorkloadType WorkloadType      `yaml:"WorkloadType"`
	EnvVars      map[string]string `yaml:"EnvVars"`
	Interval     uint64            `yaml:"Interval"`
	Timeout      uint64            `yaml:"Timeout"`
	NatsConfig   NatsConfig        `yaml:"NatsConfig"`
	ScaleConfig  ScaleConfig       `yaml:"ScaleConfig"`
}
