import { useEffect, useState } from 'react'
import {
    Button,
    Card,
    Flex,
    List,
    ListItem,
    Subtitle,
    Text,
    TextInput,
} from '@tremor/react'
import clipboardCopy from 'clipboard-copy'
import { DocumentDuplicateIcon } from '@heroicons/react/24/outline'
import { useSetAtom } from 'jotai/index'
import InformationModal from '../../../../components/Modal/InformationModal'
import { useAuthApiV1KeyCreateCreate } from '../../../../api/auth.gen'
import { notificationAtom } from '../../../../store'
import KButton from '@cloudscape-design/components/button'
import { Alert, Input, Modal, Select } from '@cloudscape-design/components'
import axios from 'axios'

interface CreateAPIKeyProps {
    close: () => void
}

const roleItems = [
    {
        value: 'admin',
        title: 'Admin',
        description: 'Have full access',
    },
    {
        value: 'editor',
        title: 'Editor',
        description: 'Can view, edit and delete data',
    },
    {
        value: 'viewer',
        title: 'Viewer',
        description: 'Member can only view the data',
    },
]

export default function CreateConnector({ close }: CreateAPIKeyProps) {
    const [isLoading, setIsLoading] = useState(false)
    const [connector, setConnector] = useState<any>({
        connector_type: 'oidc',
        connector_sub_type: undefined,
        tenant_id: undefined,
        issuer: undefined,
        client_id: undefined,
        client_secret: undefined,
    })
    const [error, setError] = useState<any>(null)
    const setNotification = useSetAtom(notificationAtom)

    const CreateConnector = () => {
        if(  !connector.client_secret || !connector.connector_sub_type){

            setError('Please fill all the fields')
            return
        }
        if(connector.connector_sub_type?.value === 'entraid' && !connector.tenant_id){
            setError('Please fill all the fields')
            return
        }
        if(connector.connector_sub_type?.value === 'general' && !connector.issuer){
            setError('Please fill all the fields')
            return
        }
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
        const body = {
            connector_type: connector.connector_type,
            connector_sub_type: connector.connector_sub_type?.value,
            issuer: connector.issuer,
            tenant_id: connector.tenant_id,
            client_id: connector.client_id,
            client_secret: connector.client_secret,
        }
        axios
            .post(`${url}/main/auth/api/v1/connector`, body, config)
            .then((res) => {
                setIsLoading(false)
                setNotification({
                    type: 'success',
                    text: 'Provider created successfully',
                })
                setNotification({
                    type: 'info',
                    text: 'Please Wait for DEX POD to restart',
                })
                close()
            })
            .catch((err) => {
                console.log(err)
                var error = err.response.data.message
                if(!error){
                    error = 'Failed to create connector'
                }
                setIsLoading(false)
                setNotification({
                    type: 'error',
                    text: error,
                })
                setError(error)
            })
    }

    return (
        <Flex
            flexDirection="col"
            alignItems="start"
            justifyContent="between"
            className="h-full"
        >
            <Flex
                flexDirection="col"
                justifyContent="start"
                alignItems="start"
                className="gap-2 w-full mb-4"
            >
                <Select
                    selectedOption={connector?.connector_sub_type}
                    inlineLabelText="OIDC Provider"
                    onChange={({ detail }) => {
                        setConnector({
                            ...connector,
                            connector_sub_type: detail.selectedOption,
                        })
                        setError(null)
                    }}
                    options={[
                        { label: 'General', value: 'general' },
                        { label: 'Entra Id', value: 'entraid' },
                        {
                            label: 'Google Workspace',
                            value: 'google-workspace',
                        },
                    ]}
                    placeholder="OIDC Provider"
                    className="w-full"
                />
                <Input
                    onChange={({ detail }) => {
                        setConnector({
                            ...connector,
                            client_id: detail.value,
                        })
                        setError(null)
                    }}
                    value={connector?.client_id}
                    placeholder="Client Id"
                    className="w-full"
                />
                <Input
                    onChange={({ detail }) => {
                        setConnector({
                            ...connector,
                            client_secret: detail.value,
                        })
                        setError(null)
                    }}
                    value={connector?.client_secret}
                    placeholder="Client Secret"
                    className="w-full"
                />
                {connector?.connector_sub_type?.value === 'general' && (
                    <>
                        <Input
                            onChange={({ detail }) => {
                                setConnector({
                                    ...connector,
                                    issuer: detail.value,
                                })
                                setError(null)
                            }}
                            value={connector?.issuer}
                            placeholder="Issuer"
                            className="w-full"
                        />
                    </>
                )}
                {connector?.connector_sub_type?.value === 'entraid' && (
                    <>
                        <Input
                            onChange={({ detail }) => {
                                setConnector({
                                    ...connector,
                                    tenant_id: detail.value,
                                })
                                setError(null)
                            }}
                            value={connector?.tenant_id}
                            placeholder="Tenant Id"
                            className="w-full"
                        />
                    </>
                )}
            </Flex>
            {error && (
                <Alert header="Attention" className="w-full mb-3" type="error">
                    {error}
                </Alert>
            )}
            <Flex justifyContent="end" className="space-x-4">
                <KButton
                    onClick={() => {
                        close()
                    }}
                >
                    Cancel
                </KButton>
                <KButton
                    variant="primary"
                    onClick={() => {
                        CreateConnector()
                    }}
                    loading={isLoading}
                >
                    Add Provider
                </KButton>
            </Flex>
        </Flex>
    )
}
