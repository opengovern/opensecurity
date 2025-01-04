// @ts-nocheck
import { useAtomValue, useSetAtom } from 'jotai'
import { ICellRendererParams, ValueFormatterParams } from 'ag-grid-community'
import { useCallback, useEffect, useMemo, useState } from 'react'
import {
    Button,
    Flex,
    MultiSelect,
    MultiSelectItem,
    Text,
    Title,
} from '@tremor/react'
import {
    ArrowPathRoundedSquareIcon,
    CloudIcon,
    PlayCircleIcon,
} from '@heroicons/react/24/outline'
import { Checkbox, useCheckboxState } from 'pretty-checkbox-react'
import { useComplianceApiV1AssignmentsBenchmarkDetail } from '../../../../../api/compliance.gen'
import {
    GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkAssignedConnection,
    GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkEvaluationSummary,
} from '../../../../../api/api'
import DrawerPanel from '../../../../../components/DrawerPanel'
import Table, { IColumn } from '../../../../../components/Table'
import { isDemoAtom, notificationAtom } from '../../../../../store'
import KFilter from '../../../../../components/Filter'
import {
    Box,
    DateRangePicker,
    Icon,
    SpaceBetween,
    Spinner,
} from '@cloudscape-design/components'

import KMulstiSelect from '@cloudscape-design/components/multiselect'
import { Fragment, ReactNode } from 'react'
import { Dialog, Transition } from '@headlessui/react'
import Modal from '@cloudscape-design/components/modal'
import KButton from '@cloudscape-design/components/button'
import axios from 'axios'
import KTable from '@cloudscape-design/components/table'
import KeyValuePairs from '@cloudscape-design/components/key-value-pairs'
import Badge from '@cloudscape-design/components/badge'
import {
    BreadcrumbGroup,
    Header,
    Link,
    Pagination,
    PropertyFilter,
} from '@cloudscape-design/components'
import { AppLayout, SplitPanel } from '@cloudscape-design/components'
import { dateTimeDisplay } from '../../../../../utilities/dateDisplay'
import StatusIndicator from '@cloudscape-design/components/status-indicator'
import SeverityBar from '../../BenchmarkCard/SeverityBar'
import { useNavigate } from 'react-router-dom'

const JOB_STATUS = {
    CANCELED: 'stopped',
    SUCCEEDED: '',
    FAILED: 'error',
    SUMMARIZER_IN_PROGRESS: 'in-progress',
    SINK_IN_PROGRESS: 'in-progress',
    RUNNERS_IN_PROGRESS: 'in-progress',
}
interface IEvaluate {
    id: string | undefined
    assignmentsCount: number
    benchmarkDetail:
        | GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkEvaluationSummary
        | undefined
    onEvaluate: (c: string[]) => void
}

export default function EvaluateTable({
    id,
    benchmarkDetail,
    assignmentsCount,
    onEvaluate,
}: IEvaluate) {
    const [open, setOpen] = useState(false)
    const isDemo = useAtomValue(isDemoAtom)
    const [openConfirm, setOpenConfirm] = useState(false)
    const [connections, setConnections] = useState<string[]>([])
    const [loading, setLoading] = useState(false)
    const [detailLoading, setDetailLoading] = useState(false)

    const [accounts, setAccounts] = useState()
    const [selected, setSelected] = useState()
    const [detail, setDetail] = useState()

    const [page, setPage] = useState(1)
    const [integrationData, setIntegrationData] = useState()
    const [loadingI, setLoadingI] = useState(false)

    const [selectedIntegrations, setSelectedIntegrations] = useState()

    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [jobStatus, setJobStatus] = useState()
    const setNotification = useSetAtom(notificationAtom)

    const navigate = useNavigate()
    const today = new Date()
    const lastWeek = new Date(
        today.getFullYear(),
        today.getMonth(),
        today.getDate() - 7
    )

    const [date, setDate] = useState({
        key: 'previous-6-hours',
        amount: 6,
        unit: 'hour',
        type: 'relative',
    })
    const GetHistory = () => {
        
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
        const temp_status = []
        if (jobStatus && jobStatus?.length > 0) {
            jobStatus.map((status) => {
                temp_status.push(status.value)
            })
        }
        const integrations = []
        if (selectedIntegrations && selectedIntegrations?.length > 0) {
            selectedIntegrations.map((integration) => {
                integrations.push({
                    integration_id: integration.value,
                })
            })
        }

        let body = {
            cursor: page,
            per_page: 20,
            job_status: temp_status,
            integration_info: integrations,
        }
        if (date) {
            if (date.type == 'relative') {
                body.interval = `${date.amount} ${date.unit}s`
            } else {
                body.start_time = date.startDate
                body.end_time = date.endDate
            }
        }
        axios
            .post(
                `${url}/main/schedule/api/v3/benchmark/${id}/run-history`,
                body,
                config
            )
            .then((res) => {
                if (!res.data.items) {
                    setAccounts([])
                    setTotalCount(0)
                    setTotalPage(0)
                } else {
                    setAccounts(res.data.items)
                    setTotalCount(res.data.total_count)
                    setTotalPage(Math.ceil(res.data.total_count / 20))
                }

                setLoading(false)
            })
            .catch((err) => {
                setLoading(false)
                console.log(err)
            })
    }
    const GetIntegrations = () => {
        
        setLoadingI(true)
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
                `${url}/main/schedule/api/v3/benchmark/run-history/integrations`,

                config
            )
            .then((res) => {
                setIntegrationData(res.data)

                setLoadingI(false)
            })
            .catch((err) => {
                setLoadingI(false)
                console.log(err)
            })
    }

    const GetDetail = () => {
        
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
                // @ts-ignore
                `${url}/main/compliance/api/v3/job-report/${selected.job_id}/summary `,
                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                navigate(`/compliance/${id}/report/${selected.job_id}`)
                // setDetail(res.data)
                // setDetailLoading(false)
            })
            .catch((err) => {
                const message = err.response.data.message
                setNotification({
                    text: `${message}`,
                    type: 'error',
                })
                // setDetailLoading(false)

                console.log(err)
                setLoading(false)
            })
    }
    useEffect(() => {
        GetHistory()
    }, [page, jobStatus, date, selectedIntegrations])

    useEffect(() => {
        GetIntegrations()
    }, [])

    useEffect(() => {
        if (selected) {
            GetDetail()
        }
    }, [selected])
    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 30 ? text.substring(0, 30) + '...' : text
        }
    }
    return (
        <>
            {/* <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                toolsHide={true}
                navigationHide={true}
                splitPanelOpen={open}
                onSplitPanelToggle={() => {
                    setOpen(!open)
                    if (open) {
                        setSelected(undefined)
                    }
                }}
                splitPanel={
                    // @ts-ignore
                    <SplitPanel
                        // @ts-ignore
                        header={
                            selected ? (
                                <>{`Job No ${selected?.job_id} Selected`}</>
                            ) : (
                                'Job not selected'
                            )
                        }
                    >
                        {detailLoading ? (
                            <>
                                <Spinner />
                            </>
                        ) : (
                            <>
                                <Flex
                                    flexDirection="col"
                                    className="w-full"
                                    alignItems="center"
                                    justifyContent="center"
                                >
                                    <KeyValuePairs
                                        className="w-full"
                                        columns={6}
                                        items={[
                                            {
                                                label: 'Job ID',
                                                value: selected?.job_id,
                                            },
                                            {
                                                label: 'Benchmark ID',
                                                value: detail?.benchmark_id,
                                            },
                                            {
                                                label: 'Benchmark Title',
                                                value: detail?.benchmark_title,
                                            },
                                            {
                                                label: 'Last Evaulated at',
                                                value: (
                                                    <>
                                                        {dateTimeDisplay(
                                                            selected?.updated_at
                                                        )}
                                                    </>
                                                ),
                                            },
                                            {
                                                label: 'Job Status',
                                                value: (
                                                    <StatusIndicator
                                                        type={
                                                            JOB_STATUS[
                                                                selected
                                                                    ?.job_status
                                                            ]
                                                        }
                                                    >
                                                        {selected?.job_status}
                                                    </StatusIndicator>
                                                ),
                                            },
                                        ]}
                                    />
                                    <Flex className="w-1/2 mt-2">
                                        <SeverityBar benchmark={detail} />
                                    </Flex>
                                </Flex>
                            </>
                        )}
                    </SplitPanel>
                }
                content={
                   <></>
                }
            /> */}
            <KTable
                className="   min-h-[450px]"
                // resizableColumns
                // variant="full-page"
                renderAriaLive={({ firstIndex, lastIndex, totalItemsCount }) =>
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
                    // @ts-ignore
                    setSelected(row)
                }}
                columnDefinitions={[
                    {
                        id: 'job_id',
                        header: 'Id',
                        cell: (item) => item.job_id,
                        sortingField: 'id',
                        isRowHeader: true,
                    },
                    {
                        id: 'updated_at',
                        header: 'Last Updated at',
                        cell: (item) => (
                            // @ts-ignore
                            <>{dateTimeDisplay(item.updated_at)}</>
                        ),
                    },

                    {
                        id: 'integration_id',
                        header: 'Integration Id',
                        cell: (item) => (
                            // @ts-ignore
                            <>{item.integration_info[0]?.integration_id}</>
                        ),
                    },

                    {
                        id: 'integration_name',
                        header: 'Integration Name',
                        cell: (item) => (
                            // @ts-ignore
                            <>{item.integration_info[0]?.name}</>
                        ),
                    },
                    {
                        id: 'integration_type',
                        header: 'Integration Type',
                        cell: (item) => (
                            // @ts-ignore
                            <>{item.integration_info[0]?.integration_type}</>
                        ),
                    },

                    {
                        id: 'job_status',
                        header: 'Job Status',
                        cell: (item) => (
                            // @ts-ignore
                            <>{item.job_status}</>
                        ),
                    },
                    {
                        id: 'action',
                        header: '',
                        cell: (item) => (
                            // @ts-ignore
                            <KButton
                                onClick={() => {
                                    setSelected(item)
                                }}
                                variant="inline-link"
                                ariaLabel={`Open Detail`}
                            >
                                See Details
                            </KButton>
                        ),
                    },
                ]}
                columnDisplay={[
                    { id: 'job_id', visible: true },
                    { id: 'updated_at', visible: true },
                    { id: 'job_status', visible: true },
                    { id: 'integration_id', visible: false },
                    { id: 'integration_name', visible: true },
                    { id: 'integration_type', visible: true },

                    // { id: 'conformanceStatus', visible: true },
                    // { id: 'severity', visible: true },
                    // { id: 'evaluatedAt', visible: true },

                    { id: 'action', visible: true },
                ]}
                enableKeyboardNavigation
                // @ts-ignore
                items={accounts}
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
                    <Flex
                        flexDirection="row"
                        justifyContent="start"
                        alignItems="start"
                        className="gap-2"
                    >
                        <KMulstiSelect
                            className="w-1/4"
                            placeholder="Filter by Job status"
                            selectedOptions={jobStatus}
                            options={[
                                {
                                    label: 'SUCCEEDED',
                                    value: 'SUCCEEDED',
                                },
                                {
                                    label: 'FAILED',
                                    value: 'FAILED',
                                },
                                {
                                    label: 'CREATED',
                                    value: 'CREATED',
                                },
                                {
                                    label: 'RUNNERS_IN_PROGRESS',
                                    value: 'RUNNERS_IN_PROGRESS',
                                },
                                {
                                    label: 'SINK_IN_PROGRESS',
                                    value: 'SINK_IN_PROGRESS',
                                },
                                {
                                    label: 'CANCELED',
                                    value: 'CANCELED',
                                },
                                {
                                    label: 'TIMEOUT',
                                    value: 'TIMEOUT',
                                },
                                {
                                    label: 'SUMMARIZER_IN_PROGRESS',
                                    value: 'SUMMARIZER_IN_PROGRESS',
                                },
                            ]}
                            onChange={({ detail }) => {
                                setJobStatus(detail.selectedOptions)
                            }}
                        />
                        <KMulstiSelect
                            className="w-1/4"
                            placeholder="Filter by Integration"
                            selectedOptions={selectedIntegrations}
                            filteringType="auto"
                            options={integrationData?.map((i) => {
                                return {
                                    label: i.id_name,
                                    value: i.integration_id,
                                    description: truncate(i.id),
                                }
                            })}
                            loadingText="Loading Integrations"
                            loading={loadingI}
                            onChange={({ detail }) => {
                                setSelectedIntegrations(detail.selectedOptions)
                            }}
                        />
                        {/* default last 24 */}
                        <DateRangePicker
                            onChange={({ detail }) => {
                                setDate(detail.value)
                            }}
                            value={date}
                            relativeOptions={[
                                {
                                    key: 'previous-5-minutes',
                                    amount: 5,
                                    unit: 'minute',
                                    type: 'relative',
                                },
                                {
                                    key: 'previous-30-minutes',
                                    amount: 30,
                                    unit: 'minute',
                                    type: 'relative',
                                },
                                {
                                    key: 'previous-1-hour',
                                    amount: 1,
                                    unit: 'hour',
                                    type: 'relative',
                                },
                                {
                                    key: 'previous-6-hours',
                                    amount: 6,
                                    unit: 'hour',
                                    type: 'relative',
                                },
                                {
                                    key: 'previous-7-days',
                                    amount: 7,
                                    unit: 'day',
                                    type: 'relative',
                                },
                            ]}
                            absoluteFormat="long-localized"
                            hideTimeOffset
                            // rangeSelectorMode={'absolute-only'}
                            isValidRange={(range) => {
                                if (range.type === 'absolute') {
                                    const [startDateWithoutTime] =
                                        range.startDate.split('T')
                                    const [endDateWithoutTime] =
                                        range.endDate.split('T')
                                    if (
                                        !startDateWithoutTime ||
                                        !endDateWithoutTime
                                    ) {
                                        return {
                                            valid: false,
                                            errorMessage:
                                                'The selected date range is incomplete. Select a start and end date for the date range.',
                                        }
                                    }
                                    if (
                                        new Date(range.startDate) -
                                            new Date(range.endDate) >
                                        0
                                    ) {
                                        return {
                                            valid: false,
                                            errorMessage:
                                                'The selected date range is invalid. The start date must be before the end date.',
                                        }
                                    }
                                }
                                return { valid: true }
                            }}
                            i18nStrings={{}}
                            placeholder="Filter by Job Range"
                        />
                    </Flex>
                }
                header={
                    <Header
                        counter={totalCount ? `(${totalCount})` : ''}
                        actions={
                            <KButton
                                onClick={() => {
                                    GetHistory()
                                }}
                            >
                                Reload
                            </KButton>
                        }
                        className="w-full"
                    >
                        Jobs{' '}
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
        </>
    )
}
