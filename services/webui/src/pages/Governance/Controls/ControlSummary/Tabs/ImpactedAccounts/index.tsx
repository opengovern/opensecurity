import { useAtomValue } from 'jotai'
import { Flex, Text } from '@tremor/react'
import { useComplianceApiV1FindingsTopDetail } from '../../../../../../api/compliance.gen'
import { PlatformEnginePkgComplianceApiGetTopFieldResponse } from '../../../../../../api/api'
import { isDemoAtom } from '../../../../../../store'
import { useState } from 'react'
import {
    AppLayout,
    Container,
    ContentLayout,
    SplitPanel,
} from '@cloudscape-design/components'
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
import CustomPagination from '../../../../../../components/Pagination'


interface IImpactedAccounts {
    controlId: string | undefined
}



export const topConnections = (
    input:
        | PlatformEnginePkgComplianceApiGetTopFieldResponse
        | undefined
) => {
    const data = []
    if (input && input.records) {
        for (let i = 0; i < input.records.length; i += 1) {
            data.push({
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-ignore
                ...input.records[i].Connection,
                count: input.records[i].count,
                totalCount: input.records[i].totalCount,
                resourceCount: input.records[i].resourceCount,
                resourceTotalCount: input.records[i].resourceTotalCount,
                controlCount: input.records[i].controlCount,
                controlTotalCount: input.records[i].controlTotalCount,
            })
        }
    }
    return data
}

export default function ImpactedAccounts({ controlId }: IImpactedAccounts) {
    const isDemo = useAtomValue(isDemoAtom)
    const { response: accounts, isLoading: accountsLoading } =
        useComplianceApiV1FindingsTopDetail('integrationID', 10000, {
            controlId: [String(controlId)],
        })
         const [page, setPage] = useState(0)
    return (
        <>
            <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                toolsHide={true}
                navigationHide={true}
                // splitPanelOpen={open}
                // onSplitPanelToggle={() => {
                //     setOpen(!open)
                // }}
                // splitPanel={
                //     <SplitPanel header="Split panel header">
                //         Split panel content
                //     </SplitPanel>
                // }
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
                            // setAccount(row)
                            // setOpen(true)
                        }}
                        columnDefinitions={[
                            {
                                id: 'providerConnectionName',
                                header: 'Integration name',
                                cell: (item) => (
                                    // @ts-ignore
                                    <>
                                        {/** @ts-ignore */}
                                        {
                                            // @ts-ignore

                                            item?.Integration?.name
                                        }
                                    </>
                                ),
                                sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 300,
                            },
                            {
                                id: 'providerConnectionID',
                                header: 'Provider ID',
                                cell: (item) => (
                                    // @ts-ignore
                                    <>
                                        {
                                            // @ts-ignore
                                            item?.Integration?.provider_id
                                        }
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
                                header: 'Resources',
                                maxWidth: 100,

                                cell: (item) => (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            justifyContent="center"
                                            className="h-full"
                                        >
                                            <Text className="text-gray-800">{`${item.count} incidents`}</Text>
                                            <Text>{`${
                                                // @ts-ignore
                                                item.totalCount - item.count
                                            } passed`}</Text>
                                        </Flex>
                                    </>
                                ),
                            },
                            // {
                            //     id: 'resourceCount',
                            //     header: 'Resources',
                            //     cell: (item) => (
                            //         <>
                            //             <Flex
                            //                 flexDirection="col"
                            //                 alignItems="start"
                            //                 justifyContent="center"
                            //                 className="h-full"
                            //             >
                            //                 <Text className="text-gray-800">
                            //                     {item.resourceCount || 0} issues
                            //                 </Text>
                            //                 <Text>
                            //                     {(item.resourceTotalCount ||
                            //                         0) -
                            //                         (item.resourceCount ||
                            //                             0)}{' '}
                            //                     passed
                            //                 </Text>
                            //             </Flex>
                            //         </>
                            //     ),
                            //     sortingField: 'title',
                            //     // minWidth: 400,
                            //     maxWidth: 200,
                            // },

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
                            { id: 'providerConnectionID', visible: true },
                            { id: 'count', visible: true },
                            // { id: 'resourceCount', visible: true },
                            // { id: 'controlCount', visible: true },

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
                            <CustomPagination
                                currentPageIndex={page + 1}
                                pagesCount={Math.ceil(
                                    // @ts-ignore
                                    accounts?.totalCount / 10
                                )}
                                onChange={({ detail }: any) =>
                                    setPage(detail.currentPageIndex - 1)
                                }
                            />
                        }
                    />
                }
            />
            {/* <CloudAccountDetail
                account={account}
                open={open}
                onClose={() => setOpen(false)}
            /> */}
        </>
    )
}
