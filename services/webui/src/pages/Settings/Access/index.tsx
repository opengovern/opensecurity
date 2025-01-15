import { Children, useEffect, useState } from 'react'
import { Flex } from '@tremor/react'
import {
    AdjustmentsVerticalIcon,
    BugAntIcon,
    Cog8ToothIcon,
    DocumentMagnifyingGlassIcon,
    DocumentTextIcon,
    FolderIcon,
    KeyIcon,
    UserIcon,
    UsersIcon,
} from '@heroicons/react/24/outline'
import { Link, useParams, useSearchParams } from 'react-router-dom'
import { useAtom, useAtomValue } from 'jotai'
import SettingsEntitlement from '../Entitlement'
import SettingsMembers from '../Members'
import SettingsWorkspaceAPIKeys from '../APIKeys'
import SettingsProfile from '../Profile'
import SettingsOrganization from '../Organization'
import { useWorkspaceApiV1WorkspaceCurrentList } from '../../../api/workspace.gen'

import { meAtom, tokenAtom } from '../../../store'
import SettingsJobs from '../Jobs'
import TopHeader from '../../../components/Layout/Header'
import SettingsParameters from '../Parameters'
import {
    useAuthApiV1MeList,
    useAuthApiV1WorkspaceRoleBindingsList,
} from '../../../api/auth.gen'
import { recordToMap } from '../../../utilities/record'
import { ApiRole } from '../../../api/api'

const navigation = [
    // {
    //     name: 'Workspace Settings',
    //     icon: DocumentTextIcon,
    //     role: ['admin', 'editor', 'viewer'],
    //     page: 'workspace-settings',
    //     children: [],
    // },
    {
        name: 'Authentication',
        page: 'authentication',
        icon: UsersIcon,
        role: ['admin'],
        children: [],
    },
    {
        name: 'API Keys',
        page: 'api-keys',
        icon: KeyIcon,
        role: ['admin'],
        children: [],
    },
    // {
    //     name: 'Customization',
    //     page: 'customization',
    //     icon: AdjustmentsVerticalIcon,
    //     role: ['admin'],
    //     children: [],
    // },
    // {
    //     name: 'Organization',
    //     icon: BuildingOfficeIcon,
    //     role: ['admin', 'editor', 'viewer'],
    //     children: [
    //         {
    //             name: 'Organization Info',
    //             page: 'org',
    //             role: ['admin', 'editor', 'viewer'],
    //         },
    //     ],
    // },
    // {
    //     name: 'Jobs',
    //     icon: BugAntIcon,
    //     page: 'jobs',
    //     role: ['admin', 'editor', 'viewer'],
    //     children: [],
    // },
    // {
    //     name: 'Metadata',
    //     icon: AdjustmentsVerticalIcon,
    //     page: 'parameters',
    //     role: ['admin'],
    //     children: [],
    //     // role: ['admin', 'editor', 'viewer'],
    //     // children: [
    //     //     {
    //     //         name: 'Parameters',
    //     //         page: 'parameters',
    //     //         role: ['admin'],
    //     //     },
    //     // ],
    // },
    // {
    //     name: 'Profile',
    //     icon: UserIcon,
    //     page: 'profile',
    //     role: ['admin', 'editor', 'viewer'],
    //     children: [],
    // },
    // {
    //     name: 'Sample data',
    //     icon: DocumentMagnifyingGlassIcon,
    //     page: 'sample-data',
    //     role: ['admin', 'editor', 'viewer'],
    //     children: [],
    // },
]

export default function SettingsAccess() {
    const [selectedTab, setSelectedTab] = useState(<SettingsEntitlement />)
    const [me, setMe] = useAtom(meAtom)

    const { response: curWorkspace, isLoading } =
        useWorkspaceApiV1WorkspaceCurrentList()
    const workspace = useParams<{ ws: string }>().ws

    const [searchParams, setSearchParams] = useSearchParams()
    const currentSubPage = searchParams.get('sp')

    useEffect(() => {
        switch (currentSubPage) {
            case 'workspace-seetings':
                setSelectedTab(<SettingsEntitlement />)
                break
            case 'authentication':
                setSelectedTab(<SettingsMembers />)
                break
            case 'api-keys':
                setSelectedTab(<SettingsWorkspaceAPIKeys />)
                break
            case 'org':
                setSelectedTab(<SettingsOrganization />)
                break
            case 'profile':
                setSelectedTab(<SettingsProfile />)
                break
            case 'jobs':
                setSelectedTab(<SettingsJobs />)
                break
            case 'parameters':
                setSelectedTab(<SettingsParameters />)
                break
            default:
                setSelectedTab(<SettingsMembers />)
                break
        }
    }, [currentSubPage])

    const getRole = () => {
        if (curWorkspace?.id) {
            if (me?.workspaceAccess) {
                return me?.workspaceAccess[curWorkspace.id] || 'viewer'
            }
        }
        return 'viewer'
    }

    return (
        <>
            {/* <TopHeader /> */}
            <Flex alignItems="start" justifyContent="start">
                <Flex flexDirection="col" alignItems="start" className="w-fit">
                    <nav className="w-56 text-sm">
                        <ul className="space-y-1.5">
                            {navigation.map((item: any) => {
                                if (
                                    !item.role.includes(getRole()) &&
                                    item.role.length > 0
                                ) {
                                    return null
                                }

                                if (item.children.length === 0) {
                                    return (
                                        <li key={item.name}>
                                            <Link
                                                to={`/settings/access?sp=${item.page}`}
                                                className={`${
                                                    item.page ===
                                                        currentSubPage ||
                                                    (!currentSubPage &&
                                                        item.page ===
                                                            'authentication')
                                                        ? 'bg-openg-100 dark:bg-openg-800  rounded-lg text-gray-800 dark:text-gray-100'
                                                        : 'text-gray-600 dark:text-gray-300'
                                                } group flex gap-x-3 pt-2 pb-0 px-2 -ml-2 font-medium`}
                                            >
                                                <Flex
                                                    justifyContent="start"
                                                    className="text-gray-800 dark:text-gray-100 font-semibold group gap-x-3 mb-2"
                                                >
                                                    {item.icon && (
                                                        <item.icon className="h-5 w-5 shrink-0" />
                                                    )}
                                                    {item.name}
                                                </Flex>
                                            </Link>
                                        </li>
                                    )
                                }

                                return (
                                    <li key={item.name}>
                                        <Flex
                                            justifyContent="start"
                                            className="text-gray-800 dark:text-gray-100 font-semibold group gap-x-3 mb-2"
                                        >
                                            {item.icon && (
                                                <item.icon className="h-5 w-5 shrink-0" />
                                            )}
                                            {item.name}
                                        </Flex>
                                        {item.children.map((child: any) => (
                                            <Link
                                                to={`/settings/about?sp=${child.page}`}
                                                className={`${
                                                    child.page ===
                                                        currentSubPage ||
                                                    (!currentSubPage &&
                                                        child.page ===
                                                            'summary')
                                                        ? 'bg-openg-100 dark:bg-openg-800  rounded-lg text-gray-800 dark:text-gray-100'
                                                        : 'text-gray-600 dark:text-gray-300'
                                                } group flex gap-x-3 py-2 px-8 font-medium`}
                                            >
                                                {child.name}
                                            </Link>
                                        ))}
                                    </li>
                                )
                            })}
                        </ul>
                    </nav>
                </Flex>
                <Flex
                    flexDirection="col"
                    justifyContent="center"
                    className="w-full"
                >
                    <Flex className="w-full h-full pl-6 max-w-7xl">
                        {selectedTab}
                    </Flex>
                </Flex>
            </Flex>
        </>
    )
}
