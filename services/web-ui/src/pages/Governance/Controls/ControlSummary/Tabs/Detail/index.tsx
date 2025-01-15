import { Flex, Text } from '@tremor/react'
import { useState } from 'react'
import MarkdownPreview from '@uiw/react-markdown-preview'
import { PlatformEnginePkgComplianceApiControl } from '../../../../../../api/api'

interface IDetail {
    control: PlatformEnginePkgComplianceApiControl | undefined
}

export default function Detail({ control }: IDetail) {
    const [selectedTab, setSelectedTab] = useState<
        'explanation' | 'nonComplianceCost' | 'usefulExample'
    >('explanation')

    return (
        <Flex alignItems="start" className="mt-6">
            <Flex flexDirection="col" alignItems="start" className="w-56 gap-3">
                {!!control?.explanation && (
                    <button
                        type="button"
                        onClick={() => setSelectedTab('explanation')}
                    >
                        <Text
                            className={`text-gray-500 cursor-pointer ${
                                selectedTab === 'explanation'
                                    ? 'text-openg-500'
                                    : ''
                            }`}
                        >
                            Explanation
                        </Text>
                    </button>
                )}
                {!!control?.nonComplianceCost && (
                    <button
                        type="button"
                        onClick={() => setSelectedTab('nonComplianceCost')}
                    >
                        <Text
                            className={`text-gray-500 cursor-pointer ${
                                selectedTab === 'nonComplianceCost'
                                    ? 'text-openg-500'
                                    : ''
                            }`}
                        >
                            Cost of non-compliance
                        </Text>
                    </button>
                )}
                {!!control?.usefulExample && (
                    <button
                        type="button"
                        onClick={() => setSelectedTab('usefulExample')}
                    >
                        <Text
                            className={`text-gray-500 cursor-pointer ${
                                selectedTab === 'usefulExample'
                                    ? 'text-openg-500'
                                    : ''
                            }`}
                        >
                            Examples of usefulness
                        </Text>
                    </button>
                )}
            </Flex>
            <div
                className="pl-8 border-l border-l-gray-200"
                style={{ width: 'calc(100% - 224px)' }}
            >
                <MarkdownPreview
                    source={control?.[selectedTab]}
                    className="!bg-transparent"
                    wrapperElement={{
                        'data-color-mode': 'light',
                    }}
                    rehypeRewrite={(node, index, parent) => {
                        if (
                            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                            // @ts-ignore
                            node.tagName === 'a' &&
                            parent &&
                            /^h(1|2|3|4|5|6)/.test(
                                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                                // @ts-ignore
                                parent.tagName
                            )
                        ) {
                            // eslint-disable-next-line no-param-reassign
                            parent.children = parent.children.slice(1)
                        }
                    }}
                />
            </div>
        </Flex>
    )
}
