import { Button, Card, CategoryBar, Col, Flex, Grid, Icon, Title } from '@tremor/react'
import { ChevronRightIcon, ShieldCheckIcon } from '@heroicons/react/24/outline'
import Compliance from './Compliance'
import { useNavigate, useParams } from 'react-router-dom'

export default function Governance() {
     const workspace = useParams<{ ws: string }>().ws
     const navigate = useNavigate()
    return (
        <Card className=" border-solid   border-2 border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:p-4 lg:p-4 ">
            <Flex justifyContent="between" className="sm:flex-row flex-col">
                <Flex justifyContent="start" className="gap-2 sm:w-fit w-full ">
                    <Icon icon={ShieldCheckIcon} className="p-0" />
                    <Title className="sm:font-semibold sm:w-fit w-full">
                        Compliance Frameworks
                    </Title>
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
                       see all
                    </Button>
                </a>
            </Flex>
            <Grid numItems={1} className="w-full gap-6 ">
                <Compliance />
                {/* <Col numColSpan={1}>
                    <Findings />
                </Col> */}
            </Grid>
        </Card>
    )
}
