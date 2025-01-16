import ReactJson from '@microlink/react-json-view'
import { Card } from '@tremor/react'
import 'ace-builds/css/ace.css'
import 'ace-builds/css/theme/cloud_editor.css'
import 'ace-builds/css/theme/cloud_editor_dark.css'
import 'ace-builds/css/theme/cloud_editor_dark.css'
import 'ace-builds/css/theme/twilight.css'
import 'ace-builds/css/theme/sqlserver.css'
import 'ace-builds/css/theme/xcode.css'

import { useEffect, useState } from 'react'
import { CodeEditor } from '@cloudscape-design/components'


interface IRenderObjectProps {
    obj: any
}

export function RenderObject({ obj }: IRenderObjectProps) {
    const [ace, setAce] = useState()
    const [preferences, setPreferences] = useState(undefined)

 useEffect(() => {
        async function loadAce() {
            const ace = await import('ace-builds')
            await import('ace-builds/webpack-resolver')
            ace.config.set('useStrictCSP', true)
            // ace.config.setMode('ace/mode/sql')
            // @ts-ignore
            // ace.edit(element, {
            //     mode: 'ace/mode/sql',
            //     selectionStyle: 'text',
            // })

            return ace
        }

        loadAce()
            .then((ace) => {
                // @ts-ignore
                setAce(ace)
            })
            .finally(() => {})
    }, [])

    return (
        /* <List>
            {Object.keys(obj).length > 0 &&
                Object.keys(obj).map((key) => {
                    if (typeof obj[key] === 'object' && obj[key] !== null) {
                        if (Object.keys(obj[key]).length === 0) {
                            return null
                        }
                        return (
                            <div>
                                {key !== '0' && (
                                    <Title className="mt-6">
                                        {changeKeysToLabel
                                            ? snakeCaseToLabel(key)
                                            : key}
                                    </Title>
                                )}
                                <RenderObject obj={obj[key]} />
                            </div>
                        )
                    }

                    return (
                        <ListItem key={key} className="py-6 flex items-start">
                            <Text>
                                {changeKeysToLabel
                                    ? snakeCaseToLabel(key)
                                    : key}
                            </Text>
                            <Text className="text-gray-800 w-3/5 whitespace-pre-wrap text-end">
                                {String(obj[key])}
                            </Text>
                        </ListItem>
                    )
                })}
        </List> */
        // <Card className="px-1.5 py-3 mb-2">
        //     <ReactJson
        //         src={obj}
        //         style={{
        //             lineBreak: 'anywhere',
        //         }}
        //     />
        // </Card>
        <CodeEditor
        // className='h-full'
            ace={ace}
            language="json"
            value={JSON.stringify(obj,null,'\t')}
            languageLabel="JSON"
            onChange={({ detail }) => {
                // setSavedQuery('')
                // setCode(detail.value)
            }}
            editorContentHeight={500}
            preferences={preferences}
            onPreferencesChange={(e) =>
                // @ts-ignore
                setPreferences(e.detail)
            }
            loading={false}
            themes={{
                light: ['xcode','cloud_editor', 'sqlserver'],
                dark: ['cloud_editor_dark', 'twilight'],
                // @ts-ignore
            }}
        />
    )
}
