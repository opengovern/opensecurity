const defaultOptions: Intl.NumberFormatOptions = {
    notation: 'compact',
    maximumSignificantDigits: 3,
}

export const numericDisplay = (
    value: string | number | undefined,
    options?: Intl.NumberFormatOptions
) => {
    if (!value) {
        return '0'
    }
    const num = parseInt(value.toString(), 10)
    const formatter = Intl.NumberFormat('en', { ...defaultOptions, ...options })
    return num ? formatter.format(num) : '0'
}

export const numberGroupedDisplay = (value: string | number | undefined) => {
    return `${parseFloat(value ? value.toString() : '0')
        .toFixed(0)
        .toString()
        .replace(/\B(?=(\d{3})+(?!\d))/g, ',')}`
}

export const exactPriceDisplay = (
    value: string | number | undefined,
    decimals = 2
) => {
    return Number(value) || Number(value) === 0
        ? `$${Number(value)
              .toFixed(decimals)
              .toString()
              .replace(/\B(?=(\d{3})+(?!\d))/g, ',')}`
        : 'Not available'
}

export const numberDisplay = (
    value: string | number | undefined,
    decPoint = 2
) => {
    return parseFloat(value ? value.toString() : '0')
        .toFixed(decPoint)
        .toString() === 'Infinity'
        ? '0'
        : `${parseFloat(value ? value.toString() : '0')
              .toFixed(decPoint)
              .toString()
              .replace(/\B(?=(\d{3})+(?!\d))/g, ',')}`
}
