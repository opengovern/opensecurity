import { Flex } from '@tremor/react'
import { ReactNode, UIEvent } from 'react'
import Footer from './Footer'
import Sidebar from './Sidebar'
import Notification from '../Notification'
import { useNavigate } from 'react-router-dom'
import { useAtomValue, useSetAtom } from 'jotai'
import { sampleAtom } from '../../store'
type IProps = {
    children: ReactNode
    onScroll?: (e: UIEvent) => void
    scrollRef?: any
}

export default function Layout({ children, onScroll, scrollRef }: IProps) {
    const url = window.location.pathname.split('/')
    const smaple = useAtomValue(sampleAtom)
    const navigate = useNavigate()
    
    let current = url[1]
    let sub_page= false
    if (url.length > 2) {
        for (let i = 2; i < url.length; i += 1) {
            current += `/${url[i]}`
        }
        sub_page= true
    }
    const showSidebar = url[1] == "callback" ? false : true
   
    return (
        <Flex
            flexDirection="row"
            className="h-screen overflow-hidden"
            justifyContent="start"
        >
            {showSidebar && (
                <Sidebar  currentPage={current} />
            )}
            <div className="z-10 w-full h-full relative">
                <Flex
                    flexDirection="col"
                    alignItems="center"
                    justifyContent="between"
                    className={`bg-gray-100 dark:bg-gray-900 h-screen ${
                        current === 'assistant' ? '' : 'overflow-y-scroll'
                    } overflow-x-hidden`}
                    id="platform-container"
                    onScroll={(e) => {
                        if (onScroll) {
                            onScroll(e)
                        }
                    }}
                    ref={scrollRef}
                >
                    <Flex
                        justifyContent="center"
                        className={`${
                            current === 'assistant'
                                ? 'h-fit'
                                : 'sm:px-6 px-2  sm:mt-16 mt-4 h-fit '
                        } ${showSidebar && ' 2xl:px-24'} ${
                            'sm:mt-6'
                        } `}
                        // pl-44
                    >
                        <div
                            className={`w-full  ${
                                current === 'dashboard' ? '' : ''
                            } ${
                                current === 'assistant'
                                    ? 'w-full max-w-full'
                                    : 'py-6'
                            }`}
                        >
                            <>
                             
                                {children}
                            </>
                        </div>
                    </Flex>
                    <Footer />
                </Flex>
            </div>
            <Notification />
        </Flex>
    )
}
