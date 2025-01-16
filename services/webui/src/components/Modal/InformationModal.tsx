import { Fragment, useState } from 'react'
import { Dialog, Transition } from '@headlessui/react'
import { CheckIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline'
import { Button, Flex } from '@tremor/react'
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/solid'
import Modal from '.'

interface InformationModalProps {
    title: string
    description: string | JSX.Element
    open: boolean
    successful?: boolean
    okButton?: string
    onClose?: () => void
}
const InformationModal: React.FC<InformationModalProps> = ({
    open,
    onClose,
    successful = true,
    title,
    description,
    okButton = 'OK',
}) => {
    const closeHandler = () => {
        if (onClose) {
            onClose()
        }
    }

    const descriptionElem = () => {
        if (description) {
            if (typeof description === 'string') {
                return (
                    <div className="mt-2">
                        <p className="text-sm text-gray-500">{description}</p>
                    </div>
                )
            }
            return description
        }
        return undefined
    }

    return (
        <Modal open={open} onClose={closeHandler}>
            <Flex flexDirection="row" justifyContent="start">
                {successful ? (
                    <CheckCircleIcon className="w-5 h-5 text-green-600" />
                ) : (
                    <XCircleIcon className="w-5 h-5 text-red-600" />
                )}
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
                className="text-gray-600 font-normal mt-4"
            >
                {descriptionElem()}
            </Flex>
            <Flex
                justifyContent="end"
                alignItems="end"
                flexDirection="row"
                className="mt-5 sm:mt-6"
            >
                <Button
                    variant="primary"
                    className="ml-1"
                    onClick={closeHandler}
                >
                    {okButton}
                </Button>
            </Flex>
        </Modal>
    )
}

export default InformationModal
