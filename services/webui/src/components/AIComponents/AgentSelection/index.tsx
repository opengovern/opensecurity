import axios from 'axios';
import { useEffect, useState } from 'react';
import { Agent } from './types';
import { useNavigate } from 'react-router-dom';
import { RiLineFill, RiLock2Line, RiUserLine } from '@remixicon/react';
import LoadingDots from '../Loading';
import Tooltip from '../Tooltip';

interface AgentProps {
  onClick: Function
}

function AgentSelection({ onClick } : AgentProps) {
    const [agents, setAgents] = useState<Agent[]>([])
    const [loading, setLoading] = useState(true)
    const navigate = useNavigate()

    const GetAgents = () => {
        setLoading(true)
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
            .get(`${url}/main/core/api/v4/chatbot/agents `, config)
            .then((res) => {
                //  const temp = []
                setAgents(res.data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }

    useEffect(() => {
        GetAgents()
    }, [])

    const truncate = (text: string | undefined, number: number) => {
        if (text) {
            return text.length > number
                ? text.substring(0, number) + '...'
                : text
        }
    }

    return (
        <>
            {loading ? (
                <>
                    <div className="mt-2  bg-slate-200 dark:bg-gray-950   max-h-[80vh] min-h-[80vh] overflow-x-hidden  h-full w-full justify-center items-center flex flex-row gap-2">
                        <span className="text-2xl font-semibold text-slate-950 dark:text-slate-200">
                            Loading
                        </span>
                        <LoadingDots />
                    </div>
                </>
            ) : (
                <>
                    <div className="  bg-slate-200 dark:bg-gray-950    h-full w-full min-h-[60vh]  justify-center items-center  flex flex-col gap-2 pt-8">
                        <div className="grid grid-cols-12 gap-6 w-fit p-2 ">

                            {agents?.map((agent,index) => {
                                return (
                                    <>
                                        <div
                                            key={agent.id}
                                            onClick={() => {
                                                onClick(agent)
                                            }}
                                            className={` ${
                                                index == 0 && '2xl:col-start-3 '
                                            } md:col-start-4 2xl:col-span-4   md:col-span-6 col-span-12  flex flex-col border border-slate-500 rounded-xl justify-start items-start gap-2 hover:dark:bg-gray-800 hover:bg-gray-400 cursor-pointer p-4 backdrop-filter`}
                                        >
                                            <span className="text-2xl font-semibold text-center w-full text-slate-950 dark:text-slate-100 flex flex-row justify-between items-center gap-2">
                                                <span className="w-full text-center">
                                                    {agent.name}
                                                </span>
                                                {/* <RiLock2Line className="w-4" /> */}
                                            </span>
                                            <span className="text-base  rounded-lg p-2 text-center  dark:text-slate-200 text-slate-800 min-h-[5rem] max-h-[5rem] overflow-y-auto">
                                                {truncate(
                                                    agent.description,
                                                    70
                                                )}
                                            </span>
                                        </div>
                                    </>
                                )
                            })}
                        </div>
                    </div>
                </>
            )}
        </>
    )
}

export default AgentSelection


/*     {agent.is_available ? (
                                            <>
                                                {agent.enabled ? (
                                                    <>
                                                        <div
                                                            key={agent.id}
                                                            onClick={() => {
                                                                onClick(agent)
                                                            }}
                                                            className=" 2xl:col-span-4 md:col-span-6 col-span-12  flex flex-col border border-slate-500 rounded-xl justify-start items-start gap-2 hover:dark:bg-gray-800 hover:bg-gray-400 cursor-pointer p-4 backdrop-filter"
                                                        >
                                                            <span className="text-2xl font-semibold text-center w-full text-slate-950 dark:text-slate-100 flex flex-row justify-between items-center gap-2">
                                                                <span className="w-full text-center">
                                                                    {agent.name}
                                                                </span>
                                                                {/* <RiLock2Line className="w-4" /> 
                                                            </span>
                                                            <span className="text-base  rounded-lg p-2 text-center  dark:text-slate-200 text-slate-800 min-h-[5rem] max-h-[5rem] overflow-y-auto">
                                                                {truncate(
                                                                    agent.description,
                                                                    70
                                                                )}
                                                            </span>
                                                        </div>
                                                    </>
                                                ) : (
                                                    <>
                                                        {' '}
                                                        <Tooltip
                                                            text="Agent is not available for you. "
                                                            className="sm:col-span-4 col-span-12"
                                                        >
                                                            <div
                                                                key={agent.id}
                                                                className=" sm:col-span-4 col-span-12  flex flex-col border border-slate-500 rounded-xl justify-start items-start gap-2 hover:dark:bg-gray-800 hover:bg-gray-400  p-4 backdrop-filter"
                                                            >
                                                                <span className="text-2xl font-semibold text-center w-full text-slate-950 dark:text-slate-100 flex flex-row justify-between items-center gap-2">
                                                                    <span className="w-full text-center">
                                                                        {
                                                                            agent.name
                                                                        }
                                                                    </span>
                                                                    <RiLock2Line className="w-4" />
                                                                </span>
                                                                <span className="text-base  rounded-lg p-2 text-center  dark:text-slate-200 text-slate-800 min-h-[5rem] max-h-[5rem] overflow-y-auto">
                                                                    {truncate(
                                                                        agent.description,
                                                                        70
                                                                    )}
                                                                </span>
                                                            </div>
                                                        </Tooltip>
                                                    </>
                                                )}{' '}
                                            </>
                                        ) : (
                                            <>
                                                <Tooltip
                                                    text="COMING SOON."
                                                    className="sm:col-span-4 col-span-12"
                                                >
                                                    <div
                                                        key={agent.id}
                                                        onClick={() => {}}
                                                        className=" sm:col-span-4 col-span-12  flex flex-col border border-dashed border-slate-500 rounded-xl justify-start items-start gap-2 hover:dark:bg-gray-800 hover:bg-gray-400 cursor-pointer p-4"
                                                    >
                                                        <span className="text-2xl font-semibold text-center w-full text-slate-950 dark:text-slate-500 flex flex-row justify-center items-center gap-2">
                                                            {/* <RiUserLine className="w-4" /> 
                                                            <span className="w-full text-center">
                                                                {agent.name}
                                                            </span>
                                                        </span>
                                                        <span className="text-base  rounded-lg p-2 text-center dark:text-slate-600 text-slate-800 min-h-[5rem] max-h-[5rem] overflow-y-auto">
                                                            {truncate(
                                                                agent.description,
                                                                70
                                                            )}
                                                        </span>
                                                    </div>
                                                </Tooltip>
                                            </>
                                        )} */