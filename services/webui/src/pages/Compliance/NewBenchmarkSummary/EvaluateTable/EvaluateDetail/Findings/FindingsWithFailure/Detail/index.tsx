import { Link, useNavigate } from 'react-router-dom'
import { useAtomValue, useSetAtom } from 'jotai'
import {
    Button,
    Card,
    Col,
    Flex,
    Grid,

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
} from '../../../../../../../../api/api'
import {

    useComplianceApiV1FindingsEventsDetail,
    useComplianceApiV1FindingsResourceCreate,
} from '../../../../../../../../api/compliance.gen'

import { dateTimeDisplay } from '../../../../../../../../utilities/dateDisplay'

import {
    useScheduleApiV1ComplianceReEvaluateDetail,
    useScheduleApiV1ComplianceReEvaluateUpdate,
} from '../../../../../../../../api/schedule.gen'
import { isDemoAtom, notificationAtom } from '../../../../../../../../store'
import { getErrorMessage } from '../../../../../../../../types/apierror'
import { searchAtom } from '../../../../../../../../utilities/urlstate'
import { KeyValuePairs, Tabs } from '@cloudscape-design/components'

import { RenderObject } from '../../../../../../../../components/RenderObject'
import { severityBadge } from '../../../../../../../../utilities/badge'

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
                                                                to={`/compliance/${finding?.benchmarkID}/${finding?.controlID}?${searchParams}`}
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
                        ]}
                    />
                </>
            ) : (
                ''
            )}
        </>
    )
}


