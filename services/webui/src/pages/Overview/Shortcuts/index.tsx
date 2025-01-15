import {
    ArrowTopRightOnSquareIcon,
    BanknotesIcon,
    ChevronRightIcon,
    CubeIcon,
    CursorArrowRaysIcon,
    PuzzlePieceIcon,
    ShieldCheckIcon,
} from '@heroicons/react/24/outline'
import { Card, Flex, Grid, Icon, Text, Title } from '@tremor/react'
import { useNavigate, useParams } from 'react-router-dom'
import Check from '../../../icons/Check.svg'
import User from '../../../icons/User.svg'
import Dollar from '../../../icons/Dollar.svg'
import Cable from '../../../icons/Cable.svg'
import Cube from '../../../icons/Cube.svg'
import { link } from 'fs'
import { useEffect, useState } from 'react'
import Evaluate from '../../Governance/Compliance/NewBenchmarkSummary/Evaluate'
import { title } from 'process'
import { Modal } from '@cloudscape-design/components'
import MemberInvite from '../../Settings/Members/MemberInvite'

const navList = [
    {
        title: 'CloudQL',
        description: 'See all workloads - from code to cloud',
        icon: Cube,
        link: 'cloudql?tab_id=0',
        new: true,
    },
    {
        title: 'Audit',
        description: 'Review and run compliance checks',
        icon: Check,
        link: 'compliance',
        new: true,
    },
    {
        title: 'Invite',
        description: 'Add new users and govern as a team',
        icon: User,
        link: 'settings/authentication?action=invite',
        new: true,
    },
    {
        title: 'Connect',
        description: 'Setup Integrations and enable visibility',
        icon: Cable,
        link: 'integrations',
        new: true,
    },

    // {
    //     title: 'Spend',
    //     description: 'See Cloud Spend across clouds, regions, and accounts',
    //     icon: Dollar,
    //     new: false,
    //     link: 'dashboard/spend-accounts',
    // },

    // {
    //     title: 'Insights',
    //     description: 'Get actionable insights',
    //     icon: DocumentChartBarIcon,
    //     link: '/:ws/insights',
    // },
]

// const SvgToComponent = (item: any) => {
//     return item.icon
// }

export default function Shortcuts() {
    const workspace = useParams<{ ws: string }>().ws
    const navigate = useNavigate()
    const [open,setOpen] = useState(false)
    const [userOpen, setUserOpen] = useState(false)


    return (
        <Card>
            <Flex justifyContent="start" className="gap-2 mb-4">
                <Icon icon={CursorArrowRaysIcon} className="p-0" />
                <Title className="font-semibold">Shortcuts</Title>
            </Flex>
            <Grid numItems={4} className="w-full mb-4 gap-4">
                {navList.map((nav, i) => (
                    <>
                        {nav?.title !== 'Audit' && nav?.title !== 'Invite' ? (
                            <>
                                <a
                                    href={`/${nav.link}`}
                                    target={nav.new ? '_blank' : '_self'}
                                >
                                    <Card className=" flex-auto  cursor-pointer  min-h-[140px] h-full pt-3 pb-3 hover:bg-gray-50 hover:dark:bg-gray-900">
                                        <Flex
                                            flexDirection="col"
                                            justifyContent="start"
                                            alignItems="start"
                                            className="gap-2"
                                        >
                                            <img
                                                className="bg-[#1164D9] rounded-[50%] p-[0.3rem] w-7 h-7"
                                                src={nav.icon}
                                            />
                                            <Text className="text-l font-semibold text-gray-900 dark:text-gray-50 text-openg-800  flex flex-row items-center gap-2">
                                                {nav.title}
                                                <ChevronRightIcon className="p-0 w-5 h-5 " />
                                            </Text>
                                            <Text className="text-sm">
                                                {nav.description}
                                            </Text>
                                        </Flex>
                                    </Card>
                                </a>
                            </>
                        ) : (
                            <>
                                <Card
                                    onClick={() => {
                                        if (nav?.title == 'Audit') {
                                            setOpen(true)
                                        } else {
                                            setUserOpen(true)
                                        }
                                    }}
                                    className="  cursor-pointer  min-h-[140px] h-full pt-3 pb-3 hover:bg-gray-50 hover:dark:bg-gray-900"
                                >
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="start"
                                        alignItems="start"
                                        className="gap-2"
                                    >
                                        <img
                                            className="bg-[#1164D9] rounded-[50%] p-[0.3rem] w-7 h-7"
                                            src={nav.icon}
                                        />
                                        <Text className="text-l font-semibold text-gray-900 dark:text-gray-50  flex flex-row items-center gap-2">
                                            {nav.title}
                                            <ChevronRightIcon className="p-0 w-5 h-5 " />
                                        </Text>
                                        <Text className="text-sm">
                                            {nav.description}
                                        </Text>
                                    </Flex>
                                </Card>
                                <Evaluate
                                    opened={open}
                                    id=""
                                    assignmentsCount={0}
                                    benchmarkDetail={undefined}
                                    setOpened={(value: boolean) => {
                                        setOpen(value)
                                    }}
                                    onEvaluate={() => {}}
                                    // complianceScore={0}
                                />
                                <Modal
                                    visible={userOpen}
                                    header={'Invite new member'}
                                    onDismiss={() => {
                                        setUserOpen(false)
                                    }}
                                >
                                    {userOpen && (
                                        <>
                                            <MemberInvite
                                                close={(refresh: boolean) => {
                                                    setUserOpen(false)
                                                }}
                                            />
                                        </>
                                    )}
                                </Modal>
                            </>
                        )}
                    </>
                ))}
            </Grid>
        </Card>
    )
}
