import React, { ChangeEvent, useEffect, useRef, useState } from 'react'
import { Chat, ChatList } from '../types'
import axios from 'axios'
import { dateTimeDisplay } from '../../../utilities/dateDisplay'
import KChatCard from '../../../components/AIComponents/ChatCard'
import KResponseCard from '../../../components/AIComponents/ResponseCard'
import KInput from '../../../components/AIComponents/Input'
import { DEVOPS, IDENTITY } from './responses'

function AIChat({ setOpen }: any) {
    const [message, setMessage] = useState('')
    const agent = JSON.parse(localStorage.getItem('agent') || '{}')

    const [chats, setChats] = useState<ChatList>()
    const [loading, setLoading] = useState(false)
    const [id, setId] = useState<string | undefined>(undefined)
    const [clarifying, setClarifying] = useState<boolean>(false)

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
    const RunQuery = (id: string, len: number) => {
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        const body = {
            chat_id: id,
        }

        axios
            .post(`${url}/main/core/api/v4/chatbot/run-query`, body, config)
            .then((res) => {
                if (res.data) {
                    const output = res?.data
                    if (output) {
                        setChats((prevChats) => {
                            const newChats = { ...prevChats }
                            newChats[`${len}`] = {
                                ...newChats[`${len}`],
                                response: output.result,
                                time: output.time_taken,
                                suggestions: output.suggestions?.suggestion,
                                text: output?.primary_interpration,
                                clarify_needed: false,
                                responseTime: `${
                                    new Date().getHours() > 12
                                        ? new Date().getHours() - 12
                                        : new Date().getHours()
                                }:${new Date().getMinutes()}${
                                    new Date().getHours() > 12 ? 'PM' : 'AM'
                                }`,
                                loading: false, // Ensure loading is set to false
                            }
                            return newChats
                        })
                        setLoading(false)
                        scroll()
                    } else {
                        setChats((prevChats) => {
                            const newChats = { ...prevChats }
                            newChats[`${len}`] = {
                                ...newChats[`${len}`],

                                loading: false, // Ensure loading is set to false
                                responseTime: `${
                                    new Date().getHours() > 12
                                        ? new Date().getHours() - 12
                                        : new Date().getHours()
                                }:${new Date().getMinutes()}${
                                    new Date().getHours() > 12 ? 'PM' : 'AM'
                                }`,
                            }
                            return newChats
                        })
                        setLoading(false)
                        scroll()
                    }
                } else {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                    setLoading(false)
                    scroll()
                }
            })
            .catch((err) => {
                console.log(err)
                if (err.response.data.error) {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            error: err?.response?.data?.error,
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                } else {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            error: 'Error in fetching data',
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                }

                setLoading(false)
                scroll()
            })
    }
    const GenerateQuery = (message: string, len: number) => {
        const body = {
            question: message,
            agent: agent.id,
            session_id: localStorage.getItem(`${agent.id}_session_id`),
            in_clarification_state: false,
        }
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        axios
            .post(
                `${url}/main/core/api/v4/chatbot/generate-query`,
                body,
                config
            )
            .then((res) => {
                if (res?.data) {
                    const output = res?.data
                    if (output) {
                        setId(output.chat_id)

                        if (output?.result?.type == 'CLARIFICATION_NEEDED') {
                            setClarifying(true)
                            setChats((prevChats) => {
                                const newChats = { ...prevChats }
                                newChats[`${len}`] = {
                                    ...newChats[`${len}`],
                                    clarify_needed: true,
                                    id: output.chat_id,
                                    clarify_questions:
                                        output?.result?.clarifying_questions,
                                    loading: false, // Ensure loading is set to false
                                }
                                return newChats
                            })
                        } else {
                            setClarifying(false)
                           
                            RunQuery(output?.chat_id, len)
                        }
                        setLoading(false)

                        scroll()
                    } else {
                        setChats((prevChats) => {
                            const newChats = { ...prevChats }
                            newChats[`${len}`] = {
                                ...newChats[`${len}`],

                                loading: false, // Ensure loading is set to false
                                responseTime: `${
                                    new Date().getHours() > 12
                                        ? new Date().getHours() - 12
                                        : new Date().getHours()
                                }:${new Date().getMinutes()}${
                                    new Date().getHours() > 12 ? 'PM' : 'AM'
                                }`,
                            }
                            return newChats
                        })
                        setLoading(false)
                        scroll()
                    }
                } else {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                    setLoading(false)
                    scroll()
                }
            })
            .catch((err) => {
                console.log(err)
                if (err.response.data.error) {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            error: err.response.data.error,
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                } else {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            error: 'Error in fetching data',
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                }

                setLoading(false)
                scroll()
            })
    }
    const ClarifyQuery = (answer: string, len: number) => {
        const body = {
            answer: answer,
            chat_id: id,
        }

        axios
            .post(
                `https://slay-router-latest.onrender.com/generate/clarify/`,
                body
            )
            .then((res) => {
                if (res?.data) {
                    const output = res?.data
                    if (output) {
                        setId(output.id)
                        if (output?.result?.result == 'CLARIFICATION_NEEDED') {
                            setClarifying(true)
                            setChats((prevChats) => {
                                const newChats = { ...prevChats }
                                newChats[`${len}`] = {
                                    ...newChats[`${len}`],
                                    clarify_needed: true,
                                    clarify_questions:
                                        output?.result?.clarifying_questions,
                                    loading: false, // Ensure loading is set to false
                                }
                                return newChats
                            })
                        } else {
                            setClarifying(false)

                            RunQuery(output.id, len)
                        }
                        setLoading(false)
                        scroll()
                    } else {
                        setChats((prevChats) => {
                            const newChats = { ...prevChats }
                            newChats[`${len}`] = {
                                ...newChats[`${len}`],

                                loading: false, // Ensure loading is set to false
                                responseTime: `${
                                    new Date().getHours() > 12
                                        ? new Date().getHours() - 12
                                        : new Date().getHours()
                                }:${new Date().getMinutes()}${
                                    new Date().getHours() > 12 ? 'PM' : 'AM'
                                }`,
                            }
                            return newChats
                        })
                        setLoading(false)
                        scroll()
                    }
                } else {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                    setLoading(false)
                    scroll()
                }
            })
            .catch((err) => {
                console.log(err)
                if (err.response.data.error) {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            error: err.response.data.error,
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                } else {
                    setChats((prevChats) => {
                        const newChats = { ...prevChats }
                        newChats[`${len}`] = {
                            ...newChats[`${len}`],

                            loading: false, // Ensure loading is set to false
                            error: 'Error in fetching data',
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                        }
                        return newChats
                    })
                }

                setLoading(false)
                scroll()
            })
    }
    const GetChats = () => {
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        url += `/main/core/api/v4/chatbot/session/`
        const session = localStorage.getItem(`${agent.id}_session_id`)
        if (session && session !== 'undefined') {
            url += session
        } else {
            url += `id`
        }

        url += `?agent=${agent.id}`
        axios
            .get(`${url}`, config)
            .then((res) => {
                if (res?.data) {
                    const output = res?.data
                    localStorage.setItem(`${agent.id}_session_id`, output.id)
                    const temp: ChatList = {
                        '0': {
                            message: '',
                            text: agent.welcome_message,
                            loading: false,
                            time: 0,
                            error: '',
                            isWelcome: true,
                            pre_loaded:
                                output?.chats && output.chats.length > 0
                                    ? true
                                    : false,
                            clarify_needed: false,
                            messageTime: '',
                            responseTime: `${
                                new Date().getHours() > 12
                                    ? new Date().getHours() - 12
                                    : new Date().getHours()
                            }:${new Date().getMinutes()}${
                                new Date().getHours() > 12 ? 'PM' : 'AM'
                            }`,
                            suggestions: agent.sample_questions?.map(
                                (question: any) => {
                                    return question
                                }
                            ),
                            response: {},
                        },
                    }
                    if (output?.chats && output.chats.length > 0) {
                        output.chats.forEach((chat: any) => {
                            if (chat.need_clarification) {
                                // for each in chat.clarifying_questions and for each one make new temp object
                                const clarifying_questions =
                                    chat.clarifying_questions
                                clarifying_questions.forEach(
                                    (question: any) => {
                                        temp[`${Object.keys(temp).length}`] = {
                                            id: chat.id,
                                            message: chat.question,
                                            messageTime: dateTimeDisplay(
                                                chat.created_at
                                            ),
                                            responseTime: dateTimeDisplay(
                                                chat.updated_at
                                            ),
                                            loading: false,
                                            time: chat.time_taken,
                                            error: chat.query_error,
                                            pre_loaded: true,

                                            isWelcome: false,
                                            clarify_needed: true,
                                            clarify_questions:
                                                question.question,
                                            clarify_answer: question.answer,
                                            suggestions: chat.suggestions,
                                            text: chat.primary_interpretation,
                                            response: chat.result,
                                        }
                                    }
                                )
                            } else {
                                temp[`${Object.keys(temp).length}`] = {
                                    id: chat.id,
                                    message: chat.question,
                                    messageTime: dateTimeDisplay(
                                        chat.created_at
                                    ),
                                    responseTime: dateTimeDisplay(
                                        chat.updated_at
                                    ),
                                    loading: false,
                                    time: chat.time_taken,
                                    error: chat.query_error,
                                    pre_loaded: true,
                                    isWelcome: false,
                                    clarify_needed: false,
                                    suggestions: chat.suggestions,
                                    text: chat.primary_interpretation,
                                    response: chat.result,
                                }
                            }
                        })
                        console.log('temp', temp)
                    }
                    setChats(temp)
                }
            })
            .catch((err) => {
                setLoading(false)
                scroll()
            })
    }

    useEffect(() => {
        scroll()
    }, [chats])
    useEffect(() => {
        GetChats()
    }, [])

    return (
        <>
            <div className=" #bg-slate-200 #dark:bg-gray-950 flex max-h-[90vh] flex-col  justify-start gap-4   items-start w-full ">
                <div
                    id="layout"
                    className=" flex justify-start  items-start overflow-y-auto  w-full  #bg-slate-200 #dark:bg-gray-950 pt-2  "
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
                                                id={id}
                                                suggestions={
                                                    chats[key].suggestions
                                                }
                                                onClickSuggestion={(
                                                    suggestion: string
                                                ) => {
                                                    const temp = chats
                                                    const len =
                                                        Object.keys(
                                                            chats
                                                        ).length
                                                    temp[`${len}`] = {
                                                        message: suggestion,
                                                        messageTime: `${
                                                            new Date().getHours() >
                                                            12
                                                                ? new Date().getHours() -
                                                                  12
                                                                : new Date().getHours()
                                                        }:${new Date().getMinutes()}${
                                                            new Date().getHours() >
                                                            12
                                                                ? 'PM'
                                                                : 'AM'
                                                        }`,
                                                        responseTime: '',
                                                        pre_loaded: false,
                                                        loading: true,
                                                        time: 0,
                                                        error: '',
                                                        isWelcome: false,
                                                        clarify_needed: false,
                                                        response: {
                                                            query: '',
                                                            result: undefined,
                                                        },
                                                    }
                                                    setChats(temp)
                                                    setLoading(true)
                                                    GenerateQuery(
                                                        suggestion,
                                                        len
                                                    )
                                                }}
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
                    disabled={loading}
                    onSend={() => {
                        const temp: any = chats
                        // @ts-ignore
                        const len = Object.keys(chats).length

                        setLoading(true)
                        if (clarifying) {
                            temp[`${len}`] = {
                                message: message,
                                messageTime: `${
                                    new Date().getHours() > 12
                                        ? new Date().getHours() - 12
                                        : new Date().getHours()
                                }:${new Date().getMinutes()}${
                                    new Date().getHours() > 12 ? 'PM' : 'AM'
                                }`,
                                responseTime: '',
                                loading: true,
                                clarify_needed: true,
                                pre_loaded: false,
                                time: 0,
                                error: '',
                                isWelcome: false,
                                response: {
                                    query: '',
                                    result: undefined,
                                },
                            }
                            setChats(temp)

                            ClarifyQuery(message, len)
                        } else {
                            temp[`${len}`] = {
                                message: message,
                                messageTime: `${
                                    new Date().getHours() > 12
                                        ? new Date().getHours() - 12
                                        : new Date().getHours()
                                }:${new Date().getMinutes()}${
                                    new Date().getHours() > 12 ? 'PM' : 'AM'
                                }`,
                                responseTime: '',
                                loading: true,
                                clarify_needed: false,
                                pre_loaded: false,
                                time: 0,
                                error: '',
                                isWelcome: false,
                                response: {
                                    query: '',
                                    result: undefined,
                                },
                            }
                            GenerateQuery(message, len)
                        }

                        setMessage('')
                    }}
                />
            </div>
        </>
    )
}

export default AIChat
