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
import Terminal from '../../../icons/Terminal.svg'

import { link } from 'fs'
import { useEffect, useState } from 'react'
import { title } from 'process'
import { Modal } from '@cloudscape-design/components'
import Evaluate from '../../../pages/Governance/Compliance/NewBenchmarkSummary/Evaluate'
import MemberInvite from '../../../pages/Settings/Members/MemberInvite'

const navList = [
    {
        title: 'CloudQL',
        description: 'Query code, cloud, and everything in between.',
        icon: Terminal,
        link: 'cloudql',
        new: true,
    },
    {
        title: 'Connect',
        description: 'Setup Integrations and enable visibility',
        icon: Cable,
        link: 'integration/plugins',
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
    const number = window.innerWidth > 768 ? 4 : 3
    const number1 = window.innerWidth > 768 ? 0 : 1


    return (
        <>
            {window.innerWidth > 640 ? (
                <>
                    {' '}
                    <Grid className="w-full grid-cols-[repeat(auto-fit,_minmax(200px,_1fr))]    2xl:gap-[20px]  sm:gap-7 gap-4">
                        {navList?.slice(number1, number).map((nav, i) => (
                            <>
                                {nav?.title !== 'Audit' &&
                                nav?.title !== 'Invite' ? (
                                    <>
                                        <a
                                            href={`/${nav.link}`}
                                            target={
                                                nav.new ? '_blank' : '_self'
                                            }
                                        >
                                            <Card className=" flex-auto  cursor-pointer  sm:min-h-[130px] h-full pt-3 sm:pb-3 pb-3 hover:bg-gray-50 hover:dark:bg-gray-900">
                                                <Flex
                                                    flexDirection="col"
                                                    justifyContent="start"
                                                    alignItems="start"
                                                    className="gap-2 sm:flex-col flex-row justify-start items-center sm:items-start"
                                                >
                                                    <img
                                                        className="bg-[#1164D9] rounded-[50%] p-[0.3rem] w-7 h-7"
                                                        src={nav.icon}
                                                    />
                                                    <Text className="text-l font-semibold  dark:text-gray-50 text-openg-800  flex flex-row items-center gap-2">
                                                        {nav.title}
                                                        <ChevronRightIcon className="p-0 w-5 h-5 " />
                                                    </Text>
                                                    <Text className="text-sm sm:inline-block hidden">
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
                                            className="  cursor-pointer  sm:min-h-[130px] h-full pt-3 sm:pb-3 pb-3 hover:bg-gray-50 hover:dark:bg-gray-900"
                                        >
                                            <Flex
                                                flexDirection="col"
                                                justifyContent="start"
                                                alignItems="start"
                                                className="gap-2 sm:flex-col flex-row justify-start sm:items-start items-center"
                                            >
                                                <img
                                                    className="bg-[#1164D9] rounded-[50%] p-[0.3rem] w-7 h-7"
                                                    src={nav.icon}
                                                />
                                                <Text className="text-l font-semibold text-gray-900 dark:text-gray-50  flex flex-row items-center gap-2">
                                                    {nav.title}
                                                    <ChevronRightIcon className="p-0 w-5 h-5 " />
                                                </Text>
                                                <Text className="text-sm sm:inline-block hidden">
                                                    {nav.description}
                                                </Text>
                                            </Flex>
                                        </Card>
                                        <Evaluate
                                            open={open}
                                            benchmarkDetail={undefined}
                                            setOpen={(value: boolean) => {
                                                setOpen(value)
                                            }}
                                            onEvaluate={() => {}}
                                            showBenchmark={true}
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
                                                        close={(
                                                            refresh: boolean
                                                        ) => {
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
                </>
            ) : (
                <>
                    {' '}
                    <div className="flex flex-row gap-2 justify-start items-center">
                        {navList?.slice(number1, number).map((nav, i) => (
                            <>
                                {nav?.title !== 'Audit' &&
                                nav?.title !== 'Invite' ? (
                                    <>
                                        <a
                                            href={`/${nav.link}`}
                                            target={
                                                nav.new ? '_blank' : '_self'
                                            }
                                        >
                                            <Card className=" flex-auto  cursor-pointer  sm:min-h-[140px] h-full pt-3 sm:pb-3 pb-3 hover:bg-gray-50 hover:dark:bg-gray-900">
                                                <Flex
                                                    flexDirection="col"
                                                    justifyContent="start"
                                                    alignItems="start"
                                                    className="gap-2 sm:flex-col flex-row justify-start items-center sm:items-start"
                                                >
                                                    <img
                                                        className="bg-[#1164D9] rounded-[50%] p-[0.3rem] w-7 h-7"
                                                        src={nav.icon}
                                                    />
                                                    <Text className="text-l font-semibold  dark:text-gray-50 text-openg-800  flex flex-row items-center gap-2">
                                                        {nav.title}
                                                        <ChevronRightIcon className="p-0 w-5 h-5 " />
                                                    </Text>
                                                    <Text className="text-sm sm:inline-block hidden">
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
                                            className="  cursor-pointer  sm:min-h-[140px] h-full pt-3 sm:pb-3 pb-3 hover:bg-gray-50 hover:dark:bg-gray-900"
                                        >
                                            <Flex
                                                flexDirection="col"
                                                justifyContent="start"
                                                alignItems="start"
                                                className="gap-2 sm:flex-col flex-row justify-start sm:items-start items-center"
                                            >
                                                <img
                                                    className="bg-[#1164D9] rounded-[50%] p-[0.3rem] w-7 h-7"
                                                    src={nav.icon}
                                                />
                                                <Text className="text-l font-semibold text-gray-900 dark:text-gray-50  flex flex-row items-center gap-2">
                                                    {nav.title}
                                                    <ChevronRightIcon className="p-0 w-5 h-5 " />
                                                </Text>
                                                <Text className="text-sm sm:inline-block hidden">
                                                    {nav.description}
                                                </Text>
                                            </Flex>
                                        </Card>
                                        <Evaluate
                                            open={open}
                                            benchmarkDetail={undefined}
                                            setOpen={(value: boolean) => {
                                                setOpen(value)
                                            }}
                                            showBenchmark={true}
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
                                                        close={(
                                                            refresh: boolean
                                                        ) => {
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
                    </div>
                </>
            )}
        </>
    )
}
