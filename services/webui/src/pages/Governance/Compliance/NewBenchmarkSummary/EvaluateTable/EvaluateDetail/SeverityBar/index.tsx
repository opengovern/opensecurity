// @ts-nocheck
import { Card, Flex, Text } from '@tremor/react'
import { numberDisplay } from '../../../../../../../utilities/numericDisplay'

export interface NewBenchmark {
    total_controls:  number;
    failed_controls: number;
    control_view:    ControlView;
}

export interface ControlView {
    by_severity: BySeverity;
}

export interface BySeverity {
    critical: Critical;
    high:     Critical;
    low:      Critical;
    medium:   Critical;
}

export interface Critical {
    total_controls:  number;
    failed_controls: number;
}


export enum Field {
    Control = 'Control',
    Resource = 'Resource',
    ResourceType = 'ResourceType',
}
export const benchmarkChecks = (ben: NewBenchmark | undefined) => {
    const critical = ben?.control_view?.by_severity?.critical?.total_controls || 0
    const high = ben?.control_view?.by_severity?.high?.total_controls || 0
    const medium = ben?.control_view?.by_severity?.medium?.total_controls || 0
    const low = ben?.control_view?.by_severity?.low?.total_controls || 0
    // const none = ben?.severity_summary_by_control?.none?.total || 0

    const total = ben?.total_controls
    const failed = ben?.failed_controls

    return {
        critical,
        high,
        medium,
        low,
        // none,
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
const calcPasses = (benchmark: NewBenchmark | undefined) => {
    // @ts-ignore
    return (benchmark?.total_controls  - benchmark?.failed_controls ) || 0
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
                ((benchmark?.control_view?.by_severity?.critical
                    ?.total_controls || 0) /
                    (benchmark?.total_controls || 1)) *
                100,
            control:
                benchmark?.control_view?.by_severity?.critical,
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
                ((benchmark?.control_view?.by_severity?.high?.total_controls ||
                    0) /
                    (benchmark?.total_controls || 1)) *
                100,
            control:
                benchmark?.control_view?.by_severity?.high,
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
                ((benchmark?.control_view?.by_severity?.medium
                    ?.total_controls || 0) /
                    (benchmark?.total_controls || 1)) *
                100,
            control:
                benchmark?.control_view?.by_severity?.medium,
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
                ((benchmark?.control_view?.by_severity?.low?.total_controls ||
                    0) /
                    (benchmark?.total_controls || 1)) *
                100,
            control:
                benchmark?.control_view?.by_severity?.low,
        },
        // {
        //     name: 'None',
        //     color: '#6B7280',
        //     percent:
        //         (benchmarkChecks(benchmark).none /
        //             (benchmarkChecks(benchmark).failed || 1)) *
        //             100 || 0,
        //     count: benchmarkChecks(benchmark).none,
        //     controlPercent:
        //         ((benchmark?.severity_summary_by_control?.none?.total || 0) /
        //             (benchmark?.total_controls ||
        //                 1)) *
        //         100,
        //     control: benchmark?.severity_summary_by_control?.none,
        //     resourcePercent:
        //         ((benchmark?.severity_summary_by_resource?.none?.total || 0) /
        //             (benchmark?.severity_summary_by_resource?.total?.total ||
        //                 1)) *
        //         100,
        //     resource: benchmark?.severity_summary_by_resource?.none,
        // },
    ]
    const passed = {
        name: 'Passed',
        color: '#54B584',
        percent:
            (((calcPasses(benchmark)) 
                 ||
                0) /
                (benchmarkChecks(benchmark).total ?? 0)) *
            100,
        count:
            (calcPasses(benchmark)) 
            ,
        controlPercent:
            ((calcPasses(benchmark)) /
                (benchmark?.total_controls || 1)) *
            100,
        control: calcPasses(benchmark),
        
    }

    return (
        <Flex flexDirection="col" alignItems="start">
            {/* @ts-ignore */}
            {benchmarkChecks(benchmark).total > 0 ? (
                <>
                    <Text className="mb-2">{`${numberDisplay(
                        (benchmark?.total_controls || 0) -
                            calcPasses(benchmark),
                        0
                    )} out of ${numberDisplay(
                        benchmark?.total_controls || 0,
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
                            <>{console.log(severity)}</>
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
                                                                    ?.total_controls 
                                                                     -
                                                                        (s
                                                                            ?.control
                                                                            ?.failed_controls 
                                                                            ),
                                                                0
                                                            )} out of ${numberDisplay(
                                                                s.control
                                                                    ?.total_controls ||
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
                                ((benchmark?.failed_controls ?? 0) /
                                    (benchmark?.total_controls || 1)) *
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
