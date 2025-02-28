import { Dialog } from '@headlessui/react'
import { Button, Flex } from '@tremor/react'
import { useState } from 'react'
import { XCircleIcon } from '@heroicons/react/24/solid'
import Modal from '.'

interface IConfirmModal {
    open: boolean
    onConfirm: () => void
    onClose?: () => void
    title?: string
    icon?: JSX.Element
    description?: string
    yesButton?: string
    noButton?: string
}
export default function ConfirmModal({
    open,
    onClose,
    onConfirm,
    title = 'Are you sure?',
    icon = <XCircleIcon className="w-5 h-5 text-red-600" />,
    description,
    yesButton = 'Yes',
    noButton = 'No',
}: IConfirmModal) {
    const [loading, setLoading] = useState<boolean>(false)
    if (!open && loading) {
        setLoading(false)
    }
    return (
        <Modal
            open={open}
            onClose={() => {
                if (onClose) {
                    onClose()
                }
            }}
        >
            <Flex flexDirection="row" justifyContent="start">
                {icon}
                <Dialog.Title
                    as="h3"
                    className="ml-3 text-base font-medium leading-6 text-gray-900"
                >
                    {title}
                </Dialog.Title>
            </Flex>
            <Flex
                flexDirection="row"
                justifyContent="start"
                className="text-gray-600 font-normal"
            >
                {description && (
                    <div className="mt-2">
                        <p className="text-sm text-gray-500">{description}</p>
                    </div>
                )}
            </Flex>
            <Flex
                justifyContent="end"
                alignItems="end"
                flexDirection="row"
                className="mt-5 sm:mt-6 gap-3"
            >
                <Button
                    variant="secondary"
                    loading={loading}
                    // className="w-1/2"
                    onClick={onClose}
                >
                    {noButton}
                </Button>
                <Button
                    variant="primary"
                    loading={loading}
                    onClick={() => {
                        setLoading(true)
                        onConfirm()
                        if (onClose) {
                            onClose()
                        }
                    }}
                >
                    {yesButton}
                </Button>
            </Flex>
        </Modal>
    )
}
