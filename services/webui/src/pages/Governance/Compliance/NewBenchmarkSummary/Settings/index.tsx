// @ts-noCheck
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
import { isDemoAtom } from '../../../../../store'
import DrawerPanel from '../../../../../components/DrawerPanel'
import Table, { IColumn } from '../../../../../components/Table'
import {
    useComplianceApiV1AssignmentsBenchmarkDetail,
    useComplianceApiV1BenchmarksSettingsCreate,
} from '../../../../../api/compliance.gen'
import Spinner from '../../../../../components/Spinner'
import KTable from '@cloudscape-design/components/table'
import KButton from '@cloudscape-design/components/button'

import Box from '@cloudscape-design/components/box'
import SpaceBetween from '@cloudscape-design/components/space-between'
import {
    FormField,
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
interface ISettings {
    id: string | undefined
    response: (x: number) => void
    autoAssign: boolean | undefined
    tracksDriftEvents: boolean | undefined
    isAutoResponse: (x: boolean) => void
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
    tracksDriftEvents,
    isAutoResponse,
    reload,
}: ISettings) {
    const [firstLoading, setFirstLoading] = useState<boolean>(true)
  
    const [allEnable, setAllEnable] = useState(autoAssign)
    const [banner, setBanner] = useState(autoAssign)
    const isDemo = useAtomValue(isDemoAtom)
    const [loading, setLoading] = useState(false)
    const [rows,setRows] = useState<any>([])
       const [page, setPage] = useState(0)
 

   

    

    // const {
    //     response: assignments,
    //     isLoading,
    //     sendNow: refreshList,
    // } = useComplianceApiV1AssignmentsBenchmarkDetail(String(id), {}, false)

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

    // useEffect(() => {
    //     if (id && !assignments) {
    //         refreshList()
    //     }
    //     if (assignments) {
    //         const count = assignments.connections?.filter((c) => c.status)
    //         response(count?.length || 0)
    //     }
    // }, [id, assignments])



 
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
               `${url}/main/compliance/api/v3/benchmark/${id}/assignments`,
               config
           )
           .then((res) => {
            setRows(res.data.items)
       setLoading(false)
              
           })
           .catch((err) => {
       setLoading(false)

               console.log(err)
           })
   }
   const ChangeStatus = (status: string) => {
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
           auto_enable: status == 'auto-enable' ? true : false,
           disable: status == 'disabled' ? true : false,
       }
       axios
           .post(
               `${url}/main/compliance/api/v3/benchmark/${id}/assign`,body,
               config
           )
           .then((res) => {
                window.location.reload()
            //    setLoading(false)

           })
           .catch((err) => {
               setLoading(false)

               console.log(err)
           })
   }
    const ChangeStatusItem = (status: string,tracker_id: string) => {
        
        setLoading(true);
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
        console.log("tracker",tracker_id)
        console.log("status",status)

        const body = {
            auto_enable: status == 'auto-enable' ? true : false,
            disable: status == 'disabled' ? true : false,
            integration: [
                {
                    integration_id: tracker_id,
                },
            ],
        }
       
        
        axios
            .post(
                `${url}/main/compliance/api/v3/benchmark/${id}/assign`,
                body,
                config
            )
            .then((res) => {
                // window.location.reload()
                GetEnabled()
            })
            .catch((err) => {
                setLoading(false)

                console.log(err)
            })
    }
    useEffect(() => {
        if (firstLoading) {
            GetEnabled()
            setFirstLoading(false)
        }
    }, [firstLoading])
  
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
                            cell: (item) => item?.integration?.integration_id,
                            sortingField: 'id',
                            isRowHeader: true,
                        },
                        {
                            id: 'id_name',
                            header: 'Name',
                            cell: (item) => item?.integration?.name,
                            sortingField: 'id',
                            isRowHeader: true,
                        },
                        {
                            id: 'provider_id',
                            header: 'Provider ID',
                            cell: (item) => item?.integration?.provider_id,
                            sortingField: 'id',
                            isRowHeader: true,
                        },
                        {
                            id: 'integration_type',
                            header: 'Integration Type',
                            cell: (item) => item?.integration?.integration_type,
                            sortingField: 'id',
                            isRowHeader: true,
                        },
                        {
                            id: 'enable',
                            header: 'Enable',
                            cell: (item) => (
                                <>
                                    <Switch
                                        disabled={banner}
                                        onChange={(e) => {
                                            ChangeStatusItem(
                                                e ? 'auto-enable' : 'disabled',
                                                item?.integration
                                                    ?.integration_id
                                            )
                                        }}
                                        checked={item?.assigned}
                                    />
                                </>
                            ),
                            sortingField: 'id',
                            isRowHeader: true,
                        },
                    ]}
                    columnDisplay={[
                        { id: 'id', visible: true },
                        { id: 'name', visible: true },
                        { id: 'provider_id', visible: true },
                        { id: 'integration_type', visible: true },
                        { id: 'enable', visible: true },
                    ]}
                    enableKeyboardNavigation
                    // @ts-ignore
                    items={rows ? rows.slice(page * 10, (page + 1) * 10) : []}
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
                            className="w-full"
                            actions={
                                <Flex className="gap-2">
                                    <KButton
                                        onClick={() => {
                                            ChangeStatus('auto-enable')
                                        }}
                                    >
                                        Enable All
                                    </KButton>
                                    <KButton
                                        onClick={() => {
                                            ChangeStatus('disabled')
                                        }}
                                    >
                                        Disable All
                                    </KButton>
                                </Flex>
                            }
                        >
                            Assigments{' '}
                            <span className=" font-medium">
                                ({rows?.length})
                            </span>
                        </Header>
                    }
                    pagination={
                        <CustomPagination
                            currentPageIndex={page + 1}
                            pagesCount={Math.ceil(rows?.length / 10)}
                            onChange={({ detail }) =>
                                setPage(detail.currentPageIndex - 1)
                            }
                        />
                    }
                />
            </div>
        </>
    )
}
