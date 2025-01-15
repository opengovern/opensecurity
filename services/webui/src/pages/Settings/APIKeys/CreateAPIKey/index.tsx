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
import { Modal } from '@cloudscape-design/components'

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

export default function CreateAPIKey({ close }: CreateAPIKeyProps) {
    const [apiKeyName, setApiKeyName] = useState<string>('')
    const [showCopied, setShowCopied] = useState<boolean>(false)
    const [copyOpen, setCopyOpen] = useState<boolean>(false)

    const [role, setRole] = useState<string>('viewer')
    const [roleValue, setRoleValue] = useState<
        'admin' | 'editor' | 'viewer' | undefined
    >()
    const setNotification = useSetAtom(notificationAtom)

    const {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow: callCreate,
    } = useAuthApiV1KeyCreateCreate(
        { name: apiKeyName, role: roleValue },
        {},
        false
    )

    useEffect(() => {
        if (isExecuted && !isLoading) {
            setNotification({
                text: 'API key successfully added',
                type: 'success',
            })
            // close()
            setCopyOpen(true)
        }
        if (error) {
            setNotification({
                text: 'Unable to add new API key',
                type: 'error',
            })
        }
    }, [isLoading, error])

    useEffect(() => {
        if (role === 'viewer' || role === 'editor' || role === 'admin') {
            setRoleValue(role)
        }
    }, [role])

    return (
        <Flex flexDirection="col" justifyContent="between" className="h-full">
            <Modal
                header={error === undefined ? 'Successful' : 'Failed'}
                visible={copyOpen}
                onDismiss={() => {
                    close()
                    setCopyOpen(false)
                    setApiKeyName('')
                }}
            >
                <>
                    {error === undefined ? (
                        <Flex
                            flexDirection="col"
                            justifyContent="start"
                            alignItems="start"
                        >
                            API key created, copy the key and keep it safe:
                            <Card
                                className="w-full cursor-pointer mt-2"
                                onClick={() => {
                                    setShowCopied(true)
                                    setTimeout(() => setShowCopied(false), 2000)
                                    clipboardCopy(response?.token || '')
                                }}
                            >
                                <Flex
                                    flexDirection="row"
                                    justifyContent="between"
                                >
                                    <div className="w-full break-all">
                                        {response?.token}
                                    </div>
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="start"
                                        className="h-5 w-5"
                                    >
                                        <DocumentDuplicateIcon className="h-5 w-5 text-openg-600 " />
                                        <Text
                                            className={`${
                                                showCopied ? '' : 'hidden'
                                            } absolute mt-6 bg-openg-600 text-white rounded-md p-1`}
                                        >
                                            Copied!
                                        </Text>
                                    </Flex>
                                </Flex>
                            </Card>
                        </Flex>
                    ) : (
                        `Failed to create the API Key`
                    )}
                </>
            </Modal>
            <List className="mt-4 h-full">
                {/* <ListItem>
                    <Text className="text-gray-900 font-medium py-2">
                        Properties
                    </Text>
                </ListItem> */}
                <ListItem>
                    <Flex>
                        <Text className="w-1/3 font-medium text-gray-800 py-2">
                            API Key Name *
                        </Text>
                        <TextInput
                            className="w-2/3"
                            onChange={(p) => {
                                setApiKeyName(p.target.value)
                            }}
                        />
                    </Flex>
                </ListItem>
                <ListItem>
                    <Flex alignItems="start">
                        <Text className="w-1/3 te font-medium text-gray-800 py-2">
                            Role *
                        </Text>
                        <Flex
                            flexDirection="col"
                            alignItems="start"
                            className="w-2/3 space-y-5"
                        >
                            {roleItems.map((item) => {
                                return (
                                    <Flex>
                                        <input
                                            name="roles"
                                            type="radio"
                                            className="h-4 w-4"
                                            onClick={() => {
                                                setRole(item.value)
                                            }}
                                            checked={item.value === role}
                                        />
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                            className="pl-7"
                                        >
                                            <Text className="font-medium text-gray-900">
                                                {item.title}
                                            </Text>
                                            <Subtitle className="text-gray-500">
                                                {item.description}
                                            </Subtitle>
                                        </Flex>
                                    </Flex>
                                )
                            })}
                        </Flex>
                    </Flex>
                </ListItem>
            </List>
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
                    disabled={apiKeyName.length === 0}
                    onClick={() => {
                        callCreate()
                    }}
                    loading={isExecuted && isLoading}
                >
                    Create API Key
                </KButton>
            </Flex>
        </Flex>
    )
}
