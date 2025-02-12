// @ts-nocheck
import React, { useEffect, useRef } from 'react'
import MonacoEditor, { loader } from '@monaco-editor/react'

const SQLEditor = ({ value, onChange, tables }) => {
    // Ref to track whether completion provider has been registered
    const providerRegistered = useRef(false)

    useEffect(() => {
        console.log(tables)
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

    return (
        <MonacoEditor
            language="sql"
            theme="vs"
            loading="Loading..."
            value={value}
            onChange={onChange}
            options={{
                automaticLayout: true,
                suggestOnTriggerCharacters: true,
                quickSuggestions: true,
                lineNumbers: 'on',
                renderLineHighlight: 'all',
                
                
            }}
        />
    )
}

export default SQLEditor
