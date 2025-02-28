import { Flex } from '@tremor/react'
import Profile from './Profile'



export default function Utilities() {
    return (
        <Flex
            flexDirection="row"
            alignItems="end"
            justifyContent="end"
            className="p-2 gap-0.5  border-t-gray-700 h-fit min-h-fit"
        >
          
            <Profile  />
        </Flex>
    )
}
