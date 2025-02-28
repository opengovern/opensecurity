import { useEffect, useState } from 'react'
import {
    Button,
    Divider,
    Flex,
    List,
    ListItem,
    Subtitle,
    Text,
    TextInput,
} from '@tremor/react'
import { useSetAtom } from 'jotai/index'
import { useAuthApiV1UserInviteCreate } from '../../../../api/auth.gen'
import { notificationAtom } from '../../../../store'
import KButton from '@cloudscape-design/components/button'
import axios from 'axios'
import { Select } from '@cloudscape-design/components'
interface MemberInviteProps {
    close: (refresh: boolean) => void
}

export default function MemberInvite({ close }: MemberInviteProps) {
    const [email, setEmail] = useState<string>('')
    const [password, setPassword] = useState<string>('')
    const [connectorID, setConnectorId] = useState<any>({
        value: 'local',
        label: 'Password (Built-in)',
    })
   
    const [emailError, setEmailError] = useState<string>('')
    const [providers, setProviders] = useState<any>([])


    const [role, setRole] = useState<string>('viewer')
    const [roleValue, setRoleValue] = useState<'viewer' | 'editor' | 'admin'>(
        'viewer'
    )
    const setNotification = useSetAtom(notificationAtom)

    const {
        isExecuted,
        isLoading,
        error,
        sendNow: createInvite,
    } = useAuthApiV1UserInviteCreate(
        // @ts-ignore
        {
            // @ts-ignore
            email_address: email || '',
            // @ts-ignore
            role: role,
            password: connectorID.value == 'local' ? password : '',
            is_active: true,
            connector_id: connectorID.value,
        },
        {},
        false
    )

    useEffect(() => {
        if (role === 'viewer' || role === 'editor' || role === 'admin') {
            setRoleValue(role)
        }
    }, [role])
     useEffect(() => {
        if(email && email != ""){
    if (!email.includes('@') || !email.includes('.')) {
        setEmailError('Invalid email address')
    } else {
        setEmailError('')
    }
        }
     
     }, [email])

    useEffect(() => {
        if (isExecuted && !isLoading) {
            setNotification({
                text: 'User successfully added',
                type: 'success',
            })
            close(true)
        }
        if (error) {
            setNotification({
                text: 'Unable to add new member',
                type: 'error',
            })
        }
    }, [isLoading, error])

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
                const temp =[]
                temp.push({value:'local',label:'Password (Built-in)'})
                data?.map((item: any) => {
                    temp.push({value:item.connector_id,label:item.name})
                })
                setProviders(temp)

            })
            .catch((err) => {
                console.log(err)
               
            })
    }

    useEffect(()=>{
        GetProviders()
    },[])
    return (
        <Flex flexDirection="col" justifyContent="between" className="h-full">
            <Flex justifyContent="between" className="truncate space-x-4 py-2">
                <Text className="font-medium text-gray-800">
                    Email
                    <span className="text-red-600 font-semibold">*</span>
                </Text>
                <TextInput
                    error={emailError != ''}
                    errorMessage={emailError}
                    placeholder="email"
                    className="font-medium w-1/2 text-gray-800"
                    onChange={(e) => {
                        setEmail(e.target.value)
                    }}
                />
            </Flex>

            <Divider />

            <Flex
                justifyContent="between"
                alignItems="start"
                flexDirection="row"
                className=" space-x-4 py-2"
            >
                <Text className="font-medium text-gray-800">
                    Identity Provider
                    <span className="text-red-600 font-semibold">*</span>
                </Text>
                <Select
                    className=" w-1/2    "
                    // h-[150px]
                    selectedOption={connectorID}
                    disabled={false}
                    onChange={({ detail }) =>
                        setConnectorId(detail.selectedOption)
                    }
                    placeholder="Identity Provider"
                    // @ts-ignore
                    options={providers}
                />
            </Flex>

            {connectorID.value == 'local' && (
                <>
                    <Divider />
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
                </>
            )}
            <Divider />
                <Flex
                    justifyContent="between"
                    alignItems="start"
                    className="truncate space-x-4"
                >
                    <Text className="font-medium text-gray-800">
                        Role
                        <span className="text-red-600 font-semibold">*</span>
                    </Text>

                    <div className="space-y-5 sm:mt-0 w-1/2">
                        {roleItems.map((item) => {
                            return (
                                <div className="relative flex items-start">
                                    <div className="absolute flex h-6 items-center">
                                        <input
                                            name="roles"
                                            type="radio"
                                            className="h-4 w-4 border-gray-300 text-openg-600 focus:ring-openg-700"
                                            onClick={() => {
                                                setRole(item.value)
                                            }}
                                            checked={item.value === role}
                                        />
                                    </div>
                                    <div className="pl-7 text-sm leading-6">
                                        <div className="font-medium text-gray-900">
                                            {item.title}
                                        </div>
                                        <p className="text-gray-500">
                                            {item.description}
                                        </p>
                                    </div>
                                </div>
                            )
                        })}
                    </div>
                </Flex>

            <Flex justifyContent="end" className="truncate space-x-4 mt-4">
                <KButton
                    disabled={isExecuted && isLoading}
                    onClick={() => close(false)}
                >
                    Cancel
                </KButton>
                <KButton
                    loading={isExecuted && isLoading}
                    disabled={email.length === 0}
                    onClick={() => createInvite()}
                    variant="primary"
                >
                    Add
                </KButton>
            </Flex>
        </Flex>
    )
}
