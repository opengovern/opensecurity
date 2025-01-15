import { useAtomValue } from 'jotai'
import { Button, Flex, Title } from '@tremor/react'
import { ReactNode, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ChevronRightIcon } from '@heroicons/react/20/solid'
import {
    kebabCaseToLabel,
    snakeCaseToLabel,
} from '../../../utilities/labelMaker'
import {
    DateRange,
    defaultTime,
    searchAtom,
    useURLParam,
    useURLState,
    useUrlDateRangeState,
} from '../../../utilities/urlstate'

// import { useIntegrationApiV1ConnectionsSummariesList } from '../../../api/integration.gen'

interface IHeader {
    supportedFilters?: string[]
    initialFilters?: string[]
    datePickerDefault?: DateRange
    children?: ReactNode
    breadCrumb?: (string | undefined)[]
    tags?: string[]
    serviceNames?: string[]
}

export default function TopHeader({
    supportedFilters = [],
    initialFilters = [],
    children,
    datePickerDefault,
    breadCrumb,
    tags,
    serviceNames,
}: IHeader) {
    const { ws } = useParams()

    const defaultActiveTimeRange = datePickerDefault || defaultTime(ws || '')
    const { value: activeTimeRange, setValue: setActiveTimeRange } =
        useUrlDateRangeState(defaultActiveTimeRange)
  
    const defaultSelectedConnectors = ''
    const [selectedConnectors, setSelectedConnectors] = useURLParam<
        '' | 'AWS' | 'Azure'
    >('provider', defaultSelectedConnectors)
    const parseConnector = (v: string) => {
        switch (v) {
            case 'AWS':
                return 'AWS'
            case 'Azure':
                return 'Azure'
            default:
                return ''
        }
    }

    const defaultSelectedSeverities = [
        'critical',
        'high',
        'medium',
        'low',
        'none',
    ]
    const [selectedSeverities, setSelectedSeverities] = useURLState<string[]>(
        defaultSelectedSeverities,
        (v) => {
            const res = new Map<string, string[]>()
            res.set('severities', v)
            return res
        },
        (v) => {
            return v.get('severities') || []
        }
    )

    const defaultSelectedCloudAccounts: string[] = []
    const [selectedCloudAccounts, setSelectedCloudAccounts] = useURLState<
        string[]
    >(
        defaultSelectedCloudAccounts,
        (v) => {
            const res = new Map<string, string[]>()
            res.set('connections', v)
            return res
        },
        (v) => {
            return v.get('connections') || []
        }
    )

    const defaultSelectedServiceNames: string[] = []
    const [selectedServiceNames, setSelectedServiceNames] = useURLState<
        string[]
    >(
        defaultSelectedServiceNames,
        (v) => {
            const res = new Map<string, string[]>()
            res.set('serviceNames', v)
            return res
        },
        (v) => {
            return v.get('serviceNames') || []
        }
    )

    const defaultSelectedScoreTags: string[] = []
    const [selectedScoreTags, setSelectedScoreTags] = useURLState<string[]>(
        defaultSelectedScoreTags,
        (v) => {
            const res = new Map<string, string[]>()
            res.set('tags', v)
            return res
        },
        (v) => {
            return v.get('tags') || []
        }
    )

    const defaultSelectedScoreCategory = ''
    const [selectedScoreCategory, setSelectedScoreCategory] =
        useURLState<string>(
            defaultSelectedScoreCategory,
            (v) => {
                const res = new Map<string, string[]>()
                res.set('score_category', [v])
                return res
            },
            (v) => {
                return (v.get('score_category') || []).at(0) || ''
            }
        )

    const calcInitialFilters = () => {
        const resp = initialFilters
        if (activeTimeRange !== defaultActiveTimeRange) {
            resp.push('Date')
        }
        if (selectedConnectors !== defaultSelectedConnectors) {
            resp.push('Connector')
        }
        if (selectedSeverities !== defaultSelectedSeverities) {
            resp.push('Severity')
        }
        if (selectedCloudAccounts !== defaultSelectedCloudAccounts) {
            resp.push('Cloud Account')
        }
        if (selectedServiceNames !== defaultSelectedServiceNames) {
            resp.push('Service Name')
        }
        if (selectedScoreTags !== defaultSelectedScoreTags) {
            resp.push('Tag')
        }
        if (selectedScoreCategory !== defaultSelectedScoreCategory) {
            resp.push('Score Category')
        }

        return resp
    }
    const [addedFilters, setAddedFilters] = useState<string[]>(
        calcInitialFilters()
    )
    const [connectionSearch, setConnectionSearch] = useState('')
    // const { response } = useIntegrationApiV1ConnectionsSummariesList({
    //     connector: selectedConnectors.length ? [selectedConnectors] : [],
    //     pageNumber: 1,
    //     pageSize: 10000,
    //     needCost: false,
    //     needResourceCount: false,
    // })

  


    const navigate = useNavigate()
    const searchParams = useAtomValue(searchAtom)
    const url = window.location.pathname.split('/')
    if (url[1] === 'ws') {
        url.shift()
    }

    const mainPage = () => {
        if (url[1] === 'billing') {
            return 'Usage & Billing'
        }
        if (url[2] === 'score') {
            return 'SCORE'
        }
        if (url[2] === 'spend-metrics') {
            return 'Services'
        }
        if (url[2] === 'infrastructure-metrics') {
            return 'Inventory'
        }
        return url[2] ? kebabCaseToLabel(url[2]) : 'opencomply'
    }

    const subPages = () => {
        const pages = []
        for (let i = 3; i < url.length; i += 1) {
            pages.push(kebabCaseToLabel(url[i]))
        }
        return pages
    }

    const goBack = (n: number) => {
        let temp = '.'
        for (let i = 0; i < n; i += 1) {
            temp += '/..'
        }
        return temp
    }

    document.title = `${mainPage()} `

    return (
        <div className="px-12 2xl:pl-48 z-10 absolute  top-0  left-0 w-full flex h-16 items-center justify-center gap-x-4 border-b border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900 shadow-sm">
            <Flex className="">
                {subPages().length > 0 ? (
                    <Flex justifyContent="start" className="w-fit">
                        <Button
                            onClick={() =>
                                navigate(
                                    `${goBack(
                                        subPages().length > 1
                                            ? subPages().length
                                            : 1
                                    )}?${searchParams}`
                                )
                            }
                            variant="light"
                            className="!text-lg mr-2 hover:text-openg-600"
                        >
                            {mainPage()}
                        </Button>
                        {subPages().map((page, i) => (
                            <Flex
                                key={page}
                                justifyContent="start"
                                className="w-fit mr-2"
                            >
                                <ChevronRightIcon className="h-5 w-5 text-gray-600" />
                                <Button
                                    onClick={() =>
                                        navigate(
                                            `${goBack(
                                                subPages().length - i - 1
                                            )}?${searchParams}`
                                        )
                                    }
                                    variant="light"
                                    className={`${
                                        i === subPages().length - 1
                                            ? 'text-black'
                                            : ''
                                    } opacity-100 ml-2 !text-lg`}
                                    disabled={i === subPages().length - 1}
                                >
                                    {i === subPages().length - 1 &&
                                    breadCrumb?.length
                                        ? breadCrumb
                                        : snakeCaseToLabel(page)}
                                </Button>
                            </Flex>
                        ))}
                    </Flex>
                ) : (
                    <Title className="font-semibold !text-xl whitespace-nowrap">
                        {mainPage()}
                    </Title>
                )}
               
            </Flex>
        </div>
    )
}
