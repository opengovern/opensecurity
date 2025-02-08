import { Card, Flex, Title } from '@tremor/react'

import { useEffect, useState } from 'react'

import axios from 'axios'

import { KeyValuePairs } from '@cloudscape-design/components'
import Spinner from '../../../../components/Spinner'
import ReactMarkdown from 'react-markdown'
import rehypeRaw from 'rehype-raw'
interface IntegrationListProps {
    name?: string
    integration_type?: string
}

export default function Setup({
    name,
    integration_type,
}: IntegrationListProps) {
    const [loading, setLoading] = useState<boolean>(false)
    const [setup, setSetup] = useState<any>()


    const GetSetup = () => {
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
                `${url}/main/integration/api/v1/integration-types/plugin/${integration_type}/setup`,
                config
            )
            .then((resp) => {
                setSetup(resp.data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)

                // params.fail()
            })
    }

   

    useEffect(() => {
        GetSetup()
    }, [])

    return (
        <>
            {loading ? (
                <>
                    <Spinner />
                </>
            ) : (
                <Flex className="flex-col gap-3 w-full">
                    <>
                        {' '}
                        {/* <h1 className=" font-bold text-2xl mb-2 w-full text-left ml-1">
                            Setup guide
                        </h1> */}
                        <Card className="p-2">
                            <Flex
                                flexDirection="col"
                                className=" p-5 justify-start w-full items-start"
                            >
                                <div className=" markdown-container">
                                    <ReactMarkdown
                                        // @ts-ignore
                                        children={setup}
                                        skipHtml={false}
                                        className={'markdown-body'}
                                        rehypePlugins={[rehypeRaw]}
                                    />
                                </div>
                               
                            </Flex>
                        </Card>
                    </>
                </Flex>
            )}
        </>
    )
}
