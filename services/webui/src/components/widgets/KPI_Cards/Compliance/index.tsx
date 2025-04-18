// @ts-nocheck
import {
    Button,
    Card,
    Flex,
    Subtitle,
    Text,
    Title,
    Divider,
    CategoryBar,
    Grid,
    ProgressCircle,
} from '@tremor/react'
import { useNavigate, useParams } from 'react-router-dom'
import { ChevronRightIcon } from '@heroicons/react/20/solid'
import { useAtomValue } from 'jotai'
import { useComplianceApiV1BenchmarksSummaryList } from '../../../../api/compliance.gen'
import { getErrorMessage } from '../../../../types/apierror'
import { searchAtom } from '../../../../utilities/urlstate'
import BenchmarkCards from '../../../Compliance/BenchmarkCard'
import { useEffect, useState } from 'react'
import axios from 'axios'
import ScoreCategoryCard from '../../../../components/Cards/ScoreCategoryCard'

const colors = [
    'fuchsia',
    'indigo',
    'slate',
    'gray',
    'zinc',
    'neutral',
    'stone',
    'red',
    'orange',
    'amber',
    'yellow',
    'lime',
    'green',
    'emerald',
    'teal',
    'cyan',
    'sky',
    'blue',
    'violet',
    'purple',
    'pink',
    'rose',
]

export default function Compliance() {
    const workspace = useParams<{ ws: string }>().ws
    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const [loading,setLoading] = useState<boolean>(false);
    const [response, setResponse] = useState()



     const GetBenchmarks = (benchmarks: string[]) => {
         setLoading(true)
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
                 setLoading(false)
                 res.data?.items?.map((item: any) => {
                     temp.push(item)
                 })
                 setResponse(temp)
             })
             .catch((err) => {
                 setLoading(false)

                 console.log(err)
             })
     }

        useEffect(() => {
              GetBenchmarks([
                  'baseline_efficiency',
                  'baseline_reliability',
                  'baseline_security',
                  'baseline_supportability',
              ])
        }, [])
   const array = window.innerWidth > 768 ? [1,2,3,4] : [1,2,3,4]

    return (
        <Grid
            numItems={1}
            numItemsSm={window.innerWidth > 1440 ? 4 : 2}
            className="2xl:gap-[20px]  sm:gap-8 gap-4 mt-4 w-full justify-items-center"
        >
            {loading || !response
                ? array.map((i) => (
                      <Flex className="gap-6 2xl:px-8 sm:px-4 px-2 py-8 w-full bg-white rounded-xl shadow-sm hover:shadow-md hover:cursor-pointer">
                          <Flex className="relative w-fit">
                              <ProgressCircle value={0} size="sm">
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
                              a.benchmark_title === 'Supportability' &&
                              b.benchmark_title === 'Efficiency'
                          ) {
                              return 1
                          }
                          if (
                              a.benchmark_title === 'Efficiency' &&
                              b.benchmark_title === 'Supportability'
                          ) {
                              return -1
                          }
                          if (
                              a.benchmark_title === 'Reliability' &&
                              b.benchmark_title === 'Efficiency'
                          ) {
                              return -1
                          }
                          if (
                              a.benchmark_title === 'Efficiency' &&
                              b.benchmark_title === 'Reliability'
                          ) {
                              return 1
                          }
                          if (
                              a.benchmark_title === 'Supportability' &&
                              b.benchmark_title === 'Reliability'
                          ) {
                              return 1
                          }
                          if (
                              a.benchmark_title === 'Security' &&
                              b.benchmark_title === 'Reliability'
                          ) {
                              return -1
                          }
                          return 0
                      })
                      .map((item: any) => {
                          return (
                              <ScoreCategoryCard
                                  title={item.framework_title || ''}
                                  percentage={
                                      (item.severity_summary_by_control.total
                                          .passed /
                                          item.severity_summary_by_control.total
                                              .total) *
                                      100
                                  }
                                  costOptimization={item.cost_optimization}
                                  value={item.issues_count}
                                  kpiText="Incidents"
                                  category={item.framework_id}
                                  varient="minimized"
                              />
                          )
                      })}
        </Grid>
    )
}


