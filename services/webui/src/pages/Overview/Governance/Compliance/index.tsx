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
} from '@tremor/react'
import { useNavigate, useParams } from 'react-router-dom'
import { ChevronRightIcon } from '@heroicons/react/20/solid'
import { useAtomValue } from 'jotai'
import { useComplianceApiV1BenchmarksSummaryList } from '../../../../api/compliance.gen'
import { getErrorMessage } from '../../../../types/apierror'
import { searchAtom } from '../../../../utilities/urlstate'
import BenchmarkCards from '../../../Governance/Compliance/BenchmarkCard'
import { useEffect, useState } from 'react'
import axios from 'axios'

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
 const [AllBenchmarks,setBenchmarks] = useState();
        const [BenchmarkDetails, setBenchmarksDetails] = useState()
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
     const body = {
         cursor: 1,
         per_page: window.innerWidth > 1920 ? 6 : window.innerWidth > 768 ? 4 :5,
         sort_by: 'incidents',
         assigned: false,
         is_baseline: false,
     }
     axios
         .post(`${url}/main/compliance/api/v3/benchmarks`, body,config)
         .then((res) => {
             //  const temp = []

            setBenchmarks(res.data.items)
         })
         .catch((err) => {
                setLoading(false)

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
         benchmarks: benchmarks
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

              console.log(err)
          })
  }
 
   useEffect(() => {

       GetCard()
   }, [])
   useEffect(() => {
    if(AllBenchmarks){
  const temp = []
  AllBenchmarks?.map((item) => {
      temp.push(item.benchmark.id)
  })
  Detail(temp)
    }
    
   }, [AllBenchmarks])
   const array = window.innerWidth > 768 ? [1,2,3,4] : [1,2,3,4,5]

    return (
        <Flex flexDirection="col" alignItems="start" justifyContent="start">
           
            {loading ? (
                <Flex  className="gap-4 flex-wrap sm:flex-row flex-col">
                    {array.map((i) => {
                        return (
                            <Card className="p-3 dark:ring-gray-500 sm:w-[calc(50%-0.5rem)] w-[calc(100%-0.5rem)] sm:h-64 h-32">
                                <Flex
                                    flexDirection="col"
                                    alignItems="start"
                                    justifyContent="start"
                                    className="animate-pulse w-full"
                                >
                                    <div className="h-5 w-24  mb-2 bg-slate-200 dark:bg-slate-700 rounded" />
                                    <div className="h-5 w-24  mb-1 bg-slate-200 dark:bg-slate-700 rounded" />
                                    <div className="h-6 w-24  bg-slate-200 dark:bg-slate-700 rounded" />
                                </Flex>
                            </Card>
                        )
                    })}
                </Flex>
            ) : (
                <Grid className="w-full gap-4 justify-items-start">
                    <BenchmarkCards
                        benchmark={BenchmarkDetails}
                        all={AllBenchmarks}
                        loading={loading}
                    />
                    
                </Grid>
            )}
        </Flex>
    )
}


