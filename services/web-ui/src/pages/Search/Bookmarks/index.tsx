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
    Grid,
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
    IServerSideGetRowsParams,
    RowClickedEvent,
    ValueFormatterParams,
} from 'ag-grid-community'
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
import Table, { IColumn } from '../../../components/Table'

import {
    GithubComKaytuIoKaytuEnginePkgInventoryApiRunQueryResponse,
    Api,
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2,
    GithubComKaytuIoKaytuEnginePkgInventoryApiListQueryRequestV2,
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItem,
} from '../../../api/api'
import { isDemoAtom, queryAtom, runQueryAtom } from '../../../store'
import AxiosAPI from '../../../api/ApiConfig'

import { snakeCaseToLabel } from '../../../utilities/labelMaker'
import { numberDisplay } from '../../../utilities/numericDisplay'
import TopHeader from '../../../components/Layout/Header'
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
    Multiselect,
    Pagination,
    PropertyFilter,
} from '@cloudscape-design/components'
import { AppLayout, SplitPanel } from '@cloudscape-design/components'
import { useIntegrationApiV1EnabledConnectorsList } from '../../../api/integration.gen'
import axios from 'axios'
import UseCaseCard from '../../../components/Cards/BookmarkCard'

export const getTable = (
    headers: string[] | undefined,
    details: any[][] | undefined,
    isDemo: boolean
) => {
    const columns: IColumn<any, any>[] = []
    const rows: any[] = []
    const headerField = headers?.map((value, idx) => {
        if (headers.filter((v) => v === value).length > 1) {
            return `${value}-${idx}`
        }
        return value
    })
    if (headers && headers.length) {
        for (let i = 0; i < headers.length; i += 1) {
            const isHide = headers[i][0] === '_'
            columns.push({
                field: headerField?.at(i),
                headerName: snakeCaseToLabel(headers[i]),
                type: 'string',
                sortable: true,
                hide: isHide,
                resizable: true,
                filter: true,
                width: 170,
                cellRenderer: (param: ValueFormatterParams) => (
                    <span className={isDemo ? 'blur-sm' : ''}>
                        {param.value}
                    </span>
                ),
            })
        }
    }
    if (details && details.length) {
        for (let i = 0; i < details.length; i += 1) {
            const row: any = {}
            for (let j = 0; j < columns.length; j += 1) {
                row[headerField?.at(j) || ''] =
                    typeof details[i][j] === 'string'
                        ? details[i][j]
                        : JSON.stringify(details[i][j])
            }
            rows.push(row)
        }
    }
    const count = rows.length
    return {
        columns,
        rows,
        count,
    }
}

const columns: IColumn<
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2,
    any
>[] = [
    {
        field: 'id',
        headerName: 'ID',
        type: 'string',
        sortable: true,
        resizable: false,
    },
    {
        field: 'title',
        headerName: 'Title',
        type: 'string',
        sortable: true,
        resizable: false,
    },
    
    // {
    //     field: 'connectors',
    //     headerName: 'Service',
    //     type: 'string',
    //     sortable: true,
    //     resizable: false,
    // },
    // {
    //     field: 'connectors',
    //     headerName: 'Primary Table',
    //     type: 'string',
    //     sortable: true,
    //     resizable: false,
    // },
    // {
    //     type: 'string',
    //     width: 130,
    //     resizable: false,
    //     sortable: false,
    //     cellRenderer: (params: any) => (
    //         <Flex
    //             justifyContent="center"
    //             alignItems="center"
    //             className="h-full"
    //         >
    //             <PlayCircleIcon className="h-5 text-openg-500 mr-1" />
    //             <Text className="text-openg-500">Run query</Text>
    //         </Flex>
    //     ),
    // },
]
export interface Props {
    setTab: Function
}

export default function Bookmarks({ setTab }: Props) {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    // const [loading, setLoading] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [code, setCode] = useState(savedQuery || '')

        useState<GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2>()
    
    const [pageSize, setPageSize] = useState(1000)
 
    const [page, setPage] = useState(1)
    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [rows, setRows] = useState<any[]>([])
    const [filterQuery, setFilterQuery] = useState({
        tokens: [],
        operation: 'and',
    })
    const [properties, setProperties] = useState<any[]>([])
    const [options, setOptions] = useState<any[]>([])
    const [selectedOptions, setSelectedOptions] = useState([])
    const [isLoading, setLoading] = useState(false)
    const [integrations, setIntegrations] = useState<any[]>([])
    const [error, setError] = useState()
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

        let body = {
            is_bookmarked: true,
        }
        // @ts-ignore
        if (selectedOptions?.length > 0) {
            body = {
                is_bookmarked: true,
                // @ts-ignore
                categories: selectedOptions.map((item: any) => item.value),
            }
        }
        axios
            .post(`${url}/main/core/api/v3/queries`, body, config)
            .then((res) => {
                if (res?.data) {
                    setRows(res.data.items)
                    setTotalCount(res.data.total_count)
                    setTotalPage(Math.ceil(res.data.total_count / pageSize))
                }
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setError(err)
                setLoading(false)
            })
    }
    const getCategories = () => {
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
                `${url}/main/core/api/v3/queries/categories`,

                config
            )
            .then((res) => {
                if (res?.data) {
                    const temp: any = []
                    res.data?.categories?.map((item: any) => {
                        temp.push({
                            label: item.category,
                            value: item.category,
                        })
                    })
                    setOptions(temp)
                }
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                // setError(err)
                setLoading(false)
            })
    }
    const getIntegrations = () => {
setLoading(true)
axios
    .get(
        'https://raw.githubusercontent.com/opengovern/opengovernance/refs/heads/main/assets/integrations/integrations.json'
    )
    .then((res) => {
        if (res.data) {
            const arr = res.data
            // arr.sort(() => Math.random() - 0.5);
            setIntegrations(arr)
           
        }
        setLoading(false)
    })
    .catch((err) => {
        setError(err)
        setLoading(false)
    })

    }
    const FindLogos = (types : string []) => {
        const temp: string[] =[]
        types.map((type) => {
            const integration = integrations.find(
                (i) => i.integration_type === type
            )
            if(integration){
                temp.push(
                    `https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/${integration?.Icon}`
                )
            }
        })
        return temp



    }

    useEffect(() => {
        getRows()
        getCategories()
        getIntegrations()
    }, [])
    useEffect(() => {
        getRows()
    }, [selectedOptions])



    return (
        <>
            <TopHeader />
            {isLoading ? (
                <Spinner />
            ) : (
                <div className='flex flex-col'>
                    <Flex
                        className="w-full mb-3 mt-2 gap-2 flex-wrap"
                        flexDirection="row"
                        justifyContent="start"
                        alignItems="center"
                    >
                        <>
                            {options?.map((item: any) => {
                                return (
                                    <>
                                        <span
                                            onClick={() => {
                                                // check if the item is already selected remove it else add it
                                                if (
                                                    // @ts-ignore
                                                    selectedOptions?.find(
                                                        (i: any) =>
                                                            i.value ===
                                                            item.value
                                                    )
                                                ) {
                                                    // @ts-ignore
                                                    setSelectedOptions(
                                                        // @ts-ignore
                                                        selectedOptions?.filter(
                                                            (i: any) =>
                                                                i.value !==
                                                                item.value
                                                        )
                                                    )
                                                } else {
                                                    // @ts-ignore

                                                    setSelectedOptions([
                                                        // @ts-ignore
                                                        ...selectedOptions,
                                                        // @ts-ignore
                                                        item,
                                                    ])
                                                }
                                            }}
                                            className={`${
                                                selectedOptions?.find(
                                                    (i: any) =>
                                                        i.value === item.value
                                                )
                                                    ? 'bg-openg-400'
                                                    : 'bg-openg-900'
                                            } cursor-pointer text-white   p-3 border  rounded-3xl w-max`}
                                        >
                                            {item.label}
                                        </span>
                                    </>
                                )
                            })}
                        </>
                        {/* <Multiselect
                            // @ts-ignore
                            selectedOptions={selectedOptions}
                            className="w-1/3"
                            placeholder="Select a category"
                            // @ts-ignore

                            onChange={({ detail }) =>
                                // @ts-ignore
                                setSelectedOptions(detail.selectedOptions)
                            }
                            // Certificates | MLOps | DevOps | Keys | Certificates | Public Endpoints | Unprotected Data | Cloud Access | WAF
                            options={options}
                            loading={isLoading}
                        /> */}
                    </Flex>
                    <Flex className="gap-4 flex-wrap justify-start items-start"
                        style={{flex: "1 1 0"}}
                    >
                        {rows?.length === 0 && (
                            <>
                                <Spinner className="mt-2" />
                            </>
                        )}
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
                                className='h-full'
                                    style={{
                                        "width": `calc(calc(100% - ${
                                            rows.length >= 4 ? '3' : (rows.length-1)
                                        }rem) / ${
                                            rows.length >= 4 ? '4' : rows.length
                                        })`,
                                    }}
                                >
                                    <UseCaseCard
                                        // @ts-ignore
                                        title={q?.title}
                                        description={q?.description}
                                        logos={FindLogos(q?.integration_types)}
                                        onClick={() => {
                                            // @ts-ignore
                                            setSavedQuery(
                                                q?.query?.query_to_execute
                                            )
                                            setTab('3')
                                        }}
                                        tag="tag1"
                                    />
                                </div>
                            ))}
                    </Flex>
                </div>
            )}
            {error && (
                <Flex
                    flexDirection="col"
                    justifyContent="between"
                    className="absolute top-0 w-full left-0 h-full backdrop-blur"
                >
                    <Flex
                        flexDirection="col"
                        justifyContent="center"
                        alignItems="center"
                    >
                        <Title className="mt-6">Failed to load component</Title>
                        <Text className="mt-2">{getErrorMessage(error)}</Text>
                    </Flex>
                    <Button
                        variant="secondary"
                        className="mb-6"
                        color="slate"
                        onClick={getRows}
                    >
                        Try Again
                    </Button>
                </Flex>
            )}
        </>
    )
}


