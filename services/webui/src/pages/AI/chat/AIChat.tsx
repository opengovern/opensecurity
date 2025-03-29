import React, { ChangeEvent, useEffect, useRef, useState } from 'react'
import { Chat, ChatList } from '../types'
import axios from 'axios'
import { dateTimeDisplay } from '../../../utilities/dateDisplay'
import KChatCard from '../../../components/AIComponents/ChatCard'
import KResponseCard from '../../../components/AIComponents/ResponseCard'
import KInput from '../../../components/AIComponents/Input'
import { DEVOPS, IDENTITY } from './responses'

function AIChat({ setOpen, }: any) {
    const [message, setMessage] = useState('')
    const agent = JSON.parse(localStorage.getItem('agent') as string)
    const [chats, setChats] = useState<ChatList>(
        agent?.id == 'devops' ? DEVOPS : IDENTITY
    )

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
                                            {chats[key].show && (
                                                <>
                                                    {!chats[key].isWelcome && (
                                                        <KChatCard
                                                            date={
                                                                chats[key]
                                                                    .messageTime
                                                            }
                                                            key={
                                                                parseInt(key) +
                                                                'chat'
                                                            }
                                                            message={
                                                                chats[key]
                                                                    .message
                                                            }
                                                        />
                                                    )}
                                                    <KResponseCard
                                                        key={
                                                            parseInt(key) +
                                                            'result'
                                                        }
                                                        ref={
                                                            key ===
                                                            (
                                                                Object.keys(
                                                                    chats
                                                                )?.length - 1
                                                            ).toString()
                                                                ? lastMessageRef
                                                                : null
                                                        }
                                                        scroll={scroll}
                                                        response={
                                                            chats[key].response
                                                        }
                                                        loading={
                                                            chats[key].loading
                                                        }
                                                        pre_loaded={
                                                            chats[key]
                                                                .pre_loaded
                                                        }
                                                        chat_id={chats[key].id}
                                                        error={chats[key].error}
                                                        time={chats[key].time}
                                                        text={chats[key].text}
                                                        isWelcome={
                                                            chats[key].isWelcome
                                                        }
                                                        date={
                                                            chats[key]
                                                                .responseTime
                                                        }
                                                        clarify_needed={
                                                            chats[key]
                                                                .clarify_needed
                                                        }
                                                        clarify_questions={
                                                            chats[key]
                                                                .clarify_questions
                                                        }
                                                        id={''}
                                                        suggestions={
                                                            chats[key]
                                                                .suggestions
                                                        }
                                                        onClickSuggestion={(
                                                            suggestion: string
                                                        ) => {
                                                            // find suggestoin index
                                                            const sug =
                                                                chats['0']
                                                                    .suggestions
                                                            const index =
                                                                sug?.indexOf(
                                                                    suggestion
                                                                )

                                                            const temp = chats
                                                            if (
                                                                index !==
                                                                undefined
                                                            ) {
                                                                temp[
                                                                    (
                                                                        index +
                                                                        1
                                                                    )?.toString()
                                                                ].show = true
                                                                setChats({
                                                                    ...temp,
                                                                })
                                                            }
                                                        }}
                                                    />
                                                </>
                                            )}
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
        </>
    )
}

export default AIChat
