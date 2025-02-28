import { ReactNode, useEffect, useRef, useState } from 'react'
import { useSpring, a, animated } from '@react-spring/web'
import useMeasure from 'react-use-measure'
import { Flex } from '@tremor/react'
import styled from 'styled-components'
import { ChevronRightIcon, ChevronUpIcon } from '@heroicons/react/24/outline'

function usePrevious<T>(value: T) {
    const ref = useRef<T>()
    // eslint-disable-next-line no-return-assign,no-void
    useEffect(() => void (ref.current = value), [value])
    return ref.current
}

interface IAnimatedAccordion {
    children: ReactNode
    header: ReactNode
    defaultOpen?: boolean
}

const Content = styled(animated.div)`
    will-change: transform, opacity, height;
    overflow: hidden;
`

export default function AnimatedAccordion({
    children,
    header,
    defaultOpen = false,
}: IAnimatedAccordion) {
    const [isOpen, setOpen] = useState(defaultOpen)
    const previous = usePrevious(isOpen)

    useEffect(() => {
        setOpen(defaultOpen)
    }, [defaultOpen])

    const [ref, { height: viewHeight }] = useMeasure()
    const { height, opacity, y } = useSpring({
        from: { height: 0, opacity: 0, y: 0 },
        to: {
            height: isOpen ? viewHeight + 4 : 0,
            opacity: isOpen ? 1 : 0,
            y: isOpen ? 0 : 20,
        },
    })

    return (
        <div className="w-full relative overflow-x-hidden">
            <Flex
                className="cursor-pointer relative"
                onClick={() => setOpen(!isOpen)}
            >
                <div className="absolute">
                    {isOpen ? (
                        <ChevronUpIcon className="w-4 text-gray-400" />
                    ) : (
                        <ChevronRightIcon className="w-4 text-gray-400" />
                    )}
                </div>
                {header}
            </Flex>
            <Content
                style={{
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-ignore
                    opacity,
                    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                    // @ts-ignore
                    height: isOpen && previous === isOpen ? 'auto' : height,
                }}
            >
                <a.div ref={ref} style={{ y }}>
                    {children}
                </a.div>
            </Content>
        </div>
    )
}
