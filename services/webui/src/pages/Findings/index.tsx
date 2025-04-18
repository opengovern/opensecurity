
import { useEffect, useState } from 'react'


import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../api/api'
import {
    DateRange,
    useURLParam,
    useURLState,
} from '../../utilities/urlstate'

import AllIncidents from './AllIncidents'


export default function Findings() {
    const [tab, setTab] = useState<number>(0);
    const [selectedGroup, setSelectedGroup] = useState<
        'findings' | 'resources' | 'controls' | 'accounts' | 'events'
    >('findings')
    useEffect(() => {
        switch (tab) {
            case 0:
                setSelectedGroup('findings')
                break
            case 1:
                setSelectedGroup('resources')
                break
            default:
                setSelectedGroup('findings')
                break
        }
    }, [tab])
    useEffect(() => {
        const url = window.location.pathname.split('/')[2]
        // setShow(false);
        
        switch (url) {
            case 'summary':
                setTab(1)
                break
            case 'controls':
                setTab(1)
                break
            case 'resources':
                setTab(2)
                break

            default:
                setTab(0)
                break
        }
    }, [window.location.pathname])
 

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
        connectionGroup: []
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
        benchmarkID: [],
        resourceTypeID: [],
        lifecycle: [true],
        activeTimeRange: undefined,
        eventTimeRange: undefined,
        jobID: [],
        connectionGroup: [],
    })

    return (
        <>
            <AllIncidents
                query={query}
                setSelectedGroup={setSelectedGroup}
                tab={tab}
                setTab={setTab}
            />
        </>
    )
}
