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
    const hasTop = () => {
        if (current) {
            if (
                (current.includes('incidents') ||
                    current.includes('integrations') ||
                    current.includes('dashboard') ||
                    current.includes('compliance') ||
                    current.includes('score') ||
                    current.includes('administration')) &&
                !sub_page
            ) {
                return false
            }
            return true
        } else {
            return true
        }
    }
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
                    id="kaytu-container"
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
                                : 'px-12 mt-16 h-fit '
                        } ${showSidebar && '2xl:px-48 '} ${
                            hasTop() ? 'mt-16' : 'mt-6'
                        } `}
                        // pl-44
                    >
                        <div
                            className={`w-full ${
                                current === 'dashboard' ? '' : ''
                            } ${
                                current === 'assistant'
                                    ? 'w-full max-w-full'
                                    : 'py-6'
                            }`}
                        >
                            <>
                                {/* {sampleAtom && (
                                    <p
                                        onClick={() => {
                                            navigate(
                                                `ws/${workspace}/settings/about`
                                            )
                                        }}
                                        className=" cursor-pointer  left-1/4 rounded-xl w-1/2 bg-[#FFEED4] border-[#ff9900] border-solid border-[1px] text-[#ff9900]  p-4  absolute top-0 mt-1"
                                    >
                                        Sample data has been loaded. Please
                                        purge it before adding your data.{' '}
                                    </p>
                                )} */}

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
