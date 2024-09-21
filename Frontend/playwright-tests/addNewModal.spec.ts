import { test, expect } from '@playwright/test';

test('INT-3: Add New Modal Has Expected HTML and Loads Properly', async ({ page }) => {
  await page.goto('/');

  await expect(page.locator('.modal-add-entity-form')).toBeHidden();

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  await expect(page.locator('.modal-add-entity-form')).toBeVisible();
  await expect(page.locator('.drawer')).toBeHidden();

  await expect(page.locator('label[for=category]')).toBeVisible();
  await expect(page.locator('label[for=category]')).toContainText("Category");
  await expect(page.locator('select#category')).toBeVisible();

  await expect(page.locator('label[for=name]')).toBeVisible();
  await expect(page.locator('label[for=name]')).toContainText("Name");
  await expect(page.locator('input#name')).toBeVisible();

  await expect(page.locator('label[for=notes]')).toBeVisible();
  await expect(page.locator('label[for=notes]')).toContainText("Notes");
  await expect(page.locator('textarea#notes')).toBeVisible();

  await expect(page.locator('footer > button.variant-ghost-surface')).toBeVisible();
  await expect(page.locator('footer > button.variant-ghost-surface')).toContainText("Cancel");

  await expect(page.locator('footer > button.variant-filled')).toBeVisible();
  await expect(page.locator('footer > button.variant-filled')).toContainText("Submit");
});
