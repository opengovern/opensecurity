import { Link, useNavigate, useParams } from 'react-router-dom'
import {
    Badge,
    Button,
    Card,
    Divider,
    Flex,
    Grid,
    Icon,
    List,
    ListItem,
    Switch,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    Title,
} from '@tremor/react'
import {
    AdjustmentsVerticalIcon,
    BookOpenIcon,
    ChevronRightIcon,
    ClockIcon,
    CodeBracketIcon,
    CodeBracketSquareIcon,
    Cog8ToothIcon,
    CommandLineIcon,
    DocumentDuplicateIcon,
    InformationCircleIcon,
    PencilIcon,
    Square2StackIcon,
} from '@heroicons/react/24/outline'
import { useSetAtom } from 'jotai/index'
import clipboardCopy from 'clipboard-copy'
import { useEffect, useState } from 'react'
import Editor from 'react-simple-code-editor'
import { highlight, languages } from 'prismjs'
import Markdown from 'react-markdown'
import MarkdownPreview from '@uiw/react-markdown-preview'
import { useComplianceApiV1ControlsSummaryDetail } from '../../../../api/compliance.gen'
import { notificationAtom, queryAtom } from '../../../../store'
import { severityBadge } from '../index'
import Spinner from '../../../../components/Spinner'
import Detail from './Tabs/Detail'
import ImpactedResources from './Tabs/ImpactedResources'
import Benchmarks from './Tabs/Benchmarks'
import ImpactedAccounts from './Tabs/ImpactedAccounts'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import TopHeader from '../../../../components/Layout/Header'
import ControlFindings from './Tabs/ControlFindings'
import { useMetadataApiV1QueryParameterList } from '../../../../api/metadata.gen'
import { toErrorMessage } from '../../../../types/apierror'
import { PlatformEnginePkgComplianceApiConformanceStatus } from '../../../../api/api'
import CodeEditor from '@cloudscape-design/components/code-editor'
import {
    Box,
    BreadcrumbGroup,
    Container,
    CopyToClipboard,
    Header,
    KeyValuePairs,
    Modal,
    SpaceBetween,
    Tabs,
} from '@cloudscape-design/components'
import KGrid from '@cloudscape-design/components/grid'
import SegmentedControl from '@cloudscape-design/components/segmented-control'
import axios from 'axios'

export default function ControlDetail() {
    const { controlId, ws } = useParams()
    const setNotification = useSetAtom(notificationAtom)
    const navigate = useNavigate()

    const [doc, setDoc] = useState('')
    const [docTitle, setDocTitle] = useState('')
    const setQuery = useSetAtom(queryAtom)
    const [params, setParams] = useState([])

    const GetParams = () => {
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

        let body: any = {
            controls: [controlDetail?.control?.id],
            cursor: 1,
            per_page: 300,
        }

        axios
            .post(`${url}/main/core/api/v1/query_parameter`, body, config)
            .then((res) => {
                const data = res.data
                setParams(data?.items)
            })
            .catch((err) => {
                console.log(err)
            })
    }

    const {
        response: controlDetail,
        isLoading,
        error: controlDetailError,
        sendNow: refreshControlDetail,
    } = useComplianceApiV1ControlsSummaryDetail(String(controlId))
    // const {
    //     response: parameters,
    //     isLoading: parametersLoading,
    //     isExecuted,
    //     error: parametersError,
    //     sendNow: refresh,
    // } = useMetadataApiV1QueryParameterList()
    // const [conformanceFilter, setConformanceFilter] = useState<
    //     | PlatformEnginePkgComplianceApiConformanceStatus[]
    //     | undefined
    // >(undefined)
    // const conformanceFilterIdx = () => {
    //     if (
    //         conformanceFilter?.length === 1 &&
    //         conformanceFilter[0] ===
    //             PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed
    //     ) {
    //         return '1'
    //     }
    //     if (
    //         conformanceFilter?.length === 1 &&
    //         conformanceFilter[0] ===
    //             PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed
    //     ) {
    //         return '2'
    //     }
    //     return '0'
    // }
    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 600 ? text.substring(0, 600) + '...' : text
        }
    }
    const GetBreadCrumb = () => {
        const temp = []
        if (window.location.pathname.includes('incident')) {
            temp.push({
                text: 'Incidents',
                href: `/incidents`,
            })
        } else if (window.location.pathname.includes('compliance')) {
            temp.push({
                text: 'Compliance',
                href: `/compliance`,
            })
            const temp2 = window.location.pathname.split('/')
            if (temp2.length > 3) {
                temp.push({
                    text: temp2[2],
                    href: `/compliance/${temp2[2]}`,
                })
            }
           
        }
        temp.push({ text: 'Control Detail', href: '#' })

        return temp
    }
    const GetKeyValue = () => {
        const temp = [
            {
                label: 'Control ID',
                value: (
                    // @ts-ignore
                    <CopyToClipboard
                        variant="inline"
                        textToCopy={controlDetail?.control?.id || ''}
                        copySuccessText="Control ID copied to clipboard"
                    />
                ),
            },
            {
                label: 'Inspected Resources:',
                value: (
                    // @ts-ignore
                    <>{controlDetail?.totalResourcesCount}</>
                ),
            },
        ]
        if (controlDetail?.resourceType?.resource_type) {
            temp.push({
                label: 'Resource type',
                value: (
                    // @ts-ignore
                    <CopyToClipboard
                        variant="inline"
                        textToCopy={
                            controlDetail?.resourceType?.resource_type || ''
                        }
                        copySuccessText="Resource type copied to clipboard"
                    />
                ),
            })
        }
        temp.push(
            {
                label: 'Compliant Resources',
                value: (
                    // @ts-ignore
                    <Text className="text-emerald-500">
                        {(controlDetail?.totalResourcesCount || 0) -
                            (controlDetail?.failedResourcesCount || 0)}
                    </Text>
                ),
            },
            {
                label: 'Last updated',
                value: (
                    // @ts-ignore
                    <>{dateTimeDisplay(controlDetail?.control?.updatedAt)}</>
                ),
            },
            {
                label: ' Non-Compliant Resources',
                value: (
                    // @ts-ignore
                    <Text className="text-rose-600">
                        {' '}
                        {controlDetail?.failedResourcesCount}
                    </Text>
                ),
            }
        )
        return temp
        // @ts-ignore
    }
    useEffect(() => {
        if (controlDetail?.control?.id) {
            GetParams()
        }
    }, [controlDetail])

    return (
        <>
         
            {isLoading ? (
                <Spinner className="mt-56" />
            ) : (
                <>
                    {controlDetail ? (
                        <>
                            <BreadcrumbGroup
                                onClick={(event) => {
                                    // event.preventDefault()
                                }}
                                items={GetBreadCrumb()}
                                ariaLabel="Breadcrumbs"
                            />
                            <Container
                                disableHeaderPaddings
                                disableContentPaddings
                                className="rounded-xl  bg-[#0f2940] p-0 text-white mt-4 "
                                header={
                                    <Header
                                        className={`bg-[#0f2940] p-4 pt-0 rounded-xl   text-white ${
                                            false ? 'rounded-b-none' : ''
                                        }`}
                                        variant="h2"
                                        description=""
                                    >
                                        <SpaceBetween
                                            size="xxxs"
                                            direction="vertical"
                                        >
                                            <Box className="rounded-xl same text-white pt-3 pl-3 pb-0">
                                                <KGrid
                                                    gridDefinition={[
                                                        {
                                                            colspan: {
                                                                default: 12,
                                                                xs: 8,
                                                                s: 9,
                                                            },
                                                        },
                                                        {
                                                            colspan: {
                                                                default: 12,
                                                                xs: 4,
                                                                s: 3,
                                                            },
                                                        },
                                                    ]}
                                                >
                                                    <div>
                                                        <Box
                                                            variant="h1"
                                                            className="text-white important w-full"
                                                        >
                                                            <span className="text-white w-full">
                                                                {
                                                                    controlDetail
                                                                        ?.control
                                                                        ?.title
                                                                }
                                                            </span>
                                                        </Box>
                                                        <Box
                                                            variant="p"
                                                            margin={{
                                                                top: 'xxs',
                                                                bottom: 's',
                                                            }}
                                                        >
                                                            <div className="group text-white important  relative flex text-wrap justify-start">
                                                                <Text className="test-start w-full text-white ">
                                                                    {/* @ts-ignore */}
                                                                    {truncate(
                                                                        controlDetail
                                                                            ?.control
                                                                            ?.description
                                                                    )}
                                                                </Text>
                                                                <Card className="absolute w-full text-wrap z-40 top-0 scale-0 transition-all p-2 group-hover:scale-100">
                                                                    <Text>
                                                                        {
                                                                            controlDetail
                                                                                ?.control
                                                                                ?.description
                                                                        }
                                                                    </Text>
                                                                </Card>
                                                            </div>
                                                        </Box>
                                                    </div>
                                                </KGrid>
                                            </Box>
                                            <Flex className="w-max pl-3">
                                                {' '}
                                                {severityBadge(
                                                    controlDetail?.control
                                                        ?.severity
                                                )}
                                            </Flex>
                                        </SpaceBetween>
                                    </Header>
                                }
                            ></Container>
                         
                            <Grid
                                numItems={2}
                                className=" w-full gap-4 mb-6 mt-4"
                            >
                                <Card className="h-fit min-h-[258px] max-h-[258px] overflow-scroll">
                                    <KeyValuePairs
                                        columns={2}
                                        items={GetKeyValue()}
                                    />
                                    <Flex
                                        flexDirection="col"
                                        className="gap-2 mt-2 justify-start items-start"
                                    >
                                        <Title>Parameters:</Title>
                                        <Flex
                                            className="gap-1 flex-row justify-start w-full flex-wrap"
                                            flexDirection="row"
                                        >
                                            <>
                                                {params?.map((item, index) => {
                                                    return (
                                                        <Badge color="severity-neutral">
                                                            <Flex
                                                                flexDirection="row"
                                                                justifyContent="start"
                                                                className="hover:cursor-pointer max-w-full w-fit  px-1"
                                                            >
                                                                <AdjustmentsVerticalIcon className="min-w-4 w-4 mr-1" />
                                                                {/* @ts-ignore */}
                                                                {`${item?.key} : ${item?.value}`}
                                                            </Flex>
                                                        </Badge>
                                                    )
                                                })}
                                                {params?.length == 0 &&
                                                    'No Parameters'}
                                            </>
                                        </Flex>
                                    </Flex>
                                
                                </Card>
                                <Card className="max-h-[258px] overflow-scroll">
                                    <Editor
                                        onValueChange={() => 1}
                                        highlight={(text) =>
                                            highlight(
                                                text,
                                                languages.sql,
                                                'sql'
                                            )
                                        }
                                        value={
                                            controlDetail?.control?.policy?.definition?.replace(
                                                '$IS_ALL_CONNECTIONS_QUERY',
                                                'true'
                                            ) || ''
                                        }
                                        className="w-full bg-white dark:bg-gray-900 dark:text-gray-50 font-mono text-sm"
                                        style={{
                                            minHeight: '228px',
                                        }}
                                        placeholder="-- write your SQL query here"
                                    />
                                    <Divider />
                                    <Flex>
                                        <Flex className="gap-4">
                                            <Button
                                                variant="light"
                                                icon={DocumentDuplicateIcon}
                                                iconPosition="left"
                                                onClick={() =>
                                                    clipboardCopy(
                                                        controlDetail?.control?.policy?.definition?.replace(
                                                            '$IS_ALL_CONNECTIONS_QUERY',
                                                            'true'
                                                        ) || ''
                                                    ).then(() =>
                                                        setNotification({
                                                            text: 'Query copied to clipboard',
                                                            type: 'info',
                                                        })
                                                    )
                                                }
                                            >
                                                Copy SQL Policy
                                            </Button>
                                            <Button
                                                variant="secondary"
                                                onClick={() => {
                                                    setQuery(
                                                        controlDetail?.control?.policy?.definition?.replace(
                                                            '$IS_ALL_CONNECTIONS_QUERY',
                                                            'true'
                                                        ) || ''
                                                    )
                                                }}
                                            >
                                                <Link to={`/cloudql`}>
                                                    Open in CloudQL
                                                </Link>
                                            </Button>
                                        </Flex>
                                    </Flex>
                                </Card>
                    
                                <Flex
                                    flexDirection="col"
                                    alignItems="start"
                                    justifyContent="start"
                                    className="h-full"
                                >
                                    <Modal
                                        header={docTitle}
                                        visible={doc.length > 0}
                                        onDismiss={() => setDoc('')}
                                    >
                                        <MarkdownPreview
                                            source={doc}
                                            className="!bg-transparent"
                                            wrapperElement={{
                                                'data-color-mode': 'light',
                                            }}
                                            rehypeRewrite={(
                                                node,
                                                index,
                                                parent
                                            ) => {
                                                if (
                                                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                                    // @ts-ignore
                                                    node.tagName === 'a' &&
                                                    parent &&
                                                    /^h(1|2|3|4|5|6)/.test(
                                                        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                                        // @ts-ignore
                                                        parent.tagName
                                                    )
                                                ) {
                                                    // eslint-disable-next-line no-param-reassign
                                                    parent.children =
                                                        parent.children.slice(1)
                                                }
                                            }}
                                        />
                                    </Modal>
                                   
                                </Flex>
                            </Grid>
                           
                            <Tabs
                                tabs={[
                                    {
                                        label: 'Impacted resources',
                                        id: '0',
                                        content: (
                                            <ImpactedResources
                                                controlId={
                                                    controlDetail?.control
                                                        ?.id || ''
                                                }
                                                linkPrefix={`/score/categories/`}
                                                // conformanceFilter={
                                                //     conformanceFilter
                                                // }
                                            />
                                        ),
                                    },
                                    {
                                        id: '1',
                                        label: 'Impacted Integrations',
                                        content: (
                                            <ImpactedAccounts
                                                controlId={
                                                    controlDetail?.control?.id
                                                }
                                            />
                                        ),
                                    },
                                  
                                    {
                                        id: '3',
                                        label: 'Frameworks',
                                        content: (
                                            <Benchmarks
                                                benchmarks={
                                                    controlDetail?.benchmarks
                                                }
                                            />
                                        ),
                                    },
                                  
                                ]}
                            />
                            {/* <TabGroup>
                        <Flex
                            flexDirection="row"
                            justifyContent="between"
                            className="mb-2"
                        >
                            <div className="w-fit">
                                <TabList>
                                    <Tab>Impacted resources</Tab>
                                    <Tab>Impacted accounts</Tab>
                                    <Tab
                                        disabled={
                                            controlDetail?.control?.explanation
                                                ?.length === 0 &&
                                            controlDetail?.control
                                                ?.nonComplianceCost?.length ===
                                                0 &&
                                            controlDetail?.control
                                                ?.usefulExample?.length === 0
                                        }
                                    >
                                        Control information
                                    </Tab>
                                    <Tab>Benchmarks</Tab>
                                    <Tab>Findings</Tab>
                                </TabList>
                            </div>
                        </Flex>
                        <TabPanels>
                            <TabPanel>
                                <ImpactedResources
                                    controlId={controlDetail?.control?.id || ''}
                                    linkPrefix={`/score/categories/`}
                                    conformanceFilter={conformanceFilter}
                                />
                            </TabPanel>
                            <TabPanel>
                                <ImpactedAccounts
                                    controlId={controlDetail?.control?.id}
                                />
                            </TabPanel>
                            <TabPanel>
                                <Detail control={controlDetail?.control} />
                            </TabPanel>
                            <TabPanel>
                                <Benchmarks
                                    benchmarks={controlDetail?.benchmarks}
                                />
                            </TabPanel>
                            <TabPanel>
                                <ControlFindings
                                    controlId={controlDetail?.control?.id}
                                />
                            </TabPanel>
                        </TabPanels>
                    </TabGroup> */}
                            {toErrorMessage(controlDetailError) && (
                                <Flex
                                    flexDirection="col"
                                    justifyContent="between"
                                    className="absolute top-0 w-full left-0 h-full backdrop-blur"
                                >
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="center"
                                        alignItems="center"
                                    >
                                        <Title className="mt-6">
                                            Failed to load component
                                        </Title>
                                        <Text className="mt-2">
                                            {toErrorMessage(controlDetailError)}
                                        </Text>
                                    </Flex>
                                    <Button
                                        variant="secondary"
                                        className="mb-6"
                                        color="slate"
                                        onClick={() => {
                                            refreshControlDetail()
                                        }}
                                    >
                                        Try Again
                                    </Button>
                                </Flex>
                            )}
                        </>
                    ) : (
                        <>
                            {toErrorMessage(controlDetailError) && (
                                <Flex
                                    flexDirection="col"
                                    justifyContent="between"
                                    className="absolute top-0 w-full left-0 h-full backdrop-blur"
                                >
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="center"
                                        alignItems="center"
                                    >
                                        <Title className="mt-6">
                                            Failed to load component
                                        </Title>
                                        <Text className="mt-2">
                                            {toErrorMessage(controlDetailError)}
                                        </Text>
                                    </Flex>
                                    <Button
                                        variant="secondary"
                                        className="mb-6"
                                        color="slate"
                                        onClick={() => {
                                            refreshControlDetail()
                                        }}
                                    >
                                        Try Again
                                    </Button>
                                </Flex>
                            )}
                        </>
                    )}
                </>
            )}
        </>
    )
}
