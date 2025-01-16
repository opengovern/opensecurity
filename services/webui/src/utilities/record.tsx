export function recordToMap<K extends string | number | symbol, T>(
    record: Record<K, T> | undefined
): Map<K, T> {
    if (record === undefined) {
        return new Map<K, T>()
    }
    return new Map<K, T>(Object.entries(record) as [K, T][])
}

export function mapToRecord<K extends string | number | symbol, T>(
    map: Map<K, T>
): Record<K, T> {
    const record = {} as Record<K, T>
    map.forEach((value, key) => {
        record[key] = value
    })
    return record
}
