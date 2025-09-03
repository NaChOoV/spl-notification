export function extractUnique<T, K>(items: T[], keySelector: (item: T) => K): K[] {
    return Array.from(new Set(items.map(keySelector)));
}
