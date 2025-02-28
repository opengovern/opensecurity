// @ts-nocheck
import { Card, Flex, Text } from '@tremor/react'
import { useNavigate } from 'react-router-dom'
import { useAtomValue } from 'jotai'
import { useComplianceApiV1FindingsTopDetail } from '../../../../../../../../api/compliance.gen'
import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../../../../../../../api/api'

import {
    DateRange,
    searchAtom,
} from '../../../../../../../../utilities/urlstate'
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
import CustomPagination from '../../../../../../../../components/Pagination'

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
        activeTimeRange: DateRange | undefined
    }
    id: string
}

export default function ControlsWithFailure({ query ,id}: ICount) {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const [queries, setQuery] = useState(query)

    const topQuery = {
        connector: query.connector.length ? [query.connector] : [],
        connectionId: query.connectionID,
        benchmarkId: query.benchmarkID,
    }
  
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
      
        GetRow('controlID', 10000, {
            integrationTypes: queries.connector.length ? queries.connector : [],
            severities: queries?.severity,
            connectionId: queries.connectionID,
            connectionGroup: queries?.connectionGroup,
            jobId: [id],
           
        })
    }, [queries])
    return (
        <div
            className="w-full"
            style={
                window.innerWidth < 768
                    ? { width: `${window.innerWidth - 80}px` }
                    : {}
            }
        >
            <KTable
                className="p-3   min-h-[450px]"
                // resizableColumns
                renderAriaLive={({ firstIndex, lastIndex, totalItemsCount }) =>
                    `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                }
                variant='full-page'
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
                                            href={`/incidents/${item.Control.id}`}
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
                        className="gap-1 mt-1 sm:flex-row flex-col"
                    >
                        <Filter
                            // @ts-ignore
                            type={'controls'}
                            onApply={(e) => {
                                // @ts-ignore
                                setQuery(e)
                            }}
                            setDate={()=>{}}
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
                    <CustomPagination
                        currentPageIndex={page + 1}
                        pagesCount={Math.ceil(controls?.totalCount / 10)}
                        onChange={({ detail }) =>
                            setPage(detail.currentPageIndex - 1)
                        }
                    />
                }
            />
        </div>
    )
}
