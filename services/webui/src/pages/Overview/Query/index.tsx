import {
    Accordion,
    AccordionBody,
    AccordionHeader,
    Button,
    Card,
    Flex,
    Icon,
    Text,
    Title,
} from '@tremor/react'
import {
    ChevronRightIcon,
    MagnifyingGlassIcon,
} from '@heroicons/react/24/outline'
import Editor from 'react-simple-code-editor'
import 'prismjs/themes/prism.css'
import { highlight, languages } from 'prismjs'
import { useNavigate, useParams } from 'react-router-dom'
import { useAtom } from 'jotai'
import { useState } from 'react'
import { useInventoryApiV1QueryList } from '../../../api/inventory.gen'
import { queryAtom } from '../../../store'
import { getErrorMessage } from '../../../types/apierror'
import { PlatformEnginePkgInventoryApiSmartQueryItem } from '../../../api/api'

interface IQuery {
    height: any
}

const getQueries = (
    response:
        | PlatformEnginePkgInventoryApiSmartQueryItem[]
        | undefined
) => {
    const data = []
    const queryId = [
        'ai_workload',
        'container_workload',
        'load_balancers',
        'server_workload',
        'cloud_networks',
    ]
    if (response) {
        for (let i = 0; i < queryId.length; i += 1) {
            const query = response?.find((q) => q.id === queryId[i])
            data.push(query)
        }
    }
    return data
}

export default function Query({ height }: IQuery) {
    const workspace = useParams<{ ws: string }>().ws
    const navigate = useNavigate()
    const [runQuery, setRunQuery] = useAtom(queryAtom)
    const [open, setOpen] = useState(0)

    const {
        response: queries,
        isLoading,
        error,
        sendNow: refresh,
    } = useInventoryApiV1QueryList({})

    return (
        <Card
            className="h-full overflow-scroll no-scrollbar"
            style={{ maxHeight: `${height}px` }}
        >
            <Flex justifyContent="between">
                <Flex justifyContent="start" className="gap-2 mb-2">
                    <Icon icon={MagnifyingGlassIcon} className="p-0" />
                    <Title className="font-semibold">
                        Bookmarked Inventory
                    </Title>
                </Flex>
                <a
                    target="__blank"
                    href={`/cloudql?tab_id=0`}
                    className=" cursor-pointer"
                >
                    <Button
                        size="xs"
                        variant="light"
                        icon={ChevronRightIcon}
                        iconPosition="right"
                        className="my-3"
                        // onClick={() => {
                        //     navigate(`/cloudql?tab_id=0`)
                        // }}
                    >
                        All Queries
                    </Button>
                </a>
            </Flex>
            {isLoading
                ? [1, 2, 3, 4, 5].map((i) => (
                      <Accordion
                          className={`w-full border-0 ${
                              i < 4 ? 'border-b border-b-gray-200' : ''
                          } !rounded-none bg-transparent ${
                              isLoading ? 'animate-pulse' : ''
                          }`}
                      >
                          <AccordionHeader className="pl-0 pr-0.5 py-4 bg-transparent flex justify-start">
                              <div className="h-5 w-32 bg-slate-200 dark:bg-slate-700 rounded" />
                          </AccordionHeader>
                      </Accordion>
                  ))
                : getQueries(queries)
                      ?.sort((a, b) => {
                          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                          // @ts-ignore
                          if (a.title < b.title) {
                              return -1
                          }
                          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                          // @ts-ignore
                          if (a.title > b.title) {
                              return 1
                          }
                          return 0
                      })
                      .map((q, i) => (
                          <Accordion
                              // eslint-disable-next-line react/no-array-index-key
                              key={`query-${i}-${open}`}
                              className={`w-full border-0 ${
                                  i < 4 ? 'border-b border-b-gray-200' : ''
                              } !rounded-none bg-transparent`}
                              defaultOpen={i === open}
                              onClick={() => {
                                  if (i !== open) {
                                      setOpen(i)
                                  }
                              }}
                          >
                              <AccordionHeader className="pl-0 pr-0.5 py-4 bg-transparent flex justify-start">
                                  <Text className="text-gray-800 !text-base line-clamp-1">
                                      {q?.title}
                                  </Text>
                              </AccordionHeader>
                              <AccordionBody className="p-0 w-full pr-0.5 cursor-default ">
                                  <Editor
                                      onValueChange={(text) => {
                                          console.log('')
                                      }}
                                      highlight={(text) =>
                                          highlight(text, languages.sql, 'sql')
                                      }
                                      value={q?.query || ''}
                                      className="w-full bg-gray-100 rounded p-5 dark:bg-gray-800 font-mono text-sm h-full no-scrollbar"
                                      style={{
                                          color: 'white !important',
                                          minHeight: '60px',
                                          overflowY: 'scroll',
                                          padding: '2rem!important',
                                      }}
                                      placeholder="-- write your SQL query here"
                                  />
                                  <Button
                                      size="xs"
                                      variant="light"
                                      icon={ChevronRightIcon}
                                      iconPosition="right"
                                      className="my-3"
                                      onClick={() => {
                                          setRunQuery(q?.query || '')
                                          navigate(
                                              `/cloudql`
                                          )
                                      }}
                                  >
                                      Run Query
                                  </Button>
                              </AccordionBody>
                          </Accordion>
                      ))}
            {error && (
                <Flex
                    flexDirection="col"
                    justifyContent="between"
                    className="absolute top-0 w-full left-0 h-full backdrop-blur"
                >
                    <Flex
                        flexDirection="col"
                        justifyContent="center"
                        alignItems="center"
                    >
                        <Title className="mt-6">Failed to load component</Title>
                        <Text className="mt-2">{getErrorMessage(error)}</Text>
                    </Flex>
                    <Button
                        variant="secondary"
                        className="mb-6"
                        color="slate"
                        onClick={refresh}
                    >
                        Try Again
                    </Button>
                </Flex>
            )}
        </Card>
    )
}
