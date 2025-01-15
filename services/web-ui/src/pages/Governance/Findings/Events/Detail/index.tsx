import { useAtomValue } from 'jotai'
import {
    Button,
    Card,
    Flex,
    Grid,
    List,
    ListItem,
    Tab,
    TabGroup,
    TabList,
    TabPanel,
    TabPanels,
    Text,
    Title,
} from '@tremor/react'
import ReactJson from '@microlink/react-json-view'
import { useEffect, useState } from 'react'
import { ArrowRightIcon } from '@heroicons/react/24/outline'
import { PlatformEnginePkgComplianceApiFindingEvent } from '../../../../../api/api'
import { getConnectorIcon } from '../../../../../components/Cards/ConnectorCard'
import DrawerPanel from '../../../../../components/DrawerPanel'
import SummaryCard from '../../../../../components/Cards/SummaryCard'
import { dateTimeDisplay } from '../../../../../utilities/dateDisplay'
import { severityBadge, statusBadge } from '../../../Controls'
import FindingDetail from '../../FindingsWithFailure/Detail'
import { useComplianceApiV1FindingsSingleDetail } from '../../../../../api/compliance.gen'
import { isDemoAtom } from '../../../../../store'
import { Tabs } from '@cloudscape-design/components'
import { RenderObject } from '../../../../../components/RenderObject'

interface IFindingDetail {
    event: PlatformEnginePkgComplianceApiFindingEvent | undefined
    open: boolean
    onClose: () => void
}

export default function EventDetail({ event, open, onClose }: IFindingDetail) {
    const [openFinding, setOpenFinding] = useState(false)
    const isDemo = useAtomValue(isDemoAtom)
    const { response: finding, sendNow } =
        useComplianceApiV1FindingsSingleDetail(
            event?.findingID || '',
            {},
            false
        )
    useEffect(() => {
        if (event) {
            sendNow()
        }
    }, [event])
    return (
        <>
            <Grid className="w-full gap-4 mb-6" numItems={2}>
                <SummaryCard
                    title="Account ID"
                    metric={event?.providerConnectionID}
                    blur={isDemo}
                    isString
                />
                <SummaryCard
                    title="Account Name"
                    metric={event?.providerConnectionName}
                    blur={isDemo}
                    isString
                />
                <SummaryCard
                    title="Region"
                    metric={event?.resourceLocation}
                    isString
                />
                <SummaryCard
                    title="Resource Type"
                    metric={event?.resourceType}
                    isString
                />
            </Grid>
            <Tabs
                tabs={[
                    {
                        label: 'Summary',
                        id: '0',
                        content: (
                            <>
                                <List>
                                    <ListItem className="py-6">
                                        <Text>Event ID</Text>
                                        <Text className="text-gray-800">
                                            {event?.id}
                                        </Text>
                                    </ListItem>
                                    <ListItem className="py-6">
                                        <Text>Event Date</Text>
                                        <Text className="text-gray-800">
                                            {dateTimeDisplay(
                                                event?.evaluatedAt
                                            )}
                                        </Text>
                                    </ListItem>
                                    <ListItem className="py-6">
                                        <Text>State</Text>
                                        <Flex
                                            flexDirection="row"
                                            className="h-full w-fit gap-2"
                                        >
                                            {statusBadge(
                                                event?.previousConformanceStatus
                                            )}
                                            <ArrowRightIcon className="w-5" />
                                            {statusBadge(
                                                event?.conformanceStatus
                                            )}
                                        </Flex>
                                    </ListItem>
                                    <ListItem className="py-6">
                                        <Text>Severity</Text>
                                        {severityBadge(event?.severity)}
                                    </ListItem>
                                    <ListItem className="py-6">
                                        <Text>Control ID</Text>
                                        <Text className="text-gray-800">
                                            {event?.controlID}
                                        </Text>
                                    </ListItem>
                                    {/* <ListItem className="py-6">
                                <Text>Finding</Text>
                                <Button
                                    variant="light"
                                    onClick={() => setOpenFinding(true)}
                                >
                                    Click to see finding detail
                                </Button>
                            </ListItem> */}
                                </List>
                            </>
                        ),
                    },
                    {
                        label: 'Resource Details',
                        id: '1',
                        content: (
                            <>
                                <Title className="mb-2">JSON</Title>
                                <Card className="px-1.5 py-3 mb-2">
                                    <RenderObject obj={event || {}} />
                                </Card>
                            </>
                        ),
                    },
                ]}
            />
            <TabGroup>
                <TabList>
                    <Tab>Summary</Tab>
                    <Tab>Resource Details</Tab>
                </TabList>
                <TabPanels>
                    <TabPanel>
                        <List>
                            <ListItem className="py-6">
                                <Text>Event ID</Text>
                                <Text className="text-gray-800">
                                    {event?.id}
                                </Text>
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>Event Date</Text>
                                <Text className="text-gray-800">
                                    {dateTimeDisplay(event?.evaluatedAt)}
                                </Text>
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>State</Text>
                                <Flex
                                    flexDirection="row"
                                    className="h-full w-fit gap-2"
                                >
                                    {statusBadge(
                                        event?.previousConformanceStatus
                                    )}
                                    <ArrowRightIcon className="w-5" />
                                    {statusBadge(event?.conformanceStatus)}
                                </Flex>
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>Severity</Text>
                                {severityBadge(event?.severity)}
                            </ListItem>
                            <ListItem className="py-6">
                                <Text>Control ID</Text>
                                <Text className="text-gray-800">
                                    {event?.controlID}
                                </Text>
                            </ListItem>
                            {/* <ListItem className="py-6">
                                <Text>Finding</Text>
                                <Button
                                    variant="light"
                                    onClick={() => setOpenFinding(true)}
                                >
                                    Click to see finding detail
                                </Button>
                            </ListItem> */}
                        </List>
                    </TabPanel>
                    <TabPanel>
                        <Title className="mb-2">JSON</Title>
                        <Card className="px-1.5 py-3 mb-2">
                            <RenderObject obj={event || {}} />
                        </Card>
                    </TabPanel>
                </TabPanels>
            </TabGroup>
        </>
    )
}
