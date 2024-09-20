import { expect, test } from '@playwright/test';

test.describe('home page layout components', () => {
	test('home page has expected h1 in layout', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('h1')).toBeVisible();
		await expect(page.locator('h1')).toContainText("Organize This!");
	});

	test('home page has hamburger menu in layout', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('#hamburgerMenu')).toBeVisible();
	});
});

test.describe('home page items pages components', () => {
	test('home page has expected h2', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('h2')).toBeVisible();
		await expect(page.locator('h2')).toContainText("All Things");
	});

	test('home page has filter button', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('#filter')).toBeVisible();
	});

	test('home page has expected table', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('table')).toBeVisible();
	});

	test('home page has expected table headers', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('#th-name')).toBeVisible();
		await expect(page.locator('#th-name')).toContainText("Name");

		await expect(page.locator('#th-location')).toBeVisible();
		await expect(page.locator('#th-location')).toContainText("Location");

	});

	test('home page has expected paginator', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('.paginator')).toBeVisible();
	});
});
