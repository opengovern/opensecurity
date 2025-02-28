import { useEffect, useState } from 'react'
import {
    Badge,
    Button,
    Flex,
    List,
    ListItem,
    MultiSelect,
    MultiSelectItem,
    Select,
    SelectItem,
    Text,
    TextInput,
} from '@tremor/react'
import { ChevronDoubleDownIcon, TrashIcon } from '@heroicons/react/24/outline'
import dayjs from 'dayjs'
import { useSetAtom } from 'jotai'
import {
    useAuthApiV1UserRoleBindingDelete,
    useAuthApiV1UserRoleBindingUpdate,
} from '../../../../api/auth.gen'
import { PlatformEnginePkgAuthApiWorkspaceRoleBinding } from '../../../../api/api'
import ConfirmModal from '../../../../components/Modal/ConfirmModal'
import { notificationAtom } from '../../../../store'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import {
    useIntegrationApiV1ConnectionsDelete,
    useIntegrationApiV1ConnectionsSummariesList,
} from '../../../../api/integration.gen'
import KButton from '@cloudscape-design/components/button'
import { Checkbox, KeyValuePairs, Modal, RadioGroup, Toggle } from '@cloudscape-design/components'
import axios from 'axios'

interface IMemberDetails {
    user?: PlatformEnginePkgAuthApiWorkspaceRoleBinding
    close: () => void
}

export default function MemberDetails({ user, close }: IMemberDetails) {
    const [deleteConfirmation, setDeleteConfirmation] = useState<boolean>(false)
    const [role, setRole] = useState<any>(user?.role_name )
    const [isActive, setIsActive] = useState<any>(
        user?.is_active 
    )

    const [password, setPassword] = useState<string>('')
    const [changePassword, setChangePassword] = useState<boolean>(false)
    const [changeProvider, setChangeProvider] = useState<boolean>(false)
    const [connectorID, setConnectorId] = useState<any>(user?.connector_id)
    const [scopedConnectionIDs, setScopedConnectionIDs] = useState<string[]>(
        user?.scopedConnectionIDs || []
    )
    const [providers, setProviders] = useState<any>([])
    const GetProviders = () => {
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
                const data = res.data
                const temp = []
                temp.push({ value: 'local', label: 'Password (Built-in)' })
                data?.map((item: any) => {
                    temp.push({ value: item.connector_id, label: item.name })
                })
                setProviders(temp)
            })
            .catch((err) => {
                console.log(err)
            })
    }
useEffect(() => {
    GetProviders()
}, [])
useEffect(() => {
   setConnectorId(user?.connector_id)
}, [user])
    const setNotification = useSetAtom(notificationAtom)

    const {
        isExecuted,
        isLoading,
        sendNow: updateRole,
    } = useAuthApiV1UserRoleBindingUpdate(
        {
            email_address: user?.email || '',
            role: role,
            is_active: isActive,
            // @ts-ignore
           password: connectorID =='local'?password : '',
            connector_id: connectorID,
        },
        {},
        false
    )

    const {
        isExecuted: deleteExecuted,
        isLoading: deleteLoading,
        sendNow: deleteRole,
    } = useAuthApiV1UserRoleBindingDelete(user?.id || 0, {}, false)


     useEffect(() => {
        setIsActive(user?.is_active)
        setRole(user?.role_name)
     }, [user])

    useEffect(() => {
        if ((isExecuted && !isLoading) || (deleteExecuted && !deleteLoading)) {
            if (isExecuted) {
                setNotification({
                    text: 'User successfully updated',
                    type: 'success',
                })
            } else {
                setNotification({
                    text: 'User successfully deleted',
                    type: 'success',
                })
            }
            close()
        }
    }, [isLoading, deleteLoading])

    if (user === undefined) {
        return <div />
    }

    const lastActivity = () => {
        if (user.last_activity === undefined) {
            return 'Never'
        }

        return dateTimeDisplay(user.last_activity)
    }

    const items = [
        {
            title: 'Email',
            value: user.email,
        },
        {
            title: 'Member Since',
            value: dateTimeDisplay(user.created_at || Date.now().toString()),
        },
        {
            title: 'Last Activity',
            value: lastActivity(),
        },
    ]

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

    const multiSelectCount = () => {
        return (
            <div className="bg-gray-200 w-9 h-9 text-center pt-2 -ml-3 rounded-l-lg">
                {scopedConnectionIDs.length}
            </div>
        )
    }
    const GetKeyItems =()=>{
        const temp =[]
         items.map((item) => {
             temp.push( {
                 label: item.title,
                 value: item.value,
             })
         })
         temp.push({
             label: 'User Status',
             value: (
                 <Toggle
                     onChange={({ detail }) => setIsActive(detail.checked)}
                     checked={isActive}
                 >
                     {isActive ? 'Active' : 'Inactive'}
                 </Toggle>
             ),
         })
         return temp
    }



    return (
        <>
            <Modal
                visible={deleteConfirmation}
                header="Delete user"
                footer={
                    <Flex justifyContent="end" className="gap-2">
                        <KButton onClick={() => setDeleteConfirmation(false)}>
                            Cancel
                        </KButton>
                        <KButton
                            loading={deleteExecuted && deleteLoading}
                            disabled={isExecuted && isLoading}
                            onClick={deleteRole}
                            variant="primary"
                        >
                            Delete
                        </KButton>
                    </Flex>
                }
                onDismiss={() => setDeleteConfirmation(false)}
            >
                <>{`Are you sure you want to delete ${user.email}?`}</>
            </Modal>

            {/* <ConfirmModal
                title="Delete user"
                description={`Are you sure you want to delete ${user.userName}?`}
                open={deleteConfirmation}
                yesButton="Delete"
                noButton="Cancel"
                onConfirm={deleteRole}
                onClose={() => setDeleteConfirmation(false)}
            /> */}
            <Flex
                flexDirection="col"
                justifyContent="between"
                className="h-full"
            >
                <KeyValuePairs
                    className="w-full"
                    columns={2}
                    // @ts-ignore
                    items={GetKeyItems()}
                />
                <Flex
                    justifyContent="between"
                    alignItems="start"
                    flexDirection="col"
                    className="truncate space-x-4 gap-2 mt-4 mb-4"
                >
                    <Text className=" font-bold text-black text-l">Role</Text>

                    <div className="space-y-5 sm:mt-0">
                        <RadioGroup
                            onChange={({ detail }) => setRole(detail.value)}
                            value={role}
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

                <Flex
                    justifyContent="start"
                    alignItems="start"
                    flexDirection="col"
                    className="mt-4 w-full mb-4 gap-2 space-x-4"
                >
                    <>{console.log(connectorID)}</>
                    <Select
                        className=" w-1/2 z-50 static  "
                        // h-[150px]
                        value={connectorID}
                        disabled={false}
                        onValueChange={setConnectorId}
                        placeholder="Identity Provider"
                    >
                        {providers.map((item: any) => {
                            return (
                                <SelectItem key={item.value} value={item.value}>
                                    {item.label}
                                </SelectItem>
                            )
                        })}
                        {/* <SelectItem className="static" value="oicd">
                                OIDC (SSO)
                            </SelectItem> */}
                    </Select>
                </Flex>
                {connectorID === 'local' && (
                    <Flex
                        justifyContent="start"
                        alignItems="start"
                        flexDirection="col"
                        className="mt-4 w-full mb-4 gap-2 space-x-4"
                    >
                        <Checkbox
                            onChange={({ detail }) =>
                                setChangePassword(detail.checked)
                            }
                            checked={changePassword}
                        >
                            Change Password
                        </Checkbox>
                        {changePassword && (
                            <TextInput
                                type="password"
                                placeholder="password"
                                className="font-medium w-1/2 text-gray-800"
                                onChange={(e) => {
                                    setPassword(e.target.value)
                                }}
                            />
                        )}
                    </Flex>
                )}
                {/* <List className="pt-4"> */}
                {/* <ListItem key="password">
                        <Flex
                            justifyContent="between"
                            className="truncate space-x-4 py-2"
                        >
                            <Text className="font-medium text-gray-800">
                                Password
                                <span className="text-red-600 font-semibold">
                                    *
                                </span>
                            </Text>
                            <TextInput
                                type="password"
                                placeholder="password"
                                className="font-medium w-1/2 text-gray-800"
                                onChange={(e) => {
                                    setPassword(e.target.value)
                                }}
                            />
                        </Flex>
                    </ListItem> */}

                {/*
                    <ListItem key="item" className="py-4">
                        <Flex
                            justifyContent="between"
                            alignItems="start"
                            className="truncate space-x-4"
                        >
                            <Text className="font-medium text-gray-500">
                                Scoped account access
                            </Text>

                            <Flex
                                flexDirection="col"
                                className="w-2/3"
                                justifyContent="start"
                                alignItems="end"
                            >
                                <MultiSelect
                                    disabled={isConnectionListLoading}
                                    className="w-96 absolute"
                                    value={scopedConnectionIDs}
                                    onValueChange={(value) =>
                                        setScopedConnectionIDs(value)
                                    }
                                    placeholder="All connections"
                                    icon={multiSelectCount}
                                >
                                    {connectionList?.connections?.map(
                                        (connection) => (
                                            <MultiSelectItem
                                                value={connection.id || ''}
                                            >
                                                <Flex
                                                    justifyContent="end"
                                                    className="truncate w-full"
                                                >
                                                    <div className="truncate p-1">
                                                        <Text className="truncate font-medium text-gray-800">
                                                            {
                                                                connection.providerConnectionID
                                                            }
                                                        </Text>
                                                        <Text className="truncate text-xs text-gray-400">
                                                            {
                                                                connection.providerConnectionName
                                                            }
                                                        </Text>
                                                    </div>
                                                </Flex>
                                            </MultiSelectItem>
                                        )
                                    )}
                                </MultiSelect>
                            </Flex>
                        </Flex>
                    </ListItem>
                    */}
                {/* </List> */}
                <Flex justifyContent="end" className="truncate space-x-4">
                    <KButton disabled={user?.id ==1 || user?.email?.includes("admin")} onClick={() => setDeleteConfirmation(true)}>
                        <TrashIcon className="h-5 w-5" color="rose" />
                    </KButton>

                    <KButton
                        loading={isExecuted && isLoading}
                        disabled={deleteExecuted && deleteLoading}
                        onClick={() => updateRole()}
                        variant="primary"
                    >
                        Update Changes
                    </KButton>
                </Flex>
            </Flex>
        </>
    )
}
