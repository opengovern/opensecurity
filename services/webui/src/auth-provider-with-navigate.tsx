import React, { useEffect, useState } from 'react'
import { Button, Card, Flex, Text, TextInput, Title } from '@tremor/react'
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from './utilities/auth'
import {  OpenGovernanceBig } from './icons/icons'
import { authHostname } from './api/ApiConfig'

interface IAuthProviderWithNavigate {
    children?: React.ReactNode
}

export const AuthProviderWithNavigate = ({
    children,
}: IAuthProviderWithNavigate) => {
    const [locationSearchParams, setSearchParams] = useSearchParams()
    const { isAuthenticated, isLoading, loginWithCode, error } = useAuth()
    const page = window.location.pathname.split('/')

    useEffect(() => {
        if (
            locationSearchParams.get('code') !== null &&
            locationSearchParams.get('code') !== '' &&
            !isLoading &&
            error.message === '' &&
            !isAuthenticated
        ) {
            loginWithCode(locationSearchParams.get('code') || '')
        }

        if (!isAuthenticated && page[1] !== 'callback') {
            const callback = `${window.location.origin}/callback`


            const searchParams = new URLSearchParams()
            searchParams.append('client_id', 'public-client')
            searchParams.append('redirect_uri', callback)
            searchParams.append('scope', 'openid email')
            searchParams.append('response_type', 'code')
            searchParams.toString()

            sessionStorage.setItem('callbackURL', window.location.href)

            window.location.href = `${authHostname()}/dex/auth?${searchParams.toString()}`
        }
    }, [])

    return (
        <>
            <span />
            {children}
        </>
    )
}
