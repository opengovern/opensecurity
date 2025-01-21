import { Card, Title } from '@tremor/react'

import { useEffect, useState } from 'react'

import axios from 'axios'

import { KeyValuePairs } from '@cloudscape-design/components'
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
            <Card className="p-5">
                <Title className='mt-2 mb-4 font-semibold'>
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
            </Card>
        </>
    )
}
