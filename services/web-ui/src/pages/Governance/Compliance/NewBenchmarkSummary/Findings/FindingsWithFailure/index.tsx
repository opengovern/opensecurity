// @ts-nocheck
import { Card, Flex, Text, Title } from '@tremor/react'
import { useEffect, useMemo, useState } from 'react'

import { useAtomValue, useSetAtom } from 'jotai'
import { isDemoAtom, notificationAtom } from '../../../../../../store'
import Table, { IColumn } from '../../../../../../components/Table'
import { dateTimeDisplay } from '../../../../../../utilities/dateDisplay'
import {
    Api,
    GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus,
    GithubComKaytuIoKaytuEnginePkgComplianceApiFinding,
    SourceType,
    TypesFindingSeverity,
} from '../../../../../../api/api'
import AxiosAPI from '../../../../../../api/ApiConfig'

// import { severityBadge, statusBadge } from '../../Controls'
import { getConnectorIcon } from '../../../../../../components/Cards/ConnectorCard'
import { DateRange } from '../../../../../../utilities/urlstate'
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
import FindingDetail from './Detail'
import dayjs from 'dayjs'


let sortKey = ''

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
        eventTimeRange: DateRange | undefined
        jobID: string[] | undefined
    }
}

export default function FindingsWithFailure({ query }: ICount) {
    const setNotification = useSetAtom(notificationAtom)

    const [open, setOpen] = useState(false)
    const [loading, setLoading] = useState(false)
    const [finding, setFinding] = useState<
        GithubComKaytuIoKaytuEnginePkgComplianceApiFinding | undefined
    >(undefined)
    const [rows, setRows] = useState<any[]>()
    const [page, setPage] = useState(1)
    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const isDemo = useAtomValue(isDemoAtom)
     const [queries, setQuery] = useState(query)
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
                 relative = `${date.amount} ${date.unit}s`
             } else {
                 // @ts-ignore

                 start = dayjs(date?.startDate)
                 // @ts-ignore

                 end = dayjs(date?.endDate)
             }
         }
        api.compliance
            .apiV1FindingsCreate({
                filters: {
                    integrationType: queries.connector.length
                        ? [query.connector]
                        : [],
                    controlID: queries.controlID,
                    integrationID: queries.connectionID,
                    benchmarkID: query.benchmarkID,
                    severity: queries.severity,
                    resourceTypeID: queries.resourceTypeID,
                    complianceStatus: queries.conformanceStatus,
                    integrationGroup: queries.connectionGroup,
                    stateActive: queries.lifecycle,
                    jobID: queries?.jobID,
                    ...(queries.eventTimeRange && {
                        lastEvent: {
                            from: queries.eventTimeRange.start.unix(),
                            to: queries.eventTimeRange.end.unix(),
                        },
                    }),
                    ...(!isRelative &&
                        date && {
                            evaluatedAt: {
                                from: start?.unix(),
                                to: end?.unix(),
                            },
                        }),
                    ...(isRelative &&
                        date && {
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
            })
            .then((resp) => {
                setLoading(false)
                // @ts-ignore

                setTotalPage(Math.ceil(resp.data.totalCount / 10))
                // @ts-ignore

                setTotalCount(resp.data.totalCount)
                // @ts-ignore
                if (resp.data.complianceResults) {
                    setRows(resp.data.complianceResults)
                } else {
                    setRows([])
                }

                // eslint-disable-next-line prefer-destructuring,@typescript-eslint/ban-ts-comment
                // @ts-ignore
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
    }, [page,queries,date])
    return (
        <>
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
                                        {getConnectorIcon(finding?.integrationType)}
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
                            if (row.platformResourceID) {
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
                                cell: (item) => item.integrationID,
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
                                                {item.resourceName}
                                            </Text>
                                            <Text>{item.resourceTypeName}</Text>
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
                                maxWidth: 100,
                                cell: (item) => (
                                    <>
                                        <Text className="text-gray-800">
                                            {truncate(
                                                item?.parentBenchmarkNames?.at(
                                                    0
                                                )
                                            )}
                                        </Text>
                                        <Text>
                                            {truncate(
                                                item?.parentBenchmarkNames?.at(
                                                    (item?.parentBenchmarkNames
                                                        ?.length || 0) - 1
                                                )
                                            )}
                                        </Text>
                                    </>
                                ),
                            },
                            {
                                id: 'controlID',
                                header: 'Control',
                                maxWidth: 100,

                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className="h-full"
                                        >
                                            <Text className="text-gray-800">
                                                {truncate(
                                                    item?.parentBenchmarkNames?.at(
                                                        0
                                                    )
                                                )}
                                            </Text>
                                            <Text>
                                                {truncate(item?.controlTitle)}
                                            </Text>
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
                                            item.complianceStatus == 'passed'
                                                ? 'green'
                                                : 'red'
                                        }`}
                                    >
                                        {item.complianceStatus}
                                    </Badge>
                                ),
                                maxWidth: 100,
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
                                maxWidth: 100,
                            },
                            {
                                id: 'evaluatedAt',
                                header: 'Last Evaluation',
                                cell: (item) => (
                                    // @ts-ignore
                                    <>{dateTimeDisplay(item.evaluatedAt)}</>
                                ),
                            },
                        ]}
                        columnDisplay={[
                            { id: 'resourceName', visible: true },
                            { id: 'benchmarkID', visible: true },
                            { id: 'controlID', visible: true },
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
