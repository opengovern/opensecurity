import Cal, { getCalApi } from '@calcom/embed-react'
import { Text, Title } from '@tremor/react'
import { useSearchParams } from 'react-router-dom'
import { useEffect } from 'react'
import TopHeader from '../../components/Layout/Header'

export default function RequestAccess() {
    const [searchParams, setSearchParams] = useSearchParams()
    const title = () => {
        if (window.location.href.indexOf('workload-optimizer') > 0) {
            return `Workload Optimizer enhances Cloud-Native and Kubernetes efficiency by analyzing trends and usage, delivering intelligent recommendations to optimize instance sizes while maintaining reliability.`
        }
        if (window.location.href.indexOf('stacks') > 0) {
            return `Dynamically group infrastructure, code repos, and entities as Apps & Environments for targeted governance.`
        }
        // if (window.location.href.indexOf('dashboard') > 0) {
        //     return `Create custom dashboards for teams, functions, and roles—In-App or through BI tools like PowerBI, Tableau, Looker, or Grafana.`
        // }
         if (window.location.href.indexOf('automation') > 0) {
             return `Automatically triggers runbooks or external APIs in response to compliance events, streamlining operations and enforcing policies efficiently.`
         }
        return `${searchParams.get(
            'connector'
        )} and 50+ others are available for Enterprise Users. Get a 30-day obligation trial now.`
    }
    const f = async () => {
        const cal = await getCalApi({ namespace: 'try-enterprise' })
        cal('ui', {
            styles: { branding: { brandColor: '#000000' } },
            hideEventTypeDetails: false,
            layout: 'month_view',
        })
    }

    useEffect(() => {
        f()
    }, [])

    return (
        <>
            <TopHeader />
            <Title className="text-black !text-xl font-bold w-full text-center mb-4">
                {title()}
            </Title>
            <Cal
                namespace="try-enterprise"
                calLink="team/opencomply/try-enterprise"
                style={{ width: '100%', height: '100%', overflow: 'scroll' }}
                config={{ layout: 'month_view' }}
            />
        </>
    )
}
