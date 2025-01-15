import { useSearchParams } from 'react-router-dom'
import { useAuth } from '../../utilities/auth'
import { useWorkspaceApiV3GetShouldSetup } from '../../api/workspace.gen'
import Spinner from '../../components/Spinner'

export const CallbackPage = () => {
    const [locationSearchParams, setSearchParams] = useSearchParams()

    const { error, isAuthenticated } = useAuth()
    
     if (isAuthenticated) {
         const c = sessionStorage.getItem('callbackURL')

         window.location.href =
             c === null || c === undefined || c === '' ? '/' : c
         return null
     }

     if (locationSearchParams.has('error_description')) {
         return <span>{locationSearchParams.get('error_description')}</span>
     }

     if (error) {
         return <span>{error.message}</span>
     }
     return null
}
