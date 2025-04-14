import React, { useEffect, useState } from 'react'
import { Button, Card, Flex, Text, Title } from '@tremor/react'
import { PlusIcon } from '@heroicons/react/24/solid'
import {
    useAuthApiV1KeyDeleteDelete,
    useAuthApiV1KeysList,
} from '../../../api/auth.gen'
import Spinner from '../../../components/Spinner'
// import CreateAPIKey from './CreateAPIKey'
// import APIKeyRecord from './APIKeyRecord'
import Notification from '../../../components/Notification'
import TopHeader from '../../../components/Layout/Header'
import {
    Alert,
    Box,
    Header,
    Input,
    KeyValuePairs,
    Link,
    Modal,
    RadioGroup,
    Select,
    SpaceBetween,
    Table,
    Toggle,
} from '@cloudscape-design/components'
import KButton from '@cloudscape-design/components/button'
import { TrashIcon } from '@heroicons/react/24/outline'
import axios from 'axios'
import CreateConnector from './CreateConnector'
import { useSetAtom } from 'jotai'
import { notificationAtom } from '../../../store'
import { dateTimeDisplay } from '../../../utilities/dateDisplay'
export default function SettingsConnectors() {
    const [drawerOpen, setDrawerOpen] = useState<boolean>(false)
    const [drawerOpenEdit, setDrawerOpenEdit] = useState<boolean>(false)
    const [deletModalOpen, setDeleteModalOpen] = useState<boolean>(false)
    const [editLoading, setEditLoading] = useState(false)
    const [selectedItem, setSelectedItem] = useState<any>()
    const [response, setResponse] = useState<any>([])
    const [isLoading, setIsLoading] = useState(false)
    const [isDeleteLoading, setIsDeleteLoading] = useState(false)
    const [error, setError] = useState('')
    const [editError, setEditError] = useState<any>()
    const setNotification = useSetAtom(notificationAtom)

    const openCreateMenu = () => {
        setDrawerOpen(true)
    }

    const EditConnector = () => {
        if (
            !selectedItem.client_id ||
            !selectedItem.client_secret ||
            !selectedItem.sub_type
        ) {
            setEditError('Please fill all the fields')
            return
        }
        if (
            (selectedItem.sub_type?.value === 'entraid' ||
                selectedItem.sub_type === 'entraid') &&
            !selectedItem.tenant_id
        ) {
            setEditError('Please fill all the fields')
            return
        }
        setEditLoading(true)
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
            connector_type: selectedItem?.type,
            connector_sub_type: selectedItem?.sub_type?.value ? selectedItem?.sub_type?.value : selectedItem?.sub_type,
            issuer: selectedItem?.issuer,
            tenant_id: selectedItem?.tenant_id,
            client_id: selectedItem?.client_id,
            client_secret: selectedItem?.client_id,
            id: selectedItem?.id,
            connector_id: selectedItem?.connector_id,
        }
        axios
            .put(`${url}/main/auth/api/v1/connector`, body, config)
            .then((res) => {
                setEditLoading(false)

                GetRows()
                setDrawerOpenEdit(false)
                setNotification({
                    type: 'success',
                    text: 'Provider updated successfully',
                })
            })
            .catch((err) => {
                console.log(err)

                var error = err.response.data.message
                if (!error) {
                    error = 'Failed to create connector'
                }
                setNotification({
                    type: 'error',
                    text: error,
                })
                setEditLoading(false)

                setEditError(error)
            })
    }
    const DeleteConnector = () => {
        setIsDeleteLoading(true)
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
            .delete(
                `${url}/main/auth/api/v1/connector/${selectedItem?.connector_id}`,

                config
            )
            .then((res) => {
                setIsDeleteLoading(false)
                GetRows()
                setDeleteModalOpen(false)
            })
            .catch((err) => {
                console.log(err)
                setIsDeleteLoading(false)
                setError('Error while deleting Provider')
            })
    }
    const GetRows = () => {
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
            .get(
                `${url}/main/auth/api/v1/connectors`,

                config
            )
            .then((res) => {
                setResponse(res.data)
                setIsLoading(false)
            })
            .catch((err) => {
                console.log(err)
                if(err.response.status === 400){
                    GetRows()
                }
                setIsLoading(false)
            })
    }
    const checkDate =(date: string) =>{
        if(date == "0001-01-01T00:00:00Z"){
            return 'Not Available'
        }
        return dateTimeDisplay(date)

    }

    useEffect(() => {
        GetRows()
    }, [])
    const FindSubtype=(value: string)=> {
        if(value === 'general'){
           return  { label: 'General', value: 'general' }
        }
        if(value === 'entraid'){
            return { label: 'Entra Id', value: 'entraid' }
        }
        if(value === 'google-workspaces'){
            return { label: 'Google Workspace', value: 'google-workspaces' }
        }

    }

    return (
        <>
            <Modal
                visible={drawerOpen}
                header="Create new OIDC Connector"
                onDismiss={() => {
                    setDrawerOpen(false)
                }}
            >
                <CreateConnector
                    close={() => {
                        setDrawerOpen(false)
                        GetRows()
                    }}
                />
            </Modal>
            <Modal
                visible={drawerOpenEdit}
                header={selectedItem?.name}
                onDismiss={() => {
                    setDrawerOpenEdit(false)
                }}
            >
                {selectedItem ? (
                    <>
                        <Flex
                            flexDirection="col"
                            justifyContent="start"
                            alignItems="start"
                            className="gap-2 w-full mb-4"
                        >
                            <KeyValuePairs
                                className="w-full"
                                columns={3}
                                items={[
                                    {
                                        label: 'ID',
                                        value: selectedItem?.connector_id,
                                    },
                                    {
                                        label: 'Name',
                                        value: selectedItem?.name,
                                    },
                                    {
                                        label: 'Type',
                                        value: selectedItem?.type.toUpperCase(),
                                    },

                                    {
                                        label: 'Created At',
                                        value: checkDate(
                                            selectedItem?.created_at
                                        ),
                                    },
                                    {
                                        label: 'Updated At',
                                        value: checkDate(
                                            selectedItem?.last_update
                                        ),
                                    },
                                    {
                                        label: 'User Count',
                                        value: selectedItem?.user_count,
                                    },
                                ]}
                            />
                            <KeyValuePairs
                                className="w-full"
                                columns={1}
                                items={[
                                    {
                                        label: 'Issuer',
                                        value: selectedItem?.issuer,
                                    },
                                ]}
                            />
                            <Select
                                selectedOption={
                                    selectedItem?.sub_type?.value
                                        ? selectedItem?.sub_type
                                        : FindSubtype(selectedItem?.sub_type)
                                }
                                onChange={({ detail }) => {
                                    setSelectedItem({
                                        ...selectedItem,
                                        sub_type: detail.selectedOption,
                                    })
                                    setEditError(null)
                                }}
                                options={[
                                    { label: 'General', value: 'general' },
                                    { label: 'Entra Id', value: 'entraid' },
                                    {
                                        label: 'Google Workspace',
                                        value: 'google-workspaces',
                                    },
                                ]}
                                placeholder="OIDC Provider"
                                inlineLabelText="OIDC Provider"
                                className="w-full"
                            />
                            <Input
                                onChange={({ detail }) => {
                                    setSelectedItem({
                                        ...selectedItem,
                                        client_id: detail.value,
                                    })
                                    setEditError(null)
                                }}
                                value={selectedItem?.client_id}
                                placeholder="Client Id"
                                className="w-full"
                            />
                            <Input
                                onChange={({ detail }) => {
                                    setSelectedItem({
                                        ...selectedItem,
                                        client_secret: detail.value,
                                    })
                                    setEditError(null)
                                }}
                                value={selectedItem?.client_secret}
                                placeholder="Client Secret"
                                className="w-full"
                            />
                            {(selectedItem?.sub_type?.value === 'general' ||
                                selectedItem?.sub_type === 'general') && (
                                <>
                                    <Input
                                        onChange={({ detail }) => {
                                            setSelectedItem({
                                                ...selectedItem,
                                                issuer: detail.value,
                                            })
                                            setEditError(null)
                                        }}
                                        value={selectedItem?.issuer}
                                        placeholder="Issuer"
                                        className="w-full"
                                    />
                                </>
                            )}
                            {(selectedItem?.sub_type?.value === 'entraid' ||
                                selectedItem?.sub_type) && (
                                <>
                                    <Input
                                        onChange={({ detail }) => {
                                            setSelectedItem({
                                                ...selectedItem,
                                                tenant_id: detail.value,
                                            })
                                            setEditError(null)
                                        }}
                                        value={selectedItem?.tenant_id}
                                        placeholder="Tenant Id"
                                        className="w-full"
                                    />
                                </>
                            )}
                        </Flex>
                        {editError && (
                            <Alert className="w-full mb-3" type="error">
                                {editError}
                            </Alert>
                        )}
                        <Flex className="w-full justify-end mt-2 gap-3">
                            <KButton
                                loading={isDeleteLoading}
                                disabled={isDeleteLoading}
                                onClick={(event) => {
                                    setDeleteModalOpen(true)
                                    setDrawerOpenEdit(false)
                                }}
                            >
                                <TrashIcon className="h-5 w-5" color="rose" />
                            </KButton>
                            <KButton
                                loading={editLoading}
                                onClick={() => EditConnector()}
                                variant="primary"
                            >
                                Update Changes
                            </KButton>
                        </Flex>
                    </>
                ) : (
                    <Spinner />
                )}
            </Modal>
            <Modal
                visible={deletModalOpen}
                header="Delete Connector"
                onDismiss={() => {
                    setDeleteModalOpen(false)
                }}
                footer={
                    <Flex className="gap-2 w-full" flexDirection="row">
                        <KButton
                            onClick={() => {
                                setDeleteModalOpen(false)
                            }}
                        >
                            Cancel
                        </KButton>
                        <KButton
                            onClick={() => {
                                DeleteConnector()
                            }}
                            variant="primary"
                        >
                            Delete
                        </KButton>
                    </Flex>
                }
            >
                <>{`Are you sure you want to delete  ${selectedItem?.name}?`}</>
                {error && error !== '' && (
                    <>
                        <Alert
                            className="mt-2 mb-2"
                            header="failed"
                            type="error"
                        >
                            Failed to delete Connector
                        </Alert>
                    </>
                )}
            </Modal>
            <Table
                className="mt-2"
                variant="full-page"
                onRowClick={(event) => {
                    const row = event.detail.item
                    if (row && row.connector_id != 'local') {
                        setSelectedItem(row)
                        setDrawerOpenEdit(true)
                    }
                }}
                resizableColumns
                columnDefinitions={[
                    {
                        id: 'id',
                        header: 'ID',
                        cell: (item: any) => item.connector_id,
                    },
                    {
                        id: 'name',
                        header: 'Name',
                        cell: (item: any) => item.name,
                    },
                    {
                        id: 'type',
                        header: 'Type',
                        cell: (item: any) => item?.type,
                    },
                    {
                        id: 'sub_type',
                        header: 'OIDC Provider',
                        cell: (item: any) => item?.sub_type,
                    },
                    {
                        id: 'client_id',
                        header: 'Client Id',
                        width: 200,
                        maxWidth: 102000,
                        cell: (item: any) => item?.client_id,
                    },
                    {
                        id: 'issuer',
                        header: 'Issuer',
                        width: 100,
                        maxWidth: 100,
                        cell: (item: any) => item?.issuer,
                    },
                    {
                        id: 'user_count',
                        header: 'User Count',
                        cell: (item: any) => item?.user_count,
                    },
                    {
                        id: 'user_count',
                        header: 'User Count',
                        cell: (item: any) => item?.user_count,
                    },
                    {
                        id: 'created_at',
                        header: 'Created At',
                        cell: (item: any) => checkDate(item?.created_at),
                    },
                    {
                        id: 'last_update',
                        header: 'Updated At',
                        cell: (item: any) => checkDate(item?.last_update),
                    },
                ]}
                columnDisplay={[
                    { id: 'id', visible: false },
                    { id: 'name', visible: true },
                    { id: 'type', visible: false },
                    { id: 'sub_type', visible: true },
                    { id: 'client_id', visible: true },
                    { id: 'issuer', visible: false },
                    { id: 'user_count', visible: true },
                    { id: 'created_at', visible: false },
                    { id: 'last_update', visible: true },
                ]}
                loading={isLoading}
                // @ts-ignore
                items={response ? response : []}
                empty={
                    <Box
                        margin={{ vertical: 'xs' }}
                        textAlign="center"
                        color="inherit"
                    >
                        <SpaceBetween size="m">
                            <b>No resources</b>
                            {/* <Button>Create resource</Button> */}
                        </SpaceBetween>
                    </Box>
                }
                header={
                    <Header
                        description={
                            <>
                                Opensecurity supports any OIDC-compatible
                                Identity Provider
                            </>
                        }
                        actions={
                            <>
                                <KButton
                                    className="float-right"
                                    variant="primary"
                                    onClick={() => {
                                        openCreateMenu()
                                    }}
                                >
                                    Add SSO Provider
                                </KButton>
                            </>
                        }
                        className="w-full"
                    >
                        SSO Configuration{' '}
                    </Header>
                }
            />
            {/* <Card key="summary">
                <Flex className="mb-6">
                    <Title className="font-semibold">API Keys</Title>
                    <Button
                        className="float-right"
                        onClick={() => {
                            openCreateMenu()
                        }}
                        icon={PlusIcon}
                    >
                        Create API Key
                    </Button>
                </Flex>
                <Flex
                    justifyContent="start"
                    flexDirection="row"
                    className="mb-6"
                >
                    <Text className="w-1/4">Key Name</Text>
                    <Text className="w-1/4">Role Name & Key</Text>
                    <Text className="w-1/4">Created by</Text>
                    <Text className="w-1/4">Create Date</Text>
                </Flex>
                {response?.map((item) => (
                    <APIKeyRecord
                        item={item}
                        refresh={() => {
                            sendNow()
                        }}
                    />
                ))}
            </Card> */}
        </>
    )
}
