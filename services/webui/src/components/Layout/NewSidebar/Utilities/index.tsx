import { Flex } from '@tremor/react'
import Profile from './Profile'

interface IUtilities {
    isCollapsed: boolean
}

export default function Utilities({ isCollapsed }: IUtilities) {
    return (
        <Flex
            flexDirection="col"
            alignItems="start"
            justifyContent="start"
            className="p-2 gap-0.5  border-t-gray-700 h-fit min-h-fit"
        >
            {/* {!isCollapsed && <Text className="my-2 !text-xs">UTILITIES</Text>}
            <JobsMenu isCollapsed={isCollapsed} workspace={workspace} />
            <CLIMenu isCollapsed={isCollapsed} /> */}
            <Profile isCollapsed={isCollapsed} />
        </Flex>
    )
}
