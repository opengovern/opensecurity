// @ts-nocheck
import {
    Card,
    Col,
    Flex,
    Grid,
    Icon,
    ProgressCircle,
    Title,
} from '@tremor/react'
import { useEffect, useState } from 'react'
import {
    DocumentTextIcon,
    PuzzlePieceIcon,
    ShieldCheckIcon,
} from '@heroicons/react/24/outline'
import { useComplianceApiV1BenchmarksSummaryList } from '../../../api/compliance.gen'
import {
    GithubComKaytuIoKaytuEnginePkgComplianceApiBenchmarkEvaluationSummary,
    SourceType,
} from '../../../api/api'
import ComplianceListCard from '../../../components/Cards/ComplianceListCard'
import TopHeader from '../../../components/Layout/Header'
import FilterGroup, { IFilter } from '../../../components/FilterGroup'
import { useURLParam, useURLState } from '../../../utilities/urlstate'
import {
    BenchmarkStateFilter,
    ConnectorFilter,
} from '../../../components/FilterGroup/FilterTypes'
import { errorHandling } from '../../../types/apierror'
import RadioSelector, {
    RadioItem,
} from '../../../components/FilterGroup/RadioSelector'
import { benchmarkChecks } from '../../../components/Cards/ComplianceCard'
import Spinner from '../../../components/Spinner'
import axios from 'axios'
import BenchmarkCard from './BenchmarkCard'
import BenchmarkCards from './BenchmarkCard'
import {
    Header,
    Pagination,
    PropertyFilter,
    Tabs,
} from '@cloudscape-design/components'
import Multiselect from '@cloudscape-design/components/multiselect'
import Select from '@cloudscape-design/components/select'
import ScoreCategoryCard from '../../../components/Cards/ScoreCategoryCard'
import AllControls from './All Controls'
import SettingsParameters from '../../Settings/Parameters'
import { useIntegrationApiV1EnabledConnectorsList } from '../../../api/integration.gen'
import AllPolicy from './All Policy'
const CATEGORY = {
    sre_efficiency: 'Efficiency',
    sre_reliability: 'Reliability',
    sre_supportability: 'Supportability',
}

export default function Compliance() {
    const defaultSelectedConnectors = ''

    const [loading, setLoading] = useState<boolean>(false)
    const [query, setQuery] = useState({
        tokens: [],
        operation: 'and',
    })
    const [connectors, setConnectors] = useState({
        label: 'Any',
        value: 'Any',
    })
    const [enable, setEnanble] = useState({
        label: 'No',
        value: false,
    })
    const [isSRE, setIsSRE] = useState({
        label: 'Compliance Benchmark',
        value: false,
    })

    const [AllBenchmarks, setBenchmarks] = useState()
    const [BenchmarkDetails, setBenchmarksDetails] = useState()
    const [page, setPage] = useState<number>(1)
    const [totalPage, setTotalPage] = useState<number>(0)
    const [totalCount, setTotalCount] = useState<number>(0)
    const [response, setResponse] = useState()
    const [isLoading, setIsLoading] = useState(false)
const {
    response: Types,
    isLoading: TypesLoading,
    isExecuted: TypesExec,
} = useIntegrationApiV1EnabledConnectorsList(0, 0)

    const getFilterOptions =() =>{
        const temp = [
          
            {
                propertyKey: 'enable',
                value: 'Yes',
            },
            {
                propertyKey: 'enable',
                value: 'No',
            },
          
        ]
        Types?.integration_types?.map((item) => {
            temp.push({
                propertyKey: 'integrationType',
                value: item.platform_name,
            })

        })


        return temp

    }
    const GetCard = () => {
        let url = ''
        setLoading(true)
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
        const connectors = []
        const enable = []
        const isSRE = []
        const title =[]
        query.tokens.map((item) => {
            if (item.propertyKey == 'integrationType') {
                connectors.push(item.value)
            }
            if (item.propertyKey == 'enable') {
                enable.push(item.value)
            }
            if (item.propertyKey == 'title_regex') {
                title.push(item.value)
            }
            
            // if(item.propertyKey == 'family'){
            //     isSRE.push(item.value)
            // }
        })
        const connector_filter = connectors.length == 1 ? connectors : []

        let sre_filter = false
        if (isSRE.length == 1) {
            if (isSRE[0] == 'SRE benchmark') {
                sre_filter = true
            }
        }

        let enable_filter = true
        if (enable.length == 1) {
            if (enable[0] == 'No') {
                enable_filter = false
            }
        }


        const body = {
            cursor: page,
            per_page: 6,
            sort_by: 'incidents',
            assigned: false,
            is_baseline: sre_filter,
            integrationType: connector_filter,
            root: true,
            title_regex: title[0],
        }

        axios
            .post(`${url}/main/compliance/api/v3/benchmarks`, body, config)
            .then((res) => {
                //  const temp = []
                if (!res.data.items) {
                    setLoading(false)
                }
                setBenchmarks(res.data.items)
                setTotalPage(Math.ceil(res.data.total_count / 6))
                setTotalCount(res.data.total_count)
            })
            .catch((err) => {
                setLoading(false)
                setBenchmarks([])

                console.log(err)
            })
    }

    const Detail = (benchmarks: string[]) => {
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
        const body = {
            benchmarks: benchmarks,
        }
        axios
            .post(
                `${url}/main/compliance/api/v3/compliance/summary/benchmark`,
                body,
                config
            )
            .then((res) => {
                //  const temp = []
                setLoading(false)
                setBenchmarksDetails(res.data)
            })
            .catch((err) => {
                setLoading(false)
                setBenchmarksDetails([])

                console.log(err)
            })
    }
    const GetBenchmarks = (benchmarks: string[]) => {
        setIsLoading(true)
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
        const body = {
            benchmarks: benchmarks,
        }
        axios
            .post(
                `${url}/main/compliance/api/v3/compliance/summary/benchmark`,
                body,
                config
            )
            .then((res) => {
                const temp = [
                    // {
                    //     benchmark_id: 'sre_supportability',
                    //     benchmark_title: 'SRE Security',
                    //     compliance_score: 0.6666666666666666,
                    //     connectors: ['Azure', 'AWS'],
                    //     severity_summary_by_control: {
                    //         total: {
                    //             total: 12,
                    //             passed: 8,
                    //             failed: 4,
                    //         },
                    //         critical: {
                    //             total: 0,
                    //             passed: 0,
                    //             failed: 0,
                    //         },
                    //         high: {
                    //             total: 7,
                    //             passed: 4,
                    //             failed: 3,
                    //         },
                    //         medium: {
                    //             total: 3,
                    //             passed: 2,
                    //             failed: 1,
                    //         },
                    //         low: {
                    //             total: 2,
                    //             passed: 2,
                    //             failed: 0,
                    //         },
                    //         none: {
                    //             total: 0,
                    //             passed: 0,
                    //             failed: 0,
                    //         },
                    //     },
                    //     severity_summary_by_resource: {
                    //         total: {
                    //             total: 19,
                    //             passed: 13,
                    //             failed: 6,
                    //         },
                    //         critical: {
                    //             total: 0,
                    //             passed: 0,
                    //             failed: 0,
                    //         },
                    //         high: {
                    //             total: 11,
                    //             passed: 7,
                    //             failed: 4,
                    //         },
                    //         medium: {
                    //             total: 6,
                    //             passed: 4,
                    //             failed: 2,
                    //         },
                    //         low: {
                    //             total: 2,
                    //             passed: 2,
                    //             failed: 0,
                    //         },
                    //         none: {
                    //             total: 0,
                    //             passed: 0,
                    //             failed: 0,
                    //         },
                    //     },
                    //     severity_summary_by_incidents: {
                    //         none: 0,
                    //         low: 0,
                    //         medium: 4,
                    //         high: 8,
                    //         critical: 0,
                    //         total: 12,
                    //     },
                    //     cost_optimization: 0,
                    //     findings_summary: {
                    //         total_count: 30,
                    //         passed: 18,
                    //         failed: 12,
                    //     },
                    //     issues_count: 12,
                    //     top_integrations: [
                    //         {
                    //             integration_info: {
                    //                 integration: 'Azure',
                    //                 type: 'azure_subscription',
                    //                 id: '75b0a9a9-3222-4290-bdf9-56127d550563',
                    //                 id_name: 'Policy Testing Subscription',
                    //                 integration_tracker:
                    //                     '1c2a6b18-ac87-4f5e-a472-1e26f8704f29',
                    //             },
                    //             issues: 4,
                    //         },
                    //         {
                    //             integration_info: {
                    //                 integration: 'AWS',
                    //                 type: 'aws_account',
                    //                 id: '861370837605',
                    //                 id_name: 'ADorigi',
                    //                 integration_tracker:
                    //                     'e6cb0afa-e624-4ca7-8b47-fa9988831137',
                    //             },
                    //             issues: 4,
                    //         },
                    //     ],
                    //     top_resources_with_issues: [
                    //         {
                    //             field: 'Resource',
                    //             key: 'o-ng68d511a2',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'Resource',
                    //             key: '861370837605',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'Resource',
                    //             key: '/subscriptions/75b0a9a9-3222-4290-bdf9-56127d550563/resourceGroups/policy-testing-us-east/providers/Microsoft.KeyVault/vaults/vm-policy-test-keyvault2',
                    //             issues: 1,
                    //         },
                    //         {
                    //             field: 'Resource',
                    //             key: '/subscriptions/75b0a9a9-3222-4290-bdf9-56127d550563/resourceGroups/policy-testing-us-east/providers/Microsoft.KeyVault/vaults/vm-policy-test-keyvault3',
                    //             issues: 1,
                    //         },
                    //         {
                    //             field: 'Resource',
                    //             key: '/subscriptions/75b0a9a9-3222-4290-bdf9-56127d550563/resourceGroups/policy-testing-us-east/providers/Microsoft.Storage/storageAccounts/kaytustorageaccounttest',
                    //             issues: 1,
                    //         },
                    //     ],
                    //     top_resource_types_with_issues: [
                    //         {
                    //             field: 'ResourceType',
                    //             key: 'aws::account::account',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'ResourceType',
                    //             key: 'microsoft.keyvault/vaults',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'ResourceType',
                    //             key: 'microsoft.storage/storageaccounts',
                    //             issues: 2,
                    //         },
                    //     ],
                    //     top_controls_with_issues: [
                    //         {
                    //             field: 'Control',
                    //             key: 'aws_cis_v120_1_11',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'Control',
                    //             key: 'aws_cis_v120_1_8',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'Control',
                    //             key: 'aws_cis_v130_1_8',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'Control',
                    //             key: 'aws_account_alternate_contact_security_registered',
                    //             issues: 2,
                    //         },
                    //         {
                    //             field: 'Control',
                    //             key: 'aws_account_alternate_contacts',
                    //             issues: 2,
                    //         },
                    //     ],
                    //     last_evaluated_at: '2024-10-06T12:21:46Z',
                    //     last_job_status: 'SUCCEEDED',
                    //     last_job_id: '70',
                    // },
                ]
                setIsLoading(false)
                res.data?.map((item) => {
                    temp.push(item)
                })
                setResponse(temp)
            })
            .catch((err) => {
                setIsLoading(false)

                console.log(err)
            })
    }
    useEffect(() => {
        GetCard()
    }, [page, query])

    useEffect(() => {
        if (AllBenchmarks) {
            const temp = []
            AllBenchmarks?.map((item) => {
                temp.push(item.benchmark.id)
            })
            Detail(temp)
        }
    }, [AllBenchmarks])
    useEffect(() => {
        GetBenchmarks([
            'baseline_efficiency',
            'baseline_reliability',
            'baseline_security',
            'baseline_supportability',
        ])
    }, [])

    return (
        <>
            {/* <TopHeader /> */}
            <Tabs
                tabs={[
                    {
                        label: 'Frameworks',
                        id: '0',
                        content: (
                            <>
                                <Flex
                                    className="bg-white w-full rounded-xl border-solid  border-2 border-gray-200   "
                                    flexDirection="col"
                                    justifyContent="center"
                                    alignItems="center"
                                >
                                    <div className="border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:p-6 lg:p-8">
                                        <header>
                                            <h1 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                                Frameworks
                                            </h1>
                                            <p className="text-tremor-default text-tremor-content dark:text-dark-tremor-content">
                                                Assign, Audit, and govern your
                                                tech stack with Compliance
                                                Frameworks.
                                            </p>
                                            <Grid
                                                numItems={
                                                    window.innerWidth > 1440
                                                        ? 4
                                                        : 2
                                                }
                                                className="2xl:gap-[30px] sm: gap-10 mt-6 w-full justify-items-center"
                                            >
                                                {isLoading || !response
                                                    ? [1, 2, 3, 4].map((i) => (
                                                          <Flex className="gap-6 2xl:px-8 sm:px-4 py-8 bg-white rounded-xl shadow-sm hover:shadow-md hover:cursor-pointer">
                                                              <Flex className="relative w-fit">
                                                                  <ProgressCircle
                                                                      value={0}
                                                                      size="sm"
                                                                  >
                                                                      <div className="animate-pulse h-3 2xl:w-8 sm:w-4 my-2 bg-slate-200 dark:bg-slate-700 rounded" />
                                                                  </ProgressCircle>
                                                              </Flex>

                                                              <Flex
                                                                  alignItems="start"
                                                                  flexDirection="col"
                                                                  className="gap-1"
                                                              >
                                                                  <div className="animate-pulse h-3 2xl:w-20 sm:w-10 my-2 bg-slate-200 dark:bg-slate-700 rounded" />
                                                              </Flex>
                                                          </Flex>
                                                      ))
                                                    : response
                                                          .sort((a, b) => {
                                                              if (
                                                                  a.benchmark_title ===
                                                                      'Supportability' &&
                                                                  b.benchmark_title ===
                                                                      'Efficiency'
                                                              ) {
                                                                  return 1
                                                              }
                                                              if (
                                                                  a.benchmark_title ===
                                                                      'Efficiency' &&
                                                                  b.benchmark_title ===
                                                                      'Supportability'
                                                              ) {
                                                                  return -1
                                                              }
                                                              if (
                                                                  a.benchmark_title ===
                                                                      'Reliability' &&
                                                                  b.benchmark_title ===
                                                                      'Efficiency'
                                                              ) {
                                                                  return -1
                                                              }
                                                              if (
                                                                  a.benchmark_title ===
                                                                      'Efficiency' &&
                                                                  b.benchmark_title ===
                                                                      'Reliability'
                                                              ) {
                                                                  return 1
                                                              }
                                                              if (
                                                                  a.benchmark_title ===
                                                                      'Supportability' &&
                                                                  b.benchmark_title ===
                                                                      'Reliability'
                                                              ) {
                                                                  return 1
                                                              }
                                                              if (
                                                                  a.benchmark_title ===
                                                                      'Security' &&
                                                                  b.benchmark_title ===
                                                                      'Reliability'
                                                              ) {
                                                                  return -1
                                                              }
                                                              return 0
                                                          })
                                                          .map((item) => {
                                                              return (
                                                                  <ScoreCategoryCard
                                                                      title={
                                                                          item.benchmark_title ||
                                                                          ''
                                                                      }
                                                                      percentage={
                                                                          (item
                                                                              .severity_summary_by_control
                                                                              .total
                                                                              .passed /
                                                                              item
                                                                                  .severity_summary_by_control
                                                                                  .total
                                                                                  .total) *
                                                                          100
                                                                      }
                                                                      costOptimization={
                                                                          item.cost_optimization
                                                                      }
                                                                      value={
                                                                          item.issues_count
                                                                      }
                                                                      kpiText="Incidents"
                                                                      category={
                                                                          item.benchmark_id
                                                                      }
                                                                      varient="minimized"
                                                                  />
                                                              )
                                                          })}
                                            </Grid>
                                            {/* <Card className="w-full md:w-7/12">
                                <div className="inline-flex items-center justify-center rounded-tremor-small border border-tremor-border p-2 dark:border-dark-tremor-border">
                                    <DocumentTextIcon
                                        className="size-5 text-tremor-content-emphasis dark:text-dark-tremor-content-emphasis"
                                        aria-hidden={true}
                                    />
                                </div>
                                <h3 className="mt-4 text-tremor-default font-medium text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                    <a
                                        href="https://docs.opengovernance.io/"
                                        target="_blank"
                                        className="focus:outline-none"
                                    >
                                     
                                        <span
                                            className="absolute inset-0"
                                            aria-hidden={true}
                                        />
                                        Documentation
                                    </a>
                                </h3>
                                <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                    Learn how to audit for compliance.
                                </p>
                            </Card> */}
                                        </header>
                                    </div>
                                    <div className="w-full">
                                        <div className="p-4 sm:p-6 lg:p-8">
                                            <main>
                                                <div className="flex items-center justify-between">
                                                    {/* <h2 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                            Available Dashboards
                                        </h2> */}
                                                    <div className="flex items-center space-x-2">
                                                        {/* <Select
                                        placeholder="Sorty by"
                                        enableClear={false}
                                        className="[&>button]:rounded-tremor-small"
                                    >
                                        <SelectItem value="1">Name</SelectItem>
                                        <SelectItem value="2">
                                            Last edited
                                        </SelectItem>
                                        <SelectItem value="3">Size</SelectItem>
                                    </Select> */}
                                                        {/* <button
                                                type="button"
                                                onClick={() => {
                                                    f()
                                                    setOpen(true)
                                                }}
                                                className="hidden h-9 items-center gap-1.5 whitespace-nowrap rounded-tremor-small bg-tremor-brand px-3 py-2.5 text-tremor-default font-medium text-tremor-brand-inverted shadow-tremor-input hover:bg-tremor-brand-emphasis dark:bg-dark-tremor-brand dark:text-dark-tremor-brand-inverted dark:shadow-dark-tremor-input dark:hover:bg-dark-tremor-brand-emphasis sm:inline-flex"
                                            >
                                                <PlusIcon
                                                    className="-ml-1 size-5 shrink-0"
                                                    aria-hidden={true}
                                                />
                                                Create new Dashboard
                                            </button> */}
                                                    </div>
                                                </div>
                                                <div className="flex items-center w-full">
                                                    <Grid
                                                        numItemsMd={1}
                                                        numItemsLg={1}
                                                        className="gap-[10px] mt-1 w-full justify-items-start"
                                                    >
                                                        {loading ? (
                                                            <Spinner />
                                                        ) : (
                                                            <>
                                                                <Grid className="w-full gap-4 justify-items-start">
                                                                    <Header className="w-full">
                                                                        Frameworks{' '}
                                                                        <span className=" font-medium">
                                                                            (
                                                                            {
                                                                                totalCount
                                                                            }
                                                                            )
                                                                        </span>
                                                                    </Header>
                                                                    <Grid
                                                                        numItems={
                                                                            9
                                                                        }
                                                                        className="gap-2 min-h-[80px]  w-full "
                                                                    >
                                                                        <Col
                                                                            numColSpan={
                                                                                4
                                                                            }
                                                                        >
                                                                            <PropertyFilter
                                                                                query={
                                                                                    query
                                                                                }
                                                                                onChange={({
                                                                                    detail,
                                                                                }) => {
                                                                                    setQuery(
                                                                                        detail
                                                                                    )
                                                                                    setPage(
                                                                                        1
                                                                                    )
                                                                                }}
                                                                                // countText="5 matches"
                                                                                // enableTokenGroups
                                                                                expandToViewport
                                                                                filteringAriaLabel="Filter Benchmarks"
                                                                                filteringOptions={getFilterOptions()}
                                                                                filteringPlaceholder="Find Frameworks"
                                                                                filteringProperties={[
                                                                                    {
                                                                                        key: 'integrationType',
                                                                                        operators:
                                                                                            [
                                                                                                '=',
                                                                                            ],
                                                                                        propertyLabel:
                                                                                            'integration Type',
                                                                                        groupValuesLabel:
                                                                                            'integration Type values',
                                                                                    },
                                                                                    {
                                                                                        key: 'enable',
                                                                                        operators:
                                                                                            [
                                                                                                '=',
                                                                                            ],
                                                                                        propertyLabel:
                                                                                            'Is Active',
                                                                                        groupValuesLabel:
                                                                                            'Is Active',
                                                                                    },
                                                                                    {
                                                                                        key: 'title_regex',
                                                                                        operators:
                                                                                            [
                                                                                                '=',
                                                                                            ],
                                                                                        propertyLabel:
                                                                                            'Title',
                                                                                        groupValuesLabel:
                                                                                            'Title',
                                                                                    },
                                                                                    // {
                                                                                    //     key: 'family',
                                                                                    //     operators: [
                                                                                    //         '=',
                                                                                    //     ],
                                                                                    //     propertyLabel:
                                                                                    //         'Family',
                                                                                    //     groupValuesLabel:
                                                                                    //         'Family values',
                                                                                    // },
                                                                                ]}
                                                                            />
                                                                        </Col>
                                                                        <Col
                                                                            numColSpan={
                                                                                5
                                                                            }
                                                                        >
                                                                            <Flex
                                                                                className="w-full"
                                                                                justifyContent="end"
                                                                            >
                                                                                <Pagination
                                                                                    currentPageIndex={
                                                                                        page
                                                                                    }
                                                                                    pagesCount={
                                                                                        totalPage
                                                                                    }
                                                                                    onChange={({
                                                                                        detail,
                                                                                    }) =>
                                                                                        setPage(
                                                                                            detail.currentPageIndex
                                                                                        )
                                                                                    }
                                                                                />
                                                                            </Flex>
                                                                        </Col>
                                                                    </Grid>
                                                                    <BenchmarkCards
                                                                        benchmark={
                                                                            BenchmarkDetails
                                                                        }
                                                                        all={
                                                                            AllBenchmarks
                                                                        }
                                                                        loading={
                                                                            loading
                                                                        }
                                                                    />
                                                                </Grid>
                                                            </>
                                                        )}
                                                    </Grid>
                                                </div>
                                            </main>
                                        </div>
                                    </div>
                                </Flex>
                            </>
                        ),
                    },
                    {
                        id: '1',
                        label: 'Controls',
                        content: <AllControls />,
                    },
                    {
                        id: '2',
                        label: 'Policy',
                        content: <AllPolicy />,
                    },
                    {
                        id: '3',
                        label: 'Parameters',
                        content: <SettingsParameters />,
                    },
                ]}
            />
        </>
    )
}
