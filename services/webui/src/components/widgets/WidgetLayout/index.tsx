// @ts-nocheck
import { useAtom } from "jotai"
import { LayoutAtom } from "../../../store"
import * as React from 'react'
import Board from '@cloudscape-design/board-components/board'
import BoardItem from '@cloudscape-design/board-components/board-item'
import Header from '@cloudscape-design/components/header'
import { ButtonDropdown } from "@cloudscape-design/components"
import { useEffect, useState } from 'react'
import TableWidget from "../table"

const  COMPONENT_MAPPING ={
    'table': TableWidget
     
}


export default function WidgetLayout() {
    const [layout, setLayout] = useAtom(LayoutAtom)
    const [items, setItems] = useState([])
    useEffect(()=>{
        if(layout){
            setItems(layout?.layout_config)
        }

    },[layout])


    
    return (
        <div className="w-full h-full">
            <Board
                renderItem={(item) => (
                    <BoardItem
                        header={<Header>{item.data.title}</Header>}
                        settings={
                            <ButtonDropdown
                                items={[
                                    {
                                        id: 'preferences',
                                        text: 'Preferences',
                                    },
                                    { id: 'remove', text: 'Remove' },
                                ]}
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
                        {COMPONENT_MAPPING["table"](item?.data?.props) }
                    </BoardItem>
                )}
                onItemsChange={(event) => {
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
                        operationAnnouncement,
                        conflicts,
                        disturbed
                    ) {
                        const conflictsAnnouncement =
                            conflicts.length > 0
                                ? `Conflicts with ${conflicts
                                      .map((c) => c.data.title)
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
                        liveAnnouncementItemRemoved: (op) =>
                            createAnnouncement(
                                `Removed item ${op.item.data.title}.`,
                                [],
                                op.disturbed
                            ),
                        navigationAriaLabel: 'Board navigation',
                        navigationAriaDescription:
                            'Click on non-empty item to move focus over',
                        navigationItemAriaLabel: (item) =>
                            item ? item.data.title : 'Empty',
                    }
                })()}
            />
        </div>
    )
}
