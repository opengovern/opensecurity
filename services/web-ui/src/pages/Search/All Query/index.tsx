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
import DrawerPanel from '../../../components/DrawerPanel'
import { RenderObject } from '../../../components/RenderObject'

import {
    GithubComKaytuIoKaytuEnginePkgInventoryApiRunQueryResponse,
    Api,
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2,
    GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2,
} from '../../../api/api'
import { isDemoAtom, queryAtom, runQueryAtom } from '../../../store'
import AxiosAPI from '../../../api/ApiConfig'

import { snakeCaseToLabel } from '../../../utilities/labelMaker'
import { numberDisplay } from '../../../utilities/numericDisplay'
import TopHeader from '../../../components/Layout/Header'
import QueryDetail from './QueryDetail'
import { array } from 'prop-types'
import KFilter from '../../../components/Filter'
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


export interface Props {
    setTab: Function
}

export default function AllQueries({ setTab }: Props) {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loading, setLoading] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [code, setCode] = useState(savedQuery || '')
    const [selectedIndex, setSelectedIndex] = useState(0)
    const [searchCategory, setSearchCategory] = useState('')
    const [selectedRow, setSelectedRow] =
        useState<GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2>()
    const [openDrawer, setOpenDrawer] = useState(false)
    const [openSlider, setOpenSlider] = useState(false)
    const [openSearch, setOpenSearch] = useState(true)
    const [query, setQuery] =
        useState<GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2>()
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

    const ConvertParams = (array: string[], key: string) => {
        return `[${array[0]}]`
        // let temp = ''
        // array.map((item,index)=>{
        //     if(index ===0){
        //         temp = temp + item
        //     }
        //     else{
        //         temp = temp +'&'+key+'='+item
        //     }
        // })
        // return temp
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
            per_page: 15,
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
                setTotalPage(Math.ceil(resp.data.total_count / 15))
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
    }, [filterExec, categoryExec, filtersLoading, categoryLoading, TypesExec,TypesLoading])

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
            <TopHeader />
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
            {/* {categoryLoading || loading ? (
                <Spinner className="mt-56" />
            ) : (
                <Flex alignItems="start" className="gap-4">
                    <DrawerPanel
                        open={false}
                        onClose={() => setOpenDrawer(false)}
                    >
                        <RenderObject obj={selectedRow} />
                    </DrawerPanel>
                    {openSearch ? (
                        <Card className="sticky w-fit">
                            <TextInput
                                className="w-56 mb-6"
                                icon={MagnifyingGlassIcon}
                                placeholder="Search..."
                                value={searchCategory}
                                onChange={(e) =>
                                    setSearchCategory(e.target.value)
                                }
                            />
                            {categories?.categories.map(
                                (cat) =>
                                    !!cat.tables?.filter((catt) =>
                                        catt.name
                                            .toLowerCase()
                                            .includes(
                                                searchCategory.toLowerCase()
                                            )
                                    ).length && (
                                        <Accordion className="w-56 border-0 rounded-none bg-transparent mb-1">
                                            <AccordionHeader className="pl-0 pr-0.5 py-1 w-full bg-transparent">
                                                <Text className="text-gray-800 text-left">
                                                    {cat.category}
                                                </Text>
                                            </AccordionHeader>
                                            <AccordionBody className="p-0 w-full pr-0.5 cursor-default bg-transparent">
                                                <Flex
                                                    flexDirection="col"
                                                    justifyContent="start"
                                                >
                                                    {cat.tables
                                                        ?.filter((catt) =>
                                                            catt.name
                                                                .toLowerCase()
                                                                .includes(
                                                                    searchCategory.toLowerCase()
                                                                )
                                                        )
                                                        .map((subCat) => (
                                                            <Flex
                                                                justifyContent="start"
                                                                onClick={() => {
                                                                    if (
                                                                        // @ts-ignore
                                                                        listofTables.includes(
                                                                            // @ts-ignore

                                                                            subCat.table
                                                                        )
                                                                    ) {
                                                                        // @ts-ignore
                                                                        setListOfTables(
                                                                            // @ts-ignore
                                                                            listofTables.filter(
                                                                                (
                                                                                    item
                                                                                ) =>
                                                                                    item !==
                                                                                    subCat.table
                                                                            )
                                                                        )
                                                                    }
                                                                    // @ts-ignore
                                                                    else {
                                                                        // @ts-ignore
                                                                        setListOfTables(
                                                                            [
                                                                                // @ts-ignore
                                                                                ...listofTables,
                                                                                // @ts-ignore
                                                                                subCat.table,
                                                                            ]
                                                                        )
                                                                    }
                                                                }}
                                                            >
                                                                <Text className="ml-4 w-full truncate text-start py-2 cursor-pointer hover:text-openg-600">
                                                                    {
                                                                        subCat.name
                                                                    }
                                                                </Text>
                                                            </Flex>
                                                        ))}
                                                </Flex>
                                            </AccordionBody>
                                        </Accordion>
                                    )
                            )}
                            {listofTables.length > 0 && (
                                <>
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="start"
                                        alignItems="start"
                                    >
                                        <Text>Selected Filters</Text>
                                        {listofTables.map((item, index) => {
                                            return (
                                                <Flex
                                                    justifyContent="start"
                                                    className="w-full"
                                                >
                                                    <Text>{item}</Text>
                                                </Flex>
                                            )
                                        })}
                                    </Flex>
                                </>
                            )}
                            <Flex justifyContent="end" className="mt-12">
                                <Button
                                    variant="light"
                                    onClick={() => setOpenSearch(false)}
                                >
                                    <ChevronDoubleLeftIcon className="h-4" />
                                </Button>
                            </Flex>
                        </Card>
                    ) : (
                        <Flex
                            flexDirection="col"
                            justifyContent="center"
                            className="min-h-full w-fit"
                        >
                            <Button
                                variant="light"
                                onClick={() => setOpenSearch(true)}
                            >
                                <Flex flexDirection="col" className="gap-4 w-4">
                                    <FunnelIcon />
                                    <Text className="rotate-90">Options</Text>
                                </Flex>
                            </Button>
                        </Flex>
                    )}*/}
            {/* <Flex flexDirection="col" className="w-full "> */}
            {/* <Transition.Root show={showEditor} as={Fragment}>
                            <Transition.Child
                                as={Fragment}
                                enter="ease-in-out duration-500"
                                enterFrom="h-0 opacity-0"
                                enterTo="h-fit opacity-100"
                                leave="ease-in-out duration-500"
                                leaveFrom="h-fit opacity-100"
                                leaveTo="h-0 opacity-0"
                            >
                                <Flex flexDirection="col" className="mb-4"> 
                                    <Card className="relative overflow-hidden">
                                        <Editor
                                            onValueChange={(text) => {
                                                setSavedQuery('')
                                                setCode(text)
                                            }}
                                            highlight={(text) =>
                                                highlight(
                                                    text,
                                                    languages.sql,
                                                    'sql'
                                                )
                                            }
                                            value={code}
                                            className="w-full bg-white dark:bg-gray-900 dark:text-gray-50 font-mono text-sm"
                                            style={{
                                                minHeight: '200px',
                                                // maxHeight: '500px',
                                                overflowY: 'scroll',
                                            }}
                                            placeholder="-- write your SQL query here"
                                        />
                                        {isLoading && isExecuted && (
                                            <Spinner className="bg-white/30 backdrop-blur-sm top-0 left-0 absolute flex justify-center items-center w-full h-full" />
                                        )}
                                    </Card>
                                    <Flex className="w-full mt-4">
                                        <Flex justifyContent="start">
                                            <Text className="mr-2 w-fit">
                                                Maximum rows:
                                            </Text>
                                            <Select
                                                enableClear={false}
                                                className="w-56"
                                                placeholder="1,000"
                                            >
                                                <SelectItem
                                                    value="1000"
                                                    onClick={() =>
                                                        setPageSize(1000)
                                                    }
                                                >
                                                    1,000
                                                </SelectItem>
                                                <SelectItem
                                                    value="3000"
                                                    onClick={() =>
                                                        setPageSize(3000)
                                                    }
                                                >
                                                    3,000
                                                </SelectItem>
                                                <SelectItem
                                                    value="5000"
                                                    onClick={() =>
                                                        setPageSize(5000)
                                                    }
                                                >
                                                    5,000
                                                </SelectItem>
                                                <SelectItem
                                                    value="10000"
                                                    onClick={() =>
                                                        setPageSize(10000)
                                                    }
                                                >
                                                    10,000
                                                </SelectItem>
                                            </Select>
                                            <Text className="mr-2 w-fit">
                                                Engine:
                                            </Text>
                                            <Select
                                                enableClear={false}
                                                className="w-56"
                                                value={engine}
                                            >
                                                <SelectItem
                                                    value="odysseus-sql"
                                                    onClick={() =>
                                                        setEngine(
                                                            'odysseus-sql'
                                                        )
                                                    }
                                                >
                                                    Odysseus SQL
                                                </SelectItem>
                                                <SelectItem
                                                    value="odysseus-rego"
                                                    onClick={() =>
                                                        setEngine(
                                                            'odysseus-rego'
                                                        )
                                                    }
                                                >
                                                    Odysseus Rego
                                                </SelectItem>
                                            </Select>
                                        </Flex>
                                        <Flex className="w-fit gap-x-3">
                                            {!!code.length && (
                                                <Button
                                                    variant="light"
                                                    color="gray"
                                                    icon={CommandLineIcon}
                                                    onClick={() => setCode('')}
                                                >
                                                    Clear editor
                                                </Button>
                                            )}
                                            <Button
                                                icon={PlayCircleIcon}
                                                onClick={() => sendNow()}
                                                disabled={!code.length}
                                                loading={
                                                    isLoading && isExecuted
                                                }
                                                loadingText="Running"
                                            >
                                                Run query
                                            </Button>
                                        </Flex>
                                    </Flex>
                                    <Flex className="w-full">
                                        {!isLoading && isExecuted && error && (
                                            <Flex
                                                justifyContent="start"
                                                className="w-fit"
                                            >
                                                <Icon
                                                    icon={ExclamationCircleIcon}
                                                    color="rose"
                                                />
                                                <Text color="rose">
                                                    {getErrorMessage(error)}
                                                </Text>
                                            </Flex>
                                        )}
                                        {!isLoading &&
                                            isExecuted &&
                                            queryResponse && (
                                                <Flex
                                                    justifyContent="start"
                                                    className="w-fit"
                                                >
                                                    {memoCount === pageSize ? (
                                                        <>
                                                            <Icon
                                                                icon={
                                                                    ExclamationCircleIcon
                                                                }
                                                                color="amber"
                                                                className="ml-0 pl-0"
                                                            />
                                                            <Text color="amber">
                                                                {`Row limit of ${numberDisplay(
                                                                    pageSize,
                                                                    0
                                                                )} reached, results are truncated`}
                                                            </Text>
                                                        </>
                                                    ) : (
                                                        <>
                                                            <Icon
                                                                icon={
                                                                    CheckCircleIcon
                                                                }
                                                                color="emerald"
                                                            />
                                                            <Text color="emerald">
                                                                Success
                                                            </Text>
                                                        </>
                                                    )}
                                                </Flex>
                                            )}
                                    </Flex>
                                </Flex>
                            </Transition.Child>
                        </Transition.Root> */}
            {/* <Flex flexDirection="row" className="gap-4">
                            <Card
                                onClick={() => {
                                    console.log('salam')
                                }}
                                className="p-3 cursor-pointer dark:ring-gray-500 hover:shadow-md"
                            >
                                <Subtitle className="font-semibold text-gray-800 mb-2">
                                    KPI
                                </Subtitle>
                            </Card>
                            <Card
                                onClick={() => {
                                    console.log('salam')
                                }}
                                className="p-3 cursor-pointer dark:ring-gray-500 hover:shadow-md"
                            >
                                <Subtitle className="font-semibold text-gray-800 mb-2">
                                    KPI
                                </Subtitle>
                            </Card>{' '}
                            <Card
                                onClick={() => {
                                    console.log('salam')
                                }}
                                className="p-3 cursor-pointer dark:ring-gray-500 hover:shadow-md"
                            >
                                <Subtitle className="font-semibold text-gray-800 mb-2">
                                    KPI
                                </Subtitle>
                            </Card>
                        </Flex> */}
            {/* <Flex
                            flexDirection="row"
                            justifyContent="start"
                            alignItems="center"
                            className="w-full  gap-1  pb-2 flex-wrap"
                            // style={{overflow:"hidden",overflowX:"scroll",overflowY: "hidden"}}
                        >
                            <Filter
                                type={'findings'}
                                // @ts-ignore
                                onApply={(e) => setQuery(e)}
                            />
                            <KFilter
                                // @ts-ignore
                                options={filters?.tags?.map((unique, index) => {
                                    return {
                                        label: unique.Key,
                                        value: unique.Key,
                                    }
                                })}
                                type="multi"
                                hasCondition={false}
                                condition={'is'}
                                //@ts-ignore
                                selectedItems={selectedFilter}
                                onChange={(values: string[]) => {
                                    // @ts-ignore
                                    setSelectedFilters(values)
                                }}
                                label={'Add filters'}
                                icon={PlusIcon}
                            />

                            {selectedFilter.map((item, index) => {
                                return (
                                    <div key={index}>
                                        {findFilters(item) &&
                                            findFilters(item)?.Key && (
                                                <>
                                                    <KFilter
                                                        // @ts-ignore
                                                        options={findFilters(
                                                            item
                                                        ).UniqueValues?.map(
                                                            (unique, index) => {
                                                                return {
                                                                    label: unique,
                                                                    value: unique,
                                                                }
                                                            }
                                                        )}
                                                        type="multi"
                                                        hasCondition={true}
                                                        condition={'is'}
                                                        //@ts-ignore
                                                        selectedItems={
                                                            query?.tags &&
                                                            //@ts-ignore
                                                            query?.tags[
                                                                //@ts-ignore
                                                                findFilters(
                                                                    item
                                                                )?.Key
                                                            ]
                                                                ? //@ts-ignore
                                                                  query?.tags[
                                                                      //@ts-ignore
                                                                      findFilters(
                                                                          item
                                                                      )?.Key
                                                                  ]
                                                                : []
                                                        }
                                                        onChange={(
                                                            values: string[]
                                                        ) => {
                                                            // @ts-ignore
                                                            if (
                                                                values.length >
                                                                0
                                                            ) {
                                                                setQuery(
                                                                    // @ts-ignore
                                                                    (
                                                                        prevSelectedItem
                                                                    ) => ({
                                                                        ...prevSelectedItem,
                                                                        tags: {
                                                                            ...prevSelectedItem?.tags,
                                                                            //@ts-ignore
                                                                            [findFilters(
                                                                                item
                                                                            )
                                                                                ?.Key]:
                                                                                values,
                                                                        },
                                                                    })
                                                                )
                                                            }
                                                        }}
                                                        //@ts-ignore
                                                        label={
                                                            //@ts-ignore
                                                            findFilters(item)
                                                                ?.Key
                                                        }
                                                        icon={TagIcon}
                                                    />
                                                </>
                                            )}
                                    </div>
                                )
                            })}
                        </Flex> */}

            {/* <Flex className="mt-2"> */}
            {/* <Table
                                id="inventory_queries"
                                columns={columns}
                                serverSideDatasource={serverSideRows}
                                loading={loading}
                                onRowClicked={(e) => {
                                    if (e.data) {
                                        setSelectedRow(e?.data)
                                    }
                                    setOpenSlider(true)
                                }}
                                options={{
                                    rowModelType: 'serverSide',
                                    serverSideDatasource: serverSideRows,
                                }}
                            /> */}
            {/* </Flex> */}
            {/* </Flex> */}
            {/* </Flex> */}
            {/* )} */}
            {/* <QueryDetail
                // type="resource"
                query={selectedRow}
                open={openSlider}
                onClose={() => setOpenSlider(false)}
                onRefresh={() => window.location.reload()}
                setTab={setTab}
            /> */}
        </>
    )
}
