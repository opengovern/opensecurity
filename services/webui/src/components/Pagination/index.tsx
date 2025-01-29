import { Box, Button, Pagination, SpaceBetween } from '@cloudscape-design/components'
import { Flex } from '@tremor/react'

interface ISpinner {
    currentPageIndex: number
    pagesCount: number
    onChange: any
}

export default function CustomPagination({
    currentPageIndex,
    pagesCount,
    onChange,
}: ISpinner) {
    return (
        <>
            {window.innerWidth > 640 ? (
                <>
                    {' '}
                    <div className="w-full flex justify-end">
                        <Pagination
                            currentPageIndex={currentPageIndex}
                            pagesCount={pagesCount}
                            onChange={onChange}
                        />
                    </div>
                </>
            ) : (
                <>
                    {' '}
                    <Flex className="w-full justify-start items-center flex-row gap-2">
                        <Button
                            iconName="angle-left"
                            variant="icon"
                            disabled={currentPageIndex === 1}
                            onClick={() =>
                                onChange({
                                    detail: {
                                        currentPageIndex: currentPageIndex - 1,
                                    },
                                })
                            }
                        >
                            Previous
                        </Button>
                        <div style={{ marginTop: '3px' }}>
                            <Box color="text-body-secondary">
                                {pagesCount === 0
                                    ? ''
                                    : `${currentPageIndex} of ${pagesCount}`}
                            </Box>
                        </div>
                        <Button
                            iconName="angle-right"
                            variant="icon"
                            disabled={
                                currentPageIndex === pagesCount ||
                                pagesCount === 0
                            }
                            onClick={() =>
                                onChange({
                                    detail: {
                                        currentPageIndex: currentPageIndex + 1,
                                    },
                                })
                            }
                        >
                            Next
                        </Button>
                    </Flex>
                </>
            )}
        </>
    )
}
