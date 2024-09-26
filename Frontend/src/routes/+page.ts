import { PUBLIC_API_URL } from '$env/static/public';

export interface GetEntitiesResponse {
    TotalCount: number,
    Entities: GetEntitiesData[]
}

export interface GetEntitiesData {
    ID: number,
    Name: string,
    Category: string,
    Location: string,
    Notes: string
}

export const _getEntities = async (offset: number, limit: number): Promise<[GetEntitiesData[], number]> => {
    let entities: GetEntitiesData[] = []
    let size: number = 0

    const response = await fetch(`${PUBLIC_API_URL}v1/entities?offset=${offset}&limit=${limit}`);

    const data = await response.json()

    if (data.message == "success") {
        entities = data.data.Entities
        size = +data.data.TotalCount
    }

    return [entities, size]
}
