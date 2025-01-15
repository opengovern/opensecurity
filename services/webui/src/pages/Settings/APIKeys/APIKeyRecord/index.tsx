import { Flex, Text } from '@tremor/react'
import { TrashIcon } from '@heroicons/react/24/outline'
import { useEffect, useState } from 'react'
import { useSetAtom } from 'jotai'
import { PlatformEnginePkgAuthApiWorkspaceApiKey } from '../../../../api/api'
import {
    useAuthApiV1KeyDeleteDelete,
    useAuthApiV1UserDetail,
} from '../../../../api/auth.gen'
import ConfirmModal from '../../../../components/Modal/ConfirmModal'
import InformationModal from '../../../../components/Modal/InformationModal'
import Spinner from '../../../../components/Spinner'
import { notificationAtom } from '../../../../store'

interface APIKeyRecordProps {
    item: PlatformEnginePkgAuthApiWorkspaceApiKey
    refresh: () => void
}

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

export default function APIKeyRecord({ item, refresh }: APIKeyRecordProps) {
    const [deleteConfirmation, setDeleteConfirmation] = useState<boolean>(false)
    const { response, isLoading } = useAuthApiV1UserDetail(
        item.creator_user_id || ''
    )
    const setNotification = useSetAtom(notificationAtom)
    const {
        response: responseDelete,
        isLoading: deleteIsLoading,
        isExecuted: deleteIsExecuted,
        error,
        sendNow: callDelete,
    } = useAuthApiV1KeyDeleteDelete((item.id || 0).toString(), {}, false)

    useEffect(() => {
        if (!deleteIsLoading && deleteIsExecuted) {
            setNotification({
                text: 'API Key successfully deleted',
                type: 'success',
            })
            refresh()
        }
    }, [deleteIsLoading])

    return (
        <>
            <ConfirmModal
                title="Delete API Key"
                description={`Are you sure you want to delete key ${item.name}?`}
                open={deleteConfirmation}
                yesButton="Delete"
                noButton="Cancel"
                onConfirm={callDelete}
                onClose={() => setDeleteConfirmation(false)}
            />
            <InformationModal
                title="Delete Failed"
                description="Failed to delete API Key"
                successful={false}
                open={error !== undefined}
            />
            <Flex
                justifyContent="start"
                flexDirection="row"
                className="mb-4 py-2 border-b"
            >
                <Text className="text-sm mt-1 w-1/4">{item.name}</Text>
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

                {isLoading ? (
                    <Flex justifyContent="start" className="w-1/4">
                        <Spinner />
                    </Flex>
                ) : (
                    <Text className="text-base w-1/4">{response?.email}</Text>
                )}

                <Flex
                    justifyContent="between"
                    flexDirection="row"
                    className="w-1/4"
                >
                    <Text className="text-base">
                        {new Date(
                            Date.parse(item.created_at || Date.now().toString())
                        ).toLocaleDateString()}
                    </Text>
                    {deleteIsLoading && deleteIsExecuted ? (
                        <Spinner />
                    ) : (
                        <TrashIcon
                            className="w-4 h-4 cursor-pointer"
                            onClick={() => {
                                setDeleteConfirmation(true)
                            }}
                        />
                    )}
                </Flex>
            </Flex>
        </>
    )
}
