import { ChevronRightIcon, PlayCircleIcon } from '@heroicons/react/24/outline'
import { PlayIcon } from '@heroicons/react/24/solid'
import {
    ComponentType,
    FunctionComponent,
    ReactNode,
    useEffect,
    useRef,
    useState,
} from 'react'

interface CardProps {
    title: string
    logos?: string[]
    tag?: string
    onClick?: () => void
    option?: string
    description?: string
}

const UseCaseCard: FunctionComponent<CardProps> = ({
    title,

    tag,
    onClick,
    option,
    description,
    logos,
}) => {
    const truncate = (text: string | undefined, number: number) => {
        if (text) {
            return text.length > number
                ? text.substring(0, number) + '...'
                : text
        }
    }

    return (
        <>
            <div
                onClick={() => {
                    onClick?.()
                }}
                className="card cursor-pointer rounded-lg border-2 border-slate-100 shadow-xl dark:border-none dark:bg-white h-full flex flex-col justify-start  py-4 w-full gap-2 "
            >
                <div className="flex flex-row justify-between rounded-xl  items-center px-4 ">
                    <div className="flex flex-row gap-2">
                        {logos?.map((logo) => {
                            return (
                                <div className=" bg-gray-300 dark:bg-slate-400 rounded p-2 flex items-center justify-center">
                                    <img
                                        src={logo}
                                        className=" h-5 w-5"
                                        onError={(e) => {
                                            e.currentTarget.onerror = null
                                            e.currentTarget.src =
                                                'https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/default.svg'
                                        }}
                                    />
                                </div>
                            )
                        })}
                    </div>
                    <div className="flex flex-row justify-center    rounded-xl p-2 items-center">
                        <div className="flex w-full text-white flex-row justify-center items-center ">
                            {/* <span>Run it</span> */}
                            <PlayIcon className="w-5" color="black" />
                        </div>
                    </div>
                </div>
                <div className=" text-start flex flex-col gap-1 text-black text-wrap px-4 h-full  ">
                    <span className=" text-base  text-ellipsis overflow-hidden w-full  text-nowrap">
                        {title}
                    </span>
                    <span className="text-sm text-gray-500 text-ellipsis overflow-hidden w-full  text-nowrap">
                        {description}
                    </span>
                </div>

                {/* <div className="flex flex-row justify-center w-full bg-openg-950 dark:bg-blue-900 rounded-b-lg px-4 py-2 items-center">
                  <div className="flex w-full text-white flex-row justify-center items-center gap-2">
                      <span>Run it</span>
                      <ChevronRightIcon className="w-5" color="white" />
                  </div>
              </div> */}
            </div>
        </>
    )
}

export default UseCaseCard
