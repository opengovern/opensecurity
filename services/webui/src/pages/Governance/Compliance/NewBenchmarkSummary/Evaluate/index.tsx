import { useAtomValue, useSetAtom } from 'jotai'
import { useCallback, useEffect, useMemo, useState } from 'react'
import {
    Button,
    Flex,
    MultiSelect,
    MultiSelectItem,
    Title,
} from '@tremor/react'
import {
    ArrowPathRoundedSquareIcon,
    CloudIcon,
    PlayCircleIcon,
} from '@heroicons/react/24/outline'
import { useComplianceApiV1AssignmentsBenchmarkDetail } from '../../../../../api/compliance.gen'
import {
    PlatformEnginePkgComplianceApiBenchmarkAssignedConnection,
    PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary,
} from '../../../../../api/api'
import { isDemoAtom, notificationAtom } from '../../../../../store'
import {
    Box,
    Checkbox,
    Icon,
    Multiselect,
    Select,
    SpaceBetween,
} from '@cloudscape-design/components'
import { Fragment, ReactNode } from 'react'
import { Dialog, Transition } from '@headlessui/react'
import Modal from '@cloudscape-design/components/modal'
import KButton from '@cloudscape-design/components/button'
import axios from 'axios'
interface IEvaluate {
    benchmarkDetail:
        | PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary
        | undefined
    onEvaluate: (c: string[],b: boolean) => void
    open: boolean 
    setOpen: Function
    showBenchmark: boolean
}

export default function Evaluate({
    benchmarkDetail,
    onEvaluate,
    open,
    setOpen,
    showBenchmark,
}: IEvaluate) {
    const [loading, setLoading] = useState(false)
    const [accounts, setAccounts] = useState()
    const isDemo = useAtomValue(isDemoAtom)
    const [openConfirm, setOpenConfirm] = useState(false)
    const [connections, setConnections] = useState<string[]>([])
    const [benchmarks, setBenchmarks] = useState<any[]>([])
    const [selectedbenchmarks, setSelectedBenchmarks] = useState<any[]>()
    const setNotification = useSetAtom(notificationAtom)
    const [withIncidents, setWithIncidents] = useState(true)

  
    const RunBenchmark = () => {
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        const body = {
            with_incidents: withIncidents,
            // with_incidents: true,

            integration_info: connections.map((c) => {
                return {
                    // @ts-ignore
                    integration_id: c.value,
                }
            }),
            // @ts-ignore
        }
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }
        //    console.log(config)
        axios
            .post(
                // @ts-ignore
                `${url}/main/schedule/api/v3/compliance/benchmark/${selectedbenchmarks?.value}/run`,
                body,
                config
            )
            .then((res) => {
                let ids = ''
                // @ts-ignore
                res.data.jobs.map((item, index) => {
                    if (index < 5) {
                        ids = ids + item.job_id + ','
                    }
                })
                setNotification({
                    text: `Run is Done You Job id is ${ids}`,
                    type: 'success',
                })
            })
            .catch((err) => {
                console.log(err)
                const message = err.response.data.message
                setNotification({
                    text: `${message}`,
                    type: 'error',
                })
            })
    }

    const GetEnabled = () => {
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
        const integrations: any = []
        if (showBenchmark) {
            // @ts-ignore
            selectedbenchmarks?.connectors?.map((c) => {
                integrations.push(c)
            })
        } else {
            // @ts-ignore
            benchmarkDetail?.integrationTypes?.map((c) => {
                integrations.push(c)
            })
        }
        const body = {
            integration_type: integrations,
        }
        axios
            .post(
                `${url}/main/integration/api/v1/integrations/list`,
                body,
                config
            )
            .then((res) => {
                setAccounts(res.data.integrations)
                setLoading(false)
            })
            .catch((err) => {
                setLoading(false)
                console.log(err)
            })
    }

    const isLoading =
        benchmarkDetail?.lastJobStatus !== 'FAILED' &&
        benchmarkDetail?.lastJobStatus !== 'SUCCEEDED' &&
        (benchmarkDetail?.lastJobStatus || '') !== ''
    const GetCard = () => {
        let url = ''
        setLoading(true)
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
            sort_by: 'incidents',
            assigned: false,
            is_baseline: false,
            connectors: [],
            root: true,
        }

        axios
            .post(`${url}/main/compliance/api/v3/benchmarks`, body, config)
            .then((res) => {
                //  const temp = []
                if (!res.data.items) {
                    setLoading(false)
                }
                setBenchmarks(res.data.items)
            })
            .catch((err) => {
                setLoading(false)
                setBenchmarks([])

                console.log(err)
            })
    }

    useEffect(() => {
        if (showBenchmark) {
            GetCard()
        } 
    }, [showBenchmark])
    useEffect(() => {
        if (selectedbenchmarks) {
            GetEnabled()
        }
    }, [selectedbenchmarks])
    useEffect(() => {
        if(!showBenchmark){
            GetEnabled()
        }
    }, [showBenchmark])

    return (
        <>
            <Modal
                onDismiss={() => {
                    setOpen(false)
                }}
                visible={open}
                footer={
                    <Box float="right">
                        <SpaceBetween direction="horizontal" size="xs">
                            <KButton onClick={() => setOpen(false)}>
                                Close
                            </KButton>
                            {/* <Button
                                variant="secondary"
                                onClick={() => {
                                    setConnections(
                                       []
                                    )
                                }}
                            >
                                DeSelect All
                            </Button> */}
                            <KButton
                                onClick={() => {
                                    if (showBenchmark) {
                                        if (!selectedbenchmarks) {
                                            setNotification({
                                                text: `Please Select Frameworks`,
                                                type: 'error',
                                            })
                                            return
                                        }
                                    }
                                    setConnections(
                                        // @ts-ignore
                                        accounts?.map((c) => {
                                            return {
                                                label: c.name,
                                                value: c.platform_name,
                                                description: c.integration_id,
                                            }
                                        })
                                    )
                                    // if (opened) {
                                    //     setSelectedBenchmarks(
                                    //         benchmarks?.map((c) => {
                                    //             return {
                                    //                 label: c?.benchmark?.title,
                                    //                 value: c?.benchmark?.id,
                                    //                 // description: c.id,
                                    //             }
                                    //         })
                                    //     )
                                    // }
                                }}
                            >
                                Select All
                            </KButton>
                            <KButton
                                variant={'primary'}
                                onClick={() => {
                                    setOpen(false)
                                    if (showBenchmark) {
                                        RunBenchmark()
                                    } else {
                                        onEvaluate(connections, withIncidents)
                                    }
                                }}
                            >
                                Run
                            </KButton>
                        </SpaceBetween>
                    </Box>
                }
                header={
                    showBenchmark
                        ? 'Select Compliance Framework for Audit'
                        : 'Select Integrations'
                }
            >
                {showBenchmark && (
                    <>
                        <Select
                            className="w-full mb-2"
                            // @ts-ignore
                            options={benchmarks?.map((c) => {
                                return {
                                    label: c?.benchmark?.title,
                                    value: c?.benchmark?.id,
                                    description: c?.benchmark?.id,
                                    connectors: c?.benchmark?.connectors,
                                }
                            })}
                            // @ts-ignore
                            selectedOption={selectedbenchmarks}
                            loadingText="Loading Frameworks"
                            // @ts-ignore
                            emptyText="No Frameworks"
                            loading={loading}
                            // tokenLimit={0}
                            filteringType="auto"
                            placeholder="Select Frameworks"
                            onChange={({ detail }) => {
                                // @ts-ignore
                                setSelectedBenchmarks(detail.selectedOption)
                                // @ts-ignore
                                setAccounts([])
                                setConnections([])
                            }}
                        />
                    </>
                )}
                <Multiselect
                    className="w-full"
                    // @ts-ignore
                    options={[
                        {
                            label: 'ACTIVE',
                            // @ts-ignore
                            options:
                                accounts
                                    // @ts-ignore
                                    ?.filter(
                                        // @ts-ignore

                                        (c) => c?.state != 'INACTIVE'
                                    )
                                    // @ts-ignore
                                    ?.map((c) => {
                                        return {
                                            label: c.name,
                                            value: c.integration_id,
                                            description: c.provider_id,
                                            labelTag: 'ACTIVE',
                                        }
                                    }) || [],
                        },
                        {
                            label: 'INACTIVE',
                            // @ts-ignore
                            options:
                                accounts
                                    // @ts-ignore
                                    ?.filter(
                                        // @ts-ignore

                                        (c) => c?.state == 'INACTIVE'
                                    )
                                    // @ts-ignore
                                    ?.map((c) => {
                                        return {
                                            label: c.name,
                                            value: c.integration_id,
                                            description: c.integration_id,
                                            labelTag: 'INACTIVE',
                                        }
                                    }) || [],
                        },
                    ]}
                    // @ts-ignore
                    selectedOptions={connections}
                    loadingText="Loading Integrations"
                    emptyText="No Integrations"
                    loading={loading}
                    tokenLimit={1}
                    filteringType="auto"
                    placeholder="Select Integrations"
                    onChange={({ detail }) => {
                        // @ts-ignore
                        setConnections(detail.selectedOptions)
                    }}
                />
                <Checkbox
                    className="mt-2 w-full"
                    onChange={({ detail }) => setWithIncidents(detail.checked)}
                    disabled={true}
                    checked={withIncidents}
                >
                    Create auditable incident
                </Checkbox>
            </Modal>
        </>
    )
}

