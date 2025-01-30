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

import {
    PlatformEnginePkgInventoryApiRunQueryResponse,
    Api,
    PlatformEnginePkgInventoryApiSmartQueryItemV2,
    PlatformEnginePkgInventoryApiListQueryRequestV2,
} from '../../../api/api'
import { isDemoAtom, queryAtom, runQueryAtom } from '../../../store'
import AxiosAPI from '../../../api/ApiConfig'

import TopHeader from '../../../components/Layout/Header'
import QueryDetail from './QueryDetail'
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
import CustomPagination from '../../../components/Pagination'
import UseCaseCard from '../../../components/Cards/BookmarkCard'
import axios from 'axios'


export interface Props {
    setTab: Function
     setOpenLayout : Function
}

export default function AllQueries({ setTab, setOpenLayout }: Props) {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loading, setLoading] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [query, setQuery] =
        useState<PlatformEnginePkgInventoryApiListQueryRequestV2>()

    const [engine, setEngine] = useState('odysseus-sql')
    const [integrations, setIntegrations] = useState<any[]>([])
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

    const recordToArray = (record?: Record<string, string[]> | undefined) => {
        if (record === undefined) {
            return []
        }

        return Object.keys(record).map((key) => {
            return {
                value: key,
                resource_types: record[key],
            }
        })
    }
    const getIntegrations = () => {
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
                `${url}/main/integration/api/v1/integration-types/plugin`,
                config
            )
            .then((res) => {
                if (res.data) {
                    const arr = res.data?.items

                    setIntegrations(arr)
                }
            })
            .catch((err) => {
                setLoading(false)
            })
    }

    const getRows = () => {
        setLoading(true)
        const api = new Api()
        api.instance = AxiosAPI

        let body = {
            //  title_filter: '',
            tags: query?.tags,
            integration_types: query?.providers,
            list_of_tables: query?.list_of_tables,
            cursor: page,
            per_page: 12,
        }
        // if (!body.integration_types) {
        //     delete body['integration_types']
        // } else {
        //     // @ts-ignore
        //     body['integration_types'] = ConvertParams(
        //         // @ts-ignore
        //         [body?.integration_types],
        //         'integration_types'
        //     )
        // }
        api.inventory
            .apiV2QueryList(body)
            .then((resp) => {
                if (resp.data.items) {
                    setRows(resp.data.items)
                } else {
                    setRows([])
                }
                setTotalCount(resp.data.total_count)
                setTotalPage(Math.ceil(resp.data.total_count / 12))
                setLoading(false)
            })
            .catch((err) => {
                setLoading(false)
            })
    }

    useEffect(() => {
        getRows()
    }, [page, query])
    useEffect(() => {
        getIntegrations()
    }, [])

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
            Types?.items?.map((item) => {
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
    const FindLogos = (types: string[]) => {
        const temp: string[] = []
        types.map((type) => {
            const integration = integrations.find((i) => i.plugin_id === type)
            if (integration) {
                temp.push(
                    `https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/${integration?.icon}`
                )
            }
        })
        return temp
    }
    return (
        <>
            <Flex className="w-full flex-col justify-start items-start gap-4">
                <Flex className="sm:flex-row flex-col gap-4 w-full sm:justify-between">
                    <Header className="w-full">
                        Queries{' '}
                        <span className=" font-medium">({totalCount})</span>
                    </Header>
                    <CustomPagination
                        currentPageIndex={page}
                        pagesCount={totalPage}
                        onChange={({ detail }: any) =>
                            setPage(detail.currentPageIndex)
                        }
                    />
                </Flex>
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
                <Flex
                    className="gap-8 flex-wrap justify-start items-start w-full"
                    // style={{flex: "1 1 0"}}
                >
                    {rows?.length === 0 || loading ? (
                        <>
                            <Spinner className="mt-2" />
                        </>
                    ) : (
                        <>
                            {rows
                                ?.sort((a, b) => {
                                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                    // @ts-ignore
                                    if (a.title < b.title) {
                                        return -1
                                    }
                                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                    // @ts-ignore
                                    if (a.title > b.title) {
                                        return 1
                                    }
                                    return 0
                                })
                                .map((q, i) => (
                                    <div
                                        className="h-full w-full"
                                        style={
                                            window.innerWidth > 768
                                                ? {
                                                      width: `calc(calc(100% - ${
                                                          rows.length >= 4
                                                              ? '6'
                                                              : ((rows.length - 1)*2)
                                                      }rem) / ${
                                                          rows.length >= 4
                                                              ? '4'
                                                              : rows.length
                                                      })`,
                                                  }
                                                : {}
                                        }
                                    >
                                        <UseCaseCard
                                            // @ts-ignore
                                            title={q?.title}
                                            description={q?.description}
                                            logos={FindLogos(
                                                q?.integration_types
                                            )}
                                            onClick={() => {
                                                // @ts-ignore
                                                setSavedQuery(
                                                    q?.query?.query_to_execute
                                                )
                                                setTab('3')
                                                setOpenLayout(false)
                                            }}
                                            tag="tag1"
                                        />
                                    </div>
                                ))}
                        </>
                    )}
                </Flex>
            </Flex>
        </>
    )
}


{/*
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
                                       
                                        <Title className="text-lg font-semibold ml-2 my-1">
                                            {selectedRow?.title}
                                        </Title>
                                    </Flex>
                                </>
                            ) : (
                                'Query not selected'
                            )
                        }
                    >
                        <>
                            {selectedRow ? (
                                <QueryDetail
                                    // type="resource"
                                    query={selectedRow}
                                    open={openSlider}
                                    onClose={() => setOpenSlider(false)}
                                    onRefresh={() => window.location.reload()}
                                    setTab={setTab}
                                />
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
                                Queries{' '}
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
    */ }