export interface RunDetail {
    es_id:              string;
    es_index:           string;
    controls:           Controls;
    compliance_summary: ComplianceSummary;
    job_summary:        JobSummary;
}

export interface ComplianceSummary {
    alarm: number;
}

export interface Controls {
    [key: string]: AzureUseByokForStorageAccountEncryption
}

export interface AzureUseByokForStorageAccountEncryption {
    severity: string;
    alarms:   number;
    oks:      number;
}

export interface JobSummary {
    job_id:          number;
    auditable:       boolean;
    framework_id:    string;
    job_started_at:  Date;
    integration_ids: null;
}

export interface ControlFilterResult {
    controls: ControlsResult
    integrations: null
    audit_summary: Summary
    job_summary: JobSummaryFilter
}

export interface Summary {
    alarm: number;
}

export interface ControlsResult {
    [key: string]: Results
}

export interface Results {
    severity: string
    control_summary: Summary
    results: Results
}

export interface Results {
    alarm: Alarm[];
}

export interface Alarm {
    resource_id:   string;
    resource_type: string;
    reason:        string;
}

export interface JobSummaryFilter {
    job_id:          number;
    auditable:       boolean;
    framework_id:    string;
    job_started_at:  Date;
    integration_ids: string[];
}
