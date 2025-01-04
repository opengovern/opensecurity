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
import {
    UncontrolledTreeEnvironment,
    Tree,
    StaticTreeDataProvider,
    ControlledTreeEnvironment,
} from 'react-complex-tree'
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
import BenchmarkChart from '../../../../components/Benchmark/Chart'
import { toErrorMessage } from '../../../../types/apierror'
import SummaryCard from '../../../../components/Cards/SummaryCard'
import Evaluate from './Evaluate'
import Table, { IColumn } from '../../../../components/Table'
import { ValueFormatterParams } from 'ag-grid-community'
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
import Filter from './Filter'
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

    const [recall, setRecall] = useState(false)
    const [focusedItem, setFocusedItem] = useState<string>()
    const [expandedItems, setExpandedItems] = useState<string[]>([])
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

    const {
        response: benchmarkKPIStart,
        isLoading: benchmarkKPIStartLoading,
        sendNow: benchmarkKPIStartSend,
    } = useComplianceApiV1BenchmarksSummaryDetail(String(benchmarkId), {
        ...topQuery,
        timeAt: activeTimeRange.start.unix(),
    })
    const {
        response: benchmarkKPIEnd,
        isLoading: benchmarkKPIEndLoading,
        sendNow: benchmarkKPIEndSend,
    } = useComplianceApiV1BenchmarksSummaryDetail(String(benchmarkId), {
        ...topQuery,
        timeAt: activeTimeRange.end.unix(),
    })

    const hideKPIs =
        (benchmarkKPIEnd?.conformanceStatusSummary?.failed || 0) +
            (benchmarkKPIEnd?.conformanceStatusSummary?.passed || 0) +
            (benchmarkKPIStart?.conformanceStatusSummary?.failed || 0) +
            (benchmarkKPIStart?.conformanceStatusSummary?.passed || 0) ===
        0

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
            ![
                'baseline_efficiency',
                'baseline_reliability',
                'baseline_security',
                'baseline_supportability',
            ].includes(benchmarkDetail?.id)
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
            {/* <TopHeader
                breadCrumb={[
                    benchmarkDetail?.title
                        ? benchmarkDetail?.title
                        : 'Benchmark summary',
                ]}
                supportedFilters={
                    enable ? ['Date', 'Cloud Account', 'Connector'] : []
                }
                initialFilters={enable ? ['Date'] : []}
            /> */}
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

                    <Container
                        disableHeaderPaddings
                        disableContentPaddings
                        className="rounded-xl  bg-[#0f2940] p-0 text-white mt-4"
                        footer={
                            false ? (
                                <>
                                    <ExpandableSection
                                        header="Additional settings"
                                        variant="footer"
                                    >
                                        <Flex
                                            justifyContent="end"
                                            className="bg-white p-4 pt-0 mb-2 w-full gap-3    rounded-xl"
                                        >
                                            <Filter
                                                type={selectedGroup}
                                                onApply={(e) => {
                                                    setAccount(e.connector)
                                                }}
                                                // id={id}
                                            />
                                            <DateRangePicker
                                                onChange={({ detail }) => {
                                                    setValue(detail.value)
                                                }}
                                                value={value}
                                                placeholder={
                                                    'Please select Date'
                                                }
                                                // disabled={true}
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
                                                        key: 'previous-7-Days',
                                                        amount: 7,
                                                        unit: 'day',
                                                        type: 'relative',
                                                    },
                                                ]}
                                                isValidRange={(range) => {
                                                    if (
                                                        range.type ===
                                                        'absolute'
                                                    ) {
                                                        const [
                                                            startDateWithoutTime,
                                                        ] =
                                                            range.startDate.split(
                                                                'T'
                                                            )
                                                        const [
                                                            endDateWithoutTime,
                                                        ] =
                                                            range.endDate.split(
                                                                'T'
                                                            )
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
                                                            new Date(
                                                                range.startDate
                                                            ) -
                                                                new Date(
                                                                    range.endDate
                                                                ) >
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
                                                absoluteFormat="long-localized"
                                                hideTimeOffset
                                                dateOnly={true}
                                                // placeholder="Filter by a date and time range"
                                            />
                                        </Flex>
                                    </ExpandableSection>
                                </>
                            ) : (
                                ''
                            )
                        }
                        header={
                            <Header
                                className={`bg-[#0f2940] p-4 pt-0 rounded-xl   text-white ${
                                    false ? 'rounded-b-none' : ''
                                }`}
                                variant="h2"
                                description=""
                            >
                                <SpaceBetween size="xxxs" direction="vertical">
                                    <Box className="rounded-xl same text-white pt-3 pl-3 pb-0">
                                        <Grid
                                            gridDefinition={[
                                                {
                                                    colspan: {
                                                        default: 12,
                                                        xs: 8,
                                                        s: 9,
                                                    },
                                                },
                                                {
                                                    colspan: {
                                                        default: 12,
                                                        xs: 4,
                                                        s: 3,
                                                    },
                                                },
                                            ]}
                                        >
                                            <div>
                                                <Box
                                                    variant="h1"
                                                    className="text-white important"
                                                    color="white"
                                                >
                                                    <span className="text-white">
                                                        {benchmarkDetail?.title}
                                                    </span>
                                                </Box>
                                                <Box
                                                    variant="p"
                                                    color="white"
                                                    margin={{
                                                        top: 'xxs',
                                                        bottom: 's',
                                                    }}
                                                >
                                                    <div className="group text-white important  relative flex text-wrap justify-start">
                                                        <Text className="test-start w-full text-white ">
                                                            {/* @ts-ignore */}
                                                            {truncate(
                                                                benchmarkDetail?.description
                                                            )}
                                                        </Text>
                                                        <Card className="absolute w-full text-wrap z-40 top-0 scale-0 transition-all p-2 group-hover:scale-100">
                                                            <Text>
                                                                {
                                                                    benchmarkDetail?.description
                                                                }
                                                            </Text>
                                                        </Card>
                                                    </div>
                                                </Box>
                                            </div>
                                        </Grid>
                                    </Box>
                                    <Flex className="w-max pl-3">
                                        <Evaluate
                                            id={benchmarkDetail?.id}
                                            benchmarkDetail={benchmarkDetail}
                                            assignmentsCount={assignments}
                                            onEvaluate={(c, b) => {
                                                RunBenchmark(c, b)
                                            }}
                                        />
                                    </Flex>
                                </SpaceBetween>
                            </Header>
                        }
                    ></Container>

                    {/* <Flex alignItems="start" className="mb-3 w-11/12">
                        <Flex
                            flexDirection="col"
                            alignItems="start"
                            justifyContent="start"
                            className="gap-2 w-full"
                        >
                            <Title className="font-semibold">
                                {benchmarkDetail?.title}
                            </Title>
                            <div className="group  relative flex text-wrap justify-start">
                                <Text className="test-start w-full ">
                                    {/* @ts-ignore 
                                    {truncate(benchmarkDetail?.description)}
                                </Text>
                                <Card className="absolute w-full text-wrap z-40 top-0 scale-0 transition-all p-2 group-hover:scale-100">
                                    <Text>{benchmarkDetail?.description}</Text>
                                </Card>
                            </div>
                        </Flex>
                        <Flex className="w-fit gap-4">
                             <Settings
                                id={benchmarkDetail?.id}
                                response={(e) => setAssignments(e)}
                                autoAssign={benchmarkDetail?.autoAssign}
                                tracksDriftEvents={
                                    benchmarkDetail?.tracksDriftEvents
                                }
                                isAutoResponse={(x) => setRecall(true)}
                                reload={() => updateDetail()}
                            /> 
                            <Evaluate
                                id={benchmarkDetail?.id}
                                benchmarkDetail={benchmarkDetail}
                                assignmentsCount={assignments}
                                onEvaluate={(c) => {
                                    RunBenchmark(c)
                                }}
                            />
                        </Flex>
                    </Flex> */}
                    <Flex flexDirection="col" className="w-full mt-4">
                        {/* {chart && enable && ( */}
                        {false && (
                            <>
                                <Flex className="bg-white  w-full border-solid border-2    rounded-xl p-4">
                                    {/* <LineChart
                                        className="h-80"
                                        data={chart?.length < 0 ? [] : chart}
                                        index="date"
                                        categories={[
                                            // 'Total',
                                            'Incidents',
                                            'Non Compliant',
                                        ]}
                                        colors={['indigo', 'rose', 'cyan']}
                                        noDataText="No Data to Display"
                                        // valueFormatter={dataFormatter}
                                        yAxisWidth={60}
                                        onValueChange={(v) => console.log(v)}
                                    /> */}
                                    <ReactEcharts
                                        // echarts={echarts}
                                        option={options()}
                                        className="w-full"
                                        onEvents={() => {}}
                                    />
                                </Flex>
                            </>
                        )}

                        <Flex className="mt-2">
                            <Tabs
                                className="mt-6 rounded-[1px] rounded-s-none rounded-e-none"
                                // variant="container"
                                tabs={find_tabs()}
                            />
                        </Flex>
                    </Flex>

                    {/* <TabGroup
                        className="mt-2"
                        // tabIndex={enable ? tab : findTab(tab)}
                        // onIndexChange={(index) => {
                        //     console.log("its index")
                        //     if(enable){
                        //         setTab(index)

                        //     }
                        //     else{
                        //         // @ts-ignore
                        //         setTab(findTab(index))
                        //     }
                        // }}
                    >
                        <TabList className="mb-4">
                            <>
                                {enable && (
                                    <>
                                        <Tab
                                            onClick={() => {
                                                setTab(0)
                                            }}
                                        >
                                            Summary
                                        </Tab>
                                    </>
                                )}

                                <Tab
                                    onClick={() => {
                                        setTab(1)
                                    }}
                                >
                                    Controls
                                </Tab>
                                {enable && (
                                    <>
                                        <Tab
                                            onClick={() => {
                                                setTab(2)
                                            }}
                                        >
                                            Incidents
                                        </Tab>
                                    </>
                                )}

                                <Tab
                                    onClick={() => {
                                        setTab(3)
                                    }}
                                >
                                    Scope Assignments
                                </Tab>
                                <Tab
                                    onClick={() => {
                                        setTab(4)
                                    }}
                                >
                                    Run History
                                </Tab>
                            </>
                        </TabList>
                        <TabPanels>
                            {enable && (
                                <>
                                    {tab == 0 && (
                                        <>
                                            <Flex
                                                className="w-full flex-wrap"
                                                flexDirection="col"
                                            >
                                                {hideKPIs ? (
                                                    ''
                                                ) : (
                                                    <Grid
                                                        numItems={4}
                                                        className="w-full gap-4 mb-4"
                                                    >
                                                        <SummaryCard
                                                            title="Security Score"
                                                            metric={
                                                                ((benchmarkKPIEnd
                                                                    ?.controlsSeverityStatus
                                                                    ?.total
                                                                    ?.passed ||
                                                                    0) /
                                                                    (benchmarkKPIEnd
                                                                        ?.controlsSeverityStatus
                                                                        ?.total
                                                                        ?.total ||
                                                                        1)) *
                                                                    100 || 0
                                                            }
                                                            metricPrev={
                                                                ((benchmarkKPIStart
                                                                    ?.controlsSeverityStatus
                                                                    ?.total
                                                                    ?.passed ||
                                                                    0) /
                                                                    (benchmarkKPIStart
                                                                        ?.controlsSeverityStatus
                                                                        ?.total
                                                                        ?.total ||
                                                                        1)) *
                                                                    100 || 0
                                                            }
                                                            isPercent
                                                            loading={
                                                                benchmarkKPIEndLoading ||
                                                                benchmarkKPIStartLoading
                                                            }
                                                        />
                                                        <SummaryCard
                                                            title="Issues"
                                                            metric={
                                                                benchmarkKPIEnd
                                                                    ?.conformanceStatusSummary
                                                                    ?.failed
                                                            }
                                                            metricPrev={
                                                                benchmarkKPIStart
                                                                    ?.conformanceStatusSummary
                                                                    ?.failed
                                                            }
                                                            loading={
                                                                benchmarkKPIEndLoading ||
                                                                benchmarkKPIStartLoading
                                                            }
                                                        />

                                                        <SummaryCard
                                                            title="Passed"
                                                            metric={
                                                                benchmarkKPIEnd
                                                                    ?.conformanceStatusSummary
                                                                    ?.passed
                                                            }
                                                            metricPrev={
                                                                benchmarkKPIStart
                                                                    ?.conformanceStatusSummary
                                                                    ?.passed
                                                            }
                                                            loading={
                                                                benchmarkKPIEndLoading ||
                                                                benchmarkKPIStartLoading
                                                            }
                                                        />

                                                        <SummaryCard
                                                            title="Accounts"
                                                            metric={
                                                                benchmarkKPIEnd
                                                                    ?.connectionsStatus
                                                                    ?.total
                                                            }
                                                            metricPrev={
                                                                benchmarkKPIStart
                                                                    ?.connectionsStatus
                                                                    ?.total
                                                            }
                                                            loading={
                                                                benchmarkKPIEndLoading ||
                                                                benchmarkKPIStartLoading
                                                            }
                                                        />

                                                        {/* <SummaryCard
                                title="Events"
                                metric={events?.count}
                                loading={eventsLoading}
                            /> 
                                                    </Grid>
                                                )}
                                                {trend === null ? (
                                                    ''
                                                ) : (
                                                    <BenchmarkChart
                                                        title="Security Score"
                                                        isLoading={trendLoading}
                                                        trend={trend}
                                                        error={toErrorMessage(
                                                            trendError
                                                        )}
                                                        onRefresh={() =>
                                                            sendTrend()
                                                        }
                                                    />
                                                )}
                                            </Flex>
                                        </>
                                    )}{' '}
                                </>
                            )}
                            {tab == 1 && (
                                <div className="w-full flex flex-row justify-start items-start ">
                                    <div className="w-11/12">
                                        <Controls
                                            id={String(benchmarkId)}
                                            assignments={trend === null ? 0 : 1}
                                            enable={enable}
                                        />
                                    </div>
                                </div>
                            )}{' '}
                            {enable && (
                                <>
                                    {' '}
                                    {tab == 2 && (
                                        <>
                                            <Findings
                                                id={
                                                    benchmarkId
                                                        ? benchmarkId
                                                        : ''
                                                }
                                            />
                                        </>
                                    )}
                                </>
                            )}
                            {tab == 3 && (
                                <>
                                    <Settings
                                        id={benchmarkDetail?.id}
                                        response={(e) => setAssignments(e)}
                                        autoAssign={benchmarkDetail?.autoAssign}
                                        tracksDriftEvents={
                                            benchmarkDetail?.tracksDriftEvents
                                        }
                                        isAutoResponse={(x) => setRecall(true)}
                                        reload={() => updateDetail()}
                                    />
                                </>
                            )}
                            {tab == 4 && (
                                <>
                                    <EvaluateTable
                                        id={benchmarkDetail?.id}
                                        benchmarkDetail={benchmarkDetail}
                                        assignmentsCount={assignments}
                                        onEvaluate={(c) => {
                                            triggerEvaluate(
                                                {
                                                    benchmark_id: [
                                                        benchmarkId || '',
                                                    ],
                                                    connection_id: c,
                                                },
                                                {}
                                            )
                                        }}
                                    />
                                </>
                            )}
                        </TabPanels>
                    </TabGroup> */}
                </>
            )}
        </>
    )
}

// {
//     label: 'Summary',
//     id: 'first',
//     content: (
//         <Flex
//             className="w-full flex-wrap"
//             flexDirection="col"
//         >
//             {/* {hideKPIs ? (
//                 ''
//             ) : (
//                 <Grid
//                     numItems={4}
//                     className="w-full gap-4 mb-4"
//                 >
//                     <SummaryCard
//                         title="Security Score"
//                         metric={
//                             ((benchmarkKPIEnd
//                                 ?.controlsSeverityStatus
//                                 ?.total?.passed ||
//                                 0) /
//                                 (benchmarkKPIEnd
//                                     ?.controlsSeverityStatus
//                                     ?.total
//                                     ?.total || 1)) *
//                                 100 || 0
//                         }
//                         metricPrev={
//                             ((benchmarkKPIStart
//                                 ?.controlsSeverityStatus
//                                 ?.total?.passed ||
//                                 0) /
//                                 (benchmarkKPIStart
//                                     ?.controlsSeverityStatus
//                                     ?.total
//                                     ?.total || 1)) *
//                                 100 || 0
//                         }
//                         isPercent
//                         loading={
//                             benchmarkKPIEndLoading ||
//                             benchmarkKPIStartLoading
//                         }
//                     />
//                     <SummaryCard
//                         title="Issues"
//                         metric={
//                             benchmarkKPIEnd
//                                 ?.conformanceStatusSummary
//                                 ?.failed
//                         }
//                         metricPrev={
//                             benchmarkKPIStart
//                                 ?.conformanceStatusSummary
//                                 ?.failed
//                         }
//                         loading={
//                             benchmarkKPIEndLoading ||
//                             benchmarkKPIStartLoading
//                         }
//                     />

//                     <SummaryCard
//                         title="Passed"
//                         metric={
//                             benchmarkKPIEnd
//                                 ?.conformanceStatusSummary
//                                 ?.passed
//                         }
//                         metricPrev={
//                             benchmarkKPIStart
//                                 ?.conformanceStatusSummary
//                                 ?.passed
//                         }
//                         loading={
//                             benchmarkKPIEndLoading ||
//                             benchmarkKPIStartLoading
//                         }
//                     />

//                     <SummaryCard
//                         title="Accounts"
//                         metric={
//                             benchmarkKPIEnd
//                                 ?.connectionsStatus
//                                 ?.total
//                         }
//                         metricPrev={
//                             benchmarkKPIStart
//                                 ?.connectionsStatus
//                                 ?.total
//                         }
//                         loading={
//                             benchmarkKPIEndLoading ||
//                             benchmarkKPIStartLoading
//                         }
//                     />
//                 </Grid>
//             )} */}
//             {trend === null ? (
//                 ''
//             ) : (
//                 // <BenchmarkChart
//                 //     title="Security Score"
//                 //     isLoading={trendLoading}
//                 //     trend={trend}
//                 //     error={toErrorMessage(
//                 //         trendError
//                 //     )}
//                 //     onRefresh={() => sendTrend()}
//                 // />
//
