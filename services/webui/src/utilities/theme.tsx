import { PlatformEnginePkgAuthApiTheme } from '../api/api'

export const parseTheme = (
    v: string
): PlatformEnginePkgAuthApiTheme => {
    switch (v) {
        case 'light':
            return PlatformEnginePkgAuthApiTheme.ThemeLight
        case 'dark':
            return PlatformEnginePkgAuthApiTheme.ThemeDark
        default:
            return PlatformEnginePkgAuthApiTheme.ThemeSystem
    }
}

export const currentTheme = () => {
    if (!('theme' in localStorage)) {
        return PlatformEnginePkgAuthApiTheme.ThemeLight
    }

    return parseTheme(localStorage.theme)
}

export const applyTheme = (v: PlatformEnginePkgAuthApiTheme) => {
    if (
        v === PlatformEnginePkgAuthApiTheme.ThemeDark ||
        (v === PlatformEnginePkgAuthApiTheme.ThemeSystem &&
            window.matchMedia('(prefers-color-scheme:dark)').matches)
    ) {
        document.documentElement.classList.add('dark')
    } else {
        document.documentElement.classList.remove('dark')
    }

    switch (v) {
        case PlatformEnginePkgAuthApiTheme.ThemeDark:
            localStorage.theme = 'dark'
            break
        case PlatformEnginePkgAuthApiTheme.ThemeLight:
            localStorage.theme = 'light'
            break
        default:
            localStorage.removeItem('theme')
    }
}
