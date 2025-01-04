import { Flex, Text } from '@tremor/react'
import { useNavigate } from 'react-router-dom'
import { ICellRendererParams, RowClickedEvent } from 'ag-grid-community'
import { useAtomValue } from 'jotai'
import { useComplianceApiV1FindingsTopDetail } from '../../../../../../api/compliance.gen'
import {
    GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../../../../../api/api'
import Table, { IColumn } from '../../../../../../components/Table'
import { topControls } from '../../TopDetails/Controls'
import { severityBadge } from '../../../../Controls'
import { DateRange, searchAtom } from '../../../../../../utilities/urlstate'

const policyColumns: IColumn<any, any>[] = [
    {
        headerName: 'Control',
        field: 'title',
        type: 'string',
        sortable: true,
        filter: true,
        resizable: true,
        cellRenderer: (param: ICellRendererParams) => (
            <Flex
                flexDirection="col"
                alignItems="start"
                justifyContent="center"
                className="h-full"
            >
                <Text className="text-gray-800">{param.value}</Text>
                <Text>{param.data.id}</Text>
            </Flex>
        ),
    },
    {
        headerName: 'Severity',
        field: 'sev',
        width: 120,
        type: 'string',
        sortable: true,
        filter: true,
        resizable: true,
        cellRenderer: (params: ICellRendererParams) => (
            <Flex
                className="h-full w-full"
                justifyContent="center"
                alignItems="center"
            >
                {severityBadge(params.data.severity)}
            </Flex>
        ),
    },
    {
        headerName: 'Findings',
        field: 'count',
        type: 'number',
        sortable: true,
        filter: true,
        resizable: true,
        width: 150,
        cellRenderer: (param: ICellRendererParams) => (
            <Flex
                flexDirection="col"
                alignItems="start"
                justifyContent="center"
                className="h-full"
            >
                <Text className="text-gray-800">{param.value || 0} issues</Text>
                <Text>
                    {(param.data.totalCount || 0) - (param.value || 0)} passed
                </Text>
            </Flex>
        ),
    },
    {
        headerName: 'Resources',
        field: 'resourceCount',
        type: 'number',
        sortable: true,
        filter: true,
        resizable: true,
        width: 150,
        cellRenderer: (param: ICellRendererParams) => (
            <Flex
                flexDirection="col"
                alignItems="start"
                justifyContent="center"
                className="h-full"
            >
                <Text className="text-gray-800">{param.value || 0} issues</Text>
                <Text>
                    {(param.data.resourceTotalCount || 0) - (param.value || 0)}{' '}
                    passed
                </Text>
            </Flex>
        ),
    },
]

interface ICount {
    query: {
        connector: SourceType
        conformanceStatus:
            | GithubComKaytuIoKaytuEnginePkgComplianceApiConformanceStatus[]
            | undefined
        severity: TypesFindingSeverity[] | undefined
        connectionID: string[] | undefined
        controlID: string[] | undefined
        benchmarkID: string[] | undefined
        resourceTypeID: string[] | undefined
        lifecycle: boolean[] | undefined
        activeTimeRange: DateRange | undefined
    }
}

export default function ControlsWithFailure({ query }: ICount) {
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const topQuery = {
        integrationType: query.connector.length ? [query.connector] : [],
        integrationID: query.connectionID,
        benchmarkId: query.benchmarkID,
    }

    const { response: controls, isLoading } =
        useComplianceApiV1FindingsTopDetail('controlID', 10000, topQuery)

    return (
        <Table
            id="compliance_policies"
            loading={isLoading}
            onGridReady={(e) => {
                if (isLoading) {
                    e.api.showLoadingOverlay()
                }
            }}
            columns={policyColumns}
            rowData={topControls(controls?.records)}
            onRowClicked={(event: RowClickedEvent) => {
                if (event.data) {
                    navigate(`${event.data.id}?${searchParams}`)
                }
            }}
            rowHeight="lg"
        />
    )
}
