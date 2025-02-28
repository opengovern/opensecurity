import { Button, Flex, Text } from '@tremor/react'
import { Radio } from 'pretty-checkbox-react'
import { SourceType } from '../../../../../../../api/api'
import { AWSIcon, AzureIcon } from '../../../../../../../icons/icons'

interface IProvider {
    value: SourceType
    defaultValue: SourceType
    onChange: (p: SourceType) => void
}
export default function Provider({ value, defaultValue, onChange }: IProvider) {
    const options = [
        { name: 'All', value: SourceType.Nil, icon: undefined },
        {
            name: 'AWS',
            value: SourceType.CloudAWS,
            icon: <img src={AWSIcon} className="w-5 rounded-full" alt="aws" />,
        },
        {
            name: 'Azure',
            value: SourceType.CloudAzure,
            icon: (
                <img src={AzureIcon} className="w-5 rounded-full" alt="azure" />
            ),
        },
    ]

    return (
        <Flex flexDirection="col" alignItems="start" className="gap-1.5">
            {options.map((o) => (
                <Radio
                    name="provider"
                    key={`provider-${o.name}`}
                    checked={value === o.value}
                    onClick={() => onChange(o.value)}
                >
                    <Flex className="gap-1 w-fit">
                        {o.icon}
                        <Text>{o.name}</Text>
                    </Flex>
                </Radio>
            ))}
            {value !== defaultValue && (
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
