import { useAtomValue } from 'jotai'
import { Card, Flex, Grid, Text, Title } from '@tremor/react'
import { useEffect } from 'react'
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline'
import DrawerPanel from '../../../../../components/DrawerPanel'
import { getConnectorIcon } from '../../../../../components/Cards/ConnectorCard'
import SummaryCard from '../../../../../components/Cards/SummaryCard'
import { useComplianceApiV1FindingsFiltersCreate } from '../../../../../api/compliance.gen'
import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    TypesFindingSeverity,
} from '../../../../../api/api'
import { isDemoAtom } from '../../../../../store'

interface IAccountDetail {
    account: any
    open: boolean
    onClose: () => void
}

export default function CloudAccountDetail({
    account,
    open,
    onClose,
}: IAccountDetail) {
    const isDemo = useAtomValue(isDemoAtom)
    const { response: filters, sendNow } =
        useComplianceApiV1FindingsFiltersCreate(
            {
                integrationID: [account?.id || ''],
            },
            {},
            false
        )

    useEffect(() => {
        if (account) {
            sendNow()
        }
    }, [account])

    const options = [
        {
            name: 'Failed',
            value:
                filters?.conformanceStatus?.find((f) => f.key === 'failed')
                    ?.count || 0,
            icon: <XCircleIcon className="h-5 text-rose-600" />,
        },
        {
            name: 'Passed',
            value:
                filters?.conformanceStatus?.find((f) => f.key === 'passed')
                    ?.count || 0,
            icon: <CheckCircleIcon className="h-5 text-emerald-500" />,
        },
    ]

    const severity = [
        {
            name: 'Critical',
            value:
                filters?.severity?.find((f) => f.key === 'critical')?.count ||
                0,
            color: '#6E120B',
        },
        {
            name: 'High',
            value: filters?.severity?.find((f) => f.key === 'high')?.count || 0,
            color: '#CA2B1D',
        },
        {
            name: 'Medium',
            value:
                filters?.severity?.find((f) => f.key === 'medium')?.count || 0,
            color: '#EE9235',
        },
        {
            name: 'Low',
            value: filters?.severity?.find((f) => f.key === 'low')?.count || 0,
            color: '#F4C744',
        },
        {
            name: 'None',
            value: filters?.severity?.find((f) => f.key === 'none')?.count || 0,
            color: '#9BA2AE',
        },
    ]

    return (
        <DrawerPanel
            open={open}
            onClose={onClose}
            title={
                <Flex justifyContent="start">
                    {getConnectorIcon(account?.connector)}
                    <Title className="text-lg font-semibold ml-2 my-1">
                        {account?.providerConnectionName}
                    </Title>
                </Flex>
            }
        >
            <Grid className="w-full gap-4 mb-6" numItems={2}>
                <SummaryCard
                    title="Account ID"
                    metric={account?.providerConnectionID}
                    blur={isDemo}
                    isString
                />
                <SummaryCard
                    title="Account name"
                    metric={account?.providerConnectionName}
                    blur={isDemo}
                    isString
                />
            </Grid>
            <Title className="border-b border-b-gray-200 py-4 mb-4">
                Findings
            </Title>
            <Grid className="w-full gap-4 mb-6" numItems={2}>
                {options.map((o) => (
                    <Card className="p-4">
                        <Flex>
                            <Flex className="gap-1.5 w-fit">
                                <Flex className="gap-1 w-fit">
                                    {o.icon}
                                    <Text>{o.name}</Text>
                                </Flex>
                            </Flex>
                            <Text className="text-gray-800">{o.value}</Text>
                        </Flex>
                    </Card>
                ))}
            </Grid>
            <Title className="border-b border-b-gray-200 py-4 mb-4">
                Issues
            </Title>
            <Grid className="w-full gap-4 mb-6" numItems={2}>
                {severity.map((s) => (
                    <Card className="p-4">
                        <Flex>
                            <Flex className="gap-1.5 w-fit">
                                <div
                                    className="h-4 w-1.5 rounded-sm"
                                    style={{
                                        backgroundColor: s.color,
                                    }}
                                />
                                <Text>{s.name} #</Text>
                            </Flex>
                            <Text className="text-gray-800">{s.value}</Text>
                        </Flex>
                    </Card>
                ))}
            </Grid>
        </DrawerPanel>
    )
}
