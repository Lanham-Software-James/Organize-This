import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createEntity } from './AddNewModal';
import { PUBLIC_API_URL } from '$env/static/public';

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
            notes: 'Test notes'
        };

        const createItemResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createItemResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entity`,
            {
                method: 'POST',
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes
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
            notes: 'Test notes'
        };

        const createContainerResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createContainerResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entity`,
            {
                method: 'POST',
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes
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
            notes: 'Test notes'
        };

        const createShelfResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createShelfResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entity`,
            {
                method: 'POST',
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes
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
            notes: 'Test notes'
        };

        const createUnitResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createUnitResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entity`,
            {
                method: 'POST',
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes
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
            notes: 'Test notes'
        };

        const createRoomResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createRoomResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entity`,
            {
                method: 'POST',
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes
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
            notes: 'Test notes'
        };

        const createBuildingResponse = {
            data: 10,
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createBuildingResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entity`,
            {
                method: 'POST',
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes
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
            notes: ''
        };

        const createBuildingResponse = {
            message: "bad request"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(createBuildingResponse));

        const [message, id] = await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entity`,
            {
                method: 'POST',
                body: JSON.stringify({
                    address: formData.address,
                    category: formData.category,
                    name: formData.name,
                    notes: formData.notes
                }),
            }
        );
        expect(message).toEqual(createBuildingResponse.message)
        expect(id).toEqual(0)
    });
});
