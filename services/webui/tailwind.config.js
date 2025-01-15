/** @type {import('tailwindcss').Config} */

module.exports = {
    darkMode: 'class',
    content: [
        './src/**/*.{js,jsx,ts,tsx}',

        // Path to the Tremor module
        './node_modules/@tremor/**/*.{js,ts,jsx,tsx}',
    ],
    theme: {
        transparent: 'transparent',
        current: 'currentColor',
        extend: {
            colors: {
                openg: {
                    50: '#EAF2FA',
                    100: '#C0D8F1',
                    200: '#96BEE8',
                    300: '#6DA4DF',
                    400: '#438AD6',
                    500: '#2970BC',
                    600: '#205792',
                    700: '#1D4F85',
                    800: '#15395F',
                    900: '#0D2239',
                    950: '#0D2239',
                },
                // light mode
                tremor: {
                    brand: {
                        faint: '#EAF2FA', // blue-50
                        muted: '#96BEE8', // blue-200
                        subtle: '#438AD6', // blue-400
                        DEFAULT: '#2970BC', // blue-500
                        emphasis: '#1D4F85', // blue-700
                        inverted: '#ffffff', // white
                    },
                    background: {
                        muted: '#f9fafb', // gray-50
                        subtle: '#f3f4f6', // gray-100
                        DEFAULT: '#ffffff', // white
                        emphasis: '#374151', // gray-700
                    },
                    border: {
                        DEFAULT: '#e5e7eb', // gray-200
                    },
                    ring: {
                        DEFAULT: '#e5e7eb', // gray-200
                    },
                    content: {
                        subtle: '#9ca3af', // gray-400
                        DEFAULT: '#6b7280', // gray-500
                        emphasis: '#374151', // gray-700
                        strong: '#111827', // gray-900
                        inverted: '#ffffff', // white
                    },
                },
                // dark mode
                'dark-tremor': {
                    brand: {
                        faint: '#0B1229', // custom
                        muted: '#172554', // blue-950
                        subtle: '#1e40af', // blue-800
                        DEFAULT: '#3b82f6', // blue-500
                        emphasis: '#60a5fa', // blue-400
                        inverted: '#FFFFFF', // gray-950
                    },
                    background: {
                        muted: '#131A2B', // custom
                        subtle: '#111827', // gray-800
                        DEFAULT: '#1F2937', // gray-900
                        emphasis: '#D0D4DA', // gray-300
                    },
                    border: {
                        DEFAULT: '#374151', // gray-800
                    },
                    ring: {
                        DEFAULT: '#374151', // gray-800
                    },
                    content: {
                        subtle: '#F2F3F5', // gray-600
                        DEFAULT: '#F2F3F5', // gray-600
                        emphasis: '#e5e7eb', // gray-200
                        strong: '#f9fafb', // gray-50
                        inverted: '#000000', // black
                    },
                },
            },
            boxShadow: {
                // light
                'tremor-input': '0 1px 2px 0 rgb(0 0 0 / 0.05)',
                'tremor-card':
                    '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
                'tremor-dropdown':
                    '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
                // dark
                'dark-tremor-input': '0 1px 2px 0 rgb(0 0 0 / 0.05)',
                'dark-tremor-card':
                    '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
                'dark-tremor-dropdown':
                    '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
            },
            borderRadius: {
                'tremor-small': '0.375rem',
                'tremor-default': '0.5rem',
                'tremor-full': '9999px',
            },
            fontSize: {
                'tremor-label': ['0.75rem'],
                'tremor-default': ['0.875rem', { lineHeight: '1.25rem' }],
                'tremor-title': ['1.125rem', { lineHeight: '1.75rem' }],
                'tremor-metric': ['1.875rem', { lineHeight: '2.25rem' }],
            },
        },
    },
    safelist: [
        {
            pattern:
                /^(bg-(?:slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-(?:50|100|200|300|400|500|600|700|800|900|950))$/,
            variants: ['hover', 'ui-selected'],
        },
        {
            pattern:
                /^(text-(?:slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-(?:50|100|200|300|400|500|600|700|800|900|950))$/,
            variants: ['hover', 'ui-selected'],
        },
        {
            pattern:
                /^(border-(?:slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-(?:50|100|200|300|400|500|600|700|800|900|950))$/,
            variants: ['hover', 'ui-selected'],
        },
        {
            pattern:
                /^(ring-(?:slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-(?:50|100|200|300|400|500|600|700|800|900|950))$/,
        },
        {
            pattern:
                /^(stroke-(?:slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-(?:50|100|200|300|400|500|600|700|800|900|950))$/,
        },
        {
            pattern:
                /^(fill-(?:slate|gray|zinc|neutral|stone|red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose)-(?:50|100|200|300|400|500|600|700|800|900|950))$/,
        },
        ...['[#000]'].flatMap((customColor) => [
            `bg-${customColor}`,
            `border-${customColor}`,
            `hover:bg-${customColor}`,
            `hover:border-${customColor}`,
            `hover:text-${customColor}`,
            `fill-${customColor}`,
            `ring-${customColor}`,
            `stroke-${customColor}`,
            `text-${customColor}`,
            `ui-selected:bg-${customColor}`,
            `ui-selected:border-${customColor}`,
            `ui-selected:text-${customColor}`,
        ]),
        ...['[#9BA2AE]'].flatMap((customColor) => [
            `bg-${customColor}`,
            `border-${customColor}`,
            `hover:bg-${customColor}`,
            `hover:border-${customColor}`,
            `hover:text-${customColor}`,
            `fill-${customColor}`,
            `ring-${customColor}`,
            `stroke-${customColor}`,
            `text-${customColor}`,
            `ui-selected:bg-${customColor}`,
            `ui-selected:border-${customColor}`,
            `ui-selected:text-${customColor}`,
        ]),
        ...['[#54B584]'].flatMap((customColor) => [
            `bg-${customColor}`,
            `border-${customColor}`,
            `hover:bg-${customColor}`,
            `hover:border-${customColor}`,
            `hover:text-${customColor}`,
            `fill-${customColor}`,
            `ring-${customColor}`,
            `stroke-${customColor}`,
            `text-${customColor}`,
            `ui-selected:bg-${customColor}`,
            `ui-selected:border-${customColor}`,
            `ui-selected:text-${customColor}`,
        ]),
        ...['[#F4C744]'].flatMap((customColor) => [
            `bg-${customColor}`,
            `border-${customColor}`,
            `hover:bg-${customColor}`,
            `hover:border-${customColor}`,
            `hover:text-${customColor}`,
            `fill-${customColor}`,
            `ring-${customColor}`,
            `stroke-${customColor}`,
            `text-${customColor}`,
            `ui-selected:bg-${customColor}`,
            `ui-selected:border-${customColor}`,
            `ui-selected:text-${customColor}`,
        ]),
        ...['[#EE9235]'].flatMap((customColor) => [
            `bg-${customColor}`,
            `border-${customColor}`,
            `hover:bg-${customColor}`,
            `hover:border-${customColor}`,
            `hover:text-${customColor}`,
            `fill-${customColor}`,
            `ring-${customColor}`,
            `stroke-${customColor}`,
            `text-${customColor}`,
            `ui-selected:bg-${customColor}`,
            `ui-selected:border-${customColor}`,
            `ui-selected:text-${customColor}`,
        ]),
        ...['[#CA2B1D]'].flatMap((customColor) => [
            `bg-${customColor}`,
            `border-${customColor}`,
            `hover:bg-${customColor}`,
            `hover:border-${customColor}`,
            `hover:text-${customColor}`,
            `fill-${customColor}`,
            `ring-${customColor}`,
            `stroke-${customColor}`,
            `text-${customColor}`,
            `ui-selected:bg-${customColor}`,
            `ui-selected:border-${customColor}`,
            `ui-selected:text-${customColor}`,
        ]),
        ...['[#6E120B]'].flatMap((customColor) => [
            `bg-${customColor}`,
            `border-${customColor}`,
            `hover:bg-${customColor}`,
            `hover:border-${customColor}`,
            `hover:text-${customColor}`,
            `fill-${customColor}`,
            `ring-${customColor}`,
            `stroke-${customColor}`,
            `text-${customColor}`,
            `ui-selected:bg-${customColor}`,
            `ui-selected:border-${customColor}`,
            `ui-selected:text-${customColor}`,
        ]),
    ],
    plugins: [
        require('@headlessui/tailwindcss'),
        require('@tailwindcss/forms'),
    ],
}
