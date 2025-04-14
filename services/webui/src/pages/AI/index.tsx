import { useEffect, useState } from "react"
import Agents from "../../components/AIComponents/Agents"
import AIChat from "./chat/AIChat"
import { ArrowLeftIcon, ArrowRightIcon } from "@heroicons/react/24/outline"
import { RiArrowLeftLine, RiArrowRightLine } from "@remixicon/react"
import { AppLayout, Button, Header, Modal } from '@cloudscape-design/components';
import Cal, { getCalApi } from '@calcom/embed-react'
import { Flex } from '@tremor/react';
import AgentSelection from "../../components/AIComponents/AgentSelection"

export default function AI() {
    const [isOpen, setIsOpen] = useState(true)
    const [isSideBarOpen, setIsSidebarOpen] = useState(true)
    const [showAgentSelection, setShowAgentSelection] = useState(false)
    const [toolsHide, setToolsHide] = useState(false)
    const selected_agent = JSON.parse(localStorage.getItem('agent') || '{}')
    useEffect(()=>{
        if (selected_agent?.id) {
            setShowAgentSelection(false)
        } else {
            setShowAgentSelection(true)
        }

    },[selected_agent])


    return (
        <>
            <AppLayout
                navigationOpen={isSideBarOpen}
                onNavigationChange={(event) => {
                    setIsSidebarOpen(event.detail.open)
                }}
                navigationHide={showAgentSelection}
                toolsHide={true}
                navigation={<Agents />}
                content={
                    <>{showAgentSelection ? <AgentSelection onClick={(agent :any)=>{
                        localStorage.setItem('agent', JSON.stringify(agent))
                        setShowAgentSelection(false)
                    }}   /> : <AIChat />}</>
                }
            />
        </>
    )
}