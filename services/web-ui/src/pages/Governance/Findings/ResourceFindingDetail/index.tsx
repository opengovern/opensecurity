import { Link } from 'react-router-dom'
import { useAtomValue, useSetAtom } from 'jotai'
import {
    Card,
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
    PlatformEnginePkgComplianceApiResourceFinding,
} from '../../../../api/api'
import DrawerPanel from '../../../../components/DrawerPanel'
import { getConnectorIcon } from '../../../../components/Cards/ConnectorCard'
import SummaryCard from '../../../../components/Cards/SummaryCard'
import { useComplianceApiV1FindingsResourceCreate } from '../../../../api/compliance.gen'
import Spinner from '../../../../components/Spinner'
import { severityBadge } from '../../Controls'
import { isDemoAtom, notificationAtom } from '../../../../store'
import Timeline from '../FindingsWithFailure/Detail/Timeline'
import { searchAtom } from '../../../../utilities/urlstate'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import {
    KeyValuePairs,
    StatusIndicator,
    Tabs,
} from '@cloudscape-design/components'
import { RenderObject } from '../../../../components/RenderObject'

interface IResourceFindingDetail {
    resourceFinding:
        | PlatformEnginePkgComplianceApiResourceFinding
        | undefined
    controlID?: string
    showOnlyOneControl: boolean
    open: boolean
    onClose: () => void
    onRefresh: () => void
    linkPrefix?: string
}

export default function ResourceFindingDetail({
    resourceFinding,
    controlID,
    showOnlyOneControl,
    open,
    onClose,
    onRefresh,
    linkPrefix = '',
}: IResourceFindingDetail) {
    const { response, isLoading, sendNow } =
        useComplianceApiV1FindingsResourceCreate(
            { platformResourceID: resourceFinding?.platformResourceID || '' },
            {},
            false
        )
    const searchParams = useAtomValue(searchAtom)

    useEffect(() => {
        if (resourceFinding && open) {
            sendNow()
        }
    }, [resourceFinding, open])

    const isDemo = useAtomValue(isDemoAtom)

    const finding = resourceFinding?.findings
        ?.filter((f) => f.controlID === controlID)
        .at(0)

    const conformance = () => {
        if (showOnlyOneControl) {
            return (finding?.complianceStatus || 0) ===
                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed ? (
                <Flex className="w-fit gap-1.5">
                    <XCircleIcon className="h-4 text-rose-600" />
                    <Text>Failed</Text>
                </Flex>
            ) : (
                <Flex className="w-fit gap-1.5">
                    <CheckCircleIcon className="h-4 text-emerald-500" />
                    <Text>Passed</Text>
                </Flex>
            )
        }

        const failingControls = new Map<string, string>()
        resourceFinding?.findings?.forEach((f) => {
            failingControls.set(f.controlID || '', '')
        })

        return failingControls.size > 0 ? (
            <Flex className="w-fit gap-1.5">
                <XCircleIcon className="h-4 text-rose-600" />
                <Text>{failingControls.size} Failing Controls</Text>
            </Flex>
        ) : (
            <Flex className="w-fit gap-1.5">
                <CheckCircleIcon className="h-4 text-emerald-500" />
                <Text>Passed</Text>
            </Flex>
        )
    }
    const [items, setItems] = useState([])

    useEffect(() => {
        if (response?.controls) {
            const temp: any = []
            response?.controls
                ?.filter((c) => {
                    if (showOnlyOneControl) {
                        return c.controlID === controlID
                    }
                    return true
                })
                .map((control) => {
                    temp.push({
                        label: 'Control: ',
                        value: control.controlTitle,
                    })
                    temp.push({
                        label: 'Status: ',
                        value: (
                            <>
                                <StatusIndicator
                                    type={
                                        control.complianceStatus ===
                                        PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed
                                            ? 'success'
                                            : 'error'
                                    }
                                >
                                    {control.complianceStatus ===
                                    PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed
                                        ? 'Passed'
                                        : 'Failed'}
                                </StatusIndicator>
                            </>
                        ),
                    })

                    temp.push({
                        label: 'Severity: ',
                        value: <>{severityBadge(control.severity)}</>,
                    })
                })

            setItems(temp)
        }
    }, [response])

    return (
        <>
            <KeyValuePairs
                columns={4}
                items={[
                    {
                        label: 'Account',
                        value: (
                            <>
                            
                                {resourceFinding?.integrationType}
                                <Text
                                    className={` w-full text-start mb-0.5 truncate`}
                                >
                                    {resourceFinding?.integrationID}
                                </Text>
                            </>
                        ),
                    },
                    {
                        label: 'Resource',
                        value: (
                            <>
                                {resourceFinding?.resourceName}
                                <Text
                                    className={` w-full text-start mb-0.5 truncate`}
                                >
                                    {resourceFinding?.platformResourceID}
                                </Text>
                            </>
                        ),
                    },
                    {
                        label: 'Resource Type',
                        value: (
                            <>
                                {resourceFinding?.resourceTypeLabel}
                                <Text
                                    className={` w-full text-start mb-0.5 truncate`}
                                >
                                    {resourceFinding?.resourceType}
                                </Text>
                            </>
                        ),
                    },
                    {
                        label: 'Conformance Status',
                        value: conformance(),
                    },
                ]}
            />
            {/* <Grid className="w-full gap-4 mb-6" numItems={2}>
                <SummaryCard
                    title="Account"
                    metric={resourceFinding?.providerConnectionName}
                    secondLine={resourceFinding?.providerConnectionID}
                    blur={isDemo}
                    blurSecondLine={isDemo}
                    isString
                />
                <SummaryCard
                    title="Resource"
                    metric={resourceFinding?.resourceName}
                    secondLine={resourceFinding?.kaytuResourceID}
                    blurSecondLine={isDemo}
                    isString
                />
                <SummaryCard
                    title="Resource Type"
                    metric={resourceFinding?.resourceTypeLabel}
                    secondLine={resourceFinding?.resourceType}
                    isString
                />
                <SummaryCard
                    title="Conformance Status"
                    metric={conformance()}
                    isString
                />
            </Grid> */}
            <Tabs
                tabs={[
                    {
                        label: 'Resource Evidence',
                        id: '1',
                        disabled: !response?.resource,
                        content: (
                            <>
                                {' '}
                                <Title className="mb-2">JSON</Title>
                                <Card className="px-1.5 py-3 mb-2">
                                    <RenderObject
                                        obj={response?.resource || {}}
                                    />
                                </Card>
                            </>
                        ),
                    },
                    {
                        label: showOnlyOneControl
                            ? 'Summary'
                            : 'Applicable Controls',
                        id: '0',
                        content: (
                            <>
                                {showOnlyOneControl ? (
                                    <List>
                                        <ListItem className="py-6">
                                            <Text>Control</Text>

                                            {isLoading ? (
                                                <div className="animate-pulse h-3 w-64 my-1 bg-slate-200 dark:bg-slate-700 rounded" />
                                            ) : (
                                                <Link
                                                    className="text-right text-openg-500 cursor-pointer underline"
                                                    to={`${linkPrefix}${finding?.controlID}?${searchParams}`}
                                                >
                                                    {response?.controls
                                                        ?.filter(
                                                            (c) =>
                                                                c.controlID ===
                                                                finding?.controlID
                                                        )
                                                        .map(
                                                            (c) =>
                                                                c.controlTitle
                                                        )}
                                                </Link>
                                            )}
                                        </ListItem>
                                        <ListItem className="py-6">
                                            <Text>Severity</Text>
                                            {severityBadge(finding?.severity)}
                                        </ListItem>
                                        <ListItem className="py-6">
                                            <Text>Last evaluated</Text>
                                            <Text className="text-gray-800">
                                                {dateTimeDisplay(
                                                    finding?.evaluatedAt
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
                                    </List>
                                ) : (
                                    <>
                                        {response ? (
                                            <KeyValuePairs
                                                columns={3}
                                                items={items}
                                            />
                                        ) : (
                                            <Spinner />
                                        )}

                                        {/* <List>
                                            {isLoading ? (
                                                <Spinner className="mt-40" />
                                            ) : (
                                                response?.controls
                                                    ?.filter((c) => {
                                                        if (
                                                            showOnlyOneControl
                                                        ) {
                                                            return (
                                                                c.controlID ===
                                                                controlID
                                                            )
                                                        }
                                                        return true
                                                    })
                                                    .map((control) => (
                                                        <ListItem>
                                                            <Flex
                                                                flexDirection="col"
                                                                alignItems="start"
                                                                className="gap-1 w-fit max-w-[80%]"
                                                            >
                                                                <Text className="text-gray-800 w-full truncate">
                                                                    {
                                                                        control.controlTitle
                                                                    }
                                                                </Text>
                                                                <Flex justifyContent="start">
                                                                    {control.conformanceStatus ===
                                                                    PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed ? (
                                                                        <Flex className="w-full gap-1.5">
                                                                            <CheckCircleIcon className="h-4 text-emerald-500" />
                                                                            <Text>
                                                                                Passed
                                                                            </Text>
                                                                        </Flex>
                                                                    ) : (
                                                                        <Flex className="w-full gap-1.5">
                                                                            <XCircleIcon className="h-4 text-rose-600" />
                                                                            <Text>
                                                                                Failed
                                                                            </Text>
                                                                        </Flex>
                                                                    )}
                                                                    <Flex className="border-l border-gray-200 ml-3 pl-3 h-full">
                                                                        <Text className="text-xs">
                                                                            SECTION:
                                                                        </Text>
                                                                    </Flex>
                                                                </Flex>
                                                            </Flex>
                                                            {severityBadge(
                                                                control.severity
                                                            )}
                                                        </ListItem>
                                                    ))
                                            )}
                                        </List> */}
                                    </>
                                )}
                            </>
                        ),
                    },
                ]}
            />
        </>
    )
}

// <Timeline data={response} isLoading={isLoading} />
//
