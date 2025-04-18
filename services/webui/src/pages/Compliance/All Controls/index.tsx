import {
    Flex,

    Title,
} from '@tremor/react'

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
    Api,
    PlatformEnginePkgControlApiListV2,
    PlatformEnginePkgControlDetailV3,
} from '../../../api/api'
import { isDemoAtom, queryAtom, runQueryAtom } from '../../../store'
import AxiosAPI from '../../../api/ApiConfig'

import ControlDetail from './ControlDetail'
import { useComplianceApiV3ControlListFilters } from '../../../api/compliance.gen'
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

import CustomPagination from '../../../components/Pagination'
import { severityBadge } from '../../../utilities/badge'



export default function AllControls() {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loading, setLoading] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [code, setCode] = useState(savedQuery || '')
    const [selectedRow, setSelectedRow] =
        useState<PlatformEnginePkgControlDetailV3>()
    const [openSlider, setOpenSlider] = useState(false)
    const [open, setOpen] = useState(false)

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
                setOpen(true)
                // setLoading(false)
            })
            .catch((err) => {
                // setLoading(false)
            })
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
            list_of_resources: query?.list_of_resources,
            primary_resource: query?.primary_resource,
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

    useEffect(() => {
        GetRows()
    }, [page,query])
    useEffect(() => {
        const temp_option = [
           
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
                key: 'list_of_resources',
                operators: ['='],
                propertyLabel: 'List of Resources',
                groupValuesLabel: 'List of Resources values',
            },
            {
                key: 'primary_resource',
                operators: ['='],
                propertyLabel: 'Primary Resources',
                groupValuesLabel: 'Primary Resources values',
            },
        ]
        filters?.provider?.map((item) => {
            temp_option.push({
                propertyKey: 'integrationType',
                value: item,
            })
        })
        filters?.parent_benchmark?.map((unique, index) => {
            temp_option.push({
                propertyKey: 'parent_benchmark',
                value: unique,
            })
        })
        filters?.list_of_resources?.map((unique, index) => {
            temp_option.push({
                propertyKey: 'list_of_resources',
                value: unique,
            })
        })
        filters?.primary_resource?.map((unique, index) => {
            temp_option.push({
                propertyKey: 'primary_resource',
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

    }, [filters])
    
     useEffect(() => {
        if(filterQuery){
            const temp_severity :any = []
            const temp_connector: any = []
            const temp_parent_benchmark: any = []
            const temp_list_of_resources: any = []
            const temp_primary_resource: any = []
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
                else if (item.propertyKey === 'list_of_resources') {
                    // @ts-ignore

                    temp_list_of_resources.push(item.value)
                }
                // @ts-ignore
                else if (item.propertyKey === 'primary_resource') {
                    // @ts-ignore

                    temp_primary_resource.push(item.value)
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
                list_of_resources:
                    temp_list_of_resources?.length > 0
                        ? temp_list_of_resources
                        : undefined,
                primary_resource:
                    temp_primary_resource?.length > 0
                        ? temp_primary_resource
                        : undefined,
                // @ts-ignore
                tags: temp_tags,
            })
        }
     }, [filterQuery])
     const [splitSize, setSplitSize] = useState(500)
     
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
                            splitPanelSize={splitSize}
                            onSplitPanelResize={({ detail }) => {
                                setSplitSize(detail.size)
                            }}
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
                                                    className="gap-2 sm:items-center sm:justify-center sm:flex-row flex-col items-start justify-start"
                                                >
                                                    <Title className="text-lg font-semibold ml-2 my-1">
                                                        {selectedRow?.title}
                                                    </Title>
                                                    {severityBadge(
                                                        selectedRow?.severity
                                                    )}
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
                                        setSelectedRow(undefined)
                                        getControlDetail(row.id)
                                        setOpen(true)
                                    }}
                                    columnDefinitions={[
                                        {
                                            id: 'id',
                                            header: 'ID',
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
                                            header: 'Policy Language',
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
                                                    {
                                                        item?.policy
                                                            ?.primary_resource
                                                    }
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
                                        // {
                                        //     id: 'parameters',
                                        //     header: 'Parametrized',
                                        //     maxWidth: 50,

                                        //     cell: (item) => (
                                        //         <>
                                        //             {item?.query?.parameters
                                        //                 .length > 0
                                        //                 ? 'True'
                                        //                 : 'False'}
                                        //         </>
                                        //     ),
                                        // },
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
                                        {
                                            id: 'severity',
                                            visible: true,
                                        },
                                        {
                                            id: 'integration_type',
                                            visible: false,
                                        },
                                        // {
                                        //     id: 'polity_type',
                                        //     visible: true,
                                        // },
                                        // { id: 'query', visible: true },

                                        { id: 'query', visible: true },
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
                                        <Header
                                            className="w-full"
                                            description={
                                                <>
                                                    <span className="">
                                                        A Control
                                                        programmatically defines
                                                        a check for a specific
                                                        rule, standard,
                                                        safeguard, or
                                                        requirement.
                                                        <br />
                                                        Aimed at managing risk,
                                                        achieving compliance, or
                                                        ensuring operational
                                                        integrity, it precisely
                                                        specifies <br />
                                                        <b>
                                                            what to verify
                                                        </b>{' '}
                                                        (e.g., maximum key age,
                                                        S3 bucket public access
                                                        status) and the exact
                                                        criteria determining its
                                                        pass or fail status.
                                                    </span>
                                                </>
                                            }
                                        >
                                            Controls{' '}
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
                    </Flex>
                </Flex>
            </Flex>
        </>
    )
}

//    getControlDetail(e.data.id)
// setOpenSlider(true)
