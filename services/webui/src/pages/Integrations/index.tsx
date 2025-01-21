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

export default function Integrations() {
    const [pageNo, setPageNo] = useState<number>(1)
    const {
        response: responseConnectors,
        isLoading: connectorsLoading,
        sendNow: getList,
    } = useIntegrationApiV1ConnectorsList(9, pageNo, undefined, 'count', 'desc')
    const [open, setOpen] = useState(false)
    const navigate = useNavigate()
    const [selected, setSelected] = useState()
    const [loading, setLoading] = useState(false)
    const [url, setUrl] = useState('')
    const connectorList = responseConnectors?.items || []
    const setNotification = useSetAtom(notificationAtom)

    // @ts-ignore

    //@ts-ignore
    const totalPages = Math.ceil(responseConnectors?.total_count / 9)
    useEffect(() => {
        getList(9, pageNo, 'count', 'desc', undefined)
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
                getList(9, pageNo, 'count', 'desc', undefined)
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
                getList(9, pageNo, 'count', 'desc', undefined)
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
                getList(9, pageNo, 'count', 'desc', undefined)
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
                getList(9, pageNo, 'count', 'desc', undefined)
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
                                        getList(9, 1, 'count', 'desc', false)
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

            {connectorsLoading ? (
                <Flex className="mt-36">
                    <Spinner />
                </Flex>
            ) : (
                <>
                    <Flex
                        className="bg-white w-[90%] rounded-xl border-solid  border-2 border-gray-200  pb-2  "
                        flexDirection="col"
                        justifyContent="center"
                        alignItems="center"
                    >
                        <div className="border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:p-6 lg:p-8">
                            <header>
                                <h1 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                    Integrations
                                </h1>
                                <p className="text-tremor-default text-tremor-content dark:text-dark-tremor-content">
                                    Create and Manage your Integrations
                                </p>
                                <div className="mt-8 w-full md:flex md:max-w-3xl md:items-stretch md:space-x-4">
                                    <Card className="w-full md:w-7/12">
                                        <div className="inline-flex items-center justify-center rounded-tremor-small border border-tremor-border p-2 dark:border-dark-tremor-border">
                                            <DocumentTextIcon
                                                className="size-5 text-tremor-content-emphasis dark:text-dark-tremor-content-emphasis"
                                                aria-hidden={true}
                                            />
                                        </div>
                                        <h3 className="mt-4 text-tremor-default font-medium text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                            <a
                                                href="https://docs.opengovernance.io/"
                                                target="_blank"
                                                className="focus:outline-none"
                                            >
                                                {/* Extend link to entire card */}
                                                <span
                                                    className="absolute inset-0"
                                                    aria-hidden={true}
                                                />
                                                Documentation
                                            </a>
                                        </h3>
                                        <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                            Learn how to add, update, remove
                                            Integrations
                                        </p>
                                    </Card>
                                </div>
                            </header>
                        </div>
                        <div className="w-full">
                            <div className="p-4 sm:p-6 lg:p-8">
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
                                                            <span>
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
                                                                    className="w-[50px] h-[50px]"
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
                                                        id: 'integrattoin',
                                                        header: 'Integrations',
                                                        content: (item) =>
                                                            item?.count
                                                                ? item.count
                                                                : '--',
                                                        width: 100,
                                                    },
                                                ],
                                            }}
                                            cardsPerRow={[
                                                { cards: 1 },
                                                { minWidth: 320, cards: 2 },
                                                { minWidth: 700, cards: 3 },
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
                        <Pagination
                            currentPageIndex={pageNo}
                            pagesCount={totalPages}
                            onChange={({ detail }) => {
                                setPageNo(detail.currentPageIndex)
                            }}
                        />
                    </Flex>
                </>
            )}
        </>
    )
}
