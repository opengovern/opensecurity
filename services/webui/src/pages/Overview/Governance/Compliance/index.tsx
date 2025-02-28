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
    Col,
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
         per_page: window.innerWidth > 1920 ? 3 : window.innerWidth > 768 ? 3 :3,
         sort_by: 'incidents',
         assigned: false,
         is_baseline: false,
     }
     axios
         .post(`${url}/main/compliance/api/v1/frameworks`, body, config)
         .then((res) => {

             setBenchmarks(res.data.items)
             setLoading(false)

         })
         .catch((err) => {
             setLoading(false)

             console.log(err)
         })
 }


 
   useEffect(() => {

       GetCard()
   }, [])

   const array = window.innerWidth > 768 ? [1,2,3] : [1,2,3,4,5]

    return (
        <Flex flexDirection="col" alignItems="start" justifyContent="start">
            {loading ? (
                <Flex className="gap-4 flex-wrap sm:flex-row flex-col">
                    {array.map((i) => {
                        return (
                            <Card className="p-3 dark:ring-gray-500 sm:w-[calc(33%-0.5rem)] w-[calc(100%-0.5rem)] sm:h-64 h-32">
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
                <Grid
                    numItems={10}
                    className="w-full gap-4 mt-4  justify-center items-center "
                >
                    {/* <Col numColSpan={0} numColSpanSm={1}></Col> */}
                    <Col numColSpan={10} numColSpanSm={10}>
                        <BenchmarkCards
                            all={AllBenchmarks}
                            loading={loading}
                        />
                    </Col>
                    {/* <Col numColSpan={0} numColSpanSm={1}></Col> */}
                </Grid>
            )}
        </Flex>
    )
}


