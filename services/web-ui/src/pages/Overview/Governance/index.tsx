import { Button, Card, CategoryBar, Col, Flex, Grid, Icon, Title } from '@tremor/react'
import { ChevronRightIcon, ShieldCheckIcon } from '@heroicons/react/24/outline'
import Compliance from './Compliance'
import { useNavigate, useParams } from 'react-router-dom'

export default function Governance() {
     const workspace = useParams<{ ws: string }>().ws
     const navigate = useNavigate()
    return (
        <Card className="border-0 ring-0 !shadow-sm">
            <Flex justifyContent="between">
                <Flex justifyContent="start" className="gap-2 ">
                    <Icon icon={ShieldCheckIcon} className="p-0" />
                    <Title className="font-semibold">Compliance Frameworks</Title>
                </Flex>
                <a
                    target="__blank"
                    href={`/compliance`}
                    className=" cursor-pointer"
                >
                    <Button
                        size="xs"
                        variant="light"
                        icon={ChevronRightIcon}
                        iconPosition="right"
                        className="my-3"
                    >
                        All Compliance Frameworks
                    </Button>
                </a>
            </Flex>
            <Grid numItems={1} className="w-full gap-6 px-2">
                <Compliance />
                {/* <Col numColSpan={1}>
                    <Findings />
                </Col> */}
            </Grid>
        </Card>
    )
}
