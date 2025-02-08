import {
    Card,
    Flex,
    Grid,
    List,
    ListItem,
    Metric,
    Switch,
    Text,
    Title,
    Tab,
    TabGroup,
    TabList,
    Button,
    TextInput,
    Divider,
} from '@tremor/react'
import { useParams } from 'react-router-dom'
import { useAtom, useSetAtom } from 'jotai'
import { useEffect, useState } from 'react'

import {
    useWorkspaceApiV3LoadSampleData,
    useWorkspaceApiV3PurgeSampleData,
} from '../../../api/metadata.gen'
import Spinner from '../../../components/Spinner'
import { numericDisplay } from '../../../utilities/numericDisplay'
import { useAuthApiV1UserDetail } from '../../../api/auth.gen'
import { dateDisplay, dateTimeDisplay } from '../../../utilities/dateDisplay'
import { PlatformEnginePkgWorkspaceApiTier } from '../../../api/api'
import { isDemoAtom, previewAtom, sampleAtom } from '../../../store'
import {
    useMetadataApiV1MetadataCreate,
    useMetadataApiV1MetadataDetail,
    useWorkspaceApiV1WorkspaceCurrentList,
} from '../../../api/metadata.gen'
import { useComplianceApiV1QueriesSyncList } from '../../../api/compliance.gen'
import { getErrorMessage } from '../../../types/apierror'
import { ConvertToBoolean } from '../../../utilities/bool'
import axios from 'axios'
import {
    Alert,
    KeyValuePairs,
    Modal,
    ProgressBar,
} from '@cloudscape-design/components'
import SettingsCustomization from '../Jobs/Customization'
import ReactMarkdown from 'react-markdown'
import rehypeRaw from 'rehype-raw'
interface ITextMetric {
    title: string
    metricId: string
    disabled?: boolean
}

function TextMetric({ title, metricId, disabled }: ITextMetric) {
    const [value, setValue] = useState<string>('')
    const [timer, setTimer] = useState<any>()

    const {
        response,
        isLoading,
        isExecuted,
        sendNow: refresh,
    } = useMetadataApiV1MetadataDetail(metricId)

    const {
        isLoading: setIsLoading,
        isExecuted: setIsExecuted,
        error,
        sendNow: sendSet,
    } = useMetadataApiV1MetadataCreate(
        {
            key: metricId,
            value,
        },
        {},
        false
    )

    useEffect(() => {
        if (isExecuted && !isLoading) {
            setValue(response?.value || '')
        }
    }, [isLoading])

    useEffect(() => {
        if (setIsExecuted && !setIsLoading) {
            refresh()
        }
    }, [setIsLoading])

    useEffect(() => {
        if (value === '' || value === response?.value) {
            return
        }

        if (timer !== undefined && timer !== null) {
            clearTimeout(timer)
        }

        const t = setTimeout(() => {
            sendSet()
        }, 1500)

        setTimer(t)
    }, [value])

    return (
        <Flex flexDirection="row" className="mb-4 sm:flex-row gap-4 justify-start flex-col">
            <Flex justifyContent="start" className="truncate space-x-4 w-fit min-w-max ">
                <div className="truncate">
                    <Text className="truncate text-sm">{title}:</Text>
                </div>
            </Flex>

            <TextInput
                value={value}
                onValueChange={(e) => setValue(String(e))}
                error={error !== undefined}
                errorMessage={getErrorMessage(error)}
                icon={isLoading ? Spinner : undefined}
                disabled={isLoading || disabled}
            />
        </Flex>
    )
}
export default function SettingsEntitlement() {
    
    const { response: currentWorkspace, isLoading: loadingCurrentWS } =
        useWorkspaceApiV1WorkspaceCurrentList()

    const [sample, setSample] = useAtom(sampleAtom)
    const [open,setOpen] = useState(false)
   
   
    const [preview, setPreview] = useAtom(previewAtom)
    const [migrations_status, setMigrationsStatus] = useState()
    const {
        response: customizationEnabled,
        isLoading: loadingCustomizationEnabled,
    } = useMetadataApiV1MetadataDetail('customization_enabled')
    const isCustomizationEnabled =
        ConvertToBoolean(
            (customizationEnabled?.value || 'false').toLowerCase()
        ) || false

    const wsTier = (v?: PlatformEnginePkgWorkspaceApiTier) => {
        switch (v) {
            // case PlatformEnginePkgWorkspaceApiTier.TierEnterprise:
            //     return 'Enterprise'
            default:
                return 'Community'
        }
    }
    const wsDetails = [
        {
            title: 'Community Version',
            // @ts-ignore
            value: currentWorkspace?.app_version,
        },
        {
            title: 'License',
            value: (
                <a
                    href="https://opencomply.io/license"
                    className="text-blue-600 underline"
                >
                    Business Source License v1.1
                </a>
            ),
        },
        {
            title: 'Privacy',
            value: (
                <a
                    onClick={()=>{
                        setOpen(true)
                    }}
                    href="#"
                    className="text-blue-600 underline"
                >
                    Usage Data Notice
                </a>
            ),
        },
        {
            title: 'Install date',
            value: dateTimeDisplay(
                // @ts-ignore
                currentWorkspace?.workspace_creation_time ||
                    Date.now().toString()
            ),
        },
        // {
        //     title: 'Edition',
        //     value: wsTier(currentWorkspace?.tier),
        // },
    ]
    const {
        isLoading: syncLoading,
        isExecuted: syncExecuted,
        error: syncError,
        sendNow: runSync,
    } = useComplianceApiV1QueriesSyncList({}, {}, false)

    const {
        isExecuted,
        isLoading: isLoadingLoad,
        error,
        sendNow: loadData,
    } = useWorkspaceApiV3LoadSampleData({}, false)
    const {
        isExecuted: isExecPurge,
        isLoading: isLoadingPurge,
        error: errorPurge,
        sendNow: PurgeData,
    } = useWorkspaceApiV3PurgeSampleData({}, false)

    const [status, setStatus] = useState()
    const [percentage, setPercentage] = useState()
    const [intervalId, setIntervalId] = useState()
    const [loaded, setLoaded] = useState()

    const GetStatus = () => {
        let url = ''
        //  setLoading(true)
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
            .get(`${url}/main/core/api/v3/migration/status `, config)
            .then((res) => {
                setStatus(res.data.status)
                setMigrationsStatus(res.data)
                setPercentage(res.data.Summary?.progress_percentage)
                if (intervalId) {
                    if (
                        res.data.status === 'SUCCEEDED' ||
                        res.data.status === 'FAILED'
                    ) {
                        clearInterval(intervalId)
                    }
                } else {
                    if (
                        res.data.status !== 'SUCCEEDED' &&
                        res.data.status !== 'FAILED'
                    ) {
                        if(!intervalId){
   const id = setInterval(GetStatus, 120000)
   // @ts-ignore
   setIntervalId(id)
                        } 
                     
                    }
                }
                //  const temp = []
                //  if (!res.data.items) {
                //      setLoading(false)
                //  }
                //  setBenchmarks(res.data.items)
                //  setTotalPage(Math.ceil(res.data.total_count / 6))
                //  setTotalCount(res.data.total_count)
            })
            .catch((err) => {
                clearInterval(intervalId)
                //  setLoading(false)
                //  setBenchmarks([])

                console.log(err)
            })
    }
    const GetSample = () => {
        let url = ''
        //  setLoading(true)
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
            .put(`${url}/main/core/api/v3/sample/loaded `, {}, config)
            .then((res) => {
                if (res.data === 'True') {
                    // @ts-ignore
                    setLoaded('True')
                    localStorage.setItem('sample', 'true')
                    setSample('true')
                } else {
                    GetSampleStatus()
                }
                //  const temp = []
                //  if (!res.data.items) {
                //      setLoading(false)
                //  }
                //  setBenchmarks(res.data.items)
                //  setTotalPage(Math.ceil(res.data.total_count / 6))
                //  setTotalCount(res.data.total_count)
            })
            .catch((err) => {
                //  setLoading(false)
                //  setBenchmarks([])

                console.log(err)
            })
    }
    const GetSampleStatus = () => {
        let url = ''
        //  setLoading(true)
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
            .get(`${url}/main/core/api/v3/sample/sync/status `, config)
            .then((res) => {
                if (res?.data?.progress !== 1 && res?.data?.progress !== 0) {
                    // @ts-ignore
                    setLoaded('True')
                    localStorage.setItem('sample', 'true')
                    setSample('true')
                } else {
                    // @ts-ignore
                    setLoaded('False')
                    localStorage.setItem('sample', 'false')
                    setSample('false')
                }
                //  const temp = []
                //  if (!res.data.items) {
                //      setLoading(false)
                //  }
                //  setBenchmarks(res.data.items)
                //  setTotalPage(Math.ceil(res.data.total_count / 6))
                //  setTotalCount(res.data.total_count)
            })
            .catch((err) => {
                //  setLoading(false)
                //  setBenchmarks([])

                console.log(err)
            })
    }
    useEffect(() => {
        GetStatus()
        GetSample()
    }, [])
    useEffect(() => {
        if (syncExecuted && !syncLoading) {
            GetStatus()
            const id = setInterval(GetStatus, 120000)
            // @ts-ignore
            setIntervalId(id)
            // setValue(response?.value || '')
            // window.location.reload()
        }
    }, [syncLoading, syncExecuted])
    return loadingCurrentWS ? (
        <Flex justifyContent="center" className="mt-56">
            <Spinner />
        </Flex>
    ) : (
        <div
            className="w-full"
            style={
                window.innerWidth < 768
                    ? { width: `${window.innerWidth - 80}px` }
                    : {}
            }
        >
            <Flex flexDirection="col" className="w-full">
                <Card key="summary" className=" w-full">
                    <Title className="font-semibold mb-2">Settings</Title>
                    <KeyValuePairs
                        columns={5}
                        // @ts-ignore
                        items={wsDetails.map((item) => {
                            return {
                                label: item.title,
                                value: item.value,
                            }
                        })}
                    />

                    <Divider />
                    <Flex className="flex-row gap-8 w-full justify-start items-start">
                        <Flex className="flex-col w-full justify-start items-start gap-4 max-w-[50%] border-r pr-8">
                            <Title className="font-semibold ">
                                Platform Configuration
                            </Title>
                            <Flex
                                justifyContent="start"
                                className="truncate space-x-4"
                            >
                                <div className="truncate">
                                    <Text className="truncate text-sm">
                                        Platform Controls, Frameworks, and
                                        Queries are sourced from Git
                                        repositories. Currently, only public
                                        repositories are supported.
                                    </Text>
                                </div>
                            </Flex>
                            <Flex
                                flexDirection="row"
                                className="mt-4 sm:flex-row flex-col sm:mb-0 mb-4"
                                alignItems="start"
                                justifyContent="start"
                            >
                                <TextMetric
                                    metricId="analytics_git_url"
                                    title="Configuration Git URL"
                                    disabled={
                                        loadingCustomizationEnabled ||
                                        isCustomizationEnabled === false
                                    }
                                />
                                <Button
                                    variant="secondary"
                                    className="ml-2"
                                    loading={syncExecuted && syncLoading}
                                    disabled={
                                        status !== 'SUCCEEDED' &&
                                        status !== 'FAILED'
                                    }
                                    onClick={() => runSync()}
                                >
                                    <Flex flexDirection="row" className="gap-2">
                                        {status !== 'SUCCEEDED' &&
                                            status !== 'FAILED' && (
                                                <Spinner className=" w-4 h-4" />
                                            )}
                                        {status === 'SUCCEEDED' ||
                                        status === 'FAILED'
                                            ? 'Re-Sync'
                                            : status}
                                    </Flex>
                                </Button>
                            </Flex>
                            {(status !== 'SUCCEEDED' ||
                                status !== 'FAILED') && (
                                <>
                                    <Flex className="w-full">
                                        <ProgressBar
                                            value={percentage}
                                            className="w-full"
                                            // additionalInfo="Additional information"
                                            // @ts-ignore
                                            description={`${status}, Last Updated: ${dateTimeDisplay(
                                                // @ts-ignore
                                                migrations_status?.updated_at
                                            )}`}
                                            resultText="Configuration done"
                                            label="Platform Configuration"
                                        />
                                    </Flex>
                                </>
                            )}
                        </Flex>
                        <Flex className="flex-col w-full justify-start items-start gap-4  ">
                            {' '}
                            <Title className="font-semibold ">
                                App configurations
                            </Title>
                            <Flex
                                flexDirection="row"
                                justifyContent="between"
                                className="w-full mt-4"
                            >
                                <Text className="font-normal">
                                    Show preview features
                                </Text>
                                <TabGroup
                                    index={preview === 'true' ? 0 : 1}
                                    onIndexChange={(idx) =>
                                        setPreview(idx === 0 ? 'true' : 'false')
                                    }
                                    className="w-fit"
                                >
                                    <TabList
                                        className="border border-gray-200"
                                        variant="solid"
                                    >
                                        <Tab>On</Tab>
                                        <Tab>Off</Tab>
                                    </TabList>
                                </TabGroup>
                            </Flex>
                            <SettingsCustomization />
                        </Flex>
                    </Flex>

                    <Divider />

                    <Title className="font-semibold mt-8">Sample Data</Title>
                    <Flex
                        justifyContent="between"
                        alignItems="center"
                        className="sm:flex-row flex-col"
                    >
                        <Text className="font-normal w-full">
                            {' '}
                            The app can be loaded with sample data, allowing you
                            to explore features without setting up integrations.
                        </Text>
                        <Flex
                            className="gap-2 sm:mt-0 mt-2"
                            justifyContent="end"
                            alignItems="center"
                        >
                            {loaded != 'True' && sample != 'true' && (
                                <Button
                                    variant="secondary"
                                    className="ml-2"
                                    loading={isLoadingLoad && isExecuted}
                                    onClick={() => {
                                        loadData()
                                        setSample('true')
                                        localStorage.setItem('sample', 'true')
                                        // @ts-ignore
                                        setLoaded('True')
                                        // window.location.reload()
                                    }}
                                >
                                    Load Sample Data
                                </Button>
                            )}

                            {loaded == 'True' && sample == 'true' && (
                                <>
                                    <Button
                                        variant="secondary"
                                        className=""
                                        loading={isLoadingPurge && isExecPurge}
                                        onClick={() => {
                                            PurgeData()
                                            setSample(false)
                                            localStorage.setItem(
                                                'sample',
                                                'false'
                                            )
                                            // @ts-ignore
                                            setLoaded('False')
                                            // window.location.reload()
                                        }}
                                    >
                                        Purge Sample Data
                                    </Button>
                                </>
                            )}
                        </Flex>
                    </Flex>
                    {((error && error !== '') ||
                        (errorPurge && errorPurge !== '')) && (
                        <>
                            <Alert className="mt-2" type="error">
                                <>
                                    {getErrorMessage(error)}
                                    {getErrorMessage(errorPurge)}
                                </>
                            </Alert>
                        </>
                    )}
                </Card>
            </Flex>
            <Modal
                visible={open}
                size="large"
                onDismiss={() => setOpen(false)}
                // header="Usage Data Notice"
            >
                <div className=" markdown-container">
                    <ReactMarkdown
                        // @ts-ignore
                        children={`# Usage Data Notice

*opencomply* Community Edition includes a built-in process to gather a minimal set of **anonymized** usage data once per day. This helps us understand how our community uses opencomply and guides future improvements. Specifically, the following information is sent to our analytics service at [https://stats.opencomply.io](https://stats.opencomply.io):

1. **Product Version**: Used to identify which releases are actively in use.  
2. **Number of Users**: Aggregated count of active users.  
3. **Plugins and Plugin Counts**: Identifies plugins being used
4. **Configured Hostname**: The hostname as configured in your environment.

## Our Commitment to Your Privacy
- The data we collect is strictly limited to the items above.  
- We do **not** collect or store any personally identifiable information.  
- **No confidential data is ever sent** from your environment.  
- All data is transmitted and stored in an anonymized format.  
- We only use this information to improve OpenComply and to better understand usage patterns and trends.

## Open Source Code
- **Client-Side Collection**: The relevant logic in opencomply that gathers and sends usage data can be found here:  
  [https://github.com/opengovern/opencomply/blob/main/jobs/checkup-job/job.go](https://github.com/opengovern/opencomply/blob/main/jobs/checkup-job/job.go)  
- **Server-Side Processing**: The server responsible for receiving and processing the data is here:  
  [https://github.com/opengovern/usage-tracker](https://github.com/opengovern/usage-tracker)

We encourage you to review this code to understand exactly what information is being sent, how it is transmitted, and how it is processed once received.

## How to Disable Usage Data Collection
If you wish to turn off this data collection process, you can simply modify and recompile opencomply with the relevant functionality disabled. The pertinent section can be found in the [job.go file](https://github.com/opengovern/opencomply/blob/main/jobs/checkup-job/job.go).

If you have any concerns or questions about this data collection—or need guidance on disabling it—please consult our documentation or reach out to us at [insert your support channel or email]. We value your input and appreciate your support in making opencomply better for everyone.`}
                        skipHtml={false}
                        className={'markdown-body'}
                        rehypePlugins={[rehypeRaw]}
                    />
                </div>
            
            </Modal>
        </div>
    )
}
