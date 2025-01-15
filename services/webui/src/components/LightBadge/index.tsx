import { Flex, Icon, Text } from '@tremor/react'
import { ArrowDownFill, ArrowRightFill, ArrowUpFill } from '../../icons/icons'

interface IProbs {
    value: number
}

export default function LightBadge({ value }: IProbs) {
    let badgeProperty = { icon: <ArrowRightFill />, color: '' }
    if (value < 0)
        badgeProperty = { icon: <ArrowDownFill />, color: 'text-rose-500' }
    else if (value > 0)
        badgeProperty = { icon: <ArrowUpFill />, color: 'text-emerald-500' }

    return (
        <Flex className="gap-0.5 p-0.5 w-fit">
            {badgeProperty.icon}
            <Text className={`!text-sm ${badgeProperty.color}`}>
                {Math.abs(value)}
                <span className="text-xs ml-0.5">%</span>
            </Text>
        </Flex>
    )
}
