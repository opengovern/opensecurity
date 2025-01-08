import { Link, useParams } from 'react-router-dom'
import { useAtomValue, useSetAtom } from 'jotai'
import {
    Button,
    Card,
    Flex,
    Grid,
    List,
    ListItem,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    Title,
} from '@tremor/react'
import { useEffect } from 'react'
import ReactJson from '@microlink/react-json-view'
import {
    CheckCircleIcon,
    PlayCircleIcon,
    TagIcon,
    XCircleIcon,
} from '@heroicons/react/24/outline'
import {
    GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus,
    GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFinding,
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItem,
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2,
} from '../../../../api/api'
import DrawerPanel from '../../../../components/DrawerPanel'
import { getConnectorIcon } from '../../../../components/Cards/ConnectorCard'
import SummaryCard from '../../../../components/Cards/SummaryCard'
import { useComplianceApiV1FindingsResourceCreate } from '../../../../api/compliance.gen'
import Spinner from '../../../../components/Spinner'
// import { severityBadge } from '../Controls'
import { isDemoAtom, notificationAtom, queryAtom } from '../../../../store'
// import Timeline from '../FindingsWithFailure/Detail/Timeline'
import { searchAtom } from '../../../../utilities/urlstate'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import Editor from 'react-simple-code-editor'
import { highlight, languages } from 'prismjs' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/components/prism-sql' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/themes/prism.css'
import { Badge, KeyValuePairs } from '@cloudscape-design/components'

interface IResourceFindingDetail {
    query:
        | any
        | undefined
    open: boolean
    onClose: () => void
    onRefresh: () => void
    linkPrefix?: string
    setTab: Function
}

export default function ViewDetail({
    query,
    open,
    onClose,
    onRefresh,
    linkPrefix = '',
    setTab,
}: IResourceFindingDetail) {
    const { ws } = useParams()
    const setQuery = useSetAtom(queryAtom)

    // const { response, isLoading, sendNow } =
    //     useComplianceApiV1FindingsResourceCreate(
    //         { kaytuResourceId: resourceFinding?.kaytuResourceID || '' },
    //         {},
    //         false
    //     )
    const searchParams = useAtomValue(searchAtom)

    // useEffect(() => {
    //     if (resourceFinding && open) {
    //         sendNow()
    //     }
    // }, [resourceFinding, open])

    const isDemo = useAtomValue(isDemoAtom)

    // const finding = resourceFinding?.findings
    //     ?.filter((f) => f.controlID === controlID)
    //     .at(0)

    // const conformance = () => {
    //     if (showOnlyOneControl) {
    //         return (finding?.conformanceStatus || 0) ===
    //             GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed ? (
    //             <Flex className="w-fit gap-1.5">
    //                 <XCircleIcon className="h-4 text-rose-600" />
    //                 <Text>Failed</Text>
    //             </Flex>
    //         ) : (
    //             <Flex className="w-fit gap-1.5">
    //                 <CheckCircleIcon className="h-4 text-emerald-500" />
    //                 <Text>Passed</Text>
    //             </Flex>
    //         )
    //     }

    //     const failingControls = new Map<string, string>()
    //     resourceFinding?.findings?.forEach((f) => {
    //         failingControls.set(f.controlID || '', '')
    //     })

    //     return failingControls.size > 0 ? (
    //         <Flex className="w-fit gap-1.5">
    //             <XCircleIcon className="h-4 text-rose-600" />
    //             <Text>{failingControls.size} Failing</Text>
    //         </Flex>
    //     ) : (
    //         <Flex className="w-fit gap-1.5">
    //             <CheckCircleIcon className="h-4 text-emerald-500" />
    //             <Text>Passed</Text>
    //         </Flex>
    //     )
    // }

    return (
        <>
            <KeyValuePairs
                columns={4}
                items={[
                    {
                        label: 'ID',
                        value: query?.id,
                    },
                    {
                        label: 'Title',
                        value: query?.title,
                    },
                    {
                        label: 'Description',
                        value: query?.description,
                    },
                    {
                        label: 'Last time rendered',
                        value: dateTimeDisplay(query?.last_time_rendered),
                    },
                    {
                        label: 'Query Engine',
                        value: query?.query?.engine,
                    },
                    {
                        label: 'Tags',
                        value: (
                            <>
                                {query?.tags && (
                                    <Flex
                                        className="gap-2 flex-wrap min-w-fit"
                                        flexDirection="row"
                                    >
                                        <>
                                            {Object.entries(
                                                // @ts-ignore
                                                query?.tags
                                            ).map((key, index) => {
                                                return (
                                                    <Badge
                                                        color="severity-neutral"
                                                        className=" min-w-fit"
                                                    >
                                                        <Flex
                                                            flexDirection="row"
                                                            justifyContent="start"
                                                            className="hover:cursor-pointer max-w-full w-fit  px-1"
                                                        >
                                                            <TagIcon className="min-w-4 w-4 mr-1" />
                                                            {`${key[0]} : ${key[1]}`}
                                                        </Flex>
                                                    </Badge>
                                                )
                                            })}
                                        </>
                                    </Flex>
                                )}
                            </>
                        ),
                    },
                ]}
            />
            {/* <Grid className="w-full gap-4 mb-6" numItems={1}>
                <Flex
                    flexDirection="row"
                    justifyContent="between"
                    alignItems="start"
                    className="mt-2"
                >
                    <Text className="w-56 font-bold">ID : </Text>
                    <Text className="w-full">{query?.id}</Text>
                </Flex>
                <Flex
                    flexDirection="row"
                    justifyContent="between"
                    alignItems="start"
                    className="mt-2"
                >
                    <Text className="w-56 font-bold">Title : </Text>
                    <Text className="w-full">{query?.title}</Text>
                </Flex>{' '}
                <Flex
                    flexDirection="row"
                    justifyContent="between"
                    alignItems="start"
                    className="mt-2"
                >
                    <Text className="w-56 font-bold">Description : </Text>
                    <Text className="w-full">{query?.description}</Text>
                </Flex>{' '}
                <Flex
                    flexDirection="row"
                    justifyContent="between"
                    alignItems="start"
                    className="mt-2"
                >
                    <Text className="w-56 font-bold">Connector : </Text>
                    <Text className="w-full">
                        {query?.connectors?.map((item, index) => {
                            return `${item} `
                        })}
                    </Text>
                </Flex>
                <Flex
                    flexDirection="row"
                    justifyContent="between"
                    alignItems="start"
                    className="mt-2"
                >
                    <Text className="w-56 font-bold">Query Engine : </Text>
                    <Text className="w-full">
                        {/* @ts-ignore 
                        {query?.query?.engine}
                    </Text>
                </Flex>
            </Grid> */}
            <Flex flexDirection="row" className="mb-2 mt-3">
                <Title className="mb-2">Query</Title>

                <Button
                    icon={PlayCircleIcon}
                    onClick={() => {
                        // @ts-ignore
                        setQuery(query?.query?.query_to_execute)
                        setTab('3')
                        // debugger
                    }}
                    disabled={false}
                    loading={false}
                    loadingText="Running"
                >
                    {/* <Link to={`/cloudql`}> */}
                    Run in Query
                    {/* </Link>{' '} */}
                </Button>
            </Flex>
            <Card className=" py-3 mb-2 relative ">
                <Editor
                    onValueChange={(text) => {}}
                    highlight={(text) => highlight(text, languages.sql, 'sql')}
                    // @ts-ignore
                    value={query?.query?.query_to_execute || ''}
                    className="w-full bg-white dark:bg-gray-900 dark:text-gray-50 font-mono text-sm"
                    style={{
                        minHeight: '200px',
                        // maxHeight: '500px',
                        overflowY: 'scroll',
                    }}
                    placeholder="-- write your SQL query here"
                    disabled={true}
                />
            </Card>
            {/* <TabGroup>
                <TabPanels>
                    <TabPanel>
                        <Flex
                            flexDirection="row"
                            alignItems="start"
                            className="gap-1 w-full flex-wrap "
                            justifyContent="start"
                        >
                            {query?.tags && (
                                <>
                                    {Object.entries(query?.tags).map(
                                        (key, index) => {
                                            return (
                                                <>
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="start"
                                                        className="hover:cursor-pointer max-w-full w-fit bg-gray-200 border-gray-300 rounded-lg border px-1"
                                                    >
                                                        <TagIcon className="min-w-4 w-4 mr-1" />
                                                        <Text className="truncate">
                                                            {/* @ts-ignore 
                                                            {key[0]}:{key[1]}
                                                        </Text>
                                                    </Flex>
                                                </>
                                            )
                                        }
                                    )}
                                </>
                            )}
                        </Flex>
                    </TabPanel>
                </TabPanels>
            </TabGroup> */}
        </>
    )
}
