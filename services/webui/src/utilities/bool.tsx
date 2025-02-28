export function ConvertToBoolean(input: string): boolean | undefined {
    try {
        return JSON.parse(input)
    } catch (e) {
        return undefined
    }
}
