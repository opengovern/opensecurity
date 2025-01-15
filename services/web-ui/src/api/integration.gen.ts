import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
    Api,
    PlatformEngineServicesIntegrationApiEntityListConnectionsSummaryResponse,
    PlatformEngineServicesIntegrationApiEntityCreateCredentialResponse,
    PlatformEngineServicesIntegrationApiEntityCredential,
    PlatformEngineServicesIntegrationApiEntityCreateAWSConnectionRequest,
    PlatformEngineServicesIntegrationApiEntityCountConnectionsResponse,
    PlatformEngineServicesIntegrationApiEntityConnection,
    PlatformEngineServicesIntegrationApiEntityConnectorCount,
    PlatformEngineServicesIntegrationApiEntityUpdateAWSCredentialRequest,
    PlatformEngineServicesIntegrationApiEntityCreateAzureCredentialRequest,
    PlatformEngineServicesIntegrationApiEntityListCredentialResponse,
    PlatformEngineServicesIntegrationApiEntityCreateAWSCredentialRequest,
    PlatformEngineServicesIntegrationApiEntityUpdateAzureCredentialRequest,
    PlatformEngineServicesIntegrationApiEntityCreateConnectionResponse,
    PlatformEngineServicesIntegrationApiEntityCatalogMetrics,
    RequestParams,
    PlatformEngineServicesIntegrationApiEntityConnectorResponse,
} from './api'

import AxiosAPI, { setWorkspace } from './ApiConfig'

interface IuseIntegrationApiV1ConnectionsAwsCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityCreateConnectionResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectionsAwsCreate = (
    request: PlatformEngineServicesIntegrationApiEntityCreateAWSConnectionRequest,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1ConnectionsAwsCreateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([request, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqrequest: PlatformEngineServicesIntegrationApiEntityCreateAWSConnectionRequest,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectionsAwsCreate(reqrequest, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([request, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([request, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, request, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, request, params)
    }

    const sendNowWithParams = (
        reqrequest: PlatformEngineServicesIntegrationApiEntityCreateAWSConnectionRequest,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqrequest, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1ConnectionsCountListState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityCountConnectionsResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectionsCountList = (
    query?: {
        connector?: string
    },
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1ConnectionsCountListState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([query, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqquery:
            | {
                  connector?: string
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectionsCountList(reqquery, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([query, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([query, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, query, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, query, params)
    }

    const sendNowWithParams = (
        reqquery:
            | {
                  connector?: string
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqquery, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1ConnectionsSummariesListState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityListConnectionsSummaryResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectionsSummariesList = (
    query?: {
        filter?: string

        connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]

        connectionId?: string[]

        resourceCollection?: string[]

        connectionGroups?: string[]

        credentialType?: (
            | 'auto-azure'
            | 'auto-aws'
            | 'manual-aws-org'
            | 'manual-azure-spn'
        )[]

        lifecycleState?:
            | 'DISABLED'
            | 'DISCOVERED'
            | 'IN_PROGRESS'
            | 'ONBOARD'
            | 'ARCHIVED'

        healthState?: 'healthy' | 'unhealthy'

        pageSize?: number

        pageNumber?: number

        startTime?: number

        endTime?: number

        needCost?: boolean

        needResourceCount?: boolean

        sortBy?:
            | 'onboard_date'
            | 'resource_count'
            | 'cost'
            | 'growth'
            | 'growth_rate'
            | 'cost_growth'
            | 'cost_growth_rate'
    },
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1ConnectionsSummariesListState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([query, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqquery:
            | {
                  filter?: string

                  connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]

                  connectionId?: string[]

                  resourceCollection?: string[]

                  connectionGroups?: string[]

                  credentialType?: (
                      | 'auto-azure'
                      | 'auto-aws'
                      | 'manual-aws-org'
                      | 'manual-azure-spn'
                  )[]

                  lifecycleState?:
                      | 'DISABLED'
                      | 'DISCOVERED'
                      | 'IN_PROGRESS'
                      | 'ONBOARD'
                      | 'ARCHIVED'

                  healthState?: 'healthy' | 'unhealthy'

                  pageSize?: number

                  pageNumber?: number

                  startTime?: number

                  endTime?: number

                  needCost?: boolean

                  needResourceCount?: boolean

                  sortBy?:
                      | 'onboard_date'
                      | 'resource_count'
                      | 'cost'
                      | 'growth'
                      | 'growth_rate'
                      | 'cost_growth'
                      | 'cost_growth_rate'
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectionsSummariesList(reqquery, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([query, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([query, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, query, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, query, params)
    }

    const sendNowWithParams = (
        reqquery:
            | {
                  filter?: string

                  connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]

                  connectionId?: string[]

                  resourceCollection?: string[]

                  connectionGroups?: string[]

                  credentialType?: (
                      | 'auto-azure'
                      | 'auto-aws'
                      | 'manual-aws-org'
                      | 'manual-azure-spn'
                  )[]

                  lifecycleState?:
                      | 'DISABLED'
                      | 'DISCOVERED'
                      | 'IN_PROGRESS'
                      | 'ONBOARD'
                      | 'ARCHIVED'

                  healthState?: 'healthy' | 'unhealthy'

                  pageSize?: number

                  pageNumber?: number

                  startTime?: number

                  endTime?: number

                  needCost?: boolean

                  needResourceCount?: boolean

                  sortBy?:
                      | 'onboard_date'
                      | 'resource_count'
                      | 'cost'
                      | 'growth'
                      | 'growth_rate'
                      | 'cost_growth'
                      | 'cost_growth_rate'
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqquery, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1ConnectionsDeleteState {
    isLoading: boolean
    isExecuted: boolean
    response?: void
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectionsDelete = (
    connectionId: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1ConnectionsDeleteState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([connectionId, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqconnectionId: string,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectionsDelete(reqconnectionId, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([connectionId, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([connectionId, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, connectionId, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, connectionId, params)
    }

    const sendNowWithParams = (
        reqconnectionId: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqconnectionId, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1ConnectionsAwsHealthcheckDetailState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityConnection
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectionsAwsHealthcheckDetail = (
    connectionId: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1ConnectionsAwsHealthcheckDetailState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([connectionId, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqconnectionId: string,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectionsAwsHealthcheckDetail(
                    reqconnectionId,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([connectionId, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([connectionId, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, connectionId, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, connectionId, params)
    }

    const sendNowWithParams = (
        reqconnectionId: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqconnectionId, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1ConnectionsAzureHealthcheckDetailState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityConnection
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectionsAzureHealthcheckDetail = (
    connectionId: string,
    query?: {
        updateMetadata?: boolean
    },
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1ConnectionsAzureHealthcheckDetailState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([connectionId, query, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqconnectionId: string,
        reqquery:
            | {
                  updateMetadata?: boolean
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectionsAzureHealthcheckDetail(
                    reqconnectionId,
                    reqquery,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (
        JSON.stringify([connectionId, query, params, autoExecute]) !== lastInput
    ) {
        setLastInput(JSON.stringify([connectionId, query, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, connectionId, query, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, connectionId, query, params)
    }

    const sendNowWithParams = (
        reqconnectionId: string,
        reqquery:
            | {
                  updateMetadata?: boolean
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqconnectionId, reqquery, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1ConnectorsListState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityConnectorResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectorsList = (
    perPage: number,
    cursor: number,
    hasIntegration?: boolean,
    sortBy?: string,
    sortOrder?: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseIntegrationApiV1ConnectorsListState>(
        {
            isLoading: true,
            isExecuted: false,
        }
    )
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
        reqperPage: number,
        reqcursor: number,
        reqSortBy: string | undefined,
        reqSortOrder: string | undefined,
        reqHasIntegration: boolean | undefined,
        abortCtrl: AbortController,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectorsList(
                    reqperPage,
                    reqcursor,
                    reqSortBy,
                    reqSortOrder,
                    reqHasIntegration,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(
                perPage,
                cursor,
                sortBy,
                sortOrder,
                hasIntegration,
                newController,
                params
            )
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = (
        reqperPage: number,
        reqCursor: number,
        reqSortBy: string,
        reqSortOrder: string,
        reqHasIntegration: boolean
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(
            reqperPage,
            reqCursor,
            reqSortBy,
            reqSortOrder,
            reqHasIntegration,
            newController,
            params
        )
    }

    const sendNowWithParams = (
        reqPerPage: number,
        reqCursor: number,
        reqSortBy: string,
        reqSortOrder: string,
        reqHasIntegration: boolean,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(
            reqPerPage,
            reqCursor,
            reqSortBy,
            reqSortOrder,
            reqHasIntegration,
            newController,
            reqparams
        )
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

export const useIntegrationApiV1EnabledConnectorsList = (
    perPage: number,
    cursor: number,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseIntegrationApiV1ConnectorsListState>(
        {
            isLoading: true,
            isExecuted: false,
        }
    )
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
        reqperPage: number,
        reqcursor: number,
        abortCtrl: AbortController,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1EnabledConnectorsList(
                    reqperPage,
                    reqcursor,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(perPage, cursor, newController, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = (reqperPage: number, reqCursor: number) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(reqperPage, reqCursor, newController, params)
    }

    const sendNowWithParams = (
        reqPerPage: number,
        reqCursor: number,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(reqPerPage, reqCursor, newController, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1ConnectorsMetricsListState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityCatalogMetrics
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1ConnectorsMetricsList = (
    query?: {
        connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]

        credentialType?: (
            | 'auto-azure'
            | 'auto-aws'
            | 'manual-aws-org'
            | 'manual-azure-spn'
        )[]
    },
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1ConnectorsMetricsListState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([query, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqquery:
            | {
                  connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]

                  credentialType?: (
                      | 'auto-azure'
                      | 'auto-aws'
                      | 'manual-aws-org'
                      | 'manual-azure-spn'
                  )[]
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1ConnectorsMetricsList(reqquery, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([query, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([query, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, query, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, query, params)
    }

    const sendNowWithParams = (
        reqquery:
            | {
                  connector?: ('' | 'AWS' | 'Azure' | 'EntraID')[]

                  credentialType?: (
                      | 'auto-azure'
                      | 'auto-aws'
                      | 'manual-aws-org'
                      | 'manual-azure-spn'
                  )[]
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqquery, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialDeleteState {
    isLoading: boolean
    isExecuted: boolean
    response?: void
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialDelete = (
    credentialId: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialDeleteState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([credentialId, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialDelete(reqcredentialId, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([credentialId, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([credentialId, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, credentialId, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, credentialId, params)
    }

    const sendNowWithParams = (
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqcredentialId, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsListState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityListCredentialResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsList = (
    query?: {
        connector?: '' | 'AWS' | 'Azure' | 'EntraID'

        health?: 'healthy' | 'unhealthy'

        credentialType?: (
            | 'auto-azure'
            | 'auto-aws'
            | 'manual-aws-org'
            | 'manual-azure-spn'
        )[]

        pageSize?: number

        pageNumber?: number
    },
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsListState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([query, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqquery:
            | {
                  connector?: '' | 'AWS' | 'Azure' | 'EntraID'

                  health?: 'healthy' | 'unhealthy'

                  credentialType?: (
                      | 'auto-azure'
                      | 'auto-aws'
                      | 'manual-aws-org'
                      | 'manual-azure-spn'
                  )[]

                  pageSize?: number

                  pageNumber?: number
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsList(reqquery, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([query, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([query, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, query, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, query, params)
    }

    const sendNowWithParams = (
        reqquery:
            | {
                  connector?: '' | 'AWS' | 'Azure' | 'EntraID'

                  health?: 'healthy' | 'unhealthy'

                  credentialType?: (
                      | 'auto-azure'
                      | 'auto-aws'
                      | 'manual-aws-org'
                      | 'manual-azure-spn'
                  )[]

                  pageSize?: number

                  pageNumber?: number
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqquery, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsAwsCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityCreateCredentialResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsAwsCreate = (
    request: PlatformEngineServicesIntegrationApiEntityCreateAWSCredentialRequest,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsAwsCreateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([request, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqrequest: PlatformEngineServicesIntegrationApiEntityCreateAWSCredentialRequest,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsAwsCreate(reqrequest, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([request, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([request, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, request, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, request, params)
    }

    const sendNowWithParams = (
        reqrequest: PlatformEngineServicesIntegrationApiEntityCreateAWSCredentialRequest,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqrequest, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsAwsUpdateState {
    isLoading: boolean
    isExecuted: boolean
    response?: void
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsAwsUpdate = (
    credentialId: string,
    config: PlatformEngineServicesIntegrationApiEntityUpdateAWSCredentialRequest,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsAwsUpdateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([credentialId, config, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqcredentialId: string,
        reqconfig: PlatformEngineServicesIntegrationApiEntityUpdateAWSCredentialRequest,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsAwsUpdate(
                    reqcredentialId,
                    reqconfig,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (
        JSON.stringify([credentialId, config, params, autoExecute]) !==
        lastInput
    ) {
        setLastInput(
            JSON.stringify([credentialId, config, params, autoExecute])
        )
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, credentialId, config, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, credentialId, config, params)
    }

    const sendNowWithParams = (
        reqcredentialId: string,
        reqconfig: PlatformEngineServicesIntegrationApiEntityUpdateAWSCredentialRequest,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqcredentialId, reqconfig, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsAwsAutoonboardCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityConnection[]
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsAwsAutoonboardCreate = (
    credentialId: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsAwsAutoonboardCreateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([credentialId, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsAwsAutoonboardCreate(
                    reqcredentialId,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([credentialId, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([credentialId, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, credentialId, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, credentialId, params)
    }

    const sendNowWithParams = (
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqcredentialId, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsAzureCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityCreateCredentialResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsAzureCreate = (
    request: PlatformEngineServicesIntegrationApiEntityCreateAzureCredentialRequest,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsAzureCreateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([request, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqrequest: PlatformEngineServicesIntegrationApiEntityCreateAzureCredentialRequest,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsAzureCreate(reqrequest, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([request, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([request, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, request, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, request, params)
    }

    const sendNowWithParams = (
        reqrequest: PlatformEngineServicesIntegrationApiEntityCreateAzureCredentialRequest,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqrequest, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsAzureUpdateState {
    isLoading: boolean
    isExecuted: boolean
    response?: void
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsAzureUpdate = (
    credentialId: string,
    config: PlatformEngineServicesIntegrationApiEntityUpdateAzureCredentialRequest,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsAzureUpdateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([credentialId, config, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqcredentialId: string,
        reqconfig: PlatformEngineServicesIntegrationApiEntityUpdateAzureCredentialRequest,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsAzureUpdate(
                    reqcredentialId,
                    reqconfig,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (
        JSON.stringify([credentialId, config, params, autoExecute]) !==
        lastInput
    ) {
        setLastInput(
            JSON.stringify([credentialId, config, params, autoExecute])
        )
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, credentialId, config, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, credentialId, config, params)
    }

    const sendNowWithParams = (
        reqcredentialId: string,
        reqconfig: PlatformEngineServicesIntegrationApiEntityUpdateAzureCredentialRequest,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqcredentialId, reqconfig, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsAzureAutoonboardCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityConnection[]
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsAzureAutoonboardCreate = (
    credentialId: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsAzureAutoonboardCreateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([credentialId, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsAzureAutoonboardCreate(
                    reqcredentialId,
                    reqparamsSignal
                )
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([credentialId, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([credentialId, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, credentialId, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, credentialId, params)
    }

    const sendNowWithParams = (
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqcredentialId, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}

interface IuseIntegrationApiV1CredentialsDetailState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEngineServicesIntegrationApiEntityCredential
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useIntegrationApiV1CredentialsDetail = (
    credentialId: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseIntegrationApiV1CredentialsDetailState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([credentialId, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        if (!api.instance.defaults.headers.common.Authorization) {
            return
        }

        setState({
            ...state,
            error: undefined,
            isLoading: true,
            isExecuted: true,
        })
        try {
            if (overwriteWorkspace) {
                setWorkspace(overwriteWorkspace)
            } else if (workspace !== undefined && workspace.length > 0) {
                setWorkspace(workspace)
            } else {
                setWorkspace('main')
            }

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.integration
                .apiV1CredentialsDetail(reqcredentialId, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err) => {
                    if (
                        err.name === 'AbortError' ||
                        err.name === 'CanceledError'
                    ) {
                        // Request was aborted
                    } else {
                        setState({
                            ...state,
                            error: err,
                            response: undefined,
                            isLoading: false,
                            isExecuted: true,
                        })
                    }
                })
        } catch (err) {
            setState({
                ...state,
                error: err,
                isLoading: false,
                isExecuted: true,
            })
        }
    }

    if (JSON.stringify([credentialId, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([credentialId, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, credentialId, params)
        }
    }, [lastInput])

    const { response } = state
    const { isLoading } = state
    const { isExecuted } = state
    const { error } = state
    const sendNow = () => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, credentialId, params)
    }

    const sendNowWithParams = (
        reqcredentialId: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqcredentialId, reqparams)
    }

    return {
        response,
        isLoading,
        isExecuted,
        error,
        sendNow,
        sendNowWithParams,
    }
}
