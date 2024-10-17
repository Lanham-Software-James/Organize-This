import { PUBLIC_API_URL } from '$env/static/public';

export interface parentData {
    ID: number
    Name: string
    Category: string
}

export interface getEntityEntity {
    ID: number
    Name: string
    Notes?: string
}

export interface getEntityParent {
    ParentID: number
    ParentCategory: string
}

export interface getEntityData {
    Entity: getEntityEntity
    Parent: getEntityParent
    Address?: string
}

export const createEntity = async (formData: { category: string; address: string; name: string; notes: string; parent: string;}): Promise<[string, number]> => {
    let message: string = ""
    let id: number = 0

    const parents = formData.parent.split('-')


    const response = await fetch(`${PUBLIC_API_URL}api/v1/entity`, {
        method: 'POST',
        body: JSON.stringify({
            address: formData.address,
            category: formData.category,
            name: formData.name,
            notes: formData.notes,
            parentID: parents[0],
            parentCategory: parents[1],
        })
    });

    const data = await response.json()

    message = data.message
    if (message == "success") {
        id = data.data
    }

    return [message, id];
}

export const editEntity = async (formData: { id: number, category: string; address: string; name: string; notes: string; parent: string;}): Promise<[string, getEntityData]> => {
    let message: string = ""
    let entity: getEntityData = {
        Entity: {
            ID: 0,
            Name: ''
        },
        Parent: {
            ParentID: 0,
            ParentCategory: ''
        }
    }

    const parents = formData.parent.split('-')


    const response = await fetch(`${PUBLIC_API_URL}api/v1/entity`, {
        method: 'PUT',
        body: JSON.stringify({
            id: formData.id,
            address: formData.address,
            category: formData.category,
            name: formData.name,
            notes: formData.notes,
            parentID: parents[0],
            parentCategory: parents[1],
        })
    });

    const data = await response.json()

    message = data.message
    if (message == "success") {
        entity = data.data
    }

    return [message, entity];
}

export const getEntity = async(id: number, category: string): Promise<[string, getEntityData]> => {
    let message: string = ""
    let entity: getEntityData = {
        Entity: {
            ID: 0,
            Name: ''
        },
        Parent: {
            ParentID: 0,
            ParentCategory: ''
        }
    }

    const response = await fetch(`${PUBLIC_API_URL}api/v1/entity/${category}/${id}`);

    try {
       const data = await response.json()

       message = data.message
        if (message == "success") {
            entity = data.data
        }
    } catch (error) {
        console.error(error)
    }




    return [message, entity]
}

export const getParents = async(category: string): Promise<[string, parentData[]]> => {
    let message: string = ""
    let parents: parentData[] = []

    const response = await fetch(`${PUBLIC_API_URL}api/v1/parents/${category}`);

    const data = await response.json()

    message = data.message
    if (message == "success") {
        parents = data.data
    }

    return [message, parents]
}
