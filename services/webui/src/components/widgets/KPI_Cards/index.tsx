import { Button, Card, CategoryBar, Col, Flex, Grid, Icon, Title } from '@tremor/react'
import { ChevronRightIcon, ShieldCheckIcon } from '@heroicons/react/24/outline'
import Compliance from './Compliance'
import { useNavigate, useParams } from 'react-router-dom'
import Infinity from '../../../icons/Infinity.svg'

export default function SRE() {

    return (
        // <Card className=" border-solid   border-2 border-b w-full rounded-xl border-tremor-border bg-tremor-background-muted p-4 dark:border-dark-tremor-border dark:bg-gray-950 sm:pt-6 sm:pb-10 sm:px-4  ">
        //     <Flex justifyContent="between" className="sm:flex-row flex-col">
        //         <Flex justifyContent="start" className="gap-2 sm:w-fit w-full ">
        //             <img
        //                 className=" w-5 h-5"
        //                 src={Infinity}
        //             />
        //             <Title className=" sm:font-semibold sm:w-fit w-full">
        //                 SRE Frameworks
        //             </Title>
        //         </Flex>
        //     </Flex>
        //     <Grid numItems={1} className="w-full  ">
        //     </Grid>
        // </Card>
        <>
            <Compliance />
        </>
    )
}
