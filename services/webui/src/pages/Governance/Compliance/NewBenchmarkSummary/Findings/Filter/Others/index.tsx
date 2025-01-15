import { useEffect, useState } from 'react'
import { Button, Flex, Text, TextInput } from '@tremor/react'
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline'
import { Checkbox, useCheckboxState } from 'pretty-checkbox-react'
import { PlatformEnginePkgComplianceApiFindingFiltersWithMetadata } from '../../../../../../../api/api'
import Spinner from '../../../../../../../components/Spinner'
import Multiselect from '@cloudscape-design/components/multiselect'


interface IOthers {
    value: string[] | undefined
    defaultValue: string[]
    data:
        | PlatformEnginePkgComplianceApiFindingFiltersWithMetadata
        | undefined
    condition: string
    type: 'benchmarkID' | 'integrationID' | 'controlID' | 'resourceTypeID'
    onChange: (o: string[]) => void
    name: string
}

export default function Others({
    value,
    defaultValue,
    condition,
    data,
    type,
    onChange,
    name,
}: IOthers) {
    const [con, setCon] = useState(condition)
    const [search, setSearch] = useState('')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const checkbox = useCheckboxState({ state: [...value] })
   const [selectedOptions, setSelectedOptions] = useState([])

    // useEffect(() => {
    //     if (
    //         condition !== con ||
    //         // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //         // @ts-ignore
    //         !compareArrays(value?.sort() || [], checkbox.state.sort())
    //     ) {
    //         if (condition === 'is') {
    //             // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //             // @ts-ignore
    //             onChange([...checkbox.state])
    //         }
    //         if (condition === 'isNot') {
    //             // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //             // @ts-ignore
    //             const arr = data[type]
    //                 .filter(
    //                     // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //                     // @ts-ignore
    //                     (x) => !checkbox.state.includes(x.key)
    //                 )
    //                 .map((x) => x.key)
    //             // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    //             // @ts-ignore
    //             onChange(arr)
    //         }
    //         setCon(condition)
    //     }
    // }, [checkbox.state, condition])

    useEffect(() => {
        if (selectedOptions.length === 0) {
            onChange(defaultValue)
            return
        } else {
            // @ts-ignore
            const temp = []
            selectedOptions.map((o) => {
                // @ts-ignore

                temp.push(o.value)
            })
            // @ts-ignore
            onChange(temp)
            // @ts-ignore
        }
    }, [selectedOptions])
    return (
        <>
            <Multiselect
                // @ts-ignore
                selectedOptions={selectedOptions}
                tokenLimit={1}
                onChange={({ detail }) =>
                    // @ts-ignore
                    setSelectedOptions(detail.selectedOptions)
                }
                options={
                    data
                        ? data[type]?.map((d) => {
                              return {
                                  label: d.displayName,
                                  value: d.key,
                                  description: d.key,
                              }
                          })
                        : []
                }
                // @ts-ignore
                loading={!data ? true : false}
                loadingText="Please Wait"
                filteringType="auto"
                placeholder={name}
            />
            {/* {data ? (
                    data[type]?.map(
                        (d, i) =>
                            i < 100 && (
                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                // @ts-ignore
                                <Checkbox
                                    shape="curve"
                                    className="!items-start"
                                    value={d.key}
                                    {...checkbox}
                                >
                                    <Flex
                                        flexDirection="col"
                                        alignItems="start"
                                    >
                                        <Text className="text-gray-800 truncate">
                                            {d.displayName}
                                        </Text>
                                        <Text className="text-xs truncate max-w-[200px]">
                                            {d.key}
                                        </Text>
                                    </Flex>
                                </Checkbox>
                            )
                    )
                ) : (
                    <Spinner />
                )} */}

            {/* // <Flex flexDirection="col" justifyContent="start" alignItems="start"> */}
            {/* <TextInput
                icon={MagnifyingGlassIcon}
                placeholder="Search..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="mb-4"
            /> */}

            {/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
            {/* @ts-ignore */}
            {/* {!compareArrays(value?.sort(), defaultValue?.sort()) && (
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
            )} */}
            {/* // </Flex> */}
        </>
    )
}
