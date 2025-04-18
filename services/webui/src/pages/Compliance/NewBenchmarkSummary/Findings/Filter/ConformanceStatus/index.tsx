import { Radio } from 'pretty-checkbox-react'
import { Button, Flex, Text } from '@tremor/react'
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline'
import { PlatformEnginePkgComplianceApiConformanceStatus } from '../../../../../../../api/api'
import Multiselect from '@cloudscape-design/components/multiselect'
import { useEffect, useState } from 'react'

interface IConformanceStatus {
    value:
        | PlatformEnginePkgComplianceApiConformanceStatus[]
        | undefined
    defaultValue:
        | PlatformEnginePkgComplianceApiConformanceStatus[]
        | undefined
    onChange: (
        c:
            | PlatformEnginePkgComplianceApiConformanceStatus[]
            | undefined
    ) => void
}

export default function ConformanceStatus({
    value,
    defaultValue,
    onChange,
}: IConformanceStatus) {
    const options = [
      
        {
            label: 'Failed',
            value: 
                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
            
            iconSvg: <XCircleIcon className="h-5 text-rose-600" />,
        },
        {
            label: 'Passed',
            value: 
                PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed,
            
            iconSvg: <CheckCircleIcon className="h-5 text-emerald-500" />,
        },
    ]
  const [selectedOptions, setSelectedOptions] = useState([
      {
          label: 'Failed',
          value: 
              PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
          
          iconSvg: <XCircleIcon className="h-5 text-rose-600" />,
      },
  ])
   useEffect(() => {
       if (selectedOptions.length === 0) {
           onChange(defaultValue)
           return
       } else {
           // @ts-ignore
           const temp = []
           selectedOptions.map((o) => {
               // @ts-ignore

               temp.push(o.value)
           })
           // @ts-ignore
           onChange(temp)
           // @ts-ignore
       }
   }, [selectedOptions])
    return (
        <>
            <Multiselect
                // @ts-ignore
                selectedOptions={selectedOptions}
                onChange={({ detail }) =>
                    // @ts-ignore
                    setSelectedOptions(detail.selectedOptions)
                }
                // @ts-ignore
                options={options}
                // filteringType="auto"
                tokenLimit={0}
                placeholder="Incident"
            />
            {/* {options.map((o) => (
                <Radio
                    name="conformance_status"
                    key={`conformance_status-${o.name}`}
                    checked={compareArrays(o.value.sort(), value?.sort() || [])}
                    onClick={() => onChange(o.value)}
                >
                    <Flex className="gap-1 w-fit">
                        {o.icon}
                        <Text>{o.name}</Text>
                    </Flex>
                </Radio>
            ))} */}
            {/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
            {/* @ts-ignore */}
            {/* {!compareArrays(value?.sort(), defaultValue?.sort()) && (
                <Flex className="pt-3 mt-3 border-t border-t-gray-200">
                    <Button
                        variant="light"
                        onClick={() => onChange(defaultValue)}
                    >
                        Reset
                    </Button>
                </Flex>
            )} */}
        </>
    )
}
