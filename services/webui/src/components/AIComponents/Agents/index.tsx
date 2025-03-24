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
    const [agents, setAgents] = useState<Agent[]>([
        {
            name: 'Identity & Access',
            description:
                'Delivers information on user identities, access controls, and activity.',
            welcome_message:
                'Hi there! This is your Identity & Access Agent. I can help you with anything related to identity management and access tools. What can I assist you with today? For example, you can ask me things like:',
            sample_questions: [
                'Get me the list of users who have access to Azure Subscriptions.',
                'Get me all SPNs with expired passwords.',
                'Show me the access activity for user John Doe.',
            ],
            id: 'identity_access',
            enabled: true,
            is_available: true,
        },
        {
            name: 'DevOps',
            description:
                'Provides data on secure code, deployment, and automation workflows.',
            welcome_message:
                'Hello! This is your DevOps Agent. I can provide insights into secure code, deployment, and automation workflows. How can I assist you today? For instance, you could ask me:',
            sample_questions: [
                'What are the latest secure code scan results?',
                'Show me the deployment status for the production environment.',
                'Provide a report on automated workflow execution times.',
            ],
            id: 'devops',
            enabled: true,
            is_available: true,
        },
    ])
 
    const [agent, setAgent] = useState<Agent | null>(localStorage.getItem('agent') ? JSON.parse(localStorage.getItem('agent') as string) : agents[0])
    return (
        <>
            <div className="  bg-slate-200 dark:bg-gray-950      h-full w-full max-w-sm  justify-start items-start  max-h-[90vh]  flex flex-col gap-2 ">
                <div
                    id="k-agent-bar"
                    className="flex flex-col gap-2 max-h-[90vh] overflow-y-scroll mt-2 "
                >
                    {agents?.map((Fagent) => {
                        return (
                            <div
                                key={Fagent.id}
                                onClick={() => {
                                    localStorage.setItem('agent', JSON.stringify(Fagent))
                                    window.location.reload()
                                }}
                                className={`rounded-sm flex flex-col justify-start items-start gap-2 hover:dark:bg-gray-700 hover:bg-gray-400 cursor-pointer p-2 ${
                                    agent?.id == Fagent.id &&
                                    ' bg-slate-400 dark:bg-slate-800'
                                }`}
                            >
                                <span className="text-base text-slate-950 dark:text-slate-200">
                                    {Fagent.name}
                                </span>
                                <span className="text-sm text-slate-500 dark:text-slate-400">
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
