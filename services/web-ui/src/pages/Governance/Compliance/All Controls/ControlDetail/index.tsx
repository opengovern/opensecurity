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
    Square2StackIcon,
    TagIcon,
    XCircleIcon,
} from '@heroicons/react/24/outline'
import {
    GithubComKaytuIoKaytuEnginePkgBenchmarkApiListV3ResponseMetaData,
    GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus,
    GithubComKaytuIoKaytuEnginePkgComplianceApiResourceFinding,
    GithubComKaytuIoKaytuEnginePkgControlDetailV3,
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItem,
    GithubComKaytuIoKaytuEnginePkgInventoryApiSmartQueryItemV2,
} from '../../../../../api/api'
import DrawerPanel from '../../../../../components/DrawerPanel'
import { getConnectorIcon } from '../../../../../components/Cards/ConnectorCard'
import SummaryCard from '../../../../../components/Cards/SummaryCard'
import { useComplianceApiV1FindingsResourceCreate } from '../../../../../api/compliance.gen'
import Spinner from '../../../../../components/Spinner'
// import { severityBadge } from '../Controls'
import { isDemoAtom, notificationAtom, queryAtom } from '../../../../../store'
// import Timeline from '../FindingsWithFailure/Detail/Timeline'
import { searchAtom } from '../../../../../utilities/urlstate'
import { dateTimeDisplay } from '../../../../../utilities/dateDisplay'
import Editor from 'react-simple-code-editor'
import { highlight, languages } from 'prismjs' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/components/prism-sql' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/themes/prism.css'
import { severityBadge } from '../../../Controls'
import { Badge, KeyValuePairs, Tabs } from '@cloudscape-design/components'

interface IResourceFindingDetail {
    selectedItem: GithubComKaytuIoKaytuEnginePkgControlDetailV3 | undefined
    open: boolean
    onClose: () => void
    onRefresh: () => void
    linkPrefix?: string
}

export default function ControlDetail({
    selectedItem,
    open,
    onClose,
    onRefresh,
    linkPrefix = '',
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
            {selectedItem ? (
                <>
                    <Tabs
                        tabs={[
                            {
                                label: 'Summary',
                                id: '0',
                                content: (
                                    <>
                                        <KeyValuePairs
                                            columns={3}
                                            items={[
                                                {
                                                    label: 'ID',
                                                    value: selectedItem?.id,
                                                },
                                                {
                                                    label: 'Title',
                                                    value: selectedItem?.title,
                                                },

                                                {
                                                    label: 'Integration Type',
                                                    value: selectedItem?.integrationType?.map(
                                                        (item, index) => {
                                                            return `${item} `
                                                        }
                                                    ),
                                                },
                                                {
                                                    label: 'Severity',
                                                    value: severityBadge(
                                                        selectedItem?.severity
                                                    ),
                                                },
                                                {
                                                    label: 'Description',
                                                    value: selectedItem?.description,
                                                },
                                                {
                                                    label: 'Tags',
                                                    value: (
                                                        <>
                                                            <Flex
                                                                className="gap-2 flex-wrap"
                                                                flexDirection="row"
                                                            >
                                                                <>
                                                                    {Object.entries(
                                                                        selectedItem?.tags
                                                                    ).map(
                                                                        (
                                                                            key,
                                                                            index
                                                                        ) => {
                                                                            return (
                                                                                <Badge color="severity-neutral">
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
                                                                        }
                                                                    )}
                                                                </>
                                                            </Flex>
                                                        </>
                                                    ),
                                                },
                                            ]}
                                        />
                                        <Grid
                                            className="w-full gap-4 mb-6"
                                            numItems={1}
                                        >
                                            {/* <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    ID :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.id}
                                                </Text>
                                            </Flex>
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Title :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.title}
                                                </Text>
                                            </Flex>{' '}
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Description :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.description}
                                                </Text>
                                            </Flex>{' '}
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Connector :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.connector?.map(
                                                        (item, index) => {
                                                            return `${item} `
                                                        }
                                                    )}
                                                </Text>
                                            </Flex>
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Severity :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {severityBadge(
                                                        selectedItem?.severity
                                                    )}
                                                </Text>
                                            </Flex> */}
                                            <Flex
                                                flexDirection="col"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Flex
                                                    flexDirection="row"
                                                    className="mb-2"
                                                >
                                                    <Title className="mb-2">
                                                        Query
                                                    </Title>

                                                    <Button
                                                        icon={PlayCircleIcon}
                                                        onClick={() => {
                                                            // @ts-ignore
                                                            setQuery(
                                                                selectedItem
                                                                    ?.policy
                                                                    ?.definition
                                                            )
                                                        }}
                                                        disabled={false}
                                                        loading={false}
                                                        loadingText="Running"
                                                    >
                                                        <Link
                                                            to={`/cloudql`}
                                                        >
                                                            Run in Query
                                                        </Link>{' '}
                                                    </Button>
                                                </Flex>
                                                <Card className=" py-3 mb-2 relative ">
                                                    <Editor
                                                        onValueChange={(
                                                            text
                                                        ) => {
                                                            console.log(text)
                                                        }}
                                                        highlight={(text) =>
                                                            highlight(
                                                                text,
                                                                languages.sql,
                                                                'sql'
                                                            )
                                                        }
                                                        // @ts-ignore
                                                        value={
                                                            selectedItem?.policy
                                                                ?.definition
                                                        }
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
                                                {/* <Flex
                                                    flexDirection="row"
                                                    alignItems="start"
                                                    className="gap-1 w-full flex-wrap "
                                                    justifyContent="start"
                                                >
                                                    {}
                                                    {Object.entries(
                                                        selectedItem?.tags
                                                    ).map((key, index) => {
                                                        return (
                                                            <>
                                                                <Flex
                                                                    flexDirection="row"
                                                                    justifyContent="start"
                                                                    className="hover:cursor-pointer max-w-full w-fit bg-gray-200 border-gray-300 rounded-lg border px-1"
                                                                >
                                                                    <TagIcon className="min-w-4 w-4 mr-1" />
                                                                    <Text className="truncate">
                                                                        {key[0]}
                                                                        :
                                                                        {key[1]}
                                                                    </Text>
                                                                </Flex>
                                                            </>
                                                        )
                                                    })}
                                                </Flex> */}
                                            </Flex>
                                        </Grid>
                                    </>
                                ),
                            },
                            {
                                label: 'Benchmark',
                                id: '1',
                                content: (
                                    <>
                                        <KeyValuePairs
                                            columns={2}
                                            items={[
                                                {
                                                    label: 'Has Root',
                                                    value:
                                                        selectedItem?.benchmarks
                                                            ?.roots?.length > 0
                                                            ? 'True'
                                                            : 'False',
                                                },
                                                {
                                                    label: 'Full Paths',
                                                    value: selectedItem?.benchmarks?.fullPath?.map(
                                                        (item, index) => {
                                                            return `${item} `
                                                        }
                                                    ),
                                                },
                                            ]}
                                        />
                                        {/* <Grid
                                            className="w-full gap-4 mb-6"
                                            numItems={1}
                                        >
                                            <>
                                                <Flex
                                                    flexDirection="row"
                                                    justifyContent="between"
                                                    alignItems="start"
                                                    className="mt-2"
                                                >
                                                    <Text className="w-56 font-bold">
                                                        Has Root :{' '}
                                                    </Text>
                                                    <Text className="w-full">
                                                        {selectedItem
                                                            ?.benchmarks?.roots
                                                            ?.length > 0
                                                            ? 'True'
                                                            : 'False'}
                                                    </Text>
                                                </Flex>
                                                <Grid
                                                    className="w-full gap-4 mb-6"
                                                    numItems={1}
                                                >
                                                    <Flex
                                                        flexDirection="row"
                                                        justifyContent="between"
                                                        alignItems="start"
                                                        className="mt-2 flex-wrap"
                                                    >
                                                        <Text className="w-56 font-bold">
                                                            Full Paths :{' '}
                                                        </Text>
                                                        {selectedItem?.benchmarks?.fullPath?.map(
                                                            (item, index) => {
                                                                return (
                                                                    <Text className="">
                                                                        {item}
                                                                    </Text>
                                                                )
                                                            }
                                                        )}
                                                    </Flex>
                                                </Grid>
                                            </>
                                        </Grid> */}
                                    </>
                                ),
                            },
                        ]}
                    />
                </>
            ) : (
                <Spinner />
            )}
        </>
    )
}
