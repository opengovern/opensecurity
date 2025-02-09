// @ts-nocheck
import {
    Card,
    Flex,
    Grid,
    Tab,
    TabGroup,
    TabList,
    Text,
    Title,
} from '@tremor/react'

import { useEffect, useState } from 'react'
import {
    ArrowDownIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    DocumentTextIcon,
    PlusIcon,
} from '@heroicons/react/24/outline'
import ConnectorCard from '../../components/Cards/ConnectorCard'
import Spinner from '../../components/Spinner'
import { useIntegrationApiV1ConnectorsList } from '../../api/integration.gen'
import TopHeader from '../../components/Layout/Header'
import {
    Box,
    Button,
    Cards,
    Input,
    Link,
    Modal,
    Pagination,
    SpaceBetween,
} from '@cloudscape-design/components'
import { PlatformEngineServicesIntegrationApiEntityTier } from '../../api/api'
import { useNavigate } from 'react-router-dom'
import { get } from 'http'
import axios from 'axios'
import { notificationAtom } from '../../store'
import { useSetAtom } from 'jotai'
import CustomPagination from '../../components/Pagination'

export default function Integrations() {
    const [pageNo, setPageNo] = useState<number>(1)
    const {
        response: responseConnectors,
        isLoading: connectorsLoading,
        sendNow: getList,
    } = useIntegrationApiV1ConnectorsList(12, pageNo, undefined, 'count', 'desc',{},false)
    const {
        response: allIntegrations,
        isLoading: IntegrationsLoading,
        sendNow: getALL,
    } = useIntegrationApiV1ConnectorsList(
        100,
        pageNo,
        undefined,
        'count',
        'desc',{},false
    )
    const [open, setOpen] = useState(false)
    const [openWait, setOpenWait] = useState(false)

    const navigate = useNavigate()
    const [selected, setSelected] = useState()
    const [loading, setLoading] = useState(false)
    const [url, setUrl] = useState('')
    const connectorList = responseConnectors?.items || []
    const setNotification = useSetAtom(notificationAtom)

    // @ts-ignore

    //@ts-ignore
    const totalPages = Math.ceil(responseConnectors?.total_count / 12)
    useEffect(() => {
        getList(12, pageNo, 'count', 'desc', undefined)
        getALL(100, pageNo, 'count', 'desc', undefined)
    }, [pageNo])
   
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
                getList(12, pageNo, 'count', 'desc', undefined)
                getALL(100, pageNo, 'count', 'desc', undefined)
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
                getList(12, pageNo, 'count', 'desc', undefined)
                getALL(100, pageNo, 'count', 'desc', undefined)
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
                getList(12, pageNo, 'count', 'desc', undefined)
                getALL(100, pageNo, 'count', 'desc', undefined)
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
                getList(12, pageNo, 'count', 'desc', undefined)
                getALL(100, pageNo, 'count', 'desc', undefined)
                setLoading(false)
            })
    }
    const CheckInstalling =()=>{
        var flag = false
        allIntegrations?.items?.map((item)=>{
            console.log(item.install_state)
            if(item.install_state =='installing'){
                flag =true
            }
        })
        return flag

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
                                        getList(12, 1, 'count', 'desc', false)
                                        getALL(
                                            100,
                                            pageNo,
                                            'count',
                                            'desc',
                                            undefined
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
            <Modal
                visible={openWait}
                onDismiss={() => setOpenWait(false)}
                header="Plugin Installation"
            >
                <div className="p-4">
                    <Text>
                        Installation is in progress. Please Wait.
                    </Text>

                    <Flex
                        justifyContent="end"
                        alignItems="center"
                        flexDirection="row"
                        className="gap-3"
                    >
                        <Button
                            // loading={loading}
                            disabled={false}
                            onClick={() => setOpenWait(false)}
                            className="mt-6"
                        >
                            Close
                        </Button>
                        <Button
                            loading={loading}
                            disabled={loading}
                            variant="primary"
                            onClick={() => {
                                getList(12, 1, 'count', 'desc', false)
                                getALL(100, pageNo, 'count', 'desc', undefined)
                                setOpenWait(false)
                            }}
                            className="mt-6"
                        >
                            Refresh
                        </Button>
                    </Flex>
                </div>
            </Modal>

            {connectorsLoading ? (
                <Flex className="mt-36">
                    <Spinner />
                </Flex>
            ) : (
                <>
                    <Flex
                        className="bg-white w-full rounded-xl border-solid  border-2 border-gray-200    "
                        flexDirection="col"
                        justifyContent="center"
                        alignItems="center"
                    >
                        <div className="flex sm:flex-row flex-col justify-between w-full p-4 sm:p-6 lg:p-8 lg:pb-0 sm:pb-0 pb-0 mb-4">
                            <span className="sm:text-2xl sm:mb-0 mb-2 sm:min-w-max text-lg font-bold">
                                Integration Plugins
                                {responseConnectors?.total_count
                                    ? ` (${responseConnectors?.total_count})`
                                    : ' ?'}
                            </span>
                            <CustomPagination
                                currentPageIndex={pageNo}
                                pagesCount={totalPages}
                                onChange={({ detail }) => {
                                    setPageNo(detail.currentPageIndex)
                                }}
                            />
                        </div>
                        <div className="w-full">
                            <div className="p-4 sm:p-6 lg:p-8 pt-0">
                                <main>
                                    <div className="flex items-center justify-between">
                                        {/* <h2 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                            Available Dashboards
                                        </h2> */}
                                        <div className="flex items-center space-x-2"></div>
                                    </div>
                                    <div className="flex items-center w-full">
                                        <Cards
                                            ariaLabels={{
                                                itemSelectionLabel: (e, t) =>
                                                    `select ${t.name}`,
                                                selectionGroupLabel:
                                                    'Item selection',
                                            }}
                                            onSelectionChange={({ detail }) => {
                                                const connector =
                                                    detail?.selectedItems[0]
                                                if (
                                                    connector.enabled ===
                                                        'disabled' ||
                                                    connector?.installed ===
                                                        'not_installed' ||
                                                    connector?.installed ===
                                                        'installing'
                                                ) {
                                                    if(CheckInstalling()){
                                                        setOpenWait(true)
                                                        return
                                                    }
                                                    setOpen(true)
                                                    setSelected(connector)
                                                    return
                                                }

                                                if (
                                                    connector.enabled ==
                                                        'enabled' &&
                                                    connector.installed ==
                                                        'installed'
                                                ) {
                                                    const name = connector?.name
                                                    const id = connector?.id
                                                    navigate(
                                                        `${connector.platform_name}`,
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
                                                    <Link
                                                        className="w-100"
                                                        onClick={() => {
                                                            // if (item.tier === 'Community') {
                                                            //     navigate(
                                                            //         '/integrations/' +
                                                            //             item.schema_id +
                                                            //             '/schema'
                                                            //     )
                                                            // } else {
                                                            //     // setOpen(true);
                                                            // }
                                                        }}
                                                    >
                                                        <div className="w-100 flex flex-row justify-between">
                                                            <span className="sm:text-base text-sm">
                                                                {item.title}
                                                            </span>
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
                                                                    src={
                                                                        item.logo
                                                                    }
                                                                    onError={(
                                                                        e
                                                                    ) => {
                                                                        e.currentTarget.onerror =
                                                                            null
                                                                        e.currentTarget.src =
                                                                            'https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/default.svg'
                                                                    }}
                                                                    alt="placeholder"
                                                                />
                                                                {/* <span>{item.status ? 'Enabled' : 'Disable'}</span> */}
                                                            </div>
                                                        ),
                                                    },

                                                    {
                                                        id: 'description',
                                                        header: (
                                                            <>
                                                                <div className="flex justify-between">
                                                                    <span className="sm:inline hidden">
                                                                        {
                                                                            'Description'
                                                                        }
                                                                    </span>
                                                                    <span>
                                                                        {
                                                                            'Integrations'
                                                                        }
                                                                    </span>
                                                                </div>
                                                            </>
                                                        ),
                                                        content: (item) => (
                                                            <>
                                                                <div className="flex justify-between gap-4">
                                                                    <span className="max-w-60 sm:inline hidden">
                                                                        {
                                                                            item.description
                                                                        }
                                                                    </span>
                                                                    <span>
                                                                        {item.count
                                                                            ? item.count
                                                                            : '--'}
                                                                    </span>
                                                                </div>
                                                            </>
                                                        ),
                                                    },
                                                ],
                                            }}
                                            cardsPerRow={[
                                                { cards: 1 },
                                                { minWidth: 320, cards: 2 },
                                                { minWidth: 700, cards: 4 },
                                            ]}
                                            // @ts-ignore
                                            items={connectorList?.map(
                                                (type) => {
                                                    return {
                                                        id: type.id,
                                                        tier: type.tier,
                                                        enabled:
                                                            type.operational_status,
                                                        installed:
                                                            type.install_state,
                                                        platform_name:
                                                            type.plugin_id,
                                                        description:
                                                            type?.description,

                                                        title: type.name,
                                                        name: type.name,
                                                        html_url: type.url,
                                                        count: type?.count
                                                            ?.total,

                                                        logo: `https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/${type.icon}`,
                                                    }
                                                }
                                            )}
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
                                    </div>
                                </main>
                            </div>
                        </div>
                    </Flex>
                </>
            )}
        </>
    )
}
