// @ts-nocheck
import { useAtomValue } from 'jotai'
import { Card, Flex, Text } from '@tremor/react'
import { useEffect, useState } from 'react'
import { useComplianceApiV1FindingsTopDetail } from '../../../../api/compliance.gen'
import {
    GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../../../api/api'
import { topConnections } from '../../Controls/ControlSummary/Tabs/ImpactedAccounts'
import { getConnectorIcon } from '../../../../components/Cards/ConnectorCard'
import CloudAccountDetail from './Detail'
import { isDemoAtom } from '../../../../store'
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

import {
    AppLayout,
    Container,
    ContentLayout,
    SplitPanel,
} from '@cloudscape-design/components'
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
    }
}

export default function FailingCloudAccounts({ query }: ICount) {
    const [account, setAccount] = useState<any>(undefined)
    const [open, setOpen] = useState(false)
    const isDemo = useAtomValue(isDemoAtom)

    const topQuery = {
        connector: query.connector.length ? [query.connector] : [],
        connectionId: query.connectionID,
        benchmarkId: query.benchmarkID,
    }

    const { response: accounts, isLoading: accountsLoading } =
        useComplianceApiV1FindingsTopDetail('integrationID', 10000, topQuery)
    const [page, setPage] = useState(0)

    return (
        <>
            <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                toolsHide={true}
                navigationHide={true}
                splitPanelOpen={open}
                onSplitPanelToggle={()=>{
                    setOpen(!open)
                }}
                splitPanel={
                    <SplitPanel header="Split panel header">
                        Split panel content
                    </SplitPanel>
                }
                content={
                    <KTable
                        className="  min-h-[450px]"
                        variant="full-page"
                        // resizableColumns
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
                            setAccount(row)
                            setOpen(true)
                        }}
                        columnDefinitions={[
                            {
                                id: 'providerConnectionName',
                                header: 'Account name',
                                cell: (item) => (
                                    <>
                                        <Flex
                                            justifyContent="start"
                                            className={`h-full gap-3 group relative ${
                                                isDemo ? 'blur-sm' : ''
                                            }`}
                                        >
                                            {getConnectorIcon(item.connector)}
                                            <Flex
                                                flexDirection="col"
                                                alignItems="start"
                                            >
                                                <Text className="text-gray-800">
                                                    {
                                                        item.providerConnectionName
                                                    }
                                                </Text>
                                                <Text>
                                                    {item.providerConnectionID}
                                                </Text>
                                            </Flex>
                                            <Card className="cursor-pointer absolute w-fit h-fit z-40 right-1 scale-0 transition-all py-1 px-4 group-hover:scale-100">
                                                <Text color="blue">Open</Text>
                                            </Card>
                                        </Flex>
                                    </>
                                ),
                                sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 300,
                            },
                            // {
                            //     id: 'severity',
                            //     header: 'Severity',
                            //     sortingField: 'severity',
                            //     cell: (item) => (
                            //         <Badge
                            //             // @ts-ignore
                            //             color={`severity-${item.Control.severity}`}
                            //         >
                            //             {item.Control.severity.charAt(0).toUpperCase() +
                            //                 item.Control.severity.slice(1)}
                            //         </Badge>
                            //     ),
                            //     maxWidth: 100,
                            // },
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
                                            <Text className="text-gray-800">{`${item.count} issues`}</Text>
                                            <Text>{`${
                                                item.totalCount - item.count
                                            } passed`}</Text>
                                        </Flex>
                                    </>
                                ),
                            },
                            {
                                id: 'resourceCount',
                                header: 'Resources',
                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className="h-full"
                                        >
                                            <Text className="text-gray-800">
                                                {item.resourceCount || 0} issues
                                            </Text>
                                            <Text>
                                                {(item.resourceTotalCount ||
                                                    0) -
                                                    (item.resourceCount ||
                                                        0)}{' '}
                                                passed
                                            </Text>
                                        </Flex>
                                    </>
                                ),
                                sortingField: 'title',
                                // minWidth: 400,
                                maxWidth: 200,
                            },
                            {
                                id: 'controlCount',
                                header: 'Controls',
                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className="h-full"
                                        >
                                            <Text className="text-gray-800">
                                                {item.controlCount || 0} issues
                                            </Text>
                                            <Text>
                                                {(item.controlTotalCount || 0) -
                                                    (item.controlCount ||
                                                        0)}{' '}
                                                passed
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
                            { id: 'providerConnectionName', visible: true },
                            { id: 'count', visible: true },
                            { id: 'resourceCount', visible: true },
                            { id: 'controlCount', visible: true },

                            // { id: 'severity', visible: true },
                            // { id: 'evaluatedAt', visible: true },

                            // { id: 'action', visible: true },
                        ]}
                        enableKeyboardNavigation
                        // @ts-ignore
                        items={accounts?.records?.slice(
                            page * 10,
                            (page + 1) * 10
                        )}
                        loading={accountsLoading}
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
                                Accounts{' '}
                                <span className=" font-medium">
                                    ({accounts?.totalCount})
                                </span>
                            </Header>
                        }
                        pagination={
                            <Pagination
                                currentPageIndex={page + 1}
                                pagesCount={Math.ceil(
                                    accounts?.totalCount / 10
                                )}
                                onChange={({ detail }) =>
                                    setPage(detail.currentPageIndex - 1)
                                }
                            />
                        }
                    />
                }
            />
            <CloudAccountDetail
                account={account}
                open={open}
                onClose={() => setOpen(false)}
            />
        </>
    )
}
