// @ts-nocheck
import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Badge,
    Card,
    Color,
    Divider,
    Flex,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    Title,
} from '@tremor/react'

import { Radio } from 'pretty-checkbox-react'
import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import {
    Api,
    PlatformEnginePkgDescribeApiJob,
} from '../../../api/api'
import AxiosAPI from '../../../api/ApiConfig'
import { useScheduleApiV1JobsCreate } from '../../../api/schedule.gen'
import DrawerPanel from '../../../components/DrawerPanel'
import KFilter from '../../../components/Filter'
import { CloudIcon } from '@heroicons/react/24/outline'
import { string } from 'prop-types'
import SettingsALLJobs from './AllJobs'
import SettingsCustomization from './Customization'
import TopHeader from '../../../components/Layout/Header'
import { Tabs } from '@cloudscape-design/components'
import DiscoveryJobs from './Discovery'
import ComplianceJobs from './Compliance'

const columns = () => {
    const temp: IColumn<any, any>[] = [
        {
            field: 'id',
            headerName: 'Job ID',
            type: 'string',
            sortable: true,
            filter: false,
            suppressMenu: true,
            resizable: true,
            hide: true,
        },
        {
            field: 'createdAt',
            headerName: 'Created At',
            type: 'string',
            sortable: true,
            filter: false,
            suppressMenu: true,
            resizable: true,
            hide: false,
            cellRenderer: (params: any) => {
                return (
                    <>{`${params.value.split('T')[0]} ${
                        params.value.split('T')[1].split('.')[0]
                    } `}</>
                )
            },
        },

        {
            field: 'type',
            headerName: 'Job Type',
            type: 'string',
            sortable: true,
            filter: false,
            suppressMenu: true,
            resizable: true,
        },
        {
            field: 'integrationID',
            headerName: 'OpenGovernance Connection ID',
            type: 'string',
            sortable: true,
            filter: false,
            suppressMenu: true,
            resizable: true,
            hide: true,
        },
        {
            field: 'connectionProviderID',
            headerName: 'Account ID',
            type: 'string',
            sortable: false,
            filter: false,
            suppressMenu: true,
            resizable: true,
            hide: true,
        },
        {
            field: 'connectionProviderName',
            headerName: 'Account Name',
            type: 'string',
            sortable: false,
            filter: false,
            resizable: true,
            suppressMenu: true,
            hide: true,
        },
        {
            field: 'title',
            headerName: 'Title',
            type: 'string',
            sortable: false,
            filter: false,
            resizable: true,
            suppressMenu: false,
        },

        {
            field: 'status',
            headerName: 'Status',
            type: 'string',
            sortable: true,
            suppressMenu: true,
            filter: false,
            resizable: true,
            cellRenderer: (
                param: ValueFormatterParams<PlatformEnginePkgDescribeApiJob>
            ) => {
                let jobStatus = ''
                let jobColor: Color = 'gray'
                switch (param.data?.status) {
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
                        jobStatus = String(param.data?.status)
                }

                return <Badge color={jobColor}>{jobStatus}</Badge>
            },
        },
        {
            field: 'updatedAt',
            headerName: 'Updated At',
            type: 'date',
            sortable: true,
            filter: false,
            suppressMenu: true,
            resizable: true,
            hide: false,
            cellRenderer: (params: any) => {
                return (
                    <>{`${params.value.split('T')[0]} ${
                        params.value.split('T')[1].split('.')[0]
                    } `}</>
                )
            },
        },
        {
            field: 'failureReason',
            headerName: 'Failure Reason',
            type: 'string',
            sortable: false,
            suppressMenu: true,
            filter: true,
            resizable: true,
            hide: true,
        },
    ]
    return temp
}

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
interface Option {
    label: string | undefined
    value: string | undefined
}
export default function SettingsJobs() {
   

    return (
        <>
            {/* <Tabs
                tabs={[
                    // {
                    //     label: 'All Jobs',
                    //     content: <SettingsALLJobs />,
                    //     id: '0',
                    // },
                    {
                        label: 'Discovery ',
                        content: <DiscoveryJobs />,
                        id: '0',
                    },
                    {
                        label: 'Compliance',
                        content: <ComplianceJobs />,
                        id: '1',
                    },

                    // {
                    //     label: 'Scheduling',
                    //     content: <SettingsCustomization />,
                    //     id: '2',
                    // },
                ]}
            /> */}
            <DiscoveryJobs />
        </>
    )
}
