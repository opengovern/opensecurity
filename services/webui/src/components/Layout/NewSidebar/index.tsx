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
    UserIcon,
} from '@heroicons/react/24/outline'
import {
    RiAdminLine,
    RiChatSmileAiLine,
    RiChatSmileLine,
    RiFileWarningFill,
    RiHome2Line,
    RiLockStarFill,
    RiPuzzleLine,
    RiRefreshLine,
    RiRobot2Line,
    RiShieldCheckLine,
    RiSlideshowLine,
    RiTaskLine,
    RiTerminalBoxLine,
} from '@remixicon/react'
import { useAtom, useAtomValue, useSetAtom } from 'jotai'
import { Popover, Transition } from '@headlessui/react'
import { Fragment, useEffect, useState } from 'react'
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
import { SideNavigation } from '@cloudscape-design/components'



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

export default function NewSidebar({ currentPage }: ISidebar) {
    const { isAuthenticated, getAccessTokenSilently } = useAuth()
    console.log(currentPage)
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
        const show_compliance =
            window.__RUNTIME_CONFIG__.REACT_APP_SHOW_COMPLIANCE
        if (show_compliance === 'false') {
            return [
                {
                    name: 'CloudQL',
                    page: ['cloudql', 'cloudql'],
                    icon: MagnifyingGlassIcon,
                    isPreview: false,
                },

                {
                    name: 'Integration',
                    page: [
                        'integration/plugins',
                        'plugins/AWS',
                        'plugins/Azure',
                        'plugins/EntraID',
                    ],
                    icon: PuzzlePieceIcon,
                    isLoading: false,
                    // count: 0,

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
                icon: RiHome2Line,
                isPreview: false,
            },

            {
                name: 'CloudQL',
                page: 'cloudql',
                icon: RiTerminalBoxLine,
                isPreview: false,
            },
            {
                name: 'Compliance',
                icon: RiShieldCheckLine,
                page: 'compliance',

                isPreview: false,
                isLoading: false,
                count: undefined,
                error: false,
            },

            {
                name: 'All Incidents',
                icon: RiFileWarningFill,
                page: 'incidents',

                isPreview: false,
            },

            {
                name: 'Integration',
                page: 'integration/plugins',

                icon: RiPuzzleLine,
                isLoading: false,
                // count: 0,

                // count: numericDisplay(connectionCount?.count) || 0,
                error: undefined,
                isPreview: false,
            },

            {
                name: 'Jobs',
                page: 'jobs',
                icon: RiTaskLine,
                isPreview: false,
            },
            {
                name: 'Administration',
                page: 'administration',
                icon: RiAdminLine,
                isPreview: false,
            },
            {
                name: 'Agent AI',
                page: 'ai',
                icon: RiRobot2Line,
                isPreview: true,
            },

            {
                name: 'Automation',
                page: 'automation',
                icon: RiRefreshLine,
                isPreview: true,
            },
            // {
            //     name: 'Dashboards',
            //     page: [
            //         'dashboards',
            //         'dashboards/infrastructure',
            //         'dashboards/spend',
            //         'dashboards/infrastructure-cloud-accounts',
            //         'dashboards/infrastructure-metrics',
            //         'dashboards/spend-accounts',
            //         'dashboards/spend-metrics',
            //     ],
            //     icon: RiSlideshowLine,
            //     isPreview: true,
            // },
        ]
    }

    return (
        <>
            <div className="flex flex-col gap-2 p-2 mt-3 w-full">
                {/* logo */}
                <Flex className="ml-6">
                    <img
                        src={require('../../../icons/logo-light.png')}
                        className=""
                    />
                </Flex>
            </div>
            <SideNavigation
                    className='w-full custom-nav'
                // @ts-ignore
                items={navigation()?.map((item) => {
                    return {
                        href: `/${item.page}`,
                        type: 'link',
                        text: item.name,
                        
                        info: item?.isPreview ? (
                            <RiLockStarFill className='w-3'  />
                        ) : (
                            ''
                        ),
                    }
                })}
                activeHref={`${currentPage}`}
            />
        </>
    )
}
