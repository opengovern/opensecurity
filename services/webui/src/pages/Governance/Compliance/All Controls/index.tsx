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
// import { highlight, languages } from 'prismjs' // eslint-disable-next-line import/no-extraneous-dependencies
// import 'prismjs/components/prism-sql' // eslint-disable-next-line import/no-extraneous-dependencies
// import 'prismjs/themes/prism.css'
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
} from '../../../../api/inventory.gen'
import Spinner from '../../../../components/Spinner'
import { getErrorMessage } from '../../../../types/apierror'
import { RenderObject } from '../../../../components/RenderObject'

import {
    PlatformEnginePkgInventoryApiRunQueryResponse,
    Api,
    PlatformEnginePkgInventoryApiSmartQueryItemV2,
    PlatformEnginePkgControlApiListV2ResponseItem,
    PlatformEnginePkgControlApiListV2ResponseItemQuery,
    PlatformEnginePkgControlApiListV2,
    PlatformEnginePkgControlDetailV3,
    TypesFindingSeverity,
} from '../../../../api/api'
import { isDemoAtom, queryAtom, runQueryAtom } from '../../../../store'
import AxiosAPI from '../../../../api/ApiConfig'

import { snakeCaseToLabel } from '../../../../utilities/labelMaker'
import { numberDisplay } from '../../../../utilities/numericDisplay'
import TopHeader from '../../../../components/Layout/Header'
import ControlDetail from './ControlDetail'
import { useComplianceApiV3ControlListFilters } from '../../../../api/compliance.gen'
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
import { useIntegrationApiV1EnabledConnectorsList } from '../../../../api/integration.gen'



export default function AllControls() {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loading, setLoading] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [code, setCode] = useState(savedQuery || '')
    const [selectedRow, setSelectedRow] =
        useState<PlatformEnginePkgControlDetailV3>()
    const [openDrawer, setOpenDrawer] = useState(false)
    const [openSlider, setOpenSlider] = useState(false)
    const [open, setOpen] = useState(false)

    const [openSearch, setOpenSearch] = useState(true)
    const [showEditor, setShowEditor] = useState(true)
    const isDemo = useAtomValue(isDemoAtom)
    const [pageSize, setPageSize] = useState(1000)
    const [autoRun, setAutoRun] = useState(false)
    const [selectedFilter, setSelectedFilters] = useState<string[]>([])
    const [engine, setEngine] = useState('odysseus-sql')
    const [query, setQuery] =
        useState<PlatformEnginePkgControlApiListV2>()
    const [rows, setRows] = useState<any[]>()
    const [page, setPage] = useState(1)
    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [properties, setProperties] = useState<any[]>([])
    const [options, setOptions] = useState<any[]>([])
    const [filterQuery, setFilterQuery] = useState({
        tokens: [
            { propertyKey: 'severity', value: 'high', operator: '=' },
            { propertyKey: 'severity', value: 'medium', operator: '=' },
            { propertyKey: 'severity', value: 'low', operator: '=' },
            { propertyKey: 'severity', value: 'critical', operator: '=' },
            { propertyKey: 'severity', value: 'none', operator: '=' },
        ],
        operation: 'or',
    })
    // const { response: categories, isLoading: categoryLoading } =
    //     useInventoryApiV2AnalyticsCategoriesList()
    // const { response: queries, isLoading: queryLoading } =
    //     useInventoryApiV2QueryList({
    //         titleFilter: '',
    //         Cursor: 0,
    //         PerPage:25
    //     })
    const { response: filters, isLoading: filtersLoading } =
        useComplianceApiV3ControlListFilters()

    const getControlDetail = (id: string) => {
        const api = new Api()
        api.instance = AxiosAPI
        // setLoading(true);
        api.compliance
            .apiV3ControlDetail(id)
            .then((resp) => {
                setSelectedRow(resp.data)
                setOpenDrawer(true)
                // setLoading(false)
            })
            .catch((err) => {
                // setLoading(false)
            })
    }
    const findFilters = (key: string) => {
        const temp = filters?.tags.filter((item, index) => {
            if (item.Key === key) {
                return item
            }
        })

        if (temp) {
            return temp[0]
        }
        return undefined
    }

    const GetRows = () => {
        // debugger;
        setLoading(true)
        const api = new Api()
        api.instance = AxiosAPI
        
        // @ts-ignore
       

        let body = {
            integration_types: query?.connector,
            severity: query?.severity,
            list_of_tables: query?.list_of_tables,
            primary_table: query?.primary_table,
            root_benchmark: query?.root_benchmark,
            parent_benchmark: query?.parent_benchmark,
            tags: query?.tags,
            cursor: page,
            per_page: 15,
        }
        // if (!body.integrationType) {
        //     delete body['integrationType']
        // } else {
        //     // @ts-ignore
        //     body['integrationType'] = [body?.integrationType]
        // }

        api.compliance
            .apiV2ControlList(body)
            .then((resp) => {
                if(resp.data?.items){
                setRows(resp.data.items)

                }
                else{
                    setRows([])
                }
                setTotalCount(resp.data?.total_count)
                setTotalPage(Math.ceil(resp.data?.total_count / 15))
                setLoading(false)
            })
            .catch((err) => {
                setLoading(false)

                console.log(err)
                // params.fail()
            })
    }
    const {
        response: Types,
        isLoading: TypesLoading,
        isExecuted: TypesExec,
    } = useIntegrationApiV1EnabledConnectorsList(0, 0)
    useEffect(() => {
        GetRows()
    }, [page,query])
    useEffect(() => {
        const temp_option = [
            { propertyKey: 'connector', value: 'AWS' },
            { propertyKey: 'connector', value: 'Azure' },
            { propertyKey: 'severity', value: 'high' },
            { propertyKey: 'severity', value: 'medium' },
            { propertyKey: 'severity', value: 'low' },
            { propertyKey: 'severity', value: 'critical' },
            { propertyKey: 'severity', value: 'none' },
        ]

        const property = [
            {
                key: 'severity',
                operators: ['='],
                propertyLabel: 'Severity',
                groupValuesLabel: 'Severity values',
            },
            {
                key: 'integrationType',
                operators: ['='],
                propertyLabel: 'integrationType',
                groupValuesLabel: 'integrationType values',
            },
            {
                key: 'parent_benchmark',
                operators: ['='],
                propertyLabel: 'Parent Benchmark',
                groupValuesLabel: 'Parent Benchmark values',
            },
            {
                key: 'list_of_tables',
                operators: ['='],
                propertyLabel: 'List of Tables',
                groupValuesLabel: 'List of Tables values',
            },
            {
                key: 'primary_table',
                operators: ['='],
                propertyLabel: 'Primary Service',
                groupValuesLabel: 'Primary Service values',
            },
        ]
        Types?.integration_types?.map((item)=>{
            temp_option.push({
                propertyKey: 'integrationType',
                value: item.platform_name,
            })
        })
        filters?.parent_benchmark?.map((unique, index) => {
            temp_option.push({
                propertyKey: 'parent_benchmark',
                value: unique,
            })
        })
        filters?.list_of_tables?.map((unique, index) => {
            temp_option.push({
                propertyKey: 'list_of_tables',
                value: unique,
            })
        })
        filters?.primary_table?.map((unique, index) => {
            temp_option.push({
                propertyKey: 'primary_table',
                value: unique,
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
        setProperties(property)
        setOptions(temp_option)

    }, [filters,Types])
    
     useEffect(() => {
        if(filterQuery){
            const temp_severity :any = []
            const temp_connector: any = []
            const temp_parent_benchmark: any = []
            const temp_list_of_tables: any = []
            const temp_primary_table: any = []
            let temp_tags = {}
            filterQuery.tokens.map((item, index) => {
                // @ts-ignore
                if (item.propertyKey === 'severity') {
                    // @ts-ignore

                    temp_severity.push(item.value)
                }
                // @ts-ignore
                else if (item.propertyKey === 'connector') {
                    // @ts-ignore

                    temp_connector.push(item.value)
                }
                // @ts-ignore
                else if (item.propertyKey === 'parent_benchmark') {
                    // @ts-ignore

                    temp_parent_benchmark.push(item.value)
                }
                // @ts-ignore
                else if (item.propertyKey === 'list_of_tables') {
                    // @ts-ignore

                    temp_list_of_tables.push(item.value)
                }
                // @ts-ignore
                else if (item.propertyKey === 'primary_table') {
                    // @ts-ignore

                    temp_primary_table.push(item.value)
                }
                
                else {
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
            setQuery({
                connector:
                    temp_connector?.length > 0 ? temp_connector : undefined,
                severity: temp_severity?.length > 0 ? temp_severity : undefined,
                parent_benchmark:
                    temp_parent_benchmark?.length > 0
                        ? temp_parent_benchmark
                        : undefined,
                list_of_tables:
                    temp_list_of_tables?.length > 0
                        ? temp_list_of_tables
                        : undefined,
                primary_table:
                    temp_primary_table?.length > 0
                        ? temp_primary_table
                        : undefined,
                // @ts-ignore
                tags: temp_tags,
            })
        }
     }, [filterQuery])
     
     
    return (
        <>

            <Flex alignItems="start">
              
                <Flex flexDirection="col" className="w-full ">
                 

                    <Flex className=" mt-2">
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
                                            'Control not selected'
                                        )
                                    }
                                >
                                    <ControlDetail
                                        // type="resource"
                                        selectedItem={selectedRow}
                                        open={openSlider}
                                        onClose={() => setOpenSlider(false)}
                                        onRefresh={() => {}}
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

                                        getControlDetail(row.id)
                                        setOpen(true)
                                    }}
                                    columnDefinitions={[
                                        {
                                            id: 'title',
                                            header: 'Title',
                                            cell: (item) => item.title,
                                            // sortingField: 'id',
                                            isRowHeader: true,
                                            maxWidth: 150,
                                        },
                                        {
                                            id: 'integration_type',
                                            header: 'Integration Type',
                                            cell: (item) =>
                                                item.integration_type,
                                            // sortingField: 'title',
                                            // minWidth: 400,
                                            maxWidth: 70,
                                        },
                                        {
                                            id: 'polity_type',
                                            header: 'Policy Type',
                                            cell: (item) =>
                                                String(item?.policy?.type)
                                                    .charAt(0)
                                                    .toUpperCase() +
                                                String(
                                                    item?.policy?.type
                                                ).slice(1),
                                            // sortingField: 'title',
                                            // minWidth: 400,
                                            maxWidth: 50,
                                        },
                                        {
                                            id: 'query',
                                            header: 'Primary Table',
                                            maxWidth: 120,
                                            cell: (item) => (
                                                <>
                                                    {item?.query?.primary_table}
                                                </>
                                            ),
                                        },
                                        {
                                            id: 'severity',
                                            header: 'Severity',
                                            // sortingField: 'severity',
                                            cell: (item) => (
                                                <Badge
                                                    // @ts-ignore
                                                    color={`severity-${item.severity}`}
                                                >
                                                    {item.severity
                                                        .charAt(0)
                                                        .toUpperCase() +
                                                        item.severity.slice(1)}
                                                </Badge>
                                            ),
                                            maxWidth: 50,
                                        },
                                        {
                                            id: 'parameters',
                                            header: 'Parametrized',
                                            maxWidth: 50,

                                            cell: (item) => (
                                                <>
                                                    {item?.query?.parameters
                                                        .length > 0
                                                        ? 'True'
                                                        : 'False'}
                                                </>
                                            ),
                                        },
                                    ]}
                                    columnDisplay={[
                                        {
                                            id: 'title',
                                            visible: true,
                                        },
                                        {
                                            id: 'integration_type',
                                            visible: true,
                                        },
                                        {
                                            id: 'polity_type',
                                            visible: true,
                                        },
                                        // { id: 'query', visible: true },
                                        {
                                            id: 'severity',
                                            visible: true,
                                        },
                                        { id: 'parameters', visible: true },
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
                                            ]}
                                            // countText="5 matches"
                                            expandToViewport
                                            filteringAriaLabel="Find Controls"
                                            filteringPlaceholder="Find Controls"
                                            filteringOptions={options}
                                            filteringProperties={properties}
                                            asyncProperties
                                            virtualScroll
                                        />
                                    }
                                    header={
                                        <Header className="w-full">
                                            Controls{' '}
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
                    </Flex>
                </Flex>
            </Flex>
        </>
    )
}

//    getControlDetail(e.data.id)
// setOpenSlider(true)
