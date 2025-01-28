// @ts-nocheck
import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Badge,
    Button,
    Card,
    Color,
    Divider,
    Flex,
    Text,
    Title,
} from '@tremor/react'

import { Radio } from 'pretty-checkbox-react'
import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import Table, { IColumn } from '../../../../components/Table'
import {
    Api,
    PlatformEnginePkgDescribeApiJob,
} from '../../../../api/api'
import AxiosAPI from '../../../../api/ApiConfig'
import { useScheduleApiV1JobsCreate } from '../../../../api/schedule.gen'
import DrawerPanel from '../../../../components/DrawerPanel'
import KFilter from '../../../../components/Filter'
import { CloudIcon } from '@heroicons/react/24/outline'
import { string } from 'prop-types'
import KTable from '@cloudscape-design/components/table'
import Box from '@cloudscape-design/components/box'
import SpaceBetween from '@cloudscape-design/components/space-between'
import KBadge from '@cloudscape-design/components/badge'
import {
    BreadcrumbGroup,
    DateRangePicker,
    Header,
    Link,
    Pagination,
    PropertyFilter,
    Select,
} from '@cloudscape-design/components'
import {
    AppLayout,
    Container,
    ContentLayout,
    SplitPanel,
} from '@cloudscape-design/components'
import KButton from '@cloudscape-design/components/button'
import KeyValuePairs from '@cloudscape-design/components/key-value-pairs'
import axios from 'axios'
import { title } from 'process'


interface Option {
    label: string | undefined
    value: string | undefined
}
interface IntegrationListProps {
    name?: string
    integration_type?: string
}

export default function DiscoveryJobs({
    name,
    integration_type,
}: IntegrationListProps) {
    const findParmas = (key: string): string[] => {
        const params = searchParams.getAll(key)
        const temp = []
        if (params) {
            params.map((item, index) => {
                temp.push(item)
            })
        }
        return temp
    }
    const [open, setOpen] = useState(false)
    const [queries, setQueries] = useState({
        tokens: [],
        operation: 'and',
    })

    const [clickedJob, setClickedJob] =
        useState<PlatformEnginePkgDescribeApiJob>()
    const [searchParams, setSearchParams] = useSearchParams()
   
    const [statusFilter, setStatusFilter] = useState<string[] | undefined>(
        findParmas('status')
    )
    const [allStatuses, setAllStatuses] = useState<Option[]>([])
    const [loading, setLoading] = useState(false)
    const [jobs, setJobs] = useState([])
    const [page, setPage] = useState(1)
    const [sort, setSort] = useState('updatedAt')
    const [sortOrder, setSortOrder] = useState(true)

    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [date, setDate] = useState({
        key: 'previous-6-hours',
        amount: 6,
        unit: 'hour',
        type: 'relative',
    })
    const [filter, setFilter] = useState()
    useEffect(() => {
        // @ts-ignore
        if (filter) {
            // @ts-ignore

            if (filter.value == '1') {
                setDate({
                    key: 'previous-7-days',
                    amount: 7,
                    unit: 'day',
                    type: 'relative',
                })
                setQueries({
                    tokens: [
                        {
                            propertyKey: 'job_status',
                            value: 'FAILED',
                            operator: '=',
                        },
                    ],
                    operation: 'and',
                })
            }
        }
    }, [filter])

    const arrayToString = (arr: string[], title: string) => {
        let temp = ``
        arr.map((item, index) => {
            if (index == 0) {
                temp += arr[index]
            } else {
                temp += `&${title}=${arr[index]}`
            }
        })
        console.log(temp)
        return temp
    }

    useEffect(() => {
        if (
          
            searchParams.get('status') !== statusFilter
        ) {
           
            if (statusFilter?.length != 0) {
                searchParams.set('status', statusFilter)
            } else {
                searchParams.delete('status')
            }
            window.history.pushState({}, '', `?${searchParams.toString()}`)
        }
    }, [ statusFilter])

    const GetRows = () => {
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
        const status_filter = []
        queries.tokens.map((item) => {
            if (item.propertyKey == 'job_status') {
                status_filter.push(item.value)
            }
        })
        let body = {
            cursor: page ,
            per_page: 15,
            integration_info: [
                {
                    integration_type: integration_type,
                },
            ],
            job_status: status_filter,

            sort_by: sort,
            sortOrder: sortOrder ? 'DESC' : 'ASC',
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
            .post(`${url}/main/schedule/api/v3/jobs/discovery`, body, config)
            .then((resp) => {

                if (resp.data.items) {
                    setJobs(resp.data.items)
                } else {
                    setJobs([])
                }

                setTotalCount(
                    resp?.data?.total_count
                )
                setTotalPage(Math.ceil(resp?.data?.total_count / 15))
                setLoading(false)

                // params.success({
                //     rowData: resp.data.jobs || [],
                //     rowCount: resp.data.summaries
                //         ?.map((v) => v.count)
                //         .reduce((prev, curr) => (prev || 0) + (curr || 0), 0),
                // })
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)

                // params.fail()
            })
    }

    useEffect(() => {
        GetRows()
    }, [queries, date, page, sort, sortOrder])

    const clickedJobDetails = [
        { title: 'ID', value: clickedJob?.job_id },
        { title: 'Title', value: clickedJob?.title },
        { title: 'Type', value: clickedJob?.type },
        { title: 'Created At', value: clickedJob?.created_at },
        { title: 'Updated At', value: clickedJob?.updated_at },

        { title: 'Status', value: clickedJob?.job_status },
        { title: 'Failure Reason', value: clickedJob?.failure_message },
       
    ]

    return (
        <>
            <AppLayout
                toolsOpen={false}
                navigationOpen={false}
                contentType="table"
                toolsHide={true}
                navigationHide={true}
                splitPanelOpen={open}
                onSplitPanelToggle={() => {
                    setOpen(!open)
                }}
                splitPanel={
                    <SplitPanel
                        header={
                            clickedJob ? clickedJob.title : 'Job not selected'
                        }
                    >
                        <Flex
                            flexDirection="col"
                            className="w-full"
                            alignItems="center"
                            justifyContent="center"
                        >
                            <KeyValuePairs
                                columns={4}
                                className="w-full"
                                items={clickedJobDetails.map((item) => {
                                    return {
                                        label: item?.title,
                                        value: item?.value,
                                    }
                                })}
                            />
                            {clickedJob?.parameters &&
                                Object.entries(clickedJob?.parameters).length >
                                    0 && (
                                    <>
                                        <Flex
                                            flexDirection="col"
                                            className="w-full gap-2 justify-start items-start mt-2"
                                        >
                                            <Title className=" font-semibold font-sans text-lg">
                                                Parameters:
                                            </Title>
                                            <KeyValuePairs
                                                columns={4}
                                                className="w-full"
                                                items={Object.entries(
                                                    clickedJob?.parameters
                                                )?.map((item) => {
                                                    return {
                                                        label: item[0],
                                                        value: item[1],
                                                    }
                                                })}
                                            />
                                        </Flex>
                                    </>
                                )}
                        </Flex>
                    </SplitPanel>
                }
                content={
                    <KTable
                        className="  min-h-[450px]"
                        variant="full-page"
                        // resizableColumns
                        renderAriaLive={({
                            firstIndex,
                            lastIndex,
                            totalItemsCount,
                        }) =>
                            `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                        }
                        onSortingChange={(event) => {
                            setSort(event.detail.sortingColumn.sortingField)
                            setSortOrder(!sortOrder)
                        }}
                        sortingColumn={sort}
                        sortingDescending={sortOrder}
                        // @ts-ignore
                        onRowClick={(event) => {
                            const row = event.detail.item
                            setClickedJob(row)
                            setOpen(true)
                        }}
                        columnDefinitions={[
                            {
                                id: 'id',
                                header: 'Id',
                                cell: (item) => <>{item.job_id}</>,
                                // sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 100,
                            },
                            {
                                id: 'createdAt',
                                header: 'Created At',
                                cell: (item) => (
                                    <>{`${item?.created_at.split('T')[0]} ${
                                        item?.created_at
                                            .split('T')[1]
                                            .split('.')[0]
                                    } `}</>
                                ),
                                sortingField: 'createdAt',
                                isRowHeader: true,
                                maxWidth: 100,
                            },
                            {
                                id: 'type',
                                header: 'Job Type',
                                cell: (item) => <>{item.type}</>,
                                sortingField: 'type',
                                isRowHeader: true,
                                maxWidth: 100,
                            },
                            {
                                id: 'title',
                                header: 'Title',
                                cell: (item) => <>{item.title}</>,
                                sortingField: 'title',
                                isRowHeader: true,
                                maxWidth: 150,
                            },
                            {
                                id: 'status',
                                header: 'Status',
                                cell: (item) => {
                                    let jobStatus = ''
                                    let jobColor: Color = 'gray'
                                    switch (item?.job_status) {
                                        case 'CREATED':
                                            jobStatus = 'created'
                                            break
                                        case 'QUEUED':
                                            jobStatus = 'queued'
                                            break
                                        case 'IN_PROGRESS':
                                            jobStatus = 'in progress'
                                            jobColor = 'orange'
                                            break
                                        case 'RUNNERS_IN_PROGRESS':
                                            jobStatus = 'in progress'
                                            jobColor = 'orange'
                                            break
                                        case 'SUMMARIZER_IN_PROGRESS':
                                            jobStatus = 'summarizing'
                                            jobColor = 'orange'
                                            break
                                        case 'SINK_IN_PROGRESS':
                                            jobStatus = 'sinking'
                                            jobColor = 'orange'
                                            break
                                        case 'OLD_RESOURCE_DELETION':
                                            jobStatus = 'summarizing'
                                            jobColor = 'orange'
                                            break
                                        case 'SUCCEEDED':
                                            jobStatus = 'succeeded'
                                            jobColor = 'emerald'
                                            break
                                        case 'COMPLETED':
                                            jobStatus = 'completed'
                                            jobColor = 'emerald'
                                            break
                                        case 'FAILED':
                                            jobStatus = 'failed'
                                            jobColor = 'red'
                                            break
                                        case 'COMPLETED_WITH_FAILURE':
                                            jobStatus = 'completed with failed'
                                            jobColor = 'red'
                                            break
                                        case 'TIMEOUT':
                                            jobStatus = 'time out'
                                            jobColor = 'red'
                                            break
                                        default:
                                            jobStatus = String(item?.status)
                                    }

                                    return (
                                        <Badge color={jobColor}>
                                            {jobStatus}
                                        </Badge>
                                    )
                                },
                                sortingField: 'status',
                                isRowHeader: true,
                                maxWidth: 50,
                            },
                            {
                                id: 'updatedAt',
                                header: 'Updated At',
                                cell: (item) => (
                                    <>{`${item?.updated_at.split('T')[0]} ${
                                        item?.updated_at
                                            .split('T')[1]
                                            .split('.')[0]
                                    } `}</>
                                ),
                                sortingField: 'updatedAt',
                                isRowHeader: true,
                                maxWidth: 100,
                            },
                        ]}
                        columnDisplay={[
                            { id: 'id', visible: true },
                            { id: 'title', visible: true },
                            { id: 'type', visible: false },
                            { id: 'status', visible: true },
                            { id: 'createdAt', visible: true },
                            { id: 'updatedAt', visible: true },
                        ]}
                        enableKeyboardNavigation
                        // @ts-ignore
                        items={jobs}
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
                                <Select
                                    // @ts-ignore
                                    selectedOption={filter}
                                    className="sm:w-1/5 w-full mt-[-9px]"
                                    inlineLabelText={'Saved Filters'}
                                    placeholder="Select Filter Set"
                                    // @ts-ignore
                                    onChange={({ detail }) =>
                                        // @ts-ignore
                                        setFilter(detail.selectedOption)
                                    }
                                    options={[
                                        {
                                            label: 'All Discovery Jobs',
                                            value: '1',
                                        },
                                    ]}
                                />
                                <PropertyFilter
                                    // @ts-ignore
                                    query={queries}
                                    // @ts-ignore
                                    onChange={({ detail }) => {
                                        // @ts-ignore
                                        setQueries(detail)
                                    }}
                                    // countText="5 matches"
                                    // enableTokenGroups
                                    expandToViewport
                                    filteringAriaLabel="Job Filters"
                                    // @ts-ignore
                                    // filteringOptions={filters}
                                    filteringPlaceholder="Job Filters"
                                    // @ts-ignore
                                    filteringOptions={[
                                        {
                                            propertyKey: 'job_status',
                                            value: 'SUCCEEDED',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'FAILED',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'IN_PROGRESS',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'TIMEOUT',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'SUMMARIZER_IN_PROGRESS',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'RUNNERS_IN_PROGRESS',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'SINK_IN_PROGRESS',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'QUEUED',
                                        },
                                        {
                                            propertyKey: 'job_status',
                                            value: 'CREATED',
                                        },
                                    ]}
                                    // @ts-ignore

                                    filteringProperties={[
                                        {
                                            key: 'job_status',
                                            operators: ['='],
                                            propertyLabel: 'Job Status',
                                            groupValuesLabel:
                                                'Job Status values',
                                        },
                                    ]}
                                    // filteringProperties={
                                    //     filterOption
                                    // }
                                />

                                <DateRangePicker
                                    onChange={({ detail }) =>
                                        setDate(detail.value)
                                    }
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
                                    hideTimeOffset
                                    absoluteFormat="long-localized"
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
                                    placeholder="Filter by a date and time range"
                                />
                            </Flex>
                        }
                        header={
                            <Header
                                counter={totalCount ? `(${totalCount})` : ''}
                                actions={
                                    <Flex className="flex-row gap-2">
                                        <KButton
                                            onClick={() => {
                                                GetRows()
                                            }}
                                            iconName="refresh"
                                        ></KButton>
                                       
                                    </Flex>
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
                }
            />
        </>
    )
}

