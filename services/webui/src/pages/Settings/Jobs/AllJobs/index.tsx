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


const jobTypes = [
    {
        label: 'Discovery',
        value: 'discovery',
    },
    {
        label: 'Compliance',
        value: 'compliance',
    },
    {
        label: 'Analytics',
        value: 'analytics',
    },
]
const ShowHours = [
    {
        label: '1h',
        value: '1',
    },
    {
        label: '3h',
        value: '3',
    },
    {
        label: '6h',
        value: '6',
    },
    {
        label: '24h',
        value: '24',
    },
    // {
    //     label: 'all',
    //     value: 'all',
    // },
]
interface Option {
    label: string | undefined
    value: string | undefined
}
export default function SettingsALLJobs() {
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
    const [jobTypeFilter, setJobTypeFilter] = useState<string[] | undefined>(
        findParmas('type')
    )
    const [statusFilter, setStatusFilter] = useState<string[] | undefined>(
        findParmas('status')
    )
    const [allStatuses, setAllStatuses] = useState<Option[]>([])
    const [loading, setLoading] = useState(false)
    const [jobs, setJobs] = useState([])
    const [page, setPage] = useState(0)
    const [sort, setSort] = useState('updatedAt')
    const [sortOrder, setSortOrder] = useState(true)

    const [totalCount, setTotalCount] = useState(0)
    const [totalPage, setTotalPage] = useState(0)
    const [propertyOptions, setPropertyOptions] = useState()
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
                    key: 'previous-6-hours',
                    amount: 6,
                    unit: 'hour',
                    type: 'relative',
                })
                setQueries({
                    tokens: [
                        {
                            propertyKey: 'job_type',
                            value: 'compliance',
                            operator: '=',
                        },
                    ],
                    operation: 'and',
                })
            }
            // @ts-ignore
            else if (filter.value == '2') {
                setDate({
                    key: 'previous-6-hours',
                    amount: 6,
                    unit: 'hour',
                    type: 'relative',
                })
                setQueries({
                    tokens: [
                        {
                            propertyKey: 'job_type',
                            value: 'compliance',
                            operator: '=',
                        },
                        {
                            propertyKey: 'job_status',
                            value: 'FAILED',
                            operator: '=',
                        },
                    ],
                    operation: 'and',
                })
            } else if (filter.value == '3') {
                setDate({
                    key: 'previous-7-days',
                    amount: 7,
                    unit: 'day',
                    type: 'relative',
                })
                setQueries({
                    tokens: [
                        {
                            propertyKey: 'job_type',
                            value: 'discovery',
                            operator: '=',
                        },
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
            searchParams.getAll('type') !== jobTypeFilter ||
            searchParams.get('status') !== statusFilter
        ) {
            if (jobTypeFilter?.length != 0) {
                searchParams.set('type', jobTypeFilter)
            } else {
                searchParams.delete('type')
            }
            if (statusFilter?.length != 0) {
                searchParams.set('status', statusFilter)
            } else {
                searchParams.delete('status')
            }
            window.history.pushState({}, '', `?${searchParams.toString()}`)
        }
    }, [jobTypeFilter, statusFilter])
    
    const GetRows = () => {
        setLoading(true)
        const api = new Api()
        api.instance = AxiosAPI
        const status_filter = []
        const jobType_filter = []
        queries.tokens.map((item) => {
            if (item.propertyKey == 'job_status') {
                status_filter.push(item.value)
            } else if (item.propertyKey == 'job_type') {
                jobType_filter.push(item.value)
            }
        })
        let body = {
            pageStart: page * 15,
            pageEnd: (page + 1) * 15,

            statusFilter: status_filter,
            typeFilters: jobType_filter,
            sortBy: sort,
            sortOrder: sortOrder ? 'DESC' : 'ASC',
        }
        if (date) {
            if (date.type == 'relative') {
                body.interval = `${date.amount} ${date.unit}s`
            } else {
                body.from = date.startDate
                body.to = date.endDate
            }
        }

        api.schedule
            .apiV1JobsCreate(body)
            .then((resp) => {
                const response = resp.data
                const temp =
                    response?.summaries
                        ?.map((v) => {
                            return { label: v.status, value: v.status }
                        })
                        .filter(
                            (thing, i, arr) =>
                                arr.findIndex(
                                    (t) => t.label === thing.label
                                ) === i
                        ) || []
                setAllStatuses(temp)
                const temp_option = []
                temp.map((item) => {
                    temp_option.push({
                        propertyKey: 'job_status',
                        value: item.value,
                    })
                })
                jobTypes?.map((item) => {
                    temp_option.push({
                        propertyKey: 'job_type',
                        value: item.value,
                    })
                })
                setPropertyOptions(temp_option)

                if (resp.data.jobs) {
                    setJobs(resp.data.jobs)
                } else {
                    setJobs([])
                }

                setTotalCount(
                    resp.data.summaries
                        ?.map((v) => v.count)
                        ?.reduce((prev, curr) => (prev || 0) + (curr || 0), 0)
                )
                setTotalPage(
                    Math.ceil(
                        resp.data.summaries
                            ?.map((v) => v.count)
                            ?.reduce(
                                (prev, curr) => (prev || 0) + (curr || 0),
                                0
                            ) / 15
                    )
                )
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
        { title: 'ID', value: clickedJob?.id },
        { title: 'Title', value: clickedJob?.title },
        { title: 'Type', value: clickedJob?.type },
        { title: 'Created At', value: clickedJob?.createdAt },
        { title: 'Updated At', value: clickedJob?.updatedAt },
        // {
        //     title: 'OpenGovernance Connection ID',
        //     value: clickedJob?.connectionID,
        // },
        // { title: 'Account ID', value: clickedJob?.connectionProviderID },
        // { title: 'Account Name', value: clickedJob?.connectionProviderName },
        { title: 'Status', value: clickedJob?.status },
        { title: 'Failure Reason', value: clickedJob?.failureReason },
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
                                cell: (item) => <>{item.id}</>,
                                // sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 100,
                            },
                            {
                                id: 'createdAt',
                                header: 'Created At',
                                cell: (item) => (
                                    <>{`${item?.createdAt.split('T')[0]} ${
                                        item?.createdAt
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
                                maxWidth: 100,
                            },
                            {
                                id: 'status',
                                header: 'Status',
                                cell: (item) => {
                                    let jobStatus = ''
                                    let jobColor: Color = 'gray'
                                    switch (item?.status) {
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
                                maxWidth: 100,
                            },
                            {
                                id: 'updatedAt',
                                header: 'Updated At',
                                cell: (item) => (
                                    <>{`${item?.updatedAt.split('T')[0]} ${
                                        item?.updatedAt
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
                            { id: 'updatedAt', visible: true },
                            { id: 'title', visible: true },
                            { id: 'type', visible: true },
                            { id: 'status', visible: true },
                            { id: 'createdAt', visible: true },
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
                                className="gap-2"
                            >
                                <Select
                                    // @ts-ignore
                                    selectedOption={filter}
                                    className="w-1/5 mt-[-9px]"
                                    inlineLabelText={'Saved Filters'}
                                    placeholder="Select Filter Set"
                                    // @ts-ignore
                                    onChange={({ detail }) =>
                                        // @ts-ignore
                                        setFilter(detail.selectedOption)
                                    }
                                    options={[
                                        {
                                            label: 'Recent Compliance Jobs',
                                            value: '1',
                                        },
                                        {
                                            label: 'Failing Compliance Jobs',
                                            value: '2',
                                        },
                                        {
                                            label: 'All Discovery Jobs',
                                            value: '3',
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
                                    filteringOptions={propertyOptions}
                                    // @ts-ignore

                                    filteringProperties={[
                                        {
                                            key: 'job_type',
                                            operators: ['='],
                                            propertyLabel: 'Job Type',
                                            groupValuesLabel: 'Job Type values',
                                        },
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
                                    <KButton
                                        onClick={() => {
                                            GetRows()
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
                                currentPageIndex={page + 1}
                                pagesCount={totalPage}
                                onChange={({ detail }) =>
                                    setPage(detail.currentPageIndex - 1)
                                }
                            />
                        }
                    />
                }
            />
        </>
    )
}

{
    /**
       <Flex flexDirection="col">
                <Flex
                    flexDirection="row"
                    alignItems="start"
                    justifyContent="between"
                    className="gap-2"
                >
                    <Flex
                        flexDirection="row"
                        alignItems="start"
                        justifyContent="start"
                        className="gap-2"
                    >
                        <KFilter
                            options={jobTypes}
                            type="multi"
                            hasCondition={true}
                            condition={jobTypeContains}
                            selectedItems={jobTypeFilter}
                            onChange={(values: string[]) => {
                                console.log(values, 'values')
                                console.log(jobTypeFilter, 'filter')

                                setJobTypeFilter(values)
                            }}
                            label="Job Types"
                            icon={CloudIcon}
                        />
                        <KFilter
                            options={allStatuses}
                            type="multi"
                            selectedItems={statusFilter}
                            condition={jobStatusContains}
                            hasCondition={true}
                            onChange={(values: string[]) => {
                                setStatusFilter(values)
                            }}
                            label="Job Status"
                            icon={CloudIcon}
                        />
                        <KFilter
                            options={ShowHours}
                            type="multi"
                            hasCondition={false}
                            selectedItems={showHoursFilter}
                            onChange={(values: string[]) => {
                                if (values.length == 0) {
                                    setShowHourFilter([])
                                } else {
                                    setShowHourFilter([values.pop()])
                                }
                            }}
                            label="Show jobs in"
                            icon={CloudIcon}
                        />
                    </Flex>

                    <Button
                        onClick={() => {
                            // @ts-ignore
                            ssr()
                        }}
                        disabled={false}
                        loading={false}
                        loadingText="Running"
                    >
                        Refresh
                    </Button>
                </Flex>
                <Card className="mt-4">
                    <Title className="font-semibold mb-5">Jobs</Title>

                    <Flex alignItems="start">
                        {/* <Card className="sticky top-6 min-w-[200px] max-w-[200px]">
                    <Accordion
                        defaultOpen
                        className="border-0 rounded-none bg-transparent mb-1"
                    >
                        <AccordionHeader className="pl-0 pr-0.5 py-1 w-full bg-transparent">
                            <Text className="font-semibold text-gray-800">
                                Job Type
                            </Text>
                        </AccordionHeader>
                        <AccordionBody className="pt-3 pb-1 px-0.5 w-full cursor-default bg-transparent">
                            <Flex
                                flexDirection="col"
                                alignItems="start"
                                className="gap-1.5"
                            >
                                {jobTypes.map((jobType) => (
                                    <Radio
                                        name="jobType"
                                        onClick={() =>
                                            setJobTypeFilter(jobType.value)
                                        }
                                        checked={
                                            jobTypeFilter === jobType.value
                                        }
                                    >
                                        {jobType.label}
                                    </Radio>
                                ))}
                            </Flex>
                        </AccordionBody>
                    </Accordion>
                    <Divider className="my-3" />
                    <Accordion
                        defaultOpen
                        className="border-0 rounded-none bg-transparent mb-1"
                    >
                        <AccordionHeader className="pl-0 pr-0.5 py-1 w-full bg-transparent">
                            <Text className="font-semibold text-gray-800">
                                Status
                            </Text>
                        </AccordionHeader>
                        <AccordionBody className="pt-3 pb-1 px-0.5 w-full cursor-default bg-transparent">
                            <Flex
                                flexDirection="col"
                                alignItems="start"
                                className="gap-1.5"
                            >
                                <Radio
                                    name="status"
                                    onClick={() => setStatusFilter('')}
                                    checked={statusFilter === ''}
                                >
                                    All
                                </Radio>
                                {allStatuses.map((status) => (
                                    <Radio
                                        name="status"
                                        onClick={() => setStatusFilter(status)}
                                        checked={statusFilter === status}
                                    >
                                        {status}
                                    </Radio>
                                ))}
                            </Flex>
                        </AccordionBody>
                    </Accordion>
                </Card> 
                        <Flex className="pl-4">
                            <Table
                                id="jobs"
                                columns={columns()}
                                serverSideDatasource={serverSideRows}
                                onCellClicked={(event) => {
                                    setClickedJob(event.data)
                                    setOpen(true)
                                }}
                                options={{
                                    rowModelType: 'serverSide',
                                    serverSideDatasource: serverSideRows,
                                }}
                            />
                        </Flex>
                    </Flex>
                    <DrawerPanel
                        open={open}
                        onClose={() => setOpen(false)}
                        title="Job Details"
                    >
                        <Flex flexDirection="col">
                            {clickedJobDetails.map((item) => {
                                return (
                                    <Flex
                                        flexDirection="row"
                                        justifyContent="between"
                                        alignItems="start"
                                        className="mt-2"
                                    >
                                        <Text className="w-56 font-bold">
                                            {item.title}
                                        </Text>
                                        <Text className="w-full">
                                            {item.value}
                                        </Text>
                                    </Flex>
                                )
                            })}
                        </Flex>
                    </DrawerPanel>
                </Card>
            </Flex>
    */
}
