
import axios from 'axios'
import { useAtom } from 'jotai'
import { useEffect, useState } from 'react'
import { meAtom } from '../../../store'
import { Tabs } from '@cloudscape-design/components'
import WidgetsTable from './Table'
import DashboardTable from './Table'

export default function Dashboards() {
    const [dashboards, setDashboards] = useState<any[]>([])
    const [tabId, setTabId] = useState('private')
    const [isLoading, setIsLoading] = useState(true)
    const [me, setMe] = useAtom(meAtom)

    const GetUserDashboards = () => {
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
                `${url}/main/core/api/v4/layout/get`,
                {
                    user_id: me?.username,
                },
                config
            )
            .then((res) => {
                if (res?.data) {
                    setDashboards(res.data)
                } else {
                    setDashboards([])
                }
                setIsLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setIsLoading(false)
            })
    }
    const GetPublicDashboards = () => {
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
            .get(`${url}/main/core/api/v4/layout/public`, config)
            .then((res) => {
                if (res?.data) {
                    setDashboards(res.data)
                } else {
                    setDashboards([])
                }
                setIsLoading(false)
            })
            .catch((err) => {
                // check for 404
                if (err?.response?.status == 404) {
                    setDashboards([])
                }

                console.log(err)
                setIsLoading(false)
            })
    }
    useEffect(() => {
        tabId == 'private' ? GetUserDashboards() : GetPublicDashboards()
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
                        label: 'My Dashboards',
                        content: (
                            <DashboardTable
                                dashboards={dashboards}
                                loading={isLoading}
                            />
                        ),
                    },
                    {
                        id: 'public',
                        label: 'Public Dashboard',
                        content: (
                            <DashboardTable
                                dashboards={dashboards}
                                loading={isLoading}
                            />
                        ),
                    },
                ]}
            />
        </>
    )
}
