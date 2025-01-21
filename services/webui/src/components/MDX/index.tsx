import type { MDXComponents } from 'mdx/types'
import { Bold, ChangelogEntry, CustomLink, H1, H2, H3, Li, P, Ul } from './mdx'

import { ChangelogImage } from './mdx'

let customComponents = {
    h1: H1,
    h2: H2,
    h3: H3,
    p: P,
    Bold: Bold,
    ul: Ul,
    a: CustomLink,
    ChangelogEntry: ChangelogEntry,
    ChangelogImage: ChangelogImage,
    li: Li,
}

export function useMDXComponents(components: MDXComponents) {
    return {
        ...customComponents,
        ...components,
    }
}
