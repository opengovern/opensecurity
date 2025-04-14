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

import {
    PlatformEnginePkgComplianceApiBenchmarkEvaluationSummary,
    SourceType,
} from '../../../../api/api'

import TopHeader from '../../../../components/Layout/Header'
import { useURLParam, useURLState } from '../../../../utilities/urlstate'

import { errorHandling } from '../../../../types/apierror'

import Spinner from '../../../../components/Spinner'
import axios from 'axios'
import BenchmarkCard from '../BenchmarkCard'
import BenchmarkCards from '../BenchmarkCard'
import {
    Header,
    Pagination,
    PropertyFilter,
    Tabs,
} from '@cloudscape-design/components'
import Multiselect from '@cloudscape-design/components/multiselect'
import Select from '@cloudscape-design/components/select'
import ScoreCategoryCard from '../../../../components/Cards/ScoreCategoryCard'
import { useIntegrationApiV1EnabledConnectorsList } from '../../../../api/integration.gen'
import { useParams, useSearchParams } from 'react-router-dom'
import CustomPagination from '../../../../components/Pagination'
const CATEGORY = {
    sre_efficiency: 'Efficiency',
    sre_reliability: 'Reliability',
    sre_supportability: 'Supportability',
}

export default function Framework() {
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
    const [page, setPage] = useState<number>(1)
    const [totalPage, setTotalPage] = useState<number>(0)
    const [totalCount, setTotalCount] = useState<number>(0)
    const [response, setResponse] = useState()
    const [isLoading, setIsLoading] = useState(false)
    const [tab,setTab] = useState('0')
    const [tab_id,setTabID] = useSearchParams()
    const number= window.innerWidth > 640 ? 9 : 4
    const {
        response: Types,
        isLoading: TypesLoading,
        isExecuted: TypesExec,
    } = useIntegrationApiV1EnabledConnectorsList(0, 0)

    const getFilterOptions = () => {
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
    const GetCard = (ids : string[] | undefined) => {
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
        const connectors :any = []
        const enable: any = []
        const isSRE: any = []
        const title: any = []
        query.tokens.map((item: any) => {
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

        let body = {
            cursor: page,
            per_page: number,
            sort_by: 'incidents',
            assigned: false,
            is_baseline: sre_filter,
            integrationType: connector_filter,
            root: true,
            title_regex: title[0],
        }
        if(ids && ids?.length > 0){
            // @ts-ignore
            body["framework_ids"]= ids
        }

        axios
            .post(`${url}/main/compliance/api/v1/frameworks`, body, config)
            .then((res) => {
                //  const temp = []
                if (!res.data.items) {
                    setLoading(false)
                }
                setBenchmarks(res.data.items)
                setTotalPage(Math.ceil(res.data.total_count / number))
                setTotalCount(res.data.total_count)
                setLoading(false)

            })
            .catch((err) => {
                setLoading(false)
                // @ts-ignore
                setBenchmarks([])

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
            framework_ids: benchmarks,
        }
        axios
            .post(`${url}/main/compliance/api/v1/frameworks`, body, config)
            .then((res) => {
                const temp: any = []
                setIsLoading(false)
                res.data?.items?.map((item: any) => {
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
        GetCard(undefined)
    }, [page, query])


    useEffect(() => {
        if(window.innerWidth > 768){
 GetBenchmarks([
     'baseline_efficiency',
     'baseline_reliability',
     'baseline_security',
     'baseline_supportability',
 ])
        }
        else{
 GetBenchmarks([
     'baseline_efficiency',
     'baseline_security',
 ])
        }
       
    }, [])
    useEffect(()=>{
        switch (tab_id?.get('tab')) {
            case 'frameworks':
                setTab('frameworks')
                break
            case 'controls':
                setTab('controls')
                break
            case 'policies':
                setTab('policies')
                break
            case 'parameters':
                setTab('parameters')
                break
            default:
                setTab('frameworks')
                break
        }
    },[tab_id])
    const array = window.innerWidth >768 ? [1,2,3,4] : [1,2]

    return (
        <>
            <Flex
                className="bg-white w-full rounded-xl border-solid  border-2 border-gray-200   "
                flexDirection="col"
                justifyContent="center"
                alignItems="center"
            >
                <div className="border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:p-6 lg:p-8">
                    <header className="flex flex-col gap-4 w-full">
                        <h1 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                            Frameworks
                        </h1>
                        <p className="text-tremor-default w-full sm:inline-block hidden text-tremor-content dark:text-dark-tremor-content">
                            Frameworks are structured guides, like practical
                            blueprints, that provide proven methods for handling
                            important areas like security or reliable
                            operations.
                            <br />
                            <br />
                            They help you:
                            <br />
                            <br />
                            <ol className=" list-decimal list-inside">
                                <li>
                                    Decide which specific checks (Controls) are
                                    necessary.
                                </li>
                                <br />

                                <li>
                                    {' '}
                                    Organize those related checks into logical
                                    Control Groups (e.g., 'Password Rules,'
                                    'Data Backup Procedures,' or 'Cloud
                                    Settings') to keep things neat and
                                    manageable as you grow.
                                </li>
                            </ol>
                        </p>
                        <Grid
                            numItems={1}
                            numItemsSm={window.innerWidth > 1440 ? 4 : 2}
                            className="2xl:gap-[30px]  sm:gap-10 gap-4 mt-6 w-full justify-items-center"
                        >
                            {isLoading || !response
                                ? array.map((i) => (
                                      <Flex className="gap-6 2xl:px-8 sm:px-4 px-2 py-8 w-full bg-white rounded-xl shadow-sm hover:shadow-md hover:cursor-pointer">
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
                                      // @ts-ignore
                                      .sort((a, b) => {
                                          if (
                                              a.benchmark_title ===
                                                  'Supportability' &&
                                              b.benchmark_title === 'Efficiency'
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
                                              b.benchmark_title === 'Efficiency'
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
                                      .map((item: any) => {
                                          return (
                                              <ScoreCategoryCard
                                                  title={
                                                      item.framework_title || ''
                                                  }
                                                  percentage={
                                                      (item
                                                          .severity_summary_by_control
                                                          .total.passed /
                                                          item
                                                              .severity_summary_by_control
                                                              .total.total) *
                                                      100
                                                  }
                                                  costOptimization={
                                                      item?.cost_optimization
                                                  }
                                                  value={item.issues_count}
                                                  kpiText="Incidents"
                                                  category={item.framework_id}
                                                  varient="minimized"
                                              />
                                          )
                                      })}
                        </Grid>
                    </header>
                </div>
                <div className="w-full">
                    <div className="p-4 sm:p-6 lg:p-8">
                        <main>
                            <div className="flex items-center justify-between">
                                <div className="flex items-center space-x-2"></div>
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
                                                        ({totalCount})
                                                    </span>
                                                </Header>
                                                <Grid
                                                    numItems={9}
                                                    className="gap-2 min-h-[80px]  w-full "
                                                >
                                                    <Col
                                                        numColSpan={9}
                                                        numColSpanSm={4}
                                                    >
                                                        <PropertyFilter
                                                            // @ts-ignore
                                                            query={query}
                                                            onChange={({
                                                                detail,
                                                            }) => {
                                                                // @ts-ignore

                                                                setQuery(detail)
                                                                setPage(1)
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
                                                                    operators: [
                                                                        '=',
                                                                    ],
                                                                    propertyLabel:
                                                                        'integration Type',
                                                                    groupValuesLabel:
                                                                        'integration Type values',
                                                                },
                                                                {
                                                                    key: 'enable',
                                                                    operators: [
                                                                        '=',
                                                                    ],
                                                                    propertyLabel:
                                                                        'Is Active',
                                                                    groupValuesLabel:
                                                                        'Is Active',
                                                                },
                                                                {
                                                                    key: 'title_regex',
                                                                    operators: [
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
                                                        numColSpan={9}
                                                        numColSpanSm={5}
                                                    >
                                                        <Flex
                                                            className="w-full"
                                                            justifyContent="end"
                                                        >
                                                            <CustomPagination
                                                                currentPageIndex={
                                                                    page
                                                                }
                                                                pagesCount={
                                                                    totalPage
                                                                }
                                                                onChange={({
                                                                    detail,
                                                                }: any) =>
                                                                    setPage(
                                                                        detail.currentPageIndex
                                                                    )
                                                                }
                                                            />
                                                        </Flex>
                                                    </Col>
                                                </Grid>
                                                <BenchmarkCards
                                                    // @ts-ignore
                                                    all={AllBenchmarks}
                                                    loading={loading}
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
    )
}
