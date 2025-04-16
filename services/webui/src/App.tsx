import { useNavigate } from 'react-router-dom'
import { ArrowLeftStartOnRectangleIcon, ArrowPathIcon, XMarkIcon } from '@heroicons/react/24/outline'
import dayjs from 'dayjs'
import { useEffect, useState } from 'react'
import { useAtom, useAtomValue } from 'jotai'
import { Button, Card, Flex, Text, Title } from '@tremor/react'
import Router from './router'
import Spinner from './components/Spinner'
import { setAuthHeader } from './api/ApiConfig'
import { colorBlindModeAtom, ForbiddenAtom, LayoutAtom, meAtom, RoleAccess, tokenAtom } from './store'
import { applyTheme } from './utilities/theme'
import { OpenGovernance } from './icons/icons'
import { useAuth } from './utilities/auth'
import { useAuthApiV1MeList, useAuthApiV1UserDetail } from './api/auth.gen'
import { PlatformEnginePkgAuthApiTheme } from './api/api'
import { Modal } from '@cloudscape-design/components'
import axios from 'axios'


export default function App() {
    const {
        isLoading,
        isAuthenticated,
        getAccessTokenSilently,
        getIdTokenClaims,
        logout
    } = useAuth()
    const [token, setToken] = useAtom(tokenAtom)
    const [me, setMe] = useAtom(meAtom)
    const [layout, setLayout] = useAtom(LayoutAtom)

    const [accessTokenLoading, setAccessTokenLoading] = useState<boolean>(true)
    const [colorBlindMode, setColorBlindMode] = useAtom(colorBlindModeAtom)
    const [expire, setExpire] = useState<number>(0)
    const [showExpired, setShowExpired] = useState<boolean>(false)
    const forbidden = useAtomValue(ForbiddenAtom)
    const [roleAccess, setRoleAccess] = useAtom(RoleAccess)
    const [layoutLoading, setLayoutLoading] = useState<boolean>(true)
    const {
        response: meResponse,
        isExecuted: meIsExecuted,
        isLoading: meIsLoading,
        error: meError,
        sendNow: getMe,
    } = useAuthApiV1MeList({}, false)

    const checkExpire = () => {
        if (expire !== 0) {
            const diff = expire - dayjs.utc().unix()
            if (diff < 0) {
                setShowExpired(true)
            }
        }
    }
const GetDefaultLayout = (meResponse: any) => {
    axios
        .get(
            `https://raw.githubusercontent.com/opengovern/platform-configuration/refs/heads/main/default_layout.json`
        )
        .then((res) => {
            setLayout(res?.data)
            // SetDefaultLayout(res?.data?.layout_config,meResponse)
            setLayoutLoading(false)
        })
        .catch((err) => {
            setLayoutLoading(false)
        })
}

const GetLayout = (meResponse :any) => {
        setLayoutLoading(true)
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
         const body = {
             user_id: meResponse?.username,
         }

         axios
             .post(`${url}/main/core/api/v4/layout/get-default`, body, config)
             .then((res) => {
                setLayout(res?.data)
                setLayoutLoading(false)
             })
             .catch((err) => {
                 console.log(err)
                //  check if error is 404
                  GetDefaultLayout(meResponse)
                    setLayoutLoading(false)

                   
             })
     }
const SetDefaultLayout = (layout: any, meResponse: any) => {
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
    const body = {
        user_id: meResponse?.username,
        layout_config: layout,
        name: 'Default',
        description: 'Default Layout',
        is_default: true,
        is_private: true,
    }
    console.log(body)

    axios
        .post(`${url}/main/core/api/v4/layout/set`, body, config)
        .then((res) => {})
        .catch((err) => {
            console.log(err)
        })
}

    useEffect(() => {
        const t = setInterval(checkExpire, 5000)
        return () => {
            clearInterval(t)
        }
    }, [expire])


    useEffect(() => {
        if (meIsExecuted && !meIsLoading) {
        
            setMe(meResponse)
            applyTheme(
                    PlatformEnginePkgAuthApiTheme.ThemeLight
            )
            setColorBlindMode( false)
            GetLayout(meResponse)
        }
    }, [meIsLoading])

    useEffect(() => {
        if (isAuthenticated && token === '') {
            getIdTokenClaims().then((v) => {
                setExpire(v?.exp || 0)
            })
            getAccessTokenSilently()
                .then((accessToken) => {
                    setToken(accessToken)
                    setAuthHeader(accessToken)
                    setAccessTokenLoading(false)
                    getMe()
                })
                .catch((err) => {
                    console.error(err)
                    setAccessTokenLoading(false)
                })
        }
    }, [isAuthenticated])

    return isLoading || accessTokenLoading || meIsLoading || layoutLoading ? (
        <Flex
            justifyContent="center"
            alignItems="center"
            className="w-screen h-screen dark:bg-gray-900"
        >
            <Spinner />
        </Flex>
    ) : (
        <>
            <Router />
            {showExpired && (
                <Flex
                    flexDirection="col"
                    className="fixed top-0 left-0 w-screen h-screen bg-gray-900/80 z-50"
                >
                    <Card className="w-1/3 mt-56">
                        <Flex
                            flexDirection="col"
                            justifyContent="center"
                            alignItems="center"
                        >
                            <OpenGovernance className="w-14 h-14 mb-6" />
                            <Title className="mb-3 text-2xl font-bold">
                                Your session has expired
                            </Title>
                            <Text className="mb-6 text-center">
                                Your session has expired. Please log in again to
                                continue accessing opensecurity platform
                            </Text>
                            <Button
                                icon={ArrowPathIcon}
                                onClick={() => {
                                    window.location.href =
                                        window.location.toString()
                                }}
                            >
                                Re-Login
                            </Button>
                        </Flex>
                    </Card>
                </Flex>
            )}

            {forbidden && (
                <Flex
                    flexDirection="col"
                    className="fixed top-0 left-0 w-screen h-screen bg-gray-900/80 z-50"
                >
                    <Card className="w-1/3 mt-56">
                        <Flex
                            flexDirection="col"
                            justifyContent="center"
                            alignItems="center"
                        >
                            <OpenGovernance className="w-14 h-14 mb-6" />
                            <Title className="mb-3 text-2xl font-bold">
                                Access Denied
                            </Title>
                            <Text className="mb-6 text-center">
                                You do not have permission to access the App.
                                Please contact your administrator for more
                                information.
                            </Text>
                            <Button
                                icon={ArrowLeftStartOnRectangleIcon}
                                onClick={() => logout()}
                            >
                                Logout
                            </Button>
                        </Flex>
                    </Card>
                </Flex>
            )}
            <Modal
                visible={roleAccess}
                onDismiss={() => {
                    setRoleAccess(false)
                }}
            >
                <Flex
                    flexDirection="col"
                    justifyContent="center"
                    alignItems="center"
                >
                    <OpenGovernance className="w-14 h-14 mb-6" />
                    <Title className="mb-3 text-2xl font-bold">
                        Access Denied
                    </Title>
                    <Text className="mb-6 text-center">
                        You do not have permission to access this page. Please
                        contact your administrator for more information.
                    </Text>
                    <Button
                        icon={XMarkIcon}
                        onClick={() => {
                            setRoleAccess(false)
                        }}
                    >
                        Close
                    </Button>
                </Flex>
            </Modal>
        </>
    )
}
