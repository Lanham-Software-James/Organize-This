import { PUBLIC_API_URL } from "$env/static/public";
import type { PageLoad } from "./$types";

export interface getEntityEntity {
    ID: number
    Name: string
    Category: string
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

//@ts-ignore
export const load = (async ({ params }) => {

    let message: string = ""
    let entity: getEntityData = {
        Entity: {
            ID: 0,
            Category: '',
            Name: '',
        },
        Parent: {
            ParentID: 0,
            ParentCategory: '',
        }
    }
    let parentName : string = ""

    const promises = [
        fetch(`${PUBLIC_API_URL}api/v1/entity/${params.category}/${params.id}`),
        fetch(`${PUBLIC_API_URL}api/v1/entity/${params.category}/${params.id}`), //Replace this with get children
    ]

    try {
        const [getEntityResponse, getChildrenResponse] = await Promise.all(promises)
        const [getEntityData, getChildrenData] = await Promise.all([getEntityResponse.json(), getChildrenResponse.json()]);

        const getParentResponse = await fetch(`${PUBLIC_API_URL}api/v1/entity/${getEntityData.data.Parent.ParentCategory}/${getEntityData.data.Parent.ParentID}`);
        const getParentData = await getParentResponse.json();

        if (getEntityData.message == getChildrenData.message &&
            getChildrenData.message == getParentData.message &&
            getParentData.message == "success") {

            entity = getEntityData.data
            entity.Entity.Category = params.category
            parentName = getParentData.data.Entity.Name
            message = "success"
        }
    } catch (error) {
        console.error(error)
    }

    return {
        message: message,
        entity: entity,
        parentName: parentName,
    }
}) satisfies PageLoad;
