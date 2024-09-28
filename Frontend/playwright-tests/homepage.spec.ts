import { expect, test } from '@playwright/test';

test('INT-1: Homepage Has Expected HTML', async ({ page }) => {
	await page.goto('/');
	await expect(page.locator('h1')).toBeVisible();
	await expect(page.locator('h1')).toContainText("Organize This!");

	await expect(page.locator('#hamburgerMenu')).toBeVisible();

	await expect(page.locator('h2')).toBeVisible();
	await expect(page.locator('h2')).toContainText("All Things");

	await expect(page.locator('#filter')).toBeVisible();

	await expect(page.locator('i.fa-spinner')).toBeVisible();

	await expect(page.locator('.paginator')).toBeVisible();
	await expect(page.locator('select')).toBeVisible();

	await expect(page.locator('option[value=5]')).toBeDefined();
	await expect(page.locator('option[value=10]')).toBeDefined();
	await expect(page.locator('option[value=15]')).toBeDefined();
	await expect(page.locator('option[value=20]')).toBeDefined();
	await expect(page.locator('option[value=25]')).toBeDefined();
});
