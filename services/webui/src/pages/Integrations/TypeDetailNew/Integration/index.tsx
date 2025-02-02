import { Flex, Title } from '@tremor/react'
import {
    useLocation,
    useNavigate,
    useParams,
    useSearchParams,
} from 'react-router-dom'
import { Cog8ToothIcon } from '@heroicons/react/24/outline'
import { useAtomValue, useSetAtom } from 'jotai'

import axios from 'axios'
import { useEffect, useState } from 'react'
import { Integration, Schema } from '../types'

import {
    AppLayout,
    Badge,
    Box,
    Button,
    Checkbox,
    FormField,
    Header,
    Input,
    KeyValuePairs,
    Modal,
    Multiselect,
    Pagination,
    SpaceBetween,
    Spinner,
    SplitPanel,
    Table,
    Tabs,
} from '@cloudscape-design/components'
import CreateIntegration from './Create'
import { Label } from '@headlessui/react/dist/components/label/label'
import { GetActions, GetDetailsActions, GetTableColumns, GetTableColumnsDefintion, GetViewFields, RenderTableField } from '../utils'
import { update } from '@react-spring/web'
import UpdateIntegration from './Update'
import { notificationAtom } from '../../../../store'
import CustomPagination from '../../../../components/Pagination'

interface IntegrationListProps {
    name?: string
    integration_type?: string
    schema?: Schema
}

const states = {
    ACTIVE: 'green',
    INACTIVE: 'red',
    ARCHIVED: 'grey',
}
export default function IntegrationList({
    name,
    integration_type,
    schema,
}: IntegrationListProps) {
    const navigate = useNavigate()
    const [row, setRow] = useState<Integration[]>([])

    const [loading, setLoading] = useState<boolean>(false)
    const [actionLoading, setActionLoading] = useState<any>({
        update: false,
        delete: false,
        health_check: false,
        discovery: false,
    })

    const [error, setError] = useState<string>('')
    const [total_count, setTotalCount] = useState<number>(0)
    const [selectedItem, setSelectedItem] = useState<Integration>()
    const [page, setPage] = useState(0)
    const [open, setOpen] = useState(false)
    const [edit, setEdit] = useState(false)
    const [openInfo, setOpenInfo] = useState(false)
    const [confirmModal, setConfirmModal] = useState(false)
    const [action, setAction] = useState()
    const setNotification = useSetAtom(notificationAtom)
    const [resourceTypes, setResourceTypes] = useState<any>([])
    const [selectedResourceType, setSelectedResourceType] = useState<any>()
    const [runOpen, setRunOpen] = useState(false)
    const [selectedIntegrations, setSelectedIntegrations] = useState<any>([])
    const [params, setParams] = useState<any>()
    const [enableSchedule, setEnableSchedule] = useState(false)

    const GetIntegrations = () => {
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
            integration_type: [integration_type],
        }
        axios
            .post(
                `${url}/main/integration/api/v1/integrations/list`,
                body,
                config
            )
            .then((res) => {
                const data = res.data

                setTotalCount(data.total_count)
                if(data.integrations){
                setRow(data.integrations)

                }
                else{
                    setRow([])
                }
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }
   
    const CheckActionsClick = (action: any) => {
        setAction(action)
        if (action.type === "update") {
            setEdit(true)
        } else if (action.type === 'delete') {
            if (action?.confirm?.message && action?.confirm?.message !== '') {
                setConfirmModal(true)
            } else {
                CheckActionsSumbit(action)
            }
        } else if (action.type == 'health_check') {
            CheckActionsSumbit(action)
        }
    }
    const CheckActionsSumbit = (action: any) => {
        if (action?.type === "update") {
            setEdit(true)
        } else if (action?.type === 'delete') {
            DeleteIntegration()
        } else if (action?.type === 'health_check') {
            HealthCheck()
        }
    }

    const DeleteIntegration = () => {
        setActionLoading({ ...actionLoading, delete: true })

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
            .delete(
                `${url}/main/integration/api/v1/integrations/${selectedItem?.integration_id}`,

                config
            )
            .then((res) => {
                GetIntegrations()
                setConfirmModal(false)
                setOpenInfo(false)
                setActionLoading({ ...actionLoading, delete: false })
            })
            .catch((err) => {
                console.log(err)
                setActionLoading({
                    ...actionLoading,
                    delete: false,
                })
            })
    }
    const HealthCheck = () => {
        setActionLoading({ ...actionLoading, health_check: true })
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
            .put(
                `${url}/main/integration/api/v1/integrations/${selectedItem?.integration_id}/healthcheck`,
                {},
                config
            )
            .then((res) => {
                GetIntegrations()
                setSelectedItem(res.data)
                setActionLoading({
                    ...actionLoading,
                    health_check: false,
                })

                setConfirmModal(false)
            })
            .catch((err) => {
                console.log(err)
                setActionLoading({
                    ...actionLoading,
                    health_check: false,
                })
            })
    }

   


    
    useEffect(() => {
        GetIntegrations()
    }, [])

    return (
        <>
            {schema ? (
                <>
                    <AppLayout
                        toolsOpen={false}
                        navigationOpen={false}
                        contentType="table"
                        toolsHide={true}
                        navigationHide={true}
                        splitPanelOpen={openInfo}
                        onSplitPanelToggle={() => {
                            setOpenInfo(!openInfo)
                        }}
                        splitPanel={
                            <SplitPanel
                                // @ts-ignore
                                header={
                                    selectedItem?.name ? selectedItem?.name : ''
                                }
                            >
                                <KeyValuePairs
                                    columns={
                                        // @ts-ignore
                                        GetViewFields(schema, 0)?.length > 4
                                            ? 4
                                            : GetViewFields(schema, 0)?.length
                                    }
                                    // @ts-ignore
                                    items={GetViewFields(schema, 0)?.map(
                                        (field) => {
                                            return {
                                                label: field.title,
                                                value: selectedItem
                                                    ? RenderTableField(
                                                          field,
                                                          selectedItem
                                                      )
                                                    : '',
                                            }
                                        }
                                    )}
                                />
                                <Flex
                                    className="mt-5 gap-2 "
                                    justifyContent="end"
                                    flexDirection="row"
                                    alignItems="end"
                                >
                                    <>
                                        {GetActions(0, schema)?.map(
                                            (action) => {
                                                if (action.type !== 'view') {
                                                    return (
                                                        <>
                                                            <Button
                                                                loading={
                                                                    actionLoading[
                                                                        action
                                                                            .type
                                                                    ]
                                                                }
                                                                onClick={() => {
                                                                    CheckActionsClick(
                                                                        action
                                                                    )
                                                                }}
                                                            >
                                                                {action.label}
                                                            </Button>
                                                        </>
                                                    )
                                                }
                                            }
                                        )}
                                    </>
                                </Flex>
                            </SplitPanel>
                        }
                        content={
                            <Table
                                className="  min-h-[450px]"
                                variant="full-page"
                                // resizableColumns
                                renderAriaLive={({
                                    firstIndex,
                                    lastIndex,
                                    totalItemsCount,
                                }) =>
                                    `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                                }
                                onSortingChange={(event) => {
                                    // setSort(event.detail.sortingColumn.sortingField)
                                    // setSortOrder(!sortOrder)
                                }}
                                // sortingColumn={sort}
                                // sortingDescending={sortOrder}
                                // sortingDescending={sortOrder == 'desc' ? true : false}
                                // @ts-ignore
                                onRowClick={(event) => {
                                    const row = event.detail.item
                                    // @ts-ignore
                                    setSelectedItem(row)
                                    setOpenInfo(true)
                                }}
                                // @ts-ignore

                                columnDefinitions={GetTableColumns(
                                    0,
                                    schema
                                )?.map((field) => {
                                    return {
                                        id: field.key,
                                        header: field.title,
                                        // @ts-ignore
                                        cell: (item) => (
                                            <>{RenderTableField(field, item)}</>
                                        ),
                                        // sortingField: 'providerConnectionID',
                                        isRowHeader: true,
                                        maxWidth: 100,
                                    }
                                })}
                                columnDisplay={GetTableColumnsDefintion(
                                    0,
                                    schema
                                )}
                                enableKeyboardNavigation
                                // @ts-ignore
                                items={row?.length > 0 ? row : []}
                                loading={loading}
                                loadingText="Loading resources"
                                // stickyColumns={{ first: 0, last: 1 }}
                                // stripedRows
                                trackBy="id"
                                empty={
                                    <Box
                                        margin={{ vertical: 'xs' }}
                                        textAlign="center"
                                        color="inherit"
                                    >
                                        <SpaceBetween size="m">
                                            <b>No resources</b>
                                        </SpaceBetween>
                                    </Box>
                                }
                                header={
                                    <Header
                                        actions={
                                            <Flex className="gap-1">
                                                <Button
                                                    // icon={PencilIcon}
                                                    onClick={() => {
                                                        GetIntegrations()
                                                    }}
                                                    iconName="refresh"
                                                ></Button>
                                                {/* <Button
                                                    // icon={PlusIcon}
                                                    // variant='primary'
                                                    onClick={() =>
                                                        setOpen(true)
                                                    }
                                                >
                                                    Add
                                                </Button> */}
                                            </Flex>
                                        }
                                        className="w-full"
                                    >
                                         Integrations{' '}
                                        <span className=" font-medium">
                                            ({total_count})
                                        </span>
                                    </Header>
                                }
                                pagination={
                                    <CustomPagination
                                        currentPageIndex={page + 1}
                                        pagesCount={Math.ceil(total_count / 10)}
                                        onChange={({ detail }: any) =>
                                            setPage(detail.currentPageIndex - 1)
                                        }
                                    />
                                }
                            />
                        }
                    />
                    
                    <UpdateIntegration
                        name={name}
                        integration_type={integration_type}
                        schema={schema}
                        open={edit}
                        onClose={() => setEdit(false)}
                        GetList={GetIntegrations}
                        selectedItem={selectedItem}
                    />
                    <Modal
                        visible={confirmModal}
                        onDismiss={() => setConfirmModal(false)}
                        // @ts-ignore
                        header={
                            // @ts-ignore

                            action?.label
                                ? // @ts-ignore
                                  action.label + ' ' + selectedItem?.name
                                : ''
                        }
                    >
                        <Box className="p-3">
                            {/* @ts-ignore */}
                            <Title>{action?.confirm?.message}</Title>
                            <Flex className="gap-2 mt-5" justifyContent="end">
                                <Button
                                    onClick={() => {
                                        setConfirmModal(false)
                                    }}
                                >
                                    Cancel
                                </Button>
                                <Button
                                    variant="primary"
                                    onClick={() => {
                                        CheckActionsSumbit(action)
                                    }}
                                >
                                    Confirm
                                </Button>
                            </Flex>
                        </Box>
                    </Modal>
                </>
            ) : (
                <Spinner />
            )}
        </>
    )
}
