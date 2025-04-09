import { Flex } from '@tremor/react'
import { ReactNode, UIEvent, useEffect, useState } from 'react'
import Footer from './Footer'
import Sidebar from './Sidebar'
import Notification from '../Notification'
import { useNavigate } from 'react-router-dom'
import { useAtom, useAtomValue, useSetAtom } from 'jotai'
import { sampleAtom } from '../../store'
import TopHeader from './Header'
import {
    AppLayoutToolbar,
    BreadcrumbGroup,
    Container,
    Flashbar,
    Header,
    HelpPanel,
    SideNavigation,
    SplitPanel,
} from '@cloudscape-design/components'
import NewSidebar from './NewSidebar'
type IProps = {
    children: ReactNode
    onScroll?: (e: UIEvent) => void
    scrollRef?: any
}
const show_compliance = window.__RUNTIME_CONFIG__.REACT_APP_SHOW_COMPLIANCE

const Mapping = {
    cloudql: 'CloudQL',
    integration: 'Integration',
    compliance: 'Compliance',
    overview: 'Overview',
    settings: 'Settings',
    tasks: 'Tasks',
    ai: 'AI',
}
const INTEGRATION_MAPPING = {
    azure_subscription: 'Microsoft Azure Subscription',
    jira_cloud: 'Atlassian JIRA Cloud',
    aws_cloud_account: 'Amazon Web Services (AWS)',
    entraid_directory: 'Microsoft EntraID Directory',
    github_account: 'GitHub',
    digitalocean_team: 'DigitalOcean',
    cloudflare_account: 'Cloudflare',
    linode_account: 'Linode (Akamai)',
    render_account: 'Render',
    fly_account: 'Fly.io',
    semgrep_account: 'Semgrep',
    kubernetes: 'Kubernetes',
    openai_integration: 'OpenAI',
    cohereai_project: 'CohereAI',
    google_workspace_account: 'Google Workspace',
    doppler_account: 'Doppler',
    tailscale_account: 'Tailscale',
    heroku_account: 'Heroku',
    oci_repository: 'OCI Repository',
    slack_account: 'Slack',
    chainguard_dev_account: 'Chainguard.dev',
    godaddy_account: 'GoDaddy',
    servicenow_account: 'ServiceNow',
    okta_account: 'Okta',
    aws_costs: 'Amazon Web Services (AWS) Costs',
    azure_costs: 'Microsoft Azure Costs',
    huggingface_account: 'HuggingFace',
    jamf_account: 'Jamf',
    jumpcloud_account: 'JumpCloud',
    gitlab_account: 'GitLab',
}
export default function Layout({ children, onScroll, scrollRef }: IProps) {
    const url = window.location.pathname.split('/')

    const showSidebarCallback = url[1] == 'callback' ? false : true
    const [showSidebar, setShowSidebar] = useState(true)
    const [breadCrumbItems, setBreadCrumbItems] = useState<any>([])
    const GetBreadCrumItems = () => {
        const temp = [
            {
                text: 'Home',
                href: '/',
            },
        ]
        const path = window.location.pathname
        console.log(path)
        if (path.includes('integration')) {
            console.log(url)
            if (url.length > 3) {
                const integration = url[3]
                // @ts-ignore
                const integrationName = INTEGRATION_MAPPING[integration]
                if (integrationName) {
                    temp.push({
                        text: 'Plugins',
                        href: '/integration/plugins',
                    })
                    temp.push({
                        text: integrationName,
                        href: path,
                    })
                }
            } else {
                temp.push({
                    text: 'Plugins',
                    href: '/integration/plugins',
                })
            }
        }

        return setBreadCrumbItems(temp)
    }
    const GetCurrentPage = () => {
        const path = window.location.pathname
        if (path.includes('cloudql')) {
            return '/cloudql'
        } else if (path.includes('integration')) {
            return '/integration/plugins'
        } else if (path.includes('compliance')) {
            return '/compliance'
        } else if (path.includes('jobs')) {
            return '/jobs'
        } else if (path.includes('administration')) {
            return '/administration'
        } else if (path.includes('ai')) {
            return '/ai'
        } else if (
            path.includes('automation') ||
            path.includes('dashboards') ||
            path.includes('request-access') ||
            path.includes('stacks') ||
            path.includes('/workload-optimizer')
        ) {
            return '/automation'
        }

        return ''
    }
    useEffect(() => {
        GetBreadCrumItems()
    }, [window.location.pathname])

    return (
        <>
            <AppLayoutToolbar
                breadcrumbs={<BreadcrumbGroup items={breadCrumbItems} />}
                navigationOpen={showSidebar}
                onNavigationChange={({ detail }) => setShowSidebar(detail.open)}
                toolsHide={true}
                navigation={
                    <>
                        {showSidebarCallback ? (
                            <NewSidebar currentPage={GetCurrentPage()} />
                        ) : (
                            ''
                        )}
                    </>
                }
                notifications={<Notification />}
                content={children}
            />
            <Footer />
        </>
    )
}
