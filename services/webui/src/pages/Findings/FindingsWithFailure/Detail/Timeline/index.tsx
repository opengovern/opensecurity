import { Flex, Icon, Text } from '@tremor/react'
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/solid'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    PlatformEnginePkgComplianceApiFindingEvent,
    PlatformEnginePkgComplianceApiGetFindingEventsByFindingIDResponse,
    PlatformEnginePkgComplianceApiGetSingleResourceFindingResponse,
} from '../../../../../api/api'
import Spinner from '../../../../../components/Spinner'
import { dateTimeDisplay } from '../../../../../utilities/dateDisplay'

dayjs.extend(relativeTime)

interface ITimeline {
    data:
        | PlatformEnginePkgComplianceApiGetSingleResourceFindingResponse
        | PlatformEnginePkgComplianceApiGetFindingEventsByFindingIDResponse
        | undefined
    isLoading: boolean
}

export default function Timeline({ data, isLoading }: ITimeline) {
    const prefix = (
        event: PlatformEnginePkgComplianceApiFindingEvent,
        idx: number
    ) => {
        const str = []
        if (event.previousStateActive !== event.stateActive) {
            if (event.stateActive === true) {
                if (idx === 1) {
                    str.push('Discovered')
                } else {
                    str.push('Re-discovered')
                }
            } else if (event.stateActive === false) {
                str.push('Resource removed')
            } else {
                str.push('Unknown state')
            }
        }

        if (event.previousConformanceStatus !== event.conformanceStatus) {
            if (
                event.conformanceStatus ===
                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed
            ) {
                str.push('Got fixed')
            } else if (
                event.conformanceStatus ===
                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed
            ) {
                str.push('Failed')
            } else {
                str.push('Unknown')
            }
        }
        return `${str.join(' & ')} at`
    }

    return isLoading ? (
        <Spinner />
    ) : (
        <Flex
            flexDirection="col"
            justifyContent="start"
            alignItems="start"
            className="gap-10 relative"
        >
            <div
                className="absolute w-0.5 bg-gray-200 z-10 top-1 left-[13px]"
                style={{ height: 'calc(100% - 30px)' }}
            />
            {data?.findingEvents?.map((tl, idx) => (
                <Flex alignItems="start" className="gap-6 z-20">
                    <Icon
                        icon={
                            tl.conformanceStatus === 'failed'
                                ? XCircleIcon
                                : CheckCircleIcon
                        }
                        color={
                            tl.conformanceStatus === 'failed'
                                ? 'rose'
                                : 'emerald'
                        }
                        size="lg"
                        className="p-0"
                    />
                    <Flex
                        flexDirection="col"
                        alignItems="start"
                        className="w-[500px] min-w-[500px] mt-1 gap-2"
                    >
                        <Text className="text-gray-800">
                            {prefix(
                                tl,
                                (data?.findingEvents?.length || 0) - idx
                            )}{' '}
                            {dateTimeDisplay(tl.evaluatedAt)}
                        </Text>
                        <Text className="text-xs">
                            about {dayjs(tl?.evaluatedAt).fromNow()}
                        </Text>
                    </Flex>
                    {/* <Flex
                        flexDirection="col"
                        alignItems="start"
                        className="gap-1 mt-1"
                    > */}
                    {/* <Text className="text-gray-800 truncate max-w-[330px]"> */}
                    {/*     {tl.controlID} */}
                    {/* </Text> */}
                    {/* <Flex className="w-fit gap-4"> */}
                    {/*     {severityBadge(tl.severity)} */}
                    {/*     <Text className="pl-4 border-l border-l-gray-200 text-xs"> */}
                    {/*         Section: */}
                    {/*     </Text> */}
                    {/* </Flex> */}
                    {/* </Flex> */}
                </Flex>
            ))}
        </Flex>
    )
}
