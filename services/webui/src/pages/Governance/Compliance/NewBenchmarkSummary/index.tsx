// @ts-nocheck
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
    ExpandableSection,
    SpaceBetween,
} from '@cloudscape-design/components'
import ReactEcharts from 'echarts-for-react'
import { numericDisplay } from '../../../../utilities/numericDisplay'

export default function NewBenchmarkSummary() {
    const { ws } = useParams()
    const { value: activeTimeRange } = useUrlDateRangeState(
        defaultTime(ws || '')
    )
    const [tab, setTab] = useState<number>(0)
    const [enable, setEnable] = useState<boolean>(false)
    const [chart, setChart] = useState()
    const options = () => {
        const confine = true
        const opt = {
            tooltip: {
                confine,
                trigger: 'axis',
                axisPointer: {
                    type: 'line',
                    label: {
                        formatter: (param: any) => {
                            let total = 0
                            if (param.seriesData && param.seriesData.length) {
                                for (
                                    let i = 0;
                                    i < param.seriesData.length;
                                    i += 1
                                ) {
                                    total += param.seriesData[i].data
                                }
                            }

                            return `${param.value} (Total: ${total.toFixed(2)})`
                        },
                        // backgroundColor: '#6a7985',
                    },
                },
                valueFormatter: (value: number | string) => {
                    return numericDisplay(value)
                },
                order: 'valueDesc',
            },
            grid: {
                left: 45,
                right: 0,
                top: 20,
                bottom: 40,
            },
            xAxis: {
                type: 'category',
                data: chart?.map((item) => {
                    return item.date
                }),
            },
            yAxis: {
                type: 'value',
            },
            series: [
                {
                    name: 'Incidents',
                    data: chart?.map((item) => {
                        return item.Incidents
                    }),
                    type: 'line',
                },
                {
                    name: 'Non Compliant',

                    data: chart?.map((item) => {
                        return item['Non Compliant']
                    }),
                    type: 'line',
                },
                {
                    name: 'High',
                    data: chart?.map((item) => {
                        return item.High
                    }),
                    type: 'line',
                },
                {
                    name: 'Medium',
                    data: chart?.map((item) => {
                        return item.Medium
                    }),
                    type: 'line',
                },
                {
                    name: 'Low',
                    data: chart?.map((item) => {
                        return item.Low
                    }),
                    type: 'line',
                },
                {
                    name: 'Critical',
                    data: chart?.map((item) => {
                        return item.Critical
                    }),
                    type: 'line',
                },
            ],
        }
        return opt
    }

    const setNotification = useSetAtom(notificationAtom)
    const [selectedGroup, setSelectedGroup] = useState<
        'findings' | 'resources' | 'controls' | 'accounts' | 'events'
    >('accounts')
    const [account, setAccount] = useState([])
    const readTemplate = (template: any, data: any = { items: {} }): any => {
        for (const [key, value] of Object.entries(template)) {
            // eslint-disable-next-line no-param-reassign
            data.items[key] = {
                index: key,
                canMove: true,
                isFolder: value !== null,
                children:
                    value !== null
                        ? Object.keys(value as Record<string, unknown>)
                        : undefined,
                data: key,
                canRename: true,
            }

            if (value !== null) {
                readTemplate(value, data)
            }
        }
        return data
    }
    const shortTreeTemplate = {
        root: {
            container: {
                item0: null,
                item1: null,
                item2: null,
                item3: {
                    inner0: null,
                    inner1: null,
                    inner2: null,
                    inner3: null,
                },
                item4: null,
                item5: null,
            },
        },
    }
    const shortTree = readTemplate(shortTreeTemplate)

    const { benchmarkId } = useParams()
    const { value: selectedConnections } = useFilterState()
    const [assignments, setAssignments] = useState(0)
    const [coverage, setCoverage] = useState([])
    const [recall, setRecall] = useState(false)
    const topQuery = {
        ...(benchmarkId && { benchmarkId: [benchmarkId] }),
        ...(selectedConnections.provider && {
            integrationType: [selectedConnections.provider],
        }),
        ...(selectedConnections.connections && {
            integrationID: selectedConnections.connections,
        }),
        ...(selectedConnections.connectionGroup && {
            connectionGroup: selectedConnections.connectionGroup,
        }),
    }

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

 

    const GetEnabled = () => {
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
                `${url}/main/compliance/api/v3/benchmark/${benchmarkId}/assignments`,
                config
            )
            .then((res) => {
                if (res.data) {
                    if (
                        res.data.status == 'enabled' ||
                        res.data.status == 'auto-enable'
                    ) {
                        setEnable(true)
                        setTab(0)
                    } else {
                        setEnable(false)
                        setTab(1)
                    }
                    // if (res.data.items.length > 0) {
                    //     setEnable(true)
                    //     setTab(0)
                    // } else {
                    //     setEnable(false)
                    //     setTab(1)
                    // }
                } else {
                    setEnable(false)
                    setTab(1)
                }
            })
            .catch((err) => {
                console.log(err)
            })
    }
     const GetCoverage = () => {
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
                 `${url}/main/compliance/api/v1/frameworks/${benchmarkId}/coverage`,
                 config
             )
             .then((res) => {
                 if (res.data) {
                   
                    setCoverage(res.data)
                }
             })
             .catch((err) => {
                 console.log(err)
             })
     }
    const GetChart = () => {
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
            .post(
                `${url}/main/compliance/api/v3/benchmarks/${benchmarkId}/trend`,
                {},
                config
            )
            .then((res) => {
                const temp = res.data
                const temp_chart = temp?.datapoints?.map((item) => {
                    if (
                        item.compliance_results_summary &&
                        item.incidents_severity_breakdown
                    ) {
                        const temp_data = {
                            date: new Date(item.timestamp)
                                .toLocaleDateString('en-US', {
                                    month: 'short',
                                    day: 'numeric',
                                    hour: 'numeric',
                                    minute: 'numeric',
                                    hour12: !1,
                                })
                                .split(',')
                                .join('\n'),
                            // Total:
                            //     item?.findings_summary?.incidents +
                            //     item?.findings_summary?.non_incidents,
                            Incidents:
                                item.compliance_results_summary?.incidents,
                            'Non Compliant':
                                item.compliance_results_summary?.non_incidents,
                            High: item.incidents_severity_breakdown.highCount,
                            Medium: item.incidents_severity_breakdown
                                .mediumCount,
                            Low: item.incidents_severity_breakdown.lowCount,
                            Critical:
                                item.incidents_severity_breakdown.criticalCount,
                        }
                        return temp_data
                    }
                })
                const new_chart = temp_chart?.filter((item) => {
                    if (item) {
                        return item
                    }
                })
                setChart(new_chart)
            })
            .catch((err) => {
                console.log(err)
            })
    }
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
        //    console.log(config)
        axios
            .post(
                `${url}/main/schedule/api/v3/compliance/benchmark/${benchmarkId}/run`,
                body,
                config
            )
            .then((res) => {
                let ids = ''
                res.data.jobs.map((item, index) => {
                    if (index < 5) {
                        ids = ids + item.job_id + ','
                    }
                })
                setNotification({
                    text: `Run is Done You Job id is ${ids}`,
                    type: 'success',
                })
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
    const today = new Date()
    const lastWeek = new Date(
        today.getFullYear(),
        today.getMonth(),
        today.getDate() - 7
    )
    const [value, setValue] = useState({
        type: 'relative',
        amount: 7,
        unit: 'day',
        key: 'previous-7-Days',
    })
    // @ts-ignore

    useEffect(() => {
        if (isExecuted || recall) {
            updateDetail()
        }
    }, [isExecuted, recall])
    useEffect(() => {
        GetEnabled()
        if (enable) {
            GetChart()
        }
    }, [])
    useEffect(() => {
        if (enable) {
            GetChart()
        }
    }, [enable])
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
                            accounts={account}
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
        if (
            true
        ) {
            tabs.push({
                label: 'Settings',
                id: 'fourth',
                content: (
                    <Settings
                        id={benchmarkDetail?.id}
                        response={(e) => setAssignments(e)}
                        autoAssign={benchmarkDetail?.autoAssign}
                        tracksDriftEvents={benchmarkDetail?.tracksDriftEvents}
                        isAutoResponse={(x) => setRecall(true)}
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

    return (
        <>
            {isLoading ? (
                <Spinner className="mt-56" />
            ) : (
                <>
                    <BreadcrumbGroup
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
                    />
                    <Header
                        className={`   rounded-xl mt-6   ${
                            false ? 'rounded-b-none' : ''
                        }`}
                        actions={
                            <Flex className="w-max ">
                                <Evaluate
                                    id={benchmarkDetail?.id}
                                    benchmarkDetail={benchmarkDetail}
                                    assignmentsCount={assignments}
                                    onEvaluate={(c, b) => {
                                        RunBenchmark(c, b)
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
                            <Button iconName="status-info" variant="icon" onClick={()=>{
                                GetCoverage()
                            }} />
                        </Flex>
                    </Header>

                    <Flex flexDirection="col" className="w-full ">
                        {/* {chart && enable && ( */}
                        {false && (
                            <>
                                <Flex className="bg-white  w-full border-solid border-2    rounded-xl p-4">
                                    <ReactEcharts
                                        // echarts={echarts}
                                        option={options()}
                                        className="w-full"
                                        onEvents={() => {}}
                                    />
                                </Flex>
                            </>
                        )}

                        <Flex className="">
                            <Tabs
                                className="mt-4 rounded-[1px] rounded-s-none rounded-e-none"
                                // variant="container"
                                tabs={find_tabs()}
                            />
                        </Flex>
                    </Flex>
                </>
            )}
        </>
    )
}

