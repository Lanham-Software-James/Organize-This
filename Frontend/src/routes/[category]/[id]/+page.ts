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
    let children: getEntityEntity[] = []

    const promises = [
        fetch(`/api/v1/entity/${params.category}/${params.id}`),
        fetch(`/api/v1/children/${params.category}/${params.id}`),
    ]

    try {
        const [getEntityResponse, getChildrenResponse] = await Promise.all(promises)
        const [getEntityData, getChildrenData] = await Promise.all([getEntityResponse.json(), getChildrenResponse.json()]);

        let getParentResponse: Response
        let getParentData: any
        let success: boolean

        if(params.category != "building"){
            getParentResponse = await fetch(`/api/v1/entity/${getEntityData.data.Parent.ParentCategory}/${getEntityData.data.Parent.ParentID}`);
            getParentData = await getParentResponse.json();
            success = getEntityData.message == getChildrenData.message &&
                getChildrenData.message == getParentData.message &&
                getParentData.message == "success"
        } else {
            getParentData = undefined
            success = getEntityData.message == getChildrenData.message &&
                getChildrenData.message == "success"
        }

        if (success) {

            entity = getEntityData.data
            entity.Entity.Category = params.category

            if(getParentData != undefined) {
                parentName = getParentData.data.Entity.Name
            }

            if(getChildrenData.data != null) {
               children = getChildrenData.data
            }

            message = "success"
        }
    } catch (error) {
        console.error(error)
    }

    return {
        message: message,
        entity: entity,
        parentName: parentName,
        children: children,
    }
}) satisfies PageLoad;
