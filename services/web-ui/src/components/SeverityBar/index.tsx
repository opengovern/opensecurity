import { Card, Flex, Text } from '@tremor/react'
import { PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary } from '../../api/api'
import { benchmarkChecks } from '../Cards/ComplianceCard'
import { numberDisplay } from '../../utilities/numericDisplay'

interface ISeverityBar {
    benchmark:
        | PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary
        | undefined
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
                ((benchmark?.controlsSeverityStatus?.critical?.total || 0) /
                    (benchmark?.controlsSeverityStatus?.total?.total || 1)) *
                100,
            control: benchmark?.controlsSeverityStatus?.critical,
            resourcePercent:
                ((benchmark?.resourcesSeverityStatus?.critical?.total || 0) /
                    (benchmark?.resourcesSeverityStatus?.total?.total || 1)) *
                100,
            resource: benchmark?.resourcesSeverityStatus?.critical,
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
                ((benchmark?.controlsSeverityStatus?.high?.total || 0) /
                    (benchmark?.controlsSeverityStatus?.total?.total || 1)) *
                100,
            control: benchmark?.controlsSeverityStatus?.high,
            resourcePercent:
                ((benchmark?.resourcesSeverityStatus?.high?.total || 0) /
                    (benchmark?.resourcesSeverityStatus?.total?.total || 1)) *
                100,
            resource: benchmark?.resourcesSeverityStatus?.high,
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
                ((benchmark?.controlsSeverityStatus?.medium?.total || 0) /
                    (benchmark?.controlsSeverityStatus?.total?.total || 1)) *
                100,
            control: benchmark?.controlsSeverityStatus?.medium,
            resourcePercent:
                ((benchmark?.resourcesSeverityStatus?.medium?.total || 0) /
                    (benchmark?.resourcesSeverityStatus?.total?.total || 1)) *
                100,
            resource: benchmark?.resourcesSeverityStatus?.medium,
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
                ((benchmark?.controlsSeverityStatus?.low?.total || 0) /
                    (benchmark?.controlsSeverityStatus?.total?.total || 1)) *
                100,
            control: benchmark?.controlsSeverityStatus?.low,
            resourcePercent:
                ((benchmark?.resourcesSeverityStatus?.low?.total || 0) /
                    (benchmark?.resourcesSeverityStatus?.total?.total || 1)) *
                100,
            resource: benchmark?.resourcesSeverityStatus?.low,
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
                ((benchmark?.controlsSeverityStatus?.none?.total || 0) /
                    (benchmark?.controlsSeverityStatus?.total?.total || 1)) *
                100,
            control: benchmark?.controlsSeverityStatus?.none,
            resourcePercent:
                ((benchmark?.resourcesSeverityStatus?.none?.total || 0) /
                    (benchmark?.resourcesSeverityStatus?.total?.total || 1)) *
                100,
            resource: benchmark?.resourcesSeverityStatus?.none,
        },
    ]
    const passed = {
        name: 'Passed',
        color: '#54B584',
        percent:
            ((benchmark?.conformanceStatusSummary?.passed || 0) /
                benchmarkChecks(benchmark).total) *
            100,
        count: benchmark?.conformanceStatusSummary?.passed || 0,
        controlPercent:
            ((benchmark?.controlsSeverityStatus?.total?.passed || 0) /
                (benchmark?.controlsSeverityStatus?.total?.total || 1)) *
            100,
        control: benchmark?.controlsSeverityStatus?.total?.passed || 0,
        resourcePercent:
            ((benchmark?.resourcesSeverityStatus?.total?.passed || 0) /
                (benchmark?.resourcesSeverityStatus?.total?.total || 1)) *
            100,
        resource: benchmark?.resourcesSeverityStatus?.total?.passed,
    }

    return (
        <Flex flexDirection="col" alignItems="start">
            <Text className="mb-2">{`${numberDisplay(
                (benchmark?.controlsSeverityStatus?.total?.total || 0) -
                    (benchmark?.controlsSeverityStatus?.total?.passed || 0),
                0
            )} out of ${numberDisplay(
                benchmark?.controlsSeverityStatus?.total?.total || 0,
                0
            )} controls failed`}</Text>
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
                                (((benchmark?.controlsSeverityStatus?.total
                                    ?.total || 0) -
                                    (benchmark?.controlsSeverityStatus?.total
                                        ?.passed || 0)) /
                                    (benchmark?.controlsSeverityStatus?.total
                                        ?.total || 1)) *
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
