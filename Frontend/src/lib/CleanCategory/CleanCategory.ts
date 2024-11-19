export function cleanCategory(category: string): string {
    var cleanedCategory = category;
    if (cleanedCategory == 'shelving_unit') {
        cleanedCategory = 'shelving unit';
    }

    return cleanedCategory;
}
