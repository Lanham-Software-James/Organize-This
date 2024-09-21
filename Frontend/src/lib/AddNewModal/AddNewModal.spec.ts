import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createEntity } from './AddNewModal';
import { PUBLIC_API_URL } from '$env/static/public';

describe("Unit Tests for createEntity()", () => {
    beforeEach(() => {
        // Mock the global fetch
        global.fetch = vi.fn();
    });

    afterEach(() => {
        vi.resetAllMocks();
    });

    it("UT-1: Build Create Item Request", async () => {
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

    it("UT-2: Build Create Container Request", async () => {
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

    it("UT-3: Build Create Shelf Request", async () => {
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

    it("UT-4: Build Create Shelving Unit Request", async () => {
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

    it("UT-5: Build Create Room Request", async () => {
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

    it("UT-6: Build Create Building Request", async () => {
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
