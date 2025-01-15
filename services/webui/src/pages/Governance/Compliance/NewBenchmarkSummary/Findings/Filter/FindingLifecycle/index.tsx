import { Button, Flex, Text } from '@tremor/react'
import { Radio } from 'pretty-checkbox-react'
import { compareArrays } from '../../../../../../../components/Layout/Header/Filter'

interface ILifecycle {
    value: boolean[]
    defaultValue: boolean[]
    onChange: (l: boolean[]) => void
}

export default function FindingLifecycle({
    value,
    defaultValue,
    onChange,
}: ILifecycle) {
    const options = [
        { name: 'All', value: [true, false] },
        {
            name: 'Active',
            value: [true],
        },
        {
            name: 'Archived',
            value: [false],
        },
    ]

    return (
        <Flex flexDirection="col" alignItems="start" className="gap-1.5">
            {options.map((o) => (
                <Radio
                    name="lifecycle"
                    key={`lifecycle-${o.name}`}
                    checked={compareArrays(value.sort(), o.value.sort())}
                    onClick={() => onChange(o.value)}
                >
                    <Text>{o.name}</Text>
                </Radio>
            ))}
            {/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
            {/* @ts-ignore */}
            {!compareArrays(value?.sort(), defaultValue?.sort()) && (
                <Flex className="pt-3 mt-3 border-t border-t-gray-200">
                    <Button
                        variant="light"
                        onClick={() => onChange(defaultValue)}
                    >
                        Reset
                    </Button>
                </Flex>
            )}
        </Flex>
    )
}
