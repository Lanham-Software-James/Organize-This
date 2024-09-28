import { test, expect } from '@playwright/test';

test('INT-3: Add New Modal Has Expected Shell', async ({ page }) => {
  await page.goto('/');

  await expect(page.locator('.modal-add-entity-form')).toBeHidden();

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  await expect(page.locator('.modal-add-entity-form')).toBeVisible();
  await expect(page.locator('.drawer')).toBeHidden();

  await expect(page.locator('header')).toBeVisible();
  await expect(page.locator('header')).toContainText("Add New");

  await expect(page.locator('article')).toBeVisible();
  await expect(page.locator('article')).toContainText("Please complete the form to add a new item, container, shelf, shelving unit, room, or building.");

  await expect(page.locator('label[for=category]')).toBeVisible();
  await expect(page.locator('label[for=category]')).toContainText("Category");

  await expect(page.locator('select#category')).toBeVisible();
  await expect(page.locator('option[value="item"]')).toBeDefined();
	await expect(page.locator('option[value="category"]')).toBeDefined();
	await expect(page.locator('option[value="shelf"]')).toBeDefined();
	await expect(page.locator('option[value="shelvingunit"]')).toBeDefined();
	await expect(page.locator('option[value="room"]')).toBeDefined();
  await expect(page.locator('option[value="building"]')).toBeDefined();

  await expect(page.locator('footer > button.variant-ghost-surface')).toBeVisible();
  await expect(page.locator('footer > button.variant-ghost-surface')).toContainText("Cancel");

  await expect(page.locator('footer > button.variant-filled')).toBeVisible();
  await expect(page.locator('footer > button.variant-filled')).toContainText("Submit");
});

test('INT-4: Add New Modal Item Category Has Fields', async ({ page }) => {
  await page.goto('/');

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  await expect(page.locator('label[for=name]')).toBeVisible();
  await expect(page.locator('label[for=name]')).toContainText("Name");
  await expect(page.locator('input#name')).toBeVisible();

  await expect(page.locator('label[for=notes]')).toBeVisible();
  await expect(page.locator('label[for=notes]')).toContainText("Notes");
  await expect(page.locator('textarea#notes')).toBeVisible();
})

test('INT-5: Add New Modal Container Category Has Fields', async ({ page }) => {
  await page.goto('/');

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  // Select Different Category
  await page.getByLabel('Category:').selectOption('container');

  await expect(page.locator('label[for=name]')).toBeVisible();
  await expect(page.locator('label[for=name]')).toContainText("Name");
  await expect(page.locator('input#name')).toBeVisible();

  await expect(page.locator('label[for=notes]')).toBeVisible();
  await expect(page.locator('label[for=notes]')).toContainText("Notes");
  await expect(page.locator('textarea#notes')).toBeVisible();
})

test('INT-6: Add New Modal Shelf Category Has Fields', async ({ page }) => {
  await page.goto('/');

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  // Select Different Category
  await page.getByLabel('Category:').selectOption('shelf');

  await expect(page.locator('label[for=name]')).toBeVisible();
  await expect(page.locator('label[for=name]')).toContainText("Name");
  await expect(page.locator('input#name')).toBeVisible();

  await expect(page.locator('label[for=notes]')).toBeVisible();
  await expect(page.locator('label[for=notes]')).toContainText("Notes");
  await expect(page.locator('textarea#notes')).toBeVisible();
})

test('INT-7: Add New Modal Shelving Unit Category Has Fields', async ({ page }) => {
  await page.goto('/');

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  // Select Different Category
  await page.getByLabel('Category:').selectOption('shelvingunit');

  await expect(page.locator('label[for=name]')).toBeVisible();
  await expect(page.locator('label[for=name]')).toContainText("Name");
  await expect(page.locator('input#name')).toBeVisible();

  await expect(page.locator('label[for=notes]')).toBeVisible();
  await expect(page.locator('label[for=notes]')).toContainText("Notes");
  await expect(page.locator('textarea#notes')).toBeVisible();
})

test('INT-8: Add New Modal Room Category Has Fields', async ({ page }) => {
  await page.goto('/');

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  // Select Different Category
  await page.getByLabel('Category:').selectOption('room');

  await expect(page.locator('label[for=name]')).toBeVisible();
  await expect(page.locator('label[for=name]')).toContainText("Name");
  await expect(page.locator('input#name')).toBeVisible();

  await expect(page.locator('label[for=notes]')).toBeVisible();
  await expect(page.locator('label[for=notes]')).toContainText("Notes");
  await expect(page.locator('textarea#notes')).toBeVisible();
})

test('INT-9: Add New Modal Building Has Fields', async ({ page }) => {
  await page.goto('/');

  // Load Nav Bar
  await page.click('#hamburgerMenu');

  // Load Add New Modal
  await page.getByRole('button', { name: '+ Add New' }).click();

  // Select Different Category
  await page.getByLabel('Category:').selectOption('building');

  await expect(page.locator('label[for=name]')).toBeVisible();
  await expect(page.locator('label[for=name]')).toContainText("Name");
  await expect(page.locator('input#name')).toBeVisible();

  await expect(page.locator('label[for=address]')).toBeVisible();
  await expect(page.locator('label[for=address]')).toContainText("Address");
  await expect(page.locator('input#address')).toBeVisible();

  await expect(page.locator('label[for=notes]')).toBeVisible();
  await expect(page.locator('label[for=notes]')).toContainText("Notes");
  await expect(page.locator('textarea#notes')).toBeVisible();
})
