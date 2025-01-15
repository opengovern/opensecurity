import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
    Api,
    PlatformEnginePkgWorkspaceApiBootstrapStatusResponse,
    PlatformEnginePkgWorkspaceApiAddCredentialRequest,
    PlatformEnginePkgWorkspaceApiCreateWorkspaceRequest,
    PlatformEnginePkgWorkspaceApiCreateWorkspaceResponse,
    PlatformEnginePkgWorkspaceApiWorkspaceResponse,
    PlatformEnginePkgWorkspaceApiWorkspaceLimitsUsage,
    RequestParams,
} from './api'

import AxiosAPI, { setWorkspace } from './ApiConfig'

interface IuseWorkspaceApiV1BootstrapDetailState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEnginePkgWorkspaceApiBootstrapStatusResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1BootstrapDetail = (
    workspaceName: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseWorkspaceApiV1BootstrapDetailState>({
        isLoading: true,
        isExecuted: false,
    })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([workspaceName, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqworkspaceName: string,
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV1BootstrapDetail(reqworkspaceName, reqparamsSignal)
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

    if (JSON.stringify([workspaceName, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([workspaceName, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, workspaceName, params)
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
        sendRequest(newController, workspaceName, params)
    }

    const sendNowWithParams = (
        reqworkspaceName: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqworkspaceName, reqparams)
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

interface IuseWorkspaceApiV1BootstrapCredentialCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: number
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1BootstrapCredentialCreate = (
    workspaceName: string,
    request: PlatformEnginePkgWorkspaceApiAddCredentialRequest,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseWorkspaceApiV1BootstrapCredentialCreateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([workspaceName, request, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqworkspaceName: string,
        reqrequest: PlatformEnginePkgWorkspaceApiAddCredentialRequest,
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV1BootstrapCredentialCreate(
                    reqworkspaceName,
                    reqrequest,
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
        JSON.stringify([workspaceName, request, params, autoExecute]) !==
        lastInput
    ) {
        setLastInput(
            JSON.stringify([workspaceName, request, params, autoExecute])
        )
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, workspaceName, request, params)
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
        sendRequest(newController, workspaceName, request, params)
    }

    const sendNowWithParams = (
        reqworkspaceName: string,
        reqrequest: PlatformEnginePkgWorkspaceApiAddCredentialRequest,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqworkspaceName, reqrequest, reqparams)
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

interface IuseWorkspaceApiV1BootstrapFinishCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: string
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1BootstrapFinishCreate = (
    workspaceName: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseWorkspaceApiV1BootstrapFinishCreateState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([workspaceName, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqworkspaceName: string,
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV1BootstrapFinishCreate(reqworkspaceName, reqparamsSignal)
                .then((resp) => {
                    setState({
                        ...state,
                        error: undefined,
                        response: resp.data,
                        isLoading: false,
                        isExecuted: true,
                    })
                })
                .catch((err: any) => {
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

    if (JSON.stringify([workspaceName, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([workspaceName, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, workspaceName, params)
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
        sendRequest(newController, workspaceName, params)
    }

    const sendNowWithParams = (
        reqworkspaceName: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqworkspaceName, reqparams)
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

interface IuseWorkspaceApiV1WorkspaceCreateState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEnginePkgWorkspaceApiCreateWorkspaceResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1WorkspaceCreate = (
    request: PlatformEnginePkgWorkspaceApiCreateWorkspaceRequest,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseWorkspaceApiV1WorkspaceCreateState>({
        isLoading: true,
        isExecuted: false,
    })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([request, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqrequest: PlatformEnginePkgWorkspaceApiCreateWorkspaceRequest,
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV1WorkspaceCreate(reqrequest, reqparamsSignal)
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
        reqrequest: PlatformEnginePkgWorkspaceApiCreateWorkspaceRequest,
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

interface IuseWorkspaceApiV1WorkspaceCurrentListState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEnginePkgWorkspaceApiWorkspaceResponse
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1WorkspaceCurrentList = (
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseWorkspaceApiV1WorkspaceCurrentListState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
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
            api.workspace
                .apiV1WorkspaceCurrentList(reqparamsSignal)
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
            sendRequest(newController, params)
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
        sendRequest(newController, params)
    }

    const sendNowWithParams = (reqparams: RequestParams) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqparams)
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

interface IuseWorkspaceApiV1WorkspaceDeleteState {
    isLoading: boolean
    isExecuted: boolean
    response?: void
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1WorkspaceDelete = (
    workspaceId: string,
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseWorkspaceApiV1WorkspaceDeleteState>({
        isLoading: true,
        isExecuted: false,
    })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([workspaceId, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqworkspaceId: string,
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV1WorkspaceDelete(reqworkspaceId, reqparamsSignal)
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

    if (JSON.stringify([workspaceId, params, autoExecute]) !== lastInput) {
        setLastInput(JSON.stringify([workspaceId, params, autoExecute]))
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, workspaceId, params)
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
        sendRequest(newController, workspaceId, params)
    }

    const sendNowWithParams = (
        reqworkspaceId: string,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqworkspaceId, reqparams)
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

interface IuseWorkspaceApiV1WorkspacesListState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEnginePkgWorkspaceApiWorkspaceResponse[]
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1WorkspacesList = (
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const workspace = useParams<{ ws: string }>().ws
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseWorkspaceApiV1WorkspacesListState>({
        isLoading: true,
        isExecuted: false,
    })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV1WorkspacesList(reqparamsSignal)
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
            sendRequest(newController, params)
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
        sendRequest(newController, params)
    }

    const sendNowWithParams = (reqparams: RequestParams) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqparams)
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

interface IuseWorkspaceApiV1WorkspacesLimitsDetailState {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEnginePkgWorkspaceApiWorkspaceLimitsUsage
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV1WorkspacesLimitsDetail = (
    workspaceName: string,
    query?: {
        ignore_usage?: boolean
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
        useState<IuseWorkspaceApiV1WorkspacesLimitsDetailState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([workspaceName, query, params, autoExecute])
    )

    const sendRequest = (
        abortCtrl: AbortController,
        reqworkspaceName: string,
        reqquery:
            | {
                  ignore_usage?: boolean
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV1WorkspacesLimitsDetail(
                    reqworkspaceName,
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
        JSON.stringify([workspaceName, query, params, autoExecute]) !==
        lastInput
    ) {
        setLastInput(
            JSON.stringify([workspaceName, query, params, autoExecute])
        )
    }

    useEffect(() => {
        if (autoExecute) {
            controller.abort()
            const newController = new AbortController()
            setController(newController)
            sendRequest(newController, workspaceName, query, params)
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
        sendRequest(newController, workspaceName, query, params)
    }

    const sendNowWithParams = (
        reqworkspaceName: string,
        reqquery:
            | {
                  ignore_usage?: boolean
              }
            | undefined,
        reqparams: RequestParams
    ) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqworkspaceName, reqquery, reqparams)
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
interface IuseWorkspaceApiV1WorkspacesSampleData {
    isLoading: boolean
    isExecuted: boolean
    response?: PlatformEnginePkgWorkspaceApiWorkspaceLimitsUsage
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

/**
 * URL:
 */
export const useWorkspaceApiV3LoadSampleData = (
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseWorkspaceApiV1WorkspacesLimitsDetailState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV3LoadSampleData(reqparamsSignal)
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
            sendRequest(newController, params)
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
        sendRequest(newController, params)
    }

    const sendNowWithParams = (reqparams: RequestParams) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqparams)
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

export const useWorkspaceApiV3PurgeSampleData = (
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] =
        useState<IuseWorkspaceApiV1WorkspacesLimitsDetailState>({
            isLoading: true,
            isExecuted: false,
        })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV3PurgeSampleData(reqparamsSignal)
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
            sendRequest(newController, params)
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
        sendRequest(newController, params)
    }

    const sendNowWithParams = (reqparams: RequestParams) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqparams)
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
interface IuseWorkspaceApiV3ShouldSetup {
    isLoading: boolean
    isExecuted: boolean
    response?: string
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

export const useWorkspaceApiV3GetShouldSetup = (
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseWorkspaceApiV3ShouldSetup>({
        isLoading: true,
        isExecuted: false,
    })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV3GetSetup(reqparamsSignal)
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
            sendRequest(newController, params)
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
        sendRequest(newController, params)
    }

    const sendNowWithParams = (reqparams: RequestParams) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqparams)
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

interface IuseWorkspaceApiV3Done {
    isLoading: boolean
    isExecuted: boolean
    response?: string
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    error?: any
}

export const useWorkspaceApiV3DoneSetup = (
    params: RequestParams = {},
    autoExecute = true,
    overwriteWorkspace: string | undefined = undefined
) => {
    const [controller, setController] = useState(new AbortController())

    const api = new Api()
    api.instance = AxiosAPI

    const [state, setState] = useState<IuseWorkspaceApiV3Done>({
        isLoading: true,
        isExecuted: false,
    })
    const [lastInput, setLastInput] = useState<string>(
        JSON.stringify([params, autoExecute])
    )

    const sendRequest = (
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
            setWorkspace('main')

            const reqparamsSignal = { ...reqparams, signal: abortCtrl.signal }
            api.workspace
                .apiV3DoneSetup(reqparamsSignal)
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
            sendRequest(newController, params)
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
        sendRequest(newController, params)
    }

    const sendNowWithParams = (reqparams: RequestParams) => {
        controller.abort()
        const newController = new AbortController()
        setController(newController)
        sendRequest(newController, reqparams)
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
