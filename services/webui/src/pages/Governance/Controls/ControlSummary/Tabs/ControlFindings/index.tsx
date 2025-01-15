// @ts-nocheck
import { useEffect, useMemo, useState } from 'react'
import { useAtomValue, useSetAtom } from 'jotai/index'

import { Flex, Text, Title } from '@tremor/react'
import { ExclamationCircleIcon } from '@heroicons/react/24/outline'
import Table, { IColumn } from '../../../../../../components/Table'
import FindingDetail from '../../../../Findings/FindingsWithFailure/Detail'
import { isDemoAtom, notificationAtom } from '../../../../../../store'
import {
    Api,
    PlatformEnginePkgComplianceApiConformanceStatus,
    PlatformEnginePkgComplianceApiFinding,
    PlatformEnginePkgComplianceApiResourceFinding,
    TypesFindingSeverity,
} from '../../../../../../api/api'
import AxiosAPI from '../../../../../../api/ApiConfig'
import { activeBadge, statusBadge } from '../../../index'
import { dateTimeDisplay } from '../../../../../../utilities/dateDisplay'
import { getConnectorIcon } from '../../../../../../components/Cards/ConnectorCard'
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
import Filter from './Filter'
import dayjs from 'dayjs'
let sortKey: any[] = []

interface IControlFindings {
    onlyFailed?: boolean
    controlId: string | undefined
}

const columns = (isDemo: boolean) => {
    const temp: IColumn<
        PlatformEnginePkgComplianceApiFinding,
        any
    >[] = [
        {
            field: 'resourceName',
            headerName: 'Resource name',
            hide: false,
            type: 'string',
            enableRowGroup: true,
            sortable: false,
            filter: true,
            resizable: true,
            flex: 1,
            cellRenderer: (
                param: ICellRendererParams<
                    PlatformEnginePkgComplianceApiFinding,
                    any
                >
            ) => (
                <Flex flexDirection="col" alignItems="start">
                    <Text className="text-gray-800">
                        {param.data?.resourceName ||
                            (param.data?.stateActive === false
                                ? 'Resource deleted'
                                : '')}
                    </Text>
                    <Text className={isDemo ? 'blur-sm' : ''}>
                        {param.data?.platformResourceID}
                    </Text>
                </Flex>
            ),
        },
        {
            field: 'resourceType',
            headerName: 'Resource type',
            type: 'string',
            enableRowGroup: true,
            sortable: true,
            hide: true,
            filter: true,
            resizable: true,
            flex: 1,
            cellRenderer: (
                param: ICellRendererParams<
                    PlatformEnginePkgComplianceApiFinding,
                    any
                >
            ) => (
                <Flex flexDirection="col" alignItems="start">
                    <Text className="text-gray-800">
                        {param.data?.resourceTypeName}
                    </Text>
                    <Text>{param.data?.resourceType}</Text>
                </Flex>
            ),
        },
        {
            field: 'benchmarkID',
            headerName: 'Benchmark',
            type: 'string',
            enableRowGroup: true,
            sortable: false,
            hide: true,
            filter: true,
            resizable: true,
            flex: 1,
            cellRenderer: (
                param: ICellRendererParams<PlatformEnginePkgComplianceApiFinding>
            ) => (
                <Flex flexDirection="col" alignItems="start">
                    <Text className="text-gray-800">
                        {param.data?.parentBenchmarkNames?.at(0)}
                    </Text>
                    <Text>{param.value}</Text>
                </Flex>
            ),
        },
        {
            field: 'providerConnectionName',
            headerName: 'Account',
            type: 'string',
            enableRowGroup: true,
            hide: false,
            sortable: false,
            filter: true,
            resizable: true,
            flex: 1,
            cellRenderer: (
                param: ICellRendererParams<
                    PlatformEnginePkgComplianceApiFinding,
                    any
                >
            ) => (
                <Flex flexDirection="row" justifyContent="start">
                    {getConnectorIcon(param.data?.connector, '-ml-2 mr-2')}
                    <Flex
                        flexDirection="col"
                        alignItems="start"
                        className={isDemo ? 'blur-sm' : ''}
                    >
                        <Text className="text-gray-800">
                            {param.data?.providerConnectionName}
                        </Text>
                        <Text>{param.data?.providerConnectionID}</Text>
                    </Flex>
                </Flex>
            ),
        },
        {
            field: 'connectionID',
            headerName: 'opencomply connection ID',
            type: 'string',
            hide: true,
            enableRowGroup: true,
            sortable: false,
            filter: true,
            resizable: true,
            flex: 1,
        },
        {
            field: 'stateActive',
            headerName: 'State',
            type: 'string',
            sortable: true,
            filter: true,
            hide: false,
            resizable: true,
            flex: 1,
            cellRenderer: (param: ValueFormatterParams) => (
                <Flex className="h-full">{activeBadge(param.value)}</Flex>
            ),
        },
        {
            field: 'conformanceStatus',
            headerName: 'Conformance status',
            type: 'string',
            sortable: true,
            filter: true,
            hide: false,
            resizable: true,
            width: 160,
            cellRenderer: (param: ValueFormatterParams) => (
                <Flex className="h-full">{statusBadge(param.value)}</Flex>
            ),
        },
        {
            field: 'evaluatedAt',
            headerName: 'Last checked',
            type: 'datetime',
            sortable: false,
            filter: true,
            resizable: true,
            width: 200,
            valueFormatter: (param: ValueFormatterParams) => {
                return param.value ? dateTimeDisplay(param.value) : ''
            },
            hide: false,
        },
    ]
    return temp
}

export default function ControlFindings({
    controlId,
    onlyFailed,
}: IControlFindings) {
    const isDemo = useAtomValue(isDemoAtom)
    const setNotification = useSetAtom(notificationAtom)

    const [open, setOpen] = useState(false)
    const [finding, setFinding] = useState<
        PlatformEnginePkgComplianceApiFinding | undefined
    >(undefined)
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)
    const [rows, setRows] = useState<any[]>()
    const [page, setPage] = useState(1)
    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [queries, setQuery] = useState<{
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
        activeTimeRange: DateRange | undefined
        eventTimeRange: DateRange | undefined
        jobID: string[] | undefined
        connectionGroup: []
    }>({
        connector: [],
      conformanceStatus:
                                onlyFailed === true
                                    ? [
                                          PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
                                      ]
                                    : [
                                          PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed,
                                          PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
                                      ],
        severity: [
            TypesFindingSeverity.FindingSeverityCritical,
            TypesFindingSeverity.FindingSeverityHigh,
            TypesFindingSeverity.FindingSeverityMedium,
            TypesFindingSeverity.FindingSeverityLow,
            TypesFindingSeverity.FindingSeverityNone,
        ],
        connectionID: [],
        controlID: [controlId || ''],
        benchmarkID: [],
        resourceTypeID: [],
        lifecycle: [true],
        activeTimeRange: undefined,
        eventTimeRange: undefined,
        jobID: [],
        connectionGroup: [],
    })
    const today = new Date()
    const lastWeek = new Date(
        today.getFullYear(),
        today.getMonth(),
        today.getDate() - 7
    )

    const [date, setDate] = useState({
        key: 'previous-3-days',
        amount: 3,
        unit: 'day',
        type: 'relative',
    })
    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 40 ? text.substring(0, 40) + '...' : text
        }
    }

    // const ssr = () => {
    //     return {
    //         getRows: (params: IServerSideGetRowsParams) => {
    //             const api = new Api()
    //             api.instance = AxiosAPI
    //             api.compliance
    //                 .apiV1FindingsCreate({
    //                     filters: {
    //                         controlID: [controlId || ''],
    //                         conformanceStatus:
    //                             onlyFailed === true
    //                                 ? [
    //                                       PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
    //                                   ]
    //                                 : [
    //                                       PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed,
    //                                       PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
    //                                   ],
    //                     },
    //                     sort: params.request.sortModel.length
    //                         ? [
    //                               {
    //                                   [params.request.sortModel[0].colId]:
    //                                       params.request.sortModel[0].sort,
    //                               },
    //                           ]
    //                         : [],
    //                     limit: 100,
    //                     afterSortKey:
    //                         params.request.startRow === 0 || sortKey.length < 1
    //                             ? []
    //                             : sortKey,
    //                 })
    //                 .then((resp) => {
    //                     params.success({
    //                         rowData: resp.data.findings || [],
    //                         rowCount: resp.data.totalCount || 0,
    //                     })

    //                     sortKey =
    //                         resp.data.findings?.at(
    //                             (resp.data.findings?.length || 0) - 1
    //                         )?.sortKey || []
    //                 })
    //                 .catch((err) => {
    //                     if (
    //                         err.message !==
    //                         "Cannot read properties of null (reading 'NaN')"
    //                     ) {
    //                         setError(err.message)
    //                     }
    //                     params.fail()
    //                 })
    //         },
    //     }
    // }
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
        let body = {
            filters: {
                connector: queries.connector.length ? [query.connector] : [],
                controlID: [controlId || ''],
                connectionID: queries.connectionID,
                benchmarkID: queries.benchmarkID,
                severity: queries.severity,
                resourceTypeID: queries.resourceTypeID,
                conformanceStatus: queries.conformanceStatus,
                stateActive: queries.lifecycle,
                jobID: queries?.jobID,
                connectionGroup: queries.connectionGroup,
                ...(queries.eventTimeRange && {
                    lastEvent: {
                        from: queries.eventTimeRange.start.unix(),
                        to: queries.eventTimeRange.end.unix(),
                    },
                }),
                ...(!isRelative
                    ? {
                          evaluatedAt: {
                              from: start.unix(),
                              to: end.unix(),
                          },
                      }
                    : {
                          interval: relative,
                      }),
            },
            // sort: params.request.sortModel.length
            //     ? [
            //           {
            //               [params.request.sortModel[0].colId]:
            //                   params.request.sortModel[0].sort,
            //           },
            //       ]
            //     : [],
            limit: 10,
            // eslint-disable-next-line prefer-destructuring,@typescript-eslint/ban-ts-comment
            // @ts-ignore
            afterSortKey: page == 1 ? [] : rows[rows?.length - 1].sortKey,
            //     params.request.startRow === 0 ||
            //     sortKey.length < 1 ||
            //     sortKey === 'none'
            //         ? []
            //         : sortKey,
        }

        api.compliance
            .apiV1FindingsCreate(body)
            .then((resp) => {
                setLoading(false)
                // @ts-ignore

                setTotalPage(Math.ceil(resp.data.totalCount / 10))
                // @ts-ignore

                setTotalCount(resp.data.totalCount)
                // @ts-ignore
                if (resp.data.findings) {
                    setRows(resp.data.findings)
                } else {
                    setRows([])
                }

                // eslint-disable-next-line prefer-destructuring,@typescript-eslint/ban-ts-comment
                // @ts-ignore
            })
            .catch((err) => {
                setLoading(false)
                if (
                    err.message !==
                    "Cannot read properties of null (reading 'NaN')"
                ) {
                    setError(err.message)
                }
                setNotification({
                    text: 'Can not Connect to Server',
                    type: 'warning',
                })
            })
    }

    useEffect(() => {
        GetRows()
    }, [page, queries, date])
    // const serverSideRows = useMemo(() => ssr(), [onlyFailed])

    return (
        <>
            {error.length > 0 && (
                <Flex className="w-fit mb-3 gap-1">
                    <ExclamationCircleIcon className="text-rose-600 h-5" />
                    <Text color="rose">{error}</Text>
                </Flex>
            )}
            <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                className="w-full"
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
                                'Incident not selected'
                            )
                        }
                    >
                        <FindingDetail
                            type="finding"
                            finding={finding}
                            open={open}
                            onClose={() => setOpen(false)}
                            onRefresh={() => window.location.reload()}
                        />
                    </SplitPanel>
                }
                content={
                    <KTable
                        className="   min-h-[450px]"
                        // resizableColumns
                        variant="full-page"
                        renderAriaLive={({
                            firstIndex,
                            lastIndex,
                            totalItemsCount,
                        }) =>
                            `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                        }
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
                            if (
                                row.platformResourceID 
                               
                            ) {
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
                                id: 'providerConnectionName',
                                header: 'Cloud Account',
                                cell: (item) => item.providerConnectionID,
                                sortingField: 'id',
                                isRowHeader: true,
                            },
                            {
                                id: 'resourceName',
                                header: 'Resource Name',
                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className={
                                                isDemo ? 'h-full' : 'h-full'
                                            }
                                        >
                                            <Text className="text-gray-800">
                                                {truncate(item.resourceName)}
                                            </Text>
                                            <Text>
                                                {truncate(
                                                    item.resourceTypeName
                                                )}
                                            </Text>
                                        </Flex>
                                    </>
                                ),
                                sortingField: 'title',
                                // minWidth: 400,
                                maxWidth: 100,
                            },
                            {
                                id: 'benchmarkID',
                                header: 'Benchmark',
                                maxWidth: 120,
                                cell: (item) => (
                                    <>
                                        <Text className="text-gray-800">
                                            {truncate(
                                                item?.parentBenchmarkNames?.at(
                                                    0
                                                )
                                            )}
                                        </Text>
                                        {/* <Text>
                                            {truncate(
                                                item?.parentBenchmarkNames?.at(
                                                    (item?.parentBenchmarkNames
                                                        ?.length || 0) - 1
                                                )
                                            )}
                                        </Text> */}
                                    </>
                                ),
                            },
                            {
                                id: 'controlID',
                                header: 'Control',
                                maxWidth: 160,

                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className="h-full"
                                        >
                                            <Text className="text-gray-800">
                                                {truncate(item?.controlTitle)}
                                                {/* {truncate(
                                                    item?.parentBenchmarkNames?.at(
                                                        0
                                                    )
                                                )} */}
                                            </Text>
                                            {/* <Text></Text> */}
                                        </Flex>
                                    </>
                                ),
                            },
                            {
                                id: 'conformanceStatus',
                                header: 'Status',
                                sortingField: 'severity',
                                cell: (item) => (
                                    <Badge
                                        // @ts-ignore
                                        color={`${
                                            item.conformanceStatus == 'passed'
                                                ? 'green'
                                                : 'red'
                                        }`}
                                    >
                                        {item.conformanceStatus}
                                    </Badge>
                                ),
                                maxWidth: 50,
                            },
                            {
                                id: 'severity',
                                header: 'Severity',
                                sortingField: 'severity',
                                cell: (item) => (
                                    <Badge
                                        // @ts-ignore
                                        color={`severity-${item.severity}`}
                                    >
                                        {item.severity.charAt(0).toUpperCase() +
                                            item.severity.slice(1)}
                                    </Badge>
                                ),
                                maxWidth: 50,
                            },
                            {
                                id: 'evaluatedAt',
                                header: 'Last Evaluation',
                                maxWidth: 100,

                                cell: (item) => (
                                    // @ts-ignore
                                    <>{dateTimeDisplay(item.evaluatedAt)}</>
                                ),
                            },
                        ]}
                        columnDisplay={[
                            { id: 'resourceName', visible: true },
                            { id: 'benchmarkID', visible: true },
                            { id: 'controlID', visible: false },
                            { id: 'conformanceStatus', visible: true },
                            { id: 'severity', visible: true },
                            { id: 'evaluatedAt', visible: true },

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
                                className="gap-1 mt-1"
                            >
                                <Filter
                                    // @ts-ignore
                                    type={'findings'}
                                    onApply={(e) => {
                                        // @ts-ignore
                                        setQuery(e)
                                    }}
                                    setDate={setDate}
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
                                            key: 'previous-7-days',
                                            amount: 7,
                                            unit: 'day',
                                            type: 'relative',
                                        },
                                    ]}
                                    hideTimeOffset
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
                                Incidents{' '}
                                <span className=" font-medium">
                                    ({totalCount})
                                </span>
                            </Header>
                        }
                        pagination={
                            <Pagination
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
