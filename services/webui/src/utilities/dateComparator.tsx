import dayjs, { Dayjs } from 'dayjs'
import { Select, SelectItem, Text } from '@tremor/react'

export const agGridDateComparator = (
    filterLocalDateAtMidnight: Date,
    cellValue: string
) => {
    const dateAsString = cellValue

    if (dateAsString == null) {
        return 0
    }

    const dateValue = new Date(Date.parse(dateAsString))
    const year = dateValue.getFullYear()
    const month = dateValue.getMonth()
    const day = dateValue.getDate()
    const cellDate = new Date(year, month, day)

    if (cellDate < filterLocalDateAtMidnight) {
        return -1
    }
    if (cellDate > filterLocalDateAtMidnight) {
        return 1
    }
    return 0
}

export const checkGranularity = (start: Dayjs, end: Dayjs) => {
    const daily = true
    let monthly = true
    let yearly = true
    let monthlyByDefault = true

    if (dayjs.utc(end).diff(dayjs.utc(start), 'month', true) < 2) {
        monthlyByDefault = false
    }
    if (dayjs.utc(end).diff(dayjs.utc(start), 'month', true) < 1) {
        monthly = false
    }
    if (dayjs.utc(end).diff(dayjs.utc(start), 'year', true) < 1) {
        yearly = false
    }

    return {
        daily,
        monthly,
        monthlyByDefault,
        yearly,
    }
}

export const generateItems = (
    s: Dayjs,
    e: Dayjs,
    placeholder: string,
    value: string,
    onValueChange: (v: string) => void
) => {
    const generateSelectItem = (vl: string, title: string) => {
        return (
            <SelectItem value={vl}>
                <Text>{title}</Text>
            </SelectItem>
        )
    }

    const items = () => {
        const i = []
        const x = checkGranularity(s, e)
        if (x.daily) {
            i.push({
                title: 'Daily',
                value: 'daily',
            })
        }
        if (x.monthly) {
            i.push({
                title: 'Monthly',
                value: 'monthly',
            })
        }
        if (x.yearly) {
            i.push({
                title: 'Yearly',
                value: 'yearly',
            })
        }

        return i.map((v) => generateSelectItem(v.value, v.title))
    }
    return (
        <div>
            <Select
                enableClear={false}
                value={value}
                placeholder={placeholder}
                onValueChange={onValueChange}
                className="w-10 asdasdsa"
            >
                {items()}
            </Select>
        </div>
    )
}
