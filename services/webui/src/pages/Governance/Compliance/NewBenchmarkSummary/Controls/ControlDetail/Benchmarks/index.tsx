import { useAtomValue } from 'jotai'
import { Badge, Card, Flex, Text, Title } from '@tremor/react'
import { ChevronRightIcon } from '@heroicons/react/24/solid'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import { PlatformEnginePkgComplianceApiBenchmark } from '../../../../../../../api/api'
import { getConnectorIcon } from '../../../../../../../components/Cards/ConnectorCard'
import { searchAtom } from '../../../../../../../utilities/urlstate'

interface IBenchmarks {
    benchmarks:
        | PlatformEnginePkgComplianceApiBenchmark[]
        | undefined
}

export default function Benchmarks({ benchmarks }: IBenchmarks) {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const workspace = useParams<{ ws: string }>().ws
                                   


    return (
        <Flex flexDirection="col" className="gap-3">
            {benchmarks?.map((bm) => (
                <Card
                    className="w-full py-4 cursor-pointer"
                    onClick={() => navigate(`/compliance/${bm.id}`)}
                >
                    <Flex>
                        <Flex
                            flexDirection="col"
                            alignItems="start"
                            className="w-3/4"
                        >
                            <Flex justifyContent="start" className="gap-3">
                                {getConnectorIcon(bm.integrationTypes)}
                                <Title className="font-semibold">
                                    {bm.title}
                                </Title>
                            </Flex>
                            <Text className="ml-12">{bm.description}</Text>
                        </Flex>
                        <Flex
                            justifyContent="end"
                            className="gap-2 flex-wrap w-1/4"
                        >
                            {bm?.tags?.category?.map((cat) => (
                                <Badge color="slate" size="xs">
                                    {cat}
                                </Badge>
                            ))}
                            {!!bm?.tags?.cis && (
                                <Badge color="sky" size="xs">
                                    CIS
                                </Badge>
                            )}
                            {!!bm?.tags?.hipaa && (
                                <Badge color="blue" size="xs">
                                    Hipaa
                                </Badge>
                            )}
                        </Flex>
                        <ChevronRightIcon className="h-5 ml-3 text-openg-500" />
                    </Flex>
                </Card>
            ))}
        </Flex>
    )
}
