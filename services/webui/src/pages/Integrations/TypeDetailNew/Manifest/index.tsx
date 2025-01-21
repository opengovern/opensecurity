import { Card, Flex, Title } from '@tremor/react'

import { useEffect, useState } from 'react'

import axios from 'axios'

import { KeyValuePairs } from '@cloudscape-design/components'
import Spinner from '../../../../components/Spinner'
interface IntegrationListProps {
    name?: string
    integration_type?: string
}

export default function Manifest({
    name,
    integration_type,
}: IntegrationListProps) {
    const [manifest, setManifest] = useState<any>()
    const [loading, setLoading] = useState<boolean>(false)

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

    useEffect(() => {
        GetManifest()
    }, [])

    return (
        <>
            {loading ? (
                <>
                    <Spinner />
                </>
            ) : (
                <Card className="">
                    <Flex flexDirection='col' className='gap-2 p-5 w-full justify-start items-start'>
                        <>
                            <Title className=" mb-4 font-semibold">
                                Plugin information
                            </Title>
                            <KeyValuePairs
                                columns={4}
                                items={[
                                    {
                                        label: 'Id',
                                        value: manifest?.IntegrationType,
                                    },
                                    {
                                        label: 'URL',
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
                                ]}
                            />
                        </>
                    </Flex>
                </Card>
            )}
        </>
    )
}
