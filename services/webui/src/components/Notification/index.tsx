import { Fragment, useEffect } from 'react'
import { Transition } from '@headlessui/react'
import { Flex } from '@tremor/react'
import { useAtom } from 'jotai'
import {
    ExclamationCircleIcon,
    QuestionMarkCircleIcon,
    XMarkIcon,
    CheckCircleIcon,
} from '@heroicons/react/24/outline'
import { notificationAtom } from '../../store'

export default function Notification() {
    const [notif, setNotif] = useAtom(notificationAtom)

    useEffect(() => {
        const timer = setTimeout(() => {
            setNotif({ text: undefined, type: undefined })
        }, 5000)
        return () => clearTimeout(timer)
    }, [notif.text])

    const color = () => {
        switch (notif.type) {
            case 'success':
                return 'text-emerald-500 bg-emerald-50 ring-emerald-100'
            case 'warning':
                return 'text-amber-500 bg-amber-50 ring-amber-100'
            case 'error':
                return 'text-rose-500 bg-rose-50 ring-rose-100'
            default:
                return 'text-openg-500 bg-openg-50 ring-openg-100'
        }
    }

    const icon = () => {
        switch (notif.type) {
            case 'success':
                return <CheckCircleIcon className="h-6" />
            case 'warning':
                return <QuestionMarkCircleIcon className="h-6" />
            case 'error':
                return <XMarkIcon className="h-6" />
            default:
                return <ExclamationCircleIcon className="h-6" />
        }
    }

    return (
        <Transition
            show={!!notif.text && !!notif.type}
            as={Fragment}
            enter="transform ease-out duration-300 transition"
            enterFrom="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
            enterTo="translate-y-0 opacity-100 sm:translate-x-0"
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
        >
            <Flex
                onClick={() => setNotif({ text: undefined, type: undefined })}
                justifyContent="start"
                className={`gap-2 z-50 fixed ${
                    notif.position === 'topRight' ||
                    notif.position === 'bottomRight'
                        ? 'right-12'
                        : 'left-12'
                } ${
                    notif.position === 'topRight' ||
                    notif.position === 'topLeft'
                        ? 'top-24'
                        : 'bottom-24'
                } w-full max-w-sm p-4 rounded-md shadow-md ring-1 ${color()}`}
            >
                {icon()}
                <span>{notif.text}</span>
            </Flex>
        </Transition>
    )
}
