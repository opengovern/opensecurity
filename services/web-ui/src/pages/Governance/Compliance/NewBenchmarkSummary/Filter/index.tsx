import { Button, Card, Flex, Icon, Text } from '@tremor/react'
import { Popover, Transition } from '@headlessui/react'
import {
    CalendarIcon,
    CheckCircleIcon,
    ChevronDownIcon,
    CloudIcon,
    PlusIcon,
    TrashIcon,
} from '@heroicons/react/24/outline'
import { Fragment, useEffect, useState } from 'react'
import dayjs from 'dayjs'
import { useParams } from 'react-router-dom'
import Provider from './Provider'
import {
    PlatformEnginePkgComplianceApiConformanceStatus,
    SourceType,
    TypesFindingSeverity,
} from '../../../../../api/api'
import ConformanceStatus from './ConformanceStatus'
import { useComplianceApiV1FindingsFiltersCreate } from '../../../../../api/compliance.gen'
import Others from './Others'
import FindingLifecycle from './FindingLifecycle'
import { compareArrays } from '../../../../../components/Layout/Header/Filter'
import ConditionDropdown from '../../../../../components/ConditionDropdown'

import {
    CloudConnect,
    Compliance,
    Control,
    Id,
    Lifecycle,
    Resources,
    SeverityIcon,
} from '../../../../../icons/icons'
import Severity from './Severity'
import Datepicker, { IDate } from './Datepicker'
import {
    DateRange,
    defaultEventTime,
    defaultFindingsTime,
    useURLParam,
    useUrlDateRangeState,
} from '../../../../../utilities/urlstate'
import { renderDateText } from '../../../../../components/Layout/Header/DatePicker'

interface IFilters {
    onApply: (obj: {
        connector: SourceType
        conformanceStatus:
            | PlatformEnginePkgComplianceApiConformanceStatus[]
            | undefined
        severity: TypesFindingSeverity[] | undefined
        connectionID: string[] | undefined
        controlID: string[] | undefined
        benchmarkID: string[] | undefined
        resourceTypeID: string[] | undefined
        lifecycle: boolean[] | undefined
        activeTimeRange: DateRange | undefined
        eventTimeRange: DateRange | undefined
    }) => void
    type: 'findings' | 'resources' | 'controls' | 'accounts' | 'events'
    id : string;
}

export default function Filter({ onApply, type ,id}: IFilters) {
    const { ws } = useParams()
    const defConnector = SourceType.Nil
    const [connector, setConnector] = useState<SourceType>(defConnector)

    const defConformanceStatus = [
        PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
        PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusPassed,
    ]
    const [conformanceStatus, setConformanceStatus] = useState<
        | PlatformEnginePkgComplianceApiConformanceStatus[]
        | undefined
    >([
        PlatformEnginePkgComplianceApiConformanceStatus.ConformanceStatusFailed,
    ])

    const defLifecycle = [true]
    const [lifecycle, setLifecycle] = useState<boolean[]>(defLifecycle)

    const defSeverity = [
        TypesFindingSeverity.FindingSeverityCritical,
        TypesFindingSeverity.FindingSeverityHigh,
        TypesFindingSeverity.FindingSeverityMedium,
        TypesFindingSeverity.FindingSeverityLow,
        TypesFindingSeverity.FindingSeverityNone,
    ]
    const [severity, setSeverity] = useState<
        TypesFindingSeverity[] | undefined
    >(defSeverity)
    const [severityCon, setSeverityCon] = useState('is')

    const [connectionID, setConnectionID] = useState<string[] | undefined>([])
    const [connectionCon, setConnectionCon] = useState('is')
    const [controlID, setControlID] = useState<string[] | undefined>([])
    const [controlCon, setControlCon] = useState('is')
    const [benchmarkID, setBenchmarkID] = useState<string[] | undefined>([id])
    const [benchmarkCon, setBenchmarkCon] = useState('is')
    const [resourceTypeID, setResourceTypeID] = useState<string[] | undefined>(
        []
    )
    const [resourceCon, setResourceCon] = useState('is')
    const [dateCon, setDateCon] = useState('isBetween')
    const [eventDateCon, setEventDateCon] = useState('isBetween')
    const { value: activeTimeRange, setValue: setActiveTimeRange } =
        useUrlDateRangeState(defaultFindingsTime(ws || ''))
    const { value: eventTimeRange, setValue: setEventTimeRange } =
        useUrlDateRangeState(
            defaultEventTime(ws || ''),
            'eventStartDate',
            'eventEndDate'
        )
    const [selectedFilters, setSelectedFilters] = useState<string[]>([
        'conformance_status',
        'provider',
        'lifecycle',
        'severity',
        'connection',
        'control',
        'benchmark',
        'resource',
        'date',
        'eventDate',
    ])

    useEffect(() => {
        onApply({
            connector,
            conformanceStatus,
            severity,
            connectionID,
            controlID,
            benchmarkID,
            resourceTypeID,
            lifecycle,
            activeTimeRange: selectedFilters.includes('date')
                ? activeTimeRange
                : undefined,
            eventTimeRange: selectedFilters.includes('eventDate')
                ? eventTimeRange
                : undefined,
        })
    }, [
        connector,
        conformanceStatus,
        severity,
        connectionID,
        controlID,
        // benchmarkID,
        resourceTypeID,
        lifecycle,
        activeTimeRange,
        eventTimeRange,
    ])

    const { response: filters } = useComplianceApiV1FindingsFiltersCreate({})

    const filterOptions = [
        // {
        //     id: 'conformance_status',
        //     name: 'Conformance Status',
        //     icon: CheckCircleIcon,
        //     component: (
        //         <ConformanceStatus
        //             value={conformanceStatus}
        //             defaultValue={defConformanceStatus}
        //             onChange={(c) => setConformanceStatus(c)}
        //         />
        //     ),
        //     conditions: ['is'],
        //     setCondition: (c: string) => console.log(c),
        //     value: conformanceStatus,
        //     defaultValue: defConformanceStatus,
        //     onDelete: undefined,
        //     types: ['findings', 'resources', 'events'],
        // },
        // {
        //     id: 'provider',
        //     name: 'Connector',
        //     icon: CloudConnect,
        //     component: (
        //         <Provider
        //             value={connector}
        //             defaultValue={defConnector}
        //             onChange={(p) => setConnector(p)}
        //         />
        //     ),
        //     conditions: ['is'],
        //     setCondition: (c: string) => undefined,
        //     value: [connector],
        //     defaultValue: [defConnector],
        //     onDelete: () => setConnector(defConnector),
        //     types: ['findings', 'resources', 'events', 'controls', 'accounts'],
        // },
        // {
        //     id: 'lifecycle',
        //     name: 'Lifecycle',
        //     icon: Lifecycle,
        //     component: (
        //         <FindingLifecycle
        //             value={lifecycle}
        //             defaultValue={defLifecycle}
        //             onChange={(l) => setLifecycle(l)}
        //         />
        //     ),
        //     conditions: ['is'],
        //     setCondition: (c: string) => console.log(c),
        //     value: lifecycle,
        //     defaultValue: defLifecycle,
        //     onDelete: () => setLifecycle(defLifecycle),
        //     types: ['findings', 'resources', 'events'],
        // },
        // {
        //     id: 'severity',
        //     name: 'Severity',
        //     icon: SeverityIcon,
        //     component: (
        //         <Severity
        //             value={severity}
        //             defaultValue={defSeverity}
        //             condition={severityCon}
        //             onChange={(s) => setSeverity(s)}
        //         />
        //     ),
        //     conditions: ['is', 'isNot'],
        //     setCondition: (c: string) => setSeverityCon(c),
        //     value: severity,
        //     defaultValue: defSeverity,
        //     onDelete: () => setSeverity(defSeverity),
        //     types: ['findings', 'resources', 'events'],
        // },
        {
            id: 'connection',
            name: 'Cloud Account',
            icon: Id,
            component: (
                <Others
                    value={connectionID}
                    defaultValue={[]}
                    data={filters}
                    condition={connectionCon}
                    type="integrationID"
                    onChange={(o) => setConnectionID(o)}
                />
            ),
            conditions: ['is', 'isNot'],
            setCondition: (c: string) => setConnectionCon(c),
            value: connectionID,
            defaultValue: [],
            onDelete: () => setConnectionID([]),
            data: filters?.integrationID,
            types: ['findings', 'resources', 'events', 'controls', 'accounts'],
        },
        // {
        //     id: 'control',
        //     name: 'Control',
        //     icon: Control,
        //     component: (
        //         <Others
        //             value={controlID}
        //             defaultValue={[]}
        //             data={filters}
        //             condition={controlCon}
        //             type="controlID"
        //             onChange={(o) => setControlID(o)}
        //         />
        //     ),
        //     conditions: ['is', 'isNot'],
        //     setCondition: (c: string) => setControlCon(c),
        //     value: controlID,
        //     defaultValue: [],
        //     onDelete: () => setControlID([]),
        //     data: filters?.controlID,
        //     types: ['findings', 'resources', 'events'],
        // },
        // {
        //     id: 'benchmark',
        //     name: 'Benchmark',
        //     icon: Compliance,
        //     component: (
        //         <Others
        //             value={benchmarkID}
        //             defaultValue={[]}
        //             data={filters}
        //             condition={benchmarkCon}
        //             type="benchmarkID"
        //             onChange={(o) => setBenchmarkID(o)}
        //         />
        //     ),
        //     conditions: ['is'],
        //     setCondition: (c: string) => setBenchmarkCon(c),
        //     value: benchmarkID,
        //     defaultValue: [],
        //     onDelete: () => setBenchmarkID([]),
        //     data: filters?.benchmarkID,
        //     types: ['findings', 'resources', 'events', 'controls', 'accounts'],
        // },
        // {
        //     id: 'resource',
        //     name: 'Resource Type',
        //     icon: Resources,
        //     component: (
        //         <Others
        //             value={resourceTypeID}
        //             defaultValue={[]}
        //             data={filters}
        //             condition={resourceCon}
        //             type="resourceTypeID"
        //             onChange={(o) => setResourceTypeID(o)}
        //         />
        //     ),
        //     conditions: ['is', 'isNot'],
        //     setCondition: (c: string) => setResourceCon(c),
        //     value: resourceTypeID,
        //     defaultValue: [],
        //     onDelete: () => setResourceTypeID([]),
        //     data: filters?.resourceTypeID,
        //     types: ['findings', 'resources', 'events'],
        // },
        // {
        //     id: 'date',
        //     name: type === 'events' ? 'Audit Period' : 'Last Evaluated',
        //     icon: CalendarIcon,
        //     component: (
        //         <Datepicker
        //             condition={dateCon}
        //             activeTimeRange={activeTimeRange}
        //             setActiveTimeRange={(v) => setActiveTimeRange(v)}
        //         />
        //     ),
        //     conditions: ['isBetween', 'isRelative'],
        //     setCondition: (c: string) => setDateCon(c),
        //     value: activeTimeRange,
        //     defaultValue: { start: dayjs.utc(), end: dayjs.utc() },
        //     onDelete: () =>
        //         setActiveTimeRange({ start: dayjs.utc(), end: dayjs.utc() }),
        //     types: ['findings', 'events'],
        // },
        // {
        //     id: 'eventDate',
        //     name: 'Last Event',
        //     icon: CalendarIcon,
        //     component: (
        //         <Datepicker
        //             condition={eventDateCon}
        //             activeTimeRange={eventTimeRange}
        //             setActiveTimeRange={(v) => setEventTimeRange(v)}
        //         />
        //     ),
        //     conditions: ['isBetween', 'isRelative'],
        //     setCondition: (c: string) => setEventDateCon(c),
        //     value: eventTimeRange,
        //     defaultValue: { start: dayjs.utc(), end: dayjs.utc() },
        //     onDelete: () =>
        //         setEventTimeRange({ start: dayjs.utc(), end: dayjs.utc() }),
        //     types: ['findings'],
        // },
    ]

    const renderFilters = selectedFilters.map((sf) => {
        const f = filterOptions.find((o) => o.id === sf)
        return (
            f?.types.includes(type) && (
                <Popover className="relative border-0" >
                    <Popover.Button
                        id={f?.id}
                        disabled={true}
                        className={`border ${
                            f?.id !== 'date' &&
                            f?.id !== 'eventDate' &&
                            compareArrays(
                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                // @ts-ignore
                                f?.value?.sort(),
                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                // @ts-ignore
                                f?.defaultValue?.sort()
                            )
                                ? 'border-gray-200 bg-white'
                                : 'border-openg-500 text-openg-500 bg-openg-50'
                        } py-1.5 px-2 rounded-md`}
                    >
                        <Flex className="w-fit">
                            <Icon
                                icon={f?.icon || CloudIcon}
                                className="w-3 p-0 mr-3 text-inherit"
                            />
                            <Text className="text-inherit whitespace-nowrap">
                                {`${f?.name}${
                                    // eslint-disable-next-line no-nested-ternary
                                    f?.id === 'date' || f?.id === 'eventDate'
                                        ? ` ${renderDateText(
                                              f?.id === 'date'
                                                  ? activeTimeRange.start
                                                  : eventTimeRange.start,
                                              f?.id === 'date'
                                                  ? activeTimeRange.end
                                                  : eventTimeRange.end
                                          )}`
                                        : compareArrays(
                                              // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                              // @ts-ignore
                                              f?.value?.sort(),
                                              // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                              // @ts-ignore
                                              f?.defaultValue?.sort()
                                          )
                                        ? ''
                                        : `${
                                              // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                              // @ts-ignore
                                              f?.value && f.value.length < 2
                                                  ? `: ${
                                                        f.data
                                                            ? f.data.find(
                                                                  (d) =>
                                                                      d.key ===
                                                                      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                                                      // @ts-ignore
                                                                      f.value[0]
                                                              )?.displayName
                                                            : f.value
                                                    }`
                                                  : // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                                    // @ts-ignore
                                                    ` (${f?.value?.length})`
                                          }`
                                }`}
                            </Text>
                            <ChevronDownIcon className="ml-1 w-3 text-inherit" />
                        </Flex>
                    </Popover.Button>
                    <Transition
                        as={Fragment}
                        enter="transition ease-out duration-200"
                        enterFrom="opacity-0 translate-y-1"
                        enterTo="opacity-100 translate-y-0"
                        leave="transition ease-in duration-150"
                        leaveFrom="opacity-100 translate-y-0"
                        leaveTo="opacity-0 translate-y-1"
                    >
                        <Popover.Panel
                            static
                            className="absolute z-50 top-full left-0"
                        >
                            <Card className="mt-2 p-4 min-w-[256px] w-fit">
                                <Flex className="mb-3">
                                    <Flex className="w-fit gap-1.5">
                                        <Text className="font-semibold">
                                            {f?.name}
                                        </Text>
                                        <ConditionDropdown
                                            onChange={(c) =>
                                                f?.setCondition
                                                    ? f?.setCondition(c)
                                                    : undefined
                                            }
                                            conditions={f?.conditions}
                                            isDate={
                                                f?.id === 'date' ||
                                                f?.id === 'eventDate'
                                            }
                                        />
                                    </Flex>
                                    {f?.onDelete && (
                                        <div className="group relative">
                                            <TrashIcon
                                                onClick={() => {
                                                    f?.onDelete()
                                                    setSelectedFilters(
                                                        (prevState) => {
                                                            return prevState.filter(
                                                                (s) =>
                                                                    s !== f?.id
                                                            )
                                                        }
                                                    )
                                                }}
                                                className="w-4 cursor-pointer hover:text-openg-500"
                                            />
                                            <Card className="absolute w-fit z-40 -top-2 left-full ml-2 scale-0 transition-all p-2 group-hover:scale-100">
                                                <Text className="whitespace-nowrap">
                                                    Remove filter
                                                </Text>
                                            </Card>
                                        </div>
                                    )}
                                </Flex>
                                {f?.component}
                            </Card>
                        </Popover.Panel>
                    </Transition>
                </Popover>
            )
        )
    })

    return (
        <Flex justifyContent="start" className=" gap-3 flex-wrap z-10 w-fit">
            {renderFilters}
            {filterOptions.filter((f) => !selectedFilters.includes(f.id))
                .length > 0 && (
                <Flex className="w-fit pl-3 border-l border-l-gray-200 h-full">
                    <Popover className="relative border-0">
                        <Popover.Button>
                            <Button
                                variant="light"
                                icon={PlusIcon}
                                className="pt-1"
                            >
                                Add Filter
                            </Button>
                        </Popover.Button>
                        <Transition
                            as={Fragment}
                            enter="transition ease-out duration-200"
                            enterFrom="opacity-0 translate-y-1"
                            enterTo="opacity-100 translate-y-0"
                            leave="transition ease-in duration-150"
                            leaveFrom="opacity-100 translate-y-0"
                            leaveTo="opacity-0 translate-y-1"
                        >
                            <Popover.Panel className="absolute z-50 top-full left-0">
                                <Card className="mt-2 p-4 w-64">
                                    <Flex
                                        flexDirection="col"
                                        justifyContent="start"
                                        alignItems="start"
                                        className="gap-1.5 max-w-full"
                                    >
                                        {filterOptions
                                            .filter(
                                                (f) =>
                                                    !selectedFilters.includes(
                                                        f.id
                                                    ) && f.types.includes(type)
                                            )
                                            .map((f) => (
                                                <Button
                                                    icon={f.icon}
                                                    color="slate"
                                                    variant="light"
                                                    className="w-full pl-1 flex justify-start"
                                                    onClick={() => {
                                                        setSelectedFilters([
                                                            ...selectedFilters,
                                                            f.id,
                                                        ])
                                                    }}
                                                >
                                                    {f.name}
                                                </Button>
                                            ))}
                                    </Flex>
                                </Card>
                            </Popover.Panel>
                        </Transition>
                    </Popover>
                </Flex>
            )}
        </Flex>
    )
}
