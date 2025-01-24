import { useEffect, useMemo, useState } from 'react'
import { useAtomValue, useSetAtom } from 'jotai/index'

import { Card, Flex, Text, Title } from '@tremor/react'
import { CheckCircleIcon, ExclamationCircleIcon, XCircleIcon } from '@heroicons/react/24/outline'
import { isDemoAtom, notificationAtom } from '../../../../../../store'
import {
    Api,
    PlatformEnginePkgComplianceApiConformanceStatus,
    PlatformEnginePkgComplianceApiResourceFinding,
} from '../../../../../../api/api'
import AxiosAPI from '../../../../../../api/ApiConfig'
import { dateTimeDisplay } from '../../../../../../utilities/dateDisplay'
import ResourceFindingDetail from '../../../../Findings/ResourceFindingDetail'
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
    SegmentedControl,
} from '@cloudscape-design/components'
import { AppLayout, SplitPanel } from '@cloudscape-design/components'
let sortKey: any[] = []

interface IImpactedResources {
    controlId: string
    // conformanceFilter?: PlatformEnginePkgComplianceApiConformanceStatus[]
    linkPrefix?: string
    isCostOptimization?: boolean
}


export const statusBadge = (
    status: PlatformEnginePkgComplianceApiConformanceStatus | undefined
) => {
    if (
        status ===
        PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed
    ) {
        return (
            <Flex className="w-fit gap-1.5">
                <CheckCircleIcon className="h-4 text-emerald-500" />
                <Text>Passed</Text>
            </Flex>
        )
    }
    return (
        <Flex className="w-fit gap-1.5">
            <XCircleIcon className="h-4 text-rose-600" />
            <Text>Failed</Text>
        </Flex>
    )
}

export default function ImpactedResources({
    controlId,
    // conformanceFilter,
    linkPrefix,
    isCostOptimization,
}: IImpactedResources) {
    const isDemo = useAtomValue(isDemoAtom)
    const setNotification = useSetAtom(notificationAtom)

    const [open, setOpen] = useState(false)
    const [finding, setFinding] = useState<
        PlatformEnginePkgComplianceApiResourceFinding | undefined
    >(undefined)
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)
 const [rows, setRows] =
     useState<PlatformEnginePkgComplianceApiResourceFinding[]>()
 const [page, setPage] = useState(1)
 const [totalCount, setTotalCount] = useState(0)
 const [totalPage, setTotalPage] = useState(0)
 const [conformanceFilter, setConformanceFilter] = useState<
     PlatformEnginePkgComplianceApiConformanceStatus[] | undefined
 >(undefined)
 const conformanceFilterIdx = () => {
     if (
         conformanceFilter?.length === 1 &&
         conformanceFilter[0] ===
             PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed
     ) {
         return '1'
     }
     if (
         conformanceFilter?.length === 1 &&
         conformanceFilter[0] ===
             PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed
     ) {
         return '2'
     }
     return '0'
 }
 
  const GetRows = () => {
      setLoading(true)
      const api = new Api()
      api.instance = AxiosAPI
      
      api.compliance
          .apiV1ResourceFindingsCreate({
              filters: {
                  controlID: [controlId || ''],
                  complianceStatus:
                      conformanceFilter === undefined
                          ? [
                                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed,
                                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
                            ]
                          : conformanceFilter,
                  // @ts-ignore
                  integrationGroup: ['active'],
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

              setTotalPage(Math.ceil(resp.data.totalCount / 15))
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
    }, [page,conformanceFilter])

    // const serverSideRows = useMemo(() => ssr(), [conformanceFilter])

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
                                        {finding?.integrationName}
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
                                            {finding?.integrationType}
                                            {/* <Flex
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
                                            </Card> */}
                                        </Flex>
                                    </>
                                ),
                            },
                            {
                                id: 'failedCount',
                                header: 'Conformance status',
                                maxWidth: 100,
                                cell: (item) => (
                                    <>
                                        {' '}
                                        <Flex className="h-full">
                                            {statusBadge(
                                                item?.findings

                                                    ?.filter(
                                                        (f) =>
                                                            f.controlID ===
                                                            controlId
                                                    )
                                                    .sort((a, b) => {
                                                        if (
                                                            (a.evaluatedAt ||
                                                                0) ===
                                                            (b.evaluatedAt || 0)
                                                        ) {
                                                            return 0
                                                        }
                                                        return (a.evaluatedAt ||
                                                            0) <
                                                            (b.evaluatedAt || 0)
                                                            ? 1
                                                            : -1
                                                    })
                                                    .map(
                                                        (f) =>
                                                            f.complianceStatus
                                                    )
                                                    .at(0)
                                            )}
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
                                            <Text className="text-gray-800">{`${item.totalCount} issues`}</Text>
                                            <Text>{`${
                                                // @ts-ignore
                                                item.totalCount -
                                                // @ts-ignore
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
                            // { id: 'providerConnectionName', visible: true },
                            // { id: 'totalCount', visible: true },
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
                            <SegmentedControl
                                selectedId={conformanceFilterIdx()}
                                onChange={({ detail }) => {
                                    switch (detail.selectedId) {
                                        case '1':
                                            setConformanceFilter([
                                                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
                                            ])
                                            break
                                        case '2':
                                            setConformanceFilter([
                                                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed,
                                            ])
                                            break
                                        default:
                                            setConformanceFilter(undefined)
                                    }
                                }}
                                label="Default segmented control"
                                options={[
                                    { text: 'All', id: '0' },
                                    { text: 'Failed', id: '1' },
                                    { text: 'Passed', id: '2' },
                                ]}
                            />
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
            {/* <ResourceFindingDetail
                resourceFinding={finding}
                controlID={controlId}
                showOnlyOneControl
                open={open}
                onClose={() => setOpen(false)}
                onRefresh={() => window.location.reload()}
                linkPrefix={linkPrefix}
            /> */}
        </>
    )
}
