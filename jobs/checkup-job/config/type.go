package config

type WorkerConfig struct {
	DoTelemetry          bool   `yaml:"do_telemetry"`
	TelemetryWorkspaceID string `yaml:"telemetry_workspace_id"`
	TelemetryHostname    string `yaml:"telemetry_hostname"`
	TelemetryBaseURL     string `yaml:"telemetry_base_url"`
}
