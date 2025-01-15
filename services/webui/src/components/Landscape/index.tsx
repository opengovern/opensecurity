import { Card, Col, Divider, Flex, Grid, Text, Title } from '@tremor/react'
import { PlatformEnginePkgInventoryApiResourceCollectionLandscape } from '../../api/api'
import Spinner from '../Spinner'

interface ILandscape {
    data:
        | PlatformEnginePkgInventoryApiResourceCollectionLandscape
        | undefined
    isLoading: boolean
}

const parentCardWidth = (num: number) => {
    if (num <= 4) {
        return 1
    }
    if (num <= 8) {
        return 2
    }
    if (num <= 12) {
        return 3
    }
    return 4
}

const parentCardSize = (num: number) => {
    if (num <= 4) {
        return 3
    }
    if (num <= 8) {
        return 5
    }
    if (num <= 12) {
        return 7
    }
    return 9
}

export default function Landscape({ data, isLoading }: ILandscape) {
    return (
        <Card>
            {isLoading ? (
                <Spinner className="my-24" />
            ) : (
                <Grid numItems={1} className="w-full gap-6">
                    {data?.categories?.map((cat) => (
                        <Flex flexDirection="col" className="gap-3">
                            <Title className="font-semibold">{cat.name}</Title>
                            <Card className="w-full bg-openg-100 p-4">
                                <Grid
                                    numItems={2}
                                    numItemsLg={4}
                                    className="w-full gap-3"
                                >
                                    {cat.subcategories
                                        ?.sort(
                                            (a, b) =>
                                                (b.items?.length || 0) -
                                                (a.items?.length || 0)
                                        )
                                        .map((sub) => (
                                            <Col
                                                numColSpan={parentCardWidth(
                                                    sub.items?.length || 0
                                                )}
                                            >
                                                <Card className="p-3 h-full">
                                                    <Flex flexDirection="col">
                                                        <Text className="font-semibold text-openg-600">
                                                            {sub.name}
                                                        </Text>
                                                        <Divider className="my-2" />
                                                        <Grid
                                                            numItems={Math.floor(
                                                                parentCardSize(
                                                                    sub.items
                                                                        ?.length ||
                                                                        0
                                                                ) / 2
                                                            )}
                                                            numItemsLg={parentCardSize(
                                                                sub.items
                                                                    ?.length ||
                                                                    0
                                                            )}
                                                            className="w-full gap-3"
                                                        >
                                                            {sub.items?.map(
                                                                (item) => (
                                                                    <Flex justifyContent="center">
                                                                        <button
                                                                            type="button"
                                                                            title={
                                                                                item.name
                                                                            }
                                                                            className="p-1 border border-white rounded-md hover:border-openg-200 transition-all"
                                                                        >
                                                                            <img
                                                                                className="h-10 w-10 rounded-md"
                                                                                src={
                                                                                    item.logo_uri
                                                                                }
                                                                                alt={
                                                                                    item.name
                                                                                }
                                                                            />
                                                                        </button>
                                                                    </Flex>
                                                                )
                                                            )}
                                                        </Grid>
                                                    </Flex>
                                                </Card>
                                            </Col>
                                        ))}
                                </Grid>
                            </Card>
                        </Flex>
                    ))}
                </Grid>
            )}
        </Card>
    )
}
