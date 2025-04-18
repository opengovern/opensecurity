import { useSetAtom } from 'jotai'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import axios from 'axios'
import { Badge, Card, Color, Flex } from '@tremor/react'
import { DocumentTextIcon } from '@heroicons/react/24/outline'
import {
    Cards,
    FormField,
    Grid,
    Input,
    KeyValuePairs,
    Link,
    Modal,
    Pagination,
    SpaceBetween,
    Table,
} from '@cloudscape-design/components'
import Spinner from '../../../components/Spinner'
import {
    BreadcrumbGroup,
    ExpandableSection,
} from '@cloudscape-design/components'
import ContentLayout from '@cloudscape-design/components/content-layout'
import Container from '@cloudscape-design/components/container'
import Header from '@cloudscape-design/components/header'
import Button from '@cloudscape-design/components/button'
import Box from '@cloudscape-design/components/box'
import {
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    Title,
    Switch,
} from '@tremor/react'
import CustomPagination from '../../../components/Pagination'
import { dateTimeDisplay } from '../../../utilities/dateDisplay'
import { notificationAtom } from '../../../store'
export interface Task {
    id:            string;
    name:          string;
    description:   string;
    image_url:     string;
    run_schedules: RunSchedule[];
    credentials:   string[];
    env_vars:      EnvVars;
    params:        string[];
    configs:       string[];
    scale_config:  ScaleConfig;
}

export interface EnvVars {
    CORE_BASEURL:                  string;
    ELASTICSEARCH_ADDRESS:         string;
    ELASTICSEARCH_ASSUME_ROLE_ARN: string;
    ELASTICSEARCH_AWS_REGION:      string;
    ELASTICSEARCH_IS_ON_AKS:       string;
    ELASTICSEARCH_IS_OPENSEARCH:   string;
    ELASTICSEARCH_PASSWORD:        string;
    ELASTICSEARCH_USERNAME:        string;
    NATS_CONSUMER:                 string;
    NATS_RESULT_TOPIC_NAME:        string;
    NATS_STREAM_NAME:              string;
    NATS_TOPIC_NAME:               string;
    NATS_URL:                      string;
}

export interface RunSchedule {
    last_run:  Date;
    params:    Params;
    frequency: number;
}

export interface Params {
    query_to_execute: string;
}

export interface ScaleConfig {
    stream:           string;
    consumer:         string;
    lag_threshold:    string;
    min_replica:      number;
    max_replica:      number;
    polling_interval: number;
    cooldown_period:  number;
}


export default function TaskDetail() {
    const { id } = useParams()
    const [loading, setLoading] = useState(false)
    const [runLoading, setRunLoading] = useState(false)
    const [task, setTask] = useState<Task>()
    const [page, setPage] = useState(1)
    const [total, setTotal] = useState(0)
    const [selected, setSelected] = useState<any>()
    const [results, setResults] = useState<any>([])
    const [detailOpen, setDetailOpen] = useState(false)
    const [scaleOpen, setScaleOpen] = useState(false)
    const [runOpen, setRunOpen] = useState(false)
    const [runParams, setRunParams] = useState<any>({})
    const [configOpen, setConfigOpen] = useState(false)
    const [config, setConfig] = useState<any>({})
    const setNotification = useSetAtom(notificationAtom)
    
    
    const getDetail = () => {
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
            .get(`${url}/main/tasks/api/v1/tasks/${id}`, config)
            .then((res) => {
                setLoading(false)
                setTask(res.data)
            })
            .catch((err) => {
                setLoading(false)
            })
    }
    const getRunResult = () => {
        setRunLoading(true)
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
            .get(
                `${url}/main/tasks/api/v1/tasks/${id}/runs?cursor=${page}&per_page=10`,
                config
            )
            .then((res) => {
                setRunLoading(false)
                if (res.data.items) {
                    setResults(res.data.items)
                }
                else{
                    setResults([])
                }
                setTotal(res.data.total_count)
                //  setTask(res.data)
            })
            .catch((err) => {
                setRunLoading(false)
            })
    }
     const getRunDetail = (id: string) => {
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
             .get(
                 `${url}/main/tasks/api/v1/tasks/run/${id}`,
                 config
             )
             .then((res) => {
                 setLoading(false)
                 if (res.data.items) {
                     setResults(res.data.items)
                 }
                 setTotal(res.data.total_count)
                 //  setTask(res.data)
             })
             .catch((err) => {
                 setLoading(false)
             })
     }

    const RunTask = () => {
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
        const body ={
            task_id:id,
            params: runParams

        }

        axios
            .post(`${url}/main/tasks/api/v1/tasks/run/`, body,config)
            .then((res) => {
                setLoading(false)
                //  setTask(res.data)
            })
            .catch((err) => {
                setLoading(false)
            })
    }
     const ConfigTask = () => {
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
         const body = {
             credentials: config,
         }

         axios
             .post(`${url}/main/tasks/api/v1/tasks/${id}/config`, body, config)
             .then((res) => {
                 setConfigOpen(false)
                 setNotification({
                        type: 'success',
                        text: 'Task Configured Successfully',
                    })
                 //  setTask(res.data)
             })
             .catch((err) => {
                 setLoading(false)
                    setNotification({
                            type: 'error',
                            text: 'Error Configuring Task',
                        })
             })
     }
    useEffect(() => {
        getRunResult()
    }, [page])
      useEffect(() => {
          getDetail()
      }, [])

    const truncate = (text: string | undefined) => {
        if (text) {
            return text.length > 600 ? text.substring(0, 600) + '...' : text
        }
    }
    return (
        <>
            {loading ? (
                <Spinner className="mt-56" />
            ) : (
                <>
                    <Modal
                        visible={detailOpen}
                        onDismiss={() => setDetailOpen(false)}
                        header="Job Detail"
                    >
                        <KeyValuePairs
                            className="mb-4"
                            columns={4}
                            items={[
                                { label: 'ID', value: selected?.id },
                                {
                                    label: 'Created At',
                                    value: `${
                                        selected?.created_at.split('T')[0]
                                    } ${
                                        selected?.created_at
                                            .split('T')[1]
                                            .split('.')[0]
                                    } `,
                                },
                                { label: 'Status', value: selected?.status },
                                {
                                    label: 'Updated At',
                                    value: `${
                                        selected?.updated_at.split('T')[0]
                                    } ${
                                        selected?.updated_at
                                            .split('T')[1]
                                            .split('.')[0]
                                    } `,
                                },
                                {
                                    label: 'Failure Reason',
                                    value: selected?.failure_message,
                                },
                            ]}
                        />
                        <Flex
                            flexDirection="col"
                            className="gap-4 justify-start items-start w-full"
                        >
                            <>
                                {selected?.params && (
                                    <>
                                        <h3 className="text-lg font-bold">
                                            Params:
                                        </h3>
                                        {/* iterate throuh params object and check if value is string show if array map through objects again */}
                                        {selected?.params &&
                                            Object.keys(selected?.params).map(
                                                (key) => {
                                                    return (
                                                        <>
                                                            <h4 className="text-md font-bold">
                                                                {key}:
                                                            </h4>
                                                            {typeof selected
                                                                ?.params[key] !=
                                                            'object' ? (
                                                                <Text>
                                                                    {
                                                                        selected
                                                                            ?.params[
                                                                            key
                                                                        ]
                                                                    }
                                                                </Text>
                                                            ) : (
                                                                <Flex
                                                                    flexDirection="col"
                                                                    className="gap-4 justify-start items-start ml-5 w-full flex-wrap"
                                                                >
                                                                    {selected?.params[
                                                                        key
                                                                    ].map(
                                                                        (
                                                                            item: any
                                                                        ) => {
                                                                            return (
                                                                                <Flex
                                                                                    flexDirection="col"
                                                                                    className="gap-4 justify-start items-start w-full"
                                                                                >
                                                                                    {Object.keys(
                                                                                        item
                                                                                    ).map(
                                                                                        (
                                                                                            key
                                                                                        ) => {
                                                                                            return (
                                                                                                <>
                                                                                                    <h5 className="text-md font-bold w-full">
                                                                                                        {
                                                                                                            key
                                                                                                        }

                                                                                                        :
                                                                                                    </h5>
                                                                                                    <Text className="w-full">
                                                                                                        {
                                                                                                            item[
                                                                                                                key
                                                                                                            ]
                                                                                                        }
                                                                                                    </Text>
                                                                                                </>
                                                                                            )
                                                                                        }
                                                                                    )}
                                                                                </Flex>
                                                                            )
                                                                        }
                                                                    )}
                                                                </Flex>
                                                            )}
                                                        </>
                                                    )
                                                }
                                            )}
                                    </>
                                )}
                            </>
                        </Flex>
                    </Modal>
                    <Modal
                        visible={scaleOpen}
                        size="large"
                        onDismiss={() => setScaleOpen(false)}
                        header="Task Detail"
                    >
                        <Flex className="flex-col justify-start items-start gap-2">
                            <KeyValuePairs
                                columns={4}
                                items={[
                                    { label: 'ID', value: task?.id },
                                    {
                                        label: 'Image Url',
                                        value: task?.image_url,
                                    },
                                    {
                                        label: 'Last Run',
                                        value: dateTimeDisplay(
                                            task?.run_schedules[0]?.last_run
                                        ),
                                    },
                                    {
                                        label: 'Frequency',
                                        value: task?.run_schedules[0]
                                            ?.frequency,
                                    },
                                ]}
                            />

                            <span className="text-base font-semibold">
                                Scale Config
                            </span>

                            <KeyValuePairs
                                columns={4}
                                items={[
                                    {
                                        label: 'Stream',
                                        value: task?.scale_config?.stream,
                                    },
                                    {
                                        label: 'Consumer',
                                        value: task?.scale_config?.consumer,
                                    },
                                    {
                                        label: 'Lag Threshold',
                                        value: task?.scale_config
                                            ?.lag_threshold,
                                    },
                                    {
                                        label: 'Min Replica',
                                        value: task?.scale_config?.min_replica,
                                    },
                                    {
                                        label: 'Max Replica',
                                        value: task?.scale_config?.max_replica,
                                    },
                                    {
                                        label: 'Polling Interval',
                                        value: task?.scale_config
                                            ?.polling_interval,
                                    },
                                    {
                                        label: 'Cooldown Period',
                                        value: task?.scale_config
                                            ?.cooldown_period,
                                    },
                                ]}
                            />
                            <Flex className='w-full gap-2 justify-end mt-2'>
                                <Button
                                    // className='bg-white'
                                    iconName="settings"

                                    onClick={() => {
                                        setScaleOpen(false)
                                        setConfigOpen(true)
                                    }}
                                >
                                    Configure Credentials
                                </Button>
                            </Flex>
                        </Flex>
                    </Modal>
                    <Modal
                        visible={runOpen}
                        onDismiss={() => setRunOpen(false)}
                        header="Run Task"
                    >
                        <Flex className="flex-col justify-start items-start gap-2">
                            <>
                                {task?.params?.map((param: any) => {
                                    return (
                                        <>
                                            <FormField
                                                label={`Value for ${param}`}
                                                className="w-full"
                                            >
                                                <Input
                                                    value={runParams[param]}
                                                    onChange={(e) => {
                                                        setRunParams({
                                                            ...runParams,
                                                            [param]:
                                                                e.detail.value,
                                                        })
                                                    }}
                                                    placeholder={`value for ${param}`}
                                                />
                                            </FormField>
                                        </>
                                    )
                                })}
                            </>

                            <Flex
                                className="w-full justify-end flex-row items-center
                            "
                            >
                                <Button
                                    onClick={() => {
                                        RunTask()
                                    }}
                                >
                                    Run
                                </Button>
                            </Flex>
                        </Flex>
                    </Modal>
                    <Modal
                        visible={configOpen}
                        onDismiss={() => setConfigOpen(false)}
                        header="Configure Task"
                    >
                        <Flex className="flex-col justify-start items-start gap-2">
                            <>
                                {task?.credentials?.map((param: any) => {
                                    return (
                                        <>
                                            <FormField
                                                label={`Value for ${param}`}
                                                className="w-full"
                                            >
                                                <Input
                                                    value={config[param]}
                                                    onChange={(e) => {
                                                        setConfig({
                                                            ...config,
                                                            [param]:
                                                                e.detail.value,
                                                        })
                                                    }}
                                                    placeholder={`value for ${param}`}
                                                />
                                            </FormField>
                                        </>
                                    )
                                })}
                            </>

                            <Flex
                                className="w-full justify-end flex-row items-center
                            "
                            >
                                <Button
                                    onClick={() => {
                                        ConfigTask()
                                    }}
                                >
                                    Save
                                </Button>
                            </Flex>
                        </Flex>
                    </Modal>

                    <Container
                        disableHeaderPaddings
                        disableContentPaddings
                        className="rounded-xl  bg-[#0f2940] p-0 text-white mt-4"
                        header={
                            <Header
                                className={`bg-[#0f2940] p-4 pt-0 rounded-xl   text-white ${
                                    false ? 'rounded-b-none' : ''
                                }`}
                                variant="h2"
                                description=""
                            >
                                <SpaceBetween size="xxxs" direction="vertical">
                                    <Box className="rounded-xl same text-white pt-3 pl-3 pb-0">
                                        <Grid
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
                                                    className="text-white important"
                                                >
                                                    <span className="text-white">
                                                        {task?.name}
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
                                                                task?.description
                                                            )}
                                                        </Text>
                                                        <Card className="absolute w-full text-wrap z-40 top-0 scale-0 transition-all p-2 group-hover:scale-100">
                                                            <Text>
                                                                {
                                                                    task?.description
                                                                }
                                                            </Text>
                                                        </Card>
                                                    </div>
                                                </Box>
                                            </div>
                                        </Grid>
                                    </Box>
                                    <Flex className="w-max gap-2 pl-3">
                                        <Button
                                            variant="primary"
                                            onClick={() => {
                                                setRunOpen(true)
                                            }}
                                        >
                                            Run
                                        </Button>
                                        <Button
                                            variant="primary"
                                            onClick={() => {
                                                setScaleOpen(true)
                                            }}
                                        >
                                            Detail
                                        </Button>
                                    </Flex>
                                </SpaceBetween>
                            </Header>
                        }
                    ></Container>

                    <Table
                        className="mt-2"
                        onRowClick={(event) => {
                            const row = event.detail.item
                            setSelected(row)
                            setDetailOpen(true)
                            getRunDetail(row.id)
                        }}
                        columnDefinitions={[
                            {
                                id: 'id',
                                header: 'Id',
                                cell: (item: any) => <>{item.id}</>,
                                // sortingField: 'id',
                                isRowHeader: true,
                                maxWidth: 100,
                            },
                            {
                                id: 'createdAt',
                                header: 'Created At',
                                cell: (item) => (
                                    <>{`${item?.created_at.split('T')[0]} ${
                                        item?.created_at
                                            .split('T')[1]
                                            .split('.')[0]
                                    } `}</>
                                ),
                                sortingField: 'createdAt',
                                isRowHeader: true,
                                maxWidth: 100,
                            },

                            {
                                id: 'status',
                                header: 'Status',
                                cell: (item) => {
                                    let jobStatus = ''
                                    let jobColor: Color = 'gray'
                                    switch (item?.status) {
                                        case 'CREATED':
                                            jobStatus = 'created'
                                            break
                                        case 'QUEUED':
                                            jobStatus = 'queued'
                                            break
                                        case 'IN_PROGRESS':
                                            jobStatus = 'in progress'
                                            jobColor = 'orange'
                                            break
                                        case 'RUNNERS_IN_PROGRESS':
                                            jobStatus = 'in progress'
                                            jobColor = 'orange'
                                            break
                                        case 'SUMMARIZER_IN_PROGRESS':
                                            jobStatus = 'summarizing'
                                            jobColor = 'orange'
                                            break
                                        case 'OLD_RESOURCE_DELETION':
                                            jobStatus = 'summarizing'
                                            jobColor = 'orange'
                                            break
                                        case 'SUCCEEDED':
                                            jobStatus = 'succeeded'
                                            jobColor = 'emerald'
                                            break
                                        case 'FINISHED':
                                            jobStatus = 'finished'
                                            jobColor = 'emerald'
                                            break
                                        case 'COMPLETED':
                                            jobStatus = 'completed'
                                            jobColor = 'emerald'
                                            break
                                        case 'FAILED':
                                            jobStatus = 'failed'
                                            jobColor = 'red'
                                            break
                                        case 'COMPLETED_WITH_FAILURE':
                                            jobStatus = 'completed with failed'
                                            jobColor = 'red'
                                            break
                                        case 'TIMEOUT':
                                            jobStatus = 'time out'
                                            jobColor = 'red'
                                            break
                                        default:
                                            jobStatus = String(item?.status)
                                    }

                                    return (
                                        <Badge color={jobColor}>
                                            {jobStatus}
                                        </Badge>
                                    )
                                },
                                sortingField: 'status',
                                isRowHeader: true,
                                maxWidth: 50,
                            },
                            {
                                id: 'updatedAt',
                                header: 'Updated At',
                                cell: (item) => (
                                    <>{`${item?.updated_at.split('T')[0]} ${
                                        item?.updated_at
                                            .split('T')[1]
                                            .split('.')[0]
                                    } `}</>
                                ),
                                sortingField: 'updatedAt',
                                isRowHeader: true,
                                maxWidth: 100,
                            },
                        ]}
                        columnDisplay={[
                            { id: 'id', visible: true },
                            { id: 'title', visible: true },
                            { id: 'type', visible: false },
                            { id: 'status', visible: true },
                            { id: 'createdAt', visible: true },
                            { id: 'updatedAt', visible: true },
                        ]}
                        loading={runLoading}
                        // @ts-ignore
                        items={results ? results : []}
                        empty={
                            <Box
                                margin={{ vertical: 'xs' }}
                                textAlign="center"
                                color="inherit"
                            >
                                <SpaceBetween size="m">
                                    <b>No resources</b>
                                    {/* <Button>Create resource</Button> */}
                                </SpaceBetween>
                            </Box>
                        }
                        header={
                            <Header
                                actions={
                                    <>
                                        <Button onClick={getRunResult}>
                                            Reload
                                        </Button>
                                    </>
                                }
                                className="w-full"
                            >
                                Results {total ?? 0}
                            </Header>
                        }
                        pagination={
                            <CustomPagination
                                currentPageIndex={page}
                                pagesCount={Math.ceil(total / 15)}
                                onChange={({ detail }: any) =>
                                    setPage(detail.currentPageIndex)
                                }
                            />
                        }
                    />
                </>
            )}
        </>
    )
}
