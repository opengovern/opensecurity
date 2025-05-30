import { Card, Flex, Tab, TabGroup, TabList, Text } from '@tremor/react'
import {
    ArrowTopRightOnSquareIcon,
    Bars2Icon,
    Bars3Icon,
    Cog6ToothIcon,
    MagnifyingGlassIcon,
    PuzzlePieceIcon,
} from '@heroicons/react/24/outline'
import { Fragment, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useAtomValue, useSetAtom } from 'jotai'
import { notificationAtom, workspaceAtom } from '../../../../../store'
import { PlatformEnginePkgAuthApiTheme } from '../../../../../api/api'
import { applyTheme, currentTheme } from '../../../../../utilities/theme'
import { useAuthApiV1UserPreferencesUpdate } from '../../../../../api/auth.gen'
import { useAuth } from '../../../../../utilities/auth'
import axios from 'axios'
import { Alert, Button, Modal, Popover } from '@cloudscape-design/components'
import FormField from '@cloudscape-design/components/form-field'
import Input from '@cloudscape-design/components/input'

export default function Profile() {
    const navigate = useNavigate()
    const { user, logout } = useAuth()

    const setNotification = useSetAtom(notificationAtom)

    const [index, setIndex] = useState(
        // eslint-disable-next-line no-nested-ternary
        currentTheme() === 'system' ? 2 : currentTheme() === 'dark' ? 1 : 0
    )
    const [isPageLoading, setIsPageLoading] = useState<boolean>(true)
    const [theme, setTheme] =
        useState<PlatformEnginePkgAuthApiTheme>(currentTheme())
    const [change, setChange] = useState<boolean>(false)
    const [password, setPassword] = useState<any>({
        current: '',
        new: '',
        confirm: '',
    })
    const [errors, setErrors] = useState<any>({
        current: '',
        new: '',
        confirm: '',
    })
     const [changeError, setChangeError] = useState()
     const [loadingChange, setLoadingChange] = useState(false)

    const { sendNow } = useAuthApiV1UserPreferencesUpdate(
        {
            theme,
        },
        {},
        false
    )
 const ChangePassword = () => {
     if (!password.current || password.current == '') {
         setErrors({ ...errors, current: 'Please enter current password' })
         return
     }
     if (!password.new || password.new == '') {
         setErrors({
             ...errors,
             new: 'Please enter new password',
         })
         return
     }
     if (!password.confirm || password.confirm == '') {
         setErrors({ ...errors, confirm: 'Please enter confirm password' })
         return
     }
     if (password.confirm !== password.new) {
         setErrors({
             ...errors,
             confirm: 'Passwords are not same',
             new: 'Passwords are not same',
         })
         return
     }
     if (password.current === password.new) {
         setErrors({
             ...errors,
             current: 'Passwords are  same',
             new: 'Passwords are  same',
         })
         return
     }
        setLoadingChange(true)

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
         current_password: password?.current,
         new_password: password?.new,
     }
     axios
         .post(`${url}/main/auth/api/v1/user/password/reset `, body, config)
         .then((res) => {
             //  const temp = []
              setNotification({
                  text: `Password Changed`,
                  type: 'success',
              })
              setChange(false)
              setLoadingChange(false)
         })
         .catch((err) => {
            console.log(err)
            setChangeError(err.response.data.message)
            setLoadingChange(false)
         })
 }
    useEffect(() => {
        if (isPageLoading) {
            setIsPageLoading(false)
            return
        }
        sendNow()
        applyTheme(theme)
    }, [theme])

    useEffect(() => {
        switch (index) {
            case 0:
                setTheme(PlatformEnginePkgAuthApiTheme.ThemeLight)
                break
            case 1:
                setTheme(PlatformEnginePkgAuthApiTheme.ThemeDark)
                break
            case 2:
                setTheme(PlatformEnginePkgAuthApiTheme.ThemeSystem)
                break
            default:
                setTheme(PlatformEnginePkgAuthApiTheme.ThemeLight)
                break
        }
    }, [index])
 const show_compliance = window.__RUNTIME_CONFIG__.REACT_APP_SHOW_COMPLIANCE
    return (
        <>
            <Modal
                header="Password Change"
                visible={change}
                onDismiss={() => {
                    setChange(false)
                }}
                footer={
                    <Flex className="w-full gap-2" justifyContent="end">
                        <Button
                            onClick={() => {
                                setChange(false)
                            }}
                        >
                            Close
                        </Button>
                        <Button
                            loading={loadingChange}
                            onClick={ChangePassword}
                            variant="primary"
                        >
                            Change Password
                        </Button>
                    </Flex>
                }
            >
                {/* <Alert type="info">
                    It's First time you logged in . Please Change your Password
                </Alert> */}
                <Flex
                    flexDirection="col"
                    className="gap-2 mt-2 mb-2 w-full"
                    justifyContent="start"
                    alignItems="start"
                >
                    <FormField
                        // description="This is a description."
                        errorText={errors?.current}
                        className=" w-full"
                        label="Current Password"
                    >
                        <Input
                            value={password?.current}
                            type="password"
                            onChange={(event) => {
                                setPassword({
                                    ...password,
                                    current: event.detail.value,
                                })
                                setErrors({
                                    ...errors,
                                    current: '',
                                })
                            }}
                        />
                    </FormField>
                    <FormField
                        // description="This is a description."
                        errorText={errors?.new}
                        className=" w-full"
                        label="New Password"
                    >
                        <Input
                            value={password?.new}
                            type="password"
                            onChange={(event) => {
                                setPassword({
                                    ...password,
                                    new: event.detail.value,
                                })
                                setErrors({
                                    ...errors,
                                    new: '',
                                })
                            }}
                        />
                    </FormField>
                    <FormField
                        // description="This is a description."
                        errorText={errors?.confirm}
                        label="Confirm Password"
                        className=" w-full"
                    >
                        <Input
                            value={password?.confirm}
                            type="password"
                            onChange={(event) => {
                                setPassword({
                                    ...password,
                                    confirm: event.detail.value,
                                })
                                setErrors({
                                    ...errors,
                                    confirm: '',
                                })
                            }}
                        />
                    </FormField>
                </Flex>
                {changeError && changeError != '' && (
                    <Alert type="error">{changeError}</Alert>
                )}
            </Modal>
            {show_compliance == 'false' && (
                <Popover
                    dismissButton={false}
                    position="bottom"
                    size="large"
                    triggerType="custom"
                    content={
                        <Flex
                            flexDirection="col"
                            alignItems="start"
                            className="pb-0 mb-0 "
                            // border-b border-b-gray-700
                        >
                            {/* <Text className="mb-1">ACCOUNT</Text> */}
                            <Flex
                                onClick={() => {
                                    // navigate(`/profile`)
                                    navigate(`/`)
                                }}
                                className="py-2 px-5 rounded-md cursor-pointer justify-start flex-row items-center gap-2 text-gray-900 hover:text-gray-800 hover:bg-gray-300"
                            >
                                <MagnifyingGlassIcon className="h-6 w-6 stroke-2" />
                                <Text className="text-inherit">CloudQL</Text>
                            </Flex>
                            <Flex
                                onClick={() => {
                                    navigate(`/integration/plugins`)
                                }}
                                className="py-2 px-5 rounded-md cursor-pointer justify-start flex-row items-center gap-2 text-gray-900 hover:text-gray-800 hover:bg-gray-300"
                            >
                                <PuzzlePieceIcon className="h-6 w-6 stroke-2" />
                                <Text className="text-inherit">
                                    Integration
                                </Text>
                            </Flex>

                            <Flex
                                onClick={() => navigate(`/administration`)}
                                className="py-2 px-5 text-gray-900 hover:text-gray-800 justify-start flex-row items-center gap-2 rounded-md cursor-pointer  hover:bg-gray-300"
                            >
                                <Cog6ToothIcon className="h-6 w-6 stroke-2" />
                                <Text className="text-inherit">
                                    Administration
                                </Text>
                            </Flex>
                        </Flex>
                    }
                >
                    <Button iconName="menu"></Button>
                </Popover>
            )}

            <Popover
                dismissButton={false}
                position="bottom"
                size="large"
                triggerType="custom"
                // className="bg-openg-950 px-4 py-2 w-64 !ring-gray-600"
                content={
                    <Flex
                        flexDirection="col"
                        alignItems="start"
                        className="pb-0 mb-0 "
                        // border-b border-b-gray-700
                    >
                        {/* <Text className="mb-1">ACCOUNT</Text> */}
                        <Flex
                            onClick={() => {
                                // navigate(`/profile`)
                                navigate(`/profile`)
                            }}
                            className="py-2 px-5 rounded-md cursor-pointer text-gray-900 hover:text-gray-800 hover:bg-gray-300"
                        >
                            <Text className="text-inherit">Profile info</Text>
                        </Flex>
                        <Flex
                            onClick={() => {
                                setChange(true)
                            }}
                            className="py-2 px-5 rounded-md cursor-pointer text-gray-900 hover:text-gray-800 hover:bg-gray-300"
                        >
                            <Text className="text-inherit">
                                Change Password
                            </Text>
                        </Flex>
                        {/* <Flex
                                onClick={() => navigate(`/ws/billing`)}
                                className="py-2 px-5 rounded-md cursor-pointer text-gray-300 hover:text-gray-50 hover:bg-gray-300"
                            >
                                <Text className="text-inherit">Billing</Text>
                            </Flex> */}
                        <Flex
                            onClick={() => logout()}
                            className="py-2 px-5 text-gray-900 rounded-md cursor-pointer hover:text-gray-800 hover:bg-gray-300"
                        >
                            <Text className="text-inherit">Logout</Text>
                            <ArrowTopRightOnSquareIcon className="w-5 text-gray-900" />
                        </Flex>
                    </Flex>
                }
            >
                <Button iconName="user-profile"></Button>
            </Popover>
        </>
    )
}
