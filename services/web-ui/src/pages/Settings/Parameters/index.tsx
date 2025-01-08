import { PlusIcon } from '@heroicons/react/24/outline'
import {
    ArrowRightCircleIcon,
    KeyIcon,
    PlusCircleIcon,
    TrashIcon,
} from '@heroicons/react/24/solid'
import { useEffect, useState } from 'react'
import {
    Button,
    Card,
    Divider,
    Flex,
    TextInput,
    Textarea,
    Title,
} from '@tremor/react'
import { useAtom, useAtomValue } from 'jotai'
import {
    useMetadataApiV1QueryParameterCreate,
    useMetadataApiV1QueryParameterList,
} from '../../../api/metadata.gen'
import { getErrorMessage } from '../../../types/apierror'
import { notificationAtom } from '../../../store'
import { searchAtom, useURLParam } from '../../../utilities/urlstate'
import TopHeader from '../../../components/Layout/Header'
import axios from 'axios'
import {
    Alert,
    Box,
    Header,
    Input,
    KeyValuePairs,
    Link,
    Modal,
    Pagination,
    PropertyFilter,
    RadioGroup,
    Select,
    SpaceBetween,
    Spinner,
    Table,
    Toggle,
} from '@cloudscape-design/components'
import KButton from '@cloudscape-design/components/button'

interface IParam {
    key: string
    value: string
}

export default function SettingsParameters() {
    const [notif, setNotif] = useAtom(notificationAtom)
    const [params, setParams] = useState([])
    const [page, setPage] = useState(1)
    const [total, setTotal] = useState(0)
    const [loading, setLoading] = useState(false)
    const [detailLoading, setDetailLoading] = useState(false)
    const [selectedItem, setSelectedItem] = useState<any>()
    const [selected, setSelected] = useState<any>()
    const [open, setOpen] = useState(false)
    const [controls, setControls] = useState([])
    const [queries, setQueries] = useState([])
    const [queryDone, setQueryDone] = useState(false)
    const [controlDone, setControlDone] = useState(false)
    const [queryToken, setQueryToken] = useState({
        tokens: [],
        operation: 'and',
    })
    const [propertyOptions, setPropertyOptions] = useState([])
    const [editValue, setEditValue] = useState({
        key: '',
        value: '',
        control_id: '',
    })
     const [addValue, setAddValue] = useState({
         key: '',
         value: '',
         control_id: {},
     })
     const [addError, setAddError] = useState('')
     const [addOpen, setAddOpen] = useState(false)
    const [sortOrder, setSortOrder] = useState('asc')
    const [sortField, setSortField] = useState<any>()

    const GetParams = () => {
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

        let body :any ={
            cursor: page,
            per_page: 15
        }
        const controls: any = []
        const queries: any = []
        const titles: any = []
        queryToken?.tokens?.map((t: any) => {
            if (t.propertyKey === 'controls') {
                controls.push(t.value)
            } 
            if (t.propertyKey === 'queries') {
                queries.push(t.value)
            }
            if (t.propertyKey === 'key_regex') {
                titles.push(t.value)
            }
        })
        if (controls.length > 0) {
            body['controls'] = controls
        }
        if(queries.length > 0){
            body['queries'] = queries
        }
        if(titles.length > 0){
            body['key_regex'] = titles[0]
        }
        if(sortField){
            body['sort_by'] = sortField?.sortingField
            body['sort_order'] = sortOrder
        }
        axios
            .post(
                `${url}/main/core/api/v1/query_parameter`,body,
                config
            )
            .then((res) => {
                const data = res.data
                setParams(data?.items)
                setTotal(data?.total_count)

                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }

    const EditParams = () => {
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
            query_parameters: [
                {
                    key: editValue.key,
                    value: editValue.value,
                    control_id: editValue?.control_id ? editValue.control_id : '',
                         
                },
            ],
        }

        axios
            .post(`${url}/main/core/api/v1/query_parameter/set`, body, config)
            .then((res) => {
                GetParams()
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)
            })
    }
     const AddParams = () => {
          if (addValue.key === '' || addValue.value === '') {
              setAddError('Key and Value cannot be empty')
              setLoading(false)
              return
          }

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
             query_parameters: [
                 {
                     key: addValue.key,
                     value: addValue.value,
                     //  @ts-ignore
                     control_id: addValue?.control_id?.value
                         ? //  @ts-ignore
                           addValue.control_id?.value
                         : '',
                 },
             ],
         }

         axios
             .post(`${url}/main/core/api/v1/query_parameter/set`, body, config)
             .then((res) => {
                 setLoading(false)
                 setAddOpen(false)
                 setAddValue({
                     key: '',
                     value: '',
                     control_id: {},
                 })
                 GetParams()
             })
             .catch((err) => {
                 console.log(err)
                 setLoading(false)
             })
     }
     const GetParamDetail = (key: string) => {
         setDetailLoading(true)
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
             .get(`${url}/main/core/api/v1/query_parameter/${key}`, config)
             .then((res) => {
                 if(res.data){
                    setSelectedItem(res.data)
                 }
                 setDetailLoading(false)
             })
             .catch((err) => {
                 console.log(err)
                 setOpen(false)
                 setDetailLoading(false)
             })
     }
     const GetControls = () => {
        //  setLoading(true)
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
            has_parameters: true
         }

         axios
             .post(`${url}/main/compliance/api/v3/controls`, body,config)
             .then((res) => {
                 if (res.data) {
                        setControls(res.data?.items)
                    //  setSelectedItem(res.data)
                    //  setOpen(true)
                 }
                 setControlDone(true)
                //  setLoading(false)
             })
             .catch((err) => {
                 console.log(err)
                //  setLoading(false)
             })
     }
const GetQueries = () => {
    // setLoading(true)
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
        has_parameters: true,
    }

    axios
        .post(`${url}/main/core/api/v3/queries `, body, config)
        .then((res) => {
            if (res.data) {
                setQueries(res.data?.items)
                //  setSelectedItem(res.data)
                //  setOpen(true)
            }
            setQueryDone(true)
            // setLoading(false)
        })
        .catch((err) => {
            console.log(err)
            // setLoading(false)
        })
}

useEffect(()=>{
    if(queryDone && controlDone){
        let options :any = []
        controls?.map((c: any) => {
            options.push({
                propertyKey: 'controls',
                value: c.id,
            })
        })
        queries?.map((c: any) => {
            options.push({
                propertyKey: 'queries',
                value: c.id,
            })
        })
        setPropertyOptions(options)
    }
},[queryDone,controlDone])

    

    useEffect(() => {
        GetControls()
        GetQueries()
    }, [page])
    useEffect(() => {
        GetParams()
    }, [queryToken,page,sortOrder,sortField])


    return (
        <>
            <Modal
                visible={open}
                onDismiss={() => setOpen(false)}
                header="Parameter Detail"
            >
                {detailLoading ? (
                    <Spinner />
                ) : (
                    <>
                        <KeyValuePairs
                            columns={4}
                            items={[
                                { label: 'Key', value: selectedItem?.key },
                                { label: 'Value', value: selectedItem?.value },
                                {
                                    label: 'Using control count',
                                    value: selected?.controls_count,
                                },
                                {
                                    label: 'Using query count',
                                    value: selected?.queries_count,
                                },
                                {
                                    label: 'Controls',
                                    value: (
                                        <>
                                            {selectedItem?.controls?.map(
                                                (c: any) => {
                                                    return (
                                                        <>
                                                            <Link
                                                                href={`/incidents/${c.id}`}
                                                            >
                                                                {c.title}
                                                            </Link>
                                                        </>
                                                    )
                                                }
                                            )}
                                        </>
                                    ),
                                },
                                {
                                    label: 'Queries',
                                    value: (
                                        <>
                                            {selectedItem?.queries?.map(
                                                (c: any) => {
                                                    return (
                                                        <>
                                                            <Link
                                                                href={`/incidents/${c.id}`}
                                                            >
                                                                {c.title}
                                                            </Link>
                                                        </>
                                                    )
                                                }
                                            )}
                                        </>
                                    ),
                                },
                            ]}
                        />
                    </>
                )}
            </Modal>
            <Modal
                visible={addOpen}
                onDismiss={() => setAddOpen(false)}
                header="Add Parameter"
            >
                <Flex
                    flexDirection="col"
                    justifyContent="start"
                    alignItems="start"
                    className="gap-3"
                >
                    <Input
                        name="Key"
                        className="w-full"
                        placeholder="Key"
                        value={addValue.key}
                        onChange={(e) => {
                            setAddValue({ ...addValue, key: e.detail.value })
                        }}
                    />
                    <Input
                        name="Value"
                        className="w-full"
                        placeholder="Value"
                        value={addValue.value}
                        onChange={(e) => {
                            setAddValue({ ...addValue, value: e.detail.value })
                        }}
                    />
                    <Select
                        selectedOption={addValue.control_id}
                        placeholder='controls'
                        inlineLabelText='Controls'
                        className="w-full"
                        onChange={({ detail }) =>
                            setAddValue({
                                ...addValue,
                                control_id: detail.selectedOption,
                            })
                        }
                        options={controls.map((c: any) => {
                            return {
                                value: c.id,
                                label: c.title,
                            }
                        })}
                    />
                    {addError && addError != '' && (
                        <>
                            <Alert
                                statusIconAriaLabel="Error"
                                className='w-full'
                                type="error"
                                header="Invalid Input"
                            >
                                {addError}
                            </Alert>
                        </>
                    )}
                    <div className='flex flex-row justify-end w-full'>
                        <KButton onClick={AddParams}>Add</KButton>
                    </div>
                </Flex>
            </Modal>
            <Table
                className="mt-2"
                sortingColumn={sortField}
                sortingDescending={sortOrder === 'desc' ? true : false}
                onSortingChange={({ detail }) => {
                    if (detail.isDescending) {
                        setSortOrder('desc')
                    } else {
                        setSortOrder('asc')
                    }
                    setSortField(detail.sortingColumn)
                }}
                columnDefinitions={[
                    {
                        id: 'key',
                        header: 'Key Name',
                        cell: (item: any) => item.key,
                        maxWidth: 150,
                        sortingField: 'key',
                    },

                    {
                        id: 'value',
                        header: 'Value',
                        cell: (item: any) => item.value,
                        sortingField: 'value',
                        maxWidth:120,
                        editConfig: {
                            ariaLabel: 'Value',
                            editIconAriaLabel: 'editable',
                            editingCell: (item, { currentValue, setValue }) => {
                                return (
                                    <Input
                                        autoFocus={true}
                                        value={currentValue ?? item.value}
                                        onChange={(event) => {
                                            setValue(event.detail.value)
                                            setEditValue({
                                                key: item.key,
                                                value: event.detail.value,
                                                control_id: item?.control_id
                                                    ? item.control_id
                                                    : '',
                                            })
                                        }}
                                    />
                                )
                            },
                        },
                    },
                    {
                        id: 'control_id',
                        header: 'Control',
                        cell: (item: any) =>
                            item.control_id ? item.control_id : 'Global',
                        maxWidth: 200,
                        sortingField: 'control_id',
                    },
                    {
                        id: 'controls_count',
                        header: 'Using control count',
                        maxWidth: 100,
                        sortingField: 'controls_count',

                        cell: (item: any) =>
                            item?.controls_count ? item?.controls_count : 0,
                    },

                    {
                        id: 'queries_count',
                        header: 'Using query count',
                        sortingField: 'queries_count',
                        maxWidth: 100,
                        cell: (item: any) =>
                            item?.queries_count ? item?.queries_count : 0,
                    },
                    {
                        id: 'action',
                        header: '',
                        maxWidth: 100,
                        cell: (item) => (
                            // @ts-ignore
                            <KButton
                                onClick={() => {
                                    GetParamDetail(item.key)
                                    setSelected(item)
                                    setOpen(true)
                                }}
                                variant="inline-link"
                                ariaLabel={`Open Detail`}
                            >
                                See details
                            </KButton>
                        ),
                    },
                ]}
                columnDisplay={[
                    { id: 'key', visible: true },
                    { id: 'value', visible: true },
                    { id: 'control_id', visible: true },
                    { id: 'controls_count', visible: true },
                    { id: 'queries_count', visible: true },
                    { id: 'action', visible: true },
                ]}
                loading={loading}
                submitEdit={async () => {
                    EditParams()
                }}
                // @ts-ignore
                items={params ? params : []}
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
                            <Flex className="gap-2">
                                <KButton
                                    onClick={() => {
                                        setAddOpen(true)
                                    }}
                                >
                                    Add
                                </KButton>
                                <KButton onClick={GetParams}>Reload</KButton>
                            </Flex>
                        }
                        className="w-full"
                    >
                        Parameters {total != 0 ? `(${total})` : ''}
                    </Header>
                }
                pagination={
                    <Pagination
                        currentPageIndex={page}
                        pagesCount={Math.ceil(total / 15)}
                        onChange={({ detail }) =>
                            setPage(detail.currentPageIndex)
                        }
                    />
                }
                filter={
                    <PropertyFilter
                        // @ts-ignore
                        query={queryToken}
                        // @ts-ignore
                        onChange={({ detail }) => {
                            // @ts-ignore
                            setQueryToken(detail)
                        }}
                        // countText="5 matches"
                        // enableTokenGroups
                        expandToViewport
                        filteringAriaLabel="Parameter Filters"
                        // @ts-ignore
                        // filteringOptions={filters}
                        filteringPlaceholder="Parameter Filters"
                        // @ts-ignore
                        filteringOptions={propertyOptions}
                        // @ts-ignore

                        filteringProperties={[
                            {
                                key: 'controls',
                                operators: ['='],
                                propertyLabel: 'Controls',
                                groupValuesLabel: 'Control values',
                            },
                            {
                                key: 'queries',
                                operators: ['='],
                                propertyLabel: 'Queries',
                                groupValuesLabel: 'Query values',
                            },
                            {
                                key: 'key_regex',
                                operators: ['='],
                                propertyLabel: 'Key',
                                groupValuesLabel: 'Key',
                            },
                        ]}
                        // filteringProperties={
                        //     filterOption
                        // }
                    />
                }
            />

            {/* <TopHeader /> */}
            {/* <Card key="summary" className="">
                <Flex>
                    <Title className="font-semibold">Variables</Title>
                    <Button
                        variant="secondary"
                        icon={PlusIcon}
                        onClick={addRow}
                    >
                        Add
                    </Button>
                </Flex>
                <Divider />

                <Flex flexDirection="col" className="mt-4">
                    {params.map((p, idx) => {
                        return (
                            <Flex flexDirection="row" className="mb-4">
                                <KeyIcon className="w-10 mr-3" />
                                <TextInput
                                    id={p.key}
                                    value={p.key}
                                    onValueChange={(e) =>
                                        updateKey(String(e), idx)
                                    }
                                    className={
                                        keyParam === p.key
                                            ? 'border-red-500'
                                            : ''
                                    }
                                />
                                <ArrowRightCircleIcon className="w-10 mx-3" />
                                <Textarea
                                    value={p.value}
                                    onValueChange={(e) =>
                                        updateValue(String(e), idx)
                                    }
                                    rows={1}
                                    className={
                                        keyParam === p.key
                                            ? 'border-red-500'
                                            : ''
                                    }
                                />
                                <TrashIcon
                                    className="w-10 ml-3 hover:cursor-pointer"
                                    onClick={() => deleteRow(idx)}
                                />
                            </Flex>
                        )
                    })}
                </Flex>
                <Flex flexDirection="row" justifyContent="end">
                    <Button
                        variant="secondary"
                        className="mx-4"
                        onClick={() => {
                            refresh()
                        }}
                        loading={isExecuted && isLoading}
                    >
                        Reset
                    </Button>
                    <Button
                        onClick={() => {
                            sendNowWithParams(
                                {
                                    queryParameters: params,
                                },
                                {}
                            )
                        }}
                        loading={updateIsExecuted && updateIsLoading}
                    >
                        Save
                    </Button>
                </Flex>
            </Card> */}
        </>
    )
}
