import { CheckCircleIcon } from '@heroicons/react/24/solid'
import { ArrowTopRightOnSquareIcon } from '@heroicons/react/24/outline'
import { Flex, Grid, Text, Title } from '@tremor/react'
import { OpenGovernanceBig } from '../../icons/icons'

export default function RequestDemo() {
    return (
        <Grid numItems={2} className="min-h-[750px]">
            <Flex
                flexDirection="col"
                justifyContent="between"
                alignItems="start"
                className="pt-36 pl-40"
            >
                <Flex
                    flexDirection="col"
                    justifyContent="start"
                    alignItems="start"
                >
                    <OpenGovernanceBig className="w-32 mb-9" />
                    <Title className="text-gray-800 text-2xl mb-6">
                        Manage Complexity
                    </Title>
                    <Text className="text-gray-600 text-sm mb-6">
                        <span className="font-bold text-gray-800">
                            Book an appointment
                        </span>{' '}
                        to access Navigate Infrastructure, Compliance, and
                        Costs, Across All Clouds with One Intuitive Platform.
                    </Text>
                    <Flex
                        flexDirection="row"
                        justifyContent="start"
                        className="mb-3 ml-4"
                    >
                        <CheckCircleIcon className="w-5 text-emerald-500 mr-1" />
                        <Text>Navigate Infrastructure</Text>
                    </Flex>
                    <Flex
                        flexDirection="row"
                        justifyContent="start"
                        className="mb-3 ml-4"
                    >
                        <CheckCircleIcon className="w-5 text-emerald-500 mr-1" />
                        <Text>Navigate Spend and Spend Metrics</Text>
                    </Flex>
                    <Flex
                        flexDirection="row"
                        justifyContent="start"
                        className="mb-3 ml-4"
                    >
                        <CheckCircleIcon className="w-5 text-emerald-500 mr-1" />
                        <Text>Manage your cloud accounts</Text>
                    </Flex>
                    <Flex
                        flexDirection="row"
                        justifyContent="start"
                        className="mb-3 ml-4"
                    >
                        <CheckCircleIcon className="w-5 text-emerald-500 mr-1" />
                        <Text>Navigate Insights</Text>
                    </Flex>
                    <Flex flexDirection="row" className=" my-8">
                        <div className="w-full mx-2 border-t border-gray-200" />
                        <Text>Or</Text>
                        <div className="w-full mx-2 border-t border-gray-200" />
                    </Flex>
                    <Text className="text-gray-600 text-sm mb-4">
                        <span className="font-bold text-gray-800">
                            Book a Demo
                        </span>{' '}
                        to navigate opensecurity product
                    </Text>
                    <Flex
                        flexDirection="row"
                        justifyContent="start"
                        className="text-sm text-openg-500"
                    >
                        <ArrowTopRightOnSquareIcon className="w-4 mr-1" />{' '}
                        Booking a Demo
                    </Flex>
                </Flex>
                <Text>
                    © 2025 open governance Inc. d/b/a opensecurity.io. All rights
                    reserved.
                </Text>
            </Flex>
            <iframe
                title="Setup a meeting"
                className="w-full h-full relative max-w-[630px]"
                src="https://cal.com/team/opensecurity/try-enterprise"
            />
        </Grid>
    )
}
