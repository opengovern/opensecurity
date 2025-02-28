import axios from 'axios'
import { atom, useAtom, useSetAtom } from 'jotai'
import { useEffect } from 'react'
import { ForbiddenAtom, RoleAccess } from '../store'


const { hostname } = window.location
export const authHostname = () => {
    if (window.location.origin === 'http://localhost:3000') {
        return window.__RUNTIME_CONFIG__.REACT_APP_AUTH_BASE_URL
    }
    return window.location.origin
}
const apiHostname = () => {
    if (window.location.origin === 'http://localhost:3000') {
        return window.__RUNTIME_CONFIG__.REACT_APP_BASE_URL
    }
    return window.location.origin
}
const instance = axios.create({
    baseURL: `${apiHostname()}${'/main/'}`,
    headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
    },

})
// @ts-ignore
const AxiosInterceptor = ({ children }) => {
    const setForbbiden = useSetAtom(ForbiddenAtom)
    const setRoleAccess = useSetAtom(RoleAccess)


    useEffect(() => {
        // @ts-ignore
        const resInterceptor = (response) => {
            return response
        }
        // @ts-ignore

        const errInterceptor = (error) => {
            if (
                error?.response?.status === 401 
            ) {
                setForbbiden(true)
            }

            if (error?.response?.status === 406) {
                setRoleAccess(true)
            }


            return Promise.reject(error)
        }

        const interceptor = instance.interceptors.response.use(
            resInterceptor,
            errInterceptor
        )

        return () => instance.interceptors.response.eject(interceptor)
    }, [])

    return children
}
export { AxiosInterceptor }
export const setAuthHeader = (authToken?: string) => {
    instance.defaults.headers.common.Authorization = `Bearer ${authToken}`
}

export const setWorkspace = (workspaceName?: string) => {
    instance.defaults.baseURL = `${apiHostname()}/${workspaceName}`
}

export default instance
