import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Button,
    Card,
    Flex,
    Grid,
    Icon,
    Select,
    SelectItem,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    TextInput,
    Title,
} from '@tremor/react'
import {
    ChevronDoubleLeftIcon,
    ChevronDownIcon,
    ChevronRightIcon,
    ChevronUpIcon,
    CommandLineIcon,
    DocumentTextIcon,
    FunnelIcon,
    MagnifyingGlassIcon,
    PlayCircleIcon,
} from '@heroicons/react/24/outline'
import { Fragment, useEffect, useMemo, useState } from 'react' // eslint-disable-next-line import/no-extraneous-dependencies
import { highlight, languages } from 'prismjs' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/components/prism-sql' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/themes/prism.css'
import Editor from 'react-simple-code-editor'
import {
    CheckCircleIcon,
    ExclamationCircleIcon,
} from '@heroicons/react/24/solid'
import { Transition } from '@headlessui/react'
import { useAtom, useAtomValue } from 'jotai'
import {
    useInventoryApiV1QueryList,
    useInventoryApiV1QueryRunCreate,
    useInventoryApiV2AnalyticsCategoriesList,
} from '../../api/inventory.gen'
import Spinner from '../../components/Spinner'
import { getErrorMessage } from '../../types/apierror'
import { RenderObject } from '../../components/RenderObject'
import {
    PlatformEnginePkgInventoryApiRunQueryResponse,
    PlatformEnginePkgInventoryApiSmartQueryItem,
} from '../../api/api'
import { isDemoAtom, queryAtom, runQueryAtom } from '../../store'
import { snakeCaseToLabel } from '../../utilities/labelMaker'
import { numberDisplay } from '../../utilities/numericDisplay'
import TopHeader from '../../components/Layout/Header'
import AllQueries from './All Query'
import Query from './Query'
import { useParams, useSearchParams } from 'react-router-dom'
import { useURLParam } from '../../utilities/urlstate'
import { URLSearchParams } from 'url'


export default function Search() {
    const [tab,setTab] = useState<number>(0);
    // find query params for tabs
    // const [searchParams, setSearchParams] = useSearchParams()
    const [show, setShow] = useState<boolean>(true)

    // useEffect(() => {
    //    const tab_id = searchParams.get('tab_id')
    //     switch (tab_id) {
    //         case '1':
    //             setShow(true)
    //             setTab(1)
    //             break
    //         case '0':
    //             setShow(true)
    //             setTab(0)
    //             break
    //         default:
    //             setTab(0)
    //             break

    //     }
    // }, [searchParams])

    return (
        <>
            {/* <TopHeader /> */}
            {show ? (
                <>
                    {' '}
                    {/* <TabGroup
                        index={tab}
                        onIndexChange={(index) => {
                            setTab(index)
                        }}
                    >
                        <TabList>
                            <Tab
                                value={0}
                                onClick={() => {
                                    setTab(0)
                                }}
                            >
                                Query Library
                            </Tab>
                            <Tab
                                value={1}
                                onClick={() => {
                                    setTab(1)
                                }}
                            >
                                CloudQL
                            </Tab>
                        </TabList>
                        <TabPanels>
                            <TabPanel>
                                {tab == 0 && (
                                    <>
                                        <AllQueries setTab={setTab} />
                                    </>
                                )}
                            </TabPanel>
                            <TabPanel>
                                {tab == 1 && (
                                    <>
                                        <Query />
                                    </>
                                )}
                            </TabPanel>
                        </TabPanels>
                    </TabGroup> */}
                    <Query />
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
                                    Finder
                                </h1>
                                <p className="text-tremor-default text-tremor-content dark:text-dark-tremor-content">
                                    Find everything you need, from Code to
                                    Cloud.
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
                                                href="https://docs.opengovernance.io/oss/platform/discovery"
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
                                            Learn how to query and find any
                                            entity, asset across clouds and
                                            platforms
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
                                                        CloudQL
                                                        <ChevronRightIcon className="w-[20px] mt-1" />
                                                    </Title>
                                                    <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                        Search Across Clouds and
                                                        Platforms with SQL: Run
                                                        Queries and Export
                                                        Results
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
                                                        Query Library
                                                        <ChevronRightIcon className="w-[20px] mt-1" />
                                                    </Title>
                                                    <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                                        See over 2k+ search
                                                        queries sourced from
                                                        Steampipe's examples,
                                                        filter them by service
                                                        and explore by
                                                        framework.
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
        </>
    )
}
