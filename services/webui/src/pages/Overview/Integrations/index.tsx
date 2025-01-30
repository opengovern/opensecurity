// @ts-nocheck

import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Button,
    Card,
    Flex,
    Icon,
    Text,
    Title,
} from '@tremor/react'
import {
    ChevronRightIcon,
    MagnifyingGlassIcon,
} from '@heroicons/react/24/outline'
import Editor from 'react-simple-code-editor'
import 'prismjs/themes/prism.css'
import { highlight, languages } from 'prismjs'
import { useNavigate, useParams } from 'react-router-dom'
import { useAtom } from 'jotai'
import { useEffect, useState } from 'react'
import { useInventoryApiV1QueryList } from '../../../api/inventory.gen'
import { runQueryAtom } from '../../../store'
import { getErrorMessage } from '../../../types/apierror'
import { PlatformEnginePkgInventoryApiSmartQueryItem, PlatformEngineServicesIntegrationApiEntityTier } from '../../../api/api'
import { useIntegrationApiV1ConnectorsList } from '../../../api/integration.gen'
import { Box, Cards, Link, Modal, SpaceBetween } from '@cloudscape-design/components'
import axios from 'axios'



const GetTierIcon = (tier: string) => {
    if (tier === 'Community') {
        return (
            <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M11.9034 2.25264C13.1418 2.20921 14.3466 2.70205 15.342 3.69744L20.3003 8.65575C21.2978 9.65321 21.7914 10.8583 21.7473 12.0969C21.7035 13.3273 21.1328 14.5083 20.1572 15.4834L15.4851 20.1555C14.5101 21.1305 13.3293 21.7013 12.0991 21.745C10.8605 21.7891 9.65543 21.2955 8.65754 20.2976L3.69923 15.3393C2.70179 14.3418 2.20839 13.1372 2.25275 11.8988C2.29681 10.6687 2.86785 9.48812 3.84318 8.51277C3.84316 8.51279 3.84319 8.51276 3.84318 8.51277L8.51445 3.84052C9.49052 2.86501 10.6731 2.29579 11.9034 2.25264ZM9.57502 4.90128C9.57495 4.90134 9.57509 4.90121 9.57502 4.90128L4.90395 9.57332C4.1474 10.3299 3.77996 11.1658 3.75179 11.9525C3.7239 12.7308 4.02445 13.5432 4.75989 14.2786L9.7182 19.2369C10.4542 19.9729 11.267 20.2737 12.0457 20.246C12.8327 20.218 13.6685 19.8508 14.4244 19.0949L19.0966 14.4226C19.8529 13.6668 20.2203 12.8307 20.2483 12.0436C20.276 11.2647 19.9751 10.4518 19.2397 9.71641L14.2814 4.7581C13.5487 4.02548 12.7361 3.72436 11.956 3.75172C11.168 3.77935 10.3308 4.146 9.57502 4.90128Z"
                    fill="black"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M11.3716 14.3594C10.9772 14.3594 10.6582 14.6784 10.6582 15.0727C10.6582 15.4671 10.9772 15.7861 11.3716 15.7861C11.7659 15.7861 12.0849 15.4671 12.0849 15.0727C12.0849 14.6784 11.7659 14.3594 11.3716 14.3594ZM9.1582 15.0727C9.1582 13.85 10.1488 12.8594 11.3716 12.8594C12.5943 12.8594 13.5849 13.85 13.5849 15.0727C13.5849 16.2955 12.5943 17.2861 11.3716 17.2861C10.1488 17.2861 9.1582 16.2955 9.1582 15.0727Z"
                    fill="black"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M15.4751 11.2891C15.0807 11.2891 14.7617 11.6081 14.7617 12.0024C14.7617 12.3968 15.0807 12.7158 15.4751 12.7158C15.8694 12.7158 16.1884 12.3968 16.1884 12.0024C16.1884 11.6081 15.8694 11.2891 15.4751 11.2891ZM13.2617 12.0024C13.2617 10.7797 14.2523 9.78906 15.4751 9.78906C16.6978 9.78906 17.6884 10.7797 17.6884 12.0024C17.6884 13.2252 16.6978 14.2158 15.4751 14.2158C14.2523 14.2158 13.2617 13.2252 13.2617 12.0024Z"
                    fill="black"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M11.3716 7.37109C10.9772 7.37109 10.6582 7.69012 10.6582 8.08446C10.6582 8.47879 10.9772 8.79782 11.3716 8.79782C11.7659 8.79782 12.0849 8.47879 12.0849 8.08446C12.0849 7.69012 11.7659 7.37109 11.3716 7.37109ZM9.1582 8.08446C9.1582 6.8617 10.1488 5.87109 11.3716 5.87109C12.5943 5.87109 13.5849 6.8617 13.5849 8.08446C13.5849 9.30722 12.5943 10.2978 11.3716 10.2978C10.1488 10.2978 9.1582 9.30722 9.1582 8.08446Z"
                    fill="black"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M7.98456 4.53877C8.27813 4.24655 8.753 4.24764 9.04522 4.5412L10.9464 6.45116C11.2386 6.74473 11.2376 7.2196 10.944 7.51182C10.6504 7.80404 10.1755 7.80295 9.88333 7.50938L7.98212 5.59942C7.68991 5.30586 7.691 4.83098 7.98456 4.53877Z"
                    fill="black"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M11.9501 8.51264C12.243 8.21975 12.7179 8.21975 13.0108 8.51264L14.9558 10.4576C15.2487 10.7505 15.2487 11.2254 14.9558 11.5183C14.6629 11.8112 14.188 11.8112 13.8951 11.5183L11.9501 9.5733C11.6572 9.28041 11.6572 8.80553 11.9501 8.51264Z"
                    fill="black"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M11.373 8.79688C11.7873 8.79688 12.123 9.13266 12.123 9.54688V13.6096C12.123 14.0239 11.7873 14.3596 11.373 14.3596C10.9588 14.3596 10.623 14.0239 10.623 13.6096V9.54688C10.623 9.13266 10.9588 8.79688 11.373 8.79688Z"
                    fill="black"
                />
            </svg>
        )
    } else if (tier === 'Premium') {
        return (
            <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    d="M3.51047 11.5274L10.8163 19.5542C11.4507 20.2518 12.5482 20.2518 13.1825 19.5542L20.4894 11.5264C21.1004 10.8561 21.1704 9.85395 20.6596 9.10381L17.7174 4.78197C17.3526 4.24685 16.7455 3.92578 16.0975 3.92578H7.91111C7.26312 3.92578 6.65698 4.24588 6.29212 4.781L3.34117 9.10381C2.8294 9.85298 2.89946 10.8561 3.51047 11.5274Z"
                    stroke="black"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
                <path
                    d="M13.9989 3.92578L15.3659 9.10381C15.5643 9.85395 15.5371 10.8561 15.2997 11.5264L12.4597 19.5542C12.2135 20.2518 11.7864 20.2518 11.5402 19.5542L8.7002 11.5274C8.4628 10.8561 8.43556 9.85298 8.63501 9.10381L10.0059 3.92578"
                    stroke="black"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
                <path
                    d="M3.00781 10.2031H20.9908"
                    stroke="black"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
            </svg>
        )
    } else if (tier === 'Free') {
        return (
            <svg
                width="24px"
                height="24px"
                viewBox="0 0 24 24"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
            >
                <title>Iconly/Two-tone/Unlock</title>
                <g
                    id="Iconly/Two-tone/Unlock"
                    stroke="none"
                    stroke-width="1"
                    fill="none"
                    fill-rule="evenodd"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <g
                        id="Unlock"
                        transform="translate(4.500000, 2.500000)"
                        stroke="#000000"
                        stroke-width="1.5"
                    >
                        <path
                            d="M11.9242,3.06203682 C11.3072,1.28003682 9.6142,0 7.6222,0 C5.1092,-0.00996317625 3.0632,2.01803682 3.0522,4.53103682 L3.0522,4.55103682 L3.0522,6.69803682"
                            id="Stroke-1"
                            opacity="0.400000006"
                        ></path>
                        <path
                            d="M11.433,18.5004368 L3.792,18.5004368 C1.698,18.5004368 1.13686838e-13,16.8024368 1.13686838e-13,14.7074368 L1.13686838e-13,10.4194368 C1.13686838e-13,8.32443682 1.698,6.62643682 3.792,6.62643682 L11.433,6.62643682 C13.527,6.62643682 15.225,8.32443682 15.225,10.4194368 L15.225,14.7074368 C15.225,16.8024368 13.527,18.5004368 11.433,18.5004368 Z"
                            id="Stroke-3"
                        ></path>
                        <line
                            x1="7.6127"
                            y1="11.4526368"
                            x2="7.6127"
                            y2="13.6746368"
                            id="Stroke-5"
                        ></line>
                    </g>
                </g>
            </svg>
        )
    } else {
        return (
            <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    d="M3.51047 11.5274L10.8163 19.5542C11.4507 20.2518 12.5482 20.2518 13.1825 19.5542L20.4894 11.5264C21.1004 10.8561 21.1704 9.85395 20.6596 9.10381L17.7174 4.78197C17.3526 4.24685 16.7455 3.92578 16.0975 3.92578H7.91111C7.26312 3.92578 6.65698 4.24588 6.29212 4.781L3.34117 9.10381C2.8294 9.85298 2.89946 10.8561 3.51047 11.5274Z"
                    stroke="black"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
                <path
                    d="M13.9989 3.92578L15.3659 9.10381C15.5643 9.85395 15.5371 10.8561 15.2997 11.5264L12.4597 19.5542C12.2135 20.2518 11.7864 20.2518 11.5402 19.5542L8.7002 11.5274C8.4628 10.8561 8.43556 9.85298 8.63501 9.10381L10.0059 3.92578"
                    stroke="black"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
                <path
                    d="M3.00781 10.2031H20.9908"
                    stroke="black"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
            </svg>
        )
    }
}
export default function Integrations() {
    const workspace = useParams<{ ws: string }>().ws
    const navigate = useNavigate()
    const [runQuery, setRunQuery] = useAtom(runQueryAtom)
    const [loading, setLoading] = useState(false)
    const [url, setUrl] = useState('')
    const number = window.innerWidth >768 ? 4 : 2

    const [open, setOpen] = useState(false)
    const {
        response: responseConnectors,
        isLoading: connectorsLoading,
        sendNow: getList,
    } = useIntegrationApiV1ConnectorsList(
        number,
        1,
        undefined,
        'count',
        'desc',
        {},
        false
    )
    const [selected, setSelected] = useState()

    useEffect(() => {
        getList(number, 1, 'count', 'desc', false)
    }, [])
     const EnableIntegration = () => {
         setLoading(true)
         let url = ''
         if (window.location.origin === 'http://localhost:3000') {
             url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
         } else {
             url = window.location.origin
         }
         // @ts-ignore
         const token = JSON.parse(localStorage.getItem('openg_auth')).token

         const config = {
             headers: {
                 Authorization: `Bearer ${token}`,
             },
         }

         axios
             .post(
                 `${url}/main/integration/api/v1/integration-types/plugin/${selected?.platform_name}/enable`,
                 {},
                 config
             )
             .then((res) => {
                 getList(number, pageNo, 'count', 'desc', undefined)
                 setLoading(false)
                 setOpen(false)
                 setNotification({
                     text: `Integration enabled`,
                     type: 'success',
                 })
             })
             .catch((err) => {
                 setNotification({
                     text: `Failed to enable integration`,
                     type: 'error',
                 })
                 getList(number, pageNo, 'count', 'desc', undefined)
                 setLoading(false)
             })
     }
     const InstallPlugin = () => {
         setLoading(true)
         let url = ''
         if (window.location.origin === 'http://localhost:3000') {
             url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
         } else {
             url = window.location.origin
         }
         let path = ''
         if (selected?.html_url) {
             path = `/main/integration/api/v1/integration-types/plugin/load/id/${selected?.platform_name}`
         } else {
             path = `/main/integration/api/v1/integration-types/plugin/load/url/${url}`
         }
         // @ts-ignore
         const token = JSON.parse(localStorage.getItem('openg_auth')).token

         const config = {
             headers: {
                 Authorization: `Bearer ${token}`,
             },
         }

         axios
             .post(`${url}${path}`, {}, config)
             .then((res) => {
                 getList(number, pageNo, 'count', 'desc', undefined)
                 setLoading(false)
                 setOpen(false)
                 setNotification({
                     text: `Plugin Installed`,
                     type: 'success',
                 })
             })
             .catch((err) => {
                 setNotification({
                     text: `Failed to install plugin`,
                     type: 'error',
                 })
                 getList(number, pageNo, 'count', 'desc', undefined)
                 setLoading(false)
             })
     }
    return (
        <>
            <Modal
                visible={open}
                onDismiss={() => setOpen(false)}
                header="Plugin Installation"
            >
                <div className="p-4">
                    <Text>
                        This Plugin is{' '}
                        {selected?.installed == 'not_installed'
                            ? 'not installed'
                            : selected?.installed == 'installing'
                            ? 'installing'
                            : 'disabled'}{' '}
                        .
                    </Text>
                    {selected?.installed == 'not_installed' &&
                        selected?.html_url == '' && (
                            <>
                                <Input
                                    className="mt-2"
                                    placeholder="Enter Plugin URL"
                                    value={url}
                                    onChange={({ detail }) =>
                                        setUrl(detail.value)
                                    }
                                />
                            </>
                        )}
                    <Flex
                        justifyContent="end"
                        alignItems="center"
                        flexDirection="row"
                        className="gap-3"
                    >
                        <Button
                            // loading={loading}
                            disabled={loading}
                            onClick={() => setOpen(false)}
                            className="mt-6"
                        >
                            Close
                        </Button>
                        {selected?.installed == 'installing' ? (
                            <>
                                <Button
                                    loading={loading}
                                    disabled={loading}
                                    variant="primary"
                                    onClick={() => {
                                        getList(
                                            number,
                                            1,
                                            'count',
                                            'desc',
                                            false
                                        )
                                        setOpen(false)
                                    }}
                                    className="mt-6"
                                >
                                    Refresh
                                </Button>
                            </>
                        ) : (
                            <>
                                {(selected?.installed == 'not_installed' ||
                                    selected?.enabled == 'disabled') && (
                                    <>
                                        <Button
                                            loading={loading}
                                            disabled={loading}
                                            variant="primary"
                                            onClick={() => {
                                                selected?.installed ==
                                                'not_installed'
                                                    ? InstallPlugin()
                                                    : EnableIntegration()
                                            }}
                                            className="mt-6"
                                        >
                                            {selected?.installed ==
                                            'not_installed'
                                                ? ' Install'
                                                : 'Enable'}
                                        </Button>
                                    </>
                                )}
                            </>
                        )}
                    </Flex>
                </div>
            </Modal>

            <Card className="sm:h-full h-fit  sm:w-full overflow-scroll no-scrollbar sm:inline-block hidden border-solid  border-2 border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:py-2 px-6">
                <Flex
                    justifyContent="between"
                    className="sm:flex-row flex-col sm:mb-0 mb-2 sm:items-center items-start"
                >
                    <Flex justifyContent="start" className="gap-2 ">
                        <Icon icon={MagnifyingGlassIcon} className="p-0" />
                        <Title className="sm:font-semibold">
                            Integrations
                        </Title>
                    </Flex>
                    <a
                        target="__blank"
                        href={`/integration/plugins`}
                        className=" cursor-pointer"
                    >
                        <Button
                            size="xs"
                            variant="light"
                            icon={ChevronRightIcon}
                            iconPosition="right"
                            className="my-3"
                            // onClick={() => {
                            //     navigate(`/cloudql?tab_id=0`)
                            // }}
                        >
                            see all
                        </Button>
                    </a>
                </Flex>

                <Cards
                    ariaLabels={{
                        itemSelectionLabel: (e, t) => `select ${t.name}`,
                        selectionGroupLabel: 'Item selection',
                    }}
                    onSelectionChange={({ detail }) => {
                        const connector = detail?.selectedItems[0]
                        if (
                            connector.enabled === 'disabled' ||
                            connector?.installed === 'not_installed'
                        ) {
                            setOpen(true)
                            setSelected(connector)
                            return
                        }

                        if (
                            connector.enabled == 'enabled' &&
                            connector.installed == 'installed'
                        ) {
                            const name = connector?.name
                            const id = connector?.id
                            navigate(
                                `integration/plugins/${connector.platform_name}`,
                                {
                                    state: {
                                        name,
                                        id,
                                    },
                                }
                            )
                            return
                        }
                        navigate(
                            `${connector.platform_name}/../../request-access?connector=${connector.title}`
                        )
                    }}
                    selectedItems={[]}
                    cardDefinition={{
                        header: (item) => (
                            <Link className="w-100">
                                <div className="w-100 flex flex-row justify-between">
                                    <span>{item.name}</span>
                                </div>
                            </Link>
                        ),
                        sections: [
                            {
                                id: 'logo',

                                content: (item) => (
                                    <div className="w-100 flex flex-row items-center  justify-between  ">
                                        <img
                                            className="sm:w-[50px] sm:h-[50px] w-[30px] h-[30px]"
                                            src={item.logo}
                                            onError={(e) => {
                                                e.currentTarget.onerror = null
                                                e.currentTarget.src =
                                                    'https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/default.svg'
                                            }}
                                            alt="placeholder"
                                        />
                                        {/* <span>{item.status ? 'Enabled' : 'Disable'}</span> */}
                                    </div>
                                ),
                            },
                            // {
                            //     id: 'integrattoin',
                            //     header: 'Description',
                            //     content: (item) => item?.description,
                            //     width: 70,
                            // },

                            // {
                            //     id: 'integrattoin',
                            //     header: 'Integrations',
                            //     content: (item) => (
                            //         <>
                            //             <span className='w-full text-end'>
                            //                 {item?.count ? item.count : '--'}
                            //             </span>
                            //         </>
                            //     ),
                            //     width: 30,
                            // },
                            {
                                id: 'description',
                                header: (
                                    <>
                                        <div className="flex justify-between">
                                            <span className="sm:inline hidden">
                                                {'Description'}
                                            </span>
                                            <span>{'Integrations'}</span>
                                        </div>
                                    </>
                                ),
                                content: (item) => (
                                    <>
                                        <div className="flex justify-between gap-4">
                                            <span className="max-w-40 sm:inline hidden">
                                                {item.description}
                                            </span>
                                            <span>
                                                {item.count ? item.count : '--'}
                                            </span>
                                        </div>
                                    </>
                                ),
                            },
                        ],
                    }}
                    cardsPerRow={[{ cards: 1 }]}
                    // @ts-ignore
                    items={responseConnectors?.items?.map((type) => {
                        return {
                            id: type.id,
                            tier: type.tier,
                            enabled: type.operational_status,
                            installed: type.install_state,
                            platform_name: type.plugin_id,
                            description: type.description,
                            title: type.name,
                            name: type.name,
                            html_url: type.url,
                            count: type?.count?.total,
                            logo: `https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/${type.icon}`,
                        }
                    })}
                    loadingText="Loading resources"
                    stickyHeader
                    entireCardClickable
                    variant="full-page"
                    selectionType="single"
                    trackBy="name"
                    empty={
                        <Box
                            margin={{ vertical: 'xs' }}
                            textAlign="center"
                            color="inherit"
                        >
                            <SpaceBetween size="m">
                                <b>No resources</b>
                            </SpaceBetween>
                        </Box>
                    }
                />
            </Card>
        </>
    )
}
