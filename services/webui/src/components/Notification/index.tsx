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
import { Alert } from '@cloudscape-design/components'

export default function Notification() {
    const [notif, setNotif] = useAtom(notificationAtom)

    useEffect(() => {
        const timer = setTimeout(() => {
            setNotif({ text: undefined, type: undefined })
        }, 5000)
        return () => clearTimeout(timer)
    }, [notif.text])


    return (
      
        <div className={`${!!notif.text && !!notif.type ? 'block' : 'hidden'}`}>
            <Alert
                dismissible
                onDismiss={() => {
                    setNotif({ text: undefined, type: undefined })
                }}
                type={notif.type}
            >
                {notif.text}
            </Alert>
        </div>
    )
}
