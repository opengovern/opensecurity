
// @ts-nocheck

import dayjs from 'dayjs'

import { useEffect, useRef, useState } from 'react'

import {
    Box,
    DateRangePicker,
    Icon,
    SpaceBetween,
    Spinner,
} from '@cloudscape-design/components'

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
    name,
}: IDatepicker) {
   
    const [val, setVal] = useState({
        startDate: dayjs(activeTimeRange.start).toISOString(),
        endDate: dayjs(activeTimeRange.end).toISOString(),
        type: 'absolute',
    })

    useEffect(() => {
        const start = val.startDate
        const end = val.endDate

        setActiveTimeRange({
            start: dayjs(val.startDate),
            end: dayjs(val.endDate),
        })
    }, [val, checked, startH, startM, endH, endM])

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
