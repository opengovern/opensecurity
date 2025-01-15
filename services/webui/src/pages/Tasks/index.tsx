import { useSetAtom } from "jotai"
import { useEffect, useState } from "react"
import { notificationAtom } from "../../store"
import { useNavigate } from "react-router-dom"
import axios from "axios"
import { Card, Flex } from "@tremor/react"
import Spinner from "../../components/Spinner"
import { DocumentTextIcon } from "@heroicons/react/24/outline"
import { Box, Cards, Link, Pagination, SpaceBetween } from "@cloudscape-design/components"


export default function Tasks() {
   const [pageNo, setPageNo] = useState<number>(0)
  
   const [open, setOpen] = useState(false)
   const navigate = useNavigate()
   const [tasks, setTasks] = useState([])
   const [loading, setLoading] = useState(false)
   const [total_count, setTotalCount] = useState(0)
   const setNotification = useSetAtom(notificationAtom)

  
   const getTasks = () => {
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

       axios
           .get(
               `${url}/main/tasks/api/v1/tasks?per_page=${10}&cursor=${pageNo}`,
               config
           )
           .then((res) => {
               setLoading(false)
               setTasks(res.data?.items)
               setTotalCount(res.data?.total_count)
           })
           .catch((err) => {
               setLoading(false)
           })
   }
   useEffect(()=>{
        getTasks()
   },[pageNo])
   return (
       <>
           {/* <TopHeader /> */}
           {/* <Grid numItems={3} className="gap-4 mb-10">
                <OnboardCard
                    title="Active Accounts"
                    active={topMetrics?.connectionsEnabled}
                    inProgress={topMetrics?.inProgressConnections}
                    healthy={topMetrics?.healthyConnections}
                    unhealthy={topMetrics?.unhealthyConnections}
                    loading={metricsLoading}
                />
            </Grid> */}
           {loading ? (
               <Flex className="mt-36">
                   <Spinner />
               </Flex>
           ) : (
               <>
                   {/* <TabGroup className='mt-4'>
                        <TabList>
                            <Tab>test</Tab>
                            <Tab>test</Tab>
                            <Tab>test</Tab>
                            <Tab>test</Tab>
                        </TabList>
                    </TabGroup> */}
                   <Flex
                       className="bg-white w-[90%] rounded-xl border-solid  border-2 border-gray-200  pb-2  "
                       flexDirection="col"
                       justifyContent="center"
                       alignItems="center"
                   >
                       <div className="border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:p-6 lg:p-8">
                           <header>
                               <h1 className="text-tremor-title font-semibold text-tremor-content-strong dark:text-dark-tremor-content-strong">
                                   Tasks
                               </h1>
                               <p className="text-tremor-default text-tremor-content dark:text-dark-tremor-content">
                                   Create and Manage your Tasks
                               </p>
                               <div className="mt-8 w-full md:flex md:max-w-3xl md:items-stretch md:space-x-4">
                                   <Card className="w-full md:w-7/12">
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
                                               {/* Extend link to entire card */}
                                               <span
                                                   className="absolute inset-0"
                                                   aria-hidden={true}
                                               />
                                               Documentation
                                           </a>
                                       </h3>
                                       <p className="dark:text-dark-tremor-cont text-tremor-default text-tremor-content">
                                           Learn how to add, update, remove
                                           Tasks
                                       </p>
                                   </Card>
                               </div>
                           </header>
                       </div>
                       <div className="w-full">
                           <div className="p-4 sm:p-6 lg:p-8">
                               <main>
                                   <div className="flex items-center justify-between">
                                       <div className="flex items-center space-x-2"></div>
                                   </div>
                                   <div className="flex items-center w-full">
                                       <Cards
                                        className="w-full"
                                           ariaLabels={{
                                               itemSelectionLabel: (e, t) =>
                                                   `select `,
                                               selectionGroupLabel:
                                                   'Item selection',
                                           }}
                                           onSelectionChange={({ detail }) => {
                                               const connector =
                                                   detail?.selectedItems[0]
                                           }}
                                           selectedItems={[]}
                                           cardDefinition={{
                                               header: (item) => (
                                                   <Link
                                                       className="w-100"
                                                       onClick={() => {
                                                              navigate(
                                                                `/tasks/${item.id}`
                                                              )
                                                       }}
                                                   >
                                                       <div className="w-100 flex flex-row justify-between">
                                                           <span>
                                                               {item.name}
                                                           </span>
                                                           {/* <div className="flex flex-row gap-1 items-center">
                                    {GetTierIcon(item.tier)}
                                    <span className="text-white">{item.tier}</span>
                                </div> */}
                                                       </div>
                                                   </Link>
                                               ),
                                               sections: [
                                                   {
                                                       id: 'description',
                                                       header: (
                                                           <>
                                                               <div className="flex justify-between">
                                                                   <span>
                                                                       {
                                                                           'Description'
                                                                       }
                                                                   </span>
                                                               </div>
                                                           </>
                                                       ),
                                                       content: (item) => (
                                                           <>
                                                               <div className="flex justify-between">
                                                                   <span className="max-w-60">
                                                                       {
                                                                           item.description
                                                                       }
                                                                   </span>
                                                               </div>
                                                           </>
                                                       ),
                                                   },
                                               ],
                                           }}
                                           cardsPerRow={[
                                               { cards: 1 },
                                               { minWidth: 540, cards: 2 },
                                               { minWidth: 750, cards: 3 },
                                           ]}
                                           // @ts-ignore
                                           items={tasks?.map((type: any) => {
                                               return {
                                                   id: type.id,
                                                   name: type.name,
                                                   description:
                                                       type.description,

                                                   // schema_id: type?.schema_ids[0],
                                                   // SourceCode: type.SourceCode,
                                               }
                                           })}
                                           loadingText="Loading resources"
                                           stickyHeader
                                           entireCardClickable
                                           variant="full-page"
                                           selectionType="single"
                                           trackBy="name"
                                           empty={
                                               <Box
                                                   margin={{ vertical: 'xs' }}
                                                   textAlign="center"
                                                   color="inherit"
                                               >
                                                   <SpaceBetween size="m">
                                                       <b>No resources</b>
                                                   </SpaceBetween>
                                               </Box>
                                           }
                                       />
                                   </div>
                               </main>
                           </div>
                       </div>
                       <Pagination
                           currentPageIndex={pageNo}
                           pagesCount={total_count}
                           onChange={({ detail }) => {
                               setPageNo(detail.currentPageIndex)
                           }}
                       />
                   </Flex>
               </>
           )}
       </>
   )
}
