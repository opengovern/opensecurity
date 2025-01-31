// @ts-nocheck
import {
    Button,
    Card,
    Col,
    Flex,
    Grid,
    TableBody,
    TableCell,
    TableHead,
    TableHeaderCell,
    TableRow,
    Text,
    Title,
} from '@tremor/react'
import { ChevronRightIcon } from '@heroicons/react/24/solid'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
    BookOpenIcon,
    CheckCircleIcon,
    CodeBracketIcon,
    Cog8ToothIcon,
    CommandLineIcon,
    XCircleIcon,
    ChevronDownIcon,
    ChevronUpIcon,
} from '@heroicons/react/24/outline'
import { useEffect, useRef, useState } from 'react'
import MarkdownPreview from '@uiw/react-markdown-preview'
import Pagination from '@cloudscape-design/components/pagination'
import DateRangePicker from '@cloudscape-design/components/date-range-picker'

import { useAtomValue } from 'jotai'
import { useComplianceApiV1BenchmarksControlsDetail } from '../../../../../api/compliance.gen'
import Spinner from '../../../../../components/Spinner'
import { numberDisplay } from '../../../../../utilities/numericDisplay'
import DrawerPanel from '../../../../../components/DrawerPanel'
import AnimatedAccordion from '../../../../../components/AnimatedAccordion'
import { searchAtom } from '../../../../../utilities/urlstate'
import {
    PlatformEnginePkgComplianceApiBenchmarkControlSummary,
    PlatformEnginePkgComplianceApiConformanceStatus,
    PlatformEnginePkgControlApiListV2,
    PlatformEnginePkgControlApiListV2ResponseItem,
} from '../../../../../api/api'
import SideNavigation from '@cloudscape-design/components/side-navigation'
import { Api } from '../../../../../api/api'
import AxiosAPI from '../../../../../api/ApiConfig'
import Table from '@cloudscape-design/components/table'
import Box from '@cloudscape-design/components/box'
import SpaceBetween from '@cloudscape-design/components/space-between'
import TextFilter from '@cloudscape-design/components/text-filter'
import Header from '@cloudscape-design/components/header'
import Badge from '@cloudscape-design/components/badge'
import KButton from '@cloudscape-design/components/button'
import axios from 'axios'
import ReactEChartsCore from 'echarts-for-react/lib/core'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
    AppLayout,
    BreadcrumbGroup,
    ContentLayout,
    Link,
    PropertyFilter,
    SplitPanel,
} from '@cloudscape-design/components'
import ControlDetail from './ControlDetail'
import { severityBadge } from '../../../Controls'
import CustomPagination from '../../../../../components/Pagination'

interface IPolicies {
    id: string | undefined
    assignments: number
    enable?: boolean
}

export const activeBadge = (status: boolean) => {
    if (status) {
        return (
            <Flex className="w-fit gap-1.5">
                <CheckCircleIcon className="h-4 text-emerald-500" />
                <Text>Active</Text>
            </Flex>
        )
    }
    return (
        <Flex className="w-fit gap-1.5">
            <XCircleIcon className="h-4 text-rose-600" />
            <Text>Inactive</Text>
        </Flex>
    )
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

export const treeRows = (
    json: PlatformEnginePkgComplianceApiBenchmarkControlSummary | undefined
) => {
    let arr: any = []
    if (json) {
        if (json.control !== null && json.control !== undefined) {
            for (let i = 0; i < json.control.length; i += 1) {
                let obj = {}
                obj = {
                    parentName: json?.benchmark?.title,
                    ...json.control[i].control,
                    ...json.control[i],
                }
                arr.push(obj)
            }
        }
        if (json.children !== null && json.children !== undefined) {
            for (let i = 0; i < json.children.length; i += 1) {
                const res = treeRows(json.children[i])
                arr = arr.concat(res)
            }
        }
    }

    return arr
}

export const groupBy = (input: any[], key: string) => {
    return input.reduce((acc, currentValue) => {
        const groupKey = currentValue[key]
        if (!acc[groupKey]) {
            acc[groupKey] = []
        }
        acc[groupKey].push(currentValue)
        return acc
    }, {})
}

export const countControls = (
    v: PlatformEnginePkgComplianceApiBenchmarkControlSummary | undefined
) => {
    const countChildren = v?.children
        ?.map((i) => countControls(i))
        .reduce((prev, curr) => prev + curr, 0)
    const total: number = (countChildren || 0) + (v?.control?.length || 0)
    return total
}

export default function Controls({
    id,
    assignments,
    enable,
    accounts,
}: IPolicies) {
    const { response: controls, isLoading } =
        useComplianceApiV1BenchmarksControlsDetail(String(id))
    const [page, setPage] = useState<number>(1)
    const [rows, setRows] = useState<
        PlatformEnginePkgControlApiListV2ResponseItem[]
    >([])
    const [selectedRow, setSelectedRow] =
        useState<PlatformEnginePkgControlDetailV3>()
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const [benchmarkId, setBenchmarkId] = useState(id)
    const [loading, setLoading] = useState(true)

    const [doc, setDoc] = useState('')
    const [docTitle, setDocTitle] = useState('')
    const [open, setOpen] = useState(false)
    const [openAllControls, setOpenAllControls] = useState(false)
    const [listofTables, setListOfTables] = useState([])
    const [totalPage, setTotalPage] = useState<number>(0)
    const [totalCount, setTotalCount] = useState<number>(0)
    const [query, setQuery] = useState<PlatformEnginePkgControlApiListV2>()
    const [queries, setQueries] = useState({
        tokens: [],
        operation: 'and',
    })

    const [tree, setTree] = useState()
    const [selected, setSelected] = useState()
    const [selectedBread, setSelectedBread] = useState()
    const [treePage, setTreePage] = useState(0)
    const [treeTotal, setTreeTotal] = useState(0)
    const [treeTotalPages, setTreeTotalPages] = useState(0)

    const [filters, setFilters] = useState([
        {
            propertyKey: 'severity',
            value: 'High',
        },
        {
            propertyKey: 'severity',
            value: 'Critical',
        },
        {
            propertyKey: 'severity',
            value: 'Low',
        },
        {
            propertyKey: 'severity',
            value: 'Medium',
        },
    ])
    const [filterOption, setFilterOptions] = useState([
        {
            key: 'severity',
            operators: ['='],
            propertyLabel: `Severity`,
            groupValuesLabel: `Severity values`,
        },
        // {
        //     key: 'severity',
        //     operators: ['='],
        //     propertyLabel: `Exclude Inactive Integration'`,
        //     groupValuesLabel: `Exclude Inactive Integration' values`,
        // },
    ])
    const [sort, setSort] = useState('noncompliant_resources')
    const [sortOrder, setSortOrder] = useState(true)

    const navigateToInsightsDetails = (id: string) => {
        navigate(`${id}?${searchParams}`)
    }
    const toggleOpen = () => {
        setOpenAllControls(!openAllControls)
    }

    const countBenchmarks = (
        v: PlatformEnginePkgComplianceApiBenchmarkControlSummary | undefined
    ) => {
        const countChildren = v?.children
            ?.map((i) => countBenchmarks(i))
            .reduce((prev, curr) => prev + curr, 0)
        const total: number = (countChildren || 0) + (v?.children?.length || 0)
        return total
    }
    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 30 ? text.substring(0, 30) + '...' : text
        }
    }
    const GetControls = (flag: boolean, id: string | undefined) => {
        setLoading(true)
        const api = new Api()
        api.instance = AxiosAPI
        //   const benchmarks = category
        //   const temp = []
        //   temp.push(`aws_score_${benchmarks}`)
        //   temp.push(`azure_score_${benchmarks}`)

        let body = {
            // list_of_tables: listofTables,
            severity: query?.severity,
            root_benchmark: flag ? [id] : [benchmarkId],
            compliance_result_summary: true,
            cursor: page,
            per_page: 10,
            sort_by: sort,
            sort_order: sortOrder ? 'desc' : 'asc',
        }
        if (listofTables.length == 0) {
            // @ts-ignore
            delete body['list_of_tables']
        }
        api.compliance
            // @ts-ignore
            .apiV2ControlList(body)
            .then((resp) => {
                setTotalPage(Math.ceil(resp.data.total_count / 10))
                setTotalCount(resp.data.total_count)
                if (resp.data.items) {
                    setRows(resp.data.items)
                }

                setLoading(false)
            })
            .catch((err) => {
                setLoading(false)
                setRows([])
            })
    }
    const GetTree = () => {
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        axios
            .get(
                `${url}/main/compliance/api/v3/benchmarks/${benchmarkId}/nested`,
                config
            )
            .then((res) => {
                const temp = []

                if (res.data.children) {
                    res.data.children.map((item) => {
                        let childs = {
                            text: truncate(item.title),
                            href: item.id,
                            type: item.children
                                ? 'expandable-link-group'
                                : 'link',
                        }
                        if (item.children && item.children.length > 0) {
                            let child_item = []
                            item.children.map((sub_item) => {
                                child_item.push({
                                    text: truncate(sub_item.title),
                                    href: sub_item.id,
                                    type: 'link',
                                    parentId: item.id,
                                    parentTitle: item.title,
                                })
                            })
                            childs['items'] = child_item
                        }
                        temp.push(childs)
                    })

                    // setTree(temp)
                } else {
                    // temp.push({
                    //     text: res.data.title,
                    //     href: res.data.id,
                    //     type: 'link',
                    // })
                }
                setTreeTotalPages(Math.ceil(temp.length / 18))
                setTreeTotal(temp.length)
                setTreePage(0)

                setTree(temp)
            })
            .catch((err) => {
                console.log(err)
            })
    }
    useEffect(() => {
        if (window.innerWidth > 640) {
            GetTree()
        }
    }, [])
    useEffect(() => {
        setPage(0)
    }, [selected])

    useEffect(() => {
        GetControls(selected ? true : false, selected)
    }, [selected, sort, sortOrder, page])
    useEffect(() => {
        let temp = {}

        queries?.tokenGroups?.map((item) => {
            // @ts-ignore
            if (temp[item.propertyKey] && temp[item.propertyKey].length > 0) {
                temp[item.propertyKey].push(item.value.toLowerCase())
            } else {
                temp[item.propertyKey] = []
                temp[item.propertyKey].push(item.value.toLowerCase())
            }
        })
        setQuery(temp)
    }, [queries])

    useEffect(() => {
        if (selected) {
            GetControls(true, selected)
        } else {
            GetControls(false)
        }
    }, [query])

    const getControlDetail = (id: string) => {
        const api = new Api()
        api.instance = AxiosAPI
        // setLoading(true);
        api.compliance
            .apiV3ControlDetail(id)
            .then((resp) => {
                setSelectedRow(resp.data)
                setOpen(true)
                // setLoading(false)
            })
            .catch((err) => {
                // setLoading(false)
            })
    }
    return (
        <AppLayout
            toolsOpen={false}
            navigationOpen={false}
            contentType="dashboard"
            className="w-full bg-transparent rounded-xl"
            toolsHide={true}
            navigationHide={
                !(tree && tree?.length > 0) || window.innerWidth < 640
            }
            disableContentPaddings={true}
            splitPanelSize={500}
            splitPanelOpen={open}
            onSplitPanelToggle={() => {
                setOpen(!open)
                if (open) {
                    setSelectedRow(undefined)
                }
            }}
            splitPanel={
                // @ts-ignore
                <SplitPanel
                    // @ts-ignore
                    header={
                        selectedRow ? (
                            <>
                                <Flex
                                    justifyContent="start"
                                    className="gap-2 items-center justify-center sm:flex-row flex-col"
                                >
                                    <Title className="text-lg font-semibold ml-2 my-1 w-full">
                                        {selectedRow?.title}
                                    </Title>
                                    {severityBadge(selectedRow?.severity)}
                                </Flex>
                            </>
                        ) : (
                            'Control not selected'
                        )
                    }
                >
                    <ControlDetail
                        // type="resource"
                        benchmarkId={id}
                        selectedItem={selectedRow}
                        open={open}
                        onClose={() => setOpen(false)}
                        onRefresh={() => {}}
                    />
                </SplitPanel>
            }
            navigationOpen={tree && tree.length > 0}
            navigationWidth={300}
            navigation={
                <>
                    {tree && tree.length > 0 && (
                        <Flex
                            className="bg-white  w-full justify-center items-center   "
                            flexDirection="col"
                        >
                            <>
                                <SideNavigation
                                    className="w-full scroll   overflow-scroll p-4 pb-0"
                                    activeHref={selected}
                                    virtualScroll
                                    header={{
                                        href: benchmarkId,
                                        text: 'Control Groups',
                                    }}
                                    onFollow={(event) => {
                                        event.preventDefault()
                                        setSelected(event.detail.href)
                                        const temp = []

                                        if (event.detail.parentId) {
                                            temp.push({
                                                text: controls?.benchmark
                                                    ?.title,
                                                href: id,
                                            })
                                            temp.push({
                                                text: event.detail.parentTitle,
                                                href: event.detail.parentId,
                                            })
                                            temp.push({
                                                text: event.detail.text,
                                                href: event.detail.href,
                                            })
                                        } else {
                                            temp.push({
                                                text: controls?.benchmark
                                                    ?.title,
                                                href: id,
                                            })
                                            if (
                                                event.detail.text !==
                                                controls?.benchmark?.title
                                            ) {
                                                temp.push({
                                                    text: event.detail.text,
                                                    href: event.detail.href,
                                                })
                                            }
                                        }
                                        setSelectedBread(temp)
                                    }}
                                    items={tree?.slice(
                                        treePage * 18,
                                        (treePage + 1) * 18
                                    )}
                                />
                            </>
                            {treeTotalPages > 1 && (
                                <>
                                    <div className="flex justify-center items-center ">
                                        <CustomPagination
                                            className="pb-2"
                                            currentPageIndex={treePage + 1}
                                            pagesCount={treeTotalPages}
                                            onChange={({ detail }) =>
                                                setTreePage(
                                                    detail.currentPageIndex - 1
                                                )
                                            }
                                        />
                                    </div>
                                </>
                            )}
                        </Flex>
                    )}
                </>
            }
            content={
                <ContentLayout
                    header={
                        <Header className="w-full">
                            Controls{' '}
                            <span className=" font-medium">({totalCount})</span>
                        </Header>
                    }
                    defaultPadding={true}
                    breadcrumbs={
                        <BreadcrumbGroup
                            onClick={(event) => {
                                event.preventDefault()
                                setSelected(event.detail.href)
                            }}
                            items={selectedBread}
                            ariaLabel="Breadcrumbs"
                        />
                    }
                >
                    <Table
                        className="   "
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
                            setSort(event.detail.sortingColumn.sortingField)
                            setSortOrder(!sortOrder)
                        }}
                        onRowClick={(event) => {
                            const row = event.detail.item
                            setSelectedRow(undefined)
                            getControlDetail(row.id)
                            setOpen(true)
                        }}
                        sortingColumn={sort}
                        sortingDescending={sortOrder}
                        // sortingDescending={sortOrder == 'desc' ? true : false}
                        columnDefinitions={[
                            {
                                id: 'id',
                                header: 'ID',
                                cell: (item) => item.id,
                                sortingField: 'id',
                                isRowHeader: true,
                            },
                            {
                                id: 'title',
                                header: 'Title',
                                cell: (item) => item.title,
                                sortingField: 'title',
                                // minWidth: 400,
                                maxWidth: 200,
                            },
                            {
                                id: 'connector',
                                header: 'Connector',
                                cell: (item) => item.connector,
                            },
                            {
                                id: 'query',
                                header: 'Primary Table',
                                cell: (item) => item?.query?.primary_table,
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
                                id: 'query.parameters',
                                header: 'Has Parametrs',
                                cell: (item) => (
                                    // @ts-ignore
                                    <>
                                        {item.query?.parameters.length > 0
                                            ? 'True'
                                            : 'False'}
                                    </>
                                ),
                            },
                            {
                                id: 'incidents',
                                header: 'Incidents',
                                // sortingField: 'incidents',

                                cell: (item) => (
                                    // @ts-ignore
                                    <>
                                        {/**@ts-ignore */}
                                        {item?.compliance_results_summary
                                            ?.incident_count
                                            ? item?.compliance_results_summary
                                                  ?.incident_count
                                            : 0}
                                    </>
                                ),
                                // minWidth: 50,
                                maxWidth: 100,
                            },
                            {
                                id: 'passing_resources',
                                header: 'Non Incidents ',

                                cell: (item) => (
                                    // @ts-ignore
                                    <>
                                        {item?.compliance_results_summary
                                            ?.non_incident_count
                                            ? item?.compliance_results_summary
                                                  ?.non_incident_count
                                            : 0}
                                    </>
                                ),
                                maxWidth: 100,
                            },
                            {
                                id: 'noncompliant_resources',
                                header: 'Non-Compliant Resources',
                                sortingField: 'noncompliant_resources',

                                cell: (item) => (
                                    // @ts-ignore
                                    <>
                                        {item?.compliance_results_summary
                                            ?.noncompliant_resources
                                            ? item?.compliance_results_summary
                                                  ?.noncompliant_resources
                                            : 0}
                                    </>
                                ),
                                maxWidth: 100,
                            },
                            {
                                id: 'waste',
                                header: 'Waste',
                                cell: (item) => (
                                    // @ts-ignore
                                    <>
                                        {item?.compliance_results_summary
                                            ?.cost_optimization
                                            ? item?.compliance_results_summary
                                                  ?.cost_optimization
                                            : 0}
                                    </>
                                ),
                                maxWidth: 100,
                            },
                            {
                                id: 'action',
                                header: 'Action',
                                cell: (item) => (
                                    // @ts-ignore
                                    <KButton
                                        onClick={() => {
                                            navigateToInsightsDetails(item.id)
                                        }}
                                        variant="inline-link"
                                        ariaLabel={`Open Detail`}
                                    >
                                        Open
                                    </KButton>
                                ),
                            },
                        ]}
                        columnDisplay={[
                            { id: 'id', visible: false },
                            { id: 'title', visible: true },
                            { id: 'connector', visible: false },
                            { id: 'query', visible: false },
                            { id: 'severity', visible: true },
                            { id: 'incidents', visible: false },
                            { id: 'passing_resources', visible: false },
                            {
                                id: 'noncompliant_resources',
                                visible: true,
                            },
                            {
                                id: 'waste',
                                visible:
                                    (benchmarkId == 'sre_efficiency' ||
                                        selected == 'sre_efficiency') &&
                                    enable
                                        ? true
                                        : false,
                            },

                            // { id: 'action', visible: true },
                        ]}
                        enableKeyboardNavigation
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
                            <PropertyFilter
                                // @ts-ignore
                                query={queries}
                                // @ts-ignore
                                onChange={({ detail }) => {
                                    // @ts-ignore
                                    setQueries(detail)
                                }}
                                // countText="5 matches"
                                enableTokenGroups
                                expandToViewport
                                filteringAriaLabel="Control Categories"
                                // @ts-ignore
                                // filteringOptions={filters}
                                filteringPlaceholder="Control Categories"
                                // @ts-ignore
                                filteringOptions={filters}
                                filteringProperties={filterOption}
                                // filteringProperties={
                                //     filterOption
                                // }
                            />
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
                </ContentLayout>
            }
        />
    )
}
