import { useEffect, useState } from 'react'
import { Button, Card, Flex, Text, TextInput, Title } from '@tremor/react'
import {
    ChevronUpDownIcon,
    MagnifyingGlassIcon,
    XMarkIcon,
} from '@heroicons/react/24/outline'
import { Checkbox, Radio, useCheckboxState } from 'pretty-checkbox-react'
import { useIntegrationApiV1ConnectionsSummariesList } from '../../../../api/integration.gen'
import Spinner from '../../../Spinner'
import { AWSIcon, AzureIcon } from '../../../../icons/icons'
import { useFilterState } from '../../../../utilities/urlstate'
import {
    GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection,
    SourceType,
} from '../../../../api/api'

export const compareArrays = (a: any[], b: any[]) => {
    if (a && b) {
        return (
            a.length === b.length &&
            a.every((element: any, index: number) => element === b[index])
        )
    }
    return undefined
}

const filteredConnectionsList = (
    connection:
        | GithubComKaytuIoKaytuEngineServicesIntegrationApiEntityConnection[]
        | undefined,
    filter: string
) => {
    const result = connection?.filter(
        (c) =>
            c?.providerConnectionName
                ?.toLowerCase()
                .includes(filter.toLowerCase()) ||
            c?.providerConnectionID
                ?.toLowerCase()
                .includes(filter.toLowerCase())
    )
    const count = result?.length || 0
    return {
        result,
        count,
    }
}

export default function Filter() {
    const [openDrawer, setOpenDrawer] = useState(false)
    const { value: selectedFilters, setValue: setSelectedFilters } =
        useFilterState()

    const { response, isLoading } = useIntegrationApiV1ConnectionsSummariesList(
        {
            connector: [],
            pageNumber: 1,
            pageSize: 10000,
            needCost: false,
            needResourceCount: false,
        }
    )
    const connectionCheckbox = useCheckboxState({
        state: [...selectedFilters.connections],
    })
    const [provider, setProvider] = useState(selectedFilters.provider)
    const [search, setSearch] = useState('')

    const findConnections = () => {
        const conn = []
        if (response) {
            for (let i = 0; i < selectedFilters.connections.length; i += 1) {
                conn.push(
                    response?.connections?.find(
                        (c) => c.id === selectedFilters.connections[i]
                    )
                )
            }
        }
        return conn
    }

    const restFilters = () => {
        setProvider(selectedFilters.provider)
        connectionCheckbox.setState([...selectedFilters.connections])
    }

    const btnDisable = () =>
        provider === selectedFilters.provider &&
        compareArrays(
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            connectionCheckbox.state.sort(),
            selectedFilters.connections.sort()
        )

    const filterText = () => {
        if (selectedFilters.connections.length > 0) {
            if (selectedFilters.connections.length === 1) {
                return findConnections()[0]?.providerConnectionName
            }
            return `${selectedFilters.connections.length} connections`
        }
        if (selectedFilters.provider !== '') {
            return selectedFilters.provider
        }
        return 'Scope selector'
    }

    useEffect(() => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        if (connectionCheckbox.state.length > 0) {
            setProvider(SourceType.Nil)
        }
    }, [connectionCheckbox.state])

    return (
        <div className="relative">
            <Button
                variant={
                    filterText() === 'Scope selector' ? 'secondary' : 'primary'
                }
                className="h-9"
                onClick={() => setOpenDrawer(true)}
                icon={ChevronUpDownIcon}
            >
                {filterText()}
            </Button>
            {openDrawer && (
                <>
                    <button
                        type="button"
                        onClick={() => {
                            restFilters()
                            setOpenDrawer(false)
                            setSearch('')
                        }}
                        className="cursor-default opacity-0 top-0 left-0 z-10 fixed w-screen h-screen"
                    >
                        filters
                    </button>
                    <Card className="p-0 border border-gray-300 shadow-lg mt-2.5 z-20 absolute top-full right-0 w-[337px]">
                        <Flex className="px-6 py-2 border-b border-b-gray-200">
                            <Title>Scope</Title>
                            <XMarkIcon
                                className="w-5"
                                onClick={() => {
                                    restFilters()
                                    setOpenDrawer(false)
                                }}
                            />
                        </Flex>
                        <Flex
                            flexDirection="col"
                            alignItems="start"
                            justifyContent="start"
                            className="p-6"
                        >
                            <Text className="mb-2">PROVIDER</Text>
                            <Flex
                                flexDirection="col"
                                alignItems="start"
                                className="gap-1.5"
                            >
                                <Radio
                                    name="scope-providoe"
                                    onClick={() => {
                                        setProvider(SourceType.Nil)
                                        connectionCheckbox.setState([])
                                    }}
                                    checked={provider === ''}
                                >
                                    <Text className="text-gray-800">All</Text>
                                </Radio>
                                <Radio
                                    name="scope-providoe"
                                    onClick={() => {
                                        setProvider(SourceType.CloudAWS)
                                        connectionCheckbox.setState([])
                                    }}
                                    checked={provider === 'AWS'}
                                >
                                    <Flex className="gap-1">
                                        <img
                                            src={AWSIcon}
                                            className="w-6 rounded-full"
                                            alt="aws"
                                        />
                                        <Text className="text-gray-800">
                                            AWS
                                        </Text>
                                    </Flex>
                                </Radio>
                                <Radio
                                    name="scope-providoe"
                                    onClick={() => {
                                        setProvider(SourceType.CloudAzure)
                                        connectionCheckbox.setState([])
                                    }}
                                    checked={provider === 'Azure'}
                                >
                                    <Flex className="gap-1">
                                        <img
                                            src={AzureIcon}
                                            className="w-6 rounded-full"
                                            alt="azure"
                                        />
                                        <Text className="text-gray-800">
                                            Azure
                                        </Text>
                                    </Flex>
                                </Radio>
                            </Flex>
                            <Flex className="mt-4 mb-2">
                                <Text>ACCOUNTS</Text>
                                {selectedFilters.connections.length > 0 && (
                                    <Text className="!text-xs">
                                        ({selectedFilters.connections.length}{' '}
                                        selected)
                                    </Text>
                                )}
                            </Flex>
                            <TextInput
                                icon={MagnifyingGlassIcon}
                                placeholder="Search account by ID or name..."
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                            />
                            <Flex
                                flexDirection="col"
                                alignItems="start"
                                justifyContent="start"
                                className="gap-2 mt-3 h-[210px] overflow-y-scroll"
                            >
                                {isLoading ? (
                                    <Spinner />
                                ) : (
                                    filteredConnectionsList(
                                        response?.connections,
                                        search
                                    ).result?.map(
                                        (con, i) =>
                                            i < 100 && (
                                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                                // @ts-ignore
                                                <Checkbox
                                                    shape="curve"
                                                    className="!items-start"
                                                    value={con.id}
                                                    {...connectionCheckbox}
                                                >
                                                    <Flex
                                                        flexDirection="col"
                                                        alignItems="start"
                                                        className="-mt-0.5"
                                                    >
                                                        <Text className="text-gray-800 truncate">
                                                            {
                                                                con.providerConnectionName
                                                            }
                                                        </Text>
                                                        <Text className="text-xs truncate">
                                                            {
                                                                con.providerConnectionID
                                                            }
                                                        </Text>
                                                    </Flex>
                                                </Checkbox>
                                            )
                                    )
                                )}
                            </Flex>
                        </Flex>
                        <Flex className="border-t border-t-gray-200 px-6 py-2">
                            {(!!selectedFilters.provider.length ||
                                !!selectedFilters.connections.length ||
                                !!selectedFilters.connectionGroup.length) && (
                                <Button
                                    variant="light"
                                    onClick={() => {
                                        setSearch('')
                                        setProvider(SourceType.Nil)
                                        connectionCheckbox.setState([])
                                        setSelectedFilters({
                                            provider: SourceType.Nil,
                                            connections: [],
                                            connectionGroup: [],
                                        })
                                        setOpenDrawer(false)
                                    }}
                                    className="whitespace-nowrap"
                                >
                                    Reset filters
                                </Button>
                            )}
                            <Flex justifyContent="end" className="gap-4">
                                <Button
                                    variant="light"
                                    onClick={() => {
                                        restFilters()
                                        setOpenDrawer(false)
                                    }}
                                >
                                    Cancel
                                </Button>
                                <Button
                                    onClick={() => {
                                        setSelectedFilters({
                                            provider,
                                            connections: [
                                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                                // @ts-ignore
                                                ...connectionCheckbox.state,
                                            ],
                                            connectionGroup: [],
                                        })
                                        setOpenDrawer(false)
                                    }}
                                    disabled={btnDisable()}
                                >
                                    Apply
                                </Button>
                            </Flex>
                        </Flex>
                    </Card>
                </>
            )}
        </div>
    )
}
