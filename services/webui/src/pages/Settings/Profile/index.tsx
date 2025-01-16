import {
    Button,
    Card,
    Divider,
    Flex,
    Select,
    SelectItem,
    Text,
    Title,
} from '@tremor/react'
import { useEffect, useState } from 'react'
import { useAtom, useAtomValue } from 'jotai'
import dayjs from 'dayjs'
import { useAuthApiV1UserPreferencesUpdate } from '../../../api/auth.gen'
import { PlatformEnginePkgAuthApiTheme } from '../../../api/api'
import { colorBlindModeAtom, isDemoAtom, meAtom } from '../../../store'
import { applyTheme, currentTheme, parseTheme } from '../../../utilities/theme'
import { useAuth } from '../../../utilities/auth'
import { KeyValuePairs } from '@cloudscape-design/components'

export default function SettingsProfile() {
    const { user } = useAuth()
    const [colorBlindMode, setColorBlindMode] = useAtom(colorBlindModeAtom)
    const [isDemo, setIsDemo] = useAtom(isDemoAtom)

    const [me, setMe] = useAtom(meAtom)

    const memberSince = me?.memberSince
    const lastLogin = me?.lastLogin

    const [enableColorBlindMode, setEnableColorBlindMode] =
        useState<boolean>(colorBlindMode)
    const [theme, setTheme] =
        useState<PlatformEnginePkgAuthApiTheme>(currentTheme())

    const { response, isLoading, isExecuted, error, sendNow } =
        useAuthApiV1UserPreferencesUpdate(
            {
                enableColorBlindMode,
                theme,
            },
            {},
            false
        )
    useEffect(() => {
        if (!isLoading && isExecuted) {
            applyTheme(theme)
            setColorBlindMode(enableColorBlindMode)
        }
    }, [isLoading])

    const dateFormat = (date: string, showOnline: boolean) => {
        if (date.length === 0) {
            return ''
        }

        const dt = dayjs.utc(date)
        const now = dayjs.utc()
        const d = now.diff(dt, 'day')
        const h = now.diff(dt, 'hour')
        const m = now.diff(dt, 'minute')
        if (d > 7) {
            return `${dt.format('YYYY-MM-DD')}`
        }
        if (h > 48) {
            return `${d} days ago`
        }
        if (m > 60) {
            return `${h} hours ago`
        }
        if (m > 15 || !showOnline) {
            return `${m} minutes ago`
        }
        return `Online`
    }

    return (
        <Card>
            {user?.picture && (
                <img
                    className="my-3 rounded-lg"
                    src={user?.picture}
                    alt={user.name}
                />
            )}
            <Title className="font-semibold mb-2">Profile Information</Title>
            {/* <Flex flexDirection="col"> */}
                {/* <Divider className="my-1 py-1" />
                <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">First Name</Text>
                    <Text className="w-1/2 text-gray-800">
                        {user?.given_name}
                    </Text>
                </Flex>
                <Divider className="my-1 py-1" />
                <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">Last Name</Text>
                    <Text className="w-1/2 text-gray-800">
                        {user?.family_name}
                    </Text>
                </Flex>
                <Divider className="my-1 py-1" />
                <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">Email</Text>
                    <Text className="w-1/2 text-gray-800">{user?.email}</Text>
                </Flex> */}
                {/* <Divider className="my-1 py-1" /> */}
                <KeyValuePairs
                    columns={3}
                    items={[
                        {
                            label: 'Email',
                            value: user?.email,
                        },
                        {
                            label: 'Member Since',
                            value: dateFormat(memberSince || '', false),
                        },
                        {
                            label: 'Last Online',
                            value: dateFormat(lastLogin || '', true),
                        },
                    ]}
                />
                {/* <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">Member Since</Text>
                    <Text className="w-1/2 text-gray-800">
                        {dateFormat(memberSince || '', false)}
                    </Text>
                </Flex>
                <Divider className="my-1 py-1" />
                <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">Last Online</Text>
                    <Text className="w-1/2 text-gray-800">
                        {dateFormat(lastLogin || '', true)}
                    </Text>
                </Flex> */}
            {/* </Flex> */}
            <Title className="font-semibold mt-10">Personalization</Title>
            <Flex flexDirection="col">
                <Divider className="my-1 py-1" />
                {/* <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">Color Theme</Text>
                    <Select
                        disabled={isExecuted && isLoading}
                        value={theme}
                        onValueChange={(v) => {
                            setTheme(parseTheme(v))
                        }}
                        className="w-1/2"
                    >
                        <SelectItem value="light">Light</SelectItem>
                        <SelectItem value="dark">Dark</SelectItem>
                        <SelectItem value="system">System</SelectItem>
                    </Select>
                </Flex>
                <Divider className="my-1 py-1" /> */}
                {/* <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">Accessibility Mode (WAI-ARIA)</Text>
                    <Select
                        disabled={isExecuted && isLoading}
                        value={String(enableColorBlindMode)}
                        onValueChange={(v) => {
                            setEnableColorBlindMode(v === 'true')
                        }}
                        className="w-1/2"
                    >
                        <SelectItem value="true">Enabled</SelectItem>
                        <SelectItem value="false">Disabled</SelectItem>
                    </Select>
                </Flex> */}
                {/* <Divider className="my-1 py-1" />
                <Flex flexDirection="row" justifyContent="between">
                    <Text className="w-1/2">Demo mode</Text>
                    <Select
                        value={String(isDemo)}
                        onValueChange={(v) => {
                            setIsDemo(v === 'true')
                            localStorage.setItem('demoMode', String(v))
                        }}
                        className="w-1/2"
                    >
                        <SelectItem value="true">True</SelectItem>
                        <SelectItem value="false">False</SelectItem>
                    </Select>
                </Flex> */}
                {/* <Flex flexDirection="row" justifyContent="end" className="mt-2">
                    <Button
                        loading={isExecuted && isLoading}
                        variant="secondary"
                        onClick={() => sendNow()}
                    >
                        Save
                    </Button>
                </Flex> */}
            </Flex>
        </Card>
    )
}
