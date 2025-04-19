import axios from 'axios'
import { useAtom } from 'jotai'
import { useEffect, useState } from 'react'
import { Box, Header, Modal, SpaceBetween, Table, Tabs } from '@cloudscape-design/components'
import CustomPagination from '../../../../components/Pagination'
import { RenderObject } from '../../../../components/RenderObject'

export default function WidgetsTable({widgets,loading} : any) {
    const [page, setPage] = useState(1)
    const [open, setOpen] = useState(false)
    const [selectedRow, setSelectedRow] = useState<any>({})


    useEffect(() => {
    }, [])

    return (
        <>
            <Table
                className="mt-2"
                variant='full-page'
                onRowClick={(event) => {
                    const row = event.detail.item
                    setSelectedRow(row)
                    setOpen(true)
                }}
                columnDefinitions={[
                    {
                        id: 'id',
                        header: 'Id',
                        cell: (item: any) => <>{item.id}</>,
                        // sortingField: 'id',
                        isRowHeader: true,
                        maxWidth: 50,
                    },
                    {
                        id: 'title',
                        header: 'Title',
                        cell: (item: any) => <>{item.title}</>,
                        sortingField: 'name',
                        isRowHeader: true,
                        maxWidth: 100,
                    },
                    // widget_type
                    {
                        id: 'type',
                        header: 'Type',
                        cell: (item: any) => <>{item.widget_type}</>,
                        sortingField: 'type',
                        isRowHeader: true,
                        maxWidth: 50,
                    },
                    {
                        id: 'widget_props',
                        header: 'Widget Props',
                        cell: (item) => (
                            <>{JSON.stringify(item?.widget_props)}</>
                        ),
                        sortingField: 'status',
                        isRowHeader: true,
                        maxWidth: 130,
                    },
                    // user_id
                    {
                        id: 'user_id',
                        header: 'User ',
                        cell: (item) => <>{item.user_id}</>,
                        sortingField: 'status',
                        isRowHeader: true,
                        maxWidth: 100,
                    },

                    {
                        id: 'createdAt',
                        header: 'Created At',
                        cell: (item) => (
                            <>{`${item?.created_at.split('T')[0]} ${
                                item?.created_at.split('T')[1].split('.')[0]
                            } `}</>
                        ),
                        sortingField: 'createdAt',
                        isRowHeader: true,
                        maxWidth: 100,
                    },

                    {
                        id: 'updatedAt',
                        header: 'Updated At',
                        cell: (item) => (
                            <>{`${item?.updated_at.split('T')[0]} ${
                                item?.updated_at.split('T')[1].split('.')[0]
                            } `}</>
                        ),
                        sortingField: 'updatedAt',
                        isRowHeader: true,
                        maxWidth: 100,
                    },
                ]}
                columnDisplay={[
                    { id: 'id', visible: true },
                    { id: 'title', visible: true },
                    { id: 'type', visible: true },
                    { id: 'widget_props', visible: true },
                    { id: 'user_id', visible: true },
                    { id: 'createdAt', visible: true },
                    { id: 'updatedAt', visible: true },
                ]}
                loading={loading}
                // @ts-ignore
                items={widgets ? widgets.slice((page-1) * 10, (page ) * 10) : []}
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
                    <Header className="w-full">
                        Results {widgets.length ?? 0}
                    </Header>
                }
                pagination={
                    <CustomPagination
                        currentPageIndex={page}
                        pagesCount={Math.ceil(widgets.length / 10)}
                        onChange={({ detail }: any) =>
                            setPage(detail.currentPageIndex)
                        }
                    />
                }
            />
             <Modal
                            visible={open}
                            onDismiss={() => setOpen(false)}
                            header="Query Result"
                            className="min-w-[500px]"
                            size="large"
                        >
                            <RenderObject obj={selectedRow} height={300} />
                        </Modal>
        </>
    )
}
