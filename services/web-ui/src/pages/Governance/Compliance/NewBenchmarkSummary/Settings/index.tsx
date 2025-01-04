// @ts-noCheck
import { ValueFormatterParams } from 'ag-grid-community'
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
    const [enableStatus,setEnableStatus] = useState('')
    const [banner, setBanner] = useState(autoAssign)
    const isDemo = useAtomValue(isDemoAtom)
    const [loading, setLoading] = useState(false)
    const [rows,setRows] = useState<any>([])
       const [page, setPage] = useState(0)
    const {
        sendNow: sendEnable,
        isLoading: enableLoading,
        isExecuted: enableExecuted,
    } = useComplianceApiV1AssignmentsConnectionCreate(
        String(id),
        { integrationID: [transfer.connectionID] },
        {},
        false
    )
    const {
        response: enableAllResponse,
        sendNow: sendEnableAll,
        isLoading: enableAllLoading,
        isExecuted: enableAllExecuted,
    } = useComplianceApiV1AssignmentsConnectionCreate(
        String(id),
        { auto_assign: !allEnable },
        {},
        false
    )

    useEffect(() => {
        if (enableAllResponse) {
            isAutoResponse(true)
            setAllEnable(!allEnable)
            window.location.reload()
        }
    }, [enableAllResponse])

    const {
        sendNow: sendDisable,
        isLoading: disableLoading,
        isExecuted: disableExecuted,
    } = useComplianceApiV1AssignmentsConnectionDelete(
        String(id),
        { integrationID: [transfer.connectionID] },
        {},
        false
    )

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

    useEffect(() => {
        if (transfer.connectionID !== '') {
            if (transfer.status) {
                sendEnable()
            } else {
                sendDisable()
            }
        }
    }, [transfer])

    useEffect(() => {
        // if (firstLoading) {
        //     refreshList()
        // }
        // setFirstLoading(false)
    }, [])

   
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
            setEnableStatus(res.data.status)
       setLoading(false)
              
           })
           .catch((err) => {
       setLoading(false)

               console.log(err)
           })
   }
   const ChangeStatus = (status: string) => {
       
       setLoading(true)
       setEnableStatus(status)
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
           })
           .catch((err) => {
               setLoading(false)

               console.log(err)
           })
   }
    const ChangeStatusItem = (status: string,tracker_id: string) => {
        
        setLoading(true);
        setEnableStatus(status);
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
            integration: [
                {
                    integration_id: tracker_id,
                },
            ],
        }
        console.log(body,"body");
        
        axios
            .post(
                `${url}/main/compliance/api/v3/benchmark/${id}/assign`,
                body,
                config
            )
            .then((res) => {
                // window.location.reload()
                getEnabled()
            })
            .catch((err) => {
                setLoading(false)

                console.log(err)
            })
    }
   useEffect(() => {
       GetEnabled()
   }, [enableExecuted, disableExecuted])
    return (
        <>
            <Flex
                flexDirection="col"
                justifyContent="start"
                alignItems="center"
            >
                <Flex className="w-full mb-3">
                    <Tiles
                        value={enableStatus}
                        className="gap-8"
                        onChange={({ detail }) => {
                            ChangeStatus(detail?.value)
                        }}
                        items={[
                            {
                                value: 'disabled',
                                label: 'Disabled',
                                description:
                                    'Makes the framework inactive, with no assignments or audits.',
                                // disabled: true,
                            },
                            {
                                value: 'enabled',
                                label: `Enabled`,
                                description:
                                    'Select integrations from the list below to enable the framework for auditing.',
                            },
                            {
                                value: 'auto-enable',
                                label: `Auto Enabled`,
                                description:
                                    'Activates the framework on all integrations including any future integrations supported by the framework',
                            },
                        ]}
                    />
                </Flex>
                <Flex className="relative">
                    {/* {allEnable && (
                        <Flex
                            justifyContent="center"
                            className="w-full h-full absolute backdrop-blur-sm z-10 top-[50px] rounded-lg"
                            style={{ backgroundColor: 'rgba(0,0,0,0.3)' }}
                        >
                            <Text className="py-2 px-4 rounded bg-white border">
                                Auto onboard enabled
                            </Text>
                        </Flex>
                    )} */}
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
                                cell: (item) =>
                                    item?.integration?.integration_id,
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
                                cell: (item) =>
                                    item?.integration?.integration_type,
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
                                                    e
                                                        ? 'auto-enable'
                                                        : 'disabled',
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
                        items={
                            rows ? rows.slice(page * 10, (page + 1) * 10) : []
                        }
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
                            <Header className="w-full">
                                Assigments{' '}
                                <span className=" font-medium">
                                    ({rows?.length})
                                </span>
                            </Header>
                        }
                        pagination={
                            <Pagination
                                currentPageIndex={page + 1}
                                pagesCount={Math.ceil(rows?.length / 10)}
                                onChange={({ detail }) =>
                                    setPage(detail.currentPageIndex - 1)
                                }
                            />
                        }
                    />
                </Flex>
                {/* {banner ? (
                    <Callout
                        title="Provider requirements"
                        className="w-full"
                        color="amber"
                    >
                        <Flex
                            flexDirection="col"
                            alignItems="start"
                            className="gap-3"
                        >
                            <Text color="amber">
                                You have auto-enabled all accounts
                            </Text>
                            <Button
                                variant="secondary"
                                color="amber"
                                onClick={() => setBanner(false)}
                            >
                                Edit
                            </Button>
                        </Flex>
                    </Callout>
                ) : (
                  
                )} */}
                {/* <Divider /> */}
                <Flex
                    className="w-full gap-2  bg-white p-7 rounded-xl mt-2"
                    justifyContent="between"
                >
                    <Text className="text-gray-800 whitespace-nowrap">
                        Maintain Detailed audit trails of Drifts Events
                    </Text>
                    {changeSettingsLoading && changeSettingsExecuted ? (
                        <Spinner />
                    ) : (
                        <>
                            <Toggle
                                onChange={({ detail }) =>
                                    changeSettings(
                                        id,
                                        { tracksDriftEvents: detail.checked },
                                        {}
                                    )
                                }
                                checked={false}
                            ></Toggle>
                        </>
                    )}
                </Flex>
            </Flex>
        </>
    )
}
