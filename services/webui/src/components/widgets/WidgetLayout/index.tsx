import { useAtom, useSetAtom } from 'jotai'
import { LayoutAtom, meAtom, notificationAtom } from '../../../store'
import * as React from 'react'
import Board from '@cloudscape-design/board-components/board'
import BoardItem from '@cloudscape-design/board-components/board-item'
import Header from '@cloudscape-design/components/header'
import {
    Alert,
    Button,
    ButtonDropdown,
    Input,
    Modal,
    Spinner,
} from '@cloudscape-design/components'
import { useEffect, useState } from 'react'
import TableWidget from '../table'
import axios from 'axios'
import ChartWidget from '../charts'
import KeyValueWidget from '../KeyValue'
import Shortcuts from '../../../pages/Overview/Shortcuts'
import Integrations from '../../../pages/Overview/Integrations'
import { array } from 'prop-types'
import SRE from '../../../pages/Overview/KPI_Cards'

const COMPONENT_MAPPING = {
    table: TableWidget,
    chart: ChartWidget,
    kpi: KeyValueWidget,
    shortcut: Shortcuts,
    integration: Integrations,
    sre: SRE
}
const NUMBER_MAPPING = {
    0: 'First',
    1: 'Second',
    2: 'Third',
    3: 'Fourth',
    4: 'Fifth',
}

export default function WidgetLayout() {
    const [layout, setLayout] = useAtom(LayoutAtom)
    const [me, setMe] = useAtom(meAtom)

    const [items, setItems] = useState([
        {
            id: 'shortcut',
            data: {
                componentId: 'shortcut',
                title: 'Shortcuts',
                description: '',
                props: {},
            },
            rowSpan: 2,
            columnSpan: 3,
            columnOffset: { '4': 0 },
        },
        {
            id: 'integration',
            data: {
                componentId: 'integration',
                title: 'Integrations',
                description: '',
                props: {},
            },
            rowSpan: 8,
            columnSpan: 1,
            columnOffset: { '4': 3 },
        },
    ])
    const [layoutLoading, setLayoutLoading] = useState<boolean>(false)
    const [addModalOpen, setAddModalOpen] = useState(false)
    const [selectedAddItem, setSelectedAddItem] = useState<any>('')
    const [widgetProps, setWidgetProps] = useState<any>({})
    const setNotification = useSetAtom(notificationAtom)
    useEffect(() => {
        if (layout) {
            console.log(layout)
            // add to ietms
            if (items.length !== layout?.layout_config.length + 2) {
                setItems([...items, ...(layout?.layout_config || [])])
            }
        }
    }, [layout])
    const GetComponent = (name: string, props: any) => {
        // @ts-ignore
        const Component = COMPONENT_MAPPING[name]
        if (Component) {
            return <Component {...props} />
        }
        return null
    }
    const SetDefaultLayout = (layout: any) => {
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
            user_id: me?.username,
            layout_config: layout,
            name: 'default',
            is_private: true,
        }

        axios
            .post(`${url}/main/core/api/v4/layout/set`, body, config)
            .then((res) => {})
            .catch((err) => {
                console.log(err)
            })
    }
    const GetDefaultLayout = () => {
        setLayoutLoading(true)
        axios
            .get(
                `https://raw.githubusercontent.com/opengovern/platform-configuration/refs/heads/main/default_layout.json`
            )
            .then((res) => {
                setLayout(res?.data)
                setLayoutLoading(false)
            })
            .catch((err) => {
                setLayoutLoading(false)
            })
    }
    const HandleRemoveItemByID = (id: string) => {
        const newItems = items.filter((item: any) => item.id !== id)
        setItems(newItems)
    }
    const HandleWidgetProps = () => {
        if (selectedAddItem == 'table') {
            return (
                <>
                    <Input
                        placeholder="Query ID"
                        value={widgetProps?.query_id}
                        onChange={(e: any) => {
                            setWidgetProps({
                                ...widgetProps,
                                query_id: e.detail.value,
                            })
                        }}
                    />
                    <Input
                        placeholder="Rows to display"
                        value={widgetProps?.display_rows}
                        onChange={(e: any) => {
                            setWidgetProps({
                                ...widgetProps,
                                display_rows: e.detail.value,
                            })
                        }}
                    />
                </>
            )
        }
        if (selectedAddItem == 'kpi') {
            return (
                <>
                    {/* map 4 times and return 3 input count_kpi info list_kpi all of them inside kpies key */}
                    {[0, 1, 2, 3].map((item, index: number) => {
                        return (
                            <div key={index} className="flex flex-col gap-2">
                                <Input
                                    placeholder={` ${
                                        //   @ts-ignore
                                        NUMBER_MAPPING[index.toString()]
                                    } KPI Name`}
                                    value={widgetProps?.kpis?.[index]?.info}
                                    onChange={(e: any) => {
                                        setWidgetProps({
                                            ...widgetProps,
                                            kpis: [
                                                ...(widgetProps?.kpis || []),
                                                {
                                                    info: e.detail.value,
                                                    count_kpi: '',
                                                    list_kpi: '',
                                                },
                                            ],
                                        })
                                    }}
                                />
                                <Input
                                    placeholder={` ${
                                        //   @ts-ignore
                                        NUMBER_MAPPING[index.toString()]
                                    } Count Query ID`}
                                    value={
                                        widgetProps?.kpis?.[index]?.count_kpi
                                    }
                                    onChange={(e: any) => {
                                        setWidgetProps({
                                            ...widgetProps,
                                            kpis: [
                                                ...(widgetProps?.kpis || []),
                                                {
                                                    info: widgetProps?.kpis?.[
                                                        index
                                                    ]?.info,
                                                    count_kpi: e.detail.value,
                                                    list_kpi: '',
                                                },
                                            ],
                                        })
                                    }}
                                />
                                <Input
                                    placeholder={` ${
                                        //   @ts-ignore
                                        NUMBER_MAPPING[index.toString()]
                                    } List Query ID`}
                                    value={widgetProps?.kpis?.[index]?.list_kpi}
                                    onChange={(e: any) => {
                                        setWidgetProps({
                                            ...widgetProps,
                                            kpis: [
                                                ...(widgetProps?.kpis || []),
                                                {
                                                    info: widgetProps?.kpis?.[
                                                        index
                                                    ]?.info,
                                                    count_kpi:
                                                        widgetProps?.kpis?.[
                                                            index
                                                        ]?.count_kpi,
                                                    list_kpi: e.detail.value,
                                                },
                                            ],
                                        })
                                    }}
                                />
                            </div>
                        )
                    })}
                </>
            )
        }
    }
    const HandleAddWidget = () => {
        if(!widgetProps?.title || !widgetProps?.description) {
            return
        }
        const newItem = {
            id: `${selectedAddItem}-${items.length}`,
            data: {
                componentId: selectedAddItem,
                props: widgetProps,
                title: widgetProps?.title,
                description: widgetProps?.description,
            },
            rowSpan: 2,
            columnSpan: 2,
            columnOffset: { '4': 0 },
        }
        // @ts-ignore
        setItems([...items, newItem])
        setAddModalOpen(false)
        setWidgetProps({})
    }
    const HandleAddProductWidgets =(id: string)=>{
        // check if id not exist in items
        const check = items.filter((item: any) => item.id === id)
        if(check.length > 0){
              setNotification({
                  text: `Widget Already exist`,
                  type: 'error',
              })
            return
        }
        if(id=='integration'){
            const new_item = {
                id: 'integration',
                data: {
                    componentId: 'integration',
                    title: 'Integrations',
                    description: '',
                    props: {},
                },
                rowSpan: 8,
                columnSpan: 1,
                columnOffset: { '4': 3 },
            }
            setItems([...items, new_item])

        }
        if(id=='shortcut'){
            const new_item = {
                id: 'shortcut',
                data: {
                    componentId: 'shortcut',
                    title: 'Shortcuts',
                    description: '',
                    props: {},
                },
                rowSpan: 2,
                columnSpan: 3,
                columnOffset: { '4': 0 },
            }
            setItems([...items, new_item])
        }
        if(id=='sre'){
            const new_item = {
                id: 'sre',
                data: {
                    componentId: 'sre',
                    title: 'SRE',
                    description: '',
                    props: {},
                },
                rowSpan: 2,
                columnSpan: 3,
                columnOffset: { '4': 0 },
            }
            setItems([...items, new_item])
        }
        return

    }

    return (
        <div className="w-full h-full flex flex-col gap-8">
            <Header
                variant="h1"
                actions={
                    <div className="flex flex-row gap-2">
                        <ButtonDropdown
                            items={[
                                {
                                    id: 'reset',
                                    text: 'Reset to default layout',
                                },
                                { id: 'save', text: 'save' },
                            ]}
                            onItemClick={(event: any) => {
                                if (event.detail.id == 'reset') {
                                    GetDefaultLayout()
                                }
                                if (event.detail.id == 'save') {
                                    SetDefaultLayout(items)
                                }
                            }}
                            ariaLabel="Board item settings"
                        >
                            Dashboard settings
                        </ButtonDropdown>
                        <ButtonDropdown
                            items={[
                                { id: 'table', text: 'Table Widget' },
                                { id: 'chart', text: 'Pie Chart Widget' },
                                { id: 'kpi', text: 'KPI Widget' },
                                { id: 'integration', text: 'Integrations' },
                                { id: 'shortcut', text: 'Shortcuts' },
                                { id: 'sre', text: 'SRE' },
                            ]}
                            onItemClick={(event: any) => {
                                if (
                                    event.detail.id == 'sre' ||
                                    event.detail.id == 'shortcut' ||
                                    event.detail.id == 'integration'
                                ) {
                                    HandleAddProductWidgets(event.detail.id)
                                } else {
                                    setSelectedAddItem(event.detail.id)
                                    setAddModalOpen(true)
                                }
                               
                            }}
                            ariaLabel="Board item settings"
                        >
                            Add Widget
                        </ButtonDropdown>
                    </div>
                }
            >
                Service Dashboard
            </Header>
            {layoutLoading ? (
                <Spinner />
            ) : (
                <Board
                    renderItem={(item: any) => (
                        <BoardItem
                            header={
                                <Header description={item?.data?.description}>
                                    {item.data.title}
                                </Header>
                            }
                            settings={
                                <ButtonDropdown
                                    items={[{ id: 'remove', text: 'Remove' }]}
                                    onItemClick={(event) => {
                                        if (event.detail.id === 'remove') {
                                            HandleRemoveItemByID(item.id)
                                        }
                                    }}
                                    ariaLabel="Board item settings"
                                    variant="icon"
                                />
                            }
                            i18nStrings={{
                                dragHandleAriaLabel: 'Drag handle',
                                dragHandleAriaDescription:
                                    'Use Space or Enter to activate drag, arrow keys to move, Space or Enter to submit, or Escape to discard. Be sure to temporarily disable any screen reader navigation feature that may interfere with the functionality of the arrow keys.',
                                resizeHandleAriaLabel: 'Resize handle',
                                resizeHandleAriaDescription:
                                    'Use Space or Enter to activate resize, arrow keys to move, Space or Enter to submit, or Escape to discard. Be sure to temporarily disable any screen reader navigation feature that may interfere with the functionality of the arrow keys.',
                            }}
                        >
                            {GetComponent(
                                item?.data?.componentId,
                                item?.data?.props
                            )}
                        </BoardItem>
                    )}
                    onItemsChange={(event: any) => {
                        setItems(event.detail.items)
                    }}
                    items={items}
                    empty={
                        <div className="flex flex-col items-center justify-center w-full h-full">
                            <span className="text-gray-500">No items</span>
                        </div>
                    }
                    i18nStrings={(() => {
                        function createAnnouncement(
                            operationAnnouncement: any,
                            conflicts: any,
                            disturbed: any
                        ) {
                            const conflictsAnnouncement =
                                //
                                conflicts?.length > 0
                                    ? `Conflicts with ${conflicts
                                          .map((c: any) => c.data.title)
                                          .join(', ')}.`
                                    : ''
                            const disturbedAnnouncement =
                                disturbed.length > 0
                                    ? `Disturbed ${disturbed.length} items.`
                                    : ''
                            return [
                                operationAnnouncement,
                                conflictsAnnouncement,
                                disturbedAnnouncement,
                            ]
                                .filter(Boolean)
                                .join(' ')
                        }
                        return {
                            liveAnnouncementDndStarted: (operationType) =>
                                operationType === 'resize'
                                    ? 'Resizing'
                                    : 'Dragging',
                            liveAnnouncementDndItemReordered: (operation) => {
                                const columns = `column ${
                                    operation.placement.x + 1
                                }`
                                const rows = `row ${operation.placement.y + 1}`
                                return createAnnouncement(
                                    `Item moved to ${
                                        operation.direction === 'horizontal'
                                            ? columns
                                            : rows
                                    }.`,
                                    operation.conflicts,
                                    operation.disturbed
                                )
                            },
                            liveAnnouncementDndItemResized: (operation) => {
                                const columnsConstraint =
                                    operation.isMinimalColumnsReached
                                        ? ' (minimal)'
                                        : ''
                                const rowsConstraint =
                                    operation.isMinimalRowsReached
                                        ? ' (minimal)'
                                        : ''
                                const sizeAnnouncement =
                                    operation.direction === 'horizontal'
                                        ? `columns ${operation.placement.width}${columnsConstraint}`
                                        : `rows ${operation.placement.height}${rowsConstraint}`
                                return createAnnouncement(
                                    `Item resized to ${sizeAnnouncement}.`,
                                    operation.conflicts,
                                    operation.disturbed
                                )
                            },
                            liveAnnouncementDndItemInserted: (operation) => {
                                const columns = `column ${
                                    operation.placement.x + 1
                                }`
                                const rows = `row ${operation.placement.y + 1}`
                                return createAnnouncement(
                                    `Item inserted to ${columns}, ${rows}.`,
                                    operation.conflicts,
                                    operation.disturbed
                                )
                            },
                            liveAnnouncementDndCommitted: (operationType) =>
                                `${operationType} committed`,
                            liveAnnouncementDndDiscarded: (operationType) =>
                                `${operationType} discarded`,
                            liveAnnouncementItemRemoved: (op: any) =>
                                createAnnouncement(
                                    `Removed item ${op.item.data.title}.`,
                                    [],
                                    op.disturbed
                                ),
                            navigationAriaLabel: 'Board navigation',
                            navigationAriaDescription:
                                'Click on non-empty item to move focus over',
                            navigationItemAriaLabel: (item: any) =>
                                item ? item.data.title : 'Empty',
                        }
                    })()}
                />
            )}
            <Modal
                visible={addModalOpen}
                onDismiss={() => {
                    setAddModalOpen(false)
                }}
                header={`Add ${
                    selectedAddItem?.charAt(0).toUpperCase() +
                    selectedAddItem?.slice(1)
                } Widget`}
            >
                <div className="flex flex-col gap-2">
                    <Input
                        placeholder="Widget Name"
                        ariaRequired={true}
                        value={widgetProps?.title}
                        onChange={(e: any) => {
                            setWidgetProps({
                                ...widgetProps,
                                title: e.detail.value,
                            })
                        }}
                    />
                    <Input
                        placeholder="Widget description"
                        ariaRequired={true}
                        value={widgetProps?.description}
                        onChange={(e: any) => {
                            setWidgetProps({
                                ...widgetProps,
                                description: e.detail.value,
                            })
                        }}
                    />
                    {HandleWidgetProps()}
                    {(!widgetProps?.title || !widgetProps?.description) && (
                        <Alert
                            type="error"
                            header="Please fill all the required fields"
                        >
                            Please fill all the required fields
                        </Alert>
                    )}
                    <div className="flex w-full justify-end items-center">
                        <Button
                            onClick={() => {
                                HandleAddWidget()
                            }}
                        >
                            Submit
                        </Button>
                    </div>
                </div>
            </Modal>
        </div>
    )
}
