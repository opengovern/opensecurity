/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export enum ApiRole {
    InternalRole = 'internal',
    KaytuAdminRole = 'kaytu-admin',
    AdminRole = 'admin',
    EditorRole = 'editor',
    ViewerRole = 'viewer',
}

export interface EchoHTTPError {
    message?: any
}

export interface EsResource {
    /** ID is the globally unique ID of the resource. */
    arn?: string
    /** Tags is the list of tags associated with the resource */
    canonical_tags?: EsTag[]
    /** CreatedAt is when the DescribeSourceJob is created */
    created_at?: number
    /** Description is the description of the resource based on the describe call. */
    description?: any
    es_id?: string
    es_index?: string
    /** ID is the globally unique ID of the resource. */
    id?: string
    /** Location is location/region of the resource */
    location?: string
    /** Metadata is arbitrary data associated with each resource */
    metadata?: Record<string, string>
    /** Name is the name of the resource. */
    name?: string
    /** ResourceGroup is the group of resource (Azure only) */
    resource_group?: string
    /** ResourceJobID is the DescribeResourceJob ID that described this resource */
    resource_job_id?: number
    /** ResourceType is the type of the resource. */
    resource_type?: string
    /** ScheduleJobID */
    schedule_job_id?: number
    /** SourceID is the Source ID that the resource belongs to */
    source_id?: string
    /** SourceJobID is the DescribeSourceJob ID that the ResourceJobID was created for */
    source_job_id?: number
    /** SourceType is the type of the source of the resource, i.e. AWS Cloud, Azure Cloud. */
    source_type?: SourceType
}

export interface EsTag {
    key?: string
    value?: string
}

export enum GithubComKaytuIoKaytuEnginePkgAnalyticsDbMetricType {
    MetricTypeAssets = 'assets',
    MetricTypeSpend = 'spend',
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiChangeUserPreferencesRequest {
    enableColorBlindMode?: boolean
    theme?: GithubComKaytuIoKaytuEnginePkgAuthApiTheme
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiCreateAPIKeyRequest {
    /** Name of the key */
    name?: string
    /**
     * Name of the role
     * @example "admin"
     */
    role?: 'admin' | 'editor' | 'viewer'
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiCreateAPIKeyResponse {
    /**
     * Activity state of the key
     * @example true
     */
    active?: boolean
    /**
     * Creation timestamp in UTC
     * @example "2023-03-31T09:36:09.855Z"
     */
    createdAt?: string
    /**
     * Unique identifier for the key
     * @example 1
     */
    id?: number
    /**
     * Name of the key
     * @example "example"
     */
    name?: string
    /**
     * Name of the role
     * @example "admin"
     */
    roleName?: 'admin' | 'editor' | 'viewer'
    /** Token of the key */
    token?: string
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiGetMeResponse {
    /**
     * Is the user blocked or not
     * @example false
     */
    blocked?: boolean
    colorBlindMode?: boolean
    /**
     * Creation timestamp in UTC
     * @example "2023-03-31T09:36:09.855Z"
     */
    createdAt?: string
    /**
     * Email address of the user
     * @example "johndoe@example.com"
     */
    email?: string
    /**
     * Is email verified or not
     * @example true
     */
    emailVerified?: boolean
    /**
     * Last activity timestamp in UTC
     * @example "2023-04-21T08:53:09.928Z"
     */
    lastActivity?: string
    lastLogin?: string
    memberSince?: string
    /**
     * Invite status
     * @example "accepted"
     */
    status?: 'accepted' | 'pending'
    theme?: GithubComKaytuIoKaytuEnginePkgAuthApiTheme
    /**
     * Unique identifier for the user
     * @example "auth|123456789"
     */
    userId?: string
    /**
     * Username
     * @example "John Doe"
     */
    userName?: string
    role?: string,
    connector_id?: string
    workspaceAccess?: Record<string, ApiRole>
    ConnectorId?: string
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiGetRoleBindingsResponse {
    /**
     * Global Access
     * @example "admin"
     */
    globalRoles?: 'admin' | 'editor' | 'viewer'
    /** List of user roles in each workspace */
    roleBindings?: GithubComKaytuIoKaytuEnginePkgAuthApiUserRoleBinding[]
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiGetUserResponse {
    /**
     * Is the user blocked or not
     * @example false
     */
    blocked?: boolean
    /**
     * Creation timestamp in UTC
     * @example "2023-03-31T09:36:09.855Z"
     */
    createdAt?: string
    /**
     * Email address of the user
     * @example "johndoe@example.com"
     */
    email?: string
    /**
     * Is email verified or not
     * @example true
     */
    emailVerified?: boolean
    /**
     * Last activity timestamp in UTC
     * @example "2023-04-21T08:53:09.928Z"
     */
    lastActivity?: string
    /**
     * Name of the role
     * @example "admin"
     */
    roleName?: 'admin' | 'editor' | 'viewer'
    /**
     * Invite status
     * @example "accepted"
     */
    status?: 'accepted' | 'pending'
    /**
     * Unique identifier for the user
     * @example "auth|123456789"
     */
    userId?: string
    /**
     * Username
     * @example "John Doe"
     */
    userName?: string
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiGetUsersRequest {
    /** @example "johndoe@example.com" */
    email?: string
    /**
     * Filter by
     * @example true
     */
    emailVerified?: boolean
    /**
     * Filter by role name
     * @example "admin"
     */
    roleName?: 'admin' | 'editor' | 'viewer'
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiGetUsersResponse {
    /**
     * Email address of the user
     * @example "johndoe@example.com"
     */
    email?: string
    /**
     * Is email verified or not
     * @example true
     */
    emailVerified?: boolean
    /**
     * Name of the role
     * @example "admin"
     */
    roleName?: 'admin' | 'editor' | 'viewer'
    /**
     * Unique identifier for the user
     * @example "auth|123456789"
     */
    userId?: string
    /**
     * Username
     * @example "John Doe"
     */
    userName?: string
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiInviteRequest {
    /**
     * User email address
     * @example "johndoe@example.com"
     */
    email_address: string
    /**
     * Name of the role
     * @example "admin"
     */
    role?: 'admin' | 'editor' | 'viewer'
    password: string
    is_active: boolean
}

export enum GithubComKaytuIoKaytuEnginePkgAuthApiInviteStatus {
    InviteStatusACCEPTED = 'accepted',
    InviteStatusPENDING = 'pending',
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiPutRoleBindingRequest {
    /** Name of the role */
    // connectionIDs?: string[]
    email_address: string

    /**
     * Name of the role
     * @example "admin"
     */
    role: 'admin' | 'editor' | 'viewer'
    /**
     * Unique identifier for the User
     * @example "auth|123456789"
     */
    // password: string
    // userId: string
    is_active: boolean
}

export enum GithubComKaytuIoKaytuEnginePkgAuthApiTheme {
    ThemeSystem = 'system',
    ThemeLight = 'light',
    ThemeDark = 'dark',
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiUserRoleBinding {
    /**
     * Name of the binding Role
     * @example "admin"
     */
    roleName?: 'admin' | 'editor' | 'viewer'
    /**
     * Unique identifier for the Workspace
     * @example "ws123456789"
     */
    workspaceID?: string
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiWorkspaceApiKey {
    /**
     * Activity state of the key
     * @example true
     */
    active?: boolean
    /**
     * Creation timestamp in UTC
     * @example "2023-03-31T09:36:09.855Z"
     */
    created_at?: string
    /**
     * Unique identifier of the user who created the key
     * @example "auth|123456789"
     */
    creator_user_id?: string
    /**
     * Unique identifier for the key
     * @example 1
     */
    id?: number
    /**
     * Masked key
     * @example "abc...de"
     */
    maskedKey?: string
    /**
     * Name of the key
     * @example "example"
     */
    name?: string
    /**
     * Name of the role
     * @example "admin"
     */
    role_name?: 'admin' | 'editor' | 'viewer'
    /**
     * Last update timestamp in UTC
     * @example "2023-04-21T08:53:09.928Z"
     */
    updated_at?: string
}

export interface GithubComKaytuIoKaytuEnginePkgAuthApiWorkspaceRoleBinding {
    /**
     * Creation timestamp in UTC
     * @example "2023-03-31T09:36:09.855Z"
     */
    created_at?: string
    /**
     * Email address of the user
     * @example "johndoe@example.com"
     */
    email?: string
    /**
     * Last activity timestamp in UTC
     * @example "2023-04-21T08:53:09.928Z"
     */
    last_activity?: string
    /**
     * Name of the role
     * @example "admin"
     */
    role_name?: 'admin' | 'editor' | 'viewer'
    scopedConnectionIDs?: string[]
    /**
     * Invite status
     * @example "accepted"
     */
    status?: 'accepted' | 'pending'
    /**
     * Unique identifier for the user
     * @example "auth|123456789"
     */
    id?: number
    /**
     * Username
     * @example "John Doe"
     */
    userName?: string
    is_active?: boolean
    connector_id?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiAccountsFindingsSummary {
    accountId?: string
    accountName?: string
    conformanceStatusesCount?: {
        error?: number
        failed?: number
        info?: number
        passed?: number
        skip?: number
    }
    lastCheckTime?: string
    securityScore?: number
    severitiesCount?: {
        critical?: number
        high?: number
        low?: number
        medium?: number
        none?: number
    }
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiAssignedBenchmark {
    benchmarkId?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmark
    /**
     * Status
     * @example true
     */
    status?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmark {
    /**
     * Whether the benchmark is auto assigned or not
     * @example true
     */
    autoAssign?: boolean
    /** Benchmark category */
    category?: string
    /**
     * Benchmark children
     * @example ["[azure_cis_v140_1"," azure_cis_v140_2]"]
     */
    children?: string[]
    /**
     * Benchmark connectors
     * @example ["[azure]"]
     */
    integrationTypes?: SourceType[]
    /**
     * Benchmark controls
     * @example ["[azure_cis_v140_1_1"," azure_cis_v140_1_2]"]
     */
    controls?: string[]
    /**
     * Benchmark creation date
     * @example "2020-01-01T00:00:00Z"
     */
    createdAt?: string
    /**
     * Benchmark description
     * @example "The CIS Microsoft Azure Foundations Security Benchmark provides prescriptive guidance for establishing a secure baseline configuration for Microsoft Azure."
     */
    description?: string
    /**
     * Benchmark document URI
     * @example "benchmarks/azure_cis_v140.md"
     */
    documentURI?: string
    /**
     * Benchmark ID
     * @example "azure_cis_v140"
     */
    id?: string
    /** Benchmark logo URI */
    logoURI?: string
    /**
     * Benchmark display code
     * @example "CIS 1.4.0"
     */
    referenceCode?: string
    /** Benchmark tags */
    tags?: Record<string, string[]>
    /**
     * Benchmark title
     * @example "Azure CIS v1.4.0"
     */
    title?: string
    /**
     * Whether the benchmark tracks drift events or not
     * @example true
     */
    tracksDriftEvents?: boolean
    /**
     * Benchmark last update date
     * @example "2020-01-01T00:00:00Z"
     */
    updatedAt?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkAssignedConnection {
    /**
     * Connection ID
     * @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"
     */
    connectionID?: string
    /**
     * Clout Provider
     * @example "Azure"
     */
    connector?: SourceType
    /**
     * Provider Connection ID
     * @example "1283192749"
     */
    providerConnectionID?: string
    /** Provider Connection Name */
    providerConnectionName?: string
    /**
     * Status
     * @example true
     */
    status?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkAssignedEntities {
    connections?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkAssignedConnection[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkAssignment {
    /** Unix timestamp */
    assignedAt?: string
    /**
     * Benchmark ID
     * @example "azure_cis_v140"
     */
    benchmarkId?: string
    /**
     * Connection ID
     * @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"
     */
    connectionId?: string
    /**
     * Resource Collection ID
     * @example "example-rc"
     */
    resourceCollectionId?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkControlSummary {
    benchmark?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmark
    children?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkControlSummary[]
    control?: GithubComKaytuIoKaytuEnginePkgComplianceApiControlSummary[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkControlsSeverityStatus {
    critical?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    high?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    low?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    medium?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    none?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    total?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkEvaluationSummary {
    /**
     * Whether the benchmark is auto assigned or not
     * @example true
     */
    autoAssign?: boolean
    /** Benchmark category */
    category?: string
    checks?: TypesSeverityResult
    /**
     * Benchmark children
     * @example ["[azure_cis_v140_1"," azure_cis_v140_2]"]
     */
    children?: string[]
    conformanceStatusSummary?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatusSummary
    connectionsStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    /**
     * Benchmark connectors
     * @example ["[azure]"]
     */
    connectors?: SourceType[]
    /**
     * Benchmark controls
     * @example ["[azure_cis_v140_1_1"," azure_cis_v140_1_2]"]
     */
    controls?: string[]
    controlsSeverityStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkControlsSeverityStatus
    costOptimization?: number
    /**
     * Benchmark creation date
     * @example "2020-01-01T00:00:00Z"
     */
    createdAt?: string
    /**
     * Benchmark description
     * @example "The CIS Microsoft Azure Foundations Security Benchmark provides prescriptive guidance for establishing a secure baseline configuration for Microsoft Azure."
     */
    description?: string
    /**
     * Benchmark document URI
     * @example "benchmarks/azure_cis_v140.md"
     */
    documentURI?: string
    /** @example "2020-01-01T00:00:00Z" */
    evaluatedAt?: string
    /**
     * Benchmark ID
     * @example "azure_cis_v140"
     */
    id?: string
    /** @example "success" */
    lastJobStatus?: string
    /** Benchmark logo URI */
    logoURI?: string
    /**
     * Benchmark display code
     * @example "CIS 1.4.0"
     */
    referenceCode?: string
    resourcesSeverityStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkResourcesSeverityStatus
    /** Benchmark tags */
    tags?: Record<string, string[]>
    /**
     * Benchmark title
     * @example "Azure CIS v1.4.0"
     */
    title?: string
    topConnections?: GithubComKaytuIoKaytuEnginePkgComplianceApiTopFieldRecord[]
    /**
     * Whether the benchmark tracks drift events or not
     * @example true
     */
    tracksDriftEvents?: boolean
    /**
     * Benchmark last update date
     * @example "2020-01-01T00:00:00Z"
     */
    updatedAt?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkRemediation {
    remediation?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkResourcesSeverityStatus {
    critical?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    high?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    low?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    medium?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    none?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
    total?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkStatusResult {
    passed?: number
    total?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkTrendDatapoint {
    checks?: TypesSeverityResult
    conformanceStatusSummary?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatusSummary
    controlsSeverityStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkControlsSeverityStatus
    /** @example "1686346668" */
    timestamp?: string
}

export enum GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus {
    ConformanceStatusFailed = 'failed',
    ConformanceStatusPassed = 'passed',
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatusSummary {
    failed?: number
    passed?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiControl {
    /** @example "To enable multi-factor authentication for a user, run the following command..." */
    cliRemediation?: string
    /** @example ["Azure"] */
    connector?: SourceType[]
    /** @example "2020-01-01T00:00:00Z" */
    createdAt?: string
    /** @example "Enable multi-factor authentication for all user credentials who have write access to Azure resources. These include roles like 'Service Co-Administrators', 'Subscription Owners', 'Contributors'." */
    description?: string
    /** @example "benchmarks/azure_cis_v140_1_1.md" */
    documentURI?: string
    /** @example true */
    enabled?: boolean
    /** @example "Multi-factor authentication adds an additional layer of security by requiring users to enter a code from a mobile device or phone in addition to their username and password when signing into Azure." */
    explanation?: string
    /** @example "To enable multi-factor authentication for a user, run the following command..." */
    guardrailRemediation?: string
    /** @example "azure_cis_v140_1_1" */
    id?: string
    /** @example true */
    managed?: boolean
    /** @example "To enable multi-factor authentication for a user, run the following command..." */
    manualRemediation?: string
    /** @example true */
    manualVerification?: boolean
    /** @example "Non-compliance to this control could result in several costs including..." */
    nonComplianceCost?: string
    /** @example "To enable multi-factor authentication for a user, run the following command..." */
    programmaticRemediation?: string
    query?: GithubComKaytuIoKaytuEnginePkgComplianceApiQuery
    /** @example "low" */
    severity?: TypesFindingSeverity
    tags?: Record<string, string[]>
    /** @example "1.1 Ensure that multi-factor authentication status is enabled for all privileged users" */
    title?: string
    /** @example "2020-01-01T00:00:00Z" */
    updatedAt?: string
    /** @example "Access to resources must be closely controlled to prevent malicious activity like data theft..." */
    usefulExample?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiControlSummary {
    benchmarks?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmark[]
    control?: GithubComKaytuIoKaytuEnginePkgComplianceApiControl
    costOptimization?: number
    evaluatedAt?: string
    failedConnectionCount?: number
    failedResourcesCount?: number
    passed?: boolean
    resourceType?: GithubComKaytuIoKaytuEnginePkgInventoryApiResourceType
    totalConnectionCount?: number
    totalResourcesCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiControlTrendDatapoint {
    failedConnectionCount?: number
    failedResourcesCount?: number
    /**
     * Time
     * @example 1686346668
     */
    timestamp?: number
    totalConnectionCount?: number
    totalResourcesCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiCountFindingEventsResponse {
    count?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiCountFindingsResponse {
    count?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata {
    /** @example 10 */
    count?: number
    /** @example "displayName" */
    displayName?: string
    /** @example "key" */
    key?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFinding {
    /** @example "azure_cis_v140" */
    benchmarkID?: string
    /** @example 1 */
    complianceJobID?: number
    /** @example "alarm" */
    complianceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    integrationID?: string
    /** @example "Azure" */
    integrationType?: SourceType
    /** @example "azure_cis_v140_7_5" */
    controlID?: string
    controlTitle?: string
    /** @example 0.5 */
    costOptimization?: number
    /** @example 1589395200 */
    evaluatedAt?: number
    /** @example "steampipe-v0.5" */
    evaluator?: string
    /** @example "1" */
    id?: string
    /** @example "/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1" */
    platformResourceID?: string
    /** @example "1589395200" */
    lastEvent?: string
    /** @example ["Azure CIS v1.4.0"] */
    parentBenchmarkNames?: string[]
    parentBenchmarkReferences?: string[]
    parentBenchmarks?: string[]
    /** @example 1 */
    parentComplianceJobID?: number
    /**
     * Connection ID
     * @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"
     */
    providerID?: string
    /**
     * Connection ID
     * @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"
     */
    integrationName?: string
    /** @example "The VM is not using managed disks" */
    reason?: string
    /** @example "/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1" */
    resourceID?: string
    /** @example "eastus" */
    resourceLocation?: string
    /** @example "vm-1" */
    resourceName?: string
    /** @example "Microsoft.Compute/virtualMachines" */
    resourceType?: string
    /** @example "Virtual Machine" */
    resourceTypeName?: string
    /** @example "low" */
    severity?: TypesFindingSeverity
    sortKey?: any[]
    /** @example true */
    stateActive?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEvent {
    /** @example "azure_cis_v140" */
    benchmarkID?: string
    complianceJobID?: number
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    connectionID?: string
    /** @example "Azure" */
    connector?: SourceType
    /** @example "azure_cis_v140_7_5" */
    controlID?: string
    evaluatedAt?: string
    findingID?: string
    /** @example "8e0f8e7a1b1c4e6fb7e49c6af9d2b1c8" */
    id?: string
    /** @example "/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1" */
    platformResourceID?: string
    parentBenchmarkReferences?: string[]
    parentComplianceJobID?: number
    previousConformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus
    previousStateActive?: boolean
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    providerConnectionID?: string
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    providerConnectionName?: string
    reason?: string
    /** @example "/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1" */
    resourceID?: string
    /** @example "eastus" */
    resourceLocation?: string
    /** @example "vm-1" */
    resourceName?: string
    /** @example "Microsoft.Compute/virtualMachines" */
    resourceType?: string
    /**
     * Fake fields (won't be stored in ES)
     * @example "Virtual Machine"
     */
    resourceTypeName?: string
    /** @example "low" */
    severity?: TypesFindingSeverity
    sortKey?: any[]
    stateActive?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEventFilters {
    /** @example ["azure_cis_v140"] */
    benchmarkID?: string[]
    /** @example ["alarm"] */
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus[]
    /** @example ["8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"] */
    connectionID?: string[]
    /** @example ["Azure"] */
    connector?: SourceType[]
    /** @example ["azure_cis_v140_7_5"] */
    controlID?: string[]
    evaluatedAt?: {
        from?: number
        to?: number
    }
    /** @example ["8e0f8e7a1b1c4e6fb7e49c6af9d2b1c8"] */
    findingID?: string[]
    /** @example ["/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"] */
    platformResourceID?: string[]
    /** @example ["8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"] */
    notConnectionID?: string[]
    /** @example ["/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines"] */
    resourceType?: string[]
    /** @example ["low"] */
    severity?: TypesFindingSeverity[]
    /** @example [true] */
    stateActive?: boolean[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEventFiltersWithMetadata {
    benchmarkID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    connectionID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    connector?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    controlID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    resourceCollection?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    resourceTypeID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    severity?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    stateActive?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEventsSort {
    benchmarkID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    connectionID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    connector?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    controlID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    platformResourceID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    resourceType?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    severity?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    stateActive?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingFilters {
    /** @example ["azure_cis_v140"] */
    benchmarkID?: string[]
    /** @example ["alarm"] */
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus[]
    /** @example ["8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"] */
    integrationID?: string[]
    /** @example ["Azure"] */
    connector?: SourceType[]
    /** @example ["azure_cis_v140_7_5"] */
    controlID?: string[]
    evaluatedAt?: {
        from?: number
        to?: number
    }
    lastEvent?: {
        from?: number
        to?: number
    }
    /** @example ["8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"] */
    notConnectionID?: string[]
    /** @example ["/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"] */
    resourceID?: string[]
    /** @example ["/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines"] */
    resourceTypeID?: string[]
    /** @example ["low"] */
    severity?: TypesFindingSeverity[]
    /** @example [true] */
    stateActive?: boolean[]
    /**@example [477] */
    jobID?: string[]
    connectionGroup?: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingFiltersWithMetadata {
    benchmarkID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    integrationID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    connector?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    controlID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    resourceCollection?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    resourceTypeID?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    severity?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
    stateActive?: GithubComKaytuIoKaytuEnginePkgComplianceApiFilterWithMetadata[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingKPIResponse {
    failedConnectionCount?: number
    failedControlCount?: number
    failedFindingsCount?: number
    failedResourceCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiFindingsSort {
    benchmarkID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    connectionID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    connector?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    controlID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    platformResourceID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    resourceID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    resourceTypeID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    severity?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    stateActive?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetAccountsFindingsSummaryResponse {
    accounts?: GithubComKaytuIoKaytuEnginePkgComplianceApiAccountsFindingsSummary[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingEventsByFindingIDResponse {
    findingEvents?: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEvent[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingEventsRequest {
    afterSortKey?: any[]
    filters?: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEventFilters
    /** @example 100 */
    limit?: number
    sort?: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEventsSort[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingEventsResponse {
    findingEvents?: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEvent[]
    /** @example 100 */
    totalCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingsRequest {
    afterSortKey?: any[]
    filters?: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingFilters
    /** @example 100 */
    limit?: number
    sort?: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingsSort[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingsResponse {
    findings?: GithubComKaytuIoKaytuEnginePkgComplianceApiFinding[]
    /** @example 100 */
    totalCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetServicesFindingsSummaryResponse {
    services?: GithubComKaytuIoKaytuEnginePkgComplianceApiServiceFindingsSummary[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetSingleResourceFindingRequest {
    /** @example "/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1" */
    platformResourceID?: string
    /** @example "Microsoft.Compute/virtualMachines" */
    resourceType?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetSingleResourceFindingResponse {
    controls?: GithubComKaytuIoKaytuEnginePkgComplianceApiFinding[]
    findingEvents?: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEvent[]
    resource?: EsResource
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiGetTopFieldResponse {
    records?: GithubComKaytuIoKaytuEnginePkgComplianceApiTopFieldRecord[]
    /** @example 100 */
    totalCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiListBenchmarksSummaryResponse {
    benchmarkSummary?: GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkEvaluationSummary[]
    totalChecks?: TypesSeverityResult
    totalConformanceStatusSummary?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatusSummary
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiListResourceFindingsRequest {
    afterSortKey?: any[]
    filters?: GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFindingFilters
    /** @example 100 */
    limit?: number
    sort?: GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFindingsSort[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiListResourceFindingsResponse {
    resourceFindings?: GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFinding[]
    totalCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiQuery {
    Global?: boolean
    /** @example ["Azure"] */
    connector?: SourceType[]
    /** @example "2023-06-07T14:00:15.677558Z" */
    createdAt?: string
    /** @example "steampipe-v0.5" */
    engine?: string
    /** @example "azure_ad_manual_control" */
    id?: string
    /** @example ["null"] */
    listOfTables?: string[]
    parameters?: GithubComKaytuIoKaytuEnginePkgComplianceApiQueryParameter[]
    /** @example "null" */
    primaryTable?: string
    /**
     * @example "select
     *   -- Required Columns
     *   'active_directory' as resource,
     *   'info' as status,
     *   'Manual verification required.' as reason;
     * "
     */
    queryToExecute?: string
    /** @example "2023-06-16T14:58:08.759554Z" */
    updatedAt?: string
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiQueryParameter {
    /** @example "key" */
    key?: string
    /** @example true */
    required?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFinding {
    /**
     * Connection ID
     * @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"
     */
    connectionID?: string
    integrationType?: string
    evaluatedAt?: string
    failedCount?: number
    findings?: GithubComKaytuIoKaytuEnginePkgComplianceApiFinding[]
    id?: string
    platformResourceID?: string
    /**
     * Connection ID
     * @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"
     */
    integrationID?: string
    /**
     * Connection ID
     * @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"
     */
    integrationName?: string
    resourceLocation?: string
    resourceName?: string
    resourceType?: string
    resourceTypeLabel?: string
    sortKey?: any[]
    totalCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFindingFilters {
    /** @example ["azure_cis_v140"] */
    benchmarkID?: string[]
    /** @example ["alarm"] */
    complianceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus[]
    /** @example ["8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"] */
    connectionID?: string[]
    /** @example ["Azure"] */
    connector?: SourceType[]
    /** @example ["azure_cis_v140_7_5"] */
    controlID?: string[]
    evaluatedAt?: {
        from?: number
        to?: number
    }
    /** @example ["8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"] */
    notConnectionID?: string[]
    /** @example ["example-rc"] */
    resourceCollection?: string[]
    /** @example ["/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"] */
    resourceID?: string[]
    /** @example ["/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines"] */
    resourceTypeID?: string[]
    /** @example ["low"] */
    severity?: TypesFindingSeverity[]
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFindingsSort {
    conformanceStatus?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    failedCount?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    platformResourceID?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    resourceLocation?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    resourceName?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
    resourceType?: GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiServiceFindingsSummary {
    conformanceStatusesCount?: {
        failed?: number
        passed?: number
    }
    securityScore?: number
    serviceLabel?: string
    serviceName?: string
    severitiesCount?: {
        critical?: number
        high?: number
        low?: number
        medium?: number
        none?: number
    }
}

export enum GithubComKaytuIoKaytuEnginePkgComplianceApiSortDirection {
    SortDirectionAscending = 'asc',
    SortDirectionDescending = 'desc',
}

export interface GithubComKaytuIoKaytuEnginePkgComplianceApiTopFieldRecord {
    connection?: GithubComKaytuIoKaytuEnginePkgOnboardApiConnection
    control?: GithubComKaytuIoKaytuEnginePkgComplianceApiControl
    controlCount?: number
    controlTotalCount?: number
    count?: number
    field?: string
    resourceCount?: number
    resourceTotalCount?: number
    resourceType?: GithubComKaytuIoKaytuEnginePkgInventoryApiResourceType
    service?: string
    totalCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgDescribeApiJob {
    connectionID?: string
    connectionProviderID?: string
    connectionProviderName?: string
    createdAt?: string
    failureReason?: string
    id?: number
    status?: string
    title?: string
    type?: GithubComKaytuIoKaytuEnginePkgDescribeApiJobType
    updatedAt?: string
}

export interface GithubComKaytuIoKaytuEnginePkgDescribeApiJobSeqCheckResponse {
    isRunning?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgDescribeApiJobSummary {
    count?: number
    status?: string
    type?: GithubComKaytuIoKaytuEnginePkgDescribeApiJobType
}

export enum GithubComKaytuIoKaytuEnginePkgDescribeApiJobType {
    JobTypeDiscovery = 'discovery',
    JobTypeAnalytics = 'analytics',
    JobTypeCompliance = 'compliance',
}

export interface GithubComKaytuIoKaytuEnginePkgDescribeApiListDiscoveryResourceTypes {
    awsResourceTypes?: string[]
    azureResourceTypes?: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgDescribeApiListJobsRequest {
    hours?: number
    pageEnd?: number
    pageStart?: number
    sortBy?: string
    sortOrder?: string
    statusFilter?: string[]
    typeFilters?: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgDescribeApiListJobsResponse {
    jobs?: GithubComKaytuIoKaytuEnginePkgDescribeApiJob[]
    summaries?: GithubComKaytuIoKaytuEnginePkgDescribeApiJobSummary[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiAnalyticsCategoriesResponse {
    categoryResourceType?: Record<string, string[]>
}
export interface GithubComKaytuIoKaytuEnginePkgInventoryApiInventoryCategoriesResponse {
    categories: GithubComKaytuIoKaytuEnginePkgInventoryApiInventoryCategoriesResponseCategory[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiInventoryCategoriesResponseCategory {
    category: string
    tables: GithubComKaytuIoKaytuEnginePkgInventoryApiInventoryCategoriesResponseTable[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiInventoryCategoriesResponseTable {
    name: string
    table: string
    resource_type: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiAnalyticsMetric {
    connectors?: SourceType[]
    finderPerConnectionQuery?: string
    finderQuery?: string
    id?: string
    name?: string
    query?: string
    tables?: string[]
    tags?: Record<string, string[]>
    type?: GithubComKaytuIoKaytuEnginePkgAnalyticsDbMetricType
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiAssetTableRow {
    connector?: SourceType
    /** @example "compute" */
    dimensionId?: string
    /** @example "Compute" */
    dimensionName?: string
    resourceCount?: Record<string, number>
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiCostMetric {
    /** @example ["Azure"] */
    connector?: SourceType[]
    /** @example "microsoft_compute_disks" */
    cost_dimension_id?: string
    /** @example "microsoft.compute/disks" */
    cost_dimension_name?: string
    /**
     * @min 0
     * @example 14118.81523108568
     */
    daily_cost_at_end_time?: number
    /**
     * @min 0
     * @example 21232.10443638001
     */
    daily_cost_at_start_time?: number
    finderPerConnectionQuery?: string
    finderQuery?: string
    /**
     * @min 0
     * @example 621041.2436112489
     */
    total_cost?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiCostStackedItem {
    category?: string[]
    cost?: number
    metricID?: string
    metricName?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiCostTrendDatapoint {
    /** @min 0 */
    cost?: number
    costStacked?: GithubComKaytuIoKaytuEnginePkgInventoryApiCostStackedItem[]
    /** @format date-time */
    date?: string
    totalConnectionCount?: number
    totalSuccessfulDescribedConnectionCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiCountAnalyticsMetricsResponse {
    connectionCount?: number
    metricCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiCountAnalyticsSpendResponse {
    connectionCount?: number
    metricCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiCountPair {
    /** @min 0 */
    count?: number
    /** @min 0 */
    old_count?: number
}

export enum GithubComKaytuIoKaytuEnginePkgInventoryApiDirectionType {
    DirectionAscending = 'asc',
    DirectionDescending = 'desc',
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiListCostCompositionResponse {
    /**
     * @min 0
     * @example 100
     */
    others?: number
    top_values?: Record<string, number>
    /**
     * @min 0
     * @example 1000
     */
    total_cost_value?: number
    /**
     * @min 0
     * @example 10
     */
    total_count?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiListCostMetricsResponse {
    metrics?: GithubComKaytuIoKaytuEnginePkgInventoryApiCostMetric[]
    /**
     * @min 0
     * @example 1000
     */
    total_cost?: number
    /**
     * @min 0
     * @example 10
     */
    total_count?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiListMetricsResponse {
    metrics?: GithubComKaytuIoKaytuEnginePkgInventoryApiMetric[]
    total_count?: number
    total_metrics?: number
    total_old_count?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequest {
    /** Specifies the Title */
    titleFilter?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiListResourceTypeCompositionResponse {
    others?: GithubComKaytuIoKaytuEnginePkgInventoryApiCountPair
    top_values?: Record<
        string,
        GithubComKaytuIoKaytuEnginePkgInventoryApiCountPair
    >
    /** @min 0 */
    total_count?: number
    /** @min 0 */
    total_value_count?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiMetric {
    /**
     * Cloud Provider
     * @example ["[Azure]"]
     */
    connectors?: SourceType[]
    /**
     * Number of Resources of this Resource Type - Metric
     * @example 100
     */
    count?: number
    /** @example "select * from kaytu_resources where resource_type = 'aws::ec2::instance' AND connection_id IN <CONNECTION_ID_LIST>" */
    finderPerConnectionQuery?: string
    /** @example "select * from kaytu_resources where resource_type = 'aws::ec2::instance'" */
    finderQuery?: string
    /** @example "vms" */
    id?: string
    /**
     * Last time the metric was evaluated
     * @example "2020-01-01T00:00:00Z"
     */
    last_evaluated?: string
    /**
     * Resource Type
     * @example "VMs"
     */
    name?: string
    /**
     * Number of Resources of this Resource Type in the past - Metric
     * @example 90
     */
    old_count?: number
    /** Tags */
    tags?: Record<string, string[]>
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiPage {
    no?: number
    size?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollection {
    connection_count?: number
    connectors?: SourceType[]
    created_at?: string
    description?: string
    filters?: KaytuResourceCollectionFilter[]
    id?: string
    last_evaluated_at?: string
    metric_count?: number
    name?: string
    resource_count?: number
    status?: GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionStatus
    tags?: Record<string, string[]>
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscape {
    categories?: GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscapeCategory[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscapeCategory {
    description?: string
    id?: string
    name?: string
    subcategories?: GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscapeSubcategory[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscapeItem {
    description?: string
    id?: string
    logo_uri?: string
    name?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscapeSubcategory {
    description?: string
    id?: string
    items?: GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscapeItem[]
    name?: string
}

export enum GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionStatus {
    ResourceCollectionStatusUnknown = '',
    ResourceCollectionStatusActive = 'active',
    ResourceCollectionStatusInactive = 'inactive',
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCountStackedItem {
    category?: string[]
    count?: number
    metricID?: string
    metricName?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceType {
    /** List supported steampipe Attributes (columns) for this resource type - Metadata (GET only) */
    attributes?: string[]
    /** List of Compliance that support this Resource Type - Metadata (GET only) */
    compliance?: string[]
    /**
     * Number of Compliance that use this Resource Type - Metadata
     * @min 0
     */
    compliance_count?: number
    /**
     * Cloud Provider
     * @example "Azure"
     */
    connector?: SourceType
    /**
     * Number of Resources of this Resource Type - Metric
     * @min 0
     * @example 100
     */
    count?: number
    /**
     * Logo URI
     * @example "https://kaytu.io/logo.png"
     */
    logo_uri?: string
    /**
     * Number of Resources of this Resource Type in the past - Metric
     * @min 0
     * @example 90
     */
    old_count?: number
    /**
     * Resource Name
     * @example "VM"
     */
    resource_name?: string
    /**
     * Resource Type
     * @example "Microsoft.Compute/virtualMachines"
     */
    resource_type?: string
    /**
     * Service Name
     * @example "compute"
     */
    service_name?: string
    /**
     * Tags
     * @example ["category:[Data and Analytics","Database","Integration","Management Governance","Storage]"]
     */
    tags?: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiResourceTypeTrendDatapoint {
    /**
     * @min 0
     * @example 100
     */
    count?: number
    countStacked?: GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCountStackedItem[]
    /** @format date-time */
    date?: string
    totalConnectionCount?: number
    totalSuccessfulDescribedConnectionCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiRunQueryRequest {
    account_id?: string
    engine?: string
    page: GithubComKaytuIoKaytuEnginePkgInventoryApiPage
    query?: string
    sorts?: GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQuerySortItem[]
    source_id?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiRunQueryResponse {
    /** Column names */
    headers?: string[]
    /** Query */
    query?: string
    /** Result of query. in order to access a specific cell please use Result[Row][Column] */
    result?: any[][]
    /** Query Title */
    title?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryHistory {
    executed_at?: string
    query?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItem {
    /** Category (Tags[category]) */
    category?: string
    /** Provider */
    connectors?: SourceType[]
    /** Query Id */
    id?: string
    /** Query */
    query?: string
    /** Tags */
    tags?: Record<string, string>
    /** Title */
    title?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQuerySortItem {
    direction?: 'asc' | 'desc'
    /** fill this with column name */
    field?: string
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSpendTableRow {
    /** @example "1239042" */
    accountID?: string
    /** @example "Compute" */
    category?: string
    /** @example "AWS" */
    connector?: SourceType
    costValue?: Record<string, number>
    /** @example "compute" */
    dimensionId?: string
    /** @example "Compute" */
    dimensionName?: string
}

export interface GithubComKaytuIoKaytuEnginePkgMetadataApiListQueryParametersResponse {
    queryParameters?: GithubComKaytuIoKaytuEnginePkgMetadataApiQueryParameter[]
}

export interface GithubComKaytuIoKaytuEnginePkgMetadataApiQueryParameter {
    key?: string
    value?: string
}

export interface GithubComKaytuIoKaytuEnginePkgMetadataApiSetConfigMetadataRequest {
    key?: string
    value?: any
}

export interface GithubComKaytuIoKaytuEnginePkgMetadataApiSetQueryParameterRequest {
    queryParameters?: GithubComKaytuIoKaytuEnginePkgMetadataApiQueryParameter[]
}

export interface GithubComKaytuIoKaytuEnginePkgMetadataModelsConfigMetadata {
    key?: GithubComKaytuIoKaytuEnginePkgMetadataModelsMetadataKey
    type?: GithubComKaytuIoKaytuEnginePkgMetadataModelsConfigMetadataType
    value?: string
}

export enum GithubComKaytuIoKaytuEnginePkgMetadataModelsConfigMetadataType {
    ConfigMetadataTypeString = 'string',
    ConfigMetadataTypeInt = 'int',
    ConfigMetadataTypeBool = 'bool',
    ConfigMetadataTypeJSON = 'json',
}

export interface GithubComKaytuIoKaytuEnginePkgMetadataModelsFilter {
    kayValue?: Record<string, string>
    name?: string
}

export enum GithubComKaytuIoKaytuEnginePkgMetadataModelsMetadataKey {
    MetadataKeyWorkspaceOwnership = 'workspace_ownership',
    MetadataKeyWorkspaceID = 'workspace_id',
    MetadataKeyWorkspaceName = 'workspace_name',
    MetadataKeyWorkspacePlan = 'workspace_plan',
    MetadataKeyWorkspaceCreationTime = 'workspace_creation_time',
    MetadataKeyWorkspaceDateTimeFormat = 'workspace_date_time_format',
    MetadataKeyWorkspaceDebugMode = 'workspace_debug_mode',
    MetadataKeyWorkspaceTimeWindow = 'workspace_time_window',
    MetadataKeyAssetManagementEnabled = 'asset_management_enabled',
    MetadataKeyComplianceEnabled = 'compliance_enabled',
    MetadataKeyProductManagementEnabled = 'product_management_enabled',
    MetadataKeyCustomIDP = 'custom_idp',
    MetadataKeyResourceLimit = 'resource_limit',
    MetadataKeyConnectionLimit = 'connection_limit',
    MetadataKeyUserLimit = 'user_limit',
    MetadataKeyAllowInvite = 'allow_invite',
    MetadataKeyWorkspaceKeySupport = 'workspace_key_support',
    MetadataKeyWorkspaceMaxKeys = 'workspace_max_keys',
    MetadataKeyAllowedEmailDomains = 'allowed_email_domains',
    MetadataKeyAutoDiscoveryMethod = 'auto_discovery_method',
    MetadataKeyDescribeJobInterval = 'describe_job_interval',
    MetadataKeyFullDiscoveryJobInterval = 'full_discovery_job_interval',
    MetadataKeyCostDiscoveryJobInterval = 'cost_discovery_job_interval',
    MetadataKeyHealthCheckJobInterval = 'health_check_job_interval',
    MetadataKeyMetricsJobInterval = 'metrics_job_interval',
    MetadataKeyComplianceJobInterval = 'compliance_job_interval',
    MetadataKeyDataRetention = 'data_retention_duration',
    MetadataKeyAnalyticsGitURL = 'analytics_git_url',
    MetadataKeyAssetDiscoveryAWSPolicyARNs = 'asset_discovery_aws_policy_arns',
    MetadataKeySpendDiscoveryAWSPolicyARNs = 'spend_discovery_aws_policy_arns',
    MetadataKeyAssetDiscoveryAzureRoleIDs = 'asset_discovery_azure_role_ids',
    MetadataKeySpendDiscoveryAzureRoleIDs = 'spend_discovery_azure_role_ids',
    MetadataKeyCustomizationEnabled = 'customization_enabled',
    MetadataKeyAWSDiscoveryRequiredOnly = 'aws_discovery_required_only',
    MetadataKeyAzureDiscoveryRequiredOnly = 'azure_discovery_required_only',
    MetadataKeyAssetDiscoveryEnabled = 'asset_discovery_enabled',
    MetadataKeySpendDiscoveryEnabled = 'spend_discovery_enabled',
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiAWSCredentialConfig {
    accessKey: string
    accountId?: string
    assumeAdminRoleName?: string
    assumeRoleName?: string
    externalId?: string
    regions?: string[]
    secretKey: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiAzureCredentialConfig {
    clientId: string
    clientSecret: string
    objectId: string
    secretId: string
    subscriptionId?: string
    tenantId: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiCatalogMetrics {
    /**
     * @min 0
     * @example 20
     */
    connectionsEnabled?: number
    /**
     * @min 0
     * @example 15
     */
    healthyConnections?: number
    /**
     * @min 0
     * @example 5
     */
    inProgressConnections?: number
    /**
     * @min 0
     * @example 20
     */
    totalConnections?: number
    /**
     * @min 0
     * @example 5
     */
    unhealthyConnections?: number
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiChangeConnectionLifecycleStateRequest {
    state?: GithubComKaytuIoKaytuEnginePkgOnboardApiConnectionLifecycleState
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiConnection {
    assetDiscovery?: boolean
    /** @example "scheduled" */
    assetDiscoveryMethod?: SourceAssetDiscoveryMethodType
    /** @example "Azure" */
    connector?: SourceType
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    cost?: number
    credential?: GithubComKaytuIoKaytuEnginePkgOnboardApiCredential
    /** @example "7r6123ac-ca1c-434f-b1a3-91w2w9d277c8" */
    credentialID?: string
    credentialName?: string
    /** @example "manual" */
    credentialType?: GithubComKaytuIoKaytuEnginePkgOnboardApiCredentialType
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    dailyCostAtEndTime?: number
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    dailyCostAtStartTime?: number
    describeJobRunning?: boolean
    /** @example "This is an example connection" */
    description?: string
    /** @example "johndoe@example.com" */
    email?: string
    healthReason?: string
    /** @example "healthy" */
    healthState?: SourceHealthStatus
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    id?: string
    /** @example "2023-05-07T00:00:00Z" */
    lastHealthCheckTime?: string
    /** @example "2023-05-07T00:00:00Z" */
    lastInventory?: string
    /** @example "enabled" */
    lifecycleState?: GithubComKaytuIoKaytuEnginePkgOnboardApiConnectionLifecycleState
    metadata?: Record<string, any>
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    oldResourceCount?: number
    /** @example "2023-05-07T00:00:00Z" */
    onboardDate?: string
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    providerConnectionID?: string
    /** @example "example-connection" */
    providerConnectionName?: string
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    resourceCount?: number
    spendDiscovery?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiConnectionGroup {
    /** @example ["[\"1e8ac3bf-c268-4a87-9374-ce04cc40a596\"]"] */
    connectionIds?: string[]
    connections?: GithubComKaytuIoKaytuEnginePkgOnboardApiConnection[]
    /** @example "UltraSightApplication" */
    name?: string
    /** @example "SELECT kaytu_id FROM kaytu_connections WHERE tags->'application' IS NOT NULL AND tags->'application' @> '"UltraSight"'" */
    query?: string
}

export enum GithubComKaytuIoKaytuEnginePkgOnboardApiConnectionLifecycleState {
    ConnectionLifecycleStateOnboard = 'ONBOARD',
    ConnectionLifecycleStateDisabled = 'DISABLED',
    ConnectionLifecycleStateDiscovered = 'DISCOVERED',
    ConnectionLifecycleStateInProgress = 'IN_PROGRESS',
    ConnectionLifecycleStateArchived = 'ARCHIVED',
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiConnectorCount {
    /** @example true */
    allowNewConnections?: boolean
    /** @example false */
    autoOnboardSupport?: boolean
    /**
     * @min 0
     * @example 1024
     */
    connection_count?: number
    /** @example "This is a long volume of words for just showing the case of the description for the demo and checking value purposes only and has no meaning whatsoever" */
    description?: string
    direction?: SourceConnectorDirectionType
    /** @example "Azure" */
    label?: string
    /** @example "https://kaytu.io/logo.png" */
    logo?: string
    /**
     * @min 0
     * @example 10000
     */
    maxConnectionLimit?: number
    /** @example "Azure" */
    name?: SourceType
    /** @example "This is a short Description for this connector" */
    shortDescription?: string
    /** @example "enabled" */
    status?: SourceConnectorStatus
    tags?: Record<string, any>
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiCreateAwsConnectionRequest {
    awsConfig?: GithubComKaytuIoKaytuEnginePkgOnboardApiV2AWSCredentialV2Config
    name?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiCreateConnectionResponse {
    id?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiCreateCredentialRequest {
    config?: any
    /** @example "Azure" */
    source_type?: SourceType
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiCreateCredentialResponse {
    id?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiCreateSourceResponse {
    id?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiCredential {
    /**
     * @min 0
     * @max 1000
     * @example 0
     */
    archived_connections?: number
    /** @example false */
    autoOnboardEnabled?: boolean
    config?: any
    connections?: GithubComKaytuIoKaytuEnginePkgOnboardApiConnection[]
    /** @example "AWS" */
    connectorType?: SourceType
    /** @example "manual-aws-org" */
    credentialType?: GithubComKaytuIoKaytuEnginePkgOnboardApiCredentialType
    /**
     * @min 0
     * @max 1000
     * @example 0
     */
    disabled_connections?: number
    /**
     * @min 0
     * @max 100
     * @example 50
     */
    discovered_connections?: number
    /** @example true */
    enabled?: boolean
    /** @example "" */
    healthReason?: string
    /** @example "healthy" */
    healthStatus?: SourceHealthStatus
    /** @example "1028642a-b22e-26ha-c5h2-22nl254678m5" */
    id?: string
    /**
     * @format date-time
     * @example "2023-06-03T12:21:33.406928Z"
     */
    lastHealthCheckTime?: string
    metadata?: Record<string, any>
    /** @example "a-1mahsl7lzk" */
    name?: string
    /**
     * @format date-time
     * @example "2023-06-03T12:21:33.406928Z"
     */
    onboardDate?: string
    /**
     * @min 0
     * @max 1000
     * @example 250
     */
    onboard_connections?: number
    spendDiscovery?: boolean
    /**
     * @min 0
     * @max 1000
     * @example 300
     */
    total_connections?: number
    /**
     * @min 0
     * @max 100
     * @example 50
     */
    unhealthy_connections?: number
    version?: number
}

export enum GithubComKaytuIoKaytuEnginePkgOnboardApiCredentialType {
    CredentialTypeAutoAzure = 'auto-azure',
    CredentialTypeAutoAws = 'auto-aws',
    CredentialTypeManualAwsOrganization = 'manual-aws-org',
    CredentialTypeManualAzureSpn = 'manual-azure-spn',
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiListConnectionSummaryResponse {
    /**
     * @min 0
     * @max 1000
     * @example 10
     */
    connectionCount?: number
    connections?: GithubComKaytuIoKaytuEnginePkgOnboardApiConnection[]
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalArchivedCount?: number
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    totalCost?: number
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalDisabledCount?: number
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalDiscoveredCount?: number
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    totalOldResourceCount?: number
    /**
     * Also includes in-progress
     * @min 0
     * @max 100
     * @example 10
     */
    totalOnboardedCount?: number
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    totalResourceCount?: number
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalUnhealthyCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiListCredentialResponse {
    credentials?: GithubComKaytuIoKaytuEnginePkgOnboardApiCredential[]
    /**
     * @min 0
     * @max 20
     * @example 5
     */
    totalCredentialCount?: number
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiSourceAwsRequest {
    config?: GithubComKaytuIoKaytuEnginePkgOnboardApiAWSCredentialConfig
    description?: string
    email?: string
    name?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiSourceAzureRequest {
    config?: GithubComKaytuIoKaytuEnginePkgOnboardApiAzureCredentialConfig
    description?: string
    name?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiUpdateCredentialRequest {
    config?: any
    /** @example "Azure" */
    connector?: SourceType
    name?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiV2AWSCredentialV2Config {
    accessKey?: string
    accountID?: string
    assumeRoleName?: string
    externalId?: string
    secretKey?: string
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiV2CreateCredentialV2Request {
    awsConfig?: GithubComKaytuIoKaytuEnginePkgOnboardApiV2AWSCredentialV2Config
    /** @example "Azure" */
    connector?: SourceType
}

export interface GithubComKaytuIoKaytuEnginePkgOnboardApiV2CreateCredentialV2Response {
    id?: string
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiAddCredentialRequest {
    awsConfig?: GithubComKaytuIoKaytuEnginePkgOnboardApiV2AWSCredentialV2Config
    azureConfig?: GithubComKaytuIoKaytuEnginePkgOnboardApiAzureCredentialConfig
    connectorType?: SourceType
    singleConnection?: boolean
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiBootstrapProgress {
    done?: number
    total?: number
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiBootstrapStatusResponse {
    analyticsStatus?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiBootstrapProgress
    complianceStatus?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiBootstrapProgress
    connection_count?: Record<string, number>
    discoveryStatus?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiBootstrapProgress
    maxConnections?: number
    minRequiredConnections?: number
    workspaceCreationStatus?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiBootstrapProgress
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiCreateWorkspaceRequest {
    name?: string
    organization_id?: number
    tier?: string
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiCreateWorkspaceResponse {
    id?: string
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiOrganization {
    address?: string
    city?: string
    companyName?: string
    contactEmail?: string
    contactName?: string
    contactPhone?: string
    country?: string
    id?: number
    state?: string
    url?: string
}

export enum GithubComKaytuIoKaytuEnginePkgWorkspaceApiStateID {
    StateIDReserving = 'RESERVING',
    StateIDReserved = 'RESERVED',
    StateIDWaitingForCredential = 'WAITING_FOR_CREDENTIAL',
    StateIDProvisioning = 'PROVISIONING',
    StateIDProvisioned = 'PROVISIONED',
    StateIDDeleting = 'DELETING',
    StateIDDeleted = 'DELETED',
}

export enum GithubComKaytuIoKaytuEnginePkgWorkspaceApiTier {
    TierFree = 'FREE',
    TierTeams = 'TEAMS',
    TierEnterprise = 'ENTERPRISE',
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceLimitsUsage {
    /** @example 100 */
    currentConnections?: number
    /** @example 10000 */
    currentResources?: number
    /** @example 10 */
    currentUsers?: number
    /** @example "ws-698542025141040315" */
    id?: string
    /** @example 1000 */
    maxConnections?: number
    /** @example 100000 */
    maxResources?: number
    /** @example 100 */
    maxUsers?: number
    /** @example "kaytu" */
    name?: string
}

export interface GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceResponse {
    /** @example "kaytu" */
    aws_unique_id?: string
    /** @example "kaytu" */
    aws_user_arn?: string
    /** @example "2023-05-17T14:39:02.707659Z" */
    createdAt?: string
    /** @example "ws-698542025141040315" */
    id?: string
    is_bootstrap_input_finished?: boolean
    is_created?: boolean
    /** @example "kaytu" */
    name?: string
    organization?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiOrganization
    /** @example "google-oauth2|204590896945502695694" */
    ownerId?: string
    /** @example "sm" */
    size?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceSize
    /** @example "PROVISIONED" */
    status?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiStateID
    /** @example "ENTERPRISE" */
    tier?: GithubComKaytuIoKaytuEnginePkgWorkspaceApiTier
    /** @example "v0.45.4" */
    version?: string
}

export enum GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceSize {
    SizeXS = 'xs',
    SizeSM = 'sm',
    SizeMD = 'md',
    SizeLG = 'lg',
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityAWSCredentialConfig {
    accessKey?: string
    accountID?: string
    assumeRoleName?: string
    externalId?: string
    secretKey?: string
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityAzureCredentialConfig {
    clientId: string
    clientSecret: string
    objectId: string
    tenantId: string
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCatalogMetrics {
    /**
     * @min 0
     * @example 20
     */
    connectionsEnabled?: number
    /**
     * @min 0
     * @example 15
     */
    healthyConnections?: number
    /**
     * @min 0
     * @example 5
     */
    inProgressConnections?: number
    /**
     * @min 0
     * @example 20
     */
    totalConnections?: number
    /**
     * @min 0
     * @example 5
     */
    unhealthyConnections?: number
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection {
    assetDiscovery?: boolean
    /** @example "scheduled" */
    assetDiscoveryMethod?: SourceAssetDiscoveryMethodType
    /** @example "Azure" */
    connector?: SourceType
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    cost?: number
    credential?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredential
    /** @example "7r6123ac-ca1c-434f-b1a3-91w2w9d277c8" */
    credentialID?: string
    credentialName?: string
    /** @example "manual" */
    credentialType?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredentialType
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    dailyCostAtEndTime?: number
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    dailyCostAtStartTime?: number
    describeJobRunning?: boolean
    /** @example "This is an example connection" */
    description?: string
    /** @example "johndoe@example.com" */
    email?: string
    healthReason?: string
    /** @example "healthy" */
    healthState?: SourceHealthStatus
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    id?: string
    /** @example "2023-05-07T00:00:00Z" */
    lastHealthCheckTime?: string
    /** @example "2023-05-07T00:00:00Z" */
    lastInventory?: string
    /** @example "enabled" */
    lifecycleState?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnectionLifecycleState
    metadata?: Record<string, any>
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    oldResourceCount?: number
    /** @example "2023-05-07T00:00:00Z" */
    onboardDate?: string
    /** @example "8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8" */
    providerConnectionID?: string
    /** @example "example-connection" */
    providerConnectionName?: string
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    resourceCount?: number
    spendDiscovery?: boolean
}

export enum GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnectionLifecycleState {
    ConnectionLifecycleStateOnboard = 'ONBOARD',
    ConnectionLifecycleStateDisabled = 'DISABLED',
    ConnectionLifecycleStateDiscovered = 'DISCOVERED',
    ConnectionLifecycleStateInProgress = 'IN_PROGRESS',
    ConnectionLifecycleStateArchived = 'ARCHIVED',
}
export interface ConnectorCountField {
    total:    number;
    active:   number;
    inactive: number;
    archived: number;
    demo:     number;
}


export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnectorCount {
    id: number
    name: string
    platform_name: string
    label: string
    title: string
    tier: string
    annotations: null
    labels: null
    short_description: string
    long_description: string
    logo: string
    enabled: boolean
    count: ConnectorCountField
}


export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnectorResponse {
    total_count: number
    integration_types: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnectorCount[]
}
export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCountConnectionsResponse {
    count?: number
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateAWSConnectionRequest {
    config?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityAWSCredentialConfig
    description?: string
    email?: string
    name?: string
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateAWSCredentialRequest {
    config?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityAWSCredentialConfig
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateAzureCredentialRequest {
    config?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityAzureCredentialConfig
    description?: string
    name?: string
    type?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredentialType
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateConnectionResponse {
    id?: string
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateCredentialResponse {
    connections?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection[]
    id?: string
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredential {
    /**
     * @min 0
     * @max 1000
     * @example 0
     */
    archived_connections?: number
    /** @example false */
    autoOnboardEnabled?: boolean
    config?: any
    connections?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection[]
    /** @example "AWS" */
    connectorType?: SourceType
    /** @example "manual-aws-org" */
    credentialType?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredentialType
    /**
     * @min 0
     * @max 1000
     * @example 0
     */
    disabled_connections?: number
    /**
     * @min 0
     * @max 100
     * @example 50
     */
    discovered_connections?: number
    /** @example true */
    enabled?: boolean
    /** @example "" */
    healthReason?: string
    /** @example "healthy" */
    healthStatus?: SourceHealthStatus
    /** @example "1028642a-b22e-26ha-c5h2-22nl254678m5" */
    id?: string
    /**
     * @format date-time
     * @example "2023-06-03T12:21:33.406928Z"
     */
    lastHealthCheckTime?: string
    metadata?: Record<string, any>
    /** @example "a-1mahsl7lzk" */
    name?: string
    /**
     * @format date-time
     * @example "2023-06-03T12:21:33.406928Z"
     */
    onboardDate?: string
    /**
     * @min 0
     * @max 1000
     * @example 250
     */
    onboard_connections?: number
    spendDiscovery?: boolean
    /**
     * @min 0
     * @max 1000
     * @example 300
     */
    total_connections?: number
    /**
     * @min 0
     * @max 100
     * @example 50
     */
    unhealthy_connections?: number
    version?: number
}

export enum GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredentialType {
    CredentialTypeAutoAzure = 'auto-azure',
    CredentialTypeAutoAws = 'auto-aws',
    CredentialTypeManualAwsOrganization = 'manual-aws-org',
    CredentialTypeManualAzureSpn = 'manual-azure-spn',
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityListConnectionsSummaryResponse {
    /**
     * @min 0
     * @max 1000
     * @example 10
     */
    connectionCount?: number
    connections?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection[]
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalArchivedCount?: number
    /**
     * @min 0
     * @max 10000000
     * @example 1000
     */
    totalCost?: number
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalDisabledCount?: number
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalDiscoveredCount?: number
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    totalOldResourceCount?: number
    /**
     * Also includes in-progress
     * @min 0
     * @max 100
     * @example 10
     */
    totalOnboardedCount?: number
    /**
     * @min 0
     * @max 1000000
     * @example 100
     */
    totalResourceCount?: number
    /**
     * @min 0
     * @max 100
     * @example 10
     */
    totalUnhealthyCount?: number
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityListCredentialResponse {
    credentials?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredential[]
    /**
     * @min 0
     * @max 20
     * @example 5
     */
    totalCredentialCount?: number
}

export enum GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityTier {
    TierCommunity = 'Community',
    TierEnterprise = 'Enterprise',
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityUpdateAWSCredentialRequest {
    config?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityAWSCredentialConfig
    name?: string
}

export interface GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityUpdateAzureCredentialRequest {
    config?: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityAzureCredentialConfig
    name?: string
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsClusterWastageRequest {
    cliVersion?: string
    cluster?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsCluster
    identification?: Record<string, string>
    instances?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRds[]
    loading?: boolean
    metrics?: Record<string, Record<string, TypesDatapoint[]>>
    preferences?: Record<string, string>
    region?: string
    requestId?: string
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsClusterWastageResponse {
    rightSizing?: Record<
        string,
        GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsRightsizingRecommendation
    >
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRds {
    availabilityZone?: string
    backupRetentionPeriod?: number
    clusterType?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsClusterType
    engine?: string
    engineVersion?: string
    hashedInstanceId?: string
    instanceType?: string
    licenseModel?: string
    performanceInsightsEnabled?: boolean
    performanceInsightsRetentionPeriod?: number
    storageIops?: number
    storageSize?: number
    storageThroughput?: number
    storageType?: string
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsCluster {
    engine?: string
    hashedClusterId?: string
}

export enum GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsClusterType {
    AwsRdsClusterTypeSingleInstance = 'Single-AZ',
    AwsRdsClusterTypeMultiAzOneInstance = 'Multi-AZ',
    AwsRdsClusterTypeMultiAzTwoInstance = 'Multi-AZ (readable standbys)',
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsRightsizingRecommendation {
    current?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingAwsRds
    description?: string
    freeMemoryBytes?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    freeStorageBytes?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    networkThroughputBytes?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    recommended?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingAwsRds
    storageIops?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    storageThroughputBytes?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    vCPU?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    volumeBytesUsed?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsWastageRequest {
    cliVersion?: string
    identification?: Record<string, string>
    instance?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRds
    loading?: boolean
    metrics?: Record<string, TypesDatapoint[]>
    preferences?: Record<string, string>
    region?: string
    requestId?: string
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsWastageResponse {
    rightSizing?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsRightsizingRecommendation
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityEBSVolumeRecommendation {
    current?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingEBSVolume
    description?: string
    iops?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    recommended?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingEBSVolume
    throughput?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2Instance {
    coreCount?: number
    ebsOptimized?: boolean
    hashedInstanceId?: string
    instanceLifecycle?: TypesInstanceLifecycleType
    instanceType?: TypesInstanceType
    monitoring?: TypesMonitoringState
    placement?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2Placement
    platform?: string
    state?: TypesInstanceStateName
    tenancy?: TypesTenancy
    threadsPerCore?: number
    usageOperation?: string
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2InstanceWastageRequest {
    cliVersion?: string
    identification?: Record<string, string>
    instance?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2Instance
    loading?: boolean
    metrics?: Record<string, TypesDatapoint[]>
    preferences?: Record<string, string>
    region?: string
    requestId?: string
    volumeCount?: number
    volumeMetrics?: Record<string, Record<string, TypesDatapoint[]>>
    volumes?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2Volume[]
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2InstanceWastageResponse {
    rightSizing?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightSizingRecommendation
    volumes?: Record<
        string,
        GithubComKaytuIoKaytuEngineServicesWastageApiEntityEBSVolumeRecommendation
    >
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2Placement {
    availabilityZone?: string
    hashedHostId?: string
    tenancy?: TypesTenancy
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2Volume {
    availabilityZone?: string
    hashedVolumeId?: string
    iops?: number
    size?: number
    throughput?: number
    volumeType?: TypesVolumeType
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightSizingRecommendation {
    current?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingEC2Instance
    description?: string
    ebsBandwidth?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    ebsIops?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    memory?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    networkThroughput?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
    recommended?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingEC2Instance
    vCPU?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingAwsRds {
    architecture?: string
    clusterType?: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsClusterType
    computeCost?: number
    computeCostComponents?: Record<string, number>
    cost?: number
    costComponents?: Record<string, number>
    engine?: string
    engineVersion?: string
    instanceType?: string
    memoryGb?: number
    processor?: string
    region?: string
    storageCost?: number
    storageCostComponents?: Record<string, number>
    storageIops?: number
    storageSize?: number
    storageThroughput?: number
    storageType?: string
    vCPU?: number
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingEBSVolume {
    baselineIOPS?: number
    baselineThroughput?: number
    cost?: number
    costComponents?: Record<string, number>
    provisionedIOPS?: number
    provisionedThroughput?: number
    tier?: TypesVolumeType
    volumeSize?: number
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityRightsizingEC2Instance {
    architecture?: string
    cost?: number
    costComponents?: Record<string, number>
    ebsBandwidth?: string
    ebsIops?: string
    enaSupported?: string
    instanceType?: string
    license?: string
    licensePrice?: number
    memory?: number
    networkThroughput?: string
    processor?: string
    region?: string
    vCPU?: number
}

export interface GithubComKaytuIoKaytuEngineServicesWastageApiEntityUsage {
    avg?: number
    last?: TypesDatapoint
    max?: number
    min?: number
}

export interface KaytuResourceCollectionFilter {
    account_ids?: string[]
    connectors?: string[]
    regions?: string[]
    resource_types?: string[]
    tags?: Record<string, string>
}
export interface GithubComKaytuIoKaytuEnginePkgInventoryApiV3ControlListFilters {
    provider: string[]
    severity: string[]
    root_benchmark: string[]
    parent_benchmark: string[]
    primary_table: string[]
    list_of_tables: string[]
    tags: GithubComKaytuIoKaytuEnginePkgInventoryApiV3ControlListFiltersTags[]
}
export interface GithubComKaytuIoKaytuEnginePkgControlApiListV2 {
    cursor: number
    per_page: number
    primary_table?: string
    severity?: string
    finding_summary?: boolean
    connector?: string[]
    parent_benchmark?: string[]
    root_benchmark?: string[]
    has_parameters?: boolean
    tags?: GithubComKaytuIoKaytuEnginePkgControlApiListV2Tags[]
    list_of_tables?: string[]
}
export interface GithubComKaytuIoKaytuEnginePkgInventoryApiV3ControlListFiltersTags {
    Key: string
    UniqueValues: string[]
}
export interface GithubComKaytuIoKaytuEnginePkgControlApiListV2Tags {
    [key: string]: string[]
}
export interface GithubComKaytuIoKaytuEnginePkgControlApiListV2Response {
    items: GithubComKaytuIoKaytuEnginePkgControlApiListV2ResponseItem[]
    total_count: number
}

export interface GithubComKaytuIoKaytuEnginePkgControlApiListV2ResponseItem {
    id: string
    title: string
    description: string
    connector: string[]
    severity: string
    tags: GithubComKaytuIoKaytuEnginePkgControlApiListV2ResponseItemTags
    query: GithubComKaytuIoKaytuEnginePkgControlApiListV2ResponseItemQuery
    findings_summary?: {
        non_incident_count?: number
        incident_count?: number
    }
}

export interface GithubComKaytuIoKaytuEnginePkgControlApiListV2ResponseItemQuery {
    primary_table: string
    list_of_tables: string[]
    parameters: any[]
}

export interface GithubComKaytuIoKaytuEnginePkgControlApiListV2ResponseItemTags {
    score_service_name: string[]
    score_tags: string[]
}
export interface GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3Response {
    items: GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseItem[]
    total_count: number
}
export interface GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3 {
    /** Specifies the Title */
    parent_benchmark_id?: string[]
    root: boolean
    // title_filter?: string
    cursor: number
    per_page: number
    primary_table?: string[]
    list_of_tables?: string[]
    tags?: GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseTags
    finding_summary?: boolean
}
export interface GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseItem {
    benchmark: GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseMetaData
    findings: null
}

export interface GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseMetaData {
    id: string
    title: string
    description: string
    connectors: string[]
    number_of_controls: number
    enabled: boolean
    track_drift_events: boolean
    primary_tables: string[]
    tags: GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseTags
    created_at: string
    updated_at: string
}

export interface GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseTags {
    [key: string]: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgControlDetailV3 {
    benchmarks: GithubComKaytuIoKaytuEnginePkgControlDetailV3Benchmarks
    integrationType: string[]
    description: string
    id: string
    query: GithubComKaytuIoKaytuEnginePkgControlDetailV3Query
    severity: string
    tags: GithubComKaytuIoKaytuEnginePkgControlDetailV3Tags
    title: string
}

export interface GithubComKaytuIoKaytuEnginePkgControlDetailV3Benchmarks {
    fullPath: string[]
    roots: string[]
}
export interface GithubComKaytuIoKaytuEnginePkgControlDetailV3Tags {
    [key: string]: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgControlDetailV3Query {
    engine: string
    listOfTables: string[]
    primaryTable: string
    queryToExecute: string
}
export interface GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2 {
    /** Specifies the Title */
    title_filter?: string
    cursor: number
    per_page: number
    providers?: string[]
    tags?: GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2TagsFilter
    list_of_tables?: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2TagsFilter {
    [key: string]: string[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2 {
    id: string
    title: string
    description: string
    integration_types: string[]
    query: GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2Query
    tags: GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2TagsFilter
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2Query {
    id: string
    queryToExecute: string
    primaryTable: string
    listOfTables: string[]
    engine: string
    parameters: string[]
    Global: boolean
    createdAt: Date
    updatedAt: Date
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2Response {
    /** List of items */
    items: GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2[]
    /** total caount of data */
    total_count: number
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryFilters {
    providers: string[]
    tags: GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryFiltersTag[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryFiltersTag {
    Key: string
    UniqueValues: string[]
}
export interface GithubComKaytuIoKaytuEnginePkgInventoryApiV3BenchmarkListFilters {
    parent_benchmark_id: string[]
    primary_table: string[]
    list_of_tables: string[]
    tags: GithubComKaytuIoKaytuEnginePkgInventoryApiV3BenchmarkListFiltersTag[]
}

export interface GithubComKaytuIoKaytuEnginePkgInventoryApiV3BenchmarkListFiltersTag {
    Key: string
    UniqueValues: string[]
}

export enum SourceAssetDiscoveryMethodType {
    AssetDiscoveryMethodTypeScheduled = 'scheduled',
}

export enum SourceConnectorDirectionType {
    ConnectorDirectionTypeIngress = 'ingress',
    ConnectorDirectionTypeEgress = 'egress',
    ConnectorDirectionTypeBoth = 'both',
}

export enum SourceConnectorStatus {
    ConnectorStatusEnabled = 'enabled',
    ConnectorStatusDisabled = 'disabled',
    ConnectorStatusComingSoon = 'coming_soon',
}

export enum SourceHealthStatus {
    HealthStatusNil = '',
    HealthStatusHealthy = 'healthy',
    HealthStatusUnhealthy = 'unhealthy',
}

export enum SourceType {
    Nil = '',
    CloudAWS = 'AWS',
    CloudAzure = 'Azure',
}

export interface TypesDatapoint {
    /** The average of the metric values that correspond to the data point. */
    average?: number
    /** The percentile statistic for the data point. */
    extendedStatistics?: Record<string, number>
    /** The maximum metric value for the data point. */
    maximum?: number
    /** The minimum metric value for the data point. */
    minimum?: number
    /**
     * The number of metric values that contributed to the aggregate value of this
     * data point.
     */
    sampleCount?: number
    /** The sum of the metric values for the data point. */
    sum?: number
    /** The time stamp used for the data point. */
    timestamp?: string
    /** The standard unit for the data point. */
    unit?: TypesStandardUnit
}

export enum TypesFindingSeverity {
    FindingSeverityNone = 'none',
    FindingSeverityLow = 'low',
    FindingSeverityMedium = 'medium',
    FindingSeverityHigh = 'high',
    FindingSeverityCritical = 'critical',
}

export enum TypesInstanceLifecycleType {
    InstanceLifecycleTypeSpot = 'spot',
    InstanceLifecycleTypeScheduled = 'scheduled',
    InstanceLifecycleTypeCapacityBlock = 'capacity-block',
}

export enum TypesInstanceStateName {
    InstanceStateNamePending = 'pending',
    InstanceStateNameRunning = 'running',
    InstanceStateNameShuttingDown = 'shutting-down',
    InstanceStateNameTerminated = 'terminated',
    InstanceStateNameStopping = 'stopping',
    InstanceStateNameStopped = 'stopped',
}

export enum TypesInstanceType {
    InstanceTypeA1Medium = 'a1.medium',
    InstanceTypeA1Large = 'a1.large',
    InstanceTypeA1Xlarge = 'a1.xlarge',
    InstanceTypeA12Xlarge = 'a1.2xlarge',
    InstanceTypeA14Xlarge = 'a1.4xlarge',
    InstanceTypeA1Metal = 'a1.metal',
    InstanceTypeC1Medium = 'c1.medium',
    InstanceTypeC1Xlarge = 'c1.xlarge',
    InstanceTypeC3Large = 'c3.large',
    InstanceTypeC3Xlarge = 'c3.xlarge',
    InstanceTypeC32Xlarge = 'c3.2xlarge',
    InstanceTypeC34Xlarge = 'c3.4xlarge',
    InstanceTypeC38Xlarge = 'c3.8xlarge',
    InstanceTypeC4Large = 'c4.large',
    InstanceTypeC4Xlarge = 'c4.xlarge',
    InstanceTypeC42Xlarge = 'c4.2xlarge',
    InstanceTypeC44Xlarge = 'c4.4xlarge',
    InstanceTypeC48Xlarge = 'c4.8xlarge',
    InstanceTypeC5Large = 'c5.large',
    InstanceTypeC5Xlarge = 'c5.xlarge',
    InstanceTypeC52Xlarge = 'c5.2xlarge',
    InstanceTypeC54Xlarge = 'c5.4xlarge',
    InstanceTypeC59Xlarge = 'c5.9xlarge',
    InstanceTypeC512Xlarge = 'c5.12xlarge',
    InstanceTypeC518Xlarge = 'c5.18xlarge',
    InstanceTypeC524Xlarge = 'c5.24xlarge',
    InstanceTypeC5Metal = 'c5.metal',
    InstanceTypeC5ALarge = 'c5a.large',
    InstanceTypeC5AXlarge = 'c5a.xlarge',
    InstanceTypeC5A2Xlarge = 'c5a.2xlarge',
    InstanceTypeC5A4Xlarge = 'c5a.4xlarge',
    InstanceTypeC5A8Xlarge = 'c5a.8xlarge',
    InstanceTypeC5A12Xlarge = 'c5a.12xlarge',
    InstanceTypeC5A16Xlarge = 'c5a.16xlarge',
    InstanceTypeC5A24Xlarge = 'c5a.24xlarge',
    InstanceTypeC5AdLarge = 'c5ad.large',
    InstanceTypeC5AdXlarge = 'c5ad.xlarge',
    InstanceTypeC5Ad2Xlarge = 'c5ad.2xlarge',
    InstanceTypeC5Ad4Xlarge = 'c5ad.4xlarge',
    InstanceTypeC5Ad8Xlarge = 'c5ad.8xlarge',
    InstanceTypeC5Ad12Xlarge = 'c5ad.12xlarge',
    InstanceTypeC5Ad16Xlarge = 'c5ad.16xlarge',
    InstanceTypeC5Ad24Xlarge = 'c5ad.24xlarge',
    InstanceTypeC5DLarge = 'c5d.large',
    InstanceTypeC5DXlarge = 'c5d.xlarge',
    InstanceTypeC5D2Xlarge = 'c5d.2xlarge',
    InstanceTypeC5D4Xlarge = 'c5d.4xlarge',
    InstanceTypeC5D9Xlarge = 'c5d.9xlarge',
    InstanceTypeC5D12Xlarge = 'c5d.12xlarge',
    InstanceTypeC5D18Xlarge = 'c5d.18xlarge',
    InstanceTypeC5D24Xlarge = 'c5d.24xlarge',
    InstanceTypeC5DMetal = 'c5d.metal',
    InstanceTypeC5NLarge = 'c5n.large',
    InstanceTypeC5NXlarge = 'c5n.xlarge',
    InstanceTypeC5N2Xlarge = 'c5n.2xlarge',
    InstanceTypeC5N4Xlarge = 'c5n.4xlarge',
    InstanceTypeC5N9Xlarge = 'c5n.9xlarge',
    InstanceTypeC5N18Xlarge = 'c5n.18xlarge',
    InstanceTypeC5NMetal = 'c5n.metal',
    InstanceTypeC6GMedium = 'c6g.medium',
    InstanceTypeC6GLarge = 'c6g.large',
    InstanceTypeC6GXlarge = 'c6g.xlarge',
    InstanceTypeC6G2Xlarge = 'c6g.2xlarge',
    InstanceTypeC6G4Xlarge = 'c6g.4xlarge',
    InstanceTypeC6G8Xlarge = 'c6g.8xlarge',
    InstanceTypeC6G12Xlarge = 'c6g.12xlarge',
    InstanceTypeC6G16Xlarge = 'c6g.16xlarge',
    InstanceTypeC6GMetal = 'c6g.metal',
    InstanceTypeC6GdMedium = 'c6gd.medium',
    InstanceTypeC6GdLarge = 'c6gd.large',
    InstanceTypeC6GdXlarge = 'c6gd.xlarge',
    InstanceTypeC6Gd2Xlarge = 'c6gd.2xlarge',
    InstanceTypeC6Gd4Xlarge = 'c6gd.4xlarge',
    InstanceTypeC6Gd8Xlarge = 'c6gd.8xlarge',
    InstanceTypeC6Gd12Xlarge = 'c6gd.12xlarge',
    InstanceTypeC6Gd16Xlarge = 'c6gd.16xlarge',
    InstanceTypeC6GdMetal = 'c6gd.metal',
    InstanceTypeC6GnMedium = 'c6gn.medium',
    InstanceTypeC6GnLarge = 'c6gn.large',
    InstanceTypeC6GnXlarge = 'c6gn.xlarge',
    InstanceTypeC6Gn2Xlarge = 'c6gn.2xlarge',
    InstanceTypeC6Gn4Xlarge = 'c6gn.4xlarge',
    InstanceTypeC6Gn8Xlarge = 'c6gn.8xlarge',
    InstanceTypeC6Gn12Xlarge = 'c6gn.12xlarge',
    InstanceTypeC6Gn16Xlarge = 'c6gn.16xlarge',
    InstanceTypeC6ILarge = 'c6i.large',
    InstanceTypeC6IXlarge = 'c6i.xlarge',
    InstanceTypeC6I2Xlarge = 'c6i.2xlarge',
    InstanceTypeC6I4Xlarge = 'c6i.4xlarge',
    InstanceTypeC6I8Xlarge = 'c6i.8xlarge',
    InstanceTypeC6I12Xlarge = 'c6i.12xlarge',
    InstanceTypeC6I16Xlarge = 'c6i.16xlarge',
    InstanceTypeC6I24Xlarge = 'c6i.24xlarge',
    InstanceTypeC6I32Xlarge = 'c6i.32xlarge',
    InstanceTypeC6IMetal = 'c6i.metal',
    InstanceTypeCc14Xlarge = 'cc1.4xlarge',
    InstanceTypeCc28Xlarge = 'cc2.8xlarge',
    InstanceTypeCg14Xlarge = 'cg1.4xlarge',
    InstanceTypeCr18Xlarge = 'cr1.8xlarge',
    InstanceTypeD2Xlarge = 'd2.xlarge',
    InstanceTypeD22Xlarge = 'd2.2xlarge',
    InstanceTypeD24Xlarge = 'd2.4xlarge',
    InstanceTypeD28Xlarge = 'd2.8xlarge',
    InstanceTypeD3Xlarge = 'd3.xlarge',
    InstanceTypeD32Xlarge = 'd3.2xlarge',
    InstanceTypeD34Xlarge = 'd3.4xlarge',
    InstanceTypeD38Xlarge = 'd3.8xlarge',
    InstanceTypeD3EnXlarge = 'd3en.xlarge',
    InstanceTypeD3En2Xlarge = 'd3en.2xlarge',
    InstanceTypeD3En4Xlarge = 'd3en.4xlarge',
    InstanceTypeD3En6Xlarge = 'd3en.6xlarge',
    InstanceTypeD3En8Xlarge = 'd3en.8xlarge',
    InstanceTypeD3En12Xlarge = 'd3en.12xlarge',
    InstanceTypeDl124Xlarge = 'dl1.24xlarge',
    InstanceTypeF12Xlarge = 'f1.2xlarge',
    InstanceTypeF14Xlarge = 'f1.4xlarge',
    InstanceTypeF116Xlarge = 'f1.16xlarge',
    InstanceTypeG22Xlarge = 'g2.2xlarge',
    InstanceTypeG28Xlarge = 'g2.8xlarge',
    InstanceTypeG34Xlarge = 'g3.4xlarge',
    InstanceTypeG38Xlarge = 'g3.8xlarge',
    InstanceTypeG316Xlarge = 'g3.16xlarge',
    InstanceTypeG3SXlarge = 'g3s.xlarge',
    InstanceTypeG4AdXlarge = 'g4ad.xlarge',
    InstanceTypeG4Ad2Xlarge = 'g4ad.2xlarge',
    InstanceTypeG4Ad4Xlarge = 'g4ad.4xlarge',
    InstanceTypeG4Ad8Xlarge = 'g4ad.8xlarge',
    InstanceTypeG4Ad16Xlarge = 'g4ad.16xlarge',
    InstanceTypeG4DnXlarge = 'g4dn.xlarge',
    InstanceTypeG4Dn2Xlarge = 'g4dn.2xlarge',
    InstanceTypeG4Dn4Xlarge = 'g4dn.4xlarge',
    InstanceTypeG4Dn8Xlarge = 'g4dn.8xlarge',
    InstanceTypeG4Dn12Xlarge = 'g4dn.12xlarge',
    InstanceTypeG4Dn16Xlarge = 'g4dn.16xlarge',
    InstanceTypeG4DnMetal = 'g4dn.metal',
    InstanceTypeG5Xlarge = 'g5.xlarge',
    InstanceTypeG52Xlarge = 'g5.2xlarge',
    InstanceTypeG54Xlarge = 'g5.4xlarge',
    InstanceTypeG58Xlarge = 'g5.8xlarge',
    InstanceTypeG512Xlarge = 'g5.12xlarge',
    InstanceTypeG516Xlarge = 'g5.16xlarge',
    InstanceTypeG524Xlarge = 'g5.24xlarge',
    InstanceTypeG548Xlarge = 'g5.48xlarge',
    InstanceTypeG5GXlarge = 'g5g.xlarge',
    InstanceTypeG5G2Xlarge = 'g5g.2xlarge',
    InstanceTypeG5G4Xlarge = 'g5g.4xlarge',
    InstanceTypeG5G8Xlarge = 'g5g.8xlarge',
    InstanceTypeG5G16Xlarge = 'g5g.16xlarge',
    InstanceTypeG5GMetal = 'g5g.metal',
    InstanceTypeHi14Xlarge = 'hi1.4xlarge',
    InstanceTypeHpc6A48Xlarge = 'hpc6a.48xlarge',
    InstanceTypeHs18Xlarge = 'hs1.8xlarge',
    InstanceTypeH12Xlarge = 'h1.2xlarge',
    InstanceTypeH14Xlarge = 'h1.4xlarge',
    InstanceTypeH18Xlarge = 'h1.8xlarge',
    InstanceTypeH116Xlarge = 'h1.16xlarge',
    InstanceTypeI2Xlarge = 'i2.xlarge',
    InstanceTypeI22Xlarge = 'i2.2xlarge',
    InstanceTypeI24Xlarge = 'i2.4xlarge',
    InstanceTypeI28Xlarge = 'i2.8xlarge',
    InstanceTypeI3Large = 'i3.large',
    InstanceTypeI3Xlarge = 'i3.xlarge',
    InstanceTypeI32Xlarge = 'i3.2xlarge',
    InstanceTypeI34Xlarge = 'i3.4xlarge',
    InstanceTypeI38Xlarge = 'i3.8xlarge',
    InstanceTypeI316Xlarge = 'i3.16xlarge',
    InstanceTypeI3Metal = 'i3.metal',
    InstanceTypeI3EnLarge = 'i3en.large',
    InstanceTypeI3EnXlarge = 'i3en.xlarge',
    InstanceTypeI3En2Xlarge = 'i3en.2xlarge',
    InstanceTypeI3En3Xlarge = 'i3en.3xlarge',
    InstanceTypeI3En6Xlarge = 'i3en.6xlarge',
    InstanceTypeI3En12Xlarge = 'i3en.12xlarge',
    InstanceTypeI3En24Xlarge = 'i3en.24xlarge',
    InstanceTypeI3EnMetal = 'i3en.metal',
    InstanceTypeIm4GnLarge = 'im4gn.large',
    InstanceTypeIm4GnXlarge = 'im4gn.xlarge',
    InstanceTypeIm4Gn2Xlarge = 'im4gn.2xlarge',
    InstanceTypeIm4Gn4Xlarge = 'im4gn.4xlarge',
    InstanceTypeIm4Gn8Xlarge = 'im4gn.8xlarge',
    InstanceTypeIm4Gn16Xlarge = 'im4gn.16xlarge',
    InstanceTypeInf1Xlarge = 'inf1.xlarge',
    InstanceTypeInf12Xlarge = 'inf1.2xlarge',
    InstanceTypeInf16Xlarge = 'inf1.6xlarge',
    InstanceTypeInf124Xlarge = 'inf1.24xlarge',
    InstanceTypeIs4GenMedium = 'is4gen.medium',
    InstanceTypeIs4GenLarge = 'is4gen.large',
    InstanceTypeIs4GenXlarge = 'is4gen.xlarge',
    InstanceTypeIs4Gen2Xlarge = 'is4gen.2xlarge',
    InstanceTypeIs4Gen4Xlarge = 'is4gen.4xlarge',
    InstanceTypeIs4Gen8Xlarge = 'is4gen.8xlarge',
    InstanceTypeM1Small = 'm1.small',
    InstanceTypeM1Medium = 'm1.medium',
    InstanceTypeM1Large = 'm1.large',
    InstanceTypeM1Xlarge = 'm1.xlarge',
    InstanceTypeM2Xlarge = 'm2.xlarge',
    InstanceTypeM22Xlarge = 'm2.2xlarge',
    InstanceTypeM24Xlarge = 'm2.4xlarge',
    InstanceTypeM3Medium = 'm3.medium',
    InstanceTypeM3Large = 'm3.large',
    InstanceTypeM3Xlarge = 'm3.xlarge',
    InstanceTypeM32Xlarge = 'm3.2xlarge',
    InstanceTypeM4Large = 'm4.large',
    InstanceTypeM4Xlarge = 'm4.xlarge',
    InstanceTypeM42Xlarge = 'm4.2xlarge',
    InstanceTypeM44Xlarge = 'm4.4xlarge',
    InstanceTypeM410Xlarge = 'm4.10xlarge',
    InstanceTypeM416Xlarge = 'm4.16xlarge',
    InstanceTypeM5Large = 'm5.large',
    InstanceTypeM5Xlarge = 'm5.xlarge',
    InstanceTypeM52Xlarge = 'm5.2xlarge',
    InstanceTypeM54Xlarge = 'm5.4xlarge',
    InstanceTypeM58Xlarge = 'm5.8xlarge',
    InstanceTypeM512Xlarge = 'm5.12xlarge',
    InstanceTypeM516Xlarge = 'm5.16xlarge',
    InstanceTypeM524Xlarge = 'm5.24xlarge',
    InstanceTypeM5Metal = 'm5.metal',
    InstanceTypeM5ALarge = 'm5a.large',
    InstanceTypeM5AXlarge = 'm5a.xlarge',
    InstanceTypeM5A2Xlarge = 'm5a.2xlarge',
    InstanceTypeM5A4Xlarge = 'm5a.4xlarge',
    InstanceTypeM5A8Xlarge = 'm5a.8xlarge',
    InstanceTypeM5A12Xlarge = 'm5a.12xlarge',
    InstanceTypeM5A16Xlarge = 'm5a.16xlarge',
    InstanceTypeM5A24Xlarge = 'm5a.24xlarge',
    InstanceTypeM5AdLarge = 'm5ad.large',
    InstanceTypeM5AdXlarge = 'm5ad.xlarge',
    InstanceTypeM5Ad2Xlarge = 'm5ad.2xlarge',
    InstanceTypeM5Ad4Xlarge = 'm5ad.4xlarge',
    InstanceTypeM5Ad8Xlarge = 'm5ad.8xlarge',
    InstanceTypeM5Ad12Xlarge = 'm5ad.12xlarge',
    InstanceTypeM5Ad16Xlarge = 'm5ad.16xlarge',
    InstanceTypeM5Ad24Xlarge = 'm5ad.24xlarge',
    InstanceTypeM5DLarge = 'm5d.large',
    InstanceTypeM5DXlarge = 'm5d.xlarge',
    InstanceTypeM5D2Xlarge = 'm5d.2xlarge',
    InstanceTypeM5D4Xlarge = 'm5d.4xlarge',
    InstanceTypeM5D8Xlarge = 'm5d.8xlarge',
    InstanceTypeM5D12Xlarge = 'm5d.12xlarge',
    InstanceTypeM5D16Xlarge = 'm5d.16xlarge',
    InstanceTypeM5D24Xlarge = 'm5d.24xlarge',
    InstanceTypeM5DMetal = 'm5d.metal',
    InstanceTypeM5DnLarge = 'm5dn.large',
    InstanceTypeM5DnXlarge = 'm5dn.xlarge',
    InstanceTypeM5Dn2Xlarge = 'm5dn.2xlarge',
    InstanceTypeM5Dn4Xlarge = 'm5dn.4xlarge',
    InstanceTypeM5Dn8Xlarge = 'm5dn.8xlarge',
    InstanceTypeM5Dn12Xlarge = 'm5dn.12xlarge',
    InstanceTypeM5Dn16Xlarge = 'm5dn.16xlarge',
    InstanceTypeM5Dn24Xlarge = 'm5dn.24xlarge',
    InstanceTypeM5DnMetal = 'm5dn.metal',
    InstanceTypeM5NLarge = 'm5n.large',
    InstanceTypeM5NXlarge = 'm5n.xlarge',
    InstanceTypeM5N2Xlarge = 'm5n.2xlarge',
    InstanceTypeM5N4Xlarge = 'm5n.4xlarge',
    InstanceTypeM5N8Xlarge = 'm5n.8xlarge',
    InstanceTypeM5N12Xlarge = 'm5n.12xlarge',
    InstanceTypeM5N16Xlarge = 'm5n.16xlarge',
    InstanceTypeM5N24Xlarge = 'm5n.24xlarge',
    InstanceTypeM5NMetal = 'm5n.metal',
    InstanceTypeM5ZnLarge = 'm5zn.large',
    InstanceTypeM5ZnXlarge = 'm5zn.xlarge',
    InstanceTypeM5Zn2Xlarge = 'm5zn.2xlarge',
    InstanceTypeM5Zn3Xlarge = 'm5zn.3xlarge',
    InstanceTypeM5Zn6Xlarge = 'm5zn.6xlarge',
    InstanceTypeM5Zn12Xlarge = 'm5zn.12xlarge',
    InstanceTypeM5ZnMetal = 'm5zn.metal',
    InstanceTypeM6ALarge = 'm6a.large',
    InstanceTypeM6AXlarge = 'm6a.xlarge',
    InstanceTypeM6A2Xlarge = 'm6a.2xlarge',
    InstanceTypeM6A4Xlarge = 'm6a.4xlarge',
    InstanceTypeM6A8Xlarge = 'm6a.8xlarge',
    InstanceTypeM6A12Xlarge = 'm6a.12xlarge',
    InstanceTypeM6A16Xlarge = 'm6a.16xlarge',
    InstanceTypeM6A24Xlarge = 'm6a.24xlarge',
    InstanceTypeM6A32Xlarge = 'm6a.32xlarge',
    InstanceTypeM6A48Xlarge = 'm6a.48xlarge',
    InstanceTypeM6GMetal = 'm6g.metal',
    InstanceTypeM6GMedium = 'm6g.medium',
    InstanceTypeM6GLarge = 'm6g.large',
    InstanceTypeM6GXlarge = 'm6g.xlarge',
    InstanceTypeM6G2Xlarge = 'm6g.2xlarge',
    InstanceTypeM6G4Xlarge = 'm6g.4xlarge',
    InstanceTypeM6G8Xlarge = 'm6g.8xlarge',
    InstanceTypeM6G12Xlarge = 'm6g.12xlarge',
    InstanceTypeM6G16Xlarge = 'm6g.16xlarge',
    InstanceTypeM6GdMetal = 'm6gd.metal',
    InstanceTypeM6GdMedium = 'm6gd.medium',
    InstanceTypeM6GdLarge = 'm6gd.large',
    InstanceTypeM6GdXlarge = 'm6gd.xlarge',
    InstanceTypeM6Gd2Xlarge = 'm6gd.2xlarge',
    InstanceTypeM6Gd4Xlarge = 'm6gd.4xlarge',
    InstanceTypeM6Gd8Xlarge = 'm6gd.8xlarge',
    InstanceTypeM6Gd12Xlarge = 'm6gd.12xlarge',
    InstanceTypeM6Gd16Xlarge = 'm6gd.16xlarge',
    InstanceTypeM6ILarge = 'm6i.large',
    InstanceTypeM6IXlarge = 'm6i.xlarge',
    InstanceTypeM6I2Xlarge = 'm6i.2xlarge',
    InstanceTypeM6I4Xlarge = 'm6i.4xlarge',
    InstanceTypeM6I8Xlarge = 'm6i.8xlarge',
    InstanceTypeM6I12Xlarge = 'm6i.12xlarge',
    InstanceTypeM6I16Xlarge = 'm6i.16xlarge',
    InstanceTypeM6I24Xlarge = 'm6i.24xlarge',
    InstanceTypeM6I32Xlarge = 'm6i.32xlarge',
    InstanceTypeM6IMetal = 'm6i.metal',
    InstanceTypeMac1Metal = 'mac1.metal',
    InstanceTypeP2Xlarge = 'p2.xlarge',
    InstanceTypeP28Xlarge = 'p2.8xlarge',
    InstanceTypeP216Xlarge = 'p2.16xlarge',
    InstanceTypeP32Xlarge = 'p3.2xlarge',
    InstanceTypeP38Xlarge = 'p3.8xlarge',
    InstanceTypeP316Xlarge = 'p3.16xlarge',
    InstanceTypeP3Dn24Xlarge = 'p3dn.24xlarge',
    InstanceTypeP4D24Xlarge = 'p4d.24xlarge',
    InstanceTypeR3Large = 'r3.large',
    InstanceTypeR3Xlarge = 'r3.xlarge',
    InstanceTypeR32Xlarge = 'r3.2xlarge',
    InstanceTypeR34Xlarge = 'r3.4xlarge',
    InstanceTypeR38Xlarge = 'r3.8xlarge',
    InstanceTypeR4Large = 'r4.large',
    InstanceTypeR4Xlarge = 'r4.xlarge',
    InstanceTypeR42Xlarge = 'r4.2xlarge',
    InstanceTypeR44Xlarge = 'r4.4xlarge',
    InstanceTypeR48Xlarge = 'r4.8xlarge',
    InstanceTypeR416Xlarge = 'r4.16xlarge',
    InstanceTypeR5Large = 'r5.large',
    InstanceTypeR5Xlarge = 'r5.xlarge',
    InstanceTypeR52Xlarge = 'r5.2xlarge',
    InstanceTypeR54Xlarge = 'r5.4xlarge',
    InstanceTypeR58Xlarge = 'r5.8xlarge',
    InstanceTypeR512Xlarge = 'r5.12xlarge',
    InstanceTypeR516Xlarge = 'r5.16xlarge',
    InstanceTypeR524Xlarge = 'r5.24xlarge',
    InstanceTypeR5Metal = 'r5.metal',
    InstanceTypeR5ALarge = 'r5a.large',
    InstanceTypeR5AXlarge = 'r5a.xlarge',
    InstanceTypeR5A2Xlarge = 'r5a.2xlarge',
    InstanceTypeR5A4Xlarge = 'r5a.4xlarge',
    InstanceTypeR5A8Xlarge = 'r5a.8xlarge',
    InstanceTypeR5A12Xlarge = 'r5a.12xlarge',
    InstanceTypeR5A16Xlarge = 'r5a.16xlarge',
    InstanceTypeR5A24Xlarge = 'r5a.24xlarge',
    InstanceTypeR5AdLarge = 'r5ad.large',
    InstanceTypeR5AdXlarge = 'r5ad.xlarge',
    InstanceTypeR5Ad2Xlarge = 'r5ad.2xlarge',
    InstanceTypeR5Ad4Xlarge = 'r5ad.4xlarge',
    InstanceTypeR5Ad8Xlarge = 'r5ad.8xlarge',
    InstanceTypeR5Ad12Xlarge = 'r5ad.12xlarge',
    InstanceTypeR5Ad16Xlarge = 'r5ad.16xlarge',
    InstanceTypeR5Ad24Xlarge = 'r5ad.24xlarge',
    InstanceTypeR5BLarge = 'r5b.large',
    InstanceTypeR5BXlarge = 'r5b.xlarge',
    InstanceTypeR5B2Xlarge = 'r5b.2xlarge',
    InstanceTypeR5B4Xlarge = 'r5b.4xlarge',
    InstanceTypeR5B8Xlarge = 'r5b.8xlarge',
    InstanceTypeR5B12Xlarge = 'r5b.12xlarge',
    InstanceTypeR5B16Xlarge = 'r5b.16xlarge',
    InstanceTypeR5B24Xlarge = 'r5b.24xlarge',
    InstanceTypeR5BMetal = 'r5b.metal',
    InstanceTypeR5DLarge = 'r5d.large',
    InstanceTypeR5DXlarge = 'r5d.xlarge',
    InstanceTypeR5D2Xlarge = 'r5d.2xlarge',
    InstanceTypeR5D4Xlarge = 'r5d.4xlarge',
    InstanceTypeR5D8Xlarge = 'r5d.8xlarge',
    InstanceTypeR5D12Xlarge = 'r5d.12xlarge',
    InstanceTypeR5D16Xlarge = 'r5d.16xlarge',
    InstanceTypeR5D24Xlarge = 'r5d.24xlarge',
    InstanceTypeR5DMetal = 'r5d.metal',
    InstanceTypeR5DnLarge = 'r5dn.large',
    InstanceTypeR5DnXlarge = 'r5dn.xlarge',
    InstanceTypeR5Dn2Xlarge = 'r5dn.2xlarge',
    InstanceTypeR5Dn4Xlarge = 'r5dn.4xlarge',
    InstanceTypeR5Dn8Xlarge = 'r5dn.8xlarge',
    InstanceTypeR5Dn12Xlarge = 'r5dn.12xlarge',
    InstanceTypeR5Dn16Xlarge = 'r5dn.16xlarge',
    InstanceTypeR5Dn24Xlarge = 'r5dn.24xlarge',
    InstanceTypeR5DnMetal = 'r5dn.metal',
    InstanceTypeR5NLarge = 'r5n.large',
    InstanceTypeR5NXlarge = 'r5n.xlarge',
    InstanceTypeR5N2Xlarge = 'r5n.2xlarge',
    InstanceTypeR5N4Xlarge = 'r5n.4xlarge',
    InstanceTypeR5N8Xlarge = 'r5n.8xlarge',
    InstanceTypeR5N12Xlarge = 'r5n.12xlarge',
    InstanceTypeR5N16Xlarge = 'r5n.16xlarge',
    InstanceTypeR5N24Xlarge = 'r5n.24xlarge',
    InstanceTypeR5NMetal = 'r5n.metal',
    InstanceTypeR6GMedium = 'r6g.medium',
    InstanceTypeR6GLarge = 'r6g.large',
    InstanceTypeR6GXlarge = 'r6g.xlarge',
    InstanceTypeR6G2Xlarge = 'r6g.2xlarge',
    InstanceTypeR6G4Xlarge = 'r6g.4xlarge',
    InstanceTypeR6G8Xlarge = 'r6g.8xlarge',
    InstanceTypeR6G12Xlarge = 'r6g.12xlarge',
    InstanceTypeR6G16Xlarge = 'r6g.16xlarge',
    InstanceTypeR6GMetal = 'r6g.metal',
    InstanceTypeR6GdMedium = 'r6gd.medium',
    InstanceTypeR6GdLarge = 'r6gd.large',
    InstanceTypeR6GdXlarge = 'r6gd.xlarge',
    InstanceTypeR6Gd2Xlarge = 'r6gd.2xlarge',
    InstanceTypeR6Gd4Xlarge = 'r6gd.4xlarge',
    InstanceTypeR6Gd8Xlarge = 'r6gd.8xlarge',
    InstanceTypeR6Gd12Xlarge = 'r6gd.12xlarge',
    InstanceTypeR6Gd16Xlarge = 'r6gd.16xlarge',
    InstanceTypeR6GdMetal = 'r6gd.metal',
    InstanceTypeR6ILarge = 'r6i.large',
    InstanceTypeR6IXlarge = 'r6i.xlarge',
    InstanceTypeR6I2Xlarge = 'r6i.2xlarge',
    InstanceTypeR6I4Xlarge = 'r6i.4xlarge',
    InstanceTypeR6I8Xlarge = 'r6i.8xlarge',
    InstanceTypeR6I12Xlarge = 'r6i.12xlarge',
    InstanceTypeR6I16Xlarge = 'r6i.16xlarge',
    InstanceTypeR6I24Xlarge = 'r6i.24xlarge',
    InstanceTypeR6I32Xlarge = 'r6i.32xlarge',
    InstanceTypeR6IMetal = 'r6i.metal',
    InstanceTypeT1Micro = 't1.micro',
    InstanceTypeT2Nano = 't2.nano',
    InstanceTypeT2Micro = 't2.micro',
    InstanceTypeT2Small = 't2.small',
    InstanceTypeT2Medium = 't2.medium',
    InstanceTypeT2Large = 't2.large',
    InstanceTypeT2Xlarge = 't2.xlarge',
    InstanceTypeT22Xlarge = 't2.2xlarge',
    InstanceTypeT3Nano = 't3.nano',
    InstanceTypeT3Micro = 't3.micro',
    InstanceTypeT3Small = 't3.small',
    InstanceTypeT3Medium = 't3.medium',
    InstanceTypeT3Large = 't3.large',
    InstanceTypeT3Xlarge = 't3.xlarge',
    InstanceTypeT32Xlarge = 't3.2xlarge',
    InstanceTypeT3ANano = 't3a.nano',
    InstanceTypeT3AMicro = 't3a.micro',
    InstanceTypeT3ASmall = 't3a.small',
    InstanceTypeT3AMedium = 't3a.medium',
    InstanceTypeT3ALarge = 't3a.large',
    InstanceTypeT3AXlarge = 't3a.xlarge',
    InstanceTypeT3A2Xlarge = 't3a.2xlarge',
    InstanceTypeT4GNano = 't4g.nano',
    InstanceTypeT4GMicro = 't4g.micro',
    InstanceTypeT4GSmall = 't4g.small',
    InstanceTypeT4GMedium = 't4g.medium',
    InstanceTypeT4GLarge = 't4g.large',
    InstanceTypeT4GXlarge = 't4g.xlarge',
    InstanceTypeT4G2Xlarge = 't4g.2xlarge',
    InstanceTypeU6Tb156Xlarge = 'u-6tb1.56xlarge',
    InstanceTypeU6Tb1112Xlarge = 'u-6tb1.112xlarge',
    InstanceTypeU9Tb1112Xlarge = 'u-9tb1.112xlarge',
    InstanceTypeU12Tb1112Xlarge = 'u-12tb1.112xlarge',
    InstanceTypeU6Tb1Metal = 'u-6tb1.metal',
    InstanceTypeU9Tb1Metal = 'u-9tb1.metal',
    InstanceTypeU12Tb1Metal = 'u-12tb1.metal',
    InstanceTypeU18Tb1Metal = 'u-18tb1.metal',
    InstanceTypeU24Tb1Metal = 'u-24tb1.metal',
    InstanceTypeVt13Xlarge = 'vt1.3xlarge',
    InstanceTypeVt16Xlarge = 'vt1.6xlarge',
    InstanceTypeVt124Xlarge = 'vt1.24xlarge',
    InstanceTypeX116Xlarge = 'x1.16xlarge',
    InstanceTypeX132Xlarge = 'x1.32xlarge',
    InstanceTypeX1EXlarge = 'x1e.xlarge',
    InstanceTypeX1E2Xlarge = 'x1e.2xlarge',
    InstanceTypeX1E4Xlarge = 'x1e.4xlarge',
    InstanceTypeX1E8Xlarge = 'x1e.8xlarge',
    InstanceTypeX1E16Xlarge = 'x1e.16xlarge',
    InstanceTypeX1E32Xlarge = 'x1e.32xlarge',
    InstanceTypeX2Iezn2Xlarge = 'x2iezn.2xlarge',
    InstanceTypeX2Iezn4Xlarge = 'x2iezn.4xlarge',
    InstanceTypeX2Iezn6Xlarge = 'x2iezn.6xlarge',
    InstanceTypeX2Iezn8Xlarge = 'x2iezn.8xlarge',
    InstanceTypeX2Iezn12Xlarge = 'x2iezn.12xlarge',
    InstanceTypeX2IeznMetal = 'x2iezn.metal',
    InstanceTypeX2GdMedium = 'x2gd.medium',
    InstanceTypeX2GdLarge = 'x2gd.large',
    InstanceTypeX2GdXlarge = 'x2gd.xlarge',
    InstanceTypeX2Gd2Xlarge = 'x2gd.2xlarge',
    InstanceTypeX2Gd4Xlarge = 'x2gd.4xlarge',
    InstanceTypeX2Gd8Xlarge = 'x2gd.8xlarge',
    InstanceTypeX2Gd12Xlarge = 'x2gd.12xlarge',
    InstanceTypeX2Gd16Xlarge = 'x2gd.16xlarge',
    InstanceTypeX2GdMetal = 'x2gd.metal',
    InstanceTypeZ1DLarge = 'z1d.large',
    InstanceTypeZ1DXlarge = 'z1d.xlarge',
    InstanceTypeZ1D2Xlarge = 'z1d.2xlarge',
    InstanceTypeZ1D3Xlarge = 'z1d.3xlarge',
    InstanceTypeZ1D6Xlarge = 'z1d.6xlarge',
    InstanceTypeZ1D12Xlarge = 'z1d.12xlarge',
    InstanceTypeZ1DMetal = 'z1d.metal',
    InstanceTypeX2Idn16Xlarge = 'x2idn.16xlarge',
    InstanceTypeX2Idn24Xlarge = 'x2idn.24xlarge',
    InstanceTypeX2Idn32Xlarge = 'x2idn.32xlarge',
    InstanceTypeX2IednXlarge = 'x2iedn.xlarge',
    InstanceTypeX2Iedn2Xlarge = 'x2iedn.2xlarge',
    InstanceTypeX2Iedn4Xlarge = 'x2iedn.4xlarge',
    InstanceTypeX2Iedn8Xlarge = 'x2iedn.8xlarge',
    InstanceTypeX2Iedn16Xlarge = 'x2iedn.16xlarge',
    InstanceTypeX2Iedn24Xlarge = 'x2iedn.24xlarge',
    InstanceTypeX2Iedn32Xlarge = 'x2iedn.32xlarge',
    InstanceTypeC6ALarge = 'c6a.large',
    InstanceTypeC6AXlarge = 'c6a.xlarge',
    InstanceTypeC6A2Xlarge = 'c6a.2xlarge',
    InstanceTypeC6A4Xlarge = 'c6a.4xlarge',
    InstanceTypeC6A8Xlarge = 'c6a.8xlarge',
    InstanceTypeC6A12Xlarge = 'c6a.12xlarge',
    InstanceTypeC6A16Xlarge = 'c6a.16xlarge',
    InstanceTypeC6A24Xlarge = 'c6a.24xlarge',
    InstanceTypeC6A32Xlarge = 'c6a.32xlarge',
    InstanceTypeC6A48Xlarge = 'c6a.48xlarge',
    InstanceTypeC6AMetal = 'c6a.metal',
    InstanceTypeM6AMetal = 'm6a.metal',
    InstanceTypeI4ILarge = 'i4i.large',
    InstanceTypeI4IXlarge = 'i4i.xlarge',
    InstanceTypeI4I2Xlarge = 'i4i.2xlarge',
    InstanceTypeI4I4Xlarge = 'i4i.4xlarge',
    InstanceTypeI4I8Xlarge = 'i4i.8xlarge',
    InstanceTypeI4I16Xlarge = 'i4i.16xlarge',
    InstanceTypeI4I32Xlarge = 'i4i.32xlarge',
    InstanceTypeI4IMetal = 'i4i.metal',
    InstanceTypeX2IdnMetal = 'x2idn.metal',
    InstanceTypeX2IednMetal = 'x2iedn.metal',
    InstanceTypeC7GMedium = 'c7g.medium',
    InstanceTypeC7GLarge = 'c7g.large',
    InstanceTypeC7GXlarge = 'c7g.xlarge',
    InstanceTypeC7G2Xlarge = 'c7g.2xlarge',
    InstanceTypeC7G4Xlarge = 'c7g.4xlarge',
    InstanceTypeC7G8Xlarge = 'c7g.8xlarge',
    InstanceTypeC7G12Xlarge = 'c7g.12xlarge',
    InstanceTypeC7G16Xlarge = 'c7g.16xlarge',
    InstanceTypeMac2Metal = 'mac2.metal',
    InstanceTypeC6IdLarge = 'c6id.large',
    InstanceTypeC6IdXlarge = 'c6id.xlarge',
    InstanceTypeC6Id2Xlarge = 'c6id.2xlarge',
    InstanceTypeC6Id4Xlarge = 'c6id.4xlarge',
    InstanceTypeC6Id8Xlarge = 'c6id.8xlarge',
    InstanceTypeC6Id12Xlarge = 'c6id.12xlarge',
    InstanceTypeC6Id16Xlarge = 'c6id.16xlarge',
    InstanceTypeC6Id24Xlarge = 'c6id.24xlarge',
    InstanceTypeC6Id32Xlarge = 'c6id.32xlarge',
    InstanceTypeC6IdMetal = 'c6id.metal',
    InstanceTypeM6IdLarge = 'm6id.large',
    InstanceTypeM6IdXlarge = 'm6id.xlarge',
    InstanceTypeM6Id2Xlarge = 'm6id.2xlarge',
    InstanceTypeM6Id4Xlarge = 'm6id.4xlarge',
    InstanceTypeM6Id8Xlarge = 'm6id.8xlarge',
    InstanceTypeM6Id12Xlarge = 'm6id.12xlarge',
    InstanceTypeM6Id16Xlarge = 'm6id.16xlarge',
    InstanceTypeM6Id24Xlarge = 'm6id.24xlarge',
    InstanceTypeM6Id32Xlarge = 'm6id.32xlarge',
    InstanceTypeM6IdMetal = 'm6id.metal',
    InstanceTypeR6IdLarge = 'r6id.large',
    InstanceTypeR6IdXlarge = 'r6id.xlarge',
    InstanceTypeR6Id2Xlarge = 'r6id.2xlarge',
    InstanceTypeR6Id4Xlarge = 'r6id.4xlarge',
    InstanceTypeR6Id8Xlarge = 'r6id.8xlarge',
    InstanceTypeR6Id12Xlarge = 'r6id.12xlarge',
    InstanceTypeR6Id16Xlarge = 'r6id.16xlarge',
    InstanceTypeR6Id24Xlarge = 'r6id.24xlarge',
    InstanceTypeR6Id32Xlarge = 'r6id.32xlarge',
    InstanceTypeR6IdMetal = 'r6id.metal',
    InstanceTypeR6ALarge = 'r6a.large',
    InstanceTypeR6AXlarge = 'r6a.xlarge',
    InstanceTypeR6A2Xlarge = 'r6a.2xlarge',
    InstanceTypeR6A4Xlarge = 'r6a.4xlarge',
    InstanceTypeR6A8Xlarge = 'r6a.8xlarge',
    InstanceTypeR6A12Xlarge = 'r6a.12xlarge',
    InstanceTypeR6A16Xlarge = 'r6a.16xlarge',
    InstanceTypeR6A24Xlarge = 'r6a.24xlarge',
    InstanceTypeR6A32Xlarge = 'r6a.32xlarge',
    InstanceTypeR6A48Xlarge = 'r6a.48xlarge',
    InstanceTypeR6AMetal = 'r6a.metal',
    InstanceTypeP4De24Xlarge = 'p4de.24xlarge',
    InstanceTypeU3Tb156Xlarge = 'u-3tb1.56xlarge',
    InstanceTypeU18Tb1112Xlarge = 'u-18tb1.112xlarge',
    InstanceTypeU24Tb1112Xlarge = 'u-24tb1.112xlarge',
    InstanceTypeTrn12Xlarge = 'trn1.2xlarge',
    InstanceTypeTrn132Xlarge = 'trn1.32xlarge',
    InstanceTypeHpc6Id32Xlarge = 'hpc6id.32xlarge',
    InstanceTypeC6InLarge = 'c6in.large',
    InstanceTypeC6InXlarge = 'c6in.xlarge',
    InstanceTypeC6In2Xlarge = 'c6in.2xlarge',
    InstanceTypeC6In4Xlarge = 'c6in.4xlarge',
    InstanceTypeC6In8Xlarge = 'c6in.8xlarge',
    InstanceTypeC6In12Xlarge = 'c6in.12xlarge',
    InstanceTypeC6In16Xlarge = 'c6in.16xlarge',
    InstanceTypeC6In24Xlarge = 'c6in.24xlarge',
    InstanceTypeC6In32Xlarge = 'c6in.32xlarge',
    InstanceTypeM6InLarge = 'm6in.large',
    InstanceTypeM6InXlarge = 'm6in.xlarge',
    InstanceTypeM6In2Xlarge = 'm6in.2xlarge',
    InstanceTypeM6In4Xlarge = 'm6in.4xlarge',
    InstanceTypeM6In8Xlarge = 'm6in.8xlarge',
    InstanceTypeM6In12Xlarge = 'm6in.12xlarge',
    InstanceTypeM6In16Xlarge = 'm6in.16xlarge',
    InstanceTypeM6In24Xlarge = 'm6in.24xlarge',
    InstanceTypeM6In32Xlarge = 'm6in.32xlarge',
    InstanceTypeM6IdnLarge = 'm6idn.large',
    InstanceTypeM6IdnXlarge = 'm6idn.xlarge',
    InstanceTypeM6Idn2Xlarge = 'm6idn.2xlarge',
    InstanceTypeM6Idn4Xlarge = 'm6idn.4xlarge',
    InstanceTypeM6Idn8Xlarge = 'm6idn.8xlarge',
    InstanceTypeM6Idn12Xlarge = 'm6idn.12xlarge',
    InstanceTypeM6Idn16Xlarge = 'm6idn.16xlarge',
    InstanceTypeM6Idn24Xlarge = 'm6idn.24xlarge',
    InstanceTypeM6Idn32Xlarge = 'm6idn.32xlarge',
    InstanceTypeR6InLarge = 'r6in.large',
    InstanceTypeR6InXlarge = 'r6in.xlarge',
    InstanceTypeR6In2Xlarge = 'r6in.2xlarge',
    InstanceTypeR6In4Xlarge = 'r6in.4xlarge',
    InstanceTypeR6In8Xlarge = 'r6in.8xlarge',
    InstanceTypeR6In12Xlarge = 'r6in.12xlarge',
    InstanceTypeR6In16Xlarge = 'r6in.16xlarge',
    InstanceTypeR6In24Xlarge = 'r6in.24xlarge',
    InstanceTypeR6In32Xlarge = 'r6in.32xlarge',
    InstanceTypeR6IdnLarge = 'r6idn.large',
    InstanceTypeR6IdnXlarge = 'r6idn.xlarge',
    InstanceTypeR6Idn2Xlarge = 'r6idn.2xlarge',
    InstanceTypeR6Idn4Xlarge = 'r6idn.4xlarge',
    InstanceTypeR6Idn8Xlarge = 'r6idn.8xlarge',
    InstanceTypeR6Idn12Xlarge = 'r6idn.12xlarge',
    InstanceTypeR6Idn16Xlarge = 'r6idn.16xlarge',
    InstanceTypeR6Idn24Xlarge = 'r6idn.24xlarge',
    InstanceTypeR6Idn32Xlarge = 'r6idn.32xlarge',
    InstanceTypeC7GMetal = 'c7g.metal',
    InstanceTypeM7GMedium = 'm7g.medium',
    InstanceTypeM7GLarge = 'm7g.large',
    InstanceTypeM7GXlarge = 'm7g.xlarge',
    InstanceTypeM7G2Xlarge = 'm7g.2xlarge',
    InstanceTypeM7G4Xlarge = 'm7g.4xlarge',
    InstanceTypeM7G8Xlarge = 'm7g.8xlarge',
    InstanceTypeM7G12Xlarge = 'm7g.12xlarge',
    InstanceTypeM7G16Xlarge = 'm7g.16xlarge',
    InstanceTypeM7GMetal = 'm7g.metal',
    InstanceTypeR7GMedium = 'r7g.medium',
    InstanceTypeR7GLarge = 'r7g.large',
    InstanceTypeR7GXlarge = 'r7g.xlarge',
    InstanceTypeR7G2Xlarge = 'r7g.2xlarge',
    InstanceTypeR7G4Xlarge = 'r7g.4xlarge',
    InstanceTypeR7G8Xlarge = 'r7g.8xlarge',
    InstanceTypeR7G12Xlarge = 'r7g.12xlarge',
    InstanceTypeR7G16Xlarge = 'r7g.16xlarge',
    InstanceTypeR7GMetal = 'r7g.metal',
    InstanceTypeC6InMetal = 'c6in.metal',
    InstanceTypeM6InMetal = 'm6in.metal',
    InstanceTypeM6IdnMetal = 'm6idn.metal',
    InstanceTypeR6InMetal = 'r6in.metal',
    InstanceTypeR6IdnMetal = 'r6idn.metal',
    InstanceTypeInf2Xlarge = 'inf2.xlarge',
    InstanceTypeInf28Xlarge = 'inf2.8xlarge',
    InstanceTypeInf224Xlarge = 'inf2.24xlarge',
    InstanceTypeInf248Xlarge = 'inf2.48xlarge',
    InstanceTypeTrn1N32Xlarge = 'trn1n.32xlarge',
    InstanceTypeI4GLarge = 'i4g.large',
    InstanceTypeI4GXlarge = 'i4g.xlarge',
    InstanceTypeI4G2Xlarge = 'i4g.2xlarge',
    InstanceTypeI4G4Xlarge = 'i4g.4xlarge',
    InstanceTypeI4G8Xlarge = 'i4g.8xlarge',
    InstanceTypeI4G16Xlarge = 'i4g.16xlarge',
    InstanceTypeHpc7G4Xlarge = 'hpc7g.4xlarge',
    InstanceTypeHpc7G8Xlarge = 'hpc7g.8xlarge',
    InstanceTypeHpc7G16Xlarge = 'hpc7g.16xlarge',
    InstanceTypeC7GnMedium = 'c7gn.medium',
    InstanceTypeC7GnLarge = 'c7gn.large',
    InstanceTypeC7GnXlarge = 'c7gn.xlarge',
    InstanceTypeC7Gn2Xlarge = 'c7gn.2xlarge',
    InstanceTypeC7Gn4Xlarge = 'c7gn.4xlarge',
    InstanceTypeC7Gn8Xlarge = 'c7gn.8xlarge',
    InstanceTypeC7Gn12Xlarge = 'c7gn.12xlarge',
    InstanceTypeC7Gn16Xlarge = 'c7gn.16xlarge',
    InstanceTypeP548Xlarge = 'p5.48xlarge',
    InstanceTypeM7ILarge = 'm7i.large',
    InstanceTypeM7IXlarge = 'm7i.xlarge',
    InstanceTypeM7I2Xlarge = 'm7i.2xlarge',
    InstanceTypeM7I4Xlarge = 'm7i.4xlarge',
    InstanceTypeM7I8Xlarge = 'm7i.8xlarge',
    InstanceTypeM7I12Xlarge = 'm7i.12xlarge',
    InstanceTypeM7I16Xlarge = 'm7i.16xlarge',
    InstanceTypeM7I24Xlarge = 'm7i.24xlarge',
    InstanceTypeM7I48Xlarge = 'm7i.48xlarge',
    InstanceTypeM7IFlexLarge = 'm7i-flex.large',
    InstanceTypeM7IFlexXlarge = 'm7i-flex.xlarge',
    InstanceTypeM7IFlex2Xlarge = 'm7i-flex.2xlarge',
    InstanceTypeM7IFlex4Xlarge = 'm7i-flex.4xlarge',
    InstanceTypeM7IFlex8Xlarge = 'm7i-flex.8xlarge',
    InstanceTypeM7AMedium = 'm7a.medium',
    InstanceTypeM7ALarge = 'm7a.large',
    InstanceTypeM7AXlarge = 'm7a.xlarge',
    InstanceTypeM7A2Xlarge = 'm7a.2xlarge',
    InstanceTypeM7A4Xlarge = 'm7a.4xlarge',
    InstanceTypeM7A8Xlarge = 'm7a.8xlarge',
    InstanceTypeM7A12Xlarge = 'm7a.12xlarge',
    InstanceTypeM7A16Xlarge = 'm7a.16xlarge',
    InstanceTypeM7A24Xlarge = 'm7a.24xlarge',
    InstanceTypeM7A32Xlarge = 'm7a.32xlarge',
    InstanceTypeM7A48Xlarge = 'm7a.48xlarge',
    InstanceTypeM7AMetal48Xl = 'm7a.metal-48xl',
    InstanceTypeHpc7A12Xlarge = 'hpc7a.12xlarge',
    InstanceTypeHpc7A24Xlarge = 'hpc7a.24xlarge',
    InstanceTypeHpc7A48Xlarge = 'hpc7a.48xlarge',
    InstanceTypeHpc7A96Xlarge = 'hpc7a.96xlarge',
    InstanceTypeC7GdMedium = 'c7gd.medium',
    InstanceTypeC7GdLarge = 'c7gd.large',
    InstanceTypeC7GdXlarge = 'c7gd.xlarge',
    InstanceTypeC7Gd2Xlarge = 'c7gd.2xlarge',
    InstanceTypeC7Gd4Xlarge = 'c7gd.4xlarge',
    InstanceTypeC7Gd8Xlarge = 'c7gd.8xlarge',
    InstanceTypeC7Gd12Xlarge = 'c7gd.12xlarge',
    InstanceTypeC7Gd16Xlarge = 'c7gd.16xlarge',
    InstanceTypeM7GdMedium = 'm7gd.medium',
    InstanceTypeM7GdLarge = 'm7gd.large',
    InstanceTypeM7GdXlarge = 'm7gd.xlarge',
    InstanceTypeM7Gd2Xlarge = 'm7gd.2xlarge',
    InstanceTypeM7Gd4Xlarge = 'm7gd.4xlarge',
    InstanceTypeM7Gd8Xlarge = 'm7gd.8xlarge',
    InstanceTypeM7Gd12Xlarge = 'm7gd.12xlarge',
    InstanceTypeM7Gd16Xlarge = 'm7gd.16xlarge',
    InstanceTypeR7GdMedium = 'r7gd.medium',
    InstanceTypeR7GdLarge = 'r7gd.large',
    InstanceTypeR7GdXlarge = 'r7gd.xlarge',
    InstanceTypeR7Gd2Xlarge = 'r7gd.2xlarge',
    InstanceTypeR7Gd4Xlarge = 'r7gd.4xlarge',
    InstanceTypeR7Gd8Xlarge = 'r7gd.8xlarge',
    InstanceTypeR7Gd12Xlarge = 'r7gd.12xlarge',
    InstanceTypeR7Gd16Xlarge = 'r7gd.16xlarge',
    InstanceTypeR7AMedium = 'r7a.medium',
    InstanceTypeR7ALarge = 'r7a.large',
    InstanceTypeR7AXlarge = 'r7a.xlarge',
    InstanceTypeR7A2Xlarge = 'r7a.2xlarge',
    InstanceTypeR7A4Xlarge = 'r7a.4xlarge',
    InstanceTypeR7A8Xlarge = 'r7a.8xlarge',
    InstanceTypeR7A12Xlarge = 'r7a.12xlarge',
    InstanceTypeR7A16Xlarge = 'r7a.16xlarge',
    InstanceTypeR7A24Xlarge = 'r7a.24xlarge',
    InstanceTypeR7A32Xlarge = 'r7a.32xlarge',
    InstanceTypeR7A48Xlarge = 'r7a.48xlarge',
    InstanceTypeC7ILarge = 'c7i.large',
    InstanceTypeC7IXlarge = 'c7i.xlarge',
    InstanceTypeC7I2Xlarge = 'c7i.2xlarge',
    InstanceTypeC7I4Xlarge = 'c7i.4xlarge',
    InstanceTypeC7I8Xlarge = 'c7i.8xlarge',
    InstanceTypeC7I12Xlarge = 'c7i.12xlarge',
    InstanceTypeC7I16Xlarge = 'c7i.16xlarge',
    InstanceTypeC7I24Xlarge = 'c7i.24xlarge',
    InstanceTypeC7I48Xlarge = 'c7i.48xlarge',
    InstanceTypeMac2M2ProMetal = 'mac2-m2pro.metal',
    InstanceTypeR7IzLarge = 'r7iz.large',
    InstanceTypeR7IzXlarge = 'r7iz.xlarge',
    InstanceTypeR7Iz2Xlarge = 'r7iz.2xlarge',
    InstanceTypeR7Iz4Xlarge = 'r7iz.4xlarge',
    InstanceTypeR7Iz8Xlarge = 'r7iz.8xlarge',
    InstanceTypeR7Iz12Xlarge = 'r7iz.12xlarge',
    InstanceTypeR7Iz16Xlarge = 'r7iz.16xlarge',
    InstanceTypeR7Iz32Xlarge = 'r7iz.32xlarge',
    InstanceTypeC7AMedium = 'c7a.medium',
    InstanceTypeC7ALarge = 'c7a.large',
    InstanceTypeC7AXlarge = 'c7a.xlarge',
    InstanceTypeC7A2Xlarge = 'c7a.2xlarge',
    InstanceTypeC7A4Xlarge = 'c7a.4xlarge',
    InstanceTypeC7A8Xlarge = 'c7a.8xlarge',
    InstanceTypeC7A12Xlarge = 'c7a.12xlarge',
    InstanceTypeC7A16Xlarge = 'c7a.16xlarge',
    InstanceTypeC7A24Xlarge = 'c7a.24xlarge',
    InstanceTypeC7A32Xlarge = 'c7a.32xlarge',
    InstanceTypeC7A48Xlarge = 'c7a.48xlarge',
    InstanceTypeC7AMetal48Xl = 'c7a.metal-48xl',
    InstanceTypeR7AMetal48Xl = 'r7a.metal-48xl',
    InstanceTypeR7ILarge = 'r7i.large',
    InstanceTypeR7IXlarge = 'r7i.xlarge',
    InstanceTypeR7I2Xlarge = 'r7i.2xlarge',
    InstanceTypeR7I4Xlarge = 'r7i.4xlarge',
    InstanceTypeR7I8Xlarge = 'r7i.8xlarge',
    InstanceTypeR7I12Xlarge = 'r7i.12xlarge',
    InstanceTypeR7I16Xlarge = 'r7i.16xlarge',
    InstanceTypeR7I24Xlarge = 'r7i.24xlarge',
    InstanceTypeR7I48Xlarge = 'r7i.48xlarge',
    InstanceTypeDl2Q24Xlarge = 'dl2q.24xlarge',
    InstanceTypeMac2M2Metal = 'mac2-m2.metal',
    InstanceTypeI4I12Xlarge = 'i4i.12xlarge',
    InstanceTypeI4I24Xlarge = 'i4i.24xlarge',
    InstanceTypeC7IMetal24Xl = 'c7i.metal-24xl',
    InstanceTypeC7IMetal48Xl = 'c7i.metal-48xl',
    InstanceTypeM7IMetal24Xl = 'm7i.metal-24xl',
    InstanceTypeM7IMetal48Xl = 'm7i.metal-48xl',
    InstanceTypeR7IMetal24Xl = 'r7i.metal-24xl',
    InstanceTypeR7IMetal48Xl = 'r7i.metal-48xl',
    InstanceTypeR7IzMetal16Xl = 'r7iz.metal-16xl',
    InstanceTypeR7IzMetal32Xl = 'r7iz.metal-32xl',
    InstanceTypeC7GdMetal = 'c7gd.metal',
    InstanceTypeM7GdMetal = 'm7gd.metal',
    InstanceTypeR7GdMetal = 'r7gd.metal',
    InstanceTypeG6Xlarge = 'g6.xlarge',
    InstanceTypeG62Xlarge = 'g6.2xlarge',
    InstanceTypeG64Xlarge = 'g6.4xlarge',
    InstanceTypeG68Xlarge = 'g6.8xlarge',
    InstanceTypeG612Xlarge = 'g6.12xlarge',
    InstanceTypeG616Xlarge = 'g6.16xlarge',
    InstanceTypeG624Xlarge = 'g6.24xlarge',
    InstanceTypeG648Xlarge = 'g6.48xlarge',
    InstanceTypeGr64Xlarge = 'gr6.4xlarge',
    InstanceTypeGr68Xlarge = 'gr6.8xlarge',
    InstanceTypeC7IFlexLarge = 'c7i-flex.large',
    InstanceTypeC7IFlexXlarge = 'c7i-flex.xlarge',
    InstanceTypeC7IFlex2Xlarge = 'c7i-flex.2xlarge',
    InstanceTypeC7IFlex4Xlarge = 'c7i-flex.4xlarge',
    InstanceTypeC7IFlex8Xlarge = 'c7i-flex.8xlarge',
    InstanceTypeU7I12Tb224Xlarge = 'u7i-12tb.224xlarge',
    InstanceTypeU7In16Tb224Xlarge = 'u7in-16tb.224xlarge',
    InstanceTypeU7In24Tb224Xlarge = 'u7in-24tb.224xlarge',
    InstanceTypeU7In32Tb224Xlarge = 'u7in-32tb.224xlarge',
    InstanceTypeU7Ib12Tb224Xlarge = 'u7ib-12tb.224xlarge',
    InstanceTypeC7GnMetal = 'c7gn.metal',
    InstanceTypeR8GMedium = 'r8g.medium',
    InstanceTypeR8GLarge = 'r8g.large',
    InstanceTypeR8GXlarge = 'r8g.xlarge',
    InstanceTypeR8G2Xlarge = 'r8g.2xlarge',
    InstanceTypeR8G4Xlarge = 'r8g.4xlarge',
    InstanceTypeR8G8Xlarge = 'r8g.8xlarge',
    InstanceTypeR8G12Xlarge = 'r8g.12xlarge',
    InstanceTypeR8G16Xlarge = 'r8g.16xlarge',
    InstanceTypeR8G24Xlarge = 'r8g.24xlarge',
    InstanceTypeR8G48Xlarge = 'r8g.48xlarge',
    InstanceTypeR8GMetal24Xl = 'r8g.metal-24xl',
    InstanceTypeR8GMetal48Xl = 'r8g.metal-48xl',
    InstanceTypeMac2M1UltraMetal = 'mac2-m1ultra.metal',
}

export enum TypesMonitoringState {
    MonitoringStateDisabled = 'disabled',
    MonitoringStateDisabling = 'disabling',
    MonitoringStateEnabled = 'enabled',
    MonitoringStatePending = 'pending',
}

export interface TypesSeverityResult {
    /** @example 1 */
    criticalCount?: number
    /** @example 1 */
    highCount?: number
    /** @example 1 */
    lowCount?: number
    /** @example 1 */
    mediumCount?: number
    /** @example 1 */
    noneCount?: number
}

export enum TypesStandardUnit {
    StandardUnitSeconds = 'Seconds',
    StandardUnitMicroseconds = 'Microseconds',
    StandardUnitMilliseconds = 'Milliseconds',
    StandardUnitBytes = 'Bytes',
    StandardUnitKilobytes = 'Kilobytes',
    StandardUnitMegabytes = 'Megabytes',
    StandardUnitGigabytes = 'Gigabytes',
    StandardUnitTerabytes = 'Terabytes',
    StandardUnitBits = 'Bits',
    StandardUnitKilobits = 'Kilobits',
    StandardUnitMegabits = 'Megabits',
    StandardUnitGigabits = 'Gigabits',
    StandardUnitTerabits = 'Terabits',
    StandardUnitPercent = 'Percent',
    StandardUnitCount = 'Count',
    StandardUnitBytesSecond = 'Bytes/Second',
    StandardUnitKilobytesSecond = 'Kilobytes/Second',
    StandardUnitMegabytesSecond = 'Megabytes/Second',
    StandardUnitGigabytesSecond = 'Gigabytes/Second',
    StandardUnitTerabytesSecond = 'Terabytes/Second',
    StandardUnitBitsSecond = 'Bits/Second',
    StandardUnitKilobitsSecond = 'Kilobits/Second',
    StandardUnitMegabitsSecond = 'Megabits/Second',
    StandardUnitGigabitsSecond = 'Gigabits/Second',
    StandardUnitTerabitsSecond = 'Terabits/Second',
    StandardUnitCountSecond = 'Count/Second',
    StandardUnitNone = 'None',
}

export enum TypesTenancy {
    TenancyDefault = 'default',
    TenancyDedicated = 'dedicated',
    TenancyHost = 'host',
}

export enum TypesVolumeType {
    VolumeTypeStandard = 'standard',
    VolumeTypeIo1 = 'io1',
    VolumeTypeIo2 = 'io2',
    VolumeTypeGp2 = 'gp2',
    VolumeTypeSc1 = 'sc1',
    VolumeTypeSt1 = 'st1',
    VolumeTypeGp3 = 'gp3',
}

import axios, {
    AxiosInstance,
    AxiosRequestConfig,
    AxiosResponse,
    HeadersDefaults,
    ResponseType,
} from 'axios'
import { useSetAtom } from 'jotai'
import { useNavigate } from 'react-router-dom'

export type QueryParamsType = Record<string | number, any>

export interface FullRequestParams
    extends Omit<
        AxiosRequestConfig,
        'data' | 'params' | 'url' | 'responseType'
    > {
    /** set parameter to `true` for call `securityWorker` for this request */
    secure?: boolean
    /** request path */
    path: string
    /** content type of request body */
    type?: ContentType
    /** query params */
    query?: QueryParamsType
    /** format of response (i.e. response.json() -> format: "json") */
    format?: ResponseType
    /** request body */
    body?: unknown
}

export type RequestParams = Omit<
    FullRequestParams,
    'body' | 'method' | 'query' | 'path'
>

export interface ApiConfig<SecurityDataType = unknown>
    extends Omit<AxiosRequestConfig, 'data' | 'cancelToken'> {
    securityWorker?: (
        securityData: SecurityDataType | null
    ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void
    secure?: boolean
    format?: ResponseType
}

export enum ContentType {
    Json = 'application/json',
    FormData = 'multipart/form-data',
    UrlEncoded = 'application/x-www-form-urlencoded',
    Text = 'text/plain',
}

export class HttpClient<SecurityDataType = unknown> {
    public instance: AxiosInstance
    private securityData: SecurityDataType | null = null
    private securityWorker?: ApiConfig<SecurityDataType>['securityWorker']
    private secure?: boolean
    private format?: ResponseType

    constructor({
        securityWorker,
        secure,
        format,
        ...axiosConfig
    }: ApiConfig<SecurityDataType> = {}) {
        this.instance = axios.create({
            ...axiosConfig,
            baseURL: axiosConfig.baseURL || 'https://api.kaytu.io',
        })
        this.secure = secure
        this.format = format
        this.securityWorker = securityWorker
      
    }

    public setSecurityData = (data: SecurityDataType | null) => {
        this.securityData = data
    }

    protected mergeRequestParams(
        params1: AxiosRequestConfig,
        params2?: AxiosRequestConfig
    ): AxiosRequestConfig {
        const method = params1.method || (params2 && params2.method)

        return {
            ...this.instance.defaults,
            ...params1,
            ...(params2 || {}),
            headers: {
                ...((method &&
                    this.instance.defaults.headers[
                        method.toLowerCase() as keyof HeadersDefaults
                    ]) ||
                    {}),
                ...(params1.headers || {}),
                ...((params2 && params2.headers) || {}),
            },
        }
    }

    protected stringifyFormItem(formItem: unknown) {
        if (typeof formItem === 'object' && formItem !== null) {
            return JSON.stringify(formItem)
        } else {
            return `${formItem}`
        }
    }

    protected createFormData(input: Record<string, unknown>): FormData {
        return Object.keys(input || {}).reduce((formData, key) => {
            const property = input[key]
            const propertyContent: any[] =
                property instanceof Array ? property : [property]

            for (const formItem of propertyContent) {
                const isFileType =
                    formItem instanceof Blob || formItem instanceof File
                formData.append(
                    key,
                    isFileType ? formItem : this.stringifyFormItem(formItem)
                )
            }

            return formData
        }, new FormData())
    }

    public request = async <T = any, _E = any>({
        secure,
        path,
        type,
        query,
        format,
        body,
        ...params
    }: FullRequestParams): Promise<AxiosResponse<T>> => {
        const secureParams =
            ((typeof secure === 'boolean' ? secure : this.secure) &&
                this.securityWorker &&
                (await this.securityWorker(this.securityData))) ||
            {}
        const requestParams = this.mergeRequestParams(params, secureParams)
        const responseFormat = format || this.format || undefined

        if (
            type === ContentType.FormData &&
            body &&
            body !== null &&
            typeof body === 'object'
        ) {
            body = this.createFormData(body as Record<string, unknown>)
        }

        if (
            type === ContentType.Text &&
            body &&
            body !== null &&
            typeof body !== 'string'
        ) {
            body = JSON.stringify(body)
        }
        const instance = this.instance
        var temp_body = body
        // remove null undefined values and empty arrays and if object serach inside object
        if (body && typeof body === 'object') {
            temp_body = JSON.parse(JSON.stringify(body, (k, v) => {
                if (v === null || v === undefined || v === '' || (Array.isArray(v) && v.length === 0)) {
                    return undefined
                }
                return v
            }))
        }
        return this.instance.request({
            ...requestParams,
            headers: {
                ...(requestParams.headers || {}),
                ...(type && type !== ContentType.FormData
                    ? { 'Content-Type': type }
                    : {}),
            },
            params: query,
            responseType: responseFormat,
            data: temp_body,
            url: path,
        })
    }
}

/**
 * @title OpenGovernance Service API
 * @version 1.0
 * @baseUrl https://api.kaytu.io
 * @contact
 */
export class Api<
    SecurityDataType extends unknown
> extends HttpClient<SecurityDataType> {
    auth = {
        /**
         * @description Creates workspace key for the defined role with the defined name in the workspace.
         *
         * @tags keys
         * @name ApiV1KeyCreateCreate
         * @summary Create Workspace Key
         * @request POST:/auth/api/v1/key/create
         * @secure
         */
        apiV1KeyCreateCreate: (
            request: GithubComKaytuIoKaytuEnginePkgAuthApiCreateAPIKeyRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgAuthApiCreateAPIKeyResponse,
                EchoHTTPError
            >({
                path: `/auth/api/v1/keys/`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Deletes the specified workspace key by ID.
         *
         * @tags keys
         * @name ApiV1KeyDeleteDelete
         * @summary Delete Workspace Key
         * @request DELETE:/auth/api/v1/key/{id}/delete
         * @secure
         */
        apiV1KeyDeleteDelete: (id: string, params: RequestParams = {}) =>
            this.request<void, any>({
                path: `/auth/api/v1/key/${id}`,
                method: 'DELETE',
                secure: true,
                ...params,
            }),

        /**
         * @description Gets list of all keys in the workspace.
         *
         * @tags keys
         * @name ApiV1KeysList
         * @summary Get Workspace Keys
         * @request GET:/auth/api/v1/keys
         * @secure
         */
        apiV1KeysList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgAuthApiWorkspaceApiKey[],
                any
            >({
                path: `/auth/api/v1/keys`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns my user details
         *
         * @tags users
         * @name ApiV1MeList
         * @summary Get Me
         * @request GET:/auth/api/v1/me
         * @secure
         */
        apiV1MeList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgAuthApiGetMeResponse,
                any
            >({
                path: `/auth/api/v1/me`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Sends an invitation to a user to join the workspace with a designated role.
         *
         * @tags users
         * @name ApiV1UserInviteCreate
         * @summary Invite User
         * @request POST:/auth/api/v3/user/create
         * @secure
         */
        apiV1UserInviteCreate: (
            request: GithubComKaytuIoKaytuEnginePkgAuthApiInviteRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/auth/api/v1/user`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Changes user color blind mode and color mode
         *
         * @tags users
         * @name ApiV1UserPreferencesUpdate
         * @summary Change User Preferences
         * @request PUT:/auth/api/v1/user/preferences
         * @secure
         */
        apiV1UserPreferencesUpdate: (
            request: GithubComKaytuIoKaytuEnginePkgAuthApiChangeUserPreferencesRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/auth/api/v1/user/preferences`,
                method: 'PUT',
                body: request,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Updates the role of a user in the workspace.
         *
         * @tags users
         * @name ApiV1UserRoleBindingUpdate
         * @summary Update User Role
         * @request PUT: /auth/api/v3/user/update
         * @secure
         */
        apiV1UserRoleBindingUpdate: (
            request: GithubComKaytuIoKaytuEnginePkgAuthApiPutRoleBindingRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/auth/api/v1/user`,
                method: 'PUT',
                body: request,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Revokes a user's access to the workspace
         *
         * @tags users
         * @name ApiV1UserRoleBindingDelete
         * @summary Revoke User Access
         * @request DELETE:/auth/api/v1/user/role/binding
         * @secure
         */
        apiV1UserRoleBindingDelete: (
            // query: {
            //     /** User ID */
            //     userId: string
            // },
            id: number,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/auth/api/v1/user/${id}`,
                method: 'DELETE',
                // query: query,
                secure: true,
                ...params,
            }),

        /**
         * @description Retrieves the roles that the user who sent the request has in all workspaces they are a member of.
         *
         * @tags users
         * @name ApiV1UserRoleBindingsList
         * @summary Get User Roles
         * @request GET:/auth/api/v1/user/role/bindings
         * @secure
         */
        apiV1UserRoleBindingsList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgAuthApiGetRoleBindingsResponse,
                any
            >({
                path: `/auth/api/v1/user/role/bindings`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns user details by specified user id.
         *
         * @tags users
         * @name ApiV1UserDetail
         * @summary Get User details
         * @request GET:/auth/api/v1/user/{userId}
         * @secure
         */
        apiV1UserDetail: (userId: string, params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgAuthApiGetUserResponse,
                any
            >({
                path: `/auth/api/v1/user/${userId}`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieves a list of users who are members of the workspace.
         *
         * @tags users
         * @name ApiV1UsersList
         * @summary List Users
         * @request GET:/auth/api/v1/users
         * @secure
         */
        apiV1UsersList: (
            request: GithubComKaytuIoKaytuEnginePkgAuthApiGetUsersRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgAuthApiGetUsersResponse[],
                any
            >({
                path: `/auth/api/v1/users`,
                method: 'GET',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Get all the RoleBindings of the workspace. RoleBinding defines the roles and actions a user can perform. There are currently three roles (admin, editor, viewer). The workspace path is based on the DNS such as (workspace1.app.kaytu.io)
         *
         * @tags users
         * @name ApiV1WorkspaceRoleBindingsList
         * @summary Workspace user roleBindings.
         * @request GET:/auth/api/v1/workspace/role/bindings
         * @secure
         */
        apiV1WorkspaceRoleBindingsList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgAuthApiWorkspaceRoleBinding[],
                any
            >({
                path: `/auth/api/v1/users`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),
    }
    compliance = {
        
        /**
         * @description API for get new control list
         *
         * @tags compliance control
         * @name ApiV3ControlList
         * @summary Get control lists
         * @request POST:/compliance/api/v2/controls
         * @secure
         */
        apiV2ControlList: (
            request: GithubComKaytuIoKaytuEnginePkgControlApiListV2,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgControlApiListV2Response,
                any
            >({
                path: `/compliance/api/v3/controls`,
                method: 'POST',
                secure: true,
                body: request,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * @description API for get new control list
         *
         * @tags compliance control
         * @name ApiV3BenchmarkList
         * @summary Get control lists
         * @request POST:/compliance/api/v3/benchmarks
         * @secure
         */
        apiV3BenchmarkList: (
            request: GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3Response,
                any
            >({
                path: `/compliance/api/v3/benchmarks`,
                method: 'POST',
                secure: true,
                body: request,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * @description API for get Control detail
         *
         * @tags compliance control
         * @name ApiV3ControlDetail
         * @summary Get control lists
         * @request compliance/api/v3/control/${control_id}
         * @secure
         */
        apiV3ControlDetail: (id: string, params: RequestParams = {}) =>
            this.request<GithubComKaytuIoKaytuEnginePkgControlDetailV3, any>({
                path: `compliance/api/v3/control/${id}?showReferences=true`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * @description API for get Control Filters
         *
         * @tags compliance control
         * @name ApiV3ControlFilters
         * @summary Get control filters
         * @request /compliance/api/v3/control/filters
         * @secure
         */
        apiV3ControlFilters: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiV3ControlListFilters,
                any
            >({
                path: `/compliance/api/v3/controls/filters`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * @description API for get Control Filters
         *
         * @tags compliance control
         * @name ApiV3ControlFilters
         * @summary Get control filters
         * @request /compliance/api/v3/benchmarks/filters
         * @secure
         */
        apiV3BenchmarkFilters: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiV3BenchmarkListFilters,
                any
            >({
                path: `/compliance/api/v3/benchmarks/filters`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving all benchmark assigned sources with benchmark id
         *
         * @tags benchmarks_assignment
         * @name ApiV1AssignmentsBenchmarkDetail
         * @summary Get benchmark assigned sources
         * @request GET:/compliance/api/v1/assignments/benchmark/{benchmark_id}
         * @secure
         */
        apiV1AssignmentsBenchmarkDetail: (
            benchmarkId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkAssignedEntities,
                any
            >({
                path: `/compliance/api/v1/assignments/benchmark/${benchmarkId}`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

       
   
        

        /**
         * No description
         *
         * @tags compliance
         * @name ApiV1BenchmarksControlsDetail
         * @summary Get benchmark controls
         * @request GET:/compliance/api/v1/benchmarks/{benchmark_id}/controls
         * @secure
         */
        apiV1BenchmarksControlsDetail: (
            benchmarkId: string,
            query?: {
                /** Connection IDs to filter by */
                connectionId?: string[]
                /** Connection groups to filter by */
                connectionGroup?: string[]
                /** timestamp for values in epoch seconds */
                timeAt?: number
                /** Key-Value tags in key=value format to filter by */
                tag?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkControlSummary,
                any
            >({
                path: `/compliance/api/v1/benchmarks/${benchmarkId}/controls`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        
        /**
         * @description Changes benchmark settings.
         *
         * @tags compliance
         * @name ApiV1BenchmarksSettingsCreate
         * @summary change benchmark settings
         * @request POST:/compliance/api/v1/benchmarks/{benchmark_id}/settings
         * @secure
         */
        apiV1BenchmarksSettingsCreate: (
            benchmarkId?: string,
            query?: {
                /** tracksDriftEvents */
                tracksDriftEvents?: boolean
            },
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/compliance/api/v1/benchmarks/${benchmarkId}/settings`,
                method: 'POST',
                query: query,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Retrieving a summary of a benchmark and its associated checks and results.
         *
         * @tags compliance
         * @name ApiV1BenchmarksSummaryDetail
         * @summary Get benchmark summary
         * @request GET:/compliance/api/v1/benchmarks/{benchmark_id}/summary
         * @secure
         */
        apiV1BenchmarksSummaryDetail: (
            benchmarkId: string,
            query?: {
                /** Connection IDs to filter by */
                connectionId?: string[]
                /** Connection groups to filter by  */
                connectionGroup?: string[]
                /** Resource collection IDs to filter by */
                resourceCollection?: string[]
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** timestamp for values in epoch seconds */
                timeAt?: number
                /**
                 * Top account count
                 * @default 3
                 */
                topAccountCount?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkEvaluationSummary,
                any
            >({
                path: `/compliance/api/v1/benchmarks/${benchmarkId}/summary`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        
       
        /**
         * No description
         *
         * @tags compliance
         * @name ApiV1ControlsSummaryDetail
         * @summary Get control summary
         * @request GET:/compliance/api/v1/controls/{controlId}/summary
         * @secure
         */
        apiV1ControlsSummaryDetail: (
            controlId: string,
            query?: {
                /** Connection IDs to filter by */
                connectionId?: string[]
                /** Connection groups to filter by  */
                connectionGroup?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiControlSummary,
                any
            >({
                path: `/compliance/api/v1/controls/${controlId}/summary`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving all compliance finding events with respect to filters.
         *
         * @tags compliance
         * @name ApiV1FindingEventsCreate
         * @summary Get finding events
         * @request POST:/compliance/api/v1/finding_events
         * @secure
         */
        apiV1FindingEventsCreate: (
            request: GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingEventsRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingEventsResponse,
                any
            >({
                path: `/compliance/api/v1/finding_events`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving all compliance run finding events count with respect to filters.
         *
         * @tags compliance
         * @name ApiV1FindingEventsCountList
         * @summary Get finding events count
         * @request GET:/compliance/api/v1/finding_events/count
         * @secure
         */
        apiV1FindingEventsCountList: (
            query?: {
                /** ConformanceStatus to filter by defaults to all conformanceStatus except passed */
                conformanceStatus?: ('failed' | 'passed')[]
                /** BenchmarkID to filter by */
                benchmarkID?: string[]
                /** StateActive to filter by defaults to all stateActives */
                stateActive?: boolean[]
                /** Start time to filter by */
                startTime?: number
                /** End time to filter by */
                endTime?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiCountFindingEventsResponse,
                any
            >({
                path: `/compliance/api/v1/finding_events/count`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving possible values for finding event filters.
         *
         * @tags compliance
         * @name ApiV1FindingEventsFiltersCreate
         * @summary Get possible values for finding event filters
         * @request POST:/compliance/api/v1/finding_events/filters
         * @secure
         */
        apiV1FindingEventsFiltersCreate: (
            request: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEventFilters,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEventFiltersWithMetadata,
                any
            >({
                path: `/compliance/api/v1/finding_events/filters`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving single finding event
         *
         * @tags compliance
         * @name ApiV1FindingEventsSingleDetail
         * @summary Get single finding event
         * @request GET:/compliance/api/v1/finding_events/single/{id}
         * @secure
         */
        apiV1FindingEventsSingleDetail: (
            findingId: string,
            id: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEvent,
                any
            >({
                path: `/compliance/api/v1/finding_events/single/${id}`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving all compliance run findings with respect to filters.
         *
         * @tags compliance
         * @name ApiV1FindingsCreate
         * @summary Get findings
         * @request POST:/compliance/api/v1/compliance_result
         * @secure
         */
        apiV1FindingsCreate: (
            request: GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingsRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingsResponse,
                any
            >({
                path: `/compliance/api/v1/compliance_result`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

      

        /**
         * @description Retrieving all compliance run finding events with respect to filters.
         *
         * @tags compliance
         * @name ApiV1FindingsEventsDetail
         * @summary Get finding events by finding ID
         * @request GET:/compliance/api/v1/compliance_result/events/{id}
         * @secure
         */
        apiV1FindingsEventsDetail: (id: string, params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiGetFindingEventsByFindingIDResponse,
                any
            >({
                path: `/compliance/api/v1/compliance_result/events/${id}`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving possible values for finding filters.
         *
         * @tags compliance
         * @name ApiV1FindingsFiltersCreate
         * @summary Get possible values for finding filters
         * @request POST:/compliance/api/v1/compliance_result/filters
         * @secure
         */
        apiV1FindingsFiltersCreate: (
            request: GithubComKaytuIoKaytuEnginePkgComplianceApiFindingFilters,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiFindingFiltersWithMetadata,
                any
            >({
                path: `/compliance/api/v1/compliance_result/filters`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

 
        /**
         * @description Retrieving a single finding
         *
         * @tags compliance
         * @name ApiV1FindingsResourceCreate
         * @summary Get finding
         * @request POST:/compliance/api/v1/compliance_result/resource
         * @secure
         */
        apiV1FindingsResourceCreate: (
            request: GithubComKaytuIoKaytuEnginePkgComplianceApiGetSingleResourceFindingRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiGetSingleResourceFindingResponse,
                any
            >({
                path: `/compliance/api/v1/compliance_result/resource`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a single finding by finding ID
         *
         * @tags compliance
         * @name ApiV1FindingsSingleDetail
         * @summary Get single finding by finding ID
         * @request GET:/compliance/api/v1/compliance_result/single/{id}
         * @secure
         */
        apiV1FindingsSingleDetail: (id: string, params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiFinding,
                any
            >({
                path: `/compliance/api/v1/compliance_result/single/${id}`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving the top field by finding count.
         *
         * @tags compliance
         * @name ApiV1FindingsTopDetail
         * @summary Get top field by finding count
         * @request GET:/compliance/api/v1/compliance_result/top/{field}/{count}
         * @secure
         */
        apiV1FindingsTopDetail: (
            field:
                | 'resourceType'
                | 'integrationID'
                | 'resourceID'
                | 'service'
                | 'controlID',
            count: number,
            query?: {
                /** Connection IDs to filter by (inclusive) */
                integrationID?: string[]
                /** Connection IDs to filter by (exclusive) */
                notConnectionId?: string[]
                /** Connection groups to filter by  */
                connectionGroup?: string[]
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** BenchmarkID */
                benchmarkId?: string[]
                /** ControlID */
                controlId?: string[]
                /** Severities to filter by defaults to all severities except passed */
                severities?: ('none' | 'low' | 'medium' | 'high' | 'critical')[]
                /** ConformanceStatus to filter by defaults to all conformanceStatus except passed */
                conformanceStatus?: ('failed' | 'passed')[]
                /** StateActive to filter by defaults to true */
                stateActive?: boolean[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiGetTopFieldResponse,
                any
            >({
                path: `/compliance/api/v1/compliance_result/top/${field}/${count}`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving list of accounts with their security score and severities findings count
         *
         * @tags compliance
         * @name ApiV1FindingsAccountsDetail
         * @summary Get accounts findings summaries
         * @request GET:/compliance/api/v1/compliance_result/{benchmarkId}/accounts
         * @secure
         */
        apiV1FindingsAccountsDetail: (
            benchmarkId: string,
            query?: {
                /** Connection IDs to filter by */
                connectionId?: string[]
                /** Connection groups to filter by  */
                connectionGroup?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiGetAccountsFindingsSummaryResponse,
                any
            >({
                path: `/compliance/api/v1/compliance_result/${benchmarkId}/accounts`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving list of services with their security score and severities findings count
         *
         * @tags compliance
         * @name ApiV1FindingsServicesDetail
         * @summary Get services findings summary
         * @request GET:/compliance/api/v1/compliance_result/{benchmarkId}/services
         * @secure
         */
        apiV1FindingsServicesDetail: (
            benchmarkId: string,
            query?: {
                /** Connection IDs to filter by */
                connectionId?: string[]
                /** Connection groups to filter by  */
                connectionGroup?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiGetServicesFindingsSummaryResponse,
                any
            >({
                path: `/compliance/api/v1/compliance_result/${benchmarkId}/services`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

      

        /**
         * @description Syncs queries with the git backend.
         *
         * @tags compliance
         * @name ApiV1QueriesSyncList
         * @summary Sync queries
         * @request GET:/compliance/api/v1/queries/sync
         * @secure
         */
        apiV1QueriesSyncList: (
            query?: {
                /** Git URL */
                configzGitURL?: string
            },
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/compliance/api/v1/queries/sync`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Retrieving list of resource findings
         *
         * @tags compliance
         * @name ApiV1ResourceFindingsCreate
         * @summary List resource findings
         * @request POST:/compliance/api/v1/resource_findings
         * @secure
         */
        apiV1ResourceFindingsCreate: (
            request: GithubComKaytuIoKaytuEnginePkgComplianceApiListResourceFindingsRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgComplianceApiListResourceFindingsResponse,
                any
            >({
                path: `/compliance/api/v1/resource_findings`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
       
    }
    integration = {
        /**
         * @description Creating AWS source [standalone]
         *
         * @tags onboard
         * @name ApiV1ConnectionsAwsCreate
         * @summary Create AWS connection [standalone]
         * @request POST:/integration/api/v1/connections/aws
         * @secure
         */
        apiV1ConnectionsAwsCreate: (
            request: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateAWSConnectionRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateConnectionResponse,
                any
            >({
                path: `/integration/api/v1/connections/aws`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Counting connections either for the given connection type or all types if not specified.
         *
         * @tags connections
         * @name ApiV1ConnectionsCountList
         * @summary Count connections
         * @request GET:/integration/api/v1/connections/count
         * @secure
         */
        apiV1ConnectionsCountList: (
            query?: {
                /** Connector */
                connector?: string
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCountConnectionsResponse,
                any
            >({
                path: `/integration/api/v1/connections/count`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a list of connections summaries
         *
         * @tags connections
         * @name ApiV1ConnectionsSummariesList
         * @summary List connections summaries
         * @request GET:/integration/api/v1/connections/summaries
         * @secure
         */
        apiV1ConnectionsSummariesList: (
            query?: {
                /** Filter costs */
                filter?: string
                /** Connector */
                connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]
                /** Connection IDs */
                connectionId?: string[]
                /** Resource collection IDs to filter by */
                resourceCollection?: string[]
                /** Connection Groups */
                connectionGroups?: string[]
                /** filter by credential type */
                credentialType?: (
                    | 'auto-azure'
                    | 'auto-aws'
                    | 'manual-aws-org'
                    | 'manual-azure-spn'
                )[]
                /** lifecycle state filter */
                lifecycleState?:
                    | 'DISABLED'
                    | 'DISCOVERED'
                    | 'IN_PROGRESS'
                    | 'ONBOARD'
                    | 'ARCHIVED'
                /** health state filter */
                healthState?: 'healthy' | 'unhealthy'
                /** page size - default is 20 */
                pageSize?: number
                /** page number - default is 1 */
                pageNumber?: number
                /** start time in unix seconds */
                startTime?: number
                /** end time in unix seconds */
                endTime?: number
                /** for quicker inquiry send this parameter as false, default: true */
                needCost?: boolean
                /** for quicker inquiry send this parameter as false, default: true */
                needResourceCount?: boolean
                /** column to sort by - default is cost */
                sortBy?:
                    | 'onboard_date'
                    | 'resource_count'
                    | 'cost'
                    | 'growth'
                    | 'growth_rate'
                    | 'cost_growth'
                    | 'cost_growth_rate'
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityListConnectionsSummaryResponse,
                any
            >({
                path: `/integration/api/v1/connections/summaries`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Deleting a single connection either AWS / Azure for the given connection id. it will delete its parent credential too, if it doesn't have any other child.
         *
         * @tags connections
         * @name ApiV1ConnectionsDelete
         * @summary Delete connection
         * @request DELETE:/onboard/api/v1/source/{sourceId}
         * @secure
         */
        apiV1ConnectionsDelete: (
            connectionId: string,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/onboard/api/v1/source/${connectionId}`,
                method: 'DELETE',
                secure: true,
                ...params,
            }),

        /**
         * @description Get live connection health status with given connection ID for AWS.
         *
         * @tags connections
         * @name ApiV1ConnectionsAwsHealthcheckDetail
         * @summary Get AWS connection health
         * @request GET:/integration/api/v1/connections/{connectionId}/aws/healthcheck
         * @secure
         */
        apiV1ConnectionsAwsHealthcheckDetail: (
            connectionId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection,
                any
            >({
                path: `/integration/api/v1/connections/${connectionId}/aws/healthcheck`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Get live connection health status with given connection ID for Azure.
         *
         * @tags connections
         * @name ApiV1ConnectionsAzureHealthcheckDetail
         * @summary Get Azure connection health
         * @request GET:/integration/api/v1/connections/{connectionId}/azure/healthcheck
         * @secure
         */
        apiV1ConnectionsAzureHealthcheckDetail: (
            connectionId: string,
            query?: {
                /**
                 * Whether to update metadata or not
                 * @default true
                 */
                updateMetadata?: boolean
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection,
                any
            >({
                path: `/integration/api/v1/connections/${connectionId}/azure/healthcheck`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns list of all connectors
         *
         * @tags connectors
         * @name ApiV1ConnectorsList
         * @summary List connectors
         * @request GET:/integration/api/v1/connectors
         * @secure
         */
        apiV1ConnectorsList: (
            per_page: number,
            cursor: number,
            sort_by?: string,
            sort_order?: string,
            has_integration?: boolean,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnectorResponse,
                any
            >({
                path: `/integration/api/v1/integrations/types?per_page=${per_page}&cursor=${cursor}&sort_by=${sort_by}&sort_order=${sort_order}&has_integration=${has_integration}`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns list of all connectors
         *
         * @tags connectors
         * @name ApiV1ConnectorsList
         * @summary List connectors
         * @request GET:/integration/api/v1/integrations/types?enabled=true
         * @secure
         */
        apiV1EnabledConnectorsList: (
            per_page: number,
            cursor: number,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnectorResponse,
                any
            >({
                path: `/integration/api/v1/integrations/types?enabled=true`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving the list of metrics for catalog page.
         *
         * @tags integration
         * @name ApiV1ConnectorsMetricsList
         * @summary List catalog metrics
         * @request GET:/integration/api/v1/connectors/metrics
         * @secure
         */
        apiV1ConnectorsMetricsList: (
            query?: {
                /** Connector */
                connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]
                /** filter by credential type */
                credentialType?: (
                    | 'auto-azure'
                    | 'auto-aws'
                    | 'manual-aws-org'
                    | 'manual-azure-spn'
                )[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCatalogMetrics,
                any
            >({
                path: `/integration/api/v1/connectors/metrics`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Remove a credential by ID
         *
         * @tags credentials
         * @name ApiV1CredentialDelete
         * @summary Delete credential
         * @request DELETE:/onboard/api/v1/credential/
         * @secure
         */
        apiV1CredentialDelete: (
            credentialId: string,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/onboard/api/v1/credential/${credentialId}`,
                method: 'DELETE',
                secure: true,
                ...params,
            }),

        /**
         * @description Retrieving list of credentials with their details
         *
         * @tags credentials
         * @name ApiV1CredentialsList
         * @summary List credentials
         * @request GET:/integration/api/v1/credentials
         * @secure
         */
        apiV1CredentialsList: (
            query?: {
                /** filter by connector type */
                connector?: '' | 'AWS' | 'Azure' | 'EntraID'
                /** filter by health status */
                health?: 'healthy' | 'unhealthy'
                /** filter by credential type */
                credentialType?: (
                    | 'auto-azure'
                    | 'auto-aws'
                    | 'manual-aws-org'
                    | 'manual-azure-spn'
                )[]
                /**
                 * page size
                 * @default 50
                 */
                pageSize?: number
                /**
                 * page number
                 * @default 1
                 */
                pageNumber?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityListCredentialResponse,
                any
            >({
                path: `/integration/api/v1/credentials`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Creating AWS credential, testing it and onboard its accounts (organization account)
         *
         * @tags credentials
         * @name ApiV1CredentialsAwsCreate
         * @summary Create AWS credential and does onboarding for its accounts (organization account)
         * @request POST:/integration/api/v1/credentials/aws
         * @secure
         */
        apiV1CredentialsAwsCreate: (
            request: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateAWSCredentialRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateCredentialResponse,
                any
            >({
                path: `/integration/api/v1/credentials/aws`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Edit an aws credential by ID
         *
         * @tags credentials
         * @name ApiV1CredentialsAwsUpdate
         * @summary Edit aws credential
         * @request PUT:/integration/api/v1/credentials/aws/{credentialId}
         * @secure
         */
        apiV1CredentialsAwsUpdate: (
            credentialId: string,
            config: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityUpdateAWSCredentialRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/integration/api/v1/credentials/aws/${credentialId}`,
                method: 'PUT',
                body: config,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Onboard all available connections for an aws credential
         *
         * @tags credentials
         * @name ApiV1CredentialsAwsAutoonboardCreate
         * @summary Onboard aws credential connections
         * @request POST:/integration/api/v1/credentials/aws/{credentialId}/autoonboard
         * @secure
         */
        apiV1CredentialsAwsAutoonboardCreate: (
            credentialId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection[],
                any
            >({
                path: `/integration/api/v1/credentials/aws/${credentialId}/autoonboard`,
                method: 'POST',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Creating Azure credential, testing it and onboard its subscriptions
         *
         * @tags integration
         * @name ApiV1CredentialsAzureCreate
         * @summary Create Azure credential and does onboarding for its subscriptions
         * @request POST:/integration/api/v1/credentials/azure
         * @secure
         */
        apiV1CredentialsAzureCreate: (
            request: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateAzureCredentialRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCreateCredentialResponse,
                any
            >({
                path: `/integration/api/v1/credentials/azure`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Edit an azure credential by ID
         *
         * @tags credentials
         * @name ApiV1CredentialsAzureUpdate
         * @summary Edit azure credential
         * @request PUT:/integration/api/v1/credentials/azure/{credentialId}
         * @secure
         */
        apiV1CredentialsAzureUpdate: (
            credentialId: string,
            config: GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityUpdateAzureCredentialRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/integration/api/v1/credentials/azure/${credentialId}`,
                method: 'PUT',
                body: config,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Onboard all available connections for an azure credential
         *
         * @tags credentials
         * @name ApiV1CredentialsAzureAutoonboardCreate
         * @summary Onboard azure credential connections
         * @request POST:/integration/api/v1/credentials/azure/{credentialId}/autoonboard
         * @secure
         */
        apiV1CredentialsAzureAutoonboardCreate: (
            credentialId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection[],
                any
            >({
                path: `/integration/api/v1/credentials/azure/${credentialId}/autoonboard`,
                method: 'POST',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving credential details by credential ID
         *
         * @tags credentials
         * @name ApiV1CredentialsDetail
         * @summary Get Credential
         * @request GET:/integration/api/v1/credentials/{credentialId}
         * @secure
         */
        apiV1CredentialsDetail: (
            credentialId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityCredential,
                any
            >({
                path: `/integration/api/v1/credentials/${credentialId}`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),
    }
    inventory = {
        /**
         * @description Retrieving list of smart queries by specified filters
         *
         * @tags smart_query
         * @name ApiV1QueryList
         * @summary List smart queries
         * @request GET:/inventory/api/v1/query
         * @secure
         */
        apiV1QueryList: (
            request: GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItem[],
                any
            >({
                path: `/inventory/api/v1/query`,
                method: 'GET',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * @description Retrieving list of smart queries by specified filters
         *
         * @tags smart_query
         * @name ApiV2QueryList
         * @summary List smart queries
         * @request GET:/inventory/api/v2/queries
         * @secure
         */
        apiV2QueryList: (
            request: GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2Response,
                any
            >({
                path: `/inventory/api/v3/queries`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * @description Retrieving Query ListFilters
         *
         * @tags smart_query
         * @name ApiV2QueryListFilters
         * @summary List smart queries Filters
         * @request GET:/inventory/api/v3/queries/filters
         * @secure
         */
        apiV3QueryListFilter: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryFilters,
                any
            >({
                path: `/inventory/api/v3/queries/filters`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Run provided smart query and returns the result.
         *
         * @tags smart_query
         * @name ApiV1QueryRunCreate
         * @summary Run query
         * @request POST:/inventory/api/v1/query/run
         * @secure
         */
        apiV1QueryRunCreate: (
            request: GithubComKaytuIoKaytuEnginePkgInventoryApiRunQueryRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiRunQueryResponse,
                any
            >({
                path: `/inventory/api/v1/query/run`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description List queries which have been run recently
         *
         * @tags smart_query
         * @name ApiV1QueryRunHistoryList
         * @summary List recently ran queries
         * @request GET:/inventory/api/v1/query/run/history
         * @secure
         */
        apiV1QueryRunHistoryList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryHistory[],
                any
            >({
                path: `/inventory/api/v1/query/run/history`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving list of categories for analytics
         *
         * @tags analytics
         * @name ApiV2AnalyticsCategoriesList
         * @summary List Analytics categories
         * @request GET:/inventory/api/v2/analytics/categories
         * @secure
         */
        apiV2AnalyticsCategoriesList: (
            query?: {
                /** Metric type, default: assets */
                metricType?: 'assets' | 'spend'
                /** For assets minimum number of resources returned resourcetype must have, default 1 */
                minCount?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiAnalyticsCategoriesResponse,
                any
            >({
                path: `/inventory/api/v2/analytics/categories`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * @description Retrieving list of categories for all query
         *
         * @tags analytics
         * @name ApiV3InventoryCategoryList
         * @summary List Analytics categories
         * @request GET:/inventory/api/v3/queries/categories
         * @secure
         */
        apiV3InventoryCategoryList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiInventoryCategoriesResponse,
                any
            >({
                path: `/inventory/api/v3/queries/categories`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving tag values with the most resources for the given key.
         *
         * @tags analytics
         * @name ApiV2AnalyticsCompositionDetail
         * @summary List analytics composition
         * @request GET:/inventory/api/v2/analytics/composition/{key}
         * @secure
         */
        apiV2AnalyticsCompositionDetail: (
            key: string,
            query: {
                /** Metric type, default: assets */
                metricType?: 'assets' | 'spend'
                /** How many top values to return default is 5 */
                top: number
                /** Connector types to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** Resource collection IDs to filter by */
                resourceCollection?: string[]
                /** timestamp for resource count in epoch seconds */
                endTime?: number
                /** timestamp for resource count change comparison in epoch seconds */
                startTime?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiListResourceTypeCompositionResponse,
                any
            >({
                path: `/inventory/api/v2/analytics/composition/${key}`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving the count of resources and connections with respect to specified filters.
         *
         * @tags analytics
         * @name ApiV2AnalyticsCountList
         * @summary Count analytics
         * @request GET:/inventory/api/v2/analytics/count
         * @secure
         */
        apiV2AnalyticsCountList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiCountAnalyticsMetricsResponse,
                any
            >({
                path: `/inventory/api/v2/analytics/count`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving list of analytics with metrics of each type based on the given input filters.
         *
         * @tags analytics
         * @name ApiV2AnalyticsMetricList
         * @summary List analytics metrics
         * @request GET:/inventory/api/v2/analytics/metric
         * @secure
         */
        apiV2AnalyticsMetricList: (
            query?: {
                /** Key-Value tags in key=value format to filter by */
                tag?: string[]
                /** Metric type, default: assets */
                metricType?: 'assets' | 'spend'
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** Resource collection IDs to filter by */
                resourceCollection?: string[]
                /** Metric IDs */
                metricIDs?: string[]
                /** timestamp for resource count in epoch seconds */
                endTime?: number
                /** timestamp for resource count change comparison in epoch seconds */
                startTime?: number
                /** Minimum number of resources with this tag value, default 0 */
                minCount?: number
                /** Sort by field - default is count */
                sortBy?: 'name' | 'count' | 'growth' | 'growth_rate'
                /** page size - default is 20 */
                pageSize?: number
                /** page number - default is 1 */
                pageNumber?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiListMetricsResponse,
                any
            >({
                path: `/inventory/api/v2/analytics/metric`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns list of metrics
         *
         * @tags analytics
         * @name ApiV2AnalyticsMetricsListList
         * @summary List metrics
         * @request GET:/inventory/api/v2/analytics/metrics/list
         * @secure
         */
        apiV2AnalyticsMetricsListList: (
            query?: {
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Metric type, default: assets */
                metricType?: 'assets' | 'spend'
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiAnalyticsMetric[],
                any
            >({
                path: `/inventory/api/v2/analytics/metrics/list`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns list of metrics
         *
         * @tags analytics
         * @name ApiV2AnalyticsMetricsDetail
         * @summary List metrics
         * @request GET:/inventory/api/v2/analytics/metrics/{metric_id}
         * @secure
         */
        apiV2AnalyticsMetricsDetail: (
            metricId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiAnalyticsMetric,
                any
            >({
                path: `/inventory/api/v2/analytics/metrics/${metricId}`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving the cost composition with respect to specified filters. Retrieving information such as the total cost for the given time range, and the top services by cost.
         *
         * @tags analytics
         * @name ApiV2AnalyticsSpendCompositionList
         * @summary List cost composition
         * @request GET:/inventory/api/v2/analytics/spend/composition
         * @secure
         */
        apiV2AnalyticsSpendCompositionList: (
            query?: {
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** How many top values to return default is 5 */
                top?: number
                /** timestamp for start in epoch seconds */
                startTime?: number
                /** timestamp for end in epoch seconds */
                endTime?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiListCostCompositionResponse,
                any
            >({
                path: `/inventory/api/v2/analytics/spend/composition`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving the count of resources and connections with respect to specified filters.
         *
         * @tags analytics
         * @name ApiV2AnalyticsSpendCountList
         * @summary Count analytics spend
         * @request GET:/inventory/api/v2/analytics/spend/count
         * @secure
         */
        apiV2AnalyticsSpendCountList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiCountAnalyticsSpendResponse,
                any
            >({
                path: `/inventory/api/v2/analytics/spend/count`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving cost metrics with respect to specified filters. The API returns information such as the total cost and costs per each service based on the specified filters.
         *
         * @tags analytics
         * @name ApiV2AnalyticsSpendMetricList
         * @summary List spend metrics
         * @request GET:/inventory/api/v2/analytics/spend/metric
         * @secure
         */
        apiV2AnalyticsSpendMetricList: (
            query?: {
                /** Filter costs */
                filter?: string
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** timestamp for start in epoch seconds */
                startTime?: number
                /** timestamp for end in epoch seconds */
                endTime?: number
                /** Sort by field - default is cost */
                sortBy?: 'dimension' | 'cost' | 'growth' | 'growth_rate'
                /** page size - default is 20 */
                pageSize?: number
                /** page number - default is 1 */
                pageNumber?: number
                /** Metric IDs */
                metricIDs?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiListCostMetricsResponse,
                any
            >({
                path: `/inventory/api/v2/analytics/spend/metric`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns spend table with respect to the dimension and granularity
         *
         * @tags analytics
         * @name ApiV2AnalyticsSpendTableList
         * @summary Get Spend Trend
         * @request GET:/inventory/api/v2/analytics/spend/table
         * @secure
         */
        apiV2AnalyticsSpendTableList: (
            query?: {
                /** timestamp for start in epoch seconds */
                startTime?: number
                /** timestamp for end in epoch seconds */
                endTime?: number
                /** Granularity of the table, default is daily */
                granularity?: 'monthly' | 'daily' | 'yearly'
                /** Dimension of the table, default is metric */
                dimension?: 'connection' | 'metric'
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** Connector */
                connector?: string[]
                /** Metrics IDs */
                metricIds?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiSpendTableRow[],
                any
            >({
                path: `/inventory/api/v2/analytics/spend/table`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a list of costs over the course of the specified time frame based on the given input filters. If startTime and endTime are empty, the API returns the last month trend.
         *
         * @tags analytics
         * @name ApiV2AnalyticsSpendTrendList
         * @summary Get Cost Trend
         * @request GET:/inventory/api/v2/analytics/spend/trend
         * @secure
         */
        apiV2AnalyticsSpendTrendList: (
            query?: {
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** Metrics IDs */
                metricIds?: string[]
                /** timestamp for start in epoch seconds */
                startTime?: number
                /** timestamp for end in epoch seconds */
                endTime?: number
                /** Granularity of the table, default is daily */
                granularity?: 'monthly' | 'daily' | 'yearly'
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiCostTrendDatapoint[],
                any
            >({
                path: `/inventory/api/v2/analytics/spend/trend`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns asset table with respect to the dimension and granularity
         *
         * @tags analytics
         * @name ApiV2AnalyticsTableList
         * @summary Get Assets Table
         * @request GET:/inventory/api/v2/analytics/table
         * @secure
         */
        apiV2AnalyticsTableList: (
            query?: {
                /** timestamp for start in epoch seconds */
                startTime?: number
                /** timestamp for end in epoch seconds */
                endTime?: number
                /** Granularity of the table, default is daily */
                granularity?: 'monthly' | 'daily' | 'yearly'
                /** Dimension of the table, default is metric */
                dimension?: 'connection' | 'metric'
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiAssetTableRow[],
                any
            >({
                path: `/inventory/api/v2/analytics/table`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a list of tag keys with their possible values for all analytic metrics.
         *
         * @tags analytics
         * @name ApiV2AnalyticsTagList
         * @summary List analytics tags
         * @request GET:/inventory/api/v2/analytics/tag
         * @secure
         */
        apiV2AnalyticsTagList: (
            query?: {
                /** Connector type to filter by */
                connector?: string[]
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** Resource collection IDs to filter by */
                resourceCollection?: string[]
                /** Minimum number of resources/spend with this tag value, default 1 */
                minCount?: number
                /** Start time in unix timestamp format, default now - 1 month */
                startTime?: number
                /** End time in unix timestamp format, default now */
                endTime?: number
                /** Metric type, default: assets */
                metricType?: 'assets' | 'spend'
            },
            params: RequestParams = {}
        ) =>
            this.request<Record<string, string[]>, any>({
                path: `/inventory/api/v2/analytics/tag`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a list of resource counts over the course of the specified time frame based on the given input filters
         *
         * @tags analytics
         * @name ApiV2AnalyticsTrendList
         * @summary Get metric trend
         * @request GET:/inventory/api/v2/analytics/trend
         * @secure
         */
        apiV2AnalyticsTrendList: (
            query?: {
                /** Key-Value tags in key=value format to filter by */
                tag?: string[]
                /** Metric type, default: assets */
                metricType?: 'assets' | 'spend'
                /** Metric IDs to filter by */
                ids?: string[]
                /** Connector type to filter by */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Connection IDs to filter by - mutually exclusive with connectionGroup */
                connectionId?: string[]
                /** Connection group to filter by - mutually exclusive with connectionId */
                connectionGroup?: string[]
                /** Resource collection IDs to filter by */
                resourceCollection?: string[]
                /** timestamp for start in epoch seconds */
                startTime?: number
                /** timestamp for end in epoch seconds */
                endTime?: number
                /** Granularity of the table, default is daily */
                granularity?: 'monthly' | 'daily' | 'yearly'
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiResourceTypeTrendDatapoint[],
                any
            >({
                path: `/inventory/api/v2/analytics/trend`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving list of resource collections by specified filters
         *
         * @tags resource_collection
         * @name ApiV2MetadataResourceCollectionList
         * @summary List resource collections
         * @request GET:/inventory/api/v2/metadata/resource-collection
         * @secure
         */
        apiV2MetadataResourceCollectionList: (
            query?: {
                /** Resource collection IDs */
                id?: string[]
                /** Resource collection status */
                status?: ('' | 'active' | 'inactive')[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollection[],
                any
            >({
                path: `/inventory/api/v2/metadata/resource-collection`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving resource collection by specified ID
         *
         * @tags resource_collection
         * @name ApiV2MetadataResourceCollectionDetail
         * @summary Get resource collection
         * @request GET:/inventory/api/v2/metadata/resource-collection/{resourceCollectionId}
         * @secure
         */
        apiV2MetadataResourceCollectionDetail: (
            resourceCollectionId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollection,
                any
            >({
                path: `/inventory/api/v2/metadata/resource-collection/${resourceCollectionId}`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving list of resource collections by specified filters with inventory data
         *
         * @tags resource_collection
         * @name ApiV2ResourceCollectionList
         * @summary List resource collections with inventory data
         * @request GET:/inventory/api/v2/resource-collection
         * @secure
         */
        apiV2ResourceCollectionList: (
            query?: {
                /** Resource collection IDs */
                id?: string[]
                /** Resource collection status */
                status?: ('' | 'active' | 'inactive')[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollection[],
                any
            >({
                path: `/inventory/api/v2/resource-collection`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving resource collection by specified ID with inventory data
         *
         * @tags resource_collection
         * @name ApiV2ResourceCollectionDetail
         * @summary Get resource collection with inventory data
         * @request GET:/inventory/api/v2/resource-collection/{resourceCollectionId}
         * @secure
         */
        apiV2ResourceCollectionDetail: (
            resourceCollectionId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollection,
                any
            >({
                path: `/inventory/api/v2/resource-collection/${resourceCollectionId}`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving resource collection landscape by specified ID
         *
         * @tags resource_collection
         * @name ApiV2ResourceCollectionLandscapeDetail
         * @summary Get resource collection landscape
         * @request GET:/inventory/api/v2/resource-collection/{resourceCollectionId}/landscape
         * @secure
         */
        apiV2ResourceCollectionLandscapeDetail: (
            resourceCollectionId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgInventoryApiResourceCollectionLandscape,
                any
            >({
                path: `/inventory/api/v2/resource-collection/${resourceCollectionId}/landscape`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),
    }
    metadata = {
        /**
         * No description
         *
         * @tags metadata
         * @name ApiV1FilterList
         * @summary list filters
         * @request GET:/metadata/api/v1/filter
         * @secure
         */
        apiV1FilterList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgMetadataModelsFilter[],
                any
            >({
                path: `/metadata/api/v1/filter`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * No description
         *
         * @tags metadata
         * @name ApiV1FilterCreate
         * @summary add filter
         * @request POST:/metadata/api/v1/filter
         * @secure
         */
        apiV1FilterCreate: (
            req: GithubComKaytuIoKaytuEnginePkgMetadataModelsFilter,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/metadata/api/v1/filter`,
                method: 'POST',
                body: req,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Sets the config metadata for the given key
         *
         * @tags metadata
         * @name ApiV1MetadataCreate
         * @summary Set key metadata
         * @request POST:/metadata/api/v1/metadata
         * @secure
         */
        apiV1MetadataCreate: (
            req: GithubComKaytuIoKaytuEnginePkgMetadataApiSetConfigMetadataRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/metadata/api/v1/metadata`,
                method: 'POST',
                body: req,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Returns the config metadata for the given key
         *
         * @tags metadata
         * @name ApiV1MetadataDetail
         * @summary Get key metadata
         * @request GET:/metadata/api/v1/metadata/{key}
         * @secure
         */
        apiV1MetadataDetail: (key: string, params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgMetadataModelsConfigMetadata,
                any
            >({
                path: `/metadata/api/v1/metadata/${key}`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns the list of query parameters
         *
         * @tags metadata
         * @name ApiV1QueryParameterList
         * @summary List query parameters
         * @request GET:/metadata/api/v1/query_parameter
         * @secure
         */
        apiV1QueryParameterList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgMetadataApiListQueryParametersResponse,
                any
            >({
                path: `/metadata/api/v1/query_parameter`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Sets the query parameters from the request body
         *
         * @tags metadata
         * @name ApiV1QueryParameterCreate
         * @summary Set query parameter
         * @request POST:/metadata/api/v1/query_parameter
         * @secure
         */
        apiV1QueryParameterCreate: (
            req: GithubComKaytuIoKaytuEnginePkgMetadataApiSetQueryParameterRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/metadata/api/v1/query_parameter`,
                method: 'POST',
                body: req,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),
    }
    onboard = {
        /**
         * @description Retrieving the list of metrics for catalog page.
         *
         * @tags onboard
         * @name ApiV1CatalogMetricsList
         * @summary List catalog metrics
         * @request GET:/onboard/api/v1/catalog/metrics
         * @secure
         */
        apiV1CatalogMetricsList: (
            query?: {
                /** Connector */
                connector?: ('' | 'AWS' | 'Azure')[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiCatalogMetrics,
                any
            >({
                path: `/onboard/api/v1/catalog/metrics`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a list of connection groups
         *
         * @tags connection-groups
         * @name ApiV1ConnectionGroupsList
         * @summary List connection groups
         * @request GET:/onboard/api/v1/connection-groups
         * @secure
         */
        apiV1ConnectionGroupsList: (
            query?: {
                /**
                 * Populate connections
                 * @default false
                 */
                populateConnections?: boolean
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiConnectionGroup[],
                any
            >({
                path: `/onboard/api/v1/connection-groups`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a connection group
         *
         * @tags connection-groups
         * @name ApiV1ConnectionGroupsDetail
         * @summary Get connection group
         * @request GET:/onboard/api/v1/connection-groups/{connectionGroupName}
         * @secure
         */
        apiV1ConnectionGroupsDetail: (
            connectionGroupName: string,
            query?: {
                /**
                 * Populate connections
                 * @default false
                 */
                populateConnections?: boolean
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiConnectionGroup,
                any
            >({
                path: `/onboard/api/v1/connection-groups/${connectionGroupName}`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Creating AWS connection
         *
         * @tags onboard
         * @name ApiV1ConnectionsAwsCreate
         * @summary Create AWS connection
         * @request POST:/onboard/api/v1/connections/aws
         * @secure
         */
        apiV1ConnectionsAwsCreate: (
            request: GithubComKaytuIoKaytuEnginePkgOnboardApiCreateAwsConnectionRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiCreateConnectionResponse,
                any
            >({
                path: `/onboard/api/v1/connections/aws`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving a list of connections summaries
         *
         * @tags connections
         * @name ApiV1ConnectionsSummaryList
         * @summary List connections summaries
         * @request GET:/onboard/api/v1/connections/summary
         * @secure
         */
        apiV1ConnectionsSummaryList: (
            query?: {
                /** Filter costs */
                filter?: string
                /** Connector */
                connector?: ('' | 'AWS' | 'Azure')[]
                /** Connection IDs */
                connectionId?: string[]
                /** Resource collection IDs to filter by */
                resourceCollection?: string[]
                /** Connection Groups */
                connectionGroups?: string[]
                /** lifecycle state filter */
                lifecycleState?:
                    | 'DISABLED'
                    | 'DISCOVERED'
                    | 'IN_PROGRESS'
                    | 'ONBOARD'
                    | 'ARCHIVED'
                /** health state filter */
                healthState?: 'healthy' | 'unhealthy'
                /** page size - default is 20 */
                pageSize?: number
                /** page number - default is 1 */
                pageNumber?: number
                /** start time in unix seconds */
                startTime?: number
                /** end time in unix seconds */
                endTime?: number
                /** for quicker inquiry send this parameter as false, default: true */
                needCost?: boolean
                /** for quicker inquiry send this parameter as false, default: true */
                needResourceCount?: boolean
                /** column to sort by - default is cost */
                sortBy?:
                    | 'onboard_date'
                    | 'resource_count'
                    | 'cost'
                    | 'growth'
                    | 'growth_rate'
                    | 'cost_growth'
                    | 'cost_growth_rate'
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiListConnectionSummaryResponse,
                any
            >({
                path: `/onboard/api/v1/connections/summary`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * No description
         *
         * @tags onboard
         * @name ApiV1ConnectionsStateCreate
         * @summary Change connection lifecycle state
         * @request POST:/onboard/api/v1/connections/{connectionId}/state
         * @secure
         */
        apiV1ConnectionsStateCreate: (
            connectionId: string,
            request: GithubComKaytuIoKaytuEnginePkgOnboardApiChangeConnectionLifecycleStateRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/onboard/api/v1/connections/${connectionId}/state`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Returns list of all connectors
         *
         * @tags onboard
         * @name ApiV1ConnectorList
         * @summary List connectors
         * @request GET:/onboard/api/v1/connector
         * @secure
         */
        apiV1ConnectorList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiConnectorCount[],
                any
            >({
                path: `/onboard/api/v1/connector`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving list of credentials with their details
         *
         * @tags onboard
         * @name ApiV1CredentialList
         * @summary List credentials
         * @request GET:/onboard/api/v1/credential
         * @secure
         */
        apiV1CredentialList: (
            query?: {
                /** filter by connector type */
                connector?: '' | 'AWS' | 'Azure'
                /** filter by health status */
                health?: 'healthy' | 'unhealthy'
                /** filter by credential type */
                credentialType?: (
                    | 'auto-azure'
                    | 'auto-aws'
                    | 'manual-aws-org'
                    | 'manual-azure-spn'
                )[]
                /**
                 * page size
                 * @default 50
                 */
                pageSize?: number
                /**
                 * page number
                 * @default 1
                 */
                pageNumber?: number
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiListCredentialResponse,
                any
            >({
                path: `/onboard/api/v1/credential`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Creating connection credentials
         *
         * @tags onboard
         * @name ApiV1CredentialCreate
         * @summary Create connection credentials
         * @request POST:/onboard/api/v1/credential
         * @secure
         */
        apiV1CredentialCreate: (
            config: GithubComKaytuIoKaytuEnginePkgOnboardApiCreateCredentialRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiCreateCredentialResponse,
                any
            >({
                path: `/onboard/api/v1/credential`,
                method: 'POST',
                body: config,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Retrieving credential details by credential ID
         *
         * @tags onboard
         * @name ApiV1CredentialDetail
         * @summary Get Credential
         * @request GET:/onboard/api/v1/credential/{credentialId}
         * @secure
         */
        apiV1CredentialDetail: (
            credentialId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiCredential,
                any
            >({
                path: `/onboard/api/v1/credential/${credentialId}`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Edit a credential by ID
         *
         * @tags onboard
         * @name ApiV1CredentialUpdate
         * @summary Edit credential
         * @request PUT:/onboard/api/v1/credential/{credentialId}
         * @secure
         */
        apiV1CredentialUpdate: (
            credentialId: string,
            config: GithubComKaytuIoKaytuEnginePkgOnboardApiUpdateCredentialRequest,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/onboard/api/v1/credential/${credentialId}`,
                method: 'PUT',
                body: config,
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Remove a credential by ID
         *
         * @tags onboard
         * @name ApiV1CredentialDelete
         * @summary Delete credential
         * @request DELETE:/onboard/api/v1/credential/{credentialId}
         * @secure
         */
        apiV1CredentialDelete: (
            credentialId: string,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/onboard/api/v1/credential/${credentialId}`,
                method: 'DELETE',
                secure: true,
                ...params,
            }),

        /**
         * @description Onboard all available connections for a credential
         *
         * @tags onboard
         * @name ApiV1CredentialAutoonboardCreate
         * @summary Onboard credential connections
         * @request POST:/onboard/api/v1/credential/{credentialId}/autoonboard
         * @secure
         */
        apiV1CredentialAutoonboardCreate: (
            credentialId: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiConnection[],
                any
            >({
                path: `/onboard/api/v1/credential/${credentialId}/autoonboard`,
                method: 'POST',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Creating AWS source
         *
         * @tags onboard
         * @name ApiV1SourceAwsCreate
         * @summary Create AWS source
         * @request POST:/onboard/api/v1/source/aws
         * @secure
         */
        apiV1SourceAwsCreate: (
            request: GithubComKaytuIoKaytuEnginePkgOnboardApiSourceAwsRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiCreateSourceResponse,
                any
            >({
                path: `/onboard/api/v1/source/aws`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Creating Azure source
         *
         * @tags onboard
         * @name ApiV1SourceAzureCreate
         * @summary Create Azure source
         * @request POST:/onboard/api/v1/source/azure
         * @secure
         */
        apiV1SourceAzureCreate: (
            request: GithubComKaytuIoKaytuEnginePkgOnboardApiSourceAzureRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiCreateSourceResponse,
                any
            >({
                path: `/onboard/api/v1/source/azure`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Deleting a single source either AWS / Azure for the given source id.
         *
         * @tags onboard
         * @name ApiV1SourceDelete
         * @summary Delete source
         * @request DELETE:/onboard/api/v1/source/{sourceId}
         * @secure
         */
        apiV1SourceDelete: (sourceId: string, params: RequestParams = {}) =>
            this.request<void, any>({
                path: `/onboard/api/v1/source/${sourceId}`,
                method: 'DELETE',
                secure: true,
                ...params,
            }),

        /**
         * @description Get live source health status with given source ID.
         *
         * @tags onboard
         * @name ApiV1SourceHealthcheckDetail
         * @summary Get source health
         * @request GET:/onboard/api/v1/source/{sourceId}/healthcheck
         * @secure
         */
        apiV1SourceHealthcheckDetail: (
            sourceId: string,
            query?: {
                /**
                 * Whether to update metadata or not
                 * @default true
                 */
                updateMetadata?: boolean
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiConnection,
                any
            >({
                path: `/onboard/api/v1/source/${sourceId}/healthcheck`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Creating connection credentials
         *
         * @tags onboard
         * @name ApiV2CredentialCreate
         * @summary Create connection credentials
         * @request POST:/onboard/api/v2/credential
         * @secure
         */
        apiV2CredentialCreate: (
            config: GithubComKaytuIoKaytuEnginePkgOnboardApiV2CreateCredentialV2Request,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgOnboardApiV2CreateCredentialV2Response,
                any
            >({
                path: `/onboard/api/v2/credential`,
                method: 'POST',
                body: config,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
    }
    schedule = {
        /**
         * @description Triggers an analytics job to run immediately
         *
         * @tags describe
         * @name ApiV1AnalyticsTriggerUpdate
         * @summary TriggerAnalyticsJob
         * @request PUT:/schedule/api/v1/analytics/trigger
         * @secure
         */
        apiV1AnalyticsTriggerUpdate: (params: RequestParams = {}) =>
            this.request<void, any>({
                path: `/schedule/api/v1/analytics/trigger`,
                method: 'PUT',
                secure: true,
                ...params,
            }),

        /**
         * @description Get re-evaluate job for the given connection and control
         *
         * @tags describe
         * @name ApiV1ComplianceReEvaluateDetail
         * @summary Get re-evaluates compliance job
         * @request GET:/schedule/api/v1/compliance/re-evaluate/{benchmark_id}
         * @secure
         */
        apiV1ComplianceReEvaluateDetail: (
            benchmarkId: string,
            query: {
                /** Connection ID */
                integrationID: string[]
                /** Control ID */
                control_id?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgDescribeApiJobSeqCheckResponse,
                any
            >({
                path: `/schedule/api/v1/compliance/re-evaluate/${benchmarkId}`,
                method: 'GET',
                query: query,
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * @description Triggers a discovery job to run immediately for the given connection then triggers compliance job
         *
         * @tags describe
         * @name ApiV1ComplianceReEvaluateUpdate
         * @summary Re-evaluates compliance job
         * @request PUT:/schedule/api/v1/compliance/re-evaluate/{benchmark_id}
         * @secure
         */
        apiV1ComplianceReEvaluateUpdate: (
            benchmarkId: string,
            query: {
                /** Connection ID */
                integrationID: string[]
                /** Control ID */
                control_id?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/schedule/api/v1/compliance/re-evaluate/${benchmarkId}`,
                method: 'PUT',
                query: query,
                secure: true,
                ...params,
            }),

        /**
         * @description Triggers a compliance job to run immediately for the given benchmark
         *
         * @tags describe
         * @name ApiV1ComplianceTriggerUpdate
         * @summary Triggers compliance job
         * @request PUT:/schedule/api/v1/compliance/trigger
         * @secure
         */
        apiV1ComplianceTriggerUpdate: (
            query: {
                /** Benchmark IDs leave empty for everything */
                benchmark_id: string[]
                /** Connection IDs leave empty for default (enabled connections) */
                connection_id?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/schedule/api/v1/compliance/trigger`,
                method: 'PUT',
                query: query,
                secure: true,
                ...params,
            }),

        /**
         * @description Triggers a compliance job to run immediately for the given benchmark
         *
         * @tags describe
         * @name ApiV1ComplianceTriggerUpdate2
         * @summary Triggers compliance job
         * @request PUT:/schedule/api/v1/compliance/trigger/{benchmark_id}
         * @originalName apiV1ComplianceTriggerUpdate
         * @duplicate
         * @secure
         */
        apiV1ComplianceTriggerUpdate2: (
            benchmarkId: string,
            query?: {
                /** Connection ID */
                connection_id?: string[]
            },
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/schedule/api/v1/compliance/trigger/${benchmarkId}`,
                method: 'PUT',
                query: query,
                secure: true,
                ...params,
            }),

        /**
         * @description Triggers a compliance job to run immediately for the given benchmark
         *
         * @tags describe
         * @name ApiV1ComplianceTriggerSummaryUpdate
         * @summary Triggers compliance job
         * @request PUT:/schedule/api/v1/compliance/trigger/{benchmark_id}/summary
         * @secure
         */
        apiV1ComplianceTriggerSummaryUpdate: (
            benchmarkId: string,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/schedule/api/v1/compliance/trigger/${benchmarkId}/summary`,
                method: 'PUT',
                secure: true,
                ...params,
            }),

        /**
         * No description
         *
         * @tags describe
         * @name ApiV1DescribeConnectionStatusUpdate
         * @summary Get connection describe status
         * @request PUT:/schedule/api/v1/describe/connection/status
         * @secure
         */
        apiV1DescribeConnectionStatusUpdate: (
            query: {
                /** Connection ID */
                connection_id: string
            },
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/schedule/api/v1/describe/connection/status`,
                method: 'PUT',
                query: query,
                secure: true,
                ...params,
            }),

        /**
         * @description Triggers a describe job to run immediately for the given connection
         *
         * @tags describe
         * @name ApiV1DescribeTriggerUpdate
         * @summary Triggers describer
         * @request PUT:/schedule/api/v1/describe/trigger/{connection_id}
         * @secure
         */
        apiV1DescribeTriggerUpdate: (
            connectionId: string,
            query?: {
                /** Force full discovery */
                force_full?: boolean
                /** Resource Type */
                resource_type?: string[]
                /** Cost discovery */
                cost_discovery?: boolean
            },
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/schedule/api/v1/describe/trigger/${connectionId}`,
                method: 'PUT',
                query: query,
                secure: true,
                ...params,
            }),

        /**
         * No description
         *
         * @tags scheduler
         * @name ApiV1DiscoveryResourcetypesListList
         * @summary List all resource types that will be discovered
         * @request GET:/schedule/api/v1/discovery/resourcetypes/list
         * @secure
         */
        apiV1DiscoveryResourcetypesListList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgDescribeApiListDiscoveryResourceTypes,
                any
            >({
                path: `/schedule/api/v1/discovery/resourcetypes/list`,
                method: 'GET',
                secure: true,
                format: 'json',
                ...params,
            }),

        /**
         * No description
         *
         * @tags scheduler
         * @name ApiV1JobsCreate
         * @summary Lists all jobs
         * @request POST:/schedule/api/v1/jobs
         * @secure
         */
        apiV1JobsCreate: (
            request: GithubComKaytuIoKaytuEnginePkgDescribeApiListJobsRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgDescribeApiListJobsResponse,
                any
            >({
                path: `/schedule/api/v1/jobs`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
    }
    wastage = {
        /**
         * @description List wastage in AWS RDS
         *
         * @tags wastage
         * @name ApiV1WastageAwsRdsCreate
         * @summary List wastage in AWS RDS
         * @request POST:/wastage/api/v1/wastage/aws-rds
         * @secure
         */
        apiV1WastageAwsRdsCreate: (
            request: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsWastageRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsRdsWastageResponse,
                any
            >({
                path: `/wastage/api/v1/wastage/aws-rds`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description List wastage in AWS RDS Cluster
         *
         * @tags wastage
         * @name ApiV1WastageAwsRdsClusterCreate
         * @summary List wastage in AWS RDS Cluster
         * @request POST:/wastage/api/v1/wastage/aws-rds-cluster
         * @secure
         */
        apiV1WastageAwsRdsClusterCreate: (
            request: GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsClusterWastageRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesWastageApiEntityAwsClusterWastageResponse,
                any
            >({
                path: `/wastage/api/v1/wastage/aws-rds-cluster`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description List wastage in EC2 Instances
         *
         * @tags wastage
         * @name ApiV1WastageEc2InstanceCreate
         * @summary List wastage in EC2 Instances
         * @request POST:/wastage/api/v1/wastage/ec2-instance
         * @secure
         */
        apiV1WastageEc2InstanceCreate: (
            request: GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2InstanceWastageRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEngineServicesWastageApiEntityEC2InstanceWastageResponse,
                any
            >({
                path: `/wastage/api/v1/wastage/ec2-instance`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
    }
    workspace = {
        /**
         * No description
         *
         * @tags workspace
         * @name ApiV1BootstrapDetail
         * @summary Get bootstrap status
         * @request GET:/workspace/api/v1/bootstrap/{workspace_name}
         * @secure
         */
        apiV1BootstrapDetail: (
            workspaceName: string,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgWorkspaceApiBootstrapStatusResponse,
                any
            >({
                path: `/workspace/api/v1/bootstrap/${workspaceName}`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * No description
         *
         * @tags workspace
         * @name ApiV1BootstrapCredentialCreate
         * @summary Add credential for workspace to be onboarded
         * @request POST:/workspace/api/v1/bootstrap/{workspace_name}/credential
         * @secure
         */
        apiV1BootstrapCredentialCreate: (
            workspaceName: string,
            request: GithubComKaytuIoKaytuEnginePkgWorkspaceApiAddCredentialRequest,
            params: RequestParams = {}
        ) =>
            this.request<number, any>({
                path: `/workspace/api/v1/bootstrap/${workspaceName}/credential`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * No description
         *
         * @tags workspace
         * @name ApiV1BootstrapFinishCreate
         * @summary finish bootstrap
         * @request POST:/workspace/api/v1/bootstrap/{workspace_name}/finish
         * @secure
         */
        apiV1BootstrapFinishCreate: (
            workspaceName: string,
            params: RequestParams = {}
        ) =>
            this.request<string, any>({
                path: `/workspace/api/v1/bootstrap/${workspaceName}/finish`,
                method: 'POST',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns workspace created
         *
         * @tags workspace
         * @name ApiV1WorkspaceCreate
         * @summary Create workspace for workspace service
         * @request POST:/workspace/api/v1/workspace
         * @secure
         */
        apiV1WorkspaceCreate: (
            request: GithubComKaytuIoKaytuEnginePkgWorkspaceApiCreateWorkspaceRequest,
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgWorkspaceApiCreateWorkspaceResponse,
                any
            >({
                path: `/workspace/api/v1/workspace`,
                method: 'POST',
                body: request,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Returns all workspaces with owner id
         *
         * @tags workspace
         * @name ApiV1WorkspaceCurrentList
         * @summary List all workspaces with owner id
         * @request GET:/workspace/api/v1/workspace/current
         * @secure
         */
        apiV1WorkspaceCurrentList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceResponse,
                any
            >({
                path: `/metadata/api/v3/about`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * @description Delete workspace with workspace id
         *
         * @tags workspace
         * @name ApiV1WorkspaceDelete
         * @summary Delete workspace for workspace service
         * @request DELETE:/workspace/api/v1/workspace/{workspace_id}
         * @secure
         */
        apiV1WorkspaceDelete: (
            workspaceId: string,
            params: RequestParams = {}
        ) =>
            this.request<void, any>({
                path: `/workspace/api/v1/workspace/${workspaceId}`,
                method: 'DELETE',
                secure: true,
                type: ContentType.Json,
                ...params,
            }),

        /**
         * @description Returns all workspaces with owner id
         *
         * @tags workspace
         * @name ApiV1WorkspacesList
         * @summary List all workspaces with owner id
         * @request GET:/workspace/api/v1/workspaces
         * @secure
         */
        apiV1WorkspacesList: (params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceResponse[],
                any
            >({
                path: `/workspace/api/v1/workspaces`,
                method: 'GET',
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),

        /**
         * No description
         *
         * @tags workspace
         * @name ApiV1WorkspacesLimitsDetail
         * @summary Get workspace limits
         * @request GET:/workspace/api/v1/workspaces/limits/{workspace_name}
         * @secure
         */
        apiV1WorkspacesLimitsDetail: (
            workspaceName: string,
            query?: {
                /** Ignore usage */
                ignore_usage?: boolean
            },
            params: RequestParams = {}
        ) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceLimitsUsage,
                any
            >({
                path: `/workspace/api/v1/workspaces/limits/${workspaceName}`,
                method: 'GET',
                query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * No Load sample data
         *
         * @tags workspace
         * @name ApiV3LoadSampleData
         * @summary Get workspace limits
         * @request PUT:/workspace/api/v3/sample/sync
         * @secure
         */
        apiV3LoadSampleData: (data: any, params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceLimitsUsage,
                any
            >({
                path: `/metadata/api/v3/sample/sync`,
                method: 'PUT',
                // query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         * No Purge sample data
         *
         * @tags workspace
         * @name ApiV3PurgeSampleData
         * @summary Get workspace limits
         * @request PUT:/workspace/api/v3/sample/purge
         * @secure
         */
        apiV3PurgeSampleData: (data: any, params: RequestParams = {}) =>
            this.request<
                GithubComKaytuIoKaytuEnginePkgWorkspaceApiWorkspaceLimitsUsage,
                any
            >({
                path: `/metadata/api/v3/sample/purge`,
                method: 'PUT',
                // query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         *  Check user should setup or not
         *
         * @tags workspace
         * @name ApiV3GetSetup
         * @summary Get workspace limits
         * @request PUT:/workspace/api/v3/configured/status
         * @secure
         */
        apiV3GetSetup: (data: any, params: RequestParams = {}) =>
            this.request<string, any>({
                path: `/metadata/api/v3/configured/status`,
                method: 'GET',
                // query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
        /**
         *  After Setup
         *
         * @tags workspace
         * @name ApiV3GetSetup
         * @summary Get workspace limits
         * @request PUT:/workspace/api/v3/configured/set
         * @secure
         */
        apiV3DoneSetup: (data: any, params: RequestParams = {}) =>
            this.request<string, any>({
                path: `/metadata/api/v3/configured/set`,
                method: 'PUT',
                // query: query,
                secure: true,
                type: ContentType.Json,
                format: 'json',
                ...params,
            }),
    }
}
