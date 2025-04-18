import { CheckCircleIcon, XCircleIcon } from "@heroicons/react/24/outline"
import { Badge, Flex, Text } from "@tremor/react"
import { PlatformEnginePkgComplianceApiBenchmarkControlSummary, PlatformEnginePkgComplianceApiConformanceStatus } from "../api/api"

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
    status: PlatformEnginePkgComplianceApiConformanceStatus | undefined
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
    json: PlatformEnginePkgComplianceApiBenchmarkControlSummary | undefined
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
