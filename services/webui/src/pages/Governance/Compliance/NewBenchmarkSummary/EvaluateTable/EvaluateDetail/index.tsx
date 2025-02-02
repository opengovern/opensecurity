// @ts-nocheck
import { useAtomValue } from 'jotai'
import { useCallback, useEffect, useMemo, useState } from 'react'
import {
    Button,
    Col,
    Flex,
    Grid,
    MultiSelect,
    MultiSelectItem,
    Text,
    Title,
} from '@tremor/react'
import {
    ArrowPathRoundedSquareIcon,
    CloudIcon,
    InformationCircleIcon,
    PlayCircleIcon,
} from '@heroicons/react/24/outline'
import { Checkbox, useCheckboxState } from 'pretty-checkbox-react'
import { useComplianceApiV1AssignmentsBenchmarkDetail } from '../../../../../api/compliance.gen'
import {
    PlatformEnginePkgComplianceApiBenchmarkAssignedConnection,
    PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary,
} from '../../../../../api/api'
import DrawerPanel from '../../../../../components/DrawerPanel'
import Table, { IColumn } from '../../../../../components/Table'
import { isDemoAtom } from '../../../../../store'
import KFilter from '../../../../../components/Filter'
import {
    Box,
    DateRangePicker,
    Icon,
    SpaceBetween,
    Spinner,
    Tabs,
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
import {
    dateTimeDisplay,
    shortDateTimeDisplayDelta,
} from '../../../../../../utilities/dateDisplay'
import StatusIndicator from '@cloudscape-design/components/status-indicator'
import SeverityBar from './SeverityBar'
import { useParams } from 'react-router-dom'
import { RunDetail } from './types'
import TopHeader from '../../../../../../components/Layout/Header'
import CustomPagination from '../../../../../../components/Pagination'
import Findings from './Findings'

const JOB_STATUS = {
    CANCELED: 'stopped',
    SUCCEEDED: '',
    FAILED: 'error',
    SUMMARIZER_IN_PROGRESS: 'in-progress',
    SINK_IN_PROGRESS: 'in-progress',
    RUNNERS_IN_PROGRESS: 'in-progress',
}

export default function EvaluateDetail() {
    const { id, benchmarkId } = useParams()
    const [detail, setDetail] = useState()
    const [detailLoading, setDetailLoading] = useState(false)
    const [runDetail, setRunDetail] = useState()
    const [page, setPage] = useState(0)
    const [resourcePage, setResourcePage] = useState(0)
    const [open, setOpen] = useState(false)
    const [selectedControl, setSelectedControl] = useState()
    const [resources, setResources] = useState()
    const [jobDetail, setJobDetail] = useState()
    const [fullLoading, setFullLoading] = useState(false)
    const [runnerLoading, setRunnerLoading] = useState(false)
    const [runners, setRunners] = useState([])
    const [selectedRunner, setSelectedRunner] = useState()
    const [runnerOpen, setRunnerOpen] = useState(false)
    const [runnerPage, setRunnerPage] = useState(0)
    const [sortOrder, setSortOrder] = useState('asc')
    const [openInfo, setOpenInfo] = useState(false)
    const [runnerTable, setRunnerTable] = useState(false)

    const [sortField, setSortField] = useState({
        id: 'severity',
        header: 'Severity',
        sortingField: 'severity',
        cell: (item) => (
            <Badge
                // @ts-ignore
                color={`severity-${item.severity}`}
            >
                {item.severity.charAt(0).toUpperCase() + item.severity.slice(1)}
            </Badge>
        ),
        maxWidth: 100,
        sortingComparator: (a, b) => {
            if (a?.severity === b?.severity) {
                return 0
            }
            if (a?.severity === 'critical') {
                return -1
            }
            if (b?.severity === 'critical') {
                return 1
            }
            if (a?.severity === 'high') {
                return -1
            }
            if (b?.severity === 'high') {
                return 1
            }
            if (a?.severity === 'medium') {
                return -1
            }
            if (b?.severity === 'medium') {
                return 1
            }
            if (a?.severity === 'low') {
                return -1
            }
            if (b?.severity === 'low') {
                return 1
            }
        },
    })
    const [integrations, setIntegrations] = useState()
    const [integrationOpen, setIntegrationOpen] = useState(false)
    const [integrationDetail, setIntegrationDetail] = useState()
    const [integrationDetailLoading, setIntegrationDetailLoading] =
        useState(false)
    const GetDetail = () => {
        setDetailLoading(true)
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
                `${url}/main/compliance/api/v3/job-report/${id}/summary`,
                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                setDetail(res.data)
                setIntegrations(res.data?.job_details?.integration_ids)

                setDetailLoading(false)
            })
            .catch((err) => {
                setDetailLoading(false)
                console.log(err)
            })
    }

    const GetControls = () => {
        setDetailLoading(true)
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
                `${url}/main/compliance/api/v3/job-report/${id}/details/by-control `,
                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                const temp = []
                Object.entries(res.data?.controls).map(([key, value]) => {
                    temp.push({
                        title: key,
                        severity: value.severity,
                        alarms: value.alarms,
                        oks: value.oks,
                    })
                })
                setRunDetail(temp)
                setDetailLoading(false)
            })
            .catch((err) => {
                setDetailLoading(false)
                console.log(err)
            })
    }
    const GetResults = () => {
        setDetailLoading(true)
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
                `${url}/main/compliance/api/v3/quick/scan/${id}?with_incidents=${detail?.with_incidents}&controls=${selectedControl.title}`,
                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                const temp = []
                const alarms =
                    res?.data?.controls[selectedControl.title]?.results

                Object.entries(alarms)?.map((key) => {
                    if (key.length > 1) {
                        key[1]?.map((alarm) => {
                            temp.push({
                                resource_id: alarm.resource_id,
                                resource_type: alarm.resource_type,
                                reason: alarm.reason,
                                type: key[0],
                            })
                        })
                    }
                })
                setResources(temp)
                setDetailLoading(false)
                setOpen(true)
            })
            .catch((err) => {
                setDetailLoading(false)
                console.log(err)
            })
    }

    const GetFullResults = () => {
        setFullLoading(true)
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
                `${url}/main/compliance/api/v3/quick/scan/${id}?with_incidents=${detail?.with_incidents}`,
                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                const json = JSON.stringify(res.data)
                var data = new Blob([json], { type: 'application/json' })
                var csvURL = window.URL.createObjectURL(data)
                var tempLink = document.createElement('a')
                tempLink.href = csvURL
                tempLink.setAttribute('download', 'result.json')
                tempLink.click()
                setFullLoading(false)
            })
            .catch((err) => {
                setFullLoading(false)
                console.log(err)
            })
    }
    const GetIntegrationDetail = () => {
        setIntegrationDetailLoading(true)
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
        const body = {
            integration_id: integrations,
        }
        axios
            .post(
                // @ts-ignore
                `${url}/main/integration/api/v1/integrations/list`,
                body,
                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                setIntegrationDetailLoading(false)
                setIntegrationDetail(res.data?.integrations)
            })
            .catch((err) => {
                console.log(err)
            })
    }
    const GetRunners = () => {
        setRunnerLoading(true)
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
                `${url}/main/schedule/api/v3/jobs/compliance/${id}/runners`,

                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                setRunners(res?.data)
                setRunnerLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setRunnerLoading(false)
            })
    }
    const GetJobDetail = () => {
        setRunnerLoading(true)
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
                `${url}/main/schedule/api/v3/job/compliance/${id}`,

                config
            )
            .then((res) => {
                //   setAccounts(res.data.integrations)
                setJobDetail(res?.data)
                setRunnerLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setRunnerLoading(false)
            })
    }

    useEffect(() => {
        GetDetail()
        GetControls()
        GetRunners()
        GetJobDetail()
    }, [])

    useEffect(() => {
        if (selectedControl) {
            GetResults()
        }
    }, [selectedControl])
    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 30 ? text.substring(0, 30) + '...' : text
        }
    }
    const getRows = () => {
        if (sortOrder == 'asc') {
            return runDetail
                ?.sort(sortField.sortingComparator)
                .slice(page * 10, (page + 1) * 10)
        }
        return runDetail
            ?.sort(sortField.sortingComparator)
            .reverse()
            .slice(page * 10, (page + 1) * 10)
    }
    const GetTabs = () => {
        const temp: any = []

        if (jobDetail?.with_incidents) {
            temp.push(
                {
                    label: 'Incidents',
                    id: '2',
                    disabled: !jobDetail?.with_incidents,
                    disabledReason: 'Job Runned without incidents',
                    content: (
                        <>
                            <Findings id={id ? id : ''} tab={0} />
                        </>
                    ),
                },
                {
                    label: 'Control summary',
                    id: '3',
                    disabled: !jobDetail?.with_incidents,
                    disabledReason: 'Job Runned without incidents',
                    content: (
                        <>
                            <Findings id={id ? id : ''} tab={2} />
                        </>
                    ),
                },
                {
                    label: 'Resource summary',
                    id: '4',
                    disabled: !jobDetail?.with_incidents,
                    disabledReason: 'Job Runned without incidents',
                    content: (
                        <>
                            <Findings id={id ? id : ''} tab={3} />
                        </>
                    ),
                }
            )
        }
        if (!jobDetail?.with_incidents) {
            temp.push({
                label: 'Controls assements',
                id: '0',
                disabled: !(runDetail && runDetail?.length > 0),
                disabledReason: 'Job is still in progress',
                content: (
                    <>
                        {' '}
                        <KTable
                            className="p-3   min-h-[550px]"
                            // resizableColumns
                            renderAriaLive={({
                                firstIndex,
                                lastIndex,
                                totalItemsCount,
                            }) =>
                                `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                            }
                            variant="full-page"
                            sortingDescending={
                                sortOrder == 'desc' ? true : false
                            }
                            sortingColumn={sortField}
                            onSortingChange={({ detail }) => {
                                const desc = detail.isDescending
                                    ? 'desc'
                                    : 'asc'
                                setSortOrder(desc)
                                setSortField(detail.sortingColumn)
                            }}
                            onRowClick={(event) => {
                                const row = event.detail.item
                                // @ts-ignore
                                setSelectedControl(row)
                                setResourcePage(0)
                            }}
                            columnDefinitions={[
                                {
                                    id: 'id',
                                    header: 'Control ID',
                                    cell: (item) => item.title,
                                    // sortingField: 'id',
                                    isRowHeader: true,
                                },

                                {
                                    id: 'severity',
                                    header: 'Severity',
                                    sortingField: 'severity',
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
                                    maxWidth: 100,
                                    sortingComparator: (a, b) => {
                                        if (a?.severity === b?.severity) {
                                            return 0
                                        }
                                        if (a?.severity === 'critical') {
                                            return -1
                                        }
                                        if (b?.severity === 'critical') {
                                            return 1
                                        }
                                        if (a?.severity === 'high') {
                                            return -1
                                        }
                                        if (b?.severity === 'high') {
                                            return 1
                                        }
                                        if (a?.severity === 'medium') {
                                            return -1
                                        }
                                        if (b?.severity === 'medium') {
                                            return 1
                                        }
                                        if (a?.severity === 'low') {
                                            return -1
                                        }
                                        if (b?.severity === 'low') {
                                            return 1
                                        }
                                    },
                                },

                                {
                                    id: 'incidents',
                                    header: 'OK',
                                    sortingField: 'oks',

                                    cell: (item) => (
                                        // @ts-ignore
                                        <>
                                            {/**@ts-ignore */}
                                            {item.oks}
                                        </>
                                    ),
                                    // minWidth: 50,
                                    maxWidth: 100,
                                    sortingComparator: (a, b) => {
                                        console.log(a)
                                        console.log(b)

                                        if (a?.oks === b?.oks) {
                                            return 0
                                        }
                                        if (a?.oks > b?.oks) {
                                            return -1
                                        }
                                        if (a?.oks < b?.oks) {
                                            return 1
                                        }
                                    },
                                },
                                {
                                    id: 'passing_resources',
                                    header: 'Alarms ',
                                    sortingField: 'alarms',

                                    cell: (item) => (
                                        // @ts-ignore
                                        <>{item.alarms}</>
                                    ),
                                    maxWidth: 100,
                                    sortingComparator: (a, b) => {
                                        if (a?.alarms === b?.alarms) {
                                            return 0
                                        }
                                        if (a?.alarms > b?.alarms) {
                                            return -1
                                        }
                                        if (a?.alarms < b?.alarms) {
                                            return 1
                                        }
                                    },
                                },

                                {
                                    id: 'action',
                                    header: '',
                                    cell: (item) => (
                                        // @ts-ignore
                                        <KButton
                                            onClick={() => {
                                                setSelectedControl(item)
                                                setResourcePage(0)
                                            }}
                                            className="w-full"
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
                                { id: 'id', visible: true },
                                { id: 'title', visible: false },
                                {
                                    id: 'connector',
                                    visible: false,
                                },
                                { id: 'query', visible: false },
                                {
                                    id: 'severity',
                                    visible: true,
                                },
                                {
                                    id: 'incidents',
                                    visible: true,
                                },
                                {
                                    id: 'passing_resources',
                                    visible: true,
                                },
                                {
                                    id: 'noncompliant_resources',
                                    visible: true,
                                },

                                { id: 'action', visible: true },
                            ]}
                            enableKeyboardNavigation
                            items={runDetail ? getRows() : []}
                            loading={detailLoading}
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
                            header={
                                <Header
                                    actions={
                                        <CustomPagination
                                            currentPageIndex={page + 1}
                                            pagesCount={Math.ceil(
                                                runDetail?.length / 10
                                            )}
                                            onChange={({ detail }) =>
                                                setPage(
                                                    detail.currentPageIndex - 1
                                                )
                                            }
                                        />
                                    }
                                    counter={
                                        runDetail?.length
                                            ? `(${runDetail?.length})`
                                            : ''
                                    }
                                    className="w-full"
                                >
                                    Controls{' '}
                                </Header>
                            }
                        />
                    </>
                ),
            })
        }
        return temp
    }
    const checkStatusRedirect = (status) => {
        switch (status) {
            case 'CREATED':
                return 0
            case 'QUEUED':
                return 0
            case 'IN_PROGRESS':
                return 1
            case 'RUNNERS_IN_PROGRESS':
                return 1
            case 'SUMMARIZER_IN_PROGRESS':
                return 1
            case 'SINK_IN_PROGRESS':
                return 1
            case 'OLD_RESOURCE_DELETION':
                return 1
            case 'SUCCEEDED':
                return 2
            case 'COMPLETED':
                return 2
            case 'FAILED':
                return 3
            case 'COMPLETED_WITH_FAILURE':
                return 3
            case 'TIMEOUT':
                return 3
            case 'CANCELED':
                return 3

            default:
                return false
        }
    }
    const GetRunnerStatus = () => {}
    return (
        <>
            <div
                className="w-full "
                style={
                    window.innerWidth < 768
                        ? { width: `${window.innerWidth - 80}px` }
                        : {}
                }
            >
                <BreadcrumbGroup
                    className="w-full"
                    onClick={(event) => {
                        // event.preventDefault()
                    }}
                    items={[
                        {
                            text: 'Compliance',
                            href: `/compliance`,
                        },
                        {
                            text: 'Frameworks',
                            href: `/compliance/${benchmarkId}`,
                        },
                        { text: 'Job Report', href: `#` },
                    ]}
                    ariaLabel="Breadcrumbs"
                />
                <Grid
                    className="w-full gap-4"
                    numItems={12}
                    style={
                        window.innerWidth > 768 ? { gridAutoRows: '1fr' } : {}
                    }
                >
                    <Col numColSpan={12} numColSpanSm={4} className="h-full">
                        <Flex
                            flexDirection="col"
                            className="w-full  bg-white p-4 border border-slate-400 rounded-lg mt-4 gap-4 h-full"
                            alignItems="start"
                            justifyContent="start"
                        >
                            <Header
                                info={
                                    <KButton
                                        onClick={() => {
                                            setOpenInfo(true)
                                        }}
                                        variant="inline-link"
                                        iconName="status-info"
                                    ></KButton>
                                }
                            >
                                Job
                            </Header>
                            <KeyValuePairs
                                className="w-full "
                                columns={2}
                                items={[
                                    {
                                        label: ' ID',
                                        value: id,
                                    },
                                    {
                                        label: 'Last Update',
                                        value: (
                                            <>
                                                {dateTimeDisplay(
                                                    jobDetail?.updated_at
                                                )}
                                            </>
                                        ),
                                    },
                                    {
                                        label: ' Status',
                                        value: jobDetail?.job_status,
                                    },

                                    {
                                        label: 'Total Time',
                                        value: (
                                            <KButton
                                                onClick={() => {
                                                    setRunnerTable(true)
                                                }}
                                                variant="inline-link"
                                                // iconName="status-info"
                                            >
                                                {shortDateTimeDisplayDelta(
                                                    jobDetail?.updated_at,
                                                    jobDetail?.created_at
                                                )}
                                            </KButton>
                                        ),
                                    },
                                ]}
                            />
                        </Flex>
                    </Col>
                    <Col numColSpan={12} numColSpanSm={8} className="h-full">
                        <Flex
                            flexDirection="col"
                            className="w-full  bg-white p-4 pb-0 rounded-lg mt-4 border border-slate-400 gap-4 h-full "
                            alignItems="start"
                            justifyContent="start"
                        >
                            <Header>Compliance </Header>
                            <KeyValuePairs
                                className="w-full"
                                columns={window.innerWidth > 640 ? 4 : 1}
                                items={[
                                    {
                                        label: ' Score',

                                        value: (
                                            <Flex className="flex-col gap-2 justify-start items-start">
                                                <Link
                                                    // variant="awsui-value-large"
                                                    fontSize="heading-xl"
                                                    variant="secondary"
                                                    href="#"
                                                    ariaLabel="Running instances (14)"
                                                >
                                                    {((
                                                        1 -
                                                        detail?.job_details
                                                            ?.job_score
                                                            ?.failed_controls /
                                                            detail?.job_details
                                                                ?.job_score
                                                                ?.total_controls
                                                    )?.toFixed(2) * 100
                                                        ? (
                                                              1 -
                                                              detail
                                                                  ?.job_details
                                                                  ?.job_score
                                                                  ?.failed_controls /
                                                                  detail
                                                                      ?.job_details
                                                                      ?.job_score
                                                                      ?.total_controls
                                                          )?.toFixed(2) * 100
                                                        : '-- ') + '%'}
                                                </Link>
                                                <Flex className="w-full ">
                                                    <SeverityBar
                                                        benchmark={
                                                            detail?.job_details
                                                                ?.job_score
                                                        }
                                                    />
                                                </Flex>
                                            </Flex>
                                        ),
                                    },
                                    {
                                        label: 'Framework',
                                        value: (
                                            <Link
                                                fontSize="heading-xl"
                                                variant="secondary"
                                                href="#"
                                                ariaLabel="Running instances (14)"
                                            >
                                                {detail?.job_details?.framework
                                                    ?.title
                                                    ? detail?.job_details
                                                          ?.framework?.title
                                                    : '--'}
                                            </Link>
                                        ),
                                    },
                                    {
                                        label: 'Incidents',
                                        value: (
                                            <Link
                                                // variant="awsui-value-large"
                                                fontSize="heading-xl"
                                                variant="secondary"
                                                href="#"
                                                ariaLabel="Running instances (14)"
                                            >
                                                {detail?.job_details?.results
                                                    ?.alarm
                                                    ? detail?.job_details
                                                          ?.results?.alarm
                                                    : '--'}
                                            </Link>
                                        ),
                                    },
                                    {
                                        label: 'Integrations',
                                        value: (
                                            <>
                                                <Link
                                                    // variant="awsui-value-large"
                                                    fontSize="heading-xl"
                                                    variant="secondary"
                                                    href="#"
                                                    ariaLabel="Running instances (14)"
                                                >
                                                    <Flex
                                                        className="gap-2 "
                                                        justifyContent="start"
                                                    >
                                                        {
                                                            jobDetail
                                                                ?.integration_info
                                                                ?.length
                                                        }
                                                        <KButton
                                                            onClick={() => {
                                                                if (
                                                                    integrations
                                                                ) {
                                                                    GetIntegrationDetail()
                                                                    setIntegrationOpen(
                                                                        true
                                                                    )
                                                                }
                                                            }}
                                                            variant="inline-link"
                                                            iconName="status-info"
                                                        ></KButton>
                                                    </Flex>
                                                </Link>
                                            </>
                                        ),
                                    },
                                    // {
                                    //     label: '',
                                    //     columnSpan: 2,
                                    //     value: <></>,
                                    // },

                                    // {
                                    //     label: 'Benchmark ID',
                                    //     value: jobDetail?.framework_id,
                                    // },

                                    // {
                                    //     label: 'Last Evaulated at',
                                    //     value: (
                                    //         <>
                                    //             {dateTimeDisplay(
                                    //                 jobDetail?.updated_at
                                    //             )}
                                    //         </>
                                    //     ),
                                    // },
                                    // {
                                    //     label: 'Total Time',
                                    //     value: shortDateTimeDisplayDelta(
                                    //         jobDetail?.updated_at,
                                    //         jobDetail?.created_at
                                    //     ),
                                    // },
                                    // {
                                    //     label: 'Create With Incidents',
                                    //     value: jobDetail?.with_incidents
                                    //         ? 'True'
                                    //         : 'False',
                                    // },

                                    // {
                                    //     label: 'Download Results',
                                    //     value: (
                                    //         <KButton
                                    //             iconName="download"
                                    //             className="mt-2"
                                    //             variant="inline-icon"
                                    //             onClick={() => GetFullResults()}
                                    //             loading={fullLoading}
                                    //         >
                                    //             Click here for results in JSON
                                    //         </KButton>
                                    //     ),
                                    // },
                                    // {
                                    //     label: 'Severity Status',
                                    //     columnSpan: 2,
                                    //     value: (
                                    //         <>
                                    //             <Flex className="w-full ">
                                    //                 <SeverityBar
                                    //                     benchmark={
                                    //                         detail?.job_details
                                    //                             ?.job_score
                                    //                     }
                                    //                 />
                                    //             </Flex>
                                    //         </>
                                    //     ),
                                    // },
                                ]}
                            />
                        </Flex>
                    </Col>
                </Grid>
                <Flex className="w-100 bg-white p-4 border border-slate-400 rounded-lg mt-8">
                    <Tabs tabs={GetTabs()} />
                </Flex>
                <Modal
                    visible={open}
                    size="max"
                    onDismiss={() => setOpen(false)}
                    header={``}
                >
                    <KTable
                        className="p-1   min-h-[550px]"
                        // resizableColumns
                        renderAriaLive={({
                            firstIndex,
                            lastIndex,
                            totalItemsCount,
                        }) =>
                            `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                        }
                        // sortingDescending={sortOrder == 'desc' ? true : false}
                        // onRowClick={(event) => {
                        //     const row = event.detail.item
                        //     // @ts-ignore
                        //     setSelectedControl(row)
                        // }}
                        columnDefinitions={[
                            {
                                id: 'resource_id',
                                header: 'Resource ID',
                                cell: (item) => item.resource_id,
                                sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: '70px',
                            },

                            {
                                id: 'resource_type',
                                header: 'Title',
                                sortingField: 'severity',
                                cell: (item) => item.resource_type,
                                maxWidth: '70px',
                            },
                            {
                                id: 'reason',
                                header: 'Reason',
                                sortingField: 'severity',
                                cell: (item) => item.reason,
                                maxWidth: 150,
                            },
                            {
                                id: 'type',
                                header: 'Type',
                                sortingField: 'severity',
                                cell: (item) => item.type,
                                maxWidth: 150,
                            },
                        ]}
                        columnDisplay={[
                            { id: 'resource_id', visible: true },
                            { id: 'resource_type', visible: true },
                            { id: 'reason', visible: true },
                            { id: 'type', visible: true },
                        ]}
                        enableKeyboardNavigation
                        items={
                            resources
                                ? resources.slice(
                                      resourcePage * 7,
                                      (resourcePage + 1) * 7
                                  )
                                : []
                        }
                        loading={detailLoading}
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
                        header={
                            <Header className="w-full">
                                Resources{' '}
                                <span className=" font-medium">
                                    ({resources?.length})
                                </span>
                            </Header>
                        }
                        pagination={
                            <CustomPagination
                                currentPageIndex={resourcePage + 1}
                                pagesCount={Math.ceil(resources?.length / 7)}
                                onChange={({ detail }) =>
                                    setResourcePage(detail.currentPageIndex - 1)
                                }
                            />
                        }
                    />
                </Modal>
                <Modal
                    visible={integrationOpen}
                    size="large"
                    onDismiss={() => setIntegrationOpen(false)}
                    header={``}
                >
                    <KTable
                        className="p-3   min-h-[550px]"
                        // resizableColumns
                        renderAriaLive={({
                            firstIndex,
                            lastIndex,
                            totalItemsCount,
                        }) =>
                            `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                        }
                        // sortingDescending={sortOrder == 'desc' ? true : false}
                        // onRowClick={(event) => {
                        //     const row = event.detail.item
                        //     // @ts-ignore
                        //     setSelectedControl(row)
                        // }}
                        columnDefinitions={[
                            {
                                id: 'integration_id',
                                header: 'Integartion id',
                                cell: (item) => item.integration_id,
                                sortingField: 'id',
                                isRowHeader: true,
                                // maxWidth: '70px',
                            },

                            {
                                id: 'provider_id',
                                header: 'Provider id',
                                sortingField: 'severity',
                                cell: (item) => item.provider_id,
                                // maxWidth: '70px',
                            },
                            {
                                id: 'name',
                                header: 'Name',
                                sortingField: 'severity',
                                cell: (item) => item.name,
                                // maxWidth: 150,
                            },
                            {
                                id: 'integration_type',
                                header: 'integration Type',
                                sortingField: 'severity',
                                cell: (item) => item.integration_type,
                                // maxWidth: 100,
                            },
                            {
                                id: 'state',
                                header: 'State',
                                sortingField: 'severity',
                                cell: (item) => (
                                    <>
                                        <Badge
                                            // @ts-ignore
                                            color={
                                                // @ts-ignore
                                                item?.state?.toLowerCase() ===
                                                'active'
                                                    ? 'green'
                                                    : 'red'
                                            }
                                            // @ts-ignore
                                        >
                                            {/* @ts-ignore */}
                                            {item.state}
                                        </Badge>
                                    </>
                                ),
                                // maxWidth: 50,
                            },
                        ]}
                        columnDisplay={[
                            { id: 'integration_id', visible: true },
                            { id: 'provider_id', visible: true },
                            { id: 'name', visible: true },
                            { id: 'integration_type', visible: true },
                            { id: 'state', visible: true },
                        ]}
                        enableKeyboardNavigation
                        items={integrationDetail ? integrationDetail : []}
                        loading={integrationDetailLoading}
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
                        header={
                            <Header className="w-full">
                                Integrations{' '}
                                <span className=" font-medium">
                                    ({integrationDetail?.length})
                                </span>
                            </Header>
                        }
                    />
                </Modal>

                <Modal
                    visible={openInfo}
                    size="large"
                    onDismiss={() => setOpenInfo(false)}
                    header={'Job Detail'}
                >
                    <KeyValuePairs
                        columns={3}
                        items={[
                            {
                                label: 'Benchmark ID',
                                value: jobDetail?.framework_id,
                            },
                            // {
                            //     label: 'Last Evaulated at',
                            //     value: (
                            //         <>
                            //             {dateTimeDisplay(jobDetail?.updated_at)}
                            //         </>
                            //     ),
                            // },
                            // {
                            //     label: 'Total Time',
                            //     value: shortDateTimeDisplayDelta(
                            //         jobDetail?.updated_at,
                            //         jobDetail?.created_at
                            //     ),
                            // },
                            {
                                label: 'Create With Incidents',
                                value: jobDetail?.with_incidents
                                    ? 'True'
                                    : 'False',
                            },

                            // {
                            //     label: 'Severity Status',
                            //     columnSpan: 2,
                            //     value: (
                            //         <>
                            //             <Flex className="w-full ">
                            //                 <SeverityBar
                            //                     benchmark={
                            //                         detail?.job_details
                            //                             ?.job_score
                            //                     }
                            //                 />
                            //             </Flex>
                            //         </>
                            //     ),
                            // },
                            {
                                label: 'Download Results',
                                value: (
                                    <KButton
                                        iconName="download"
                                        className="mt-2"
                                        variant="inline-icon"
                                        onClick={() => GetFullResults()}
                                        loading={fullLoading}
                                    >
                                        Click here for results in JSON
                                    </KButton>
                                ),
                            },
                        ]}
                    />
                </Modal>
                <Modal
                    visible={runnerTable}
                    size="large"
                    onDismiss={() => setRunnerTable(false)}
                    header={'Execution Detail'}
                >
                    <AppLayout
                        toolsOpen={false}
                        navigationOpen={false}
                        contentType="full-page"
                        // className="w-full"
                        toolsHide={true}
                        navigationHide={true}
                        splitPanelOpen={runnerOpen}
                        onSplitPanelToggle={() => {
                            setRunnerOpen(!runnerOpen)
                        }}
                        splitPanel={
                            // @ts-ignore
                            <SplitPanel
                                // @ts-ignore
                                header={selectedRunner?.runner_id}
                            >
                                <KeyValuePairs
                                    columns={3}
                                    items={[
                                        {
                                            label: 'Compliance Job ID',
                                            value: selectedRunner?.compliance_job_id,
                                        },
                                        {
                                            label: 'Control ID',
                                            value: selectedRunner?.control_id,
                                        },
                                        {
                                            label: 'Integration ID',
                                            value: selectedRunner?.integration_id,
                                        },

                                        {
                                            label: 'Queued At',
                                            value: dateTimeDisplay(
                                                selectedRunner?.queued_at
                                            ),
                                        },
                                        {
                                            label: 'Executed At',
                                            value: dateTimeDisplay(
                                                selectedRunner?.executed_at
                                            ),
                                        },
                                        {
                                            label: 'Completed At',
                                            value: dateTimeDisplay(
                                                selectedRunner?.completed_at
                                            ),
                                        },
                                        {
                                            label: 'Worker Pod Name',
                                            value: selectedRunner?.worker_pod_name,
                                        },
                                        {
                                            label: 'Status',
                                            value: selectedRunner?.status,
                                        },
                                        {
                                            label: 'Failure Message',
                                            value:
                                                selectedRunner?.failure_message ||
                                                'N/A',
                                        }, // Default if empty
                                        {
                                            label: 'Trigger Type',
                                            value: selectedRunner?.trigger_type,
                                        },
                                    ]}
                                />
                            </SplitPanel>
                        }
                        content={
                            <KTable
                                className=""
                                // resizableColumns
                                variant="full-page"
                                // variant="table"
                                renderAriaLive={({
                                    firstIndex,
                                    lastIndex,
                                    totalItemsCount,
                                }) =>
                                    `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                                }
                                onRowClick={(event) => {
                                    const row = event.detail.item
                                    // @ts-ignore
                                    setSelectedRunner(row)
                                    setRunnerOpen(true)
                                }}
                                columnDefinitions={[
                                    {
                                        id: 'id',
                                        header: 'Runner ID',
                                        cell: (item) => item.runner_id,
                                        // sortingField: 'id',
                                        isRowHeader: true,
                                    },

                                    {
                                        id: 'control_id',
                                        header: 'Control ID',
                                        sortingField: 'severity',
                                        cell: (item) => item?.control_id,
                                        maxWidth: 100,
                                    },

                                    {
                                        id: 'integration_id',
                                        header: 'Integration ID',
                                        sortingField: 'oks',

                                        cell: (item) => item?.integration_id,
                                        // minWidth: 50,
                                        maxWidth: 100,
                                    },
                                    {
                                        id: 'queued_at',
                                        header: 'Queued At',
                                        sortingField: 'oks',

                                        cell: (item) =>
                                            dateTimeDisplay(item?.queued_at),
                                        // minWidth: 50,
                                        maxWidth: 100,
                                    },
                                    {
                                        id: 'executed_at',
                                        header: 'Executed At',
                                        sortingField: 'oks',

                                        cell: (item) =>
                                            dateTimeDisplay(item?.executed_at),
                                        // minWidth: 50,
                                        maxWidth: 100,
                                    },
                                    {
                                        id: 'completed_at',
                                        header: 'Completed At',
                                        sortingField: 'oks',

                                        cell: (item) =>
                                            dateTimeDisplay(item?.completed_at),
                                        // minWidth: 50,
                                        maxWidth: 100,
                                    },
                                    {
                                        id: 'wait',
                                        header: 'Wait time',
                                        sortingField: 'oks',

                                        cell: (item) =>
                                            shortDateTimeDisplayDelta(
                                                item?.executed_at,
                                                item?.queued_at
                                            ),
                                        // minWidth: 50,
                                        maxWidth: 100,
                                    },
                                    {
                                        id: 'execution_time',
                                        header: 'Total time',
                                        sortingField: 'oks',

                                        cell: (item) =>
                                            shortDateTimeDisplayDelta(
                                                item?.completed_at,
                                                item?.queued_at
                                            ),
                                        // minWidth: 50,
                                        maxWidth: 100,
                                    },

                                    {
                                        id: 'status',
                                        header: 'status',
                                        sortingField: 'oks',

                                        cell: (item) => item?.status,
                                        // minWidth: 50,
                                        maxWidth: 100,
                                    },
                                    {
                                        id: 'action',
                                        header: '',
                                        cell: (item) => (
                                            // @ts-ignore
                                            <KButton
                                                onClick={() => {
                                                    setSelectedRunner(item)
                                                    setRunnerOpen(true)
                                                }}
                                                className="w-full"
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
                                    { id: 'id', visible: true },
                                    {
                                        id: 'integration_id',
                                        visible: true,
                                    },
                                    {
                                        id: 'control_id',
                                        visible: true,
                                    },
                                    {
                                        id: 'queued_at',
                                        visible: false,
                                    },
                                    {
                                        id: 'executed_at',
                                        visible: false,
                                    },
                                    {
                                        id: 'completed_at',
                                        visible: false,
                                    },
                                    {
                                        id: 'execution_time',
                                        visible: true,
                                    },

                                    {
                                        id: 'status',
                                        visible: true,
                                    },
                                    {
                                        id: 'action',
                                        visible: true,
                                    },
                                ]}
                                enableKeyboardNavigation
                                items={
                                    // @prettie
                                    runners && runners.length > 0
                                        ? runners?.slice(
                                              runnerPage * 20,
                                              (runnerPage + 1) * 20
                                          )
                                        : []
                                }
                                loading={detailLoading}
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
                                header={
                                    <Header
                                        counter={
                                            runners?.length
                                                ? `(${runners?.length})`
                                                : ''
                                        }
                                        actions={
                                            <KButton
                                                onClick={() => {
                                                    GetRunners()
                                                }}
                                                iconName="refresh"
                                            ></KButton>
                                        }
                                        className="w-full"
                                    >
                                        Runners{' '}
                                    </Header>
                                }
                                pagination={
                                    <CustomPagination
                                        currentPageIndex={runnerPage + 1}
                                        pagesCount={Math.ceil(
                                            runners?.length / 20
                                        )}
                                        onChange={({ detail }) =>
                                            setRunnerPage(
                                                detail.currentPageIndex - 1
                                            )
                                        }
                                    />
                                }
                            />
                        }
                    />
                </Modal>
            </div>
        </>
    )
}
