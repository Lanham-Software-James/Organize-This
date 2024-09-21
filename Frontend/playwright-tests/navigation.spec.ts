import { test, expect } from '@playwright/test';

test('navbar loads', async ({ page }) => {
  await page.goto('/');

  await expect(page.locator('.drawer')).toBeHidden();
  await expect(page.locator('.modal-add-entity-form')).toBeHidden();

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  await expect(page.locator('.drawer')).toBeVisible();

  await expect(page.locator('#pages')).toBeVisible();
  await expect(page.locator('#pages > p')).toBeVisible();
  await expect(page.locator('#pages > p')).toContainText("Overview");
  await expect(page.locator('#pages > li > a > button')).toBeVisible();

  await expect(page.locator('#tools')).toBeVisible();
  await expect(page.locator('#tools > p')).toContainText("Tools");
  await expect(page.locator('#tools > li > button')).toBeVisible();

  await expect(page.locator('#account')).toBeVisible();
  await expect(page.locator('#account > p')).toContainText("Account");
});
