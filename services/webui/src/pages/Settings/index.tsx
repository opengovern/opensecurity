
import SettingsEntitlement from './Entitlement'
import SettingsMembers from './Members'
import SettingsWorkspaceAPIKeys from './APIKeys'

import { Tabs } from '@cloudscape-design/components'
import SettingsConnectors from './Connectors'


export default function Settings() {

    return (
        <>
            
            <div
                className="w-full"
                style={
                    window.innerWidth < 768
                        ? { width: `${window.innerWidth - 80}px` }
                        : {}
                }
            >
                <Tabs
                    tabs={[
                        {
                            label: 'Settings',
                            content: (
                                <>
                                    <SettingsEntitlement />
                                </>
                            ),
                            id: '0',
                        },

                        {
                            label: 'Authentication',
                            content: (
                                <>
                                    <SettingsMembers />
                                </>
                            ),
                            id: '1',
                        },
                        {
                            label: 'SSO Configuration',
                            content: (
                                <>
                                    <SettingsConnectors />
                                </>
                            ),
                            id: '2',
                        },
                        {
                            label: 'API',
                            content: (
                                <>
                                    <SettingsWorkspaceAPIKeys />
                                </>
                            ),
                            id: '3',
                        },
                    ]}
                />
            </div>
        </>
    )
}
