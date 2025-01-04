// @ts-nocheck
import { Card, Flex, Text } from '@tremor/react'
import { useNavigate } from 'react-router-dom'
import { ICellRendererParams, RowClickedEvent } from 'ag-grid-community'
import { useAtomValue } from 'jotai'
import { useComplianceApiV1FindingsTopDetail } from '../../../../api/compliance.gen'
import {
    GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../../../api/api'
import Table, { IColumn } from '../../../../components/Table'
import { severityBadge } from '../../Controls'
import { DateRange, searchAtom } from '../../../../utilities/urlstate'
import { useEffect, useState } from 'react'
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
import Filter from '../Filter'
import dayjs from 'dayjs'

const policyColumns: IColumn<any, any>[] = [
    {
        headerName: 'Control',
        field: 'title',
        type: 'string',
        sortable: true,
        filter: true,
        resizable: true,
        cellRenderer: (param: ICellRendererParams) => (
            <Flex
                flexDirection="col"
                alignItems="start"
                justifyContent="center"
                className="h-full"
            >
                <Text className="text-gray-800">{param.value}</Text>
                <Text>{param.data.id}</Text>
            </Flex>
        ),
    },
    {
        headerName: 'Severity',
        field: 'sev',
        width: 120,
        type: 'string',
        sortable: true,
        filter: true,
        resizable: true,
        cellRenderer: (params: ICellRendererParams) => (
            <Flex
                className="h-full w-full"
                justifyContent="center"
                alignItems="center"
            >
                {severityBadge(params.data.severity)}
            </Flex>
        ),
    },
    {
        headerName: 'Findings',
        field: 'count',
        type: 'number',
        sortable: true,
        filter: true,
        resizable: true,
        width: 150,
        cellRenderer: (param: ICellRendererParams) => (
            <Flex
                flexDirection="col"
                alignItems="start"
                justifyContent="center"
                className="h-full"
            >
                <Text className="text-gray-800">{param.value || 0} issues</Text>
                <Text>
                    {(param.data.totalCount || 0) - (param.value || 0)} passed
                </Text>
            </Flex>
        ),
    },
    {
        headerName: 'Resources',
        field: 'resourceCount',
        type: 'number',
        sortable: true,
        filter: true,
        resizable: true,
        width: 150,
        cellRenderer: (param: ICellRendererParams) => (
            <Flex
                flexDirection="col"
                alignItems="start"
                justifyContent="center"
                className="h-full"
            >
                <Text className="text-gray-800">{param.value || 0} issues</Text>
                <Text>
                    {(param.data.resourceTotalCount || 0) - (param.value || 0)}{' '}
                    passed
                </Text>
            </Flex>
        ),
    },
]

interface ICount {
    query: {
        connector: SourceType
        conformanceStatus:
            | GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus[]
            | undefined
        severity: TypesFindingSeverity[] | undefined
        connectionID: string[] | undefined
        controlID: string[] | undefined
        benchmarkID: string[] | undefined
        resourceTypeID: string[] | undefined
        lifecycle: boolean[] | undefined
        activeTimeRange: DateRange | undefined
    }
}

export default function ControlsWithFailure({ query }: ICount) {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const [queries, setQuery] = useState(query)

    const topQuery = {
        connector: query.connector.length ? [query.connector] : [],
        connectionId: query.connectionID,
        benchmarkId: query.benchmarkID,
    }
    const [date, setDate] = useState({
        key: 'previous-3-days',
        amount: 3,
        unit: 'day',
        type: 'relative',
    })
    const {
        response: controls,
        isLoading,
        sendNowWithParams: GetRow,
    } = useComplianceApiV1FindingsTopDetail(
        'controlID',
        10000,
        {
            connector: queries.connector.length ? queries.connector : [],
            severities: queries?.severity,
            integrationID: queries.connectionID,
            integrationGroup: queries?.connectionGroup,
        },
        {},
        false
    )
    const [page, setPage] = useState(0)

    useEffect(() => {
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
        GetRow('controlID', 10000, {
            integrationType: queries.connector.length ? queries.connector : [],
            severities: queries?.severity,
            connectionId: queries.connectionID,
            connectionGroup: queries?.connectionGroup,
            ...(!isRelative &&
                date && {
                    startTime: start?.unix(),
                    endTime: end?.unix(),
                }),
            ...(isRelative &&
                date && {
                    interval: relative,
                }),
        })
    }, [queries, date])
    return (
        <KTable
            className="p-3   min-h-[450px]"
            // resizableColumns
            renderAriaLive={({ firstIndex, lastIndex, totalItemsCount }) =>
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
                // const row = event.detail.item
                // if (row) {
                //     navigate(`${row?.Control.id}?${searchParams}`)
                // }
            }}
            columnDefinitions={[
                {
                    id: 'title',
                    header: 'Control',
                    cell: (item) => (
                        <>
                            <Flex
                                flexDirection="col"
                                alignItems="start"
                                justifyContent="center"
                                className="h-full"
                            >
                                <Text className="text-gray-800">
                                    <Link
                                        href={`${window.location}/${item.Control.id}`}
                                        target="__blank"
                                    >
                                        {item.Control.title}
                                    </Link>
                                </Text>
                                <Text>{item.Control.id}</Text>
                            </Flex>
                        </>
                    ),
                    sortingField: 'id',
                    isRowHeader: true,
                    maxWidth: 300,
                },
                {
                    id: 'severity',
                    header: 'Severity',
                    sortingField: 'severity',
                    cell: (item) => (
                        <Badge
                            // @ts-ignore
                            color={`severity-${item.Control.severity}`}
                        >
                            {item.Control.severity.charAt(0).toUpperCase() +
                                item.Control.severity.slice(1)}
                        </Badge>
                    ),
                    maxWidth: 100,
                },
                {
                    id: 'count',
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
                                <Text className="text-gray-800">{`${item.count} Incidents`}</Text>
                                <Text>{`${
                                    item.totalCount - item.count
                                } passed`}</Text>
                            </Flex>
                        </>
                    ),
                },
                {
                    id: 'resourceCount',
                    header: 'Impacted Resources',
                    cell: (item) => (
                        <>
                            <Flex
                                flexDirection="col"
                                alignItems="start"
                                justifyContent="center"
                                className="h-full"
                            >
                                <Text className="text-gray-800">
                                    {item.resourceCount || 0} failing
                                </Text>
                                <Text>
                                    {(item.resourceTotalCount || 0) -
                                        (item.resourceCount || 0)}{' '}
                                    passing
                                </Text>
                            </Flex>
                        </>
                    ),
                    sortingField: 'title',
                    // minWidth: 400,
                    maxWidth: 200,
                },
                // {
                //     id: 'providerConnectionName',
                //     header: 'Cloud account',
                //     maxWidth: 100,
                //     cell: (item) => (
                //         <>
                //             <Flex
                //                 justifyContent="start"
                //                 className={`h-full gap-3 group relative ${
                //                     isDemo ? 'blur-sm' : ''
                //                 }`}
                //             >
                //                 {getConnectorIcon(item.connector)}
                //                 <Flex flexDirection="col" alignItems="start">
                //                     <Text className="text-gray-800">
                //                         {item.providerConnectionName}
                //                     </Text>
                //                     <Text>{item.providerConnectionID}</Text>
                //                 </Flex>
                //                 <Card className="cursor-pointer absolute w-fit h-fit z-40 right-1 scale-0 transition-all py-1 px-4 group-hover:scale-100">
                //                     <Text color="blue">Open</Text>
                //                 </Card>
                //             </Flex>
                //         </>
                //     ),
                // },

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
                { id: 'title', visible: true },
                { id: 'severity', visible: true },
                // { id: 'count', visible: true },
                { id: 'resourceCount', visible: true },
                // { id: 'severity', visible: true },
                // { id: 'evaluatedAt', visible: true },

                // { id: 'action', visible: true },
            ]}
            enableKeyboardNavigation
            // @ts-ignore
            items={controls?.records?.slice(page * 10, (page + 1) * 10)}
            loading={isLoading}
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
                        type={'controls'}
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
                    Controls{' '}
                    <span className=" font-medium">
                        ({controls?.totalCount})
                    </span>
                </Header>
            }
            pagination={
                <Pagination
                    currentPageIndex={page + 1}
                    pagesCount={Math.ceil(controls?.totalCount / 10)}
                    onChange={({ detail }) =>
                        setPage(detail.currentPageIndex - 1)
                    }
                />
            }
        />
    )
}
