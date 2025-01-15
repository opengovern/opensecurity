import {
    Badge,
    Button,
    Card,
    Flex,
    Icon,
    Subtitle,
    Text,
    Title,
} from '@tremor/react'
import { ChevronRightIcon, LinkIcon } from '@heroicons/react/24/outline'
import { useNavigate } from 'react-router-dom'
import { useAtomValue } from 'jotai'
import { numericDisplay } from '../../../utilities/numericDisplay'
import { AWSAzureIcon, AWSIcon, AzureIcon } from '../../../icons/icons'
import {
    PlatformEngineServicesIntegrationApiEntityTier,
    SourceType,
} from '../../../api/api'
import { searchAtom } from '../../../utilities/urlstate'
import './style.css'
interface IConnectorCard {
    connector: string | undefined
    title: string | undefined
    status: boolean | undefined
    count: number | undefined
    description: string | undefined
    tier?: PlatformEngineServicesIntegrationApiEntityTier
    logo?: string
    onClickCard?: Function
    name?: string
    id?: number
}
export const getConnectorsIcon = (connector: SourceType[], className = '') => {
    if (connector?.length >= 2) {
        return (
            <img
                src={AWSAzureIcon}
                alt="connector"
                className="min-w-[36px] w-9 h-9 rounded-full"
            />
        )
    }

    const connectorIcon = () => {
        if (connector[0] === SourceType.CloudAzure) {
            return AzureIcon
        }
        if (connector[0] === SourceType.CloudAWS) {
            return AWSIcon
        }
        return undefined
    }

    return (
        <Flex className={`w-9 h-9 gap-1 ${className}`}>
            <img
                src={connectorIcon()}
                alt="connector"
                className="min-w-[36px] w-9 h-9 rounded-full"
            />
        </Flex>
    )
}

export const getConnectorIcon = (
    connector: string | SourceType[] | SourceType | undefined | string[],
    className = ''
) => {
    const connectorIcon = () => {
        if (String(connector).toLowerCase() === 'azure_subscription') {
            return AzureIcon
        }
        if (String(connector).toLowerCase() === 'aws_cloud_account') {
            return AWSIcon
        }
        if (connector?.length && connector?.length > 0) {
            if (String(connector[0]).toLowerCase() === 'azure_subscription') {
                return AzureIcon
            }
            if (String(connector[0]).toLowerCase() === 'aws_cloud_account') {
                return AWSIcon
            }
        }
        return undefined
    }

    return (
        <Flex className={`w-9 h-9 gap-1 ${className}`}>
            <img
                src={connectorIcon()}
                alt="connector"
                className="min-w-[36px] w-9 h-9 rounded-full"
            />
        </Flex>
    )
}

export default function ConnectorCard({
    connector,
    title,
    status,
    count,
    description,
    name,
    tier,
    logo,
    onClickCard,
    id,
}: IConnectorCard) {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)

    const onClick = () => {
        onClickCard?.()
    }

    return (
        <>
            <Card
                key={connector}
                className={`cursor-pointer integration-card  w-[210px] h-[140px] ${
                    tier ==
                        PlatformEngineServicesIntegrationApiEntityTier.TierEnterprise &&
                    'enterprise'
                } `}
                onClick={() => onClick()}
            >
                <Flex
                    flexDirection="col"
                    justifyContent="center"
                    alignItems="center"
                    className="gap-[16px] h-full"
                >
                    {logo === undefined || logo === '' ? (
                        <LinkIcon className="w-[50px] h-[35px]  " />
                    ) : (
                        <Flex className="w-[50px] h-[35px] ">
                            <img
                                src={logo}
                                alt="Connector Logo"
                                className="w-[50px] h-[35px] connector-logo "
                            />
                        </Flex>
                    )}
                    <Title className="integration-name text-center w-full">
                        {title}
                    </Title>
                    {/* <Flex
                        flexDirection="col"
                        justifyContent="center"
                        alignItems="center"
                        // className="gap-[8px] "
                    >
                        {' '}
                        {logo === undefined || logo === '' ? (
                            <LinkIcon className="w-[50px] h-[35px]  " />
                        ) : (
                            <Flex className="w-[50px] h-[35px] ">
                                <img
                                    src={logo}
                                    alt="Connector Logo"
                                    className="w-[50px] h-[35px] "
                                />
                            </Flex>
                        )}
                        <Title className="integration-name text-center w-full">{title}</Title>
                    </Flex> */}

                    {/* {tier ==
                    GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityTier.TierEnterprise ? (
                        <>
                            <div className="w-fit rounded-[32px] px-[12px] py-[2px] border-[#F90] border-[1px] bg-[#FFEED4]">
                                <Flex
                                    flexDirection="row"
                                    justifyContent="center"
                                    alignItems="center"
                                    className="gap-1 "
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        width="17"
                                        height="17"
                                        viewBox="0 0 17 17"
                                        fill="none"
                                    >
                                        <path
                                            d="M14.5026 8.5C14.5026 11.8151 11.8151 14.5025 8.50006 14.5025C5.18497 14.5025 2.49756 11.8151 2.49756 8.5C2.49756 5.18491 5.18497 2.4975 8.50006 2.4975"
                                            stroke="#FF9900"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                        <path
                                            fill-rule="evenodd"
                                            clip-rule="evenodd"
                                            d="M12.7287 2.09203C12.8166 1.93103 12.9854 1.83087 13.1688 1.83087C13.3523 1.83087 13.5211 1.93103 13.609 2.09203L13.9979 2.80233C14.0438 2.88696 14.1133 2.95646 14.1979 3.00241L14.9076 3.39057C15.0689 3.47859 15.1692 3.64767 15.1692 3.83142C15.1692 4.01518 15.0689 4.18426 14.9076 4.27227L14.1979 4.6611C14.1132 4.70687 14.0436 4.77642 13.9979 4.86119L13.609 5.57081C13.5211 5.73182 13.3523 5.83198 13.1688 5.83198C12.9854 5.83198 12.8166 5.73182 12.7287 5.57081L12.3398 4.86052C12.2941 4.77575 12.2245 4.7062 12.1397 4.66044L11.4301 4.27227C11.2688 4.18426 11.1685 4.01518 11.1685 3.83142C11.1685 3.64767 11.2688 3.47859 11.4301 3.39057L12.1397 3.00241C12.2244 2.95646 12.2939 2.88696 12.3398 2.80233L12.7287 2.09203Z"
                                            stroke="#FF9900"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                        <path
                                            d="M9.3335 7.16608H11.0009C11.185 7.16608 11.3343 7.31538 11.3343 7.49955V10.8343C11.3343 11.0184 11.185 11.1677 11.0009 11.1677H9.3335"
                                            stroke="#FF9900"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                        <path
                                            d="M7.33286 11.1678H5.6655C5.48133 11.1678 5.33203 11.0185 5.33203 10.8343V8.83347C5.33203 8.6493 5.48133 8.5 5.6655 8.5H7.33286"
                                            stroke="#FF9900"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                        <path
                                            fill-rule="evenodd"
                                            clip-rule="evenodd"
                                            d="M9.33384 11.1678H7.33301V5.49876C7.33301 5.31458 7.48231 5.16528 7.66648 5.16528H9.00037C9.18454 5.16528 9.33384 5.31458 9.33384 5.49876V11.1678Z"
                                            stroke="#FF9900"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                    </svg>
                                    <Text className="text-[#FF9900]">
                                        Enterprise
                                    </Text>
                                </Flex>
                            </div>
                        </>
                    ) : (
                        <>
                            {count !== 0 && (
                                <>
                                    {' '}
                                    <div className="w-fit rounded-[32px] px-[12px] py-[2px] border-[#164085] border-[1px] bg-[#DBE9FF]   ">
                                        <Flex
                                            flexDirection="row"
                                            justifyContent="center"
                                            alignItems="center"
                                            className="gap-1 "
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                width="13"
                                                height="13"
                                                viewBox="0 0 13 13"
                                                fill="none"
                                            >
                                                <path
                                                    d="M7.04 10.44L5.93335 11.5667C5.64002 11.86 5.28665 12.1 4.89998 12.26C4.51332 12.42 4.10002 12.5 3.68669 12.5C3.27335 12.5 2.85334 12.42 2.46667 12.26C2.08001 12.1 1.72668 11.8667 1.43335 11.5667C1.14002 11.2733 0.89999 10.92 0.73999 10.5333C0.57999 10.1467 0.5 9.73331 0.5 9.31331C0.5 8.89331 0.57999 8.48001 0.73999 8.09334C0.89999 7.70668 1.14002 7.35335 1.43335 7.06002L2.55334 5.93335"
                                                    stroke="#164085"
                                                    stroke-miterlimit="10"
                                                    stroke-linecap="round"
                                                />
                                                <path
                                                    d="M5.93311 2.56002L7.05977 1.43335C7.65977 0.840016 8.4664 0.5 9.31307 0.5C10.1597 0.5 10.9664 0.83335 11.5664 1.43335C12.1597 2.03335 12.4998 2.84002 12.4998 3.68669C12.4998 4.53335 12.1597 5.33998 11.5664 5.93998L10.4397 7.06665"
                                                    stroke="#164085"
                                                    stroke-miterlimit="10"
                                                    stroke-linecap="round"
                                                />
                                                <path
                                                    d="M4.22021 8.75336L8.72689 4.25336"
                                                    stroke="#164085"
                                                    stroke-miterlimit="10"
                                                    stroke-linecap="round"
                                                />
                                            </svg>
                                            <Text className="text-[#164085]">
                                                {count}
                                            </Text>
                                        </Flex>
                                    </div>
                                </>
                            )}
                        </>
                    )} */}
                </Flex>
                {count !== 0 && (
                    <button className="integration-button " onClick={onClick}>
                        <>
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="13"
                                height="13"
                                viewBox="0 0 13 13"
                                fill="none"
                            >
                                <path
                                    d="M7.04 10.44L5.93335 11.5667C5.64002 11.86 5.28665 12.1 4.89998 12.26C4.51332 12.42 4.10002 12.5 3.68669 12.5C3.27335 12.5 2.85334 12.42 2.46667 12.26C2.08001 12.1 1.72668 11.8667 1.43335 11.5667C1.14002 11.2733 0.89999 10.92 0.73999 10.5333C0.57999 10.1467 0.5 9.73331 0.5 9.31331C0.5 8.89331 0.57999 8.48001 0.73999 8.09334C0.89999 7.70668 1.14002 7.35335 1.43335 7.06002L2.55334 5.93335"
                                    stroke="#164085"
                                    stroke-miterlimit="10"
                                    stroke-linecap="round"
                                />
                                <path
                                    d="M5.93311 2.56002L7.05977 1.43335C7.65977 0.840016 8.4664 0.5 9.31307 0.5C10.1597 0.5 10.9664 0.83335 11.5664 1.43335C12.1597 2.03335 12.4998 2.84002 12.4998 3.68669C12.4998 4.53335 12.1597 5.33998 11.5664 5.93998L10.4397 7.06665"
                                    stroke="#164085"
                                    stroke-miterlimit="10"
                                    stroke-linecap="round"
                                />
                                <path
                                    d="M4.22021 8.75336L8.72689 4.25336"
                                    stroke="#164085"
                                    stroke-miterlimit="10"
                                    stroke-linecap="round"
                                />
                            </svg>
                            {count}
                        </>
                    </button>
                )}
                {tier ==
                    PlatformEngineServicesIntegrationApiEntityTier.TierEnterprise && (
                    <>
                        <p className="coming-soon"> ENTERPRISE </p>
                        <p className="add-btn ">+</p>
                    </>
                )}
            </Card>
        </>
    )
}

// <Card
//     key={connector}
//     className="cursor-pointer"
//     onClick={() => onClick()}
// >
//     <Flex flexDirection="row" className="mb-3">
//         {logo === undefined || logo === '' ? (
//             <LinkIcon className="w-9 h-9 gap-1" />
//         ) : (
//             <Flex className="w-9 h-9 gap-1">
//                 <img
//                     src={logo}
//                     alt="Connector Logo"
//                     className="min-w-[36px] w-9 h-9 rounded-full"
//                 />
//             </Flex>
//         )}
//         <Badge color={getTierBadgeColor(tier)}>
//             {tier ===
//             GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityTier.TierCommunity ? (
//                 <Text color="emerald">Community</Text>
//             ) : (
//                 <Text color="violet">Enterprise</Text>
//             )}
//         </Badge>
//         {/* <Badge color={getTierBadgeColor(tier)}>
//             {tier ===
//             GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityTier.TierCommunity ? (
//                 <Text color="emerald">Community</Text>
//             ) : (
//                 <Text color="emerald">Enterprise</Text>
//             )}
//         </Badge>
//         <Badge color={getBadgeColor(status)}>
//             {status === 'enabled' ? (
//                 <Text color="emerald">Active</Text>
//             ) : (
//                 <Text color="rose">InActive</Text>
//             )}
//         </Badge> */}
//     </Flex>
//     <Flex flexDirection="row" className="mb-1">
//         <Title className="font-semibold">{title}</Title>
//         {(count || 0) !== 0 && (
//             <Title className="font-semibold">
//                 {numericDisplay(count)}
//             </Title>
//         )}
//     </Flex>
//     <Subtitle>{description}</Subtitle>
//     <Flex flexDirection="row" justifyContent="end">
//         {button()}
//     </Flex>
// </Card>
