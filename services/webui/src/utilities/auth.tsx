import dayjs from 'dayjs'
import { atom, useAtom } from 'jotai'
import jwtDecode from 'jwt-decode'
import { authHostname } from '../api/ApiConfig'

interface IAuthModel {
    token: string
    isLoading: boolean
    isSuccessful: boolean
    error: string
    resp: any
}

const sessionAuth = localStorage.getItem('openg_auth')
const sessionAuthModel =
    sessionAuth && sessionAuth.length > 0
        ? (JSON.parse(sessionAuth) as IAuthModel)
        : undefined
const authAtom = atom<IAuthModel>(
    sessionAuthModel || {
        token: '',
        isLoading: false,
        isSuccessful: false,
        error: '',
        resp: {},
    }
)

interface JwtPayload {
    iss?: string
    sub?: string
    aud?: string[] | string
    exp?: number
    nbf?: number
    iat?: number
    jti?: string
    email?: string
}

let authLoading = false
export function useAuth() {
    const [auth, setAuth] = useAtom(authAtom)
    const decodedToken =
        auth.token === undefined || auth.token === ''
            ? undefined
            : jwtDecode<JwtPayload>(auth.token)

    return {
        isLoading: auth.isLoading,
        isAuthenticated:
            auth.isSuccessful &&
            auth.token !== undefined &&
            auth.token !== '' &&
            dayjs.unix(decodedToken?.exp || 0).isAfter(dayjs()),
        getAccessTokenSilently: () => {
            if (auth.isSuccessful) {
                return Promise.resolve(auth.token)
            }
            return Promise.reject(Error('not authenticated'))
        },
        getIdTokenClaims: () => {
            if (auth.isSuccessful) {
                return Promise.resolve({
                    exp: decodedToken?.exp,
                })
            }
            return Promise.reject(Error('not authenticated'))
        },
        error: {
            message: auth.error,
        },
        user: {
            given_name: '', // TODO-Saleh
            family_name: '',
            name: '',
            email: decodedToken?.email,
            picture: '',
        },
        logout: () => {
            const newAuth = {
                token: '',
                isLoading: false,
                isSuccessful: false,
                error: '',
                resp: {},
            }
            setAuth(newAuth)
            localStorage.setItem('openg_auth', JSON.stringify(newAuth))
            window.location.href = '/'
        },
        loginWithCode: (code: string) => {
            if (code.length === 0) {
                return Promise.resolve()
            }

            if (authLoading) {
                return Promise.resolve()
            }
            authLoading = true

            const getToken = async (retryCount: number) => {
                setAuth({
                    ...auth,
                    isLoading: true,
                    isSuccessful: false,
                    error: '',
                })

                const callback = `${window.location.origin}/callback`
                const url = `${authHostname()}/dex/token`
                const headers = new Headers()
                headers.append(
                    'Content-Type',
                    'application/x-www-form-urlencoded'
                )

                const body = new URLSearchParams()
                body.append('grant_type', 'authorization_code')
                body.append('client_id', 'public-client')
                body.append('client_secret', '')
                body.append('code', code)
                body.append('redirect_uri', callback)

                try {
                    const response = await fetch(url, {
                        method: 'POST',
                        headers,
                        body,
                    })

                    const data = await response.json()
                    if (data.error) {
                        if (retryCount < 3) {
                            getToken(retryCount + 1)
                        } else {
                            console.log(
                                `Failed to fetch token due to ${data.error}`
                            )
                            setAuth({
                                ...auth,
                                isLoading: false,
                                isSuccessful: false,
                                error: data.error_description,
                            })
                        }
                    } else {
                        const newAuth = {
                            token: data.access_token,
                            isLoading: false,
                            isSuccessful: true,
                            error: '',
                            resp: data,
                        }
                        setAuth(newAuth)
                        localStorage.setItem(
                            'openg_auth',
                            JSON.stringify(newAuth)
                        )
                    }
                } catch (error) {
                    if (retryCount < 3) {
                        getToken(retryCount + 1)
                    } else {
                        console.log(`Failed to fetch token due to ${error}`)
                        setAuth({
                            ...auth,
                            isLoading: false,
                            isSuccessful: false,
                            error: `Failed to fetch token due to ${error}`,
                        })
                    }
                }
            }

            return getToken(0).finally(() => {
                authLoading = false
            })
        },
    }
}
