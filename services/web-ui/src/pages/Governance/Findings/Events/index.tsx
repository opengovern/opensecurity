import { Card, Flex, Text, Title } from '@tremor/react'
import { useEffect, useMemo, useState } from 'react'

import { useAtomValue, useSetAtom } from 'jotai'
import { ArrowRightIcon } from '@heroicons/react/24/outline'
import { isDemoAtom, notificationAtom } from '../../../../store'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import {
    Api,
    GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus,
    GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEvent,
    SourceType,
    TypesFindingSeverity,
} from '../../../../api/api'
import AxiosAPI from '../../../../api/ApiConfig'
import { severityBadge, statusBadge } from '../../Controls'
import { getConnectorIcon } from '../../../../components/Cards/ConnectorCard'
import { DateRange } from '../../../../utilities/urlstate'
import EventDetail from './Detail'
import KTable from '@cloudscape-design/components/table'
import Box from '@cloudscape-design/components/box'
import SpaceBetween from '@cloudscape-design/components/space-between'
import Badge from '@cloudscape-design/components/badge'
import {
    BreadcrumbGroup,
    Header,
    Link,
    Pagination,
    PropertyFilter,
} from '@cloudscape-design/components'
import { AppLayout, SplitPanel } from '@cloudscape-design/components'




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
    }
}

export default function Events({ query }: ICount) {
    const setNotification = useSetAtom(notificationAtom)
    const [loading, setLoading] = useState(false)
    const [open, setOpen] = useState(false)
    const [finding, setFinding] = useState<
        GithubComKaytuIoKaytuEnginePkgComplianceApiFindingEvent | undefined
    >(undefined)
    const [rows, setRows] = useState<any[]>()
    const [page, setPage] = useState(1)
    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const isDemo = useAtomValue(isDemoAtom)

    const GetRows = () => {
        setLoading(true)
        const api = new Api()
        api.instance = AxiosAPI
        api.compliance
            .apiV1FindingEventsCreate({
                filters: {
                    connector: query.connector.length ? [query.connector] : [],
                    controlID: query.controlID,
                    connectionID: query.connectionID,
                    benchmarkID: query.benchmarkID,
                    severity: query.severity,
                    resourceType: query.resourceTypeID,
                    conformanceStatus: query.conformanceStatus,
                    stateActive: query.lifecycle,
                    ...(query.activeTimeRange && {
                        evaluatedAt: {
                            from: query.activeTimeRange.start.unix(),
                            to: query.activeTimeRange.end.unix(),
                        },
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
                // limit: 100,
                // // eslint-disable-next-line prefer-destructuring,@typescript-eslint/ban-ts-comment
                // @ts-ignore
                afterSortKey: page == 1 ? [] : rows[rows?.length - 1].sortKey,
                // afterSortKey:
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

                setRows(resp.data.findingEvents)
                // eslint-disable-next-line prefer-destructuring,@typescript-eslint/ban-ts-comment
                // @ts-ignore
                // sortKey =
                //     // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                //     // @ts-ignore
                //     // eslint-disable-next-line no-unsafe-optional-chaining
                //     resp.data.findingEvents[
                //         // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                //         // @ts-ignore
                //         // eslint-disable-next-line no-unsafe-optional-chaining
                //         resp.data.findingEvents?.length - 1
                //     ].sortKey
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
    }, [page])
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
                                            {finding?.providerConnectionName}
                                        </Title>
                                    </Flex>
                                </>
                            ) : (
                                'Event not selected'
                            )
                        }
                    >
                        <EventDetail
                            event={finding}
                            open={open}
                            onClose={() => setOpen(false)}
                        />
                    </SplitPanel>
                }
                content={
                    <KTable
                        className="p-3   min-h-[450px]"
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
                                id: 'id',
                                header: 'Event ID',
                                cell: (item) => item.id,
                                sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 200,
                            },
                            {
                                id: 'resourceType',
                                header: 'Resource info',
                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className={
                                                isDemo
                                                    ? 'h-full blur-sm'
                                                    : 'h-full'
                                            }
                                        >
                                            <Text className="text-gray-800">
                                                {item.resourceType}
                                            </Text>
                                            <Text>{item.resourceID}</Text>
                                        </Flex>
                                    </>
                                ),
                                sortingField: 'title',
                                // minWidth: 400,
                                maxWidth: 200,
                            },
                            {
                                id: 'conformanceStatus',
                                header: 'State Change',
                                maxWidth: 100,
                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="row"
                                            className="h-full w-fit gap-2"
                                        >
                                            <Badge
                                                // @ts-ignore
                                                color={`${
                                                    item.previousConformanceStatus ==
                                                    'passed'
                                                        ? 'green'
                                                        : 'red'
                                                }`}
                                            >
                                                {item.previousConformanceStatus}
                                            </Badge>

                                            <ArrowRightIcon className="w-5" />
                                            <Badge
                                                // @ts-ignore
                                                color={`${
                                                    item.conformanceStatus ==
                                                    'passed'
                                                        ? 'green'
                                                        : 'red'
                                                }`}
                                            >
                                                {item.conformanceStatus}
                                            </Badge>
                                        </Flex>
                                    </>
                                ),
                            },
                            // {
                            //     id: 'controlID',
                            //     header: 'Control',
                            //     maxWidth: 100,

                            //     cell: (item) => (
                            //         <>
                            //             <Flex
                            //                 flexDirection="col"
                            //                 alignItems="start"
                            //                 justifyContent="center"
                            //                 className="h-full"
                            //             >
                            //                 <Text className="text-gray-800">
                            //                     {truncate(
                            //                         item?.parentBenchmarkNames?.at(0)
                            //                     )}
                            //                 </Text>
                            //                 <Text>{truncate(item?.controlTitle)}</Text>
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
                            { id: 'id', visible: true },
                            { id: 'resourceType', visible: true },
                            // { id: 'controlID', visible: true },
                            { id: 'conformanceStatus', visible: true },
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
                            ''
                            // <PropertyFilter
                            //     // @ts-ignore
                            //     query={undefined}
                            //     // @ts-ignore
                            //     onChange={({ detail }) => {
                            //         // @ts-ignore
                            //         setQueries(detail)
                            //     }}
                            //     // countText="5 matches"
                            //     enableTokenGroups
                            //     expandToViewport
                            //     filteringAriaLabel="Control Categories"
                            //     // @ts-ignore
                            //     // filteringOptions={filters}
                            //     filteringPlaceholder="Control Categories"
                            //     // @ts-ignore
                            //     filteringOptions={undefined}
                            //     // @ts-ignore

                            //     filteringProperties={undefined}
                            //     // filteringProperties={
                            //     //     filterOption
                            //     // }
                            // />
                        }
                        header={
                            <Header className="w-full">
                                Events{' '}
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
