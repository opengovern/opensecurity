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
import { useAtomValue, useSetAtom } from 'jotai'

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
import {
    BreadcrumbGroup,
    KeyValuePairs,
    Spinner,
    Tabs,
} from '@cloudscape-design/components'

import IntegrationList from './Integration'
import CredentialsList from './Credentials'
import { OpenGovernance } from '../../../icons/icons'
import DiscoveryJobs from './Discovery'
import Configuration from './Configuration'
import Setup from './Setup'
import ButtonDropdown from '@cloudscape-design/components/button-dropdown'
import { notificationAtom } from '../../../store'

export default function TypeDetail() {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const { type } = useParams()
    const [manifest, setManifest] = useState<any>()
    const { state } = useLocation()
    const [shcema, setSchema] = useState<Schema>()
    const [loading, setLoading] = useState<boolean>(false)
    const [status, setStatus] = useState<string>()
    const setNotification = useSetAtom(notificationAtom)

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
    const GetManifest = () => {
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
                `${url}/main/integration/api/v1/integration-types/plugin/${type}/manifest`,
                config
            )
            .then((resp) => {
                setManifest(resp.data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)

                // params.fail()
            })
    }
    const GetStatus = () => {
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
            .post(
                `${url}/main/integration/api/v1/integration-types/plugin/${type}/healthcheck`,
                {},
                config
            )
            .then((resp) => {
                setStatus(resp.data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)

                // params.fail()
            })
    }
    const UpdatePlugin = () => {
        setLoading(true)
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        let path = ''
        path = `/main/integration/api/v1/integration-types/plugin/load/id/${type}`

        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }

        axios
            .post(`${url}${path}`, {}, config)
            .then((res) => {
                setLoading(false)

                setNotification({
                    text: `Plugin Updated`,
                    type: 'success',
                })
            })
            .catch((err) => {
                 setNotification({
                     text: `Error: ${err.response.data.message}`,
                     type: 'error',
                 })

                setLoading(false)
            })
    }
    const UnInstallPlugin = () => {
        setLoading(true)
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        let path = ''
        path = `/main/integration/api/v1/integration-types/plugin/uninstall/id/${type}`
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }

        axios
            .post(`${url}${path}`, {}, config)
            .then((res) => {
                setLoading(false)
                setNotification({
                    text: `Plugin Uninstalled`,
                    type: 'success',
                })
                 navigate('/plugins')
            })
            .catch((err) => {
                 setNotification({
                     text: `Error: ${err.response.data.message}`,
                     type: 'error',
                 })
                setLoading(false)
            })
    }

    const DisablePlugin = () => {
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
            .post(
                `${url}/main/integration/api/v1/integration-types/plugin/${type}/disable`,
                {},
                config
            )
            .then((res) => {
                setLoading(false)
                navigate('/plugins')
            })
            .catch((err) => {
                setLoading(false)
                setNotification({
                    text: `Error: ${err.response.data.message}`,
                    type: 'error',
                })
            })
    }
    useEffect(() => {
        GetSchema()
        GetStatus()
        GetManifest()
    }, [])

    return (
        <>
            {/* <TopHeader breadCrumb={[state?.name]} /> */}

            {shcema && shcema?.integration_type_id ? (
                <>
                    <Flex className="flex-col w-full justify-start items-start gap-4">
                        <BreadcrumbGroup
                            className="w-full"
                            items={[
                                {
                                    text: 'Plugins',
                                    href: '/plugins',
                                },
                                {
                                    // @ts-ignore
                                    text: state?.name,
                                    href: `/plugins/${type}`,
                                },
                            ]}
                        />
                        <Flex className="flex-col gap-3 justify-start items-start w-full">
                            <Flex className="flex-row justify-between w-full gap-8">
                                <h1 className=" font-bold text-2xl mb-2  text-left ml-1">
                                    {state?.name} plugin
                                </h1>
                                <ButtonDropdown
                                onItemClick={({detail})=>{
                                    const id = detail.id
                                    switch (id){
                                        case 'update':
                                            UpdatePlugin()
                                            break;
                                        case 'disable':
                                            DisablePlugin()
                                            break;
                                        case 'uninstall':
                                            UnInstallPlugin()
                                            break;
                                        default:
                                            break
                                    }
                                }}
                                    variant="primary"
                                    items={[
                                        {
                                            text: 'Settings',
                                            items: [
                                                {
                                                    text: 'Update',
                                                    id: 'update',
                                                },
                                                {
                                                    text: 'Disable',
                                                    id: 'disable',
                                                },
                                                {
                                                    text: 'Uninstall',
                                                    id: 'uninstall',
                                                },
                                            ],
                                        },
                                      
                                    ]}
                                >
                                    Actions
                                </ButtonDropdown>
                            </Flex>
                            <Card className="">
                                <Flex
                                    flexDirection="col"
                                    className="gap-2 p-2 w-full justify-start items-start"
                                >
                                    <>
                                        <KeyValuePairs
                                            columns={4}
                                            items={[
                                                {
                                                    label: 'Id',
                                                    value: manifest?.IntegrationType,
                                                },
                                                {
                                                    label: 'Artifact URL',
                                                    value: manifest?.DescriberURL,
                                                },
                                                {
                                                    label: 'Version',
                                                    value: manifest?.DescriberTag,
                                                },
                                                {
                                                    label: 'Publisher',
                                                    value: manifest?.Publisher,
                                                },
                                                {
                                                    label: 'Author',
                                                    value: manifest?.Author,
                                                },
                                                {
                                                    label: 'Supported Platform Version',
                                                    value: manifest?.SupportedPlatformVersion,
                                                },
                                                {
                                                    label: 'Update date',
                                                    value: manifest?.UpdateDate,
                                                },
                                                {
                                                    label: 'Operational Status',
                                                    // @ts-ignore
                                                    value: status
                                                        ? status
                                                              ?.charAt(0)
                                                              .toUpperCase() +
                                                          status?.slice(1)
                                                        : '',
                                                },
                                            ]}
                                        />
                                    </>
                                </Flex>
                            </Card>
                            <></>
                        </Flex>
                        <Tabs
                            tabs={[
                                {
                                    id: '3',
                                    label: 'Resource Types',
                                    content: (
                                        <Configuration
                                            name={state?.name}
                                            integration_type={type}
                                        />
                                    ),
                                },
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
                                    id: '4',
                                    label: 'Setup Guide',
                                    content: (
                                        <Setup
                                            name={state?.name}
                                            integration_type={type}
                                        />
                                    ),
                                },
                            ]}
                        />
                    </Flex>
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
                                                navigate('/plugins')
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
