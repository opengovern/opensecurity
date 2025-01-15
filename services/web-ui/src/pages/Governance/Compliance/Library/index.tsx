import { Card, Divider, Flex, Grid, Icon, Tab, TabGroup, TabList, TabPanel, TabPanels, Title } from '@tremor/react'
import TopHeader from '../../../../components/Layout/Header'
import { useState } from 'react'
import AllControls from '../All Controls';
import SettingsParameters from '../../../Settings/Parameters';
import { ChevronRightIcon, DocumentTextIcon } from '@heroicons/react/24/outline';
import { Tabs } from '@cloudscape-design/components';

export default function Library() {
    const [tab, setTab] = useState<number>(0)
    const [show,setShow] = useState<boolean>(false)
    return (
        <>
            <TopHeader />
            {show ? (
                <>
                    {' '}
                    <Tabs
                        activeTabId={tab.toString()}
                        onChange={({ detail }) => {
                            setTab(parseInt(detail.activeTabId))
                        }}
                        tabs={[
                            {
                                id: '0',
                                label: 'Controls',
                                content: <AllControls />,
                            },
                            {
                                id: '1',
                                label: 'Parameters',
                                content: <SettingsParameters />,
                            },
                        ]}
                    />
                    {/*       <AllBenchmarks /> */}
                </>
            ) : (
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
                                    Compliance Library
                                </h1>
                                <p className="text-tremor-default text-tremor-content dark:text-dark-tremor-content">
                                    View all Compliance Frameworks, Controls and
                                    Parameter controls.
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
                                            Learn how to audit, customize and
                                            create your own Compliance
                                            Frameworks
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
                                        numItems={2}
                                        // justifyContent="center"
                                        // alignItems="center"
                                        className="mt-5 gap-8 flex-col w-full"
                                    >
                                        {/* <Card
                                            className=" cursor-pointer flex justify-center items-center"
                                            onClick={() => {
                                                setTab(0)
                                                setShow(true)
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
                                                        Benchmarks
                                                        <ChevronRightIcon className="w-[20px] mt-1" />
                                                    </Title>
                                                    <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                        See the full list of
                                                        Benchmarks, filter them
                                                        by objective, and
                                                        explore the underlying
                                                        controls.
                                                    </p>
                                                </Flex>
                                                <Flex
                                                    flexDirection="row"
                                                    justifyContent="end"
                                                    className="h-full"
                                                >
                                                    <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                        50{' '}
                                                        <span className="font-semibold text-blue-600">
                                                            +
                                                        </span>
                                                    </Title>
                                                </Flex>
                                            </Flex>
                                        </Card> */}
                                        <Card
                                            className=" cursor-pointer flex justify-center items-center"
                                            onClick={() => {
                                                setTab(0)
                                                setShow(true)
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
                                                        Controls
                                                        <ChevronRightIcon className="w-[20px] mt-1" />
                                                    </Title>
                                                    <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                        See the full list of
                                                        controls, filter them by
                                                        service and explore by
                                                        framework.
                                                    </p>
                                                </Flex>
                                                <Flex
                                                    flexDirection="row"
                                                    justifyContent="end"
                                                    className="h-full"
                                                >
                                                    <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                        2.5K{' '}
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
                                                setTab(1)
                                                setShow(true)
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
                                                        Parameters
                                                        <ChevronRightIcon className="w-[20px] mt-1" />
                                                    </Title>
                                                    <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                        Manage your variables
                                                        and parameters used
                                                        across controls &
                                                        compliance.
                                                    </p>
                                                </Flex>
                                                <Flex
                                                    flexDirection="row"
                                                    justifyContent="end"
                                                    className="h-full"
                                                >
                                                    <Title className=" font-bold  border-solid w-fit h-full  border-l-2 border-black pl-2 h-full">
                                                        100{' '}
                                                        <span className="font-semibold text-blue-600">
                                                            +
                                                        </span>
                                                    </Title>
                                                </Flex>
                                            </Flex>
                                        </Card>
                                    </Grid>
                                </main>
                            </div>
                        </div>
                    </Flex>
                </>
            )}
        </>
    )
}

