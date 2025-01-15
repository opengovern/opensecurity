// @ts-nocheck

import { getLocalTimeZone, parseDate, today } from '@internationalized/date'
import dayjs from 'dayjs'
import { AriaDateRangePickerProps, DateValue } from '@react-aria/datepicker'
import { useDateRangePickerState } from 'react-stately'
import { useEffect, useRef, useState } from 'react'
import { useDateRangePicker } from 'react-aria'
import { Checkbox } from 'pretty-checkbox-react'
import { ClockIcon } from '@heroicons/react/24/outline'
import { Flex, Select, SelectItem, Text, Title } from '@tremor/react'
import { renderDateText } from '../../../../../../../components/Layout/Header/DatePicker'
import { RangeCalendar } from '../../../../../../../components/Layout/Header/DatePicker/Calendar/RangePicker/RangeCalendar'
import {
    Box,
    DateRangePicker,
    Icon,
    SpaceBetween,
    Spinner,
} from '@cloudscape-design/components'
function CustomDatePicker(props: AriaDateRangePickerProps<DateValue>) {
    const state = useDateRangePickerState(props)
    const ref = useRef(null)

    const { calendarProps } = useDateRangePicker(props, state, ref)

    return <RangeCalendar {...calendarProps} />
}

export interface IDate {
    start: dayjs.Dayjs
    end: dayjs.Dayjs
}

interface IDatepicker {
    condition: string
    activeTimeRange: IDate
    setActiveTimeRange: (v: IDate) => void
    name: string
}

export default function Datepicker({
    condition,
    activeTimeRange,
    setActiveTimeRange,
    name
}: IDatepicker) {
    const [startH, setStartH] = useState(activeTimeRange.start.hour())
    const [startM, setStartM] = useState(activeTimeRange.start.minute())
    const [endH, setEndH] = useState(activeTimeRange.end.hour())
    const [endM, setEndM] = useState(activeTimeRange.end.minute())
    const [checked, setChecked] = useState(
        startH !== 0 || startM !== 0 || endH !== 23 || endM !== 59
    )
    const [val, setVal] = useState({
        startDate: dayjs(activeTimeRange.start).toISOString(),
        endDate: dayjs(activeTimeRange.end).toISOString(),
        type: 'absolute'
    })

    useEffect(() => {
        const start = val.startDate
        const end = val.endDate
     
        setActiveTimeRange({
            start: dayjs(val.startDate),
            end: dayjs(val.endDate),
        })
    }, [val, checked, startH, startM, endH, endM])

    const minValue = () => {
        return parseDate('2022-12-01')
    }
    const maxValue = () => {
        return today(getLocalTimeZone())
    }

    const last7Days = () => {
        return {
            start: dayjs().utc().subtract(1, 'week').startOf('day'),
            end: dayjs().utc().endOf('day'),
        }
    }

    const last30Days = () => {
        return {
            start: dayjs().utc().subtract(1, 'month').startOf('day'),
            end: dayjs().utc().endOf('day'),
        }
    }

    const thisMonth = () => {
        return {
            start: dayjs().utc().startOf('month').startOf('day'),
            end: dayjs().utc().endOf('day'),
        }
    }

    const lastMonth = () => {
        return {
            start: dayjs()
                .utc()
                .subtract(1, 'month')
                .startOf('month')
                .startOf('day'),
            end: dayjs().utc().subtract(1, 'month').endOf('month').endOf('day'),
        }
    }

    const thisQuarter = () => {
        return {
            start: dayjs().utc().startOf('quarter').startOf('day'),
            end: dayjs().utc().endOf('day'),
        }
    }

    const lastQuarter = () => {
        return {
            start: dayjs()
                .utc()
                .subtract(1, 'quarter')
                .startOf('quarter')
                .startOf('day'),
            end: dayjs()
                .utc()
                .subtract(1, 'quarter')
                .endOf('quarter')
                .endOf('day'),
        }
    }

    const thisYear = () => {
        return {
            start: dayjs().utc().startOf('year').startOf('day'),
            end: dayjs().utc().endOf('day'),
        }
    }

    return (
        <>
            <DateRangePicker
                onChange={({ detail }) => {
                    setVal(detail.value)
                }}
                value={val}
                absoluteFormat="long-localized"
                hideTimeOffset
                rangeSelectorMode={'absolute-only'}
                isValidRange={(range) => {
                    if (range.type === 'absolute') {
                        const [startDateWithoutTime] =
                            range.startDate.split('T')
                        const [endDateWithoutTime] = range.endDate.split('T')
                        if (!startDateWithoutTime || !endDateWithoutTime) {
                            return {
                                valid: false,
                                errorMessage:
                                    'The selected date range is incomplete. Select a start and end date for the date range.',
                            }
                        }
                        if (
                            new Date(range.startDate) -
                                new Date(range.endDate) >
                            0
                        ) {
                            return {
                                valid: false,
                                errorMessage:
                                    'The selected date range is invalid. The start date must be before the end date.',
                            }
                        }
                    }
                    return { valid: true }
                }}
                i18nStrings={{}}
                placeholder={name}
            />
        </>
    )
}
