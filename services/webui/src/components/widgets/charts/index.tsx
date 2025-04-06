import { Button, Grid } from '@tremor/react'
import { useInventoryApiV1QueryRunCreate } from '../../../api/inventory.gen'
import { snakeCaseToLabel } from '../../../utilities/labelMaker'
import {
    Box,
    Header,
    Modal,
    SpaceBetween,
    Table,
} from '@cloudscape-design/components'
import { useMemo, useState } from 'react'
import CustomPagination from '../../Pagination'
import { RenderObject } from '../../RenderObject'
import PieChart from '@cloudscape-design/components/pie-chart'
export interface TableProps {
    query: string
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

export default function ChartWidget({ query }: TableProps) {
    const [page, setPage] = useState(0)
    const [openDrawer, setOpenDrawer] = useState(false)
    const [selectedRow, setSelectedRow] = useState({})

    const {
        response: queryResponse,
        isLoading,
        isExecuted,
        sendNow,
        sendNowWithParams,
        error,
    } = useInventoryApiV1QueryRunCreate(
        {
            page: { no: 1, size: 1000 },
            // @ts-ignore
            engine: 'cloudql',
            query: query,
        },
        {},
        true
    )

    return (
        <>
            {' '}
            <Modal
                visible={openDrawer}
                onDismiss={() => setOpenDrawer(false)}
                header="Query Result"
                className="min-w-[500px]"
                size="large"
            >
                <RenderObject obj={selectedRow} />
            </Modal>
            <Grid numItems={1} className="w-full">
                <PieChart
                    data={[
                        {
                            title: 'Running',
                            value: 60,
                            lastUpdate: 'Dec 7, 2020',
                        },
                        {
                            title: 'Failed',
                            value: 30,
                            lastUpdate: 'Dec 6, 2020',
                        },
                        {
                            title: 'In-progress',
                            value: 10,
                            lastUpdate: 'Dec 6, 2020',
                        },
                        {
                            title: 'Pending',
                            value: 0,
                            lastUpdate: 'Dec 7, 2020',
                        },
                    ]}
                    detailPopoverContent={(datum, sum) => [
                        { key: 'Resource count', value: datum.value },
                        {
                            key: 'Percentage',
                            value: `${((datum.value / sum) * 100).toFixed(0)}%`,
                        },
                        { key: 'Last update on', value: datum.lastUpdate },
                    ]}
                    segmentDescription={(datum, sum) =>
                        `${datum.value} units, ${(
                            (datum.value / sum) *
                            100
                        ).toFixed(0)}%`
                    }
                    ariaDescription="Pie chart showing how many resources are currently in which state."
                    ariaLabel="Pie chart"
                    empty={
                        <Box textAlign="center" color="inherit">
                            <b>No data available</b>
                            <Box variant="p" color="inherit">
                                There is no data available
                            </Box>
                        </Box>
                    }
                    noMatch={
                        <Box textAlign="center" color="inherit">
                            <b>No matching data</b>
                            <Box variant="p" color="inherit">
                                There is no matching data to display
                            </Box>
                            <Button>Clear filter</Button>
                        </Box>
                    }
                />
            </Grid>
        </>
    )
}
