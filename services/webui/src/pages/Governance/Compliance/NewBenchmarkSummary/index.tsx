import { useParams } from 'react-router-dom'
import {
    Card,
    Flex,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    Title,
    Switch,
} from '@tremor/react'

import Tabs from '@cloudscape-design/components/tabs'
import Box from '@cloudscape-design/components/box'
// import Button from '@cloudscape-design/components/button'
import Grid from '@cloudscape-design/components/grid'
import DateRangePicker from '@cloudscape-design/components/date-range-picker'

import { useEffect, useState } from 'react'
import {
    useComplianceApiV1BenchmarksSummaryDetail,
    useComplianceApiV1FindingEventsCountList,
} from '../../../../api/compliance.gen'
import { useScheduleApiV1ComplianceTriggerUpdate } from '../../../../api/schedule.gen'
import Spinner from '../../../../components/Spinner'
import Controls from './Controls'
import Settings from './Settings'
import TopHeader from '../../../../components/Layout/Header'
import {
    defaultTime,
    useFilterState,
    useUrlDateRangeState,
} from '../../../../utilities/urlstate'

import { toErrorMessage } from '../../../../types/apierror'

import Evaluate from './Evaluate'

import Findings from './Findings'
import axios from 'axios'
import { get } from 'http'
import EvaluateTable from './EvaluateTable'
import { notificationAtom } from '../../../../store'
import { useSetAtom } from 'jotai'
import ContentLayout from '@cloudscape-design/components/content-layout'
import Container from '@cloudscape-design/components/container'
import Header from '@cloudscape-design/components/header'
import Link from '@cloudscape-design/components/link'
import Button from '@cloudscape-design/components/button'
// import { LineChart } from '@tremor/react'
import {
    BreadcrumbGroup,
    ButtonDropdown,
    ExpandableSection,
    SpaceBetween,
} from '@cloudscape-design/components'
import ReactEcharts from 'echarts-for-react'
import { numericDisplay } from '../../../../utilities/numericDisplay'

export default function NewBenchmarkSummary() {
    const { ws } = useParams()

    const [enable, setEnable] = useState<boolean>(false)
    const [is_baseline, setIs_baseline] = useState<boolean>(false)
    const setNotification = useSetAtom(notificationAtom)
    const { benchmarkId } = useParams()
    const [assignments, setAssignments] = useState(0)
    const [runLoading, setRunLoading] = useState(false)
    const [open, setOpen] = useState(false)

    const {
        response: benchmarkDetail,
        isLoading,
        sendNow: updateDetail,
    } = useComplianceApiV1BenchmarksSummaryDetail(String(benchmarkId))
    const { sendNowWithParams: triggerEvaluate, isExecuted } =
        useScheduleApiV1ComplianceTriggerUpdate(
            {
                benchmark_id: [benchmarkId ? benchmarkId : ''],
                connection_id: [],
            },
            {},
            false
        )

    const RunBenchmark = (c: any[], b: boolean) => {
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        const body = {
            // with_incidents: true,
            with_incidents: b,

            integration_info: c.map((c) => {
                return {
                    integration_id: c.value,
                }
            }),
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        setRunLoading(true)
        axios
            .post(
                `${url}/main/schedule/api/v3/compliance/benchmark/${benchmarkId}/run`,
                body,
                config
            )
            .then((res) => {
                let ids = ''
                res.data.jobs.map((item: any, index: number) => {
                    if (index < 5) {
                        ids = ids + item.job_id + ','
                    }
                })
                setNotification({
                    text: `Run is Done You Job id is ${ids}`,
                    type: 'success',
                })
                setRunLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setRunLoading(false)
            })
    }
    const UpdateBenchmark = (is_baseline: boolean,enabled:boolean) => {
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        const body = {
            is_baseline,
            enabled
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        axios
            .put(
                `${url}/main/compliance/api/v1/frameworks/${benchmarkId}`,
                body,
                config
            )
            .then((res) => {
              
                setNotification({
                    text: `Framework changed successfully`, 
                    type: 'success',
                })
                window.location.reload()
            })
            .catch((err) => {
                console.log(err)
            })
    }
    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 600 ? text.substring(0, 600) + '...' : text
        }
    }

    useEffect(() => {
        if (isExecuted) {
            updateDetail()
        }
    }, [isExecuted])

    useEffect(() => {
        // @ts-ignore
        setEnable(benchmarkDetail?.enabled)
        // @ts-ignore
        setIs_baseline(benchmarkDetail?.isBaseline)
    }, [benchmarkDetail])
    const find_tabs = () => {
        const tabs = []
        tabs.push({
            label: 'Controls',
            id: 'second',
            content: (
                <div className="w-full flex flex-row justify-start items-start ">
                    <div className="w-full">
                        <Controls
                            id={String(benchmarkId)}
                            assignments={1}
                            enable={enable}
                        />
                    </div>
                </div>
            ),
        })
        tabs.push({
            label: 'Framework-Specific Incidents',
            id: 'third',
            content: <Findings id={benchmarkId ? benchmarkId : ''} />,
            disabled: false,
            disabledReason:
                'This is available when the Framework has at least one assignments.',
        })
        if (enable) {
            tabs.push({
                label: 'Settings',
                id: 'fourth',
                content: (
                    <Settings
                        id={benchmarkDetail?.id}
                        response={(e) => setAssignments(e)}
                        autoAssign={benchmarkDetail?.autoAssign}
                        reload={() => updateDetail()}
                    />
                ),
                disabled: false,
            })
        }
        tabs.push({
            label: 'Run History',
            id: 'fifth',
            content: (
                <EvaluateTable
                    id={benchmarkDetail?.id}
                    benchmarkDetail={benchmarkDetail}
                    assignmentsCount={assignments}
                    onEvaluate={(c) => {
                        triggerEvaluate(
                            {
                                benchmark_id: [benchmarkId || ''],
                                connection_id: c,
                            },
                            {}
                        )
                    }}
                />
            ),
            // disabled: true,
            // disabledReason: 'COMING SOON',
        })
        return tabs
    }
    const GetActions = () =>{
        const temp = []
        if(enable){
            temp.push({
                text: 'Disable framework',
                id: 'disable',
            })
        }
        else{
            temp.push({
                text: 'Enable framework',
                id: 'enable',
            })
        }
        if(!is_baseline){
            temp.push({
                text: 'Set as baseline',
                id: 'add_baseline',
            })
        }
        else{
            temp.push({
                text: 'Remove as baseline',
                id: 'remove_baseline',
            })
        }
        return temp
    }

    return (
        <>
            {isLoading ? (
                <Spinner className="mt-56" />
            ) : (
                <>
                    {/* <BreadcrumbGroup
                        onClick={(event) => {
                            // event.preventDefault()
                        }}
                        items={[
                            {
                                text: 'Compliance',
                                href: `/compliance`,
                            },
                            { text: 'Frameworks', href: '#' },
                        ]}
                        ariaLabel="Breadcrumbs"
                    /> */}
                    <Header
                        className={`   rounded-xl mt-6   ${
                            false ? 'rounded-b-none' : ''
                        }`}
                        actions={
                            <Flex className="w-max ">
                                <ButtonDropdown
                                    onItemClick={({ detail }) => {
                                        const id = detail.id
                                        switch (id) {
                                            case 'enable':
                                                UpdateBenchmark(is_baseline,true)
                                                break
                                            case 'disable':
                                                UpdateBenchmark(
                                                    is_baseline,
                                                    false
                                                )

                                                break
                                            case 'add_baseline':
                                                UpdateBenchmark(
                                                    true,
                                                    enable
                                                )

                                                break
                                            case 'remove_baseline':
                                                UpdateBenchmark(
                                                    false,
                                                    enable
                                                )

                                                break

                                            default:
                                                break
                                        }
                                    }}
                                    variant="primary"
                                    items={[
                                        {
                                            text: 'Settings',
                                            items: GetActions()
                                        },
                                    ]}
                                    mainAction={{
                                        text: 'Run ',

                                        onClick: () => {
                                            setOpen(true)
                                        },
                                        loading: runLoading,
                                    }}
                                    
                                />
                            </Flex>
                        }
                        variant="h2"
                        description={
                            <div className="group  important text-black  relative sm:flex hidden text-wrap justify-start">
                                <Text className="test-start w-full text-black  ">
                                    {/* @ts-ignore */}
                                    {truncate(benchmarkDetail?.description)}
                                </Text>
                                <Card className="absolute w-full text-wrap text-black z-40 top-0 scale-0 transition-all p-2 group-hover:scale-100">
                                    <Text>{benchmarkDetail?.description}</Text>
                                </Card>
                            </div>
                        }
                    >
                        <Flex className="gap-2">
                            <span>{benchmarkDetail?.title}</span>
                            {/* <Button iconName="status-info" variant="icon" onClick={()=>{
                                GetCoverage()
                            }} /> */}
                        </Flex>
                    </Header>

                    <Flex flexDirection="col" className="w-full ">
                        <Flex className="">
                            <Tabs
                                className="mt-4 rounded-[1px] rounded-s-none rounded-e-none"
                                // variant="container"
                                tabs={find_tabs()}
                            />
                        </Flex>
                    </Flex>
                    <Evaluate
                        open={open}
                        setOpen={setOpen}
                        benchmarkDetail={benchmarkDetail}
                        onEvaluate={(c, b) => {
                            RunBenchmark(c, b)
                        }}
                        showBenchmark={false}
                    />
                </>
            )}
        </>
    )
}
