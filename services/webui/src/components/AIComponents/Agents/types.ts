export interface Agent {
    id: string
    name: string
    description: string
    is_available: boolean
    enabled: boolean
    welcome_message: string
    sample_questions: string[]
}
