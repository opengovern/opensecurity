import { Button, Flex, Title } from '@tremor/react'

import Utilities from './Utilities'

import { OpenGovernanceBig } from '../../../icons/icons'

export default function TopHeader() {
 

   

    // document.title = `${mainPage()} `

    return (
        <div className="    w-full flex  items-center justify-center  bg-[#0f1b2a] ">
            <Flex className=" flex-row items-center justify-between w-full">
                <Flex className="w-full">
                    <OpenGovernanceBig />
                </Flex>
                <Flex className="w-full flex-row items-center justify-end">
                    <Utilities />
                </Flex>
            </Flex>
        </div>
    )
}
