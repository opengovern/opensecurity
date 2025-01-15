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
import DrawerPanel from '../../../../components/DrawerPanel'
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
import PolicyDetail from './PolicyDetail'
import { useComplianceApiV3ControlListFilters } from '../../../../api/compliance.gen'
import KFilter from '../../../../components/Filter'
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
import { getConnectorIcon } from '../../../../components/Cards/ConnectorCard'
import { useIntegrationApiV1EnabledConnectorsList } from '../../../../api/integration.gen'
import axios from 'axios'


export default function AllPolicy() {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loading, setLoading] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [code, setCode] = useState(savedQuery || '')
    const [selectedRow, setSelectedRow] =
        useState<any>()
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

    const getPolicyDetail = (id: string) => {
        
        // setLoading(true);
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
            .get(`${url}/main/compliance/api/v3/policies/${id}`, config)

            .then((resp) => {
                setSelectedRow(resp.data)
                setOpenDrawer(true)
                // setLoading(false)
            })
            .catch((err) => {
                // setLoading(false)
            })
    }
    

    const GetRows = () => {
        // debugger;
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
            .get(`${url}/main/compliance/api/v3/policies?cursor=${page}&per_page=15`, config)
            .then((resp) => {
                setRows(resp.data.policies)
                setTotalCount(resp.data.total_count)
                setTotalPage(Math.ceil(resp.data.total_count / 15))
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
   
   
     
    return (
        <>
            {/* <TopHeader /> */}

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
                                                        {selectedRow?.id}
                                                    </Title>
                                                </Flex>
                                            </>
                                        ) : (
                                            'Policy not selected'
                                        )
                                    }
                                >
                                    <PolicyDetail
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

                                        getPolicyDetail(row.ID)
                                        setOpen(true)
                                    }}
                                    columnDefinitions={[
                                        {
                                            id: 'id',
                                            header: 'Id',
                                            cell: (item) => item.ID,
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
                                            id: 'type',
                                            header: 'Type',
                                            cell: (item) =>
                                                item?.type ?  String(item?.type)
                                                    .charAt(0)
                                                    .toUpperCase() +
                                                String(item?.type).slice(1) : '',
                                            // sortingField: 'title',
                                            // minWidth: 400,
                                            maxWidth: 70,
                                        },
                                        {
                                            id: 'language',
                                            header: 'Language',
                                            cell: (item) => item?.language,
                                            maxWidth: 50,
                                        },
                                        {
                                            id: 'controls_count',
                                            header: 'Controls Ccount',
                                            maxWidth: 120,
                                            cell: (item) => (
                                                <>{item?.controls_count}</>
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
                                        { id: 'id', visible: true },
                                        // { id: 'title', visible: true },
                                        { id: 'type', visible: true },
                                        { id: 'language', visible: true },
                                        { id: 'controls_count', visible: true },
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
                                    // filter={
                                    //     <PropertyFilter
                                    //         // @ts-ignore
                                    //         query={filterQuery}
                                    //         tokenLimit={2}
                                    //         onChange={({ detail }) =>
                                    //             // @ts-ignore
                                    //             setFilterQuery(detail)
                                    //         }
                                    //         customGroupsText={[
                                    //             {
                                    //                 properties: 'Tags',
                                    //                 values: 'Tag values',
                                    //                 group: 'tags',
                                    //             },
                                    //         ]}
                                    //         // countText="5 matches"
                                    //         expandToViewport
                                    //         filteringAriaLabel="Find Controls"
                                    //         filteringPlaceholder="Find Controls"
                                    //         filteringOptions={options}
                                    //         filteringProperties={properties}
                                    //         asyncProperties
                                    //         virtualScroll
                                    //     />
                                    // }
                                    header={
                                        <Header className="w-full">
                                            Policies{' '}
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
