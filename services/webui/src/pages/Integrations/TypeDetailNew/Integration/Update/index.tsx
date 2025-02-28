import { Flex, Title } from '@tremor/react'
import {
    useLocation,
    useNavigate,
    useParams,
    useSearchParams,
} from 'react-router-dom'
import { Cog8ToothIcon } from '@heroicons/react/24/outline'
import { useAtomValue } from 'jotai'

import axios from 'axios'
import { useEffect, useState } from 'react'
import {  Integration, Schema } from '../../types'

import {
    Alert,
    AppLayout,
    Box,
    Button,
    Header,
    Modal,
    Multiselect,
    Pagination,
    SpaceBetween,
    SplitPanel,
    Table,
    Tabs,
} from '@cloudscape-design/components'
import { GetDiscover, GetDiscoverField, GetEditField, RenderInputField } from '../../utils'

interface UpdateIntegrationProps {
    name?: string
    integration_type?: string
    schema?: Schema
    open: boolean
    onClose: () => void
    GetList: Function
    selectedItem?: Integration
    actionFunc? : Function
}

export default function UpdateIntegration({
    name,
    integration_type,
    schema,
    open,
    onClose,
    GetList,
    selectedItem,
    actionFunc

}: UpdateIntegrationProps) {
    const navigate = useNavigate()
    const [row, setRow] = useState<Integration[]>([])

    const [loading, setLoading] = useState<boolean>(false)
    const [error, setError] = useState<string>('')
    const [selectedCredential, setSelectedCredential] = useState<number>(-1)
    const [integration, setIntegration] = useState<any>({})
    const [providers, setProviders] = useState<any>([])
    const [selectedProviders, setSelectedProviders] = useState<any>([])
    const [credentialId, setCredentialId] = useState<string>('')
    const [described, setDescribed] = useState<boolean>(false)
    const DiscoverIntegrations = () => {
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

        const body = {
            integration_type: integration_type,
            credentials: integration,
        }
        axios
            .post(
                `${url}/main/integration/api/v1/integrations/discover`,
                body,
                config
            )
            .then((res) => {
                const data = res.data

                setCredentialId(data.credential_id)
                setProviders(data.integrations)
                setDescribed(true)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }
    const AddIntegration = () => {
        if (selectedProviders.length === 0) {
            setError('Please select provider')
            return
        }
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

        const body = {
            integration_type: integration_type,
            provider_ids: selectedProviders?.map((provider: any) => {
                return provider.value
            }),
            credential_id: credentialId,
        }
        axios
            .post(
                `${url}/main/integration/api/v1/integrations/add`,
                body,
                config
            )
            .then((res) => {
                GetList()
                onClose()
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }

    useEffect(() => {}, [])

    return (
        <>
            <Modal
                visible={open}
                onDismiss={onClose}
                header="Update Integration"
                className="p-4"
            >
                <>
                    {described ? (
                        <>
                            <Title className="mb-4">
                                Please Select Provider
                            </Title>
                            <Flex
                                className="mt-4 mb-2 gap-4 w-full"
                                justifyContent="center"
                            >
                                <>
                                    <Multiselect
                                        className="w-full"
                                        selectedOptions={selectedProviders}
                                        options={providers.map(
                                            (provider: any) => {
                                                return {
                                                    label: provider?.name,
                                                    value: provider?.provider_id,
                                                    description:
                                                        provider?.provider_id,
                                                }
                                            }
                                        )}
                                        onChange={({ detail }) => {
                                            setSelectedProviders(
                                                detail.selectedOptions
                                            )
                                        }}
                                        placeholder="Select Provider"
                                    />
                                </>
                            </Flex>
                            <Flex
                                className="gap-2 w-full mt-4"
                                justifyContent="end"
                                alignItems="center"
                            >
                                <Button
                                    onClick={() => {
                                        setDescribed(false)
                                    }}
                                >
                                    Back
                                </Button>

                                <Button
                                    onClick={() => {
                                        AddIntegration()
                                    }}
                                    disabled={
                                        loading ||
                                        selectedProviders.length === 0
                                    }
                                    loading={loading}
                                    variant="primary"
                                >
                                    Add
                                </Button>
                            </Flex>
                            {error && error != '' && (
                                <>
                                    <Alert
                                        type="error"
                                        header={'Please select provider'}
                                        className="mt-4"
                                    >
                                        {error}
                                    </Alert>
                                </>
                            )}
                        </>
                    ) : (
                        <>
                            <Flex
                                className="mt-4 mb-2 gap-4 w-full"
                                justifyContent="start"
                                alignItems="start"
                                flexDirection="col"
                            >
                                <>
                                    {/* {GetEditField(
                                        schema,
                                        0
                                    )?.map((field) => {
                                        return (
                                            <>
                                                {RenderInputField(
                                                    field,
                                                    (value: any) => {
                                                        setIntegration({
                                                            ...integration,
                                                            [field.name]: value,
                                                        })
                                                    },
                                                    integration[field.name]
                                                )}
                                            </>
                                        )
                                    })} */}
                                </>
                            </Flex>
                            <Flex
                                className="gap-2 w-full mt-4"
                                justifyContent="end"
                                alignItems="center"
                            >
                                <Button
                                    onClick={() => {
                                        setSelectedCredential(-1)
                                    }}
                                >
                                    Back
                                </Button>

                                <Button
                                    onClick={() => {
                                        DiscoverIntegrations()
                                    }}
                                    loading={loading}
                                    variant="primary"
                                >
                                    See Providers
                                </Button>
                            </Flex>
                        </>
                    )}
                </>
            </Modal>
        </>
    )
}
