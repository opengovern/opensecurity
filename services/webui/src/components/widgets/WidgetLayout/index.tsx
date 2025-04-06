import { useAtom } from "jotai"
import { LayoutAtom } from "../../../store"
import * as React from 'react'
import Board from '@cloudscape-design/board-components/board'
import BoardItem from '@cloudscape-design/board-components/board-item'
import Header from '@cloudscape-design/components/header'




export default function WidgetLayout() {
    const [layout, setLayout] = useAtom(LayoutAtom)
    
    return (
       <></>
    )
}
