import React, { useEffect, useState } from 'react'
import { Button, Card, Flex, Text, Title } from '@tremor/react'
import { PlusIcon } from '@heroicons/react/24/solid'
import { useAuthApiV1KeyDeleteDelete, useAuthApiV1KeysList } from '../../../api/auth.gen'
import Spinner from '../../../components/Spinner'
import DrawerPanel from '../../../components/DrawerPanel'
import CreateAPIKey from './CreateAPIKey'
import APIKeyRecord from './APIKeyRecord'
import Notification from '../../../components/Notification'
import TopHeader from '../../../components/Layout/Header'
import {
    Alert,
    Box,
    Header,
    KeyValuePairs,
    Link,
    Modal,
    RadioGroup,
    SpaceBetween,
    Table,
    Toggle,
} from '@cloudscape-design/components'
import KButton from '@cloudscape-design/components/button'
import { TrashIcon } from '@heroicons/react/24/outline'
import axios from 'axios'
export default function SettingsWorkspaceAPIKeys() {
    const [drawerOpen, setDrawerOpen] = useState<boolean>(false)
    const [drawerOpenEdit, setDrawerOpenEdit] = useState<boolean>(false)

    const [deletModalOpen, setDeleteModalOpen] = useState<boolean>(false)
    const [editLoading,setEditLoading] = useState(false)
    const [selectedItem, setSelectedItem] = useState<any>()

   const {
       response: responseDelete,
       isLoading: deleteIsLoading,
       isExecuted: deleteIsExecuted,
       error,
       sendNow: callDelete,
   } = useAuthApiV1KeyDeleteDelete((selectedItem?.id || 0).toString(), {}, false)


    const { response, isLoading, sendNow } = useAuthApiV1KeysList()
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

const fixRole = (role: string) => {
    switch (role) {
        case 'admin':
            return 'Admin'
        case 'editor':
            return 'Editor'
        case 'viewer':
            return 'Viewer'
        default:
            return role
    }
}   
  
    const openCreateMenu = () => {
        setDrawerOpen(true)
    }
  
       useEffect(() => {
            if(deleteIsExecuted && !deleteIsLoading && error!='') {
                sendNow()
                setDeleteModalOpen(false)
            }
       }, [responseDelete,deleteIsExecuted,deleteIsLoading])
       
     const EditKey = () => {
         
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
             is_active: selectedItem?.active,
             role: selectedItem?.role_name,
         }
         axios
             .put(
                 `${url}/main/auth/api/v1/key/${selectedItem?.id}`,
                 body,
                 config
             )
             .then((res) => {
                    setEditLoading(false)
                    sendNow()
                    setDrawerOpenEdit(false)
             })
             .catch((err) => {
                 console.log(err)
                    setEditLoading(false)

             })
     }
   
    return (
        <>
            {/* <TopHeader /> */}

            <Modal
                visible={drawerOpen}
                header="Create new API Key"
                onDismiss={() => {
                    setDrawerOpen(false)
                }}
            >
                <CreateAPIKey
                    close={() => {
                        setDrawerOpen(false)
                        sendNow()
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
                            className="gap-4 w-full"
                        >
                            <KeyValuePairs
                                columns={2}
                                className="w-full"
                                items={[
                                    {
                                        label: 'Name',
                                        value: selectedItem?.name,
                                    },
                                    {
                                        label: 'Creator User',
                                        value:
                                            selectedItem?.creator_user_id?.split(
                                                '|'
                                            )[1] || '',
                                    },
                                ]}
                            />
                            <KeyValuePairs
                                className="w-full"
                                columns={2}
                                items={[
                                    {
                                        label: 'Status',
                                        value: '',
                                    },
                                    {
                                        label: '',
                                        value: (
                                            <Toggle
                                                onChange={({ detail }) =>
                                                    setSelectedItem({
                                                        ...selectedItem,
                                                        active: detail.checked,
                                                    })
                                                }
                                                checked={selectedItem?.active}
                                            >
                                                {selectedItem?.active
                                                    ? 'Active'
                                                    : 'Inactive'}
                                            </Toggle>
                                        ),
                                    },
                                ]}
                            />
                        </Flex>
                        {/* <Flex className="mt-2">
                            <Text className=" font-bold text-black text-l">
                                Status
                            </Text>
                        </Flex> */}
                        <Flex
                            justifyContent="between"
                            alignItems="start"
                            flexDirection='col'
                            className="truncate space-x-4 gap-2 mt-4 mb-4"
                        >
                            <Text className=" font-bold text-black text-l">
                                Role
                            </Text>

                            <div className="space-y-5 sm:mt-0">
                                <RadioGroup
                                    onChange={({ detail }) =>
                                        setSelectedItem({
                                            ...selectedItem,
                                            role_name: detail.value,
                                        })
                                    }
                                    value={selectedItem.role_name}
                                    items={roleItems.map((item) => {
                                        return {
                                            value: item.value,
                                            label: item.title,
                                            description: item.description,
                                        }
                                    })}
                                />
                            </div>
                        </Flex>
                        <Flex className="w-full justify-end mt-2 gap-3">
                            <KButton
                                loading={deleteIsLoading && deleteIsExecuted}
                                disabled={deleteIsExecuted && deleteIsLoading}
                                onClick={(event) => {
                                    setDeleteModalOpen(true)
                                    setDrawerOpenEdit(false)
                                }}
                            >
                                <TrashIcon className="h-5 w-5" color="rose" />
                            </KButton>
                            <KButton
                                loading={editLoading}
                                onClick={() => EditKey()}
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
                header="Delete API Key"
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
                                callDelete()
                            }}
                            variant="primary"
                        >
                            Delete
                        </KButton>
                    </Flex>
                }
            >
                <>{`Are you sure you want to delete key ${selectedItem?.name}?`}</>
                {error && error !== '' && (
                    <>
                        <Alert header="failed" type="error">
                            Failed to delete API Key
                        </Alert>
                    </>
                )}
            </Modal>
            <Table
                className="mt-2"
                onRowClick={(event) => {
                    const row = event.detail.item
                    if (row) {
                        setSelectedItem(row)
                        setDrawerOpenEdit(true)
                    }
                }}
                columnDefinitions={[
                    {
                        id: 'name',
                        header: 'Key Name',
                        cell: (item: any) => item.name,
                    },
                    {
                        id: 'role',
                        header: 'Role Name & Key',
                        cell: (item: any) => (
                            <Flex
                                alignItems="start"
                                justifyContent="start"
                                flexDirection="col"
                                className="w-1/4"
                            >
                                <Text className="text-sm font-medium">
                                    {fixRole(item.role_name || '')}
                                </Text>
                                <Text className="text-xs">
                                    {item.maskedKey?.replace('...', '*******')}
                                </Text>
                            </Flex>
                        ),
                    },
                    {
                        id: 'created_by',
                        header: 'Created by',
                        cell: (item: any) =>
                            item.creator_user_id?.split('|')[1] || '',
                    },
                    {
                        id: 'create_date',
                        header: 'Create Date',
                        cell: (item: any) => (
                            <Flex
                                justifyContent="start"
                                className="truncate w-full"
                            >
                                {new Date(
                                    Date.parse(
                                        item.created_at || Date.now().toString()
                                    )
                                ).toLocaleDateString()}
                            </Flex>
                        ),
                    },
                    {
                        id: 'active',
                        header: 'Active',
                        cell: (item: any) => (
                            <Flex
                                justifyContent="start"
                                className="truncate w-full"
                            >
                                {item.active ? 'Yes' : 'No'}
                            </Flex>
                        ),
                    },
                    {
                        id: 'action',
                        header: 'Actions',
                        cell: (item: any) => <></>,
                    },
                ]}
                columnDisplay={[
                    { id: 'name', visible: true },
                    { id: 'role', visible: true },
                    { id: 'created_by', visible: true },
                    { id: 'created_date', visible: true },
                    { id: 'active', visible: true },

                    // { id: 'action', visible: true },

                    // { id: 'action', visible: true },
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
                        actions={
                            <>
                                <KButton
                                    className="float-right"
                                    variant="primary"
                                    onClick={() => {
                                        openCreateMenu()
                                    }}
                                >
                                    Create API Key
                                </KButton>
                            </>
                        }
                        className="w-full"
                    >
                        API Keys{' '}
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
