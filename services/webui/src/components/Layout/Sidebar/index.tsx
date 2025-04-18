import { Badge, Card, Flex, Text } from '@tremor/react'
import { Link, useNavigate } from 'react-router-dom'

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

export default function Sidebar({ currentPage }: ISidebar) {
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
                    page: 'cloudql',
                    icon: RiTerminalBoxLine,
                    isPreview: false,
                },

                {
                    name: 'Integration',
                    page: 'integration/plugins',
                    icon: RiPuzzleLine,
                    isLoading: false,
                    error: undefined,
                    isPreview: false,
                },

                {
                    name: 'Administration',
                    page: 'administration',
                    icon: RiAdminLine,
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
                children: [
                    {
                        name: 'Overview',
                        page: '',
                        icon: RiHome2Line,
                        isPreview: false,
                    },
                ],
            },
            {
                name: 'Find',
                page: 'cloudql',
                icon: RiTerminalBoxLine,
                isPreview: false,
                children: [
                    {
                        name: 'CloudQL',
                        page: 'cloudql',
                        icon: RiTerminalBoxLine,
                        isPreview: false,
                    },
                    {
                        name: 'OPS Agents',
                        page: 'ai',
                        icon: RiRobot2Line,
                        isPreview: false,
                    },
                ],
            },

            {
                name: 'Compliance',
                icon: RiShieldCheckLine,
                page: 'compliance',
                children: [
                    {
                        name: 'Frameworks',
                        page: 'compliance/frameworks',
                        icon: RiShieldCheckLine,
                        isPreview: false,
                    },
                    {
                        name: 'Controls',
                        page: 'compliance/controls',
                        icon: RiShieldCheckLine,
                        isPreview: false,
                    },
                    // {
                    //     name: 'Policies',
                    //     page: 'compliance/policies',
                    //     icon: RiShieldCheckLine,
                    //     isPreview: false,
                    // },
                    {
                        name: 'Parameters',
                        page: 'compliance/parameters',
                        icon: RiShieldCheckLine,
                        isPreview: false,
                    },
                    {
                        name: 'Compliance Jobs',
                        page: 'compliance/jobs',
                        icon: RiTaskLine,
                        isPreview: false,
                    },
                ],

                isPreview: false,
                isLoading: false,
                count: undefined,
                error: false,
            },
            // {
            //     name: 'Tasks',
            //     icon: RiShieldCheckLine,
            //     page: 'task',
            //     children: [

            //     ],

            //     isPreview: false,
            //     isLoading: false,
            //     count: undefined,
            //     error: false,
            // },

            {
                name: 'Incidents',
                icon: RiFileWarningFill,
                page: 'incidents',
                isPreview: false,
                children: [
                    {
                        name: 'All Incidents',
                        icon: RiFileWarningFill,
                        page: 'incidents',
                        isPreview: false,
                    },
                    // {
                    //     name: 'Control Summary',
                    //     icon: RiFileWarningFill,
                    //     page: 'incidents/controls',
                    //     isPreview: false,
                    // },
                    // {
                    //     name: 'Resource Incident',
                    //     icon: RiFileWarningFill,
                    //     page: 'incidents/resources',
                    //     isPreview: false,
                    // },
                ],
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
                children: [
                    {
                        name: 'Plugins',
                        page: 'integration/plugins',

                        icon: RiPuzzleLine,
                        isLoading: false,
                        // count: 0,

                        // count: numericDisplay(connectionCount?.count) || 0,
                        error: undefined,
                        isPreview: false,
                    },
                    {
                        name: 'Task',
                        page: 'tasks',
                        icon: RiShieldCheckLine,
                        isPreview: false,
                    },

                    {
                        name: 'Discovery Jobs',
                        page: 'integration/jobs',
                        icon: RiTaskLine,
                        isPreview: false,
                    },
                ],
            },

            // {
            //     name: 'Jobs',
            //     page: 'jobs',
            //     icon: RiTaskLine,
            //     isPreview: false,
            // },
            {
                name: 'Administration',
                page: 'administration',
                icon: RiAdminLine,
                isPreview: false,
                children: [
                    {
                        name: 'Settings',
                        page: 'administration/settings',
                        icon: RiAdminLine,
                        isPreview: false,
                    },
                    {
                        name: 'Access',
                        page: 'administration/access',
                        icon: RiAdminLine,
                        isPreview: false,
                    },
                ],
            },

            {
                name: 'Automation',
                page: 'automation',
                icon: RiRefreshLine,
                isPreview: true,
            },
        ]
    }

    return (
        <>
            <div className="flex flex-col gap-2 p-2 mt-3 w-full">
                {/* logo */}
                <Flex className="">
                    <img
                        src={require('../../../icons/logo-light.png')}
                        className="ml-4"
                    />
                </Flex>
            </div>
            <SideNavigation
                className="w-full custom-nav"
                // @ts-ignore
                items={navigation()?.map((item) => {
                    return {
                        href: `/${item.page}`,
                        type: item?.children ? 'section' : 'link',
                        text: item.name,

                        info: item?.isPreview ? (
                            <RiLockStarFill className="w-3" />
                        ) : (
                            ''
                        ),
                        items: item?.children
                            ? item?.children.map((child) => {
                                  return {
                                      href: `/${child.page}`,
                                      type: 'link',
                                      text: child.name,

                                    //   info: child?.isPreview ? (
                                    //       <RiLockStarFill className="w-3" />
                                    //   ) : (
                                    //       ''
                                    //   )
                                      
                                  }
                              })
                            : [],
                    }
                })}
                activeHref={`${currentPage}`}
            />
        </>
    )
}
