export const capitalizeFirstLetter = (string: string) => {
    const splitStr = string.toLowerCase().split(' ')
    for (let i = 0; i < splitStr.length; i += 1) {
        // You do not need to check if i is larger than splitStr length, as your for does that for you
        // Assign it back to the array
        splitStr[i] =
            splitStr[i].charAt(0).toUpperCase() + splitStr[i].substring(1)
    }
    // Directly return the joined string
    return splitStr.join(' ')
}

export const snakeCaseToLabel = (string: string) =>
    capitalizeFirstLetter(
        string
            .toLowerCase()
            .replace(/([-_][a-z])/g, (group) => group.replace('_', ' '))
    )
export const kebabCaseToLabel = (string: string) =>
    capitalizeFirstLetter(
        string
            .toLowerCase()
            .replace(/([-_][a-z])/g, (group) => group.replace('-', ' '))
    )

export const camelCaseToLabel = (s: string) => {
    const result = s.replace(/([A-Z])/g, ' $1')
    return result.charAt(0).toUpperCase() + result.slice(1)
}
