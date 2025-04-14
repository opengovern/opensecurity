import axios from 'axios';
import { useEffect, useState } from 'react';
import { Agent } from './types';
import { useNavigate } from 'react-router-dom';
import LoadingDots from '../Loading';
import Tooltip from '../Tooltip';
import { Button, Modal } from '@cloudscape-design/components'
import Cal, { getCalApi } from '@calcom/embed-react'
import { Flex } from '@tremor/react'
function Agents() {
      const [agents, setAgents] = useState<Agent[]>([])
      const [loading, setLoading] = useState(true)
      const selected_agent = JSON.parse(localStorage.getItem('agent') || '{}')
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
      
      return (
        <>
            <div className="  #bg-slate-200 #dark:bg-gray-950      h-full w-full max-w-sm  justify-start items-start  max-h-[90vh]  flex flex-col gap-2 ">
                <div
                    id="k-agent-bar"
                    className="flex flex-col gap-2 max-h-[90vh] overflow-y-scroll mt-2 "
                >
                    {agents?.map((Fagent) => {
                        return (
                            <div
                                key={Fagent.id}
                                onClick={() => {
                                    localStorage.setItem(
                                        'agent',
                                        JSON.stringify(Fagent)
                                    )
                                    window.location.reload()
                                }}
                                className={`rounded-sm flex flex-col w-full justify-start items-start gap-2 hover:dark:bg-gray-700 hover:bg-gray-400 cursor-pointer p-2 ${
                                    selected_agent?.id == Fagent.id &&
                                    ' bg-slate-400 dark:bg-slate-800'
                                }`}
                            >
                                <span className="text-base text-slate-950 dark:text-slate-200">
                                    {Fagent.name}
                                </span>
                                <span className="text-sm text-slate-800 dark:text-slate-400">
                                    {Fagent.description}
                                </span>
                            </div>
                        )
                    })}
                </div>
            </div>
        </>
    )
}

export default Agents;
