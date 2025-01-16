import { ChevronRightIcon } from '@heroicons/react/24/outline'
import { ChevronRightIcon as ChevronRightIconSolid } from '@heroicons/react/20/solid'
import { useAtomValue } from 'jotai'

import {
    Flex,
    ProgressCircle,
    Text,
    Title,
    Icon,
    Subtitle,
    Button,
    Card,
} from '@tremor/react'
import { useNavigate, useParams } from 'react-router-dom'
import { numericDisplay } from '../../../utilities/numericDisplay'
import { searchAtom } from '../../../utilities/urlstate'

interface IScoreCategoryCard {
    title: string
    percentage: number
    value: number
    kpiText: string
    costOptimization: number
    varient?: 'minimized' | 'default'
    category: string
}

export default function ScoreCategoryCard({
    title,
    percentage,
    value,
    kpiText,
    costOptimization,
    varient,
    category,
}: IScoreCategoryCard) {
    const { ws } = useParams()
    const navigate = useNavigate()
    // const { response, isLoading } =
    //     useComplianceApiV1BenchmarksControlsDetail(controlID)
    const searchParams = useAtomValue(searchAtom)

    let color = 'blue'
    if (percentage >= 75) {
        color = 'emerald'
    } else if (percentage >= 50 && percentage < 75) {
        color = 'lime'
    } else if (percentage >= 25 && percentage < 50) {
        color = 'yellow'
    } else if (percentage >= 0 && percentage < 25) {
        color = 'red'
    }
    return (
        <Card
            onClick={() =>
                navigate(`/compliance/${category}${searchParams}`)
            }
            className={` ${
                varient === 'default'
                    ? 'gap-6 2xl:px-8 sm:px-4 py-8 rounded-xl'
                    : 'pl-5 pr-4 py-6 rounded-lg'
            } ${
                varient === 'default' ? 'items-center' : 'items-start'
            } flex bg-white dark:bg-openg-950 shadow-sm  hover:shadow-lg hover:cursor-pointer`}
        >
            <Flex className="relative w-fit">
                <ProgressCircle color={color} value={percentage} size="md" >
                    <Text>{percentage? percentage.toFixed(1) : 0}%</Text>
                </ProgressCircle>
            </Flex>
            <Flex justifyContent="between" className="h-full">
                <Flex
                    alignItems="start"
                    flexDirection="col"
                    className={varient === 'default' ? 'gap-2' : 'pl-5 gap-1.5'}
                >
                    <Title
                        className={
                            varient === 'default'
                                ? 'text-xl'
                                : '!text-base font-bold'
                        }
                    >
                        {title}
                    </Title>

                    {costOptimization > 0 || title == 'Efficiency' ? (
                        // <Text>${costOptimization} Waste</Text>
                        <Text>
                            <Flex className="gap-1">
                                <span className="text-gray-900">{value}</span>
                                <span>{kpiText}</span>
                            </Flex>
                        </Text>
                    ) : (
                        <Text>
                            <Flex className="gap-1">
                                <span className="text-gray-900">{value}</span>
                                <span>{kpiText}</span>
                            </Flex>
                        </Text>
                    )}
                    {/* <BadgeDeltaSimple change={change}>
                    from previous time period
                </BadgeDeltaSimple> */}
                </Flex>
                {varient === 'default' ? (
                    <Icon size="md" icon={ChevronRightIcon} />
                ) : (
                    <ChevronRightIconSolid className="w-6 text-gray-300" />
                )}
            </Flex>
        </Card>
    )
}
