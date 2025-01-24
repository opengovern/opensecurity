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
                label: 'Policy Language',
                value: selectedItem?.policy?.language,
            },
        ]
        //    if (selectedItem?.policy.type == 'external') {
        //        items.push({
        //            label: 'Policy Id',
        //            value: selectedItem?.policy?.reference,
        //        })
        //    }
        items.push({
            label: 'Tables',
            //    @ts-ignore
            value: (
                <>
                    <Flex className="gap-2 flex-wrap" flexDirection="row">
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
                                                        {key}
                                                    </Popover>
                                                </>
                                            ) : (
                                                <>{key}</>
                                            )}
                                        </>
                                    )
                                }
                            )}
                        </>
                    </Flex>
                </>
            ),
        })
    
        return items
    }

    return (
        <>
            {selectedItem ? (
                <>
                    <KeyValuePairs className='mb-8' columns={4} items={getItems()} />
                    <Tabs
                        tabs={[
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
                                                    className="mb-2 justify-end w-full"
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
                                label: 'Metadata',
                                id: '1',
                                content: (
                                    <>
                                        <KeyValuePairs
                                            columns={2}
                                            items={[
                                                // {
                                                //     label: 'Has Root',
                                                //     value:
                                                //         selectedItem?.benchmarks
                                                //             ?.roots?.length > 0
                                                //             ? 'True'
                                                //             : 'False',
                                                // },
                                                // {
                                                //     label: 'Full Paths',
                                                //     value: selectedItem?.benchmarks?.fullPath?.map(
                                                //         (item, index) => {
                                                //             return `${item} `
                                                //         }
                                                //     ),
                                                // },
                                                {
                                                    label: 'Description',
                                                    value: selectedItem?.description,
                                                },
                                                 {
                                                    label: 'last updatet at',
                                                    value: selectedItem?.description
                                                },
                                                {
                                                    label: 'Parameters',
                                                    value: (
                                                        <>
                                                            <Flex
                                                                flexDirection="col"
                                                                className="gap-2 mt-2 justify-start items-start"
                                                            >
                                                                {/* <Title>
                                                                    Parameters:
                                                                </Title> */}
                                                                <Flex
                                                                    className="gap-1 flex-wrap w-full justify-start "
                                                                    flexDirection="row"
                                                                >
                                                                    <>
                                                                        {params?.map(
                                                                            (
                                                                                item,
                                                                                index
                                                                            ) => {
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
                                                                            }
                                                                        )}
                                                                        {params?.length ==
                                                                            0 &&
                                                                            'No Parameters'}
                                                                    </>
                                                                </Flex>
                                                            </Flex>
                                                        </>
                                                    ),
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
                                                                    {/* @ts-ignore */}
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
                                    </>
                                ),
                            },
                            {
                                                                    label: 'Impacted resources',
                                                                    id: '2',
                                                                    content: (
                                                                        <ImpactedResources
                                                                            controlId={
                                                                              selectedItem?.id || ''
                                                                            }
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
                                                              
                                                                // {
                                                                //     id: '4',
                                                                //     label: 'Frameworks',
                                                                //     content: (
                                                                //         <Benchmarks
                                                                //             benchmarks={
                                                                //                 selectedItem?.benchmarks
                                                                //             }
                                                                //         />
                                                                //     ),
                                                                // },
                        ]}
                    />
                </>
            ) : (
                <Spinner />
            )}
        </>
    )
}
