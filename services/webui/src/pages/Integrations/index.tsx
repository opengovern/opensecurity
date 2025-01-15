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
    const connectorList = responseConnectors?.integration_types || []
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
            .put(
                `${url}/main/integration/api/v1/integrations/types/${selected?.platform_name}/enable`,
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
    return (
        <>
            <Modal
                visible={open}
                onDismiss={() => setOpen(false)}
                header="Integration Disabled"
            >
                <div className="p-8">
                    <Text>This integration is disabled.</Text>
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
                        <Button
                            loading={loading}
                            disabled={loading}
                            variant="primary"
                            onClick={() => EnableIntegration()}
                            className="mt-6"
                        >
                            Enable
                        </Button>
                    </Flex>
                </div>
            </Modal>
            {/* <TopHeader /> */}
            {/* <Grid numItems={3} className="gap-4 mb-10">
                <OnboardCard
                    title="Active Accounts"
                    active={topMetrics?.connectionsEnabled}
                    inProgress={topMetrics?.inProgressConnections}
                    healthy={topMetrics?.healthyConnections}
                    unhealthy={topMetrics?.unhealthyConnections}
                    loading={metricsLoading}
                />
            </Grid> */}
            {connectorsLoading ? (
                <Flex className="mt-36">
                    <Spinner />
                </Flex>
            ) : (
                <>
                    {/* <TabGroup className='mt-4'>
                        <TabList>
                            <Tab>test</Tab>
                            <Tab>test</Tab>
                            <Tab>test</Tab>
                            <Tab>test</Tab>
                        </TabList>
                    </TabGroup> */}
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
                                        <div className="flex items-center space-x-2">
                                            {/* <Select
                                        placeholder="Sorty by"
                                        enableClear={false}
                                        className="[&>button]:rounded-tremor-small"
                                    >
                                        <SelectItem value="1">Name</SelectItem>
                                        <SelectItem value="2">
                                            Last edited
                                        </SelectItem>
                                        <SelectItem value="3">Size</SelectItem>
                                    </Select> */}
                                            {/* <button
                                                type="button"
                                                onClick={() => {
                                                    f()
                                                    setOpen(true)
                                                }}
                                                className="hidden h-9 items-center gap-1.5 whitespace-nowrap rounded-tremor-small bg-tremor-brand px-3 py-2.5 text-tremor-default font-medium text-tremor-brand-inverted shadow-tremor-input hover:bg-tremor-brand-emphasis dark:bg-dark-tremor-brand dark:text-dark-tremor-brand-inverted dark:shadow-dark-tremor-input dark:hover:bg-dark-tremor-brand-emphasis sm:inline-flex"
                                            >
                                                <PlusIcon
                                                    className="-ml-1 size-5 shrink-0"
                                                    aria-hidden={true}
                                                />
                                                Create new Dashboard
                                            </button> */}
                                        </div>
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
                                                        false &&
                                                    connector?.tier ===
                                                        PlatformEngineServicesIntegrationApiEntityTier.TierCommunity
                                                ) {
                                                    setOpen(true)
                                                    setSelected(connector)
                                                    return
                                                }

                                                if (
                                                    connector?.tier ===
                                                    PlatformEngineServicesIntegrationApiEntityTier.TierCommunity
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
                                                            {/* <div className="flex flex-row gap-1 items-center">
                                    {GetTierIcon(item.tier)}
                                    <span className="text-white">{item.tier}</span>
                                </div> */}
                                                        </div>
                                                    </Link>
                                                ),
                                                sections: [
                                                    {
                                                        id: 'logo',
                                                        // header :(<>
                                                        //     <div className="flex justify-end">
                                                        //         <span>{'Status'}</span>
                                                        //     </div>
                                                        // </>),

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
                                                    // {
                                                    //     id: 'description',
                                                    //     header: (
                                                    //         <>
                                                    //             <div className="flex justify-between">
                                                    //                 <span>{'Description'}</span>
                                                    //                 <span>{'Table'}</span>
                                                    //             </div>
                                                    //         </>
                                                    //     ),
                                                    //     content: (item) => (
                                                    //         <>
                                                    //             <div className="flex justify-between">
                                                    //                 <span className="max-w-60">
                                                    //                     {item.description}
                                                    //                 </span>
                                                    //                 <span>
                                                    //                     {item.count ? item.count : '--'}
                                                    //                 </span>
                                                    //             </div>
                                                    //         </>
                                                    //     ),
                                                    // },
                                                    // {
                                                    //     id: 'status',
                                                    //     header: 'Status',
                                                    //     content: (item) =>
                                                    //         item.status ? 'Enabled' : 'Disabled',
                                                    //     width: 70,
                                                    // },
                                                    {
                                                        id: 'integrattoin',
                                                        header: 'Integrations',
                                                        content: (item) =>
                                                            item?.count
                                                                ? item.count
                                                                : '--',
                                                        width: 100,
                                                    },
                                                    // {
                                                    //   id: "tier",
                                                    //   header: "Tier",
                                                    //   content: (item) => item.tier,
                                                    //   width: 85,
                                                    // },
                                                    // {
                                                    //   id: "tables",
                                                    //   header: "Table",
                                                    //   content: (item) => (item.count ? item.count : "--"),
                                                    //   width: 15,
                                                    // },
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
                                                        enabled: type.enabled,
                                                        platform_name:
                                                            type.platform_name,
                                                        // description: type.Description,
                                                        title: type.title,
                                                        name: type.name,
                                                        count: type?.count
                                                            ?.total,
                                                        // schema_id: type?.schema_ids[0],
                                                        // SourceCode: type.SourceCode,
                                                        logo: `https://raw.githubusercontent.com/opengovern/website/main/connectors/icons/${type.logo}`,
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
