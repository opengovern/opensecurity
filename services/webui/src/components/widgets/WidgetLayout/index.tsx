import { useAtom } from "jotai"
import { LayoutAtom, meAtom } from "../../../store"
import * as React from 'react'
import Board from '@cloudscape-design/board-components/board'
import BoardItem from '@cloudscape-design/board-components/board-item'
import Header from '@cloudscape-design/components/header'
import { Button, ButtonDropdown, Spinner } from "@cloudscape-design/components"
import { useEffect, useState } from 'react'
import TableWidget from "../table"
import axios from "axios"
import ChartWidget from "../charts"

const  COMPONENT_MAPPING ={
    'table': TableWidget,
    'chart': ChartWidget
     
}


export default function WidgetLayout() {
    const [layout, setLayout] = useAtom(LayoutAtom)
        const [me, setMe] = useAtom(meAtom)
    
    const [items, setItems] = useState([])
    const [layoutLoading, setLayoutLoading] = useState<boolean>(false)
    useEffect(()=>{
        if(layout){
            setItems(layout?.layout_config)
        }

    },[layout])
    const GetComponent =(name:string,props :any)=>{
        // @ts-ignore
        const Component = COMPONENT_MAPPING[name]
        if(Component){
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
                        <Button>Add widget</Button>
                    </div>
                }
            >
                opensecurity Dashboard
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
                                    items={[
                                        
                                        { id: 'remove', text: 'Remove' },
                                    ]}
                                    onItemClick={(event)=>{
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
                        console.log(event.detail.items)
                        console.log(event.detail)

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
        </div>
    )
}
