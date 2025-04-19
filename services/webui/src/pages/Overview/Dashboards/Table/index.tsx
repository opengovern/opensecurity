import axios from 'axios'
import { useAtom } from 'jotai'
import { useEffect, useState } from 'react'
import { Box, Header, Modal, SpaceBetween, Table, Tabs } from '@cloudscape-design/components'
import CustomPagination from '../../../../components/Pagination'
import { RenderObject } from '../../../../components/RenderObject'

export default function DashboardTable({ dashboards, loading }: any) {
    const [page, setPage] = useState(1)
    const [open, setOpen] = useState(false)
    const [selectedRow, setSelectedRow] = useState<any>({})

    useEffect(() => {}, [])

    return (
        <>
            <Table
                className="mt-2"
                variant="full-page"
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
                        id: 'name',
                        header: 'Name',
                        cell: (item: any) => <>{item.name}</>,
                        sortingField: 'name',
                        isRowHeader: true,
                        maxWidth: 100,
                    },
                   
                   
                    // user_id
                    {
                        id: 'user_id',
                        header: 'User ',
                        cell: (item) => <>{item?.user_id}</>,
                        sortingField: 'status',
                        isRowHeader: true,
                        maxWidth: 100,
                    },
                    {
                        id: 'is_default',
                        header: 'Default ',
                        cell: (item) => <>{item?.is_default ? 'True' : 'False'}</>,
                        sortingField: 'status',
                        isRowHeader: true,
                        maxWidth: 100,
                    },

                    {
                        id: 'updatedAt',
                        header: 'Updated At',
                        cell: (item) => (
                            <>{`${item?.updated_at?.split('T')[0]} ${
                                item?.updated_at?.split('T')[1]?.split('.')[0]
                            } `}</>
                        ),
                        sortingField: 'updatedAt',
                        isRowHeader: true,
                        maxWidth: 100,
                    },
                ]}
                columnDisplay={[
                    { id: 'id', visible: true },
                    { id: 'name', visible: true },
                    { id: 'type', visible: true },
                    { id: 'widget_props', visible: true },
                    { id: 'user_id', visible: true },
                    { id: 'is_default', visible: true },
                    { id: 'updatedAt', visible: true },
                ]}
                loading={loading}
                // @ts-ignore
                items={
                    dashboards
                        ? dashboards.slice((page - 1) * 10, page * 10)
                        : []
                }
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
                        Results {dashboards.length ?? 0}
                    </Header>
                }
                pagination={
                    <CustomPagination
                        currentPageIndex={page}
                        pagesCount={Math.ceil(dashboards.length / 10)}
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
