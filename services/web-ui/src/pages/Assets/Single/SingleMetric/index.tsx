import dayjs, { Dayjs } from 'dayjs'
import { useAtomValue, useSetAtom } from 'jotai'
import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Button,
    Callout,
    Card,
    Divider,
    Flex,
    Icon,
    List,
    ListItem,
    Select,
    SelectItem,
    Text,
    Title,
} from '@tremor/react'
import { useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { DocumentDuplicateIcon } from '@heroicons/react/24/outline'
import clipboardCopy from 'clipboard-copy'
import { highlight, languages } from 'prismjs'
import Editor from 'react-simple-code-editor'
import { RowClickedEvent } from 'ag-grid-community'
import {
    CheckCircleIcon,
    ExclamationCircleIcon,
} from '@heroicons/react/24/solid'
import {
    useInventoryApiV1QueryRunCreate,
    useInventoryApiV2AnalyticsMetricsDetail,
    useInventoryApiV2AnalyticsTrendList,
} from '../../../../api/inventory.gen'
import { isDemoAtom, notificationAtom, queryAtom } from '../../../../store'
import { resourceTrendChart } from '../../index'
import SummaryCard from '../../../../components/Cards/SummaryCard'
import { numberDisplay } from '../../../../utilities/numericDisplay'
import Table from '../../../../components/Table'
import { getTable } from '../../../Search/Query'
import { getConnectorIcon } from '../../../../components/Cards/ConnectorCard'
import { dateTimeDisplay } from '../../../../utilities/dateDisplay'
import Modal from '../../../../components/Modal'
import DrawerPanel from '../../../../components/DrawerPanel'
import { getErrorMessage } from '../../../../types/apierror'
import Tag from '../../../../components/Tag'
import Trends from '../../../../components/Trends'
import { useFilterState } from '../../../../utilities/urlstate'

interface ISingle {
    activeTimeRange: { start: Dayjs; end: Dayjs }
    metricId: string | undefined
    resourceId: string | undefined
}

export default function SingleMetric({
    activeTimeRange,
    metricId,
    resourceId,
}: ISingle) {
    const { value: selectedConnections } = useFilterState()
    const { ws, id, metric } = useParams()
    const isDemo = useAtomValue(isDemoAtom)
    const [modalData, setModalData] = useState('')
    const setNotification = useSetAtom(notificationAtom)
    const setQuery = useSetAtom(queryAtom)

    const [openDrawer, setOpenDrawer] = useState(false)
    const [selectedRow, setSelectedRow] = useState<any>(null)
    const [pageSize, setPageSize] = useState(1000)

    const query = {
        ...(selectedConnections.provider && {
            connector: [selectedConnections.provider],
        }),
        connectionId: metric
            ? [String(id).replace('account_', '')]
            : selectedConnections.connections,
        ...(metricId && { ids: [metricId] }),
        ...(resourceId && { resourceCollection: [resourceId] }),
        ...(selectedConnections.connectionGroup && {
            connectionGroup: selectedConnections.connectionGroup,
        }),
        ...(activeTimeRange.start && {
            startTime: activeTimeRange.start.unix(),
        }),
        ...(activeTimeRange.end && {
            endTime: activeTimeRange.end.unix(),
        }),
    }
    const { response: resourceTrend, isLoading: resourceTrendLoading } =
        useInventoryApiV2AnalyticsTrendList(query)
    const { response: metricDetail } = useInventoryApiV2AnalyticsMetricsDetail(
        metricId || ''
    )
    const generateQuery = () => {
        let q = ''
        if (metric) {
            q =
                metricDetail?.finderPerConnectionQuery?.replace(
                    '<CONNECTION_ID_LIST>',
                    [String(id).replace('account_', '')]
                        .map((a) => `'${a}'`)
                        .join(',')
                ) || ''
        } else if (selectedConnections.connections.length > 0) {
            q =
                metricDetail?.finderPerConnectionQuery?.replace(
                    '<CONNECTION_ID_LIST>',
                    selectedConnections.connections
                        .map((a) => `'${a}'`)
                        .join(',')
                ) || ''
        } else {
            q = metricDetail?.finderQuery || ''
        }
        return q
    }

    const {
        response: queryResponse,
        isLoading,
        isExecuted,
        error,
        sendNow,
    } = useInventoryApiV1QueryRunCreate(
        {
            page: { no: 1, size: pageSize },
            query: generateQuery(),
        },
        {},
        false
    )

    useEffect(() => {
        if (metricDetail && metricDetail.finderQuery) {
            sendNow()
        }
    }, [selectedConnections.connections, metricDetail, pageSize])

    const memoColumns = useMemo(
        () =>
            getTable(queryResponse?.headers, queryResponse?.result)
                .columns,
        [queryResponse, isDemo]
    )
    const memoCount = useMemo(
        () =>
            getTable(queryResponse?.headers, queryResponse?.result)
                .count,
        [queryResponse, isDemo]
    )

    const showTable = () => {
        return (
            !resourceId &&
            activeTimeRange.end.format('DD-MM-YYYY') ===
                dayjs().utc().format('DD-MM-YYYY')
        )
    }

    return (
        <>
            <Flex className="mb-6">
                <Flex alignItems="start" className="gap-2">
                    {getConnectorIcon(
                        metricDetail?.connectors
                            ? metricDetail?.connectors[0]
                            : undefined
                    )}
                    <Flex
                        flexDirection="col"
                        alignItems="start"
                        justifyContent="start"
                    >
                        <Title className="font-semibold whitespace-nowrap">
                            {metricDetail?.name}
                        </Title>
                        <Text>{metricDetail?.id}</Text>
                    </Flex>
                </Flex>
                <Button
                    variant="secondary"
                    onClick={() =>
                        setModalData(
                            generateQuery().replace(
                                '$IS_ALL_CONNECTIONS_QUERY',
                                'true'
                            ) || ''
                        )
                    }
                >
                    See query
                </Button>
            </Flex>
            <Modal open={!!modalData.length} onClose={() => setModalData('')}>
                <Title className="font-semibold">Metric query</Title>
                <Card className="my-4">
                    <Editor
                        onValueChange={() => 1}
                        highlight={(text) =>
                            highlight(text, languages.sql, 'sql')
                        }
                        value={modalData}
                        className="w-full bg-white dark:bg-gray-900 dark:text-gray-50 font-mono text-sm"
                        style={{
                            minHeight: '200px',
                        }}
                        placeholder="-- write your SQL query here"
                    />
                </Card>
                <Flex>
                    <Button
                        variant="light"
                        icon={DocumentDuplicateIcon}
                        iconPosition="left"
                        onClick={() =>
                            clipboardCopy(modalData).then(() =>
                                setNotification({
                                    text: 'Query copied to clipboard',
                                    type: 'info',
                                })
                            )
                        }
                    >
                        Copy
                    </Button>
                    <Flex className="w-fit gap-4">
                        <Button
                            variant="secondary"
                            onClick={() => {
                                setQuery(modalData)
                            }}
                        >
                            <Link to={`/cloudql`}>Open in Query</Link>
                        </Button>
                        <Button onClick={() => setModalData('')}>Close</Button>
                    </Flex>
                </Flex>
            </Modal>
            <Trends
                activeTimeRange={activeTimeRange}
                trend={resourceTrend}
                trendName="Resources"
                firstKPI={
                    <SummaryCard
                        title="Resource count"
                        metric={
                            resourceTrend
                                ? resourceTrend[resourceTrend.length - 1]?.count
                                : 0
                        }
                        isExact
                        metricPrev={resourceTrend ? resourceTrend[0]?.count : 0}
                        loading={resourceTrendLoading}
                        border={false}
                    />
                }
                secondKPI={
                    <SummaryCard
                        border={false}
                        title="Results in"
                        loading={resourceTrendLoading}
                        metric={
                            resourceTrend
                                ? resourceTrend[resourceTrend.length - 1]
                                      ?.totalConnectionCount
                                : 0
                        }
                        unit="Cloud accounts"
                    />
                }
                labels={resourceTrendChart(resourceTrend, 'daily').label}
                chartData={resourceTrendChart(resourceTrend, 'daily').data}
                loading={resourceTrendLoading}
            />
            <div className="mt-4">
                {showTable() ? (
                    <Table
                        title="Resource list"
                        id="metric_table"
                        loading={isLoading}
                        onGridReady={(e) => {
                            if (isLoading) {
                                e.api.showLoadingOverlay()
                            }
                        }}
                        columns={memoColumns}
                        rowData={
                            getTable(
                                queryResponse?.headers,
                                queryResponse?.result,
                                
                            ).rows
                        }
                        downloadable
                        onRowClicked={(event: RowClickedEvent) => {
                            setSelectedRow(event.data)
                            setOpenDrawer(true)
                        }}
                    >
                        <Flex
                            flexDirection="row-reverse"
                            justifyContent="between"
                            className="pl-3 gap-4"
                        >
                            <Flex
                                className="w-fit"
                                flexDirection="row"
                                alignItems="center"
                                justifyContent="start"
                            >
                                <Text className="mr-2">Maximum rows:</Text>
                                <Select
                                    enableClear={false}
                                    className="w-56"
                                    placeholder="1,000"
                                >
                                    <SelectItem
                                        value="1000"
                                        onClick={() => setPageSize(1000)}
                                    >
                                        1,000
                                    </SelectItem>
                                    <SelectItem
                                        value="3000"
                                        onClick={() => setPageSize(3000)}
                                    >
                                        3,000
                                    </SelectItem>
                                    <SelectItem
                                        value="5000"
                                        onClick={() => setPageSize(5000)}
                                    >
                                        5,000
                                    </SelectItem>
                                    <SelectItem
                                        value="10000"
                                        onClick={() => setPageSize(10000)}
                                    >
                                        10,000
                                    </SelectItem>
                                </Select>
                            </Flex>
                            {!isLoading && isExecuted && error && (
                                <Flex justifyContent="start" className="w-fit">
                                    <Icon
                                        icon={ExclamationCircleIcon}
                                        color="rose"
                                    />
                                    <Text color="rose">
                                        {getErrorMessage(error)}
                                    </Text>
                                </Flex>
                            )}
                            {!isLoading && isExecuted && queryResponse && (
                                <Flex justifyContent="start" className="w-fit">
                                    {memoCount === pageSize ? (
                                        <>
                                            <Icon
                                                icon={ExclamationCircleIcon}
                                                color="amber"
                                            />
                                            <Text color="amber">
                                                {`Row limit of ${numberDisplay(
                                                    pageSize,
                                                    0
                                                )} reached, results are truncated`}
                                            </Text>
                                        </>
                                    ) : (
                                        <>
                                            <Icon
                                                icon={CheckCircleIcon}
                                                color="emerald"
                                            />
                                            <Text color="emerald">Success</Text>
                                        </>
                                    )}
                                </Flex>
                            )}
                        </Flex>
                    </Table>
                ) : (
                    <Callout title="We only support LIVE data" color="amber">
                        To see the resource table you have to check if your end
                        date is set to today and remove all filters the you have
                        applied (like visiting this page from Resource
                        Collection).
                    </Callout>
                )}
            </div>
            <DrawerPanel
                title="Resource detail"
                open={openDrawer}
                onClose={() => {
                    setOpenDrawer(false)
                    // setSelectedRow(null)
                }}
            >
                <Accordion
                    className="w-full p-0 !rounded-none border-0"
                    defaultOpen
                >
                    <AccordionHeader className="w-full p-0 border-0">
                        <Title>Summary</Title>
                    </AccordionHeader>
                    <AccordionBody className="w-full p-0 border-0">
                        <List>
                            <ListItem className="py-6">
                                <Text>Resource ID</Text>
                                <Text className="text-gray-800 w-3/5 whitespace-pre-wrap text-end">
                                    {selectedRow?.resource_id}
                                </Text>
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>Resource type</Text>
                                <Text className="text-gray-800">
                                    {selectedRow?.resource_type}
                                </Text>
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>Cloud provider</Text>
                                <Text className="text-gray-800">
                                    {selectedRow?.connector}
                                </Text>
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>Resource name</Text>
                                <Text className="text-gray-800">
                                    {selectedRow?.name}
                                </Text>
                            </ListItem>
                        </List>
                    </AccordionBody>
                </Accordion>
                <Divider />
                <Accordion className="w-full p-0 !rounded-none border-0">
                    <AccordionHeader className="w-full p-0 border-0">
                        <Title>Details</Title>
                    </AccordionHeader>
                    <AccordionBody className="w-full p-0 border-0">
                        <List>
                            <ListItem className="py-6">
                                <Text>Tags</Text>
                                <Tag text={selectedRow?.tags || ''} />
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>Last discovered</Text>
                                <Text className="text-gray-800">
                                    {dateTimeDisplay(selectedRow?.created_at)}
                                </Text>
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>opencomply connection ID</Text>
                                <Text className="text-gray-800">
                                    {selectedRow?.connection_id}
                                </Text>
                            </ListItem>
                        </List>
                    </AccordionBody>
                </Accordion>
            </DrawerPanel>
        </>
    )
}
