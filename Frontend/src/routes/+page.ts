import { PUBLIC_API_URL } from '$env/static/public';

export interface GetEntitiesResponse {
    TotalCount: number,
    Entities: GetEntitiesData[]
}

export interface GetEntitiesData {
    ID: number,
    Name: string,
    Category: string,
    Parent: Parent[],
    Notes: string,
    Address?: string,
}

export interface Parent {
    ID: number,
    Name: string,
    Category: string,
}

export const _getEntities = async (offset: number, limit: number, search: string, filters: {[key: string]: boolean}): Promise<[GetEntitiesData[], number]> => {
    let entities: GetEntitiesData[] = []
    let size: number = 0
    let filterQuery: string = ''

    for(let key in filters) {
        if(filters[key] && filterQuery.length > 0) {
            filterQuery += ',' + key
        } else if(filters[key]) {
            filterQuery += key
        }
    }

    try{
        const response = await fetch(`${PUBLIC_API_URL}api/v1/entities?offset=${offset}&limit=${limit}&search=${search}&filter=${filterQuery}`);
        const data = await response.json()

        if (data.message == "success") {
            if(data.data.Entities != null) {
                entities = data.data.Entities
            }

            size = +data.data.TotalCount
        }

    } catch(error) {
        console.error(error);
    }

    // Give empty default row for table to iterate over
    if(entities.length === 0) {
        entities.push({
            ID: 0,
            Name: " ",
            Category: " ",
            Parent: [],
            Notes: " "
        })
    }

    return [entities, size]
}
