export interface Integration {
    integration_id:   string;
    provider_id:      string;
    name:             string;
    integration_type: string;
    annotations:      Annotations;
    labels:           Annotations;
    credential_id:    string;
    state:            string;
    last_check:       Date;
}

export interface Annotations {
}

export interface Credentials {
    id:               string;
    secret:           string;
    integration_type: string;
    metadata:         Metadata;
    masked_secret:    MaskedSecret;
    integration_count: number;
    description:      string;
    created_at:       Date;
    updated_at:       Date;
}

export interface MaskedSecret {
     [key: string]: string[]

}
export interface Metadata {
}

export interface Schema {
    integration_type_id:    string;
    integration_name:       string;
    help_text_md:           string;
    platform_documentation: string;
    provider_documentation: string;
    icon:                   string;
    discover:               Discover;
    render:                 Render;
    actions:                Actions;
}

export interface Actions {
    credentials:  ActionsCredential[];
    integrations: ActionsIntegration[];
}

export interface ActionsCredential {
    type: string
    label: string
    editableFields?: string[]
    confirm?: CredentialConfirm
    tooltip?: string
}

export interface CredentialConfirm {
    message:   string;
    condition: Condition;
}

export interface Condition {
    field:        string;
    operator:     string;
    value:        number;
    errorMessage: string;
}

export interface ActionsIntegration {
    type: string
    label: string
    confirm?: IntegrationConfirm
    editableFields?: string[]
    tooltip?: string
}

export interface IntegrationConfirm {
    message: string;
}

export interface Discover {
    credentials:  DiscoverCredential[];
    integrations: DiscoverIntegration[];
}

export interface DiscoverCredential {
    type:     string;
    label:    string;
    priority: number;
    fields:   CredentialField[];
}

export interface CredentialField {
    name:               string;
    label:              string;
    inputType:          string;
    required:           boolean;
    order:              number;
    validation:         Validation;
    info:               string;
    external_help_url?: string;
}

export interface Validation {
    pattern?:       string;
    errorMessage:   string;
    fileTypes?:     string[];
    maxFileSizeMB?: number;
}

export interface DiscoverIntegration {
    label:  string;
    type:   string;
    fields: IntegrationField[];
}

export interface IntegrationField {
    name:           string;
    label:          string;
    fieldType:      string;
    required:       boolean;
    order:          number;
    info:           string;
    valueMap?:      ValueMap;
    statusOptions?: StatusOption[];
}

export interface StatusOption {
    value: string;
    label: string;
    color: string;
}

export interface ValueMap {
    spn_password_based: string;
    spn_certificate:    string;
}

export interface Render {
    credentials:  Credentials;
    integrations: Credentials;
}

export interface Credentials {
    defaultPageSize: number;
    fields:          CredentialsField[];
}

export interface CredentialsField {
    name:           string;
    label:          string;
    fieldType:      string;
    order:          number;
    sortable?:      boolean;
    show?:          boolean;
    filterable?:    boolean;
    info:           string;
    detail:         boolean;
    detail_order:   number;
    required?:      boolean;
    valueMap?:      ValueMap;
    statusOptions?: StatusOption[];
}
