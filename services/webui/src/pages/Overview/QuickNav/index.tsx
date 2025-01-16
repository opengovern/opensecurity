import {
    ArrowTopRightOnSquareIcon,
    BanknotesIcon,
    CubeIcon,
    CursorArrowRaysIcon,
    PuzzlePieceIcon,
    ShieldCheckIcon,
} from '@heroicons/react/24/outline'
import { Card, Flex, Grid, Icon, Text, Title } from '@tremor/react'
import { useNavigate, useParams } from 'react-router-dom'

const navList = [
    {
        title: 'Infrastructure',
        description: 'Discover and query cloud assets across multiple clouds.',
        icon: CubeIcon,
        children: [
            { name: 'Overview', link: 'dashboard/infrastructure' },
            {
                name: 'Infrastructure metrics',
                link: 'dashboard/infrastructure-metrics',
            },
            {
                name: 'Cloud account details',
                link: 'dashboard/infrastructure-cloud-accounts',
            },
        ],
    },
    {
        title: 'Spend',
        description: 'Understand cloud spend, reduce wastage.',
        icon: BanknotesIcon,
        children: [
            { name: 'Summary', link: 'dashboard/spend' },
            { name: 'Spend metrics', link: 'dashboard/spend-metrics' },
            {
                name: 'Spend by accounts',
                link: 'dashboard/spend-accounts',
            },
        ],
    },
    {
        title: 'Security',
        description: 'Audit cloud configurations against benchmarks.',
        icon: ShieldCheckIcon,
        children: [
            { name: 'Benchmarks', link: 'compliance' },
            { name: 'Security findings', link: 'incidents' },
        ],
    },
    {
        title: 'Integrations',
        description:
            'Connect your cloud environments and tools to OpenGovernance to gain insights.',
        icon: PuzzlePieceIcon,
        children: [
            { name: 'AWS accounts', link: 'integrations/AWS' },
            { name: 'Azure subscriptions', link: 'integrations/Azure' },
        ],
    },
    // {
    //     title: 'Insights',
    //     description: 'Get actionable insights',
    //     icon: DocumentChartBarIcon,
    //     link: '/:ws/insights',
    // },
]

export default function QuickNav() {
    const workspace = useParams<{ ws: string }>().ws
    const navigate = useNavigate()

    return (
        <Card>
            <Flex justifyContent="start" className="gap-2 mb-4">
                <Icon icon={CursorArrowRaysIcon} className="p-0" />
                <Title className="font-semibold">Quick Navigation</Title>
            </Flex>
            <Grid numItems={4} className="w-full mb-4">
                {navList.map((nav, i) => (
                    <div
                        className={
                            i < navList.length - 1
                                ? 'border-r border-r-gray-200 dark:border-r-gray-700'
                                : ''
                        }
                    >
                        <Card
                            key={nav.title}
                            className="border-0 ring-0 !shadow-none py-0"
                        >
                            <Flex justifyContent="start" className="gap-2 mb-2">
                                <Icon icon={nav.icon} className="p-0" />
                                <Title className="text-gray-800 line-clamp-1">
                                    {nav.title}
                                </Title>
                            </Flex>
                            <Text className="line-clamp-2 mb-5 h-10">
                                {nav.description}
                            </Text>
                            <Flex
                                flexDirection="col"
                                justifyContent="start"
                                alignItems="start"
                            >
                                {nav.children.map((c) => (
                                    <Flex
                                        justifyContent="start"
                                        className="gap-2 py-1.5 cursor-pointer hover:text-openg-50"
                                        onClick={() =>
                                            navigate(
                                                c.link.replaceAll(
                                                    ':ws',
                                                    workspace || ''
                                                )
                                            )
                                        }
                                    >
                                        <ArrowTopRightOnSquareIcon className="h-5 text-openg-500" />
                                        <Text className="text-inherit dark:hover:text-openg-800">
                                            {c.name}
                                        </Text>
                                    </Flex>
                                ))}
                            </Flex>
                        </Card>
                    </div>
                ))}
            </Grid>
        </Card>
    )
}
