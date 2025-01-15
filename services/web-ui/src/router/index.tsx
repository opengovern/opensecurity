import { Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import { useEffect } from 'react'
import NotFound from '../pages/Errors'
import { CallbackPage } from '../pages/Callback'
import Settings from '../pages/Settings'
import Logout from '../pages/Logout'
import Integrations from '../pages/Integrations'
import Compliance from '../pages/Governance/Compliance'
import Overview from '../pages/Overview'
import Stack from '../pages/Stack'
// import Boostrap from '../pages/Workspaces/Bootstrap'
import ControlDetail from '../pages/Governance/Controls/ControlSummary'
import Findings from '../pages/Governance/Findings'
import Layout from '../components/Layout'
import RequestDemo from '../pages/RequestDemo'
import RequestAccess from '../pages/Integrations/RequestAccess'
import SettingsJobs from '../pages/Settings/Jobs'
import AllControls from '../pages/Governance/Compliance/All Controls'
import SettingsWorkspaceAPIKeys from '../pages/Settings/APIKeys'
import SettingsParameters from '../pages/Settings/Parameters'
import SettingsMembers from '../pages/Settings/Members'
import NewBenchmarkSummary from '../pages/Governance/Compliance/NewBenchmarkSummary'
import Dashboard from '../pages/Dashboard'
import Library from '../pages/Governance/Compliance/Library'
import Search from '../pages/Search'
import SettingsAccess from '../pages/Settings/Access'
import SettingsProfile from '../pages/Settings/Profile'
import SearchLanding from '../pages/Search/landing'
import TypeDetail from '../pages/Integrations/TypeDetailNew'
import EvaluateDetail from '../pages/Governance/Compliance/NewBenchmarkSummary/EvaluateTable/EvaluateDetail/inde'
import Tasks from '../pages/Tasks'
import TaskDetail from '../pages/Tasks/TaskDetail'

const authRoutes = [
    // {
    //     key: 'url',
    //     path: '/',
    //     element: <Navigate to="/ws/workspaces?onLogin" replace />,
    //     noAuth: true,
    // },
    // {
    //     key: 'ws name',
    //     path: '/',
    //     element: <Navigate to="overview" />,
    //     noAuth: true,
    // },
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
    // {
    //     key: 'workspaces',
    //     path: '/ws/workspaces',
    //     element: <Workspaces />,
    // },
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
        key: 'dashboards',
        path: '/dashboard',
        element: <Dashboard />,
    },
    
   
    
    {
        key: 'integrations',
        path: '/integrations',
        element: <Integrations />,
    },
    {
        key: 'request-access',
        path: '/request-access',
        element: <RequestAccess />,
    },

    {
        key: 'connector detail',
        path: '/integrations/:type',
        element: <TypeDetail />,
    },

    {
        key: 'settings page',
        path: '/administration',
        element: <Settings />,
    },
    {
        key: 'Profile',
        path: '/profile',
        element: <SettingsProfile />,
    },
    {
        key: 'settings Jobs',
        path: '/jobs',
        element: <SettingsJobs />,
    },
    {
        key: 'settings APi Keys',
        path: '/settings/api-keys',
        element: <SettingsWorkspaceAPIKeys />,
    },
    // {
    //     key: 'settings variables',
    //     path: '/settings/variables',
    //     element: <SettingsParameters />,
    // },
    {
        key: 'settings Authentications',
        path: '/settings/authentication',
        element: <SettingsMembers />,
    },
    {
        key: 'settings Access',
        path: '/settings/access',
        element: <SettingsAccess />,
    },
   
    {
        key: 'Compliance',
        path: '/compliance',
        element: <Compliance />,
    },

    {
        key: 'benchmark summary 2',
        path: '/compliance/:benchmarkId',
        element: <NewBenchmarkSummary />,
    },
    {
        key: 'allControls',
        path: '/compliance/library',
        element: <Library />,
    },
    {
        key: 'allControls',
        path: '/compliance/library/parameters',
        element: <SettingsParameters />,
    },
    // {
    //     key: 'allBenchmarks',
    //     path: '/compliance/benchmarks',
    //     element: <AllBenchmarks />,
    // },
    {
        key: 'benchmark summary',
        path: '/compliance/:benchmarkId/:controlId',
        element: <ControlDetail />,
    },
   
    {
        key: 'benchmark single connection',
        path: '/compliance/:benchmarkId/report/:id',
        element: <EvaluateDetail />,
    },
    {
        key: 'Incidents control',
        path: '/incidents',
        element: <Findings />,
    },
    // {
    //     key: 'Resource summary',
    //     path: '/incidents/resource-summary',
    //     element: <Findings />,
    // },
    {
        key: ' summary',
        path: '/incidents/summary',
        element: <Findings />,
    },

    // {
    //     key: 'Drift Events',
    //     path: '/incidents/drift-events',
    //     element: <Findings />,
    // },
    {
        key: 'Account Posture',
        path: '/incidents/account-posture',
        element: <Findings />,
    },
    // {
    //     key: 'Control Summary',
    //     path: '/incidents/control-summary',
    //     element: <Findings />,
    // },
    {
        key: 'incidents',
        path: '/incidents/:controlId',
        element: <ControlDetail />,
    },
    
    {
        key: 'home',
        path: '/',
        element: <Overview />,
    },
    {
        key: 'deployment',
        path: '/deployment',
        element: <Stack />,
    },
    // {
    //     key: 'query',
    //     path: '/query',
    //     element: <Query />,
    // },
    // {
    //     key: 'bootstrap',
    //     path: '/bootstrap',
    //     element: <Boostrap />,
    // },
    // {
    //     key: 'new-ws',
    //     path: '/ws/new-ws',
    //     element: <Boostrap />,
    // },

    
    // {
    //     key: 'resource collection assets metrics',
    //     path: '/:ws/resource-collection/:resourceId/assets-details',
    //     component: AssetDetails,
    // },
  
    {
        key: 'request a demo',
        path: '/ws/requestdemo',
        element: <RequestDemo />,
    },

    {
        key: 'Search',
        path: '/cloudql',
        element: <Search />,
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
    // {
    //     key: 'Search Main',
    //     path: '/cloudql-dashboard',
    //     element: <SearchLanding />,
    // },
    // {
    //     key: 'test',
    //     path: '/test',
    //     element: <Test />,
    // },
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
