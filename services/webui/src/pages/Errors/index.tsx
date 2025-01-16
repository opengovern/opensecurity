import { Flex, Grid, Metric, Text } from '@tremor/react'
import notFoundImg from '../../icons/404.png'
import Layout from '../../components/Layout'
import TopHeader from '../../components/Layout/Header'

export default function NotFound() {
    return (
        <>
            <TopHeader />
            <Flex justifyContent="center" className="md:mt-32">
                <Grid numItems={1} numItemsMd={2}>
                    <img src={notFoundImg} alt="stack" />
                    <Flex
                        flexDirection="col"
                        alignItems="start"
                        justifyContent="center"
                        className="p-12"
                    >
                        <Metric>Remain calm</Metric>
                        <Text className="mt-3 mb-6">Something is missing!</Text>
                    </Flex>
                </Grid>
            </Flex>
        </>
    )
}
