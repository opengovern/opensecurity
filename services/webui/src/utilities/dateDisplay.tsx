import dayjs, { Dayjs } from 'dayjs'
import LocalizedFormat from 'dayjs/plugin/localizedFormat'
import timezone from 'dayjs/plugin/timezone'
import advancedFormat from 'dayjs/plugin/advancedFormat'
import relativeTime from 'dayjs/plugin/relativeTime'
import duration from 'dayjs/plugin/duration'
dayjs.extend(LocalizedFormat)
dayjs.extend(timezone)
dayjs.extend(advancedFormat)
dayjs.extend(relativeTime)
dayjs.extend(duration)

export const dateDisplay = (
    date: Dayjs | Date | number | string | undefined,
    subtract?: number
) => {
    const s = subtract || 0
    if ((typeof date).toString() === 'Dayjs') {
        return (date as Dayjs).subtract(s, 'day').format('MMM DD, YYYY')
    }
    if (date) {
        return dayjs.utc(date).subtract(s, 'day').format('MMM DD, YYYY')
    }
    return 'Not available'
}

export const monthDisplay = (
    date: Dayjs | Date | number | string | undefined,
    subtract?: number
) => {
    const s = subtract || 0
    if ((typeof date).toString() === 'Dayjs') {
        return (date as Dayjs).subtract(s, 'day').format('MMM, YYYY')
    }
    if (date) {
        return dayjs.utc(date).subtract(s, 'day').format('MMM, YYYY')
    }
    return 'Not available'
}

export const dateTimeDisplay = (
    date: Dayjs | Date | number | string | undefined
) => {
    // tz(dayjs.tz.guess())
    if ((typeof date).toString() === 'Dayjs') {
        return (date as Dayjs).format('MMM DD, YYYY kk:mm:ss UTC')
    }
    const regexp = /^\d+$/g
    const isNumber = regexp.test(String(date))

    if (isNumber) {
        const v = parseInt(String(date), 10)
        const value = v > 17066236800 ? v / 1000 : v
        return dayjs.unix(value).utc().format('MMM DD, YYYY kk:mm:ss UTC')
    }
    if (date) {
      
        return dayjs.utc(date).format('MMM DD, YYYY kk:mm:ss UTC')
    }
    return 'Not available'
}

export const dateTimeDisplayAgo = (
    date: Dayjs | Date | number | string | undefined
): string => {
    if (!date) return 'Not available'

    let parsedDate: Dayjs

    // Check if it's already a Dayjs object
    if (dayjs.isDayjs(date)) {
        parsedDate = date
    } else {
        // If it's a number, determine if it's in seconds or milliseconds
        const isNumber = /^\d+$/g.test(String(date))
        if (isNumber) {
            const v = parseInt(String(date), 10)
            const value = v > 17066236800 ? v / 1000 : v // Convert if timestamp is too large
            parsedDate = dayjs.unix(value).utc()
        } else {
            // Parse as a normal date
            parsedDate = dayjs.utc(date)
        }
    }

    // Calculate exact time difference
    const now = dayjs.utc()
    const diffMs = now.diff(parsedDate, 'milliseconds')
    const diffDuration = dayjs.duration(diffMs)

    const days = diffDuration.days()
    const hours = diffDuration.hours()
    const minutes = diffDuration.minutes()
    const seconds = diffDuration.seconds()

    let result = ''
    if (days > 0) result += `${days} day${days !== 1 ? 's' : ''} `
    if (hours > 0) result += `${hours} h `
    if (minutes > 0) result += `${minutes} m `
    if (seconds > 0 && result === '')
        result += `${seconds} s ` // Show seconds only if no larger unit

    return result.trim() + ' ago'
}

export const shortDateTimeDisplay = (
    date: Dayjs | Date | number | string | undefined
) => {
    // tz(dayjs.tz.guess())
    if ((typeof date).toString() === 'Dayjs') {
        return (date as Dayjs).format('MM-DD-YYYY HH:mm')
    }
    if (date) {
        return dayjs.utc(date).format('MM-DD-YYYY HH:mm')
    }
    return 'Not available'
}


export const shortDateTimeDisplayDelta = (
    date: Dayjs | Date | number | string | undefined,date2: Dayjs | Date | number | string | undefined
) => {
    
    
    if (date && date2) {
        const d1 = dayjs.utc(date)
        const d2 = dayjs.utc(date2)

        const diff = d1.diff(d2, 'ms')
        const minutes = Math.floor(diff / 60000) // 1 minute = 60000 ms
        const seconds = Math.floor((diff % 60000) / 1000) // Remaining seconds
        const milliseconds = diff % 1000 // Remaining milliseconds
        if (minutes > 0) {
            return `${minutes} min ${seconds} sec`
        }
        if(seconds > 0){
            return `${seconds} sec`
        }
        return `${milliseconds} ms`
    }
    return 'Not available'
}
