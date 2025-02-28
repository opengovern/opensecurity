import { Fragment, ReactNode } from 'react'
import { Dialog, Transition } from '@headlessui/react'

interface IModal {
    open: boolean
    onClose: () => void
    children: ReactNode
}

export default function Modal({ open, onClose, children }: IModal) {
    return (
        <Transition.Root show={open} as={Fragment}>
            <Dialog
                as="div"
                className="relative z-10"
                onClose={() => {
                    onClose()
                }}
            >
                <Transition.Child
                    as={Fragment}
                    enter="ease-out duration-300"
                    enterFrom="opacity-0"
                    enterTo="opacity-100"
                    leave="ease-in duration-200"
                    leaveFrom="opacity-100"
                    leaveTo="opacity-0"
                >
                    <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
                </Transition.Child>

                <div
                    className="fixed mx-auto inset-0 z-50 overflow-y-auto"
                    style={{ width: '1000px' }}
                >
                    <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                        <Transition.Child
                            as={Fragment}
                            enter="ease-out duration-300"
                            enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                            enterTo="opacity-100 translate-y-0 sm:scale-100"
                            leave="ease-in duration-200"
                            leaveFrom="opacity-100 translate-y-0 sm:scale-100"
                            leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        >
                            <Dialog.Panel className="bg-gray-50 relative transform overflow-hidden rounded-lg px-4 pb-4 pt-5 text-left shadow-xl transition-all lg:max-w-5xl md:max-w-3xl sm:my-8 sm:w-full sm:max-w-sm sm:p-6">
                                {children}
                            </Dialog.Panel>
                        </Transition.Child>
                    </div>
                </div>
            </Dialog>
        </Transition.Root>
    )
}
