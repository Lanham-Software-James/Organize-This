import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createEntity } from './AddNewModal';
import { PUBLIC_API_URL } from '$env/static/public';

describe("testing send request", () => {
    beforeEach(() => {
        // Mock the global fetch
        global.fetch = vi.fn();
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    it("item request", async () => {
        const formData = {
            Category: 'item',
            Name: 'Test item',
            Address: '',
            Notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity-management/${formData.Category}`,
            {
                method: 'POST',
                body: JSON.stringify({
                    Address: formData.Address,
                    Name: formData.Name,
                    Notes: formData.Notes
                }),
            }
        );
    });

    it("container request", async () => {
        const formData = {
            Category: 'container',
            Name: 'Test container',
            Address: '',
            Notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity-management/${formData.Category}`,
            {
                method: 'POST',
                body: JSON.stringify({
                    Address: formData.Address,
                    Name: formData.Name,
                    Notes: formData.Notes
                }),
            }
        );
    });

    it("shelf request", async () => {
        const formData = {
            Category: 'shelf',
            Name: 'Test shelf',
            Address: '',
            Notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity-management/${formData.Category}`,
            {
                method: 'POST',
                body: JSON.stringify({
                    Address: formData.Address,
                    Name: formData.Name,
                    Notes: formData.Notes
                }),
            }
        );
    });

    it("shelving unit request", async () => {
        const formData = {
            Category: 'shelvingunit',
            Name: 'Test shelving unit',
            Address: '',
            Notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity-management/${formData.Category}`,
            {
                method: 'POST',
                body: JSON.stringify({
                    Address: formData.Address,
                    Name: formData.Name,
                    Notes: formData.Notes
                }),
            }
        );
    });

    it("room request", async () => {
        const formData = {
            Category: 'room',
            Name: 'Test room',
            Address: '',
            Notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity-management/${formData.Category}`,
            {
                method: 'POST',
                body: JSON.stringify({
                    Address: formData.Address,
                    Name: formData.Name,
                    Notes: formData.Notes
                }),
            }
        );
    });

    it("building request", async () => {
        const formData = {
            Category: 'building',
            Name: 'Test building',
            Address: '888 test road',
            Notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity-management/${formData.Category}`,
            {
                method: 'POST',
                body: JSON.stringify({
                    Address: formData.Address,
                    Name: formData.Name,
                    Notes: formData.Notes
                }),
            }
        );
    });
});
