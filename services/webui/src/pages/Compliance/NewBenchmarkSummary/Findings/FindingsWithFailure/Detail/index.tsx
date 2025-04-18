import { Link, useNavigate } from 'react-router-dom'
import { useAtomValue, useSetAtom } from 'jotai'
import {
    Button,
    Card,
    Col,
    Flex,
    Grid,
    List,
    ListItem,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    Title,
} from '@tremor/react'
import { useEffect, useState } from 'react'
import ReactJson from '@microlink/react-json-view'
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline'
import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    PlatformEnginePkgComplianceApiFinding,
    PlatformEnginePkgComplianceApiResourceFinding,
} from '../../../../../../../api/api'
import {
    useComplianceApiV1BenchmarksControlsDetail,
    useComplianceApiV1ControlsSummaryDetail,
    useComplianceApiV1FindingsEventsDetail,
    useComplianceApiV1FindingsResourceCreate,
} from '../../../../../../../api/compliance.gen'
import Spinner from '../../../../../../../components/Spinner'

import { dateTimeDisplay } from '../../../../../../../utilities/dateDisplay'
import Timeline from './Timeline'
import {
    useScheduleApiV1ComplianceReEvaluateDetail,
    useScheduleApiV1ComplianceReEvaluateUpdate,
} from '../../../../../../../api/schedule.gen'
import { isDemoAtom, notificationAtom } from '../../../../../../../store'
import { getErrorMessage } from '../../../../../../../types/apierror'
import { searchAtom } from '../../../../../../../utilities/urlstate'
import { KeyValuePairs, Tabs } from '@cloudscape-design/components'
import { severityBadge } from '../../../../../Controls'
import { RenderObject } from '../../../../../../../components/RenderObject'


interface IFindingDetail {
    finding: PlatformEnginePkgComplianceApiFinding | undefined
    type: 'finding' | 'resource'
    open: boolean
    onClose: () => void
    onRefresh: () => void
}

const renderStatus = (state: boolean | undefined) => {
    if (state) {
        return (
            <Flex className="w-fit gap-2">
                <div className="w-2 h-2 bg-emerald-500 rounded-full" />
                <Text className="text-gray-800">Active</Text>
            </Flex>
        )
    }
    return (
        <Flex className="w-fit gap-2">
            <div className="w-2 h-2 bg-rose-600 rounded-full" />
            <Text className="text-gray-800">Not active</Text>
        </Flex>
    )
}

export default function FindingDetail({
    finding,
    type,
    open,
    onClose,
    onRefresh,
}: IFindingDetail) {
    const { response, isLoading, sendNow } =
        useComplianceApiV1FindingsResourceCreate(
            { platformResourceID: finding?.platformResourceID || '' },
            {},
            false
        )
    const {
        response: findingTimeline,
        isLoading: findingTimelineLoading,
        sendNow: findingTimelineSend,
    } = useComplianceApiV1FindingsEventsDetail(finding?.id || '', {}, false)

    useEffect(() => {
        if (finding && open) {
            sendNow()
            if (type === 'finding') {
                findingTimelineSend()
            }
        }
    }, [finding, open])

    const failedEvents =
        findingTimeline?.findingEvents?.filter(
            (v) =>
                v.conformanceStatus ===
                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed
        ) || []

    const {
        response: getReevaluateResp,
        error: getReevaluateError,
        isLoading: isGetReevaluateLoading,
        isExecuted: isGetReevaluateExecuted,
        sendNow: refreshReEvaluate,
    } = useScheduleApiV1ComplianceReEvaluateDetail(
        finding?.benchmarkID || '',
        {
            integrationID: [finding?.integrationID || ''],
            control_id: [finding?.controlID || ''],
        },
        {},
        open
    )

    const {
        error: reevaluateError,
        isLoading: isReevaluateLoading,
        isExecuted: isReevaluateExecuted,
        sendNow: Reelavuate,
    } = useScheduleApiV1ComplianceReEvaluateUpdate(
        finding?.benchmarkID || '',
        {
            integrationID: [finding?.integrationID || ''],
            control_id: [finding?.controlID || ''],
        },
        {},
        false
    )

    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)

    const setNotification = useSetAtom(notificationAtom)
    useEffect(() => {
        if (isReevaluateExecuted && !isReevaluateLoading) {
            refreshReEvaluate()
            const err = getErrorMessage(reevaluateError)
            if (err.length > 0) {
                setNotification({
                    text: `Failed to re-evaluate due to ${err}`,
                    type: 'error',
                    position: 'bottomLeft',
                })
            } else {
                setNotification({
                    text: 'Re-evaluate job triggered',
                    type: 'success',
                    position: 'bottomLeft',
                })
            }
        }
    }, [isReevaluateLoading])

    const [wasReEvaluating, setWasReEvaluating] = useState<boolean>(false)
    useEffect(() => {
        if (isGetReevaluateExecuted && !isGetReevaluateLoading) {
            if (getReevaluateResp?.isRunning === true) {
                setTimeout(() => {
                    refreshReEvaluate()
                }, 5000)
            } else if (wasReEvaluating) {
                onRefresh()
            }
            setWasReEvaluating(getReevaluateResp?.isRunning || false)
        }
    }, [isGetReevaluateLoading])

    const reEvaluateLoading =
        (isReevaluateExecuted && isReevaluateLoading) ||
        (isGetReevaluateExecuted && isGetReevaluateLoading) ||
        getReevaluateResp?.isRunning === true

    const isDemo = useAtomValue(isDemoAtom)

    return (
        <>
            {finding ? (
                <>
                    <Grid className="w-full gap-4 mb-6" numItems={1}>
                        <KeyValuePairs
                            columns={4}
                            items={[
                                {
                                    label: 'Account',
                                    value: (
                                        <>
                                            {finding?.integrationName}
                                            <Text
                                                className={` w-full text-start mb-0.5 truncate`}
                                            >
                                                {finding?.integrationID}
                                            </Text>
                                        </>
                                    ),
                                },
                                {
                                    label: 'Resource',
                                    value: (
                                        <>
                                            {finding?.resourceName}
                                            <Text
                                                className={` w-full text-start mb-0.5 truncate`}
                                            >
                                                {finding?.resourceID}
                                            </Text>
                                        </>
                                    ),
                                },
                                {
                                    label: 'Resource Type',
                                    value: (
                                        <>
                                            {finding?.resourceTypeName}
                                            <Text
                                                className={` w-full text-start mb-0.5 truncate`}
                                            >
                                                {finding?.resourceType}
                                            </Text>
                                        </>
                                    ),
                                },
                                {
                                    label: 'Severity',
                                    value: severityBadge(finding?.severity),
                                },
                            ]}
                        />
                        {/* <SummaryCard
                            title="Account"
                            metric={finding?.providerConnectionName}
                            secondLine={finding?.providerConnectionID}
                            blur={isDemo}
                            blurSecondLine={isDemo}
                            isString
                        />
                        <SummaryCard
                            title="Resource"
                            metric={finding?.resourceName}
                            secondLine={finding?.resourceID}
                            blurSecondLine={isDemo}
                            isString
                        />
                        <SummaryCard
                            title="Resource Type"
                            metric={finding?.resourceTypeName}
                            secondLine={finding?.resourceType}
                            isString
                        />
                        <SummaryCard
                            title="Severity"
                            metric={severityBadge(finding?.severity)}
                            isString
                        /> */}
                        {/* <Button
                    color="orange"
                    variant="secondary"
                    disabled={reEvaluateLoading}
                    onClick={() => {
                        Reelavuate()
                    }}
                >
                    <Flex flexDirection="row">
                        {reEvaluateLoading && (
                            <Spinner
                                className="w-6 h-6"
                                color="fill-orange-600"
                            />
                        )}
                        Re-evaluate
                    </Flex>
                </Button> */}
                    </Grid>
                    <Tabs
                        tabs={[
                            {
                                label: 'Summary',
                                id: '0',
                                content: (
                                    <>
                                        <KeyValuePairs
                                            columns={5}
                                            items={[
                                                {
                                                    label: 'Control',
                                                    value: (
                                                        <>
                                                            <Link
                                                                className="text-openg-500 cursor-pointer underline"
                                                                to={`${finding?.controlID}?${searchParams}`}
                                                            >
                                                                {
                                                                    finding?.controlTitle
                                                                }
                                                            </Link>
                                                        </>
                                                    ),
                                                },
                                                {
                                                    label: 'Compliance Status',
                                                    value: (
                                                        <>
                                                            {finding?.complianceStatus ===
                                                            PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed ? (
                                                                <Flex className="w-fit gap-1.5">
                                                                    <CheckCircleIcon className="h-4 text-emerald-500" />
                                                                    <Text>
                                                                        Passed
                                                                    </Text>
                                                                </Flex>
                                                            ) : (
                                                                <Flex className="w-fit gap-1.5">
                                                                    <XCircleIcon className="h-4 text-rose-600" />
                                                                    <Text>
                                                                        Failed
                                                                    </Text>
                                                                </Flex>
                                                            )}
                                                        </>
                                                    ),
                                                },
                                                {
                                                    label: 'Findings state',
                                                    value: (
                                                        <>
                                                            {renderStatus(
                                                                finding?.stateActive
                                                            )}
                                                        </>
                                                    ),
                                                },
                                                {
                                                    label: 'First discovered',
                                                    value: (
                                                        <>
                                                            {dateTimeDisplay(
                                                                failedEvents.at(
                                                                    failedEvents.length -
                                                                        1
                                                                )?.evaluatedAt
                                                            )}
                                                        </>
                                                    ),
                                                },
                                                {
                                                    label: 'Reason',
                                                    value: (
                                                        <>{finding?.reason}</>
                                                    ),
                                                },
                                            ]}
                                        />
                                        {/* <List>
                                            <ListItem className="py-6">
                                                <Text>Control</Text>
                                                <Link
                                                    className="text-openg-500 cursor-pointer underline"
                                                    to={`${finding?.controlID}?${searchParams}`}
                                                >
                                                    {finding?.controlTitle}
                                                </Link>
                                            </ListItem>
                                            <ListItem className="py-6">
                                                <Text>Conformance Status</Text>
                                                {finding?.conformanceStatus ===
                                                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed ? (
                                                    <Flex className="w-fit gap-1.5">
                                                        <CheckCircleIcon className="h-4 text-emerald-500" />
                                                        <Text>Passed</Text>
                                                    </Flex>
                                                ) : (
                                                    <Flex className="w-fit gap-1.5">
                                                        <XCircleIcon className="h-4 text-rose-600" />
                                                        <Text>Failed</Text>
                                                    </Flex>
                                                )}
                                            </ListItem>
                                            <ListItem className="py-6">
                                                <Text>Findings state</Text>
                                                {renderStatus(
                                                    finding?.stateActive
                                                )}
                                            </ListItem>
                                            <ListItem className="py-6">
                                                <Text>Last evaluated</Text>
                                                <Text className="text-gray-800">
                                                    {dateTimeDisplay(
                                                        finding?.evaluatedAt
                                                    )}
                                                </Text>
                                            </ListItem>
                                            <ListItem className="py-6">
                                                <Text>First discovered</Text>
                                                <Text className="text-gray-800">
                                                    {dateTimeDisplay(
                                                        failedEvents.at(
                                                            failedEvents.length -
                                                                1
                                                        )?.evaluatedAt
                                                    )}
                                                </Text>
                                            </ListItem>

                                            <ListItem className="py-6 space-x-5">
                                                <Flex
                                                    flexDirection="row"
                                                    justifyContent="between"
                                                    alignItems="start"
                                                    className="w-full"
                                                >
                                                    <Text className="w-1/4">
                                                        Reason
                                                    </Text>
                                                    <Text className="text-gray-800 text-end w-3/4 whitespace-break-spaces h-fit">
                                                        {finding?.reason}
                                                    </Text>
                                                </Flex>
                                            </ListItem>
                                        </List> */}
                                    </>
                                ),
                            },
                            {
                                label: 'Evidence',
                                disabled: !response?.resource,
                                id: '1',
                                content: (
                                    <>
                                        <Title className="mb-2">JSON</Title>
                                        <Card className="px-1.5 py-3 mb-2">
                                            <RenderObject
                                                obj={response?.resource || {}}
                                            />
                                        </Card>
                                    </>
                                ),
                            },
                            // {
                            //     label: 'Timeline',
                            //     id: '2',
                            //     content: (
                            //         <>
                            //             <Timeline
                            //                 data={
                            //                     type === 'finding'
                            //                         ? findingTimeline
                            //                         : response
                            //                 }
                            //                 isLoading={
                            //                     type === 'finding'
                            //                         ? findingTimelineLoading
                            //                         : isLoading
                            //                 }
                            //             />
                            //         </>
                            //     ),
                            // },
                        ]}
                    />
                </>
            ) : (
                ''
            )}
        </>
    )
}

// 
{
    /**
    
       <TabGroup>
                <Flex
                    flexDirection="row"
                    justifyContent="between"
                    alignItems="end"
                >
                    <TabList className="w-full">
                        {type === 'finding' ? (
                            <>
                                <Tab>Summary</Tab>
                                <Tab disabled={!response?.resource}>
                                    Resource Details
                                </Tab>
                                <Tab>Timeline</Tab>
                            </>
                        ) : (
                            <>
                                <Tab>Applicable Controls</Tab>
                                <Tab disabled={!response?.resource}>
                                    Resource Details
                                </Tab>
                            </>
                        )}
                    </TabList>
                    <Button
                        color="orange"
                        variant="secondary"
                        disabled={reEvaluateLoading}
                        onClick={() => {
                            Reelavuate()
                        }}
                    >
                        <Flex flexDirection="row">
                            {reEvaluateLoading && (
                                <Spinner
                                    className="w-6 h-6"
                                    color="fill-orange-600"
                                />
                            )}
                            Re-evaluate
                        </Flex>
                    </Button>
                </Flex>

                <TabPanels>
                    {type === 'finding' ? (
                        <TabPanel>
                            <List>
                                <ListItem className="py-6">
                                    <Text>Control</Text>
                                    <Link
                                        className="text-openg-500 cursor-pointer underline"
                                        to={`${finding?.controlID}?${searchParams}`}
                                    >
                                        {finding?.controlTitle}
                                    </Link>
                                </ListItem>
                                <ListItem className="py-6">
                                    <Text>Conformance Status</Text>
                                    {finding?.conformanceStatus ===
                                    PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed ? (
                                        <Flex className="w-fit gap-1.5">
                                            <CheckCircleIcon className="h-4 text-emerald-500" />
                                            <Text>Passed</Text>
                                        </Flex>
                                    ) : (
                                        <Flex className="w-fit gap-1.5">
                                            <XCircleIcon className="h-4 text-rose-600" />
                                            <Text>Failed</Text>
                                        </Flex>
                                    )}
                                </ListItem>
                                <ListItem className="py-6">
                                    <Text>Findings state</Text>
                                    {renderStatus(finding?.stateActive)}
                                </ListItem>
                                <ListItem className="py-6">
                                    <Text>Last evaluated</Text>
                                    <Text className="text-gray-800">
                                        {dateTimeDisplay(finding?.evaluatedAt)}
                                    </Text>
                                </ListItem>
                                <ListItem className="py-6">
                                    <Text>First discovered</Text>
                                    <Text className="text-gray-800">
                                        {dateTimeDisplay(
                                            failedEvents.at(
                                                failedEvents.length - 1
                                            )?.evaluatedAt
                                        )}
                                    </Text>
                                </ListItem>

                                <ListItem className="py-6 space-x-5">
                                    <Flex
                                        flexDirection="row"
                                        justifyContent="between"
                                        alignItems="start"
                                        className="w-full"
                                    >
                                        <Text className="w-1/4">Reason</Text>
                                        <Text className="text-gray-800 text-end w-3/4 whitespace-break-spaces h-fit">
                                            {finding?.reason}
                                        </Text>
                                    </Flex>
                                </ListItem>
                            </List>
                        </TabPanel>
                    ) : (
                        <TabPanel>
                            {isLoading ? (
                                <Spinner className="mt-12" />
                            ) : (
                                <List>
                                    {response?.controls?.map((control) => (
                                        <ListItem>
                                            <Flex
                                                flexDirection="col"
                                                alignItems="start"
                                                className="gap-1 w-fit max-w-[80%]"
                                            >
                                                <Text className="text-gray-800 w-full truncate">
                                                    {control.controlTitle}
                                                </Text>
                                                <Flex justifyContent="start">
                                                    {control.conformanceStatus ===
                                                    PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed ? (
                                                        <Flex className="w-fit gap-1.5">
                                                            <CheckCircleIcon className="h-4 text-emerald-500" />
                                                            <Text>Passed</Text>
                                                        </Flex>
                                                    ) : (
                                                        <Flex className="w-fit gap-1.5">
                                                            <XCircleIcon className="h-4 text-rose-600" />
                                                            <Text>Failed</Text>
                                                        </Flex>
                                                    )}
                                                    <Flex className="border-l border-gray-200 ml-3 pl-3 h-full">
                                                        <Text className="text-xs">
                                                            SECTION:
                                                        </Text>
                                                    </Flex>
                                                </Flex>
                                            </Flex>
                                            {severityBadge(control.severity)}
                                        </ListItem>
                                    ))}
                                </List>
                            )}
                        </TabPanel>
                    )}
                    <TabPanel>
                        <Title className="mb-2">JSON</Title>
                        <Card className="px-1.5 py-3 mb-2">
                            <ReactJson src={response?.resource || {}} />
                        </Card>
                    </TabPanel>
                    <TabPanel className="pt-8">
                        <Timeline
                            data={
                                type === 'finding' ? findingTimeline : response
                            }
                            isLoading={
                                type === 'finding'
                                    ? findingTimelineLoading
                                    : isLoading
                            }
                        />
                    </TabPanel>
                </TabPanels>
            </TabGroup>
    */
}