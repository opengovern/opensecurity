import { Card, Flex, Title, Text } from '@tremor/react'
import {
    useLocation,
    useNavigate,
    useParams,
    useSearchParams,
} from 'react-router-dom'
import {
    ArrowLeftStartOnRectangleIcon,
    Cog8ToothIcon,
} from '@heroicons/react/24/outline'
import { useAtomValue, useSetAtom } from 'jotai'

import {
    useIntegrationApiV1ConnectorsMetricsList,
    useIntegrationApiV1CredentialsList,
} from '../../../api/integration.gen'
import TopHeader from '../../../components/Layout/Header'
import {
    defaultTime,
    searchAtom,
    useUrlDateRangeState,
} from '../../../utilities/urlstate'
import axios from 'axios'
import { useEffect, useState } from 'react'
import { Integration, Schema } from './types'
import {
    BreadcrumbGroup,
    Button,
    Checkbox,
    FormField,
    Input,
    KeyValuePairs,
    Modal,
    Multiselect,
    Spinner,
    Tabs,
} from '@cloudscape-design/components'

import IntegrationList from './Integration'
import CredentialsList from './Credentials'
import { OpenGovernance } from '../../../icons/icons'
import DiscoveryJobs from './Discovery'
import Resources from './Resources'
import Setup from './Setup'
import ButtonDropdown from '@cloudscape-design/components/button-dropdown'
import { notificationAtom } from '../../../store'
import CreateIntegration from './Integration/Create'

export default function TypeDetail() {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const { type } = useParams()
    const [manifest, setManifest] = useState<any>()
    const { state } = useLocation()
    const [shcema, setSchema] = useState<Schema>()
    const [loading, setLoading] = useState<boolean>(false)
    const [status, setStatus] = useState<string>()
    const [row, setRow] = useState<Integration[]>([])
    const setNotification = useSetAtom(notificationAtom)
    const [actionLoading, setActionLoading] = useState<any>({
        discovery: false,
    })
    const [resourceTypes, setResourceTypes] = useState<any>([])
    const [selectedResourceType, setSelectedResourceType] = useState<any>()
    const [runOpen, setRunOpen] = useState(false)
    const [selectedIntegrations, setSelectedIntegrations] = useState<any>([])
    const [params, setParams] = useState<any>()
    const [enableSchedule, setEnableSchedule] = useState(false)
     const [open, setOpen] = useState(false)
    const GetSchema = () => {
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
                `${url}/main/integration/api/v1/integrations/types/${type}/ui/spec `,
                config
            )
            .then((res) => {
                const data = res.data
                setSchema(data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }
    const GetManifest = () => {
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
                `${url}/main/integration/api/v1/integration-types/plugin/${type}/manifest`,
                config
            )
            .then((resp) => {
                setManifest(resp.data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)

                // params.fail()
            })
    }
    const GetStatus = () => {
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
            .post(
                `${url}/main/integration/api/v1/integration-types/plugin/${type}/healthcheck`,
                {},
                config
            )
            .then((resp) => {
                setStatus(resp.data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)

                // params.fail()
            })
    }
    const UpdatePlugin = () => {
        setLoading(true)
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        let path = ''
        path = `/main/integration/api/v1/integration-types/plugin/load/id/${type}`

        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }

        axios
            .post(`${url}${path}`, {}, config)
            .then((res) => {
                setLoading(false)

                setNotification({
                    text: `Plugin Updated`,
                    type: 'success',
                })
                navigate('/integration/plugins')
            })
            .catch((err) => {
                setNotification({
                    text: `Error: ${err.response.data.message}`,
                    type: 'error',
                })

                setLoading(false)
            })
    }
    const UnInstallPlugin = () => {
        setLoading(true)
        let url = ''
        if (window.location.origin === 'http://localhost:3000') {
            url = window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
        } else {
            url = window.location.origin
        }
        let path = ''
        path = `/main/integration/api/v1/integration-types/plugin/uninstall/id/${type}`
        // @ts-ignore
        const token = JSON.parse(localStorage.getItem('openg_auth')).token

        const config = {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }

        axios
            .delete(`${url}${path}`, config)
            .then((res) => {
                setLoading(false)
                setNotification({
                    text: `Plugin Uninstalled`,
                    type: 'success',
                })
                navigate('/integration/plugins')
            })
            .catch((err) => {
                setNotification({
                    text: `Error: ${err.response.data.message}`,
                    type: 'error',
                })
                setLoading(false)
            })
    }

    const DisablePlugin = () => {
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
            .post(
                `${url}/main/integration/api/v1/integration-types/plugin/${type}/disable`,
                {},
                config
            )
            .then((res) => {
                setLoading(false)
                setNotification({
                    text: `Plugin Disabled`,
                    type: 'success',
                })
                navigate('/integration/plugins')
            })
            .catch((err) => {
                setLoading(false)
                setNotification({
                    text: `Error: ${err.response.data.message}`,
                    type: 'error',
                })
            })
    }
    const RunDiscovery = () => {
        setActionLoading({ ...actionLoading, discovery: true })
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
        let body = {}
        body = {
            integration_info: selectedIntegrations?.map((item: any) => {
                return {
                    integration_type: type,
                    provider_id: item.provider_id,
                    integration_id: item.integration_id,
                    name: item.name,
                }
            }),
        }
        if (
            selectedResourceType?.length > 0 &&
            selectedResourceType?.length < resourceTypes?.length
        ) {
            // @ts-ignore
            body['resource_types'] = selectedResourceType?.map((item: any) => {
                if (selectedResourceType?.length == 1) {
                    if (selectedResourceType[0]?.params?.length > 0) {
                        if (params) {
                            // @ts-ignore
                            return {
                                resource_type: item.value,
                                parameters: params,
                                enable_schedule: enableSchedule,
                            }
                        }
                    }
                }
                return {
                    resource_type: item.value,
                    enable_schedule: enableSchedule,
                }
            })
        }

        axios
            .post(`${url}/main/schedule/api/v3/discovery/run`, body, config)
            .then((res) => {
                GetIntegrations()
                setActionLoading({
                    ...actionLoading,
                    discovery: false,
                })
                setRunOpen(false)
                setNotification({
                    text: `Discovery started`,
                    type: 'success',
                })
                setParams({})
            })
            .catch((err) => {
                console.log(err)
                setActionLoading({
                    ...actionLoading,
                    discovery: false,
                })
                setNotification({
                    text: `Error: ${err.response.data.message}`,
                    type: 'error',
                })
            })
    }
    const GetResourceTypes = () => {
        setActionLoading({ ...actionLoading, discovery: true })
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

        // const body = {
        //     integration_type: [integration_type],
        // }
        axios
            .get(
                `${url}/main/integration/api/v1/integrations/types/${type}/resource_types`,

                config
            )
            .then((res) => {
                const data = res.data
                setResourceTypes(data?.integration_types)
                const temp: any = []
                data?.integration_types?.map((item: any) => {
                    temp.push({
                        label: item?.name,
                        value: item?.name,
                        params: item?.params,
                    })
                })
                setSelectedResourceType(temp)
                setActionLoading({ ...actionLoading, discovery: false })
            })
            .catch((err) => {
                console.log(err)
                setActionLoading({ ...actionLoading, discovery: false })
            })
    }
    const GetIntegrations = () => {
        setActionLoading({ ...actionLoading, discovery: true })
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
            integration_type: [type],
        }
        axios
            .post(
                `${url}/main/integration/api/v1/integrations/list`,
                body,
                config
            )
            .then((res) => {
                const data = res.data

                if (data.integrations) {
                    setRow(data.integrations)
                } else {
                    setRow([])
                }
                setActionLoading({ ...actionLoading, discovery: false })
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
                setActionLoading({ ...actionLoading, discovery: false })
            })
    }
    useEffect(() => {
        GetSchema()
        GetStatus()
        GetManifest()
    }, [])
    const GetItems = () => {
        const temp = []
        temp.push({
            label: 'Id',
            value: manifest?.IntegrationType,
        })
        if (window.innerWidth > 640) {
            temp.push(
                {
                    label: 'Version',
                    value: manifest?.DescriberTag,
                },
                {
                    label: 'Artifact URL',
                    value: manifest?.DescriberURL,
                }
            )
        }
        temp.push({
            label: 'Publisher',
            value: manifest?.Publisher,
        })
        temp.push({
            label: 'Operational Status',
            // @ts-ignore
            value: status
                ? status?.charAt(0).toUpperCase() + status?.slice(1)
                : '',
        })
        if (window.innerWidth > 640) {
            temp.push(
                {
                    label: 'Update date',
                    value: manifest?.UpdateDate,
                },
                {
                    label: 'Supported Platform Version',
                    value: manifest?.SupportedPlatformVersion,
                },

                {
                    label: 'Author',
                    value: manifest?.Author,
                }
            )
        }

        return temp
    }

    return (
        <>
            {/* <TopHeader breadCrumb={[state?.name]} /> */}

            {shcema && !loading && shcema?.integration_type_id ? (
                <>
                    <div
                        className="w-full flex justify-center items-center"
                        style={
                            window.innerWidth < 768
                                ? { width: `${window.innerWidth - 80}px` }
                                : {}
                        }
                    >
                        <Flex className="flex-col w-full justify-start items-start gap-4">
                            <BreadcrumbGroup
                                className="w-full"
                                items={[
                                    {
                                        text: 'Plugins',
                                        href: '/integration/plugins',
                                    },
                                    {
                                        // @ts-ignore
                                        text: state?.name,
                                        href: `/integration/plugins/${type}`,
                                    },
                                ]}
                            />
                            <Flex className="flex-col gap-3 justify-start items-start w-full">
                                <Flex className="sm:flex-row flex-col sm:justify-between justify-start sm:items-center items-start w-full sm:gap-8 gap-2">
                                    <h1 className=" font-bold text-2xl mb-2  text-left ml-1">
                                        {state?.name} plugin
                                    </h1>
                                    <ButtonDropdown
                                        onItemClick={({ detail }) => {
                                            const id = detail.id
                                            switch (id) {
                                                case 'update':
                                                    UpdatePlugin()
                                                    break
                                                case 'disable':
                                                    DisablePlugin()
                                                    break
                                                case 'uninstall':
                                                    UnInstallPlugin()
                                                    break
                                                case 'healthckeck':
                                                    GetStatus()
                                                    break
                                                case 'discovery':
                                                        GetIntegrations()
                                                        GetResourceTypes()
                                                        setRunOpen(true)
                                                    break
                                                default:
                                                    break
                                            }
                                        }}
                                        variant="primary"
                                        items={[
                                            {
                                                text: 'Run Discovery',
                                                id: 'discovery',
                                            },
                                            {
                                                text: 'Settings',
                                                items: [
                                                    {
                                                        text: 'Update',
                                                        id: 'update',
                                                    },
                                                    {
                                                        text: 'Disable',
                                                        id: 'disable',
                                                    },
                                                    {
                                                        text: 'Uninstall',
                                                        id: 'uninstall',
                                                    },
                                                ],
                                            },
                                            {
                                                text: 'Run Health Check',
                                                id: 'healthckeck',
                                            },
                                        ]}
                                        mainAction={{
                                            text: 'Add new integration',
                                            
                                            onClick: () => {
                                                setOpen(true)
                                            },
                                            loading: actionLoading['discovery'],
                                        }}
                                    ></ButtonDropdown>
                                </Flex>
                                <div
                                    className="w-full"
                                    style={
                                        window.innerWidth < 768
                                            ? {
                                                  width: `${
                                                      window.innerWidth - 80
                                                  }px`,
                                              }
                                            : {}
                                    }
                                >
                                    <Card className="w-full">
                                        <>
                                            <KeyValuePairs
                                                columns={4}
                                                items={GetItems()}
                                            />
                                        </>
                                    </Card>
                                </div>

                                <></>
                            </Flex>
                            <Tabs
                                tabs={[
                                    {
                                        id: '3',
                                        label: ' Discovered Resources',
                                        content: (
                                            <Resources
                                                name={state?.name}
                                                integration_type={type}
                                            />
                                        ),
                                    },
                                    {
                                        id: '0',
                                        label: 'Integrations',
                                        content: (
                                            <IntegrationList
                                                schema={shcema}
                                                name={state?.name}
                                                integration_type={type}
                                            />
                                        ),
                                    },
                                    {
                                        id: '1',
                                        label: 'Credentials',
                                        content: (
                                            <CredentialsList
                                                schema={shcema}
                                                name={state?.name}
                                                integration_type={type}
                                            />
                                        ),
                                    },
                                    {
                                        id: '2',
                                        label: 'Discovery Jobs',
                                        content: (
                                            <DiscoveryJobs
                                                name={state?.name}
                                                integration_type={type}
                                            />
                                        ),
                                    },

                                    {
                                        id: '4',
                                        label: 'Setup Guide',
                                        content: (
                                            <Setup
                                                name={state?.name}
                                                integration_type={type}
                                            />
                                        ),
                                    },
                                ]}
                            />
                        </Flex>
                        <Modal
                            visible={runOpen}
                            onDismiss={() => {
                                setRunOpen(false)
                                setSelectedIntegrations([])
                                setSelectedResourceType([])
                                setParams({})
                            }}
                            // @ts-ignore
                            header={'Run Discovery'}
                            footer={
                                <Flex className="gap-3" justifyContent="end">
                                    <Button
                                        onClick={() => {
                                            setRunOpen(false)
                                            setSelectedIntegrations([])
                                            setSelectedResourceType([])
                                            setParams({})
                                        }}
                                    >
                                        Cancel
                                    </Button>
                                    <Button
                                        onClick={() => {
                                            if (
                                                selectedResourceType?.length ==
                                                resourceTypes?.length
                                            ) {
                                                setSelectedResourceType([])
                                                return
                                            }
                                            const temp: any = []
                                            resourceTypes?.map((item: any) => {
                                                temp.push({
                                                    label: item?.name,
                                                    value: item?.name,
                                                    params: item?.params,
                                                })
                                            })
                                            setSelectedResourceType(temp)
                                        }}
                                    >
                                        {selectedResourceType?.length ==
                                        resourceTypes?.length
                                            ? 'Unselect all types'
                                            : 'Select all types'}
                                    </Button>
                                    <Button
                                        variant="primary"
                                        loading={actionLoading['discovery']}
                                        onClick={() => {
                                            RunDiscovery()
                                        }}
                                    >
                                        Confirm
                                    </Button>
                                </Flex>
                            }
                        >
                            <Flex
                                className="gap-5 w-full justify-start items-start"
                                flexDirection="col"
                            >
                                <Multiselect
                                    className="w-full"
                                    options={row?.map((item: any) => {
                                        return {
                                            label: item?.name,
                                            value: item?.name,
                                            provider_id: item.provider_id,
                                            integration_id: item.integration_id,
                                            name: item.name,
                                        }
                                    })}
                                    selectedOptions={selectedIntegrations}
                                    onChange={({ detail }) => {
                                        setSelectedIntegrations(
                                            detail.selectedOptions
                                        )
                                    }}
                                    tokenLimit={5}
                                    placeholder="Select Integration"
                                />
                                <Multiselect
                                    className="w-full"
                                    options={resourceTypes?.map((item: any) => {
                                        return {
                                            label: item?.name,
                                            value: item?.name,
                                            params: item?.params,
                                        }
                                    })}
                                    selectedOptions={selectedResourceType}
                                    onChange={({ detail }) => {
                                        setSelectedResourceType(
                                            detail.selectedOptions
                                        )
                                    }}
                                    tokenLimit={0}
                                    placeholder="Select resource type"
                                />
                                <Checkbox
                                    onChange={({ detail }) =>
                                        setEnableSchedule(detail.checked)
                                    }
                                    checked={enableSchedule}
                                >
                                    Make this a recurring Discovery Job
                                </Checkbox>
                                {selectedResourceType?.length == 1 && (
                                    <>
                                        {/* show params to input */}
                                        {selectedResourceType[0]?.params?.map(
                                            (item: any) => {
                                                return (
                                                    <FormField
                                                        className="w-full"
                                                        label={`${item.name} (Optional)`}
                                                        description={
                                                            item.description
                                                        }
                                                    >
                                                        <Input
                                                            className="w-full"
                                                            value={
                                                                params?.[
                                                                    item.name
                                                                ]
                                                            }
                                                            type={'text'}
                                                            onChange={({
                                                                detail,
                                                            }) =>
                                                                setParams({
                                                                    ...params,
                                                                    [item.name]:
                                                                        detail.value,
                                                                })
                                                            }
                                                        />
                                                    </FormField>
                                                )
                                            }
                                        )}
                                    </>
                                )}
                            </Flex>
                        </Modal>
                        <CreateIntegration
                            name={state?.name}
                            integration_type={type}
                            schema={shcema}
                            open={open}
                            onClose={() => setOpen(false)}
                            GetList={() => {
                                window.location.reload()
                            }}
                        />
                    </div>
                </>
            ) : (
                <>
                    {loading ? (
                        <>
                            <Spinner />
                        </>
                    ) : (
                        <>
                            <Flex
                                flexDirection="col"
                                className="fixed top-0 left-0 w-screen h-screen bg-gray-900/80 z-50"
                            >
                                <Card className="w-1/3 mt-56">
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="center"
                                        alignItems="center"
                                    >
                                        <OpenGovernance className="w-14 h-14 mb-6" />
                                        <Title className="mb-3 text-2xl font-bold">
                                            Data not found
                                        </Title>
                                        <Text className="mb-6 text-center">
                                            Json schema not found for this
                                            integration
                                        </Text>
                                        <Button
                                            // icon={ArrowLeftStartOnRectangleIcon}
                                            onClick={() => {
                                                navigate('/integration/plugins')
                                            }}
                                        >
                                            Back
                                        </Button>
                                    </Flex>
                                </Card>
                            </Flex>
                        </>
                    )}
                </>
            )}
        </>
    )
}
