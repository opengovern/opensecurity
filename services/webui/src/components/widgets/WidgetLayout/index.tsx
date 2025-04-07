import { useAtom } from 'jotai'
import { LayoutAtom, meAtom } from '../../../store'
import * as React from 'react'
import Board from '@cloudscape-design/board-components/board'
import BoardItem from '@cloudscape-design/board-components/board-item'
import Header from '@cloudscape-design/components/header'
import {
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

const COMPONENT_MAPPING = {
    table: TableWidget,
    chart: ChartWidget,
    'kpi': KeyValueWidget,
    'shortcut': Shortcuts,
    'integration': Integrations,
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
                props: {
                 
                },
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
    useEffect(() => {
        if (layout) {
            // add to ietms
            setItems([...items, ...(layout?.layout_config || [])])
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
            user_id: me?.id,
            layout_config: layout,
            name: 'default',
            is_private: true,
        }

        axios
            .post(`${url}/main/core/api/v4/layout/set `, body, config)
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
    }
    const HandleAddWidget = () => {
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

    return (
        <div className="w-full h-full flex flex-col gap-2">
            <Header
                actions={
                    <div className="flex flex-row gap-2">
                        <Button
                            onClick={() => {
                                GetDefaultLayout()
                            }}
                        >
                            Reset to default layout
                        </Button>
                        <Button
                            onClick={() => {
                                SetDefaultLayout(items)
                            }}
                        >
                            save
                        </Button>
                        <ButtonDropdown
                            items={[
                                { id: 'table', text: 'Table Widget' },
                                { id: 'chart', text: 'Pie Chart Widget' },
                                { id: 'kpi', text: 'KPI Widget' },
                            ]}
                            onItemClick={(event: any) => {
                                setSelectedAddItem(event.detail.id)
                                setAddModalOpen(true)
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
                        value={widgetProps?.description}
                        onChange={(e: any) => {
                            setWidgetProps({
                                ...widgetProps,
                                description: e.detail.value,
                            })
                        }}
                    />
                    {HandleWidgetProps()}
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
