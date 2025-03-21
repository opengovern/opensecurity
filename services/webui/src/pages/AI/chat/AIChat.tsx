import React, { ChangeEvent, useEffect, useRef, useState } from 'react';
import { Chat, ChatList } from '../types';
import axios from 'axios';
import { dateTimeDisplay } from '../../../utilities/dateDisplay';
import KChatCard from '../../../components/AIComponents/ChatCard';
import KResponseCard from '../../../components/AIComponents/ResponseCard';
import KInput from '../../../components/AIComponents/Input'
import { Button, Modal } from '@cloudscape-design/components';
import Cal, { getCalApi } from '@calcom/embed-react'
import { Flex } from '@tremor/react';
function AIChat() {
    const [message, setMessage] = useState('')
    const [open,setOpen] = useState(false)
    const [openCal, setOpenCal] = useState(false)


    const [chats, setChats] = useState<ChatList>({
        '0': {
            message: '',
            text: 'Hi there! This is your Identity & Access Agent. I can help you with anything related to identity management and access tools. What can I assist you with today? For example, you can ask me things like:',
            loading: false,
            time: 0,
            error: '',
            isWelcome: true,
            pre_loaded: false,
            clarify_needed: false,
            messageTime: '',
            responseTime: '1:5AM',
            suggestions: [
                'Get me the list of users who have access to Azure Subscriptions.',
                'Get me all SPNs with expired passwords.',
                'Show me the access activity for user John Doe.',
            ],
            response: {},
        },
    })
   
    const lastMessageRef = useRef(null)
    const scroll = () => {
        const layout = document.getElementById('layout')
        if (layout) {
            const start = layout.scrollTop
            const end = layout.scrollHeight
            const duration = 1500 // Adjust duration in milliseconds
            let startTime: any = null
            const animateScroll = (timestamp: any) => {
                if (!startTime) startTime = timestamp
                const progress = Math.min((timestamp - startTime) / duration, 1)
                layout.scrollTop = start + (end - start) * progress
                if (progress < 1) {
                    requestAnimationFrame(animateScroll)
                }
            }
            requestAnimationFrame(animateScroll)
            // layout.scrollTop = layout?.scrollHeight+400;
        }
        //  if (lastMessageRef.current) {
        //   // @ts-ignore
        //    lastMessageRef.current.scrollIntoView({ behavior: "smooth" });
        //  }
    }

    useEffect(() => {
        scroll()
    }, [chats])

    return (
        <>
            <div className=" bg-slate-200 dark:bg-gray-950 flex max-h-[65vh] flex-col  justify-start   items-start w-full ">
                <div
                    id="layout"
                    className=" flex justify-start  items-start overflow-y-auto  w-full  bg-slate-200 dark:bg-gray-950 pt-2  "
                >
                    <div className="  w-full relative ">
                        <section className="chat-section h-full     flex flex-col relative gap-8 w-full max-w-[95%]   ">
                            {chats &&
                                Object.keys(chats).map((key) => {
                                    return (
                                        <>
                                            {!chats[key].isWelcome && (
                                                <KChatCard
                                                    date={
                                                        chats[key].messageTime
                                                    }
                                                    key={parseInt(key) + 'chat'}
                                                    message={chats[key].message}
                                                />
                                            )}
                                            <KResponseCard
                                                key={parseInt(key) + 'result'}
                                                ref={
                                                    key ===
                                                    (
                                                        Object.keys(chats)
                                                            ?.length - 1
                                                    ).toString()
                                                        ? lastMessageRef
                                                        : null
                                                }
                                                scroll={scroll}
                                                response={chats[key].response}
                                                loading={chats[key].loading}
                                                pre_loaded={
                                                    chats[key].pre_loaded
                                                }
                                                chat_id={chats[key].id}
                                                error={chats[key].error}
                                                time={chats[key].time}
                                                text={chats[key].text}
                                                isWelcome={chats[key].isWelcome}
                                                date={chats[key].responseTime}
                                                clarify_needed={
                                                    chats[key].clarify_needed
                                                }
                                                clarify_questions={
                                                    chats[key].clarify_questions
                                                }
                                                id={''}
                                                suggestions={
                                                    chats[key].suggestions
                                                }
                                                onClickSuggestion={(
                                                    suggestion: string
                                                ) => {}}
                                            />
                                        </>
                                    )
                                })}
                        </section>
                    </div>
                </div>
                <KInput
                    value={message}
                    chats={chats}
                    onChange={(e: any) => {
                        setMessage(e?.target?.value)
                    }}
                    onSend={() => {
                        setOpen(true)
                    }}
                />
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

export default AIChat;
