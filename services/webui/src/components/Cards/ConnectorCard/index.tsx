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
