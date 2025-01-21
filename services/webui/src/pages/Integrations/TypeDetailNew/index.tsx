import { Button, Card, Flex, Title, Text } from '@tremor/react'
import {
    useLocation,
    useNavigate,
    useParams,
    useSearchParams,
} from 'react-router-dom'
import {
    ArrowLeftStartOnRectangleIcon,
    Cog8ToothIcon,
} from '@heroicons/react/24/outline'
import { useAtomValue } from 'jotai'


import {
    useIntegrationApiV1ConnectorsMetricsList,
    useIntegrationApiV1CredentialsList,
} from '../../../api/integration.gen'
import TopHeader from '../../../components/Layout/Header'
import {
    defaultTime,
    searchAtom,
    useUrlDateRangeState,
} from '../../../utilities/urlstate'
import axios from 'axios'
import { useEffect, useState } from 'react'
import { Schema } from './types'
import { Spinner, Tabs } from '@cloudscape-design/components'

import IntegrationList from './Integration'
import CredentialsList from './Credentials'
import { OpenGovernance } from '../../../icons/icons'
import DiscoveryJobs from './Discovery'
import Setup from './Setup'
import Manifest from './Manifest'

export default function TypeDetail() {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const { type } = useParams()
    const { state } = useLocation()
    const [shcema, setSchema] = useState<Schema>()
    const [loading, setLoading] = useState<boolean>(false)

    const GetSchema = () => {
        setLoading(true)
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
            .get(
                `${url}/main/integration/api/v1/integrations/types/${type}/ui/spec `,
                config
            )
            .then((res) => {
                const data = res.data
                setSchema(data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }

    useEffect(() => {
        GetSchema()
    }, [])

    return (
        <>
            <TopHeader breadCrumb={[type]} />

            {shcema && shcema?.integration_type_id ? (
                <>
                    <Tabs
                        tabs={[
                            {
                                id: '0',
                                label: 'Integrations',
                                content: (
                                    <IntegrationList
                                        schema={shcema}
                                        name={state?.name}
                                        integration_type={type}
                                    />
                                ),
                            },
                            {
                                id: '1',
                                label: 'Credentials',
                                content: (
                                    <CredentialsList
                                        schema={shcema}
                                        name={state?.name}
                                        integration_type={type}
                                    />
                                ),
                            },
                            {
                                id: '2',
                                label: 'Discovery Jobs',
                                content: (
                                    <DiscoveryJobs
                                        name={state?.name}
                                        integration_type={type}
                                    />
                                ),
                            },
                            {
                                id: '3',
                                label: 'Setup Guide',
                                content: (
                                    <Setup
                                        name={state?.name}
                                        integration_type={type}
                                    />
                                ),
                            },
                            {
                                id: '4',
                                label: 'Plugin information',
                                content: (
                                    <Manifest
                                        name={state?.name}
                                        integration_type={type}
                                    />
                                ),
                            },
                        ]}
                    />
                </>
            ) : (
                <>
                    {loading ? (
                        <>
                            <Spinner />
                        </>
                    ) : (
                        <>
                            <Flex
                                flexDirection="col"
                                className="fixed top-0 left-0 w-screen h-screen bg-gray-900/80 z-50"
                            >
                                <Card className="w-1/3 mt-56">
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="center"
                                        alignItems="center"
                                    >
                                        <OpenGovernance className="w-14 h-14 mb-6" />
                                        <Title className="mb-3 text-2xl font-bold">
                                            Data not found
                                        </Title>
                                        <Text className="mb-6 text-center">
                                            Json schema not found for this
                                            integration
                                        </Text>
                                        <Button
                                            icon={ArrowLeftStartOnRectangleIcon}
                                            onClick={() => {
                                                navigate('/integrations')
                                            }}
                                        >
                                            Back
                                        </Button>
                                    </Flex>
                                </Card>
                            </Flex>
                        </>
                    )}
                </>
            )}
        </>
    )
}
