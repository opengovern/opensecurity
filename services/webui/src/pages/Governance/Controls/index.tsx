import {
    Badge,
    Button,
    Card,
    Col,
    Flex,
    Grid,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeaderCell,
    TableRow,
    Text,
    Title,
} from '@tremor/react'
import { ChevronRightIcon } from '@heroicons/react/24/solid'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
    BookOpenIcon,
    CheckCircleIcon,
    CodeBracketIcon,
    Cog8ToothIcon,
    CommandLineIcon,
    XCircleIcon,
    ChevronDownIcon,
    ChevronUpIcon,
} from '@heroicons/react/24/outline'
import { useState } from 'react'
import MarkdownPreview from '@uiw/react-markdown-preview'
import { useAtomValue } from 'jotai'
import { useComplianceApiV1BenchmarksControlsDetail } from '../../../api/compliance.gen'
import Spinner from '../../../components/Spinner'
import { numberDisplay } from '../../../utilities/numericDisplay'
import AnimatedAccordion from '../../../components/AnimatedAccordion'
import { searchAtom } from '../../../utilities/urlstate'
import {
    PlatformEnginePkgComplianceApiBenchmarkControlSummary,
    PlatformEnginePkgComplianceApiConformanceStatus,
} from '../../../api/api'
import { Modal } from '@cloudscape-design/components'

interface IPolicies {
    id: string | undefined
    assignments: number
    enable? :boolean
}

export const severityBadge = (severity: any) => {
    const style = {
        color: '#fff',
        borderRadius: '8px',
        width: '64px',
    }
    if (severity) {
        if (severity === 'critical') {
            return (
                <Badge style={{ backgroundColor: '#6E120B', ...style }}>
                    Critical
                </Badge>
            )
        }
        if (severity === 'high') {
            return (
                <Badge style={{ backgroundColor: '#CA2B1D', ...style }}>
                    High
                </Badge>
            )
        }
        if (severity === 'medium') {
            return (
                <Badge style={{ backgroundColor: '#EE9235', ...style }}>
                    Medium
                </Badge>
            )
        }
        if (severity === 'low') {
            return (
                <Badge style={{ backgroundColor: '#F4C744', ...style }}>
                    Low
                </Badge>
            )
        }
        if (severity === 'none') {
            return (
                <Badge style={{ backgroundColor: '#9BA2AE', ...style }}>
                    None
                </Badge>
            )
        }
        return (
            <Badge style={{ backgroundColor: '#54B584', ...style }}>
                Passed
            </Badge>
        )
    }
    return <Badge style={{ backgroundColor: '#9BA2AE', ...style }}>None</Badge>
}

export const activeBadge = (status: boolean) => {
    if (status) {
        return (
            <Flex className="w-fit gap-1.5">
                <CheckCircleIcon className="h-4 text-emerald-500" />
                <Text>Active</Text>
            </Flex>
        )
    }
    return (
        <Flex className="w-fit gap-1.5">
            <XCircleIcon className="h-4 text-rose-600" />
            <Text>Inactive</Text>
        </Flex>
    )
}

export const statusBadge = (
    status:
        | PlatformEnginePkgComplianceApiConformanceStatus
        | undefined
) => {
    if (
        status ===
        PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed
    ) {
        return (
            <Flex className="w-fit gap-1.5">
                <CheckCircleIcon className="h-4 text-emerald-500" />
                <Text>Passed</Text>
            </Flex>
        )
    }
    return (
        <Flex className="w-fit gap-1.5">
            <XCircleIcon className="h-4 text-rose-600" />
            <Text>Failed</Text>
        </Flex>
    )
}

export const treeRows = (
    json:
        | PlatformEnginePkgComplianceApiBenchmarkControlSummary
        | undefined
) => {
    let arr: any = []
    if (json) {
        if (json.control !== null && json.control !== undefined) {
            for (let i = 0; i < json.control.length; i += 1) {
                let obj = {}
                obj = {
                    parentName: json?.benchmark?.title,
                    ...json.control[i].control,
                    ...json.control[i],
                }
                arr.push(obj)
            }
        }
        if (json.children !== null && json.children !== undefined) {
            for (let i = 0; i < json.children.length; i += 1) {
                const res = treeRows(json.children[i])
                arr = arr.concat(res)
            }
        }
    }

    return arr
}

export const groupBy = (input: any[], key: string) => {
    return input.reduce((acc, currentValue) => {
        const groupKey = currentValue[key]
        if (!acc[groupKey]) {
            acc[groupKey] = []
        }
        acc[groupKey].push(currentValue)
        return acc
    }, {})
}

export const countControls = (
    v:
        | PlatformEnginePkgComplianceApiBenchmarkControlSummary
        | undefined
) => {
    const countChildren = v?.children
        ?.map((i) => countControls(i))
        .reduce((prev, curr) => prev + curr, 0)
    const total: number = (countChildren || 0) + (v?.control?.length || 0)
    return total
}

export default function Controls({ id, assignments, enable }: IPolicies) {
    const { response: controls, isLoading } =
        useComplianceApiV1BenchmarksControlsDetail(String(id))
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const [doc, setDoc] = useState('')
    const [docTitle, setDocTitle] = useState('')
    const [openAllControls, setOpenAllControls] = useState(false)

    const toggleOpen = () => {
        setOpenAllControls(!openAllControls)
    }

    const countBenchmarks = (
        v:
            | PlatformEnginePkgComplianceApiBenchmarkControlSummary
            | undefined
    ) => {
        const countChildren = v?.children
            ?.map((i) => countBenchmarks(i))
            .reduce((prev, curr) => prev + curr, 0)
        const total: number = (countChildren || 0) + (v?.children?.length || 0)
        return total
    }

    const sections = Object.entries(groupBy(treeRows(controls), 'parentName'))

    return (
        <Flex flexDirection="col" className="gap-4">
            <Flex
                className="w-full"
                flexDirection="row"
                justifyContent="between"
            >
                {/* <Title>Controls ({countControls(controls)})</Title> */}
                {isLoading ? (
                    ''
                ) : (
                    <Button
                        variant="light"
                        icon={openAllControls ? ChevronUpIcon : ChevronDownIcon}
                        onClick={toggleOpen}
                    >
                        {openAllControls
                            ? 'Collapse all'
                            : `Expand ${sections.length} section${
                                  sections.length > 1 ? 's' : ''
                              }`}
                    </Button>
                )}
            </Flex>

            <Modal
                header={docTitle}
                visible={doc.length > 0}
                onDismiss={() => setDoc('')}
            >
                <MarkdownPreview
                    source={doc}
                    className="!bg-transparent"
                    wrapperElement={{
                        'data-color-mode': 'light',
                    }}
                    rehypeRewrite={(node, index, parent) => {
                        if (
                            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                            // @ts-ignore
                            node.tagName === 'a' &&
                            parent &&
                            /^h(1|2|3|4|5|6)/.test(
                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                // @ts-ignore
                                parent.tagName
                            )
                        ) {
                            // eslint-disable-next-line no-param-reassign
                            parent.children = parent.children.slice(1)
                        }
                    }}
                />
            </Modal>
            {isLoading ? (
                <Spinner className="mt-20" />
            ) : (
                Object.entries(groupBy(treeRows(controls), 'parentName'))?.map(
                    ([name, value]: any[]) => (
                        <Card>
                            <AnimatedAccordion
                                header={
                                    <Flex className="pl-6">
                                        <Flex>
                                            <Title className="font-semibold">
                                                {name}
                                            </Title>
                                            {enable !== false && (
                                                <>
                                                    {assignments > 0 && (
                                                        <Flex
                                                            justifyContent="start"
                                                            className="gap-2"
                                                            style={{
                                                                width: '230px',
                                                            }}
                                                        >
                                                            {value?.filter(
                                                                (c: any) =>
                                                                    c.passed
                                                            ).length ===
                                                            value?.length ? (
                                                                <CheckCircleIcon className="w-5 min-w-[20px] text-emerald-500" />
                                                            ) : (
                                                                <XCircleIcon className="w-5 min-w-[20px] text-rose-600" />
                                                            )}
                                                            <Text className="font-semibold whitespace-nowrap">{`Passed controls: ${numberDisplay(
                                                                value?.filter(
                                                                    (c: any) =>
                                                                        c.passed
                                                                ).length,
                                                                0
                                                            )}/${numberDisplay(
                                                                value?.length,
                                                                0
                                                            )} (${Math.floor(
                                                                // eslint-disable-next-line no-unsafe-optional-chaining
                                                                (value?.filter(
                                                                    (c: any) =>
                                                                        c.passed
                                                                ).length /
                                                                    (value?.length ||
                                                                        0)) *
                                                                    100
                                                            )}%)`}</Text>
                                                        </Flex>
                                                    )}
                                                </>
                                            )}
                                        </Flex>
                                    </Flex>
                                }
                                defaultOpen={openAllControls}
                            >
                                <Table className="max-w-full">
                                    <TableHead className="max-w-full">
                                        <TableRow className="max-w-full">
                                            <TableHeaderCell className="w-24">
                                                Control
                                            </TableHeaderCell>
                                            <TableHeaderCell>
                                                Title
                                            </TableHeaderCell>
                                            {enable !== false && (
                                                <>
                                                    <TableHeaderCell className="w-40">
                                                        Remediation
                                                    </TableHeaderCell>
                                                    {assignments > 0 && (
                                                        <TableHeaderCell className="w-48">
                                                            Passed resources
                                                        </TableHeaderCell>
                                                    )}
                                                </>
                                            )}

                                            <TableHeaderCell className="w-5" />
                                        </TableRow>
                                    </TableHead>
                                    <TableBody className="max-w-full">
                                        {value.map((v: any, i: number) => (
                                            <TableRow
                                                className="max-w-full cursor-pointer hover:bg-openg-50 dark:hover:bg-gray-900"
                                                key={v?.id}
                                            >
                                                <TableCell
                                                    className="w-24 min-w-[96px]"
                                                    onClick={() =>
                                                        navigate(
                                                            `${String(
                                                                v?.id
                                                            )}?${searchParams}`
                                                        )
                                                    }
                                                >{`${name.substring(
                                                    0,
                                                    name.indexOf(' ')
                                                )}.${i + 1}`}</TableCell>
                                                <TableCell
                                                    onClick={() =>
                                                        navigate(
                                                            `${String(
                                                                v?.id
                                                            )}?${searchParams}`
                                                        )
                                                    }
                                                >
                                                    <Grid numItems={12}>
                                                        <Col numColSpan={2}>
                                                            {severityBadge(
                                                                v?.severity
                                                            )}
                                                        </Col>
                                                        <Col
                                                            numColSpan={10}
                                                            className="-ml-8"
                                                        >
                                                            <Text className="truncate">
                                                                {v?.title}
                                                            </Text>
                                                        </Col>
                                                    </Grid>
                                                </TableCell>
                                                {enable !== false && (
                                                    <>
                                                        <TableCell>
                                                            <Flex
                                                                justifyContent="start"
                                                                className="gap-1.5"
                                                            >
                                                                {v?.cliRemediation &&
                                                                    v
                                                                        ?.cliRemediation
                                                                        .length >
                                                                        0 && (
                                                                        <div className="group relative flex justify-center">
                                                                            <CommandLineIcon
                                                                                className="text-openg-500 w-5"
                                                                                onClick={() => {
                                                                                    setDoc(
                                                                                        v?.cliRemediation
                                                                                    )
                                                                                    setDocTitle(
                                                                                        `Command line (CLI) remediation for '${v?.title}'`
                                                                                    )
                                                                                }}
                                                                            />
                                                                            <Card className="absolute -top-2.5 left-6 w-fit z-40 scale-0 transition-all rounded p-2 group-hover:scale-100">
                                                                                <Text>
                                                                                    Command
                                                                                    line
                                                                                    (CLI)
                                                                                </Text>
                                                                            </Card>
                                                                        </div>
                                                                    )}
                                                                {v?.manualRemediation &&
                                                                    v
                                                                        ?.manualRemediation
                                                                        .length >
                                                                        0 && (
                                                                        <div className="group relative flex justify-center">
                                                                            <BookOpenIcon
                                                                                className="text-openg-500 w-5"
                                                                                onClick={() => {
                                                                                    setDoc(
                                                                                        v?.manualRemediation
                                                                                    )
                                                                                    setDocTitle(
                                                                                        `Manual remediation for '${v?.title}'`
                                                                                    )
                                                                                }}
                                                                            />
                                                                            <Card className="absolute -top-2.5 left-6 w-fit z-40 scale-0 transition-all rounded p-2 group-hover:scale-100">
                                                                                <Text>
                                                                                    Manual
                                                                                </Text>
                                                                            </Card>
                                                                        </div>
                                                                    )}
                                                                {v?.programmaticRemediation &&
                                                                    v
                                                                        ?.programmaticRemediation
                                                                        .length >
                                                                        0 && (
                                                                        <div className="group relative flex justify-center">
                                                                            <CodeBracketIcon
                                                                                className="text-openg-500 w-5"
                                                                                onClick={() => {
                                                                                    setDoc(
                                                                                        v?.programmaticRemediation
                                                                                    )
                                                                                    setDocTitle(
                                                                                        `Programmatic remediation for '${v?.title}'`
                                                                                    )
                                                                                }}
                                                                            />
                                                                            <Card className="absolute -top-2.5 left-6 w-fit z-40 scale-0 transition-all rounded p-2 group-hover:scale-100">
                                                                                <Text>
                                                                                    Programmatic
                                                                                </Text>
                                                                            </Card>
                                                                        </div>
                                                                    )}
                                                                {v?.guardrailRemediation &&
                                                                    v
                                                                        ?.guardrailRemediation
                                                                        .length >
                                                                        0 && (
                                                                        <div className="group relative flex justify-center">
                                                                            <Cog8ToothIcon
                                                                                className="text-openg-500 w-5"
                                                                                onClick={() => {
                                                                                    setDoc(
                                                                                        v?.guardrailRemediation
                                                                                    )
                                                                                    setDocTitle(
                                                                                        `Guard rails remediation for '${v?.title}'`
                                                                                    )
                                                                                }}
                                                                            />
                                                                            <Card className="absolute -top-2.5 left-6 w-fit z-40 scale-0 transition-all rounded p-2 group-hover:scale-100">
                                                                                <Text>
                                                                                    Guard
                                                                                    rail
                                                                                </Text>
                                                                            </Card>
                                                                        </div>
                                                                    )}
                                                            </Flex>
                                                        </TableCell>
                                                        {assignments > 0 && (
                                                            <TableCell
                                                                onClick={() =>
                                                                    navigate(
                                                                        `${String(
                                                                            v?.id
                                                                        )}?${searchParams}`
                                                                    )
                                                                }
                                                            >
                                                                <Flex
                                                                    justifyContent="start"
                                                                    className="gap-2"
                                                                >
                                                                    {(v?.totalResourcesCount ||
                                                                        0) -
                                                                        (v?.failedResourcesCount ||
                                                                            0) ===
                                                                    (v?.totalResourcesCount ||
                                                                        0) ? (
                                                                        <CheckCircleIcon className="w-5 min-w-[20px] text-emerald-500" />
                                                                    ) : (
                                                                        <XCircleIcon className="w-5 min-w-[20px] text-rose-600" />
                                                                    )}
                                                                    {`${numberDisplay(
                                                                        (v?.totalResourcesCount ||
                                                                            0) -
                                                                            (v?.failedResourcesCount ||
                                                                                0),
                                                                        0
                                                                    )}/${numberDisplay(
                                                                        v?.totalResourcesCount ||
                                                                            0,
                                                                        0
                                                                    )} (${Math.floor(
                                                                        (((v?.totalResourcesCount ||
                                                                            0) -
                                                                            (v?.failedResourcesCount ||
                                                                                0)) /
                                                                            (v?.totalResourcesCount ||
                                                                                1)) *
                                                                            100
                                                                    )}%)`}
                                                                </Flex>
                                                            </TableCell>
                                                        )}
                                                    </>
                                                )}

                                                <TableCell
                                                    onClick={() =>
                                                        navigate(
                                                            `${String(
                                                                v?.id
                                                            )}?${searchParams}`
                                                        )
                                                    }
                                                >
                                                    <ChevronRightIcon className="h-5 text-openg-500" />
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </AnimatedAccordion>
                        </Card>
                    )
                )
            )}
        </Flex>
    )
}
