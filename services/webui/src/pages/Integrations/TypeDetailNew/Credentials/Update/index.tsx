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
import {
    DiscoverCredential,
    Integration,
    Schema,
    Credentials,
} from '../../types'

import {
    Alert,
    AppLayout,
    Box,
    Button,
    FormField,
    Header,
    Input,
    Modal,
    Multiselect,
    Pagination,
    SpaceBetween,
    SplitPanel,
    Table,
    Tabs,
} from '@cloudscape-design/components'
import {
    GetDiscover,
    GetDiscoverField,
    GetUpdateCredentialFields,
    RenderInputField,
} from '../../utils'

interface CreateIntegrationProps {
    name?: string
    integration_type?: string
    schema?: Schema
    open: boolean
    onClose: () => void
    GetList: Function
    selectedItem?: Credentials
}

export default function UpdateCredentials({
    name,
    integration_type,
    schema,
    open,
    onClose,
    GetList,
    selectedItem,
}: CreateIntegrationProps) {
    const navigate = useNavigate()
    const [row, setRow] = useState<Integration[]>([])

    const [loading, setLoading] = useState<boolean>(false)
    const [error, setError] = useState<string>('')
    const [selectedCredential, setSelectedCredential] = useState<number>(-1)
    const [credential, setCredentials] = useState<any>({})
    const [providers, setProviders] = useState<any>([])
    const [selectedProviders, setSelectedProviders] = useState<any>([])
    const [credentialId, setCredentialId] = useState<string>('')
    const [described, setDescribed] = useState<boolean>(false)
    const [credentialType, setCredentialType] = useState<string>('')
    const [description, setDescription] = useState<string>('')

    const UpdateCredentials = () => {
        setLoading(true)
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token
        var has_file = false
        Object.keys(credential).forEach((key) => {
            if (credential[key] instanceof File) {
                has_file = true
            }
        })
        if (has_file) {
            const formData = new FormData()
            // @ts-ignore
            formData.append('credential_type', credentialType)
            formData.append('description', description)

            Object.keys(credential).forEach((key) => {
                formData.append(`credentials.${key}`, credential[key])
            })

            axios
                .post(
                    `${url}/main/integration/api/v1/credentials/${selectedItem?.id}`,
                    formData,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                            'Content-Type': 'multipart/form-data',
                        },
                    }
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
        } else {
             const config = {
                 headers: {
                     Authorization: `Bearer ${token}`,
                 },
             }

             const body = {
                 credentials: credential,
                 credential_type: credentialType,
                    description: description
             }
             axios
                 .put(
                     `${url}/main/integration/api/v1/credentials/${selectedItem?.id}`,
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
      
    }

    useEffect(() => {}, [])

    return (
        <>
            <Modal
                visible={open}
                onDismiss={onClose}
                header="Update Credential"
                className="p-4"
            >
                {selectedCredential === -1 ? (
                    <>
                        <Title className="mb-4">
                            Please Select Credential Input method
                        </Title>
                        <Flex
                            className="mt-4 mb-2 gap-4"
                            justifyContent="center"
                        >
                            <>
                                {GetDiscover(schema)?.map(
                                    (
                                        cred: DiscoverCredential,
                                        index: number
                                    ) => {
                                        return (
                                            <>
                                                <Button
                                                    onClick={() => {
                                                        setSelectedCredential(
                                                            index
                                                        )
                                                        setCredentials({
                                                            ...credential,
                                                            credential_type:
                                                                cred.label,
                                                        })
                                                        setCredentialType(
                                                            cred.label
                                                        )
                                                    }}
                                                >
                                                    {cred.label}
                                                </Button>
                                            </>
                                        )
                                    }
                                )}
                            </>
                        </Flex>
                    </>
                ) : (
                    <>
                        <Title className="mb-4">
                            {GetDiscover(schema)?.[selectedCredential].label}
                        </Title>
                        <Flex
                            className="mt-4 mb-2 gap-4 w-full"
                            justifyContent="start"
                            alignItems="start"
                            flexDirection="col"
                        >
                            <>
                                {GetUpdateCredentialFields(
                                    schema,
                                    selectedCredential
                                )?.map((field) => {
                                    return (
                                        <>
                                            {RenderInputField(
                                                field,
                                                (value: any) => {
                                                    setCredentials({
                                                        ...credential,
                                                        [field.name]: value,
                                                    })
                                                },
                                                credential[field.name]
                                            )}
                                        </>
                                    )
                                })}
                                <FormField
                                    className="w-full"
                                    label={'Description'}
                                    errorText={''}
                                >
                                    <Input
                                        className="w-full"
                                        value={description}
                                        type={'text'}
                                        onChange={({ detail }) =>
                                            setDescription(detail.value)
                                        }
                                    />
                                </FormField>
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
                                    UpdateCredentials()
                                }}
                                loading={loading}
                                variant="primary"
                            >
                                Update
                            </Button>
                        </Flex>
                    </>
                )}
            </Modal>
        </>
    )
}
