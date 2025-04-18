import {
    Card,
    Col,
    Flex,
    Grid,
    Icon,
    ProgressCircle,
    Title,
} from '@tremor/react'
import { useEffect, useState } from 'react'

import {
    Header,
    Pagination,
    PropertyFilter,
    Tabs,
} from '@cloudscape-design/components'

import AllControls from './All Controls'
import SettingsParameters from '../Settings/Parameters'
import AllPolicy from './All Policy'
import { useParams, useSearchParams } from 'react-router-dom'
import Framework from './FrameWorks'


export default function Compliance() {



    const [tab,setTab] = useState('0')
    const [tab_id,setTabID] = useSearchParams()
 
    useEffect(()=>{
        switch (tab_id?.get('tab')) {
            case 'frameworks':
                setTab('frameworks')
                break
            case 'controls':
                setTab('controls')
                break
            case 'policies':
                setTab('policies')
                break
            case 'parameters':
                setTab('parameters')
                break
            default:
                setTab('frameworks')
                break
        }
    },[tab_id])

    return (
        <>
            <Tabs
                activeTabId={tab}
                className={`w-[270px]`}
                onChange={({ detail }) => {
                console.log(detail.activeTabId)
                    tab_id.set('tab', detail.activeTabId)
                    setTabID(tab_id)
                    setTab(detail.activeTabId)
                }}
                tabs={[
                    {
                        label: 'Frameworks',
                        id: 'frameworks',
                        content: (
                            <Framework />
                        ),
                    },
                    {
                        id: 'controls',
                        label: 'Controls',
                        content: <AllControls />,
                    },
                    {
                        id: 'policies',
                        label: 'Policies',
                        content: <AllPolicy />,
                    },
                    {
                        id: 'parameters',
                        label: 'Parameters',
                        content: <SettingsParameters />,
                    },
                ]}
            />
        </>
    )
}
