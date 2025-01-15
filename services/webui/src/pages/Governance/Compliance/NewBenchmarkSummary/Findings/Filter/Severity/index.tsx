import { Checkbox, useCheckboxState } from 'pretty-checkbox-react'
import { Button, Flex, Text } from '@tremor/react'
import { useEffect, useState } from 'react'
import { TypesFindingSeverity } from '../../../../../../../api/api'
import { compareArrays } from '../../../../../../../components/Layout/Header/Filter'
import Multiselect from '@cloudscape-design/components/multiselect'

interface ISeverity {
    value: TypesFindingSeverity[] | undefined
    defaultValue: TypesFindingSeverity[]
    condition: string
    onChange: (s: TypesFindingSeverity[]) => void
}

export default function Severity({
    value,
    defaultValue,
    condition,
    onChange,
}: ISeverity) {
    const [selectedOptions, setSelectedOptions] = useState([
        {
            label: 'Critical',
            value: TypesFindingSeverity.FindingSeverityCritical,
            color: '#6E120B',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#6E120B',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'High',
            value: TypesFindingSeverity.FindingSeverityHigh,
            color: '#CA2B1D',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#ca2b1d',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'Medium',
            value: TypesFindingSeverity.FindingSeverityMedium,
            color: '#EE9235',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#ee9235',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'Low',
            value: TypesFindingSeverity.FindingSeverityLow,
            color: '#F4C744',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#f4c744',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'None',
            value: TypesFindingSeverity.FindingSeverityNone,
            color: '#9BA2AE',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#9ba2ae',
                        }}
                    />
                </>
            ),
        },
    ])

    const options = [
        {
            label: 'Critical',
            value: TypesFindingSeverity.FindingSeverityCritical,
            color: '#6E120B',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#6E120B',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'High',
            value: TypesFindingSeverity.FindingSeverityHigh,
            color: '#CA2B1D',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#ca2b1d',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'Medium',
            value: TypesFindingSeverity.FindingSeverityMedium,
            color: '#EE9235',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#ee9235',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'Low',
            value: TypesFindingSeverity.FindingSeverityLow,
            color: '#F4C744',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#f4c744',
                        }}
                    />
                </>
            ),
        },
        {
            label: 'None',
            value: TypesFindingSeverity.FindingSeverityNone,
            color: '#9BA2AE',
            iconSvg: (
                <>
                    <div
                        className="h-4 w-1.5 rounded-sm"
                        style={{
                            backgroundColor: '#9ba2ae',
                        }}
                    />
                </>
            ),
        },
    ]

    useEffect(() => {
        if (selectedOptions.length === 0) {
            onChange(defaultValue)
            return
        }
        else {
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
                tokenLimit={0}
                onChange={({ detail }) =>
                    // @ts-ignore
                    setSelectedOptions(detail.selectedOptions)
                }
                options={options}
                // filteringType="auto"
                placeholder=" Severity"
                virtualScroll
            />
            {/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
            {/* @ts-ignore */}
        </>
    )
}

{
    /**
    
                    <Checkbox
                    shape="curve"
                    className="!items-start w-full"
                    value={o.value}
                    {...severityCheckbox}
                >
                    <Flex className="gap-1.5">
                        <div
                            className="h-4 w-1.5 rounded-sm"
                            style={{
                                backgroundColor: o.color,
                            }}
                        />
                        <Text>{o.name}</Text>
                    </Flex>
                </Checkbox>
            ))}

    */
}
