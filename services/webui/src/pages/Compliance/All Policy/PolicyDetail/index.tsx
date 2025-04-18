import { Link, useParams } from 'react-router-dom'
import { useAtomValue, useSetAtom } from 'jotai'
import {
    Button,
    Card,
    Flex,
    Grid,

    Title,
} from '@tremor/react'
import { useEffect, useState } from 'react'
import ReactJson from '@microlink/react-json-view'
import {

    PlayCircleIcon,
   
} from '@heroicons/react/24/outline'

import Spinner from '../../../../components/Spinner'
// import { severityBadge } from '../Controls'
import { isDemoAtom, notificationAtom, queryAtom } from '../../../../store'
// import Timeline from '../FindingsWithFailure/Detail/Timeline'

import Editor from 'react-simple-code-editor'
import { highlight, languages } from 'prismjs' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/components/prism-sql' // eslint-disable-next-line import/no-extraneous-dependencies
import 'prismjs/themes/prism.css'
import { Badge, KeyValuePairs, Tabs } from '@cloudscape-design/components'
import axios from 'axios'

interface IResourceFindingDetail {
    selectedItem: any | undefined
    open: boolean
    onClose: () => void
    onRefresh: () => void
    linkPrefix?: string
}

export default function PolicyDetail({
    selectedItem,
    open,
    onClose,
    onRefresh,
    linkPrefix = '',
}: IResourceFindingDetail) {
    const { ws } = useParams()
    const setQuery = useSetAtom(queryAtom)
    const [params, setParams] = useState([])

 
    const getItems = () => {
        const items = [
            {
                label: 'ID',
                value: selectedItem?.ID,
            },
           

            {
                label: 'Type',
                value: selectedItem?.type,
            },
            
           
        ]
        if(selectedItem?.type == 'external'){
            items.push({
                label: 'Title',
                value: selectedItem?.title,
            })
        }
        items.push(
            {
                label: 'Description',
                value: selectedItem?.description,
            },
            {
                label: 'Language',
                value: selectedItem?.language,
            },
            {
                label: 'Control count',
                value: selectedItem?.controls_count,
            },
            {
                label: 'Controls list',
                value: selectedItem?.list_of_controls?.join(','),
            }
        )
        return items
    }


    return (
        <>
            {selectedItem ? (
                <>
                    <>
                        <KeyValuePairs columns={4} items={getItems()} />

                        <Grid className="w-full gap-4 mb-6" numItems={1}>
                            {/* <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    ID :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.id}
                                                </Text>
                                            </Flex>
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Title :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.title}
                                                </Text>
                                            </Flex>{' '}
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Description :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.description}
                                                </Text>
                                            </Flex>{' '}
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Connector :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {selectedItem?.connector?.map(
                                                        (item, index) => {
                                                            return `${item} `
                                                        }
                                                    )}
                                                </Text>
                                            </Flex>
                                            <Flex
                                                flexDirection="row"
                                                justifyContent="between"
                                                alignItems="start"
                                                className="mt-2"
                                            >
                                                <Text className="w-56 font-bold">
                                                    Severity :{' '}
                                                </Text>
                                                <Text className="w-full">
                                                    {severityBadge(
                                                        selectedItem?.severity
                                                    )}
                                                </Text>
                                            </Flex> */}
                            <Flex
                                flexDirection="col"
                                justifyContent="between"
                                alignItems="start"
                                className="mt-2"
                            >
                                <Flex flexDirection="row" className="mb-2">
                                    <Title className="mb-2">Definition</Title>

                                    <Button
                                        icon={PlayCircleIcon}
                                        onClick={() => {
                                            // @ts-ignore
                                            setQuery(selectedItem?.definition)
                                        }}
                                        disabled={false}
                                        loading={false}
                                        loadingText="Running"
                                    >
                                        <Link to={`/cloudql`}>
                                            Open in CloudQL
                                        </Link>{' '}
                                    </Button>
                                </Flex>
                                <Card className=" py-3 mb-2 relative ">
                                    <Editor
                                        onValueChange={(text) => {
                                            console.log(text)
                                        }}
                                        highlight={(text) =>
                                            highlight(
                                                text,
                                                languages.sql,
                                                'sql'
                                            )
                                        }
                                        // @ts-ignore
                                        value={selectedItem?.definition}
                                        className="w-full bg-white dark:bg-gray-900 dark:text-gray-50 font-mono text-sm"
                                        style={{
                                            minHeight: '200px',
                                            // maxHeight: '500px',
                                            overflowY: 'scroll',
                                        }}
                                        placeholder="-- write your SQL query here"
                                        disabled={true}
                                    />
                                </Card>
                                {/* <Flex
                                                    flexDirection="row"
                                                    alignItems="start"
                                                    className="gap-1 w-full flex-wrap "
                                                    justifyContent="start"
                                                >
                                                    {}
                                                    {Object.entries(
                                                        selectedItem?.tags
                                                    ).map((key, index) => {
                                                        return (
                                                            <>
                                                                <Flex
                                                                    flexDirection="row"
                                                                    justifyContent="start"
                                                                    className="hover:cursor-pointer max-w-full w-fit bg-gray-200 border-gray-300 rounded-lg border px-1"
                                                                >
                                                                    <TagIcon className="min-w-4 w-4 mr-1" />
                                                                    <Text className="truncate">
                                                                        {key[0]}
                                                                        :
                                                                        {key[1]}
                                                                    </Text>
                                                                </Flex>
                                                            </>
                                                        )
                                                    })}
                                                </Flex> */}
                            </Flex>
                        </Grid>
                    </>
                </>
            ) : (
                <Spinner />
            )}
        </>
    )
}
