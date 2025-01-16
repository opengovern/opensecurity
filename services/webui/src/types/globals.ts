export {}

declare global {
    interface Window {
        __RUNTIME_CONFIG__: {
            NODE_ENV: string
            REACT_APP_RUNTIME_ENVIRONMENT: string
            REACT_APP_BASE_URL: string
            REACT_APP_AUTH_BASE_URL: string
        }
    }
}
