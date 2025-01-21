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
    ProgressBar,
} from '@cloudscape-design/components'

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
        <Flex flexDirection="row" className="mb-4">
            <Flex justifyContent="start" className="truncate space-x-4 ">
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
            title: 'Version',
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
            title: 'Install date',
            value: dateTimeDisplay(
                // @ts-ignore
                currentWorkspace?.workspace_creation_time ||
                    Date.now().toString()
            ),
        },
        {
            title: 'Edition',
            value: wsTier(currentWorkspace?.tier),
        },
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
    return   loadingCurrentWS ? (
        <Flex justifyContent="center" className="mt-56">
            <Spinner />
        </Flex>
    ) : (
        <Flex flexDirection="col">
            
            <Card key="summary" className=" w-full">
                <Title className="font-semibold mb-2">Settings</Title>
                <KeyValuePairs
                    columns={4}
                    items={wsDetails.map((item) => {
                        return {
                            label: item.title,
                            value: item.value,
                        }
                    })}
                />
              
                <Divider />
                <Title className="font-semibold mt-8">
                    Platform Configuration
                </Title>
                <Flex justifyContent="start" className="truncate space-x-4">
                    <div className="truncate">
                        <Text className="truncate text-sm">
                            Platform Controls, Frameworks, and Queries are
                            sourced from Git repositories. Currently, only
                            public repositories are supported.
                        </Text>
                    </div>
                </Flex>
                <Flex
                    flexDirection="row"
                    className="mt-4"
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
                        disabled={status !== 'SUCCEEDED' && status !== 'FAILED'}
                        onClick={() => runSync()}
                    >
                        <Flex flexDirection="row" className="gap-2">
                            {status !== 'SUCCEEDED' && status !== 'FAILED' && (
                                <Spinner className=" w-4 h-4" />
                            )}
                            {status === 'SUCCEEDED' || status === 'FAILED'
                                ? 'Re-Sync'
                                : status}
                        </Flex>
                    </Button>
                </Flex>
                {(status !== 'SUCCEEDED' || status !== 'FAILED') && (
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
                <Divider />

                <Title className="font-semibold mt-8">App configurations</Title>

              
                <Flex
                    flexDirection="row"
                    justifyContent="between"
                    className="w-full mt-4"
                >
                    <Text className="font-normal">Show preview features</Text>
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
                <Divider />

                <Title className="font-semibold mt-8">Sample Data</Title>
                <Flex justifyContent="between" alignItems="center">
                    <Text className="font-normal w-full">
                        {' '}
                        The app can be loaded with sample data, allowing you to
                        explore features without setting up integrations.
                    </Text>
                    <Flex
                        className="gap-2"
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
                                        localStorage.setItem('sample', 'false')
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
    )
}
