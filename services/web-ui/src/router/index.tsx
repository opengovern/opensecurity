import { Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import { useEffect } from 'react'
import Assets from '../pages/Assets'
import NotFound from '../pages/Errors'
import { CallbackPage } from '../pages/Callback'
import Settings from '../pages/Settings'
import Workspaces from '../pages/Workspaces'
import Logout from '../pages/Logout'
import Integrations from '../pages/Integrations'
import Compliance from '../pages/Governance/Compliance'
import Overview from '../pages/Overview'
import Stack from '../pages/Stack'
import Single from '../pages/Assets/Single'
import SingleSpend from '../pages/Spend/Single'
// import Boostrap from '../pages/Workspaces/Bootstrap'
import ControlDetail from '../pages/Governance/Controls/ControlSummary'
import Findings from '../pages/Governance/Findings'
import { SpendOverview } from '../pages/Spend/Overview'
import { SpendMetrics } from '../pages/Spend/Metric'
import { SpendAccounts } from '../pages/Spend/Account'
import Layout from '../components/Layout'
import RequestDemo from '../pages/RequestDemo'
import AssetAccounts from '../pages/Assets/Account'
import AssetMetrics from '../pages/Assets/Metric'
import WorkloadOptimizer from '../pages/WorkloadOptimizer'
import RequestAccess from '../pages/Integrations/RequestAccess'
import SettingsJobs from '../pages/Settings/Jobs'
import AllControls from '../pages/Governance/Compliance/All Controls'
import AllBenchmarks from '../pages/Governance/Compliance/All Benchmarks'
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
        key: 'infrastructure',
        path: '/dashboard/infrastructure',
        element: <Assets />,
    },
    {
        key: 'infrastructure single',
        path: '/dashboard/infrastructure/:id',
        element: <Single />,
    },
    {
        key: 'infrastructure single metric',
        path: '/dashboard/infrastructure/:id/:metric',
        element: <Single />,
    },
    {
        key: 'infrastructure single metric',
        path: '/dashboard/infrastructure-cloud-account/:id/:metric',
        element: <Single />,
    },
    {
        key: 'infrastructure account detail',
        path: '/dashboard/infrastructure-cloud-accounts',
        element: <AssetAccounts />,
    },
    {
        key: 'infrastructure account detail single',
        path: '/dashboard/infrastructure-cloud-accounts/:id/:metric',
        element: <Single />,
    },
    {
        key: 'infrastructure account detail single',
        path: '/dashboard/infrastructure-cloud-accounts/:id',
        element: <Single />,
    },
    {
        key: 'infrastructure metric detail',
        path: '/dashboard/infrastructure-metrics',
        element: <AssetMetrics />,
    },
    {
        key: 'infrastructure single 2',
        path: '/dashboard/infrastructure-metrics/:id',
        element: <Single />,
    },
    {
        key: 'infrastructure single metric 2',
        path: '/dashboard/infrastructure-metrics/:id/:metric',
        element: <Single />,
    },
    {
        key: 'spend',
        path: '/dashboard/spend',
        element: <SpendOverview />,
    },
    {
        key: 'spend single 1',
        path: '/dashboard/spend/:id',
        element: <SingleSpend />,
    },
    {
        key: 'spend single metric 1',
        path: '/dashboard/spend/:id/:metric',
        element: <SingleSpend />,
    },
    {
        key: 'spend',
        path: '/dashboard/spend-metrics',
        element: <SpendMetrics />,
    },
    {
        key: 'spend',
        path: '/dashboard/spend-accounts',
        element: <SpendAccounts />,
    },
    {
        key: 'spend single',
        path: '/dashboard/spend-accounts/:id',
        element: <SingleSpend />,
    },
    {
        key: 'spend single metric',
        path: '/dashboard/spend-accounts/:id/:metric',
        element: <SingleSpend />,
    },
    {
        key: 'spend single',
        path: '/dashboard/spend-metrics/:id',
        element: <SingleSpend />,
    },
    {
        key: 'spend single metric',
        path: '/dashboard/spend-metrics/:id/:metric',
        element: <SingleSpend />,
    },
    {
        key: 'spend single 2',
        path: '/dashboard/spend/spend-details/:id',
        element: <SingleSpend />,
    },
    {
        key: 'spend single metric 2',
        path: '/dashboard/spend/spend-details/:id/:metric',
        element: <SingleSpend />,
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
        key: 'resource collection infrastructure single 2',
        path: '/resource-collection/:resourceId/infrastructure-details/:id',
        element: <Single />,
    },
    {
        key: 'resource collection infrastructure single metric 2',
        path: '/resource-collection/:resourceId/infrastructure-details/:id/:metric',
        element: <Single />,
    },
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
