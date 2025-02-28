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
import { useEffect, useState } from 'react'
import ReactJson from '@microlink/react-json-view'
import {
    AdjustmentsVerticalIcon,
    CheckCircleIcon,
    PlayCircleIcon,
    Square2StackIcon,
    TagIcon,
    VariableIcon,
    XCircleIcon,
} from '@heroicons/react/24/outline'
import {
    PlatformEnginePkgBenchmarkApiListV3ResponseMetaData,
    PlatformEnginePkgComplianceApiConformanceStatus,
    PlatformEnginePkgComplianceApiResourceFinding,
    PlatformEnginePkgControlDetailV3,
    PlatformEnginePkgInventoryApiSmartQueryItem,
    PlatformEnginePkgInventoryApiSmartQueryItemV2,
} from '../../../../../api/api'
import { useComplianceApiV1FindingsResourceCreate } from '../../../../../api/compliance.gen'
import Spinner from '../../../../../components/Spinner'
// import { severityBadge } from '../Controls'
import { isDemoAtom, notificationAtom, queryAtom } from '../../../../../store'
// import Timeline from '../FindingsWithFailure/Detail/Timeline'
import { searchAtom } from '../../../../../utilities/urlstate'
import { dateTimeDisplay } from '../../../../../utilities/dateDisplay'
import Editor from 'react-simple-code-editor'
import { severityBadge } from '../../../Controls'
import { Badge, KeyValuePairs, Popover, Tabs } from '@cloudscape-design/components'
import axios from 'axios'
import { RenderObject } from '../../../../../components/RenderObject'
import ImpactedResources from './ImpactedResources'
import Benchmarks from './Benchmarks'

interface IResourceFindingDetail {
    selectedItem: PlatformEnginePkgControlDetailV3 | undefined
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
            controls: [selectedItem?.id],
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

    useEffect(() => {
        if (selectedItem) {
            console.log(selectedItem)
            GetParams()
            // @ts-ignore
            // setParams(selectedItem?.policy?.parameters)
        }
    }, [selectedItem])
    const getItems = () => {
        const items = [
            {
                label: 'ID',
                value: selectedItem?.id,
            },
            {
                label: 'Resource Tables',
                //    @ts-ignore
                value: (
                    <>
                        <Flex
                            className="gap-2 flex-wrap w-full justify-start items-center"
                            flexDirection="row"
                        >
                            <>
                                {/* @ts-ignore */}
                                {selectedItem?.policy?.list_of_resources?.map(
                                    (key, index) => {
                                        return (
                                            <>
                                                {key ===
                                                selectedItem?.policy
                                                    ?.primary_resource ? (
                                                    <>
                                                        <Popover
                                                            content={
                                                                'This is the table used to record and track incidents related to this control. '
                                                            }
                                                            position="bottom"
                                                        >
                                                            {key},
                                                        </Popover>
                                                    </>
                                                ) : (
                                                    <>{key},</>
                                                )}
                                            </>
                                        )
                                    }
                                )}
                            </>
                        </Flex>
                    </>
                ),
            },
            {
                label: 'Policy Language',
                value: selectedItem?.policy?.language,
            },
            {
                label: 'Frameworks',
                value: selectedItem?.frameworks.length,
            },
        ]
       
    
        return items
    }

    return (
        <>
            {selectedItem ? (
                <>
                    <KeyValuePairs
                        className="mb-8"
                        columns={window.innerWidth > 768 ? 4 : 1}
                        items={getItems()}
                    />
                    <Tabs
                        tabs={[
                            {
                                label: 'Details',
                                id: '1',
                                content: (
                                    <>
                                        <KeyValuePairs
                                            columns={1}
                                            items={[
                                                {
                                                    label: 'Description',
                                                    value: selectedItem?.description
                                                        ? selectedItem?.description
                                                        : 'Not Available',
                                                },
                                                {
                                                    label: 'Parameters',
                                                    value: (
                                                        <>
                                                            {params?.length >
                                                            0 ? (
                                                                <>
                                                                    <Flex
                                                                        flexDirection="col"
                                                                        className="gap-2 mt-2 justify-start items-start"
                                                                    >
                                                                        {/* <Title>
                                                                    Parameters:
                                                                </Title> */}
                                                                        <Flex
                                                                            className="gap-2 flex-wrap w-full justify-start items-center "
                                                                            flexDirection="row"
                                                                        >
                                                                            <>
                                                                                {params?.map(
                                                                                    (
                                                                                        item,
                                                                                        index
                                                                                    ) => {
                                                                                        return (
                                                                                            <span className="inline-flex text-lg items-center gap-x-2.5 rounded-tremor-small bg-tremor-background  pl-2.5 pr-2.5 text-tremor-label  text-tremor-content-strong ring-1 ring-tremor-ring dark:bg-dark-tremor-background dark:text-dark-tremor-content dark:ring-dark-tremor-ring">
                                                                                                {
                                                                                                    // @ts-ignore

                                                                                                    item?.key
                                                                                                }
                                                                                                <span className="h-4 w-px bg-tremor-ring dark:bg-dark-tremor-ring" />
                                                                                                <span className="font-medium text-tremor-content dark:text-dark-tremor-content-emphasis">
                                                                                                    {
                                                                                                        // @ts-ignore
                                                                                                        item?.value
                                                                                                    }
                                                                                                </span>
                                                                                            </span>
                                                                                        )
                                                                                    }
                                                                                )}
                                                                                {params?.length ==
                                                                                    0 &&
                                                                                    'No Parameters'}
                                                                            </>
                                                                        </Flex>
                                                                    </Flex>
                                                                </>
                                                            ) : (
                                                                <>None</>
                                                            )}
                                                        </>
                                                    ),
                                                },
                                                {
                                                    label: 'Policy source',
                                                    value: selectedItem?.has_inline_policy
                                                        ? 'Inline'
                                                        : selectedItem?.policy
                                                              ?.id,
                                                },

                                                {
                                                    label: 'Tags',
                                                    value: (
                                                        <>
                                                            <Flex
                                                                className="gap-2 mt-2 flex-wrap w-full justify-start items-center"
                                                                flexDirection="row"
                                                            >
                                                                <>
                                                                    {/* @ts-ignore */}
                                                                    {Object.entries(
                                                                        selectedItem?.tags
                                                                    ).map(
                                                                        (
                                                                            key,
                                                                            index
                                                                        ) => {
                                                                            return (
                                                                                <>
                                                                                    <span className="inline-flex items-center gap-x-2.5 rounded-tremor-full bg-tremor-background py-1 pl-2.5 pr-2.5 text-tremor-label text-tremor-content ring-1 ring-inset ring-tremor-ring dark:bg-dark-tremor-background dark:text-dark-tremor-content dark:ring-dark-tremor-ring">
                                                                                        {
                                                                                            key[0]
                                                                                        }
                                                                                        <span className="h-4 w-px bg-tremor-ring dark:bg-dark-tremor-ring" />
                                                                                        <span className="font-medium text-tremor-content-strong dark:text-dark-tremor-content-emphasis">
                                                                                            {
                                                                                                key[1]
                                                                                            }
                                                                                        </span>
                                                                                    </span>
                                                                                </>
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
                                    </>
                                ),
                            },
                            {
                                label: 'Policy Definition',
                                id: '0',
                                content: (
                                    <>
                                        <Grid
                                            className="w-full gap-4 mb-6"
                                            numItems={1}
                                        >
                                            <Flex
                                                flexDirection="col"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                {/* <Card className=" py-3 mb-2 relative "> */}
                                                <RenderObject
                                                    obj={
                                                        selectedItem?.policy
                                                            ?.definition
                                                    }
                                                    className="w-full bg-white dark:bg-gray-900 dark:text-gray-50 font-mono text-sm"
                                                    language="sql"
                                                />
                                                {/* </Card> */}

                                                <Flex
                                                    flexDirection="row"
                                                    className=" justify-end w-full mt-2"
                                                >
                                                    {/* <Title className="mb-2">
                                                        Definition
                                                    </Title> */}

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
                                                        <Link to={`/cloudql`}>
                                                            Open in CloudQL
                                                        </Link>{' '}
                                                    </Button>
                                                </Flex>
                                            </Flex>
                                        </Grid>
                                    </>
                                ),
                            },

                            {
                                label: 'Impacted resources',
                                id: '2',
                                content: (
                                    <ImpactedResources
                                        controlId={selectedItem?.id || ''}
                                        linkPrefix={`/score/categories/`}
                                        // conformanceFilter={
                                        //     conformanceFilter
                                        // }
                                    />
                                ),
                            },
                            // {
                            //     id: '3',
                            //     label: 'Impacted Integrations',
                            //     content: (
                            //         <ImpactedAccounts
                            //             controlId={
                            //                 controlDetail?.control?.id
                            //             }
                            //         />
                            //     ),
                            // },

                            {
                                id: '4',
                                label: 'Frameworks',
                                content: (
                                    <Benchmarks
                                        benchmarks={selectedItem?.frameworks}
                                    />
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
