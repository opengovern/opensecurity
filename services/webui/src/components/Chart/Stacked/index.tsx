import ReactEcharts from 'echarts-for-react'
import { useAtomValue } from 'jotai'
import {
    exactPriceDisplay,
    numericDisplay,
} from '../../../utilities/numericDisplay'
import { colorBlindModeAtom } from '../../../store'

export interface StackItem {
    label: string
    value: number
}

interface IChart {
    labels: string[]
    labelType?: 'category' | 'time' | 'value' | 'log'
    chartData: StackItem[][]
    chartType?: 'bar' | 'line'
    isCost: boolean
    loading?: boolean
    error?: string
    onClick?: (param?: any) => void
    confine?: boolean
    colors?: string[]
}

export default function StackedChart({
    labels,
    labelType = 'category',
    chartData,
    chartType = 'bar',
    isCost,
    loading,
    error,
    onClick,
    confine = true,
    colors,
}: IChart) {
    const colorBlindMode = useAtomValue(colorBlindModeAtom)

    const uniqueStackLabels = chartData
        .flatMap((v) => v.map((i) => i.label))
        .filter((l, idx, arr) => arr.indexOf(l) === idx)
    const series = uniqueStackLabels.map((label) => {
        return {
            name: label,
            type: chartType === 'bar' ? 'bar' : 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
                focus: 'series',
            },
            data: chartData.map(
                (v) =>
                    v
                        .filter((i) => i.label === label)
                        .map((i) => i.value)
                        .at(0) || 0
            ),
            itemStyle: {
                borderWidth: 0.5,
                borderType: 'solid',
                borderColor: '#73c0de',
            },
        }
    })

    const options = () => {
        const opt = {
            aria: {
                enable: colorBlindMode,
                decal: {
                    show: colorBlindMode,
                },
            },
            tooltip: {
                confine,
                trigger: 'axis',
                axisPointer: {
                    type: 'line',
                    label: {
                        formatter: (param: any) => {
                            let total = 0
                            if (param.seriesData && param.seriesData.length) {
                                for (
                                    let i = 0;
                                    i < param.seriesData.length;
                                    i += 1
                                ) {
                                    total += param.seriesData[i].data
                                }
                            }

                            return `${param.value} (Total: ${
                                isCost
                                    ? exactPriceDisplay(total)
                                    : total.toFixed(2)
                            })`
                        },
                        // backgroundColor: '#6a7985',
                    },
                },
                valueFormatter: (value: number | string) => {
                    if (isCost) {
                        return exactPriceDisplay(value)
                    }
                    return numericDisplay(value)
                },
                order: 'valueDesc',
            },
            grid: {
                left: 45,
                right: 0,
                top: 20,
                bottom: 40,
            },
            xAxis: {
                type: labelType,
                data: labels,
            },
            yAxis: {
                type: 'value',
                axisLabel: {
                    formatter: (value: number | string) => {
                        if (isCost) {
                            return `$${numericDisplay(value)}`
                        }
                        return numericDisplay(value)
                    },
                },
            },
            animation: false,
            barWidth: '40%',
            // color: false
            //     ? [
            //           '#780000',
            //           '#DC0000',
            //           '#FD8C00',
            //           '#FDC500',
            //           '#10B880',
            //           '#D0D4DA',
            //       ]
            //     : [
            //           '#1E7CE0',
            //           '#2ECC71',
            //           '#FFA500',
            //           '#9B59B6',
            //           '#D0D4DA',
            //           '#D0D4DA',
            //       ],
            series,
        }

        return {
            ...opt,
            ...(colors && colors.length > 0 && { color: colors }),
        }
    }

    const onEvents = {
        click: (params: any) => (onClick ? onClick(params) : undefined),
    }

    if (loading || (error !== undefined && error !== '')) {
        return (
            <div
                className={`${
                    error !== undefined && error !== '' ? '' : 'animate-pulse'
                } h-72 mb-2 w-full bg-slate-200 dark:bg-slate-700 rounded`}
            />
        )
    }

    return (
        <ReactEcharts
            option={options()}
            showLoading={loading}
            className="w-full"
            onEvents={
                chartType === 'bar' || chartType === 'line'
                    ? onEvents
                    : undefined
            }
        />
    )
}
