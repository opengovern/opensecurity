import { useState } from "react"
import Agents from "../../components/AIComponents/Agents"
import AIChat from "./chat/AIChat"
import { ArrowLeftIcon, ArrowRightIcon } from "@heroicons/react/24/outline"
import { RiArrowLeftLine, RiArrowRightLine } from "@remixicon/react"
import { Button, Modal } from '@cloudscape-design/components';
import Cal, { getCalApi } from '@calcom/embed-react'
import { Flex } from '@tremor/react';

export default function AI() {
    const [isOpen, setIsOpen] = useState(true)
    const is_ai_page = true
        const [open, setOpen] = useState(false)
        const [openCal, setOpenCal] = useState(false)


    return (
        <>
            <div className=" sm:px-12 px-2      w-full flex  items-start justify-center flex-row    overflow-x-hidden  ">
                <div
                    className={` rounded-xl border-slate-500 p-4 pt-8      bg-slate-200 dark:bg-gray-950  w-full max-w-[75rem]  px-2 ${
                        is_ai_page
                            ? ` grid grid-cols-12 border ${
                                  isOpen ? '   gap-12' : '  gap-1'
                              }`
                            : 'flex  items-start justify-center flex-row'
                    }   max-h-[90vh] overflow-x-hidden `}
                >
                    {is_ai_page && (
                        <div
                            className={`sm:inline-block hidden  bg-slate-200 ${
                                isOpen
                                    ? 'col-span-4 border-r-2 border-slate-500'
                                    : 'col-span-1'
                            } dark:bg-gray-950  w-full max-w-sm pr-2 h-full max-h-[75vh]  overflow-hidden `}
                        >
                            <div className="w-full">
                                {isOpen ? (
                                    <>
                                        <span
                                            className="text-slate-950 dark:text-slate-200 w-full justify-end flex pr-2 cursor-pointer"
                                            onClick={() => setIsOpen(false)}
                                        >
                                            <RiArrowLeftLine className="w-5" />
                                        </span>
                                    </>
                                ) : (
                                    <>
                                        <span
                                            className="text-slate-950 dark:text-slate-200 w-full justify-start flex pr-2 cursor-pointer"
                                            onClick={() => setIsOpen(true)}
                                        >
                                            <RiArrowRightLine className="w-5" />
                                        </span>
                                    </>
                                )}
                            </div>
                            {isOpen && <Agents />}
                        </div>
                    )}
                    <div
                        className={`w-full overflow-x-scroll relative max-h-[75vh] min-h-[75vh] ${
                            is_ai_page &&
                            `${isOpen ? 'col-span-8' : 'col-span-10'}`
                        } `}
                    >
                        <AIChat />
                    </div>
                </div>
            </div>

            <Modal
                size="medium"
                visible={open}
                onDismiss={() => setOpen(false)}
                header="Not available"
            >
                <Flex className="flex-col gap-2">
                    <span>
                        {' '}
                        This feature is only available on commercial version.
                    </span>
                    <Button
                        onClick={() => {
                            setOpenCal(true)
                        }}
                    >
                        Contact us
                    </Button>
                </Flex>
            </Modal>
            <Modal
                size="large"
                visible={openCal}
                onDismiss={() => setOpenCal(false)}
                header="Not available"
            >
                <Cal
                    namespace="try-enterprise"
                    calLink="team/clearcompass/try-enterprise"
                    style={{
                        width: '100%',
                        height: '100%',
                        overflow: 'scroll',
                    }}
                    config={{ layout: 'month_view' }}
                />
            </Modal>
        </>
    )
}