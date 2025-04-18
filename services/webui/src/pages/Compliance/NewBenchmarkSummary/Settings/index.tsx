// @ts-nocheck

import { useAtomValue } from 'jotai'
import {
    Button,
    Callout,
    Divider,
    Flex,
    Grid,
    Switch,
    Text,
} from '@tremor/react'
import { useEffect, useState } from 'react'
import { Cog6ToothIcon } from '@heroicons/react/24/outline'
import { isDemoAtom } from '../../../../store'

import {
    useComplianceApiV1AssignmentsBenchmarkDetail,
    useComplianceApiV1BenchmarksSettingsCreate,
} from '../../../../api/compliance.gen'

import KTable from '@cloudscape-design/components/table'
import KButton from '@cloudscape-design/components/button'

import Box from '@cloudscape-design/components/box'
import SpaceBetween from '@cloudscape-design/components/space-between'
import {
    FormField,
    Modal,
    Multiselect,
    RadioGroup,
    Tiles,
    Toggle,
} from '@cloudscape-design/components'
import axios from 'axios'
import {
    BreadcrumbGroup,
    Header,
    Link,
    Pagination,
    PropertyFilter,
} from '@cloudscape-design/components'
import CustomPagination from '../../../../components/Pagination'
interface ISettings {
    id: string | undefined
    response: (x: number) => void
    autoAssign: boolean | undefined
    reload: () => void
}

interface ITransferState {
    connectionID: string
    status: boolean
}

export default function Settings({
    id,
    response,
    autoAssign,
    reload,
}: ISettings) {
    const [firstLoading, setFirstLoading] = useState<boolean>(true)

    const [allEnable, setAllEnable] = useState(autoAssign)
    const [banner, setBanner] = useState(autoAssign)
    const isDemo = useAtomValue(isDemoAtom)
    const [loading, setLoading] = useState(false)
    const [rows, setRows] = useState<any>([])
    const [page, setPage] = useState(1)
    const [totalPages, setTotalPages] = useState(0)
    const [available, setAvailable] = useState([])
    const [availableLoading, setAvailableLoading] = useState(false)
    const [openAdd, setOpenAdd] = useState(false)
    const [openDelete, setOpenDelete] = useState(false)
    const [selected, setSelected] = useState<any>([])
    const [selectedDelete, setSelectedDelete] = useState<any>([])
    const [addLoading, setAddLoading] = useState(false)
    const [deleteLoading, setDeleteLoading] = useState(false)
    const [err, setError] = useState('')

    useComplianceApiV1AssignmentsBenchmarkDetail(String(id), {}, false)

    const {
        isLoading: changeSettingsLoading,
        isExecuted: changeSettingsExecuted,
        sendNowWithParams: changeSettings,
    } = useComplianceApiV1BenchmarksSettingsCreate(String(id), {}, {}, false)

    useEffect(() => {
        if (!changeSettingsLoading) {
            reload()
        }
    }, [changeSettingsLoading])

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
        axios
            .get(
                `${url}/main/compliance/api/v1/frameworks/${id}/assignments?page_size=10&page=${page}`,
                config
            )
            .then((res) => {
                setRows(res.data?.data)
                setTotalPages(res.data?.page_info.total_items)
                setLoading(false)
            })
            .catch((err) => {
                setLoading(false)

                console.log(err)
            })
    }
    const GetAvailable = () => {
        setAvailableLoading(true)
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
                `${url}/main/compliance/api/v1/frameworks/${id}/assignments/available`,
                config
            )
            .then((res) => {
                if (res?.data?.items) {
                    const temp = res?.data?.items?.map((item) => {
                        return {
                            label: item?.name,
                            value: item?.integration_id,
                            description: item?.provider_id,
                        }
                    })
                    setAvailable(temp)
                } else {
                    setAvailable([])
                }
                setAvailableLoading(false)
            })
            .catch((err) => {
                setAvailableLoading(false)
                setAvailable([])
                setError(err?.response?.data?.message)
                console.log(err)
            })
            .catch((err) => {
                setAvailableLoading(false)
                setAvailable([])
                console.log(err)
            })
    }
    const Add = () => {
        setAddLoading(true)
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
            integrations: selected.map((item) => {
                return item?.value
            }),
        }
        axios
            .put(
                `${url}/main/compliance/api/v1/frameworks/${id}/assignments`,
                body,
                config
            )
            .then((res) => {
                setAddLoading(false)
                setOpenAdd(false)
                GetEnabled()
                setSelected([])
            })
            .catch((err) => {
                setAddLoading(false)

                console.log(err)
            })
    }
    const Delete = () => {
        setDeleteLoading(true)
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
                `${url}/main/compliance/api/v1/frameworks/${id}/assignments/${selectedDelete.integration_id}`,
                config
            )
            .then((res) => {
                setDeleteLoading(false)
                setOpenDelete(false)
                setSelectedDelete()
                GetEnabled()
            })
            .catch((err) => {
                setDeleteLoading(false)

                console.log(err)
            })
    }

    useEffect(() => {
        GetEnabled()
    }, [page])

    return (
        <>
            <div
                className="w-full"
                style={
                    window.innerWidth < 768
                        ? { width: `${window.innerWidth - 80}px` }
                        : {}
                }
            >
                <KTable
                    className="   min-h-[450px]"
                    // resizableColumns
                    // variant="full-page"
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
                        // console.log(event)
                        // const row = event.detail.item
                    }}
                    columnDefinitions={[
                        {
                            id: 'id',
                            header: 'Id',
                            cell: (item) => item?.integration_id,
                            sortingField: 'id',
                            isRowHeader: true,
                            // maxWidth
                        },
                        {
                            id: 'id_name',
                            header: 'Name',
                            cell: (item) => item?.integration_name,
                        },
                        {
                            id: 'provider_id',
                            header: 'Provider ID',
                            cell: (item) => item?.integration_provider_id,
                        },
                        {
                            id: 'integration_type',
                            header: 'Integration Type',
                            cell: (item) => item?.plugin_id,
                        },
                        {
                            id: 'type',
                            header: ' Type',
                            cell: (item) => item?.assignment_type,
                        },
                        {
                            id: 'enable',
                            header: '',
                            cell: (item) => (
                                <>
                                    <KButton
                                        // variant="icon"
                                        iconName="remove"
                                        variant="inline-icon"
                                        onClick={() => {
                                            setSelectedDelete(item)
                                            setOpenDelete(true)
                                        }}
                                    />
                                </>
                            ),
                        },
                    ]}
                    columnDisplay={[
                        { id: 'id', visible: true },
                        { id: 'name', visible: true },
                        { id: 'provider_id', visible: true },
                        { id: 'integration_type', visible: true },
                        { id: 'type', visible: true },

                        { id: 'enable', visible: true },
                    ]}
                    enableKeyboardNavigation
                    // @ts-ignore
                    items={rows ? rows : []}
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
                            className="w-full"
                            actions={
                                <Flex className="gap-2">
                                    <KButton
                                        onClick={() => {
                                            GetAvailable()
                                            setOpenAdd(true)
                                        }}
                                    >
                                        Add
                                    </KButton>
                                </Flex>
                            }
                        >
                            Assigments{' '}
                            <span className=" font-medium">({totalPages})</span>
                        </Header>
                    }
                    pagination={
                        <CustomPagination
                            currentPageIndex={page}
                            pagesCount={Math.ceil(totalPages / 10)}
                            onChange={({ detail }) =>
                                setPage(detail.currentPageIndex)
                            }
                        />
                    }
                />
            </div>
            <Modal
                visible={openAdd}
                onDismiss={() => {
                    setOpenAdd(false)
                    setSelected([])
                }}
                header="Add Assignments"
                footer={
                    <Box float="right">
                        <SpaceBetween direction="horizontal" size="xs">
                            <KButton
                                onClick={() => {
                                    setOpenAdd(false)
                                    setSelected([])
                                }}
                            >
                                Close
                            </KButton>

                            {selected.length == available.length ? (
                                <>
                                    <KButton
                                        onClick={() => {
                                            setSelected([])
                                        }}
                                    >
                                        UnSelect All
                                    </KButton>
                                </>
                            ) : (
                                <>
                                    <KButton
                                        onClick={() => {
                                            setSelected(
                                                available?.map((item) => {
                                                    return {
                                                        label: item?.name,
                                                        value: item?.integration_id,
                                                        description:
                                                            item?.provider_id,
                                                    }
                                                })
                                            )
                                        }}
                                    >
                                        Select All
                                    </KButton>
                                </>
                            )}

                            <KButton
                                variant={'primary'}
                                loading={addLoading}
                                onClick={() => {
                                    Add()
                                }}
                            >
                                Add
                            </KButton>
                        </SpaceBetween>
                    </Box>
                }
            >
                <Multiselect
                    className="w-full"
                    options={available}
                    selectedOptions={selected}
                    loadingText="Loading Assignment"
                    emptyText="No assignment"
                    empty={err && err != '' ? err : 'No Assignment'}
                    loading={false}
                    tokenLimit={1}
                    // filteringType="auto"
                    placeholder="Select Assignment"
                    onChange={({ detail }) => {
                        setSelected(detail.selectedOptions)
                    }}
                />
            </Modal>
            <Modal
                visible={openDelete}
                onDismiss={() => {
                    setOpenDelete(false)
                    setSelectedDelete()
                }}
                header="Delete Assignments"
                footer={
                    <Box float="right">
                        <SpaceBetween direction="horizontal" size="xs">
                            <KButton
                                onClick={() => {
                                    setOpenDelete(false)
                                    setSelectedDelete()
                                }}
                            >
                                Close
                            </KButton>

                            <KButton
                                variant={'primary'}
                                loading={deleteLoading}
                                onClick={() => {
                                    Delete()
                                }}
                            >
                                Delete
                            </KButton>
                        </SpaceBetween>
                    </Box>
                }
            >
                Are you sure you want to delete assignment{' '}
                {selectedDelete?.integration_name}?
            </Modal>
        </>
    )
}
