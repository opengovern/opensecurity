import {
    Card,
    Flex,

} from '@tremor/react'

import { useEffect, useState } from 'react'

import axios from 'axios'
import ReactMarkdown from 'react-markdown'
import rehypeRaw from 'rehype-raw'
import { useMDXComponents } from '../../../../components/MDX'
import Spinner from '../../../../components/Spinner'
interface IntegrationListProps {
    name?: string
    integration_type?: string
}


export default function Setup({
    name,
    integration_type,
}: IntegrationListProps) {
    const [setup,setSetup] = useState<any>()
    const [loading, setLoading] = useState<boolean>(false)

   
  
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

                <Card className="p-2">
                    <Flex flexDirection='col' className=' p-5 justify-start w-full items-start'>
                        <></>
                    <ReactMarkdown
                        children={setup}
                        skipHtml={false}
                        rehypePlugins={[rehypeRaw]}
                        // @ts-ignore
                        components={useMDXComponents({})}
                    />
                    </Flex>
                </Card>
            )}
        </>
    )
}

