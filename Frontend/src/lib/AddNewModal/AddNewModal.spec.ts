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
            category: 'item',
            name: 'Test item',
            address: '',
            notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity`,
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
    });

    it("container request", async () => {
        const formData = {
            category: 'container',
            name: 'Test container',
            address: '',
            notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity`,
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
    });

    it("shelf request", async () => {
        const formData = {
            category: 'shelf',
            name: 'Test shelf',
            address: '',
            notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity`,
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
    });

    it("shelving unit request", async () => {
        const formData = {
            category: 'shelvingunit',
            name: 'Test shelving unit',
            address: '',
            notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity`,
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
    });

    it("room request", async () => {
        const formData = {
            category: 'room',
            name: 'Test room',
            address: '',
            notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity`,
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
    });

    it("building request", async () => {
        const formData = {
            name: 'Test building',
            category: 'building',
            address: '888 test road',
            notes: 'Test notes'
        };

        await createEntity(formData);

        expect(global.fetch).toHaveBeenCalledWith(
            `${PUBLIC_API_URL}v1/entity`,
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
    });
});
