import { useAtomValue } from 'jotai'
import { useEffect, useState } from 'react'
import { Button, Flex, Text, TextInput } from '@tremor/react'
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline'
import { Checkbox, useCheckboxState } from 'pretty-checkbox-react'
import { PlatformEngineServicesIntegrationApiEntityListConnectionsSummaryResponse } from '../../../../../../api/api'
import Spinner from '../../../../../Spinner'
import { isDemoAtom } from '../../../../../../store'

export const compareArrays = (a: any[], b: any[]) => {
    if (a && b) {
        return (
            a.length === b.length &&
            a.every((element: any, index: number) => element === b[index])
        )
    }
    return undefined
}

interface IOthers {
    value: string[] | undefined
    defaultValue: string[]
    data:
        | PlatformEngineServicesIntegrationApiEntityListConnectionsSummaryResponse
        | undefined
    condition: string
    onChange: (o: string[]) => void
}

export default function CloudAccounts({
    value,
    defaultValue,
    condition,
    data,
    onChange,
}: IOthers) {
    const [con, setCon] = useState(condition)
    const [search, setSearch] = useState('')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const checkbox = useCheckboxState({ state: [...value] })
    const isDemo = useAtomValue(isDemoAtom)

    useEffect(() => {
        if (
            condition !== con ||
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            !compareArrays(value?.sort() || [], checkbox.state.sort())
        ) {
            if (condition === 'is') {
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-ignore
                onChange([...checkbox.state])
            }
            if (condition === 'isNot') {
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-ignore
                const arr = data?.connections
                    ?.filter(
                        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                        // @ts-ignore
                        (x) => x.id !== checkbox.state.includes(x.id)
                    )
                    .map((x) => x.id)
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-ignore
                onChange(arr)
            }
            setCon(condition)
        }
    }, [checkbox.state, condition])

    return (
        <Flex flexDirection="col" justifyContent="start" alignItems="start">
            <TextInput
                icon={MagnifyingGlassIcon}
                placeholder="Search..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="mb-4"
            />
            <Flex
                flexDirection="col"
                justifyContent="start"
                alignItems="start"
                className="gap-1.5 max-h-[200px] overflow-y-scroll no-scroll max-w-full"
            >
                {data ? (
                    data.connections
                        ?.sort((a, b) => {
                            if (value?.includes(a.id || '')) {
                                return -1
                            }
                            if (value?.includes(b.id || '')) {
                                return 1
                            }
                            return 0
                        })
                        ?.filter(
                            (d) =>
                                d.providerConnectionName
                                    ?.toLowerCase()
                                    .includes(search.toLowerCase()) ||
                                d.providerConnectionID
                                    ?.toLowerCase()
                                    .includes(search.toLowerCase())
                        )
                        .map(
                            (d, i) =>
                                i < 100 && (
                                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                    // @ts-ignore
                                    <Checkbox
                                        shape="curve"
                                        className="!items-start"
                                        value={d.id}
                                        {...checkbox}
                                    >
                                        <Flex
                                            flexDirection="col"
                                            alignItems="start"
                                        >
                                            <Text
                                                className={`${
                                                    isDemo ? 'blur-sm' : ''
                                                } text-gray-800 truncate`}
                                            >
                                                {d.providerConnectionName}
                                            </Text>
                                            <Text
                                                className={`${
                                                    isDemo ? 'blur-sm' : ''
                                                } text-xs truncate max-w-[200px]`}
                                            >
                                                {d.providerConnectionID}
                                            </Text>
                                        </Flex>
                                    </Checkbox>
                                )
                        )
                ) : (
                    <Spinner />
                )}
            </Flex>
            {/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
            {/* @ts-ignore */}
            {!compareArrays(value?.sort(), defaultValue?.sort()) && (
                <Flex className="pt-3 mt-3 border-t border-t-gray-200">
                    <Button
                        variant="light"
                        onClick={() => {
                            onChange(defaultValue)
                            checkbox.setState(defaultValue)
                        }}
                    >
                        Reset
                    </Button>
                </Flex>
            )}
        </Flex>
    )
}
