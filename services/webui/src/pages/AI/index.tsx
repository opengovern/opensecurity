import { useState } from "react"
import Agents from "../../components/AIComponents/Agents"
import AIChat from "./chat/AIChat"
import { ArrowLeftIcon, ArrowRightIcon } from "@heroicons/react/24/outline"
import { RiArrowLeftLine, RiArrowRightLine } from "@remixicon/react"
import { AppLayout, Button, Header, Modal } from '@cloudscape-design/components';
import Cal, { getCalApi } from '@calcom/embed-react'
import { Flex } from '@tremor/react';

export default function AI() {
    const [isOpen, setIsOpen] = useState(true)
    const is_ai_page = true
    const [isSideBarOpen, setIsSidebarOpen] = useState(true)
      


    return (
        <>
            <AppLayout
                navigationOpen={isSideBarOpen}
                onNavigationChange={(event) => {
                    setIsSidebarOpen(event.detail.open)
                }}
                toolsHide={true}
                navigation={<Agents />}
                content={
                    <>
                        <AIChat />
                    </>
                }
            />
        </>
    )
}