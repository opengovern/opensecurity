// @ts-nocheck
import { Card, Flex, Text, Title } from '@tremor/react'

import { useEffect, useMemo, useState } from 'react'
import { useAtomValue, useSetAtom } from 'jotai/index'
import {
    Api,
    PlatformEnginePkgComplianceApiConformanceStatus,
    PlatformEnginePkgComplianceApiResourceFinding,
    SourceType,
    TypesFindingSeverity,
} from '../../../../api/api'
import AxiosAPI from '../../../../api/ApiConfig'
import { isDemoAtom, notificationAtom } from '../../../../store'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import { getConnectorIcon } from '../../../../components/Cards/ConnectorCard'
import ResourceFindingDetail from '../ResourceFindingDetail'
import KTable from '@cloudscape-design/components/table'
import Box from '@cloudscape-design/components/box'
import SpaceBetween from '@cloudscape-design/components/space-between'
import Badge from '@cloudscape-design/components/badge'
import {
    BreadcrumbGroup,
    DateRangePicker,
    Header,
    Link,
    Pagination,
    PropertyFilter,
} from '@cloudscape-design/components'
import { AppLayout, SplitPanel } from '@cloudscape-design/components'
import Filter from '../Filter'
import dayjs from 'dayjs'
import CustomPagination from '../../../../components/Pagination'


let sortKey: any[] = []

interface ICount {
    query: {
        connector: SourceType
        conformanceStatus:
            | PlatformEnginePkgComplianceApiConformanceStatus[]
            | undefined
        severity: TypesFindingSeverity[] | undefined
        connectionID: string[] | undefined
        controlID: string[] | undefined
        benchmarkID: string[] | undefined
        resourceTypeID: string[] | undefined
        lifecycle: boolean[] | undefined
    }
}

export default function ResourcesWithFailure({ query }: ICount) {
    const setNotification = useSetAtom(notificationAtom)
    const [loading, setLoading] = useState(false)

    const [open, setOpen] = useState(false)
    const [finding, setFinding] = useState<
        PlatformEnginePkgComplianceApiResourceFinding | undefined
    >(undefined)
    const [rows, setRows] = useState<any[]>()
    const [page, setPage] = useState(1)
    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [queries, setQuery] = useState(query)

    const isDemo = useAtomValue(isDemoAtom)

    const [date, setDate] = useState({
        key: 'previous-3-days',
        amount: 3,
        unit: 'day',
        type: 'relative',
    })
    const GetRows = () => {
        setLoading(true)
        const api = new Api()
        api.instance = AxiosAPI
        let isRelative = false
        let relative = ''
        let start = ''
        let end = ''
        if (date) {
            if (date.type == 'relative') {
                // @ts-ignore
                isRelative = true
                relative = `${date.amount}${date.unit}s`
            } else {
                // @ts-ignore

                start = dayjs(date?.startDate)
                // @ts-ignore

                end = dayjs(date?.endDate)
            }
        }
        api.compliance
            .apiV1ResourceFindingsCreate({
                filters: {
                    // @ts-ignore
                    integrationTypes: queries.connector.length
                        ? queries.connector
                        : [],
                    controlID: queries.controlID,
                    integrationID: queries.connectionID,
                    benchmarkID: queries.benchmarkID,
                    severity: queries.severity,
                    resourceTypeID: queries.resourceTypeID,
                    conformanceStatus: queries.conformanceStatus,
                    compliance_job_id: queries?.jobID,

                    ...(!isRelative &&
                        date && {
                            evaluatedAt: {
                                // @ts-ignore

                                from: start?.unix(),
                                // @ts-ignore

                                to: end?.unix(),
                            },
                        }),
                    ...(isRelative &&
                        date && {
                            interval: relative,
                        }),
                },
                // sort: [],
                limit: 15,
                // @ts-ignore
                afterSortKey: page == 1 ? [] : rows[rows?.length - 1].sortKey,

                // afterSortKey:
                //    [],
            })
            .then((resp) => {
                setLoading(false)
                if (resp.data.resourceFindings) {
                    setRows(resp.data.resourceFindings)
                } else {
                    setRows([])
                }
                // @ts-ignore

                setTotalPage(Math.ceil(resp.data.totalCount / 10))
                // @ts-ignore

                setTotalCount(resp.data.totalCount)
                // @ts-ignore
                // sortKey =
                //     resp.data?.resourceFindings?.at(
                //         (resp.data.resourceFindings?.length || 0) - 1
                //     )?.sortKey || []
            })
            .catch((err) => {
                setLoading(false)
                setNotification({
                    text: 'Can not Connect to Server',
                    type: 'warning',
                })
            })
    }

    useEffect(() => {
        GetRows()
    }, [page, queries, date])

    return (
        <>
            <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                toolsHide={true}
                navigationHide={true}
                splitPanelOpen={open}
                onSplitPanelToggle={() => {
                    setOpen(!open)
                    if (open) {
                        setFinding(undefined)
                    }
                }}
                splitPanel={
                    // @ts-ignore
                    <SplitPanel
                        // @ts-ignore
                        header={
                            finding ? (
                                <>
                                    <Flex justifyContent="start">
                                        {getConnectorIcon(finding?.connector)}
                                        <Title className="text-lg font-semibold ml-2 my-1">
                                            {finding?.resourceName}
                                        </Title>
                                    </Flex>
                                </>
                            ) : (
                                'Resource not selected'
                            )
                        }
                    >
                        <ResourceFindingDetail
                            // type="resource"
                            resourceFinding={finding}
                            open={open}
                            showOnlyOneControl={false}
                            onClose={() => setOpen(false)}
                            onRefresh={() => window.location.reload()}
                        />
                    </SplitPanel>
                }
                content={
                    <KTable
                        className="min-h-[450px]"
                        // resizableColumns
                        renderAriaLive={({
                            firstIndex,
                            lastIndex,
                            totalItemsCount,
                        }) =>
                            `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                        }
                        variant="full-page"
                        onSortingChange={(event) => {
                            // setSort(event.detail.sortingColumn.sortingField)
                            // setSortOrder(!sortOrder)
                        }}
                        // sortingColumn={sort}
                        // sortingDescending={sortOrder}
                        // sortingDescending={sortOrder == 'desc' ? true : false}
                        // @ts-ignore
                        onRowClick={(event) => {
                            const row = event.detail.item
                            if (row?.platformResourceID) {
                                setFinding(row)
                                setOpen(true)
                            } else {
                                setNotification({
                                    text: 'Detail for this finding is currently not available',
                                    type: 'warning',
                                })
                            }
                        }}
                        columnDefinitions={[
                            {
                                id: 'resourceName',
                                header: 'Resource name',
                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className="h-full"
                                        >
                                            <Text className="text-gray-800">
                                                {item.resourceName}
                                            </Text>
                                            <Text
                                                className={
                                                    isDemo ? 'blur-sm' : ''
                                                }
                                            >
                                                {item.platformResourceID}
                                            </Text>
                                        </Flex>
                                    </>
                                ),
                                sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 200,
                            },
                            {
                                id: 'resourceType',
                                header: 'Resource type',
                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                        >
                                            <Text className="text-gray-800">
                                                {item.resourceType}
                                            </Text>
                                            <Text>
                                                {item.resourceTypeLabel}
                                            </Text>
                                        </Flex>
                                    </>
                                ),
                                sortingField: 'title',
                                // minWidth: 400,
                                maxWidth: 200,
                            },
                            {
                                id: 'providerConnectionName',
                                header: 'Integration Type',
                                maxWidth: 100,
                                cell: (item) => (
                                    <>
                                        <Flex
                                            justifyContent="start"
                                            className={`h-full gap-3 group relative ${
                                                isDemo ? 'blur-sm' : ''
                                            }`}
                                        >
                                            {item.integrationType}
                                            {/* <Flex
                                                flexDirection="col"
                                                alignItems="start"
                                            >
                                                <Text className="text-gray-800">
                                                    {item.integrationID}
                                                </Text>
                                                <Text>
                                                    {item.integrationName}
                                                </Text>
                                            </Flex> */}
                                            {/* <Card className="cursor-pointer absolute w-fit h-fit z-40 right-1 scale-0 transition-all py-1 px-4 group-hover:scale-100">
                                                <Text color="blue">Open</Text>
                                            </Card> */}
                                        </Flex>
                                    </>
                                ),
                            },
                            {
                                id: 'totalCount',
                                header: 'Findings',
                                maxWidth: 100,

                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className="h-full"
                                        >
                                            <Text className="text-gray-800">{`${item.totalCount} incidents`}</Text>
                                            <Text>{`${
                                                item.totalCount -
                                                item.failedCount
                                            } passed`}</Text>
                                        </Flex>
                                    </>
                                ),
                            },
                            // {
                            //     id: 'conformanceStatus',
                            //     header: 'Status',
                            //     sortingField: 'severity',
                            //     cell: (item) => (
                            //         <Badge
                            //             // @ts-ignore
                            //             color={`${
                            //                 item.conformanceStatus == 'passed'
                            //                     ? 'green'
                            //                     : 'red'
                            //             }`}
                            //         >
                            //             {item.conformanceStatus}
                            //         </Badge>
                            //     ),
                            //     maxWidth: 100,
                            // },
                            // {
                            //     id: 'severity',
                            //     header: 'Severity',
                            //     sortingField: 'severity',
                            //     cell: (item) => (
                            //         <Badge
                            //             // @ts-ignore
                            //             color={`severity-${item.severity}`}
                            //         >
                            //             {item.severity.charAt(0).toUpperCase() +
                            //                 item.severity.slice(1)}
                            //         </Badge>
                            //     ),
                            //     maxWidth: 100,
                            // },
                            // {
                            //     id: 'evaluatedAt',
                            //     header: 'Last Evaluation',
                            //     cell: (item) => (
                            //         // @ts-ignore
                            //         <>{dateTimeDisplay(item.value)}</>
                            //     ),
                            // },
                        ]}
                        columnDisplay={[
                            { id: 'resourceName', visible: true },
                            { id: 'resourceType', visible: true },
                            { id: 'providerConnectionName', visible: true },
                            { id: 'totalCount', visible: true },
                            // { id: 'severity', visible: true },
                            // { id: 'evaluatedAt', visible: true },

                            // { id: 'action', visible: true },
                        ]}
                        enableKeyboardNavigation
                        // @ts-ignore
                        items={rows}
                        loading={loading}
                        loadingText="Loading resources"
                        // stickyColumns={{ first: 0, last: 1 }}
                        // stripedRows
                        trackBy="id"
                        empty={
                            <Box
                                margin={{ vertical: 'xs' }}
                                textAlign="center"
                                color="inherit"
                            >
                                <SpaceBetween size="m">
                                    <b>No resources</b>
                                </SpaceBetween>
                            </Box>
                        }
                        filter={
                            <Flex
                                flexDirection="row"
                                justifyContent="start"
                                alignItems="start"
                                className="gap-1 mt-1 sm:flex-row flex-col"
                            >
                                <Filter
                                    // @ts-ignore
                                    type={'resources'}
                                    onApply={(e) => {
                                        // @ts-ignore
                                        setQuery(e)
                                    }}
                                    setDate={() => {}}
                                />
                                <DateRangePicker
                                    onChange={({ detail }) =>
                                        // @ts-ignore
                                        setDate(detail.value)
                                    }
                                    // @ts-ignore

                                    value={date}
                                    relativeOptions={[
                                        {
                                            key: 'previous-5-minutes',
                                            amount: 5,
                                            unit: 'minute',
                                            type: 'relative',
                                        },
                                        {
                                            key: 'previous-30-minutes',
                                            amount: 30,
                                            unit: 'minute',
                                            type: 'relative',
                                        },
                                        {
                                            key: 'previous-1-hour',
                                            amount: 1,
                                            unit: 'hour',
                                            type: 'relative',
                                        },
                                        {
                                            key: 'previous-6-hours',
                                            amount: 6,
                                            unit: 'hour',
                                            type: 'relative',
                                        },
                                        {
                                            key: 'previous-3-days',
                                            amount: 3,
                                            unit: 'day',
                                            type: 'relative',
                                        },
                                        {
                                            key: 'previous-7-days',
                                            amount: 7,
                                            unit: 'day',
                                            type: 'relative',
                                        },
                                    ]}
                                    hideTimeOffset
                                    // showClearButton={false}
                                    absoluteFormat="long-localized"
                                    isValidRange={(range) => {
                                        if (range.type === 'absolute') {
                                            const [startDateWithoutTime] =
                                                range.startDate.split('T')
                                            const [endDateWithoutTime] =
                                                range.endDate.split('T')
                                            if (
                                                !startDateWithoutTime ||
                                                !endDateWithoutTime
                                            ) {
                                                return {
                                                    valid: false,
                                                    errorMessage:
                                                        'The selected date range is incomplete. Select a start and end date for the date range.',
                                                }
                                            }
                                            if (
                                                new Date(range.startDate) -
                                                    new Date(range.endDate) >
                                                0
                                            ) {
                                                return {
                                                    valid: false,
                                                    errorMessage:
                                                        'The selected date range is invalid. The start date must be before the end date.',
                                                }
                                            }
                                        }
                                        return { valid: true }
                                    }}
                                    i18nStrings={{}}
                                    placeholder="Filter by a date and time range"
                                />
                            </Flex>
                        }
                        header={
                            <Header className="w-full">
                                Resources{' '}
                                <span className=" font-medium">
                                    ({totalCount})
                                </span>
                            </Header>
                        }
                        pagination={
                            <CustomPagination
                                currentPageIndex={page}
                                pagesCount={totalPage}
                                onChange={({ detail }) =>
                                    setPage(detail.currentPageIndex)
                                }
                            />
                        }
                    />
                }
            />
        </>
    )
}
