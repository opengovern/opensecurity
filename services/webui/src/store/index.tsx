import { atom } from 'jotai'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
import {
    PlatformEnginePkgAuthApiGetMeResponse,
    PlatformEnginePkgWorkspaceApiWorkspaceResponse,
} from '../api/api'

dayjs.extend(utc)

interface INotification {
    text: string | undefined
    type: 'success' | 'warning' | 'error' | 'info' | undefined
    position?: 'topLeft' | 'topRight' | 'bottomRight' | 'bottomLeft'
}

export const notificationAtom = atom<INotification>({
    text: undefined,
    type: undefined,
    position: undefined,
})

export const sideBarCollapsedAtom = atom(
    window.innerWidth < 768 ?true : localStorage.collapse ? localStorage.collapse === 'true' : true 
)
export const complianceOpenAtom = atom(false)
export const automationOpenAtom = atom(false)
export const queryAtom = atom('')
export const sampleAtom = atom(localStorage.getItem('sample') || false)
export const ForbiddenAtom = atom(false)
export const RoleAccess = atom(false)


export const isDemoAtom = atom(localStorage.demoMode === 'true')
export const workspaceAtom = atom<{
    list: PlatformEnginePkgWorkspaceApiWorkspaceResponse[]
    current:
        | PlatformEnginePkgWorkspaceApiWorkspaceResponse
        | undefined
}>({ list: [], current: undefined })
export const previewAtom = atom(
    localStorage.preview === 'true' ||
        localStorage.preview === null ||
        localStorage.preview === undefined
        ? 'true'
        : 'false'
)
export const runQueryAtom = atom('')
export const LayoutAtom = atom(undefined)

export const tokenAtom = atom<string>('')
export const colorBlindModeAtom = atom<boolean>(false)
export const meAtom = atom<PlatformEnginePkgAuthApiGetMeResponse>()
