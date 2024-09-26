import { expect, test } from '@playwright/test';

test('INT-1: Homepage Has Expected HTML', async ({ page }) => {
	await page.goto('/');
	await expect(page.locator('h1')).toBeVisible();
	await expect(page.locator('h1')).toContainText("Organize This!");

	await expect(page.locator('#hamburgerMenu')).toBeVisible();

	await expect(page.locator('h2')).toBeVisible();
	await expect(page.locator('h2')).toContainText("All Things");

	await expect(page.locator('#filter')).toBeVisible();

	await expect(page.locator('#th-name')).toBeVisible();
	await expect(page.locator('#th-name')).toContainText("Name");

	await expect(page.locator('#th-location')).toBeVisible();
	await expect(page.locator('#th-location')).toContainText("Location");

	await expect(page.locator('.paginator')).toBeVisible();
});
