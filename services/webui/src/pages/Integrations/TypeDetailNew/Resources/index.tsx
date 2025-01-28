import { Card, Flex, Title } from '@tremor/react'

import { useEffect, useState } from 'react'

import axios from 'axios'

import { AppLayout, Box, Button, Header, KeyValuePairs, Pagination, SpaceBetween, SplitPanel, Table } from '@cloudscape-design/components'
import Spinner from '../../../../components/Spinner'
import ReactMarkdown from 'react-markdown'
import rehypeRaw from 'rehype-raw'
import { useMDXComponents } from '../../../../components/MDX'
interface IntegrationListProps {
    name?: string
    integration_type?: string
}

export default function Resources({
    name,
    integration_type,
}: IntegrationListProps) {
  
    const [loading, setLoading] = useState<boolean>(false)
    const [resourceLoading, setResourceLoading] = useState<boolean>(false)
   
const [resourceTypes, setResourceTypes] = useState<any>([])
const [total_count, setTotalCount] = useState<number>(0)
    const [page, setPage] = useState(0)
    const[selected,setSelected]= useState<any>();
    const [open,setOpen] = useState(false)

      const GetResourceTypes = () => {
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

          setResourceLoading(true)
          axios
              .get(
                  `${url}/main/integration/api/v1/integrations/types/${integration_type}/resource_types`,

                  config
              )
              .then((res) => {
                  const data = res.data
                  setResourceTypes(data?.integration_types)
                  setTotalCount(data?.total_count)

          setResourceLoading(false)
        
              })
              .catch((err) => {
                  console.log(err)
                            setResourceLoading(false)

              })
      }


    useEffect(() => {
    
        GetResourceTypes()
    }, [])
                           


    return (
        <>
            {loading ? (
                <>
                    <Spinner />
                </>
            ) : (
                <AppLayout
                    toolsOpen={false}
                    navigationOpen={false}
                    contentType="table"
                    toolsHide={true}
                    navigationHide={true}
                    splitPanelOpen={open}
                    onSplitPanelToggle={() => {
                        setOpen(!open)
                    }}
                    splitPanel={
                        <SplitPanel
                            // @ts-ignore
                            header={selected?.name ? selected?.name : ''}
                        >
                            <Flex className="flex-col gap-3 justify-start items-start">
                                <Title> Parameters: </Title>
                                <KeyValuePairs
                                    columns={2}
                                    // @ts-ignore
                                    items={
                                        selected
                                            ? selected?.params?.map(
                                                  (param: any) => {
                                                      return {
                                                          label: param?.name,
                                                          value: param?.description,
                                                      }
                                                  }
                                              )
                                            : []
                                    }
                                />
                            </Flex>
                        </SplitPanel>
                    }
                    content={
                        <Table
                            className="  min-h-[450px]"
                            variant="full-page"
                            // resizableColumns
                            renderAriaLive={({
                                firstIndex,
                                lastIndex,
                                totalItemsCount,
                            }) =>
                                `Displaying items ${firstIndex} to ${lastIndex} of ${totalItemsCount}`
                            }
                            onRowClick={(event) => {
                                const row = event.detail.item
                                // @ts-ignore
                                if (row.params.length > 0) {
                                    setSelected(row)
                                    setOpen(true)
                                }
                            }}
                            columnDefinitions={[
                                {
                                    id: 'name',
                                    header: 'Name',
                                    cell: (item: any) => item.name,
                                },
                                {
                                    id: 'table',
                                    header: 'Table',
                                    cell: (item: any) => item.table,
                                },
                                {
                                    id: 'has_parameter',
                                    header: 'Parameterized',
                                    cell: (item: any) =>
                                        item.params.length > 0 ? 'Yes' : 'No',
                                },
                            ]}
                            columnDisplay={[
                                {
                                    id: 'name',
                                    visible: true,
                                },
                                {
                                    id: 'table',
                                    visible: true,
                                },
                                {
                                    id: 'has_parameter',
                                    visible: true,
                                },
                            ]}
                            enableKeyboardNavigation
                            // @ts-ignore
                            items={resourceTypes?.slice(
                                page * 15,
                                (page + 1) * 15
                            )}
                            loading={resourceLoading}
                            loadingText="Loading resources"
                            // stickyColumns={{ first: 0, last: 1 }}
                            // stripedRows
                            trackBy="id"
                            empty={
                                <Box
                                    margin={{ vertical: 'xs' }}
                                    textAlign="center"
                                    color="inherit"
                                >
                                    <SpaceBetween size="m">
                                        <b>No resources</b>
                                    </SpaceBetween>
                                </Box>
                            }
                            header={
                                <Header
                                    className="w-full"
                                    actions={
                                        <Pagination
                                            // @ts-ignore
                                            className="min-w-fit"
                                            currentPageIndex={page + 1}
                                            pagesCount={Math.ceil(
                                                total_count / 15
                                            )}
                                            onChange={({ detail }) =>
                                                setPage(
                                                    detail.currentPageIndex - 1
                                                )
                                            }
                                        />
                                    }
                                >
                                    <Flex className="flex-row justify-between items-center w-full resource-types ">
                                        <div className="w-full">
                                            {' '}
                                            Resources{' '}
                                            <span className=" font-medium">
                                                ({total_count})
                                            </span>
                                        </div>
                                    </Flex>
                                </Header>
                            }
                            // pagination={

                            // }
                        />
                    }
                />
            )}
        </>
    )
}
