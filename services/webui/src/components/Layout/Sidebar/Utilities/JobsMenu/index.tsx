import { Fragment, useEffect, useState } from 'react'
import { Popover, Transition } from '@headlessui/react'
import {
    ArrowPathRoundedSquareIcon,
    CheckIcon,
    ChevronDownIcon,
    ChevronRightIcon,
    ChevronUpIcon,
    ClipboardDocumentListIcon,
    ExclamationTriangleIcon,
} from '@heroicons/react/24/outline'
import { BarList, Button, Card, Color, Flex, Text, Title } from '@tremor/react'
import { useNavigate, useParams } from 'react-router-dom'
import {
    PlatformEnginePkgDescribeApiJobSummary,
    PlatformEnginePkgDescribeApiJobType,
} from '../../../../../api/api'
import { numberDisplay } from '../../../../../utilities/numericDisplay'
import { useScheduleApiV1JobsCreate } from '../../../../../api/schedule.gen'
import Spinner from '../../../../Spinner'
import { useAuth } from '../../../../../utilities/auth'

interface IJobCategoryItem {
    title: string
    jobType: string
    summaries?: PlatformEnginePkgDescribeApiJobSummary[]
}

const inProgressStatuses = [
    'CREATED',
    'QUEUED',
    'IN_PROGRESS',
    'RUNNERS_IN_PROGRESS',
    'SUMMARIZER_IN_PROGRESS',
    'OLD_RESOURCE_DELETION',
]
const failedStatuses = ['FAILED', 'TIMEOUT', 'COMPLETED_WITH_FAILURE']
const succeededStatuses = ['SUCCEEDED', 'COMPLETED']

const checkStatus = (v: string, arr: string[]) => {
    let exists = false
    arr.forEach((element) => {
        if (element === v) {
            exists = true
        }
    })
    return exists
}

function JobCategoryItem({ title, jobType, summaries }: IJobCategoryItem) {
    const navigate = useNavigate()
    const url = window.location.pathname.split('/')
    if (url[1] === 'ws') {
        url.shift()
    }

    const workspace = url[1]

    const result = () => {
        const inProgressJobs =
            summaries?.filter((job) =>
                checkStatus(job.status || '', inProgressStatuses)
            ) || []

        const failedJobs =
            summaries?.filter((job) =>
                checkStatus(job.status || '', failedStatuses)
            ) || []

        const succeededJobs =
            summaries?.filter((job) =>
                checkStatus(job.status || '', succeededStatuses)
            ) || []

        const unknownJobs =
            summaries?.filter(
                (job) =>
                    !checkStatus(job.status || '', succeededStatuses) &&
                    !checkStatus(job.status || '', inProgressStatuses) &&
                    !checkStatus(job.status || '', failedStatuses)
            ) || []

        if (inProgressJobs.length > 0) {
            const color: Color = 'neutral'

            return {
                status: 'In Progress',
                count:
                    inProgressJobs
                        ?.map((v) => v.count)
                        .reduce((prev, curr) => (prev || 0) + (curr || 0), 0) ||
                    0,
                icon: ArrowPathRoundedSquareIcon,
                color,
            }
        }
        if (failedJobs.length > 0) {
            const color: Color = 'rose'
            return {
                status: 'Failed',
                count:
                    failedJobs
                        ?.map((v) => v.count)
                        .reduce((prev, curr) => (prev || 0) + (curr || 0), 0) ||
                    0,
                icon: ExclamationTriangleIcon,
                color,
            }
        }
        if (unknownJobs.length > 0) {
            const color: Color = 'neutral'
            return {
                status: 'Unknown',
                count:
                    unknownJobs
                        ?.map((v) => v.count)
                        .reduce((prev, curr) => (prev || 0) + (curr || 0), 0) ||
                    0,
                icon: ArrowPathRoundedSquareIcon,
                color,
            }
        }
        const color: Color = 'emerald'
        return {
            status: 'Succeeded',
            count:
                (succeededJobs?.length || 0) === 0
                    ? 0
                    : succeededJobs
                          ?.map((v) => v.count)
                          .reduce(
                              (prev, curr) => (prev || 0) + (curr || 0),
                              0
                          ) || 0,
            icon: CheckIcon,
            color,
        }
    }

    const { status, count, icon, color } = result()

    const totalJobs =
        summaries
            ?.map((v) => v.count)
            .reduce((prev, current) => {
                return (prev || 0) + (current || 0)
            }, 0) || 0

    const data =
        summaries?.map((v) => {
            return { name: String(v.status), value: v.count || 0 }
        }) || []

    const fullTitle = () => {
        const percentage = numberDisplay((count / totalJobs) * 100)
        if (status === 'Failed') {
            return `${percentage}% of ${title} activities failed`
        }
        if (status === 'In Progress') {
            return `${percentage}% of ${title} activities are running`
        }
        return `${title} activities finished successfully`
    }
    const [open, setOpen] = useState(false)

    return (
        <Card
            decoration="left"
            decorationColor={color}
            key={title}
            className="h-fit w-96 m-2 p-3 px-3"
        >
            <Flex
                justifyContent="between"
                className="cursor-pointer"
                onClick={() => setOpen(!open)}
            >
                <Flex justifyContent="start" className="space-x-4">
                    <Title className="!text-base truncate">{fullTitle()}</Title>
                </Flex>
                {open ? (
                    <ChevronUpIcon height={20} color="text-blue-500" />
                ) : (
                    <ChevronDownIcon height={20} color="text-blue-500" />
                )}
            </Flex>
            {open && (
                <BarList
                    key={title}
                    data={data}
                    className="mt-2 cursor-pointer"
                    color={color}
                    onClick={(e) => {
                        navigate(
                            `/jobs?type=${jobType}`
                        )
                    }}
                    // valueFormatter={dataFormatter}
                />
            )}
        </Card>
    )
}

interface IJobsMenu {
    isCollapsed: boolean
    workspace: string | undefined
}

export default function JobsMenu({ isCollapsed, workspace }: IJobsMenu) {
    const { isAuthenticated, getAccessTokenSilently } = useAuth()
    const navigate = useNavigate()

    const {
        response: jobs,
        isLoading,
        sendNow,
    } = useScheduleApiV1JobsCreate(
        {
            pageStart: 0,
            pageEnd: 1,
            hours: 24,
        },
        {},
        true,
        workspace
    )
    useEffect(() => {
        if (isAuthenticated) {
            getAccessTokenSilently()
                .then((res) => {
                    console.log('')
                })
                .then((res) => {
                    sendNow()
                })
        }
    }, [isAuthenticated, workspace])

    if (workspace === undefined || workspace === '') {
        return null
    }

    return (
        <Popover className="relative z-50 border-0 w-full">
            <Popover.Button
                className={`w-full px-6 py-2 flex items-center rounded-md gap-2.5 text-gray-50 hover:bg-openg-800 ${
                    isCollapsed ? '!p-2' : ''
                }`}
                id="Jobs"
            >
                <ClipboardDocumentListIcon className="h-5 w-5 stroke-2 text-gray-400" />
                {!isCollapsed && (
                    <Text className="text-inherit">Activity Log</Text>
                )}
            </Popover.Button>
            <Transition
                as={Fragment}
                enter="transition ease-out duration-200"
                enterFrom="opacity-0 translate-y-1"
                enterTo="opacity-100 translate-y-0"
                leave="transition ease-in duration-150"
                leaveFrom="opacity-100 translate-y-0"
                leaveTo="opacity-0 translate-y-1"
            >
                <Popover.Panel
                    className={`absolute ${
                        isCollapsed ? 'left-[57px]' : 'left-[292px]'
                    } bottom-0 z-10`}
                >
                    <Card className="w-fit dark:!ring-gray-600">
                        <Flex justifyContent="between">
                            <Title className="font-bold text-gray-800">
                                Jobs in last 24 hours
                            </Title>
                            <Button
                                size="xs"
                                variant="light"
                                icon={ChevronRightIcon}
                                iconPosition="right"
                                onClick={() =>
                                    navigate(
                                        `/jobs`
                                    )
                                }
                            >
                                See all
                            </Button>
                        </Flex>
                        {isLoading ? (
                            <Flex className="min-w-[400px] min-h-[256px]">
                                <Spinner />
                            </Flex>
                        ) : (
                            <Flex
                                flexDirection="col"
                                justifyContent="between"
                                alignItems="start"
                                className="mt-6"
                            >
                                <JobCategoryItem
                                    title="Discovery"
                                    jobType="discovery"
                                    summaries={jobs?.summaries?.filter(
                                        (v) =>
                                            v.type ===
                                            PlatformEnginePkgDescribeApiJobType.JobTypeDiscovery
                                    )}
                                />
                                <JobCategoryItem
                                    title="Metric"
                                    jobType="analytics"
                                    summaries={jobs?.summaries?.filter(
                                        (v) =>
                                            v.type ===
                                            PlatformEnginePkgDescribeApiJobType.JobTypeAnalytics
                                    )}
                                />
                                <JobCategoryItem
                                    title="Governance"
                                    jobType="compliance"
                                    summaries={jobs?.summaries?.filter(
                                        (v) =>
                                            v.type ===
                                            PlatformEnginePkgDescribeApiJobType.JobTypeCompliance
                                    )}
                                />
                            </Flex>
                        )}
                    </Card>
                </Popover.Panel>
            </Transition>
        </Popover>
    )
}
