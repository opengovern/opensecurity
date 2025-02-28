import { useEffect } from 'react'
import { useAuth } from '../../utilities/auth'

const Logout = () => {
    const { logout } = useAuth()
    useEffect(() => {
        logout()
    }, [logout])
    return <>Logging out</>
}

export default Logout
