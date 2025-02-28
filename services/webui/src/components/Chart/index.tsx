import ReactEcharts from 'echarts-for-react'
import { useAtomValue } from 'jotai'
import {
    exactPriceDisplay,
    numberDisplay,
    numericDisplay,
} from '../../utilities/numericDisplay'
import { colorBlindModeAtom } from '../../store'

interface IChart {
    labels: string[]
    labelType?: 'category' | 'time' | 'value' | 'log'
    chartData:
        | (string | number | undefined)[]
        | (
              | {
                    name: string
                    value: string
                    itemStyle?: undefined
                    label?: undefined
                }
              | {
                    value: number
                    name: string
                    itemStyle: { color: string; decal: { symbol: string } }
                    label: { show: boolean }
                }
          )[]
        | (
              | {
                    name: string
                    value: number | undefined
                    itemStyle?: undefined
                    label?: undefined
                }
              | {
                    value: number
                    name: string
                    itemStyle: { color: string; decal: { symbol: string } }
                    label: { show: boolean }
                }
          )[]
        | undefined
    chartType: 'bar' | 'line' | 'doughnut' | 'half-doughnut'
    chartAggregation?: 'trend' | 'cumulative'
    visualMap?: any
    markArea?: any
    isCost?: boolean
    isPercent?: boolean
    loading?: boolean
    error?: string
    onClick?: (param?: any) => void
    colorful?: boolean
}

export default function Chart({
    labels,
    labelType = 'category',
    chartData,
    chartType,
    chartAggregation,
    isCost = false,
    isPercent = false,
    markArea,
    visualMap,
    loading,
    error,
    onClick,
    colorful = false,
}: IChart) {
    const colorBlindMode = useAtomValue(colorBlindModeAtom)
    const options = () => {
        if (chartType === 'line' || chartType === 'bar') {
            if (chartAggregation === 'trend') {
                return {
                    aria: {
                        enable: colorBlindMode,
                        decal: {
                            show: colorBlindMode,
                        },
                    },
                    xAxis: {
                        type: labelType,
                        data: labels,
                    },
                    yAxis: {
                        type: 'value',
                        axisLabel: {
                            formatter: (value: string | number) => {
                                if (isCost) {
                                    return `$${numericDisplay(value)}`
                                }
                                if (isPercent) {
                                    return `${numericDisplay(value)} %`
                                }
                                return numericDisplay(value)
                            },
                        },
                    },
                    visualMap,
                    animation: false,
                    series: [
                        chartType === 'bar' && {
                            data: chartData,
                            type: chartType,
                            areaStyle: { opacity: 0 },
                        },
                        chartType === 'line' && {
                            data: chartData,
                            markArea,
                            type: chartType,
                            areaStyle: { opacity: 0 },
                        },
                    ],
                    grid: {
                        left: 45,
                        right: 0,
                        top: 20,
                        bottom: 40,
                    },
                    tooltip: {
                        show: true,
                        trigger: 'axis',
                        valueFormatter: (value: string | number) => {
                            if (isCost) {
                                return `$${numberDisplay(Number(value), 2)}`
                            }
                            if (isPercent) {
                                return `${numericDisplay(value)} %`
                            }
                            return numberDisplay(Number(value), 0)
                        },
                    },
                    color: colorful
                        ? [
                              '#780000',
                              '#DC0000',
                              '#FD8C00',
                              '#FDC500',
                              '#10B880',
                              '#D0D4DA',
                          ]
                        : [
                              '#1D4F85',
                              '#2970BC',
                              '#6DA4DF',
                              '#96BEE8',
                              '#C0D8F1',
                              '#D0D4DA',
                          ],
                }
            }
            if (chartAggregation === 'cumulative') {
                return {
                    aria: {
                        enable: colorBlindMode,
                        decal: {
                            show: colorBlindMode,
                        },
                    },
                    xAxis: {
                        type: labelType,
                        data: labels,
                    },
                    yAxis: {
                        type: 'value',
                        axisLabel: {
                            formatter: (value: string | number) => {
                                if (isCost) {
                                    return `$${numericDisplay(value)}`
                                }
                                if (isPercent) {
                                    return `${numericDisplay(value)} %`
                                }
                                return numericDisplay(value)
                            },
                        },
                    },
                    visualMap,
                    animation: false,
                    series: [
                        chartType === 'bar' && {
                            data: chartData,
                            type: chartType,
                            areaStyle: { opacity: 0 },
                        },
                        chartType === 'line' && {
                            data: chartData,
                            markArea,
                            type: chartType,
                            areaStyle: { opacity: 0.7 },
                        },
                    ],
                    grid: {
                        left: 45,
                        right: 0,
                        top: 20,
                        bottom: 40,
                    },
                    tooltip: {
                        show: true,
                        trigger: 'axis',
                        valueFormatter: (value: string | number) => {
                            if (isCost) {
                                return `$${numberDisplay(Number(value), 2)}`
                            }
                            if (isPercent) {
                                return `${numericDisplay(value)} %`
                            }
                            return numberDisplay(Number(value), 0)
                        },
                    },
                    color: colorful
                        ? [
                              '#780000',
                              '#DC0000',
                              '#FD8C00',
                              '#FDC500',
                              '#10B880',
                              '#D0D4DA',
                          ]
                        : [
                              '#1D4F85',
                              '#2970BC',
                              '#6DA4DF',
                              '#96BEE8',
                              '#C0D8F1',
                              '#D0D4DA',
                          ],
                }
            }
        }

        if (chartType === 'doughnut') {
            return {
                aria: {
                    enable: colorBlindMode,
                    decal: {
                        show: colorBlindMode,
                    },
                },
                tooltip: {
                    trigger: 'item',
                    formatter: (params: any) => {
                        return `${params.data.name} (${params.percent}%):\n\n${
                            isCost
                                ? exactPriceDisplay(params.data.value)
                                : numberDisplay(params.data.value, 0)
                        }`
                    },
                },
                series: [
                    {
                        type: 'pie',
                        radius: ['47%', '70%'],
                        // center: ['50%', '50%'],
                        avoidLabelOverlap: false,
                        label: {
                            show: true,
                            position: 'center',
                            formatter: () => {
                                let total = 0
                                for (
                                    let i = 0;
                                    i < (chartData ? chartData?.length : 0);
                                    i += 1
                                ) {
                                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                    // @ts-ignore
                                    total += Number(chartData[i].value)
                                }

                                return `Total: ${
                                    isCost
                                        ? exactPriceDisplay(total)
                                        : numberDisplay(total, 0)
                                }`
                            },
                        },
                        itemStyle: {
                            borderRadius: 4,
                            borderColor: '#fff',
                            borderWidth: 1,
                        },
                        data: chartData,
                        left: '-5%',
                        width: '70%',
                    },
                ],
                legend: {
                    right: 12,
                    top: 'middle',
                    icon: 'circle',
                    orient: 'vertical',
                    textStyle: {
                        width: 140,
                        overflow: 'truncate',
                    },
                },
                color: colorful
                    ? [
                          '#780000',
                          '#DC0000',
                          '#FD8C00',
                          '#FDC500',
                          '#10B880',
                          '#D0D4DA',
                      ]
                    : [
                          '#1E7CE0',
                          '#2ECC71',
                          '#FFA500',
                          '#9B59B6',
                          '#D0D4DA',
                          // '#D0D4DA',
                      ],
            }
        }
        if (chartType === 'half-doughnut') {
            return {
                aria: {
                    enable: colorBlindMode,
                    decal: {
                        show: colorBlindMode,
                    },
                },
                tooltip: {
                    trigger: 'item',
                },
                series: [
                    {
                        type: 'pie',
                        radius: ['30%', '50%'],
                        center: ['40%', '63%'],
                        // adjust the start angle
                        startAngle: 180,
                        label: {
                            show: false,
                        },
                        data: chartData,
                    },
                ],
                itemStyle: {
                    borderRadius: 4,
                    borderColor: '#fff',
                    borderWidth: 2,
                },
                color: ['#C0D8F1', '#0D2239'],
            }
        }
        return {}
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
