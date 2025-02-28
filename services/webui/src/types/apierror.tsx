import { ExclamationTriangleIcon } from '@heroicons/react/24/outline'
import { Button, Flex, Text, Title } from '@tremor/react'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const getErrorMessage = (err: any) => {
    if (err) {
        try {
            return `Error: ${String(err.response.data.message)}`
        } catch {
            return String(err)
        }
    }
    return ''
}

export const toErrorMessage = (...errs: any[]) => {
    const msgs = errs
        .map((err) => {
            const msg = getErrorMessage(err)
            if (msg.length > 0) {
                return msg
            }
            return undefined
        })
        .filter((v) => v !== undefined)

    return msgs.length > 0 ? msgs.join(',') : undefined
}

export const errorHandlingWithErrorMessage = (
    onRefresh: (() => void) | undefined,
    msg: string | undefined
) => {
    return msg === undefined || msg === '' ? null : (
        <Flex
            flexDirection="col"
            justifyContent="between"
            className="absolute top-0 w-full left-0 h-full backdrop-blur"
        >
            <Flex
                flexDirection="col"
                justifyContent="center"
                alignItems="center"
            >
                <Flex
                    flexDirection="row"
                    justifyContent="center"
                    className="mt-6"
                >
                    <ExclamationTriangleIcon className="text-red-600 w-6 mr-1 mt-1" />
                    <Title className="text-red-600">
                        Failed to load component
                    </Title>
                </Flex>
                <Text className="mt-2 text-red-600">{msg}</Text>
            </Flex>
            <Button
                variant="secondary"
                className="mb-6"
                color="rose"
                onClick={() => {
                    if (onRefresh) {
                        onRefresh()
                    }
                }}
            >
                Try Again
            </Button>
        </Flex>
    )
}

export const errorHandling = (onRefresh: () => void, ...errs: any[]) => {
    const msg = toErrorMessage(errs)
    return errorHandlingWithErrorMessage(onRefresh, msg)
}
