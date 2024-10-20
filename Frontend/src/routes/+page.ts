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
    Notes: string
}

export interface Parent {
    ID: number,
    Name: string,
    Category: string,
}

export const _getEntities = async (offset: number, limit: number): Promise<[GetEntitiesData[], number]> => {
    let entities: GetEntitiesData[] = []
    let size: number = 0

    try{
        const response = await fetch(`${PUBLIC_API_URL}api/v1/entities?offset=${offset}&limit=${limit}`);
        const data = await response.json()

        if (data.message == "success") {
            entities = data.data.Entities
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
