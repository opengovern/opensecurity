import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Button,
    Card,
    Flex,
    Grid,
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
} from '@tremor/react'
import {
    ChevronDoubleLeftIcon,
    ChevronDownIcon,
    ChevronUpIcon,
    CommandLineIcon,
    FunnelIcon,
    MagnifyingGlassIcon,
    PlayCircleIcon,
    TableCellsIcon,
} from '@heroicons/react/24/outline'
import { Fragment, useEffect, useMemo, useState } from 'react' // eslint-disable-next-line import/no-extraneous-dependencies
import { highlight, languages } from 'prismjs' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/components/prism-sql' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/themes/prism.css'
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
} from '../../../api/inventory.gen'
import Spinner from '../../../components/Spinner'
import { getErrorMessage } from '../../../types/apierror'
import { RenderObject } from '../../../components/RenderObject'

import { isDemoAtom, queryAtom, runQueryAtom } from '../../../store'
import { snakeCaseToLabel } from '../../../utilities/labelMaker'
import { numberDisplay } from '../../../utilities/numericDisplay'
import TopHeader from '../../../components/Layout/Header'
import KTable from '@cloudscape-design/components/table'
import {
    AppLayout,
    Box,
    ExpandableSection,
    Header,
    Modal,
    Pagination,
    SpaceBetween,
    SplitPanel,
    Tabs,
} from '@cloudscape-design/components'
import AceEditor from 'react-ace-builds'
// import 'ace-builds/src-noconflict/theme-github'
import 'ace-builds/css/ace.css'
import 'ace-builds/css/theme/cloud_editor.css'
import 'ace-builds/css/theme/cloud_editor_dark.css'
import 'ace-builds/css/theme/cloud_editor_dark.css'
import 'ace-builds/css/theme/twilight.css'
import 'ace-builds/css/theme/sqlserver.css'
import 'ace-builds/css/theme/xcode.css'

import CodeEditor from '@cloudscape-design/components/code-editor'
import KButton from '@cloudscape-design/components/button'
import AllQueries from '../All Query'
import View from '../View'
import Bookmarks from '../Bookmarks'
import axios from 'axios'
import CustomPagination from '../../../components/Pagination'
export const getTable = (
    headers: string[] | undefined,
    details: any[][] | undefined
) => {
    const columns: any[] = []
    const rows: any[] = []
    const column_def: any[] = []
    const headerField = headers?.map((value, idx) => {
        if (headers.filter((v) => v === value).length > 1) {
            return `${value}-${idx}`
        }
        return value
    })
    if (headers && headers.length) {
        for (let i = 0; i < headers.length; i += 1) {
            const isHide = headers[i][0] === '_'
            // columns.push({
            //     field: headerField?.at(i),
            //     headerName: snakeCaseToLabel(headers[i]),
            //     type: 'string',
            //     sortable: true,
            //     hide: isHide,
            //     resizable: true,
            //     filter: true,
            //     width: 170,
            //     cellRenderer: (param: ValueFormatterParams) => (
            //         <span className={isDemo ? 'blur-sm' : ''}>
            //             {param.value}
            //         </span>
            //     ),
            // })
            columns.push({
                id: headerField?.at(i),
                header: snakeCaseToLabel(headers[i]),
                // @ts-ignore
                cell: (item: any) => (
                    <>
                        {/* @ts-ignore */}
                        {typeof item[headerField?.at(i)] == 'string'
                            ? // @ts-ignore
                              item[headerField?.at(i)]
                            : // @ts-ignore
                              JSON.stringify(item[headerField?.at(i)])}
                    </>
                ),
                maxWidth: '200px',
                // sortingField: 'id',
                // isRowHeader: true,
                // maxWidth: 150,
            })
            column_def.push({
                id: headerField?.at(i),
                visible: !isHide,
            })
        }
    }
    if (details && details.length) {
        for (let i = 0; i < details.length; i += 1) {
            const row: any = {}
            for (let j = 0; j < columns.length; j += 1) {
                row[headerField?.at(j) || ''] = details[i][j]
                //     typeof details[i][j] === 'string'
                //         ? details[i][j]
                //         : JSON.stringify(details[i][j])
            }
            rows.push(row)
        }
    }
    const count = rows.length

    return {
        columns,
        column_def,
        rows,
        count,
    }
}

export default function Query() {
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loaded, setLoaded] = useState(false)
    const [savedQuery, setSavedQuery] = useAtom(queryAtom)
    const [code, setCode] = useState(savedQuery ? savedQuery : '')
    const [selectedIndex, setSelectedIndex] = useState(0)
    const [searchCategory, setSearchCategory] = useState('')
    const [selectedRow, setSelectedRow] = useState({})
    const [openDrawer, setOpenDrawer] = useState(false)
    const [openSearch, setOpenSearch] = useState(true)
    const [showEditor, setShowEditor] = useState(true)
    const [pageSize, setPageSize] = useState(1000)
    const [autoRun, setAutoRun] = useState(false)

    const [page, setPage] = useState(0)

    const [tab, setTab] = useState('0')
    const [preferences, setPreferences] = useState(undefined)
    const [integrations, setIntegrations] = useState([])
    const [selectedIntegration, setSelectedIntegration] = useState('')
    const [tables, setTables] = useState([])
    const [selectedTable, setSelectedTable] = useState('')
    const [columns, setColumns] = useState([])
    const [schemaLoading, setSchemaLoading] = useState(false)
    const [schemaLoading1, setSchemaLoading1] = useState(false)
    const [schemaLoading2, setSchemaLoading2] = useState(false)
    const [expanded, setExpanded] = useState(-1)
    const [expanded1, setExpanded1] = useState(-1)
    const [openIntegration, setOpenIntegration] = useState(false)
    const [openLayout, setOpenLayout] = useState(true)

    // const { response: categories, isLoading: categoryLoading } =
    //     useInventoryApiV2AnalyticsCategoriesList()

    const {
        response: queryResponse,
        isLoading,
        isExecuted,
        sendNow,
        error,
    } = useInventoryApiV1QueryRunCreate(
        {
            page: { no: 1, size: pageSize },
            // @ts-ignore
            engine: 'cloudql',
            query: code,
        },
        {},
        autoRun
    )

    useEffect(() => {
        if (autoRun) {
            setAutoRun(false)
        }
        if (queryResponse?.query?.length) {
            setSelectedIndex(2)
        } else setSelectedIndex(0)
    }, [queryResponse])

    useEffect(() => {
        if (!loaded && code.length > 0) {
            sendNow()
            setLoaded(true)
        }
    }, [page])

    useEffect(() => {
        if (code.length) setShowEditor(true)
    }, [code])

    const [ace, setAce] = useState()

    useEffect(() => {
        async function loadAce() {
            const ace = await import('ace-builds')
            await import('ace-builds/webpack-resolver')
            ace.config.set('useStrictCSP', true)
            // ace.config.setMode('ace/mode/sql')
            // @ts-ignore
            // ace.edit(element, {
            //     mode: 'ace/mode/sql',
            //     selectionStyle: 'text',
            // })

            return ace
        }

        loadAce()
            .then((ace) => {
                // @ts-ignore
                setAce(ace)
            })
            .finally(() => {})
    }, [])

    const memoCount = useMemo(
        () => getTable(queryResponse?.headers, queryResponse?.result).count,
        [queryResponse]
    )

    useEffect(() => {
        if (savedQuery.length > 0 && savedQuery !== '') {
            setCode(savedQuery)
            setAutoRun(true)
        }
    }, [savedQuery])

    const getIntegrations = () => {
        setSchemaLoading(true)
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
                    const temp: any = []
                    // arr.sort(() => Math.random() - 0.5);
                    arr?.map((integration: any) => {
                        if (integration.source_code !== '') {
                            temp.push(integration)
                        }
                    })
                    setIntegrations(temp)
                }
                setSchemaLoading(false)
            })
            .catch((err) => {
                setSchemaLoading(false)
            })
    }
    const getMasterSchema = (id: string) => {
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
        setSchemaLoading1(true)
        axios
            .get(
                `${url}/main/integration/api/v1/integration-types/${id}/table`,
                config
            )
            .then((res) => {
                if (res.data) {
                    setTables(res.data?.tables)
                }
                setSchemaLoading1(false)
            })
            .catch((err) => {
                setSchemaLoading1(false)
            })
    }

    useEffect(() => {
        getIntegrations()
    }, [])

    return (
        <>
            <Header
                className={`   rounded-xl mb-4    ${
                    false ? 'rounded-b-none' : ''
                }`}
                variant="h1"
                description={
                    <div className="group  important text-black  relative sm:flex hidden text-wrap justify-start">
                        Query all discovered assets across clouds and
                        integrations in SQL.
                    </div>
                }
            >
                CloudQL
            </Header>
            <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                className="w-full"
                toolsHide={true}
                navigationHide={true}
                splitPanelOpen={openLayout}
                onSplitPanelToggle={() => {
                    setOpenLayout(!openLayout)
                }}
                splitPanelSize={1200}
                splitPanel={
                    // @ts-ignore
                    <SplitPanel
                        // @ts-ignore
                        header={<>Saved Queries</>}
                    >
                        <>
                            <AllQueries
                                setTab={setTab}
                                setOpenLayout={setOpenLayout}
                            />
                        </>
                    </SplitPanel>
                }
                content={
                    <>
                        <Flex
                            className="w-full"
                            alignItems="start"
                            flexDirection="col"
                        >
                            <Flex
                                flexDirection="row"
                                className="gap-5 "
                                justifyContent="start"
                                alignItems="start"
                                style={{ flex: '1 1 0' }}
                            >
                                <Modal
                                    visible={openDrawer}
                                    onDismiss={() => setOpenDrawer(false)}
                                    header="Query Result"
                                    className="min-w-[500px]"
                                    size="large"
                                >
                                    <RenderObject obj={selectedRow} />
                                </Modal>
                                {openSearch ? (
                                    <>
                                        <Card className="p-3 rounded-xl w-1/4 h-full sm:flex hidden  ">
                                            <Flex
                                                flexDirection="col"
                                                justifyContent="start"
                                                alignItems="start"
                                                className="gap-2 overflow-y-scroll w-full max-h-[500px] "
                                            >
                                                <Text className="text-base text-black flex flex-row justify-between w-full">
                                                    <span className='w-full'>Plugin schema</span>
                                                    <Flex
                                                        justifyContent="end"
                                                        // className="mt-12"
                                                    >
                                                        <Button
                                                            variant="light"
                                                            onClick={() =>
                                                                setOpenSearch(
                                                                    false
                                                                )
                                                            }
                                                        >
                                                            <ChevronDoubleLeftIcon className="h-4" />
                                                        </Button>
                                                    </Flex>
                                                </Text>
                                                <>
                                                    {schemaLoading ? (
                                                        <>
                                                            <Spinner />
                                                        </>
                                                    ) : (
                                                        <>
                                                            {integrations?.map(
                                                                (
                                                                    integration: any,
                                                                    index
                                                                ) => {
                                                                    return (
                                                                        <>
                                                                            {/*   prettier-ignore */}
                                                                            {
                                                                                //  prettier-ignore
                                                                                (integration.install_state ==
                                                                    'installed' &&
                                                                integration.operational_status ==
                                                                    'enabled') ? (
                                                                    <>
                                                                        <ExpandableSection
                                                                            expanded={
                                                                                expanded ==
                                                                                index
                                                                            }
                                                                            onChange={({
                                                                                detail,
                                                                            }) => {
                                                                                if (
                                                                                    detail.expanded
                                                                                ) {
                                                                                    setExpanded(
                                                                                        index
                                                                                    )
                                                                                    setSelectedIntegration(
                                                                                        integration
                                                                                    )
                                                                                    getMasterSchema(
                                                                                        integration.plugin_id
                                                                                    )
                                                                                } else {
                                                                                    setExpanded(
                                                                                        -1
                                                                                    )
                                                                                }
                                                                            }}
                                                                            headerText={
                                                                                <span className=" text-sm font-normal ">
                                                                                    {
                                                                                        integration?.name
                                                                                    }
                                                                                </span>
                                                                            }
                                                                        >
                                                                            <>
                                                                                {schemaLoading1 ? (
                                                                                    <>
                                                                                        <Spinner />
                                                                                    </>
                                                                                ) : (
                                                                                    <div className="ml-4">
                                                                                        {' '}
                                                                                        <>
                                                                                            {tables &&
                                                                                                Object.entries(
                                                                                                    tables
                                                                                                )?.map(
                                                                                                    (
                                                                                                        item: any,
                                                                                                        index1
                                                                                                    ) => {
                                                                                                        return (
                                                                                                            <>
                                                                                                                <ExpandableSection
                                                                                                                    expanded={
                                                                                                                        expanded1 ==
                                                                                                                        index1
                                                                                                                    }
                                                                                                                    onChange={({
                                                                                                                        detail,
                                                                                                                    }) => {
                                                                                                                        if (
                                                                                                                            detail.expanded
                                                                                                                        ) {
                                                                                                                            setExpanded1(
                                                                                                                                index1
                                                                                                                            )
                                                                                                                            setSelectedTable(
                                                                                                                                item[0]
                                                                                                                            )
                                                                                                                            setColumns(
                                                                                                                                item[1]
                                                                                                                            )
                                                                                                                        } else {
                                                                                                                            setExpanded1(
                                                                                                                                -1
                                                                                                                            )
                                                                                                                        }
                                                                                                                    }}
                                                                                                                    headerText={
                                                                                                                        <span
                                                                                                                            onClick={(
                                                                                                                                e
                                                                                                                            ) => {
                                                                                                                                e.preventDefault()
                                                                                                                                e.stopPropagation()
                                                                                                                                setCode(
                                                                                                                                    code +
                                                                                                                                        `${item[0]}`
                                                                                                                                )
                                                                                                                            }}
                                                                                                                            className=" text-sm font-normal"
                                                                                                                        >
                                                                                                                            {
                                                                                                                                item[0]
                                                                                                                            }
                                                                                                                        </span>
                                                                                                                    }
                                                                                                                >
                                                                                                                    <>
                                                                                                                        {schemaLoading2 ? (
                                                                                                                            <>
                                                                                                                                <Spinner />
                                                                                                                            </>
                                                                                                                        ) : (
                                                                                                                            <>
                                                                                                                                {columns?.map(
                                                                                                                                    (
                                                                                                                                        column: any,
                                                                                                                                        index2
                                                                                                                                    ) => {
                                                                                                                                        return (
                                                                                                                                            <>
                                                                                                                                                <Flex className="pl-6 w-full">
                                                                                                                                                    <span className=" font-normal text-sm">
                                                                                                                                                        {
                                                                                                                                                            column.Name
                                                                                                                                                        }
                                                                                                                                                    </span>
                                                                                                                                                    <span>
                                                                                                                                                        (
                                                                                                                                                        {
                                                                                                                                                            column.Type
                                                                                                                                                        }

                                                                                                                                                        )
                                                                                                                                                    </span>
                                                                                                                                                </Flex>
                                                                                                                                            </>
                                                                                                                                        )
                                                                                                                                    }
                                                                                                                                )}
                                                                                                                            </>
                                                                                                                        )}
                                                                                                                    </>
                                                                                                                </ExpandableSection>
                                                                                                            </>
                                                                                                        )
                                                                                                    }
                                                                                                )}
                                                                                        </>
                                                                                    </div>
                                                                                )}
                                                                            </>
                                                                        </ExpandableSection>
                                                                    </>
                                                                ) : (
                                                                    <>
                                                                      <span  onClick={(e)=>{
                                                                        setSelectedIntegration(
                                                                            integration
                                                                        )
                                                                        setOpenIntegration(true)
                                                                          
                                                                      }} className=" text-sm text-gray-400  ml-5 cursor-pointer">
                                                                                    {
                                                                                        integration?.name
                                                                                    }
                                                                                </span>
                                                                    </>
                                                                )
                                                                            }
                                                                        </>
                                                                    )
                                                                }
                                                            )}
                                                        </>
                                                    )}
                                                </>
                                            </Flex>
                                        </Card>
                                    </>
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
                                            <Flex
                                                flexDirection="col"
                                                className="gap-4 w-4"
                                            >
                                                <TableCellsIcon />
                                                <Text className="rotate-90">
                                                    Plugin schema
                                                </Text>
                                            </Flex>
                                        </Button>
                                    </Flex>
                                )}

                                <Flex className="h-full">
                                    <CodeEditor
                                        ace={ace}
                                        language="sql"
                                        value={code}
                                        languageLabel="SQL"
                                        onChange={({ detail }) => {
                                            setSavedQuery('')
                                            setCode(detail.value)
                                            if (tab !== '3') {
                                                setTab('3')
                                            }
                                        }}
                                        preferences={preferences}
                                        onPreferencesChange={(e) =>
                                            // @ts-ignore
                                            setPreferences(e.detail)
                                        }
                                        loading={false}
                                        themes={{
                                            light: [
                                                'xcode',
                                                'cloud_editor',
                                                'sqlserver',
                                            ],
                                            dark: [
                                                'cloud_editor_dark',
                                                'twilight',
                                            ],
                                            // @ts-ignore
                                        }}
                                    />
                                </Flex>
                            </Flex>
                            <Flex flexDirection="col" className="w-full ">
                                <Flex flexDirection="col" className="mb-4">
                                    =
                                    <Flex className="w-full mt-4">
                                        <Flex
                                            justifyContent="start"
                                            className="gap-1"
                                        >
                                            <Text className="mr-2 w-fit">
                                                Maximum rows:
                                            </Text>
                                            <Select
                                                enableClear={false}
                                                className="w-44"
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
                                        </Flex>
                                        <Flex className="w-max gap-x-3">
                                            {!!code.length && (
                                                <KButton
                                                    className="  w-max min-w-max  "
                                                    onClick={() => setCode('')}
                                                    iconSvg={
                                                        <CommandLineIcon className="w-5 " />
                                                    }
                                                >
                                                    Clear editor
                                                </KButton>
                                            )}
                                            <KButton
                                                // icon={PlayCircleIcon}
                                                variant="primary"
                                                className="w-max  min-w-[300px]  "
                                                onClick={() => {
                                                    sendNow()
                                                    setLoaded(true)
                                                    setPage(0)
                                                }}
                                                disabled={!code.length}
                                                loading={
                                                    isLoading && isExecuted
                                                }
                                                loadingText="Running"
                                                iconSvg={
                                                    <PlayCircleIcon className="w-5 " />
                                                }
                                            >
                                                Run
                                            </KButton>
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
                                <Grid numItems={1} className="w-full">
                                    <KTable
                                        className="   min-h-[450px]   "
                                        // resizableColumns
                                        // variant="full-page"
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
                                        // stickyHeader={true}
                                        resizableColumns={true}
                                        // stickyColumns={
                                        //  {   first:1,
                                        //     last: 1}
                                        // }
                                        onRowClick={(event) => {
                                            const row = event.detail.item
                                            // @ts-ignore
                                            setSelectedRow(row)
                                            setOpenDrawer(true)
                                        }}
                                        columnDefinitions={
                                            getTable(
                                                queryResponse?.headers,
                                                queryResponse?.result
                                            ).columns
                                        }
                                        columnDisplay={
                                            getTable(
                                                queryResponse?.headers,
                                                queryResponse?.result
                                            ).column_def
                                        }
                                        enableKeyboardNavigation
                                        // @ts-ignore
                                        items={getTable(
                                            queryResponse?.headers,
                                            queryResponse?.result
                                        ).rows?.slice(
                                            page * 10,
                                            (page + 1) * 10
                                        )}
                                        loading={isLoading}
                                        loadingText="Loading resources"
                                        // stickyColumns={{ first: 0, last: 1 }}
                                        // stripedRows
                                        trackBy="id"
                                        empty={
                                            <Box
                                                margin={{
                                                    vertical: 'xs',
                                                }}
                                                textAlign="center"
                                                color="inherit"
                                            >
                                                <SpaceBetween size="m">
                                                    <b>No Results</b>
                                                </SpaceBetween>
                                            </Box>
                                        }
                                        header={
                                            <Header className="w-full">
                                                Results{' '}
                                                <span className=" font-medium">
                                                    {isLoading && isExecuted
                                                        ? '(?)'
                                                        : `(${memoCount})`}{' '}
                                                </span>
                                            </Header>
                                        }
                                        pagination={
                                            <CustomPagination
                                                currentPageIndex={page + 1}
                                                pagesCount={
                                                    // prettier-ignore
                                                    (isLoading &&
                                                            isExecuted)
                                                                ? 0
                                                                : Math.ceil(
                                                                      // @ts-ignore
                                                                      getTable(
                                                                          queryResponse?.headers,
                                                                          queryResponse?.result
                                                                      ).rows
                                                                          .length /
                                                                          10
                                                                  )
                                                }
                                                onChange={({ detail }: any) =>
                                                    setPage(
                                                        detail.currentPageIndex -
                                                            1
                                                    )
                                                }
                                            />
                                        }
                                    />
                                </Grid>
                            </Flex>
                        </Flex>
                    </>
                }
            />

            <Modal
                visible={openIntegration}
                onDismiss={() => setOpenIntegration(false)}
                header="Plugin Installation"
            >
                <div className="p-4">
                    <Text>
                        This plugin is not available. Plugins need to be
                        {/* @ts-ignore */}
                        {selectedIntegration?.install_state == 'not_installed'
                            ? ' installed'
                            : ' enabled'}{' '}
                        to fetch the schema.
                    </Text>

                    <Flex
                        justifyContent="end"
                        alignItems="center"
                        flexDirection="row"
                        className="gap-3"
                    >
                        <Button
                            // loading={loading}
                            disabled={false}
                            onClick={() => setOpenIntegration(false)}
                            className="mt-6"
                        >
                            Close
                        </Button>
                    </Flex>
                </div>
            </Modal>
        </>
    )
}
