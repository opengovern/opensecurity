import { Card, Flex, Text } from '@tremor/react'
import { numberDisplay } from '../../../../../utilities/numericDisplay'

export interface NewBenchmark {
    benchmark_id: string
    compliance_score: number
    severity_summary_by_control: SeveritySummaryBy
    severity_summary_by_resource: SeveritySummaryBy
    findings_summary: FindingsSummary
    issues_count: number
    top_integrations: null
    top_resources_with_issues: TopSWithIssue[]
    top_resource_types_with_issues: TopSWithIssue[]
    top_controls_with_issues: TopSWithIssue[]
    last_evaluated_at: Date
    last_job_status: string
    last_job_id: string
}

export interface FindingsSummary {
    total_count: number
    passed: number
    failed: number
}

export interface SeveritySummaryBy {
    total: Critical
    critical: Critical
    high: Critical
    medium: Critical
    low: Critical
    none: Critical
}

export interface Critical {
    total: number
    passed: number
    failed: number
}

export interface TopSWithIssue {
    field: Field
    key: string
    issues: number
}

export enum Field {
    Control = 'Control',
    Resource = 'Resource',
    ResourceType = 'ResourceType',
}
export const benchmarkChecks = (ben: NewBenchmark | undefined) => {
    const critical = ben?.severity_summary_by_control?.critical?.total || 0
    const high = ben?.severity_summary_by_control?.high?.total || 0
    const medium = ben?.severity_summary_by_control?.medium?.total || 0
    const low = ben?.severity_summary_by_control?.low?.total || 0
    const none = ben?.severity_summary_by_control?.none?.total || 0

    const total = ben?.severity_summary_by_control?.total?.total
    const failed = ben?.severity_summary_by_control?.total?.failed

    return {
        critical,
        high,
        medium,
        low,
        none,
        total,
        failed,
    }
}
interface ISeverityBar {
    benchmark: NewBenchmark | undefined
}
const JOB_STATUS = {
    CANCELED: 'Job is canceled',
    SUCCEEDED: 'Job Completed',
    FAILED: 'Job Failed',
    SUMMARIZER_IN_PROGRESS: 'In progress',
    SINK_IN_PROGRESS: 'In progress',
    RUNNERS_IN_PROGRESS : 'In progress',
}
export default function SeverityBar({ benchmark }: ISeverityBar) {
    const severity = [
        {
            name: 'Critical',
            color: '#6E120B',
            percent:
                (benchmarkChecks(benchmark).critical /
                    (benchmarkChecks(benchmark).failed || 1)) *
                    100 || 0,
            count: benchmarkChecks(benchmark).critical,
            controlPercent:
                ((benchmark?.severity_summary_by_control?.critical?.total ||
                    0) /
                    (benchmark?.severity_summary_by_control?.total?.total ||
                        1)) *
                100,
            control: benchmark?.severity_summary_by_control?.critical,
            resourcePercent:
                ((benchmark?.severity_summary_by_resource?.critical?.total ||
                    0) /
                    (benchmark?.severity_summary_by_resource?.total?.total ||
                        1)) *
                100,
            resource: benchmark?.severity_summary_by_resource?.critical,
        },
        {
            name: 'High',
            color: '#CA2B1D',
            percent:
                (benchmarkChecks(benchmark).high /
                    (benchmarkChecks(benchmark).failed || 1)) *
                    100 || 0,
            count: benchmarkChecks(benchmark).high,
            controlPercent:
                ((benchmark?.severity_summary_by_control?.high?.total || 0) /
                    (benchmark?.severity_summary_by_control?.total?.total ||
                        1)) *
                100,
            control: benchmark?.severity_summary_by_control?.high,
            resourcePercent:
                ((benchmark?.severity_summary_by_resource?.high?.total || 0) /
                    (benchmark?.severity_summary_by_resource?.total?.total ||
                        1)) *
                100,
            resource: benchmark?.severity_summary_by_resource?.high,
        },
        {
            name: 'Medium',
            color: '#EE9235',
            percent:
                (benchmarkChecks(benchmark).medium /
                    (benchmarkChecks(benchmark).failed || 1)) *
                    100 || 0,
            count: benchmarkChecks(benchmark).medium,
            controlPercent:
                ((benchmark?.severity_summary_by_control?.medium?.total || 0) /
                    (benchmark?.severity_summary_by_control?.total?.total || 1)) *
                100,
            control: benchmark?.severity_summary_by_control?.medium,
            resourcePercent:
                ((benchmark?.severity_summary_by_resource?.medium?.total || 0) /
                    (benchmark?.severity_summary_by_resource?.total?.total || 1)) *
                100,
            resource: benchmark?.severity_summary_by_resource?.medium,
        },
        {
            name: 'Low',
            color: '#F4C744',
            percent:
                (benchmarkChecks(benchmark).low /
                    (benchmarkChecks(benchmark).failed || 1)) *
                    100 || 0,
            count: benchmarkChecks(benchmark).low,
            controlPercent:
                ((benchmark?.severity_summary_by_control?.low?.total || 0) /
                    (benchmark?.severity_summary_by_control?.total?.total || 1)) *
                100,
            control: benchmark?.severity_summary_by_control?.low,
            resourcePercent:
                ((benchmark?.severity_summary_by_resource?.low?.total || 0) /
                    (benchmark?.severity_summary_by_resource?.total?.total || 1)) *
                100,
            resource: benchmark?.severity_summary_by_resource?.low,
        },
        {
            name: 'None',
            color: '#6B7280',
            percent:
                (benchmarkChecks(benchmark).none /
                    (benchmarkChecks(benchmark).failed || 1)) *
                    100 || 0,
            count: benchmarkChecks(benchmark).none,
            controlPercent:
                ((benchmark?.severity_summary_by_control?.none?.total || 0) /
                    (benchmark?.severity_summary_by_control?.total?.total || 1)) *
                100,
            control: benchmark?.severity_summary_by_control?.none,
            resourcePercent:
                ((benchmark?.severity_summary_by_resource?.none?.total || 0) /
                    (benchmark?.severity_summary_by_resource?.total?.total || 1)) *
                100,
            resource: benchmark?.severity_summary_by_resource?.none,
        },
    ]
    const passed = {
        name: 'Passed',
        color: '#54B584',
        percent:
            (((benchmark?.severity_summary_by_control?.total?.passed ?? 0) +
                (benchmark?.severity_summary_by_resource?.total?.passed ?? 0) ||
                0) /
                (benchmarkChecks(benchmark).total ?? 0)) *
            100,
        count:
            (benchmark?.severity_summary_by_control?.total?.passed ?? 0) +
            (benchmark?.severity_summary_by_resource?.total?.passed ?? 0),
        controlPercent:
            ((benchmark?.severity_summary_by_control?.total?.passed || 0) /
                (benchmark?.severity_summary_by_control?.total?.total || 1)) *
            100,
        control: benchmark?.severity_summary_by_control?.total?.passed || 0,
        resourcePercent:
            ((benchmark?.severity_summary_by_resource?.total?.passed || 0) /
                (benchmark?.severity_summary_by_resource?.total?.total || 1)) *
            100,
        resource: benchmark?.severity_summary_by_resource?.total?.passed,
    }

    return (
        <Flex flexDirection="col" alignItems="start">
            {/* @ts-ignore */}
            {benchmarkChecks(benchmark).total > 0 ? (
                <>
                    <Text className="mb-2">{`${numberDisplay(
                        (benchmark?.severity_summary_by_control?.total?.total ||
                            0) -
                            (benchmark?.severity_summary_by_control?.total
                                ?.passed || 0),
                        0
                    )} out of ${numberDisplay(
                        benchmark?.severity_summary_by_control?.total?.total ||
                            0,
                        0
                    )} controls failed`}</Text>
                </>
            ) : (
                <>
                    {/* @ts-ignore */}
                    <span>{JOB_STATUS[benchmark?.last_job_status]}</span>
                </>
            )}

            {/* @ts-ignore */}
            {benchmarkChecks(benchmark).total > 0 ? (
                <Flex alignItems="start" style={{ gap: '3px' }}>
                    <Flex flexDirection="col">
                        <Flex className="h-5" style={{ gap: '3px' }}>
                            {severity.map(
                                (s, i) =>
                                    s.controlPercent > 0 && (
                                        <div
                                            className="group h-full relative"
                                            style={{
                                                width: `${s.controlPercent}%`,
                                                minWidth: '2.5%',
                                            }}
                                        >
                                            <div
                                                className={`h-full w-full ${
                                                    i === 0 ? '' : ''
                                                }`}
                                                style={{
                                                    backgroundColor: s.color,
                                                }}
                                            />
                                            <Card
                                                className="absolute w-72 z-10 top-7 scale-0 transition-all p-2 group-hover:scale-100"
                                                style={{
                                                    border: `1px solid ${s.color}`,
                                                }}
                                            >
                                                <Flex
                                                    flexDirection="col"
                                                    alignItems="start"
                                                >
                                                    <Text
                                                        className={`text-[${s.color}] font-semibold mb-1`}
                                                    >
                                                        {s.name}
                                                    </Text>
                                                    <Flex>
                                                        <Text>Controls</Text>
                                                        <Text>
                                                            {`${numberDisplay(
                                                                s.control
                                                                    ?.passed ||
                                                                    0,
                                                                0
                                                            )} out of ${numberDisplay(
                                                                s.control
                                                                    ?.total ||
                                                                    0,
                                                                0
                                                            )} passed`}
                                                        </Text>
                                                    </Flex>
                                                    <Flex>
                                                        <Text>Issues</Text>
                                                        <Text>
                                                            {`${numberDisplay(
                                                                s.count,
                                                                0
                                                            )} (${s.percent.toFixed(
                                                                2
                                                            )}%)`}
                                                        </Text>
                                                    </Flex>
                                                    <Flex>
                                                        <Text>Resources</Text>
                                                        <Text>
                                                            {`${
                                                                s.resource
                                                                    ?.passed ||
                                                                0
                                                            } out of ${
                                                                s.resource
                                                                    ?.total || 0
                                                            } passed`}
                                                        </Text>
                                                    </Flex>
                                                </Flex>
                                            </Card>
                                        </div>
                                    )
                            )}
                        </Flex>
                        <Flex flexDirection="col" className="mt-2">
                            <Flex className="border-x-2 h-1.5 border-x-gray-400">
                                <div className="w-full h-0.5 bg-gray-400" />
                            </Flex>
                            <Text className="mt-1 text-xs">{`${(
                                (((benchmark?.severity_summary_by_control?.total
                                    ?.total || 0) -
                                    (benchmark?.severity_summary_by_control
                                        ?.total?.passed || 0)) /
                                    (benchmark?.severity_summary_by_control
                                        ?.total?.total || 1)) *
                                100
                            ).toFixed(2)}% failed`}</Text>
                        </Flex>
                    </Flex>
                    {passed.controlPercent > 0 && (
                        <div
                            className="group h-5 relative"
                            style={{
                                width: `${passed.controlPercent}%`,
                                minWidth: '2.5%',
                            }}
                        >
                            <div
                                className="h-full w-full"
                                style={{
                                    backgroundColor: passed.color,
                                }}
                            />
                            <Card
                                className="absolute w-72 z-10 top-7 scale-0 transition-all p-2 group-hover:scale-100"
                                style={{
                                    border: `1px solid ${passed.color}`,
                                }}
                            >
                                <Flex flexDirection="col" alignItems="start">
                                    <Text
                                        className={`text-[${passed.color}] font-semibold mb-1`}
                                    >
                                        Passed
                                    </Text>
                                    <Flex>
                                        <Text>Controls</Text>
                                        <Text>{`${numberDisplay(
                                            passed.control,
                                            0
                                        )}`}</Text>
                                    </Flex>
                                    <Flex>
                                        <Text>Issues</Text>
                                        <Text>
                                            {`${numberDisplay(
                                                passed.count,
                                                0
                                            )} (${passed.percent.toFixed(2)}%)`}
                                        </Text>
                                    </Flex>
                                    <Flex>
                                        <Text>Resources</Text>
                                        <Text>{`${numberDisplay(
                                            passed.resource,
                                            0
                                        )}`}</Text>
                                    </Flex>
                                </Flex>
                            </Card>
                        </div>
                    )}
                </Flex>
            ) : (
                <div className="bg-gray-200 h-5 rounded-md" />
            )}
        </Flex>
    )
}
