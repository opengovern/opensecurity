import {
    Button,
    Card,
    Col,
    Flex,
    Grid,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
} from '@tremor/react'
import { useEffect, useState } from 'react'
import FindingsWithFailure from './FindingsWithFailure'
import TopHeader from '../../../../../components/Layout/Header'
// import Filter from './Filter'
import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../../../../api/api'
import {
    DateRange,
    useURLParam,
    useURLState,
} from '../../../../../utilities/urlstate'
import Spinner from '../../../../../components/Spinner'
interface Props {
    id: string
}
export default function Findings({ id }: Props) {
    const [tab, setTab] = useState<number>(0)
    const [selectedGroup, setSelectedGroup] = useState<
        'findings' | 'resources' | 'controls' | 'accounts' | 'events'
    >('findings')
    useEffect(() => {
        switch (tab) {
            case 0:
                setSelectedGroup('findings')
                break
            case 1:
                setSelectedGroup('events')
                break
            case 2:
                setSelectedGroup('controls')
                break
            case 3:
                setSelectedGroup('resources')
                break
            case 4:
                setSelectedGroup('accounts')
                break
            default:
                setSelectedGroup('findings')
                break
        }
    }, [tab])

    const findComponent = () => {
        switch (tab) {
            case 0:
                return <FindingsWithFailure query={query} />
            default:
                return <Spinner />
        }
    }

    const [query, setQuery] = useState<{
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
        jobID: string[] | undefined
    }>({
        connector: SourceType.Nil,
        conformanceStatus: [
            PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
        ],
        severity: [
            TypesFindingSeverity.FindingSeverityCritical,
            TypesFindingSeverity.FindingSeverityHigh,
            TypesFindingSeverity.FindingSeverityMedium,
            TypesFindingSeverity.FindingSeverityLow,
            TypesFindingSeverity.FindingSeverityNone,
        ],
        connectionID: [],
        controlID: [],
        benchmarkID: [id],
        resourceTypeID: [],
        lifecycle: [true],
        activeTimeRange: undefined,
        eventTimeRange: undefined,
        jobID: [],
    })

    return (
        <>
            {/* <TopHeader /> */}
            {/* @ts-ignore */}
            {/* <Filter type={selectedGroup} onApply={(e) => setQuery(e)} id={id} /> */}
            {/* <Flex className="mt-2">{findComponent()}</Flex> */}
            <Grid numItems={6} className="mt-2 gap-2">
               
                <Col numColSpan={6}>{findComponent()}</Col>
            </Grid>
        </>
    )
}
