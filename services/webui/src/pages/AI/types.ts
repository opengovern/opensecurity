
export interface ChatList {
    // uknown key
    [key: string]: Chat;
}

export interface Chat {
    message?: string
    pre_loaded?: boolean
    id?: string
    loading: boolean
    suggestions?: string[]
    messageTime?: string
    responseTime?: string
    text?: string
    isWelcome?: boolean
    error?: string
    response?: any
    show?: boolean
    time?: number
    clarify_needed?: boolean
    clarify_answer?: string
    clarify_questions?: ClarifyQuestion[]
}


interface ClarifyQuestion {
    question: string
}
