import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
import { atom, useAtom } from 'jotai'
import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { SourceType } from '../api/api'

dayjs.extend(utc)

export interface DateRange {
    start: dayjs.Dayjs
    end: dayjs.Dayjs
}

export function defaultTime(wsName: string) {
    if (wsName === 'genco-olive') {
        const v: DateRange = {
            start: dayjs.utc('2024-01-07').startOf('day'),
            end: dayjs.utc('2024-01-14').endOf('day'),
        }
        return v
    }
    return {
        start: dayjs.utc().add(-7, 'days').startOf('day'),
        end: dayjs.utc().endOf('day'),
    }
}

export function defaultFindingsTime(wsName: string) {
    if (wsName === 'genco-olive') {
        const v: DateRange = {
            start: dayjs.utc('2024-01-07').startOf('day'),
            end: dayjs.utc('2024-01-14').endOf('day'),
        }
        return v
    }
    return {
        start: dayjs.utc().add(-7, 'days').startOf('day'),
        end: dayjs.utc().endOf('day'),
    }
}
export function defaultEventTime(wsName: string) {
    if (wsName === 'genco-olive') {
        const v: DateRange = {
            start: dayjs.utc('2024-01-07').startOf('day'),
            end: dayjs.utc('2024-01-14').endOf('day'),
        }
        return v
    }
    return {
        start: dayjs.utc().add(-7, 'days').startOf('day'),
        end: dayjs.utc().endOf('day'),
    }
}

export function defaultSpendTime(wsName: string) {
    if (wsName === 'genco-olive') {
        const v: DateRange = {
            start: dayjs.utc('2024-01-07').startOf('day'),
            end: dayjs.utc('2024-01-14').endOf('day'),
        }
        return v
    }
    return {
        start: dayjs.utc().add(-30, 'days').startOf('day'),
        end: dayjs.utc().add(-2, 'days').endOf('day'),
    }
}

export function defaultHomepageTime() {
    const v: DateRange = {
        start: dayjs.utc().add(-7, 'days').startOf('day'),
        end: dayjs.utc().endOf('day'),
    }
    return v
}

export interface IFilter {
    provider: SourceType
    connections: string[]
    connectionGroup: string[]
}

const getLocationSearch = () =>
    window.location.search.at(0) === '?'
        ? window.location.search.slice(1, 100000)
        : window.location.search

export const searchAtom = atom<string>(getLocationSearch())
export const oldUrlAtom = atom<string>('')
export const nextUrlAtom = atom<string>('')




export function useURLState<T>(
    defaultValue: T,
    serialize: (v: T) => Map<string, string[]>,
    deserialize: (v: Map<string, string[]>) => T
): [T, (v: T) => void] {
    const [searchParams] = useSearchParams()
    const [search, setSearch] = useAtom(searchAtom)

    const currentParams = () => {
        const params = new URLSearchParams()
        getLocationSearch()
            .split('&')
            .map((v2) =>
                v2
                    .split('=')
                    .map((v3) => decodeURIComponent(v3.replace(/\+/g, '%20')))
            )
            .forEach((i) => {
                if (i[0] !== '') {
                    params.append(i[0], i[1])
                }
            })
        return params
    }

    const currentValue = () => {
        const serialized = serialize(defaultValue)
        const params = currentParams()

        const v: [string, string[]][] = []
        serialized.forEach((defValue, key) => {
            const value = params.has(key) ? params.getAll(key) : defValue
            v.push([key, value])
        })

        return deserialize(new Map(v))
    }

    const [state, setState] = useState<T>(currentValue())

    const updateSearch = () => {
        // keeping search in sync with url
        const nSearch = getLocationSearch()
        if (search !== nSearch) {
            setSearch(nSearch)
        }
    }

    const setValue = (v: T) => {
        const defaultVal = serialize(defaultValue)
        const serialized = serialize(v)
        const newParams = new URLSearchParams()
        currentParams().forEach((value, key) => newParams.append(key, value))

        serialized.forEach((value, key) => {
            const defaultArray = defaultVal.get(key) || []

            newParams.delete(key)

            const isDefault =
                false &&
                defaultArray.length === value.length &&
                defaultArray.filter((item) => value.indexOf(item) >= 0)
                    .length === defaultArray.length

            if (!isDefault) {
                value.forEach((item) => {
                    newParams.append(key, item)
                })
            }
        })

        window.history.pushState({}, '', `?${newParams.toString()}`)
        updateSearch()
    }

    useEffect(() => {
        setState(currentValue())
    }, [search])

    useEffect(() => {
        updateSearch()
    }, [searchParams])

    return [state, setValue]
}

export function useURLParam<T>(
    urlParam: string,
    defaultValue: T,
    serialize?: (v: T) => string,
    deserialize?: (v: string) => T
) {
    const serializeFn =
        serialize !== undefined ? serialize : (v: T) => String(v)
    const deserializeFn =
        deserialize !== undefined ? deserialize : (v: string) => v as T
    return useURLState<T>(
        defaultValue,
        (v) => {
            const res = new Map<string, string[]>()
            res.set(urlParam, [serializeFn(v)])
            return res
        },
        (v) => {
            const m = v.get(urlParam) || []
            return deserializeFn(m[0])
        }
    )
}

export function useURLStringState(urlParam: string, defaultValue: string) {
    const [state, setState] = useURLState<string>(
        defaultValue,
        (v) => {
            const res = new Map<string, string[]>()
            res.set(urlParam, [v])
            return res
        },
        (v) => {
            const m = v.get(urlParam) || []
            return m[0]
        }
    )
    return {
        value: state,
        setValue: setState,
    }
}

export function useFilterState() {
    const [state, setState] = useURLState<IFilter>(
        {
            provider: SourceType.Nil,
            connections: [],
            connectionGroup: [],
        },
        (v) => {
            const res = new Map<string, string[]>()
            res.set(
                'provider',
                v.provider !== SourceType.Nil ? [v.provider] : []
            )
            res.set('connections', v.connections)
            return res
        },
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        (v) => {
            return {
                provider: (v.get('provider')?.at(0) as SourceType) || '',
                connections: v.get('connections') || [],
                connectionGroup: [],
            }
        }
    )
    return {
        value: state,
        setValue: setState,
    }
}

export function useUrlDateRangeState(
    defaultValue: DateRange,
    startDateParam = 'startDate',
    endDateParam = 'endDate'
) {
    const parseValue = (v: string) => {
        return dayjs.utc(v.replaceAll('+', ' '), 'YYYY-MM-DD HH:mm:ss')
    }
    const toString = (v: dayjs.Dayjs) => {
        return v.format('YYYY-MM-DD HH:mm:ss')
    }
    const [state, setState] = useURLState<DateRange>(
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        defaultValue,
        (v) => {
            const res = new Map<string, string[]>()
            res.set(startDateParam, [toString(v.start)])
            res.set(endDateParam, [toString(v.end)])
            return res
        },
        (v) => {
            return {
                start: parseValue((v.get(startDateParam) || [])[0]),
                end: parseValue((v.get(endDateParam) || [])[0]),
            }
        }
    )
    return {
        value: state,
        setValue: setState,
    }
}
