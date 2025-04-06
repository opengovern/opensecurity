import { Grid } from "@tremor/react"
import { useInventoryApiV1QueryRunCreate } from "../../../api/inventory.gen"
import { snakeCaseToLabel } from "../../../utilities/labelMaker"
import { Box, Header, Modal, SpaceBetween, Table } from "@cloudscape-design/components"
import { useMemo, useState } from "react"
import CustomPagination from "../../Pagination"
import { RenderObject } from "../../RenderObject"

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


export default function TableWidget({query}: TableProps) {
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
      const memoCount = useMemo(
          () => getTable(queryResponse?.headers, queryResponse?.result).count,
          [queryResponse]
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
                <Table
                    // resizableColumns
                    variant="full-page"
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
                    // stickyHeader={true}
                    resizableColumns={true}
                    // stickyColumns={
                    //  {   first:1,
                    //     last: 1}
                    // }
                    onRowClick={(event) => {
                        const row = event.detail.item
                        // @ts-ignore
                        setSelectedRow(row)
                        setOpenDrawer(true)
                    }}
                    columnDefinitions={
                        getTable(queryResponse?.headers, queryResponse?.result)
                            .columns
                    }
                    columnDisplay={
                        getTable(queryResponse?.headers, queryResponse?.result)
                            .column_def
                    }
                    enableKeyboardNavigation
                    // @ts-ignore
                    items={getTable(
                        queryResponse?.headers,
                        queryResponse?.result
                    ).rows?.slice(page * 10, (page + 1) * 10)}
                    loading={isLoading}
                    loadingText="Loading resources"
                    // stickyColumns={{ first: 0, last: 1 }}
                    // stripedRows
                    trackBy="id"
                    empty={
                        <Box
                            margin={{
                                vertical: 'xs',
                            }}
                            textAlign="center"
                            color="inherit"
                        >
                            <SpaceBetween size="m">
                                <b>No Results</b>
                            </SpaceBetween>
                        </Box>
                    }
                    // header={
                    //     <Header className="w-full">
                    //         Results{' '}
                    //         <span className=" font-medium">
                    //             {isLoading && isExecuted
                    //                 ? '(?)'
                    //                 : `(${memoCount})`}{' '}
                    //         </span>
                    //     </Header>
                    // }
                    pagination={
                        <CustomPagination
                            currentPageIndex={page + 1}
                            pagesCount={
                                // prettier-ignore
                                (isLoading &&
                                                            isExecuted)
                                                                ? 0
                                                                : Math.ceil(
                                                                      // @ts-ignore
                                                                      getTable(
                                                                          queryResponse?.headers,
                                                                          queryResponse?.result
                                                                      ).rows
                                                                          .length /
                                                                          10
                                                                  )
                            }
                            onChange={({ detail }: any) =>
                                setPage(detail.currentPageIndex - 1)
                            }
                        />
                    }
                />
            </Grid>
        </>
    )
}
