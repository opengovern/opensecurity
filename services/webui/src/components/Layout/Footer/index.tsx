import { Flex, Text } from '@tremor/react'
import { useAtomValue } from 'jotai'
import { sampleAtom } from '../../../store'

export default function Footer() {
    const smaple = useAtomValue(sampleAtom)

    return (
        <Flex
            justifyContent="center"
            className="px-12 py-3 border-t border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900 shadow-sm"
        >
            <Flex
                flexDirection="row"
                justifyContent="center"
                className="max-w-7xl w-full"
            >
                {/* eslint-disable-next-line jsx-a11y/anchor-is-valid */}
                <Text>
                    © 2025 open governance Inc. All rights
                    reserved.
                </Text>{' '}
                {/* {smaple && <>{" "}
                Demo data loaded</>} */}
                {/* eslint-disable-next-line jsx-a11y/anchor-is-valid */}
            </Flex>
        </Flex>
    )
}
