import { Card, Flex, Title } from '@tremor/react'

import { useEffect, useState } from 'react'

import axios from 'axios'

import { KeyValuePairs } from '@cloudscape-design/components'
import Spinner from '../../../../components/Spinner'
import ReactMarkdown from 'react-markdown'
import rehypeRaw from 'rehype-raw'
import { useMDXComponents } from '../../../../components/MDX'
interface IntegrationListProps {
    name?: string
    integration_type?: string
}

export default function Configuration({
    name,
    integration_type,
}: IntegrationListProps) {
    const [manifest, setManifest] = useState<any>()
    const [loading, setLoading] = useState<boolean>(false)
    const [setup, setSetup] = useState<any>()
    const [status, setStatus] = useState<string>()


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

    const GetManifest = () => {
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
                `${url}/main/integration/api/v1/integration-types/plugin/${integration_type}/manifest`,
                config
            )
            .then((resp) => {
                setManifest(resp.data)
                setLoading(false)
            })
            .catch((err) => {
                console.log(err)
                setLoading(false)

                // params.fail()
            })
    }
     const GetStatus = () => {
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
             .post(
                 `${url}/main/integration/api/v1/integration-types/plugin/${integration_type}/healthcheck`,{},
                 config
             )
             .then((resp) => {
                 setStatus(resp.data)
                 setLoading(false)
             })
             .catch((err) => {
                 console.log(err)
                 setLoading(false)

                 // params.fail()
             })
     }

    useEffect(() => {
        GetStatus()
        GetSetup()
        GetManifest()
    }, [])

    return (
        <>
            {loading ? (
                <>
                    <Spinner />
                </>
            ) : (
                <Flex className="flex-col gap-3 w-full">
                    <h1 className=" font-bold text-2xl mb-2 w-full text-left ml-1">
                        Information
                    </h1>
                    <Card className="">
                        <Flex
                            flexDirection="col"
                            className="gap-2 p-2 w-full justify-start items-start"
                        >
                            <>
                                <KeyValuePairs
                                    columns={4}
                                    items={[
                                        {
                                            label: 'Id',
                                            value: manifest?.IntegrationType,
                                        },
                                        {
                                            label: 'Artifact URL',
                                            value: manifest?.DescriberURL,
                                        },
                                        {
                                            label: 'Version',
                                            value: manifest?.DescriberTag,
                                        },
                                        {
                                            label: 'Publisher',
                                            value: manifest?.Publisher,
                                        },
                                        {
                                            label: 'Author',
                                            value: manifest?.Author,
                                        },
                                        {
                                            label: 'Supported Platform Version',
                                            value: manifest?.SupportedPlatformVersion,
                                        },
                                        {
                                            label: 'Update date',
                                            value: manifest?.UpdateDate,
                                        },
                                        {
                                            label: 'Operational Status',
                                            // @ts-ignore
                                            value: status ? status?.charAt(0).toUpperCase()+status?.slice(1) : '',
                                        },
                                    ]}
                                />
                            </>
                        </Flex>
                    </Card>
                    <>
                        {' '}
                        <h1 className=" font-bold text-2xl mb-2 w-full text-left ml-1">
                            Setup guide
                        </h1>
                        <Card className="p-2">
                            <Flex
                                flexDirection="col"
                                className=" p-5 justify-start w-full items-start"
                            >
                                <ReactMarkdown
                                    children={setup}
                                    skipHtml={false}
                                    rehypePlugins={[rehypeRaw]}
                                    // @ts-ignore
                                    components={useMDXComponents({})}
                                />
                            </Flex>
                        </Card>
                    </>
                </Flex>
            )}
        </>
    )
}
