import { Flex, Title } from '@tremor/react'
import {
    useLocation,
    useNavigate,
    useParams,
    useSearchParams,
} from 'react-router-dom'
import { Cog8ToothIcon } from '@heroicons/react/24/outline'
import { useAtomValue } from 'jotai'

import {
    AppLayout,
    Badge,
    Box,
    Button,
    Header,
    KeyValuePairs,
    Modal,
    Pagination,
    SpaceBetween,
    Spinner,
    SplitPanel,
    Table,
    Tabs,
} from '@cloudscape-design/components'

import axios from 'axios'
import { useEffect, useState } from 'react'
import { Credentials, Integration, Schema } from '../types'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import {
    GetActions,
    GetTableColumns,
    GetTableColumnsDefintion,
    GetViewFields,
    RenderTableField,
} from '../utils'
import UpdateCredentials from './Update'

interface CredentialsListProps {
    name?: string
    integration_type?: string
    schema?: Schema
}

export default function CredentialsList({
    name,
    integration_type,
    schema,
}: CredentialsListProps) {
    const navigate = useNavigate()
    const [row, setRow] = useState<Credentials[]>([])

    const [loading, setLoading] = useState<boolean>(false)
    const [actionLoading, setActionLoading] = useState<any>({
        update: false,
        delete: false,
        health_check: false,
    })

    const [error, setError] = useState<string>('')
    const [total_count, setTotalCount] = useState<number>(0)
    const [selectedItem, setSelectedItem] = useState<Credentials>()
    const [page, setPage] = useState(0)
    const [open, setOpen] = useState(false)
    const [edit, setEdit] = useState(false)
    const [openInfo, setOpenInfo] = useState(false)
    const [confirmModal, setConfirmModal] = useState(false)
    const [action, setAction] = useState()

    const GetCredentials = () => {
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
                `${url}/main/integration/api/v1/credentials/list`,
                body,
                config
            )
            .then((res) => {
                const data = res.data

                setTotalCount(data.total_count)
                setRow(data.credentials)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }
    const CheckActionsClick = (action: any) => {
        setAction(action)
        if (action.type === 'update') {
            setEdit(true)
        } else if (action.type === 'delete') {
            if (action?.confirm?.message && action?.confirm?.message !== '') {
                setConfirmModal(true)
            } else {
                CheckActionsSumbit(action)
            }
        }
    }
    const CheckActionsSumbit = (action: any) => {
        if (action?.type === 'update') {
            setEdit(true)
        } else if (action?.type === 'delete') {
            DeleteIntegration()
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
                `${url}/main/integration/api/v1/credentials/${selectedItem?.id}`,

                config
            )
            .then((res) => {
                GetCredentials()
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
    const GetMarkedSecrets = () => {
        // iterate over key value pair
        // @ts-ignore
        if (selectedItem?.masked_secret) {
            return Object.keys(selectedItem?.masked_secret)?.map((key) => {
                return {
                    label: key,
                    value: selectedItem?.masked_secret[key],
                }
            })
        }
        return []
    }

    useEffect(() => {
        GetCredentials()
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
                                    selectedItem?.id ? selectedItem?.id : ''
                                }
                            >
                                <KeyValuePairs
                                    columns={
                                        // @ts-ignore
                                        GetViewFields(schema, 1)?.length > 4
                                            ? 4
                                            : GetViewFields(schema, 1)?.length
                                    }
                                    // @ts-ignore
                                    items={GetViewFields(schema, 1)?.map(
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
                                <KeyValuePairs
                                    className='mt-5'
                                    columns={
                                        GetMarkedSecrets().length > 4
                                            ? 4
                                            : GetMarkedSecrets().length
                                    }
                                    items={GetMarkedSecrets()}
                                />
                                <Flex
                                    className="mt-5 gap-2 "
                                    justifyContent="end"
                                    flexDirection="row"
                                    alignItems="end"
                                >
                                    <>
                                        {GetActions(1, schema)?.map(
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
                                    1,
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
                                    1,
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
                                filter={
                                    ''
                                    // <PropertyFilter
                                    //     // @ts-ignore
                                    //     query={undefined}
                                    //     // @ts-ignore
                                    //     onChange={({ detail }) => {
                                    //         // @ts-ignore
                                    //         setQueries(detail)
                                    //     }}
                                    //     // countText="5 matches"
                                    //     enableTokenGroups
                                    //     expandToViewport
                                    //     filteringAriaLabel="Control Categories"
                                    //     // @ts-ignore
                                    //     // filteringOptions={filters}
                                    //     filteringPlaceholder="Control Categories"
                                    //     // @ts-ignore
                                    //     filteringOptions={undefined}
                                    //     // @ts-ignore

                                    //     filteringProperties={undefined}
                                    //     // filteringProperties={
                                    //     //     filterOption
                                    //     // }
                                    // />
                                }
                                header={
                                    <Header
                                        actions={
                                            <Flex className="gap-1">
                                                {/* <Button
                                                      // icon={PlusIcon}
                                                      onClick={() =>
                                                          setOpen(true)
                                                      }
                                                  >
                                                      Add New {`${name}`}{' '}
                                                      Integration
                                                  </Button> */}
                                                {/* <Button
                                            // icon={PencilIcon}
                                            onClick={() => setEdit(true)}
                                        >
                                            Edit Integration
                                        </Button> */}
                                                <Button
                                                    // icon={PencilIcon}
                                                    onClick={() => {
                                                        GetCredentials()
                                                    }}
                                                >
                                                    Reload
                                                </Button>
                                            </Flex>
                                        }
                                        className="w-full"
                                    >
                                        {name} Credentials{' '}
                                        <span className=" font-medium">
                                            ({total_count})
                                        </span>
                                    </Header>
                                }
                                pagination={
                                    <Pagination
                                        currentPageIndex={page + 1}
                                        pagesCount={Math.ceil(total_count / 10)}
                                        onChange={({ detail }) =>
                                            setPage(detail.currentPageIndex - 1)
                                        }
                                    />
                                }
                            />
                        }
                    />
                    <UpdateCredentials
                        name={name}
                        integration_type={integration_type}
                        schema={schema}
                        open={edit}
                        onClose={() => setEdit(false)}
                        GetList={GetCredentials}
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
