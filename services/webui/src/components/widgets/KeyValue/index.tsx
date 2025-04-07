import { Button, Grid } from '@tremor/react'
import { useInventoryApiV1QueryRunCreate } from '../../../api/inventory.gen'
import { snakeCaseToLabel } from '../../../utilities/labelMaker'
import {
    Box,
    Header,
    KeyValuePairs,
    Modal,
    SpaceBetween,
    Table,
} from '@cloudscape-design/components'
import { useEffect, useMemo, useState } from 'react'
import axios from 'axios'

export interface KPIProps {
    kpis: Kpi[]
}
export interface Kpi {
    info:      string;
    count_kpi: string;
    list_kpi:  string;
}


export const getTable = (
    headers: string[] | undefined,
    details: any[][] | undefined
) => {
    const columns: any[] = []
    const rows: any[] = []
    const column_def: any[] = []
    const headerField = headers?.map((value, idx) => {
        if (headers.filter((v) => v === value).length > 1) {
            return `${value}-${idx}`
        }
        return value
    })
    if (headers && headers.length) {
        for (let i = 0; i < headers.length; i += 1) {
            const isHide = headers[i][0] === '_'
            // columns.push({
            //     field: headerField?.at(i),
            //     headerName: snakeCaseToLabel(headers[i]),
            //     type: 'string',
            //     sortable: true,
            //     hide: isHide,
            //     resizable: true,
            //     filter: true,
            //     width: 170,
            //     cellRenderer: (param: ValueFormatterParams) => (
            //         <span className={isDemo ? 'blur-sm' : ''}>
            //             {param.value}
            //         </span>
            //     ),
            // })
            columns.push({
                id: headerField?.at(i),
                header: snakeCaseToLabel(headers[i]),
                // @ts-ignore
                cell: (item: any) => (
                    <>
                        {/* @ts-ignore */}
                        {typeof item[headerField?.at(i)] == 'string'
                            ? // @ts-ignore
                              item[headerField?.at(i)]
                            : // @ts-ignore
                              JSON.stringify(item[headerField?.at(i)])}
                    </>
                ),
                maxWidth: '200px',
                // sortingField: 'id',
                // isRowHeader: true,
                // maxWidth: 150,
            })
            column_def.push({
                id: headerField?.at(i),
                visible: !isHide,
            })
        }
    }
    if (details && details.length) {
        for (let i = 0; i < details.length; i += 1) {
            const row: any = {}
            for (let j = 0; j < columns.length; j += 1) {
                row[headerField?.at(j) || ''] = details[i][j]
                //     typeof details[i][j] === 'string'
                //         ? details[i][j]
                //         : JSON.stringify(details[i][j])
            }
            rows.push(row)
        }
    }
    const count = rows.length

    return {
        columns,
        column_def,
        rows,
        count,
    }
}

export default function KeyValueWidget({ kpis }: KPIProps) {
    const [items, setItems] = useState<any[]>([])

   const RunQuery = (query_id: string) => {
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
           page: { no: 1, size: 1000 },
           // @ts-ignore
           engine: 'cloudql',
           query_id: query_id,
           use_cache: true,
       }

      return  axios
           .post(`${url}/main/core/api/v1/query/run `, body, config)
           
   }

    useEffect(()=>{
        if(kpis.length >0){
            const temp_items:any= []
            kpis?.map((item)=>{
                RunQuery(item.count_kpi).then((res) => {
                    const { columns, rows } = getTable(
                        res.data?.headers,
                        res.data?.details
                    )
                    temp_items.push({
                        title: item.info,
                        value: rows[0][columns[0].id],
                    })
                })
            })
            setItems(temp_items)

        }


    },[kpis])

    return (
        <>
            <KeyValuePairs
            columns={kpis.length}
            items={items}
            />
        
        </>
    )
}
