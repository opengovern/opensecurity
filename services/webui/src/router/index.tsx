import { Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import { useEffect } from 'react'
import NotFound from '../pages/Errors'
import { CallbackPage } from '../pages/Callback'
import Settings from '../pages/Settings'
import Logout from '../pages/Logout'
import Integrations from '../pages/Integrations'
import Compliance from '../pages/Compliance'
import Overview from '../pages/Overview'
// import Boostrap from '../pages/Workspaces/Bootstrap'

import Findings from '../pages/Findings'
import Layout from '../components/Layout'
import RequestDemo from '../pages/RequestDemo'
import RequestAccess from '../pages/Integrations/RequestAccess'
import SettingsJobs from '../pages/Settings/Jobs'
import AllControls from '../pages/Compliance/All Controls'
import SettingsWorkspaceAPIKeys from '../pages/Settings/APIKeys'
import SettingsParameters from '../pages/Settings/Parameters'
import SettingsMembers from '../pages/Settings/Members'
import NewBenchmarkSummary from '../pages/Compliance/NewBenchmarkSummary'

import SettingsProfile from '../pages/Settings/Profile'
import TypeDetail from '../pages/Integrations/TypeDetailNew'
import EvaluateDetail from '../pages/Compliance/NewBenchmarkSummary/EvaluateTable/EvaluateDetail'
import Tasks from '../pages/Tasks'
import TaskDetail from '../pages/Tasks/TaskDetail'
import AI from '../pages/AI'
import AllPolicy from '../pages/Compliance/All Policy'
import Framework from '../pages/Compliance/FrameWorks'
import ComplianceJobs from '../pages/Settings/Jobs/Compliance'
import Query from '../pages/Search/Query'

const show_compliance = window.__RUNTIME_CONFIG__.REACT_APP_SHOW_COMPLIANCE
console.log(show_compliance)

const authRoutes = [
    {
        key: 'callback',
        path: '/callback',
        element: <CallbackPage />,
        noAuth: true,
    },
    {
        key: 'logout',
        path: '/logout',
        element: <Logout />,
        noAuth: true,
    },
    {
        key: '*',
        path: '*',
        element: <NotFound />,
        noAuth: true,
    },

    {
        key: 'workload optimizer',
        path: '/workload-optimizer',
        element: <RequestAccess />,
    },
    {
        key: 'stacks',
        path: '/stacks',
        element: <RequestAccess />,
    },
    {
        key: 'Automation',
        path: '/automation',
        element: <RequestAccess />,
    },
    {
        key: 'Automation',
        path: '/dashboards',
        element: <RequestAccess />,
    },

    {
        key: 'integrations',
        path: 'integration/plugins',
        element: <Integrations />,
    },
    {
        key: 'request-access',
        path: '/request-access',
        element: <RequestAccess />,
    },

    {
        key: 'connector detail',
        path: 'integration/plugins/:type',
        element: <TypeDetail />,
    },

    {
        key: 'settings page',
        path: '/administration/settings',
        element: <Settings />,
    },
    {
        key: 'settings page',
        path: '/administration/access',
        element: <SettingsMembers />,
    },
    {
        key: 'Profile',
        path: '/profile',
        element: <SettingsProfile />,
    },
    {
        key: 'settings Jobs',
        path: 'integration/jobs',
        element: <SettingsJobs />,
    },
    {
        key: 'settings APi Keys',
        path: '/settings/api-keys',
        element: <SettingsWorkspaceAPIKeys />,
    },

    {
        key: 'settings Authentications',
        path: '/settings/authentication',
        element: <SettingsMembers />,
    },
    {
        key: 'Compliance',
        path: '/compliance',
        element: <Compliance />,
    },
    {
        key: 'Compliance',
        path: '/compliance/jobs',
        element: <ComplianceJobs />,
    },
    {
        key: 'Compliance',
        path: '/compliance/frameworks',
        element: <Framework />,
    },
    {
        key: 'Compliance',
        path: '/compliance/controls',
        element: <Compliance />,
    },
    {
        key: 'Compliance',
        path: '/compliance/policies',
        element: <AllPolicy />,
    },
    {
        key: 'Compliance',
        path: '/compliance/parameters',
        element: <Compliance />,
    },

    {
        key: 'benchmark summary 2',
        path: '/compliance/frameworks/:benchmarkId',
        element: <NewBenchmarkSummary />,
    },

    {
        key: 'allControls',
        path: '/compliance/library/parameters',
        element: <SettingsParameters />,
    },

    {
        key: 'benchmark single connection',
        path: '/compliance/frameworks/:benchmarkId/report/:id',
        element: <EvaluateDetail />,
    },
    {
        key: 'Incidents control',
        path: '/incidents',
        element: <Findings />,
    },
    {
        key: 'Incidents control',
        path: '/incidents/controls',
        element: <Findings />,
    },

    {
        key: 'Incidents control',
        path: '/incidents/resources',
        element: <Findings />,
    },

    {
        key: ' summary',
        path: '/incidents/summary',
        element: <Findings />,
    },

    {
        key: 'Account Posture',
        path: '/incidents/account-posture',
        element: <Findings />,
    },

    {
        key: 'home',
        path: '/',
        element: show_compliance !== 'false' ? <Overview /> : <Query />,
    },

    {
        key: 'request a demo',
        path: '/ws/requestdemo',
        element: <RequestDemo />,
    },

    {
        key: 'Search',
        path: '/cloudql',
        element: <Query />,
    },
    {
        key: 'Tasks',
        path: '/tasks',
        element: <Tasks />,
    },
    {
        key: 'Tasks',
        path: '/tasks/:id',
        element: <TaskDetail />,
    },
    {
        key: 'Ai',
        path: '/ai',
        element: <AI />,
    },
]

export default function Router() {
    const navigate = useNavigate()

    const url = window.location.pathname.split('/')
  

    useEffect(() => {
        if (url[1] === 'undefined') {
            navigate('/')
        }
    }, [url])

    return (
        <Layout>
            <Routes>
                {authRoutes.map((route) => (
                    <Route
                        key={route.key}
                        path={route.path}
                        element={route.element}
                    />
                ))}
            </Routes>
        </Layout>
    )
}
