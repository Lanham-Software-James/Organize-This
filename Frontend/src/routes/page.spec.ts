import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { _getEntities as getEntities } from './+page';
import { PUBLIC_API_URL } from '$env/static/public';

function createFetchResponse(data: unknown) {
    return { json: () => new Promise((resolve) => resolve(data)) }
}

// Define a type for the mock fetch function
type FetchMock = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;

describe("Unit Tests for _getEntities()", () => {
    beforeEach(() => {
        // Mock the global fetch
        global.fetch = vi.fn() as FetchMock
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    it("FEUT-8: Get Entities Successful Response", async () => {
        let offset = 15;
        let limit = 25;

        const getEntitiesResponse =
        {
            data: {
                TotalCount: 117,
                Entities: [
                        {
                            ID: 25,
                            Name: "Test Item",
                            Category: "Item",
                            Parent: [],
                            Notes: "Test notes"
                        },
                        {
                            ID: 26,
                            Name: "Test Container",
                            Category: "Container",
                            Parent: [],
                            Notes: "Test notes"
                        },
                        {
                            ID: 27,
                            Name: "Test Shelf",
                            Category: "Shelf",
                            Parent: [],
                            Notes: "Test notes"
                        },
                        {
                            ID: 28,
                            Name: "Test Shelving Unit",
                            Category: "Shelving Unit",
                            Parent: [],
                            Notes: "Test notes"
                        },
                        {
                            ID: 29,
                            Name: "Test Room",
                            Category: "Room",
                            Parent: [],
                            Notes: "Test notes"
                        },
                        {
                            ID: 30,
                            Name: "Test Building",
                            Category: "Building",
                            Parent: [],
                            Notes: "Test notes"
                        },
                ],
            },
            message: "success"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(getEntitiesResponse));

        let [entities, size] = await getEntities(offset, limit);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entities?offset=${offset}&limit=${limit}`
        );

        expect(size).toEqual(getEntitiesResponse.data.TotalCount)

        expect(entities).toEqual(getEntitiesResponse.data.Entities)
    });

    it("FEUT-9: Get Entities Bad Response", async () => {
        let offset = 15;
        let limit = 25;

        const getEntitiesResponse =
        {
            data: {},
            message: "bad request"
        };

        (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(createFetchResponse(getEntitiesResponse));

        let [entities, size] = await getEntities(offset, limit);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}api/v1/entities?offset=${offset}&limit=${limit}`
        );

        expect(size).toEqual(0)

        expect(entities).toEqual([{
            ID: 0,
            Name: " ",
            Category: " ",
            Parent: [],
            Notes: " "
        }])
    });
});
