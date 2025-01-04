import { Badge, Card, Flex, Text } from '@tremor/react'
import { Link, useNavigate } from 'react-router-dom'
import {
    BanknotesIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    Cog6ToothIcon,
    CubeIcon,
    DocumentChartBarIcon,
    ExclamationCircleIcon,
    Squares2X2Icon,
    MagnifyingGlassIcon,
    PuzzlePieceIcon,
    RectangleStackIcon,
    ShieldCheckIcon,
    ClipboardDocumentCheckIcon,
    DocumentMagnifyingGlassIcon,
    ArrowUpCircleIcon,
    PresentationChartBarIcon,
    CubeTransparentIcon,
    BoltIcon,
    ArrowUpIcon,
    ChevronDoubleUpIcon,
    CalendarDateRangeIcon,
    CommandLineIcon,
} from '@heroicons/react/24/outline'
import { useAtom, useAtomValue, useSetAtom } from 'jotai'
import { Popover, Transition } from '@headlessui/react'
import { Fragment, useEffect } from 'react'
import { previewAtom, sideBarCollapsedAtom } from '../../../store'
import { OpenGovernance, OpenGovernanceBig } from '../../../icons/icons'
import Utilities from './Utilities'
import {
    useInventoryApiV2AnalyticsCountList,
    useInventoryApiV2AnalyticsSpendCountList,
} from '../../../api/inventory.gen'
import { useIntegrationApiV1ConnectionsCountList } from '../../../api/integration.gen'
import { numericDisplay } from '../../../utilities/numericDisplay'
import AnimatedAccordion from '../../AnimatedAccordion'
import { setAuthHeader } from '../../../api/ApiConfig'
import {
    searchAtom,
    oldUrlAtom,
    nextUrlAtom,
} from '../../../utilities/urlstate'
import { useAuth } from '../../../utilities/auth'

const badgeStyle = {
    color: '#fff',
    borderRadius: '8px',
    backgroundColor: '#15395F80',
}

interface ISidebar {
   
    currentPage: string
}

interface ISidebarItem {
    name: string
    page: string | string[]
    icon?: any
    isLoading?: boolean
    count?: number | string
    error?: any
    isPreview?: boolean
    children?: ISidebarItem[]
    selected?: string
}

export default function Sidebar({  currentPage }: ISidebar) {
    const navigate = useNavigate()
    const { isAuthenticated, getAccessTokenSilently } = useAuth()
    const [collapsed, setCollapsed] = useAtom(sideBarCollapsedAtom)
    const preview = useAtomValue(previewAtom)
    

    const searchParams = useAtomValue(searchAtom)
    const setOldUrl = useSetAtom(oldUrlAtom)
    const oldUrl = useAtomValue(oldUrlAtom)
    const nextUrl = useAtomValue(nextUrlAtom)
    const setNextUrl = useSetAtom(nextUrlAtom)

    const isCurrentPage = (page: string | string[] | undefined): boolean => {
        if (Array.isArray(page)) {
            return page.map((p) => isCurrentPage(p)).includes(true)
        }

        if (page?.includes('?')) {
            const pageParams = new URLSearchParams(
                page?.substring(page?.indexOf('?'))
            )

            const locUrl = new URL(window.location.href)
            const locParams = new URLSearchParams(locUrl.search)

            let ok = true
            pageParams.forEach((value, key) => {
                if (locParams.get(key) !== value) {
                    ok = false
                }
            })
            return currentPage === page?.substring(0, page?.indexOf('?')) && ok
        }
        // if(page?.includes(":")){
        //     return currentPage.split("/")[0] === page?.substring(0, page?.indexOf(':')).split("/")[0]
        // }
        if (page == '') {
            return currentPage == ''
        }
        // @ts-ignore

        return currentPage.includes(page) && page !==''
    }
    const findPage = (page: string | string[], item: ISidebarItem): string => {
        if (Array.isArray(item.page)) {
            if (
                oldUrl.includes(item.page[0]) ||
                nextUrl.includes(item.page[0])
            ) {
                if (Array.isArray(page)) {
                    return `/${page[0]}?${searchParams}`
                } else {
                    return `/${page}?${searchParams}`
                }
            } else {
                if (Array.isArray(page)) {
                    return `/${page[0]}`
                } else {
                    return `/${page}`
                }
            }
        } else {
            if (oldUrl.includes(item.page) || nextUrl.includes(item.page)) {
                return `/${page}?${searchParams}`
            } else {
                return `/${page}`
            }
        }

        // if (page.includes('?')) {
        //     return `/${page}`
        // }
        // // if (searchParams) {
        // //     return `/${page}?${searchParams}`
        // // }
        // if (page.includes('/')) {
        //     return `/${page}`
        // }
        // return `/${page}?${searchParams}`
    }
    useEffect(() => {
        if (isAuthenticated) {
            getAccessTokenSilently()
                .then((accessToken) => {
                    setAuthHeader(accessToken)
                    // sendSpend()
                    // sendAssets()
                    // sendFindings()
                    // sendConnections()
                    // fetchDashboardToken()
                })
                .catch((e) => {
                    console.log('====> failed to get token due to', e)
                })
        }
    }, [isAuthenticated])

    const navigation: () => ISidebarItem[] = () => {
        if (
            // connectionsIsExecuted &&
            // !connectionsIsLoading &&
            // (connectionCount?.count || 0) === 0
            false
        ) {
            if (
                !currentPage.includes('integrations') &&
                currentPage !== 'administration'
            ) {
                navigate(`/integrations`)
            }
            return [
                {
                    name: 'Integrations',
                    page: 'integrations',
                    icon: PuzzlePieceIcon,
                    // isLoading: connectionsIsLoading,
                    isLoading: false,
                    count:  0,

                    // count: numericDisplay(connectionCount?.count) || 0,
                    error: undefined,
                    isPreview: false,
                },
                {
                    name: 'Administration',
                    page: ['administration'],
                    icon: Cog6ToothIcon,
                    isPreview: false,
                },
            ]
        }
        return [
            {
                name: 'Overview',
                page: '',
                icon: Squares2X2Icon,
                isPreview: false,
            },

            // {
            //     name: 'Cloud Dashboards',
            //     page: [
            //         'infrastructure',
            //         'infrastructure-cloud-accounts',
            //         'infrastructure-metrics',
            //         'spend',
            //         'spend-accounts',
            //         'spend-metrics',
            //     ],
            //     icon: CubeIcon,
            //     children: [
            //         {
            //             name: 'Infrastructure Summary',
            //             page: 'infrastructure',
            //             isPreview: false,
            //             isLoading: false,
            //             count: undefined,
            //             error: false,
            //         },
            //         {
            //             name: 'Spend Summary',
            //             page: 'spend',
            //             isPreview: false,
            //             isLoading: false,
            //             count: undefined,
            //             error: false,
            //         },
            //         {
            //             name: 'Infrastructure - Cloud Accounts',
            //             page: 'infrastructure-cloud-accounts',
            //             isPreview: false,
            //             isLoading: assetsIsLoading,
            //             count: numericDisplay(assetCount?.connectionCount) || 0,
            //             error: assetCountErr,
            //         },
            //         {
            //             name: 'Infrastructure - Cloud Services',
            //             page: 'infrastructure-metrics',
            //             isPreview: false,
            //             isLoading: assetsIsLoading,
            //             count: numericDisplay(assetCount?.metricCount) || 0,
            //             error: assetCountErr,
            //         },

            //         {
            //             name: 'Spend - Cloud Accounts',
            //             page: 'spend-accounts',
            //             isPreview: false,
            //             isLoading: spendCountIsLoading,
            //             count: numericDisplay(spendCount?.connectionCount) || 0,
            //             error: spendCountErr,
            //         },
            //         {
            //             name: 'Spend - Cloud Services',
            //             page: 'spend-metrics',
            //             isPreview: false,
            //             isLoading: spendCountIsLoading,
            //             count: numericDisplay(spendCount?.metricCount) || 0,
            //             error: spendCountErr,
            //         },
            //     ],
            //     isPreview: false,
            // },
            // {
            //     name: 'Spend',
            //     page: ['spend', 'spend-accounts', 'spend-metrics'],
            //     icon: BanknotesIcon,
            //     children: [
            //         {
            //             name: 'Summary',
            //             page: 'spend',
            //             isPreview: false,
            //             isLoading: false,
            //             count: undefined,
            //             error: false,
            //         },
            //         {
            //             name: 'Cloud Accounts',
            //             page: 'spend-accounts',
            //             isPreview: false,
            //             isLoading: spendCountIsLoading,
            //             count: numericDisplay(spendCount?.connectionCount) || 0,
            //             error: spendCountErr,
            //         },
            //         {
            //             name: 'Services',
            //             page: 'spend-metrics',
            //             isPreview: false,
            //             isLoading: spendCountIsLoading,
            //             count: numericDisplay(spendCount?.metricCount) || 0,
            //             error: spendCountErr,
            //         },
            //     ],
            //     isPreview: false,
            // },
            {
                name: 'Compliance',
                icon: ShieldCheckIcon,
                page: [
                    'compliance',
                    'compliance/:benchmarkId',
                    'compliance/controls',
                    'compliance/benchmarks',
                ],
                isPreview: false,
                isLoading: false,
                count: undefined,
                error: false,
                // children: [
                // {
                //     name: 'Frameworks',
                //     page: ['compliance','compliance/:benchmarkId'],
                //     isPreview: false,
                //     isLoading: false,
                //     count: undefined,
                //     error: false,
                // },
                // {
                //     name: 'SCORE',
                //     page: 'score',
                //     isPreview: false,
                //     isLoading: false,
                //     count: undefined,
                //     error: false,
                // },
                // {
                //     name: ' Library',
                //     page: 'compliance/library',
                //     isPreview: false,
                //     isLoading: false,
                //     count: undefined,
                //     error: false,
                // },
                // {
                //     name: ' Benchmarks',
                //     page: 'compliance/benchmarks',
                //     isPreview: false,
                //     isLoading: false,
                //     count: undefined,
                //     error: false,
                // },
                // {
                //     name: ' Controls',
                //     page: 'compliance/controls',
                //     isPreview: false,
                //     isLoading: false,
                //     count: undefined,
                //     error: false,
                // },
                // ],
            },
            // {
            //     name: 'SRE',
            //     page: ['score', 'score/categories'],
            //     icon: ChevronDoubleUpIcon,
            //     isPreview: false,
            // },

            {
                name: 'All Incidents',
                icon: DocumentMagnifyingGlassIcon,
                page: [
                    'incidents',
                    'incidents/summary',
                    // 'incidents/drift-events',
                ],
                isPreview: false,
                // children: [
                //     {
                //         name: 'Summary',
                //         page: 'incidents/summary',
                //     },
                //     {
                //         name: 'Incidents',
                //         page: 'incidents',
                //         isLoading: findingsIsLoading,
                //         // count: numericDisplay(findingsCount?.count) || 0,
                //         error: findingsErr,
                //     },

                //     // {
                //     //     name: 'Resource Summary',
                //     //     page: 'incidents/resource-summary',
                //     // },
                //     // {
                //     //     name: 'Drift Events',
                //     //     page: 'incidents/drift-events',
                //     // },
                //     // {
                //     //     name: 'Account Posture',
                //     //     page: 'incidents/account-posture',
                //     // },
                //     // {
                //     //     name: 'Control Summary',
                //     //     page: 'incidents/control-summary',
                //     // },
                // ],
            },

            // {
            //     name: 'Security',
            //     icon: ShieldCheckIcon,
            //     page: 'security-overview',
            //     isPreview: false,
            //     isLoading: false,
            //     count: undefined,
            //     error: false,
            // },

            {
                name: 'CloudQL',
                page: ['cloudql', 'cloudql'],
                icon: MagnifyingGlassIcon,
                isPreview: false,
                // children: [
                //     {
                //         name: ' Query',
                //         page: 'query',
                //         isPreview: false,
                //         isLoading: false,
                //         count: undefined,
                //         error: false,
                //     },
                //     {
                //         name: 'All Queries',
                //         page: 'query/all-queries',
                //         isPreview: false,
                //         isLoading: false,
                //         count: undefined,
                //         error: false,
                //     },
                // ],
            },
            // {
            //     name: 'Tasks',
            //     page: ['tasks', 'tasks'],
            //     icon: CommandLineIcon,
            //     isPreview: false,
               
            // },

            {
                name: 'Integrations',
                page: [
                    'integrations',
                    'integrations/AWS',
                    'integrations/Azure',
                    'integrations/EntraID',
                ],
                icon: PuzzlePieceIcon,
                // isLoading: connectionsIsLoading,
                isLoading: false,
                // count: 0,

                // count: numericDisplay(connectionCount?.count) || 0,
                error: undefined,
                isPreview: false,
            },
            // {
            //     name: 'Workload Optimizer',
            //     page: 'workload-optimizer',
            //     icon: ArrowUpCircleIcon,
            //     isPreview: true,
            //     isLoading: false,
            //     count: undefined,
            //     // error: spendCountErr,
            // },
            // {
            //     name: 'Stacks',
            //     page: 'stacks',
            //     icon: RectangleStackIcon,
            //     isPreview: true,
            // },

            {
                name: 'Jobs',
                page: 'jobs',
                icon: CalendarDateRangeIcon,
                isPreview: false,
            },
            {
                name: 'Administration',
                page: ['administration'],
                icon: Cog6ToothIcon,
                isPreview: false,
            },
            {
                name: 'Automation',
                page: 'automation',
                icon: BoltIcon,
                isPreview: true,
            },
            {
                name: 'Dashboards',
                page: [
                    'dashboard',
                    'dashboard/infrastructure',
                    'dashboard/spend',
                    'dashboard/infrastructure-cloud-accounts',
                    'dashboard/infrastructure-metrics',
                    'dashboard/spend-accounts',
                    'dashboard/spend-metrics',
                ],
                icon: PresentationChartBarIcon,
                isPreview: true,
            },
        ]
    }

    return (
        <Flex
            flexDirection="col"
            alignItems="start"
            className={`z-50 !max-h-screen h-full  ${
                collapsed ? 'w-20' : ' 2xl:w-60 sm:w-44'
            }
             pt-4 bg-openg-950  dark:bg-gray-950 relative border-r border-r-gray-700`}
        >
            <Flex
                flexDirection="col"
                justifyContent="start"
                className={`h-full ${collapsed ? 'w-full' : 'w-full'}`}
            >
                <a
                    className={`cursor-pointer ${
                        collapsed ? '' : 'pl-4 mr-6'
                    } w-full`}
                    href={`/`}
                >
                    <Flex
                        justifyContent={collapsed ? 'center' : 'between'}
                        className={`pb-[17px] pt-[6px]  `}
                    >
                        {collapsed ? <OpenGovernance /> : <OpenGovernanceBig />}

                        {!collapsed && (
                            <ChevronLeftIcon
                                className="h-6 w-6 text-gray-400 cursor-pointer "
                                onClick={() => {
                                    setCollapsed(true)
                                    localStorage.collapse = 'true'
                                }}
                            />
                        )}
                    </Flex>
                </a>

                {/* <Flex
                    justifyContent="center"
                    className={collapsed ? 'p-0' : 'p-2'}
                >
                    <Workspaces isCollapsed={collapsed} />
                </Flex> */}
                <Flex
                    flexDirection="col"
                    justifyContent="between"
                    className="h-full max-h-full"
                >
                    <div
                        className={`w-full p-2  ${
                            collapsed
                                ? 'flex justify-start flex-col items-center'
                                : 'overflow-y-scroll'
                        } h-full no-scrollbar`}
                        style={{ maxHeight: 'calc(100vh - 130px)' }}
                    >
                        {!collapsed && (
                            <Text className="my-2 !text-xs">OVERVIEW</Text>
                        )}
                        {collapsed && (
                            <ChevronRightIcon
                                className="m-2 h-6 text-gray-400 cursor-pointer"
                                onClick={() => {
                                    setCollapsed(false)
                                    localStorage.collapse = 'false'
                                }}
                            />
                        )}
                        {navigation()
                            .filter((item) =>
                                preview === 'true'
                                    ? item
                                    : String(item.isPreview) === String(preview)
                            )
                            .map((item) =>
                                // eslint-disable-next-line no-nested-ternary
                                item.children && !collapsed ? (
                                    <div className="w-full my-1">
                                        <AnimatedAccordion
                                            defaultOpen={
                                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                                // @ts-ignore
                                                // item.children.filter(
                                                //     (c: any) =>
                                                //         isCurrentPage(c.page) ||
                                                //         isCurrentPage(
                                                //             c.selected
                                                //         )
                                                // ).length > 0
                                                item.name !== 'Settings'
                                            }
                                            header={
                                                <div
                                                    className={`w-full text-gray-50 ${
                                                        collapsed
                                                            ? 'px-2'
                                                            : 'px-6'
                                                    } py-2`}
                                                >
                                                    <Flex
                                                        justifyContent="start"
                                                        className="h-full gap-2.5"
                                                    >
                                                        {item.isPreview ===
                                                        true ? (
                                                            <item.icon
                                                                className={`h-6 w-6 stroke-2 ${
                                                                    collapsed &&
                                                                    isCurrentPage(
                                                                        item.page
                                                                    )
                                                                        ? 'text-orange-200'
                                                                        : 'text-orange-400'
                                                                }`}
                                                            />
                                                        ) : (
                                                            <item.icon
                                                                className={`h-6 w-6 stroke-2 ${
                                                                    collapsed &&
                                                                    isCurrentPage(
                                                                        item.page
                                                                    )
                                                                        ? 'text-gray-200'
                                                                        : 'text-gray-400'
                                                                }`}
                                                            />
                                                        )}

                                                        <Text className="text-inherit">
                                                            {item.name}
                                                        </Text>
                                                        {item.isPreview &&
                                                            !collapsed && (
                                                                <Badge
                                                                    className="absolute right-2 top-1.5"
                                                                    style={
                                                                        badgeStyle
                                                                    }
                                                                >
                                                                    Preview
                                                                </Badge>
                                                            )}
                                                    </Flex>
                                                </div>
                                            }
                                        >
                                            {item.children.map((i) => (
                                                <Link
                                                    onClick={() => {
                                                        setOldUrl(
                                                            window.location.href
                                                        )
                                                        setNextUrl(
                                                            findPage(
                                                                i.page,
                                                                item
                                                            )
                                                        )
                                                    }}
                                                    to={findPage(i.page, item)}
                                                    className={`my-0.5 py-2 flex rounded-md relative
                                                    ${
                                                        isCurrentPage(i.page)
                                                            ? 'bg-openg-500 text-gray-200 font-semibold'
                                                            : 'text-gray-50 hover:bg-openg-800'
                                                    }`}
                                                >
                                                    <Text className="ml-[54px] text-inherit">
                                                        {i.name}
                                                    </Text>
                                                    {i.count && (
                                                        <Badge
                                                            className="absolute right-2 top-1.5"
                                                            style={badgeStyle}
                                                        >
                                                            {/* eslint-disable-next-line no-nested-ternary */}
                                                            {i.isLoading ? (
                                                                <div className="animate-pulse h-1 w-4 my-2 bg-gray-700 rounded-md" />
                                                            ) : i.error ? (
                                                                <ExclamationCircleIcon className="h-6" />
                                                            ) : (
                                                                i.count
                                                            )}
                                                        </Badge>
                                                    )}
                                                    {i.isPreview && (
                                                        <div className="absolute right-2 top-1.5">
                                                            <ArrowUpCircleIcon
                                                                height={18}
                                                                color="orange"
                                                            />
                                                        </div>
                                                    )}
                                                </Link>
                                            ))}
                                        </AnimatedAccordion>
                                    </div>
                                ) : item.children && collapsed ? (
                                    <div className="w-full my-1 ">
                                        <Popover className="relative z-50 border-0 w-full h-[36px]">
                                            <div
                                                className={`group relative ${
                                                    collapsed
                                                        ? 'justify-center flex'
                                                        : ''
                                                }`}
                                            >
                                                <Popover.Button id={item.name}>
                                                    <div
                                                        className={`w-full rounded-md p-2 
                                                    ${
                                                        isCurrentPage(item.page)
                                                            ? 'bg-openg-500 text-gray-200 font-semibold'
                                                            : 'text-gray-50 hover:bg-openg-800'
                                                    }`}
                                                    >
                                                        <item.icon
                                                            className={`h-6 w-6 stroke-2 ${
                                                                isCurrentPage(
                                                                    item.page
                                                                )
                                                                    ? 'text-gray-200'
                                                                    : 'text-gray-400'
                                                            }`}
                                                        />
                                                    </div>
                                                </Popover.Button>
                                                <div
                                                    className="absolute z-50 scale-0 transition-all rounded p-2 shadow-md bg-openg-950 group-hover:scale-100"
                                                    style={{
                                                        left: '50px',
                                                        top: 0,
                                                    }}
                                                >
                                                    <Text className="text-white">
                                                        {item.name}
                                                    </Text>
                                                </div>
                                            </div>
                                            <Transition
                                                as={Fragment}
                                                enter="transition ease-out duration-200"
                                                enterFrom="opacity-0 translate-y-1"
                                                enterTo="opacity-100 translate-y-0"
                                                leave="transition ease-in duration-150"
                                                leaveFrom="opacity-100 translate-y-0"
                                                leaveTo="opacity-0 translate-y-1"
                                            >
                                                <Popover.Panel className="absolute left-[163px] top-[1px] z-40 flex w-screen max-w-max -translate-x-1/2 px-4">
                                                    <Card className="z-50 rounded p-2 shadow-md !ring-gray-600 w-56 bg-openg-950">
                                                        <Text className="ml-1 mb-3 text-white">
                                                            {item.name}
                                                        </Text>
                                                        {item.children.map(
                                                            (i) => (
                                                                <Link
                                                                    onClick={() => {
                                                                        setOldUrl(
                                                                            window
                                                                                .location
                                                                                .href
                                                                        )
                                                                        setNextUrl(
                                                                            findPage(
                                                                                i.page,
                                                                                item
                                                                            )
                                                                        )
                                                                    }}
                                                                    to={findPage(
                                                                        i.page,
                                                                        item
                                                                    )}
                                                                    className={`my-0.5 py-2 px-4 flex  rounded-md relative 
                                                    ${
                                                        isCurrentPage(i.page)
                                                            ? 'bg-openg-500 text-gray-200 font-semibold'
                                                            : 'text-gray-50 hover:bg-openg-800'
                                                    }`}
                                                                >
                                                                    <Text className="text-inherit">
                                                                        {i.name}
                                                                    </Text>
                                                                    {i.count && (
                                                                        <Badge
                                                                            className="absolute right-2 top-1.5"
                                                                            style={
                                                                                badgeStyle
                                                                            }
                                                                        >
                                                                            {/* eslint-disable-next-line no-nested-ternary */}
                                                                            {i.isLoading ? (
                                                                                <div className="animate-pulse h-1 w-4 my-2 bg-gray-700 rounded-md" />
                                                                            ) : i.error ? (
                                                                                <ExclamationCircleIcon className="h-6" />
                                                                            ) : (
                                                                                i.count
                                                                            )}
                                                                        </Badge>
                                                                    )}
                                                                    {i.isPreview && (
                                                                        <Badge
                                                                            className="absolute right-2 top-1.5"
                                                                            style={
                                                                                badgeStyle
                                                                            }
                                                                        >
                                                                            Preview
                                                                        </Badge>
                                                                    )}
                                                                </Link>
                                                            )
                                                        )}
                                                    </Card>
                                                </Popover.Panel>
                                            </Transition>
                                        </Popover>
                                    </div>
                                ) : (
                                    <div className="w-full my-1">
                                        {/* eslint-disable-next-line jsx-a11y/anchor-is-valid */}
                                        <Link
                                            onClick={() => {
                                                setOldUrl(window.location.href)
                                                setNextUrl(
                                                    findPage(item.page, item)
                                                )
                                            }}
                                            to={findPage(item.page, item)}
                                            className={`w-full relative px-6 py-2 flex items-center gap-2.5 rounded-md ${
                                                collapsed
                                                    ? 'justify-center'
                                                    : ''
                                            }
                                                        ${
                                                            isCurrentPage(
                                                                item.page
                                                            ) ||
                                                            (collapsed &&
                                                                isCurrentPage(
                                                                    item.page
                                                                ))
                                                                ? 'bg-openg-500 text-gray-200 font-semibold'
                                                                : 'text-gray-50 hover:bg-openg-800'
                                                        }
                                                        ${
                                                            collapsed
                                                                ? '!p-2'
                                                                : ''
                                                        }`}
                                        >
                                            <div className="group relative">
                                                {item.isPreview === true ? (
                                                    <item.icon
                                                        className={`h-6 w-6 stroke-2 ${
                                                            isCurrentPage(
                                                                item.page
                                                            ) ||
                                                            (collapsed &&
                                                                isCurrentPage(
                                                                    item.page
                                                                ))
                                                                ? 'text-orange-200'
                                                                : 'text-orange-400'
                                                        }`}
                                                    />
                                                ) : (
                                                    <item.icon
                                                        className={`h-6 w-6 stroke-2 ${
                                                            isCurrentPage(
                                                                item.page
                                                            ) ||
                                                            (collapsed &&
                                                                isCurrentPage(
                                                                    item.page
                                                                ))
                                                                ? 'text-gray-200'
                                                                : 'text-gray-400'
                                                        }`}
                                                    />
                                                )}
                                                {collapsed && (
                                                    <div
                                                        className="absolute z-50 scale-0 transition-all rounded p-2 shadow-md bg-openg-950 group-hover:scale-100"
                                                        style={{
                                                            left: '43px',
                                                            top: '-8px',
                                                        }}
                                                    >
                                                        <Text className="text-white">
                                                            {item.name}
                                                        </Text>
                                                    </div>
                                                )}
                                            </div>
                                            {!collapsed && (
                                                <Text className="text-inherit">
                                                    {item.name}
                                                </Text>
                                            )}
                                            {item.count && !collapsed && (
                                                <Badge
                                                    className="absolute right-2 top-1.5"
                                                    style={badgeStyle}
                                                >
                                                    {/* eslint-disable-next-line no-nested-ternary */}
                                                    {item.isLoading ? (
                                                        <div className="animate-pulse h-1 w-4 my-2 bg-gray-700 rounded-md" />
                                                    ) : item.error ? (
                                                        <ExclamationCircleIcon className="h-5" />
                                                    ) : (
                                                        item.count
                                                    )}
                                                </Badge>
                                            )}
                                        </Link>
                                    </div>
                                )
                            )}
                    </div>
                </Flex>
                <Utilities isCollapsed={collapsed} />
            </Flex>
        </Flex>
    )
}
