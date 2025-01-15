import { Card, Col, Divider, Flex, Grid, Tab, TabGroup, TabList, TabPanel, TabPanels, Title } from '@tremor/react'
import { useEffect, useState } from 'react'
import FindingsWithFailure from './FindingsWithFailure'
import TopHeader from '../../../components/Layout/Header'
import Filter from './Filter'
import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../../api/api'
import {
    DateRange,
    useURLParam,
    useURLState,
} from '../../../utilities/urlstate'
import Spinner from '../../../components/Spinner'
import AllIncidents from './AllIncidents'
import { ChevronRightIcon, DocumentTextIcon } from '@heroicons/react/24/outline'

export default function Findings() {
    const [tab, setTab] = useState<number>(0);
    const [secondTab, setSecondTab] = useState<number>(0)
    const [show,setShow] = useState<boolean>(true)
    const [selectedGroup, setSelectedGroup] = useState<
        'findings' | 'resources' | 'controls' | 'accounts' | 'events'
    >('findings')
    useEffect(() => {
        switch (tab) {
            case 0:
                setSelectedGroup('findings')
                break
            case 1:
                setSelectedGroup('resources')
                break
            default:
                setSelectedGroup('findings')
                break
        }
    }, [tab])
    useEffect(() => {
        const url = window.location.pathname.split('/')[4]
        // setShow(false);
        
        switch (url) {
            case 'summary':
                setTab(1)
                break

            default:
                setTab(0)
                break
        }
    }, [window.location.pathname])
 

    const [query, setQuery] = useState<{
        connector: SourceType
        conformanceStatus:
            | PlatformEnginePkgComplianceApiConformanceStatus[]
            | undefined
        severity: TypesFindingSeverity[] | undefined
        connectionID: string[] | undefined
        controlID: string[] | undefined
        benchmarkID: string[] | undefined
        resourceTypeID: string[] | undefined
        lifecycle: boolean[] | undefined
        activeTimeRange: DateRange | undefined
        eventTimeRange: DateRange | undefined
        jobID: string[] | undefined
        connectionGroup: []
    }>({
        connector: SourceType.Nil,
        conformanceStatus: [
            PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
        ],
        severity: [
            TypesFindingSeverity.FindingSeverityCritical,
            TypesFindingSeverity.FindingSeverityHigh,
            TypesFindingSeverity.FindingSeverityMedium,
            TypesFindingSeverity.FindingSeverityLow,
            TypesFindingSeverity.FindingSeverityNone,
        ],
        connectionID: [],
        controlID: [],
        benchmarkID: [],
        resourceTypeID: [],
        lifecycle: [true],
        activeTimeRange: undefined,
        eventTimeRange: undefined,
        jobID: [],
        connectionGroup: [],
    })

    return (
        <>
            {/* <TopHeader /> */}
            {show ? (
                <>
                    {/* <Filter type={selectedGroup} onApply={(e) => {
                        // @ts-ignore
                        setQuery(e)}} /> */}

                    <Flex className="mt-2 w-full">
                        <>
                          
                            {tab == 0 && (
                                <AllIncidents
                                    query={query}
                                    setSelectedGroup={setSelectedGroup}
                                    tab={secondTab}
                                    setTab={setSecondTab}
                                />
                            )}
                        </>
                    </Flex>
                </>
            ) : (
                <>
                    {tab == 1 && (
                        <>
                            <Flex
                                className="bg-white w-[90%] rounded-xl border-solid  border-2 border-gray-200   "
                                flexDirection="col"
                                justifyContent="center"
                                alignItems="center"
                            >
                                <div className="border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:p-6 lg:p-8">
                                    <header>
                                        <h1 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                            Posture Summary
                                        </h1>
                                        <p className="text-tremor-default text-tremor-content dark:text-dark-tremor-content">
                                            Get a summarized view of posture by
                                            Cloud Accounts, Controls, and
                                            Severity.
                                        </p>
                                        <div className="mt-8 w-full md:flex md:max-w-3xl md:items-stretch md:space-x-4">
                                            <Card className="w-full md:w-7/12">
                                                <div className="inline-flex items-center justify-center rounded-tremor-small border border-tremor-border p-2 dark:border-dark-tremor-border">
                                                    <DocumentTextIcon
                                                        className="size-5 text-tremor-content-emphasis dark:text-dark-tremor-content-emphasis"
                                                        aria-hidden={true}
                                                    />
                                                </div>
                                                <h3 className="mt-4 text-tremor-default font-medium text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                                    <a
                                                        href="https://docs.opengovernance.io/oss/how-to-guide/audit-for-compliance"
                                                        className="focus:outline-none"
                                                        target="_blank"
                                                    >
                                                        {/* Extend link to entire card */}
                                                        <span
                                                            className="absolute inset-0"
                                                            aria-hidden={true}
                                                        />
                                                        Documentation
                                                    </a>
                                                </h3>
                                                <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                    How-to guide on posture
                                                    summary.
                                                </p>
                                            </Card>
                                        </div>
                                    </header>
                                </div>
                                <div className="w-full">
                                    <div className="p-4 sm:p-6 lg:p-8">
                                        <main>
                                            <Grid
                                                // flexDirection="row"
                                                numItems={3}
                                                // justifyContent="center"
                                                // alignItems="center"
                                                className="mt-5 gap-8 flex-col w-full"
                                            >
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setSecondTab(0)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'accounts'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex flex-row gap-1 justify-center align-center items-center">
                                                                Cloud Accounts
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                View all
                                                                incidents and
                                                                drift events
                                                                across clouds,
                                                                accounts,
                                                                regions and
                                                                platforms.
                                                            </p>
                                                        </Flex>
                                                        <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                1K{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex>
                                                    </Flex>
                                                </Card>
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setSecondTab(1)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'resources'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex flex-row gap-1 justify-center align-center items-center">
                                                                Resource Posture
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                Get a summarized
                                                                view of posture
                                                                by Asset,
                                                                Entity, or
                                                                Resource Type.
                                                            </p>
                                                        </Flex>
                                                        <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                2K{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex>
                                                    </Flex>
                                                </Card>{' '}
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setSecondTab(2)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'controls'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex flex-row gap-1 justify-center align-center items-center">
                                                                Control Posture
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                Get an overview
                                                                of conformance
                                                                by controls and
                                                                identify
                                                                problematic
                                                                ones.
                                                            </p>
                                                        </Flex>
                                                        <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                2K{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex>
                                                    </Flex>
                                                </Card>{' '}
                                            </Grid>
                                        </main>
                                    </div>
                                </div>
                            </Flex>
                        </>
                    )}
                    {tab == 0 && (
                        <>
                            <Flex
                                className="bg-white w-[90%] rounded-xl border-solid  border-2 border-gray-200   "
                                flexDirection="col"
                                justifyContent="center"
                                alignItems="center"
                            >
                                <div className="border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:p-6 lg:p-8">
                                    <header>
                                        <h1 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                            Compliance Incidents
                                        </h1>
                                        <p className="text-tremor-default text-tremor-content dark:text-dark-tremor-content">
                                            See all incidents, drift incidents,
                                            evidence and summaries.
                                        </p>
                                        <div className="mt-8 w-full md:flex md:max-w-3xl md:items-stretch md:space-x-4">
                                            <Card className="w-full md:w-7/12">
                                                <div className="inline-flex items-center justify-center rounded-tremor-small border border-tremor-border p-2 dark:border-dark-tremor-border">
                                                    <DocumentTextIcon
                                                        className="size-5 text-tremor-content-emphasis dark:text-dark-tremor-content-emphasis"
                                                        aria-hidden={true}
                                                    />
                                                </div>
                                                <h3 className="mt-4 text-tremor-default font-medium text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                                    <a
                                                        href="https://docs.opengovernance.io/oss/platform/compliance"
                                                        className="focus:outline-none"
                                                        target="_blank"
                                                    >
                                                        {/* Extend link to entire card */}
                                                        <span
                                                            className="absolute inset-0"
                                                            aria-hidden={true}
                                                        />
                                                        Documentation
                                                    </a>
                                                </h3>
                                                <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                    Learn how to review events
                                                    and incidents tied to
                                                    Compliance Frameworks.
                                                    Access all incidents, drift
                                                    events, evidence, and
                                                    summaries.
                                                </p>
                                            </Card>
                                        </div>
                                    </header>
                                </div>
                                <div className="w-full">
                                    <div className="p-4  sm:p-4 lg:p-8 lg:pt-2">
                                        <main>
                                            <Grid
                                                // flexDirection="row"
                                                numItems={3}
                                                // justifyContent="center"
                                                // alignItems="center"
                                                className="mt-5 gap-8 flex-col w-full"
                                            >
                                                <Col numColSpan={3}>
                                                    <Title>
                                                        Incident Summary
                                                    </Title>
                                                </Col>
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setTab(1)
                                                        setSecondTab(0)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'accounts'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex w-max flex-row gap-1 justify-center align-center items-center">
                                                                By Cloud
                                                                Account
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                View all
                                                                incidents and
                                                                drift events
                                                                across clouds,
                                                                accounts,
                                                                regions and
                                                                platforms.
                                                            </p>
                                                        </Flex>
                                                        {/* <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                1K{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex> */}
                                                    </Flex>
                                                </Card>
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setTab(1)
                                                        setSecondTab(1)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'resources'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex w-max flex-row gap-1 justify-center align-center items-center">
                                                                By Resource
                                                                Type
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                Get a summarized
                                                                view of posture
                                                                by Asset,
                                                                Entity, or
                                                                Resource Type.
                                                            </p>
                                                        </Flex>
                                                        {/* <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                2K{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex> */}
                                                    </Flex>
                                                </Card>{' '}
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setTab(1)

                                                        setSecondTab(2)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'controls'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex w-max flex-row gap-1 justify-center align-center items-center">
                                                                By Control
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                Get an overview
                                                                of conformance
                                                                by controls and
                                                                identify
                                                                problematic
                                                                ones.
                                                            </p>
                                                        </Flex>
                                                        {/* <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                2K{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex> */}
                                                    </Flex>
                                                </Card>{' '}
                                            </Grid>
                                            {/* <Divider className="mt-10 mb-10" /> */}

                                            <Grid
                                                // flexDirection="row"
                                                numItems={2}
                                                // justifyContent="center"
                                                // alignItems="center"
                                                className="mt-5 gap-8 flex-col w-full"
                                            >
                                                <Col numColSpan={2}>
                                                    <Title>All Events</Title>
                                                </Col>
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setSecondTab(0)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'findings'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex flex-row gap-1 justify-center align-center items-center">
                                                                All Incidents
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                View all
                                                                incidents and
                                                                drift events
                                                                across clouds,
                                                                accounts,
                                                                regions and
                                                                platforms.
                                                            </p>
                                                        </Flex>
                                                        {/* <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                1K{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex> */}
                                                    </Flex>
                                                </Card>
                                                <Card
                                                    className=" cursor-pointer flex justify-center items-center"
                                                    onClick={() => {
                                                        setSecondTab(1)
                                                        setShow(true)
                                                        setSelectedGroup(
                                                            'events'
                                                        )
                                                    }}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        className="h-100"
                                                    >
                                                        <Flex
                                                            flexDirection="col"
                                                            alignItems="start"
                                                            justifyContent="center"
                                                            className="gap-3 w-full"
                                                        >
                                                            <Title className="flex flex-row gap-1 justify-center align-center items-center">
                                                                Drift Events
                                                                <ChevronRightIcon className="w-[20px] mt-1" />
                                                            </Title>
                                                            <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                                See all
                                                                Compliance
                                                                Drifts over time
                                                                across
                                                                configured
                                                                Benchmarks
                                                            </p>
                                                        </Flex>
                                                        {/* <Flex
                                                            flexDirection="row"
                                                            justifyContent="end"
                                                            className="h-full"
                                                        >
                                                            <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                                500{' '}
                                                                <span className="font-semibold text-blue-600">
                                                                    +
                                                                </span>
                                                            </Title>
                                                        </Flex> */}
                                                    </Flex>
                                                </Card>{' '}
                                            </Grid>
                                        </main>
                                    </div>
                                </div>
                            </Flex>
                        </>
                    )}
                </>
            )}
        </>
    )
}
