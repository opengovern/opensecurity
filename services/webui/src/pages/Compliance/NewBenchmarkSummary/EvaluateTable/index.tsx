// @ts-nocheck
import { useAtomValue, useSetAtom } from 'jotai'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { Button, Flex } from '@tremor/react'
import {
    InformationCircleIcon,
    PlayCircleIcon,
} from '@heroicons/react/24/outline'
import { Checkbox, useCheckboxState } from 'pretty-checkbox-react'

import { PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary } from '../../../../api/api'

import { isDemoAtom, notificationAtom } from '../../../../store'

import {
    Box,
    DateRangePicker,
    Icon,
    SpaceBetween,
    Spinner,
} from '@cloudscape-design/components'

import KMulstiSelect from '@cloudscape-design/components/multiselect'

import KButton from '@cloudscape-design/components/button'
import axios from 'axios'
import KTable from '@cloudscape-design/components/table'

import {
    BreadcrumbGroup,
    Header,
    Link,
    Pagination,
    PropertyFilter,
} from '@cloudscape-design/components'

import {
    dateTimeDisplay,
    dateTimeDisplayAgo,
} from '../../../../utilities/dateDisplay'

import { useNavigate } from 'react-router-dom'
import CustomPagination from '../../../../components/Pagination'

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
        | PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary
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
        key: 'previous-7-days',
        amount: 7,
        unit: 'day',
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
                    setAccounts(res.data?.items)
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

    useEffect(() => {
        GetHistory()
    }, [page, jobStatus, date, selectedIntegrations])

    useEffect(() => {
        GetIntegrations()
    }, [])

    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 30 ? text.substring(0, 30) + '...' : text
        }
    }
    const checkStatusRedirect = (status) => {
        switch (status) {
            case 'CREATED':
                return false
            case 'QUEUED':
                return false
            case 'IN_PROGRESS':
                return true
            case 'RUNNERS_IN_PROGRESS':
                return true
            case 'SUMMARIZER_IN_PROGRESS':
                return true
            case 'SINK_IN_PROGRESS':
                return true
            case 'OLD_RESOURCE_DELETION':
                return true
            case 'SUCCEEDED':
                return true
            case 'COMPLETED':
                return true
            case 'FAILED':
                return true
            case 'COMPLETED_WITH_FAILURE':
                return true
            case 'TIMEOUT':
                return false
            case 'CANCELED':
                return false

            default:
                return false
        }
    }
    return (
        <>
            <div
                className="w-full"
                style={
                    window.innerWidth < 768
                        ? { width: `${window.innerWidth - 80}px` }
                        : {}
                }
            >
                {' '}
                <KTable
                    className="   min-h-[450px] w-full"
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
                    onRowClick={(event) => {
                        const row = event.detail.item
                        // @ts-ignore
                        // setSelected(row)
                        if (checkStatusRedirect(row.job_status)) {
                            navigate(`/compliance/${id}/report/${row.job_id}`)
                        }
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
                                <>{dateTimeDisplayAgo(item.updated_at)}</>
                            ),
                        },

                        {
                            id: 'integration_id',
                            header: 'Integration Id',
                            cell: (item) => (
                                // @ts-ignore
                                <>{item?.integration_info[0]?.integration_id}</>
                            ),
                        },

                        {
                            id: 'integration_name',
                            header: 'Integration(s)',
                            cell: (item) => {
                                const names = []
                                item?.integration_info?.map((i) => {
                                    names.push(i.name)
                                })
                                var unique = names.filter(
                                    (value, index, array) =>
                                        array.indexOf(value) === index
                                )
                                const length = unique.length

                                return (
                                    // @ts-ignore
                                    <>
                                        {length > 2 ? (
                                            <>
                                                <>
                                                    {unique[0]}
                                                    {length > 2 &&
                                                        ` + ${length - 1} more`}
                                                </>
                                            </>
                                        ) : (
                                            <>{unique.join(', ')}</>
                                        )}
                                    </>
                                )
                            },
                        },
                        {
                            id: 'integration_type',
                            header: 'Integration Type',
                            cell: (item) => {
                                const types = []
                                item?.integration_info?.map((i) => {
                                    types.push(i.integration_type)
                                })
                                var unique = types.filter(
                                    (value, index, array) =>
                                        array.indexOf(value) === index
                                )
                                return (
                                    // @ts-ignore
                                    <>{unique.join(', ')}</>
                                )
                            },
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
                                        // setSelected(item)
                                        navigate(
                                            `/compliance/${id}/report/${item?.job_id}`
                                        )
                                    }}
                                    variant="inline-link"
                                    ariaLabel={`Open Detail`}
                                >
                                    {window.innerWidth > 768 ? (
                                        'See details'
                                    ) : (
                                        <InformationCircleIcon className="w-5" />
                                    )}
                                </KButton>
                            ),
                        },
                    ]}
                    columnDisplay={[
                        { id: 'job_id', visible: true },
                        { id: 'updated_at', visible: true },
                        { id: 'job_status', visible: true },
                        { id: 'integration_id', visible: false },
                        {
                            id: 'integration_name',
                            visible: window.innerWidth > 640 ? true : false,
                        },
                        {
                            id: 'integration_type',
                            visible: window.innerWidth > 640 ? true : false,
                        },

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
                            className="gap-2 sm:flex-row flex-col"
                        >
                            <KMulstiSelect
                                className="sm:w-1/4 w-full"
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
                                className="sm:w-1/4 w-full"
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
                                    setSelectedIntegrations(
                                        detail.selectedOptions
                                    )
                                }}
                            />
                            {/* default last 24 */}
                            <DateRangePicker
                                className="sm:w-fit w-full"
                                onChange={({ detail }) => {
                                    setDate(detail.value)
                                }}
                                value={date}
                                relativeOptions={[
                                    {
                                        key: 'previous-30-minutes',
                                        amount: 30,
                                        unit: 'minute',
                                        type: 'relative',
                                    },
                                    {
                                        key: 'previous-3-hour',
                                        amount: 3,
                                        unit: 'hour',
                                        type: 'relative',
                                    },
                                    {
                                        key: 'previous-8-hours',
                                        amount: 8,
                                        unit: 'hour',
                                        type: 'relative',
                                    },
                                    {
                                        key: 'previous-1-days',
                                        amount: 1,
                                        unit: 'day',
                                        type: 'relative',
                                    },
                                    {
                                        key: 'previous-3-days',
                                        amount: 3,
                                        unit: 'day',
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
                                    iconName="refresh"
                                ></KButton>
                            }
                            className="w-full"
                        >
                            Jobs{' '}
                        </Header>
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
            </div>
        </>
    )
}
