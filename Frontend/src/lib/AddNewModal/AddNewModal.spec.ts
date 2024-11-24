import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createEntity, deleteEntity, editEntity, getEntity, getParents } from './AddNewModal';

// Define a type for the mock fetch function
type FetchMock = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;

function createFetchResponse(data: unknown) {
    return { json: () => new Promise((resolve) => resolve(data)) }
}

describe("Unit Tests for createEntity()", () => {
    beforeEach(() => {
        // Mock the global fetch
        global.fetch = vi.fn() as FetchMock;
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    it("FEUT-1: Build Create Item Request", async () => {
        const formData = {
            category: 'item',
            name: 'Test item',
            address: '',
            notes: 'Test notes',
            parent: '1-container',
        };

        const parentData = formData.parent.split('-')

        const createItemResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createItemResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `/api/v1/entity`,
            {
                method: 'POST',
                headers:{
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes,
                    parentID: parentData[0],
                    parentCategory: parentData[1],
                }),
            }
        );
        expect(message).toEqual(createItemResponse.message)
        expect(id).toEqual(createItemResponse.data)
    });

    it("FEUT-2: Build Create Container Request", async () => {
        const formData = {
            category: 'container',
            name: 'Test container',
            address: '',
            notes: 'Test notes',
            parent: '1-shelf',
        };

        const parentData = formData.parent.split('-')

        const createContainerResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createContainerResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `/api/v1/entity`,
            {
                method: 'POST',
                headers:{
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes,
                    parentID: parentData[0],
                    parentCategory: parentData[1],
                }),
            }
        );
        expect(message).toEqual(createContainerResponse.message)
        expect(id).toEqual(createContainerResponse.data)
    });

    it("FEUT-3: Build Create Shelf Request", async () => {
        const formData = {
            category: 'shelf',
            name: 'Test shelf',
            address: '',
            notes: 'Test notes',
            parent: '1-shelving_unit',
        };

        const parentData = formData.parent.split('-')

        const createShelfResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createShelfResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `/api/v1/entity`,
            {
                method: 'POST',
                headers:{
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes,
                    parentID: parentData[0],
                    parentCategory: parentData[1],
                }),
            }
        );
        expect(message).toEqual(createShelfResponse.message)
        expect(id).toEqual(createShelfResponse.data)
    });

    it("FEUT-4: Build Create Shelving Unit Request", async () => {
        const formData = {
            category: 'shelvingunit',
            name: 'Test shelving unit',
            address: '',
            notes: 'Test notes',
            parent: '1-room',
        };

        const parentData = formData.parent.split('-')

        const createUnitResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createUnitResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `/api/v1/entity`,
            {
                method: 'POST',
                headers:{
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes,
                    parentID: parentData[0],
                    parentCategory: parentData[1],
                }),
            }
        );
        expect(message).toEqual(createUnitResponse.message)
        expect(id).toEqual(createUnitResponse.data)
    });

    it("FEUT-5: Build Create Room Request", async () => {
        const formData = {
            category: 'room',
            name: 'Test room',
            address: '',
            notes: 'Test notes',
            parent: '1-building',
        };

        const parentData = formData.parent.split('-')

        const createRoomResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createRoomResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `/api/v1/entity`,
            {
                method: 'POST',
                headers:{
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes,
                    parentID: parentData[0],
                    parentCategory: parentData[1],
                }),
            }
        );
        expect(message).toEqual(createRoomResponse.message)
        expect(id).toEqual(createRoomResponse.data)
    });

    it("FEUT-6: Build Create Building Request", async () => {
        const formData = {
            name: 'Test building',
            category: 'building',
            address: '888 test road',
            notes: 'Test notes',
            parent: '',
        };

        const parentData = formData.parent.split('-')

        const createBuildingResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createBuildingResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `/api/v1/entity`,
            {
                method: 'POST',
                headers:{
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes,
                    parentID: parentData[0],
                    parentCategory: parentData[1],
                }),
            }
        );
        expect(message).toEqual(createBuildingResponse.message)
        expect(id).toEqual(createBuildingResponse.data)
    });

    it("FEUT-7: Bad Request", async () => {
        const formData = {
            name: '',
            category: '',
            address: '',
            notes: '',
            parent: '1-container',
        };

        const parentData = formData.parent.split('-')

        const createBuildingResponse = {
            message: "bad request"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createBuildingResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `/api/v1/entity`,
            {
                method: 'POST',
                headers:{
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes,
                    parentID: parentData[0],
                    parentCategory: parentData[1],
                }),
            }
        );
        expect(message).toEqual(createBuildingResponse.message)
        expect(id).toEqual(0)
    });
});

describe("Unit Tests for Entity Functions", () => {
    beforeEach(() => {
        // Mock the global fetch
        global.fetch = vi.fn() as FetchMock;
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    // Existing createEntity tests...

    describe("editEntity function", () => {
        it("FEUT-8: Edit Entity Request", async () => {
            const formData = {
                id: 1,
                category: 'item',
                name: 'Updated Test item',
                address: '',
                notes: 'Updated Test notes',
                parent: '2-container',
            };

            const parentData = formData.parent.split('-');

            const editEntityResponse = {
                message: "success",
                data: {
                    Entity: {
                        ID: 1,
                        Name: 'Updated Test item',
                        Notes: 'Updated Test notes'
                    },
                    Parent: {
                        ParentID: 2,
                        ParentCategory: 'container'
                    }
                }
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(editEntityResponse));

            const [message, entity] = await editEntity(formData);

            expect(global.fetch).toHaveBeenCalledWith(
                `/api/v1/entity`,
                {
                    method: 'PUT',
                    headers:{
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        id: formData.id,
                        address: formData.address,
                        category: formData.category,
                        name: formData.name,
                        notes: formData.notes,
                        parentID: parentData[0],
                        parentCategory: parentData[1],
                    }),
                }
            );
            expect(message).toEqual(editEntityResponse.message);
            expect(entity).toEqual(editEntityResponse.data);
        });

        it("FEUT-9: Edit Entity Bad Request", async () => {
            const formData = {
                id: 1,
                category: '',
                name: '',
                address: '',
                notes: '',
                parent: '',
            };

            const editEntityResponse = {
                message: "bad request",
                data: null
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(editEntityResponse));

            const [message, entity] = await editEntity(formData);

            expect(message).toEqual(editEntityResponse.message);
            expect(entity).toEqual({
                Entity: { ID: 0, Name: '' },
                Parent: { ParentID: 0, ParentCategory: '' }
            });
        });
    });

    describe("getEntity function", () => {
        it("FEUT-10: Get Entity Request", async () => {
            const id = 1;
            const category = 'item';

            const getEntityResponse = {
                message: "success",
                data: {
                    Entity: {
                        ID: 1,
                        Name: 'Test item',
                        Notes: 'Test notes'
                    },
                    Parent: {
                        ParentID: 2,
                        ParentCategory: 'container'
                    }
                }
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(getEntityResponse));

            const [message, entity] = await getEntity(id, category);

            expect(global.fetch).toHaveBeenCalledWith(
                `/api/v1/entity/${category}/${id}`
            );
            expect(message).toEqual(getEntityResponse.message);
            expect(entity).toEqual(getEntityResponse.data);
        });

        it("FEUT-11: Get Entity Not Found", async () => {
            const id = 999;
            const category = 'item';

            const getEntityResponse = {
                message: "not found",
                data: null
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(getEntityResponse));

            const [message, entity] = await getEntity(id, category);

            expect(message).toEqual(getEntityResponse.message);
            expect(entity).toEqual({
                Entity: { ID: 0, Name: '' },
                Parent: { ParentID: 0, ParentCategory: '' }
            });
        });
    });

    describe("getParents function", () => {
        it("FEUT-12: Get Parents Request", async () => {
            const category = 'item';

            const getParentsResponse = {
                message: "success",
                data: [
                    { ID: 1, Name: 'Parent 1', Category: 'container' },
                    { ID: 2, Name: 'Parent 2', Category: 'container' }
                ]
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(getParentsResponse));

            const [message, parents] = await getParents(category);

            expect(global.fetch).toHaveBeenCalledWith(
                `/api/v1/parents/${category}`
            );
            expect(message).toEqual(getParentsResponse.message);
            expect(parents).toEqual(getParentsResponse.data);
        });

        it("FEUT-13: Get Parents No Results", async () => {
            const category = 'item';

            const getParentsResponse = {
                message: "no results",
                data: []
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(getParentsResponse));

            const [message, parents] = await getParents(category);

            expect(message).toEqual(getParentsResponse.message);
            expect(parents).toEqual([]);
        });
    });

    describe("deleteEntity function", () => {
        it("FEUT-75: Delete Entity Request - Successful", async () => {
            const id = 1;
            const category = 'item';

            const deleteEntityResponse = {
                message: "success",
                data: null
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(deleteEntityResponse));

            const [message, error] = await deleteEntity(id, category);

            expect(global.fetch).toHaveBeenCalledWith(
                `/api/v1/entity/${category}/${id}`,
                {
                    method: 'DELETE',
                    headers:{
                        'Content-Type': 'application/json',
                    },
                }
            );
            expect(message).toEqual(deleteEntityResponse.message);
            expect(error).toEqual("");
        });

        it("FEUT-76: Delete Entity Request - Unsuccess", async () => {
            const id = 999;
            const category = 'item';

            const deleteEntityResponse = {
                message: "error",
                data: "Entity not found"
            };

            (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(deleteEntityResponse));

            const [message, error] = await deleteEntity(id, category);

            expect(global.fetch).toHaveBeenCalledWith(
                `/api/v1/entity/${category}/${id}`,
                {
                    method: 'DELETE',
                    headers:{
                        'Content-Type': 'application/json',
                    },
                }
            );
            expect(message).toEqual(deleteEntityResponse.message);
            expect(error).toEqual(deleteEntityResponse.data);
        });
    });
});
