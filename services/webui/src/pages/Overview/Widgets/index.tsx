
import axios from 'axios'
import { useAtom } from 'jotai'
import { useEffect, useState } from 'react'
import { meAtom } from '../../../store'
import { Tabs } from '@cloudscape-design/components'
import WidgetsTable from './Table'

export default function Widgets() {
    const [widgets, setWidgets] = useState<any[]>([])
    const [tabId, setTabId] = useState('private')
    const [isLoading, setIsLoading] = useState(true)
    const [me, setMe] = useAtom(meAtom)
    

    const GetUserWidgets = (
      
    ) => {
        setIsLoading(true)
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        

        axios
            .post(
                `${url}/main/core/api/v4/layout/widget/get`,
                {
                    user_id: me?.username,
                },
                config
            )
            .then((res) => {
                 if (res?.data) {
                     setWidgets(res.data)
                 } else {
                     setWidgets([])
                 }
                 setIsLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setIsLoading(false)
            })
    }
    const GetPublicWidgets = () => {
        setIsLoading(true)
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }

        axios
            .post(
                `${url}/main/core/api/v4/layout/widget/get/public`,
                {},
                config
            )
            .then((res) => {
                if(res?.data){
                    setWidgets(res.data)
                }
                else{
                    setWidgets([])
                }
                setIsLoading(false)

            })
            .catch((err) => {
                console.log(err)
                setIsLoading(false)
            })
    }
    useEffect(() => {
        tabId == 'private' ? GetUserWidgets() : GetPublicWidgets()
    }, [tabId])
  
    return (
        <>
            <Tabs
                activeTabId={tabId}
                onChange={(event) => {
                    setTabId(event.detail.activeTabId)
                }}
                tabs={[
                    {
                        id: 'private',
                        label: 'My Widgets',
                        content: (
                            <WidgetsTable
                                widgets={widgets}
                                loading={isLoading}
                            />
                        ),
                    },
                    {
                        id: 'public',
                        label: 'Public Widgets',
                        content: (
                            <WidgetsTable
                                widgets={widgets}
                                loading={isLoading}
                            />
                        ),
                    },
                ]}
            />
        </>
    )
}
