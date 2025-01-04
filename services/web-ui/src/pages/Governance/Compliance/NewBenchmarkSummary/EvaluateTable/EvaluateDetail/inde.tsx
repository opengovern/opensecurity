// @ts-nocheck
import { useAtomValue } from 'jotai'
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
import { isDemoAtom } from '../../../../../store'
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
import { dateTimeDisplay } from '../../../../../../utilities/dateDisplay'
import StatusIndicator from '@cloudscape-design/components/status-indicator'
import SeverityBar from './SeverityBar'
import { useParams } from 'react-router-dom'
import { RunDetail } from './types'
import TopHeader from '../../../../../../components/Layout/Header'

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
    const [fullLoading, setFullLoading] = useState(false)
    const [sortOrder, setSortOrder] = useState('asc')
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
    const [integrationOpen,setIntegrationOpen]= useState(false)
    const [integrationDetail,setIntegrationDetail]= useState()
    const [integrationDetailLoading,setIntegrationDetailLoading]= useState(false)
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
                    res?.data?.controls[selectedControl.title]?.results?.alarm
                alarms.map((alarm) => {
                    temp.push({
                        resource_id: alarm.resource_id,
                        resource_type: alarm.resource_type,
                        reason: alarm.reason,
                    })
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
        const body ={
            integration_id:integrations
        }
         axios
             .post(
                 // @ts-ignore
                 `${url}/main/integration/api/v1/integrations/list`,body,
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

    useEffect(() => {
        GetDetail()
        GetControls()
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
    const getRows = () =>{
        if(sortOrder == 'asc'){
            return runDetail?.sort(sortField.sortingComparator).slice(page * 10, (page + 1) * 10)
        }
        return runDetail
            ?.sort(sortField.sortingComparator)
            .reverse()
            .slice(page * 10, (page + 1) * 10)

    }
    return (
        <>
            {/* <TopHeader /> */}
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
            <Flex
                className="w-full bg-white p-4 rounded-lg mt-4"
                flexDirection="col"
                alignItems="start"
            >
                <Flex
                    flexDirection="col"
                    className="w-full mt-4"
                    alignItems="center"
                    justifyContent="center"
                >
                    <KeyValuePairs
                        className="w-full"
                        columns={6}
                        items={[
                            {
                                label: 'Job ID',
                                value: id,
                            },
                            {
                                label: 'Benchmark ID',
                                value: detail?.job_details?.framework?.id,
                            },
                            {
                                label: 'Benchmark Title',
                                value: detail?.job_details?.framework?.title,
                            },
                            {
                                label: 'Control score',
                                value:
                                    (
                                        1 -
                                        detail?.job_details?.job_score
                                            ?.failed_controls /
                                            detail?.job_details?.job_score
                                                ?.total_controls
                                    )?.toFixed(2) *
                                        100 +
                                    '%',
                            },
                            {
                                label: 'Integrations',
                                value: (
                                    <>
                                        <Flex
                                            className="gap-2 mt-2"
                                            justifyContent="start"
                                        >
                                            {integrations?.length}
                                            <KButton
                                                onClick={() => {
                                                    if (integrations) {
                                                        GetIntegrationDetail()
                                                        setIntegrationOpen(true)
                                                    }
                                                }}
                                                variant="inline-link"
                                                ariaLabel={`Open Detail`}
                                            >
                                                See details
                                            </KButton>
                                        </Flex>
                                    </>
                                ),
                            },
                            {
                                label: 'Last Evaulated at',
                                value: (
                                    <>
                                        {dateTimeDisplay(
                                            detail?.job_details?.updated_at
                                        )}
                                    </>
                                ),
                            },
                            {
                                label: 'Severity Status',
                                columnSpan: 2,
                                value: (
                                    <>
                                        <Flex className="w-full ">
                                            <SeverityBar
                                                benchmark={
                                                    detail?.job_details
                                                        ?.job_score
                                                }
                                            />
                                        </Flex>
                                    </>
                                ),
                            },
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
                            // {
                            //     label: 'Job Status',
                            //     value: (
                            //         <StatusIndicator
                            //             type={JOB_STATUS[detail?.job_status]}
                            //         >
                            //             {detail?.job_status}
                            //         </StatusIndicator>
                            //     ),
                            // },
                        ]}
                    />
                    {/* <Flex className="w-1/2 mt-2">
                    <SeverityBar benchmark={detail} />
                </Flex> */}
                </Flex>
            </Flex>
            <Flex className="w-100 bg-white p-4 rounded-lg mt-4">
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
                    sortingDescending={sortOrder == 'desc' ? true : false}
                    sortingColumn={sortField}
                    onSortingChange={({ detail }) => {
                        const desc = detail.isDescending ? 'desc' : 'asc'
                        setSortOrder(desc)
                        setSortField(detail.sortingColumn)
                    }}
                    onRowClick={(event) => {
                        const row = event.detail.item
                        // @ts-ignore
                        setSelectedControl(row)
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
                                    {item.severity.charAt(0).toUpperCase() +
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
                                    }}
                                    variant="inline-link"
                                    ariaLabel={`Open Detail`}
                                >
                                    See details
                                </KButton>
                            ),
                        },
                    ]}
                    columnDisplay={[
                        { id: 'id', visible: true },
                        { id: 'title', visible: false },
                        { id: 'connector', visible: false },
                        { id: 'query', visible: false },
                        { id: 'severity', visible: true },
                        { id: 'incidents', visible: true },
                        { id: 'passing_resources', visible: true },
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
                        <Header className="w-full">
                            Controls{' '}
                            <span className=" font-medium">
                                ({runDetail?.length})
                            </span>
                        </Header>
                    }
                    pagination={
                        <Pagination
                            currentPageIndex={page + 1}
                            pagesCount={Math.ceil(runDetail?.length / 10)}
                            onChange={({ detail }) =>
                                setPage(detail.currentPageIndex - 1)
                            }
                        />
                    }
                />
            </Flex>
            <Modal
                visible={open}
                size="max"
                onDismiss={() => setOpen(false)}
                header={`Alarms`}
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
                    ]}
                    columnDisplay={[
                        { id: 'resource_id', visible: true },
                        { id: 'resource_type', visible: true },
                        { id: 'reason', visible: true },
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
                        <Pagination
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
                                            item.state === 'active'
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
                    // pagination={
                    //     <Pagination
                    //         currentPageIndex={resourcePage + 1}
                    //         pagesCount={Math.ceil(resources?.length / 7)}
                    //         onChange={({ detail }) =>
                    //             setResourcePage(detail.currentPageIndex - 1)
                    //         }
                    //     />
                    // }
                />
            </Modal>
        </>
    )
}
