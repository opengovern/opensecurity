import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Button,
    Card,
    Flex,
    Icon,
    Select,
    SelectItem,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    TextInput,
    Subtitle,
    Title,
} from '@tremor/react'
import {
    ChevronDoubleLeftIcon,
    ChevronDownIcon,
    ChevronUpIcon,
    CloudIcon,
    CommandLineIcon,
    FunnelIcon,
    MagnifyingGlassIcon,
    PlayCircleIcon,
    PlusIcon,
    TagIcon,
} from '@heroicons/react/24/outline'
import { Fragment, useEffect, useMemo, useState } from 'react' // eslint-disable-next-line import/no-extraneous-dependencies
import { highlight, languages } from 'prismjs' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/components/prism-sql' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/themes/prism.css'
import Editor from 'react-simple-code-editor'

import {
    CheckCircleIcon,
    ExclamationCircleIcon,
} from '@heroicons/react/24/solid'
import { Transition } from '@headlessui/react'
import { useAtom, useAtomValue } from 'jotai'
import {
    useInventoryApiV1QueryList,
    useInventoryApiV1QueryRunCreate,
    useInventoryApiV2AnalyticsCategoriesList,
    useInventoryApiV2QueryList,
    useInventoryApiV3AllQueryCategory,
    useInventoryApiV3QueryFiltersList,
} from '../../../api/inventory.gen'
import Spinner from '../../../components/Spinner'
import { getErrorMessage } from '../../../types/apierror'
import { RenderObject } from '../../../components/RenderObject'

import {
    PlatformEnginePkgInventoryApiRunQueryResponse,
    Api,
    PlatformEnginePkgInventoryApiSmartQueryItemV2,
    PlatformEnginePkgInventoryApiListQueryRequestV2,
} from '../../../api/api'
import { isDemoAtom, queryAtom, runQueryAtom } from '../../../store'
import AxiosAPI from '../../../api/ApiConfig'

import { snakeCaseToLabel } from '../../../utilities/labelMaker'
import { numberDisplay } from '../../../utilities/numericDisplay'
import TopHeader from '../../../components/Layout/Header'
import { array } from 'prop-types'
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
import { useIntegrationApiV1EnabledConnectorsList } from '../../../api/integration.gen'
import axios from 'axios'
import ViewDetail from './detail'
import CustomPagination from '../../../components/Pagination'




export interface Props {
    setTab: Function
    setOpenLayout: Function
}

export default function View({ setTab, setOpenLayout }: Props) {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loading, setLoading] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [code, setCode] = useState(savedQuery || '')
    const [selectedIndex, setSelectedIndex] = useState(0)
    const [searchCategory, setSearchCategory] = useState('')
    const [selectedRow, setSelectedRow] =
        useState<PlatformEnginePkgInventoryApiSmartQueryItemV2>()
    const [openDrawer, setOpenDrawer] = useState(false)
    const [openSlider, setOpenSlider] = useState(false)
    const [openSearch, setOpenSearch] = useState(true)
    const [query, setQuery] =
        useState<PlatformEnginePkgInventoryApiListQueryRequestV2>()
    const [selectedFilter, setSelectedFilters] = useState<string[]>([])

    const [showEditor, setShowEditor] = useState(true)
    const isDemo = useAtomValue(isDemoAtom)
    const [pageSize, setPageSize] = useState(1000)
    const [autoRun, setAutoRun] = useState(false)
    const [listofTables, setListOfTables] = useState([])

    const [engine, setEngine] = useState('odysseus-sql')
    const [page, setPage] = useState(1)
    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [rows, setRows] = useState<any[]>()
    const [filterQuery, setFilterQuery] = useState({
        tokens: [],
        operation: 'and',
    })
    const [properties, setProperties] = useState<any[]>([])
    const [options, setOptions] = useState<any[]>([])

    const {
        response: categories,
        isLoading: categoryLoading,
        isExecuted: categoryExec,
    } = useInventoryApiV3AllQueryCategory()

    const {
        response: filters,
        isLoading: filtersLoading,
        isExecuted: filterExec,
    } = useInventoryApiV3QueryFiltersList()

    const {
        response: Types,
        isLoading: TypesLoading,
        isExecuted: TypesExec,
    } = useIntegrationApiV1EnabledConnectorsList(0, 0)

    // const { response: queries, isLoading: queryLoading } =
    //     useInventoryApiV2QueryList({
    //         titleFilter: '',
    //         Cursor: 0,
    //         PerPage:25
    //     })

    const getRows = () => {
        setLoading(true)
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
                `${url}/main/core/api/v3/views?per_page=10&cursor=${page}`,
                config
            )
            .then((res) => {
                if (res.data) {
                    setRows(res.data.views)
                    setTotalCount(res.data.total_count)
                    setTotalPage(Math.ceil(res.data.total_count / 10))
                }
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }

    useEffect(() => {
        getRows()
    }, [page, query])

    useEffect(() => {
        if (
            filterExec &&
            categoryExec &&
            TypesExec &&
            !TypesLoading &&
            !filtersLoading &&
            !categoryLoading
        ) {
            const temp_option: any = []
            Types?.integration_types?.map((item) => {
                temp_option.push({
                    propertyKey: 'integrationType',
                    value: item.platform_name,
                })
            })

            const property: any = [
                {
                    key: 'integrationType',
                    operators: ['='],
                    propertyLabel: 'integration Type',
                    groupValuesLabel: 'integrationType values',
                },
            ]
            categories?.categories?.map((item) => {
                property.push({
                    key: `list_of_table${item.category}`,
                    operators: ['='],
                    propertyLabel: item.category,
                    groupValuesLabel: `${item.category} values`,
                    group: 'category',
                })
                item?.tables?.map((sub) => {
                    temp_option.push({
                        propertyKey: `list_of_table${item.category}`,
                        value: sub.table,
                    })
                })
            })
            filters?.tags?.map((unique, index) => {
                property.push({
                    key: unique.Key,
                    operators: ['='],
                    propertyLabel: unique.Key,
                    groupValuesLabel: `${unique.Key} values`,
                    // @ts-ignore
                    group: 'tags',
                })
                unique.UniqueValues?.map((value, idx) => {
                    temp_option.push({
                        propertyKey: unique.Key,
                        value: value,
                    })
                })
            })
            setOptions(temp_option)
            setProperties(property)
        }
    }, [
        filterExec,
        categoryExec,
        filtersLoading,
        categoryLoading,
        TypesExec,
        TypesLoading,
    ])

    useEffect(() => {
        if (filterQuery) {
            const temp_provider: any = []
            const temp_tables: any = []
            const temp_tags = {}
            filterQuery.tokens.map((item, index) => {
                // @ts-ignore
                if (item.propertyKey === 'integrationType') {
                    // @ts-ignore

                    temp_provider.push(item.value)
                }
                // @ts-ignore
                else if (item.propertyKey.includes('list_of_table')) {
                    // @ts-ignore

                    temp_tables.push(item.value)
                } else {
                    // @ts-ignore

                    if (temp_tags[item.propertyKey]) {
                        // @ts-ignore

                        temp_tags[item.propertyKey].push(item.value)
                    } else {
                        // @ts-ignore

                        temp_tags[item.propertyKey] = [item.value]
                    }
                }
            })
            // @ts-ignore
            setQuery({
                providers: temp_provider.length > 0 ? temp_provider : undefined,
                list_of_tables:
                    temp_tables.length > 0 ? temp_tables : undefined,
                // @ts-ignore
                tags: temp_tags,
            })
        }
    }, [filterQuery])

    return (
        <>
            
            <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                className="w-full"
                toolsHide={true}
                navigationHide={true}
                splitPanelOpen={openSlider}
                onSplitPanelToggle={() => {
                    setOpenSlider(!openSlider)
                    if (openSlider) {
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
                                    <Flex justifyContent="start">
                                        {/* {getConnectorIcon(
                                            selectedRow?.connector
                                        )} */}
                                        <Title className="text-lg font-semibold ml-2 my-1">
                                            {selectedRow?.title}
                                        </Title>
                                    </Flex>
                                </>
                            ) : (
                                'View not selected'
                            )
                        }
                    >
                        <>
                            {selectedRow ? (
                                <>
                                    <ViewDetail
                                        // type="resource"
                                        query={selectedRow}
                                        open={openSlider}
                                        onClose={() => setOpenSlider(false)}
                                        onRefresh={() =>
                                            window.location.reload()
                                        }
                                        setTab={setTab}
                                        setOpenLayout={setOpenLayout}
                                    />
                                </>
                            ) : (
                                <Spinner />
                            )}
                        </>
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

                            setSelectedRow(row)
                            setOpenSlider(true)
                        }}
                        columnDefinitions={[
                            {
                                id: 'id',
                                header: 'Id',
                                cell: (item) => item.id,
                                // sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 150,
                            },
                            {
                                id: 'title',
                                header: 'Title',
                                cell: (item) => item.title,
                                // sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 150,
                            },
                            {
                                id: 'description',
                                header: 'Description',
                                cell: (item) => item.description,
                                // sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 150,
                            },
                        ]}
                        columnDisplay={[
                            {
                                id: 'id',
                                visible: true,
                            },
                            {
                                id: 'title',
                                visible: true,
                            },

                            { id: 'description', visible: true },
                            // {
                            //     id: 'severity',
                            //     visible: true,
                            // },
                            // { id: 'parameters', visible: true },
                            // {
                            //     id: 'evaluatedAt',
                            //     visible: true,
                            // },

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
                            <PropertyFilter
                                // @ts-ignore
                                query={filterQuery}
                                tokenLimit={2}
                                onChange={({ detail }) =>
                                    // @ts-ignore
                                    setFilterQuery(detail)
                                }
                                customGroupsText={[
                                    {
                                        properties: 'Tags',
                                        values: 'Tag values',
                                        group: 'tags',
                                    },
                                    {
                                        properties: 'Category',
                                        values: 'Category values',
                                        group: 'category',
                                    },
                                ]}
                                // countText="5 matches"
                                expandToViewport
                                filteringAriaLabel="Find Query"
                                filteringPlaceholder="Find Query"
                                filteringOptions={options}
                                filteringProperties={properties}
                                asyncProperties
                                virtualScroll
                            />
                        }
                        header={
                            <Header className="w-full">
                                Views{' '}
                                <span className=" font-medium">
                                    ({totalCount})
                                </span>
                            </Header>
                        }
                        pagination={
                            <CustomPagination
                                currentPageIndex={page}
                                pagesCount={totalPage}
                                onChange={({ detail }: any) =>
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
