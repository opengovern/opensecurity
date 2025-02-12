// @ts-nocheck
import React, { useEffect, useRef } from 'react'
import MonacoEditor, { loader } from '@monaco-editor/react'

const Map = {
    aws_: 'aws_cloud_account',
    azure_: 'azure_subscription',
    entraid_: 'entraid_directory',
    github_: 'github_account',
    digitalocean_: 'digitalocean_team',
    cloudflare_: 'cloudflare_account',
    openai_: 'openai_integration',
    linode_: 'linode_account',
    cohereai_: 'cohereai_project',
    google_: 'google_workspace_account',
    oci_: 'oci_repository',
    render_: 'render_account',
    doppler_: 'doppler_account',
    jira_: 'jira_cloud',
    tailscale_: 'tailscale_account',
    heroku_: 'heroku_account',
    fly_: 'fly_account',
    semgrep_: 'semgrep_account',
    kubernetes_: 'kubernetes',
}

const handleName = (text: string) => {
    const table = text.trim().split(' ').pop()
    for (const key in Map) {
        if (table == key) {
            return Map[key]
        }
    }
    return undefined
}

const SQLEditor = ({ value, onChange, tables,tableFetch,run }) => {
    // Ref to track whether completion provider has been registered
    const providerRegistered = useRef(false)
      const editorRef = useRef(null)
    useEffect(() => {
        loader.init().then((monaco) => {
            if (providerRegistered.current) return // Prevent duplicate registration
            providerRegistered.current = true // Mark as registered

            // Define SQL keywords and commands
            const sqlKeywords = [
                'SELECT',
                'FROM',
                'WHERE',
                'INSERT INTO',
                'UPDATE',
                'DELETE',
                'JOIN',
                'INNER JOIN',
                'LEFT JOIN',
                'RIGHT JOIN',
                'ON',
                'GROUP BY',
                'HAVING',
                'ORDER BY',
                'LIMIT',
                'OFFSET',
                'DISTINCT',
                'CREATE TABLE',
                'DROP TABLE',
                'ALTER TABLE',
                'ADD',
                'AND',
                'OR',
                'NOT',
                'IN',
                'EXISTS',
                'BETWEEN',
                'LIKE',
                'IS NULL',
                'IS NOT NULL',
                'UNION',
                'UNION ALL',
                'CASE',
                'WHEN',
                'THEN',
                'ELSE',
                'END',
            ]

            // Registering SQL completion provider
            monaco.languages.registerCompletionItemProvider('sql', {
                triggerCharacters: [' ', '.', ','],
                provideCompletionItems: (model, position) => {
                    const textUntilPosition = model.getValueInRange(
                        new monaco.Range(
                            1,
                            1,
                            position.lineNumber,
                            position.column
                        )
                    )
                    const table_name = handleName(textUntilPosition)
                    if (table_name) {
                        tableFetch(table_name)
                    }
                    

                    const suggestions = []

                    // Suggest column names first
                    tables.forEach((table) => {
                        if (textUntilPosition.includes(table.table)) {
                            table.columns.forEach((column) => {
                                if (
                                    !suggestions.some(
                                        (s) =>
                                            s.label === column &&
                                            s.detail ===
                                                `Column in ${table.table}`
                                    )
                                ) {
                                    suggestions.push({
                                        label: column,
                                        kind: monaco.languages
                                            .CompletionItemKind.Field,
                                        insertText: column,
                                        detail: `Column in ${table.table}`,
                                    })
                                }
                            })
                        }
                    })

                    // Suggest table names next
                    tables.forEach((table) => {
                        if (
                            !suggestions.some(
                                (s) =>
                                    s.label === table.table &&
                                    s.detail === 'Table'
                            )
                        ) {
                            suggestions.push({
                                label: table.table,
                                kind: monaco.languages.CompletionItemKind
                                    .Keyword,
                                insertText: table.table,
                                detail: 'Table',
                            })
                        }
                    })

                    // Add SQL keywords last
                    sqlKeywords.forEach((keyword) => {
                        if (
                            !suggestions.some(
                                (s) =>
                                    s.label === keyword &&
                                    s.detail === 'SQL Keyword'
                            )
                        ) {
                            suggestions.push({
                                label: keyword,
                                kind: monaco.languages.CompletionItemKind
                                    .Keyword,
                                insertText: keyword,
                                detail: 'SQL Keyword',
                            })
                        }
                    })

                    return { suggestions }
                },
            })
        })
    }, [tables]) // Effect runs only when `tables` changes
     const handleShiftEnter = () => {
        if (editorRef.current) {
            const currentValue = editorRef.current.getValue() // Fetch the latest value
            run(currentValue) // Pass the value to the `run` function
        }
     }
      const handleEditorDidMount = (editor, monaco) => {
          editorRef.current = editor

          // Register Shift+Enter command
          editor.addCommand(
              monaco.KeyMod.Shift | monaco.KeyCode.Enter, // Detect Shift+Enter
              () => {
                  // Add your custom behavior here
                  handleShiftEnter()
              }
          )
      }

    return (
        <MonacoEditor
            language="sql"
            theme="vs"
            loading="Loading..."
            onMount={handleEditorDidMount}
            value={value}
            onChange={onChange}
            options={{
                automaticLayout: true,
                suggestOnTriggerCharacters: true,
                quickSuggestions: true,
                lineNumbers: 'on',
                renderLineHighlight: 'all',
                fontFamily: `var(--font-family-base-dnvic8, "Open Sans", "Helvetica Neue", Roboto, Arial, sans-serif)`,
                fontSize: 16,
                fontWeight: 400,

                // fontLigatures: true,
            }}
            
        />
    )
}

export default SQLEditor
