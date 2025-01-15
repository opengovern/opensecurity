import {
    Flex,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
} from '@tremor/react'
import { useEffect, useState } from 'react'

import ResourcesWithFailure from '../ResourcesWithFailure'
import ControlsWithFailure from '../ControlsWithFailure'

import Spinner from '../../../../components/Spinner'
import {
    SourceType,
    TypesFindingSeverity,
    PlatformEnginePkgComplianceApiConformanceStatus,
} from '../../../../api/api'
import { DateRange } from '../../../../utilities/urlstate'
import FailingCloudAccounts from '../FailingCloudAccounts'
import Events from '../Events'
import Tabs from '@cloudscape-design/components/tabs'

interface ICount {
    connector: SourceType
    conformanceStatus:
        | PlatformEnginePkgComplianceApiConformanceStatus[]
        | undefined
    severity: TypesFindingSeverity[] | undefined
    connectionID: string[] | undefined
    controlID: string[] | undefined
    benchmarkID: string[] | undefined
    resourceTypeID: string[] | undefined
    lifecycle: boolean[] | undefined
    activeTimeRange: DateRange | undefined
    eventTimeRange: DateRange | undefined
}
interface Props {
    query: ICount
    setSelectedGroup: Function
    tab: number
    setTab: Function
}
const GROUPS = { 0: 'accounts', 1: 'resources', 2: 'controls' }

export default function Summary({
    query,
    setSelectedGroup,
    tab,
    setTab,
}: Props) {
    return (
        <>
            <Tabs
                onChange={({ detail }) => {
                    setTab(parseInt(detail.activeTabId))
                    // @ts-ignore
                    setSelectedGroup(GROUPS[parseInt(detail.activeTabId)])
                }}
                activeTabId={tab.toString()}
                tabs={[
                    {
                        label: ' Account Summary',
                        id: '0',
                        content: (
                            <>
                                {tab == 0 && (
                                    <>
                                        <FailingCloudAccounts query={query} />
                                    </>
                                )}
                            </>
                        ),
                    },
                    {
                        label: 'Resource Incident Summary',
                        id: '1',
                        content: (
                            <>
                                {' '}
                                {tab == 1 && (
                                    <>
                                        <ResourcesWithFailure query={query} />
                                    </>
                                )}
                            </>
                        ),
                    },
                    {
                        label: 'Control Incident Summary',
                        id: '2',
                        content: (
                            <>
                                {' '}
                                {tab == 2 && (
                                    <ControlsWithFailure query={query} />
                                )}
                            </>
                        ),
                    },
                ]}
            />
        </>
    )
}
