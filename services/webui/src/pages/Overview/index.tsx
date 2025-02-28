import { useParams } from 'react-router-dom'
import { Col, Flex, Grid } from '@tremor/react'
import Governance from './Governance'
import TopHeader from '../../components/Layout/Header'
import { defaultHomepageTime } from '../../utilities/urlstate'
import Query from './Query'
import QuickNav from './QuickNav'
import Shortcuts from './Shortcuts'
import axios from 'axios'
import { useEffect, useState } from 'react'
import { Alert, Button, Modal } from '@cloudscape-design/components'
import FormField from '@cloudscape-design/components/form-field'
import Input from '@cloudscape-design/components/input'
import { error } from 'console'
import { useAtom, useSetAtom } from 'jotai'
import { ForbiddenAtom, meAtom, notificationAtom } from '../../store'
import { useAuth } from '../../utilities/auth'
import { useAuthApiV1UserInviteCreate } from '../../api/auth.gen'
import Integrations from './Integrations'
import { useComplianceApiV1QueriesSyncList } from '../../api/compliance.gen'
import SRE from './KPI_Cards'
import { useWorkspaceApiV3LoadSampleData } from '../../api/metadata.gen'

export default function Overview() {
    
    const element = document.getElementById('myDIV')?.offsetHeight
    const [change, setChange] = useState<boolean>(false)
    const [userModal, setUserModal] = useState<boolean>(false)
    const [userData, setUserData] = useState<any>({
        email: '',
        password: '',
        confirm: '',
    })
    const [userErrors, setUserErrors] = useState({
        email: '',
        password: '',
        success: '',
        failed: '',
    })
    const { user, logout } = useAuth()
    const [me, setMe] = useAtom(meAtom)

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
    const setForbbiden = useSetAtom(ForbiddenAtom)
    const [changeError, setChangeError] = useState()
    const {
        isExecuted,
        isLoading,
        error,
        sendNow: createInvite,
    } = useAuthApiV1UserInviteCreate(
        {
            email_address: userData?.email || '',
            role: 'admin',
            password: userData?.password,
            is_active: true
        },
        {},
        false
    )
      const {
            isExecuted:isExecuted1,
            isLoading: isLoadingLoad,
            error:error1,
            sendNow: loadData,
        } = useWorkspaceApiV3LoadSampleData({}, false)
    const setNotification = useSetAtom(notificationAtom)
    const [loadingChange, setLoadingChange] = useState(false)
    const PassCheck = () => {
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

        axios
            .get(`${url}/main/auth/api/v1/user/password/check`, config)
            .then((res) => {
                //  const temp = []
                if (res.data == 'CHANGE_REQUIRED') {
                    setChange(true)
                    if (me?.email == 'admin@opensecurity.sh') {
                        runSync()
                        loadData()
                    }
                }
            })
            .catch((err) => {
                if( err?.response?.status === 401){
                        setForbbiden(true)
                }

                console.log(err)
            })
    }
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
                setChange(false)
                setLoadingChange(false)
                setNotification({
                    text: `Password Changed`,
                    type: 'success',
                })
                logout()
            })
            .catch((err) => {
                console.log(err)
                setChangeError(err.response.data.message)
                setLoadingChange(false)
            })
    }
const {
    isLoading: syncLoading,
    isExecuted: syncExecuted,
    error: syncError,
    sendNow: runSync,
} = useComplianceApiV1QueriesSyncList({}, {}, false)

    useEffect(() => {
       
         if (me?.connector_id === 'local') {
             PassCheck()
            
         }
         
    }, [me])
    const CheckEmail = () => {
        if (!userData?.email || userData?.email == '') {
            setUserErrors({
                ...userErrors,
                email: 'Please enter email',
            })
            return
        }
        if (!userData?.password || userData?.password == '') {
            setUserErrors({
                ...userErrors,
                password: 'Please enter  password',
            })
            return
        }
        if (!userData?.confirm || userData?.confirm == '') {
            setUserErrors({
                ...userErrors,
                password: 'Please enter  password',
            })
            return
        }
        if (!userData?.email.includes('@')) {
            setUserErrors({
                ...userErrors,
                email: 'Please enter a valid email',
            })
            return
        }
        if (userData?.password !== userData?.confirm) {
            setUserErrors({
                ...userErrors,
                password: 'Passwords are not same',
            })
            return
        }

        createInvite()
    }
    useEffect(() => {
        if (!isLoading && isExecuted && (!error || error != '')) {
            setUserErrors({ ...userErrors, success: 'Loggin out' })
            setInterval(logout, 3000)
        }
    }, [isLoading, isExecuted])
    return (
        <>
            <Modal
                header="First Time Login"
                visible={change}
                onDismiss={() => {}}
                footer={
                    <Flex className="w-full" justifyContent="end">
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
                <Alert type="info">
                    It's First time you logged in . Please Change your Password
                </Alert>
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
            <Modal
                header="Update Login Credentials"
                visible={userModal}
                onDismiss={() => {}}
                footer={
                    <Flex className="w-full" justifyContent="end">
                        <Button
                            loading={isLoading && isExecuted}
                            disabled={isLoading && isExecuted}
                            onClick={CheckEmail}
                            variant="primary"
                        >
                            Sumbit
                        </Button>
                    </Flex>
                }
            >
                <Alert type="info">
                    You logged in with default credentials. Please create new
                    ones.
                </Alert>
                <Flex
                    flexDirection="col"
                    className="gap-2 mt-2 mb-2 w-full"
                    justifyContent="start"
                    alignItems="start"
                >
                    <FormField
                        // description="This is a description."
                        errorText={userErrors?.email}
                        className=" w-full"
                        label="Email"
                    >
                        <Input
                            value={userData?.email}
                            type="email"
                            onChange={(event) => {
                                setUserData({
                                    ...userData,
                                    email: event.detail.value,
                                })
                                setUserErrors({
                                    ...userErrors,
                                    email: '',
                                })
                            }}
                        />
                    </FormField>
                    <FormField
                        // description="This is a description."
                        errorText={userErrors?.password}
                        className=" w-full"
                        label="Password"
                    >
                        <Input
                            value={userData?.password}
                            type="password"
                            onChange={(event) => {
                                setUserData({
                                    ...userData,
                                    password: event.detail.value,
                                })
                                setUserErrors({
                                    ...userErrors,
                                    password: '',
                                })
                            }}
                        />
                    </FormField>
                    <FormField
                        // description="This is a description."
                        errorText={userErrors?.password}
                        className=" w-full"
                        label="Confirm Password"
                    >
                        <Input
                            value={userData?.confirm}
                            type="password"
                            onChange={(event) => {
                                setUserData({
                                    ...userData,
                                    confirm: event.detail.value,
                                })
                                setUserErrors({
                                    ...userErrors,
                                    password: '',
                                })
                            }}
                        />
                    </FormField>
                </Flex>
                {error && error != '' && <Alert type="error">{error}</Alert>}
                {userErrors?.success && userErrors?.success != '' && (
                    <Alert header="User Created" type="success">
                        Logging out...
                    </Alert>
                )}
            </Modal>

            <Grid
                numItems={11}
                className="w-full gap-8  h-fit "
                style={window.innerWidth > 768 ? { gridAutoRows: '1fr' } : {}}
            >
                <Col numColSpan={11} numColSpanSm={8}>
                    <Flex
                        flexDirection="col"
                        alignItems="start"
                        className="gap-4 h-full"
                        id="myDIV"
                    >
                        <Grid numItems={6} className="w-full gap-6 h-fit ">
                            <Col numColSpan={6} className="h-fit">
                                <Shortcuts />
                            </Col>
                            <Col numColSpan={6} className="h-fit">
                                <SRE />
                            </Col>
                            <Col numColSpan={6} className="h-full">
                                <Governance />
                            </Col>
                        </Grid>
                    </Flex>
                </Col>
                <Col
                    numColSpan={11}
                    numColSpanSm={3}
                    className="sm:h-full h-fit"
                >
                    <Integrations />
                </Col>
            </Grid>
        </>
    )
}
